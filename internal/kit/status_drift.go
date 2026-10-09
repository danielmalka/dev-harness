package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/danielmalka/dev-harness/internal/dashboard"
	"github.com/danielmalka/dev-harness/internal/harness"
)

var driftPRDFile = regexp.MustCompile(`^(PRD-\d+)`)

// checkStatusDrift (PRD-012, R19) fails when tickets and their PRD disagree:
// (a) an open ticket links to a PRD whose Status starts with "entregue em";
// (b) every linked ticket (at least one) is done but the PRD is not "entregue em".
// Classification and linking come from internal/dashboard. Tickets and extra PRDs are read from
// the resolved harness folder (<repo>/.harness or <home>/projects/<name>, PRD-014 R19); silent
// when it has no tasks (CI never has them).
func checkStatusDrift(root string, errors *[]string) {
	cfg, _ := harness.LoadConfig(harness.Home()) // a broken config only means "not registered"
	dir := harness.Resolve(root, cfg).Dir
	if dir == "" {
		return
	}
	tasks, _ := filepath.Glob(filepath.Join(dir, "tasks", "*", "TASK.md"))
	if len(tasks) == 0 {
		return
	}
	status := map[string]string{}
	for _, pat := range []string{filepath.Join(root, "docs", "prd", "PRD-*.md"), filepath.Join(dir, "prd", "PRD-*.md")} {
		files, _ := filepath.Glob(pat)
		for _, f := range files {
			base := filepath.Base(f)
			id := driftPRDFile.FindString(base)
			if id == "" || strings.HasSuffix(base, ".review.md") {
				continue
			}
			if _, seen := status[id]; seen {
				continue
			}
			if b, err := os.ReadFile(f); err == nil {
				status[id] = dashboard.PRDStatus(string(b))
			}
		}
	}
	type counts struct {
		total, done int
		open        []string
	}
	byPRD := map[string]*counts{}
	for _, f := range tasks {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		id := dashboard.PRDLink(string(b))
		if _, known := status[id]; id == "" || !known {
			continue
		}
		c := byPRD[id]
		if c == nil {
			c = &counts{}
			byPRD[id] = c
		}
		c.total++
		if dashboard.TaskStatus(string(b)) == "done" {
			c.done++
		} else {
			c.open = append(c.open, filepath.Base(filepath.Dir(f)))
		}
	}
	ids := make([]string, 0, len(byPRD))
	for id := range byPRD {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		c, delivered := byPRD[id], strings.HasPrefix(status[id], "entregue em")
		sort.Strings(c.open)
		if delivered && len(c.open) > 0 {
			*errors = append(*errors, fmt.Sprintf("%s: expected no open ticket for a PRD marked \"entregue em\", found %s open", id, strings.Join(c.open, ", ")))
		}
		if !delivered && c.total > 0 && c.done == c.total {
			*errors = append(*errors, fmt.Sprintf("%s: expected Status starting with \"entregue em\" because all %d linked tickets are done, found %q", id, c.total, status[id]))
		}
	}
}
