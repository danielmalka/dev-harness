package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLinkCreateRelinkOff(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	r, err := Link(home, repo, "", false)
	if err != nil || r.Mode != "global" || r.Dir != filepath.Join(home, "projects", "repo") {
		t.Fatalf("%+v %v", r, err)
	}
	if entries, _ := os.ReadDir(repo); len(entries) != 0 {
		t.Errorf("repo touched: %v", entries)
	}
	os.WriteFile(filepath.Join(r.Dir, "MEMORY.md"), []byte("x"), 0o600)
	// idempotent, and "." style paths do not duplicate
	if _, err := Link(home, filepath.Join(repo, "."), "", false); err != nil {
		t.Fatal(err)
	}
	if c, _ := LoadConfig(home); len(c.Projects) != 1 {
		t.Errorf("projects: %v", c.Projects)
	}
	// --off keeps the folder
	r, err = Link(home, repo, "", true)
	if err != nil || r.Mode != "none" {
		t.Fatalf("off: %+v %v", r, err)
	}
	if _, err := os.Stat(filepath.Join(home, "projects", "repo", "MEMORY.md")); err != nil {
		t.Error("folder not preserved")
	}
	// relink reuses it
	if r, err = Link(home, repo, "", false); err != nil || r.Mode != "global" {
		t.Fatalf("relink: %+v %v", r, err)
	}
	if b, _ := os.ReadFile(filepath.Join(r.Dir, "MEMORY.md")); string(b) != "x" {
		t.Error("memory lost")
	}
}

func TestLinkRepoModeOnlyRegisters(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	os.MkdirAll(filepath.Join(repo, ".harness"), 0o755)
	r, err := Link(home, repo, "", false)
	if err != nil || r.Mode != "repo" {
		t.Fatalf("%+v %v", r, err)
	}
	if isDir(filepath.Join(home, "projects", "repo")) {
		t.Error("global folder created in repo mode")
	}
	if c, _ := LoadConfig(home); c.Projects[repo] != "repo" {
		t.Error("not registered")
	}
}

func TestLinkNameCollision(t *testing.T) {
	home, _ := setup(t)
	EnsureHome()
	base := t.TempDir()
	a := filepath.Join(base, "a", "app")
	b := filepath.Join(base, "b", "app")
	os.MkdirAll(a, 0o755)
	os.MkdirAll(b, 0o755)
	if _, err := Link(home, a, "", false); err != nil {
		t.Fatal(err)
	}
	_, err := Link(home, b, "", false)
	if err == nil || !strings.Contains(err.Error(), "explicit name") {
		t.Fatalf("want refusal, got %v", err)
	}
	if r, err := Link(home, b, "app-b", false); err != nil || r.Mode != "global" {
		t.Fatalf("explicit: %+v %v", r, err)
	}
	// old path gone -> replaced
	os.RemoveAll(a)
	if _, err := Link(home, b, "app", false); err != nil {
		t.Fatalf("stale replace: %v", err)
	}
	c, _ := LoadConfig(home)
	if _, ok := c.Projects[a]; ok || c.Projects[b] != "app" {
		t.Errorf("projects: %v", c.Projects)
	}
}

func TestLinkRejectsBadNames(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	for _, n := range []string{"../x", "a/b", `a\b`, "..", "."} {
		if _, err := Link(home, repo, n, false); err == nil {
			t.Errorf("name %q accepted", n)
		}
	}
	if _, err := Link(home, filepath.Join(repo, "missing"), "", false); err == nil {
		t.Error("missing dir accepted")
	}
}

func TestLinkOffOnStalePath(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	Link(home, repo, "", false)
	os.RemoveAll(repo)
	if _, err := Link(home, repo, "", true); err != nil {
		t.Fatalf("off on stale path: %v", err)
	}
	if c, _ := LoadConfig(home); len(c.Projects) != 0 {
		t.Error("entry kept")
	}
}

func TestLinkSymlinkSharesEntry(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	ln := filepath.Join(t.TempDir(), "ln")
	if err := os.Symlink(repo, ln); err != nil {
		t.Skip(err)
	}
	Link(home, repo, "", false)
	r, err := Link(home, ln, "", false)
	if err != nil || r.Mode != "global" {
		t.Fatalf("%+v %v", r, err)
	}
	if c, _ := LoadConfig(home); len(c.Projects) != 1 {
		t.Errorf("projects: %v", c.Projects)
	}
	c, _ := LoadConfig(home)
	if r := Resolve(ln, c); r.Mode != "global" {
		t.Errorf("resolve via symlink: %+v", r)
	}
}

func TestLinkNonNotExistStatErrorRefuses(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	// a path under a regular file gives ENOTDIR, not ErrNotExist
	file := filepath.Join(t.TempDir(), "f")
	os.WriteFile(file, nil, 0o600)
	c, _ := LoadConfig(home)
	c.Projects[filepath.Join(file, "sub")] = "repo"
	SaveConfig(home, c)
	if _, err := Link(home, repo, "", false); err == nil || !strings.Contains(err.Error(), "explicit name") {
		t.Fatalf("want refusal, got %v", err)
	}
}

func TestLinkSymlinkSpellings(t *testing.T) {
	home, real := t.TempDir(), t.TempDir()
	link := filepath.Join(t.TempDir(), "ln")
	if err := os.Symlink(real, link); err != nil {
		t.Skip("symlink unsupported")
	}
	// legacy entry saved under the symlinked spelling
	if err := os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects:\n  \""+link+"\": p\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Link(home, real, "p", false); err != nil {
		t.Fatalf("same repo, other spelling must replace: %v", err)
	}
	c, _ := LoadConfig(home)
	if len(c.Projects) != 1 || c.Projects[mustCanon(t, real)] != "p" {
		t.Fatalf("projects = %v", c.Projects)
	}
	if err := os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects:\n  \""+link+"\": p\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Link(home, link, "", true); err != nil {
		t.Fatal(err)
	}
	if c, _ := LoadConfig(home); len(c.Projects) != 0 {
		t.Fatalf("--off left %v", c.Projects)
	}
}

func mustCanon(t *testing.T, p string) string {
	t.Helper()
	c, err := canon(p)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestLinkConcurrent(t *testing.T) {
	home := t.TempDir()
	const n = 12
	repos := make([]string, n)
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := range repos {
		repos[i] = t.TempDir()
		wg.Add(1)
		go func(r string, i int) {
			defer wg.Done()
			_, err := Link(home, r, fmt.Sprintf("p%d", i), false)
			errs <- err
		}(repos[i], i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	c, _ := LoadConfig(home)
	if len(c.Projects) != n {
		t.Fatalf("want %d entries, got %d", n, len(c.Projects))
	}
	if _, err := os.Stat(filepath.Join(home, "config.lock")); err == nil {
		t.Fatal("lock left behind")
	}
}

func TestLockStaleAndTimeout(t *testing.T) {
	home := t.TempDir()
	path := filepath.Join(home, "config.lock")
	os.WriteFile(path, nil, 0o600)
	old := time.Now().Add(-time.Minute)
	os.Chtimes(path, old, old)
	unlock, err := lock(home)
	if err != nil {
		t.Fatalf("stale lock must be removed: %v", err)
	}
	unlock()
	os.WriteFile(path, nil, 0o600) // fresh, held
	defer func(w time.Duration) { lockWait = w }(lockWait)
	lockWait = 200 * time.Millisecond
	if _, err := lock(home); err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("want timeout, got %v", err)
	}
}

func TestLockOwnership(t *testing.T) {
	home := t.TempDir()
	defer func(s time.Duration) { lockStale = s }(lockStale)
	lockStale = 50 * time.Millisecond
	unlockA, err := lock(home)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(120 * time.Millisecond)
	unlockB, err := lock(home) // takes over A's stale lock
	if err != nil {
		t.Fatal(err)
	}
	unlockA() // must not remove B's lock
	if _, err := os.Stat(filepath.Join(home, "config.lock")); err != nil {
		t.Fatal("A's unlock removed B's lock")
	}
	unlockB()
	if _, err := os.Stat(filepath.Join(home, "config.lock")); err == nil {
		t.Fatal("B's unlock left the lock")
	}
}

func TestLinkRefusesCaseInsensitiveDuplicate(t *testing.T) {
	home := t.TempDir()
	a, b := t.TempDir(), t.TempDir()
	if _, err := Link(home, a, "App", false); err != nil {
		t.Fatal(err)
	}
	if _, err := Link(home, b, "app", false); err == nil {
		t.Fatal("want refusal for App/app")
	}
}

func TestLinkIgnoresRelativeKey(t *testing.T) {
	home := t.TempDir()
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects:\n  \"rel/x\": same\n"), 0o600)
	if _, err := Link(home, t.TempDir(), "same", false); err != nil {
		t.Fatalf("relative key must be skipped, got %v", err)
	}
}
