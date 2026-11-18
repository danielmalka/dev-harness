package harness

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// providerKinds are the entries of secretPatterns that have the shape of a key or a credential of a real provider. A test or a fixture must not spell one of
// them out: GitHub Secret Scanning and gitleaks read the source text, so a synthetic key written as a literal opens an alert even when it is fake (alert #1 of
// the repository, a Google key in facts_test.go). A test that needs one builds it at run time from pieces (see fakeSecrets).
var providerKinds = map[string]bool{
	"chave privada": true, "chave de acesso da AWS": true, "token do GitHub": true, "token do Slack": true, "chave de API (sk-)": true,
	"chave de API do Google": true, "JWT": true, "credencial em URL": true,
}

// providerShapes returns the provider kinds whose pattern the line matches.
func providerShapes(line string) []string {
	var kinds []string
	for _, p := range secretPatterns {
		if providerKinds[p.kind] && p.re.MatchString(line) {
			kinds = append(kinds, p.kind)
		}
	}
	return kinds
}

// The detector itself: every provider kind that fakeSecrets builds at run time is recognized, and the text of this repository's own regexes and look-alikes is not.
func TestProviderShapesDetectsTheFakeSecrets(t *testing.T) {
	found := map[string]bool{}
	for kind, secret := range fakeSecrets() {
		if !providerKinds[kind] {
			continue
		}
		got := providerShapes("valor " + secret + " visto")
		if len(got) == 0 {
			t.Errorf("%s: the detector does not recognize the fake secret", kind)
		}
		for _, g := range got {
			found[g] = true
		}
	}
	for kind := range providerKinds {
		if !found[kind] {
			t.Errorf("kind %q is never detected: fakeSecrets lacks a value for it", kind)
		}
	}
	for _, ok := range []string{"AKIA", "ghp_", "sk-curto", "risk-001-grader-02-known-gap", "xox", "AIza", "eyJ", "-----BEGIN", "https://example.com/a", "postgres://db.local/x"} {
		if got := providerShapes(ok); len(got) != 0 {
			t.Errorf("%q must not be a provider shape, got %v", ok, got)
		}
	}
}

// skipDirs are the folders the guard does not enter: git itself, the generated package, installed modules and the other worktrees of the checkout. The recorded
// eval results (evals/baselines, evals/results, evals/archive) are skipped by path. Nothing else is skipped, and a folder or a file the guard cannot read is a failure.
var skipDirs = map[string]bool{".git": true, "dist": true, "node_modules": true, ".worktrees": true}

// isTestOrFixture says whether rel (slash-separated, relative to the repository root) is a test file or a fixture: a Go or a TypeScript test, anything under a
// `testdata` folder at ANY level (the root included), and the eval cases, fixtures and not-run cases.
func isTestOrFixture(rel string) bool {
	return strings.HasSuffix(rel, "_test.go") || strings.HasSuffix(rel, ".test.ts") || strings.Contains("/"+rel, "/testdata/") ||
		strings.HasPrefix(rel, "evals/cases/") || strings.HasPrefix(rel, "evals/fixtures/") || strings.HasPrefix(rel, "evals/not-run/")
}

// scanProviderShapes walks root and returns one message per problem: a line with the shape of a provider key in a test or a fixture, and every path it could not
// walk or read (a guard that cannot read must fail, not pass in silence).
func scanProviderShapes(root string) []string {
	var problems []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		rel = filepath.ToSlash(rel)
		if err != nil {
			problems = append(problems, "cannot walk "+rel+": "+err.Error())
			return nil // keep going, so that every unreadable path is listed
		}
		if d.IsDir() {
			if skipDirs[d.Name()] || rel == "evals/baselines" || rel == "evals/results" || rel == "evals/archive" {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTestOrFixture(rel) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			problems = append(problems, "cannot read "+rel+": "+err.Error())
			return nil
		}
		for i, line := range strings.Split(string(data), "\n") {
			if kinds := providerShapes(line); len(kinds) > 0 {
				problems = append(problems, fmt.Sprintf("%s:%d spells out a provider-shaped secret (%s): build it at run time from pieces (fakeSecrets)", rel, i+1, strings.Join(kinds, ", ")))
			}
		}
		return nil
	})
	if err != nil {
		problems = append(problems, "walk failed: "+err.Error())
	}
	return problems
}

// Every test file and every fixture of the repository: no line has the shape of a provider key, and every path is readable.
func TestNoProviderShapedSecretLiteralsInTests(t *testing.T) {
	for _, problem := range scanProviderShapes(filepath.Join("..", "..")) {
		t.Error(problem)
	}
}

func TestIsTestOrFixture(t *testing.T) {
	for _, rel := range []string{"internal/harness/facts_test.go", "adapters/claude-code/plugin/hooks/panel.test.ts", "testdata/a.txt", "testdata/sub/b.json",
		"internal/dashboard/testdata/legacy-en/TASK.md", "a/b/testdata/c", "evals/cases/x/prompt.md", "evals/fixtures/slice-01/BRIEF.md", "evals/not-run/y/README.md"} {
		if !isTestOrFixture(rel) {
			t.Errorf("%s must be scanned", rel)
		}
	}
	for _, rel := range []string{"internal/harness/facts.go", "README.md", "mytestdata/a.txt", "a/testdata-old/b", "testdatax/c", "evals/baselines/a.json", "docs/prd/PRD-035.md"} {
		if isTestOrFixture(rel) {
			t.Errorf("%s must not be scanned", rel)
		}
	}
}

// The scanner itself, on a temporary tree: a planted literal is found in a root `testdata` folder, in a nested one and in a test file; the skipped folders and the
// files that are neither tests nor fixtures are not read; an unreadable test file and a missing root are failures, not silent passes.
func TestScanProviderShapes(t *testing.T) {
	planted := fakeSecrets()["chave de API do Google"] // built at run time: this file has no literal of the shape
	write := func(root, rel, body string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	root := t.TempDir()
	write(root, "testdata/a.txt", "ok\nkey "+planted+"\n")
	write(root, "pkg/testdata/b.json", "{\"k\": \""+planted+"\"}")
	write(root, "pkg/x_test.go", "package x\n\nconst k = \""+planted+"\"\n")
	write(root, "evals/cases/c/fixtures/d.md", planted)
	for _, ignored := range []string{"pkg/code.go", ".git/x_test.go", "dist/y_test.go", "node_modules/z_test.go", ".worktrees/w/x_test.go", "evals/baselines/b.json", "evals/results/r.json", "evals/archive/a.json", "README.md"} {
		write(root, ignored, planted)
	}
	got := strings.Join(scanProviderShapes(root), "\n")
	for _, want := range []string{"testdata/a.txt:2", "pkg/testdata/b.json:1", "pkg/x_test.go:3", "evals/cases/c/fixtures/d.md:1"} {
		if !strings.Contains(got, want+" spells out a provider-shaped secret (chave de API do Google)") {
			t.Errorf("want a finding at %s in:\n%s", want, got)
		}
	}
	for _, not := range []string{"code.go", ".git/", "dist/", "node_modules/", ".worktrees/", "evals/baselines", "evals/results", "evals/archive", "README.md"} {
		if strings.Contains(got, not) {
			t.Errorf("%s must not be scanned:\n%s", not, got)
		}
	}
	if strings.Contains(got, planted) {
		t.Error("a finding echoes the secret")
	}
	// A clean tree has no problem.
	clean := t.TempDir()
	write(clean, "testdata/a.txt", "nada aqui\n")
	if p := scanProviderShapes(clean); len(p) != 0 {
		t.Errorf("clean tree: %v", p)
	}
	// A missing root is a failure.
	if p := scanProviderShapes(filepath.Join(t.TempDir(), "nao-existe")); len(p) == 0 {
		t.Error("a missing root must be reported, not pass in silence")
	}
}

func TestScanProviderShapesFailsOnUnreadablePaths(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads everything")
	}
	root := t.TempDir()
	locked := filepath.Join(root, "pkg", "locked_test.go")
	if err := os.MkdirAll(filepath.Dir(locked), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(locked, []byte("package x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(locked, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(locked, 0o644) })
	closed := filepath.Join(root, "testdata", "closed")
	if err := os.MkdirAll(closed, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(closed, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(closed, 0o755) })
	got := strings.Join(scanProviderShapes(root), "\n")
	for _, want := range []string{"cannot read pkg/locked_test.go", "cannot walk testdata/closed"} {
		if !strings.Contains(got, want) {
			t.Errorf("want %q in:\n%s", want, got)
		}
	}
}
