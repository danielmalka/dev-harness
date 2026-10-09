//go:build !linux && !windows

package snapshot

import (
	"os"
	"syscall"
)

// openIn: the stdlib syscall package has no Openat outside Linux, so open the directory fd (no symlink), open
// dir/name by path with O_NOFOLLOW, then require that the directory path still names the same directory.
// ponytail: this detects a directory swap that is still in place after the open; a swap-and-restore inside that
// window is not caught here (Linux uses openat and has no window). Upgrade: golang.org/x/sys/unix Openat.
func openIn(dir, name string) (*os.File, error) {
	dfd, err := syscall.Open(dir, syscall.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: dir, Err: err}
	}
	df := os.NewFile(uintptr(dfd), dir)
	defer df.Close()
	di, err := df.Stat()
	if err != nil || !di.IsDir() {
		return nil, ErrNotRegular
	}
	f, err := openNoFollow(dir + string(os.PathSeparator) + name)
	if err != nil {
		return nil, err
	}
	if now, err := os.Lstat(dir); err != nil || !os.SameFile(di, now) {
		f.Close()
		return nil, ErrNotRegular
	}
	return f, nil
}
