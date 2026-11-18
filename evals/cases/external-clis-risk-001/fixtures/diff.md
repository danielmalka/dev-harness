# Diff em revisão

Mudança fictícia, usada só como fixture. Escopo da etapa de revisão `code` em todo caso `external-clis-*`.

- Escopo: `internal/service/retry.go`
- Comparação: `main..fix/retry-backoff`

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

O laço agora sai pela condição do `for` em vez de um break explícito, e o `time.Sleep(backoff(attempt))` final deixa de rodar antes de a função retornar na última tentativa.
