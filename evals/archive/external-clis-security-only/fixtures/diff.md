# Diff under review

Fictional change, used as a fixture only. Scope for the `security` review
stage of `external-clis-security-only`.

- Scope: `src/api/invoices.ts`, `src/repo/invoices.ts`
- Comparison: `main..feat/invoice-lookup`

```diff
--- a/src/api/invoices.ts
+++ b/src/api/invoices.ts
@@ -18,6 +18,10 @@ router.get('/invoices/:id', requireSession, async (req, res) => {
+  const invoice = await repository.findByIdForAccount(
+    req.params.id,
+    req.session.accountId,
+  );
+  if (!invoice) return res.status(404).end();
+  res.json(invoice);
 });
```

The handler requires a session (`requireSession`) and filters the lookup by
the caller's own `accountId` before returning the invoice.
