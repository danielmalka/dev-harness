package harness

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func factsDir(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	if content != "" {
		if err := os.WriteFile(filepath.Join(dir, FactsFile), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func leftoversExcept(t *testing.T, dir, keep string) []string {
	t.Helper()
	ents, _ := os.ReadDir(dir)
	var left []string
	for _, e := range ents {
		if e.Name() != keep {
			left = append(left, e.Name())
		}
	}
	return left
}

// R11: the closed list has exactly the two records, and any other name is refused with both valid names in the message.
func TestRecordFilesClosedList(t *testing.T) {
	if len(RecordFiles) != 2 || RecordFiles[0] != "MEMORY.md" || RecordFiles[1] != "FACTS.md" {
		t.Fatalf("RecordFiles = %v", RecordFiles)
	}
	for _, name := range []string{"MEMORY.md", "FACTS.md"} {
		if err := CheckRecordFile(name); err != nil {
			t.Errorf("%s refused: %v", name, err)
		}
	}
	for _, name := range []string{"", "EPOCHAL.md", "RISKS.md", "REPO-MAP.md", "project.yaml", "facts.md", "memory.md", "../FACTS.md", "sub/FACTS.md", "FACTS.md ", "MEMORY.md.lock"} {
		err := CheckRecordFile(name)
		if err == nil {
			t.Errorf("%q must be refused", name)
			continue
		}
		for _, want := range []string{"MEMORY.md", "FACTS.md"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%q: message %q lacks %q", name, err, want)
			}
		}
	}
}

// R11: neither FileHash nor WriteFile touches a name outside the list: no file, no lock and no temp file is created.
func TestFileFunctionsRefuseNamesOutsideTheList(t *testing.T) {
	dir := t.TempDir()
	if _, err := FileHash(dir, "EPOCHAL.md"); err == nil {
		t.Error("FileHash(EPOCHAL.md) must fail")
	}
	for _, name := range []string{"EPOCHAL.md", "RISKS.md", "../escape.md", "project.yaml"} {
		if _, err := WriteFile(dir, name, MemoryAbsent, []byte("x")); err == nil || IsMemoryChanged(err) {
			t.Errorf("WriteFile(%q) must fail with a plain error, got %v", name, err)
		}
	}
	if left := leftoversExcept(t, dir, ""); len(left) != 0 {
		t.Fatalf("a refused name left files behind: %v", left)
	}
}

// R11: hash and write on FACTS.md follow the same rules as MEMORY.md, and the lock is FACTS.md.lock (nothing left behind).
func TestFactsHashAndWrite(t *testing.T) {
	dir := factsDir(t, "")
	if got, err := FileHash(dir, FactsFile); err != nil || got != MemoryAbsent {
		t.Fatalf("absent = %q, %v", got, err)
	}
	if _, err := WriteFile(dir, FactsFile, MemoryAbsent, []byte("# Fatos\n")); err != nil {
		t.Fatal(err)
	}
	if got, err := FileHash(dir, FactsFile); err != nil || got != sha("# Fatos\n") {
		t.Fatalf("hash = %q, %v", got, err)
	}
	got, err := WriteFile(dir, FactsFile, sha("# Fatos\n"), []byte("# Fatos\nmais\n"))
	if err != nil || got != sha("# Fatos\nmais\n") {
		t.Fatalf("write = %q, %v", got, err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, FactsFile)); string(b) != "# Fatos\nmais\n" {
		t.Fatalf("file = %q", b)
	}
	if left := leftoversExcept(t, dir, FactsFile); len(left) != 0 {
		t.Fatalf("leftovers: %v", left)
	}
	if _, err := os.Stat(filepath.Join(dir, MemoryFile)); err == nil {
		t.Fatal("writing FACTS.md created MEMORY.md")
	}
}

// R10: a different hash writes nothing and the message names FACTS.md, not MEMORY.md.
func TestFactsWriteConflictNamesTheFile(t *testing.T) {
	dir := factsDir(t, "v1")
	read := sha("v1")
	if err := os.WriteFile(filepath.Join(dir, FactsFile), []byte("v2 by someone else"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := WriteFile(dir, FactsFile, read, []byte("my edit of v1"))
	if !IsMemoryChanged(err) {
		t.Fatalf("want conflict, got %v", err)
	}
	for _, want := range []string{"FACTS.md mudou desde que foi lido", "Leia o FACTS.md de novo", read, sha("v2 by someone else"), "Nada foi gravado"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("message %q lacks %q", err, want)
		}
	}
	if strings.Contains(err.Error(), "MEMORY.md") {
		t.Errorf("message %q names MEMORY.md", err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, FactsFile)); string(b) != "v2 by someone else" {
		t.Fatalf("the other writer's content was lost: %q", b)
	}
	if _, err := WriteFile(dir, FactsFile, MemoryAbsent, []byte("x")); !IsMemoryChanged(err) {
		t.Fatalf("absent on an existing FACTS.md must conflict, got %v", err)
	}
	if _, err := WriteFile(dir, FactsFile, sha("v2 by someone else"), nil); err == nil || IsMemoryChanged(err) || !strings.Contains(err.Error(), "FACTS.md vazio") {
		t.Fatalf("empty content must be a plain error naming FACTS.md, got %v", err)
	}
}

// The two records are independent: a MEMORY.md write does not change the FACTS.md hash, and the other way round, and a held FACTS.md lock does not block MEMORY.md.
func TestFactsAndMemoryAreIndependent(t *testing.T) {
	dir := t.TempDir()
	if _, err := WriteFile(dir, FactsFile, MemoryAbsent, []byte("facts")); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteMemory(dir, MemoryAbsent, []byte("memory")); err != nil {
		t.Fatalf("MEMORY.md write after FACTS.md: %v", err)
	}
	if h, _ := FileHash(dir, FactsFile); h != sha("facts") {
		t.Fatalf("FACTS.md hash changed by the MEMORY.md write: %s", h)
	}
	release, err := lockNamed(dir, FactsFile+".lock")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	done := make(chan error, 1)
	go func() { _, err := WriteMemory(dir, sha("memory"), []byte("memory 2")); done <- err }()
	if err := <-done; err != nil {
		t.Fatalf("MEMORY.md write blocked or failed while FACTS.md.lock is held: %v", err)
	}
}

// R4 of PRD-034 for FACTS.md: N writers read the same hash and race: exactly one wins, the others get the conflict, the file is the winner's byte for byte.
func TestFactsConcurrentWriters(t *testing.T) {
	const n = 16
	for round := 0; round < 5; round++ {
		dir := factsDir(t, "start")
		expect := sha("start")
		start := make(chan struct{})
		type result struct {
			content string
			err     error
		}
		results := make([]result, n)
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c := strings.Repeat("candidate-"+string(rune('a'+i))+"\n", 100)
				<-start
				_, err := WriteFile(dir, FactsFile, expect, []byte(c))
				results[i] = result{c, err}
			}(i)
		}
		close(start)
		wg.Wait()
		winners, conflicts := 0, 0
		var winner string
		for _, r := range results {
			switch {
			case r.err == nil:
				winners++
				winner = r.content
			case IsMemoryChanged(r.err):
				conflicts++
			default:
				t.Fatalf("unexpected error: %v", r.err)
			}
		}
		if winners != 1 || conflicts != n-1 {
			t.Fatalf("round %d: winners=%d conflicts=%d, want 1 and %d", round, winners, conflicts, n-1)
		}
		if b, _ := os.ReadFile(filepath.Join(dir, FactsFile)); string(b) != winner {
			t.Fatalf("round %d: the file is not the winner's content", round)
		}
		if left := leftoversExcept(t, dir, FactsFile); len(left) != 0 {
			t.Fatalf("round %d: leftovers: %v", round, left)
		}
	}
}
