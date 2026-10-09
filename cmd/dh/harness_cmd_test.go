package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func dh(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := run(args, strings.NewReader(""), &out, &errb)
	return code, out.String(), errb.String()
}

func env(t *testing.T) (home, repo string) {
	base := t.TempDir()
	home = filepath.Join(base, "nothere", "h") // nonexistent
	t.Setenv("DH_HOME", home)
	repo = filepath.Join(base, "repo")
	os.MkdirAll(repo, 0o755)
	return
}

func TestHarnessPathCreatesHomeAndReportsNone(t *testing.T) {
	home, repo := env(t)
	code, out, _ := dh(t, "harness-path", "--json", repo)
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	var got struct {
		Mode, Dir string
		Defaults  map[string]string
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err, out)
	}
	if got.Mode != "none" || got.Dir != "" || len(got.Defaults) != 3 {
		t.Errorf("%+v", got)
	}
	if _, err := os.Stat(filepath.Join(home, "sessions")); err != nil {
		t.Error("sessions missing")
	}
	if entries, _ := os.ReadDir(repo); len(entries) != 0 {
		t.Error("repo touched")
	}
	code, out, _ = dh(t, "harness-path", repo)
	if code != 0 || out != "mode: none\n" {
		t.Errorf("text: %d %q", code, out)
	}
	if code, _, _ := dh(t, "harness-path", filepath.Join(repo, "missing")); code != 1 {
		t.Error("missing dir should exit 1")
	}
}

func TestDefaultsVerbatimAndRegistryNotWritten(t *testing.T) {
	home, repo := env(t)
	dh(t, "harness-path", repo)
	cfg := "language: pt-br\nmode: global\nreviewers:\n  code:\n    - claude\nprojects:\n"
	p := filepath.Join(home, "config.yaml")
	os.WriteFile(p, []byte(cfg), 0o600)
	_, out, _ := dh(t, "harness-path", "--json", repo)
	var got struct{ Defaults map[string]string }
	json.Unmarshal([]byte(out), &got)
	if got.Defaults["language"] != "pt-br" || got.Defaults["mode"] != "global" ||
		got.Defaults["reviewers"] != "reviewers:\n  code:\n    - claude" {
		t.Errorf("%v", got.Defaults)
	}
	if b, _ := os.ReadFile(p); string(b) != cfg {
		t.Error("config rewritten")
	}
}

func TestLinkProjectsAndZeroTrace(t *testing.T) {
	home, repo := env(t)
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	if b, err := exec.Command("git", "-C", repo, "init", "-q").CombinedOutput(); err != nil {
		t.Skip(string(b))
	}
	code, out, _ := dh(t, "link", "--json", repo)
	want := filepath.Join(home, "projects", "repo")
	var r struct{ Mode, Dir string }
	json.Unmarshal([]byte(out), &r)
	if code != 0 || r.Mode != "global" || r.Dir != want {
		t.Fatalf("%d %s", code, out)
	}
	st, _ := exec.Command("git", "-C", repo, "status", "--porcelain", "-uall").Output()
	if len(st) != 0 {
		t.Errorf("git status not empty: %s", st)
	}
	if _, err := os.Stat(filepath.Join(repo, ".harness")); err == nil {
		t.Error(".harness created")
	}
	_, out, _ = dh(t, "projects", "--json")
	var ps []map[string]string
	json.Unmarshal([]byte(out), &ps)
	if len(ps) != 1 || ps[0]["name"] != "repo" || ps[0]["mode"] != "global" || ps[0]["harness"] != want || ps[0]["path"] != repo {
		t.Errorf("%v", ps)
	}
	_, out, _ = dh(t, "projects")
	if out != "repo\tglobal\t"+repo+"\t"+want+"\n" {
		t.Errorf("text %q", out)
	}
	code, out, _ = dh(t, "link", "--off", repo)
	if code != 0 || out != "mode: none\n" {
		t.Errorf("off: %d %q", code, out)
	}
	if _, out, _ = dh(t, "projects", "--json"); strings.TrimSpace(out) != "[]" {
		t.Errorf("after off: %q", out)
	}
}

func TestProjectsEmptyWithoutConfig(t *testing.T) {
	env(t)
	code, out, _ := dh(t, "projects", "--json")
	if code != 0 || strings.TrimSpace(out) != "[]" {
		t.Errorf("%d %q", code, out)
	}
}

func TestLinkCollisionExitsOne(t *testing.T) {
	_, repo := env(t)
	other := filepath.Join(t.TempDir(), "repo")
	os.MkdirAll(other, 0o755)
	if code, _, _ := dh(t, "link", repo); code != 0 {
		t.Fatal("first link")
	}
	code, _, e := dh(t, "link", other)
	if code != 1 || !strings.Contains(e, "explicit name") {
		t.Errorf("%d %q", code, e)
	}
	if code, _, _ := dh(t, "link", "other-name", other); code != 0 {
		t.Error("explicit name should work")
	}
}

func TestUsageListsNewCommands(t *testing.T) {
	_, _, e := dh(t)
	for _, c := range []string{"harness-path", "link", "projects"} {
		if !strings.Contains(e, c) {
			t.Errorf("usage lacks %s", c)
		}
	}
}

func TestHarnessPathPdocsField(t *testing.T) {
	home, repo := env(t)
	get := func() (mode, dir, pdocs string) {
		code, out, _ := dh(t, "harness-path", "--json", repo)
		var g struct{ Mode, Dir, Pdocs string }
		if code != 0 || json.Unmarshal([]byte(out), &g) != nil {
			t.Fatalf("%d %q", code, out)
		}
		return g.Mode, g.Dir, g.Pdocs
	}
	if m, d, p := get(); m != "none" || d != "" || p != "" {
		t.Errorf("unregistered: %s %q %q", m, d, p)
	}
	name := "proj"
	if c, _, e := dh(t, "link", name, repo); c != 0 {
		t.Fatal(e)
	}
	want := filepath.Join(home, "projects", name, "pdocs")
	if m, _, p := get(); m != "global" || p != want {
		t.Errorf("global: %s %q want %q", m, p, want)
	}
	os.MkdirAll(filepath.Join(repo, ".harness"), 0o755)
	if m, d, p := get(); m != "repo" || p != want || d != filepath.Join(repo, ".harness") {
		t.Errorf("repo: %s %q %q", m, d, p)
	}
	if _, out, _ := dh(t, "harness-path", repo); strings.Contains(out, "pdocs") {
		t.Errorf("text output has pdocs: %q", out)
	}
	if _, out, _ := dh(t, "link", "--json", repo); strings.Contains(out, "pdocs") {
		t.Errorf("link --json has pdocs: %q", out)
	}
}
