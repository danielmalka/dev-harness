import type { EngineInterface, On } from 'claude-code'
import { homeFrom } from './runtime-guard.ts'

// PRD-012 R11/R13/R13b/R14/R15: the mod writes the per-session snapshot (schema 2) that
// `dh snapshot event` used to write, plus activity/rate_limits, and keeps updated_at alive.
const SAFE_ID = /^[\w-]{1,128}$/
const BEAT_MS = 30_000
const JOB_MS = 2_000

type Ev = Record<string, any>
type Snap = Record<string, any>
type Limits = { five_hour?: number; seven_day?: number }

// state is the schema 1 vocabulary; activity is the dashboard's (R13 table).
const MAP: Record<string, { activity: string; state: string }> = {
  SessionStart: { activity: 'idle', state: 'idle' },
  UserPromptSubmit: { activity: 'working', state: 'active' },
  PostToolUse: { activity: 'working', state: 'active' },
  SubagentStart: { activity: 'working', state: 'active' },
  SubagentStop: { activity: 'working', state: 'active' },
  PermissionRequest: { activity: 'waiting', state: 'active' },
  Notification: { activity: 'waiting', state: 'active' }, // permission_prompt only
  StopFailure: { activity: 'error', state: 'idle' },
  Stop: { activity: 'done', state: 'idle' },
  SessionEnd: { activity: 'idle', state: 'closed' },
}

export function iso(ms: number): string {
  return new Date(ms).toISOString().replace(/\.\d+Z$/, 'Z')
}

export function limitsFrom(rl: { kind: string; percentUsed: number }[]): Limits | undefined {
  const out: Limits = {}
  for (const r of rl) {
    if (r.kind === 'five_hour') out.five_hour = r.percentUsed
    if (r.kind === 'seven_day') out.seven_day = r.percentUsed
  }
  return Object.keys(out).length ? out : undefined
}

// Pure: previous file content + one classic hook event -> next content, or undefined when the event writes nothing.
// Unknown fields of `prev` (tasks, ...) ride along.
export function applyEvent(prev: Snap, name: string, e: Ev, now: string, limits?: Limits): Snap | undefined {
  const m = MAP[name]
  if (!m || (name === 'Notification' && e.notification_type !== 'permission_prompt')) return undefined
  // A subagent finishing says nothing about the main thread: SubagentStop keeps the current activity.
  const activity = name === 'SubagentStop' ? (prev.activity ?? m.activity) : m.activity
  const s: Snap = { ...prev, schema: 2, session_id: e.session_id, state: m.state, activity, updated_at: now }
  if (prev.activity !== activity || !prev.activity_at) s.activity_at = now
  if (name === 'SessionStart') {
    const title = e.session_name ?? e.session_title
    if (typeof title === 'string') s.session_name = title
    if (!s.started_at) s.started_at = now
    if (typeof e.cwd === 'string') s.cwd = e.cwd
    if (typeof e.agent_type === 'string') s.agent = e.agent_type
    if (typeof e.model === 'string') s.model = { ...(prev.model ?? {}), id: e.model }
  } else if (!s.cwd && typeof e.cwd === 'string') {
    s.cwd = e.cwd // mod loaded mid-session: still attributable to a project
  }
  if (name === 'SubagentStart' || name === 'SubagentStop') {
    const rec: Ev = { at: now, event: name }
    if (typeof e.agent_type === 'string') rec.agent_type = e.agent_type
    if (typeof e.agent_id === 'string') rec.agent_id = e.agent_id
    s.events = [...(Array.isArray(prev.events) ? prev.events : []), rec].slice(-50)
  }
  if (limits) s.rate_limits = limits
  return s
}

// ponytail: module-level state; a mod reload resets it (current session only).
let chain: Promise<unknown> = Promise.resolve()
let beatSession = ''
let beat: { cancel: () => void } | undefined
const SKIP_LIMIT = 3
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
  try {
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
    const next = applyEvent(prev, name, e, iso(await $.clock.now()), limits)
    if (!next) return
    await $.fs.write(path, JSON.stringify(next, null, 2) + '\n')
    if (name !== 'SessionEnd') {
      beatSession = e.session_id
      beat ??= $.clock.every(BEAT_MS, () => { chain = chain.then(() => withTimeout($, touch($))).catch(() => undefined) })
    }
  } finally {
    if (name === 'SessionEnd' && e.session_id === beatSession) {
      beat?.cancel()
      beat = undefined
      beatSession = ''
    }
  }
}

// Events are serialized so two quick hooks never read-modify-write the same file out of order.
// A failed write is silent: a snapshot is a convenience, never a reason to disturb the session.
async function record($: EngineInterface, name: string, e: Ev): Promise<void> {
  const job = chain.then(() => withTimeout($, write($, name, e)))
  chain = job.catch(() => undefined)
  await job.catch(() => undefined)
}

// ponytail: a mod failure never disturbs the session (fail-open).
const failOpen = (_$: unknown, e: any, next: any) => next(e)

export function registerSnapshot(on: On): void {
  on('classic.SessionStart', async ($, e, next) => { await record($, 'SessionStart', e); return next(e) }).catch(failOpen)
  on('classic.UserPromptSubmit', async ($, e, next) => { void record($, 'UserPromptSubmit', e); return next(e) }).catch(failOpen)
  on('classic.PostToolUse', async ($, e, next) => { void record($, 'PostToolUse', e); return next(e) }).catch(failOpen)
  // Hot-path events (prompt, tool, subagent) do not wait for the write; chain keeps the order.
  // PermissionRequest/Notification record before next(): next waits for the person's answer.
  on('classic.PermissionRequest', async ($, e, next) => { await record($, 'PermissionRequest', e); return next(e) }).catch(failOpen)
  on('classic.Notification', async ($, e, next) => { await record($, 'Notification', e); return next(e) }).catch(failOpen)
  on('classic.StopFailure', async ($, e, next) => { await record($, 'StopFailure', e); return next(e) }).catch(failOpen)
  on('classic.Stop', async ($, e, next) => { await record($, 'Stop', e); return next(e) }).catch(failOpen)
  on('classic.SubagentStart', async ($, e, next) => { void record($, 'SubagentStart', e); return next(e) }).catch(failOpen)
  on('classic.SubagentStop', async ($, e, next) => { void record($, 'SubagentStop', e); return next(e) }).catch(failOpen)
  on('classic.SessionEnd', async ($, e, next) => { await record($, 'SessionEnd', e); return next(e) }).catch(failOpen)
}
