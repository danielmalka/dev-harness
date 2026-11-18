# Diff em revisão

Mudança fictícia, usada apenas como fixture. Escopo do caso de revisão de segurança
`security-reviewer-object-authz-miss`.

- Escopo: `src/api/notes.ts`, `src/repo/notes.ts`
- Comparação: `main..feat/notes-export`

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

`repository` e `auditLog` são módulos já importados em outro ponto de
`src/api/notes.ts`; este diff não muda esses imports. `req.session`
traz `userId` (o chamador autenticado) e `defaultNoteId` (a nota criada
automaticamente para toda conta no cadastro). As duas rotas validam
que `:id` é um UUID antes de fazer qualquer outra coisa, e as duas gravam em
`auditLog` depois de agir.
