package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const rmPT = "docs/roadmap.html"

// Card titles copied from the real docs/roadmap.html: the 0.10.0 title holds a <code> element, and "Comando de planejamento em loop" is a detail card with no version prefix.
const (
	cardReal010 = `0.10.0 · <code>/dh:plan-loop</code> e o planejador CLI`
	cardDetail  = `Comando de planejamento em loop`
)

const rmChangelog = "# Changelog\n\n## 0.10.0 — 2026-09-25\n\n- x\n\n## 0.9.0 — 2026-09-24\n\n- y\n\n## 0.8.0 — 2026-09-23\n"

type rmOpts struct {
	date, status, rangeEnd string
	cards                  []string // <h3> inner text, in order
}

func rmPage(o rmOpts) string {
	if o.date == "" {
		o.date = "2026-09-25"
	}
	if o.status == "" {
		o.status = "0.10.0 · vivo"
	}
	if o.rangeEnd == "" {
		o.rangeEnd = "0.10.0"
	}
	if o.cards == nil {
		o.cards = []string{cardReal010, "0.9.0 · B", "0.8.0 · C"}
	}
	var b strings.Builder
	b.WriteString(`<span class="chip">data <b>` + o.date + "</b></span>\n")
	b.WriteString(`<span class="chip">status <b>` + o.status + "</b></span>\n")
	b.WriteString(`<span class="chip">path <b>docs/roadmap.html</b></span>` + "\n")
	b.WriteString(`<h2 id="estrutura">S</h2>` + "\n")
	b.WriteString(`<h2 id="entregue"><span class="n">02</span>E</h2>` + "\n")
	for _, c := range o.cards {
		b.WriteString("<h3>" + c + "</h3>\n<p>x</p>\n")
	}
	b.WriteString(`<h2 id="planejado">P</h2>` + "\n<h3>9.9.9 · não entregue</h3>\n")
	b.WriteString(`<h2 id="fontes">F</h2>` + "\n")
	b.WriteString(`<li><code>CHANGELOG.md</code> — releases 0.1.0 a ` + o.rangeEnd + "</li>\n")
	return b.String()
}

// rmRoot writes the scenario files; an empty changelog or page is not written. `extra` writes more files (relative path -> content).
func rmRoot(t *testing.T, changelog, page string, extra map[string]string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{"CHANGELOG.md": changelog, rmPT: page}
	for rel, content := range extra {
		files[rel] = content
	}
	for rel, content := range files {
		if content == "" {
			continue
		}
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func rmErrs(t *testing.T, changelog string, page string) []string {
	t.Helper()
	var errs []string
	checkRoadmaps(rmRoot(t, changelog, page, nil), &errs)
	return errs
}

func wantErrs(t *testing.T, errs []string, fragments ...string) {
	t.Helper()
	if len(errs) != len(fragments) {
		t.Fatalf("%d erros, esperado %d:\n%s", len(errs), len(fragments), strings.Join(errs, "\n"))
	}
	for i, w := range fragments {
		if !strings.Contains(errs[i], w) {
			t.Errorf("erro %d = %q, deveria conter %q", i, errs[i], w)
		}
	}
}

func TestCheckRoadmaps(t *testing.T) {
	cards := func(c ...string) string { return rmPage(rmOpts{cards: c}) }
	tests := []struct {
		name      string
		changelog string
		page      string
		want      []string // substrings that must all appear; empty means no error
	}{
		{name: "R1 roadmap válido passa", changelog: rmChangelog, page: rmPage(rmOpts{})},
		{name: "R1 sem roadmap passa (pacote ou projeto consumidor)", changelog: rmChangelog},
		{name: "R1 sem CHANGELOG com o roadmap presente", changelog: "-", page: rmPage(rmOpts{}),
			want: []string{"CHANGELOG.md:1: esperado o CHANGELOG.md porque o roadmap existe"}},
		{name: "R1 CHANGELOG sem título de release", changelog: "# Changelog\n", page: rmPage(rmOpts{}),
			want: []string{"CHANGELOG.md:1: esperado um título como"}},
		{name: "R2 chip de data diferente", changelog: rmChangelog, page: rmPage(rmOpts{date: "2026-09-01"}),
			want: []string{rmPT + ":1: esperado chip de data \"2026-09-25\", achei \"2026-09-01\""}},
		{name: "R3 chip de status diferente", changelog: rmChangelog, page: rmPage(rmOpts{status: "0.9.0 · vivo"}),
			want: []string{rmPT + ":2: esperado chip de status \"0.10.0 · vivo\", achei \"0.9.0 · vivo\""}},
		{name: "R3 chip de status com a palavra em inglês", changelog: rmChangelog, page: rmPage(rmOpts{status: "0.10.0 · live"}),
			want: []string{rmPT + ":2: esperado chip de status \"0.10.0 · vivo\", achei \"0.10.0 · live\""}},
		{name: "R4 fim do intervalo de Fontes diferente", changelog: rmChangelog, page: rmPage(rmOpts{rangeEnd: "0.9.0"}),
			want: []string{rmPT + ":", "esperado fim do intervalo de fontes \"0.10.0\", achei \"0.9.0\""}},
		{name: "R5 card faltando cita a linha do h2 entregue", changelog: rmChangelog, page: cards(cardReal010, "0.9.0 · B"),
			want: []string{rmPT + ":5: esperado um card de Entregue para 0.8.0, achei nenhum"}},
		{name: "R5 card duplicado", changelog: rmChangelog, page: cards(cardReal010, "0.9.0 · B", "0.9.0 · B de novo", "0.8.0 · C"),
			want: []string{rmPT + ":", "esperado exatamente um card para 0.9.0, achei um duplicado"}},
		{name: "R5 card extra sem título no CHANGELOG", changelog: rmChangelog, page: cards("0.11.0 · X", cardReal010, "0.9.0 · B", "0.8.0 · C"),
			want: []string{rmPT + ":", "esperado nenhum card para 0.11.0 (sem título no CHANGELOG.md), achei um"}},
		{name: "R6 ordem semântica: 0.9.0 acima de 0.10.0", changelog: rmChangelog, page: cards("0.9.0 · B", cardReal010, "0.8.0 · C"),
			want: []string{rmPT + ":", "achei o card 0.10.0 depois do card 0.9.0"}},
		{name: "R6 card de detalhe entre cards de versão passa", changelog: rmChangelog, page: cards(cardReal010, cardDetail, "0.9.0 · B", "0.8.0 · C")},
		{name: "R6 card de detalhe no topo, no meio e no fim passa", changelog: rmChangelog, page: cards(cardDetail, cardReal010, "0.9.0 · B", cardDetail, "0.8.0 · C", cardDetail)},
		{name: "chip de data com a palavra em inglês não é mais aceito", changelog: rmChangelog, page: strings.Replace(rmPage(rmOpts{}), `chip">data`, `chip">date`, 1),
			want: []string{rmPT + ":1: esperado chip de data \"2026-09-25\", achei nenhum"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cl := tc.changelog
			if cl == "-" {
				cl = ""
			}
			var errs []string
			checkRoadmaps(rmRoot(t, cl, tc.page, nil), &errs)
			got := strings.Join(errs, "\n")
			if len(tc.want) == 0 {
				if got != "" {
					t.Fatalf("esperado nenhum erro, achei:\n%s", got)
				}
				return
			}
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("falta %q em:\n%s", w, got)
				}
			}
		})
	}
}

// ADR-011: the English roadmap is no longer checked. A stale docs/en/roadmap.html (removed in 0.45.0) raises no error and is not compared with the pt one.
func TestCheckRoadmapsIgnoresTheEnglishMirror(t *testing.T) {
	var errs []string
	stale := strings.Replace(rmPage(rmOpts{date: "2020-01-01", cards: []string{"0.1.0 · velho"}}), "data <b>", "date <b>", 1)
	checkRoadmaps(rmRoot(t, rmChangelog, rmPage(rmOpts{}), map[string]string{"docs/en/roadmap.html": stale}), &errs)
	wantErrs(t, errs)
	// only the English roadmap, no pt one: nothing to check
	root := rmRoot(t, rmChangelog, "", map[string]string{"docs/en/roadmap.html": stale})
	checkRoadmaps(root, &errs)
	wantErrs(t, errs)
}

// R6: the comparison is numeric, not lexical (0.10.0 is newer than 0.9.0).
func TestSemverLess(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		want bool
	}{
		{"0.9.0", "0.10.0", true},
		{"0.10.0", "0.9.0", false},
		{"0.10.0", "0.10.0", false},
		{"0.12.0", "0.12.1", true},
		{"1.0.0", "0.99.99", false},
	} {
		if got := semverLess(tc.a, tc.b); got != tc.want {
			t.Errorf("semverLess(%s, %s) = %v, esperado %v", tc.a, tc.b, got, tc.want)
		}
	}
}

// R5/R6: the real 0.10.0 <h3> holds <code>; the parser must still see a card and skip the unversioned detail card that follows it.
func TestParseRoadmapRealMarkup(t *testing.T) {
	d := parseRoadmap(rmPT, rmPage(rmOpts{cards: []string{cardReal010, cardDetail, "0.9.0 · B"}}))
	var versions []string
	for _, c := range d.cards {
		versions = append(versions, c.version)
	}
	if got := strings.Join(versions, ","); got != "0.10.0,0.9.0" {
		t.Fatalf("cards = %q, esperado 0.10.0,0.9.0", got)
	}
	if d.entregueLine != 5 {
		t.Errorf("entregueLine = %d, esperado 5", d.entregueLine)
	}
}

// R6: every failure is "<file>:<line>: esperado ..., achei ...", and a wrong date is reported once, against the changelog.
func TestCheckRoadmapsErrorFormat(t *testing.T) {
	errs := rmErrs(t, rmChangelog, rmPage(rmOpts{date: "2026-01-01"}))
	if len(errs) != 1 || !strings.HasPrefix(errs[0], rmPT+":1: esperado ") {
		t.Fatalf("erros = %q", errs)
	}
}

// A version-shaped <h3> in a non-strict form must be reported, not ignored.
func TestCheckRoadmapsVariantCards(t *testing.T) {
	variants := map[string]struct{ h3, found string }{
		"atributo":   {`<h3 class="x">0.9.0 · B</h3>`, `"0.9.0 · B"`},
		"quebrado":   {"<h3>\n0.9.0 · B</h3>", `" 0.9.0 · B"`},
		"quebrado2":  {"<h3>0.9.0\n· B</h3>", `"0.9.0 · B"`},
		"espaço":     {`<h3> 0.9.0 · B</h3>`, `" 0.9.0 · B"`},
		"nbsp":       {"<h3>\u00a00.9.0 · B</h3>", `"\u00a00.9.0 · B"`},
		"htmlnbsp":   {`<h3>&nbsp;0.9.0 · B</h3>`, `"&nbsp;0.9.0 · B"`},
		"travessão":  {`<h3>0.9.0 — B</h3>`, `"0.9.0 — B"`},
		"code antes": {`<h3><code>0.9.0</code> · B</h3>`, `"0.9.0 · B"`},
	}
	for name, v := range variants {
		t.Run(name, func(t *testing.T) {
			page := strings.Replace(rmPage(rmOpts{}), "<h3>0.8.0 · C</h3>", v.h3+"\n<h3>0.8.0 · C</h3>", 1)
			wantErrs(t, rmErrs(t, rmChangelog, page), rmPT+":10: esperado \"<h3>N.N.N · \", achei "+v.found)
		})
	}
	page := strings.Replace(rmPage(rmOpts{}), "<h3>0.8.0 · C</h3>", "<h3>Card de detalhe</h3>\n<h3>0.8.0 · C</h3>", 1)
	wantErrs(t, rmErrs(t, rmChangelog, page)) // a detail card without a version stays ignored
}

func TestCheckRoadmapsMissingPieces(t *testing.T) {
	t.Run("sem h2 entregue", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog, strings.Replace(rmPage(rmOpts{}), `<h2 id="entregue">`, `<h2 id="outro">`, 1)),
			rmPT+`:1: esperado <h2 id="entregue">, achei nenhum`)
	})
	t.Run("chip de status ausente", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog, strings.Replace(rmPage(rmOpts{}), `<span class="chip">status`, `<span class="chp">status`, 1)),
			rmPT+`:1: esperado chip de status "0.10.0 · vivo", achei nenhum`)
	})
	t.Run("chip de caminho ausente", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog, strings.Replace(rmPage(rmOpts{}), `<span class="chip">path`, `<span class="chp">path`, 1)),
			rmPT+":1: esperado o chip de caminho, achei nenhum")
	})
	t.Run("card depois do h2 de fechamento não conta", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog, strings.Replace(rmPage(rmOpts{}), "<h3>9.9.9 · não entregue</h3>", "<h3>0.9.0 · de novo</h3>", 1)))
	})
	t.Run("arquivos CRLF", func(t *testing.T) {
		wantErrs(t, rmErrs(t, strings.ReplaceAll(rmChangelog, "\n", "\r\n"), strings.ReplaceAll(rmPage(rmOpts{}), "\n", "\r\n")))
	})
	t.Run("roadmap ilegível", func(t *testing.T) {
		root := rmRoot(t, rmChangelog, "x", nil)
		if err := os.Remove(filepath.Join(root, "docs", "roadmap.html")); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(root, "docs", "roadmap.html"), 0o755); err != nil {
			t.Fatal(err)
		}
		var errs []string
		checkRoadmaps(root, &errs)
		wantErrs(t, errs, rmPT+":1: esperado ler o arquivo do roadmap")
	})
}

func TestCheckRoadmapsChangelogHeadings(t *testing.T) {
	page := rmPage(rmOpts{})
	t.Run("título repetido mantém a primeira data", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## 0.10.0 — 2020-01-01\n", page), "CHANGELOG.md:13: esperado um título por release, achei 0.10.0 de novo")
	})
	t.Run("títulos com colchetes e com v", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## [0.13.0] — 2026-09-26\n", page), "CHANGELOG.md:13: esperado um título como")
		wantErrs(t, rmErrs(t, rmChangelog+"\n## v0.13.0 — 2026-09-26\n", page), "CHANGELOG.md:13: esperado um título como")
	})
	t.Run("título Unreleased é permitido", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## [Unreleased]\n", page))
	})
	t.Run("título quase certo", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## 0.7.0 - 2026-09-20\n", page), "CHANGELOG.md:13: esperado um título como")
	})
}

// The real repository must satisfy its own check.
func TestCheckRoadmapsRealRepo(t *testing.T) {
	var errs []string
	checkRoadmaps("../..", &errs)
	if len(errs) != 0 {
		t.Fatalf("erros do roadmap real:\n%s", strings.Join(errs, "\n"))
	}
}
