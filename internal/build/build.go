package build

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/danielmalka/dev-harness/internal/kit"
)

const defaultVersion = "0.2.0"

type Target struct {
	OS   string
	Arch string
}

var DefaultTargets = []Target{
	{OS: "linux", Arch: "amd64"},
	{OS: "linux", Arch: "arm64"},
	{OS: "darwin", Arch: "amd64"},
	{OS: "darwin", Arch: "arm64"},
	{OS: "windows", Arch: "amd64"},
}

type Options struct {
	Output             string
	Targets            []Target
	NoBinaries         bool
	SkipMinimumCounts  bool
	UpdateRootManifest bool
}

type binaryManifest struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type manifest struct {
	Name        string           `json:"name"`
	Version     string           `json:"version"`
	License     string           `json:"license"`
	Generated   bool             `json:"generated"`
	GeneratedAt string           `json:"generated_at"`
	Source      string           `json:"source"`
	Runtime     runtimeInfo      `json:"runtime"`
	Inventory   inventory        `json:"inventory"`
	Binaries    []binaryManifest `json:"binaries"`
}

type runtimeInfo struct {
	Primary        string `json:"primary"`
	Load           string `json:"load"`
	CommandPrefix  string `json:"command_prefix"`
	MinimumVersion string `json:"minimum_version"`
}

type inventory struct {
	Agents    []string `json:"agents"`
	Commands  []string `json:"commands"`
	Skills    []string `json:"skills"`
	Profiles  []string `json:"profiles"`
	Templates []string `json:"templates"`
}

// Build stages and validates a Claude Code package before replacing its output.
func Build(root string, options Options, out io.Writer) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if options.Output == "" {
		options.Output = filepath.Join(root, "dist", "claude-code", "dev-harness")
	}
	if !filepath.IsAbs(options.Output) {
		options.Output, err = filepath.Abs(options.Output)
		if err != nil {
			return err
		}
	}
	targets := options.Targets
	if len(targets) == 0 {
		targets = append([]Target(nil), DefaultTargets...)
	}
	if err := validateTargets(targets); err != nil {
		return err
	}

	sourceReport, err := kit.Validate(root, kit.Options{
		SourceOnly:        true,
		SkipMinimumCounts: options.SkipMinimumCounts,
	})
	if err != nil {
		return err
	}
	if !sourceReport.Passed() {
		return fmt.Errorf("source validation failed: %s", strings.Join(sourceReport.Errors, "; "))
	}

	parent := filepath.Dir(options.Output)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(parent, ".dev-harness-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	pkg := filepath.Join(stage, "dev-harness")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		return err
	}

	for _, pair := range []struct{ source, destination string }{
		{source: ".agents", destination: "agents"},
		{source: ".skills", destination: "skills"},
		{source: ".commands", destination: "commands"},
		{source: "templates", destination: "templates"},
		{source: "profiles", destination: "profiles"},
	} {
		if err := copyTree(filepath.Join(root, pair.source), filepath.Join(pkg, pair.destination)); err != nil {
			return err
		}
	}
	pluginSource := filepath.Join(root, "adapters", "claude-code", "plugin")
	if err := copyTree(pluginSource, pkg); err != nil {
		return err
	}

	if !options.NoBinaries {
		for _, target := range targets {
			if err := buildBinary(root, pkg, target); err != nil {
				return err
			}
		}
	}

	version := os.Getenv("DEV_HARNESS_VERSION")
	if version == "" {
		version = defaultVersion
	}
	binaries, err := binaryEntries(pkg, targets, options.NoBinaries)
	if err != nil {
		return err
	}
	generatedAt := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	manifestData := manifest{
		Name:        "dev-harness",
		Version:     version,
		License:     "MIT",
		Generated:   true,
		GeneratedAt: generatedAt,
		Source:      "dh build",
		Runtime: runtimeInfo{
			Primary:        "claude-code",
			Load:           "claude --plugin-dir dist/claude-code/dev-harness --agent coordinator",
			CommandPrefix:  "/dev-harness:",
			MinimumVersion: "unverified",
		},
		Inventory: inventoryFor(pkg),
		Binaries:  binaries,
	}
	manifestBytes, err := json.MarshalIndent(manifestData, "", "  ")
	if err != nil {
		return err
	}
	manifestBytes = append(manifestBytes, '\n')
	if err := writeFile(filepath.Join(pkg, "harness-manifest.json"), manifestBytes, 0o644); err != nil {
		return err
	}
	plugin := map[string]any{
		"name":        "dev-harness",
		"version":     version,
		"description": "Portable AI-assisted development kit: coordinator plus specialist agents, skills and work commands with owner-authorized delivery.",
		"author":      map[string]string{"name": "malka"},
		"license":     "MIT",
		"keywords":    []string{"development", "agents", "coordinator", "qa", "review"},
	}
	pluginBytes, err := json.MarshalIndent(plugin, "", "  ")
	if err != nil {
		return err
	}
	pluginBytes = append(pluginBytes, '\n')
	if err := writeFile(filepath.Join(pkg, ".claude-plugin", "plugin.json"), pluginBytes, 0o644); err != nil {
		return err
	}
	generated := []byte("Generated by dh build. Do not edit. Change canonical sources and rebuild.\n")
	if err := writeFile(filepath.Join(pkg, "GENERATED.txt"), generated, 0o644); err != nil {
		return err
	}

	for _, pair := range []struct{ source, destination string }{
		{source: ".agents", destination: "agents"},
		{source: ".skills", destination: "skills"},
		{source: ".commands", destination: "commands"},
		{source: "templates", destination: "templates"},
		{source: "profiles", destination: "profiles"},
	} {
		if err := compareTrees(filepath.Join(root, pair.source), filepath.Join(pkg, pair.destination)); err != nil {
			return fmt.Errorf("copy mismatch: %s -> %s: %w", pair.source, pair.destination, err)
		}
	}
	if err := compareCopiedTree(pluginSource, pkg); err != nil {
		return fmt.Errorf("copy mismatch: adapters/claude-code/plugin -> package: %w", err)
	}

	packageReport, err := kit.Validate(pkg, kit.Options{SkipMinimumCounts: options.SkipMinimumCounts})
	if err != nil {
		return err
	}
	if !packageReport.Passed() {
		return fmt.Errorf("package validation failed: %s", strings.Join(packageReport.Errors, "; "))
	}

	if options.UpdateRootManifest {
		if err := writeFile(filepath.Join(root, "harness-manifest.json"), manifestBytes, 0o644); err != nil {
			return err
		}
	}
	if err := os.RemoveAll(options.Output); err != nil {
		return err
	}
	if err := os.Rename(pkg, options.Output); err != nil {
		return err
	}

	fmt.Fprintf(out, "manifest     version %s\n", version)
	fmt.Fprintf(out, "agents       %d files\n", len(manifestData.Inventory.Agents))
	fmt.Fprintf(out, "commands     %d files\n", len(manifestData.Inventory.Commands))
	fmt.Fprintf(out, "skills       %d files\n", len(manifestData.Inventory.Skills))
	fmt.Fprintf(out, "profiles     %d files\n", len(manifestData.Inventory.Profiles))
	fmt.Fprintf(out, "templates    %d files\n", len(manifestData.Inventory.Templates))
	fmt.Fprintf(out, "binaries     %d\n", len(manifestData.Binaries))
	fmt.Fprintf(out, "wrote        %s\n", options.Output)
	return nil
}

func validateTargets(targets []Target) error {
	seen := make(map[string]bool, len(targets))
	for _, target := range targets {
		if target.OS != "linux" && target.OS != "darwin" && target.OS != "windows" {
			return fmt.Errorf("unsupported target OS: %s", target.OS)
		}
		if target.Arch != "amd64" && target.Arch != "arm64" {
			return fmt.Errorf("unsupported target architecture: %s", target.Arch)
		}
		if target.OS == "windows" && target.Arch != "amd64" {
			return fmt.Errorf("unsupported Windows target: %s/%s", target.OS, target.Arch)
		}
		key := target.OS + "/" + target.Arch
		if seen[key] {
			return fmt.Errorf("duplicate target: %s", key)
		}
		seen[key] = true
	}
	return nil
}

func buildBinary(root, pkg string, target Target) error {
	platform := target.OS + "_" + target.Arch
	name := "dh"
	if target.OS == "windows" {
		name += ".exe"
	}
	output := filepath.Join(pkg, "bin", platform, name)
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}
	command := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w -buildid=", "-o", output, "./cmd/dh")
	command.Dir = root
	command.Env = withEnvironment(os.Environ(), map[string]string{
		"GOOS":        target.OS,
		"GOARCH":      target.Arch,
		"CGO_ENABLED": "0",
	})
	outputBytes, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(outputBytes))
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("go is required to build binaries")
		}
		if message != "" {
			return fmt.Errorf("build %s/%s: %w: %s", target.OS, target.Arch, err, message)
		}
		return fmt.Errorf("build %s/%s: %w", target.OS, target.Arch, err)
	}
	return nil
}

func binaryEntries(pkg string, targets []Target, noBinaries bool) ([]binaryManifest, error) {
	if noBinaries {
		return []binaryManifest{}, nil
	}
	entries := make([]binaryManifest, 0, len(targets))
	for _, target := range targets {
		name := "dh"
		if target.OS == "windows" {
			name += ".exe"
		}
		relative := filepath.ToSlash(filepath.Join("bin", target.OS+"_"+target.Arch, name))
		path := filepath.Join(pkg, filepath.FromSlash(relative))
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(data)
		entries = append(entries, binaryManifest{
			OS:     target.OS,
			Arch:   target.Arch,
			Path:   relative,
			SHA256: hex.EncodeToString(sum[:]),
		})
	}
	return entries, nil
}

func inventoryFor(pkg string) inventory {
	return inventory{
		Agents:    names(filepath.Join(pkg, "agents"), "*.md"),
		Commands:  names(filepath.Join(pkg, "commands"), "*.md"),
		Skills:    skillNames(filepath.Join(pkg, "skills")),
		Profiles:  names(filepath.Join(pkg, "profiles"), "*.yaml"),
		Templates: templateNames(filepath.Join(pkg, "templates")),
	}
}

func names(directory, pattern string) []string {
	paths, _ := filepath.Glob(filepath.Join(directory, pattern))
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
	}
	sort.Strings(result)
	return result
}

func skillNames(directory string) []string {
	entries, _ := os.ReadDir(directory)
	result := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && isRegular(filepath.Join(directory, entry.Name(), "SKILL.md")) {
			result = append(result, entry.Name())
		}
	}
	sort.Strings(result)
	return result
}

func templateNames(directory string) []string {
	result := make([]string, 0)
	_ = filepath.WalkDir(directory, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			return nil
		}
		relative, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}
		result = append(result, filepath.ToSlash(relative))
		return nil
	})
	sort.Strings(result)
	return result
}

func copyTree(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return os.MkdirAll(destination, 0o755)
		}
		if entry.Name() == ".gitkeep" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(destination, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink not supported: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := writeFile(target, data, info.Mode().Perm()); err != nil {
			return err
		}
		return nil
	})
}

func compareTrees(source, destination string) error {
	left := treeFiles(source)
	right := treeFiles(destination)
	for relative := range left {
		if !right[relative] {
			return fmt.Errorf("missing %s", relative)
		}
		leftData, err := os.ReadFile(filepath.Join(source, relative))
		if err != nil {
			return err
		}
		rightData, err := os.ReadFile(filepath.Join(destination, relative))
		if err != nil {
			return err
		}
		if string(leftData) != string(rightData) {
			return fmt.Errorf("differs %s", relative)
		}
	}
	for relative := range right {
		if !left[relative] {
			return fmt.Errorf("extra %s", relative)
		}
	}
	return nil
}

func compareCopiedTree(source, destination string) error {
	for relative := range treeFiles(source) {
		leftData, err := os.ReadFile(filepath.Join(source, relative))
		if err != nil {
			return err
		}
		rightData, err := os.ReadFile(filepath.Join(destination, relative))
		if err != nil {
			return fmt.Errorf("missing %s", relative)
		}
		if string(leftData) != string(rightData) {
			return fmt.Errorf("differs %s", relative)
		}
	}
	return nil
}

func treeFiles(root string) map[string]bool {
	files := make(map[string]bool)
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Name() != ".gitkeep" {
			relative, _ := filepath.Rel(root, path)
			files[relative] = true
		}
		return nil
	})
	return files
}

func writeFile(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, mode.Perm()); err != nil {
		return err
	}
	return os.Chmod(path, mode.Perm())
}

func isRegular(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func withEnvironment(environment []string, values map[string]string) []string {
	result := make([]string, 0, len(environment)+len(values))
	for _, item := range environment {
		name, _, ok := strings.Cut(item, "=")
		if !ok || values[name] == "" {
			result = append(result, item)
		}
	}
	for name, value := range values {
		result = append(result, name+"="+value)
	}
	return result
}
