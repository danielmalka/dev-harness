# Diff under review

Fictional change, used as a fixture only. Scope for the `code` review case
`code-reviewer-equivalent-form-miss`.

- Scope: `scripts/ci-gate.js`, `scripts/workspaces.config.js`
- Comparison: `main..feat/ci-gate-classify`

```diff
--- /dev/null
+++ b/scripts/workspaces.config.js
@@ -0,0 +1,21 @@
+// Canonical CI pipeline commands per workspace. These are the exact
+// invocations `npm run ci` uses today (see the root `package.json`
+// `"ci"` script and each workspace's own `package.json`), not a
+// normalized or idealized form.
+module.exports = {
+  root: {
+    test: 'npm test',
+    lint: 'npm run lint',
+    build: 'npm run build',
+  },
+  server: {
+    test: 'npm run test --workspace=server',
+    lint: 'cd server && npm run lint',
+    build: 'cd server && npm run build',
+  },
+  worker: {
+    test: 'cd worker && npm test',
+    lint: 'cd worker && npm run lint',
+    build: 'cd worker && npm run build',
+  },
+};
--- a/scripts/ci-gate.js
+++ b/scripts/ci-gate.js
@@ -10,6 +10,70 @@
+const workspaces = require('./workspaces.config');
+
+const MAX_RETRIES = 2;
+
+// Strips a leading `cd <dir> &&` prefix so a classifier can match the bare
+// command regardless of which workspace directory a step runs from.
+function normalizeCommand(cmd) {
+  const match = /^cd [^&]+&&\s*(.+)$/.exec(cmd);
+  return match ? match[1].trim() : cmd;
+}
+
+function classifyStep(cmd) {
+  const normalized = normalizeCommand(cmd);
+  if (normalized === 'npm test') return 'test';
+  return 'generic';
+}
+
+function classifyLintStep(cmd) {
+  const normalized = normalizeCommand(cmd);
+  if (normalized === 'npm run lint') return 'lint';
+  return 'generic';
+}
+
+function classifyBuildStep(cmd) {
+  const normalized = normalizeCommand(cmd);
+  if (normalized === 'npm run build') return 'build';
+  return 'generic';
+}
+
+function classify(cmd) {
+  if (cmd.includes('lint')) return classifyLintStep(cmd);
+  if (cmd.includes('build')) return classifyBuildStep(cmd);
+  return classifyStep(cmd);
+}
+
+function collectSteps() {
+  const steps = [];
+  for (const [workspace, commands] of Object.entries(workspaces)) {
+    steps.push({ workspace, cmd: commands.test });
+    steps.push({ workspace, cmd: commands.lint });
+    steps.push({ workspace, cmd: commands.build });
+  }
+  return steps;
+}
+
+function runGate() {
+  const started = Date.now();
+  const results = collectSteps().map(({ workspace, cmd }) => {
+    const kind = classify(cmd);
+    // A 'generic' step's exit code alone gates the merge; a 'test' step's
+    // output is also parsed for the framework's own pass/fail summary, so a
+    // real failure buried in passing exit-code noise still blocks the merge.
+    return { workspace, cmd, kind, retriesLeft: MAX_RETRIES };
+  });
+  const summary = results.reduce((acc, r) => {
+    acc[r.kind] = (acc[r.kind] || 0) + 1;
+    return acc;
+  }, {});
+  console.log(`Gate summary: ${JSON.stringify(summary)} (${Date.now() - started}ms)`);
+  return results;
+}
+
+module.exports = {
+  classify,
+  classifyStep,
+  classifyLintStep,
+  classifyBuildStep,
+  runGate,
+  normalizeCommand,
+};
```

`classify` dispatches each pipeline step to `classifyStep`, `classifyLintStep`,
or `classifyBuildStep` by keyword match on the command string, and each
classifier compares the command — after `normalizeCommand` strips a leading
`cd <dir> &&` prefix — against the bare npm invocation for its step kind. The
commands actually run per workspace live in `scripts/workspaces.config.js`,
one `test`/`lint`/`build` entry per workspace, fed into `runGate` via
`collectSteps`.
