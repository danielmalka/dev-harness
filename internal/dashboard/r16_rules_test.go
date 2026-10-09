package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// PRD-016 rule tests (QA, T-1702-05). Existing builder tests are chained; the cases below fill the gaps.

func readF(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestR1_TemplatesAndStatusVocabulary(t *testing.T) {
	TestStatusCasesFixture(t)
	TestStatusCasesTSCopyMatches(t)
	for lang, p := range map[string][2]string{"pt-br": {"Criado / atualizado", "Concluído em"}, "en": {"Created / updated", "Done on"}} {
		md := readF(t, "../../templates/"+lang+"/TASK.md")
		lines := strings.Split(md, "\n")
		found := false
		for i, l := range lines {
			if strings.HasPrefix(l, "| "+p[0]+" |") {
				found = strings.HasPrefix(lines[i+1], "| "+p[1]+" |")
			}
		}
		if !found {
			t.Errorf("%s: %q row not right after %q", lang, p[1], p[0])
		}
		if got := TaskStatus(md); got == "done" {
			t.Errorf("%s: template reads as done", lang)
		}
	}
}

// R3: old project fields keep name and value; metrics is added.
func TestR3_AdditiveStateContract(t *testing.T) {
	home, repo := t.TempDir(), t.TempDir()
	register(t, home, repo)
	mk(t, filepath.Join(repo, ".harness/tasks/T-1/TASK.md"), tk("2026-10-01 / x", "2026-10-02", "", "en"))
	mk(t, filepath.Join(repo, "docs/prd/PRD-001-a.md"), "| Status | delivered on 2026-10-02 |\n")
	mk(t, filepath.Join(repo, "docs/prd/PRD-002-b.md"), "| Status | approved |\n")
	_, raw := getState(t, Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	var s struct {
		Projects []map[string]json.RawMessage `json:"projects"`
	}
	if err := json.Unmarshal([]byte(raw), &s); err != nil || len(s.Projects) != 1 {
		t.Fatalf("%v %s", err, raw)
	}
	keys := []string{}
	for k := range s.Projects[0] {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if want := []string{"delivered", "docs", "harness", "metrics", "mode", "name", "open", "path"}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("keys %v want %v", keys, want)
	}
	if string(s.Projects[0]["delivered"]) != `["PRD-001"]` {
		t.Fatalf("delivered: %s", s.Projects[0]["delivered"])
	}
	var m Metrics
	if err := json.Unmarshal(s.Projects[0]["metrics"], &m); err != nil || len(m.PRDs) != 2 || m.PRDs[0].Status != "delivered" {
		t.Fatalf("metrics: %v %s", err, s.Projects[0]["metrics"])
	}
}

func TestR3b_IncidentsRealTemplates(t *testing.T) {
	for _, lang := range []string{"pt-br", "en"} {
		if n := countIncidents(readF(t, "../../templates/"+lang+"/RISKS.md")); n != 0 {
			t.Errorf("%s template: %d", lang, n)
		}
	}
	TestIncidentsMemoryEpochal(t)
}

func TestR3c_EpochalRealTemplates(t *testing.T) {
	for _, lang := range []string{"pt-br", "en"} {
		h := t.TempDir()
		mk(t, filepath.Join(h, "EPOCHAL.md"), readF(t, "../../templates/"+lang+"/EPOCHAL.md"))
		e, err := scanEpochal(filepath.Join(h, "EPOCHAL.md"))
		if err != nil || *e.Batches != 0 || e.LastAt != nil {
			t.Errorf("%s: %+v %v", lang, e, err)
		}
	}
	// en markers, over 1 MiB, three batches in the same file
	h := t.TempDir()
	big := "## Archived batches\n### a\n- 2026-01-01T00:00:00Z\n### b\n<!-- BEGIN-RAW x -->\n## fake\n### fake\n" +
		strings.Repeat("y", 1<<20+1) + "\n<!-- END-RAW x -->\n### c\n- 2026-03-03T03:03:03-03:00\n"
	mk(t, filepath.Join(h, "EPOCHAL.md"), big)
	e, err := scanEpochal(filepath.Join(h, "EPOCHAL.md"))
	if err != nil || *e.Batches != 3 || *e.LastAt != "2026-03-03T03:03:03-03:00" {
		t.Fatalf("%+v %v", e, err)
	}
	TestEpochal(t)
}

func TestR4_UndatedAndNoDatedNull(t *testing.T) {
	repo, h := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(h, "tasks/T-1/TASK.md"), tk("2026-10-01 / x", "", "", "pt"))
	mk(t, filepath.Join(h, "tasks/T-2/TASK.md"), tk("2026-10-01 / x", "", "in progress", "en"))
	mk(t, filepath.Join(repo, "docs/prd/PRD-001-a.md"), "| Status | approved |\n")
	m, _ := metricsOf(t, repo, h)
	p := m.PRDs[0]
	if p.TicketDaysAvg != nil || p.TicketDaysN != 0 || p.TicketUndated != 1 || p.Tickets != 2 || p.Done != 1 {
		t.Fatalf("%+v", p)
	}
	TestTicketAndPRDDates(t)
}

func TestR5_SameDayIsZero(t *testing.T) {
	repo, h := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(h, "tasks/T-1/TASK.md"), tk("2026-10-01 / x", "2026-10-01", "", "pt"))
	mk(t, filepath.Join(h, "tasks/T-2/TASK.md"), tk("2026-10-01 / x", "2026-10-01", "", "en"))
	mk(t, filepath.Join(repo, "docs/prd/PRD-001-a.md"), "| Criado / atualizado | 2026-10-01 / z |\n| Status | entregue em 2026-10-01 |\n")
	m, _ := metricsOf(t, repo, h)
	p := m.PRDs[0]
	if p.TicketDaysN != 2 || p.TicketUndated != 0 || *p.TicketDaysAvg != 0 || *p.Days != 0 || m.PRDDaysN != 1 || *m.PRDDaysAvg != 0 {
		t.Fatalf("%+v avg=%v", p, m.PRDDaysAvg)
	}
}

func TestR6_MissingAndGlobalNoWarning(t *testing.T) {
	m, w := metricsOf(t, t.TempDir(), t.TempDir())
	if !m.Memory.Missing || m.Incidents != nil || *m.Epochal.Batches != 0 || m.Epochal.LastAt != nil || len(w) != 0 {
		t.Fatalf("%+v %v", m, w)
	}
	TestIncidentsMemoryEpochal(t)
	TestMetricsInStateGlobalMode(t)
}

func TestR8_RouteAndFileSizes(t *testing.T) {
	TestMetricsJSRoute(t)
	for _, f := range []string{"web/index.html", "web/metrics.js"} {
		if n := strings.Count(readF(t, f), "\n"); n >= 500 {
			t.Errorf("%s: %d lines", f, n)
		}
	}
}

func TestR10_RouteBothModesAndNoEcho(t *testing.T) {
	TestMemoryRoute(t)
	home, repo := t.TempDir(), t.TempDir()
	register(t, home, repo)
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale, DoneDecay: time.Minute})
	g := filepath.Join(home, "projects", filepath.Base(repo))
	mk(t, filepath.Join(g, "RISKS.md"), "global risks\n")
	if w := do(h, "GET", "127.0.0.1:1", "/api/memory?project="+filepath.Base(repo)+"&file=risks"); w.Code != 200 || w.Body.String() != "global risks\n" {
		t.Fatalf("global risks: %d %q", w.Code, w.Body.String())
	}
	for _, u := range []string{
		"/api/memory?project=zzqq-unreg&file=epochal", "/api/memory?project=a%5Cb&file=zzqq-file",
		"/api/memory?project=a%2Fb&file=memory", "/api/memory?project=" + filepath.Base(repo) + "&file=memory", // memory absent
	} {
		w := do(h, "GET", "127.0.0.1:1", u)
		b := w.Body.String()
		if w.Code != 404 || strings.Contains(b, "zzqq") || strings.Contains(b, "/") && !strings.Contains(b, "404 page not found") || strings.Contains(b, home) || strings.Contains(b, repo) {
			t.Errorf("%s: %d %q", u, w.Code, b)
		}
	}
}

func TestR12_GuardOnNewRoutes(t *testing.T) {
	TestMemoryRoute(t) // POST/OPTIONS/HEAD 405, Host 403, loopback Origin 200
}

func TestR14_NoNewDependencies(t *testing.T) {
	for _, f := range []string{"../../go.mod"} {
		if strings.Contains(strings.ToLower(readF(t, f)), "sqlite") {
			t.Errorf("%s mentions sqlite", f)
		}
	}
}

func TestR2_DateReading(t *testing.T) {
	TestTicketAndPRDDates(t) // valid, placeholders (AAAA/YYYY), absent, inverted + warning w/o paths, en "delivered on", pt/en
	TestMetricsInStateGlobalMode(t)
	TestWarningsUseRegisteredNameAndEpochalReadError(t)
}
