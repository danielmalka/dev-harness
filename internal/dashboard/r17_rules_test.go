package dashboard

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// PRD-017 rule tests (QA, T-1703-05). Builder tests are chained; the cases below fill the gaps.
// Names carry PRD017 where they could collide with r16_rules_test.go.

const skill = "../../.skills/doc-template-html/"

var pdocsIDs = []string{"motivacao", "como-rodar", "planejado-desenvolvido", "changelog", "diagrama-macro", "dominios", "referencias"}
var mermaidTag = regexp.MustCompile(`mermaid@[0-9]+\.[0-9]+\.[0-9]+/`)

func r17Home(t *testing.T) (Config, string) {
	home, repoA, repoG := t.TempDir(), t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repoA: "alpha", repoG: "glob"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(repoA, ".harness/MEMORY.md"), "m\n")
	mk(t, filepath.Join(repoG, ".harness/MEMORY.md"), "m\n")
	mk(t, filepath.Join(home, "projects/alpha/pdocs/index.html"), "<html>hi</html>")
	mk(t, filepath.Join(home, "projects/alpha/pdocs/d.svg"), "<svg/>")
	return Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale}, home
}

func TestR9_PdocsRouteServesAndRefuses(t *testing.T) {
	TestPdocsRouteAndDocsFlag(t)
	TestPdocsHomeSymlinkAccepted(t)
	// 404 bodies never echo name, file or reason; the route never touches warnings.
	cfg, _ := r17Home(t)
	h := Handler(cfg)
	before, _ := getState(t, cfg)
	for _, u := range []string{"/pdocs/nope/index.html", "/pdocs/alpha/missing.html", "/pdocs/alpha/x.css", "/pdocs/alpha/x.HTML"} {
		b := do(h, "GET", "127.0.0.1:1", u).Body.String()
		for _, s := range []string{"alpha", "nope", "missing", "x."} {
			if strings.Contains(b, s) {
				t.Errorf("%s echoes %q: %q", u, s, b)
			}
		}
	}
	if after, _ := getState(t, cfg); strings.Join(after.Warnings, "|") != strings.Join(before.Warnings, "|") {
		t.Errorf("warnings changed: %v -> %v", before.Warnings, after.Warnings)
	}
}

func TestR9b_PdocsStateDocsAndWarnings(t *testing.T) {
	TestPdocsStateDocsField(t)
}

func TestR10_PdocsCSPOnlyOnRoute(t *testing.T) {
	cfg, _ := r17Home(t)
	h := Handler(cfg)
	want := "default-src 'self' 'unsafe-inline'; img-src 'self' data:; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net/npm/mermaid@11.17.2/dist/mermaid.min.js; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; connect-src 'self'; sandbox allow-scripts allow-popups"
	if pdocsCSP != want {
		t.Fatalf("pdocsCSP drifted: %q", pdocsCSP)
	}
	common := "default-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'"
	for _, u := range []string{"/pdocs/alpha/index.html", "/pdocs/alpha/d.svg", "/pdocs/alpha/x.js"} {
		if got := do(h, "GET", "127.0.0.1:1", u).Header().Get("Content-Security-Policy"); got != want {
			t.Errorf("%s: %q", u, got)
		}
	}
	for _, u := range []string{"/", "/api/state", "/api/memory?project=alpha&file=memory", "/metrics.js", "/sprite/frente"} {
		if got := do(h, "GET", "127.0.0.1:1", u).Header().Get("Content-Security-Policy"); got != common {
			t.Errorf("%s: %q", u, got)
		}
	}
}

// Every external host in non-test Go source is on the allow-list (existing ones are not from the route).
func TestR10_NoNewExternalDomain(t *testing.T) {
	re := regexp.MustCompile(`https?://[a-z0-9.-]+`)
	ok := map[string]bool{"https://cdn.jsdelivr.net": true, "https://fonts.googleapis.com": true, "https://fonts.gstatic.com": true}
	for _, root := range []string{"../../internal", "../../cmd"} {
		filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() || strings.HasSuffix(p, "_test.go") || !strings.HasSuffix(p, ".go") {
				return nil
			}
			for _, m := range re.FindAllString(readF(t, p), -1) {
				if !ok[m] && !strings.HasPrefix(m, "http://127.0.0.1") && !strings.HasPrefix(m, "http://localhost") && !strings.HasPrefix(m, "http://[") {
					t.Errorf("%s: %s", p, m)
				}
			}
			return nil
		})
	}
}

func TestR11_PdocsGETOnlyAndHeaders(t *testing.T) {
	cfg, _ := r17Home(t)
	h := Handler(cfg)
	for _, m := range []string{"POST", "PUT", "HEAD"} {
		w := do(h, m, "127.0.0.1:1", "/pdocs/alpha/index.html")
		if w.Code != 405 || w.Header().Get("Allow") != "GET" {
			t.Errorf("%s: %d allow=%q", m, w.Code, w.Header().Get("Allow"))
		}
	}
	if w := do(h, "GET", "exemplo.com", "/pdocs/alpha/index.html"); w.Code != 403 {
		t.Errorf("host %d", w.Code)
	}
	w := do(h, "GET", "127.0.0.1:1", "/pdocs/alpha/index.html")
	if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("X-Dh-Dashboard") != "1" {
		t.Errorf("headers: %v", w.Header())
	}
}

func TestR12_DocsFieldAdditive(t *testing.T) {
	cfg, _ := r17Home(t)
	st, raw := getState(t, cfg)
	m := map[string]bool{}
	for _, p := range st.Projects {
		m[p.Name] = p.Docs
	}
	if !m["alpha"] || m["glob"] || len(m) != 2 || strings.Count(raw, `"docs":`) != 2 {
		t.Errorf("docs: %v in %s", m, raw)
	}
	// old fields keep their keys (R3 of PRD-016 chained) and CLI struct has no docs
	TestR3_AdditiveStateContract(t)
	if _, ok := reflect.TypeOf(harness.Project{}).FieldByName("Docs"); ok {
		t.Error("harness.Project has Docs")
	}
}

func TestR4_PdocsSkeletonsAndModels(t *testing.T) {
	var tags []string
	for _, lang := range []string{"pt-br", "en"} {
		files := []string{skill + "assets/skeletons/" + lang + "/projeto.html", skill + "assets/projeto-modelo." + lang + ".html"}
		for i, f := range files {
			s := readF(t, f)
			var got []string
			for _, m := range regexp.MustCompile(`<h2[^>]*\sid="([^"]+)"`).FindAllStringSubmatch(s, -1) {
				got = append(got, m[1])
			}
			if strings.Join(got, ",") != strings.Join(pdocsIDs, ",") {
				t.Errorf("%s h2 ids: %v", f, got)
			}
			if i == 1 {
				seg := s[strings.Index(s, `id="diagrama-macro"`):]
				if j := strings.Index(seg[10:], "<h2"); j > 0 {
					seg = seg[:j+10]
				}
				if !strings.Contains(seg, `<pre class="mermaid">`) {
					t.Errorf("%s: no mermaid block in diagrama-macro", f)
				}
				tags = append(tags, regexp.MustCompile(`<script[^>]*mermaid[^>]*>`).FindString(s)+mermaidTag.FindString(s))
			}
		}
	}
	if len(tags) != 2 || tags[0] != tags[1] || !mermaidTag.MatchString(tags[0]) {
		t.Errorf("model script tags differ: %q", tags)
	}
}

func TestR7_MermaidPinnedOnceAndFallback(t *testing.T) {
	stamp := readF(t, skill+"scripts/stamp.sh")
	pin := mermaidTag.FindAllString(stamp, -1)
	if len(pin) != 1 {
		t.Fatalf("stamp.sh pins: %v", pin)
	}
	for _, f := range []string{"projeto-modelo.pt-br.html", "projeto-modelo.en.html"} {
		if !strings.Contains(readF(t, skill+"assets/"+f), pin[0]) {
			t.Errorf("%s lacks %s", f, pin[0])
		}
	}
	if regexp.MustCompile(`mermaid@[0-9]`).MatchString(readF(t, skill+"references/types.md")) {
		t.Error("types.md repeats the version")
	}
	css := readF(t, skill+"assets/shell.css")
	i := strings.Index(css, "pre.mermaid")
	if i < 0 {
		t.Fatal("no pre.mermaid rule")
	}
	rule := css[i:]
	rule = rule[:strings.Index(rule, "}")]
	if strings.Contains(rule, "display:none") || strings.Contains(rule, "display: none") || strings.Contains(rule, "visibility:hidden") || strings.Contains(rule, "visibility: hidden") {
		t.Errorf("pre.mermaid hides the fallback: %s", rule)
	}
}

func TestR14_OneCDNRule(t *testing.T) {
	for _, f := range []string{"../../AGENTS.md", skill + "SKILL.md", skill + "references/types.md"} {
		s := readF(t, f)
		if regexp.MustCompile(`(?i)única exceção|only CDN exception`).MatchString(s) {
			t.Errorf("%s keeps the old single-exception phrase", f)
		}
		if !strings.Contains(s, "Mermaid") || !strings.Contains(s, "Google Fonts") {
			t.Errorf("%s does not name Mermaid and Google Fonts", f)
		}
	}
	t.Log("docs/tutorial.html and docs/en/tutorial.html: not-run (docs step)")
}

func TestR15_OutOfScope(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	b, err := exec.Command("git", "diff", "main", "--", "../../go.mod", "../../go.sum").Output()
	if err != nil {
		t.Skip("git diff main unavailable")
	}
	if len(b) != 0 {
		t.Errorf("go.mod/go.sum changed:\n%s", b)
	}
}
