package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func shaOf(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

// memRepo makes a repo with a repo-mode harness dir and returns the repo and the MEMORY.md path.
func memRepo(t *testing.T, content string) (repo, mem string) {
	t.Helper()
	_, repo = env(t)
	hd := filepath.Join(repo, ".harness")
	if err := os.MkdirAll(hd, 0o755); err != nil {
		t.Fatal(err)
	}
	mem = filepath.Join(hd, "MEMORY.md")
	if content != "" {
		if err := os.WriteFile(mem, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repo, mem
}

func dhStdin(t *testing.T, stdin string, args ...string) (int, string, string) {
	t.Helper()
	var out, errb bytes.Buffer
	code := run(args, strings.NewReader(stdin), &out, &errb)
	return code, out.String(), errb.String()
}

// R1: hash prints the sha256 of the file, `absent` without it, and fails with a cause in mode none.
func TestMemoryHashCommand(t *testing.T) {
	repo, _ := memRepo(t, "# memory\n")
	if code, out, _ := dh(t, "memory", "hash", repo); code != 0 || strings.TrimSpace(out) != shaOf("# memory\n") {
		t.Fatalf("hash: %d %q", code, out)
	}
	repo2, _ := memRepo(t, "")
	if code, out, _ := dh(t, "memory", "hash", repo2); code != 0 || strings.TrimSpace(out) != "absent" {
		t.Fatalf("absent: %d %q", code, out)
	}
	_, none := env(t) // a folder with no harness dir and no registration: mode none
	if code, _, e := dh(t, "memory", "hash", none); code != 1 || !strings.Contains(e, "mode none") {
		t.Fatalf("mode none: %d %q", code, e)
	}
}

// R2 and R3 through the command: stdin and --from, equal hash writes, a stale hash exits 3 and keeps the other content.
func TestMemoryWriteCommand(t *testing.T) {
	repo, mem := memRepo(t, "v1")
	code, out, e := dhStdin(t, "v2 from stdin", "memory", "write", "--expect", shaOf("v1"), repo)
	if code != 0 || strings.TrimSpace(out) != "written "+shaOf("v2 from stdin") {
		t.Fatalf("write stdin: %d %q %q", code, out, e)
	}
	draft := filepath.Join(t.TempDir(), "draft.md")
	os.WriteFile(draft, []byte("v3 from file"), 0o644)
	if code, out, e := dh(t, "memory", "write", "--expect", shaOf("v2 from stdin"), "--from", draft, repo); code != 0 || !strings.HasPrefix(out, "written ") {
		t.Fatalf("write --from: %d %q %q", code, out, e)
	}
	// someone else writes between my read and my write
	os.WriteFile(mem, []byte("v4 by another session"), 0o644)
	code, _, e = dh(t, "memory", "write", "--expect", shaOf("v3 from file"), "--from", draft, repo)
	if code != 3 || !strings.Contains(e, "alguém escreveu no meio") || !strings.Contains(e, shaOf("v4 by another session")) {
		t.Fatalf("conflict: %d %q", code, e)
	}
	if b, _ := os.ReadFile(mem); string(b) != "v4 by another session" {
		t.Fatalf("the other session's content was lost: %q", b)
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(mem), ".MEMORY.md.tmp-*")); len(left) != 0 {
		t.Fatalf("temp files: %v", left)
	}
}

func TestMemoryWriteAbsentAndUsage(t *testing.T) {
	repo, mem := memRepo(t, "")
	if code, _, e := dhStdin(t, "first memory", "memory", "write", "--expect", "absent", repo); code != 0 {
		t.Fatalf("absent: %d %q", code, e)
	}
	if b, _ := os.ReadFile(mem); string(b) != "first memory" {
		t.Fatalf("file: %q", b)
	}
	if code, _, _ := dhStdin(t, "again", "memory", "write", "--expect", "absent", repo); code != 3 {
		t.Fatalf("absent on an existing file: exit %d, want 3", code)
	}
	for name, args := range map[string][]string{
		"no subcommand": {"memory"},
		"unknown":       {"memory", "nope"},
		"no --expect":   {"memory", "write", repo},
		"bad --expect":  {"memory", "write", "--expect", "abc", repo},
	} {
		if code, _, _ := dhStdin(t, "x", args...); code != 1 {
			t.Errorf("%s: exit %d, want 1", name, code)
		}
	}
	if code, _, e := dhStdin(t, "", "memory", "write", "--expect", shaOf("first memory"), repo); code != 1 || !strings.Contains(e, "vazio") {
		t.Errorf("empty content: %d %q", code, e)
	}
	if b, _ := os.ReadFile(mem); string(b) != "first memory" {
		t.Fatalf("file changed by a refused write: %q", b)
	}
}

// R4, separate processes: the test binary re-executes itself as N writers that wait for a start file and then run
// `dh memory write` with the SAME expected hash and different contents. One wins, the rest exit 3, the file is the winner's.
func TestMemoryHelperProcess(t *testing.T) {
	if os.Getenv("DH_MEMORY_HELPER") != "1" {
		t.Skip("helper process")
	}
	for i := 0; i < 400 && func() bool { _, err := os.Stat(os.Getenv("DH_MEMORY_GO")); return err != nil }(); i++ {
		time.Sleep(25 * time.Millisecond)
	}
	args := []string{"memory", "write", "--expect", os.Getenv("DH_MEMORY_EXPECT"), "--from", os.Getenv("DH_MEMORY_FROM"), os.Getenv("DH_MEMORY_REPO")}
	if file := os.Getenv("DH_MEMORY_FILE"); file != "" {
		args = append([]string{"memory", "write", "--file", file}, args[2:]...)
	}
	code := run(args, strings.NewReader(""), os.Stdout, os.Stderr)
	os.Exit(code)
}

func TestMemoryConcurrentProcesses(t *testing.T) {
	const writers = 8
	repo, mem := memRepo(t, "start")
	home := os.Getenv("DH_HOME")
	goFile := filepath.Join(t.TempDir(), "go")
	var cmds []*exec.Cmd
	var contents []string
	for i := 0; i < writers; i++ {
		c := strings.Repeat(fmt.Sprintf("process-%d\n", i), 300)
		contents = append(contents, c)
		draft := filepath.Join(t.TempDir(), fmt.Sprintf("draft-%d.md", i))
		os.WriteFile(draft, []byte(c), 0o644)
		cmd := exec.Command(os.Args[0], "-test.run=TestMemoryHelperProcess")
		cmd.Env = append(os.Environ(), "DH_MEMORY_HELPER=1", "DH_HOME="+home, "DH_MEMORY_GO="+goFile, "DH_MEMORY_EXPECT="+shaOf("start"),
			"DH_MEMORY_FROM="+draft, "DH_MEMORY_REPO="+repo)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		cmds = append(cmds, cmd)
	}
	os.WriteFile(goFile, nil, 0o644) // release them together
	wins, conflicts := 0, 0
	for _, cmd := range cmds {
		err := cmd.Wait()
		switch {
		case err == nil:
			wins++
		default:
			if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 3 {
				conflicts++
			} else {
				t.Fatalf("unexpected exit: %v", err)
			}
		}
	}
	if wins != 1 || conflicts != writers-1 {
		t.Fatalf("wins=%d conflicts=%d, want 1 and %d", wins, conflicts, writers-1)
	}
	got, _ := os.ReadFile(mem)
	ok := false
	for _, c := range contents {
		if string(got) == c {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("the final file is not the content of any single writer (%d bytes)", len(got))
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(mem), "*.lock")); len(left) != 0 {
		t.Fatalf("lock left: %v", left)
	}
	if left, _ := filepath.Glob(filepath.Join(filepath.Dir(mem), ".MEMORY.md.tmp-*")); len(left) != 0 {
		t.Fatalf("temp left: %v", left)
	}
}
