// Pure part of the snapshot writer; the engine-facing part lives in panel.ts, because the loader never
// follows `$` across an import (BUG-001). PRD-012 R11/R13/R13b/R14/R15: the mod writes the per-session snapshot (schema 2) that
// `dh snapshot event` used to write, plus activity/rate_limits, and keeps updated_at alive.
export type Ev = Record<string, any>
export type Snap = Record<string, any>
export type Limits = { five_hour?: number; seven_day?: number }

// Keys are the snapshot's own event kinds (the names the dashboard and `dh snapshot event` know), not engine
// events: BUG-001 maps session.start/prompt.submit/tool.call/agent.spawn/turn.complete/session.end onto them,
// because the engine floors every `classic.*` hook on Team accounts. There is no un-floored "waiting" source.
// state is the schema 1 vocabulary; activity is the dashboard's (R13 table).
const MAP: Record<string, { activity: string; state: string }> = {
  SessionStart: { activity: 'idle', state: 'idle' },
  UserPromptSubmit: { activity: 'working', state: 'active' },
  PostToolUse: { activity: 'working', state: 'active' },
  SubagentStart: { activity: 'working', state: 'active' },
  SubagentStop: { activity: 'working', state: 'active' },
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

// Pure: previous file content + one snapshot event kind -> next content, or undefined when the event writes nothing.
// Unknown fields of `prev` (tasks, ...) ride along.
export function applyEvent(prev: Snap, name: string, e: Ev, now: string, limits?: Limits): Snap | undefined {
  const m = MAP[name]
  if (!m) return undefined
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
