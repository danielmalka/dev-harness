package snapshot

import "os"

// ponytail: Windows keeps os.Open; the f.Stat regular check in ReadCapped is the only extra guard.
func openNoFollow(path string) (*os.File, error) { return os.Open(path) }

// openIn: Lstat the directory (no symlink), open dir/name, and check the directory is still the same one.
// ponytail: Windows has no openat in the stdlib; a swap between these steps is not excluded here. The
// dashboard serves only the user's own files on loopback.
func openIn(dir, name string) (*os.File, error) {
	di, err := os.Lstat(dir)
	if err != nil {
		return nil, err
	}
	if !di.IsDir() {
		return nil, ErrNotRegular
	}
	f, err := OpenRegular(dir + string(os.PathSeparator) + name) // Lstat + SameFile: a final-file symlink is refused
	if err != nil {
		return nil, err
	}
	if now, err := os.Lstat(dir); err != nil || !os.SameFile(di, now) {
		f.Close()
		return nil, ErrNotRegular
	}
	return f, nil
}
