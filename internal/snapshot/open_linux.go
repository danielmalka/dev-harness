//go:build linux

package snapshot

import (
	"os"
	"syscall"
)

// openIn opens dir/name relative to a directory fd opened without following symlinks, so a swap of dir for a
// symlink after the check cannot redirect the read; name is opened with O_NOFOLLOW and never blocks on a FIFO.
func openIn(dir, name string) (*os.File, error) {
	dfd, err := syscall.Open(dir, syscall.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: dir, Err: err}
	}
	defer syscall.Close(dfd)
	fd, err := syscall.Openat(dfd, name, syscall.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	return os.NewFile(uintptr(fd), name), nil
}
