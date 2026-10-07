import type { EngineInterface, Register } from 'claude-code'
import { denySpawn, guardFile } from './runtime-guard.ts'
import { DASHBOARD_COMMAND, registerDashboard } from './dashboard-command.ts'
import { registerSnapshot } from './snapshot-writer.ts'
import { driftMessage, findDrift, statusValue } from './status-drift.ts'

// Commit/PR text that mentions AI: an empty `attribution` covers the default footer, this covers the rest.
const AI_MENTION = /co-authored-by:\s*claude|generated with \[?claude|🤖|anthropic\.com|\bclaude code\b/i
const GIT_C = String.raw`(?:-C\s+(?:"[^"]*"|'[^']*'|\S+)\s+)?`
const SHIPS = new RegExp(String.raw`\b(git\s+${GIT_C}commit|gh\s+pr\s+(create|edit)|gh\s+release\s+create)\b`)
const COMMIT = new RegExp(String.raw`\bgit\s+${GIT_C}commit(?![\w-])`)
const GENERIC_GATE = /\b(make\s+(check|test|lint|ci|stan|gate)|go\s+(test|vet)|golangci-lint|pytest|ruff\s+check|cargo\s+(test|clippy)|(npm|pnpm|yarn|bun)\s+(run\s+)?(test|lint|check|typecheck)|tsc\b|phpstan|phpunit|pint)\b/
const GATE_KEYS = /^ {2}(test|lint|vet|check|typecheck|stan|fmt[\w-]*):\s*\n\s+value:\s*(.+)$/gm

export function bar(pct: number, width = 10): string {
  const full = Math.round((Math.min(Math.max(pct, 0), 100) / 100) * width)
  return '▰'.repeat(full) + '▱'.repeat(width - full)
}

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

// D5: dh bookkeeping under .harness/ is not code churn, so it does not dirty the gate.
export function isHarnessPath(p: unknown): boolean {
  const n = String(p ?? '').replace(/\\/g, '/')
  return n.startsWith('.harness/') || n.includes('/.harness/')
}

const unquote = (s: string) => s.replace(/^(["'])(.*)\1$/, '$2')

// Normalizes `p` against `base`: absolute result, no `.`/`..`/empty segments. A Windows drive path
// (`C:\\x`, `c:/x`) is absolute too and reads as `/C:/x`, so git's `C:/repo` and the engine's `C:\\repo` meet.
export function resolvePath(base: string, p: string): string {
  const abs = (x: string) => x.replace(/\\/g, '/').replace(/^\/?([A-Za-z]):(?=\/|$)/, (_m, d: string) => `/${d.toUpperCase()}:`)
  const n = abs(p)
  const out: string[] = []
  for (const seg of (n.startsWith('/') ? n : `${abs(base)}/${n}`).split('/')) {
    if (seg === '..') out.pop()
    else if (seg && seg !== '.') out.push(seg)
  }
  return `/${out.join('/')}`
}

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
function gateAt(command: string, projectGates: string[]): number {
  const hits = [GENERIC_GATE.exec(command)?.index ?? -1, ...projectGates.map(g => command.indexOf(g))].filter(i => i >= 0)
  return hits.length ? Math.min(...hits) : -1
}

const MSG = {
  en: {
    ai: 'dev-harness: commit/PR text mentions AI/Claude. Remove the mention and retry.',
    gate: (hint: string) => `dev-harness: files were edited after the last green gate in this session. Run ${hint} and commit again.`,
    gateDefault: "the project's quality gate (tests + lint)",
    ctx: (p: number, hard: boolean) => `Context at ${p}% - ${hard ? 'time to /compact or /dh:handoff' : 'consider /compact at the next pause'}`,
  },
  'pt-br': {
    ai: 'dev-harness: commit/PR menciona IA/Claude. Remova a menção e tente de novo.',
    gate: (hint: string) => `dev-harness: houve edição depois do último gate verde nesta sessão. Rode ${hint} e commite de novo.`,
    gateDefault: 'o quality gate do projeto (testes + lint)',
    ctx: (p: number, hard: boolean) => `Contexto em ${p}% - ${hard ? 'hora de /compact ou /dh:handoff' : 'considere /compact no próximo ponto de pausa'}`,
  },
}

export function taskStatus(taskMd: string): 'done' | 'blocked' | 'open' {
  const s = statusValue(taskMd)
  if (/^(conclu|pronta, entregue|done)/.test(s)) return 'done'
  if (/^bloque/.test(s)) return 'blocked'
  return 'open'
}

// ponytail: module-level state; a mod reload resets gate/agents (current session only).
// Absolute paths edited since the last green gate of their repository.
const dirty = new Set<string>()
const inside = (f: string, root: string) => f === root || f.startsWith(root === '/' ? '/' : `${root}/`)
const running = new Map<string, string>()
const warned = new Set<number>()

async function home($: EngineInterface): Promise<string> {
  try { return (await $.env.get('HOME')) ?? '' } catch { return '' }
}

// Repository root of `dir` from git; `dir` itself when git does not answer.
async function repoRoot($: EngineInterface, dir: string): Promise<string> {
  try {
    const r = await $.process.run(['git', '-C', native(dir), 'rev-parse', '--show-toplevel'], { timeoutMs: 2_000 })
    const top = r.exitCode === 0 ? r.stdout.trim() : ''
    return /^(\/|[A-Za-z]:)/.test(top) ? resolvePath('/', top) : dir
  } catch {
    return dir
  }
}

async function rootOf($: EngineInterface, command: string, at?: number): Promise<string> {
  return repoRoot($, commandDir(command, await $.session.cwd(), await home($), at))
}

async function readProject($: EngineInterface, root?: string): Promise<{ gates: string[]; lang: 'en' | 'pt-br' }> {
  try {
    const y = await $.fs.read(`${root ?? await $.session.cwd()}/.harness/project.yaml`) as string
    return { gates: parseGateCommands(y), lang: parseLanguage(y) }
  } catch {
    return { gates: [], lang: 'en' }
  }
}

async function dhProgress($: EngineInterface): Promise<string | undefined> {
  const dir = `${await $.session.cwd()}/.harness/tasks`
  try {
    let done = 0, blocked = 0, total = 0
    for (const entry of await $.fs.list(dir)) {
      if (entry.kind !== 'dir') continue
      const file = `${dir}/${entry.name}/TASK.md`
      if (!(await $.fs.exists(file))) continue
      const st = taskStatus(await $.fs.read(file) as string)
      total++
      if (st === 'done') done++
      if (st === 'blocked') blocked++
    }
    if (!total) return undefined
    return `dh ${bar((done / total) * 100, 5)} ${done}/${total}${blocked ? ` (${blocked} blocked)` : ''}`
  } catch {
    return undefined
  }
}

async function refresh($: EngineInterface): Promise<void> {
  const u = await $.session.usage()
  const pct = u.context.percent ?? 0
  const parts = [`ctx ${bar(pct)} ${pct}%`]
  if (running.size) parts.push(`agents ${running.size}`)
  const dh = await dhProgress($)
  if (dh) parts.push(dh)
  const cwd = await $.session.cwd()
  if ([...dirty].some(f => inside(f, resolvePath('/', cwd)))) parts.push('gate pending')
  for (const r of u.rateLimits) {
    const label = r.kind === 'five_hour' ? '5h' : r.kind === 'seven_day' ? 'wk' : undefined
    if (label) parts.push(`${label} ${Math.round(r.percentUsed)}%${r.percentUsed >= 80 ? ' 🔴' : ''}`)
  }
  if (u.cost) parts.push(`$${u.cost.usd.toFixed(2)}`)
  $.ui.status(parts.join(' · '))
  for (const t of [60, 80]) {
    if (pct >= t && !warned.has(t)) {
      warned.add(t)
      $.ui.toast(MSG[(await readProject($)).lang].ctx(pct, t >= 80))
    }
  }
}

// R20: fail-open on any read error; no tickets (or no readable PRDs) means no drift.
async function driftNow($: EngineInterface, cwd: string): Promise<ReturnType<typeof findDrift>> {
  try {
    const tickets: { id: string; md: string }[] = []
    for (const entry of await $.fs.list(`${cwd}/.harness/tasks`)) {
      const file = `${cwd}/.harness/tasks/${entry.name}/TASK.md`
      if (entry.kind === 'dir' && (await $.fs.exists(file))) tickets.push({ id: entry.name, md: await $.fs.read(file) as string })
    }
    if (!tickets.length) return []
    const prds = new Map<string, string>()
    for (const dir of ['docs/prd', '.harness/prd']) {
      try {
        for (const f of (await $.fs.list(`${cwd}/${dir}`)).sort((a, b) => (a.name < b.name ? -1 : 1))) {
          const id = /^(PRD-\d+)/.exec(f.name)?.[1]
          if (!id || f.kind !== 'file' || !f.name.endsWith('.md') || f.name.endsWith('.review.md') || prds.has(id)) continue
          prds.set(id, statusValue(await $.fs.read(`${cwd}/${dir}/${f.name}`) as string))
        }
      } catch { /* directory absent */ }
    }
    return findDrift(tickets, prds, taskStatus)
  } catch {
    return []
  }
}

async function guardBash($: EngineInterface, command: string): Promise<string | undefined> {
  const aiHit = deniesAiMention(command)
  const commit = COMMIT.exec(command)
  if (!aiHit && !commit) return undefined
  // The commit's repository, not the session's, decides the gate, the language and the drift check.
  const root = commit ? await rootOf($, command, commit.index) : await $.session.cwd()
  const { gates, lang } = await readProject($, native(root))
  if (aiHit) return MSG[lang].ai
  if ([...dirty].some(f => inside(f, root))) return MSG[lang].gate(gates.length ? gates.join(' && ') : MSG[lang].gateDefault)
  const drift = await driftNow($, native(root))
  return drift.length ? driftMessage(drift, lang, native(root)) : undefined
}

async function markDirty<R extends { deny?: unknown; isError?: boolean }>($: EngineInterface, r: R, e: any): Promise<R> {
  const p = e.file_path ?? e.notebook_path
  if (r.deny === undefined && r.isError !== true && p && !isHarnessPath(p)) {
    dirty.add(resolvePath(await $.session.cwd().catch(() => '/'), String(p)))
  }
  return r
}

// ponytail: a mod failure never blocks a tool (fail-open); the trap is a convenience, not security.
const failOpen = ($: EngineInterface, e: any, next: any) => next(e)

// hooks.json `modules` accepts one entry per plugin and a duplicate on(event) is refused,
// so panel.ts is the entry and chains the runtime-guard checks (ADR-005) ahead of its own.
// They ignore $, and the loader refuses $ passed across an import, hence `{}`.
export const register: Register = on => {
  registerSnapshot(on)
  registerDashboard(on)

  on('session.start', async ($, e, next) => {
    const r = await next(e)
    await $.command.register(DASHBOARD_COMMAND).catch(() => undefined)
    await refresh($).catch(() => undefined)
    return r
  })

  on('turn.complete', async ($, e, next) => {
    const r = await next(e)
    if (e.agentId) running.delete(e.agentId)
    await refresh($).catch(() => undefined)
    return r
  })

  on('agent.spawn', ($, e, next) => denySpawn({}, e, async (x: typeof e) => {
    const r = await next(x)
    if (r.agentId) {
      running.set(r.agentId, x.description)
      await refresh($).catch(() => undefined)
    }
    return r
  })).catch(failOpen)

  on('tool.call', { tool: 'Edit' }, ($, e, next) => guardFile({}, e, async (x: typeof e) => markDirty($, await next(x), x))).catch(failOpen)
  on('tool.call', { tool: 'Write' }, ($, e, next) => guardFile({}, e, async (x: typeof e) => markDirty($, await next(x), x))).catch(failOpen)
  on('tool.call', { tool: 'NotebookEdit' }, async ($, e, next) => markDirty($, await next(e), e)).catch(failOpen)

  on('tool.call', { tool: 'Bash' }, async ($, e, next) => {
    const denied = await guardBash($, e.command).catch(() => undefined)
    if (denied) return { deny: denied }
    const r = await next(e)
    // ponytail: pre-filter by generic + session gates (no git spawn on plain Bash); a custom gate that only
    // another repo's project.yaml names does not clear that repo.
    const at = r.deny === undefined && r.isError !== true && dirty.size ? gateAt(e.command, (await readProject($)).gates) : -1
    if (at >= 0) {
      const root = await rootOf($, e.command, at)
      for (const f of dirty) if (inside(f, root)) dirty.delete(f)
    }
    return r
  }).catch(failOpen)
}
