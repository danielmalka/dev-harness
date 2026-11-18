package harness

import (
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

// Every test file and every fixture of the repository: no line has the shape of a provider key. The recorded eval results and the generated package are outside.
func TestNoProviderShapedSecretLiteralsInTests(t *testing.T) {
	root := filepath.Join("..", "..")
	skipDirs := map[string]bool{".git": true, "dist": true, "node_modules": true, ".worktrees": true}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if skipDirs[d.Name()] || rel == "evals/baselines" || rel == "evals/results" || rel == "evals/archive" {
				return filepath.SkipDir
			}
			return nil
		}
		isTest := strings.HasSuffix(rel, "_test.go") || strings.HasSuffix(rel, ".test.ts") || strings.Contains(rel, "/testdata/") ||
			strings.HasPrefix(rel, "evals/cases/") || strings.HasPrefix(rel, "evals/fixtures/") || strings.HasPrefix(rel, "evals/not-run/")
		if !isTest {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for i, line := range strings.Split(string(data), "\n") {
			if kinds := providerShapes(line); len(kinds) > 0 {
				t.Errorf("%s:%d spells out a provider-shaped secret (%s): build it at run time from pieces (fakeSecrets)", rel, i+1, strings.Join(kinds, ", "))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
