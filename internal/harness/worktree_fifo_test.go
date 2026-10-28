//go:build !windows

package harness

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// R26 hardening: each small pointer file (gitfile, commondir, back-pointer gitdir) that is a FIFO
// is refused without blocking. Limit: readSmall Lstats first, so a FIFO placed at the path is
// refused by the Lstat check; the O_NONBLOCK open (FIFO swapped in AFTER Lstat) needs a race
// that cannot be hit deterministically, so that line is not covered by this test.
// (Separate file: syscall.Mkfifo does not exist on windows, so the file is !windows.)
func TestR26_FifoPointerFilesDoNotHang(t *testing.T) {
	for _, name := range []string{"gitfile", "commondir", "back-pointer"} {
		t.Run(name, func(t *testing.T) {
			f := newWT(t)
			p := map[string]string{
				"gitfile":      filepath.Join(f.wt, ".git"),
				"commondir":    filepath.Join(f.admin, "commondir"),
				"back-pointer": filepath.Join(f.admin, "gitdir"),
			}[name]
			if err := os.Remove(p); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(p, 0o600); err != nil {
				t.Skip("mkfifo unavailable: ", err)
			}
			done := make(chan *LinkedWT, 1)
			go func() { done <- LinkedWorktree(f.wt) }()
			select {
			case got := <-done:
				if got != nil {
					t.Errorf("accepted: %+v", *got)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("LinkedWorktree blocked on a FIFO")
			}
			// readSmall directly: same refusal, also bounded.
			rd := make(chan bool, 1)
			go func() { _, ok := readSmall(p); rd <- ok }()
			select {
			case ok := <-rd:
				if ok {
					t.Error("readSmall accepted a FIFO")
				}
			case <-time.After(3 * time.Second):
				t.Fatal("readSmall blocked on a FIFO")
			}
		})
	}
}
