package kit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	MinAgents   = 18
	MinCommands = 17
	MinSkills   = 16
)

var (
	modelPattern    = regexp.MustCompile(`(?m)^model:\s*(\S+)`)
	fieldPattern    = regexp.MustCompile(`(?m)^([A-Za-z][A-Za-z0-9-]*):\s*(.*)$`)
	templatePattern = regexp.MustCompile(`templates/(?:(<[a-z-]+>|[a-z]{2}(?:-[a-z]{2})?)/)?([A-Za-z0-9_.-]+\.md)`)
	linkPattern     = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	privatePattern  = regexp.MustCompile(`(/home/|/Users/|[A-Za-z]:\\)`)
)

var templateLanguages = []string{"en", "pt-br"}

var allowedModels = map[string]bool{
	"haiku":  true,
	"sonnet": true,
	"opus":   true,
}

// Options controls checks that are useful to callers embedding the validator.
// The command line always enforces the production minimum counts.
type Options struct {
	SourceOnly        bool
	SkipMinimumCounts bool
}

type Inventory struct {
	Agents    []string
	Commands  []string
	Skills    []string
	Profiles  []string
	Templates []string
}

type Report struct {
	Target    string
	Layout    string
	Inventory Inventory
	Errors    []string
}

func (r Report) Passed() bool {
	return len(r.Errors) == 0
}

// Validate checks a source tree, a kit checkout, or a generated package.
func Validate(target string, options Options) (Report, error) {
	resolved, err := resolveTarget(target)
	if err != nil {
		return Report{}, err
	}
	layout, err := detectLayout(resolved, options.SourceOnly)
	if err != nil {
		return Report{Target: resolved}, err
	}
	dirs := layoutDirs(resolved, layout)
	report := Report{Target: resolved, Layout: layout}

	checkAgents(dirs, &report.Errors)
	checkSkills(dirs, &report.Errors)
	checkCommands(dirs, &report.Errors)
	checkProfiles(dirs, &report.Errors)
	report.Inventory = inventory(dirs)
	if !options.SkipMinimumCounts {
		checkCounts(report.Inventory, &report.Errors)
	}
	checkLinks(dirs["root"], &report.Errors)
	checkTemplateRefs(dirs, &report.Errors)
	checkDocumentedCounts(dirs["root"], report.Inventory, &report.Errors)
	checkPrivatePaths(dirs["root"], &report.Errors)
	if isDir(dirs["agents"]) && isDir(dirs["skills"]) {
		checkAntiDelegationClause(dirs, &report.Errors)
	}
	if isDir(dirs["evals"]) {
		checkEvalReferences(dirs, &report.Errors)
		checkEmbeddedAgentBodies(dirs, &report.Errors)
	}

	if layout == "package" {
		checkPlugin(resolved, &report.Errors)
		checkRequiredPackageFiles(resolved, &report.Errors)
	} else if layout == "kit" {
		packageRoot := dirs["package"]
		checkPlugin(packageRoot, &report.Errors)
		checkCoverage(dirs, packageRoot, &report.Errors)
		checkEvalsCoverage(dirs, packageRoot, &report.Errors)
		checkRequiredPackageFiles(packageRoot, &report.Errors)
		license := filepath.Join(resolved, "LICENSE")
		if _, err := os.Stat(license); err != nil {
			report.Errors = append(report.Errors, "LICENSE missing at kit root")
		} else if text, err := readText(license); err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %v", license, err))
		} else if !strings.Contains(text, "MIT License") {
			report.Errors = append(report.Errors, "LICENSE is not MIT")
		}
	}
	return report, nil
}

func WriteReport(w io.Writer, report Report) {
	fmt.Fprintf(w, "target: %s\n", report.Target)
	fmt.Fprintf(w, "layout: %s\n", report.Layout)
	writeInventory(w, "agents", report.Inventory.Agents)
	writeInventory(w, "commands", report.Inventory.Commands)
	writeInventory(w, "skills", report.Inventory.Skills)
	writeInventory(w, "profiles", report.Inventory.Profiles)
	writeInventory(w, "templates", report.Inventory.Templates)
	fmt.Fprintf(w, "errors: %d\n", len(report.Errors))
	if report.Passed() {
		fmt.Fprintln(w, "status: passed")
		return
	}
	for _, problem := range report.Errors {
		fmt.Fprintf(w, "  - %s\n", problem)
	}
	fmt.Fprintln(w, "status: failed")
}

func Run(target string, options Options, w io.Writer) int {
	report, err := Validate(target, options)
	if err != nil {
		fmt.Fprintln(w, err)
		return 1
	}
	WriteReport(w, report)
	if !report.Passed() {
		return 1
	}
	return 0
}

func resolveTarget(target string) (string, error) {
	if target == "" {
		target = "."
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("missing target: %s", abs)
		}
		return "", err
	}
	if _, err := os.Stat(resolved); err != nil {
		return "", fmt.Errorf("missing target: %s", abs)
	}
	return resolved, nil
}

func detectLayout(target string, sourceOnly bool) (string, error) {
	if isDir(filepath.Join(target, "agents")) && isDir(filepath.Join(target, ".claude-plugin")) {
		return "package", nil
	}
	if isDir(filepath.Join(target, ".agents")) && isDir(filepath.Join(target, ".commands")) {
		packageRoot := filepath.Join(target, "dist", "claude-code", "dev-harness")
		if !sourceOnly && isDir(packageRoot) && isDir(filepath.Join(packageRoot, ".claude-plugin")) {
			return "kit", nil
		}
		return "source", nil
	}
	return "", fmt.Errorf("unrecognized layout: %s", target)
}

func layoutDirs(target, layout string) map[string]string {
	if layout == "source" || layout == "package" {
		names := map[string]string{
			"agents":   ".agents",
			"commands": ".commands",
			"skills":   ".skills",
		}
		if layout == "package" {
			names["agents"] = "agents"
			names["commands"] = "commands"
			names["skills"] = "skills"
		}
		return map[string]string{
			"root":      target,
			"agents":    filepath.Join(target, names["agents"]),
			"commands":  filepath.Join(target, names["commands"]),
			"skills":    filepath.Join(target, names["skills"]),
			"templates": filepath.Join(target, "templates"),
			"profiles":  filepath.Join(target, "profiles"),
			"scripts":   filepath.Join(target, "scripts"),
			"evals":     filepath.Join(target, "evals"),
		}
	}
	packageRoot := filepath.Join(target, "dist", "claude-code", "dev-harness")
	return map[string]string{
		"root":      target,
		"agents":    filepath.Join(target, ".agents"),
		"commands":  filepath.Join(target, ".commands"),
		"skills":    filepath.Join(target, ".skills"),
		"templates": filepath.Join(target, "templates"),
		"profiles":  filepath.Join(target, "profiles"),
		"scripts":   filepath.Join(target, "scripts"),
		"evals":     filepath.Join(target, "evals"),
		"package":   packageRoot,
	}
}

func inventory(dirs map[string]string) Inventory {
	agents := make([]string, 0)
	for _, path := range listMarkdown(dirs["agents"]) {
		agents = append(agents, strings.TrimSuffix(filepath.Base(path), ".md"))
	}
	commands := make([]string, 0)
	for _, path := range listMarkdown(dirs["commands"]) {
		commands = append(commands, strings.TrimSuffix(filepath.Base(path), ".md"))
	}
	skills := make([]string, 0)
	for _, path := range listSkills(dirs["skills"]) {
		skills = append(skills, filepath.Base(filepath.Dir(path)))
	}
	profiles := make([]string, 0)
	for _, path := range listProfiles(dirs["profiles"]) {
		profiles = append(profiles, strings.TrimSuffix(filepath.Base(path), ".yaml"))
	}
	templates := make([]string, 0)
	for lang, names := range listTemplates(dirs["templates"]) {
		for name := range names {
			templates = append(templates, lang+"/"+name)
		}
	}
	sort.Strings(templates)
	return Inventory{Agents: agents, Commands: commands, Skills: skills, Profiles: profiles, Templates: templates}
}

func writeInventory(w io.Writer, name string, values []string) {
	fmt.Fprintf(w, "%s: %d\n", name, len(values))
	if len(values) == 0 {
		fmt.Fprintln(w, "  (none)")
		return
	}
	fmt.Fprintf(w, "  %s\n", strings.Join(values, ", "))
}

func listMarkdown(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	paths := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" && entry.Name() != ".gitkeep" {
			paths = append(paths, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(paths)
	return paths
}

func listSkills(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	paths := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(directory, entry.Name(), "SKILL.md")
			if isFile(path) {
				paths = append(paths, path)
			}
		}
	}
	sort.Strings(paths)
	return paths
}

func listTemplates(directory string) map[string]map[string]bool {
	found := make(map[string]map[string]bool)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return found
	}
	for _, lang := range entries {
		if !lang.IsDir() || !contains(templateLanguages, lang.Name()) {
			continue
		}
		names := make(map[string]bool)
		files, err := os.ReadDir(filepath.Join(directory, lang.Name()))
		if err != nil {
			continue
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
				names[file.Name()] = true
			}
		}
		found[lang.Name()] = names
	}
	return found
}

func listProfiles(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	paths := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yaml" {
			paths = append(paths, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(paths)
	return paths
}

func checkCounts(inv Inventory, errors *[]string) {
	if len(inv.Agents) < MinAgents {
		*errors = append(*errors, fmt.Sprintf("agents: %d < %d", len(inv.Agents), MinAgents))
	}
	if len(inv.Commands) < MinCommands {
		*errors = append(*errors, fmt.Sprintf("commands: %d < %d", len(inv.Commands), MinCommands))
	}
	if len(inv.Skills) < MinSkills {
		*errors = append(*errors, fmt.Sprintf("skills: %d < %d", len(inv.Skills), MinSkills))
	}
	for _, required := range []string{"base", "go-api", "typescript-web"} {
		if !contains(inv.Profiles, required) {
			*errors = append(*errors, "missing profile: "+required)
		}
	}
}

func checkAgents(dirs map[string]string, errors *[]string) {
	for _, path := range listMarkdown(dirs["agents"]) {
		text, err := readText(path)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		fm := frontmatter(text)
		if fm == "" {
			*errors = append(*errors, path+": missing frontmatter")
			continue
		}
		name := field(fm, "name")
		stem := strings.TrimSuffix(filepath.Base(path), ".md")
		if name != "" && name != stem {
			*errors = append(*errors, fmt.Sprintf("%s: name %q != %q", path, name, stem))
		}
		if name == "" {
			*errors = append(*errors, path+": missing name")
		}
		model := modelPattern.FindStringSubmatch(fm)
		if len(model) == 0 {
			*errors = append(*errors, path+": missing model")
		} else if !allowedModels[model[1]] {
			*errors = append(*errors, fmt.Sprintf("%s: model %q not in [haiku opus sonnet]", path, model[1]))
		}
		checkAuthorAndTools(path, fm, errors)
	}
}

func checkSkills(dirs map[string]string, errors *[]string) {
	for _, path := range listSkills(dirs["skills"]) {
		text, err := readText(path)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		fm := frontmatter(text)
		folder := filepath.Base(filepath.Dir(path))
		if field(fm, "name") != folder {
			*errors = append(*errors, fmt.Sprintf("%s: name %q != folder %q", path, field(fm, "name"), folder))
		}
		checkAuthor(path, fm, errors)
		if !strings.HasPrefix(description(fm), "Use when") {
			*errors = append(*errors, path+`: description must start with "Use when"`)
		}
	}
}

func checkCommands(dirs map[string]string, errors *[]string) {
	for _, path := range listMarkdown(dirs["commands"]) {
		text, err := readText(path)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		fm := frontmatter(text)
		if fm == "" {
			*errors = append(*errors, path+": missing frontmatter")
			continue
		}
		stem := strings.TrimSuffix(filepath.Base(path), ".md")
		if name := field(fm, "name"); name != stem {
			*errors = append(*errors, fmt.Sprintf("%s: name %q != %q", path, name, stem))
		}
		checkAuthor(path, fm, errors)
		for _, key := range []string{"argument-hint", "metadata", "roles", "skills", "writes"} {
			if key == "roles" || key == "skills" || key == "writes" {
				if !metadataField(fm, key) {
					*errors = append(*errors, fmt.Sprintf("%s: missing metadata.%s", path, key))
				}
			} else if field(fm, key) == "" && !hasKey(fm, key) {
				*errors = append(*errors, fmt.Sprintf("%s: missing %s", path, key))
			}
		}
		for _, role := range listField(fm, "roles") {
			if !containsStem(listMarkdown(dirs["agents"]), role) {
				*errors = append(*errors, fmt.Sprintf("%s: missing role %s", path, role))
			}
		}
		for _, skill := range listField(fm, "skills") {
			if !containsSkill(listSkills(dirs["skills"]), skill) {
				*errors = append(*errors, fmt.Sprintf("%s: missing skill %s", path, skill))
			}
		}
	}
}

func checkProfiles(dirs map[string]string, errors *[]string) {
	for _, path := range listProfiles(dirs["profiles"]) {
		text, err := readText(path)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
			continue
		}
		id := regexp.MustCompile(`(?m)^id:\s*(\S+)`).FindStringSubmatch(text)
		stem := strings.TrimSuffix(filepath.Base(path), ".yaml")
		if len(id) == 0 {
			*errors = append(*errors, path+": missing id")
		} else if id[1] != stem {
			*errors = append(*errors, fmt.Sprintf("%s: id %q != %q", path, id[1], stem))
		}
	}
}

// checkAuthorAndTools accepts an absent tools field (the agent inherits every
// tool) but rejects a present, empty one, which fails to launch at runtime.
func checkAuthorAndTools(path, fm string, errors *[]string) {
	checkAuthor(path, fm, errors)
	if regexp.MustCompile(`(?m)^tools:`).MatchString(fm) && !hasListField(fm, "tools") {
		*errors = append(*errors, path+": empty tools")
	}
}

func checkAuthor(path, fm string, errors *[]string) {
	if field(fm, "author") == "" {
		*errors = append(*errors, path+": missing author")
	}
}

func checkLinks(root string, errors *[]string) {
	for _, path := range iterTextFiles(root) {
		if filepath.Ext(path) != ".md" {
			continue
		}
		text, err := readText(path)
		if err != nil {
			continue
		}
		for _, match := range linkPattern.FindAllStringSubmatch(text, -1) {
			href := match[1]
			raw := strings.TrimSpace(strings.SplitN(href, "#", 2)[0])
			if raw == "" || strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "mailto:") || strings.HasPrefix(raw, "#") || regexp.MustCompile(`^[a-z]+:`).MatchString(raw) {
				continue
			}
			target := filepath.Clean(filepath.Join(filepath.Dir(path), raw))
			rootAbs, _ := filepath.Abs(root)
			targetAbs, _ := filepath.Abs(target)
			if !strings.HasPrefix(targetAbs, rootAbs) {
				if strings.HasPrefix(href, "../") || strings.HasPrefix(href, "./") || strings.HasPrefix(href, "/") || strings.Contains(href, "/") {
					*errors = append(*errors, fmt.Sprintf("%s: link escapes tree: %s", path, href))
				}
				continue
			}
			if !isFileOrDir(target) {
				*errors = append(*errors, fmt.Sprintf("%s: unresolved link %s", path, href))
			}
		}
	}
}

func checkTemplateRefs(dirs map[string]string, errors *[]string) {
	byLang := listTemplates(dirs["templates"])
	for _, lang := range templateLanguages {
		if byLang[lang] == nil {
			*errors = append(*errors, "templates: missing language directory "+lang+"/")
		}
	}
	langs := make([]string, 0)
	for _, lang := range templateLanguages {
		if byLang[lang] != nil {
			langs = append(langs, lang)
		}
	}
	if len(langs) > 1 {
		base := byLang[langs[0]]
		for _, lang := range langs[1:] {
			for _, name := range sortedKeys(symmetricDifference(base, byLang[lang])) {
				*errors = append(*errors, fmt.Sprintf("templates: %s is not present in every language (%s vs %s)", name, langs[0], lang))
			}
		}
	}
	for _, root := range []string{dirs["agents"], dirs["commands"], dirs["skills"]} {
		if !isDir(root) {
			continue
		}
		for _, path := range iterTextFiles(root) {
			if filepath.Ext(path) != ".md" && filepath.Ext(path) != ".yaml" {
				continue
			}
			text, err := readText(path)
			if err != nil {
				continue
			}
			for _, match := range templatePattern.FindAllStringSubmatch(text, -1) {
				lang, name := match[1], match[2]
				targets := langs
				if lang != "" && !strings.HasPrefix(lang, "<") {
					targets = []string{lang}
				}
				for _, target := range targets {
					if !byLang[target][name] {
						*errors = append(*errors, fmt.Sprintf("%s: missing template %s/%s", path, target, name))
					}
				}
			}
		}
	}
}

func checkPrivatePaths(root string, errors *[]string) {
	for _, path := range iterTextFiles(root) {
		if filepath.Ext(path) == ".html" || filepath.Base(path) == "validate.py" {
			continue
		}
		rel, _ := filepath.Rel(root, path)
		relSlash := filepath.ToSlash(rel)
		if strings.HasPrefix(relSlash, "evals/baselines/") || strings.HasPrefix(relSlash, "evals/results/") {
			continue
		}
		text, err := readText(path)
		if err != nil || !privatePattern.MatchString(text) {
			continue
		}
		if (strings.Contains(text, "code.claude.com")) && !strings.Contains(text, "/home/") && !strings.Contains(text, "/Users/") {
			continue
		}
		for lineNumber, line := range strings.Split(text, "\n") {
			if privatePattern.MatchString(line) && !strings.Contains(line, "code.claude.com") {
				*errors = append(*errors, fmt.Sprintf("%s:%d: private path", rel, lineNumber+1))
			}
		}
	}
}

func checkPlugin(packageRoot string, errors *[]string) map[string]any {
	path := filepath.Join(packageRoot, ".claude-plugin", "plugin.json")
	if !isFile(path) {
		*errors = append(*errors, path+": missing")
		return nil
	}
	text, err := readText(path)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
		return nil
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: invalid JSON: %v", path, err))
		return nil
	}
	if data["name"] != "dev-harness" {
		*errors = append(*errors, "plugin.json name must be dev-harness")
	}
	if data["license"] != "MIT" {
		*errors = append(*errors, "plugin.json license must be MIT")
	}
	if value, ok := data["version"].(string); !ok || value == "" {
		*errors = append(*errors, "plugin.json missing version")
	}
	return data
}

func checkRequiredPackageFiles(packageRoot string, errors *[]string) {
	for _, required := range []string{"settings.json", "hooks/hooks.json", "harness-manifest.json", "GENERATED.txt"} {
		if !isFile(filepath.Join(packageRoot, filepath.FromSlash(required))) {
			*errors = append(*errors, "package missing "+required)
		}
	}
	wrapper := filepath.Join(packageRoot, "bin", "dh")
	if !isExecutable(wrapper) {
		*errors = append(*errors, "package missing executable bin/dh")
	}
	manifestPath := filepath.Join(packageRoot, "harness-manifest.json")
	text, err := readText(manifestPath)
	if err != nil {
		return
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: invalid JSON: %v", manifestPath, err))
		return
	}
	rawBinaries, exists := data["binaries"]
	if !exists {
		return
	}
	binaries, ok := rawBinaries.([]any)
	if !ok {
		*errors = append(*errors, "harness-manifest.json binaries must be an array")
		return
	}
	for _, raw := range binaries {
		entry, ok := raw.(map[string]any)
		if !ok {
			*errors = append(*errors, "harness-manifest.json contains an invalid binary entry")
			continue
		}
		path, ok := entry["path"].(string)
		path = filepath.ToSlash(path)
		parts := strings.Split(path, "/")
		if !ok || len(parts) != 3 || parts[0] != "bin" || !strings.Contains(parts[1], "_") || (parts[2] != "dh" && parts[2] != "dh.exe") {
			*errors = append(*errors, "harness-manifest.json binary path must be bin/<os>_<arch>/dh[.exe]")
			continue
		}
		if !isExecutable(filepath.Join(packageRoot, filepath.FromSlash(path))) {
			*errors = append(*errors, "package missing executable "+path)
		}
	}
	if len(binaries) == 0 {
		return
	}
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0
}

func checkCoverage(sourceDirs map[string]string, packageRoot string, errors *[]string) {
	pairs := []struct {
		source string
		dest   string
	}{
		{sourceDirs["agents"], filepath.Join(packageRoot, "agents")},
		{sourceDirs["commands"], filepath.Join(packageRoot, "commands")},
		{sourceDirs["skills"], filepath.Join(packageRoot, "skills")},
		{sourceDirs["templates"], filepath.Join(packageRoot, "templates")},
		{sourceDirs["profiles"], filepath.Join(packageRoot, "profiles")},
	}
	for _, pair := range pairs {
		sourceFiles := relativeFiles(pair.source)
		if !isDir(pair.dest) {
			*errors = append(*errors, "package missing "+pair.dest)
			continue
		}
		destFiles := relativeFiles(pair.dest)
		for rel := range difference(sourceFiles, destFiles) {
			*errors = append(*errors, fmt.Sprintf("package missing %s/%s", filepath.Base(pair.source), filepath.ToSlash(rel)))
		}
		for rel := range difference(destFiles, sourceFiles) {
			*errors = append(*errors, fmt.Sprintf("package extra %s/%s", filepath.Base(pair.dest), filepath.ToSlash(rel)))
		}
		for rel := range intersection(sourceFiles, destFiles) {
			sourcePath := filepath.Join(pair.source, rel)
			destPath := filepath.Join(pair.dest, rel)
			if !sameFileHash(sourcePath, destPath) {
				*errors = append(*errors, fmt.Sprintf("package differs %s/%s", filepath.Base(pair.source), filepath.ToSlash(rel)))
			}
		}
	}
	for _, name := range []string{"doctor.sh", "validate.py"} {
		sourcePath := filepath.Join(sourceDirs["scripts"], name)
		destPath := filepath.Join(packageRoot, "scripts", name)
		if isFile(sourcePath) && isFile(destPath) && !sameFileHash(sourcePath, destPath) {
			*errors = append(*errors, "package differs scripts/"+name)
		}
	}
	srcSkills := skillNames(sourceDirs["skills"])
	dstSkills := skillNames(filepath.Join(packageRoot, "skills"))
	for name := range difference(srcSkills, dstSkills) {
		*errors = append(*errors, "package missing skill "+name)
	}
	for name := range difference(dstSkills, srcSkills) {
		*errors = append(*errors, "package extra skill "+name)
	}
	sort.Strings(*errors)
}

func relativeFiles(root string) map[string]bool {
	files := make(map[string]bool)
	if !isDir(root) {
		return files
	}
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}
		if entry.Name() != ".gitkeep" {
			rel, _ := filepath.Rel(root, path)
			files[rel] = true
		}
		return nil
	})
	return files
}

func sameFileHash(left, right string) bool {
	leftHash, ok := fileHash(left)
	if !ok {
		return false
	}
	rightHash, ok := fileHash(right)
	return ok && leftHash == rightHash
}

func fileHash(path string) (string, bool) {
	data, err := readFileGuarded(path)
	if err != nil {
		return "", false
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), true
}

func iterTextFiles(root string) []string {
	paths := make([]string, 0)
	if !isDir(root) {
		return paths
	}
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			if contains([]string{".git", "dist", "node_modules", "__pycache__", ".codegraph", ".understand-anything", ".pytest_cache"}, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if contains([]string{".md", ".yaml", ".yml", ".html", ".py", ".sh", ".json", ".txt"}, ext) {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths
}

func frontmatter(text string) string {
	if !strings.HasPrefix(text, "---\n") {
		return ""
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return ""
	}
	return text[4 : 4+end]
}

func field(fm, name string) string {
	for _, match := range fieldPattern.FindAllStringSubmatch(fm, -1) {
		if match[1] != name {
			continue
		}
		value := strings.TrimSpace(match[2])
		if value == "|" || value == ">" {
			continue
		}
		return strings.Trim(strings.TrimSpace(value), "\"'")
	}
	return ""
}

func description(fm string) string {
	value := field(fm, "description")
	if value != "" {
		return value
	}
	lines := strings.Split(fm, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "description: |" && strings.TrimSpace(line) != "description: >" {
			continue
		}
		for _, next := range lines[i+1:] {
			next = strings.TrimSpace(next)
			if next != "" {
				return next
			}
		}
	}
	return ""
}

func hasKey(fm, name string) bool {
	for _, match := range fieldPattern.FindAllStringSubmatch(fm, -1) {
		if match[1] == name {
			return true
		}
	}
	return false
}

func metadataField(fm, name string) bool {
	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, name+":") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, name+":")) != ""
		}
	}
	return false
}

func hasListField(fm, name string) bool {
	lines := strings.Split(fm, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, name+": [") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(trimmed, "]"), name+": [")) != ""
		}
		if trimmed != name+":" {
			continue
		}
		for _, next := range lines[i+1:] {
			next = strings.TrimSpace(next)
			if next == "" {
				continue
			}
			return strings.HasPrefix(next, "-")
		}
	}
	return false
}

func listField(fm, name string) []string {
	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimSpace(line)
		prefix := name + ":"
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		value = strings.TrimPrefix(value, "[")
		value = strings.TrimSuffix(value, "]")
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.Trim(strings.TrimSpace(part), "\"'")
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}
	return nil
}

// maxReadableFileSize bounds every text read and hash comparison the
// validator performs. The largest file it legitimately reads today, across
// .agents/, .skills/, .commands/, templates/, profiles/ and evals/ (excluding
// evals/results/ and evals/baselines/, which no check here ever opens), is
// docs/plano-produto.html at ~96KiB (measured 2026-09-22). 5MiB leaves about
// 50x headroom while still bounding memory and I/O against an oversized or
// adversarial file — such as a case a third party contributed under
// evals/cases/ (see readFileGuarded).
const maxReadableFileSize = 5 * 1024 * 1024 // 5MiB

// readFileGuarded is the one place the validator turns a path into bytes.
// Several of its callers resolve paths declared inside an eval case a third
// party may have authored (evals/cases/<id>/case.yaml, graders/*.md,
// prompt.md, fixtures/) — an untrusted-input boundary — so every read stays
// inside two limits regardless of caller:
//
//  1. Lstat, not Stat: a symlink (or a FIFO, device, or other special file)
//     is rejected outright rather than followed, because a kit package is
//     meant to be self-contained and portable and a symlink inside it has no
//     legitimate use.
//  2. maxReadableFileSize, enforced twice: once from the size Lstat reports,
//     before anything is allocated, and once from what the read actually
//     returns, so a file that grows after Lstat cannot slip past the ceiling.
func readFileGuarded(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%s: not a regular file (refusing to follow a symlink or read a special file)", path)
	}
	if info.Size() > maxReadableFileSize {
		return nil, fmt.Errorf("%s: %d bytes exceeds the %d byte read limit", path, info.Size(), maxReadableFileSize)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }() // read-only: nothing to flush, nothing to report
	data, err := io.ReadAll(io.LimitReader(file, maxReadableFileSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxReadableFileSize {
		return nil, fmt.Errorf("%s: exceeds the %d byte read limit", path, maxReadableFileSize)
	}
	return data, nil
}

func readText(path string) (string, error) {
	data, err := readFileGuarded(path)
	if err != nil {
		return "", err
	}
	// Normalize CRLF once, here, so every line-based check downstream sees the
	// same shape. A file saved on Windows otherwise slips past a key match or a
	// frontmatter prefix and the check that should have failed skips in silence.
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func isFileOrDir(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func containsStem(paths []string, wanted string) bool {
	for _, path := range paths {
		if strings.TrimSuffix(filepath.Base(path), ".md") == wanted {
			return true
		}
	}
	return false
}

func containsSkill(paths []string, wanted string) bool {
	for _, path := range paths {
		if filepath.Base(filepath.Dir(path)) == wanted {
			return true
		}
	}
	return false
}

func skillNames(root string) map[string]bool {
	names := make(map[string]bool)
	for _, path := range listSkills(root) {
		names[filepath.Base(filepath.Dir(path))] = true
	}
	return names
}

func difference[T comparable](left, right map[T]bool) map[T]bool {
	result := make(map[T]bool)
	for value := range left {
		if !right[value] {
			result[value] = true
		}
	}
	return result
}

func intersection[T comparable](left, right map[T]bool) map[T]bool {
	result := make(map[T]bool)
	for value := range left {
		if right[value] {
			result[value] = true
		}
	}
	return result
}

func symmetricDifference[T comparable](left, right map[T]bool) map[T]bool {
	result := difference(left, right)
	for value := range difference(right, left) {
		result[value] = true
	}
	return result
}

func sortedKeys[T ~string](values map[T]bool) []T {
	keys := make([]T, 0, len(values))
	for value := range values {
		keys = append(keys, value)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
