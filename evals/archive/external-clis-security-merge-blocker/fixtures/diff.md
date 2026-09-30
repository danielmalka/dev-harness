# Diff under review

Fictional change, used as a fixture only. Scope for the `security` review
stage of `external-clis-security-merge-blocker`. Deliberately vulnerable:
the lookup filters by id alone, with no check that the invoice belongs to
the caller's own account, to exercise the merge rule in
`.agents/coordinator.md` ("External CLI reviewers").

- Scope: `src/api/invoices.ts`, `src/repo/invoices.ts`
- Comparison: `main..feat/invoice-lookup`

```diff
--- a/src/api/invoices.ts
+++ b/src/api/invoices.ts
@@ -18,6 +18,9 @@ router.get('/invoices/:id', requireSession, async (req, res) => {
+  const invoice = await repository.findById(req.params.id);
+  if (!invoice) return res.status(404).end();
+  res.json(invoice);
 });
```

The handler requires a session (`requireSession`) but looks the invoice up
by id alone, with no filter on the caller's own account.
