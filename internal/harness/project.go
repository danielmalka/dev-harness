package harness

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CodeLanguages are the values accepted in `conventions.code_language` of project.yaml: the language of the PROJECT's commit messages and code comments
// (ADR-011). It is not the language option that was removed: the kit's conversation, documents and records are always pt-BR.
var CodeLanguages = []string{"en", "pt-br"}

// ProjectCheck is what the kit reads from the project.yaml of a harness dir. The file is not interpreted as a whole: only these two things matter to `dh`.
type ProjectCheck struct {
	// LegacyLanguage is true when the file still has the `language:` key from before 0.42.0. It is read, ignored and never rewritten.
	LegacyLanguage bool
	// CodeLanguage is the value of `conventions.code_language` ("" when the key is absent).
	CodeLanguage string
	// Problems lists what fails the file (today, only a `conventions.code_language` outside CodeLanguages).
	Problems []string
}

var (
	topKeyRe   = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_-]*)[ \t]*:(.*)$`)
	codeLangRe = regexp.MustCompile(`^[ \t]+code_language[ \t]*:(.*)$`)
)

// yamlScalar strips a trailing comment and quotes from a simple one-line value.
func yamlScalar(raw string) string {
	v := strings.TrimSpace(raw)
	if i := strings.Index(v, " #"); i >= 0 {
		v = strings.TrimSpace(v[:i])
	}
	if n := len(v); n >= 2 && (v[0] == '"' || v[0] == '\'') && v[n-1] == v[0] {
		v = v[1 : n-1]
	}
	return v
}

// ReadProjectCheck reads <dir>/project.yaml. A missing file is not an error (zero result). It reads line by line and only looks at the top-level keys `language` and
// `conventions`, so a project.yaml with any other content (lists, maps, comments) never fails here.
func ReadProjectCheck(path string) (ProjectCheck, error) {
	var out ProjectCheck
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxConfigBytes {
		return out, fmt.Errorf("%s não é um arquivo comum de até %d bytes", path, maxConfigBytes)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	inConventions := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if t := strings.TrimSpace(line); t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		if line[0] != ' ' && line[0] != '\t' {
			m := topKeyRe.FindStringSubmatch(line)
			inConventions = m != nil && m[1] == "conventions"
			out.LegacyLanguage = out.LegacyLanguage || (m != nil && m[1] == "language")
			continue
		}
		if m := codeLangRe.FindStringSubmatch(line); inConventions && m != nil {
			out.CodeLanguage = yamlScalar(m[1])
			if !isCodeLanguage(out.CodeLanguage) {
				out.Problems = append(out.Problems, fmt.Sprintf("conventions.code_language %q inválido (%s)", out.CodeLanguage, strings.Join(CodeLanguages, ", ")))
			}
		}
	}
	return out, nil
}

func isCodeLanguage(v string) bool {
	for _, ok := range CodeLanguages {
		if v == ok {
			return true
		}
	}
	return false
}

// HasLegacyLanguage says whether the text of a config.yaml has the top-level `language:` key from before 0.42.0.
func HasLegacyLanguage(configYAML string) bool {
	for _, line := range strings.Split(configYAML, "\n") {
		if m := topKeyRe.FindStringSubmatch(strings.TrimRight(line, "\r")); m != nil && m[1] == "language" {
			return true
		}
	}
	return false
}
