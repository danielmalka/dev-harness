package snapshot

import (
	"io"
	"os"
	"strings"
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

// OpenRegular opens a regular file for streaming with the same guards as ReadCapped (no symlink, no special
// file, same file as Lstat saw) but without the size cap. The caller closes it.
func OpenRegular(path string) (*os.File, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, ErrNotRegular
	}
	f, err := openNoFollow(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil || !fi.Mode().IsRegular() || !os.SameFile(info, fi) {
		f.Close()
		return nil, ErrNotRegular
	}
	return f, nil
}

// OpenRegularIn opens dir/name for streaming through openIn (the directory is not followed if it is a symlink,
// and the file is opened relative to it), and requires a regular file. name must be a single path element.
func OpenRegularIn(dir, name string) (*os.File, error) {
	if name == "" || strings.ContainsAny(name, `/\`) || name == ".." {
		return nil, ErrNotRegular
	}
	f, err := openIn(dir, name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil || !fi.Mode().IsRegular() {
		f.Close()
		return nil, ErrNotRegular
	}
	return f, nil
}

// ReadCappedIn is ReadCapped for dir/name with the directory-safe open of OpenRegularIn and an explicit cap.
func ReadCappedIn(dir, name string, max int64) ([]byte, error) {
	f, err := OpenRegularIn(dir, name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > max {
		return nil, ErrTooLarge
	}
	return b, nil
}
