package build

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/danielmalka/dev-harness/internal/kit"
)

func TestBuildNoBinariesCreatesValidatedPackage(t *testing.T) {
	root := minimalRoot(t)
	output := filepath.Join(t.TempDir(), "dev-harness")
	var log bytes.Buffer
	if err := Build(root, Options{
		Output:            output,
		NoBinaries:        true,
		SkipMinimumCounts: true,
	}, &log); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"agents/alpha.md",
		"skills/sample/SKILL.md",
		"commands/run.md",
		"templates/en/TASK.md",
		"templates/pt-br/TASK.md",
		"profiles/base.yaml",
		"settings.json",
		"hooks/hooks.json",
		"bin/dh",
		".claude-plugin/plugin.json",
		"harness-manifest.json",
		"GENERATED.txt",
	} {
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(path))); err != nil {
			t.Errorf("missing package file %s: %v", path, err)
		}
	}
	info, err := os.Stat(filepath.Join(output, "bin", "dh"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatal("bin/dh is not executable")
	}
	var data map[string]any
	manifestBytes, err := os.ReadFile(filepath.Join(output, "harness-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(manifestBytes, &data); err != nil {
		t.Fatal(err)
	}
	if data["source"] != "dh build" {
		t.Fatalf("source = %v, want dh build", data["source"])
	}
	if binaries, ok := data["binaries"].([]any); !ok || len(binaries) != 0 {
		t.Fatalf("binaries = %v, want empty array", data["binaries"])
	}
	if !strings.Contains(log.String(), "binaries     0") {
		t.Fatalf("build summary missing binary count: %s", log.String())
	}
}

func TestBuildValidatorFailureLeavesPreviousPackage(t *testing.T) {
	root := minimalRoot(t)
	if err := os.Remove(filepath.Join(root, "adapters", "claude-code", "plugin", "settings.json")); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "dev-harness")
	if err := os.MkdirAll(output, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(output, "sentinel")
	if err := os.WriteFile(sentinel, []byte("previous"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, Options{
		Output:            output,
		NoBinaries:        true,
		SkipMinimumCounts: true,
	}, &bytes.Buffer{}); err == nil {
		t.Fatal("Build succeeded with missing plugin setting")
	}
	data, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "previous" {
		t.Fatalf("previous package changed: %q", data)
	}
}

func TestRepositoryBuildNoBinariesPassesValidator(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	output := filepath.Join(t.TempDir(), "dev-harness")
	if err := Build(root, Options{
		Output:     output,
		NoBinaries: true,
	}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	report, err := kit.Validate(output, kit.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed() {
		t.Fatalf("package validation failed: %v", report.Errors)
	}
}

func TestBuildCopiesEvalsWithoutResultsBaselinesNotRun(t *testing.T) {
	root := minimalRoot(t)
	for _, path := range []string{
		"evals/cases/x/fixtures/results",
		"evals/results",
		"evals/baselines",
		"evals/not-run",
		"evals/fixtures/slice-01/__pycache__",
	} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(path)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeBuildFixture(t, root, "evals/cases/x/case.yaml", "id: x\n", 0o644)
	writeBuildFixture(t, root, "evals/cases/x/fixtures/results/kept.txt", "kept\n", 0o644)
	writeBuildFixture(t, root, "evals/results/run.json", "{}\n", 0o644)
	writeBuildFixture(t, root, "evals/baselines/base.json", "{}\n", 0o644)
	writeBuildFixture(t, root, "evals/not-run/skip.json", "{}\n", 0o644)
	writeBuildFixture(t, root, "evals/fixtures/slice-01/__pycache__/x.pyc", "cache\n", 0o644)
	writeBuildFixture(t, root, "evals/fixtures/slice-01/keep.py", "print(1)\n", 0o644)

	output := filepath.Join(t.TempDir(), "dev-harness")
	if err := Build(root, Options{
		Output:            output,
		NoBinaries:        true,
		SkipMinimumCounts: true,
	}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		"evals/cases/x/case.yaml",
		"evals/cases/x/fixtures/results/kept.txt",
		"evals/fixtures/slice-01/keep.py",
	} {
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(path))); err != nil {
			t.Errorf("missing package file %s: %v", path, err)
		}
	}
	for _, path := range []string{
		"evals/results",
		"evals/baselines",
		"evals/not-run",
		"evals/fixtures/slice-01/__pycache__",
	} {
		if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Errorf("excluded path present: %s", path)
		}
	}
}

func TestBuildWithoutEvalsDirectory(t *testing.T) {
	root := minimalRoot(t)
	output := filepath.Join(t.TempDir(), "dev-harness")
	if err := Build(root, Options{
		Output:            output,
		NoBinaries:        true,
		SkipMinimumCounts: true,
	}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(output, "evals")); !os.IsNotExist(err) {
		t.Errorf("evals should not exist in the package when absent from source")
	}
}

func minimalRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{
		".agents", ".commands", ".skills/sample", "profiles",
		"templates/en", "templates/pt-br", "adapters/claude-code/plugin/hooks",
		"adapters/claude-code/plugin/bin",
	} {
		if err := os.MkdirAll(filepath.Join(root, path), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeBuildFixture(t, root, ".agents/alpha.md", "---\nname: alpha\ndescription: Use when testing\nauthor: test\nmodel: haiku\ntools:\n  - Read\n---\nagent\n", 0o644)
	writeBuildFixture(t, root, ".commands/run.md", "---\nname: run\ndescription: Run test\nauthor: test\nargument-hint: \"\"\nmetadata:\n  roles: [alpha]\n  skills: [sample]\n  writes: none\n---\ncommand\n", 0o644)
	writeBuildFixture(t, root, ".skills/sample/SKILL.md", "---\nname: sample\ndescription: Use when testing\nauthor: test\nmetadata:\n  provenance: test\n---\nskill\n", 0o644)
	writeBuildFixture(t, root, "profiles/base.yaml", "id: base\n", 0o644)
	writeBuildFixture(t, root, "templates/en/TASK.md", "# Task\n", 0o644)
	writeBuildFixture(t, root, "templates/pt-br/TASK.md", "# Task\n", 0o644)
	writeBuildFixture(t, root, "adapters/claude-code/plugin/settings.json", "{}\n", 0o644)
	writeBuildFixture(t, root, "adapters/claude-code/plugin/hooks/hooks.json", "{}\n", 0o644)
	writeBuildFixture(t, root, "adapters/claude-code/plugin/bin/dh", "#!/bin/sh\n", 0o755)
	writeBuildFixture(t, root, "adapters/claude-code/plugin/bin/dh.cmd", "@echo off\n", 0o644)
	return root
}

func writeBuildFixture(t *testing.T, root, relative, contents string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}
