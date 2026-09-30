# Diff under review

Fictional change, used as a fixture only. Scope for the `code` review stage in every `external-clis-*` case.

- Scope: `internal/service/retry.go`
- Comparison: `main..fix/retry-backoff`

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

The loop now exits through the `for` condition instead of an explicit break, and the trailing `time.Sleep(backoff(attempt))` no longer runs before the function returns on the last attempt.
