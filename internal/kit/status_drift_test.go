package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func driftTree(t *testing.T, prdStatus string, tickets map[string]string) string {
	t.Helper()
	t.Setenv("DH_HOME", t.TempDir())
	root := t.TempDir()
	w := func(rel, body string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	w("docs/prd/PRD-003-x.md", "| Status | "+prdStatus+" |\n")
	w("docs/prd/PRD-003-x.review.md", "| Status | entregue em 1 |\n")
	for name, row := range tickets {
		w(".harness/tasks/"+name+"/TASK.md", row)
	}
	return root
}

func drift(root string) []string {
	var errs []string
	checkStatusDrift(root, &errs)
	return errs
}

func TestStatusDrift(t *testing.T) {
	cases := []struct {
		name, prd string
		tickets   map[string]string
		want      string // substring of the single expected error; "" = clean
	}{
		{"a open ticket on delivered PRD", "entregue em 0.1.0", map[string]string{
			"T-1": "| Status | pronta |\n| Story / PRD | PRD-003 |\n"}, "T-1 open"},
		{"b all done, PRD open", "aprovado", map[string]string{
			"T-1": "| Status | concluída |\n| PRD (RF-<n>) | PRD-003 (R1) |\n"}, "all 1 linked tickets are done"},
		{"clean open", "aprovado", map[string]string{
			"T-1": "| Status | concluída |\n| Story / PRD | PRD-003 |\n",
			"T-2": "| Status | pronta |\n| Story / PRD | PRD-003 |\n"}, ""},
		{"clean delivered", "entregue em 0.1.0", map[string]string{
			"T-1": "| Status | concluída |\n| Story / PRD | PRD-003 |\n"}, ""},
		{"fora does not link", "entregue em 0.1.0", map[string]string{
			"T-1": "| Status | pronta |\n| Story / PRD | fora do PRD-003 (x) |\n"}, ""},
		{"no tickets", "aprovado", nil, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			errs := drift(driftTree(t, c.prd, c.tickets))
			if c.want == "" {
				if len(errs) != 0 {
					t.Fatalf("want clean, got %v", errs)
				}
				return
			}
			if len(errs) != 1 || !strings.Contains(errs[0], c.want) {
				t.Fatalf("want one error containing %q, got %v", c.want, errs)
			}
		})
	}
}

func TestStatusDriftNoTasksDir(t *testing.T) {
	t.Setenv("DH_HOME", t.TempDir())
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "docs", "prd"), 0o755)
	os.WriteFile(filepath.Join(root, "docs", "prd", "PRD-001-a.md"), []byte("| Status | entregue em 1 |\n"), 0o644)
	if errs := drift(root); len(errs) != 0 {
		t.Fatalf("got %v", errs)
	}
}

func TestStatusDriftGlobalProject(t *testing.T) {
	home, root := t.TempDir(), t.TempDir()
	t.Setenv("DH_HOME", home)
	w := func(p, body string) {
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	w(filepath.Join(root, "docs", "prd", "PRD-003-x.md"), "| Status | aprovado |\n")
	w(filepath.Join(home, "projects", "p", "tasks", "T-1", "TASK.md"), "| Status | concluída |\n| Story / PRD | PRD-003 |\n")
	if errs := drift(root); len(errs) != 0 { // not registered: none, silent
		t.Fatalf("unregistered: got %v", errs)
	}
	w(filepath.Join(home, "config.yaml"), "projects:\n  \""+root+"\": p\n")
	if errs := drift(root); len(errs) != 1 || !strings.Contains(errs[0], "all 1 linked tickets are done") {
		t.Fatalf("global drift: got %v", errs)
	}
	os.RemoveAll(filepath.Join(home, "projects", "p", "tasks"))
	if errs := drift(root); len(errs) != 0 {
		t.Fatalf("no tasks: got %v", errs)
	}
}
