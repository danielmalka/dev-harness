package dashboard

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// Metrics is the additive per-project block of /api/state, derived on read from the project files.
type Metrics struct {
	PRDs       []PRDMetric `json:"prds"`
	PRDDaysAvg *float64    `json:"prd_days_avg"`
	PRDDaysN   int         `json:"prd_days_n"`
	Incidents  *int        `json:"incidents"`
	Memory     MemoryStat  `json:"memory"`
	Epochal    EpochalStat `json:"epochal"`
}

type PRDMetric struct {
	ID            string   `json:"id"`
	Status        string   `json:"status"` // "open" | "delivered"
	Tickets       int      `json:"tickets"`
	Done          int      `json:"done"`
	Started       *string  `json:"started"`
	Delivered     *string  `json:"delivered"`
	Days          *int     `json:"days"`
	TicketDaysAvg *float64 `json:"ticket_days_avg"`
	TicketDaysN   int      `json:"ticket_days_n"`
	TicketUndated int      `json:"ticket_undated"`
}

type MemoryStat struct {
	Lines   *int `json:"lines"`
	Bytes   *int `json:"bytes"`
	Missing bool `json:"missing"`
}

type EpochalStat struct {
	Batches *int    `json:"batches"`
	LastAt  *string `json:"last_at"`
}

// counts is the per-PRD ticket tally ProjectProgress fills in its single pass over the tickets.
type counts struct{ done, blocked, total, daysSum, daysN, undated int }

var (
	createdRow   = regexp.MustCompile(`(?mi)^\|\s*(?:Criado / atualizado|Created / updated)\s*\|\s*([^|]*)\|`)
	doneOnRow    = regexp.MustCompile(`(?mi)^\|\s*(?:Concluído em|Done on)\s*\|\s*([^|]*)\|`)
	dateOnly     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	isoDate      = regexp.MustCompile(`\d{4}-\d{2}-\d{2}(?:[T ]\d{2}:\d{2}(?::\d{2})?(?:Z|[+-]\d{2}:\d{2})?)?`)
	incidentLine = regexp.MustCompile(`^(?:###\s+|[-*+]\s+)\W*[A-Z]+-\d+`)
	htmlComment  = regexp.MustCompile(`(?s)<!--.*?-->`)
)

// parseDay returns s as a day when it is exactly a valid AAAA-MM-DD (placeholders and other text are absent).
func parseDay(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if !dateOnly.MatchString(s) {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", s)
	return t, err == nil
}

// firstDay reads the first date of the "Criado / atualizado" row ("a / b" -> a).
func firstDay(md string) (time.Time, bool) {
	m := createdRow.FindStringSubmatch(md)
	if m == nil {
		return time.Time{}, false
	}
	return parseDay(strings.SplitN(m[1], "/", 2)[0])
}

func dayCount(from, to time.Time) int { return int(to.Sub(from).Hours() / 24) }

// isDelivered is the PRD delivered rule (lowercased Status), in both template languages.
func isDelivered(st string) bool {
	return strings.HasPrefix(st, "entregue em") || strings.HasPrefix(st, "delivered on")
}

// addDone tallies one finished ticket; it returns a warning when the end date precedes the start.
func (c *counts) addDone(md, project, ticket string) string {
	start, ok1 := firstDay(md)
	end, ok2 := time.Time{}, false
	if m := doneOnRow.FindStringSubmatch(md); m != nil {
		end, ok2 = parseDay(m[1])
	}
	switch {
	case !ok1 || !ok2:
		c.undated++
	case end.Before(start):
		c.undated++
		return fmt.Sprintf("%s: ticket %s: Done on is before Created: counted as undated.", project, ticket)
	default:
		c.daysSum += dayCount(start, end)
		c.daysN++
	}
	return ""
}

// prdDates returns the Created day and, for a delivered Status, the delivery day ("" when absent).
func prdDates(md string) (started, delivered string) {
	if t, ok := firstDay(md); ok {
		started = t.Format("2006-01-02")
	}
	st := statusValue(md)
	for _, p := range []string{"entregue em", "delivered on"} {
		if rest, ok := strings.CutPrefix(st, p); ok {
			if f := strings.Fields(rest); len(f) > 0 {
				if t, ok := parseDay(f[0]); ok {
					delivered = t.Format("2006-01-02")
				}
			}
		}
	}
	return
}

func avg1(sum, n int) *float64 {
	if n == 0 {
		return nil
	}
	v := math.Round(float64(sum)/float64(n)*10) / 10
	return &v
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// prdMetrics builds Metrics.PRDs and the PRD-level average from the single pass of ProjectProgress.
func prdMetrics(status map[string]string, dates map[string][2]string, byPRD map[string]*counts, project string) (Metrics, []string) {
	var warns []string
	m := Metrics{PRDs: []PRDMetric{}}
	ids := make([]string, 0, len(status))
	for id := range status {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	sum := 0
	for _, id := range ids {
		pm := PRDMetric{ID: id, Status: "open", Started: strPtr(dates[id][0])}
		if c := byPRD[id]; c != nil {
			pm.Tickets, pm.Done = c.total, c.done
			pm.TicketDaysN, pm.TicketUndated = c.daysN, c.undated
			pm.TicketDaysAvg = avg1(c.daysSum, c.daysN)
		}
		if isDelivered(status[id]) {
			pm.Status = "delivered"
			pm.Delivered = strPtr(dates[id][1])
			s, e := parseISO(dates[id][0]), parseISO(dates[id][1])
			if pm.Started != nil && pm.Delivered != nil {
				if e.Before(s) {
					warns = append(warns, fmt.Sprintf("%s: %s: delivery date is before Created: no duration.", project, id))
				} else {
					d := dayCount(s, e)
					pm.Days = &d
					sum += d
					m.PRDDaysN++
				}
			}
		}
		m.PRDs = append(m.PRDs, pm)
	}
	m.PRDDaysAvg = avg1(sum, m.PRDDaysN)
	return m, warns
}

func parseISO(s string) time.Time { t, _ := time.Parse("2006-01-02", s); return t }

// fileMetrics fills incidents, memory and epochal from the resolved harness folder; it returns warnings.
func fileMetrics(m *Metrics, project, harnessDir string) []string {
	var warns []string
	if symlinked(harnessDir) {
		m.Memory = MemoryStat{}
		m.Epochal = EpochalStat{}
		return []string{fmt.Sprintf("%s: harness folder is a symlink: metrics skipped.", project)}
	}
	refuse := func(file string) {
		warns = append(warns, fmt.Sprintf("%s: %s unreadable (not a regular file or over 1 MiB).", project, file))
	}
	if b, err := readCappedIn(harnessDir, "RISKS.md"); err == nil {
		n := countIncidents(string(b))
		m.Incidents = &n
	} else if !errors.Is(err, fs.ErrNotExist) {
		refuse("RISKS.md")
	}
	if b, err := readCappedIn(harnessDir, "MEMORY.md"); err == nil {
		l, n := bytes.Count(b, []byte("\n")), len(b)
		if n > 0 && b[n-1] != '\n' {
			l++ // last line without a final newline
		}
		m.Memory = MemoryStat{Lines: &l, Bytes: &n}
	} else if errors.Is(err, fs.ErrNotExist) {
		m.Memory.Missing = true
	} else {
		refuse("MEMORY.md")
	}
	b0 := 0
	m.Epochal = EpochalStat{Batches: &b0}
	if e, err := scanEpochal(filepath.Join(harnessDir, "EPOCHAL.md")); err == nil {
		m.Epochal = e
	} else if !errors.Is(err, fs.ErrNotExist) {
		m.Epochal = EpochalStat{}
		warns = append(warns, fmt.Sprintf("%s: EPOCHAL.md unreadable.", project))
	}
	return warns
}

// symlinked reports whether the resolved harness folder itself is a symlink (refused like the files in it).
func symlinked(dir string) bool {
	fi, err := os.Lstat(dir)
	return err == nil && fi.Mode()&fs.ModeSymlink != 0
}

// epochalCache keeps the last scan per path while size and mtime are unchanged, so each state build does not
// re-read a large append-only file. ponytail: one entry per project path, never evicted (a handful of projects).
var (
	epochalMu    sync.Mutex
	epochalCache = map[string]epochalEntry{}
	epochalScans int // test counter: full scans performed
)

type epochalEntry struct {
	size int64
	mod  time.Time
	stat EpochalStat
}

func scanEpochal(path string) (EpochalStat, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return EpochalStat{}, err
	}
	if !fi.Mode().IsRegular() {
		return EpochalStat{}, snapshot.ErrNotRegular
	}
	epochalMu.Lock()
	e, ok := epochalCache[path]
	epochalMu.Unlock()
	if ok && e.size == fi.Size() && e.mod.Equal(fi.ModTime()) {
		return e.stat, nil
	}
	st, err := scanEpochalFile(path)
	if err != nil {
		return st, err
	}
	epochalMu.Lock()
	epochalScans++
	epochalCache[path] = epochalEntry{fi.Size(), fi.ModTime(), st}
	epochalMu.Unlock()
	return st, nil
}

// countIncidents counts ### titles and list items starting with an id, inside the Incidents section.
func countIncidents(md string) int {
	n, in := 0, false
	for _, line := range strings.Split(htmlComment.ReplaceAllString(md, ""), "\n") {
		line = strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(line, "## ") {
			t := strings.TrimSpace(line[3:])
			in = t == "Incidentes" || t == "Incidents"
			continue
		}
		if in && incidentLine.MatchString(line) {
			n++
		}
	}
	return n
}

// scanEpochal streams EPOCHAL.md line by line (no size cap, constant memory): batches are ### titles under
// the archived-batches section outside the raw copies; last_at is the first ISO date after the last title.
// ponytail: only the first 4 KiB of each line is inspected; titles and dates sit at the start of their lines.
func scanEpochalFile(path string) (EpochalStat, error) {
	f, err := snapshot.OpenRegularIn(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return EpochalStat{}, err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 4096)
	batches, inSection, inRaw, want := 0, false, false, false
	var last *string
	for {
		chunk, isPrefix, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return EpochalStat{}, err
			}
			break
		}
		line := string(chunk)
		for isPrefix { // drop the rest of an over-long line
			_, isPrefix, err = r.ReadLine()
			if err != nil && err != io.EOF {
				return EpochalStat{}, err
			}
		}
		switch {
		case hasAnyPrefix(line, "<!-- INICIO-BRUTO", "<!-- BEGIN-RAW"):
			inRaw, want = true, false
		case hasAnyPrefix(line, "<!-- FIM-BRUTO", "<!-- END-RAW"):
			inRaw = false
		case inRaw:
		case strings.HasPrefix(line, "## "):
			t := strings.TrimSpace(line[3:])
			inSection, want = t == "Lotes arquivados" || t == "Archived batches", false
		case inSection && strings.HasPrefix(line, "### "):
			batches++
			last, want = nil, true
		case want:
			if d := isoDate.FindString(line); d != "" {
				last, want = &d, false
			}
		}
	}
	return EpochalStat{Batches: &batches, LastAt: last}, nil
}

func hasAnyPrefix(s string, ps ...string) bool {
	for _, p := range ps {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
