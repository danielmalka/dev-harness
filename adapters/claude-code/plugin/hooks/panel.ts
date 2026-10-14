import type { EngineInterface, Register } from 'claude-code'
import { denySpawn, guardFile, homeFrom, resolvePath } from './runtime-guard.ts'
import { DASHBOARD_COMMAND, registerDashboard } from './dashboard-command.ts'
import { applyEvent, iso, limitsFrom, type Ev, type Limits, type Snap } from './snapshot-writer.ts'
import { driftMessage, findDrift, statusValue } from './status-drift.ts'
import {
  COMMIT, commandDir, deniesAiMention, gateAt, type Harness, harnessArgv, HARNESS_SECTION, harnessText, isGateCommand,
  isHarnessPath, MSG, native, parseGateCommands, parseHarness, parseLanguage, taskStatus,
} from './lib/gate.ts'
export { commandDir, deniesAiMention, harnessArgv, HARNESS_SECTION, harnessText, isGateCommand, isHarnessPath, native, parseGateCommands, parseHarness, parseLanguage, resolvePath, taskStatus }

// ponytail: module-level state; a mod reload resets the gate (current session only).
// Absolute paths edited since the last green gate of their repository.
const dirty = new Set<string>()
const inside = (f: string, root: string) => f === root || f.startsWith(root === '/' ? '/' : `${root}/`)

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

// <home> by the ADR-007 rule, plus the user home for `~`; literal names so `claude plugin validate` lists them.
async function homes($: EngineInterface): Promise<{ home: string; user: string; cwd: string }> {
  try {
    const cwd = await $.session.cwd().catch(() => '/')
    const dh = await $.env.get('DH_HOME'), hm = await $.env.get('HOME'), up = await $.env.get('USERPROFILE')
    return { home: homeFrom(dh, hm, up, cwd), user: hm || up || '', cwd }
  } catch {
    return { home: '', user: '', cwd: '/' }
  }
}

// R8: <home> is read here and passed in; guardFile cannot take $ across the import.
async function guardWrite($: EngineInterface, e: any, next: any): Promise<any> {
  const { home, user, cwd } = await homes($)
  return guardFile(home, e, async (x: typeof e) => markDirty($, await next(x), x), cwd, user)
}

// R20: `dh harness-path --json <dir>`, the only resolver (R2: the mod never reads the global registry file).
// undefined = dh failed (missing binary, exit != 0, bad JSON, timeout).
const harnesses = new Map<string, Harness | undefined>()

async function resolveHarness($: EngineInterface, dir: string): Promise<Harness | undefined> {
  try {
    const argv = harnessArgv($.plugin.root, dir)
    if (!argv) return undefined
    const r = await $.process.run(argv, { timeoutMs: 2_000 })
    return r.exitCode === 0 ? parseHarness(r.stdout) : undefined
  } catch {
    return undefined // fail open
  }
}

// Cached per directory. repo/global stick for the session; none and failures are dropped by every
// prompt.compose (so /dh:setup or dh link apply without a restart) and re-resolved on the next read.
// ponytail: readers between two composes reuse a none/failed answer, so a gate never spawns dh per call.
// Resolved from the repository root (as the commit gate does), so a session opened in a subfolder
// finds its repo's harness; a non-git dir resolves as itself. Roots are cached per directory.
// ponytail: a `git init` mid-session is not seen for a dir already cached as non-git.
const roots = new Map<string, string>()
async function rootFor($: EngineInterface, dir: string): Promise<string> {
  const key = resolvePath('/', dir)
  let root = roots.get(key)
  if (root === undefined) roots.set(key, root = await repoRoot($, key))
  return root
}

async function harnessOf($: EngineInterface, dir: string): Promise<{ root: string; h: Harness | undefined }> {
  const root = await rootFor($, dir)
  if (!harnesses.has(root)) harnesses.set(root, await resolveHarness($, root))
  return { root, h: harnesses.get(root) }
}

// The folder holding project.yaml, tasks/ and prd/ for `dir`: the resolved one, else `<dir>/.harness` when it exists, else `<root>/.harness`.
async function harnessDir($: EngineInterface, dir: string): Promise<string> {
  try {
    const { root, h } = await harnessOf($, dir)
    if (h && h.mode !== 'none') return h.dir
    // 0.20.0 read `<cwd>/.harness`; keep it when that folder exists below the git toplevel.
    const own = `${dir}/.harness`
    return (await $.fs.exists(own).catch(() => false)) ? own : `${native(root)}/.harness`
  } catch {
    return `${dir}/.harness`
  }
}

async function composeSection($: EngineInterface): Promise<{ id: string; text: string; scope: 'session' } | undefined> {
  for (const [k, v] of harnesses) if (!v || v.mode === 'none') harnesses.delete(k)
  const { h } = await harnessOf($, await $.session.cwd())
  return h ? { id: HARNESS_SECTION, text: harnessText(h), scope: 'session' } : undefined
}

async function readProject($: EngineInterface, root?: string): Promise<{ gates: string[]; lang: 'en' | 'pt-br' }> {
  try {
    const y = await $.fs.read(`${await harnessDir($, root ?? await $.session.cwd())}/project.yaml`) as string
    return { gates: parseGateCommands(y), lang: parseLanguage(y) }
  } catch {
    return { gates: [], lang: 'en' }
  }
}

// R20: fail-open on any read error; no tickets (or no readable PRDs) means no drift.
async function driftNow($: EngineInterface, cwd: string): Promise<ReturnType<typeof findDrift>> {
  try {
    const hdir = await harnessDir($, cwd)
    const tickets: { id: string; md: string }[] = []
    for (const entry of await $.fs.list(`${hdir}/tasks`)) {
      const file = `${hdir}/tasks/${entry.name}/TASK.md`
      if (entry.kind === 'dir' && (await $.fs.exists(file))) tickets.push({ id: entry.name, md: await $.fs.read(file) as string })
    }
    if (!tickets.length) return []
    const prds = new Map<string, string>()
    for (const dir of [`${cwd}/docs/prd`, `${hdir}/prd`]) {
      try {
        for (const f of (await $.fs.list(dir)).sort((a, b) => (a.name < b.name ? -1 : 1))) {
          const id = /^(PRD-\d+)/.exec(f.name)?.[1]
          if (!id || f.kind !== 'file' || !f.name.endsWith('.md') || f.name.endsWith('.review.md') || prds.has(id)) continue
          prds.set(id, statusValue(await $.fs.read(`${dir}/${f.name}`) as string))
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
  if (r.deny === undefined && r.isError !== true && p && !isHarnessPath(p, (await homes($)).home)) {
    dirty.add(resolvePath(await $.session.cwd().catch(() => '/'), String(p)))
  }
  return r
}

// BUG-001 snapshot writer (engine events, no classic.*): per-session file, serialized writes.
const SAFE_ID = /^[\w-]{1,128}$/
const BEAT_MS = 30_000
const JOB_MS = 2_000
const END_MS = 1_000 // session.end's whole budget is 1.5s: wait at most this long for the close write

// ponytail: module-level state; a mod reload resets it (current session only).
let chain: Promise<unknown> = Promise.resolve()
let beatSession = ''
let beat: { cancel: () => void } | undefined
const SKIP_LIMIT = 3
const closedHere = new Set<string>() // ids this module instance closed via session.end: a late record must not reopen them
const skips = new Map<string, number>() // consecutive skipped writes per session (unparseable file)

// R10: <home>/sessions, <home> by the ADR-007 rule (DH_HOME is the only control); '' when no home is known.
// $.fs.write creates missing parent directories (engine type doc), so a clean machine needs no `dh` call first.
async function snapshotDir($: EngineInterface): Promise<string> {
  const home = homeFrom(await $.env.get('DH_HOME'), await $.env.get('HOME'), await $.env.get('USERPROFILE'), await $.session.cwd().catch(() => '/'))
  return home ? `${home.replace(/^\/([A-Z]):/, '$1:')}/sessions` : ''
}

// {} when the file does not exist; undefined when it exists but cannot be read or parsed (the caller skips the write).
async function readPrev($: EngineInterface, path: string): Promise<Snap | undefined> {
  try {
    if (!(await $.fs.exists(path))) return {}
    const o = JSON.parse(await $.fs.read(path) as string)
    return o && typeof o === 'object' && !Array.isArray(o) ? o : undefined
  } catch {
    return undefined
  }
}

// A hung fs/usage call must never hold the session: each serialized job gets 2s.
function withTimeout<T>($: EngineInterface, p: Promise<T>): Promise<T> {
  let t: { cancel: () => void } | undefined
  const limit = new Promise<never>((_, reject) => { t = $.clock.after(JOB_MS, () => reject(new Error('snapshot job timed out'))) })
  return Promise.race([p, limit]).finally(() => t?.cancel())
}

async function usageLimits($: EngineInterface): Promise<Limits | undefined> {
  try { return limitsFrom((await $.session.usage()).rateLimits) } catch { return undefined } // no usage: keep previous
}

async function touch($: EngineInterface): Promise<void> {
  const id = beatSession
  const dir = await snapshotDir($)
  if (!id || !dir) return
  const path = `${dir}/${id}.json`
  const prev = await readPrev($, path)
  const limits = await usageLimits($) // a limit's age must reflect when it was measured
  const now = iso(await $.clock.now())
  // Re-checked after the last await, inside the serialized job: an ended session is never resurrected.
  if (!prev || beatSession !== id || !prev.session_id || prev.state === 'closed') return
  await $.fs.write(path, JSON.stringify({ ...prev, updated_at: now, ...(limits ? { rate_limits: limits } : {}) }, null, 2) + '\n')
}

async function write($: EngineInterface, name: string, e: Ev): Promise<void> {
  if (!SAFE_ID.test(String(e.session_id ?? ''))) return
  {
    const dir = await snapshotDir($)
    if (!dir) return
    const path = `${dir}/${e.session_id}.json`
    const limits = await usageLimits($)
    let prev = await readPrev($, path)
    if (!prev) {
      // exists but unreadable/partial: do not overwrite what another writer is mid-way through,
      // unless it stays unparseable for SKIP_LIMIT writes in a row (then it is corrupt, not partial)
      const n = (skips.get(e.session_id) ?? 0) + 1
      if (n < SKIP_LIMIT) { skips.set(e.session_id, n); return }
      prev = {}
    }
    skips.delete(e.session_id)
    // A late record of a session this instance ended never reopens it. A SessionStart or a prompt (always
    // user-initiated, e.g. an in-process /resume of the same id) lifts the mark; the file's own state is no signal.
    if (name === 'SessionStart' || name === 'UserPromptSubmit') closedHere.delete(e.session_id)
    else if (name !== 'SessionEnd' && closedHere.has(e.session_id)) return
    const now = iso(await $.clock.now())
    const next = applyEvent(prev, name, e, now, limits)
    if (!next) return
    // After /clear or a resume the process continues under a new id and no session.start fires:
    // fill what only that event would have written, so the dashboard can attribute the file to a project.
    if (name !== 'SessionEnd') {
      if (!next.cwd) { const c = await $.session.cwd().catch(() => undefined); if (typeof c === 'string' && c) next.cwd = c }
      if (!next.model?.id) { const m = await $.session.model().catch(() => undefined); if (typeof m === 'string' && m) next.model = { ...(next.model ?? {}), id: m } }
      next.started_at ||= now
    }
    await $.fs.write(path, JSON.stringify(next, null, 2) + '\n')
    if (name !== 'SessionEnd') {
      beatSession = e.session_id
      beat ??= $.clock.every(BEAT_MS, () => { chain = chain.then(() => withTimeout($, touch($))).catch(() => undefined) })
    }
  }
}

// Events are serialized so two quick hooks never read-modify-write the same file out of order.
// A failed write is silent: a snapshot is a convenience, never a reason to disturb the session.
// The session id is read when the event happens, before queueing, so a late record never lands on the next id.
async function record($: EngineInterface, name: string, e0: Ev): Promise<void> {
  const e: Ev = { ...e0, session_id: e0.session_id ?? await $.session.id().catch(() => undefined) }
  const job = chain.then(() => withTimeout($, write($, name, e)))
  chain = job.catch(() => undefined)
  await job.catch(() => undefined)
}

// What panel.ts calls from its engine hooks. Each is fail-open and returns a promise that never rejects.
// Hot-path events (prompt, tool, subagent) do not wait for the write; chain keeps the order.
async function recordStart($: EngineInterface, e: { cwd?: string }): Promise<void> {
  const model = await $.session.model().catch(() => undefined) // agent type: the engine does not expose it
  await record($, 'SessionStart', { cwd: e.cwd, model })
}
function recordPrompt($: EngineInterface): void { void record($, 'UserPromptSubmit', {}) }
function recordTool($: EngineInterface): void { void record($, 'PostToolUse', {}) }
function recordSpawn($: EngineInterface, subagentType: string, agentId: string): void {
  void record($, 'SubagentStart', { agent_type: subagentType, agent_id: agentId })
}
// A subagent's turn.complete is SubagentStop (keeps activity, does not wait); the main loop's is Stop/StopFailure.
function recordTurn($: EngineInterface, e: { agentId?: string; reason?: string }): Promise<void> {
  if (e.agentId) { void record($, 'SubagentStop', { agent_id: e.agentId }); return Promise.resolve() }
  return record($, e.reason === 'error' ? 'StopFailure' : 'Stop', {})
}
// The close write stays last in the chain (ordering kept), but session.end has 1.5s in all: when the chain
// is slow (a hung fs, a long queue) stop waiting after END_MS and let the job finish on its own, instead of
// being cut by the engine and replayed through failOpen.
async function recordEnd($: EngineInterface, e: { sessionId?: string }): Promise<void> {
  // Synchronously, before anything is queued: a hung fs call in the close job cannot leave the beat running,
  // and records still queued for this id are dropped as late (the close write follows them).
  if (e.sessionId) {
    closedHere.add(e.sessionId)
    if (e.sessionId === beatSession) { beat?.cancel(); beat = undefined; beatSession = '' }
  }
  let t: { cancel: () => void } | undefined
  const cap = new Promise<void>(resolve => { t = $.clock.after(END_MS, () => resolve()) })
  await Promise.race([record($, 'SessionEnd', { session_id: e.sessionId }), cap])
  t?.cancel()
}

// ponytail: a mod failure never blocks a tool (fail-open); the trap is a convenience, not security.
const failOpen = ($: EngineInterface, e: any, next: any) => next(e)

// hooks.json `modules` accepts one entry per plugin and a duplicate on(event) is refused,
// so panel.ts is the entry and chains the runtime-guard checks (ADR-005) ahead of its own.
// The loader refuses $ passed across an import: denySpawn gets `{}`, guardFile gets <home> as a string.
export const register: Register = on => {
  registerDashboard(on)

  on('session.start', async ($, e, next) => {
    const r = await next(e)
    await $.command.register(DASHBOARD_COMMAND).catch(() => undefined)
    await recordStart($, e)
    return r
  })

  on('turn.complete', async ($, e, next) => {
    const r = await next(e)
    await recordTurn($, e)
    return r
  })

  on('agent.spawn', ($, e, next) => denySpawn({}, e, async (x: typeof e) => {
    const r = await next(x)
    if (r.agentId) recordSpawn($, x.subagentType, r.agentId)
    return r
  })).catch(failOpen)

  on('prompt.submit', ($, e, next) => { recordPrompt($); return next(e) }).catch(failOpen)

  // BUG-001: every other tool call is activity too; recorded after the guards beneath have answered.
  on('tool.call', async ($, e, next) => {
    const r = await next(e)
    recordTool($)
    return r
  }).catch(failOpen)

  on('session.end', async ($, e, next) => {
    await recordEnd($, e)
    return next(e)
  }).catch(failOpen)

  on('tool.call', { tool: 'Edit' }, guardWrite).catch(failOpen)
  on('tool.call', { tool: 'Write' }, guardWrite).catch(failOpen)
  on('tool.call', { tool: 'NotebookEdit' }, async ($, e, next) => markDirty($, await next(e), e)).catch(failOpen)

  // R20: tell the session where its harness lives; a dh failure adds nothing (fail-open).
  on('prompt.compose', async ($, e, next) => {
    const r = await next(e)
    const section = await composeSection($).catch(() => undefined)
    return section ? { sections: [...r.sections, section] } : r
  }).catch(failOpen)

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
