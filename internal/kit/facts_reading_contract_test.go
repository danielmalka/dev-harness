package kit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// PRD-035 F4 contract (R17 to R21): the reading lives in prose too, so a table, one dispatch line and one answer line are fixed here, and each of them is checked in every
// file that has to repeat it. The Verificador proves in a real session that the agents follow the prose; this test keeps the files from drifting apart.

const (
	factsDispatchLine = "Fatos: <harness dir>/FACTS.md · categorias: <lista> · só confirmado"
	factsLabel        = "· rótulo: dica, não regra"
	contradictionLine = "Fato contradito: F-<n>, <caminho:linha>"
)

// factsReaders are the roles of the table of R17, each with the file that has to carry the reading rule.
var factsReaders = []string{"repo-scout", "implementation-planner", "architect", "builder", "qa-verifier", "reviewer", "debugger"}

// factsRoleTable is the table of R17 as the PRD writes it. The coordinator's table must be this one, row by row.
func factsRoleTable(t *testing.T, file string) map[string]string {
	t.Helper()
	text := lerDoRepo(t, file)
	i := strings.Index(text, "| Papel | Categorias |")
	if i < 0 {
		t.Fatalf("%s: the table `| Papel | Categorias |` is missing", file)
	}
	rows := map[string]string{}
	for _, line := range strings.Split(text[i:], "\n")[2:] {
		if !strings.HasPrefix(line, "|") {
			break
		}
		cells := strings.Split(strings.Trim(line, "| "), "|")
		if len(cells) != 2 {
			t.Fatalf("%s: row with %d cells: %q", file, len(cells), line)
		}
		rows[strings.TrimSpace(cells[0])] = strings.TrimSpace(cells[1])
	}
	return rows
}

// R17: the Coordinator's table is the PRD's table, and every category in it is one of the 11 that `dh facts check` accepts.
func TestFactsReadingTableMatchesThePRDAndTheCategories(t *testing.T) {
	prd := factsRoleTable(t, "docs/prd/PRD-035-arquivo-de-fatos.md")
	coord := factsRoleTable(t, ".agents/coordinator.md")
	if len(prd) != 7 {
		t.Fatalf("the table of R17 has 7 rows, found %d: %v", len(prd), prd)
	}
	for role, cats := range prd {
		if coord[role] != cats {
			t.Errorf("coordinator.md row %q is %q, the PRD says %q", role, coord[role], cats)
		}
	}
	for role := range coord {
		if _, ok := prd[role]; !ok {
			t.Errorf("coordinator.md has a row %q that the PRD does not", role)
		}
	}
	valid := map[string]bool{}
	for _, c := range harness.FactCategories {
		valid[strings.ToLower(c)] = true
	}
	for role, cats := range coord {
		list := strings.TrimSuffix(strings.SplitN(cats, ", com o rótulo", 2)[0], ".")
		if list == "todas" {
			continue
		}
		for _, c := range strings.Split(list, ", ") {
			if !valid[strings.ToLower(strings.TrimSpace(c))] {
				t.Errorf("row %q: %q is not one of the 11 categories of R4", role, c)
			}
		}
	}
	// a row per kind of builder, so a lane never reads another lane's categories
	for _, row := range []string{"`builder`, lane `frontend`", "`builder`, lanes `backend`, `dados`, `infra`"} {
		if coord[row] == "" {
			t.Errorf("coordinator.md must have the builder row %s", row)
		}
	}
	// the reviewer reads every category, and the rest of the table never does
	if coord["`reviewer` (modo `code`)"] != `todas, com o rótulo "dica, não regra"` || coord["`repo-scout`"] != "todas" {
		t.Errorf("only the repo-scout and the reviewer read all the categories: %v", coord)
	}
}

// R17, R18, R21: the Coordinator puts one line in `sources` of the roles of the table, reads nothing, and keeps FACTS.md out of the document bundle.
func TestCoordinatorDispatchesTheFactsLineByRole(t *testing.T) {
	mustContain(t, ".agents/coordinator.md",
		"ponha em `sources` do despacho de cada papel desta tabela uma linha `"+factsDispatchLine+"`", "não leia o arquivo para isso", factsLabel,
		"(`dh facts check` imprime `ok …` e não `absent`", "Quem lê usa só fatos `confirmado`", "o fato é dica, e o código vence",
		"`reviewer` em modo `security`", "o pacote de uma etapa de documento e o de um revisor `cli:` não levam a linha",
		"O `FACTS.md` não entra nesse pacote (PRD-035, R21)")
	// the pre-existing rule of the bundle stays at four items, and the facts file is not one of them
	mustContain(t, ".agents/coordinator.md", "O pacote de uma etapa de documento tem exatamente quatro itens: caminho do documento, material-fonte")
	mustNotContain(t, ".skills/document-review/SKILL.md", "FACTS.md")
	mustNotContain(t, ".agents/document-validator.md", "FACTS.md")
}

// R17 to R20: every reader repeats the same input bullet and the same output line, so the dispatch line means one thing to all of them.
func TestFactsReadersShareTheReadingRule(t *testing.T) {
	const input = "- Fatos, opcional: a linha `Fatos: <caminho do FACTS.md> · categorias: <lista> · só confirmado`. Leia o arquivo e use só as categorias nomeadas, e só os fatos com `Estado: confirmado`; um `candidato` ou uma `hipótese` não é fato para você. O fato é dica, não regra: o código vence (PRD-035)."
	const output = "Fato contradito: quando o código que você leu diz outra coisa que um fato que você usou, acrescente ao fim da resposta uma linha por fato, `" + contradictionLine + "` (a evidência no código). Você não edita o `FACTS.md`: o Coordenador grava a linha e o curador decide na rodada seguinte (PRD-035, R20)."
	for _, role := range factsReaders {
		file := ".agents/" + role + ".md"
		text := lerDoRepo(t, file)
		in := strings.Index(text, input)
		out := strings.Index(text, output)
		if in < 0 || out < 0 {
			t.Errorf("%s must carry the reading bullet (found %v) and the `Fato contradito` line (found %v)", file, in >= 0, out >= 0)
			continue
		}
		if sec := strings.Index(text, "## Entradas mínimas"); in < sec || in > strings.Index(text[sec:], "\n## Procedimento")+sec {
			t.Errorf("%s: the reading bullet belongs to `## Entradas mínimas`", file)
		}
		if sec := strings.Index(text, "## Formato de saída"); out < sec {
			t.Errorf("%s: the `Fato contradito` line belongs to `## Formato de saída`", file)
		}
	}
	// the roles outside the table do not carry it: they get no `Fatos:` line (docs-guide curates, it does not read for a task)
	for _, role := range []string{"docs-guide", "release-manager", "document-validator", "product-discovery", "harness-maintainer"} {
		mustNotContain(t, ".agents/"+role+".md", input)
	}
	// nobody but the Coordinator and the curator writes the file: the readers say so in the output line itself, and none of them lists FACTS.md in its write set
	for _, role := range factsReaders {
		text := lerDoRepo(t, ".agents/"+role+".md")
		if regexp.MustCompile("(?i)(arquivos permitidos|conjunto de escrita|write set)[^\\n]*FACTS\\.md").MatchString(text) {
			t.Errorf(".agents/%s.md must not put FACTS.md in a write set", role)
		}
	}
}

// R19: the reviewer does not open a finding only because the code departs from a fact; it files `Fato contradito` when the code contradicts one.
func TestReviewerDoesNotOpenAFindingOnlyBecauseOfAFact(t *testing.T) {
	mustContain(t, ".agents/reviewer.md",
		"- O fato do `FACTS.md` é dica, não regra: não abra achado só porque o código desvia de um fato.",
		"Um achado cita evidência no código (local e cenário)", "devolva a linha `Fato contradito`, não um achado (PRD-035, R19)")
	mustContain(t, ".agents/coordinator.md", factsLabel)
}

// R20: the Coordinator records the contradiction as a line of `## Mudanças` that `dh facts check` accepts, and the curator reads it in the next round.
func TestContradictionIsRecordedAndSettledByTheCurator(t *testing.T) {
	mustContain(t, ".agents/coordinator.md",
		"Fato contradito (R20): quando a resposta de um leitor traz a linha `"+contradictionLine+"`", "no fim de `## Mudanças`, a linha `- <AAAA-MM-DD> · contradito · F-<n>, <caminho:linha> · <papel>, <tarefa>`",
		"sem repetir uma contradição que já está aberta", "sem inventar um id que não está no arquivo", "Não mude o estado do fato nem o remova: quem decide é o curador, na rodada seguinte")
	mustContain(t, ".skills/facts-curation/SKILL.md",
		"uma linha `· contradito · F-<n>, <caminho:linha> ·` de `## Mudanças` está aberta quando nenhuma linha posterior diz `F-<n> saiu` ou `F-<n> mantido`",
		"`F-<n> mantido: <motivo>; evidência: <caminho:linha>`", "Os contraditos não contam para o limite de 10", "As linhas de contradição de `## Mudanças` ainda abertas")
	// the line the Coordinator writes is free text for the checker, even with a path that looks like a secret-free fact, and the file stays valid
	doc := "# Fatos do sistema\n\nÚltima curadoria: nenhuma\n\n## Candidatos\n\n## Mudanças\n- 2026-11-18 · contradito · F-004, src/ui/Button.tsx:8 · builder, T-31\n- 2026-11-20 · sha 7a1b2c3 · entraram 0 · saíram 0 · fundidos 0\n  - F-004 mantido: o `Button` continua no código; evidência: src/ui/Button.tsx:8\n"
	if r := harness.CheckFacts([]byte(doc), func(string) bool { return true }); len(r.Problems) != 0 {
		t.Errorf("the contradiction and the `mantido` lines must pass the checker: %v", r.Problems)
	}
}
