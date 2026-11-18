# Diff em revisão

Mudança fictícia, usada apenas como fixture. Escopo do caso de revisão `code`
`code-reviewer-equivalent-form-miss`.

- Escopo: `scripts/ci-gate.js`, `scripts/workspaces.config.js`
- Comparação: `main..feat/ci-gate-classify`

```diff
--- /dev/null
+++ b/scripts/workspaces.config.js
@@ -0,0 +1,21 @@
+// Comandos canônicos do pipeline de CI por workspace. São as invocações
+// exatas que o `npm run ci` usa hoje (veja o script `"ci"` do
+// `package.json` da raiz e o `package.json` de cada workspace), não uma
+// forma normalizada ou idealizada.
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
+// Remove um prefixo `cd <dir> &&` inicial para que um classificador case o
+// comando puro, qualquer que seja o diretório de workspace em que o passo roda.
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
+    // Só o código de saída de um passo 'generic' barra o merge; a saída de um passo
+    // 'test' também é analisada atrás do resumo de passou/falhou do próprio framework,
+    // então uma falha real escondida em saídas com código 0 ainda barra o merge.
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

`classify` despacha cada passo do pipeline para `classifyStep`, `classifyLintStep`
ou `classifyBuildStep` por correspondência de palavra-chave na string do comando, e
cada classificador compara o comando (depois que `normalizeCommand` remove um
prefixo `cd <dir> &&` inicial) com a invocação npm pura do seu tipo de passo. Os
comandos de fato rodados por workspace ficam em `scripts/workspaces.config.js`,
uma entrada `test`/`lint`/`build` por workspace, repassada a `runGate` por
`collectSteps`.
