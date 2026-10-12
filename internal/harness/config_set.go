package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigKeys are the scalar keys `dh config` can read and write, in display order.
var ConfigKeys = []string{"language", "mode", "dashboard.sprites"}

// SetConfigValue sets (or, with value "" and unset, clears) one scalar key of <home>/config.yaml.
// The load-modify-save runs under config.lock like Link; an invalid existing file is an error and
// is never overwritten. dashboard.sprites is stored as typed, so a leading "~/" stays portable.
func SetConfigValue(home, key, value string, unset bool) error {
	switch key {
	case "reviewers":
		return fmt.Errorf("reviewers is not editable here: edit config.yaml or .harness/project.yaml by hand")
	case "projects":
		return fmt.Errorf("projects is not editable here: use dh link")
	case "language", "mode", "dashboard.sprites":
	default:
		return fmt.Errorf("unknown key %q (language, mode, dashboard.sprites)", key)
	}
	if unset {
		value = ""
	} else if err := validConfigValue(key, value); err != nil {
		return err
	}
	unlock, err := lock(home)
	if err != nil {
		return err
	}
	defer unlock()
	c, err := LoadConfig(home)
	if err != nil {
		return err
	}
	switch key {
	case "language":
		c.Language = value
	case "mode":
		c.Mode = value
	default:
		c.Sprites = value
	}
	return SaveConfig(home, c)
}

func validConfigValue(key, v string) error {
	switch key {
	case "language":
		if v != "en" && v != "pt-br" {
			return fmt.Errorf("invalid language %q (en, pt-br)", v)
		}
	case "mode":
		if v != "repo" && v != "global" {
			return fmt.Errorf("invalid mode %q (repo, global)", v)
		}
	default: // same rule as the dashboard's spritesDirOf
		p := v
		if strings.HasPrefix(p, "~/") {
			if uh, err := os.UserHomeDir(); err == nil && uh != "" {
				p = filepath.Join(uh, p[2:])
			}
		}
		if fi, err := os.Stat(p); !filepath.IsAbs(p) || err != nil || !fi.IsDir() {
			return fmt.Errorf("dashboard.sprites %q is not an absolute existing directory", v)
		}
	}
	return nil
}
