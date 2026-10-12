package harness

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const cfgBody = "# mine\nlanguage: en\nmode: repo\ndashboard:\n  sprites: /old\nreviewers:\n  code: [claude]\nprojects:\n  \"/a\": a\n"

func cfgHome(t *testing.T, body string) (string, string) {
	home := t.TempDir()
	p := filepath.Join(home, "config.yaml")
	if body != "" {
		os.WriteFile(p, []byte(body), 0o600)
	}
	return home, p
}

func TestSetConfigValue(t *testing.T) {
	dir := t.TempDir()
	home, p := cfgHome(t, cfgBody)
	for _, c := range []struct{ k, v string }{{"language", "pt-br"}, {"mode", "global"}, {"dashboard.sprites", dir}} {
		if err := SetConfigValue(home, c.k, c.v, false); err != nil {
			t.Fatalf("%s: %v", c.k, err)
		}
	}
	got, _ := LoadConfig(home)
	if got.Language != "pt-br" || got.Mode != "global" || got.Sprites != dir || got.Projects["/a"] != "a" {
		t.Fatalf("%+v", got)
	}
	b, _ := os.ReadFile(p)
	if !strings.Contains(string(b), "# mine\n") || !strings.Contains(string(b), "reviewers:\n  code: [claude]\n") {
		t.Errorf("unmodeled lines lost:\n%s", b)
	}
	// "~/" is stored as typed
	t.Setenv("HOME", dir)
	os.Mkdir(filepath.Join(dir, "sp"), 0o755)
	if err := SetConfigValue(home, "dashboard.sprites", "~/sp", false); err != nil {
		t.Fatal(err)
	}
	if got, _ := LoadConfig(home); got.Sprites != "~/sp" {
		t.Errorf("sprites %q", got.Sprites)
	}
}

func TestSetConfigValueRefusedLeavesFileIntact(t *testing.T) {
	home, p := cfgHome(t, cfgBody)
	for _, c := range []struct{ k, v, msg string }{
		{"language", "fr", "invalid language"}, {"mode", "foo", "invalid mode"},
		{"dashboard.sprites", "relativo", "not an absolute"}, {"dashboard.sprites", "/no/such/dir", "not an absolute"},
		{"reviewers", "x", "by hand"}, {"projects", "x", "dh link"}, {"nope", "x", "unknown key"},
	} {
		for _, unset := range []bool{false, true} {
			if unset && c.k != "reviewers" && c.k != "projects" && c.k != "nope" {
				continue
			}
			err := SetConfigValue(home, c.k, c.v, unset)
			if err == nil || !strings.Contains(err.Error(), c.msg) {
				t.Errorf("%s=%s unset=%v: %v", c.k, c.v, unset, err)
			}
		}
	}
	if b, _ := os.ReadFile(p); string(b) != cfgBody {
		t.Errorf("file changed:\n%s", b)
	}
	// an invalid existing file is an error, never overwritten
	bad := "mode: nope\n"
	os.WriteFile(p, []byte(bad), 0o600)
	if err := SetConfigValue(home, "language", "en", false); err == nil {
		t.Error("want error")
	}
	if b, _ := os.ReadFile(p); string(b) != bad {
		t.Errorf("invalid file overwritten: %s", b)
	}
}

func TestUnsetAndCreate(t *testing.T) {
	home, _ := cfgHome(t, cfgBody)
	for _, k := range ConfigKeys {
		if err := SetConfigValue(home, k, "", true); err != nil {
			t.Fatal(err)
		}
	}
	if c, _ := LoadConfig(home); c.Language != "" || c.Mode != "" || c.Sprites != "" || c.Projects["/a"] != "a" {
		t.Errorf("%+v", c)
	}
	home, _ = cfgHome(t, "") // no config.yaml yet
	if err := SetConfigValue(home, "language", "en", false); err != nil {
		t.Fatal(err)
	}
	if c, _ := LoadConfig(home); c.Language != "en" {
		t.Errorf("%+v", c)
	}
}
