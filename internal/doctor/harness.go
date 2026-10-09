package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/danielmalka/dev-harness/internal/harness"
)

// printHarness reports the resolved harness for the current directory (PRD-014 R14, R9, R6).
// Read-only: ~/.claude/settings.json is only ever read.
func printHarness(out io.Writer) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(out, "harness: unavailable")
		return
	}
	var settings string
	if h, err := os.UserHomeDir(); err == nil {
		settings = filepath.Join(h, ".claude", "settings.json")
	}
	printHarnessFor(out, cwd, harness.Home(), settings)
}

func printHarnessFor(out io.Writer, cwd, home, settings string) {
	repo := absolutePath(cwd)
	if r, err := filepath.EvalSymlinks(repo); err == nil {
		repo = r
	}
	cfg, cfgErr := harness.LoadConfig(home)
	if cfgErr != nil {
		fmt.Fprintf(out, "warning: %v\n", cfgErr)
	}
	res := harness.Resolve(repo, cfg)
	if home == "" {
		fmt.Fprintln(out, "home: unknown")
		fmt.Fprintln(out, "hint: cannot determine the harness home: set DH_HOME, HOME or USERPROFILE")
	} else {
		fmt.Fprintf(out, "home: %s\n", home)
	}
	fmt.Fprintf(out, "mode: %s\n", res.Mode)
	if res.Mode != "none" {
		fmt.Fprintf(out, "dir: %s\n", res.Dir)
		for _, name := range []string{"project.yaml", "MEMORY.md", "EPOCHAL.md", "RISKS.md"} {
			state := "missing"
			if isFile(filepath.Join(res.Dir, name)) {
				state = "present"
			}
			fmt.Fprintf(out, "%s: %s\n", name, state)
		}
	}
	name, registered := cfg.Projects[repo]
	if res.Mode == "repo" {
		if registered && harness.IsHarness(filepath.Join(home, "projects", name)) {
			fmt.Fprintf(out, "warning: both %s and %s exist: the internal one wins\n",
				res.Dir, filepath.Join(home, "projects", name))
		}
		if !registered {
			fmt.Fprintf(out, "warning: %s is not registered in %s: not shown in the dashboard until dh link\n",
				repo, filepath.Join(home, "config.yaml"))
		}
		if isFile(filepath.Join(res.Dir, "local.yaml")) && notIgnored(repo, ".harness/local.yaml") {
			fmt.Fprintln(out, "warning: .harness/local.yaml is not ignored by git")
		}
	}
	if res.Mode == "none" {
		fmt.Fprintln(out, "hint: no harness for this directory: run /dh:setup (or dh link)")
	}
	if home != "" && !within(repo, home) {
		if msg := settingsHint(settings, home); msg != "" {
			fmt.Fprintln(out, msg)
		}
	}
}

// notIgnored is true only when git answers "not ignored" (exit 1); no git or no repo is silent.
func notIgnored(dir, rel string) bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	c := exec.Command(git, "check-ignore", "-q", rel)
	c.Dir = dir
	var ee *exec.ExitError
	return errors.As(c.Run(), &ee) && ee.ExitCode() == 1
}

func within(root, p string) bool {
	rel, err := filepath.Rel(root, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// settingsHint returns the R9 hint, or "" when additionalDirectories covers home.
func settingsHint(settings, home string) string {
	suggest := "~/.harness"
	if uh, err := os.UserHomeDir(); err != nil || home != filepath.Join(uh, ".harness") {
		suggest = home // custom DH_HOME
	}
	hint := func(why string) string {
		return fmt.Sprintf("hint: ~/.claude/settings.json %s covering %s; add %q (dh never writes that file)", why, home, suggest)
	}
	missing := hint("has no permissions.additionalDirectories entry")
	data, err := os.ReadFile(settings)
	if err != nil {
		return missing
	}
	var doc struct {
		Permissions struct {
			AdditionalDirectories []string `json:"additionalDirectories"`
		} `json:"permissions"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return hint("could not be parsed, so no permissions.additionalDirectories entry was found")
	}
	userHome, _ := os.UserHomeDir()
	for _, e := range doc.Permissions.AdditionalDirectories {
		switch {
		case e == "~":
			e = userHome
		case strings.HasPrefix(e, "~/"):
			e = filepath.Join(userHome, e[2:])
		}
		e = os.ExpandEnv(e)
		if abs, err := filepath.Abs(e); err == nil && within(abs, home) {
			return ""
		}
	}
	return missing
}
