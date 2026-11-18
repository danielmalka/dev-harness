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
			want: "author ausente",
		},
		{
			name: "skill description",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".skills", "sample", "SKILL.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "Use when", "Apply when", 1))
			},
			want: "description deve começar com",
		},
		{
			name: "skill description in pt-BR passes",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".skills", "sample", "SKILL.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "Use when", "Use quando", 1))
			},
		},
		{
			name:   "skill description still in the old English form passes during the transition",
			mutate: func(t *testing.T, root string) {},
		},
		{
			name: "template language folder en/ is rejected",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, "templates", "en", "TASK.md"), "# Task\n")
			},
			want: "a subpasta de idioma en/ não existe mais",
		},
		{
			name: "template language folder pt-br/ is rejected",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, "templates", "pt-br", "TASK.md"), "# Task\n")
			},
			want: "a subpasta de idioma pt-br/ não existe mais",
		},
		{
			name: "PRD and REVIEW live in templates/ like any other template",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, "templates", "PRD.md"), "# PRD\n")
				writeFixture(t, filepath.Join(root, "templates", "REVIEW.md"), "# Review\n")
			},
		},
		{
			name: "reference to a missing template",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, ".agents", "alpha.md"), agentFixture("alpha")+"Use templates/NOPE.md.\n")
			},
			want: "template ausente templates/NOPE.md",
		},
		{
			name: "reference with the old language folder checks the file in templates/",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, ".agents", "alpha.md"), agentFixture("alpha")+"Use templates/<lang>/TASK.md e templates/en/TASK.md e templates/pt-br/TASK.md.\n")
			},
		},
		{
			name: "reference with the old language folder to a missing file",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, ".agents", "alpha.md"), agentFixture("alpha")+"Use templates/en/PRD.md.\n")
			},
			want: "template ausente templates/PRD.md",
		},
		{
			name: "missing command role",
			mutate: func(t *testing.T, root string) {
				path := filepath.Join(root, ".commands", "run.md")
				text := readFixture(t, path)
				writeFixture(t, path, strings.Replace(text, "roles: [alpha]", "roles: [missing-role]", 1))
			},
			want: "papel missing-role ausente",
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
	writeFixture(t, filepath.Join(packageRoot, ".claude-plugin", "plugin.json"), `{"name":"dh","license":"MIT","version":"0.1.0"}`)
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
	if report.Passed() || !strings.Contains(strings.Join(report.Errors, "\n"), "pacote difere em") {
		t.Fatalf("errors = %v, want package differs", report.Errors)
	}
}

func TestCoverageEvalsExclusions(t *testing.T) {
	root := syntheticRoot(t)
	packageRoot := filepath.Join(root, "dist", "claude-code", "dev-harness")
	copyFixtureTree(t, filepath.Join(root, ".agents"), filepath.Join(packageRoot, "agents"))
	copyFixtureTree(t, filepath.Join(root, ".commands"), filepath.Join(packageRoot, "commands"))
	copyFixtureTree(t, filepath.Join(root, ".skills"), filepath.Join(packageRoot, "skills"))
	copyFixtureTree(t, filepath.Join(root, "templates"), filepath.Join(packageRoot, "templates"))
	copyFixtureTree(t, filepath.Join(root, "profiles"), filepath.Join(packageRoot, "profiles"))
	writePackageFixtures(t, packageRoot)

	// Source-only case dir: expect "package missing".
	writeFixture(t, filepath.Join(root, "evals", "cases", "x", "case.yaml"), "id: x\n")
	// Exclusions present only on one side must not fail coverage.
	writeFixture(t, filepath.Join(root, "evals", "baselines", "b.json"), "{}")
	writeFixture(t, filepath.Join(root, "evals", "archive", "old", "case.yaml"), "id: old\n")
	writeFixture(t, filepath.Join(packageRoot, "evals", "results", "r.json"), "{}")
	writeFixture(t, filepath.Join(root, "evals", "fixtures", "slice", "__pycache__", "x.pyc"), "cache")

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if !strings.Contains(joined, "pacote sem evals/cases/x/case.yaml") {
		t.Fatalf("errors = %v, want package missing evals/cases/x/case.yaml", report.Errors)
	}
	if strings.Contains(joined, "baselines") || strings.Contains(joined, "archive") || strings.Contains(joined, "results") || strings.Contains(joined, "__pycache__") {
		t.Fatalf("errors = %v, want exclusions to be silent", report.Errors)
	}

	// Now add the case only on the package side: expect "package extra".
	root2 := syntheticRoot(t)
	packageRoot2 := filepath.Join(root2, "dist", "claude-code", "dev-harness")
	copyFixtureTree(t, filepath.Join(root2, ".agents"), filepath.Join(packageRoot2, "agents"))
	copyFixtureTree(t, filepath.Join(root2, ".commands"), filepath.Join(packageRoot2, "commands"))
	copyFixtureTree(t, filepath.Join(root2, ".skills"), filepath.Join(packageRoot2, "skills"))
	copyFixtureTree(t, filepath.Join(root2, "templates"), filepath.Join(packageRoot2, "templates"))
	copyFixtureTree(t, filepath.Join(root2, "profiles"), filepath.Join(packageRoot2, "profiles"))
	writePackageFixtures(t, packageRoot2)
	writeFixture(t, filepath.Join(packageRoot2, "evals", "cases", "y", "case.yaml"), "id: y\n")

	report2, err := Validate(root2, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(report2.Errors, "\n"), "pacote com sobra evals/cases/y/case.yaml") {
		t.Fatalf("errors = %v, want package extra evals/cases/y/case.yaml", report2.Errors)
	}
}

func writePackageFixtures(t *testing.T, packageRoot string) {
	t.Helper()
	writeFixture(t, filepath.Join(packageRoot, ".claude-plugin", "plugin.json"), `{"name":"dh","license":"MIT","version":"0.1.0"}`)
	writeFixture(t, filepath.Join(packageRoot, "settings.json"), "{}")
	writeFixture(t, filepath.Join(packageRoot, "hooks", "hooks.json"), "{}")
	writeFixture(t, filepath.Join(packageRoot, "harness-manifest.json"), `{"binaries":[]}`)
	writeFixture(t, filepath.Join(packageRoot, "GENERATED.txt"), "generated")
	writeFixture(t, filepath.Join(packageRoot, "bin", "dh"), "#!/bin/sh\n")
	if err := os.Chmod(filepath.Join(packageRoot, "bin", "dh"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestPrivatePathsSkipEvalsBaselinesAndResults(t *testing.T) {
	root := syntheticRoot(t)
	writeFixture(t, filepath.Join(root, "evals", "baselines", "b.md"), "path /home/someone/project\n")
	writeFixture(t, filepath.Join(root, "evals", "results", "r.md"), "path /home/someone/project\n")
	writeFixture(t, filepath.Join(root, "evals", "cases", "x", "notes.md"), "path /home/someone/project\n")

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(report.Errors, "\n")
	if strings.Contains(joined, "evals/baselines") || strings.Contains(joined, "evals/results") {
		t.Fatalf("errors = %v, want baselines/results skipped", report.Errors)
	}
	if !strings.Contains(joined, "evals/cases/x/notes.md") {
		t.Fatalf("errors = %v, want private path reported outside baselines/results", report.Errors)
	}
}

// A development worktree lives in <repo>/.worktrees/<id> (PRD-033) and holds a full copy of the repository: the walkers must not
// read it as part of the repository, or every private path, link and secret check runs twice and fails on the copy.
func TestValidateSkipsWorktreesFolder(t *testing.T) {
	root := syntheticRoot(t)
	writeFixture(t, filepath.Join(root, "evals", "cases", "x", "notes.md"), "clean\n")
	writeFixture(t, filepath.Join(root, ".worktrees", "T-1", "evals", "cases", "x", "notes.md"), "path /home/someone/project\n")
	writeFixture(t, filepath.Join(root, ".worktrees", "T-1", "docs", "a.md"), "[broken](missing.md)\n")

	report, err := Validate(root, Options{SkipMinimumCounts: true})
	if err != nil {
		t.Fatal(err)
	}
	if joined := strings.Join(report.Errors, "\n"); strings.Contains(joined, ".worktrees") {
		t.Fatalf("errors = %v, want the .worktrees folder skipped", report.Errors)
	}
}

func TestAntiDelegationClauseParity(t *testing.T) {
	block := "## External CLI reviewers\n\n```\nrule one\nrule two\n```\n"
	skillBlock := "## A cláusula anti-delegação\n\n```\nrule one\nrule two\n```\n"
	whitespaceBlock := "## A cláusula anti-delegação\n\n```\n rule one \n  rule two\n```\n"
	divergentBlock := "## A cláusula anti-delegação\n\n```\nrule one\nrule three\n```\n"

	tests := []struct {
		name           string
		coordinator    string
		skill          string
		writeSkillFile bool
		want           string
	}{
		{name: "identical", coordinator: block, skill: skillBlock, writeSkillFile: true, want: ""},
		{name: "divergent", coordinator: block, skill: divergentBlock, writeSkillFile: true, want: "a cláusula anti-delegação difere entre .agents/coordinator.md e .skills/external-clis/SKILL.md"},
		{name: "whitespace only", coordinator: block, skill: whitespaceBlock, writeSkillFile: true, want: ""},
		{name: "skill file absent", coordinator: block, writeSkillFile: false, want: ""},
		// S2: the English heading of the skill is no longer read, so an old-heading skill with a different block stays silent (the pairing needs both headings).
		{name: "old English heading is not read", coordinator: block, skill: strings.Replace(divergentBlock, "## A cláusula anti-delegação", "## The anti-delegation clause", 1), writeSkillFile: true, want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := syntheticRoot(t)
			writeFixture(t, filepath.Join(root, ".agents", "coordinator.md"), "---\nname: coordinator\ndescription: Use when testing\nauthor: test\nmodel: haiku\n---\n"+test.coordinator)
			if test.writeSkillFile {
				writeFixture(t, filepath.Join(root, ".skills", "external-clis", "SKILL.md"), "---\nname: external-clis\ndescription: Use when testing\nauthor: test\n---\n"+test.skill)
			}
			report, err := Validate(root, Options{SkipMinimumCounts: true})
			if err != nil {
				t.Fatal(err)
			}
			joined := strings.Join(report.Errors, "\n")
			if test.want == "" {
				if strings.Contains(joined, "a cláusula anti-delegação difere") {
					t.Fatalf("errors = %v, want no clause-parity error", report.Errors)
				}
				return
			}
			if !strings.Contains(joined, test.want) {
				t.Fatalf("errors = %v, want substring %q", report.Errors, test.want)
			}
		})
	}
}

func TestEvalReferences(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(t *testing.T, root string)
		want    string
		noError bool
	}{
		{
			name: "existing add_dirs",
			mutate: func(t *testing.T, root string) {
				caseDir := filepath.Join(root, "evals", "cases", "x")
				writeFixture(t, filepath.Join(caseDir, "fixtures", ".gitkeep"), "")
				writeFixture(t, filepath.Join(caseDir, "case.yaml"), "context:\n  add_dirs:\n    - fixtures\n")
			},
			noError: true,
		},
		{
			name: "missing add_dirs",
			mutate: func(t *testing.T, root string) {
				caseDir := filepath.Join(root, "evals", "cases", "x")
				writeFixture(t, filepath.Join(caseDir, "case.yaml"), "context:\n  add_dirs:\n    - fixtures/nao-existe\n")
			},
			want: "context.add_dirs não resolvido: fixtures/nao-existe",
		},
		{
			name: "escaping add_dirs",
			mutate: func(t *testing.T, root string) {
				writeFixture(t, filepath.Join(root, "evals", "cases", "outro-caso", "fixtures", ".gitkeep"), "")
				caseDir := filepath.Join(root, "evals", "cases", "x")
				writeFixture(t, filepath.Join(caseDir, "case.yaml"), "context:\n  add_dirs:\n    - ../outro-caso/fixtures\n")
			},
			want: "context.add_dirs sai da pasta do caso: ../outro-caso/fixtures",
		},
		{
			name: "grader target.path outside case",
			mutate: func(t *testing.T, root string) {
				caseDir := filepath.Join(root, "evals", "cases", "x")
				writeFixture(t, filepath.Join(caseDir, "case.yaml"), "id: x\n")
				writeFixture(t, filepath.Join(caseDir, "graders", "g.md"), "---\ntarget:\n  source: file\n  path: ../../../outside.txt\n---\ngrader\n")
			},
			want: "target.path sai da pasta do caso: ../../../outside.txt",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := syntheticRoot(t)
			test.mutate(t, root)
			report, err := Validate(root, Options{SkipMinimumCounts: true})
			if err != nil {
				t.Fatal(err)
			}
			joined := strings.Join(report.Errors, "\n")
			if test.noError {
				if strings.Contains(joined, "não resolvido") || strings.Contains(joined, "sai da pasta do caso") {
					t.Fatalf("errors = %v, want no eval-reference error", report.Errors)
				}
				return
			}
			if !strings.Contains(joined, test.want) {
				t.Fatalf("errors = %v, want substring %q", report.Errors, test.want)
			}
		})
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
		"templates",
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
	for _, profile := range []string{"base", "go-api", "typescript-web", "kotlin", "python", "typescript-api"} {
		writeFixture(t, filepath.Join(root, "profiles", profile+".yaml"), "id: "+profile+"\n")
	}
	writeFixture(t, filepath.Join(root, "templates", "TASK.md"), "# Task\n")
	return root
}

func agentFixture(name string) string {
	return "---\nname: " + name + "\ndescription: Use when testing\nauthor: test\nmodel: haiku\n---\nagent\n"
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

// TestCheckCountsAgentMinimum pins the agent floor at 13 (PRD-010 R5): 13
// agents pass, 12 fail.
func TestCheckCountsAgentMinimum(t *testing.T) {
	if MinAgents != 13 {
		t.Fatalf("MinAgents = %d, want 13", MinAgents)
	}
	inventory := func(agents int) Inventory {
		return Inventory{
			Agents:   make([]string, agents),
			Commands: make([]string, MinCommands),
			Skills:   make([]string, MinSkills),
			Profiles: []string{"base", "go-api", "typescript-web", "kotlin", "python", "typescript-api"},
		}
	}
	var errs []string
	checkCounts(inventory(13), &errs)
	if len(errs) != 0 {
		t.Fatalf("13 agents: errors = %v, want none", errs)
	}
	checkCounts(inventory(12), &errs)
	if len(errs) != 1 || errs[0] != "agentes: 12 < 13" {
		t.Fatalf("12 agents: errors = %v, want [agents: 12 < 13]", errs)
	}
}

// TestCheckCountsRequiresNewProfiles pins PRD-015 R8: dropping any of the
// three language profiles yields "missing profile: <name>".
func TestCheckCountsRequiresNewProfiles(t *testing.T) {
	all := []string{"base", "go-api", "typescript-web", "kotlin", "python", "typescript-api"}
	for _, drop := range []string{"kotlin", "python", "typescript-api"} {
		var profiles []string
		for _, p := range all {
			if p != drop {
				profiles = append(profiles, p)
			}
		}
		inv := Inventory{
			Agents:   make([]string, MinAgents),
			Commands: make([]string, MinCommands),
			Skills:   make([]string, MinSkills),
			Profiles: profiles,
		}
		var errs []string
		checkCounts(inv, &errs)
		if len(errs) != 1 || errs[0] != "perfil ausente: "+drop {
			t.Fatalf("drop %s: errors = %v", drop, errs)
		}
	}
}
