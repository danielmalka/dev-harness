package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type hpOut struct {
	Mode, Dir, Pdocs string
	Defaults         map[string]string
	LinkedWorktree   *struct{ Main, Gitdir, Commondir string } `json:"linked_worktree"`
	raw              map[string]json.RawMessage
}

func hp(t *testing.T, dir string) hpOut {
	t.Helper()
	code, out, e := dh(t, "harness-path", "--json", dir)
	if code != 0 {
		t.Fatalf("harness-path %s: exit %d %s", dir, code, e)
	}
	var g hpOut
	if err := json.Unmarshal([]byte(out), &g); err != nil {
		t.Fatal(err, out)
	}
	json.Unmarshal([]byte(out), &g.raw)
	return g
}

func wf(t *testing.T, p, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// fakeWT builds <base>/main (with a .git directory) and a valid linked worktree <base>/wts/T-1.
func fakeWT(t *testing.T) (main, wt, admin string) {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	main = filepath.Join(base, "main")
	g := filepath.Join(main, ".git")
	wf(t, filepath.Join(g, "HEAD"), "ref: refs/heads/main\n")
	os.MkdirAll(filepath.Join(g, "objects"), 0o755)
	os.MkdirAll(filepath.Join(g, "refs"), 0o755)
	wt = filepath.Join(base, "wts", "T-1")
	admin = filepath.Join(g, "worktrees", "T-1")
	wf(t, filepath.Join(admin, "commondir"), "../..\n")
	wf(t, filepath.Join(admin, "gitdir"), filepath.Join(wt, ".git")+"\n")
	wf(t, filepath.Join(wt, ".git"), "gitdir: "+admin+"\n")
	return
}

func cfgPath(t *testing.T) string {
	return filepath.Join(os.Getenv("DH_HOME"), "config.yaml")
}

// R1, R2, R3: dh config set/unset worktree.
func TestR1_R2_R3_ConfigSetUnsetWorktree(t *testing.T) {
	env(t)
	p := cfgPath(t)
	for _, v := range []string{"always", "on-request"} {
		if code, _, e := dh(t, "config", "set", "worktree", v); code != 0 {
			t.Fatalf("set %s: %s", v, e)
		}
		if b, _ := os.ReadFile(p); !strings.Contains(string(b), "worktree: "+v+"\n") {
			t.Errorf("R1 %s not stored:\n%s", v, b)
		}
	}
	before, _ := os.ReadFile(p)
	code, _, e := dh(t, "config", "set", "worktree", "sometimes")
	if code != 1 || !strings.Contains(e, "always") || !strings.Contains(e, "on-request") {
		t.Errorf("R2: %d %q", code, e)
	}
	if after, _ := os.ReadFile(p); string(after) != string(before) {
		t.Errorf("R2 file changed:\n%s", after)
	}
	if code, _, _ := dh(t, "config", "unset", "worktree"); code != 0 {
		t.Fatal("unset")
	}
	if b, _ := os.ReadFile(p); strings.Contains(string(b), "worktree") {
		t.Errorf("R3 key left:\n%s", b)
	}
	if g := hp(t, t.TempDir()); g.Defaults["worktree"] != "on-request" {
		t.Errorf("R3 effective value after unset: %v", g.Defaults)
	}
}

// R4: text and JSON output of dh config carry worktree (empty when absent).
func TestR4_ConfigShowsWorktree(t *testing.T) {
	env(t)
	_, out, _ := dh(t, "config")
	if !strings.HasSuffix(out, "worktree: \n") {
		t.Errorf("absent text: %q", out)
	}
	_, out, _ = dh(t, "config", "--json")
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatal(err, out)
	}
	if v, ok := m["worktree"]; !ok || v != "" {
		t.Errorf("absent json: %q", out)
	}
	dh(t, "config", "set", "worktree", "always")
	_, out, _ = dh(t, "config")
	if !strings.HasSuffix(out, "worktree: always\n") {
		t.Errorf("set text: %q", out)
	}
	_, out, _ = dh(t, "config", "--json")
	json.Unmarshal([]byte(out), &m)
	if m["worktree"] != "always" {
		t.Errorf("set json: %q", out)
	}
}

// R5, R7 (Go side): defaults.worktree is the global value, else on-request; local.yaml never enters.
func TestR5_R7_DefaultsWorktree(t *testing.T) {
	home, repo := env(t)
	if g := hp(t, repo); g.Defaults["worktree"] != "on-request" {
		t.Errorf("no config.yaml: %v", g.Defaults)
	}
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("worktree: always\n"), 0o600)
	if g := hp(t, repo); g.Defaults["worktree"] != "always" {
		t.Errorf("global: %v", g.Defaults)
	}
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("language: en\n"), 0o600)
	harnessDir := filepath.Join(repo, ".harness")
	wf(t, filepath.Join(harnessDir, "project.yaml"), "language: en\n")
	wf(t, filepath.Join(harnessDir, "local.yaml"), "worktree: always\n")
	if g := hp(t, repo); g.Mode != "repo" || g.Defaults["worktree"] != "on-request" {
		t.Errorf("local.yaml leaked: %+v", g)
	}
	// the harness dir of a global-mode project too
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("worktree: on-request\n"), 0o600)
	if g := hp(t, repo); g.Defaults["worktree"] != "on-request" {
		t.Errorf("%v", g.Defaults)
	}
}

// R8 (Go side): any value loads; an unquoted ':' aborts like language does; the quoted form loads.
func TestR8_HarnessPathScalarRules(t *testing.T) {
	home, repo := env(t)
	os.MkdirAll(home, 0o700)
	p := filepath.Join(home, "config.yaml")
	os.WriteFile(p, []byte("worktree: sometimes\n"), 0o600)
	if g := hp(t, repo); g.Defaults["worktree"] != "sometimes" {
		t.Errorf("sometimes: %v", g.Defaults)
	}
	os.WriteFile(p, []byte("worktree: \"a:b\"\n"), 0o600)
	if g := hp(t, repo); g.Defaults["worktree"] != "a:b" {
		t.Errorf("quoted: %v", g.Defaults)
	}
	for _, key := range []string{"worktree", "language"} {
		os.WriteFile(p, []byte(key+": a:b\n"), 0o600)
		code, _, e := dh(t, "harness-path", "--json", repo)
		if code == 0 || !strings.Contains(e, "':'") {
			t.Errorf("%s: a:b exit %d %q", key, code, e)
		}
	}
}

// R61: the field exists and is null for a normal checkout, a plain dir and a registered project.
// R21: those resolve exactly as before (R27 Go side: linked_worktree null).
func TestR21_R27_R61_NotAWorktreeIsNullAndUnchanged(t *testing.T) {
	_, repo := env(t)
	g := hp(t, repo)
	if g.Mode != "none" || g.Dir != "" || g.LinkedWorktree != nil || string(g.raw["linked_worktree"]) != "null" {
		t.Errorf("plain dir: %+v raw=%s", g, g.raw["linked_worktree"])
	}
	// normal checkout: .git is a directory
	wf(t, filepath.Join(repo, ".git", "HEAD"), "ref: refs/heads/main\n")
	os.MkdirAll(filepath.Join(repo, ".git", "objects"), 0o755)
	os.MkdirAll(filepath.Join(repo, ".git", "refs"), 0o755)
	wf(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	g = hp(t, repo)
	if g.Mode != "repo" || g.Dir != filepath.Join(repo, ".harness") || string(g.raw["linked_worktree"]) != "null" {
		t.Errorf("normal checkout: %+v", g)
	}
	code, out, _ := dh(t, "harness-path", repo)
	if code != 0 || out != "mode: repo\ndir: "+filepath.Join(repo, ".harness")+"\n" {
		t.Errorf("text output changed: %q", out)
	}
}

// R27 (Go side), R61: refused layouts report null and keep today's resolution.
func TestR27_R61_RefusedLayoutIsNull(t *testing.T) {
	env(t)
	main, wt, admin := fakeWT(t)
	wf(t, filepath.Join(main, ".harness", "project.yaml"), "main\n")
	wf(t, filepath.Join(wt, ".harness", "project.yaml"), "copy\n")
	wf(t, filepath.Join(admin, "gitdir"), filepath.Join(wt, "other", ".git")+"\n") // back-pointer mismatch
	g := hp(t, wt)
	if g.LinkedWorktree != nil || string(g.raw["linked_worktree"]) != "null" || g.Dir != filepath.Join(wt, ".harness") {
		t.Errorf("mismatch: %+v", g)
	}
	// symlinked .git
	wf(t, filepath.Join(admin, "gitdir"), filepath.Join(wt, ".git")+"\n")
	real := filepath.Join(filepath.Dir(wt), "gitfile")
	os.Rename(filepath.Join(wt, ".git"), real)
	if err := os.Symlink(real, filepath.Join(wt, ".git")); err != nil {
		t.Skip("symlink unavailable")
	}
	g = hp(t, wt)
	if g.LinkedWorktree != nil || g.Dir != filepath.Join(wt, ".harness") {
		t.Errorf("symlinked .git: %+v", g)
	}
}

// R18, R19, R26, R61: repo mode; the worktree's own .harness copy is ignored.
func TestR18_R19_R26_R61_RepoModeLinkedWorktree(t *testing.T) {
	_, _ = env(t)
	main, wt, admin := fakeWT(t)
	wf(t, filepath.Join(main, ".harness", "project.yaml"), "main\n")
	wf(t, filepath.Join(wt, ".harness", "project.yaml"), "copy\n")
	gm, gw := hp(t, main), hp(t, wt)
	if gw.Mode != "repo" || gw.Dir != filepath.Join(main, ".harness") || gw.Dir != gm.Dir || gw.Pdocs != gm.Pdocs {
		t.Errorf("wt %+v main %+v", gw, gm)
	}
	lw := gw.LinkedWorktree
	if lw == nil || lw.Main != main || lw.Gitdir != admin || lw.Commondir != filepath.Join(main, ".git") {
		t.Fatalf("linked_worktree %+v", lw)
	}
	if gm.LinkedWorktree != nil {
		t.Error("main checkout reports a linked worktree")
	}
}

// R18, R61: global mode; dir and pdocs equal the main checkout's, and pdocs is the registered name's.
func TestR18_R61_GlobalModeLinkedWorktree(t *testing.T) {
	home, _ := env(t)
	main, wt, admin := fakeWT(t)
	if c, _, e := dh(t, "link", "proj", main); c != 0 {
		t.Fatal(e)
	}
	gm, gw := hp(t, main), hp(t, wt)
	if gm.Mode != "global" || gw.Mode != "global" || gw.Dir != gm.Dir || gw.Dir != filepath.Join(home, "projects", "proj") {
		t.Errorf("dir: wt %+v main %+v", gw, gm)
	}
	if gw.Pdocs != gm.Pdocs || gw.Pdocs != filepath.Join(home, "projects", "proj", "pdocs") {
		t.Errorf("pdocs: %q vs %q", gw.Pdocs, gm.Pdocs)
	}
	if gw.LinkedWorktree == nil || gw.LinkedWorktree.Gitdir != admin || gw.LinkedWorktree.Main != main {
		t.Errorf("%+v", gw.LinkedWorktree)
	}
	// a worktree whose main is NOT registered stays none, with the field filled
	dh(t, "link", "--off", main)
	if g := hp(t, wt); g.Mode != "none" || g.Pdocs != "" || g.LinkedWorktree == nil {
		t.Errorf("unregistered: %+v", g)
	}
}

// R18, R19, R26, R61 with a real `git worktree add`, fully isolated from the user's git config.
func TestR18_R19_R26_R61_RealGitWorktree(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("not-run: git not in PATH")
	}
	base, _ := filepath.EvalSymlinks(t.TempDir())
	empty := filepath.Join(base, "gitconfig")
	wf(t, empty, "")
	t.Setenv("GIT_CONFIG_GLOBAL", empty)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("HOME", base)
	for _, k := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE"} {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
	g := func(dir string, a ...string) {
		t.Helper()
		c := exec.Command(git, append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false"}, a...)...)
		if b, err := c.CombinedOutput(); err != nil {
			t.Skipf("not-run: git %v: %v %s", a, err, b)
		}
	}
	for _, mode := range []string{"repo", "global"} {
		t.Run(mode, func(t *testing.T) {
			home, _ := env(t)
			main := filepath.Join(base, mode, "main")
			os.MkdirAll(main, 0o755)
			g(main, "init", "-q")
			wf(t, filepath.Join(main, "README.md"), "x\n")
			if mode == "repo" {
				wf(t, filepath.Join(main, ".harness", "project.yaml"), "x\n") // committed: the worktree gets a copy
			}
			g(main, "add", ".")
			g(main, "commit", "-q", "-m", "i")
			wt := filepath.Join(base, mode, "main.worktrees", "T-X")
			g(main, "worktree", "add", "-q", "-b", "t-x", wt)
			if mode == "global" {
				if c, _, e := dh(t, "link", "proj", main); c != 0 {
					t.Fatal(e)
				}
			} else if _, err := os.Stat(filepath.Join(wt, ".harness", "project.yaml")); err != nil {
				t.Fatal("fixture: worktree has no .harness copy")
			}
			gm, gw := hp(t, main), hp(t, wt)
			if gw.Mode != mode || gw.Dir != gm.Dir || gw.Pdocs != gm.Pdocs || gm.LinkedWorktree != nil {
				t.Fatalf("wt %+v main %+v", gw, gm)
			}
			if mode == "global" && gw.Dir != filepath.Join(home, "projects", "proj") {
				t.Errorf("dir %s", gw.Dir)
			}
			lw := gw.LinkedWorktree
			if lw == nil || lw.Main != main || lw.Commondir != filepath.Join(main, ".git") ||
				lw.Gitdir != filepath.Join(main, ".git", "worktrees", "T-X") {
				t.Errorf("linked_worktree %+v", lw)
			}
		})
	}
}
