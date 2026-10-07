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
