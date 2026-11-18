package kit

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// legacyMarker is one marker of the pre-0.42.0 English format. The list is fixed on purpose (ADR-011, S5): a language heuristic would give false positives on
// code and technical names, while these strings only ever appeared in the old format.
type legacyMarker struct {
	name    string
	pattern *regexp.Regexp
}

// legacyMarkers: the opening of the old skill description, the old template and doc paths, the old language key, and the headings and labels that skills and agents
// used to write in English. Report and plan labels that readers accepted in both forms until 0.46.0 are here too: since 0.46.0 only the pt-BR form is read.
var legacyMarkers = []legacyMarker{
	{`"Use when"`, regexp.MustCompile(`Use when`)},
	{`"Do not use when"`, regexp.MustCompile(`Do not use when`)},
	{"templates/en, templates/pt-br ou templates/<lang>", regexp.MustCompile(`templates/(?:en|pt-br|<lang>)(?:/|\b)`)},
	{"docs/en", regexp.MustCompile(`docs/en\b`)},
	{"README.en", regexp.MustCompile(`README\.en\b`)},
	{"language: en ou language: pt-br", regexp.MustCompile(`(?:^|[^_\w])language: (?:en|pt-br)\b`)},
	{`"## External CLI reviewers"`, regexp.MustCompile(`## External CLI reviewers`)},
	{`"## The anti-delegation clause"`, regexp.MustCompile(`## The anti-delegation clause`)},
	{`"## Decisions needed"`, regexp.MustCompile(`## Decisions needed`)},
	{`"## Plan" (título sem "o")`, regexp.MustCompile(`## Plan\b`)},
	{`"### Slice"`, regexp.MustCompile(`### Slice\b`)},
	{`"Review status:"`, regexp.MustCompile(`Review status:`)},
	{`"## Findings"`, regexp.MustCompile(`## Findings`)},
	{`"### <SEVERITY>"`, regexp.MustCompile(`### <SEVERITY>`)},
	{`"Reported by:"`, regexp.MustCompile(`Reported by:`)},
	{`"## Comprehension"`, regexp.MustCompile(`## Comprehension`)},
	{`"## Not raised"`, regexp.MustCompile(`## Not raised`)},
	{`"Guide rule:"`, regexp.MustCompile(`Guide rule:`)},
	{`"coverage incomplete"`, regexp.MustCompile(`coverage incomplete`)},
	{`"Checks:"`, regexp.MustCompile(`Checks:`)},
}

// checkLegacyMarkers fails on every legacy marker in the kit sources: agents, commands, skills, templates, profiles and the eval cases. History in docs/, the
// CHANGELOG, the glossary and the recorded eval results quote the old format on purpose and are not read. Each error names the file, the line and the marker.
func checkLegacyMarkers(dirs map[string]string, errors *[]string) {
	roots := []string{dirs["agents"], dirs["commands"], dirs["skills"], dirs["templates"], dirs["profiles"]}
	if dirs["evals"] != "" {
		roots = append(roots, filepath.Join(dirs["evals"], "cases"))
	}
	for _, root := range roots {
		for _, path := range iterTextFiles(root) {
			text, err := readText(path)
			if err != nil {
				continue
			}
			for number, line := range strings.Split(text, "\n") {
				for _, marker := range legacyMarkers {
					if marker.pattern.MatchString(line) {
						*errors = append(*errors, fmt.Sprintf("%s:%d: marcador do formato antigo em inglês %s (ADR-011: o kit é só pt-BR; trava da 0.46.0)", path, number+1, marker.name))
					}
				}
			}
		}
	}
}
