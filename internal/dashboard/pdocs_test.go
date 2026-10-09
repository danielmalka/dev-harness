package dashboard

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/harness"
)

func TestPdocsRouteAndDocsFlag(t *testing.T) {
	home, repoA, repoG := t.TempDir(), t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repoA: "alpha", repoG: "glob"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(repoA, ".harness/MEMORY.md"), "m\n")
	mk(t, filepath.Join(repoG, ".harness/MEMORY.md"), "m\n")
	pd := filepath.Join(home, "projects/alpha/pdocs")
	mk(t, filepath.Join(pd, "index.html"), "<html>hi</html>")
	mk(t, filepath.Join(pd, "d.svg"), "<svg/>")
	mk(t, filepath.Join(pd, "p.png"), "png")
	mk(t, filepath.Join(pd, "x.js"), "alert(1)")
	mk(t, filepath.Join(pd, "UP.HTML"), "x")
	mk(t, filepath.Join(pd, "big.html"), strings.Repeat("x", pdocsCap+1))
	mk(t, filepath.Join(pd, "edge.html"), strings.Repeat("x", pdocsCap))
	mk(t, filepath.Join(home, "secret.html"), "secret")
	symOK := os.Symlink(filepath.Join(home, "secret.html"), filepath.Join(pd, "ln.html")) == nil
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})

	ok := map[string][2]string{
		"/pdocs/alpha/index.html": {"<html>hi</html>", "text/html; charset=utf-8"},
		"/pdocs/alpha/":           {"<html>hi</html>", "text/html; charset=utf-8"},
		"/pdocs/alpha/d.svg":      {"<svg/>", "image/svg+xml"},
		"/pdocs/alpha/p.png":      {"png", "image/png"},
	}
	for u, want := range ok {
		w := do(h, "GET", "127.0.0.1:1", u)
		if w.Code != 200 || w.Body.String() != want[0] || w.Header().Get("Content-Type") != want[1] ||
			w.Header().Get("Cache-Control") != "no-cache" || w.Header().Get("Content-Security-Policy") != pdocsCSP ||
			w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatalf("%s: %d %v", u, w.Code, w.Header())
		}
	}
	if w := do(h, "GET", "127.0.0.1:1", "/pdocs/alpha/edge.html"); w.Code != 200 {
		t.Errorf("exactly 8 MiB: %d", w.Code)
	}
	bad := []string{"/pdocs/alpha", "/pdocs/alpha/../x.html", "/pdocs/alpha/%2e%2e/x.html", "/pdocs/alpha/a%2Fb.html",
		"/pdocs/alpha/a%5Cb.html", "/pdocs/alpha/index.html:x.html", "/pdocs/alpha/a:b.html", "/pdocs/alpha/x.js", "/pdocs/alpha/UP.HTML", "/pdocs/nope/index.html",
		"/pdocs/alpha/missing.html", "/pdocs/alpha/big.html", "/pdocs/glob/", "/pdocs/", "/pdocs//index.html"}
	if symOK {
		bad = append(bad, "/pdocs/alpha/ln.html")
	}
	for _, u := range bad {
		w := do(h, "GET", "127.0.0.1:1", u)
		if w.Code != 404 || strings.TrimSpace(w.Body.String()) != "404 page not found" || w.Header().Get("Content-Security-Policy") != pdocsCSP {
			t.Fatalf("%s: %d %q csp=%q", u, w.Code, w.Body.String(), w.Header().Get("Content-Security-Policy"))
		}
	}
	// other routes keep the common CSP; guard behavior intact
	common := "default-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'"
	for _, u := range []string{"/", "/api/state", "/api/memory", "/metrics.js"} {
		if got := do(h, "GET", "127.0.0.1:1", u).Header().Get("Content-Security-Policy"); got != common {
			t.Errorf("%s csp %q", u, got)
		}
	}
	if w := do(h, "POST", "127.0.0.1:1", "/pdocs/alpha/index.html"); w.Code != 405 || w.Header().Get("Content-Security-Policy") != common {
		t.Errorf("POST %d", w.Code)
	}
	if w := do(h, "HEAD", "127.0.0.1:1", "/pdocs/alpha/index.html"); w.Code != 405 {
		t.Errorf("HEAD %d", w.Code)
	}
	if w := do(h, "GET", "exemplo.com", "/pdocs/alpha/index.html"); w.Code != 403 {
		t.Errorf("host %d", w.Code)
	}
	// route never writes warnings
	if st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale}); len(st.Warnings) != 0 {
		t.Errorf("warnings: %v", st.Warnings)
	}
}

func TestPdocsStateDocsField(t *testing.T) {
	home, repoA, repoG := t.TempDir(), t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repoA: "alpha", repoG: "glob"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(repoA, ".harness/MEMORY.md"), "m\n")
	mk(t, filepath.Join(repoG, ".harness/MEMORY.md"), "m\n")
	cfg := Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale}
	docs := func() map[string]bool {
		st, raw := getState(t, cfg)
		if !strings.Contains(raw, `"docs":`) {
			t.Fatal("docs missing")
		}
		m := map[string]bool{}
		for _, p := range st.Projects {
			m[p.Name] = p.Docs
		}
		if strings.Contains(raw, home+"/projects/alpha/pdocs") {
			t.Errorf("absolute path leaked")
		}
		return m
	}
	if m := docs(); m["alpha"] || m["glob"] {
		t.Errorf("no pdocs: %v", m)
	}
	idx := filepath.Join(home, "projects/alpha/pdocs/index.html")
	mk(t, idx, "<html/>")
	if m := docs(); !m["alpha"] || m["glob"] {
		t.Errorf("with index: %v", m)
	}
	mk(t, idx, strings.Repeat("x", pdocsCap+1))
	if m := docs(); m["alpha"] {
		t.Error("oversized should be false")
	}
	if st, _ := getState(t, cfg); !strings.Contains(strings.Join(st.Warnings, "|"), "alpha: pdocs/index.html is over 8 MiB") {
		t.Errorf("no warning: %v", st.Warnings)
	}
	os.Remove(idx)
	if os.Symlink(filepath.Join(home, "config.yaml"), idx) == nil {
		if m := docs(); m["alpha"] {
			t.Error("symlink should be false")
		}
		if st, _ := getState(t, cfg); !strings.Contains(strings.Join(st.Warnings, "|"), "alpha: pdocs/index.html is a symlink") {
			t.Errorf("no symlink warning: %v", st.Warnings)
		}
	}
}

func TestPdocsHomeSymlinkAccepted(t *testing.T) {
	real, repo := t.TempDir(), t.TempDir()
	link := filepath.Join(t.TempDir(), "h")
	if os.Symlink(real, link) != nil {
		t.Skip("symlinks unavailable")
	}
	harness.SaveConfig(link, harness.Config{Projects: map[string]string{repo: "alpha"}})
	mk(t, filepath.Join(repo, ".harness/MEMORY.md"), "m\n")
	mk(t, filepath.Join(real, "projects/alpha/pdocs/index.html"), "ok")
	h := Handler(Config{Home: link, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	if w := do(h, "GET", "127.0.0.1:1", "/pdocs/alpha/"); w.Code != 200 || w.Body.String() != "ok" {
		t.Errorf("%d", w.Code)
	}
}

func TestAPIRefusesCrossOriginBrowsers(t *testing.T) {
	h := testHandler(t, "")
	get := func(path string, hdr map[string]string) int {
		r := httptest.NewRequest("GET", "http://127.0.0.1:1"+path, nil)
		r.Host = "127.0.0.1:1"
		for k, v := range hdr {
			r.Header.Set(k, v)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Code
	}
	for _, p := range []string{"/api/state", "/api/memory?project=x&file=memory"} {
		for _, hdr := range []map[string]string{{"Origin": "null"}, {"Sec-Fetch-Site": "cross-site"}} {
			if c := get(p, hdr); c != 403 {
				t.Errorf("%s %v: %d", p, hdr, c)
			}
		}
		if c := get(p, map[string]string{"Origin": "http://127.0.0.1:1", "Sec-Fetch-Site": "same-origin"}); c == 403 {
			t.Errorf("%s same-origin refused", p)
		}
	}
	if c := get("/api/state", nil); c != 200 {
		t.Errorf("no origin: %d", c)
	}
}
