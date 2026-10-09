package harness

import (
	"os"
	"path/filepath"
	"strings"
)

// Resolved is the answer of the single resolver. Mode is repo, global or none; Dir is empty in none.
type Resolved struct {
	Mode string `json:"mode"`
	Dir  string `json:"dir"`
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

// validName accepts a single path segment: no separators, not "." or "..", no control chars.
func validName(n string) bool {
	if n == "" || n == "." || n == ".." || len(n) > 255 || strings.ContainsAny(n, `/\:`) ||
		strings.HasSuffix(n, ".") || strings.HasSuffix(n, " ") {
		return false
	}
	for _, r := range n {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func resolveIn(home, repo string, c Config) Resolved {
	if cr, err := canon(repo); err == nil {
		repo = cr
	}
	if d := filepath.Join(repo, ".harness"); isDir(d) {
		return Resolved{"repo", d}
	}
	if name, ok := c.Projects[repo]; ok && validName(name) && home != "" {
		if d := filepath.Join(home, "projects", name); IsHarness(d) {
			return Resolved{"global", d}
		}
	}
	return Resolved{"none", ""}
}

// Resolve: (1) <repo>/.harness if it is a directory; (2) <home>/projects/<name> through the
// registry, if it exists; (3) none. repo must be absolute and clean.
func Resolve(repo string, c Config) Resolved { return resolveIn(Home(), repo, c) }
