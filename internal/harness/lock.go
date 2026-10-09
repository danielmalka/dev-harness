package harness

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	lockWait  = 5 * time.Second
	lockStale = 30 * time.Second
)

// lock takes <home>/config.lock (O_CREATE|O_EXCL), retrying every 50 ms up to lockWait. A lock
// older than lockStale is removed once. The returned func releases it.
func lock(home string) (func(), error) {
	if err := os.MkdirAll(home, dirPerm); err != nil {
		return nil, err
	}
	path := filepath.Join(home, "config.lock")
	deadline := time.Now().Add(lockWait)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	nonce := hex.EncodeToString(b)
	removedStale := false
	for {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
		if err == nil {
			_, werr := f.WriteString(nonce)
			f.Close()
			if werr != nil {
				os.Remove(path)
				return nil, werr
			}
			return func() { // remove only our own lock (a stale-takeover may have replaced it)
				if got, err := os.ReadFile(path); err == nil && string(got) == nonce {
					os.Remove(path)
				}
			}, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if !removedStale && staleLock(path) {
			removedStale = true
			os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for %s (another dh link running? remove the file if not)", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// staleLock: older than lockStale, with mtime and content unchanged across two reads.
// ponytail: a holder that replaces the file between the second read and Remove still loses its
// lock; closing that needs flock/LockFileEx. Fine for a 30 s stale window on a local registry.
func staleLock(path string) bool {
	fi, err := os.Stat(path)
	if err != nil || time.Since(fi.ModTime()) <= lockStale {
		return false
	}
	c1, _ := os.ReadFile(path)
	fi2, err := os.Stat(path)
	c2, _ := os.ReadFile(path)
	return err == nil && fi2.ModTime().Equal(fi.ModTime()) && string(c1) == string(c2)
}
