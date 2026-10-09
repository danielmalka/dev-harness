package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigRoundTripPreservesBytes(t *testing.T) {
	home := t.TempDir()
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte(sample), 0o600)
	c, err := LoadConfig(home)
	if err != nil {
		t.Fatal(err)
	}
	if c.Language != "pt-br" || c.Mode != "global" || c.Sprites != `C:\sprites dir\x` ||
		c.Projects[`C:\Users\me\repo`] != "win" || c.Projects["/a/b"] != "ab" {
		t.Fatalf("parsed %+v", c)
	}
	want := "reviewers:\n  code:\n    - claude\n    - cli:codex/gpt\n  # note"
	if c.ReviewersRaw != want {
		t.Fatalf("reviewers = %q", c.ReviewersRaw)
	}
	if err := SaveConfig(home, c); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(home, "config.yaml"))
	if string(b) != sample {
		t.Fatalf("unchanged save altered bytes:\n%s", b)
	}
	// add, rename, remove a project; change language; foreign lines survive
	c.Language = "en"
	c.Projects["/c/d"] = "cd"
	c.Projects["/a/b"] = "ab2"
	delete(c.Projects, `C:\Users\me\repo`)
	if err := SaveConfig(home, c); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(filepath.Join(home, "config.yaml"))
	got := string(b)
	for _, must := range []string{"language: en\n", "  other: keep\n", "    - cli:codex/gpt\n", "extra: 1\n", "  # first\n", "  \"/a/b\": ab2\n", "  \"/c/d\": cd\n"} {
		if !strings.Contains(got, must) {
			t.Errorf("missing %q in\n%s", must, got)
		}
	}
	if strings.Contains(got, "win") {
		t.Error("removed entry still there")
	}
	back, err := LoadConfig(home)
	if err != nil || back.Projects["/c/d"] != "cd" || back.ReviewersRaw != want {
		t.Fatalf("reload %+v %v", back, err)
	}
	fi, _ := os.Stat(filepath.Join(home, "config.yaml"))
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("perm %v", fi.Mode().Perm())
	}
}

func TestSaveNewConfigQuotesWindowsPath(t *testing.T) {
	home := filepath.Join(t.TempDir(), "new")
	p := `C:\Users\me\repo`
	if err := SaveConfig(home, Config{Language: "en", Projects: map[string]string{p: "r"}}); err != nil {
		t.Fatal(err)
	}
	c, err := LoadConfig(home)
	if err != nil || c.Projects[p] != "r" || c.Language != "en" {
		t.Fatalf("%+v %v", c, err)
	}
}

func TestInvalidConfigIsReadableAndNotOverwritten(t *testing.T) {
	cases := map[string]string{
		"tab":       "projects:\n\t\"/a\": x\n",
		"dupkey":    "mode: repo\nmode: global\n",
		"dupproj":   "projects:\n  \"/a\": x\n  \"/a\": y\n",
		"colon":     "language: a:b\n",
		"unquoted":  "projects:\n  /a: x\n",
		"badmode":   "mode: sideways\n",
		"toplevel":  "- item\n",
		"orphan":    "  indented: 1\n",
		"noclosing": "projects:\n  \"/a: x\n",
	}
	for name, body := range cases {
		home := t.TempDir()
		p := filepath.Join(home, "config.yaml")
		os.WriteFile(p, []byte(body), 0o600)
		_, err := LoadConfig(home)
		if err == nil || !strings.Contains(err.Error(), "line ") {
			t.Errorf("%s: err = %v", name, err)
		}
		if err := SaveConfig(home, Config{Projects: map[string]string{"/z": "z"}}); err == nil {
			t.Errorf("%s: save should fail", name)
		}
		if b, _ := os.ReadFile(p); string(b) != body {
			t.Errorf("%s: file overwritten", name)
		}
	}
}

func TestResolveTable(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	glob := filepath.Join(home, "projects", "repo")
	c := Config{Projects: map[string]string{repo: "repo"}}
	// none
	if r := Resolve(repo, Config{}); r.Mode != "none" || r.Dir != "" {
		t.Errorf("none: %+v", r)
	}
	// registered but folder missing -> none
	if r := Resolve(repo, c); r.Mode != "none" {
		t.Errorf("registered/no folder: %+v", r)
	}
	// global only
	os.MkdirAll(glob, 0o700)
	if r := Resolve(repo, c); r.Mode != "global" || r.Dir != glob {
		t.Errorf("global: %+v", r)
	}
	// both: repo wins
	os.MkdirAll(filepath.Join(repo, ".harness"), 0o755)
	if r := Resolve(repo, c); r.Mode != "repo" || r.Dir != filepath.Join(repo, ".harness") {
		t.Errorf("both: %+v", r)
	}
	// repo only (unregistered)
	if r := Resolve(repo, Config{}); r.Mode != "repo" {
		t.Errorf("repo only: %+v", r)
	}
	// hostile hand-edited name cannot escape projects/
	if r := Resolve(repo+"x", Config{Projects: map[string]string{repo + "x": ".."}}); r.Mode != "none" {
		t.Errorf("traversal: %+v", r)
	}
}

func TestProjectsListSkipsAndCaps(t *testing.T) {
	home, repo := setup(t)
	EnsureHome()
	if list, w := Projects(home); len(list) != 0 || len(w) != 0 {
		t.Errorf("empty: %v %v", list, w)
	}
	Link(home, repo, "", false)
	c, _ := LoadConfig(home)
	c.Projects["/nope/gone"] = "gone"
	SaveConfig(home, c)
	list, w := Projects(home)
	if len(list) != 1 || list[0].Mode != "global" || list[0].Name != "repo" || len(w) != 1 ||
		!strings.Contains(w[0], "1 registered project(s) skipped") {
		t.Fatalf("%+v %v", list, w)
	}
	// cap
	base := t.TempDir()
	for i := 0; i < MaxProjects+1; i++ {
		d := filepath.Join(base, "p"+strings.Repeat("x", i%7)+string(rune('a'+i/26))+string(rune('a'+i%26)))
		os.MkdirAll(filepath.Join(d, ".harness"), 0o755)
		c.Projects[d] = filepath.Base(d)
	}
	SaveConfig(home, c)
	list, w = Projects(home)
	if len(list) != MaxProjects || len(w) < 2 {
		t.Errorf("cap: %d %v", len(list), w)
	}
	for i := 1; i < len(list); i++ {
		if list[i-1].Path >= list[i].Path {
			t.Fatal("not sorted by path")
		}
	}
}

func TestSaveKeepsProjectsHeaderComment(t *testing.T) {
	home := t.TempDir()
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects: # note\n  \"/a\": a\n"), 0o600)
	c, err := LoadConfig(home)
	if err != nil {
		t.Fatal(err)
	}
	c.Projects["/b"] = "b"
	if err := SaveConfig(home, c); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(home, "config.yaml"))
	if string(b) != "projects: # note\n  \"/a\": a\n  \"/b\": b\n" {
		t.Errorf("%q", b)
	}
}

func TestSaveConfigNilProjectsKeepsBlock(t *testing.T) {
	home := t.TempDir()
	body := "projects:\n  \"/a\": x\n"
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte(body), 0o600)
	c, _ := LoadConfig(home)
	c.Projects = nil
	if err := SaveConfig(home, c); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(home, "config.yaml")); string(b) != body {
		t.Fatalf("got %q", b)
	}
}

func TestValidNameStricter(t *testing.T) {
	for _, n := range []string{"a:b", "x.", "x ", strings.Repeat("a", 256)} {
		if validName(n) {
			t.Errorf("%q accepted", n)
		}
	}
	if !validName("ok-1.x") {
		t.Error("ok-1.x rejected")
	}
}

func TestProjectsSkipsBadKeys(t *testing.T) {
	home := t.TempDir()
	os.WriteFile(filepath.Join(home, "config.yaml"), []byte("projects:\n  \"rel/x\": a\n  \"/tmp/../tmp\": b\n"), 0o600)
	list, warns := Projects(home)
	if len(list) != 0 || len(warns) == 0 || !strings.Contains(strings.Join(warns, ";"), "2 registry key(s) ignored") {
		t.Fatalf("list=%v warns=%v", list, warns)
	}
}
