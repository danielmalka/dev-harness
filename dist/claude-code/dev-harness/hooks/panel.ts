import type { EngineInterface, Register } from 'claude-code'
import { denySpawn, guardFile } from './runtime-guard.ts'
import { DASHBOARD_COMMAND, registerDashboard } from './dashboard-command.ts'
import { registerSnapshot } from './snapshot-writer.ts'
import { driftMessage, findDrift, statusValue } from './status-drift.ts'

// Commit/PR text that mentions AI: an empty `attribution` covers the default footer, this covers the rest.
const AI_MENTION = /co-authored-by:\s*claude|generated with \[?claude|🤖|anthropic\.com|\bclaude code\b/i
const SHIPS = /\b(git\s+commit|gh\s+pr\s+(create|edit)|gh\s+release\s+create)\b/
const COMMIT = /\bgit\s+commit(?![\w-])/
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
let dirtySinceGate = false
const running = new Map<string, string>()
const warned = new Set<number>()

async function readProject($: EngineInterface): Promise<{ gates: string[]; lang: 'en' | 'pt-br' }> {
  try {
    const y = await $.fs.read(`${await $.session.cwd()}/.harness/project.yaml`) as string
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
  if (dirtySinceGate) parts.push('gate pending')
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
async function driftNow($: EngineInterface): Promise<ReturnType<typeof findDrift>> {
  try {
    const cwd = await $.session.cwd()
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
  const commit = COMMIT.test(command)
  if (!aiHit && !commit) return undefined
  const { gates, lang } = await readProject($)
  if (aiHit) return MSG[lang].ai
  if (dirtySinceGate) return MSG[lang].gate(gates.length ? gates.join(' && ') : MSG[lang].gateDefault)
  const drift = await driftNow($)
  return drift.length ? driftMessage(drift, lang, await $.session.cwd()) : undefined
}

function markDirty<R extends { deny?: unknown; isError?: boolean }>(r: R, e: any): R {
  if (r.deny === undefined && r.isError !== true && !isHarnessPath(e.file_path ?? e.notebook_path)) dirtySinceGate = true
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

  on('tool.call', { tool: 'Edit' }, ($, e, next) => guardFile({}, e, async (x: typeof e) => markDirty(await next(x), x))).catch(failOpen)
  on('tool.call', { tool: 'Write' }, ($, e, next) => guardFile({}, e, async (x: typeof e) => markDirty(await next(x), x))).catch(failOpen)
  on('tool.call', { tool: 'NotebookEdit' }, async ($, e, next) => markDirty(await next(e), e)).catch(failOpen)

  on('tool.call', { tool: 'Bash' }, async ($, e, next) => {
    const denied = await guardBash($, e.command).catch(() => undefined)
    if (denied) return { deny: denied }
    const r = await next(e)
    if (r.deny === undefined && r.isError !== true && isGateCommand(e.command, (await readProject($)).gates)) {
      dirtySinceGate = false
    }
    return r
  }).catch(failOpen)
}
