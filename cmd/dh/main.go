package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielmalka/dev-harness/internal/build"
	"github.com/danielmalka/dev-harness/internal/dashboard"
	"github.com/danielmalka/dev-harness/internal/doctor"
	"github.com/danielmalka/dev-harness/internal/harness"
	"github.com/danielmalka/dev-harness/internal/kit"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printError(stderr, usage())
		return 1
	}
	if args[0] == "validate" {
		return runValidate(args[1:], stdout, stderr)
	}
	if args[0] == "doctor" {
		if len(args) > 2 {
			printError(stderr, "usage: dh doctor [plugin-dir-or-kit-root]")
			return 1
		}
		target := "."
		if len(args) == 2 {
			target = args[1]
		}
		return doctor.Run(target, stdout)
	}
	if args[0] == "dashboard" {
		return dashboard.Run(args[1:], stdout, stderr)
	}
	switch args[0] {
	case "harness-path":
		return runHarnessPath(args[1:], stdout, stderr)
	case "link":
		return runLink(args[1:], stdout, stderr)
	case "projects":
		return runProjects(args[1:], stdout, stderr)
	}
	if args[0] == "build" {
		return runBuild(args[1:], stdout, stderr)
	}
	if len(args) < 2 || args[0] != "snapshot" {
		printError(stderr, usage())
		return 1
	}

	dir := snapshot.SnapshotDir()
	switch args[1] {
	case "event":
		if err := snapshot.Event(dir, stdin); err != nil {
			printError(stderr, err.Error())
		}
		return 0
	case "subagents":
		if err := snapshot.Subagents(dir, stdin, stdout); err != nil {
			printError(stderr, err.Error())
		}
		return 0
	case "statusline":
		if err := snapshot.Statusline(dir, stdin, stdout); err != nil {
			printError(stderr, err.Error())
		}
		return 0
	case "prune":
		days, err := pruneDays(args[2:])
		if err != nil {
			printError(stderr, err.Error())
			return 1
		}
		count, err := snapshot.Prune(dir, days)
		if err != nil {
			printError(stderr, err.Error())
			return 1
		}
		fmt.Fprintf(stdout, "pruned %d\n", count)
		return 0
	default:
		printError(stderr, "unknown snapshot subcommand: "+args[1])
		return 1
	}
}

func runBuild(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("dh build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	targets := flags.String("targets", "", "comma-separated GOOS/GOARCH targets")
	noBinaries := flags.Bool("no-binaries", false, "skip cross-compiled binaries")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		if err == nil {
			printError(stderr, "usage: dh build [--targets os/arch,...] [--no-binaries]")
		}
		return 1
	}
	parsedTargets, err := parseTargets(*targets)
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	err = build.Build(".", build.Options{
		Targets:            parsedTargets,
		NoBinaries:         *noBinaries,
		UpdateRootManifest: true,
	}, stdout)
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	return 0
}

func parseTargets(value string) ([]build.Target, error) {
	if value == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	result := make([]build.Target, 0, len(parts))
	for _, part := range parts {
		pieces := strings.Split(part, "/")
		if len(pieces) != 2 || pieces[0] == "" || pieces[1] == "" {
			return nil, fmt.Errorf("invalid target: %s", part)
		}
		result = append(result, build.Target{OS: pieces[0], Arch: pieces[1]})
	}
	return result, nil
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	sourceOnly := false
	target := "."
	targetSet := false
	for _, arg := range args {
		if arg == "--source-only" {
			if sourceOnly {
				printError(stderr, "usage: dh validate [--source-only] [path]")
				return 1
			}
			sourceOnly = true
			continue
		}
		if targetSet {
			printError(stderr, "usage: dh validate [--source-only] [path]")
			return 1
		}
		target = arg
		targetSet = true
	}
	report, err := kit.Validate(target, kit.Options{SourceOnly: sourceOnly})
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	kit.WriteReport(stdout, report)
	if !report.Passed() {
		return 1
	}
	return 0
}

func usage() string {
	return "usage: dh <snapshot <event|subagents|statusline|prune>|validate [--source-only] [path]|build [--targets os/arch,...] [--no-binaries]|doctor [plugin-dir-or-kit-root]|dashboard [--port N] [--stale D] [--done-decay D] [--detach]|harness-path [dir] [--json]|link [name] [--off] [dir] [--json]|projects [--json]>"
}

func pruneDays(args []string) (int, error) {
	if len(args) == 0 {
		return 7, nil
	}
	if len(args) != 2 || args[0] != "--days" {
		return 0, fmt.Errorf("usage: dh snapshot prune [--days N]")
	}
	days, err := strconv.Atoi(args[1])
	if err != nil || days < 0 {
		return 0, fmt.Errorf("invalid days: %s", args[1])
	}
	return days, nil
}

func printError(w io.Writer, message string) {
	message = strings.ReplaceAll(strings.ReplaceAll(message, "\n", " "), "\r", " ")
	fmt.Fprintln(w, message)
}

// splitArgs separates flags from positionals so they may appear in any order.
func splitArgs(args []string, bools ...string) (pos []string, set map[string]bool, ok bool) {
	set = map[string]bool{}
	for _, a := range args {
		if strings.HasPrefix(a, "--") {
			found := false
			for _, b := range bools {
				if a == b {
					set[b], found = true, true
				}
			}
			if !found {
				return nil, nil, false
			}
			continue
		}
		pos = append(pos, a)
	}
	return pos, set, true
}

func writeJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
}

func printResolved(w io.Writer, r harness.Resolved) {
	fmt.Fprintln(w, "mode: "+r.Mode)
	if r.Dir != "" {
		fmt.Fprintln(w, "dir: "+r.Dir)
	}
}

func runHarnessPath(args []string, stdout, stderr io.Writer) int {
	pos, f, ok := splitArgs(args, "--json")
	if !ok || len(pos) > 1 {
		printError(stderr, "usage: dh harness-path [dir] [--json]")
		return 1
	}
	dir := "."
	if len(pos) == 1 {
		dir = pos[0]
	}
	repo, err := filepath.Abs(dir)
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	if st, err := os.Stat(repo); err != nil || !st.IsDir() {
		printError(stderr, repo+" is not a directory")
		return 1
	}
	home, err := harness.EnsureHome()
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	c, err := harness.LoadConfig(home)
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	r := harness.Resolve(repo, c)
	if !f["--json"] {
		printResolved(stdout, r)
		return 0
	}
	type defaults struct {
		Language  string `json:"language"`
		Mode      string `json:"mode"`
		Reviewers string `json:"reviewers"`
	}
	key := repo
	if cr, err := filepath.EvalSymlinks(repo); err == nil {
		key = cr
	}
	writeJSON(stdout, struct {
		harness.Resolved
		Pdocs    string   `json:"pdocs"`
		Defaults defaults `json:"defaults"`
	}{r, harness.PdocsDir(home, c.Projects[key]), defaults{c.Language, c.Mode, c.ReviewersRaw}})
	return 0
}

// runLink: with one positional, a value with a path separator or "."/".." is the directory,
// anything else is the name; with two, they are name then dir.
func runLink(args []string, stdout, stderr io.Writer) int {
	usageMsg := "usage: dh link [name] [--off] [dir] [--json]"
	pos, f, ok := splitArgs(args, "--off", "--json")
	if !ok || len(pos) > 2 {
		printError(stderr, usageMsg)
		return 1
	}
	name, dir := "", "."
	switch len(pos) {
	case 2:
		name, dir = pos[0], pos[1]
	case 1:
		if pos[0] == "." || pos[0] == ".." || strings.ContainsAny(pos[0], `/\`) {
			dir = pos[0]
		} else {
			name = pos[0]
		}
	}
	home, err := harness.EnsureHome()
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	r, err := harness.Link(home, dir, name, f["--off"])
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	if f["--json"] {
		writeJSON(stdout, r)
	} else {
		printResolved(stdout, r)
	}
	return 0
}

func runProjects(args []string, stdout, stderr io.Writer) int {
	pos, f, ok := splitArgs(args, "--json")
	if !ok || len(pos) != 0 {
		printError(stderr, "usage: dh projects [--json]")
		return 1
	}
	home, err := harness.EnsureHome()
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	if _, err := harness.LoadConfig(home); err != nil {
		printError(stderr, err.Error())
		return 1
	}
	list, warns := harness.Projects(home)
	for _, w := range warns {
		printError(stderr, "warning: "+w)
	}
	if f["--json"] {
		if list == nil {
			list = []harness.Project{}
		}
		writeJSON(stdout, list)
		return 0
	}
	for _, p := range list {
		fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", p.Name, p.Mode, p.Path, p.Harness)
	}
	return 0
}
