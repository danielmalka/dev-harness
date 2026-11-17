package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func sha(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

func memFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	if content != "" {
		if err := os.WriteFile(filepath.Join(dir, MemoryFile), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func leftovers(t *testing.T, dir string) []string {
	t.Helper()
	ents, _ := os.ReadDir(dir)
	var left []string
	for _, e := range ents {
		if e.Name() != MemoryFile {
			left = append(left, e.Name())
		}
	}
	return left
}

// R1: the digest is the sha256 of the file, or "absent".
func TestMemoryHash(t *testing.T) {
	dir := memFile(t, "# memory\nline\n")
	if got, err := MemoryHash(dir); err != nil || got != sha("# memory\nline\n") {
		t.Fatalf("hash = %q, %v", got, err)
	}
	if got, err := MemoryHash(memFile(t, "")); err != nil || got != MemoryAbsent {
		t.Fatalf("absent = %q, %v", got, err)
	}
}

// R2: equal hash writes the new content, returns the new digest and leaves no temp or lock file.
func TestWriteMemoryEqualHash(t *testing.T) {
	dir := memFile(t, "old")
	got, err := WriteMemory(dir, sha("old"), []byte("new content"))
	if err != nil || got != sha("new content") {
		t.Fatalf("write = %q, %v", got, err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, MemoryFile)); string(b) != "new content" {
		t.Fatalf("file = %q", b)
	}
	if left := leftovers(t, dir); len(left) != 0 {
		t.Fatalf("leftovers: %v", left)
	}
}

// R2: "absent" creates the first MEMORY.md; R3: "absent" on an existing file is a conflict.
func TestWriteMemoryAbsent(t *testing.T) {
	dir := memFile(t, "")
	if _, err := WriteMemory(dir, MemoryAbsent, []byte("first")); err != nil {
		t.Fatal(err)
	}
	_, err := WriteMemory(dir, MemoryAbsent, []byte("second"))
	if !IsMemoryChanged(err) {
		t.Fatalf("want conflict, got %v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, MemoryFile)); string(b) != "first" {
		t.Fatalf("file = %q", b)
	}
}

// R3: a different hash writes nothing, keeps the other writer's content and names both digests and the cause.
func TestWriteMemoryChangedInTheMiddle(t *testing.T) {
	dir := memFile(t, "v1")
	read := sha("v1")                                                                                         // the caller read v1 ...
	if err := os.WriteFile(filepath.Join(dir, MemoryFile), []byte("v2 by someone else"), 0o644); err != nil { // ... another session wrote v2
		t.Fatal(err)
	}
	_, err := WriteMemory(dir, read, []byte("my edit of v1"))
	if !IsMemoryChanged(err) {
		t.Fatalf("want conflict, got %v", err)
	}
	for _, want := range []string{read, sha("v2 by someone else"), "alguém escreveu no meio", "Nada foi gravado"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("message %q lacks %q", err, want)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(dir, MemoryFile)); string(b) != "v2 by someone else" {
		t.Fatalf("the other writer's content was lost: %q", b)
	}
	if left := leftovers(t, dir); len(left) != 0 {
		t.Fatalf("leftovers: %v", left)
	}
}

func TestWriteMemoryRefusesBadInput(t *testing.T) {
	dir := memFile(t, "x")
	for name, c := range map[string]struct {
		expect  string
		content string
	}{
		"empty content": {sha("x"), ""},
		"short expect":  {"abc", "y"},
		"upper-case":    {strings.ToUpper(sha("x")), "y"},
		"empty expect":  {"", "y"},
	} {
		if _, err := WriteMemory(dir, c.expect, []byte(c.content)); err == nil || IsMemoryChanged(err) {
			t.Errorf("%s: want a plain error, got %v", name, err)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(dir, MemoryFile)); string(b) != "x" {
		t.Fatalf("file changed: %q", b)
	}
}

// R4: N writers read the same hash and race with different contents: exactly one wins, the others get the conflict,
// the file is the winner's byte for byte, and no temp or lock file is left.
func TestWriteMemoryConcurrentWriters(t *testing.T) {
	const n = 16
	for round := 0; round < 5; round++ {
		dir := memFile(t, "start")
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
				c := strings.Repeat("writer-"+string(rune('a'+i))+"\n", 200)
				<-start
				_, err := WriteMemory(dir, expect, []byte(c))
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
		if b, _ := os.ReadFile(filepath.Join(dir, MemoryFile)); string(b) != winner {
			t.Fatalf("round %d: the file is not the winner's content", round)
		}
		if left := leftovers(t, dir); len(left) != 0 {
			t.Fatalf("round %d: leftovers: %v", round, left)
		}
	}
}
