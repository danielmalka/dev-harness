package snapshot

import "os"

// ponytail: Windows keeps os.Open; the f.Stat regular check in ReadCapped is the only extra guard.
func openNoFollow(path string) (*os.File, error) { return os.Open(path) }
