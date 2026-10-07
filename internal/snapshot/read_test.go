package snapshot

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadPartial(t *testing.T) {
	d := t.TempDir()
	for name, body := range map[string]string{"a.json": `{"schema":2,"session_id":"a","sta`, "b.json": ``} {
		if _, err := Load(write(t, d, name, body)); !errors.Is(err, ErrPartial) {
			t.Fatalf("%s: want ErrPartial, got %v", name, err)
		}
	}
	if len(LoadDir(d)) != 0 {
		t.Fatal("partial files must be skipped")
	}
}

func TestLoadSchema1DerivesActivity(t *testing.T) {
	d := t.TempDir()
	cases := map[string]string{"active": ActivityWorking, "idle": ActivityIdle, "closed": ActivityIdle}
	for state, want := range cases {
		s, err := Load(write(t, d, state+".json", `{"schema":1,"session_id":"s","state":"`+state+`","tasks":[{"id":"x"}],"extra":1}`))
		if err != nil || s.Activity != want {
			t.Fatalf("%s: got %q err %v, want %q", state, s.Activity, err, want)
		}
	}
}

func TestLoadSchema2(t *testing.T) {
	d := t.TempDir()
	s, err := Load(write(t, d, "s.json", `{"schema":2,"session_id":"s","cwd":"/p","state":"active","activity":"waiting","activity_at":"2026-10-07T10:00:00Z","updated_at":"2026-10-07T10:01:00Z","rate_limits":[{"kind":"five_hour","percentUsed":42},{"kind":"seven_day","percentUsed":81.5}]}`))
	if err != nil || s.Activity != ActivityWaiting || s.CWD != "/p" || s.ActivityAt.IsZero() || s.UpdatedAt.IsZero() {
		t.Fatalf("bad session %+v err %v", s, err)
	}
	if *s.Limits.FiveHour != 42 || *s.Limits.SevenDay != 81.5 {
		t.Fatalf("bad limits %+v", s.Limits)
	}
	s, _ = Load(write(t, d, "o.json", `{"schema":2,"session_id":"o","activity":"idle","rate_limits":{"five_hour":{"used_percentage":10},"seven_day":20}}`))
	if *s.Limits.FiveHour != 10 || *s.Limits.SevenDay != 20 {
		t.Fatalf("object-shaped limits: %+v", s.Limits)
	}
	s, _ = Load(write(t, d, "n.json", `{"schema":2,"session_id":"n","activity":"bogus","state":"active"}`))
	if s.Activity != ActivityWorking || s.Limits.Has() {
		t.Fatalf("unknown activity must derive from state: %+v", s)
	}
}

func TestDoneDecay(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	s := Session{Activity: ActivityDone, ActivityAt: now.Add(-6 * time.Minute), UpdatedAt: now.Add(-5 * time.Second)}
	if got := s.EffectiveActivity(now, DefaultDoneDecay); got != ActivityIdle {
		t.Fatalf("6 min old done must read idle, got %s", got)
	}
	s.ActivityAt = now.Add(-4 * time.Minute)
	if got := s.EffectiveActivity(now, DefaultDoneDecay); got != ActivityDone {
		t.Fatalf("4 min old done must stay done, got %s", got)
	}
}

func TestDoneDecayMissingAndFuture(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	if (Session{Activity: ActivityDone, UpdatedAt: now}).EffectiveActivity(now, DefaultDoneDecay) != ActivityIdle {
		t.Fatal("missing activity_at must read idle")
	}
	if (Session{Activity: ActivityDone, ActivityAt: now.Add(time.Hour)}).EffectiveActivity(now, DefaultDoneDecay) != ActivityDone {
		t.Fatal("future activity_at clamps to now: still done")
	}
}

func TestWriterPreservesUnknownFields(t *testing.T) {
	d := t.TempDir()
	write(t, d, "s.json", `{"schema":2,"session_id":"s","state":"active","activity":"waiting","activity_at":"2026-10-07T10:00:00Z","future":{"a":1}}`)
	if err := Subagents(d, strings.NewReader(`{"session_id":"s","tasks":[{"id":"t","name":"n","status":"running"}]}`), io.Discard); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(d, "s.json"))
	for _, want := range []string{`"activity":"waiting"`, `"activity_at":"2026-10-07T10:00:00Z"`, `"future":{"a":1}`, `"tasks"`} {
		if !strings.Contains(string(b), want) {
			t.Fatalf("missing %s in %s", want, b)
		}
	}
	if err := Event(d, strings.NewReader(`{"session_id":"s","hook_event_name":"Stop"}`)); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(filepath.Join(d, "s.json"))
	if !strings.Contains(string(b), `"activity":"waiting"`) {
		t.Fatalf("Event dropped fields: %s", b)
	}
}
