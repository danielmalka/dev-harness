package doctor

import (
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

// Run prints the read-only diagnosis and returns the validator exit status.
func Run(target string, out io.Writer) int {
	if target == "" {
		target = "."
	}
	targetForDisplay := absolutePath(target)

	fmt.Fprintln(out, "## Mode")
	fmt.Fprintln(out, "diagnosis (no writes)")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Kit")
	fmt.Fprintf(out, "target: %s\n", targetForDisplay)
	fmt.Fprintln(out, "validator: embedded dh validate")

	report, err := kit.Validate(target, kit.Options{})
	validateStatus := 0
	if err != nil {
		fmt.Fprintln(out, err)
		validateStatus = 1
	} else {
		kit.WriteReport(out, report)
		if !report.Passed() {
			validateStatus = 1
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Environment")
	fmt.Fprintln(out, "| Layer | Check | Command or observation | Status | Evidence |")
	fmt.Fprintln(out, "|---|---|---|---|---|")
	probe(out, "runtime", "python3", "python3", "--version")
	probe(out, "runtime", "git", "git", "--version")
	probe(out, "runtime", "claude", "claude", "--version")

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Harness records in current directory")
	printHarnessRecords(out)

	fmt.Fprintln(out)
	fmt.Fprintln(out, "## Limits")
	fmt.Fprintln(out, "- Nothing was installed, upgraded or configured.")
	fmt.Fprintln(out, "- Claude Code discovery of commands and agents is unverified unless this session loaded the plugin.")
	fmt.Fprintln(out, "- A missing claude binary blocks loading the plugin; it does not block reading the docs.")
	printRF07(out, targetForDisplay)

	if validateStatus == 0 {
		fmt.Fprintln(out, "status: passed")
	} else {
		fmt.Fprintln(out, "status: failed (validator errors above)")
	}
	return validateStatus
}

func probe(out io.Writer, layer, check string, args ...string) {
	command := strings.Join(args, " ")
	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(out, "| %s | %s | %s | missing | command not found |\n", layer, check, command)
		return
	}
	result := exec.Command(path, args[1:]...)
	data, err := result.CombinedOutput()
	if err != nil {
		fmt.Fprintf(out, "| %s | %s | %s | unknown | probe failed |\n", layer, check, command)
		return
	}
	trimmed := strings.ReplaceAll(strings.TrimSpace(string(data)), "\n", " ")
	if len(trimmed) > 120 {
		trimmed = trimmed[:120]
	}
	fmt.Fprintf(out, "| %s | %s | %s | present | %s |\n", layer, check, command, trimmed)
}

func printHarnessRecords(out io.Writer) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(out, ".harness/: unavailable")
		return
	}
	harness := filepath.Join(cwd, ".harness")
	if !isDir(harness) {
		fmt.Fprintln(out, ".harness/: missing")
		return
	}
	fmt.Fprintln(out, ".harness/: present")
	for _, name := range []string{"project.yaml", "MEMORY.md", "EPOCHAL.md", "RISKS.md"} {
		if isFile(filepath.Join(harness, name)) {
			fmt.Fprintf(out, ".harness/%s: present\n", name)
		} else {
			fmt.Fprintf(out, ".harness/%s: missing\n", name)
		}
	}
	if ignoredByGit(cwd) {
		fmt.Fprintln(out, "warning: .harness/ is ignored by git")
	}
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
			fmt.Fprintf(out, "- binary %s: present (%s)\n", platform, binary)
		} else {
			fmt.Fprintf(out, "- binary %s: missing (%s)\n", platform, binary)
		}
		printPluginFile(out, "hooks", filepath.Join(packageRoot, "hooks", "hooks.json"))
		printPluginFile(out, "settings", filepath.Join(packageRoot, "settings.json"))
	}

	snapshotDir := snapshot.SnapshotDir()
	if err := probeSnapshotDir(snapshotDir); err != nil {
		fmt.Fprintf(out, "- snapshot directory: not writable (%s: %v)\n", snapshotDir, err)
	} else {
		fmt.Fprintf(out, "- snapshot directory: writable (%s)\n", snapshotDir)
	}
	if ignored, _ := harnessIgnored(); ignored {
		fmt.Fprintln(out, "- warning: .harness/ is ignored by git")
	}
}

func printPluginFile(out io.Writer, name, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(out, "- %s: missing (%s)\n", name, path)
		return
	}
	if snapshotMention.Match(data) {
		fmt.Fprintf(out, "- %s: present and mentions dh snapshot (%s)\n", name, path)
	} else {
		fmt.Fprintf(out, "- %s: present but does not mention dh snapshot (%s)\n", name, path)
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

func ignoredByGit(directory string) bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	command := exec.Command(git, "check-ignore", "-q", ".harness")
	command.Dir = directory
	return command.Run() == nil
}

func harnessIgnored() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil || !isDir(filepath.Join(cwd, ".harness")) {
		return false, err
	}
	return ignoredByGit(cwd), nil
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
