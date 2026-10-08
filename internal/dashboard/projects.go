// Package dashboard is the read-only core of `dh dashboard`: projects, progress, sessions, limits, sprites.
package dashboard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

const RootsEnv = "DH_DASHBOARD_ROOTS"

// Bounds on what one poll reads (security round 1): per-file bytes and entry counts.
const (
	MaxProjects       = 100
	MaxTicketsPerProj = 1000
	MaxSnapshots      = 500
)

// readCapped is the shared guarded read (regular file, at most MaxFileBytes).
func readCapped(path string) ([]byte, error) { return snapshot.ReadCapped(path) }

// skipped reports whether err means the file was refused (symlink/special/oversized), as opposed to absent or unreadable.
func skipped(err error) bool {
	return errors.Is(err, snapshot.ErrTooLarge) || errors.Is(err, snapshot.ErrNotRegular)
}

type Project struct {
	Name string
	Path string
}

// RootList splits the roots value (separator ';'): trimmed, absolute, empties dropped, env order kept.
func RootList(roots string) []string {
	out := []string{}
	for _, root := range strings.Split(roots, ";") {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		if abs, err := filepath.Abs(root); err == nil {
			root = abs // Clean is implied; symlinks and WSL/Windows path mapping are out of scope
		}
		out = append(out, root)
	}
	return out
}

// Projects lists direct subfolders of each root (separator ';') that contain .harness/.
// A missing/empty variable gives zero projects and a readable warning; a nonexistent root is skipped.
func Projects(roots string) ([]Project, string) {
	if strings.TrimSpace(roots) == "" {
		return nil, "DH_DASHBOARD_ROOTS is not set: no projects to show. Set it to folders separated by ';'."
	}
	var out []Project
	seen := map[string]bool{}
	readable := 0
	for _, root := range RootList(roots) {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		readable++
		for _, e := range entries {
			p := filepath.Join(root, e.Name())
			if seen[p] {
				continue
			}
			if st, err := os.Stat(filepath.Join(p, ".harness")); err == nil && st.IsDir() {
				seen[p] = true
				out = append(out, Project{Name: e.Name(), Path: p})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	if readable == 0 {
		return nil, "None of the folders in DH_DASHBOARD_ROOTS could be read: no projects to show."
	}
	if len(out) > MaxProjects {
		return out[:MaxProjects], fmt.Sprintf("More than %d projects found: showing the first %d.", MaxProjects, MaxProjects)
	}
	return out, ""
}

var (
	statusRow  = regexp.MustCompile(`(?m)^\|\s*Status\s*\|\s*([^|]+)\|`)
	statusBold = regexp.MustCompile(`\*\*Status:?\*\*:?\s*([^\n]+)`)
	prdField   = regexp.MustCompile(`(?m)^\|\s*(?:Story / PRD|PRD \(RF-[^)]*\))\s*\|\s*([^|]+)\|`)
	prdID      = regexp.MustCompile(`PRD-\d+`)
	prdFile    = regexp.MustCompile(`^(PRD-\d+)`)
)

func statusValue(md string) string {
	if m := statusRow.FindStringSubmatch(md); m != nil {
		return strings.TrimSpace(strings.ToLower(m[1]))
	}
	if m := statusBold.FindStringSubmatch(md); m != nil {
		return strings.TrimSpace(strings.ToLower(m[1]))
	}
	return ""
}

// TaskStatus mirrors taskStatus in adapters/claude-code/plugin/hooks/panel.ts: done, blocked or open.
// Shared fixture: testdata/status-cases.json.
func TaskStatus(md string) string {
	s := statusValue(md)
	switch {
	case strings.HasPrefix(s, "conclu"), strings.HasPrefix(s, "pronta, entregue"), strings.HasPrefix(s, "done"):
		return "done"
	case strings.HasPrefix(s, "bloque"):
		return "blocked"
	}
	return "open"
}

// PRDLink returns the PRD id a ticket links to, or "" ("fora"/"nenhum" or no PRD-n).
func PRDLink(md string) string {
	m := prdField.FindStringSubmatch(md)
	if m == nil {
		return ""
	}
	v := strings.TrimSpace(strings.ToLower(m[1]))
	if strings.HasPrefix(v, "fora") || strings.HasPrefix(v, "nenhum") {
		return ""
	}
	return prdID.FindString(m[1])
}

type Bar struct {
	PRD       string
	Done      int
	Blocked   int
	Total     int
	NoTickets bool // render "sem tickets", not 0%
}

type Progress struct {
	Warnings  []string
	Open      []Bar    // one per open PRD, sorted by id
	Delivered []string // PRD ids whose Status starts with "entregue em" (100%)
}

// ProjectProgress builds the bars for one project from docs/prd, .harness/prd and .harness/tasks.
func ProjectProgress(projectPath string) Progress {
	type counts struct{ done, blocked, total int }
	byPRD := map[string]*counts{}
	var warns []string
	name := filepath.Base(projectPath)
	tasks, _ := filepath.Glob(filepath.Join(projectPath, ".harness", "tasks", "*", "TASK.md"))
	if len(tasks) > MaxTicketsPerProj {
		warns = append(warns, fmt.Sprintf("%s: more than %d tickets: counting the first %d.", name, MaxTicketsPerProj, MaxTicketsPerProj))
		tasks = tasks[:MaxTicketsPerProj]
	}
	for _, f := range tasks {
		b, err := readCapped(f)
		if err != nil {
			if skipped(err) {
				warns = append(warns, fmt.Sprintf("%s: skipped %s (not a regular file or over 1 MiB).", name, filepath.Base(filepath.Dir(f))))
			}
			continue
		}
		md := string(b)
		id := PRDLink(md)
		if id == "" {
			continue
		}
		c := byPRD[id]
		if c == nil {
			c = &counts{}
			byPRD[id] = c
		}
		c.total++
		switch TaskStatus(md) {
		case "done":
			c.done++
		case "blocked":
			c.blocked++
		}
	}
	status := map[string]string{}
	for _, pat := range []string{"docs/prd/PRD-*.md", ".harness/prd/PRD-*.md"} {
		files, _ := filepath.Glob(filepath.Join(projectPath, filepath.FromSlash(pat)))
		files = slices.DeleteFunc(files, func(f string) bool { return strings.HasSuffix(f, ".review.md") }) // before the cap
		if len(files) > MaxTicketsPerProj {
			warns = append(warns, fmt.Sprintf("%s: more than %d PRD files: reading the first %d.", name, MaxTicketsPerProj, MaxTicketsPerProj))
			files = files[:MaxTicketsPerProj]
		}
		for _, f := range files {
			id := prdFile.FindString(filepath.Base(f))
			b, err := readCapped(f)
			if skipped(err) {
				warns = append(warns, fmt.Sprintf("%s: skipped %s (not a regular file or over 1 MiB).", name, filepath.Base(f)))
			}
			if id == "" || err != nil {
				continue
			}
			if _, seen := status[id]; !seen {
				status[id] = statusValue(string(b))
			}
		}
	}
	p := Progress{Warnings: warns}
	for id, st := range status {
		if strings.HasPrefix(st, "entregue em") {
			p.Delivered = append(p.Delivered, id)
			continue
		}
		bar := Bar{PRD: id, NoTickets: true}
		if c := byPRD[id]; c != nil {
			bar = Bar{PRD: id, Done: c.done, Blocked: c.blocked, Total: c.total}
		}
		p.Open = append(p.Open, bar)
	}
	sort.Strings(p.Delivered)
	sort.Slice(p.Open, func(i, j int) bool { return p.Open[i].PRD < p.Open[j].PRD })
	return p
}

// PRDStatus returns the lowercased header Status of a PRD (or ticket) markdown, or "".
func PRDStatus(md string) string { return statusValue(md) }
