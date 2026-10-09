//go:build windows

package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadCappedInRefusesFinalSymlinkWindows(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "a"), filepath.Join(dir, "l")); err != nil {
		t.Skip("symlinks not available")
	}
	if _, err := ReadCappedIn(dir, "l", 10); err == nil {
		t.Fatal("final symlink read")
	}
}
