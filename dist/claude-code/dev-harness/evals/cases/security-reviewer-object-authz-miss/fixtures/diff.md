# Diff under review

Fictional change, used as a fixture only. Scope for the security review case
`security-reviewer-object-authz-miss`.

- Scope: `src/api/notes.ts`, `src/repo/notes.ts`
- Comparison: `main..feat/notes-export`

```diff
--- /dev/null
+++ b/src/repo/notes.ts
@@ -0,0 +1,13 @@
+import { db } from '../db';
+
+export async function findById(id: string) {
+  return db.notes.findOne({ id });
+}
+
+export async function findByIdAndOwner(id: string, ownerId: string) {
+  return db.notes.findOne({ id, ownerId });
+}
+
+export async function deleteById(id: string) {
+  return db.notes.deleteOne({ id });
+}
--- a/src/api/notes.ts
+++ b/src/api/notes.ts
@@ -22,7 +22,38 @@ router.delete('/notes/:id', requireSession, async (req, res) => {
+router.get('/notes/:id/export', requireSession, async (req, res) => {
+  if (!/^[a-f0-9-]{36}$/.test(req.params.id)) {
+    return res.status(400).json({ error: 'invalid note id' });
+  }
+  const note = await repository.findById(req.params.id);
+  if (!note) return res.status(404).end();
+  // A caller must own a note before it can be exported. `defaultNoteId` is
+  // the note created for every account at signup, so this reuses the same
+  // ownership check used elsewhere for that resource.
+  const ownerNote = await repository.findByIdAndOwner(req.session.defaultNoteId, req.session.userId);
+  if (!ownerNote) return res.status(403).end();
+  auditLog.record('note.export', { userId: req.session.userId, noteId: note.id });
+  res.setHeader('Content-Type', 'text/markdown');
+  res.send(note.body);
+});
+
 router.delete('/notes/:id', requireSession, async (req, res) => {
-  await repository.deleteById(req.params.id);
+  if (!/^[a-f0-9-]{36}$/.test(req.params.id)) {
+    return res.status(400).json({ error: 'invalid note id' });
+  }
+  const note = await repository.findByIdAndOwner(req.params.id, req.session.userId);
+  if (!note) return res.status(404).end();
+  await repository.deleteById(note.id);
+  auditLog.record('note.delete', { userId: req.session.userId, noteId: note.id });
   res.status(204).end();
 });
```

`repository` and `auditLog` are already-imported modules elsewhere in
`src/api/notes.ts`; this diff does not change those imports. `req.session`
carries `userId` (the authenticated caller) and `defaultNoteId` (the note
created automatically for every account at signup). Both routes validate
that `:id` is a UUID before doing anything else, and both write to
`auditLog` after acting.
