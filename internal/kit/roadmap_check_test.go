package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	rmPT = "docs/roadmap.html"
	rmEN = "docs/en/roadmap.html"
)

// Card titles copied from the real docs/roadmap.html: the 0.10.0 title holds a
// <code> element, and "Comando de planejamento em loop" is a detail card with
// no version prefix.
const (
	cardReal010 = `0.10.0 · <code>/dh:plan-loop</code> e o planejador CLI`
	cardDetail  = `Comando de planejamento em loop`
)

const rmChangelog = "# Changelog\n\n## 0.10.0 — 2026-09-25\n\n- x\n\n## 0.9.0 — 2026-09-24\n\n- y\n\n## 0.8.0 — 2026-09-23\n"

type rmOpts struct {
	date, status, rangeEnd string
	cards                  []string // <h3> inner text, in order
}

func rmPage(lang string, o rmOpts) string {
	dateW, word, to := "data", "vivo", "a"
	if lang == "en" {
		dateW, word, to = "date", "live", "to"
	}
	if o.date == "" {
		o.date = "2026-09-25"
	}
	if o.status == "" {
		o.status = "0.10.0 · " + word
	}
	if o.rangeEnd == "" {
		o.rangeEnd = "0.10.0"
	}
	if o.cards == nil {
		o.cards = []string{cardReal010, "0.9.0 · B", "0.8.0 · C"}
	}
	var b strings.Builder
	b.WriteString(`<span class="chip">` + dateW + ` <b>` + o.date + "</b></span>\n")
	b.WriteString(`<span class="chip">status <b>` + o.status + "</b></span>\n")
	b.WriteString(`<span class="chip">path <b>docs/roadmap.html</b></span>` + "\n")
	b.WriteString(`<h2 id="estrutura">S</h2>` + "\n")
	b.WriteString(`<h2 id="entregue"><span class="n">02</span>E</h2>` + "\n")
	for _, c := range o.cards {
		b.WriteString("<h3>" + c + "</h3>\n<p>x</p>\n")
	}
	b.WriteString(`<h2 id="planejado">P</h2>` + "\n<h3>9.9.9 · not delivered</h3>\n")
	b.WriteString(`<h2 id="fontes">F</h2>` + "\n")
	b.WriteString(`<li><code>CHANGELOG.md</code> — releases 0.1.0 ` + to + ` ` + o.rangeEnd + "</li>\n")
	return b.String()
}

// rmRoot writes the fixture files; a nil page or empty changelog is not written.
func rmRoot(t *testing.T, changelog, pt, en string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range map[string]string{"CHANGELOG.md": changelog, rmPT: pt, rmEN: en} {
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

func rmRun(root string) string {
	var errs []string
	checkRoadmaps(root, &errs)
	return strings.Join(errs, "\n")
}

func TestCheckRoadmaps(t *testing.T) {
	ok := func(o rmOpts) (string, string) { return rmPage("pt", o), rmPage("en", o) }
	cards := func(c ...string) rmOpts { return rmOpts{cards: c} }
	tests := []struct {
		name      string
		changelog string
		pt, en    string
		want      []string // substrings that must all appear; empty means no error
	}{
		{name: "R1 valid pair passes", changelog: rmChangelog},
		{name: "R1 both roadmaps absent passes", changelog: rmChangelog, pt: "-", en: "-"},
		{name: "R1 only pt present names the missing en file", changelog: rmChangelog, en: "-",
			want: []string{rmEN + ":1: expected the roadmap file to exist"}},
		{name: "R1 only en present names the missing pt file", changelog: rmChangelog, pt: "-",
			want: []string{rmPT + ":1: expected the roadmap file to exist"}},
		{name: "R1 missing CHANGELOG with roadmaps present", changelog: "-",
			want: []string{"CHANGELOG.md:1: expected CHANGELOG.md to exist"}},
		{name: "R1 CHANGELOG without a release heading", changelog: "# Changelog\n",
			want: []string{"CHANGELOG.md:1: expected a heading"}},
		{name: "R2 pt date chip differs", changelog: rmChangelog,
			pt: "date", want: []string{rmPT + ":1: expected date chip \"2026-09-25\", found \"2026-09-01\""}},
		{name: "R2 en date chip differs", changelog: rmChangelog,
			en: "date", want: []string{rmEN + ":1: expected date chip \"2026-09-25\", found \"2026-09-01\""}},
		{name: "R3 pt status chip differs", changelog: rmChangelog,
			pt: "status", want: []string{rmPT + ":2: expected status chip \"0.10.0 · vivo\", found \"0.9.0 · vivo\""}},
		{name: "R3 en status chip word differs", changelog: rmChangelog,
			en: "statusword", want: []string{rmEN + ":2: expected status chip \"0.10.0 · live\", found \"0.10.0 · vivo\""}},
		{name: "R4 pt Fontes range differs", changelog: rmChangelog,
			pt: "range", want: []string{rmPT + ":", "expected Sources range end \"0.10.0\", found \"0.9.0\""}},
		{name: "R4 en Fontes range differs", changelog: rmChangelog,
			en: "range", want: []string{rmEN + ":", "expected Sources range end \"0.10.0\", found \"0.9.0\""}},
		{name: "R5 missing card cites the entregue h2 line", changelog: rmChangelog,
			pt: "missing", want: []string{rmPT + ":5: expected a Delivered card for 0.8.0, found none"}},
		{name: "R5 duplicate card", changelog: rmChangelog,
			pt: "dup", want: []string{rmPT + ":", "expected exactly one card for 0.9.0, found a duplicate"}},
		{name: "R5 extra card without CHANGELOG heading", changelog: rmChangelog,
			en: "extra", want: []string{rmEN + ":", "expected no card for 0.11.0 (no CHANGELOG.md heading), found one"}},
		{name: "R6 semantic order 0.9.0 above 0.10.0", changelog: rmChangelog,
			pt: "swap", want: []string{rmPT + ":", "found card 0.10.0 after card 0.9.0"}},
		{name: "R6 detail card between version cards passes", changelog: rmChangelog,
			pt: "detail-ok", en: "detail-ok"},
		{name: "R6 detail card at the top and bottom passes", changelog: rmChangelog,
			pt: "detail-edges", en: "detail-edges"},
		{name: "R7 en has one card fewer", changelog: rmChangelog,
			en: "missing", want: []string{rmEN + ":5: expected the same version cards as docs/roadmap.html, found pt [0.10.0, 0.9.0, 0.8.0] and en [0.10.0, 0.9.0]"}},
		{name: "R7 en Fontes start differs from pt", changelog: rmChangelog,
			en: "rangestart", want: []string{rmEN + ":", "expected Sources range start \"0.1.0\" as in docs/roadmap.html, found \"0.2.0\""}},
		{name: "R7 same versions in another order", changelog: rmChangelog,
			en: "swap", want: []string{rmEN + ":5: expected the same version cards", "en [0.9.0, 0.10.0, 0.8.0]"}},
		{name: "R8 pt and en both wrong, both cited", changelog: rmChangelog,
			pt: "date", en: "date", want: []string{rmPT + ":1: expected date chip", rmEN + ":1: expected date chip"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ptPage, enPage := ok(rmOpts{})
			for lang, mode := range map[string]string{"pt": tc.pt, "en": tc.en} {
				var o rmOpts
				word := map[string]string{"pt": "vivo", "en": "live"}[lang]
				switch mode {
				case "date":
					o.date = "2026-09-01"
				case "status":
					o.status = "0.9.0 · " + word
				case "statusword":
					o.status = "0.10.0 · vivo"
				case "range":
					o.rangeEnd = "0.9.0"
				case "rangestart":
					// rewritten below on the rendered page
				case "missing":
					o = cards(cardReal010, "0.9.0 · B")
				case "dup":
					o = cards(cardReal010, "0.9.0 · B", "0.9.0 · B again", "0.8.0 · C")
				case "extra":
					o = cards("0.11.0 · X", cardReal010, "0.9.0 · B", "0.8.0 · C")
				case "swap":
					o = cards("0.9.0 · B", cardReal010, "0.8.0 · C")
				case "detail-ok":
					o = cards(cardReal010, cardDetail, "0.9.0 · B", "0.8.0 · C")
				case "detail-edges":
					o = cards(cardDetail, cardReal010, "0.9.0 · B", cardDetail, "0.8.0 · C", cardDetail)
				case "", "-":
					continue
				default:
					t.Fatalf("unknown mode %q", mode)
				}
				page := rmPage(lang, o)
				if mode == "rangestart" {
					page = strings.Replace(page, "releases 0.1.0", "releases 0.2.0", 1)
				}
				if lang == "pt" {
					ptPage = page
				} else {
					enPage = page
				}
			}
			if tc.pt == "-" {
				ptPage = ""
			}
			if tc.en == "-" {
				enPage = ""
			}
			cl := tc.changelog
			if cl == "-" {
				cl = ""
			}
			got := rmRun(rmRoot(t, cl, ptPage, enPage))
			if len(tc.want) == 0 {
				if got != "" {
					t.Fatalf("expected no error, got:\n%s", got)
				}
				return
			}
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("missing %q in:\n%s", w, got)
				}
			}
		})
	}
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
			t.Errorf("semverLess(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

// R5/R6: the real 0.10.0 <h3> holds <code>; the parser must still see a card
// and must skip the unversioned detail card that follows it.
func TestParseRoadmapRealMarkup(t *testing.T) {
	page := rmPage("pt", rmOpts{cards: []string{cardReal010, cardDetail, "0.9.0 · B"}})
	d := parseRoadmap(rmPT, page)
	if got := strings.Join(roadmapVersions(d), ","); got != "0.10.0,0.9.0" {
		t.Fatalf("cards = %q, want 0.10.0,0.9.0", got)
	}
	if d.entregueLine != 5 {
		t.Errorf("entregueLine = %d, want 5", d.entregueLine)
	}
}

// R8: every failure is "<file>:<line>: expected ..., found ...". A wrong pt date
// is reported once, against the changelog; the pt/en pair is not re-compared.
func TestCheckRoadmapsErrorFormat(t *testing.T) {
	var errs []string
	root := rmRoot(t, rmChangelog, rmPage("pt", rmOpts{date: "2026-01-01"}), rmPage("en", rmOpts{}))
	checkRoadmaps(root, &errs)
	if len(errs) != 1 || !strings.HasPrefix(errs[0], rmPT+":1: expected ") {
		t.Fatalf("errors = %q", errs)
	}
}

// rmErrs returns the error list for a fixture where pt/en are edited by fn.
func rmErrs(t *testing.T, changelog string, fn func(lang, page string) string) []string {
	t.Helper()
	pt, en := rmPage("pt", rmOpts{}), rmPage("en", rmOpts{})
	if fn != nil {
		pt, en = fn("pt", pt), fn("en", en)
	}
	var errs []string
	checkRoadmaps(rmRoot(t, changelog, pt, en), &errs)
	return errs
}

func wantErrs(t *testing.T, errs []string, prefixes ...string) {
	t.Helper()
	if len(errs) != len(prefixes) {
		t.Fatalf("got %d errors, want %d:\n%s", len(errs), len(prefixes), strings.Join(errs, "\n"))
	}
	for i, w := range prefixes {
		if !strings.Contains(errs[i], w) {
			t.Errorf("error %d = %q, want it to contain %q", i, errs[i], w)
		}
	}
}

// M1: a version-shaped <h3> in a non-strict form must be reported, not ignored.
func TestCheckRoadmapsVariantCards(t *testing.T) {
	variants := map[string]struct{ h3, found string }{
		"attribute":  {`<h3 class="x">0.9.0 · B</h3>`, `"0.9.0 · B"`},
		"wrapped":    {"<h3>\n0.9.0 · B</h3>", `" 0.9.0 · B"`},
		"wrapped2":   {"<h3>0.9.0\n· B</h3>", `"0.9.0 · B"`},
		"space":      {`<h3> 0.9.0 · B</h3>`, `" 0.9.0 · B"`},
		"nbsp":       {"<h3>\u00a00.9.0 · B</h3>", `"\u00a00.9.0 · B"`},
		"htmlnbsp":   {`<h3>&nbsp;0.9.0 · B</h3>`, `"&nbsp;0.9.0 · B"`},
		"emdash":     {`<h3>0.9.0 — B</h3>`, `"0.9.0 — B"`},
		"code first": {`<h3><code>0.9.0</code> · B</h3>`, `"0.9.0 · B"`},
	}
	for name, v := range variants {
		t.Run(name, func(t *testing.T) {
			// A variant duplicate of an existing 0.9.0 card, in en only.
			errs := rmErrs(t, rmChangelog, func(lang, page string) string {
				if lang == "pt" {
					return page
				}
				return strings.Replace(page, "<h3>0.8.0 · C</h3>", v.h3+"\n<h3>0.8.0 · C</h3>", 1)
			})
			wantErrs(t, errs, rmEN+":10: expected \"<h3>N.N.N · \", found "+v.found)
		})
	}
	// Detail cards without a version stay ignored.
	wantErrs(t, rmErrs(t, rmChangelog, func(_, page string) string {
		return strings.Replace(page, "<h3>0.8.0 · C</h3>", "<h3>Detail card</h3>\n<h3>0.8.0 · C</h3>", 1)
	}))
}

func TestCheckRoadmapsMissingPieces(t *testing.T) {
	t.Run("no entregue h2", func(t *testing.T) {
		errs := rmErrs(t, rmChangelog, func(lang, page string) string {
			if lang == "en" {
				return strings.Replace(page, `<h2 id="entregue">`, `<h2 id="other">`, 1)
			}
			return page
		})
		wantErrs(t, errs, rmEN+`:1: expected <h2 id="entregue">, found none`)
	})
	t.Run("chip found none", func(t *testing.T) {
		errs := rmErrs(t, rmChangelog, func(lang, page string) string {
			if lang == "pt" {
				return strings.Replace(page, `<span class="chip">status`, `<span class="chp">status`, 1)
			}
			return page
		})
		wantErrs(t, errs, rmPT+`:1: expected status chip "0.10.0 · vivo", found none`)
	})
	t.Run("missing path chip on one side is not blamed on en", func(t *testing.T) {
		errs := rmErrs(t, rmChangelog, func(lang, page string) string {
			if lang == "pt" {
				return strings.Replace(page, `<span class="chip">path`, `<span class="chp">path`, 1)
			}
			return page
		})
		wantErrs(t, errs, rmPT+":1: expected path chip, found none")
	})
	t.Run("missing path chip on both sides", func(t *testing.T) {
		errs := rmErrs(t, rmChangelog, func(_, page string) string {
			return strings.Replace(page, `<span class="chip">path`, `<span class="chp">path`, 1)
		})
		wantErrs(t, errs, rmPT+":1: expected path chip", rmEN+":1: expected path chip")
	})
	t.Run("card after the closing h2 is not counted", func(t *testing.T) {
		errs := rmErrs(t, rmChangelog, func(lang, page string) string {
			if lang == "en" {
				return strings.Replace(page, "<h3>9.9.9 · not delivered</h3>", "<h3>0.9.0 · again</h3>", 1)
			}
			return page
		})
		wantErrs(t, errs)
	})
	t.Run("CRLF files", func(t *testing.T) {
		crlf := func(_, page string) string { return strings.ReplaceAll(page, "\n", "\r\n") }
		wantErrs(t, rmErrs(t, strings.ReplaceAll(rmChangelog, "\n", "\r\n"), crlf))
	})
}

func TestCheckRoadmapsLanguageWords(t *testing.T) {
	errs := rmErrs(t, rmChangelog, func(lang, page string) string {
		if lang == "en" {
			return strings.Replace(strings.Replace(page, "chip\">date", "chip\">data", 1), " to 0.10.0", " a 0.10.0", 1)
		}
		return page
	})
	wantErrs(t, errs, rmEN+":1: expected date chip word \"date\", found \"data\"", "expected Sources range word \"to\", found \"a\"")
}

func TestCheckRoadmapsChangelogHeadings(t *testing.T) {
	t.Run("duplicate heading keeps the first date", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## 0.10.0 — 2020-01-01\n", nil),
			"CHANGELOG.md:13: expected one heading per release, found 0.10.0 again")
	})
	t.Run("bracketed and v-prefixed headings", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## [0.13.0] — 2026-09-26\n", nil), "CHANGELOG.md:13: expected a heading like")
		wantErrs(t, rmErrs(t, rmChangelog+"\n## v0.13.0 — 2026-09-26\n", nil), "CHANGELOG.md:13: expected a heading like")
	})
	t.Run("Unreleased heading is allowed", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## [Unreleased]\n", nil))
	})
	t.Run("near-miss heading", func(t *testing.T) {
		wantErrs(t, rmErrs(t, rmChangelog+"\n## 0.7.0 - 2026-09-20\n", nil),
			"CHANGELOG.md:13: expected a heading like")
	})
}

// The real repository must satisfy its own check.
func TestCheckRoadmapsRealRepo(t *testing.T) {
	var errs []string
	checkRoadmaps("../..", &errs)
	if len(errs) != 0 {
		t.Fatalf("real repo roadmap errors:\n%s", strings.Join(errs, "\n"))
	}
}

func TestCheckRoadmapsPathChipDiffers(t *testing.T) {
	errs := rmErrs(t, rmChangelog, func(lang, page string) string {
		if lang == "en" {
			return strings.Replace(page, "<b>docs/roadmap.html</b>", "<b>docs/other.html</b>", 1)
		}
		return page
	})
	wantErrs(t, errs, rmEN+`:3: expected path chip "docs/roadmap.html" as in docs/roadmap.html, found "docs/other.html"`)
}
