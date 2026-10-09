package harness

import (
	"fmt"
	"path/filepath"
	"sort"
)

// MaxProjects caps the list `dh projects` and the dashboard show.
const MaxProjects = 100

// Project is one registered repository with a resolvable harness.
type Project struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	Harness string `json:"harness"`
}

// Projects lists registered projects whose path is a directory and whose resolution is not
// none, ordered by path, capped at MaxProjects. Warnings are returned, never printed here.
func Projects(home string) ([]Project, []string) {
	c, err := LoadConfig(home)
	if err != nil {
		return nil, []string{err.Error()}
	}
	paths := make([]string, 0, len(c.Projects))
	bad := 0
	for p := range c.Projects {
		if !filepath.IsAbs(p) || p != filepath.Clean(p) { // hand-edited key: never resolve against the cwd
			bad++
			continue
		}
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var out []Project
	var warns []string
	skipped := 0
	for _, p := range paths {
		r := resolveIn(home, p, c)
		if !isDir(p) || r.Mode == "none" {
			skipped++
			continue
		}
		out = append(out, Project{Name: c.Projects[p], Path: p, Mode: r.Mode, Harness: r.Dir})
	}
	if bad > 0 {
		warns = append(warns, fmt.Sprintf("%d registry key(s) ignored: not an absolute, clean path", bad))
	}
	if skipped > 0 {
		warns = append(warns, fmt.Sprintf("%d registered project(s) skipped: path missing or no harness (run dh link)", skipped))
	}
	if len(out) > MaxProjects {
		warns = append(warns, fmt.Sprintf("%d projects registered: showing the first %d", len(out), MaxProjects))
		out = out[:MaxProjects]
	}
	return out, warns
}
