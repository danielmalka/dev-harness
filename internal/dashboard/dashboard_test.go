package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

func mk(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestStatusCasesFixture(t *testing.T) {
	b, err := os.ReadFile("testdata/status-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ Name, MD, Want string }
	if err := json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		if got := TaskStatus(c.MD); got != c.Want {
			t.Errorf("%s: got %s want %s", c.Name, got, c.Want)
		}
	}
}

func TestProjects(t *testing.T) {
	if p, w := Projects(""); len(p) != 0 || w == "" {
		t.Fatal("empty variable must warn with zero projects")
	}
	r1, r2 := t.TempDir(), t.TempDir()
	mk(t, filepath.Join(r1, "a", ".harness", "x"), "")
	mk(t, filepath.Join(r1, "nope", "file"), "")
	mk(t, filepath.Join(r2, "b", ".harness", "x"), "")
	p, w := Projects(r1 + ";" + filepath.Join(r2, "missing") + ";" + r2)
	if w != "" || len(p) != 2 {
		t.Fatalf("got %v warn %q", p, w)
	}
}

func TestPRDLink(t *testing.T) {
	cases := map[string]string{
		"| Story / PRD | PRD-003 (RF-1) |":              "PRD-003",
		"| PRD (RF-<n>) | PRD-011 (R1, R2) |":           "PRD-011",
		"| PRD (RF-<n>) | fora do PRD-008 (kit only) |": "",
		"| Story / PRD | nenhum |":                      "",
		"| Story / PRD | story X |":                     "",
		"| Tipo | PRD-001 |":                            "",
	}
	for md, want := range cases {
		if got := PRDLink(md + "\n"); got != want {
			t.Errorf("%q: got %q want %q", md, got, want)
		}
	}
}

func TestProjectProgress(t *testing.T) {
	p := t.TempDir()
	mk(t, filepath.Join(p, "docs/prd/PRD-001-a.md"), "| Status | entregue em 2026-01-01 |\n")
	mk(t, filepath.Join(p, "docs/prd/PRD-002-b.md"), "| Status | aprovado |\n")
	mk(t, filepath.Join(p, ".harness/prd/PRD-003-c.md"), "**Status:** rascunho\n")
	mk(t, filepath.Join(p, ".harness/tasks/T-1/TASK.md"), "| Status | pronta |\n| Story / PRD | PRD-001 |\n") // open ticket, delivered PRD: still 100%
	mk(t, filepath.Join(p, ".harness/tasks/T-2/TASK.md"), "| Status | concluída |\n| PRD (RF-<n>) | PRD-002 |\n")
	mk(t, filepath.Join(p, ".harness/tasks/T-3/TASK.md"), "| Status | bloqueada |\n| PRD (RF-<n>) | PRD-002 |\n")
	mk(t, filepath.Join(p, ".harness/tasks/T-4/TASK.md"), "| Status | pronta |\n| PRD (RF-<n>) | fora do PRD-003 |\n")
	pr := ProjectProgress(p)
	if len(pr.Delivered) != 1 || pr.Delivered[0] != "PRD-001" || len(pr.Open) != 2 {
		t.Fatalf("%+v", pr)
	}
	if b := pr.Open[0]; b.PRD != "PRD-002" || b.Done != 1 || b.Blocked != 1 || b.Total != 2 || b.NoTickets {
		t.Fatalf("bar %+v", b)
	}
	if b := pr.Open[1]; b.PRD != "PRD-003" || !b.NoTickets || b.Total != 0 {
		t.Fatalf("bar %+v", b)
	}
}

func f(v float64) *float64 { return &v }

func TestOpenSessionsAndLimits(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	projs := []Project{{Name: "p", Path: "/w/p"}, {Name: "sub", Path: "/w/p/sub"}, {Name: "q", Path: "/w/pq"}}
	all := []snapshot.Session{
		{SessionID: "fresh", State: "active", Activity: "working", CWD: "/w/p/src/deep", UpdatedAt: now.Add(-10 * time.Second)},
		{SessionID: "old", State: "active", Activity: "working", CWD: "/w/p", UpdatedAt: now.Add(-3 * time.Minute)},
		{SessionID: "closed", State: "closed", Activity: "idle", CWD: "/w/p", UpdatedAt: now, Limits: snapshot.Limits{FiveHour: f(10)}},
		{SessionID: "nested", State: "active", Activity: "done", ActivityAt: now.Add(-6 * time.Minute), CWD: "/w/p/sub", UpdatedAt: now},
		{SessionID: "nolimit", State: "idle", Activity: "idle", CWD: "/elsewhere", UpdatedAt: now.Add(-time.Second)},
		{SessionID: "oldlimit", State: "idle", Activity: "idle", UpdatedAt: now.Add(-2 * time.Hour), Limits: snapshot.Limits{FiveHour: f(90)}},
	}
	got := OpenSessions(all, projs, now, DefaultStale, snapshot.DefaultDoneDecay)
	byID := map[string]OpenSession{}
	for _, s := range got {
		byID[s.SessionID] = s
	}
	if len(got) != 3 || byID["old"].SessionID != "" || byID["closed"].SessionID != "" {
		t.Fatalf("stale/closed leaked: %d", len(got))
	}
	if byID["fresh"].Project != "/w/p" || byID["nested"].Project != "/w/p/sub" || byID["nolimit"].Project != "" {
		t.Fatalf("attribution wrong: %+v", got)
	}
	if byID["nested"].Activity != snapshot.ActivityIdle {
		t.Fatal("decayed done must read idle")
	}
	lim := LatestLimits(all, now) // newest with limits is "closed" (now), not the 2 h old one
	if !lim.OK || *lim.FiveHour != 10 || lim.FiveHourAge != 0 {
		t.Fatalf("%+v", lim)
	}
	if LatestLimits(all[:2], now).OK {
		t.Fatal("no limits must be OK=false")
	}
}

func TestAvatarPriorityPairs(t *testing.T) {
	order := []string{StateWaiting, StateError, StateWorking, StateDone, StateAttn, StateIdle}
	act := map[string]string{StateWaiting: "waiting", StateError: "error", StateWorking: "working", StateDone: "done", StateIdle: "idle"}
	state := func(s string) (ss []OpenSession, l AccountLimits) {
		if s == StateAttn {
			return nil, AccountLimits{SevenDay: f(80), OK: true}
		}
		return []OpenSession{{Activity: act[s]}}, AccountLimits{}
	}
	for i := 0; i+1 < len(order); i++ {
		hi, lo := order[i], order[i+1]
		s1, l1 := state(hi)
		s2, l2 := state(lo)
		s := append(s1, s2...)
		l := l1
		if l2.OK {
			l = l2
		}
		if got := AvatarState(s, l); got != hi {
			t.Errorf("%s vs %s: got %s", hi, lo, got)
		}
	}
	if AvatarState(nil, AccountLimits{FiveHour: f(79.9), OK: true}) != StateIdle {
		t.Error("79.9 must not alert")
	}
	if AvatarState(nil, AccountLimits{}) != StateIdle {
		t.Error("nothing = idle")
	}
}

func TestLimitsPerFieldAndStale(t *testing.T) {
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	all := []snapshot.Session{
		{UpdatedAt: now.Add(-time.Minute), Limits: snapshot.Limits{FiveHour: f(30)}},
		{UpdatedAt: now.Add(-6 * time.Hour), Limits: snapshot.Limits{FiveHour: f(95), SevenDay: f(90)}},
	}
	l := LatestLimits(all, now)
	if *l.FiveHour != 30 || *l.SevenDay != 90 || l.SevenDayAge != 6*time.Hour || l.FiveHourAge != time.Minute {
		t.Fatalf("%+v", l)
	}
	if AvatarState(nil, l) != StateIdle {
		t.Fatal("limit older than 5h must not raise attention")
	}
	l.SevenDayAge = time.Hour
	if AvatarState(nil, l) != StateAttn {
		t.Fatal("fresh 90% must raise attention")
	}
}

func TestUnknownActivityIsIdle(t *testing.T) {
	got := AvatarState([]OpenSession{{Activity: ""}, {Activity: "bogus"}, {Activity: "done"}}, AccountLimits{})
	if got != StateDone || !HigherPriority(StateDone, "") || HigherPriority("", StateIdle) {
		t.Fatalf("got %s", got)
	}
}

func TestProjectsAbsDedupeAndAllFail(t *testing.T) {
	r := t.TempDir()
	mk(t, filepath.Join(r, "a", ".harness", "x"), "")
	p, _ := Projects(r + ";" + r + string(filepath.Separator) + ".")
	if len(p) != 1 {
		t.Fatalf("dedupe: %v", p)
	}
	if p, w := Projects(filepath.Join(r, "no1") + ";" + filepath.Join(r, "no2")); len(p) != 0 || w == "" {
		t.Fatal("all roots failing must warn")
	}
}

func TestPosesJSONJevmonStyle(t *testing.T) {
	d := t.TempDir()
	mk(t, filepath.Join(d, "poses.json"), `{"voz_escuta":{"pose":"fone","fps":0},"voz_fala":{"pose":"explicando","fps":8},"sem_conexao":{"pose":"assustado"},"erro":{"pose":"triste","fps":2}}`)
	s := LoadSprites(d)
	if len(s.Warnings) != 0 || s.PoseFor(StateError) != (PoseSpec{Pose: "triste", FPS: 2}) {
		t.Fatalf("jevmon-style file must apply: %v %v", s.Warnings, s.PoseFor(StateError))
	}
}

func TestPosesJSONUnreadable(t *testing.T) {
	d := t.TempDir()
	mk(t, filepath.Join(d, "poses.json", "x"), "") // a directory, not a file
	if s := LoadSprites(d); len(s.Warnings) != 1 {
		t.Fatalf("read error must warn: %v", s.Warnings)
	}
}

func TestSprites(t *testing.T) {
	// no folder: kit examples only
	s := LoadSprites("")
	for _, st := range []string{StateWaiting, StateError, StateWorking, StateDone, StateAttn, StateIdle} {
		if _, ok := s.Frame(s.PoseFor(st).Pose, 0); !ok {
			t.Errorf("no example for %s", st)
		}
	}
	// "../x" and unknown poses are never served
	for _, n := range []string{"../x", "x", ""} {
		if _, ok := s.Frame(n, 0); ok {
			t.Errorf("%q must 404", n)
		}
	}
	// empty folder = examples; partial folder = owner first, example for the rest
	d := t.TempDir()
	s = LoadSprites(d)
	if _, ok := s.Frame("frente", 0); !ok {
		t.Fatal("empty folder must fall back to example")
	}
	mk(t, filepath.Join(d, "celular_00.png"), "A")
	mk(t, filepath.Join(d, "celular_01.png"), "B")
	mk(t, filepath.Join(d, "frente.png"), "F")
	if fr, _ := s.Frames("celular"); len(fr) != 2 || string(fr[1]) != "B" {
		t.Fatalf("frames %v", fr)
	}
	if b, _ := s.Frame("frente", 0); string(b) != "F" {
		t.Fatal("owner pose must win")
	}
	if a, _ := s.Frame("caneca", 0); string(a) != "F" {
		t.Fatal("pose with no file anywhere must fall back to frente")
	}
	if b, _ := s.Frame("duvida", 0); len(b) < 8 {
		t.Fatal("missing owner pose must fall to example")
	}
	if _, ok := s.Frame("celular", 2); ok {
		t.Fatal("out of range frame")
	}
}

func TestPosesJSON(t *testing.T) {
	d := t.TempDir()
	mk(t, filepath.Join(d, "poses.json"), `{"trabalhando":{"pose":"explicando","fps":12}}`)
	if s := LoadSprites(d); s.PoseFor(StateWorking) != (PoseSpec{Pose: "explicando", FPS: 12}) || len(s.Warnings) != 0 {
		t.Fatalf("valid override: %+v %v", s.table, s.Warnings)
	}
	for _, bad := range []string{`{`, `{"trabalhando":{"pose":"explicando","fps":61}}`, `{"trabalhando":{"pose":"nope"}}`,
		`{"voando":{"pose":"frente"}}`, `{"trabalhando":{"pose":"frente","x":1}}`} {
		mk(t, filepath.Join(d, "poses.json"), bad)
		s := LoadSprites(d)
		if len(s.Warnings) != 1 || s.PoseFor(StateWorking).Pose != "celular" || !strings.Contains(s.Warnings[0], "poses.json") {
			t.Errorf("%s: want default + warning, got %v %v", bad, s.PoseFor(StateWorking), s.Warnings)
		}
	}
}

func TestEmbeddedSize(t *testing.T) {
	if n := EmbeddedBytes(); n == 0 || n > 50_000 {
		t.Fatalf("embedded sprites %d bytes", n)
	}
}
