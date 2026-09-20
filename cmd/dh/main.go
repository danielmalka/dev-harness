package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/danielmalka/dev-harness/internal/build"
	"github.com/danielmalka/dev-harness/internal/doctor"
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
	return "usage: dh <snapshot <event|subagents|statusline|prune>|validate [--source-only] [path]|build [--targets os/arch,...] [--no-binaries]|doctor [plugin-dir-or-kit-root]>"
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
