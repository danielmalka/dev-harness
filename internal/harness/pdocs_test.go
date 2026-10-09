package harness

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPdocsDir(t *testing.T) {
	if got := PdocsDir("/h", "p"); got != filepath.Join("/h", "projects", "p", "pdocs") {
		t.Error(got)
	}
	for _, c := range [][2]string{{"", "p"}, {"/h", ""}, {"/h", ".."}, {"/h", "a/b"}} {
		if PdocsDir(c[0], c[1]) != "" {
			t.Error(c)
		}
	}
}

func TestIsHarnessAndResolvePdocsOnly(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		want  bool
	}{
		{"empty", nil, true},
		{"pdocs only", []string{"pdocs/"}, false},
		{"pdocs+project.yaml", []string{"pdocs/", "project.yaml"}, true},
		{"pdocs+tasks", []string{"pdocs/", "tasks/"}, true},
		{"pdocs+prd", []string{"pdocs/", "prd/"}, true},
		{"pdocs+MEMORY", []string{"pdocs/", "MEMORY.md"}, true},
		{"records only", []string{"RISKS.md"}, true},
	}
	for _, c := range cases {
		for _, withRepo := range []bool{false, true} {
			home, repo := t.TempDir(), t.TempDir()
			repo, _ = canon(repo)
			d := filepath.Join(home, "projects", "p")
			os.MkdirAll(d, 0o755)
			for _, it := range c.items {
				if it[len(it)-1] == '/' {
					os.MkdirAll(filepath.Join(d, it), 0o755)
				} else {
					os.WriteFile(filepath.Join(d, it), nil, 0o644)
				}
			}
			if got := IsHarness(d); got != c.want {
				t.Errorf("%s: IsHarness=%v", c.name, got)
			}
			if withRepo {
				os.MkdirAll(filepath.Join(repo, ".harness"), 0o755)
			}
			want := "none"
			switch {
			case withRepo:
				want = "repo"
			case c.want:
				want = "global"
			}
			r := resolveIn(home, repo, Config{Projects: map[string]string{repo: "p"}})
			if r.Mode != want {
				t.Errorf("%s repo=%v: mode %s want %s", c.name, withRepo, r.Mode, want)
			}
		}
	}
	if IsHarness(filepath.Join(t.TempDir(), "missing")) {
		t.Error("missing dir")
	}
}
