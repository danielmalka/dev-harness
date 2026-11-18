package kit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Label contract (ADR-011, 0.43.0): the skills and the agents that produce a plan or a review report write the labels in pt-BR; whoever READS one of them
// (commands, the coordinator, external-clis) also accepts the English form of a plan or report written before 0.43.0, until 0.46.0. docs/glossario.md section 7 is
// the human-readable table; this test pins it to the files, so a rename in one place that forgets the readers fails here (the bug: implementation-planning wrote
// "## Decisões necessárias" while .commands/plan.md still required "## Decisions needed", and the approval page was never generated).

type rotuloPar struct {
	pt, en     string
	leitores   []string // files that must accept both forms
	escritores []string // files that must write the pt-BR form
}

var paresDeRotulo = []rotuloPar{
	{"## Decisões necessárias", "## Decisions needed", []string{".commands/plan.md"}, []string{".skills/implementation-planning/SKILL.md", ".agents/implementation-planner.md"}},
	{"### Fatia <n> (<id>)", "### Slice <n> (<id>)", []string{".commands/review.md"}, []string{".agents/implementation-planner.md"}},
	{"### Fatia <n> (<id>): <nome>", "### Slice <n> (<id>): <name>", nil, []string{".skills/implementation-planning/SKILL.md"}},
	{"`## Plano`", "`## Plan`", []string{".skills/external-clis/SKILL.md", ".skills/implementation-planning/SKILL.md"}, nil},
	{"## Plano", "## Plan\n", nil, []string{".skills/implementation-planning/SKILL.md", ".agents/implementation-planner.md"}},
	{"Status da revisão:", "Review status:", []string{".agents/coordinator.md", ".skills/external-clis/SKILL.md"}, []string{".skills/code-review/SKILL.md"}},
	{"## Achados", "## Findings", []string{".agents/coordinator.md", ".skills/external-clis/SKILL.md"}, []string{".skills/security-review/SKILL.md", ".skills/document-review/SKILL.md", ".agents/reviewer.md", ".agents/document-validator.md"}},
	{"### <SEVERIDADE>: <título", "### <SEVERITY>: <title", []string{".agents/coordinator.md", ".skills/external-clis/SKILL.md"}, []string{".skills/security-review/SKILL.md"}},
	{"Reportado por:", "Reported by:", []string{".agents/coordinator.md"}, []string{".skills/code-review/SKILL.md", ".skills/security-review/SKILL.md", ".skills/document-review/SKILL.md"}},
	{"## Compreensão", "## Comprehension", []string{".agents/coordinator.md"}, []string{".skills/document-review/SKILL.md", ".agents/document-validator.md"}},
	{"## Não levantado", "## Not raised", nil, []string{".skills/document-review/SKILL.md", ".agents/document-validator.md"}},
	{"Regra do guia:", "Guide rule:", []string{".commands/review.md"}, []string{".skills/code-review/SKILL.md", ".agents/reviewer.md"}},
	{"Checagens:", "Checks:", []string{".commands/plan-loop.md"}, []string{".skills/implementation-planning/SKILL.md"}},
	{"cobertura incompleta", "coverage incomplete", nil, []string{".skills/code-review/SKILL.md", ".skills/security-review/SKILL.md", ".skills/document-review/SKILL.md", ".agents/reviewer.md", ".agents/document-validator.md"}},
}

func lerDoRepo(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestLabelContractReadersAcceptBothFormsAndWritersUsePtBR(t *testing.T) {
	for _, par := range paresDeRotulo {
		for _, f := range par.leitores {
			text := lerDoRepo(t, f)
			if !strings.Contains(text, par.pt) || !strings.Contains(text, strings.TrimSuffix(par.en, "\n")) {
				t.Errorf("%s reads %q and must accept both forms (%q and %q)", f, par.pt, par.pt, par.en)
			}
		}
		for _, f := range par.escritores {
			if !strings.Contains(lerDoRepo(t, f), par.pt) {
				t.Errorf("%s writes %q and must use the pt-BR label %q", f, par.en, par.pt)
			}
		}
	}
}

// A writer must not keep emitting the old English template: the English form may only appear in a line that says it is the legacy form.
func TestLabelContractWritersDoNotTemplateTheEnglishLabels(t *testing.T) {
	legacy := regexp.MustCompile(`(?i)0\.43\.0|older|legacy|antes da 0\.43|anterior à 0\.43|antigo`)
	for _, par := range paresDeRotulo {
		en := strings.TrimSuffix(par.en, "\n")
		for _, f := range par.escritores {
			for i, line := range strings.Split(lerDoRepo(t, f), "\n") {
				if strings.Contains(line, en) && !strings.HasPrefix(en, "## Plan") && !legacy.MatchString(line) {
					t.Errorf("%s:%d still writes the English label %q: %s", f, i+1, en, strings.TrimSpace(line))
				}
			}
		}
	}
}

func TestDocumentStageBundleHasFourItemsEverywhere(t *testing.T) {
	for f, want := range map[string]string{
		".agents/coordinator.md":           "exatamente quatro itens",
		".agents/document-validator.md":    "Estes quatro itens",
		".skills/document-review/SKILL.md": "exatamente quatro itens",
		".skills/external-clis/SKILL.md":   "exatamente quatro itens",
		"AGENTS.md":                        "quatro itens",
	} {
		text := lerDoRepo(t, f)
		if !strings.Contains(text, want) {
			t.Errorf("%s must say %q", f, want)
		}
		for _, old := range []string{"exactly five items", "These five items", "exatamente cinco itens", "cinco itens"} {
			if strings.Contains(text, old) {
				t.Errorf("%s still says %q", f, old)
			}
		}
	}
	if strings.Contains(lerDoRepo(t, ".agents/coordinator.md"), "`language`, harness dir") {
		t.Error("coordinator.md still lists `language` in the document-stage bundle")
	}
}

// Reference reader of the approval-page rule in .commands/plan.md: the headings the command text says it reads are taken FROM THE COMMAND FILE, so removing one
// of the two forms from the command makes the matching plan below stop generating the page.
var decisoesAceitas = regexp.MustCompile("leia a seção `(## [^`]+)` do `PLAN.md` \\(um plano escrito antes da 0\\.43\\.0 tem `(## [^`]+)`:")

func titulosDoComando(t *testing.T) []string {
	t.Helper()
	m := decisoesAceitas.FindStringSubmatch(lerDoRepo(t, ".commands/plan.md"))
	if m == nil {
		t.Fatal(".commands/plan.md no longer states the two accepted headings of the approval-page rule")
	}
	return []string{m[1], m[2]}
}

// generatesApprovalPage is the rule of the command: at least one non-empty item (not just "none"/"nenhuma") under an accepted heading.
func generatesApprovalPage(plan string, titulos []string) bool {
	lines := strings.Split(plan, "\n")
	for i, line := range lines {
		accepted := false
		for _, title := range titulos {
			accepted = accepted || strings.TrimSpace(line) == title
		}
		if !accepted {
			continue
		}
		for _, next := range lines[i+1:] {
			trimmed := strings.TrimSpace(next)
			if strings.HasPrefix(trimmed, "## ") {
				break
			}
			item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if strings.HasPrefix(trimmed, "- ") && item != "" && !strings.EqualFold(item, "none") && !strings.EqualFold(item, "nenhuma") {
				return true
			}
		}
	}
	return false
}

func TestApprovalPageIsGeneratedForBothHeadingForms(t *testing.T) {
	titulos := titulosDoComando(t)
	if titulos[0] != "## Decisões necessárias" || titulos[1] != "## Decisions needed" {
		t.Fatalf("headings read by the command: %q", titulos)
	}
	plan := func(heading, item string) string {
		return "# PLAN\n## Plano: x\n### Fatia 1 (T-1): a\n- Objetivo: y\n" + heading + "\n" + item + "\n## Premissas\n- p\n"
	}
	casos := []struct {
		nome string
		plan string
		want bool
	}{
		{"pt-BR heading with an item", plan("## Decisões necessárias", "- D1: usar fila?, bloqueia a fatia 1, recomendado: sim"), true},
		{"English heading with an item (plan written before 0.43.0)", plan("## Decisions needed", "- D1: use a queue?, blocks slice 1, recommended: yes"), true},
		{"pt-BR heading with only 'nenhuma'", plan("## Decisões necessárias", "- nenhuma"), false},
		{"English heading with only 'none'", plan("## Decisions needed", "- none"), false},
		{"no heading at all", "# PLAN\n## Plano: x\n- Objetivo: y\n", false},
		{"another heading is not read", plan("## Decisões", "- D1: x"), false},
	}
	for _, c := range casos {
		if got := generatesApprovalPage(c.plan, titulos); got != c.want {
			t.Errorf("%s: generates the page = %v, want %v", c.nome, got, c.want)
		}
	}
	// The reference reader only knows what the command tells it: with the English form removed from the accepted list, an old plan stops generating the page.
	if generatesApprovalPage(plan("## Decisions needed", "- D1: x"), titulos[:1]) {
		t.Error("the reference reader must depend on the headings listed in the command")
	}
}
