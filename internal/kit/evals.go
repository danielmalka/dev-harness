package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// checkAntiDelegationClause compares the anti-delegation clause repeated in
// .agents/coordinator.md and .skills/external-clis/SKILL.md. It stays silent
// when either file, or the fenced block it expects inside it, is absent —
// that pairing is only enforced once both sides exist to compare.
func checkAntiDelegationClause(dirs map[string]string, errors *[]string) {
	coordinatorPath := filepath.Join(dirs["agents"], "coordinator.md")
	skillPath := filepath.Join(dirs["skills"], "external-clis", "SKILL.md")
	if !isFile(coordinatorPath) || !isFile(skillPath) {
		return
	}
	coordinatorText, err := readText(coordinatorPath)
	if err != nil {
		return
	}
	skillText, err := readText(skillPath)
	if err != nil {
		return
	}
	coordinatorBlock, ok := extractFencedBlockAfterHeading(coordinatorText, "## External CLI reviewers")
	if !ok {
		return
	}
	skillBlock, ok := extractFencedBlockAfterHeading(skillText, "## The anti-delegation clause")
	if !ok {
		return
	}
	if !sameTrimmedLines(coordinatorBlock, skillBlock) {
		*errors = append(*errors, fmt.Sprintf(
			"anti-delegation clause differs between %s and %s",
			rootRelative(dirs["root"], coordinatorPath), rootRelative(dirs["root"], skillPath),
		))
	}
}

// extractFencedBlockAfterHeading returns the body of the first fenced code
// block (```...\n...\n```) that follows the given heading line.
func extractFencedBlockAfterHeading(text, heading string) (string, bool) {
	idx := strings.Index(text, heading)
	if idx < 0 {
		return "", false
	}
	rest := text[idx+len(heading):]
	fenceStart := strings.Index(rest, "```")
	if fenceStart < 0 {
		return "", false
	}
	afterOpen := rest[fenceStart+3:]
	newline := strings.Index(afterOpen, "\n")
	if newline < 0 {
		return "", false
	}
	body := afterOpen[newline+1:]
	fenceEnd := strings.Index(body, "```")
	if fenceEnd < 0 {
		return "", false
	}
	return body[:fenceEnd], true
}

// sameTrimmedLines compares two blocks line by line after trimming each
// line, so reformatted whitespace does not count as drift.
func sameTrimmedLines(a, b string) bool {
	linesA := strings.Split(strings.TrimRight(a, "\n"), "\n")
	linesB := strings.Split(strings.TrimRight(b, "\n"), "\n")
	if len(linesA) != len(linesB) {
		return false
	}
	for i := range linesA {
		if strings.TrimSpace(linesA[i]) != strings.TrimSpace(linesB[i]) {
			return false
		}
	}
	return true
}

func rootRelative(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

// checkEvalReferences resolves the relative paths that evals/cases/<id>/case.yaml
// and its graders/*.md files point at, and fails when a target is missing or
// escapes the case directory.
func checkEvalReferences(dirs map[string]string, errors *[]string) {
	casesRoot := filepath.Join(dirs["evals"], "cases")
	entries, err := os.ReadDir(casesRoot)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		caseDir := filepath.Join(casesRoot, entry.Name())
		checkCaseYAML(caseDir, errors)
		checkCaseGraders(caseDir, errors)
	}
}

func checkCaseYAML(caseDir string, errors *[]string) {
	caseYAML := filepath.Join(caseDir, "case.yaml")
	if !isFile(caseYAML) {
		return
	}
	text, err := readText(caseYAML)
	if err != nil {
		*errors = append(*errors, fmt.Sprintf("%s: %v", caseYAML, err))
		return
	}
	addDirs, scaffoldScript, historyFile := contextFields(text)
	for _, dir := range addDirs {
		checkFieldPath(caseYAML, caseDir, "context.add_dirs", dir, errors)
	}
	checkFieldPath(caseYAML, caseDir, "context.scaffold_script", scaffoldScript, errors)
	checkFieldPath(caseYAML, caseDir, "context.history_file", historyFile, errors)
}

func checkCaseGraders(caseDir string, errors *[]string) {
	if symlinkedSubdir(caseDir, "graders") {
		*errors = append(*errors, fmt.Sprintf(
			"%s: graders is a symbolic link; a case directory holds its own files",
			filepath.Join(caseDir, "graders"),
		))
		return
	}
	graders, _ := filepath.Glob(filepath.Join(caseDir, "graders", "*.md"))
	sort.Strings(graders)
	for _, graderPath := range graders {
		text, err := readText(graderPath)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", graderPath, err))
			continue
		}
		fm := frontmatter(text)
		if fm == "" {
			continue
		}
		target := yamlBlockLines(fm, "target")
		source, _ := yamlScalar(target, "source")
		path, hasPath := yamlScalar(target, "path")
		if source == "file" && hasPath {
			checkFieldPath(graderPath, caseDir, "target.path", path, errors)
		}
	}
}

// checkFieldPath resolves value relative to caseDir and fails when it
// escapes caseDir lexically, when a symlink sits anywhere between caseDir and
// the resolved path, or, if neither applies, when the path does not resolve
// to an existing file or directory. An empty value (field absent) is not an
// error.
func checkFieldPath(referencingPath, caseDir, field, value string, errors *[]string) {
	if value == "" {
		return
	}
	caseDirClean := filepath.Clean(caseDir)
	resolved := filepath.Clean(filepath.Join(caseDirClean, value))
	if resolved != caseDirClean && !strings.HasPrefix(resolved, caseDirClean+string(filepath.Separator)) {
		*errors = append(*errors, fmt.Sprintf("%s: %s escapes case directory: %s", referencingPath, field, value))
		return
	}
	if hasSymlinkComponent(caseDirClean, resolved) {
		*errors = append(*errors, fmt.Sprintf("%s: %s is or passes through a symlink: %s", referencingPath, field, value))
		return
	}
	if !isFileOrDir(resolved) {
		*errors = append(*errors, fmt.Sprintf("%s: %s unresolved: %s", referencingPath, field, value))
	}
}

// hasSymlinkComponent reports whether any path component between caseDir and
// resolved is a symlink, checked with Lstat component by component rather
// than filepath.EvalSymlinks — the latter would resolve the link and could
// silently allow exactly the escape this check exists to catch. A case
// directory is meant to be self-contained, so a symlink anywhere under it,
// pointing anywhere, has no legitimate use and is rejected outright rather
// than followed and re-checked for containment.
// symlinkedSubdir reports whether <parent>/<name> is a symbolic link. A case
// directory's own `graders/` or `fixtures/` is globbed by name, and a glob
// walks through a symlinked named segment: the files it then matches are
// regular files outside the case, which the per-file guard cannot tell from
// legitimate ones. The containment has to be decided on the directory, before
// anything under it is enumerated.
func symlinkedSubdir(parent, name string) bool {
	info, err := os.Lstat(filepath.Join(parent, name))
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func hasSymlinkComponent(caseDir, resolved string) bool {
	rel, err := filepath.Rel(caseDir, resolved)
	if err != nil || rel == "." {
		return false
	}
	current := caseDir
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			// Missing path: not a symlink escape. checkFieldPath's
			// existence check reports this case.
			return false
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

// contextFields extracts the known scalar and list keys from the `context:`
// block of a case.yaml file. It is a deliberately narrow, line-based reader
// for these specific keys, not a general YAML parser (the module has no YAML
// dependency and this task does not warrant adding one).
func contextFields(text string) (addDirs []string, scaffoldScript, historyFile string) {
	block := yamlBlockLines(text, "context")
	for i := 0; i < len(block); i++ {
		line := block[i]
		switch {
		case line == "add_dirs:":
			for i+1 < len(block) && strings.HasPrefix(block[i+1], "- ") {
				i++
				addDirs = append(addDirs, strings.Trim(strings.TrimSpace(strings.TrimPrefix(block[i], "-")), " \"'"))
			}
		case strings.HasPrefix(line, "add_dirs: ["):
			inline := strings.TrimSuffix(strings.TrimPrefix(line, "add_dirs: ["), "]")
			for _, part := range strings.Split(inline, ",") {
				part = strings.Trim(strings.TrimSpace(part), "\"'")
				if part != "" {
					addDirs = append(addDirs, part)
				}
			}
		default:
			if value, ok := scalarField(line, "scaffold_script"); ok {
				scaffoldScript = value
			}
			if value, ok := scalarField(line, "history_file"); ok {
				historyFile = value
			}
		}
	}
	return
}

func scalarField(trimmedLine, key string) (string, bool) {
	prefix := key + ":"
	if !strings.HasPrefix(trimmedLine, prefix) {
		return "", false
	}
	return strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmedLine, prefix)), "\"'"), true
}

// yamlBlockLines returns the trimmed lines of the indented block that
// follows a "key:" line at the same indentation, stopping at the first line
// indented at or above that key's own indentation.
func yamlBlockLines(text, key string) []string {
	lines := strings.Split(text, "\n")
	var block []string
	collecting := false
	keyIndent := 0
	for _, raw := range lines {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		if !collecting {
			if trimmed == key+":" {
				collecting = true
				keyIndent = indent
			}
			continue
		}
		if indent <= keyIndent {
			break
		}
		block = append(block, trimmed)
	}
	return block
}

func yamlScalar(lines []string, key string) (string, bool) {
	prefix := key + ":"
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, prefix)), "\"'"), true
		}
	}
	return "", false
}

// checkEvalsCoverage compares evals/ between source and package the same way
// checkCoverage compares the other trees, except results/, baselines/, and
// not-run/ are excluded only at the first level of evals/, and __pycache__/
// and *.pyc are excluded at any depth. It runs only when evals/ exists on at
// least one side.
func checkEvalsCoverage(sourceDirs map[string]string, packageRoot string, errors *[]string) {
	source := sourceDirs["evals"]
	dest := filepath.Join(packageRoot, "evals")
	if !isDir(source) && !isDir(dest) {
		return
	}
	sourceFiles := relativeFilesSkipping(source, evalsExcluded)
	if !isDir(dest) {
		*errors = append(*errors, "package missing "+dest)
		return
	}
	destFiles := relativeFilesSkipping(dest, evalsExcluded)
	for rel := range difference(sourceFiles, destFiles) {
		*errors = append(*errors, fmt.Sprintf("package missing evals/%s", filepath.ToSlash(rel)))
	}
	for rel := range difference(destFiles, sourceFiles) {
		*errors = append(*errors, fmt.Sprintf("package extra evals/%s", filepath.ToSlash(rel)))
	}
	for rel := range intersection(sourceFiles, destFiles) {
		if !sameFileHash(filepath.Join(source, rel), filepath.Join(dest, rel)) {
			*errors = append(*errors, fmt.Sprintf("package differs evals/%s", filepath.ToSlash(rel)))
		}
	}
}

// evalsExcluded mirrors the exclusion rule `dh build` applies when copying
// evals/ into the package: results/, baselines/, and not-run/ only at the
// first level of evals/; __pycache__/ and *.pyc at any depth.
func evalsExcluded(relative string) bool {
	parts := strings.Split(filepath.ToSlash(relative), "/")
	if len(parts) > 0 {
		switch parts[0] {
		case "results", "baselines", "not-run":
			return true
		}
	}
	for _, part := range parts {
		if part == "__pycache__" {
			return true
		}
	}
	return strings.HasSuffix(relative, ".pyc")
}

func relativeFilesSkipping(root string, skip func(relative string) bool) map[string]bool {
	files := relativeFiles(root)
	for rel := range files {
		if skip(rel) {
			delete(files, rel)
		}
	}
	return files
}

// checkEmbeddedAgentBodies compares each evals/cases/<id>/prompt.md's
// append_system_prompt block against the body of the agent it copies. The
// runner accepts no file reference for that field, so the copy is the only
// text a case actually runs against: when it drifts from .agents/<role>.md,
// every run of that case silently measures the old instructions. It stays
// silent when the block is absent or no shipped agent body matches it.
func checkEmbeddedAgentBodies(dirs map[string]string, errors *[]string) {
	bodies := agentBodies(dirs["agents"])
	if len(bodies) == 0 {
		return
	}
	prompts, _ := filepath.Glob(filepath.Join(dirs["evals"], "cases", "*", "prompt.md"))
	sort.Strings(prompts)
	for _, promptPath := range prompts {
		text, err := readText(promptPath)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", promptPath, err))
			continue
		}
		embedded, ok := yamlBlockScalar(text, "append_system_prompt")
		if !ok {
			continue
		}
		name, body, ok := matchingAgentBody(bodies, embedded)
		if !ok {
			continue
		}
		if !sameTrimmedLines(embedded, body) {
			*errors = append(*errors, fmt.Sprintf(
				"%s: append_system_prompt no longer matches %s",
				rootRelative(dirs["root"], promptPath),
				rootRelative(dirs["root"], filepath.Join(dirs["agents"], name+".md")),
			))
		}
	}
	checkFixtureAgentCopies(dirs, bodies, errors)
}

// checkFixtureAgentCopies applies the same rule to a case fixture that ships a
// copy of an agent body, so a case that hands the model a role prompt measures
// the current one.
func checkFixtureAgentCopies(dirs map[string]string, bodies map[string]string, errors *[]string) {
	caseDirs, _ := filepath.Glob(filepath.Join(dirs["evals"], "cases", "*"))
	sort.Strings(caseDirs)
	fixtures := []string{}
	for _, caseDir := range caseDirs {
		if !isDir(caseDir) {
			continue
		}
		if symlinkedSubdir(caseDir, "fixtures") {
			*errors = append(*errors, fmt.Sprintf(
				"%s: fixtures is a symbolic link; a case directory holds its own files",
				rootRelative(dirs["root"], filepath.Join(caseDir, "fixtures")),
			))
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(caseDir, "fixtures", "*.md"))
		fixtures = append(fixtures, matches...)
	}
	sort.Strings(fixtures)
	for _, fixturePath := range fixtures {
		text, err := readText(fixturePath)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("%s: %v", fixturePath, err))
			continue
		}
		name, body, ok := matchingAgentBody(bodies, text)
		if !ok {
			continue
		}
		if !sameTrimmedLines(text, body) {
			*errors = append(*errors, fmt.Sprintf(
				"%s: copy of %s is out of date",
				rootRelative(dirs["root"], fixturePath),
				rootRelative(dirs["root"], filepath.Join(dirs["agents"], name+".md")),
			))
		}
	}
}

// agentBodies reads every .agents/<role>.md and returns its body, frontmatter
// removed, keyed by role name.
func agentBodies(agentsDir string) map[string]string {
	bodies := map[string]string{}
	paths, _ := filepath.Glob(filepath.Join(agentsDir, "*.md"))
	for _, path := range paths {
		text, err := readText(path)
		if err != nil {
			continue
		}
		body, ok := bodyAfterFrontmatter(text)
		if !ok {
			continue
		}
		bodies[strings.TrimSuffix(filepath.Base(path), ".md")] = body
	}
	return bodies
}

// matchingAgentBody picks the agent whose body an embedded copy was taken from.
// Two independent paths claim a copy, because each one alone goes silent on a
// different drift: the opening line still matching a shipped agent, and a high
// share of lines in common. Identity by opening line alone misses a copy whose
// first line drifted with the rest; identity by overlap alone misses a copy
// that drifted past the threshold. A copy that lost both its opening line and
// most of its lines is no longer attributable to any agent and is left to the
// case author, which the skill records as a known limit.
func matchingAgentBody(bodies map[string]string, embedded string) (string, string, bool) {
	const minOverlap = 0.6
	names := make([]string, 0, len(bodies))
	for name := range bodies {
		names = append(names, name)
	}
	sort.Strings(names)
	first := firstNonEmptyLine(embedded)
	bestName, bestScore := "", 0.0
	for _, name := range names {
		if score := lineOverlap(embedded, bodies[name]); score > bestScore {
			bestName, bestScore = name, score
		}
	}
	if bestName != "" && bestScore >= minOverlap {
		return bestName, bodies[bestName], true
	}
	// Overlap did not claim it: fall back to the opening line, which is what
	// caught a partially rewritten copy before overlap existed.
	if first != "" {
		for _, name := range names {
			if firstNonEmptyLine(bodies[name]) == first {
				return name, bodies[name], true
			}
		}
	}
	return "", "", false
}

func firstNonEmptyLine(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

// lineOverlap is the fraction of the candidate's non-empty lines that also
// appear in the text, trimmed. Membership, not order: a body whose sections
// were reordered still identifies as the same agent, and the exact comparison
// that follows is what decides whether it drifted.
func lineOverlap(text, candidate string) float64 {
	present := map[string]bool{}
	for _, line := range strings.Split(text, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			present[trimmed] = true
		}
	}
	total, shared := 0, 0
	for _, line := range strings.Split(candidate, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		total++
		if present[trimmed] {
			shared++
		}
	}
	if total == 0 {
		return 0
	}
	return float64(shared) / float64(total)
}

// bodyAfterFrontmatter strips a leading YAML frontmatter block.
func bodyAfterFrontmatter(text string) (string, bool) {
	if !strings.HasPrefix(text, "---\n") {
		return "", false
	}
	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return "", false
	}
	return strings.TrimSpace(rest[end+len("\n---\n"):]), true
}

// yamlBlockScalar reads a literal block scalar (key: |) and returns it with the
// block indentation removed.
func yamlBlockScalar(text, key string) (string, bool) {
	lines := strings.Split(text, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimRight(line, " ") == key+": |" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return "", false
	}
	indent := ""
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent = line[:len(line)-len(strings.TrimLeft(line, " "))]
		break
	}
	if indent == "" {
		return "", false
	}
	collected := []string{}
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			collected = append(collected, "")
			continue
		}
		if !strings.HasPrefix(line, indent) {
			break
		}
		collected = append(collected, strings.TrimPrefix(line, indent))
	}
	return strings.TrimSpace(strings.Join(collected, "\n")), true
}
