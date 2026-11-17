package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigKeys are the scalar keys `dh config` can read and write, in display order.
var ConfigKeys = []string{"mode", "dashboard.sprites", "worktree"}

// SetConfigValue sets (or, with value "" and unset, clears) one scalar key of <home>/config.yaml.
// The load-modify-save runs under config.lock like Link; an invalid existing file is an error and
// is never overwritten. dashboard.sprites is stored as typed, so a leading "~/" stays portable.
func SetConfigValue(home, key, value string, unset bool) error {
	switch key {
	case "reviewers":
		return fmt.Errorf("reviewers não é editável aqui: edite config.yaml ou .harness/project.yaml à mão")
	case "projects":
		return fmt.Errorf("projects não é editável aqui: use dh link")
	case "mode", "dashboard.sprites", "worktree":
	default:
		// `language` lands here like any nonexistent key: the option is gone (ADR-011), with no special message.
		return fmt.Errorf("chave desconhecida %q (%s)", key, strings.Join(ConfigKeys, ", "))
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
	case "mode":
		c.Mode = value
	case "worktree":
		c.Worktree = value
	default:
		c.Sprites = value
	}
	return SaveConfig(home, c)
}

func validConfigValue(key, v string) error {
	switch key {
	case "mode":
		if v != "repo" && v != "global" {
			return fmt.Errorf("mode inválido %q (repo, global)", v)
		}
	case "worktree":
		if v != "always" && v != "on-request" {
			return fmt.Errorf("worktree inválido %q (always, on-request)", v)
		}
	default: // same rule as the dashboard's spritesDirOf
		p := v
		if strings.HasPrefix(p, "~/") {
			if uh, err := os.UserHomeDir(); err == nil && uh != "" {
				p = filepath.Join(uh, p[2:])
			}
		}
		if fi, err := os.Stat(p); !filepath.IsAbs(p) || err != nil || !fi.IsDir() {
			return fmt.Errorf("dashboard.sprites %q não é um diretório absoluto que exista", v)
		}
	}
	return nil
}
