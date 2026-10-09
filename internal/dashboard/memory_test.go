package dashboard

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/harness"
)

func TestMemoryRoute(t *testing.T) {
	home, repoA, repoG := t.TempDir(), t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repoA: "alpha", repoG: "glob"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(repoA, ".harness/MEMORY.md"), "repo memory\n")
	mk(t, filepath.Join(repoA, ".harness/RISKS.md"), "repo risks\n")
	mk(t, filepath.Join(repoA, ".harness/EPOCHAL.md"), "secret epochal\n")
	mk(t, filepath.Join(home, "projects/glob/MEMORY.md"), "global memory\n")
	mk(t, filepath.Join(home, "projects/glob/big/x"), "")
	big := filepath.Join(home, "projects/glob/RISKS.md")
	mk(t, big, strings.Repeat("x", 1<<20+1))
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})

	ok := map[string]string{
		"/api/memory?project=alpha&file=memory": "repo memory\n",
		"/api/memory?project=alpha&file=risks":  "repo risks\n",
		"/api/memory?project=glob&file=memory":  "global memory\n",
	}
	for u, want := range ok {
		w := do(h, "GET", "127.0.0.1:1", u)
		if w.Code != 200 || w.Body.String() != want || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" ||
			w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatalf("%s: %d %q %v", u, w.Code, w.Body.String(), w.Header())
		}
	}
	for _, u := range []string{
		"/api/memory?project=alpha&file=epochal", "/api/memory?project=alpha&file=MEMORY.md",
		"/api/memory?project=..%2Falpha&file=memory", "/api/memory?project=a%5Calpha&file=memory",
		"/api/memory?project=..&file=memory", "/api/memory?project=nope&file=memory",
		"/api/memory?project=glob&file=risks", // oversized
		"/api/memory?project=alpha", "/api/memory",
	} {
		w := do(h, "GET", "127.0.0.1:1", u)
		if w.Code != 404 || strings.TrimSpace(w.Body.String()) != "404 page not found" {
			t.Fatalf("%s: %d %q", u, w.Code, w.Body.String())
		}
	}
	// symlink and absent
	os.Remove(filepath.Join(repoA, ".harness/RISKS.md"))
	if w := do(h, "GET", "127.0.0.1:1", "/api/memory?project=alpha&file=risks"); w.Code != 404 {
		t.Fatalf("absent: %d", w.Code)
	}
	if err := os.Symlink(filepath.Join(repoA, ".harness/EPOCHAL.md"), filepath.Join(repoA, ".harness/RISKS.md")); err == nil {
		if w := do(h, "GET", "127.0.0.1:1", "/api/memory?project=alpha&file=risks"); w.Code != 404 || strings.Contains(w.Body.String(), "secret") {
			t.Fatalf("symlink: %d %q", w.Code, w.Body.String())
		}
	}
	// method, host, Origin, state carries no text
	for _, p := range []string{"/api/memory?project=alpha&file=memory", "/metrics.js"} {
		for _, m := range []string{"POST", "OPTIONS", "HEAD"} {
			if w := do(h, m, "127.0.0.1:1", p); w.Code != 405 || w.Header().Get("Allow") != "GET" {
				t.Fatalf("%s %s: %d", m, p, w.Code)
			}
		}
		if w := do(h, "GET", "example.com", p); w.Code != 403 {
			t.Fatalf("host %s: %d", p, w.Code)
		}
	}
	r := httptest.NewRequest("GET", "http://127.0.0.1:1/api/memory?project=alpha&file=memory", nil)
	r.Header.Set("Origin", "http://127.0.0.1:1")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("origin: %d", w.Code)
	}
	if w := do(h, "GET", "127.0.0.1:1", "/api/state"); strings.Contains(w.Body.String(), "repo memory") {
		t.Fatal("state leaks memory text")
	}
}

func TestMetricsJSRoute(t *testing.T) {
	w := do(testHandler(t, ""), "GET", "127.0.0.1:1", "/metrics.js")
	if w.Code != 200 || w.Header().Get("Content-Type") != "text/javascript; charset=utf-8" ||
		w.Header().Get("Cache-Control") != "no-cache" || w.Header().Get("X-Content-Type-Options") != "nosniff" || !strings.Contains(w.Body.String(), "use strict") {
		t.Fatalf("%d %v %q", w.Code, w.Header(), w.Body.String())
	}
}

func TestWarningsUseRegisteredNameAndEpochalReadError(t *testing.T) {
	home, repo := t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repo: "custom"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(repo, ".harness/tasks/T-9/TASK.md"), tk("2026-10-05 / x", "2026-10-01", "", "en"))
	st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	if !strings.Contains(strings.Join(st.Warnings, "\n"), "custom: ticket T-9") {
		t.Fatalf("warnings: %v", st.Warnings)
	}
	// a non-ErrNotExist open error (parent is a file) is "unreadable", not "not a regular file"
	f := filepath.Join(t.TempDir(), "file")
	mk(t, f, "x")
	var m Metrics
	w := strings.Join(fileMetrics(&m, "custom", f), "\n")
	if m.Epochal.Batches != nil || !strings.Contains(w, "custom: EPOCHAL.md unreadable.") || strings.Contains(w, "EPOCHAL.md unreadable (not") {
		t.Fatalf("epochal: %+v %q", m.Epochal, w)
	}
}

func TestSymlinkedHarnessDirRefused(t *testing.T) {
	home, repo, real := t.TempDir(), t.TempDir(), t.TempDir()
	if err := harness.SaveConfig(home, harness.Config{Projects: map[string]string{repo: "lnk"}}); err != nil {
		t.Fatal(err)
	}
	mk(t, filepath.Join(real, "MEMORY.md"), "outside memory\n")
	mk(t, filepath.Join(real, "EPOCHAL.md"), "## Archived batches\n### a\n")
	if err := os.Symlink(real, filepath.Join(repo, ".harness")); err != nil {
		t.Skip("no symlinks")
	}
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	if w := do(h, "GET", "127.0.0.1:1", "/api/memory?project=lnk&file=memory"); w.Code != 404 || strings.Contains(w.Body.String(), "outside") {
		t.Fatalf("memory via symlinked dir: %d %q", w.Code, w.Body.String())
	}
	st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	for _, p := range st.Projects {
		if p.Name == "lnk" && (p.Metrics.Memory.Lines != nil || p.Metrics.Epochal.Batches != nil || p.Metrics.Incidents != nil) {
			t.Fatalf("metrics read through symlink: %+v", p.Metrics)
		}
	}
	if !strings.Contains(strings.Join(st.Warnings, "\n"), "lnk: harness folder is a symlink") {
		t.Fatalf("warnings: %v", st.Warnings)
	}
}

func TestEpochalScanCached(t *testing.T) {
	h := t.TempDir()
	ep := filepath.Join(h, "EPOCHAL.md")
	mk(t, ep, "## Archived batches\n### a\n- 2026-01-01\n")
	before := epochalScans
	for i := 0; i < 3; i++ {
		e, err := scanEpochal(ep)
		if err != nil || *e.Batches != 1 {
			t.Fatalf("%v %+v", err, e)
		}
	}
	if epochalScans-before != 1 {
		t.Fatalf("scans: %d", epochalScans-before)
	}
	mk(t, ep, "## Archived batches\n### a\n### b\n")
	if e, _ := scanEpochal(ep); *e.Batches != 2 {
		t.Fatalf("changed file not rescanned: %+v", e)
	}
}
