import { resolvePath } from '../runtime-guard.ts'
import { dhPath } from '../dashboard-command.ts'
import { statusValue } from '../status-drift.ts'

// Pure helpers of the commit gate and the harness resolver, used by panel.ts.
export type Harness = { mode: 'repo' | 'global' | 'none'; dir: string }

// Commit/PR text that mentions AI: an empty `attribution` covers the default footer, this covers the rest.
export const AI_MENTION = /co-authored-by:\s*claude|generated with \[?claude|🤖|anthropic\.com|\bclaude code\b/i
export const GIT_C = String.raw`(?:-C\s+(?:"[^"]*"|'[^']*'|\S+)\s+)?`
export const SHIPS = new RegExp(String.raw`\b(git\s+${GIT_C}commit|gh\s+pr\s+(create|edit)|gh\s+release\s+create)\b`)
export const COMMIT = new RegExp(String.raw`\bgit\s+${GIT_C}commit(?![\w-])`)
export const GENERIC_GATE = /\b(make\s+(check|test|lint|ci|stan|gate)|go\s+(test|vet)|golangci-lint|pytest|ruff\s+check|cargo\s+(test|clippy)|(npm|pnpm|yarn|bun)\s+(run\s+)?(test|lint|check|typecheck)|tsc\b|phpstan|phpunit|pint)\b/
export const GATE_KEYS = /^ {2}(test|lint|vet|check|typecheck|stan|fmt[\w-]*):\s*\n\s+value:\s*(.+)$/gm

export function deniesAiMention(command: string): boolean {
  return SHIPS.test(command) && AI_MENTION.test(command)
}

export function parseGateCommands(projectYaml: string): string[] {
  return [...projectYaml.matchAll(GATE_KEYS)]
    .map(m => (m[2] ?? '').trim().replace(/^["']|["']$/g, ''))
    .filter(v => v && v !== 'null')
}

export function isGateCommand(command: string, projectGates: string[]): boolean {
  return GENERIC_GATE.test(command) || projectGates.some(g => command.includes(g))
}

export function parseLanguage(projectYaml: string): 'en' | 'pt-br' {
  return /^language:\s*(\S+)/m.exec(projectYaml)?.[1]?.replace(/^["']|["']$/g, '').toLowerCase() === 'pt-br' ? 'pt-br' : 'en'
}

// D5: dh bookkeeping under .harness/ (or, global mode, under <home>/projects/) is not code churn,
// so it does not dirty the gate. `home` is resolvePath-normalized; '' skips the global check.
export function isHarnessPath(p: unknown, home = ''): boolean {
  const n = String(p ?? '').replace(/\\/g, '/')
  const abs = /^(\/|[A-Za-z]:)/.test(n)
  const r = resolvePath('/', n) // `.harness/../src/a.go` is code, not bookkeeping
  if (r.includes('/.harness/')) return true
  return !!home && abs && r.startsWith(`${home === '/' ? '' : home}/projects/`)
}

const unquote = (s: string) => s.replace(/^(["'])(.*)\1$/, '$2')


// Back to the platform's form for git and fs calls: `/C:/x` -> `C:/x`; POSIX paths unchanged.
export const native = (p: string) => p.replace(/^\/([A-Z]):/, '$1:')

// Directory a git commit (or a gate) at offset `at` runs in: each `cd <dir>` before `at` moves it
// (relative to the previous one; a bare `cd` goes home), then `git -C <dir>` on the commit itself
// resolves against that. `home` expands a leading `~`.
// ponytail: shell parsing by regex; pushd, `cd -`, env-var dirs and `--git-dir` resolve wrong and fail open.
export function commandDir(command: string, cwd: string, home: string, at = command.length): string {
  let dir = cwd
  const expand = (d: string) => resolvePath(dir, unquote(d).replace(/^~(?=\/|$)/, home))
  for (const m of command.slice(0, at).matchAll(/(?:^|[;&|(\n]\s*|\bthen\s+)cd(?:\s+("[^"]*"|'[^']*'|[^\s;&|)]+))?(?=\s*(?:$|[;&|)\n]))/g)) {
    dir = m[1] === undefined ? (home || dir) : expand(m[1])
  }
  const c = /^git\s+-C\s+("[^"]*"|'[^']*'|\S+)\s+/.exec(command.slice(at))
  return c ? expand(c[1] ?? '.') : dir
}

// Offset of the first gate command in `command`, or -1.
export function gateAt(command: string, projectGates: string[]): number {
  const hits = [GENERIC_GATE.exec(command)?.index ?? -1, ...projectGates.map(g => command.indexOf(g))].filter(i => i >= 0)
  return hits.length ? Math.min(...hits) : -1
}

export const MSG = {
  en: {
    ai: 'dev-harness: commit/PR text mentions AI/Claude. Remove the mention and retry.',
    gate: (hint: string) => `dev-harness: files were edited after the last green gate in this session. Run ${hint} and commit again.`,
    gateDefault: "the project's quality gate (tests + lint)",
  },
  'pt-br': {
    ai: 'dev-harness: commit/PR menciona IA/Claude. Remova a menção e tente de novo.',
    gate: (hint: string) => `dev-harness: houve edição depois do último gate verde nesta sessão. Rode ${hint} e commite de novo.`,
    gateDefault: 'o quality gate do projeto (testes + lint)',
  },
}

export function taskStatus(taskMd: string): 'done' | 'blocked' | 'open' {
  const s = statusValue(taskMd)
  if (/^(conclu|pronta, entregue|done)/.test(s)) return 'done'
  if (/^bloque/.test(s)) return 'blocked'
  return 'open'
}

// The argv for `dh harness-path`, or undefined when spawning is unsafe: on Windows dhPath is a .cmd that
// cmd.exe parses, so a directory with a cmd metacharacter is never passed (treated as a dh failure).
export function harnessArgv(pluginRoot: string, dir: string): string[] | undefined {
  const dh = dhPath(pluginRoot)
  const d = native(dir)
  if (dh.endsWith('.cmd') && /[&|<>^%!"()\r\n]/.test(d)) return undefined
  return [dh, 'harness-path', '--json', d]
}

// What `dh harness-path --json` printed, or undefined when it is not a usable answer: the dir must be
// absolute, at most 4096 chars, and free of control/line-separator characters (it is shown to the model).
export function parseHarness(stdout: string): Harness | undefined {
  try {
    const o = JSON.parse(stdout)
    if (o?.mode === 'none') return { mode: 'none', dir: '' }
    const d = o?.dir
    if ((o?.mode === 'repo' || o?.mode === 'global') && typeof d === 'string' && d.length <= 4096
      && /^(\/|[A-Za-z]:|\\\\)/.test(d) && !/[\u0000-\u001f\u007f\u2028\u2029]/.test(d)) return { mode: o.mode, dir: d }
  } catch { /* fail open */ }
  return undefined
}

export const HARNESS_SECTION = 'dh:harness'

export function harnessText(h: Harness): string {
  return h.mode === 'none'
    ? 'Dev Harness: no harness for this project (mode none). Run /dh:setup to create one, or dh link if it already exists.'
    : `Dev Harness: mode ${h.mode}; harness dir ${JSON.stringify(h.dir)}. Records (MEMORY.md, EPOCHAL.md, RISKS.md, project.yaml, tasks/, prd/) live there.`
}

