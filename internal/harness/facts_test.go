package harness

import (
	"fmt"
	"strings"
	"testing"
)

var existing = map[string]bool{"cmd/dh/main.go": true, "internal/harness/memory.go": true, "README.md": true}

func exists(rel string) bool { return existing[rel] }

const goodFact = "- F-001 · O binário usa lock por arquivo · Evidência: internal/harness/memory.go:40 · Visto: 2026-11-18, builder, T-12 · sha: abc1234 · Estado: confirmado"

func factsDoc(sections ...string) string {
	return "# Fatos do sistema\n\nÚltima curadoria: 2026-11-18, sha abc1234\n\n" + strings.Join(sections, "\n")
}

func check(content string) FactsReport { return CheckFacts([]byte(content), exists) }

func messages(r FactsReport) string {
	var b strings.Builder
	for _, p := range r.Problems {
		fmt.Fprintf(&b, "%d: %s\n", p.Line, p.Msg)
	}
	return b.String()
}

func wantProblem(t *testing.T, name, content, line, want string) {
	t.Helper()
	r := check(content)
	for _, p := range r.Problems {
		if (line == "" || fmt.Sprint(p.Line) == line) && strings.Contains(p.Msg, want) {
			return
		}
	}
	t.Errorf("%s: want a problem on line %q containing %q, got:\n%s", name, line, want, messages(r))
}

// R27: a good file passes with the counts; an empty file passes too.
func TestCheckFactsOK(t *testing.T) {
	doc := factsDoc("## Arquitetura", goodFact, "", "## Armadilhas",
		strings.Replace(strings.Replace(goodFact, "F-001", "F-002", 1), "confirmado", "hipótese", 1), "## Candidatos",
		strings.Replace(strings.Replace(goodFact, "F-001", "F-003", 1), "confirmado", "candidato", 1), "", "## Mudanças", "- 2026-11-18 · abc1234 · 1 entrou, 0 saíram, 0 fundidos", "texto livre aqui")
	r := check(doc)
	if len(r.Problems) != 0 || r.Facts != 3 || r.Candidates != 1 {
		t.Fatalf("facts=%d candidates=%d problems:\n%s", r.Facts, r.Candidates, messages(r))
	}
	if r := check(""); len(r.Problems) != 0 || r.Facts != 0 {
		t.Fatalf("empty: %v", r)
	}
	if r := check(strings.ReplaceAll(doc, "\n", "\r\n")); len(r.Problems) != 0 {
		t.Fatalf("CRLF: %s", messages(r))
	}
}

// R3: each missing field is reported with the line and the field name.
func TestCheckFactsMissingFields(t *testing.T) {
	cases := map[string]struct{ line, field string }{
		"frase":     {"- F-001 · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · sha: abc1234 · Estado: confirmado", "campo frase ausente"},
		"Evidência": {"- F-001 · frase · Visto: 2026-11-18, builder, T-1 · sha: abc1234 · Estado: confirmado", "campo Evidência: ausente"},
		"Visto":     {"- F-001 · frase · Evidência: README.md:1 · sha: abc1234 · Estado: confirmado", "campo Visto: ausente"},
		"sha":       {"- F-001 · frase · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · Estado: confirmado", "campo sha ausente"},
		"Estado":    {"- F-001 · frase · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · sha: abc1234", "campo Estado: ausente"},
	}
	for name, c := range cases {
		wantProblem(t, name, factsDoc("## Testes", c.line), "6", c.field)
	}
	wantProblem(t, "id", factsDoc("## Testes", "- F-1 · frase · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · sha: abc1234 · Estado: confirmado"), "6", "campo id")
	wantProblem(t, "out of order", factsDoc("## Testes", "- F-001 · frase · Visto: 2026-11-18, builder, T-1 · Evidência: README.md:1 · sha: abc1234 · Estado: confirmado"), "6", "campos fora de ordem")
	if r := check(factsDoc("## Testes", "- F-001 · frase · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · sha abc1234 · Estado: confirmado")); len(r.Problems) != 0 {
		t.Errorf("sha without colon must pass: %s", messages(r))
	}
}

func TestCheckFactsFieldValues(t *testing.T) {
	line := func(evidence, seen, sha, state string) string {
		return fmt.Sprintf("- F-001 · frase · Evidência: %s · Visto: %s · sha: %s · Estado: %s", evidence, seen, sha, state)
	}
	good := [4]string{"README.md:1", "2026-11-18, builder, T-1", "abc1234", "confirmado"}
	mk := func(i int, v string) string { g := good; g[i] = v; return line(g[0], g[1], g[2], g[3]) }
	wantProblem(t, "evidence none", factsDoc("## Testes", mk(0, "")), "6", "de 1 a 3")
	wantProblem(t, "evidence four", factsDoc("## Testes", mk(0, "README.md:1, README.md:2, README.md:3, README.md:4")), "6", "tem 4")
	wantProblem(t, "evidence without line", factsDoc("## Testes", mk(0, "README.md")), "6", "caminho:linha")
	wantProblem(t, "evidence line zero", factsDoc("## Testes", mk(0, "README.md:0")), "6", "linha maior que zero")
	wantProblem(t, "evidence absolute", factsDoc("## Testes", mk(0, "/etc/passwd:1")), "6", "relativo ao repositório")
	wantProblem(t, "evidence dotdot", factsDoc("## Testes", mk(0, "../x/README.md:1")), "6", "relativo ao repositório")
	wantProblem(t, "evidence backslash", factsDoc("## Testes", mk(0, `cmd\dh\main.go:1`)), "6", "relativo ao repositório")
	wantProblem(t, "evidence missing path", factsDoc("## Testes", mk(0, "cmd/dh/nao-existe.go:3")), "6", "não existe no repositório")
	if r := check(factsDoc("## Testes", mk(0, "README.md:1, cmd/dh/main.go:2, internal/harness/memory.go:3"))); len(r.Problems) != 0 {
		t.Errorf("three paths must pass: %s", messages(r))
	}
	wantProblem(t, "seen short", factsDoc("## Testes", mk(1, "2026-11-18, builder")), "6", "campo Visto:")
	wantProblem(t, "seen bad date", factsDoc("## Testes", mk(1, "2026-13-40, builder, T-1")), "6", "fora do formato AAAA-MM-DD")
	wantProblem(t, "seen dd/mm", factsDoc("## Testes", mk(1, "18/11/2026, builder, T-1")), "6", "fora do formato AAAA-MM-DD")
	wantProblem(t, "sha upper", factsDoc("## Testes", mk(2, "ABC1234")), "6", "não é um sha curto")
	wantProblem(t, "sha short", factsDoc("## Testes", mk(2, "abc12")), "6", "não é um sha curto")
	wantProblem(t, "state", factsDoc("## Testes", mk(3, "pendente")), "6", "campo Estado:")
	for _, s := range FactStates {
		if r := check(factsDoc("## Testes", mk(3, s))); len(r.Problems) != 0 {
			t.Errorf("state %q must pass: %s", s, messages(r))
		}
	}
}

func TestCheckFactsStructure(t *testing.T) {
	wantProblem(t, "unknown section", factsDoc("## Qualquer coisa", goodFact), "5", "não é uma das categorias")
	wantProblem(t, "duplicate id", factsDoc("## Testes", goodFact, "## Armadilhas", goodFact), "8", "F-001 repetido")
	wantProblem(t, "duplicate section", factsDoc("## Testes", goodFact, "## Testes"), "7", "seção repetida")
	wantProblem(t, "out of order", factsDoc("## Fronteiras", "## Arquitetura"), "6", "fora de ordem")
	wantProblem(t, "not a fact line", factsDoc("## Testes", "- texto solto sem id"), "6", "linha que não é um fato")
	wantProblem(t, "prose in a category", factsDoc("## Testes", "um parágrafo qualquer"), "6", "linha que não é um fato")
	for _, c := range FactCategories {
		heading := "## " + strings.ToUpper(c[:1]) + c[1:]
		if r := check(factsDoc(heading)); len(r.Problems) != 0 {
			t.Errorf("category %q must pass: %s", heading, messages(r))
		}
	}
	if r := check(factsDoc("## TESTES", "## Mudanças", "qualquer texto, - F-xx · solto", "## Candidatos")); len(r.Problems) != 0 {
		t.Errorf("case-insensitive headings and free text in Mudanças: %s", messages(r))
	}
}

// R5: each ceiling passes at the limit and fails one past it, and the message names the ceiling.
func TestCheckFactsCeilings(t *testing.T) {
	fact := func(i int) string { return strings.Replace(goodFact, "F-001", fmt.Sprintf("F-%03d", i), 1) }
	var facts []string
	for i := 1; i <= FactsMaxPerCategory; i++ {
		facts = append(facts, fact(i))
	}
	if r := check(factsDoc(append([]string{"## Testes"}, facts...)...)); len(r.Problems) != 0 {
		t.Fatalf("15 facts in a category must pass: %s", messages(r))
	}
	wantProblem(t, "16 facts", factsDoc(append(append([]string{"## Testes"}, facts...), fact(16))...), "", "teto de fatos por categoria excedido")

	// 300 characters per fact: the whole line counts.
	pad := func(n int) string {
		base := "- F-001 · %s · Evidência: README.md:1 · Visto: 2026-11-18, builder, T-1 · sha: abc1234 · Estado: confirmado"
		empty := len([]rune(fmt.Sprintf(base, "")))
		return fmt.Sprintf(base, strings.Repeat("á", n-empty))
	}
	if l := pad(300); len([]rune(l)) != 300 {
		t.Fatalf("test helper: %d", len([]rune(l)))
	}
	if r := check(factsDoc("## Testes", pad(300))); len(r.Problems) != 0 {
		t.Fatalf("300 characters must pass: %s", messages(r))
	}
	wantProblem(t, "301 characters", factsDoc("## Testes", pad(301)), "6", "teto de caracteres por fato excedido")

	// 150 lines in the file.
	body := func(total int) string {
		lines := []string{"## Mudanças"}
		for len(lines) < total {
			lines = append(lines, "- linha de mudança")
		}
		return strings.Join(lines, "\n") + "\n"
	}
	if r := check(body(150)); len(r.Problems) != 0 {
		t.Fatalf("150 lines must pass: %s", messages(r))
	}
	wantProblem(t, "151 lines", body(151), "0", "teto de linhas excedido")

	// 12 KB in the file.
	big := func(size int) string {
		s := "## Mudanças\n"
		for len(s) < size {
			s += strings.Repeat("m", 99) + "\n"
		}
		return s[:size]
	}
	if r := check(big(FactsMaxBytes)); len(r.Problems) != 0 {
		t.Fatalf("12 KB must pass: %s", messages(r))
	}
	wantProblem(t, "12 KB + 1", big(FactsMaxBytes+1), "0", "teto de tamanho excedido")
}

// R6: every secret shape is reported with its kind and line, and the value is never echoed.
func TestCheckFactsSecrets(t *testing.T) {
	cases := map[string]string{
		"chave privada":              "-----BEGIN RSA PRIVATE KEY-----",
		"chave de acesso da AWS":     "AKIAABCDEFGHIJKLMNOP",
		"token do GitHub":            "ghp_abcdefghijklmnopqrstuvwxyz0123456789",
		"token do Slack":             "xoxb-123456789012-abcdefghij",
		"chave de API (sk-)":         "sk-abcdefghijklmnopqrstuvwx",
		"chave de API do Google":     "AIzaSyA1234567890abcdefghijklmnopqrstuv",
		"JWT":                        "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abcdefghijklmnop",
		"credencial em URL":          "postgres://app:s3cr3tpass@db.local/x",
		"valor de segredo atribuído": "API_KEY=abc12345xyz",
		"e-mail":                     "fulano@example.com",
		"CPF":                        "123.456.789-09",
		"CNPJ":                       "12.345.678/0001-95",
	}
	for kind, secret := range cases {
		doc := factsDoc("## Armadilhas", strings.Replace(goodFact, "O binário usa lock por arquivo", "valor "+secret+" visto", 1))
		r := check(doc)
		found := false
		for _, p := range r.Problems {
			if strings.Contains(p.Msg, "padrão de segredo") && strings.Contains(p.Msg, kind) && p.Line == 6 {
				found = true
			}
			if strings.Contains(p.Msg, secret) {
				t.Errorf("%s: the message echoes the secret: %s", kind, p.Msg)
			}
		}
		if !found {
			t.Errorf("%s not reported:\n%s", kind, messages(r))
		}
	}
	// A secret in a free-text line of Mudanças is reported too: the file is versioned.
	wantProblem(t, "secret in Mudanças", factsDoc("## Mudanças", "- trocado token: abcd1234efgh"), "6", "padrão de segredo")
	// Look-alikes that are not secrets pass.
	for _, ok := range []string{"o campo senha: bcrypt", "usa a variável API_KEY sem valor", "token: <valor>", "sk-curto", "veja ghp_ no texto", "um sha: abc1234 qualquer"} {
		if r := check(factsDoc("## Armadilhas", strings.Replace(goodFact, "O binário usa lock por arquivo", ok, 1))); len(r.Problems) != 0 {
			t.Errorf("%q must not be reported: %s", ok, messages(r))
		}
	}
	if r := check(factsDoc("## Armadilhas", strings.Replace(goodFact, "O binário usa lock por arquivo", "valor [redigido] no lugar do segredo", 1))); len(r.Problems) != 0 {
		t.Errorf("[redigido] must pass: %s", messages(r))
	}
}
