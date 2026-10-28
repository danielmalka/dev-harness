package harness

import (
	"os"
	"strings"
	"testing"
)

// SaveConfig always writes a projects: block when none exists (pre-existing behavior), so the
// expectations below end with it.
const wtCfg = "language: en\n# keep lang\nworktree: always\n# keep me\n\nreviewers:\n  code: [claude]\nprojects:\n"

func readCfg(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// R3: unset drops only the key line; comment and blank lines after it survive.
func TestR3_UnsetWorktreeKeepsCommentAndBlank(t *testing.T) {
	home, p := cfgHome(t, strings.TrimSuffix(wtCfg, "projects:\n"))
	if err := SetConfigValue(home, "worktree", "", true); err != nil {
		t.Fatal(err)
	}
	got := readCfg(t, p)
	want := "language: en\n# keep lang\n# keep me\n\nreviewers:\n  code: [claude]\nprojects:\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
	if strings.Contains(got, "worktree") {
		t.Error("worktree key left")
	}
	if c, _ := LoadConfig(home); c.Worktree != "" {
		t.Errorf("%+v", c)
	}
}

// R3 regression: the same guard (scalarKeep) covers language and mode.
func TestR3_UnsetLanguageAndModeKeepComments(t *testing.T) {
	body := "language: en\n# keep lang\nmode: repo\n# keep mode\n\nworktree: always\nprojects:\n"
	home, p := cfgHome(t, body)
	if err := SetConfigValue(home, "language", "", true); err != nil {
		t.Fatal(err)
	}
	if got := readCfg(t, p); got != "# keep lang\nmode: repo\n# keep mode\n\nworktree: always\nprojects:\n" {
		t.Errorf("unset language: %q", got)
	}
	if err := SetConfigValue(home, "mode", "", true); err != nil {
		t.Fatal(err)
	}
	if got := readCfg(t, p); got != "# keep lang\n# keep mode\n\nworktree: always\nprojects:\n" {
		t.Errorf("unset mode: %q", got)
	}
}

// R1: set stores the value, byte-preserving every other line (a trailing comment line stays).
func TestR1_SetWorktreeValuesPreserveOtherLines(t *testing.T) {
	home, p := cfgHome(t, strings.TrimSuffix(wtCfg, "projects:\n"))
	if err := SetConfigValue(home, "worktree", "on-request", false); err != nil {
		t.Fatal(err)
	}
	if got := readCfg(t, p); got != strings.Replace(wtCfg, "worktree: always", "worktree: on-request", 1) {
		t.Errorf("%q", got)
	}
	if err := SetConfigValue(home, "worktree", "always", false); err != nil {
		t.Fatal(err)
	}
	if got := readCfg(t, p); got != wtCfg {
		t.Errorf("round trip changed bytes: %q", got)
	}
	// absent key: appended, existing lines untouched
	home, p = cfgHome(t, "# mine\nlanguage: en\nprojects:\n  \"/a\": a\n")
	if err := SetConfigValue(home, "worktree", "always", false); err != nil {
		t.Fatal(err)
	}
	got := readCfg(t, p)
	if !strings.HasPrefix(got, "# mine\nlanguage: en\nprojects:\n  \"/a\": a\n") || !strings.Contains(got, "worktree: always\n") {
		t.Errorf("%q", got)
	}
	if c, _ := LoadConfig(home); c.Worktree != "always" || c.Projects["/a"] != "a" {
		t.Errorf("%+v", c)
	}
}

// R2: any other value is refused, lists both valid values, and leaves the bytes unchanged.
func TestR2_InvalidWorktreeRefusedFileUnchanged(t *testing.T) {
	home, p := cfgHome(t, wtCfg)
	for _, v := range []string{"sometimes", "", "Always", "on_request", "always "} {
		err := SetConfigValue(home, "worktree", v, false)
		if err == nil || !strings.Contains(err.Error(), "always") || !strings.Contains(err.Error(), "on-request") {
			t.Errorf("%q: err %v", v, err)
		}
		if got := readCfg(t, p); got != wtCfg {
			t.Errorf("%q changed the file: %q", v, got)
		}
	}
}

// R8 (Go side): the stored value is not compared with always/on-request on load.
func TestR8_LoadAcceptsAnyScalarRejectsColon(t *testing.T) {
	home, _ := cfgHome(t, "worktree: sometimes\n")
	c, err := LoadConfig(home)
	if err != nil || c.Worktree != "sometimes" {
		t.Fatalf("%+v %v", c, err)
	}
	home, _ = cfgHome(t, "worktree: \"a:b\"\n")
	if c, err := LoadConfig(home); err != nil || c.Worktree != "a:b" {
		t.Errorf("quoted: %+v %v", c, err)
	}
	for _, v := range []string{"a:b", "{x", "[x", "|x", ">x", "&x", "*x", "!x"} {
		home, _ = cfgHome(t, "worktree: "+v+"\n")
		if _, err := LoadConfig(home); err == nil {
			t.Errorf("worktree: %s should abort the parse", v)
		}
		home, _ = cfgHome(t, "language: "+v+"\n")
		if _, err := LoadConfig(home); err == nil {
			t.Errorf("language: %s should abort the parse (control)", v)
		}
	}
}
