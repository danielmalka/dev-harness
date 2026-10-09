package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tk(created, done, status string, lang string) string {
	c, d, s := "Criado / atualizado", "Concluído em", "concluída"
	if lang == "en" {
		c, d, s = "Created / updated", "Done on", "done"
	}
	if status != "" {
		s = status
	}
	out := "| " + c + " | " + created + " |\n| Status | " + s + " |\n| PRD (RF-<n>) | PRD-001 |\n"
	if done != "" {
		out += "| " + d + " | " + done + " |\n"
	}
	return out
}

func metricsOf(t *testing.T, repo, h string) (Metrics, []string) {
	t.Helper()
	pr := ProjectProgress(repo, h)
	m := pr.Metrics
	w := append(pr.Warnings, fileMetrics(&m, "proj", h)...)
	return m, w
}

func TestTicketAndPRDDates(t *testing.T) {
	repo, h := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(h, "tasks/T-1/TASK.md"), tk("2026-10-01 / 2026-10-03", "2026-10-03", "", "pt")) // 2 days
	mk(t, filepath.Join(h, "tasks/T-2/TASK.md"), tk("2026-10-01 / 2026-10-01", "2026-10-02", "", "en")) // 1 day
	mk(t, filepath.Join(h, "tasks/T-3/TASK.md"), tk("2026-10-01 / x", "AAAA-MM-DD", "", "pt"))          // placeholder
	mk(t, filepath.Join(h, "tasks/T-4/TASK.md"), tk("2026-10-05 / x", "2026-10-01", "", "en"))          // inverted
	mk(t, filepath.Join(h, "tasks/T-5/TASK.md"), tk("2026-10-01 / x", "", "", "en"))                    // missing
	mk(t, filepath.Join(h, "tasks/T-6/TASK.md"), tk("2026-10-01 / x", "", "in progress", "en"))         // open: nowhere
	mk(t, filepath.Join(h, "tasks/T-7/TASK.md"), tk("2026-10-01 / x", "YYYY-MM-DD", "", "en"))          // placeholder
	mk(t, filepath.Join(repo, "docs/prd/PRD-001-a.md"), "| Criado / atualizado | 2026-09-20 / z |\n| Status | entregue em 2026-09-25 |\n")
	mk(t, filepath.Join(repo, "docs/prd/PRD-002-b.md"), "| Created / updated | 2026-09-25 / z |\n| Status | delivered on 2026-09-25 |\n")
	mk(t, filepath.Join(repo, "docs/prd/PRD-003-c.md"), "| Created / updated | 2026-09-25 / z |\n| Status | delivered on 2026-09-20 |\n")
	mk(t, filepath.Join(repo, "docs/prd/PRD-004-d.md"), "| Created / updated | AAAA-MM-DD / z |\n| Status | approved |\n")
	m, w := metricsOf(t, repo, h)
	p := m.PRDs[0]
	if p.ID != "PRD-001" || p.Status != "delivered" || p.Tickets != 7 || p.Done != 6 || *p.Started != "2026-09-20" || *p.Delivered != "2026-09-25" || *p.Days != 5 {
		t.Fatalf("prd1: %+v", p)
	}
	if p.TicketDaysN != 2 || p.TicketUndated != 4 || *p.TicketDaysAvg != 1.5 {
		t.Fatalf("ticket stats: n=%d undated=%d avg=%v", p.TicketDaysN, p.TicketUndated, p.TicketDaysAvg)
	}
	if m.PRDs[1].Status != "delivered" || *m.PRDs[1].Days != 0 {
		t.Fatalf("same day: %+v", m.PRDs[1])
	}
	if m.PRDs[2].Days != nil || m.PRDs[3].Status != "open" || m.PRDs[3].Started != nil || m.PRDs[3].TicketDaysAvg != nil {
		t.Fatalf("inverted/open: %+v %+v", m.PRDs[2], m.PRDs[3])
	}
	if m.PRDDaysN != 2 || *m.PRDDaysAvg != 2.5 {
		t.Fatalf("prd avg: %v n=%d", m.PRDDaysAvg, m.PRDDaysN)
	}
	joined := strings.Join(w, "\n")
	if !strings.Contains(joined, "T-4") || !strings.Contains(joined, "PRD-003") || strings.Contains(joined, repo) || strings.Contains(joined, h) {
		t.Fatalf("warnings: %v", w)
	}
	pr := ProjectProgress(repo, h)
	if len(pr.Delivered) != 3 || len(pr.Open) != 1 { // delivered on counts as delivered
		t.Fatalf("progress: %+v", pr)
	}
}

func TestIncidentsMemoryEpochal(t *testing.T) {
	repo, h := t.TempDir(), t.TempDir()
	m, w := metricsOf(t, repo, h)
	if m.Incidents != nil || !m.Memory.Missing || m.Memory.Lines != nil || *m.Epochal.Batches != 0 || m.Epochal.LastAt != nil || len(w) != 0 {
		t.Fatalf("absent: %+v %v", m, w)
	}
	for hd, body := range map[string]string{
		"## Incidents\n\nNo incidents.\n\n<!-- - RISK-001 x\n### RISK-002 y -->\n": "0",
		"## Incidentes\n\nNenhum.\n": "0",
		"## Incidents\n\n### RISK-001 a\n- RISK-002 b\n- not an id\n## Other\n- RISK-009 z\n": "2",
	} {
		mk(t, filepath.Join(h, "RISKS.md"), "# RISKS\n\n"+hd)
		m, _ = metricsOf(t, repo, h)
		if body != "0" && *m.Incidents != 2 || body == "0" && *m.Incidents != 0 {
			t.Fatalf("incidents %q: %v", hd, *m.Incidents)
		}
	}
	mk(t, filepath.Join(h, "MEMORY.md"), "a\nb\n")
	m, _ = metricsOf(t, repo, h)
	if *m.Memory.Lines != 2 || *m.Memory.Bytes != 4 || m.Memory.Missing {
		t.Fatalf("memory: %+v", m.Memory)
	}
	mk(t, filepath.Join(h, "MEMORY.md"), "a\nb")
	m, _ = metricsOf(t, repo, h)
	if *m.Memory.Lines != 2 {
		t.Fatalf("no final newline: %d", *m.Memory.Lines)
	}
	mk(t, filepath.Join(h, "MEMORY.md"), "")
	m, _ = metricsOf(t, repo, h)
	if *m.Memory.Lines != 0 {
		t.Fatalf("empty: %d", *m.Memory.Lines)
	}
	mk(t, filepath.Join(h, "MEMORY.md"), strings.Repeat("x", 1<<20+1))
	mk(t, filepath.Join(h, "RISKS.md"), strings.Repeat("x", 1<<20+1))
	m, w = metricsOf(t, repo, h)
	if m.Memory.Lines != nil || m.Memory.Missing || m.Incidents != nil || len(w) != 2 {
		t.Fatalf("oversize: %+v %v", m, w)
	}
	out, _ := json.Marshal(m)
	if !strings.Contains(string(out), `"memory":{"lines":null,"bytes":null,"missing":false}`) {
		t.Fatalf("json: %s", out)
	}
}

const epRaw = "## Lotes arquivados\n\n<!-- template\n-->\n\n### Lote A\n\n- Consolidação: 2026-09-24T22:03:47-03:00\n\n" +
	"<!-- INICIO-BRUTO a -->\n## Regras\n### dentro\n<!-- FIM-BRUTO a -->\n\n### Lote B\n\n- Consolidação: 2026-09-30T14:08:31-03:00\n- later 2026-12-01\n\n" +
	"<!-- INICIO-BRUTO b -->\n## Regras\n### dentro 2026-01-01\n<!-- FIM-BRUTO b -->\n"

func TestEpochal(t *testing.T) {
	repo, h := t.TempDir(), t.TempDir()
	ep := filepath.Join(h, "EPOCHAL.md")
	mk(t, ep, "# E\n\n## Archived batches\n\n<!-- tpl -->\n")
	m, w := metricsOf(t, repo, h)
	if *m.Epochal.Batches != 0 || m.Epochal.LastAt != nil || len(w) != 0 {
		t.Fatalf("empty: %+v %v", m.Epochal, w)
	}
	mk(t, ep, "# E\n\n"+epRaw)
	m, _ = metricsOf(t, repo, h)
	if *m.Epochal.Batches != 2 || *m.Epochal.LastAt != "2026-09-30T14:08:31-03:00" {
		t.Fatalf("two batches: %+v", m.Epochal)
	}
	// over 1 MiB with an over-long line inside the raw copy: three batches, no warning
	big := "## Archived batches\n### a\n### b\n<!-- BEGIN-RAW x -->\n" + strings.Repeat("y", 1<<20+1) + "\n<!-- END-RAW x -->\n### c\n"
	mk(t, ep, big)
	m, w = metricsOf(t, repo, h)
	if *m.Epochal.Batches != 3 || len(w) != 0 {
		t.Fatalf("big: %+v %v", m.Epochal, w)
	}
	// symlink: null + warning, for all three files
	for _, f := range []string{"EPOCHAL.md", "MEMORY.md", "RISKS.md"} {
		os.Remove(filepath.Join(h, f))
		if err := os.Symlink(ep+".real", filepath.Join(h, f)); err != nil {
			t.Skip("no symlinks")
		}
	}
	m, w = metricsOf(t, repo, h)
	if m.Epochal.Batches != nil || m.Memory.Lines != nil || m.Memory.Missing || m.Incidents != nil || len(w) != 3 {
		t.Fatalf("symlink: %+v %v", m, w)
	}
}

func TestMetricsInStateGlobalMode(t *testing.T) {
	home, repo := t.TempDir(), t.TempDir()
	register(t, home, repo)
	g := filepath.Join(home, "projects", filepath.Base(repo))
	mk(t, filepath.Join(g, "MEMORY.md"), "one\n")
	mk(t, filepath.Join(g, "prd/PRD-001-a.md"), "| Status | delivered on 2026-09-25 |\n")
	st, raw := getState(t, Config{Home: home, SnapshotDir: t.TempDir(), Stale: DefaultStale})
	p := st.Projects[0]
	if p.Mode != "global" || len(p.Delivered) != 1 || *p.Metrics.Memory.Bytes != 4 || p.Metrics.PRDs[0].Status != "delivered" || !strings.Contains(raw, `"metrics":{"prds":[`) {
		t.Fatalf("state: %s", raw)
	}
}
