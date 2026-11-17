package doctor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/danielmalka/dev-harness/internal/kit"
	"github.com/danielmalka/dev-harness/internal/snapshot"
)

var snapshotMention = regexp.MustCompile(`bin/dh(\.cmd)?\\?"? +snapshot`)

// Run prints the read-only diagnosis and returns the exit status (validator and harness).
func Run(target string, out io.Writer) int {
	if target == "" {
		target = "."
	}
	targetForDisplay := absolutePath(target)

	fmt.Fprintln(out, "## Modo")
	fmt.Fprintln(out, "diagnóstico (somente leitura; cria <home> e <home>/sessions só quando faltam)")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Kit")
	fmt.Fprintf(out, "alvo: %s\n", targetForDisplay)
	fmt.Fprintln(out, "validador: dh validate embutido")

	report, err := kit.Validate(target, kit.Options{})
	validateStatus := 0
	if errors.Is(err, kit.ErrUnrecognizedLayout) {
		fmt.Fprintln(out, "validação do kit: não se aplica (não é um checkout do kit)")
	} else if err != nil {
		fmt.Fprintln(out, "a validação do kit falhou:", err)
		validateStatus = 1
	} else {
		kit.WriteReport(out, report)
		if !report.Passed() {
			validateStatus = 1
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Ambiente")
	fmt.Fprintln(out, "| Camada | Checagem | Comando ou observação | Estado | Evidência |")
	fmt.Fprintln(out, "|---|---|---|---|---|")
	probe(out, "runtime", "python3", "python3", "--version")
	probe(out, "runtime", "git", "git", "--version")
	probe(out, "runtime", "claude", "claude", "--version")

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Harness")
	harnessProblems := printHarness(out)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Limites")
	fmt.Fprintln(out, "- Nada foi instalado, atualizado nem configurado.")
	fmt.Fprintln(out, "- A descoberta de comandos e agentes pelo Claude Code não foi verificada, a menos que esta sessão tenha carregado o plugin.")
	fmt.Fprintln(out, "- A falta do binário claude impede carregar o plugin; não impede ler a documentação.")
	printRF07(out, targetForDisplay)

	switch {
	case validateStatus != 0:
		fmt.Fprintln(out, "status: reprovado (erros do validador acima)")
	case harnessProblems != 0:
		fmt.Fprintln(out, "status: reprovado (erros do harness acima)")
	default:
		fmt.Fprintln(out, "status: aprovado")
	}
	if harnessProblems != 0 {
		return 1
	}
	return validateStatus
}

func probe(out io.Writer, layer, check string, args ...string) {
	command := strings.Join(args, " ")
	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(out, "| %s | %s | %s | ausente | comando não encontrado |\n", layer, check, command)
		return
	}
	result := exec.Command(path, args[1:]...)
	data, err := result.CombinedOutput()
	if err != nil {
		fmt.Fprintf(out, "| %s | %s | %s | desconhecido | a sonda falhou |\n", layer, check, command)
		return
	}
	trimmed := strings.ReplaceAll(strings.TrimSpace(string(data)), "\n", " ")
	if len(trimmed) > 120 {
		trimmed = trimmed[:120]
	}
	fmt.Fprintf(out, "| %s | %s | %s | presente | %s |\n", layer, check, command, trimmed)
}

func printRF07(out io.Writer, target string) {
	packageRoot := target
	if !isDir(filepath.Join(packageRoot, "agents")) || !isDir(filepath.Join(packageRoot, ".claude-plugin")) {
		candidate := filepath.Join(target, "dist", "claude-code", "dev-harness")
		if isDir(filepath.Join(candidate, "agents")) && isDir(filepath.Join(candidate, ".claude-plugin")) {
			packageRoot = candidate
		} else {
			packageRoot = ""
		}
	}
	if packageRoot != "" {
		platform := runtime.GOOS + "_" + runtime.GOARCH
		binary := filepath.Join(packageRoot, "bin", platform, "dh")
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		if isFile(binary) {
			fmt.Fprintf(out, "- binário %s: presente (%s)\n", platform, binary)
		} else {
			fmt.Fprintf(out, "- binário %s: ausente (%s)\n", platform, binary)
		}
		printPluginFile(out, "hooks", filepath.Join(packageRoot, "hooks", "hooks.json"))
		printPluginFile(out, "settings", filepath.Join(packageRoot, "settings.json"))
	}

	snapshotDir := snapshot.SnapshotDir()
	if snapshotDir == "" {
		fmt.Fprintln(out, "- pasta de snapshots: ignorada (home do harness desconhecido)")
	} else if err := probeSnapshotDir(snapshotDir); err != nil {
		fmt.Fprintf(out, "- pasta de snapshots: sem permissão de escrita (%s: %v)\n", snapshotDir, err)
	} else {
		fmt.Fprintf(out, "- pasta de snapshots: gravável (%s)\n", snapshotDir)
	}
}

func printPluginFile(out io.Writer, name, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(out, "- %s: ausente (%s)\n", name, path)
		return
	}
	if snapshotMention.Match(data) {
		fmt.Fprintf(out, "- %s: presente e cita dh snapshot (%s)\n", name, path)
	} else {
		fmt.Fprintf(out, "- %s: presente, mas não cita dh snapshot (%s)\n", name, path)
	}
}

func probeSnapshotDir(directory string) error {
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, ".doctor-*")
	if err != nil {
		return err
	}
	name := temporary.Name()
	if err := temporary.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	return os.Remove(name)
}

func absolutePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return abs
	}
	return resolved
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}
