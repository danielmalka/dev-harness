//go:build !windows

package snapshot

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestReadCappedRefusesFIFOSwappedAfterLstat(t *testing.T) {
	p := filepath.Join(t.TempDir(), "s.json")
	if err := os.WriteFile(p, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	afterLstat = func() {
		_ = os.Remove(p)
		if err := syscall.Mkfifo(p, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	defer func() { afterLstat = nil }()
	done := make(chan error, 1)
	go func() { _, err := ReadCapped(p); done <- err }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("swapped FIFO was read")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ReadCapped blocked on a FIFO")
	}
}

func TestReadCappedRefusesSymlinkSwappedAfterLstat(t *testing.T) {
	p := filepath.Join(t.TempDir(), "s.json")
	if err := os.WriteFile(p, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	afterLstat = func() { _ = os.Remove(p); _ = os.Symlink("/dev/zero", p) }
	defer func() { afterLstat = nil }()
	if _, err := ReadCapped(p); err == nil {
		t.Fatal("swapped symlink was read")
	}
}

func TestReadCappedIn(t *testing.T) {
	dir, real := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadCappedIn(dir, "a", 10); err != nil || string(b) != "hi" {
		t.Fatalf("normal: %q %v", b, err)
	}
	if _, err := ReadCappedIn(dir, "a", 1); err != ErrTooLarge {
		t.Fatalf("cap: %v", err)
	}
	if _, err := ReadCappedIn(dir, "../a", 10); err != ErrNotRegular {
		t.Fatalf("separator: %v", err)
	}
	if err := os.Symlink(filepath.Join(dir, "a"), filepath.Join(dir, "l")); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCappedIn(dir, "l", 10); err == nil {
		t.Fatal("final symlink read")
	}
	if err := os.WriteFile(filepath.Join(real, "a"), []byte("out"), 0o644); err != nil {
		t.Fatal(err)
	}
	ld := filepath.Join(t.TempDir(), "ld")
	if err := os.Symlink(real, ld); err != nil {
		t.Fatal(err)
	}
	if b, err := ReadCappedIn(ld, "a", 10); err == nil {
		t.Fatalf("symlinked dir read: %q", b)
	}
	if _, err := OpenRegularIn(ld, "a"); err == nil {
		t.Fatal("symlinked dir streamed")
	}
}

func TestReadCappedInFIFOAsDirFailsFast(t *testing.T) {
	p := filepath.Join(t.TempDir(), "d")
	if err := syscall.Mkfifo(p, 0o644); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, err := ReadCappedIn(p, "a", 10); done <- err }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("FIFO dir read")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocked on a FIFO in place of the dir")
	}
}
