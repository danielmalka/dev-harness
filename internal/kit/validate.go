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
	"sync"
)

const (
	MinAgents   = 13
	MinCommands = 17
	MinSkills   = 16
)

var (
	modelPattern    = lazyRe(`(?m)^model:\s*(\S+)`)
	fieldPattern    = lazyRe(`(?m)^([A-Za-z][A-Za-z0-9-]*):\s*(.*)$`)
	templatePattern = lazyRe(`templates/(?:(<[a-z-]+>|[a-z]{2}(?:-[a-z]{2})?)/)?([A-Za-z0-9_.-]+\.md)`)
	linkPattern     = lazyRe(`\[[^\]]*\]\(([^)]+)\)`)
	privatePattern  = lazyRe(`(/home/|/Users/|[A-Za-z]:\\)`)
)

// legacyTemplateLanguageDirs are the language subfolders `templates/` had up to 0.41.0. Since 0.42.0 (ADR-011) the templates sit directly in `templates/`, all in
// pt-BR, and one of these folders coming back fails the validator.
var legacyTemplateLanguageDirs = []string{"en", "pt-br"}

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
	checkLegacyMarkers(dirs, &report.Errors)
	checkDocumentedCounts(dirs["root"], report.Inventory, &report.Errors)
	checkRoadmaps(dirs["root"], &report.Errors)
	checkStatusDrift(dirs["root"], &report.Errors)
	checkPrivatePaths(dirs["root"], &report.Errors)
	if isDir(dirs["agents"]) && isDir(dirs["skills"]) {
		checkAntiDelegationClause(dirs, &report.Errors)
	}
	if isDir(dirs["evals"]) {
		checkEvalReferences(dirs, &report.Errors)
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
			report.Errors = append(report.Errors, "LICENSE ausente na raiz do kit")
		} else if text, err := readText(license); err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %v", license, err))
		} else if !strings.Contains(text, "MIT License") {
			report.Errors = append(report.Errors, "LICENSE não é MIT")
		}
	}
	return report, nil
}

func WriteReport(w io.Writer, report Report) {
	fmt.Fprintf(w, "alvo: %s\n", report.Target)
	fmt.Fprintf(w, "layout: %s\n", report.Layout)
	writeInventory(w, "agentes", report.Inventory.Agents)
	writeInventory(w, "comandos", report.Inventory.Commands)
	writeInventory(w, "skills", report.Inventory.Skills)
	writeInventory(w, "perfis", report.Inventory.Profiles)
	writeInventory(w, "templates", report.Inventory.Templates)
	fmt.Fprintf(w, "erros: %d\n", len(report.Errors))
	if report.Passed() {
		fmt.Fprintln(w, "status: aprovado")
		return
	}
	for _, problem := range report.Errors {
		fmt.Fprintf(w, "  - %s\n", problem)
	}
	fmt.Fprintln(w, "status: reprovado")
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
			return "", fmt.Errorf("alvo ausente: %s", abs)
		}
		return "", err
	}
	if _, err := os.Stat(resolved); err != nil {
		return "", fmt.Errorf("alvo ausente: %s", abs)
	}
	return resolved, nil
}

// ErrUnrecognizedLayout: the target is neither the kit source nor the generated package.
var ErrUnrecognizedLayout = errors.New("layout não reconhecido")

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
	return "", fmt.Errorf("%w: %s (esperado .agents/ e .commands/ para o layout de fontes, ou agents/ e .claude-plugin/ para o pacote)", ErrUnrecognizedLayout, target)
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
	names, _ := listTemplates(dirs["templates"])
	for name := range names {
		templates = append(templates, name)
	}
	sort.Strings(templates)
	return Inventory{Agents: agents, Commands: commands, Skills: skills, Profiles: profiles, Templates: templates}
}

func writeInventory(w io.Writer, name string, values []string) {
	fmt.Fprintf(w, "%s: %d\n", name, len(values))
	if len(values) == 0 {
		fmt.Fprintln(w, "  (nenhum)")
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

// listTemplates returns the template names (`templates/*.md`, no subfolder) and any legacy language subfolders that still exist.
func listTemplates(directory string) (names map[string]bool, legacyDirs []string) {
	names = make(map[string]bool)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return names, nil
	}
	for _, entry := range entries {
		switch {
		case entry.IsDir() && contains(legacyTemplateLanguageDirs, entry.Name()):
			legacyDirs = append(legacyDirs, entry.Name())
		case !entry.IsDir() && filepath.Ext(entry.Name()) == ".md":
			names[entry.Name()] = true
		}
	}
	return names, legacyDirs
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
		*errors = append(*errors, fmt.Sprintf("agentes: %d < %d", len(inv.Agents), MinAgents))
	}
	if len(inv.Commands) < MinCommands {
		*errors = append(*errors, fmt.Sprintf("comandos: %d < %d", len(inv.Commands), MinCommands))
	}
	if len(inv.Skills) < MinSkills {
		*errors = append(*errors, fmt.Sprintf("skills: %d < %d", len(inv.Skills), MinSkills))
	}
	for _, required := range []string{"base", "go-api", "typescript-web", "kotlin", "python", "typescript-api"} {
		if !contains(inv.Profiles, required) {
			*errors = append(*errors, "perfil ausente: "+required)
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
			*errors = append(*errors, path+": frontmatter ausente")
			continue
		}
		name := field(fm, "name")
		stem := strings.TrimSuffix(filepath.Base(path), ".md")
		if name != "" && name != stem {
			*errors = append(*errors, fmt.Sprintf("%s: name %q diferente de %q", path, name, stem))
		}
		if name == "" {
			*errors = append(*errors, path+": name ausente")
		}
		model := modelPattern().FindStringSubmatch(fm)
		if len(model) == 0 {
			*errors = append(*errors, path+": model ausente")
		} else if !allowedModels[model[1]] {
			*errors = append(*errors, fmt.Sprintf("%s: model %q fora de [haiku opus sonnet]", path, model[1]))
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
			*errors = append(*errors, fmt.Sprintf("%s: name %q diferente da pasta %q", path, field(fm, "name"), folder))
		}
		checkAuthor(path, fm, errors)
		// Since 0.45.0 only the pt-BR form passes: the old English "Use when" is rejected (ADR-011).
		if d := description(fm); !strings.HasPrefix(d, "Use quando") {
			*errors = append(*errors, path+`: description deve começar com "Use quando"`)
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
			*errors = append(*errors, path+": frontmatter ausente")
			continue
		}
		stem := strings.TrimSuffix(filepath.Base(path), ".md")
		if name := field(fm, "name"); name != stem {
			*errors = append(*errors, fmt.Sprintf("%s: name %q diferente de %q", path, name, stem))
		}
		checkAuthor(path, fm, errors)
		for _, key := range []string{"argument-hint", "metadata", "roles", "skills", "writes"} {
			if key == "roles" || key == "skills" || key == "writes" {
				if !metadataField(fm, key) {
					*errors = append(*errors, fmt.Sprintf("%s: metadata.%s ausente", path, key))
				}
			} else if field(fm, key) == "" && !hasKey(fm, key) {
				*errors = append(*errors, fmt.Sprintf("%s: %s ausente", path, key))
			}
		}
		for _, role := range listField(fm, "roles") {
			if !containsStem(listMarkdown(dirs["agents"]), role) {
				*errors = append(*errors, fmt.Sprintf("%s: papel %s ausente", path, role))
			}
		}
		for _, skill := range listField(fm, "skills") {
			if !containsSkill(listSkills(dirs["skills"]), skill) {
				*errors = append(*errors, fmt.Sprintf("%s: skill %s ausente", path, skill))
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
			*errors = append(*errors, path+": id ausente")
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
		*errors = append(*errors, path+": tools vazio")
	}
}

func checkAuthor(path, fm string, errors *[]string) {
	if field(fm, "author") == "" {
		*errors = append(*errors, path+": author ausente")
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
		for _, match := range linkPattern().FindAllStringSubmatch(text, -1) {
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
					*errors = append(*errors, fmt.Sprintf("%s: link sai da árvore: %s", path, href))
				}
				continue
			}
			if !isFileOrDir(target) {
				*errors = append(*errors, fmt.Sprintf("%s: link não resolvido %s", path, href))
			}
		}
	}
}

// checkTemplateRefs checks that `templates/` is a single folder with no language subfolder, and that every `templates/<file>.md` cited in the sources exists in it. Citations
// with the old subfolder (`templates/en/X.md`, `templates/pt-br/X.md`, `templates/<lang>/X.md`) are still accepted during the transition (0.42.0 to 0.44.0, until agents,
// commands and skills are rewritten in slices S2 and S3): the language of the citation is ignored and the file name is checked in `templates/`. The guard against the old format is S5.
func checkTemplateRefs(dirs map[string]string, errors *[]string) {
	names, legacyDirs := listTemplates(dirs["templates"])
	for _, dir := range legacyDirs {
		*errors = append(*errors, fmt.Sprintf("templates: a subpasta de idioma %s/ não existe mais: os templates ficam direto em templates/ (ADR-011)", dir))
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
			for _, match := range templatePattern().FindAllStringSubmatch(text, -1) {
				if name := match[2]; !names[name] {
					*errors = append(*errors, fmt.Sprintf("%s: template ausente templates/%s", path, name))
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
		if err != nil || !privatePattern().MatchString(text) {
			continue
		}
		if (strings.Contains(text, "code.claude.com")) && !strings.Contains(text, "/home/") && !strings.Contains(text, "/Users/") {
			continue
		}
		for lineNumber, line := range strings.Split(text, "\n") {
			if privatePattern().MatchString(line) && !strings.Contains(line, "code.claude.com") {
				*errors = append(*errors, fmt.Sprintf("%s:%d: caminho privado", rel, lineNumber+1))
			}
		}
	}
}

func checkPlugin(packageRoot string, errors *[]string) map[string]any {
	path := filepath.Join(packageRoot, ".claude-plugin", "plugin.json")
	if !isFile(path) {
		*errors = append(*errors, path+": ausente")
		return nil
	}
	text, err := readText(path)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: %v", path, err))
		return nil
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: JSON inválido: %v", path, err))
		return nil
	}
	if data["name"] != "dh" {
		*errors = append(*errors, "plugin.json: name deve ser dh")
	}
	if data["license"] != "MIT" {
		*errors = append(*errors, "plugin.json: license deve ser MIT")
	}
	if value, ok := data["version"].(string); !ok || value == "" {
		*errors = append(*errors, "plugin.json: version ausente")
	}
	return data
}

func checkRequiredPackageFiles(packageRoot string, errors *[]string) {
	for _, required := range []string{"settings.json", "hooks/hooks.json", "harness-manifest.json", "GENERATED.txt"} {
		if !isFile(filepath.Join(packageRoot, filepath.FromSlash(required))) {
			*errors = append(*errors, "pacote sem "+required)
		}
	}
	wrapper := filepath.Join(packageRoot, "bin", "dh")
	if !isExecutable(wrapper) {
		*errors = append(*errors, "pacote sem o executável bin/dh")
	}
	manifestPath := filepath.Join(packageRoot, "harness-manifest.json")
	text, err := readText(manifestPath)
	if err != nil {
		return
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: JSON inválido: %v", manifestPath, err))
		return
	}
	rawBinaries, exists := data["binaries"]
	if !exists {
		return
	}
	binaries, ok := rawBinaries.([]any)
	if !ok {
		*errors = append(*errors, "harness-manifest.json: binaries deve ser uma lista")
		return
	}
	for _, raw := range binaries {
		entry, ok := raw.(map[string]any)
		if !ok {
			*errors = append(*errors, "harness-manifest.json: entrada de binário inválida")
			continue
		}
		path, ok := entry["path"].(string)
		path = filepath.ToSlash(path)
		parts := strings.Split(path, "/")
		if !ok || len(parts) != 3 || parts[0] != "bin" || !strings.Contains(parts[1], "_") || (parts[2] != "dh" && parts[2] != "dh.exe") {
			*errors = append(*errors, "harness-manifest.json: o caminho do binário deve ser bin/<os>_<arch>/dh[.exe]")
			continue
		}
		if !isExecutable(filepath.Join(packageRoot, filepath.FromSlash(path))) {
			*errors = append(*errors, "pacote sem o executável "+path)
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
			*errors = append(*errors, "pacote sem "+pair.dest)
			continue
		}
		destFiles := relativeFiles(pair.dest)
		for rel := range difference(sourceFiles, destFiles) {
			*errors = append(*errors, fmt.Sprintf("pacote sem %s/%s", filepath.Base(pair.source), filepath.ToSlash(rel)))
		}
		for rel := range difference(destFiles, sourceFiles) {
			*errors = append(*errors, fmt.Sprintf("pacote com sobra %s/%s", filepath.Base(pair.dest), filepath.ToSlash(rel)))
		}
		for rel := range intersection(sourceFiles, destFiles) {
			sourcePath := filepath.Join(pair.source, rel)
			destPath := filepath.Join(pair.dest, rel)
			if !sameFileHash(sourcePath, destPath) {
				*errors = append(*errors, fmt.Sprintf("pacote difere em %s/%s", filepath.Base(pair.source), filepath.ToSlash(rel)))
			}
		}
	}
	for _, name := range []string{"doctor.sh", "validate.py"} {
		sourcePath := filepath.Join(sourceDirs["scripts"], name)
		destPath := filepath.Join(packageRoot, "scripts", name)
		if isFile(sourcePath) && isFile(destPath) && !sameFileHash(sourcePath, destPath) {
			*errors = append(*errors, "pacote difere em scripts/"+name)
		}
	}
	srcSkills := skillNames(sourceDirs["skills"])
	dstSkills := skillNames(filepath.Join(packageRoot, "skills"))
	for name := range difference(srcSkills, dstSkills) {
		*errors = append(*errors, "pacote sem a skill "+name)
	}
	for name := range difference(dstSkills, srcSkills) {
		*errors = append(*errors, "pacote com skill sobrando "+name)
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
			if contains([]string{".git", "dist", "node_modules", "__pycache__", ".codegraph", ".understand-anything", ".pytest_cache", ".worktrees"}, entry.Name()) {
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
	for _, match := range fieldPattern().FindAllStringSubmatch(fm, -1) {
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
	for _, match := range fieldPattern().FindAllStringSubmatch(fm, -1) {
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
		return nil, fmt.Errorf("%s: não é um arquivo comum (recusa seguir symlink ou ler arquivo especial)", path)
	}
	if info.Size() > maxReadableFileSize {
		return nil, fmt.Errorf("%s: %d bytes passam do limite de leitura de %d bytes", path, info.Size(), maxReadableFileSize)
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
		return nil, fmt.Errorf("%s: passa do limite de leitura de %d bytes", path, maxReadableFileSize)
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

// lazyRe compiles pattern on first use, so a package-level regexp costs nothing at init.
func lazyRe(pattern string) func() *regexp.Regexp {
	return sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(pattern) })
}
