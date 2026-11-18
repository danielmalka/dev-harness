package harness

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false"}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("git is not usable here: %v %s", err, out)
	}
	return strings.TrimSpace(string(out))
}

// HeadSHA gives the short sha git itself reports, from the root and from a subfolder, and says "no" for a folder outside any repository, a repository with no
// commit, a folder that does not exist and a machine without git.
func TestHeadSHA(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	repo := t.TempDir()
	runGit(t, repo, "init", "-q")
	if sha, ok := HeadSHA(repo); ok {
		t.Fatalf("a repository with no commit has no HEAD, got %q", sha)
	}
	if err := os.MkdirAll(filepath.Join(repo, "sub", "deep"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "sub", "deep", "f.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "-A")
	runGit(t, repo, "commit", "-q", "-m", "x")
	want := runGit(t, repo, "rev-parse", "--short", "HEAD")
	for _, dir := range []string{repo, filepath.Join(repo, "sub", "deep")} {
		if sha, ok := HeadSHA(dir); !ok || sha != want {
			t.Errorf("HeadSHA(%s) = %q, %v; want %q", dir, sha, ok, want)
		}
	}
	if sha, ok := HeadSHA(filepath.Join(repo, "nao-existe")); ok {
		t.Errorf("a missing folder has no sha, got %q", sha)
	}
	if sha, ok := HeadSHA(t.TempDir()); ok {
		// a temporary folder may sit inside some other work tree on an odd machine: only a hex sha is acceptable then
		if !headSHA.MatchString(sha) {
			t.Errorf("not a sha: %q", sha)
		}
	}
	t.Setenv("PATH", "")
	if sha, ok := HeadSHA(repo); ok {
		t.Errorf("without git on PATH there is no sha, got %q", sha)
	}
}
