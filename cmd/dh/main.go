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
			printError(stderr, "uso: dh doctor [pasta-do-plugin-ou-raiz-do-kit]")
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
	case "config":
		return runConfig(args[1:], stdout, stderr)
	case "memory":
		return runMemory(args[1:], stdin, stdout, stderr)
	case "facts":
		return runFacts(args[1:], stdout, stderr)
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
		fmt.Fprintf(stdout, "removidos %d\n", count)
		return 0
	default:
		printError(stderr, "subcomando de snapshot desconhecido: "+args[1])
		return 1
	}
}

func runBuild(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("dh build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	targets := flags.String("targets", "", "alvos GOOS/GOARCH separados por vírgula")
	noBinaries := flags.Bool("no-binaries", false, "não gera os binários multiplataforma")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		if err == nil {
			printError(stderr, "uso: dh build [--targets os/arch,...] [--no-binaries]")
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
			return nil, fmt.Errorf("alvo inválido: %s", part)
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
				printError(stderr, "uso: dh validate [--source-only] [caminho]")
				return 1
			}
			sourceOnly = true
			continue
		}
		if targetSet {
			printError(stderr, "uso: dh validate [--source-only] [caminho]")
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

// runMemory is the hash guard of the records MEMORY.md and FACTS.md (PRD-034, PRD-035): `hash` prints the sha256 (or "absent"); `write` is a compare-and-swap.
// --file picks the record from the closed list (MEMORY.md, the default, or FACTS.md); any other name exits 1 and lists both.
// Exit codes: 0 ok, 1 usage or I/O error, 3 the file changed since it was read (nothing written).
// The tokens `absent` and `written <hash>` are a contract the kit skills document, so they stay as they are; only the errors are in pt-BR.
func runMemory(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	const memUsage = "uso: dh memory hash [--file MEMORY.md|FACTS.md] [pasta] | dh memory write --expect <sha256|absent> [--file MEMORY.md|FACTS.md] [--from <arquivo>] [pasta]"
	if len(args) == 0 {
		printError(stderr, memUsage)
		return 1
	}
	switch args[0] {
	case "hash":
		flags := flag.NewFlagSet("dh memory hash", flag.ContinueOnError)
		flags.SetOutput(stderr)
		file := flags.String("file", harness.MemoryFile, "registro: MEMORY.md (padrão) ou FACTS.md")
		if err := flags.Parse(args[1:]); err != nil || flags.NArg() > 1 {
			if err == nil {
				printError(stderr, memUsage)
			}
			return 1
		}
		if err := harness.CheckRecordFile(*file); err != nil {
			printError(stderr, err.Error())
			return 1
		}
		dir, ok := memoryDir(flags.Args(), stderr)
		if !ok {
			return 1
		}
		sum, err := harness.FileHash(dir, *file)
		if err != nil {
			printError(stderr, err.Error())
			return 1
		}
		fmt.Fprintln(stdout, sum)
		return 0
	case "write":
		flags := flag.NewFlagSet("dh memory write", flag.ContinueOnError)
		flags.SetOutput(stderr)
		expect := flags.String("expect", "", "sha256 do registro como foi lido, ou absent")
		file := flags.String("file", harness.MemoryFile, "registro: MEMORY.md (padrão) ou FACTS.md")
		from := flags.String("from", "", "arquivo com o novo conteúdo (padrão: entrada padrão)")
		if err := flags.Parse(args[1:]); err != nil || *expect == "" || flags.NArg() > 1 {
			if err == nil {
				printError(stderr, memUsage)
			}
			return 1
		}
		if err := harness.CheckRecordFile(*file); err != nil {
			printError(stderr, err.Error())
			return 1
		}
		dir, ok := memoryDir(flags.Args(), stderr)
		if !ok {
			return 1
		}
		var content []byte
		var err error
		if *from != "" {
			content, err = os.ReadFile(*from)
		} else {
			content, err = io.ReadAll(stdin)
		}
		if err != nil {
			printError(stderr, err.Error())
			return 1
		}
		sum, err := harness.WriteFile(dir, *file, *expect, content)
		if harness.IsMemoryChanged(err) {
			printError(stderr, err.Error())
			return 3
		}
		if err != nil {
			printError(stderr, err.Error())
			return 1
		}
		fmt.Fprintln(stdout, "written "+sum)
		return 0
	default:
		printError(stderr, memUsage)
		return 1
	}
}

// runFacts is `dh facts check [pasta]` (PRD-035 R27): it resolves the harness dir like `dh harness-path`, checks <harness dir>/FACTS.md against R3 to R6 and that
// every evidence path exists under the project folder, and writes nothing. Exit codes: 0 the file is fine or absent, 1 a violation, mode none or an I/O error.
// The tokens `ok` and `absent` on stdout are a contract the kit documents.
func runFacts(args []string, stdout, stderr io.Writer) int {
	const factsUsage = "uso: dh facts check [pasta]"
	if len(args) == 0 || args[0] != "check" || len(args) > 2 {
		printError(stderr, factsUsage)
		return 1
	}
	repo, dir, ok := projectHarness(args[1:], stderr)
	if !ok {
		return 1
	}
	content, err := os.ReadFile(filepath.Join(dir, harness.FactsFile))
	if os.IsNotExist(err) {
		fmt.Fprintln(stdout, "absent")
		return 0
	}
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	rep := harness.CheckFacts(content, func(rel string) bool {
		_, err := os.Stat(filepath.Join(repo, filepath.FromSlash(rel)))
		return err == nil
	})
	if len(rep.Problems) > 0 {
		for _, p := range rep.Problems {
			if p.Line > 0 {
				fmt.Fprintf(stderr, "%s:%d: %s\n", harness.FactsFile, p.Line, p.Msg)
			} else {
				fmt.Fprintf(stderr, "%s: %s\n", harness.FactsFile, p.Msg)
			}
		}
		printError(stderr, fmt.Sprintf("%s reprovado: %d problema(s)", harness.FactsFile, len(rep.Problems)))
		return 1
	}
	fmt.Fprintf(stdout, "ok %d fatos %d candidatos\n", rep.Facts, rep.Candidates)
	return 0
}

// memoryDir resolves the harness dir of the optional project dir argument, like `dh harness-path`; mode none is an error.
func memoryDir(pos []string, stderr io.Writer) (string, bool) {
	_, dir, ok := projectHarness(pos, stderr)
	return dir, ok
}

// projectHarness resolves the project folder (the optional argument, else the current one) and its harness dir; mode none is an error.
func projectHarness(pos []string, stderr io.Writer) (repo, dir string, ok bool) {
	arg := "."
	if len(pos) == 1 {
		arg = pos[0]
	}
	repo, err := filepath.Abs(arg)
	if err != nil {
		printError(stderr, err.Error())
		return "", "", false
	}
	if st, err := os.Stat(repo); err != nil || !st.IsDir() {
		printError(stderr, repo+" não é um diretório")
		return "", "", false
	}
	home, err := harness.EnsureHome()
	if err != nil {
		printError(stderr, err.Error())
		return "", "", false
	}
	c, err := harness.LoadConfig(home)
	if err != nil {
		printError(stderr, err.Error())
		return "", "", false
	}
	r := harness.Resolve(repo, c)
	if r.Mode == "none" || r.Dir == "" {
		printError(stderr, "não há harness dir para "+repo+" (mode none): rode /dh:setup antes")
		return "", "", false
	}
	return repo, r.Dir, true
}

func usage() string {
	return "uso: dh <memory <hash|write>|facts check [pasta]|snapshot <event|subagents|statusline|prune>|validate [--source-only] [caminho]|build [--targets os/arch,...] [--no-binaries]|doctor [pasta-do-plugin-ou-raiz-do-kit]|dashboard [--port N] [--stale D] [--done-decay D] [--detach]|harness-path [pasta] [--json]|link [nome] [--off] [pasta] [--json]|projects [--json]|config [--json]|config set <chave> <valor>|config unset <chave>>"
}

func pruneDays(args []string) (int, error) {
	if len(args) == 0 {
		return 7, nil
	}
	if len(args) != 2 || args[0] != "--days" {
		return 0, fmt.Errorf("uso: dh snapshot prune [--days N]")
	}
	days, err := strconv.Atoi(args[1])
	if err != nil || days < 0 {
		return 0, fmt.Errorf("dias inválidos: %s", args[1])
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
		printError(stderr, "uso: dh harness-path [pasta] [--json]")
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
		printError(stderr, repo+" não é um diretório")
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
		Mode      string `json:"mode"`
		Reviewers string `json:"reviewers"`
		Worktree  string `json:"worktree"`
	}
	key := repo
	if cr, err := filepath.EvalSymlinks(repo); err == nil {
		key = cr
	}
	lw := harness.LinkedWorktree(repo)
	if lw != nil {
		key = lw.Main
	}
	writeJSON(stdout, struct {
		harness.Resolved
		Pdocs          string            `json:"pdocs"`
		Defaults       defaults          `json:"defaults"`
		LinkedWorktree *harness.LinkedWT `json:"linked_worktree"`
	}{r, harness.PdocsDir(home, c.Projects[key]), defaults{c.Mode, c.ReviewersRaw, worktreeOr(c.Worktree)}, lw})
	return 0
}

// worktreeOr is the effective global worktree policy: the stored value, or "on-request" when absent.
func worktreeOr(v string) string {
	if v == "" {
		return "on-request"
	}
	return v
}

// runLink: with one positional, a value with a path separator or "."/".." is the directory,
// anything else is the name; with two, they are name then dir.
func runLink(args []string, stdout, stderr io.Writer) int {
	usageMsg := "uso: dh link [nome] [--off] [pasta] [--json]"
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
	if d := harness.LinkedDir(home, dir); r.Mode == "none" && d != "" {
		fmt.Fprintf(stderr, "aviso: link gravado, mas %s ainda não tem marcador de harness (project.yaml, MEMORY.md, EPOCHAL.md, RISKS.md, tasks/ ou prd/); rode /dh:setup para criar o harness\n", d)
	}
	return 0
}

func runProjects(args []string, stdout, stderr io.Writer) int {
	pos, f, ok := splitArgs(args, "--json")
	if !ok || len(pos) != 0 {
		printError(stderr, "uso: dh projects [--json]")
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
		printError(stderr, "aviso: "+w)
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

func runConfig(args []string, stdout, stderr io.Writer) int {
	pos, f, ok := splitArgs(args, "--json")
	if !ok || len(pos) > 3 || (len(pos) > 0 && (pos[0] != "set" || len(pos) != 3) && (pos[0] != "unset" || len(pos) != 2)) {
		printError(stderr, "uso: dh config [--json] | dh config set <chave> <valor> | dh config unset <chave>")
		return 1
	}
	home, err := harness.EnsureHome()
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	if len(pos) > 0 {
		value := ""
		if pos[0] == "set" {
			value = pos[2]
		}
		if err := harness.SetConfigValue(home, pos[1], value, pos[0] == "unset"); err != nil {
			printError(stderr, err.Error())
			return 1
		}
		return 0
	}
	c, err := harness.LoadConfig(home)
	if err != nil {
		printError(stderr, err.Error())
		return 1
	}
	if f["--json"] {
		writeJSON(stdout, struct {
			Mode     string `json:"mode"`
			Sprites  string `json:"dashboard.sprites"`
			Worktree string `json:"worktree"`
		}{c.Mode, c.Sprites, c.Worktree})
		return 0
	}
	fmt.Fprintf(stdout, "mode: %s\ndashboard.sprites: %s\nworktree: %s\n", c.Mode, c.Sprites, c.Worktree)
	return 0
}
