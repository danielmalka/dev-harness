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
	graders, _ := filepath.Glob(filepath.Join(caseDir, "graders", "*.md"))
	sort.Strings(graders)
	for _, graderPath := range graders {
		text, err := readText(graderPath)
		if err != nil {
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
// escapes caseDir, or, if it stays inside, when it does not resolve to an
// existing file or directory. An empty value (field absent) is not an error.
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
	if !isFileOrDir(resolved) {
		*errors = append(*errors, fmt.Sprintf("%s: %s unresolved: %s", referencingPath, field, value))
	}
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
