package kit

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// Each fixed marker of the old English format fails in a kit source with the file and the line, and the same text without the marker passes (ADR-011, S5).
func TestLegacyMarkersAreRejectedWithFileAndLine(t *testing.T) {
	markers := []struct{ name, text string }{
		{`"Use when"`, "Use when testing"},
		{`"Do not use when"`, "Do not use when the work is trivial"},
		{"templates/en", "Copie templates/en/PRD.md"},
		{"templates/<lang>", "Copie templates/<lang>/TASK.md"},
		{"docs/en", "Veja docs/en/roadmap.html"},
		{"README.en", "Veja README.en.md"},
		{"language: en", "language: en"},
		{"language: pt-br", "language: pt-br"},
		{"External CLI reviewers", "## External CLI reviewers"},
		{"The anti-delegation clause", "## The anti-delegation clause"},
		{"Decisions needed", "## Decisions needed"},
		{"Plan", "## Plan"},
		{"Slice", "### Slice 1 (T-1)"},
		{"Review status", "- Review status: approve"},
		{"Findings", "## Findings"},
		{"SEVERITY", "### <SEVERITY>: <title>"},
		{"Reported by", "- Reported by: claude"},
		{"Comprehension", "## Comprehension"},
		{"Not raised", "## Not raised"},
		{"Guide rule", "- Guide rule: none applies"},
		{"coverage incomplete", "coverage incomplete"},
		{"Checks", "- Checks: `go test`"},
	}
	// the marker is planted on line 8 of an agent, a command, a skill, a template, a profile and an eval case
	places := []struct{ name, rel, prefix string }{
		{"agent", filepath.Join(".agents", "alpha.md"), agentFixture("alpha")},
		{"command", filepath.Join(".commands", "run.md"), commandFixture()},
		{"skill", filepath.Join(".skills", "sample", "SKILL.md"), skillFixture()},
		{"template", filepath.Join("templates", "NOTE.md"), "# Nota\n"},
		{"profile", filepath.Join("profiles", "base.yaml"), "name: base\n"},
		{"eval case", filepath.Join("evals", "cases", "case-a", "prompt.md"), "---\nname: case-a\n---\n"},
	}
	for _, marker := range markers {
		for _, place := range places {
			t.Run(marker.name+" in "+place.name, func(t *testing.T) {
				root := syntheticRoot(t)
				path := filepath.Join(root, place.rel)
				before := ""
				if isFile(path) {
					before = readFixture(t, path)
				} else {
					before = place.prefix
				}
				writeFixture(t, path, before+marker.text+"\n")
				line := strings.Count(before, "\n") + 1
				errs := legacyErrors(t, root)
				want := fmt.Sprintf("%s:%d: marcador do formato antigo em inglês", path, line)
				if !strings.Contains(strings.Join(errs, "\n"), want) {
					t.Fatalf("errors = %v, want substring %q", errs, want)
				}
				writeFixture(t, path, before)
				if errs := legacyErrors(t, root); len(errs) != 0 {
					t.Fatalf("without the marker the source must pass, got %v", errs)
				}
			})
		}
	}
}

// Only the fixed list counts: pt-BR labels, a key that merely ends in "language", the plan title in pt-BR and prose that talks about "uso" do not.
func TestLegacyMarkersIgnoreTheNewFormatAndLookalikes(t *testing.T) {
	root := syntheticRoot(t)
	writeFixture(t, filepath.Join(root, ".agents", "alpha.md"), agentFixture("alpha")+strings.Join([]string{
		"## Plano", "### Fatia 1 (T-1)", "## Achados", "## Compreensão", "## Não levantado", "- Status da revisão: approve", "- Reportado por: claude",
		"- Regra do guia: nenhuma", "- Checagens: `go test`", "conventions.code_language: en | pt-br", "  code_language: pt-br", "Use quando algo falhar",
		"Planejamento e Slices de trabalho", "templates/TASK.md e docs/roadmap.html",
	}, "\n")+"\n")
	if errs := legacyErrors(t, root); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

// History, the CHANGELOG, the glossary and recorded eval results quote the old format on purpose: they are outside the scanned folders.
func TestLegacyMarkersDoNotReadHistory(t *testing.T) {
	root := syntheticRoot(t)
	writeFixture(t, filepath.Join(root, "docs", "prd", "PRD-001.md"), "## Findings\nUse when\n")
	writeFixture(t, filepath.Join(root, "CHANGELOG.md"), "# Changelog\n\n- docs/en removido\n")
	writeFixture(t, filepath.Join(root, "evals", "results", "run.md"), "Review status: approve\n")
	writeFixture(t, filepath.Join(root, "evals", "baselines", "b.md"), "Reported by: claude\n")
	if errs := legacyErrors(t, root); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func legacyErrors(t *testing.T, root string) []string {
	t.Helper()
	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	out := make([]string, 0)
	for _, e := range report.Errors {
		if strings.Contains(e, "marcador do formato antigo") {
			out = append(out, e)
		}
	}
	return out
}
