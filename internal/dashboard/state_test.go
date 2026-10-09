package dashboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// writeCfg writes <home>/config.yaml with dashboard.sprites and the given repos (name = base dir).
func writeCfg(t *testing.T, home, sprites string, repos ...string) {
	t.Helper()
	c := harness.Config{Sprites: sprites, Projects: map[string]string{}}
	for _, r := range repos {
		c.Projects[r] = filepath.Base(r)
	}
	if err := harness.SaveConfig(home, c); err != nil {
		t.Fatal(err)
	}
}

func register(t *testing.T, home string, repos ...string) {
	t.Helper()
	writeCfg(t, home, "", repos...)
}

type apiState struct {
	Config   jsonConfig    `json:"config"`
	Projects []jsonProject `json:"projects"`
	Sessions []jsonSession `json:"sessions"`
	Warnings []string      `json:"warnings"`
}

func getState(t *testing.T, cfg Config) (apiState, string) {
	t.Helper()
	w := do(Handler(cfg), "GET", "127.0.0.1:1", "/api/state")
	var st apiState
	if err := json.Unmarshal(w.Body.Bytes(), &st); err != nil {
		t.Fatal(err)
	}
	return st, w.Body.String()
}

func TestGlobalProjectSessionsAndProgress(t *testing.T) {
	home, snaps, repo := t.TempDir(), t.TempDir(), t.TempDir()
	name := filepath.Base(repo)
	register(t, home, repo)
	g := filepath.Join(home, "projects", name)
	mk(t, filepath.Join(g, "tasks/T-1/TASK.md"), "| Status | concluída |\n| PRD (RF-<n>) | PRD-002 |\n")
	mk(t, filepath.Join(g, "prd/PRD-003-g.md"), "| Status | rascunho |\n")
	mk(t, filepath.Join(repo, "docs/prd/PRD-002-b.md"), "| Status | aprovado |\n")
	now := time.Now()
	for id, cwd := range map[string]string{"root": repo, "sub": filepath.Join(repo, "src", "deep"), "out": t.TempDir()} {
		mk(t, filepath.Join(snaps, id+".json"), fmt.Sprintf(`{"schema":2,"session_id":%q,"state":"active","activity":"working","cwd":%q,"updated_at":%q}`,
			id, cwd, now.UTC().Format(time.RFC3339)))
	}
	st, raw := getState(t, Config{Home: home, SnapshotDir: snaps, Stale: DefaultStale, DoneDecay: time.Minute, Port: 4747})
	if st.Config.Home != home || st.Config.Port != 4747 || st.Config.Sprites != "" || strings.Contains(raw, `"roots"`) {
		t.Fatalf("config: %s", raw)
	}
	if len(st.Projects) != 1 {
		t.Fatalf("projects: %s", raw)
	}
	p := st.Projects[0]
	if p.Name != name || p.Path != repo || p.Mode != "global" || p.Harness != g {
		t.Fatalf("project: %+v", p)
	}
	if len(p.Open) != 2 || p.Open[0].PRD != "PRD-002" || p.Open[0].Done != 1 || p.Open[0].Total != 1 || p.Open[1].PRD != "PRD-003" || !p.Open[1].NoTickets {
		t.Fatalf("bars: %+v", p.Open)
	}
	byID := map[string]string{}
	for _, s := range st.Sessions {
		byID[s.ID] = s.Project
	}
	if byID["root"] != repo || byID["sub"] != repo || byID["out"] != "" {
		t.Fatalf("attribution: %v", byID)
	}
}

func TestProjectsParityWithRegistryAndOmissions(t *testing.T) {
	home := t.TempDir()
	repoA, repoB, none := t.TempDir(), t.TempDir(), t.TempDir()
	mk(t, filepath.Join(repoA, ".harness", "x"), "")
	register(t, home, repoA, repoB, filepath.Join(repoB, "missing"), none)
	if err := os.MkdirAll(filepath.Join(home, "projects", filepath.Base(repoB)), 0o700); err != nil {
		t.Fatal(err)
	}
	st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir()})
	want, _ := harness.Projects(home)
	if len(want) != 2 || len(st.Projects) != len(want) {
		t.Fatalf("got %+v want %+v", st.Projects, want)
	}
	for i, w := range want {
		g := st.Projects[i]
		if g.Name != w.Name || g.Path != w.Path || g.Mode != w.Mode || g.Harness != w.Harness || (g.Mode != "repo" && g.Mode != "global") {
			t.Fatalf("item %d: %+v vs %+v", i, g, w)
		}
	}
	if len(st.Warnings) != 1 || !strings.Contains(st.Warnings[0], "2 registered project(s) skipped") {
		t.Fatalf("warnings: %v", st.Warnings)
	}
}

func TestProjectsCap(t *testing.T) {
	home, base := t.TempDir(), t.TempDir()
	var repos []string
	for i := 0; i < harness.MaxProjects+1; i++ {
		r := filepath.Join(base, fmt.Sprintf("p%03d", i))
		mk(t, filepath.Join(r, ".harness", "x"), "")
		repos = append(repos, r)
	}
	register(t, home, repos...)
	st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir()})
	if len(st.Projects) != harness.MaxProjects || len(st.Warnings) != 1 || !strings.Contains(st.Warnings[0], "showing the first") {
		t.Fatalf("got %d projects, warnings %v", len(st.Projects), st.Warnings)
	}
}

func TestConfigRereadWithServerUp(t *testing.T) {
	home, a, b := t.TempDir(), t.TempDir(), t.TempDir()
	mk(t, filepath.Join(a, ".harness", "x"), "")
	mk(t, filepath.Join(b, ".harness", "x"), "")
	register(t, home, a)
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir()}) // CacheTTL 0: every request rebuilds
	count := func() (n int, sprites string) {
		var st apiState
		_ = json.Unmarshal(do(h, "GET", "127.0.0.1:1", "/api/state").Body.Bytes(), &st)
		return len(st.Projects), st.Config.Sprites
	}
	if n, _ := count(); n != 1 {
		t.Fatalf("before: %d", n)
	}
	sp := t.TempDir()
	mk(t, filepath.Join(sp, "frente.png"), "NEW")
	writeCfg(t, home, sp, a, b)
	n, got := count()
	if n != 2 || got != sp {
		t.Fatalf("after edit: %d projects, sprites %q", n, got)
	}
	if w := do(h, "GET", "127.0.0.1:1", "/sprite/frente"); w.Body.String() != "NEW" {
		t.Fatalf("sprites not reloaded: %q", w.Body.String())
	}
	writeCfg(t, home, "", a)
	if n, got := count(); n != 1 || got != "" {
		t.Fatalf("after revert: %d %q", n, got)
	}
}

func TestTokenLivesInHomeDashboard(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DH_HOME", home)
	_, p, err := newStopToken(4803)
	if err != nil {
		t.Fatal(err)
	}
	if p != filepath.Join(home, "dashboard", "stop-4803.token") {
		t.Fatalf("path %s", p)
	}
}

func TestSpritesPathRules(t *testing.T) {
	home := t.TempDir()
	uh := t.TempDir()
	t.Setenv("HOME", uh)
	t.Setenv("USERPROFILE", uh)
	mk(t, filepath.Join(uh, "sp", "x"), "")
	get := func(v string) apiState {
		writeCfg(t, home, v)
		st, _ := getState(t, Config{Home: home, SnapshotDir: t.TempDir()})
		return st
	}
	if st := get("~/sp"); st.Config.Sprites != filepath.Join(uh, "sp") || len(st.Warnings) != 0 {
		t.Fatalf("tilde: %+v", st)
	}
	for _, bad := range []string{"rel/sp", filepath.Join(uh, "missing")} {
		st := get(bad)
		if st.Config.Sprites != "" || len(st.Warnings) != 1 || !strings.Contains(st.Warnings[0], "is not an absolute existing directory") {
			t.Fatalf("%q: %+v", bad, st)
		}
	}
}

func TestNestedProjectAttribution(t *testing.T) {
	home, snaps, base := t.TempDir(), t.TempDir(), t.TempDir()
	a, ab, abc := filepath.Join(base, "a"), filepath.Join(base, "a", "b"), filepath.Join(base, "a", "bc")
	for _, d := range []string{a, ab, abc} {
		mk(t, filepath.Join(d, ".harness", "x"), "")
	}
	register(t, home, a, ab)
	now := time.Now().UTC().Format(time.RFC3339)
	for id, cwd := range map[string]string{"deep": filepath.Join(ab, "x"), "up": filepath.Join(a, "y"), "sib": abc} {
		mk(t, filepath.Join(snaps, id+".json"), fmt.Sprintf(`{"schema":2,"session_id":%q,"state":"active","activity":"working","cwd":%q,"updated_at":%q}`, id, cwd, now))
	}
	st, _ := getState(t, Config{Home: home, SnapshotDir: snaps, Stale: DefaultStale, DoneDecay: time.Minute})
	got := map[string]string{}
	for _, s := range st.Sessions {
		got[s.ID] = s.Project
	}
	if got["deep"] != ab || got["up"] != a || got["sib"] != a {
		t.Fatalf("%v", got)
	}
}

func TestSiblingPrefixNotAttributed(t *testing.T) {
	if within("/a/b", "/a/bc") || !within("/a/b", "/a/b/x") {
		t.Fatal("prefix match must respect path boundaries")
	}
}

func TestConcurrentStateWhileConfigRewritten(t *testing.T) {
	home, a, b := t.TempDir(), t.TempDir(), t.TempDir()
	mk(t, filepath.Join(a, ".harness", "x"), "")
	mk(t, filepath.Join(b, ".harness", "x"), "")
	register(t, home, a)
	h := Handler(Config{Home: home, SnapshotDir: t.TempDir()})
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				if c := do(h, "GET", "127.0.0.1:1", "/api/state").Code; c != 200 {
					t.Errorf("status %d", c)
				}
			}
		}()
	}
	go func() { wg.Wait(); close(done) }()
	for i := 0; ; i++ {
		select {
		case <-done:
			return
		default:
			if i%2 == 0 {
				register(t, home, a, b)
			} else {
				register(t, home, a)
			}
		}
	}
}
