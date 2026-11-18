package build

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

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
		"templates/TASK.md",
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
	if !strings.Contains(log.String(), "binários      0") {
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
		"evals/archive/old",
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
	writeBuildFixture(t, root, "evals/archive/old/case.yaml", "id: old\n", 0o644)
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
		"evals/archive",
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
		"templates", "adapters/claude-code/plugin/hooks",
		"adapters/claude-code/plugin/bin",
	} {
		if err := os.MkdirAll(filepath.Join(root, path), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeBuildFixture(t, root, ".agents/alpha.md", "---\nname: alpha\ndescription: Use quando testando\nauthor: test\nmodel: haiku\ntools:\n  - Read\n---\nagent\n", 0o644)
	writeBuildFixture(t, root, ".commands/run.md", "---\nname: run\ndescription: Run test\nauthor: test\nargument-hint: \"\"\nmetadata:\n  roles: [alpha]\n  skills: [sample]\n  writes: none\n---\ncommand\n", 0o644)
	writeBuildFixture(t, root, ".skills/sample/SKILL.md", "---\nname: sample\ndescription: Use quando testando\nauthor: test\nmetadata:\n  provenance: test\n---\nskill\n", 0o644)
	writeBuildFixture(t, root, "profiles/base.yaml", "id: base\n", 0o644)
	writeBuildFixture(t, root, "templates/TASK.md", "# Task\n", 0o644)
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

func TestBuildKeepsManifestStampUntilSourcesChange(t *testing.T) {
	root := minimalRoot(t)
	output := filepath.Join(t.TempDir(), "dev-harness")
	options := Options{Output: output, NoBinaries: true, SkipMinimumCounts: true}
	read := func() string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(output, "harness-manifest.json"))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	first := read()
	var parsed manifest
	if err := json.Unmarshal([]byte(first), &parsed); err != nil {
		t.Fatal(err)
	}
	// Age the stamp: an unchanged tree must keep it, so the time source cannot leak in.
	aged := strings.Replace(first, parsed.GeneratedAt, "2000-01-01T00:00:00Z", 1)
	if err := os.WriteFile(filepath.Join(output, "harness-manifest.json"), []byte(aged), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	if got := read(); got != aged {
		t.Fatalf("manifest changed on an unchanged tree:\n%s", got)
	}
	skill := filepath.Join(root, ".skills", "sample", "SKILL.md")
	body, err := os.ReadFile(skill)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, append(body, []byte("\nChanged.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	changed := read()
	if changed == aged || strings.Contains(changed, "2000-01-01T00:00:00Z") {
		t.Fatalf("editing a skill did not refresh the manifest:\n%s", changed)
	}
}

func TestPackageDigestIgnoresCRLFInTextButNotInBinaries(t *testing.T) {
	pkg := t.TempDir()
	write := func(name string, data []byte) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(pkg, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	digest := func() string {
		t.Helper()
		sum, err := packageDigest(pkg)
		if err != nil {
			t.Fatal(err)
		}
		return sum
	}
	write("SKILL.md", []byte("---\nname: sample\n---\nBody line.\n"))
	write("blob.dat", []byte{0, 1, '\r', '\n', 2})
	lf := digest()
	write("SKILL.md", []byte("---\r\nname: sample\r\n---\r\nBody line.\r\n"))
	if got := digest(); got != lf {
		t.Fatalf("CRLF text changed the digest: %s != %s", got, lf)
	}
	write("SKILL.md", []byte("---\nname: sample\n---\nBody line changed.\n"))
	if digest() == lf {
		t.Fatal("changed text kept the digest")
	}
	write("SKILL.md", []byte("---\nname: sample\n---\nBody line.\n"))
	write("blob.dat", []byte{0, 1, '\n', 2})
	if digest() == lf {
		t.Fatal("binary content with \\r\\n removed kept the digest")
	}
}

func TestBuildEditAndUndoGoesBackToCommittedStamp(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := minimalRoot(t)
	output := filepath.Join(t.TempDir(), "dev-harness")
	options := Options{Output: output, NoBinaries: true, SkipMinimumCounts: true, UpdateRootManifest: true}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	rootManifest := filepath.Join(root, "harness-manifest.json")
	data, err := os.ReadFile(rootManifest)
	if err != nil {
		t.Fatal(err)
	}
	var parsed manifest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}
	const committed = "2000-01-01T00:00:00Z"
	if err := os.WriteFile(rootManifest, []byte(strings.Replace(string(data), parsed.GeneratedAt, committed, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.name=t", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false"}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	git("add", "harness-manifest.json")
	git("commit", "-q", "-m", "manifest")
	stamp := func() string {
		t.Helper()
		data, err := os.ReadFile(rootManifest)
		if err != nil {
			t.Fatal(err)
		}
		var m manifest
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatal(err)
		}
		return m.GeneratedAt
	}
	skill := filepath.Join(root, ".skills", "sample", "SKILL.md")
	original, err := os.ReadFile(skill)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, append(original, []byte("edit\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	if stamp() == committed {
		t.Fatal("editing a skill kept the committed stamp")
	}
	if err := os.WriteFile(skill, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, options, io.Discard); err != nil {
		t.Fatal(err)
	}
	if got := stamp(); got != committed {
		t.Fatalf("undoing the edit gave stamp %s, want the committed %s", got, committed)
	}
}

func TestStableGeneratedAtWithoutGitFallsBackToNow(t *testing.T) {
	t.Setenv("PATH", "")
	got := stableGeneratedAt(manifest{}, t.TempDir(), filepath.Join(t.TempDir(), "none.json"))
	if _, err := time.Parse("2006-01-02T15:04:05Z", got); err != nil {
		t.Fatalf("stamp %q is not a timestamp: %v", got, err)
	}
}
