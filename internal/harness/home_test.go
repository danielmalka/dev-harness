package harness

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setup(t *testing.T) (home, repo string) {
	t.Helper()
	base := t.TempDir()
	home = filepath.Join(base, "h")
	t.Setenv("DH_HOME", home)
	repo = filepath.Join(base, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	return
}

func TestEnsureHomeCreatesAndIsIdempotent(t *testing.T) {
	home, _ := setup(t)
	for i := 0; i < 2; i++ {
		got, err := EnsureHome()
		if err != nil || got != home {
			t.Fatalf("EnsureHome = %q, %v", got, err)
		}
	}
	for _, d := range []string{"projects", "sessions", "dashboard"} {
		if !isDir(filepath.Join(home, d)) {
			t.Errorf("%s missing", d)
		}
	}
	b, _ := os.ReadFile(filepath.Join(home, "config.yaml"))
	if string(b) != "projects:\n" {
		t.Errorf("config = %q", b)
	}
	if SessionsDir() != filepath.Join(home, "sessions") || DashboardDir() != filepath.Join(home, "dashboard") {
		t.Error("dirs")
	}
}

func TestHomeDefault(t *testing.T) {
	t.Setenv("DH_HOME", "")
	fake := t.TempDir()
	t.Setenv("HOME", fake)
	t.Setenv("USERPROFILE", fake)
	if got := Home(); got != filepath.Join(fake, ".harness") {
		t.Errorf("Home = %q", got)
	}
}

const sample = `# my config
language: pt-br
mode: global
dashboard:
  sprites: "C:\\sprites dir\\x"
  other: keep
reviewers:
  code:
    - claude
    - cli:codex/gpt
  # note
extra: 1
projects:
  # first
  "C:\\Users\\me\\repo": win
  "/a/b": ab
`

func TestNoHome(t *testing.T) {
	t.Setenv("DH_HOME", "")
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	if Home() != "" || SessionsDir() != "" || DashboardDir() != "" {
		t.Fatalf("want empty dirs, got %q %q %q", Home(), SessionsDir(), DashboardDir())
	}
	if _, err := EnsureHome(); err == nil {
		t.Fatal("EnsureHome: want error")
	}
}

func TestHomeOrder(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	t.Setenv("DH_HOME", "")
	t.Setenv("HOME", a)
	t.Setenv("USERPROFILE", b)
	if got := Home(); got != filepath.Join(a, ".harness") {
		t.Errorf("HOME must win, got %q", got)
	}
	t.Setenv("HOME", "")
	if got := Home(); got != filepath.Join(b, ".harness") {
		t.Errorf("USERPROFILE fallback, got %q", got)
	}
}

func TestEnsureHomeTightensModes(t *testing.T) {
	home := filepath.Join(t.TempDir(), "h")
	t.Setenv("DH_HOME", home)
	os.MkdirAll(filepath.Join(home, "sessions"), 0o755)
	os.Chmod(home, 0o755)
	os.Chmod(filepath.Join(home, "sessions"), 0o755)
	if _, err := EnsureHome(); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{home, filepath.Join(home, "sessions")} {
		if fi, _ := os.Stat(d); fi.Mode().Perm()&0o077 != 0 {
			t.Errorf("%s mode %v", d, fi.Mode().Perm())
		}
	}
}

func TestExistingModesTightened(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX modes")
	}
	home := filepath.Join(t.TempDir(), "h")
	t.Setenv("DH_HOME", home)
	proj := filepath.Join(home, "projects", "old")
	os.MkdirAll(proj, 0o755)
	os.Chmod(proj, 0o755)
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects:\n"), 0o644)
	os.Chmod(filepath.Join(home, "config.yaml"), 0o644)
	rec := filepath.Join(proj, "MEMORY.md")
	os.WriteFile(rec, nil, 0o644)
	os.Chmod(rec, 0o644)
	if _, err := EnsureHome(); err != nil {
		t.Fatal(err)
	}
	mode := func(p string) os.FileMode { fi, _ := os.Stat(p); return fi.Mode().Perm() }
	if mode(filepath.Join(home, "config.yaml")) != 0o600 || mode(proj) != 0o700 {
		t.Errorf("config %v project %v", mode(filepath.Join(home, "config.yaml")), mode(proj))
	}
	if mode(rec) != 0o644 {
		t.Errorf("record touched: %v", mode(rec))
	}
	// Link tightens an existing project folder too
	repo := t.TempDir()
	d := filepath.Join(home, "projects", "other")
	os.MkdirAll(d, 0o755)
	os.Chmod(d, 0o755)
	if _, err := Link(home, repo, "other", false); err != nil {
		t.Fatal(err)
	}
	if mode(d) != 0o700 {
		t.Errorf("link left %v", mode(d))
	}
}

func TestEnsureHomeSymlinkedHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX modes")
	}
	real := filepath.Join(t.TempDir(), "real")
	os.Mkdir(real, 0o755)
	os.Chmod(real, 0o755)
	link := filepath.Join(t.TempDir(), "lnk")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("symlink unsupported")
	}
	t.Setenv("DH_HOME", link)
	if _, err := EnsureHome(); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Stat(real); fi.Mode().Perm() != 0o700 {
		t.Errorf("target mode %v", fi.Mode().Perm())
	}
}

func TestTightenIgnoresChmodFailure(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs a POSIX non-root user")
	}
	if fi, err := os.Stat("/etc"); err != nil || fi.Mode().Perm()&0o077 == 0 {
		t.Skip("/etc not usable as a foreign-owned directory")
	}
	if err := tighten("/etc", dirPerm); err != nil {
		t.Fatalf("chmod failure must be skipped: %v", err)
	}
}
