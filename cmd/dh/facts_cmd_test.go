package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const cmdFact = "- F-001 · O comando grava com trava de hash · Evidência: src/a.go:3 · Visto: 2026-11-18, builder, T-1 · sha: abc1234 · Estado: %s"

// factsRepo makes a repo-mode project with a source file for the evidence and, when content is not empty, a FACTS.md.
func factsRepo(t *testing.T, content string) (repo, facts string) {
	t.Helper()
	repo, _ = memRepo(t, "")
	if err := os.MkdirAll(filepath.Join(repo, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "src", "a.go"), []byte("package a\n\nfunc A() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	facts = filepath.Join(repo, ".harness", "FACTS.md")
	if content != "" {
		if err := os.WriteFile(facts, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo, facts
}

func factsDocFor(sections ...string) string {
	return "# Fatos do sistema\n\nÚltima curadoria: 2026-11-18, sha abc1234\n\n" + strings.Join(sections, "\n") + "\n"
}

// R11: --file FACTS.md hashes and writes FACTS.md with the same contract as MEMORY.md; without --file the file is still MEMORY.md.
func TestMemoryFileFlagWorksOnFacts(t *testing.T) {
	repo, facts := factsRepo(t, "")
	mem := filepath.Join(repo, ".harness", "MEMORY.md")
	if code, out, _ := dh(t, "memory", "hash", "--file", "FACTS.md", repo); code != 0 || strings.TrimSpace(out) != "absent" {
		t.Fatalf("hash absent: %d %q", code, out)
	}
	code, out, e := dhStdin(t, "# Fatos\n", "memory", "write", "--file", "FACTS.md", "--expect", "absent", repo)
	if code != 0 || !strings.HasPrefix(out, "written "+shaOf("# Fatos\n")) {
		t.Fatalf("first write: %d %q %q", code, out, e)
	}
	if b, _ := os.ReadFile(facts); string(b) != "# Fatos\n" {
		t.Fatalf("FACTS.md = %q", b)
	}
	if _, err := os.Stat(mem); err == nil {
		t.Fatal("writing FACTS.md created MEMORY.md")
	}
	if code, out, _ := dh(t, "memory", "hash", "--file", "FACTS.md", repo); code != 0 || strings.TrimSpace(out) != shaOf("# Fatos\n") {
		t.Fatalf("hash after write: %d %q", code, out)
	}
	if code, out, _ := dh(t, "memory", "hash", repo); code != 0 || strings.TrimSpace(out) != "absent" {
		t.Fatalf("default is still MEMORY.md (absent): %d %q", code, out)
	}
	// A stale hash exits 3, names FACTS.md and keeps the other writer's content.
	code, _, e = dhStdin(t, "mine", "memory", "write", "--file", "FACTS.md", "--expect", shaOf("outra leitura"), repo)
	if code != 3 || !strings.Contains(e, "FACTS.md mudou desde que foi lido") || strings.Contains(e, "MEMORY.md") {
		t.Fatalf("stale: %d %q", code, e)
	}
	if b, _ := os.ReadFile(facts); string(b) != "# Fatos\n" {
		t.Fatalf("FACTS.md changed by a stale write: %q", b)
	}
	// --file may come after --expect, and --from works.
	draft := filepath.Join(t.TempDir(), "draft.md")
	os.WriteFile(draft, []byte("# Fatos\nv2\n"), 0o644)
	if code, out, e := dh(t, "memory", "write", "--expect", shaOf("# Fatos\n"), "--from", draft, "--file", "FACTS.md", repo); code != 0 || !strings.HasPrefix(out, "written ") {
		t.Fatalf("flag order: %d %q %q", code, out, e)
	}
}

// R11: any name outside the closed list exits 1, lists both valid names and writes nothing, for hash and for write.
func TestMemoryFileFlagClosedList(t *testing.T) {
	repo, _ := factsRepo(t, "")
	for _, name := range []string{"EPOCHAL.md", "RISKS.md", "REPO-MAP.md", "facts.md", "../x.md", "project.yaml"} {
		code, _, e := dh(t, "memory", "hash", "--file", name, repo)
		if code != 1 || !strings.Contains(e, "MEMORY.md") || !strings.Contains(e, "FACTS.md") {
			t.Errorf("hash --file %s: %d %q", name, code, e)
		}
		code, _, e = dhStdin(t, "x", "memory", "write", "--file", name, "--expect", "absent", repo)
		if code != 1 || !strings.Contains(e, "MEMORY.md") || !strings.Contains(e, "FACTS.md") {
			t.Errorf("write --file %s: %d %q", name, code, e)
		}
	}
	ents, _ := os.ReadDir(filepath.Join(repo, ".harness"))
	if len(ents) != 0 {
		t.Fatalf("a refused name left files in the harness dir: %v", ents)
	}
	if code, _, e := dh(t, "memory", "hash", "--file"); code != 1 {
		t.Fatalf("--file without value: %d %q", code, e)
	}
}

// R4 of PRD-034 for FACTS.md, separate processes: N writers, one winner, the rest exit 3, no lock or temp file left.
func TestFactsConcurrentProcesses(t *testing.T) {
	const writers = 8
	repo, facts := factsRepo(t, "start")
	home := os.Getenv("DH_HOME")
	goFile := filepath.Join(t.TempDir(), "go")
	var cmds []*exec.Cmd
	var contents []string
	for i := 0; i < writers; i++ {
		c := strings.Repeat(fmt.Sprintf("facts-process-%d\n", i), 300)
		contents = append(contents, c)
		draft := filepath.Join(t.TempDir(), fmt.Sprintf("draft-%d.md", i))
		os.WriteFile(draft, []byte(c), 0o644)
		cmd := exec.Command(os.Args[0], "-test.run=TestMemoryHelperProcess")
		cmd.Env = append(os.Environ(), "DH_MEMORY_HELPER=1", "DH_HOME="+home, "DH_MEMORY_GO="+goFile, "DH_MEMORY_EXPECT="+shaOf("start"),
			"DH_MEMORY_FROM="+draft, "DH_MEMORY_REPO="+repo, "DH_MEMORY_FILE=FACTS.md")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, cmd)
	}
	os.WriteFile(goFile, nil, 0o644)
	wins, conflicts := 0, 0
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err == nil {
			wins++
		} else if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 3 {
			conflicts++
		} else {
			t.Fatalf("unexpected exit: %v", err)
		}
	}
	if wins != 1 || conflicts != writers-1 {
		t.Fatalf("wins=%d conflicts=%d, want 1 and %d", wins, conflicts, writers-1)
	}
	got, _ := os.ReadFile(facts)
	ok := false
	for _, c := range contents {
		ok = ok || string(got) == c
	}
	if !ok {
		t.Fatalf("the final FACTS.md is not the content of any single writer (%d bytes)", len(got))
	}
	dir := filepath.Dir(facts)
	if left, _ := filepath.Glob(filepath.Join(dir, "*.lock")); len(left) != 0 {
		t.Fatalf("lock left: %v", left)
	}
	if left, _ := filepath.Glob(filepath.Join(dir, ".FACTS.md.tmp-*")); len(left) != 0 {
		t.Fatalf("temp left: %v", left)
	}
}

// R27: absent exits 0 and says `absent`; a good file exits 0 and says `ok`, the facts and the candidates; nothing is written.
func TestFactsCheckOKAndAbsent(t *testing.T) {
	repo, facts := factsRepo(t, "")
	if code, out, e := dh(t, "facts", "check", repo); code != 0 || strings.TrimSpace(out) != "absent" {
		t.Fatalf("absent: %d %q %q", code, out, e)
	}
	doc := factsDocFor("## Testes", fmt.Sprintf(cmdFact, "confirmado"), "## Candidatos", strings.Replace(fmt.Sprintf(cmdFact, "candidato"), "F-001", "F-002", 1))
	os.WriteFile(facts, []byte(doc), 0o644)
	before, _ := os.ReadDir(filepath.Dir(facts))
	code, out, e := dh(t, "facts", "check", repo)
	if code != 0 || strings.TrimSpace(out) != "ok 2 fatos 1 candidatos" || e != "" {
		t.Fatalf("ok: %d %q %q", code, out, e)
	}
	after, _ := os.ReadDir(filepath.Dir(facts))
	if len(before) != len(after) {
		t.Fatalf("the check wrote files: %v -> %v", before, after)
	}
	if b, _ := os.ReadFile(facts); string(b) != doc {
		t.Fatal("the check changed FACTS.md")
	}
}

// R27: a violation exits 1, names the line and the cause on stderr; an evidence path that is not in the repository is a violation.
func TestFactsCheckReportsViolations(t *testing.T) {
	repo, facts := factsRepo(t, "")
	cases := map[string]struct{ doc, want string }{
		"missing field": {factsDocFor("## Testes", "- F-001 · frase · Evidência: src/a.go:3 · Visto: 2026-11-18, builder, T-1 · sha: abc1234"), "FACTS.md:6: campo Estado: ausente"},
		"unknown path":  {factsDocFor("## Testes", strings.Replace(fmt.Sprintf(cmdFact, "confirmado"), "src/a.go", "src/b.go", 1)), `não é um arquivo do repositório`},
		"category name": {factsDocFor("## Outra coisa", fmt.Sprintf(cmdFact, "confirmado")), "não é uma das categorias"},
		"secret":        {factsDocFor("## Testes", strings.Replace(fmt.Sprintf(cmdFact, "confirmado"), "O comando", "token: abc12345xyz O comando", 1)), "padrão de segredo"},
		"ceiling":       {factsDocFor("## Mudanças") + strings.Repeat("- mudança\n", 150), "teto de linhas excedido"},
	}
	for name, c := range cases {
		os.WriteFile(facts, []byte(c.doc), 0o644)
		code, out, e := dh(t, "facts", "check", repo)
		if code != 1 || !strings.Contains(e, c.want) || !strings.Contains(e, "reprovado") || out != "" {
			t.Errorf("%s: %d out=%q err=%q (want %q)", name, code, out, e, c.want)
		}
		if strings.Contains(e, "abc12345xyz") {
			t.Errorf("%s: the secret was echoed: %q", name, e)
		}
	}
}

// R27: mode none names the cause and exits 1; bad usage exits 1.
func TestFactsCheckModeNoneAndUsage(t *testing.T) {
	_, none := env(t)
	if code, _, e := dh(t, "facts", "check", none); code != 1 || !strings.Contains(e, "mode none") {
		t.Fatalf("mode none: %d %q", code, e)
	}
	repo, _ := factsRepo(t, "")
	for name, args := range map[string][]string{
		"no subcommand": {"facts"},
		"unknown":       {"facts", "nope"},
		"two folders":   {"facts", "check", repo, repo},
		"not a folder":  {"facts", "check", filepath.Join(repo, "nao-existe")},
	} {
		if code, _, e := dh(t, args...); code != 1 || e == "" {
			t.Errorf("%s: %d %q", name, code, e)
		}
	}
}

// Global mode: the file lives in <home>/projects/<name>/FACTS.md and the evidence is checked against the project folder.
func TestFactsCheckGlobalMode(t *testing.T) {
	repo, _ := factsRepo(t, "")
	home := os.Getenv("DH_HOME")
	if code, _, e := dh(t, "link", "proj", repo); code != 0 {
		t.Fatalf("link: %d %q", code, e)
	}
	os.RemoveAll(filepath.Join(repo, ".harness")) // global mode only: no in-repo harness dir
	dir := filepath.Join(home, "projects", "proj")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir, "FACTS.md"), []byte(factsDocFor("## Testes", fmt.Sprintf(cmdFact, "confirmado"))), 0o644)
	if code, out, e := dh(t, "facts", "check", repo); code != 0 || strings.TrimSpace(out) != "ok 1 fatos 0 candidatos" {
		t.Fatalf("global: %d %q %q", code, out, e)
	}
}

// R27 (review of F1): the evidence must be a FILE. A folder passes os.Stat but is not evidence; a symlink counts only when it points to a file inside the project.
func TestFactsCheckEvidenceMustBeAFile(t *testing.T) {
	repo, facts := factsRepo(t, "")
	outside := filepath.Join(t.TempDir(), "outside.go")
	os.WriteFile(outside, []byte("package outside\n"), 0o644)
	linkOK := os.Symlink(filepath.Join(repo, "src", "a.go"), filepath.Join(repo, "src", "link-file.go")) == nil
	os.Symlink(filepath.Join(repo, "src"), filepath.Join(repo, "link-dir"))
	os.Symlink(outside, filepath.Join(repo, "src", "link-outside.go"))
	check := func(evidence string) (int, string) {
		doc := factsDocFor("## Testes", strings.Replace(fmt.Sprintf(cmdFact, "confirmado"), "src/a.go:3", evidence, 1))
		os.WriteFile(facts, []byte(doc), 0o644)
		code, _, e := dh(t, "facts", "check", repo)
		return code, e
	}
	if code, e := check("src/a.go:3"); code != 0 {
		t.Fatalf("a regular file must pass: %d %q", code, e)
	}
	for _, dir := range []string{"src:1", ".harness:2", "link-dir:1", "src/:1"} {
		code, e := check(dir)
		if code != 1 || !strings.Contains(e, "não é um arquivo do repositório") {
			t.Errorf("folder %q must fail as not a file: %d %q", dir, code, e)
		}
	}
	if linkOK {
		if code, e := check("src/link-file.go:1"); code != 0 {
			t.Errorf("a symlink to a file inside the project must pass: %d %q", code, e)
		}
		if code, e := check("src/link-outside.go:1"); code != 1 || !strings.Contains(e, "não é um arquivo do repositório") {
			t.Errorf("a symlink to a file outside the project must fail: %d %q", code, e)
		}
	}
	if code, e := check("src/nao-existe.go:1"); code != 1 || !strings.Contains(e, "não é um arquivo do repositório") {
		t.Errorf("a missing path must fail: %d %q", code, e)
	}
}
