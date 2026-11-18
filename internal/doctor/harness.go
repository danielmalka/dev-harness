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

// printHarness reports the resolved harness for the current directory (PRD-014 R14, R9, R6) and returns how many problems fail the diagnosis
// (today, only an invalid `conventions.code_language`). Read-only: ~/.claude/settings.json is only ever read.
func printHarness(out io.Writer) int {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(out, "harness: indisponível")
		return 0
	}
	var settings string
	if h, err := os.UserHomeDir(); err == nil {
		settings = filepath.Join(h, ".claude", "settings.json")
	}
	return printHarnessFor(out, cwd, harness.Home(), settings)
}

func printHarnessFor(out io.Writer, cwd, home, settings string) int {
	repo := absolutePath(cwd)
	if r, err := filepath.EvalSymlinks(repo); err == nil {
		repo = r
	}
	cfg, cfgErr := harness.LoadConfig(home)
	if cfgErr != nil {
		fmt.Fprintf(out, "aviso: %v\n", cfgErr)
	}
	res := harness.Resolve(repo, cfg)
	if home == "" {
		fmt.Fprintln(out, "home: desconhecido")
		fmt.Fprintln(out, "dica: não foi possível determinar a pasta do harness: defina DH_HOME, HOME ou USERPROFILE")
	} else {
		fmt.Fprintf(out, "home: %s\n", home)
	}
	fmt.Fprintf(out, "mode: %s\n", res.Mode)
	if res.Mode != "none" {
		fmt.Fprintf(out, "dir: %s\n", res.Dir)
		for _, name := range []string{"project.yaml", "MEMORY.md", "EPOCHAL.md", "RISKS.md"} {
			state := "ausente"
			if isFile(filepath.Join(res.Dir, name)) {
				state = "presente"
			}
			fmt.Fprintf(out, "%s: %s\n", name, state)
			if name == "MEMORY.md" && state == "presente" {
				printMemorySize(out, filepath.Join(res.Dir, name))
			}
		}
		if isFile(filepath.Join(res.Dir, harness.FactsFile)) {
			fmt.Fprintf(out, "%s: presente\n", harness.FactsFile)
			printFactsSize(out, filepath.Join(res.Dir, harness.FactsFile))
		}
	}
	problems := printLegacyAndConventions(out, home, res)
	name, registered := cfg.Projects[repo]
	if res.Mode == "repo" {
		if registered && harness.IsHarness(filepath.Join(home, "projects", name)) {
			fmt.Fprintf(out, "aviso: %s e %s existem os dois: vale o interno\n",
				res.Dir, filepath.Join(home, "projects", name))
		}
		if !registered {
			fmt.Fprintf(out, "aviso: %s não está registrado em %s: não aparece no dashboard até o dh link\n",
				repo, filepath.Join(home, "config.yaml"))
		}
		if isFile(filepath.Join(res.Dir, "local.yaml")) && notIgnored(repo, ".harness/local.yaml") {
			fmt.Fprintln(out, "aviso: .harness/local.yaml não está ignorado pelo git")
		}
	}
	if res.Mode == "none" {
		fmt.Fprintln(out, "dica: não há harness para este diretório: rode /dh:setup (ou dh link)")
	}
	if home != "" && !within(repo, home) {
		if msg := settingsHint(settings, home); msg != "" {
			fmt.Fprintln(out, msg)
		}
	}
	return problems
}

// printLegacyAndConventions warns about the legacy `language:` line (config.yaml in home and project.yaml in the harness dir: read, ignored, never rewritten, exit unchanged)
// and fails an invalid `conventions.code_language`. It returns how many problems fail.
func printLegacyAndConventions(out io.Writer, home string, res harness.Resolved) int {
	legacyHint := func(file string) {
		fmt.Fprintf(out, "aviso: `language` não existe mais e é ignorada (o kit é só pt-BR); pode apagar a linha de %s\n", file)
	}
	if home != "" {
		path := filepath.Join(home, "config.yaml")
		if data, err := os.ReadFile(path); err == nil && harness.HasLegacyLanguage(string(data)) {
			legacyHint(path)
		}
	}
	if res.Dir == "" {
		return 0
	}
	path := filepath.Join(res.Dir, "project.yaml")
	check, err := harness.ReadProjectCheck(path)
	if err != nil {
		fmt.Fprintf(out, "aviso: não foi possível ler %s: %v\n", path, err)
		return 0
	}
	if check.LegacyLanguage {
		legacyHint(path)
	}
	for _, problem := range check.Problems {
		fmt.Fprintf(out, "erro: %s em %s\n", problem, path)
	}
	return len(check.Problems)
}

// memoryBudget is the MEMORY.md file budget in bytes (PRD-032 R7, R8).
const memoryBudget = 40000

func printMemorySize(out io.Writer, path string) {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(out, "aviso: não foi possível ler o tamanho do MEMORY.md: %v\n", err)
		return
	}
	fmt.Fprintf(out, "tamanho do MEMORY.md: %d bytes\n", fi.Size())
	if fi.Size() > memoryBudget {
		fmt.Fprintf(out, "aviso: o MEMORY.md tem %d bytes, acima do orçamento de %d: rode /dh:consolidate-memory\n", fi.Size(), memoryBudget)
	}
}

// printFactsSize reports the size of FACTS.md and warns when it passes the ceilings of PRD-035 R5 (150 lines, 12 KB), like printMemorySize does for MEMORY.md.
// The file is listed only when it exists: a project that does not use it yet gets no line.
func printFactsSize(out io.Writer, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(out, "aviso: não foi possível ler o tamanho do %s: %v\n", harness.FactsFile, err)
		return
	}
	lines := 0
	if len(data) > 0 {
		lines = strings.Count(strings.TrimSuffix(string(data), "\n"), "\n") + 1
	}
	fmt.Fprintf(out, "tamanho do %s: %d linhas, %d bytes\n", harness.FactsFile, lines, len(data))
	if lines > harness.FactsMaxLines {
		fmt.Fprintf(out, "aviso: o %s tem %d linhas, acima do teto de %d: peça ao Coordenador a curadoria do arquivo\n", harness.FactsFile, lines, harness.FactsMaxLines)
	}
	if len(data) > harness.FactsMaxBytes {
		fmt.Fprintf(out, "aviso: o %s tem %d bytes, acima do teto de %d (12 KB): peça ao Coordenador a curadoria do arquivo\n", harness.FactsFile, len(data), harness.FactsMaxBytes)
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
		return fmt.Sprintf("dica: ~/.claude/settings.json %s que cubra %s; acrescente %q (o dh nunca escreve esse arquivo)", why, home, suggest)
	}
	missing := hint("não tem entrada em permissions.additionalDirectories")
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
		return hint("não pôde ser lido, então não se achou entrada em permissions.additionalDirectories")
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
