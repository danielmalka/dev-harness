package kit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Label contract (ADR-011): the skills and the agents that produce a plan or a review report write the labels in pt-BR, and whoever reads one of them (commands,
// the coordinator, external-clis) reads only the pt-BR form. From 0.43.0 to 0.45.0 the readers also accepted the English form; 0.46.0 ended that and the validate lock
// (legacy_markers.go) rejects the old English labels in the kit sources. docs/glossario.md section 7 is the human-readable table; this test pins it to the files, so a
// rename in one place that forgets the readers fails here (the bug: implementation-planning wrote "## Decisões necessárias" while .commands/plan.md still required
// "## Decisions needed", and the approval page was never generated).

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

func TestLabelContractReadersAndWritersUseOnlyPtBR(t *testing.T) {
	for _, par := range paresDeRotulo {
		for _, f := range append(append([]string{}, par.leitores...), par.escritores...) {
			text := lerDoRepo(t, f)
			if !strings.Contains(text, par.pt) {
				t.Errorf("%s uses the label %q and must write it in pt-BR: %q", f, par.en, par.pt)
			}
			if rotuloEmIngles(par).MatchString(text) {
				t.Errorf("%s still carries the old English label %q: since 0.46.0 only the pt-BR form %q is read and written", f, strings.TrimSuffix(par.en, "\n"), par.pt)
			}
		}
	}
}

// rotuloEmIngles matches the English form of a pair; the "## Plan" title is matched as a whole word so that "## Plano" does not count.
func rotuloEmIngles(par rotuloPar) *regexp.Regexp {
	en := strings.TrimSuffix(par.en, "\n")
	if en == "## Plan" {
		return regexp.MustCompile(`## Plan\b`)
	}
	return regexp.MustCompile(regexp.QuoteMeta(en))
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

// Reference reader of the approval-page rule in .commands/plan.md: the heading the command text says it reads is taken FROM THE COMMAND FILE, so changing the heading in
// the command makes the plans below stop generating the page. Since 0.46.0 only the pt-BR heading is read.
var decisoesAceitas = regexp.MustCompile("leia a seção `(## [^`]+)` do `PLAN.md`")

func tituloDoComando(t *testing.T) string {
	t.Helper()
	m := decisoesAceitas.FindStringSubmatch(lerDoRepo(t, ".commands/plan.md"))
	if m == nil {
		t.Fatal(".commands/plan.md no longer states the heading of the approval-page rule")
	}
	return m[1]
}

// generatesApprovalPage is the rule of the command: at least one non-empty item (not just "nenhuma") under the heading.
func generatesApprovalPage(plan string, titulo string) bool {
	lines := strings.Split(plan, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != titulo {
			continue
		}
		for _, next := range lines[i+1:] {
			trimmed := strings.TrimSpace(next)
			if strings.HasPrefix(trimmed, "## ") {
				break
			}
			item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if strings.HasPrefix(trimmed, "- ") && item != "" && !strings.EqualFold(item, "nenhuma") {
				return true
			}
		}
	}
	return false
}

func TestApprovalPageIsGeneratedOnlyForThePtBRHeading(t *testing.T) {
	titulo := tituloDoComando(t)
	if titulo != "## Decisões necessárias" {
		t.Fatalf("heading read by the command: %q", titulo)
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
		{"pt-BR heading with only 'nenhuma'", plan("## Decisões necessárias", "- nenhuma"), false},
		{"English heading of a plan written before 0.43.0 is no longer read", plan("## Decisions needed", "- D1: use a queue?, blocks slice 1, recommended: yes"), false},
		{"no heading at all", "# PLAN\n## Plano: x\n- Objetivo: y\n", false},
		{"another heading is not read", plan("## Decisões", "- D1: x"), false},
	}
	for _, c := range casos {
		if got := generatesApprovalPage(c.plan, titulo); got != c.want {
			t.Errorf("%s: generates the page = %v, want %v", c.nome, got, c.want)
		}
	}
}
