package kit

import (
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// PRD-035 F2 contract: the instructions that make the scout, the builder, the QA, the Coordinator and /dh:setup produce, write and reference the facts file
// live in prose, so a rename in one place that forgets the others would silently break the flow. This test pins the exact texts, and it runs the examples of
// the skill through the real `dh facts check` so that the format the agents are told to write is the one the checker accepts.

const candidateFormat = "- F-000 · <frase> · Evidência: <caminho:linha>[, <caminho:linha>...] · Visto: <AAAA-MM-DD>, <papel>, <tarefa> · sha: <sha curto> · Estado: candidato"

const (
	repoRefLine   = "Fatos do sistema: .harness/FACTS.md (dica com evidência; o código vence)."
	globalRefLine = `Fatos do sistema: mantidos pelo dev-harness fora deste repositório; caminho em "dh harness-path --json" (campo harness) mais FACTS.md.`
)

func mustContain(t *testing.T, file string, wants ...string) {
	t.Helper()
	text := lerDoRepo(t, file)
	for _, want := range wants {
		if !strings.Contains(text, want) {
			t.Errorf("%s must say %q", file, want)
		}
	}
}

// R8, R9 and the format E1: the candidate line is spelled the same in the skill the scout loads and in the two agents that may propose up to 2.
func TestCandidateFormatIsIdenticalWhereAgentsWriteIt(t *testing.T) {
	for _, f := range []string{".skills/repository-mapping/SKILL.md", ".agents/builder.md", ".agents/qa-verifier.md"} {
		mustContain(t, f, candidateFormat)
	}
}

func TestScoutBuilderAndQAProposeCandidatesButNeverWriteTheFile(t *testing.T) {
	mustContain(t, ".agents/repo-scout.md", "`## Fatos candidatos`", "no máximo 5", "nunca escreve o `FACTS.md`", "grava os 5 primeiros e descarta o resto")
	mustContain(t, ".agents/builder.md", "`## Fatos candidatos` com no máximo 2 linhas", "papel `builder`", "Você não escreve o `FACTS.md`")
	mustContain(t, ".agents/qa-verifier.md", "`## Fatos candidatos` com no máximo 2 linhas", "papel `qa-verifier`", "Você não escreve o `FACTS.md`")
	mustContain(t, ".agents/docs-guide.md", "`Convenções: ver FACTS.md`", "no máximo 5", "nunca escreve o `FACTS.md`")
	// the evidence may be a line or a range N-M, and the three places that tell an agent how to write it say both forms
	for _, f := range []string{".skills/repository-mapping/SKILL.md", ".agents/builder.md", ".agents/qa-verifier.md"} {
		mustContain(t, f, "um intervalo `N-M`", "`src/", ":5-6`")
	}
	mustContain(t, ".skills/repository-mapping/SKILL.md", "no máximo 5 num mapa do `repo-scout`", "no máximo 2 numa resposta do `builder` ou do `qa-verifier`", "## Fatos candidatos")
}

// R10, R1, R2: the Coordinator is the only writer, through the hash lock, with the closed command, the skeleton in the order of R4 and the exit 3 rule.
func TestCoordinatorWritesCandidatesThroughTheLock(t *testing.T) {
	const file = ".agents/coordinator.md"
	mustContain(t, file,
		"dh memory hash --file FACTS.md", "dh memory write --file FACTS.md --expect <hash>", "`--expect absent`", "dh facts check",
		"só na seção `## Candidatos`", "grave os 5 primeiros e descarte o resto", "`F-000`", "próximo id livre de 3 dígitos", "`sha: ?`",
		"Nunca crie `FACTS.md` na raiz do repositório", "No modo `none` não o crie", "A saída 3 significa", "nunca sobrescreva",
		"`<repo>/.harness/FACTS.md`", "`<home>/projects/<nome>/FACTS.md`", "Candidato nunca vai para o `MEMORY.md`")
	// the skeleton lists the 11 categories in the order the checker enforces
	text := lerDoRepo(t, file)
	last := -1
	for _, c := range harness.FactCategories {
		name := strings.ToUpper(c[:1]) + c[1:]
		i := strings.Index(text, name)
		if i < 0 || i < last {
			t.Fatalf("%s: category %q is missing from the skeleton or out of the order of PRD-035 R4", file, name)
		}
		last = i
	}
	if !strings.Contains(text, "`## Candidatos` e `## Mudanças`") {
		t.Errorf("%s: the skeleton must end with the Candidatos and Mudanças sections", file)
	}
}

// The example of the skill passes the real checker, and so does the skeleton the Coordinator is told to create plus the lines it appends.
func TestSkillExampleAndSkeletonPassFactsCheck(t *testing.T) {
	text := lerDoRepo(t, ".skills/repository-mapping/SKILL.md")
	var example []string
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "- F-000 · ") && !strings.Contains(line, "<frase>") {
			example = append(example, line)
		}
	}
	if len(example) != 2 {
		t.Fatalf("the skill example must have 2 candidate lines, found %d", len(example))
	}
	// what the Coordinator writes: the skeleton, with F-000 replaced by free ids and `?` by a sha
	var doc strings.Builder
	doc.WriteString("# Fatos do sistema\n\nÚltima curadoria: nenhuma\n\n")
	for _, c := range harness.FactCategories {
		doc.WriteString("## " + strings.ToUpper(c[:1]) + c[1:] + "\n\n")
	}
	doc.WriteString("## Candidatos\n")
	for i, line := range example {
		doc.WriteString(strings.Replace(line, "F-000", "F-00"+string(rune('1'+i)), 1) + "\n")
	}
	doc.WriteString("\n## Mudanças\n")
	rep := harness.CheckFacts([]byte(doc.String()), func(string) bool { return true })
	if len(rep.Problems) != 0 || rep.Facts != 2 || rep.Candidates != 2 {
		t.Fatalf("the skeleton plus the example must pass dh facts check; facts=%d candidates=%d problems=%v", rep.Facts, rep.Candidates, rep.Problems)
	}
	// with `sha: ?` as a specialist may leave it, the checker flags the line: the Coordinator must fill it before writing
	bad := strings.Replace(doc.String(), "sha: 4fc0ac3", "sha: ?", 1)
	if r := harness.CheckFacts([]byte(bad), func(string) bool { return true }); len(r.Problems) == 0 {
		t.Error("a `sha: ?` line must not pass dh facts check")
	}
}

// R22 to R26: the two reference lines are spelled the same in the command and in the skill, and the rules around them are there.
func TestSetupReferenceLineAndBootstrapSeed(t *testing.T) {
	for _, f := range []string{".commands/setup.md", ".skills/project-onboarding/SKILL.md"} {
		mustContain(t, f, repoRefLine, globalRefLine, "`@import`", "<!-- dh:facts:inicio -->", "`Fatos do sistema:`", "`CLAUDE.md`", "`AGENTS.md`")
	}
	mustContain(t, ".commands/setup.md", "pergunte ao dono, uma vez", "Sem o sim do dono, não toque em arquivo nenhum", "O setup não cria o `FACTS.md`", "`Convenções: ver FACTS.md`", "no máximo 5", "`--expect absent`")
	// the six parts of REPO-MAP.md have fixed titles, so the conventions part is always "## 6. Convenções"
	mustContain(t, ".skills/project-onboarding/SKILL.md", "`## 2. Estrutura de nível superior`", "`## 3. Pontos de entrada e comandos`", "`## 4. Arquitetura`", "`## 5. Áreas de risco`", "`## 6. Convenções`")
	mustContain(t, ".commands/setup.md", "`## 6. Convenções`")
	mustContain(t, ".skills/project-onboarding/SKILL.md", "Sem o sim explícito do dono", "um `AGENTS.md` mínimo", "não acrescente outra", "`Convenções: ver FACTS.md`", "cada um com evidência")
}
