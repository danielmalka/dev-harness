package doctor

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func mk(t *testing.T, p, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// report runs the harness section for repo with an isolated home and returns the home too.
func report(t *testing.T, repo, home, settings string) string {
	t.Helper()
	t.Setenv("DH_HOME", home)
	var b bytes.Buffer
	printHarnessFor(&b, repo, home, settings)
	return b.String()
}

func TestHarnessRepoRegistered(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	out := report(t, repo, home, filepath.Join(t.TempDir(), "none.json"))
	for _, want := range []string{"home: " + home, "mode: repo", "dir: " + repo + "/.harness", "project.yaml: present", "MEMORY.md: missing"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	for _, bad := range []string{"is not registered", "both ", "ignored by git"} {
		if strings.Contains(out, bad) {
			t.Errorf("unexpected %q in:\n%s", bad, out)
		}
	}
}

func TestHarnessUnregisteredAndBoth(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	out := report(t, repo, home, "")
	if !strings.Contains(out, "warning: "+repo+" is not registered in "+home+"/config.yaml: not shown in the dashboard until dh link") {
		t.Errorf("no unregistered warning:\n%s", out)
	}
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p"), 0o755)
	out = report(t, repo, home, "")
	if !strings.Contains(out, "mode: repo") || !strings.Contains(out, "warning: both "+repo+"/.harness and "+home+"/projects/p exist: the internal one wins") {
		t.Errorf("no both warning:\n%s", out)
	}
}

func TestHarnessGlobalAndNone(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	mk(t, filepath.Join(home, "projects", "p", "project.yaml"), "x\n")
	out := report(t, repo, home, "")
	if !strings.Contains(out, "mode: global") || !strings.Contains(out, "dir: "+home+"/projects/p") || strings.Contains(out, "warning") {
		t.Errorf("global:\n%s", out)
	}
	moved := t.TempDir() // R6: repo moved, old path stays in the map
	out = report(t, moved, home, "")
	if !strings.Contains(out, "mode: none") || strings.Contains(out, "dir:") || !strings.Contains(out, "hint: no harness for this directory: run /dh:setup (or dh link)") {
		t.Errorf("none:\n%s", out)
	}
}

func TestHarnessLocalYamlIgnore(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo, home := t.TempDir(), t.TempDir()
	if err := exec.Command("git", "-C", repo, "init", "-q").Run(); err != nil {
		t.Skip("git init failed")
	}
	mk(t, filepath.Join(repo, ".harness", "local.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	const warn = "warning: .harness/local.yaml is not ignored by git"
	if out := report(t, repo, home, ""); !strings.Contains(out, warn) {
		t.Errorf("want warning:\n%s", out)
	}
	mk(t, filepath.Join(repo, ".gitignore"), ".harness/local.yaml\n")
	if out := report(t, repo, home, ""); strings.Contains(out, warn) {
		t.Errorf("unexpected warning:\n%s", out)
	}
	mk(t, filepath.Join(repo, ".gitignore"), ".harness/\n") // whole folder ignored: no warning either
	if out := report(t, repo, home, ""); strings.Contains(out, warn) || strings.Contains(out, "ignored by git") {
		t.Errorf("unexpected warning:\n%s", out)
	}
}

func TestHarnessSettingsHint(t *testing.T) {
	repo, home, cfgDir := t.TempDir(), t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	const hint = "hint: ~/.claude/settings.json has no permissions.additionalDirectories entry covering " // + home
	settings := filepath.Join(cfgDir, "settings.json")
	if out := report(t, repo, home, settings); !strings.Contains(out, hint+home+`; add "`+home+`" (dh never writes that file)`) {
		t.Errorf("missing file: want hint:\n%s", out)
	}
	mk(t, settings, `{"permissions":{"additionalDirectories":["/elsewhere"]}}`)
	if out := report(t, repo, home, settings); !strings.Contains(out, hint) {
		t.Errorf("no entry: want hint:\n%s", out)
	}
	mk(t, settings, `{not json`)
	if out := report(t, repo, home, settings); !strings.Contains(out, "could not be parsed") {
		t.Errorf("bad json: want parse hint:\n%s", out)
	}
	mk(t, settings, `{"permissions":{"additionalDirectories":["`+home+`"]}}`)
	if out := report(t, repo, home, settings); strings.Contains(out, "hint:") {
		t.Errorf("covered: unexpected hint:\n%s", out)
	}
	// home inside the repository: never a hint
	inside := filepath.Join(repo, "dh-home")
	os.MkdirAll(inside, 0o755)
	if out := report(t, repo, inside, filepath.Join(cfgDir, "absent.json")); strings.Contains(out, "settings.json") {
		t.Errorf("home inside repo: unexpected hint:\n%s", out)
	}
	// "~/..." entries expand against the user home
	uh := t.TempDir()
	t.Setenv("HOME", uh)
	mk(t, settings, `{"permissions":{"additionalDirectories":["~/.harness"]}}`)
	if out := report(t, repo, filepath.Join(uh, ".harness"), settings); strings.Contains(out, "hint:") {
		t.Errorf("tilde entry: unexpected hint:\n%s", out)
	}
}

func TestHarnessUnknownHome(t *testing.T) {
	repo := t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	t.Setenv("DH_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	var b bytes.Buffer
	printHarnessFor(&b, repo, "", "")
	if out := b.String(); !strings.Contains(out, "home: unknown") || !strings.Contains(out, "set DH_HOME") {
		t.Errorf("unknown home:\n%s", out)
	}
}

func TestSettingsHintSuggestsCustomHome(t *testing.T) {
	uh, custom := t.TempDir(), t.TempDir()
	t.Setenv("HOME", uh)
	if h := settingsHint(filepath.Join(uh, "none.json"), custom); !strings.Contains(h, `add "`+custom+`"`) {
		t.Errorf("custom home: %s", h)
	}
	if h := settingsHint(filepath.Join(uh, "none.json"), filepath.Join(uh, ".harness")); !strings.Contains(h, `add "~/.harness"`) {
		t.Errorf("default home: %s", h)
	}
}

func TestHarnessPdocsOnlyNoBothWarning(t *testing.T) {
	repo, home := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repo, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+repo+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p", "pdocs"), 0o755)
	out := report(t, repo, home, "")
	if strings.Contains(out, "both ") {
		t.Errorf("unexpected both warning:\n%s", out)
	}
	mk(t, filepath.Join(home, "projects", "p", "project.yaml"), "x\n")
	if out = report(t, repo, home, ""); !strings.Contains(out, "warning: both ") {
		t.Errorf("no both warning:\n%s", out)
	}
}

func TestHarnessSymlinkedCwd(t *testing.T) {
	real, home := t.TempDir(), t.TempDir()
	real, _ = filepath.EvalSymlinks(real)
	link := filepath.Join(t.TempDir(), "ln")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("symlinks unavailable")
	}
	mk(t, filepath.Join(real, ".harness", "project.yaml"), "x\n")
	mk(t, filepath.Join(home, "config.yaml"), "projects:\n  \""+real+"\": p\n")
	os.MkdirAll(filepath.Join(home, "projects", "p"), 0o755)
	out := report(t, link, home, "")
	if strings.Contains(out, "is not registered") || !strings.Contains(out, "warning: both ") {
		t.Errorf("symlinked cwd:\n%s", out)
	}
}
