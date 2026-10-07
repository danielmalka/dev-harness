package snapshot

import (
	"io"
	"os"
)

// afterLstat is a test hook that runs between Lstat and Open (nil in production).
var afterLstat func()

// ReadCapped reads a regular file of at most MaxFileBytes. Lstat rejects symlinks and special files (a /dev/zero
// link reports size 0); the open does not follow symlinks or block on a FIFO (unix); the opened file must still be
// the same regular file Lstat saw, so a swap in between is refused; and the read itself is bounded.
func ReadCapped(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, ErrNotRegular
	}
	if info.Size() > MaxFileBytes {
		return nil, ErrTooLarge
	}
	if afterLstat != nil {
		afterLstat()
	}
	f, err := openNoFollow(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() || !os.SameFile(info, fi) {
		return nil, ErrNotRegular
	}
	b, err := io.ReadAll(io.LimitReader(f, MaxFileBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > MaxFileBytes {
		return nil, ErrTooLarge
	}
	return b, nil
}
