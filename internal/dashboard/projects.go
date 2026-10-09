// Package dashboard is the read-only core of `dh dashboard`: projects, progress, sessions, limits, sprites.
package dashboard

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// Bounds on what one poll reads (security round 1): per-file bytes and entry counts.
const (
	MaxTicketsPerProj = 1000
	MaxSnapshots      = 500
)

// readCapped is the shared guarded read (regular file, at most MaxFileBytes).
func readCapped(path string) ([]byte, error) { return snapshot.ReadCapped(path) }

// readCappedIn is readCapped for dir/name where dir may be swapped for a symlink: the open is relative to dir.
func readCappedIn(dir, name string) ([]byte, error) {
	return snapshot.ReadCappedIn(dir, name, snapshot.MaxFileBytes)
}

// skipped reports whether err means the file was refused (symlink/special/oversized), as opposed to absent or unreadable.
func skipped(err error) bool {
	return errors.Is(err, snapshot.ErrTooLarge) || errors.Is(err, snapshot.ErrNotRegular)
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
	Delivered []string // PRD ids whose Status starts with "entregue em" or "delivered on" (100%)
	Metrics   Metrics  // PRDs and PRD averages; buildState adds the file metrics
}

// ProjectProgress builds the bars for one project: tickets from <harnessDir>/tasks, PRDs from
// <repoPath>/docs/prd and <harnessDir>/prd. harnessDir is the resolved harness folder (repo or global mode).
func ProjectProgress(repoPath, harnessDir string) Progress {
	return projectProgress(filepath.Base(repoPath), repoPath, harnessDir)
}

// projectProgress is ProjectProgress with the registered project name used in warnings.
func projectProgress(name, repoPath, harnessDir string) Progress {
	byPRD := map[string]*counts{}
	var warns []string
	tasks, _ := filepath.Glob(filepath.Join(harnessDir, "tasks", "*", "TASK.md"))
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
			if w := c.addDone(md, name, filepath.Base(filepath.Dir(f))); w != "" {
				warns = append(warns, w)
			}
		case "blocked":
			c.blocked++
		}
	}
	status := map[string]string{}
	dates := map[string][2]string{}
	for _, pat := range []string{filepath.Join(repoPath, "docs", "prd", "PRD-*.md"), filepath.Join(harnessDir, "prd", "PRD-*.md")} {
		files, _ := filepath.Glob(pat)
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
				st, dl := prdDates(string(b))
				dates[id] = [2]string{st, dl}
			}
		}
	}
	p := Progress{Warnings: warns}
	var mw []string
	p.Metrics, mw = prdMetrics(status, dates, byPRD, name)
	p.Warnings = append(p.Warnings, mw...)
	for id, st := range status {
		if isDelivered(st) {
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
