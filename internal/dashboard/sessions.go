package dashboard

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// DefaultStale is the liveness limit (T-1301: heartbeat every 30 s, H3 true).
const DefaultStale = 2 * time.Minute

type OpenSession struct {
	snapshot.Session
	Project  string // project path, "" when no project contains the cwd
	Activity string // effective activity (done decayed)
}

// OpenSessions drops closed and stale sessions and attributes the rest to the deepest project containing the cwd.
func OpenSessions(all []snapshot.Session, projects []Project, now time.Time, stale, decay time.Duration) []OpenSession {
	var out []OpenSession
	for _, s := range all {
		if s.State == "closed" || s.UpdatedAt.IsZero() || now.Sub(s.UpdatedAt) > stale {
			continue
		}
		o := OpenSession{Session: s, Activity: s.EffectiveActivity(now, decay)}
		for _, p := range projects {
			if within(p.Path, s.CWD) && len(p.Path) > len(o.Project) {
				o.Project = p.Path
			}
		}
		out = append(out, o)
	}
	return out
}

func within(dir, cwd string) bool {
	if cwd == "" {
		return false
	}
	rel, err := filepath.Rel(dir, cwd)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// AccountLimits merges five_hour and seven_day per field from the newest snapshot that has each, with per-field age.
type AccountLimits struct {
	FiveHour, SevenDay       *float64
	FiveHourAge, SevenDayAge time.Duration
	OK                       bool // false = "sem dado"
}

// LimitFreshness: older limit values still show (with age) but no longer drive the attention state.
const LimitFreshness = 5 * time.Hour

// LatestLimits takes each limit from the newest snapshot (any state; limits belong to the account) that has it.
func LatestLimits(all []snapshot.Session, now time.Time) AccountLimits {
	var out AccountLimits
	var fiveAt, sevenAt time.Time
	for _, s := range all {
		if s.Limits.FiveHour != nil && (out.FiveHour == nil || s.UpdatedAt.After(fiveAt)) {
			out.FiveHour, fiveAt = s.Limits.FiveHour, s.UpdatedAt
		}
		if s.Limits.SevenDay != nil && (out.SevenDay == nil || s.UpdatedAt.After(sevenAt)) {
			out.SevenDay, sevenAt = s.Limits.SevenDay, s.UpdatedAt
		}
	}
	out.OK = out.FiveHour != nil || out.SevenDay != nil
	out.FiveHourAge, out.SevenDayAge = now.Sub(fiveAt), now.Sub(sevenAt)
	return out
}

// Avatar states (jevmon ids). Priority: esperando > erro > trabalhando > concluido > atencao > parado.
const (
	StateWaiting = "esperando"
	StateError   = "erro"
	StateWorking = "trabalhando"
	StateDone    = "concluido"
	StateAttn    = "atencao"
	StateIdle    = "parado"
)

var statePriority = map[string]int{StateWaiting: 1, StateError: 2, StateWorking: 3, StateDone: 4, StateAttn: 5, StateIdle: 6}

var activityState = map[string]string{
	snapshot.ActivityWaiting: StateWaiting, snapshot.ActivityError: StateError,
	snapshot.ActivityWorking: StateWorking, snapshot.ActivityDone: StateDone, snapshot.ActivityIdle: StateIdle,
}

// LimitAlert is the panel.ts cutoff for attention.
const LimitAlert = 80.0

// AvatarState picks the highest-priority state among open sessions and the limit alert.
func AvatarState(sessions []OpenSession, lim AccountLimits) string {
	best := StateIdle
	pick := func(s string) {
		if prio(s) < prio(best) {
			best = s
		}
	}
	for _, s := range sessions {
		st, ok := activityState[s.Activity]
		if !ok {
			st = StateIdle // unknown/empty activity never outranks a real state
		}
		pick(st)
	}
	hot := func(v *float64, age time.Duration) bool { return v != nil && *v >= LimitAlert && age < LimitFreshness }
	if hot(lim.FiveHour, lim.FiveHourAge) || hot(lim.SevenDay, lim.SevenDayAge) {
		pick(StateAttn)
	}
	return best
}

// HigherPriority reports whether state a outranks b.
func HigherPriority(a, b string) bool { return prio(a) < prio(b) }

func prio(s string) int {
	if p, ok := statePriority[s]; ok {
		return p
	}
	return statePriority[StateIdle]
}
