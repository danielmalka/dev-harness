# Diff under review

Fictional change, used as a fixture only. Scope for the `code` review stage
in this case, spanning two unrelated files at unrelated lines so two
reviewers' findings cannot plausibly refer to the same location — the
negative counterpart to `external-clis-attribution`'s single-location
positive-merge coverage.

- Scope: `internal/service/retry.go`, `internal/cache/lru.go`
- Comparison: `main..fix/retry-and-cache`

```diff
--- a/internal/service/retry.go
+++ b/internal/service/retry.go
@@ -38,13 +38,12 @@ func WithRetry(ctx context.Context, fn func() error) error {
 	var lastErr error
-	for attempt := 0; ; attempt++ {
+	for attempt := 0; attempt < maxAttempts; attempt++ {
 		lastErr = fn()
 		if lastErr == nil {
 			return nil
 		}
-		if attempt >= maxAttempts {
-			return lastErr
-		}
 		time.Sleep(backoff(attempt))
 	}
+	return lastErr
 }
```

```diff
--- a/internal/cache/lru.go
+++ b/internal/cache/lru.go
@@ -80,11 +80,10 @@ func (c *LRU) Set(key string, value []byte) {
 	c.mu.Lock()
 	defer c.mu.Unlock()
-	if len(c.entries) >= c.capacity {
-		c.evictOldest()
-	}
 	c.entries[key] = value
+	c.touch(key)
 }
```

The loop in `retry.go` now exits through the `for` condition instead of an
explicit break, and the trailing `time.Sleep(backoff(attempt))` no longer
runs before the function returns on the last attempt. Separately and
unrelatedly, `lru.go`'s `Set` no longer evicts the oldest entry before
inserting a new one when the cache is already at capacity, so the cache
grows without bound.
