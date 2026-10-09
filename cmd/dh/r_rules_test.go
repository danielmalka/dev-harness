package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/dashboard"
	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// TestMain keeps every test in this package off the real ~/.harness.
func TestMain(m *testing.M) {
	d, err := os.MkdirTemp("", "dh-cmd-test-")
	if err != nil {
		panic(err)
	}
	os.Setenv("DH_HOME", d)
	code := m.Run()
	os.RemoveAll(d)
	os.Exit(code)
}

func gitRepo(t *testing.T, dir string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("x\n"), 0o644)
	for _, a := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "i"}} {
		if b, err := exec.Command("git", append([]string{"-C", dir}, a...)...).CombinedOutput(); err != nil {
			t.Skip(string(b))
		}
	}
}

func TestR1_HomeFromDHHomeCreatedOnFirstUse(t *testing.T) {
	home, repo := env(t) // DH_HOME nonexistent
	for _, args := range [][]string{{"harness-path", "--json", repo}, {"projects"}, {"link", repo}} {
		os.RemoveAll(home)
		if code, _, e := dh(t, args...); code != 0 {
			t.Fatalf("%v: exit %d %s", args, code, e)
		}
		for _, d := range []string{"projects", "sessions", "dashboard"} {
			if st, err := os.Stat(filepath.Join(home, d)); err != nil || !st.IsDir() {
				t.Errorf("%v: %s missing", args, d)
			}
		}
		if _, err := os.Stat(filepath.Join(home, "config.yaml")); err != nil {
			t.Errorf("%v: config.yaml missing", args)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, ".harness")); err == nil {
		t.Error("repo touched")
	}
}

func TestR1_HomeDefaultUnderUserHome(t *testing.T) {
	fake := t.TempDir()
	t.Setenv("DH_HOME", "")
	t.Setenv("HOME", fake)
	t.Setenv("USERPROFILE", fake)
	repo := t.TempDir()
	if code, _, e := dh(t, "harness-path", repo); code != 0 {
		t.Fatalf("exit %d %s", code, e)
	}
	if st, err := os.Stat(filepath.Join(fake, ".harness", "sessions")); err != nil || !st.IsDir() {
		t.Error("default home not created under simulated user home")
	}
}

func TestR4_RelinkStaleRefuseLive(t *testing.T) {
	_, a := env(t)
	b := filepath.Join(t.TempDir(), "repo") // same base name
	os.MkdirAll(b, 0o755)
	if code, _, _ := dh(t, "link", a); code != 0 {
		t.Fatal("first link")
	}
	if code, _, e := dh(t, "link", b); code != 1 || !strings.Contains(e, "explicit name") {
		t.Errorf("live path must refuse: %d %q", code, e)
	}
	if code, _, _ := dh(t, "link", "repo-b", b); code != 0 {
		t.Error("explicit name must work")
	}
	os.RemoveAll(a) // original repo gone: same name relinks (replaces stale entry)
	c := filepath.Join(t.TempDir(), "repo")
	os.MkdirAll(c, 0o755)
	if code, _, e := dh(t, "link", c); code != 0 {
		t.Fatalf("stale entry should be replaced: %d %q", code, e)
	}
	_, out, _ := dh(t, "projects", "--json")
	var ps []harness.Project
	json.Unmarshal([]byte(out), &ps)
	for _, p := range ps {
		if p.Path == a {
			t.Errorf("stale entry still listed: %v", ps)
		}
	}
	if len(ps) != 2 {
		t.Errorf("want 2 projects, got %v", ps)
	}
}

func TestR5_ResolveFourSituations(t *testing.T) {
	home, base := env(t)
	_ = base
	root := t.TempDir()
	mk := func(n string) string { p := filepath.Join(root, n); os.MkdirAll(p, 0o755); return p }
	only, glob, both, none := mk("only"), mk("glob"), mk("both"), mk("none")
	os.MkdirAll(filepath.Join(only, ".harness"), 0o755)
	os.MkdirAll(filepath.Join(both, ".harness"), 0o755)
	dh(t, "link", glob)
	dh(t, "link", both)
	os.MkdirAll(filepath.Join(home, "projects", "both"), 0o755)
	cases := []struct{ dir, mode, want string }{
		{only, "repo", filepath.Join(only, ".harness")},
		{glob, "global", filepath.Join(home, "projects", "glob")},
		{both, "repo", filepath.Join(both, ".harness")},
		{none, "none", ""},
	}
	for _, c := range cases {
		code, out, _ := dh(t, "harness-path", "--json", c.dir)
		var r harness.Resolved
		if code != 0 || json.Unmarshal([]byte(out), &r) != nil || r.Mode != c.mode || r.Dir != c.want {
			t.Errorf("%s: %d %s", c.dir, code, out)
		}
		_, txt, _ := dh(t, "harness-path", c.dir)
		wantTxt := "mode: " + c.mode + "\n"
		if c.want != "" {
			wantTxt += "dir: " + c.want + "\n"
		}
		if txt != wantTxt {
			t.Errorf("%s text: %q want %q", c.dir, txt, wantTxt)
		}
	}
}

func TestR6_LinkCreateRelinkOffKeepsFolder(t *testing.T) {
	home, repo := env(t)
	dh(t, "link", repo)
	folder := filepath.Join(home, "projects", "repo")
	os.WriteFile(filepath.Join(folder, "MEMORY.md"), []byte("m"), 0o644)
	moved := filepath.Join(filepath.Dir(repo), "moved")
	os.Rename(repo, moved)
	if _, out, _ := dh(t, "harness-path", moved); out != "mode: none\n" {
		t.Errorf("after move: %q", out)
	}
	dh(t, "link", "repo", "--off", repo) // drop stale entry
	if code, _, e := dh(t, "link", "repo", moved); code != 0 {
		t.Fatalf("relink: %d %s", code, e)
	}
	if b, _ := os.ReadFile(filepath.Join(folder, "MEMORY.md")); string(b) != "m" {
		t.Error("existing folder not reused")
	}
	if _, out, _ := dh(t, "harness-path", moved); !strings.Contains(out, "mode: global") {
		t.Errorf("after relink: %q", out)
	}
	dh(t, "link", "--off", moved)
	if _, err := os.Stat(filepath.Join(folder, "MEMORY.md")); err != nil {
		t.Error("--off deleted the project folder")
	}
}

func TestR7_ZeroTraceInRepo(t *testing.T) {
	home, repo := env(t)
	gitRepo(t, repo)
	gi := filepath.Join(repo, ".gitignore")
	before, _ := os.ReadFile(gi)
	if code, _, e := dh(t, "link", repo); code != 0 {
		t.Fatalf("link: %d %s", code, e)
	}
	g := filepath.Join(home, "projects", "repo") // simulate a session writing memory and a task
	os.MkdirAll(filepath.Join(g, "tasks", "T-1"), 0o755)
	os.WriteFile(filepath.Join(g, "MEMORY.md"), []byte("m"), 0o644)
	os.WriteFile(filepath.Join(g, "tasks", "T-1", "TASK.md"), []byte("t"), 0o644)
	if st, _ := exec.Command("git", "-C", repo, "status", "--porcelain", "-uall").Output(); len(st) != 0 {
		t.Errorf("git status not empty: %s", st)
	}
	if _, err := os.Stat(filepath.Join(repo, ".harness")); err == nil {
		t.Error(".harness created in repo")
	}
	after, _ := os.ReadFile(gi)
	if string(before) != string(after) {
		t.Error(".gitignore changed")
	}
}

func TestR9_DoctorNeverWritesSettings(t *testing.T) {
	// the settings reader lives in harness.go; doctor.go's CreateTemp is the pre-existing writability probe
	for _, f := range []string{"../../internal/doctor/harness.go"} {
		b, _ := os.ReadFile(f)
		if strings.Contains(string(b), "WriteFile") || strings.Contains(string(b), "os.Create") || strings.Contains(string(b), "OpenFile") {
			t.Errorf("%s writes files", f)
		}
	}
}

func TestR10_SnapshotDirOnlyDHHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DH_HOME", home)
	t.Setenv("DEV_HARNESS_SNAPSHOT_DIR", t.TempDir())
	t.Setenv("CLAUDE_CONFIG_DIR", t.TempDir())
	if got := snapshot.SnapshotDir(); got != filepath.Join(home, "sessions") {
		t.Errorf("SnapshotDir = %q", got)
	}
}

func TestR18_ProjectsJSONParityWithAPIState(t *testing.T) {
	home, _ := env(t)
	if code, out, _ := dh(t, "projects", "--json"); code != 0 || strings.TrimSpace(out) != "[]" {
		t.Fatalf("no config: %d %q", code, out)
	}
	root := t.TempDir()
	for _, n := range []string{"b", "a", "c"} {
		p := filepath.Join(root, n)
		os.MkdirAll(p, 0o755)
		if n == "c" {
			os.MkdirAll(filepath.Join(p, ".harness"), 0o755)
		}
		dh(t, "link", p)
	}
	gone := filepath.Join(root, "gone")
	os.MkdirAll(gone, 0o755)
	dh(t, "link", gone)
	os.RemoveAll(gone) // registered but missing: omitted from both
	_, out, _ := dh(t, "projects", "--json")
	var cli []harness.Project
	if err := json.Unmarshal([]byte(out), &cli); err != nil || len(cli) != 3 {
		t.Fatalf("%v %q", err, out)
	}
	r := httptest.NewRequest("GET", "http://127.0.0.1:1/api/state", nil)
	r.Host = "127.0.0.1:1"
	w := httptest.NewRecorder()
	dashboard.Handler(dashboard.Config{Home: home, SnapshotDir: t.TempDir()}).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}
	var st struct {
		Projects []harness.Project `json:"projects"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(st.Projects, cli) {
		t.Errorf("api %+v\ncli %+v", st.Projects, cli)
	}
}

func TestR18_NoNewSlashCommands(t *testing.T) {
	ents, err := os.ReadDir("../../.commands")
	if err != nil || len(ents) != 19 {
		t.Fatalf("%v: %d files in .commands, want 19", err, len(ents))
	}
	for _, e := range ents {
		n := strings.TrimSuffix(e.Name(), ".md")
		if n == "link" || n == "projects" || n == "harness-path" {
			t.Errorf("unexpected command %s", n)
		}
	}
}

func TestR15_NoMigrateCommand(t *testing.T) {
	b, _ := os.ReadFile("main.go")
	if strings.Contains(string(b), `"migrate"`) {
		t.Error("dh migrate exists")
	}
}
