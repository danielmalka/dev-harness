package dashboard

import (
	"os"
	"sync"
	"time"

	"github.com/danielmalka/dev-harness/internal/snapshot"
)

// progressCache keeps what projectProgress parsed from a ticket or PRD file while its path, size and
// modification time are unchanged. It is consulted only after the Lstat, regular-file and size checks that
// snapshot.ReadCapped applies, so a refused file is refused with or without it. In memory only; safe for
// concurrent builds. ponytail: mtime granularity is the staleness ceiling (a same-size rewrite inside one
// timestamp tick, or a preserved mtime, is served stale until the next change); add a content hash if that bites.
type progressCache struct {
	mu   sync.Mutex
	m    map[string]cacheEntry
	seen map[string]struct{}
}

type cacheEntry struct {
	size int64
	mod  time.Time
	v    any
}

func newProgressCache() *progressCache {
	return &progressCache{m: map[string]cacheEntry{}, seen: map[string]struct{}{}}
}

// sweep drops the entries for paths not read since the last sweep (call once at the end of a state build).
func (c *progressCache) sweep() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for p := range c.m {
		if _, ok := c.seen[p]; !ok {
			delete(c.m, p)
		}
	}
	c.seen = map[string]struct{}{}
}

// readParsed reads path with readCapped's guards and returns parse(content), from the cache when the file is unchanged.
// A nil cache reads and parses every time.
func readParsed[T any](c *progressCache, path string, parse func(md string) T) (T, error) {
	var zero T
	if c == nil {
		b, err := readCapped(path)
		if err != nil {
			return zero, err
		}
		return parse(string(b)), nil
	}
	info, err := os.Lstat(path)
	if err != nil {
		return zero, err
	}
	if !info.Mode().IsRegular() {
		return zero, snapshot.ErrNotRegular
	}
	if info.Size() > snapshot.MaxFileBytes {
		return zero, snapshot.ErrTooLarge
	}
	c.mu.Lock()
	e, ok := c.m[path]
	c.mu.Unlock()
	if ok && e.size == info.Size() && e.mod.Equal(info.ModTime()) {
		if v, isT := e.v.(T); isT {
			c.mu.Lock()
			c.seen[path] = struct{}{}
			c.mu.Unlock()
			return v, nil
		}
	}
	b, err := readCapped(path)
	if err != nil {
		return zero, err
	}
	v := parse(string(b))
	c.mu.Lock()
	c.m[path] = cacheEntry{info.Size(), info.ModTime(), v}
	c.seen[path] = struct{}{}
	c.mu.Unlock()
	return v, nil
}
