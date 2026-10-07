package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ErrPartial means the file is empty or truncated JSON (a writer is mid-write); callers skip it for this cycle.
var ErrPartial = errors.New("partial snapshot")

// ErrTooLarge means the file is over MaxFileBytes and is not read (a snapshot is a few KB).
var ErrTooLarge = errors.New("snapshot too large")

// ErrNotRegular means the path is a symlink, FIFO, device or directory: never followed or read.
var ErrNotRegular = errors.New("not a regular file")

// MaxFileBytes bounds what a poll reads from any one snapshot file.
const MaxFileBytes = 1 << 20

// Activity values of schema 2.
const (
	ActivityIdle    = "idle"
	ActivityWorking = "working"
	ActivityWaiting = "waiting"
	ActivityError   = "error"
	ActivityDone    = "done"
)

// DefaultDoneDecay is how long `done` stays visible before reading as `idle`.
const DefaultDoneDecay = 5 * time.Minute

// Limits are account rate-limit percentages; nil means the snapshot had no value.
type Limits struct {
	FiveHour *float64
	SevenDay *float64
}

// Has reports whether at least one limit value is present.
func (l Limits) Has() bool { return l.FiveHour != nil || l.SevenDay != nil }

// Session is the tolerant reader's view of a schema 1 or 2 snapshot.
type Session struct {
	Schema      int
	SessionID   string
	SessionName string
	CWD         string
	Agent       string
	State       string // schema 1 vocabulary: active, idle, closed
	Activity    string // schema 2, or derived from State for schema 1
	ActivityAt  time.Time
	UpdatedAt   time.Time
	Limits      Limits
}

type rawSession struct {
	Schema      int             `json:"schema"`
	SessionID   string          `json:"session_id"`
	SessionName string          `json:"session_name"`
	CWD         string          `json:"cwd"`
	Agent       string          `json:"agent"`
	State       string          `json:"state"`
	Activity    string          `json:"activity"`
	ActivityAt  string          `json:"activity_at"`
	UpdatedAt   string          `json:"updated_at"`
	RateLimits  json.RawMessage `json:"rate_limits"`
}

// Load reads one snapshot file. Truncated or empty JSON returns ErrPartial.
// Unknown fields are ignored (the file also carries tasks, events, ...).
func Load(path string) (Session, error) {
	b, err := ReadCapped(path)
	if err != nil {
		return Session{}, err
	}
	var r rawSession
	if err := json.Unmarshal(b, &r); err != nil {
		var syn *json.SyntaxError
		if len(strings.TrimSpace(string(b))) == 0 || errors.As(err, &syn) || strings.Contains(err.Error(), "unexpected end") {
			return Session{}, fmt.Errorf("%s: %w", filepath.Base(path), ErrPartial)
		}
		return Session{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	s := Session{Schema: r.Schema, SessionID: r.SessionID, SessionName: r.SessionName, CWD: r.CWD,
		Agent: r.Agent, State: r.State, Activity: r.Activity, Limits: parseLimits(r.RateLimits)}
	s.ActivityAt, _ = time.Parse(time.RFC3339, r.ActivityAt)
	s.UpdatedAt, _ = time.Parse(time.RFC3339, r.UpdatedAt)
	switch s.Activity {
	case ActivityIdle, ActivityWorking, ActivityWaiting, ActivityError, ActivityDone:
	default: // schema 1, or an unknown value: derive from state
		s.Activity = ActivityIdle
		if s.State == "active" {
			s.Activity = ActivityWorking
		}
	}
	return s, nil
}

// LoadDir loads every *.json in dir, skipping partial and unreadable files. A missing dir yields none.
func LoadDir(dir string) []Session {
	out, _ := LoadDirMax(dir, 0)
	return out
}

// LoadDirMax is LoadDir keeping only the max newest files by mtime (max <= 0 = no cap); skipped counts the files left out.
func LoadDirMax(dir string, max int) (out []Session, skipped int) {
	files, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if max > 0 && len(files) > max {
		mt := map[string]time.Time{}
		for _, f := range files {
			if st, err := os.Stat(f); err == nil {
				mt[f] = st.ModTime()
			}
		}
		sort.Slice(files, func(i, j int) bool { return mt[files[i]].After(mt[files[j]]) })
		skipped, files = len(files)-max, files[:max]
	}
	for _, f := range files {
		if s, err := Load(f); err == nil {
			out = append(out, s)
		}
	}
	return out, skipped
}

// EffectiveActivity decays `done` to `idle` once now - ActivityAt >= decay (a missing/unparseable ActivityAt counts as decayed; a future one clamps to now).
func (s Session) EffectiveActivity(now time.Time, decay time.Duration) string {
	if s.Activity != ActivityDone {
		return s.Activity
	}
	if s.ActivityAt.IsZero() { // missing or unparseable activity_at: cannot prove it is fresh
		return ActivityIdle
	}
	at := s.ActivityAt
	if at.After(now) {
		at = now
	}
	if now.Sub(at) >= decay {
		return ActivityIdle
	}
	return ActivityDone
}

// parseLimits accepts {"five_hour": 42} , {"five_hour": {"used_percentage"|"percent_used"|"percentUsed": 42}}
// and [{"kind":"five_hour","percentUsed":42}] (the shape of $.session.usage().rateLimits).
func parseLimits(raw json.RawMessage) Limits {
	var l Limits
	if len(raw) == 0 {
		return l
	}
	set := func(kind string, v any) {
		f := pct(v)
		if f == nil {
			return
		}
		switch kind {
		case "five_hour":
			l.FiveHour = f
		case "seven_day":
			l.SevenDay = f
		}
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) == nil {
		for k, v := range obj {
			set(k, v)
		}
		return l
	}
	var arr []map[string]any
	if json.Unmarshal(raw, &arr) == nil {
		for _, e := range arr {
			k, _ := e["kind"].(string)
			set(k, e)
		}
	}
	return l
}

func pct(v any) *float64 {
	switch t := v.(type) {
	case float64:
		return &t
	case map[string]any:
		for _, k := range []string{"used_percentage", "percent_used", "percentUsed"} {
			if f, ok := t[k].(float64); ok {
				return &f
			}
		}
	}
	return nil
}
