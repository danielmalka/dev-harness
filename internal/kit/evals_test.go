package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadFileGuardedRejectsSymlinkAndOversizedFile exercises the shared read
// guard directly: a symlink is rejected regardless of what it points at, an
// oversized file is rejected with a message naming the limit, and a normal
// file still reads through unchanged.
func TestReadFileGuardedRejectsSymlinkAndOversizedFile(t *testing.T) {
	root := t.TempDir()
	real := filepath.Join(root, "real.txt")
	writeFixture(t, real, "hello")
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	if _, err := readFileGuarded(link); err == nil {
		t.Fatal("want error for symlink")
	}

	big := filepath.Join(root, "big.txt")
	if err := os.WriteFile(big, make([]byte, maxReadableFileSize+1), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := readFileGuarded(big)
	if err == nil {
		t.Fatal("want error for oversized file")
	}
	if !strings.Contains(err.Error(), "read limit") {
		t.Fatalf("err = %v, want a message naming the read limit", err)
	}

	small := filepath.Join(root, "small.txt")
	writeFixture(t, small, "hello")
	data, err := readFileGuarded(small)
	if err != nil {
		t.Fatalf("unexpected error for a normal file: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("data = %q, want %q", data, "hello")
	}
}

// TestEvalFieldPathRejectsSymlinks covers MEDIUM 1: a case that declares a
// symlinked directory under context.add_dirs stays lexically inside the case
// directory, so containment alone would let it through and os.Stat would
// follow it. It must be rejected as a symlink, not as a lexical escape.
func TestEvalFieldPathRejectsSymlinks(t *testing.T) {
	root := syntheticRoot(t)
	caseDir := filepath.Join(root, "evals", "cases", "x")
	writeFixture(t, filepath.Join(caseDir, "case.yaml"), "context:\n  add_dirs:\n    - fixtures/outside\n")

	outside := filepath.Join(root, "outside-target")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, filepath.Join(outside, "secret.txt"), "top secret")

	if err := os.MkdirAll(filepath.Join(caseDir, "fixtures"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(caseDir, "fixtures", "outside")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "context.add_dirs is or passes through a symlink: fixtures/outside") {
		t.Fatalf("errors = %v, want symlink rejection", report.Errors)
	}
	if strings.Contains(joined, "escapes case directory") {
		t.Fatalf("errors = %v, want symlink rejection reported instead of a lexical-escape message", report.Errors)
	}
}

// TestEvalGraderTargetRejectsSymlink covers the same containment bypass for a
// grader's target.path.
func TestEvalGraderTargetRejectsSymlink(t *testing.T) {
	root := syntheticRoot(t)
	caseDir := filepath.Join(root, "evals", "cases", "x")
	writeFixture(t, filepath.Join(caseDir, "case.yaml"), "id: x\n")
	outside := filepath.Join(root, "outside.txt")
	writeFixture(t, outside, "secret")
	writeFixture(t, filepath.Join(caseDir, "graders", "g.md"), "---\ntarget:\n  source: file\n  path: linked.txt\n---\ngrader\n")
	if err := os.Symlink(outside, filepath.Join(caseDir, "linked.txt")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "target.path is or passes through a symlink: linked.txt") {
		t.Fatalf("errors = %v, want symlink rejection", report.Errors)
	}
}

// TestEvalCaseYAMLSymlinkIsReportedNotSkipped covers the direct-read half of
// MEDIUM 1: a symlinked case.yaml must surface as a validation error, not
// disappear as a silent skip.
func TestEvalCaseYAMLSymlinkIsReportedNotSkipped(t *testing.T) {
	root := syntheticRoot(t)
	caseDir := filepath.Join(root, "evals", "cases", "x")
	if err := os.MkdirAll(caseDir, 0o755); err != nil {
		t.Fatal(err)
	}
	real := filepath.Join(root, "real-case.yaml")
	writeFixture(t, real, "id: x\n")
	if err := os.Symlink(real, filepath.Join(caseDir, "case.yaml")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "case.yaml") || !strings.Contains(joined, "not a regular file") {
		t.Fatalf("errors = %v, want case.yaml symlink reported", report.Errors)
	}
}

// TestEvalGraderSymlinkIsReportedNotSkipped covers the same for a grader.
func TestEvalGraderSymlinkIsReportedNotSkipped(t *testing.T) {
	root := syntheticRoot(t)
	caseDir := filepath.Join(root, "evals", "cases", "x")
	writeFixture(t, filepath.Join(caseDir, "case.yaml"), "id: x\n")
	real := filepath.Join(root, "real-grader.md")
	writeFixture(t, real, "grader\n")
	if err := os.MkdirAll(filepath.Join(caseDir, "graders"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(caseDir, "graders", "g.md")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "graders/g.md") || !strings.Contains(joined, "not a regular file") {
		t.Fatalf("errors = %v, want grader symlink reported", report.Errors)
	}
}

// TestEvalPromptSymlinkIsReportedNotSkipped covers the prompt.md read inside
// checkEmbeddedAgentBodies.
func TestEvalPromptSymlinkIsReportedNotSkipped(t *testing.T) {
	root := syntheticRoot(t)
	real := filepath.Join(root, "real-prompt.md")
	writeFixture(t, real, "---\nname: c\n---\nDo the thing.\n")
	caseDir := filepath.Join(root, "evals", "cases", "c")
	if err := os.MkdirAll(caseDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(caseDir, "prompt.md")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "prompt.md") || !strings.Contains(joined, "not a regular file") {
		t.Fatalf("errors = %v, want prompt.md symlink reported", report.Errors)
	}
}

// TestEvalFixtureSymlinkIsReportedNotSkipped covers the fixture read inside
// checkFixtureAgentCopies.
func TestEvalFixtureSymlinkIsReportedNotSkipped(t *testing.T) {
	root := syntheticRoot(t)
	real := filepath.Join(root, "real-fixture.md")
	writeFixture(t, real, "body\n")
	fixtureDir := filepath.Join(root, "evals", "cases", "c", "fixtures")
	if err := os.MkdirAll(fixtureDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(fixtureDir, "coordinator.md")); err != nil {
		t.Fatal(err)
	}

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "fixtures/coordinator.md") || !strings.Contains(joined, "not a regular file") {
		t.Fatalf("errors = %v, want fixture symlink reported", report.Errors)
	}
}

// TestEvalGraderOversizeIsReportedNotSkipped covers MEDIUM 2 for an eval read:
// an oversized grader file must be a validation error naming the path and the
// limit, not a silent skip.
func TestEvalGraderOversizeIsReportedNotSkipped(t *testing.T) {
	root := syntheticRoot(t)
	caseDir := filepath.Join(root, "evals", "cases", "x")
	writeFixture(t, filepath.Join(caseDir, "case.yaml"), "id: x\n")
	writeFixture(t, filepath.Join(caseDir, "graders", "g.md"), strings.Repeat("a", maxReadableFileSize+1))

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "graders/g.md") || !strings.Contains(joined, "read limit") {
		t.Fatalf("errors = %v, want oversized grader reported with path and limit", report.Errors)
	}
}

func TestSymlinkedCaseSubdirectoriesRejected(t *testing.T) {
	// A glob walks through a symlinked named segment, so the files it matches
	// are regular files outside the case and the per-file guard cannot tell
	// them apart. Containment is decided on the directory.
	for _, name := range []string{"graders", "fixtures"} {
		t.Run(name, func(t *testing.T) {
			root := syntheticRoot(t)
			outside := filepath.Join(t.TempDir(), "elsewhere")
			if err := os.MkdirAll(outside, 0o755); err != nil {
				t.Fatal(err)
			}
			writeFixture(t, filepath.Join(outside, "leaked.md"), "---\ntype: regex\nweight: 1\n---\nq\n")
			caseDir := filepath.Join(root, "evals", "cases", "c")
			writeFixture(t, filepath.Join(caseDir, "prompt.md"), "---\nname: c\n---\nDo it.\n")
			if err := os.Symlink(outside, filepath.Join(caseDir, name)); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			report, err := Validate(root, Options{SkipMinimumCounts: true})
			if err != nil {
				t.Fatal(err)
			}
			joined := strings.Join(report.Errors, "\n")
			if !strings.Contains(joined, name+" is a symbolic link") {
				t.Fatalf("errors = %v, want %s rejected as a symbolic link", report.Errors, name)
			}
			if strings.Contains(joined, "leaked.md") {
				t.Fatalf("errors = %v, want no file under the symlinked directory to be read", report.Errors)
			}
		})
	}
}
