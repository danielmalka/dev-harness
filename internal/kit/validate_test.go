package kit

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestSyntheticKitMutations(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(t *testing.T, root string)
		want   string
	}{
		{name: "passes", want: ""},
		{
			name: "missing author",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".agents", "alpha.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "author: test\n", "", 1))
			},
			want: "missing author",
		},
		{
			name: "skill description",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".skills", "sample", "SKILL.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "Use when", "Apply when", 1))
			},
			want: "description must start with",
		},
		{
			name: "template parity",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "templates", "pt-br", "TASK.md")); err != nil {
					t.Fatal(err)
				}
			},
			want: "is not present in every language",
		},
		{
			name: "missing command role",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".commands", "run.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "roles: [alpha]", "roles: [missing-role]", 1))
			},
			want: "missing role missing-role",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := syntheticRoot(t)
			if test.mutate != nil {
				test.mutate(t, root)
			}
			report, err := Validate(root, Options{SkipMinimumCounts: true})
			if err != nil {
				t.Fatal(err)
			}
			if test.want == "" {
				if !report.Passed() {
					t.Fatalf("unexpected errors: %v", report.Errors)
				}
				return
			}
			if report.Passed() || !strings.Contains(strings.Join(report.Errors, "\n"), test.want) {
				t.Fatalf("errors = %v, want substring %q", report.Errors, test.want)
			}
		})
	}
}

func TestPackageCoverageChecksContent(t *testing.T) {
	root := syntheticRoot(t)
	packageRoot := filepath.Join(root, "dist", "claude-code", "dev-harness")
	copyFixtureTree(t, filepath.Join(root, ".agents"), filepath.Join(packageRoot, "agents"))
	copyFixtureTree(t, filepath.Join(root, ".commands"), filepath.Join(packageRoot, "commands"))
	copyFixtureTree(t, filepath.Join(root, ".skills"), filepath.Join(packageRoot, "skills"))
	copyFixtureTree(t, filepath.Join(root, "templates"), filepath.Join(packageRoot, "templates"))
	copyFixtureTree(t, filepath.Join(root, "profiles"), filepath.Join(packageRoot, "profiles"))
	writeFixture(t, filepath.Join(packageRoot, ".claude-plugin", "plugin.json"), `{"name":"dev-harness","license":"MIT","version":"0.1.0"}`)
	writeFixture(t, filepath.Join(packageRoot, "settings.json"), "{}")
	writeFixture(t, filepath.Join(packageRoot, "hooks", "hooks.json"), "{}")
	writeFixture(t, filepath.Join(packageRoot, "harness-manifest.json"), `{"binaries":[]}`)
	writeFixture(t, filepath.Join(packageRoot, "GENERATED.txt"), "generated")
	writeFixture(t, filepath.Join(packageRoot, "bin", "dh"), "#!/bin/sh\n")
	if err := os.Chmod(filepath.Join(packageRoot, "bin", "dh"), 0o755); err != nil {
		t.Fatal(err)
	}

	writeFixture(t, filepath.Join(packageRoot, "agents", "alpha.md"), "different")
	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed() || !strings.Contains(strings.Join(report.Errors, "\n"), "package differs") {
		t.Fatalf("errors = %v, want package differs", report.Errors)
	}
}

func TestRepositoryPasses(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	report, err := Validate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed() {
		t.Fatalf("repository validation failed: %v", report.Errors)
	}
}

func syntheticRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{
		".agents", ".commands", ".skills/sample", "profiles",
		"templates/en", "templates/pt-br",
	} {
		if err := os.MkdirAll(filepath.Join(root, path), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFixture(t, filepath.Join(root, ".agents", "alpha.md"), agentFixture("alpha"))
	writeFixture(t, filepath.Join(root, ".agents", "beta.md"), agentFixture("beta"))
	writeFixture(t, filepath.Join(root, ".commands", "run.md"), commandFixture())
	writeFixture(t, filepath.Join(root, ".commands", "check.md"), commandFixtureWithName("check"))
	writeFixture(t, filepath.Join(root, ".skills", "sample", "SKILL.md"), skillFixture())
	for _, profile := range []string{"base", "go-api", "typescript-web"} {
		writeFixture(t, filepath.Join(root, "profiles", profile+".yaml"), "id: "+profile+"\n")
	}
	for _, lang := range []string{"en", "pt-br"} {
		writeFixture(t, filepath.Join(root, "templates", lang, "TASK.md"), "# Task\n")
	}
	return root
}

func agentFixture(name string) string {
	return "---\nname: " + name + "\ndescription: Use when testing\nauthor: test\nmodel: haiku\ntools:\n  - Read\n---\nagent\n"
}

func commandFixture() string {
	return commandFixtureWithName("run")
}

func commandFixtureWithName(name string) string {
	return "---\nname: " + name + "\ndescription: Run test\nauthor: test\nargument-hint: \"\"\nmetadata:\n  roles: [alpha]\n  skills: [sample]\n  writes: none\n---\ncommand\n"
}

func skillFixture() string {
	return "---\nname: sample\ndescription: Use when testing\nauthor: test\nmetadata:\n  provenance: test\n---\nskill\n"
}

func readFixture(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeFixture(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func copyFixtureTree(t *testing.T, source, destination string) {
	t.Helper()
	if err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
}
