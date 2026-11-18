package harness

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var headSHA = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

// HeadSHA returns the short sha of HEAD in the git work tree that contains dir, and false when there is none: git is not installed, dir is not inside a repository,
// or the repository has no commit yet. `dh facts check` prints it so that the Coordinator gets the sha from a `dh` command it already runs and does not call git
// itself (PRD-035 R15, verification of F3). It only reads, and it gives up after 5 seconds.
func HeadSHA(dir string) (string, bool) {
	git, err := exec.LookPath("git")
	if err != nil {
		return "", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, git, "-C", dir, "rev-parse", "--short", "HEAD")
	cmd.Env = append(cmd.Environ(), "LC_ALL=C", "GIT_OPTIONAL_LOCKS=0")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	sha := strings.TrimSpace(string(out))
	if !headSHA.MatchString(sha) {
		return "", false
	}
	return sha, true
}
