import { expect, test, mock } from 'claude-code/testing'
import { applyEvent } from './snapshot-writer'

// Engine-level tests of the mod's snapshot writer, /dh:dashboard command and R20 commit trap.
// The world is an in-memory filesystem under the plugin; mock.clock moves only when told.
const CWD = '/proj'
const SNAP = '/snap'
const SID = 'sess-1'
const FILE = `${SNAP}/${SID}.json`
const T0 = Date.parse('2026-10-07T10:00:00Z')

const CLASSIC = ['SessionStart', 'UserPromptSubmit', 'PostToolUse', 'PermissionRequest', 'Notification', 'StopFailure', 'Stop', 'SubagentStart', 'SubagentStop', 'SessionEnd']

type Opts = { files?: Record<string, string>; limits?: { kind: string; percentUsed: number }[]; writeFails?: boolean; procFails?: boolean }

function world(on: any, o: Opts = {}) {
  const files = new Map(Object.entries(o.files ?? {}))
  const registered: any[] = []
  const runs: string[][] = []
  mock.env(on, { DEV_HARNESS_SNAPSHOT_DIR: SNAP })
  const clock = mock.clock(on, { now: T0 })
  on('session.cwd', () => ({ value: CWD }) as any)
  on('session.usage', () => ({ value: { startedAt: 0, context: { window: 100, percent: 1, tokens: 1 }, rateLimits: o.limits ?? [] } }) as any)
  on('fs.read', (_$: any, e: any) => { if (!files.has(e.path)) throw new Error('ENOENT'); return { value: files.get(e.path) } as any })
  on('fs.exists', (_$: any, e: any) => ({ value: files.has(e.path) }) as any)
  on('fs.write', (_$: any, e: any) => { if (o.writeFails) throw new Error('EACCES'); files.set(e.path, e.text); return { value: undefined } as any })
  on('fs.list', (_$: any, e: any) => {
    const prefix = e.path + '/'
    const seen = new Map<string, string>()
    for (const k of files.keys()) if (k.startsWith(prefix)) {
      const rest = k.slice(prefix.length)
      seen.set(rest.split('/')[0]!, rest.includes('/') ? 'dir' : 'file')
    }
    if (!seen.size) throw new Error('ENOENT')
    return { value: [...seen].map(([name, kind]) => ({ name, kind, size: 0, mtimeMs: 0 })) } as any
  })
  on('command.register', (_$: any, e: any) => { registered.push(e); return { value: undefined } as any })
  on('process.run', (_$: any, e: any) => {
    runs.push(e.argv)
    if (o.procFails) throw new Error('spawn failed')
    return { value: { exitCode: 0, stdout: 'http://127.0.0.1:4747/\n', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } as any
  })
  on('ui.status', () => ({ value: undefined }) as any)
  on('ui.toast', () => ({ value: undefined }) as any)
  on('session.start', () => ({ cwd: CWD }) as any)
  on('tool.call', () => ({ result: 'ok' }) as any)
  for (const n of CLASSIC) on(`classic.${n}` as any, () => ({}) as any)
  return { files, registered, runs, clock }
}

// Hot-path hooks (prompt, tool, subagent) do not wait for their snapshot write: give the queue real time to drain.
const settled = () => new Promise(r => setTimeout(r, 20))
const snap = (w: any) => JSON.parse(w.files.get(FILE))
const start = ($: any, extra: any = {}) => $.classic.SessionStart({ session_id: SID, source: 'startup', cwd: CWD, agent_type: 'coordinator', model: 'opus-x', ...extra })

// R13: one test per row of the event table.
const ROWS: [string, ($: any) => Promise<any>, string, string][] = [
  ['SessionStart', $ => start($), 'idle', 'idle'],
  ['UserPromptSubmit', $ => $.classic.UserPromptSubmit({ session_id: SID, prompt: 'x' }), 'working', 'active'],
  ['PostToolUse (tool activity)', $ => $.classic.PostToolUse({ session_id: SID, tool_name: 'Bash', tool_input: {}, tool_response: {}, tool_use_id: 't' }), 'working', 'active'],
  ['PermissionRequest', $ => $.classic.PermissionRequest({ session_id: SID, tool_name: 'Bash', tool_input: {} }), 'waiting', 'active'],
  ['Notification permission_prompt', $ => $.classic.Notification({ session_id: SID, message: 'm', notification_type: 'permission_prompt' }), 'waiting', 'active'],
  ['StopFailure', $ => $.classic.StopFailure({ session_id: SID, error: 'rate_limit' }), 'error', 'idle'],
  ['Stop', $ => $.classic.Stop({ session_id: SID, stop_hook_active: false }), 'done', 'idle'],
  ['SubagentStart', $ => $.classic.SubagentStart({ session_id: SID, agent_id: 'a1', agent_type: 'builder' }), 'working', 'active'],
  ['SubagentStop', $ => $.classic.SubagentStop({ session_id: SID, agent_id: 'a1', agent_type: 'builder', stop_hook_active: false, agent_transcript_path: '' }), 'working', 'active'],
  ['SessionEnd', $ => $.classic.SessionEnd({ session_id: SID, reason: 'other' }), 'idle', 'closed'],
]
for (const [name, fire, activity, state] of ROWS) {
  test(`R13 ${name} -> activity ${activity}, state ${state}`, async ($, on) => {
    const w = world(on)
    await fire($)
    await settled()
    const s = snap(w)
    expect(s.schema).toBe(2)
    expect(s.session_id).toBe(SID)
    expect(s.activity).toBe(activity)
    expect(s.state).toBe(state)
    expect(s.activity_at).toBe('2026-10-07T10:00:00Z')
    expect(s.updated_at).toBe('2026-10-07T10:00:00Z')
  })
}

test('R13 Notification of another type writes nothing', async ($, on) => {
  const w = world(on)
  await $.classic.Notification({ session_id: SID, message: 'm', notification_type: 'idle_prompt' })
  expect(w.files.has(FILE)).toBe(false)
})

test('waiting leaves on the next activity', async ($, on) => {
  const w = world(on)
  await $.classic.PermissionRequest({ session_id: SID, tool_name: 'Bash', tool_input: {} })
  expect(snap(w).activity).toBe('waiting')
  await $.classic.PostToolUse({ session_id: SID, tool_name: 'Bash', tool_input: {}, tool_response: {}, tool_use_id: 't' }); await settled()
  expect(snap(w).activity).toBe('working')
})

test('activity_at moves only when activity changes; updated_at moves every event', async ($, on) => {
  const w = world(on)
  await $.classic.UserPromptSubmit({ session_id: SID, prompt: 'x' }); await settled()
  await w.clock.advance(5000)
  await $.classic.PostToolUse({ session_id: SID, tool_name: 'Bash', tool_input: {}, tool_response: {}, tool_use_id: 't' }); await settled()
  expect(snap(w).activity_at).toBe('2026-10-07T10:00:00Z')
  expect(snap(w).updated_at).toBe('2026-10-07T10:00:05Z')
  await w.clock.advance(5000)
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  expect(snap(w).activity_at).toBe('2026-10-07T10:00:10Z')
})

test('preserves unknown fields (tasks) and recovers from partial JSON', async ($, on) => {
  const w = world(on, { files: { [FILE]: JSON.stringify({ schema: 1, session_id: SID, tasks: [{ id: 't1' }], future: { a: 1 }, started_at: 'T0' }) } })
  await $.classic.UserPromptSubmit({ session_id: SID, prompt: 'x' }); await settled()
  expect(snap(w).tasks).toEqual([{ id: 't1' }])
  expect(snap(w).future).toEqual({ a: 1 })
  expect(snap(w).started_at).toBe('T0')
})

test('partial JSON is left alone (write skipped), a missing file starts fresh', async ($, on) => {
  const w = world(on, { files: { [FILE]: '{"schema":2,"sess' } })
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  expect(w.files.get(FILE)).toBe('{"schema":2,"sess')
  await $.classic.Stop({ session_id: 'fresh', stop_hook_active: false })
  expect(JSON.parse(w.files.get(`${SNAP}/fresh.json`)!).activity).toBe('done')
})

test('a file that stays unparseable is overwritten on the third skipped write', async ($, on) => {
  const w = world(on, { files: { [FILE]: 'garbage' } })
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  expect(w.files.get(FILE)).toBe('garbage')
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  expect(snap(w).activity).toBe('done')
})

test('SubagentStop keeps the current activity but still records the event', async ($, on) => {
  const w = world(on)
  await $.classic.Stop({ session_id: SID, stop_hook_active: false })
  await $.classic.SubagentStop({ session_id: SID, agent_id: 'a1', agent_type: 'builder', stop_hook_active: false, agent_transcript_path: '' }); await settled()
  expect(snap(w).activity).toBe('done')
  expect(snap(w).events.length).toBe(1)
})

test('a heartbeat never resurrects a closed session or reopens a missing file', async ($, on) => {
  const w = world(on)
  await start($)
  await $.classic.SessionEnd({ session_id: SID, reason: 'other' })
  w.files.delete(FILE)
  await w.clock.advance(60_000)
  expect(w.files.has(FILE)).toBe(false)
})

test('git commit-graph / commit-tree are not the commit trap', async ($, on) => {
  world(on, { files: TASKS('pronta', 'entregue em 2026-01-01') })
  const r: any = await $.tool.call({ tool: 'Bash', command: 'git commit-graph write' } as any)
  expect(String(r.deny ?? r.text)).not.toContain('drift')
})

test('first PRD file per id is the lowest name', async ($, on) => {
  world(on, { files: TASKS('pronta', 'entregue em 2026-01-01', { [`${CWD}/docs/prd/PRD-001-a.md`]: '| Status | aprovado |\n' }) })
  // PRD-001-a sorts before PRD-001-x: it is the one read, and it is not delivered, so no drift.
  const r: any = await commit($)
  expect(String(r.deny ?? r.text)).not.toContain('drift')
})

test('R13b started_at only the first time; events ring of 50', async ($, on) => {
  const w = world(on)
  await start($)
  await w.clock.advance(1000)
  await start($)
  expect(snap(w).started_at).toBe('2026-10-07T10:00:00Z')
  for (let i = 0; i < 60; i++) await $.classic.SubagentStart({ session_id: SID, agent_id: `a${i}`, agent_type: 'builder' })
  await settled()
  const ev = snap(w).events
  expect(ev.length).toBe(50)
  expect(ev[49].agent_id).toBe('a59')
  expect(ev[0].agent_id).toBe('a10')
})

test('R14 rate_limits come from session usage and are kept when usage has none', async ($, on) => {
  const w = world(on, { limits: [{ kind: 'five_hour', percentUsed: 23.5 }, { kind: 'seven_day', percentUsed: 7 }, { kind: 'spend_limit', percentUsed: 1 }] })
  await start($)
  expect(snap(w).rate_limits).toEqual({ five_hour: 23.5, seven_day: 7 })
})

test('R14 rate_limits are kept when usage has none', async ($, on) => {
  const w = world(on, { files: { [FILE]: JSON.stringify({ rate_limits: { five_hour: 3 } }) } })
  await start($)
  expect(snap(w).rate_limits).toEqual({ five_hour: 3 })
})

test('write failure and unsafe session id are silent no-ops', async ($, on) => {
  const w = world(on, { writeFails: true })
  await start($)
  expect(w.files.has(FILE)).toBe(false)
})

test('unsafe session id writes nothing', async ($, on) => {
  const w = world(on)
  await $.classic.Stop({ session_id: '../evil', stop_hook_active: false })
  expect([...w.files.keys()]).toEqual([])
})

// R13b parity with `dh snapshot event`. GOLDEN is the time-stripped output of the Go binary for this same
// input sequence (printf '<json>' | dh snapshot event, DEV_HARNESS_SNAPSHOT_DIR=tmp), captured 2026-10-07.
const SEQ: [string, any][] = [
  ['SessionStart', { session_id: 's1', cwd: '/p', agent_type: 'coordinator', model: 'opus-x', session_name: 'nm' }],
  ['UserPromptSubmit', { session_id: 's1' }],
  ['SubagentStart', { session_id: 's1', agent_id: 'a1', agent_type: 'builder' }],
  ['SubagentStop', { session_id: 's1', agent_id: 'a1', agent_type: 'builder' }],
  ['Stop', { session_id: 's1' }],
]
const GO_STATE = ['idle', 'active', 'active', 'active', 'idle']
const GO_FINAL = {
  session_id: 's1', session_name: 'nm', cwd: '/p', agent: 'coordinator', model: { id: 'opus-x' }, state: 'idle',
  events: [{ event: 'SubagentStart', agent_type: 'builder', agent_id: 'a1' }, { event: 'SubagentStop', agent_type: 'builder', agent_id: 'a1' }],
}
const strip = (s: any) => {
  const { schema, activity, activity_at, updated_at, started_at, rate_limits, ...rest } = s
  return { ...rest, events: s.events?.map(({ at, ...r }: any) => r) }
}

test('R13b mod output equals dh snapshot event (time fields excluded)', () => {
  let s: any = {}
  SEQ.forEach(([name, e], i) => {
    s = applyEvent(s, name, e, `2026-10-07T10:00:0${i}Z`)
    expect(s.state).toBe(GO_STATE[i])
  })
  expect(strip(s)).toEqual(GO_FINAL)
  expect(s.started_at).toBe('2026-10-07T10:00:00Z')
  expect(applyEvent(s, 'SessionEnd', { session_id: 's1' }, 'x')!.state).toBe('closed')
})

test('heartbeat refreshes only updated_at every 30s and stops at SessionEnd', async ($, on) => {
  const w = world(on)
  await start($)
  const before = snap(w)
  await w.clock.advance(30_000)
  const after = snap(w)
  expect(after.updated_at).toBe('2026-10-07T10:00:30Z')
  expect({ ...after, updated_at: '' }).toEqual({ ...before, updated_at: '' })
  await $.classic.SessionEnd({ session_id: SID, reason: 'other' })
  const closed = snap(w)
  await w.clock.advance(120_000)
  expect(snap(w)).toEqual(closed)
})

test('R6 mod registers the dashboard command and runs dh dashboard --detach, returning the URL', async ($, on) => {
  const w = world(on)
  await $.session.start({ source: 'startup', cwd: CWD })
  expect(w.registered.map((c: any) => c.name)).toEqual(['dashboard'])
  const r: any = await $.command.run({ command: 'dashboard', args: '', origin: { kind: 'composer' }, presentation: { layout: 'main', columns: 80 } } as any)
  expect(r.text).toBe('http://127.0.0.1:4747/')
  expect(w.runs[0]!.slice(1)).toEqual(['dashboard', '--detach'])
  expect(w.runs[0]![0]).toMatch(/bin\/dh(\.cmd)?$/)
})

test('heartbeat refreshes rate_limits from usage and keeps them when usage fails', async ($, on) => {
  const w = world(on, { files: { [FILE]: JSON.stringify({ schema: 2, session_id: SID, state: 'active', activity: 'working', rate_limits: { five_hour: 3 } }) }, limits: [{ kind: 'five_hour', percentUsed: 41 }] })
  await $.classic.UserPromptSubmit({ session_id: SID, prompt: 'x' }); await settled()
  expect(snap(w).rate_limits).toEqual({ five_hour: 41 })
  const st = snap(w)
  w.files.set(FILE, JSON.stringify({ ...st, rate_limits: { five_hour: 3 } })) // stale limit on disk
  await w.clock.advance(30_000)
  expect(snap(w).rate_limits).toEqual({ five_hour: 41 })
  expect(snap(w).updated_at).toBe('2026-10-07T10:00:30Z')
})

test('R6 extra args with shell metacharacters run nothing', async ($, on) => {
  const w = world(on)
  for (const args of ['& calc', '--port=1;rm', '$(x)', '--a|b', '"q"']) {
    const r: any = await $.command.run({ command: 'dashboard', args, origin: { kind: 'composer' }, presentation: { layout: 'main', columns: 80 } } as any)
    expect(r.text).toContain('unsupported argument')
  }
  expect(w.runs.length).toBe(0)
  const ok: any = await $.command.run({ command: 'dashboard', args: '--port=4748', origin: { kind: 'composer' }, presentation: { layout: 'main', columns: 80 } } as any)
  expect(ok.text).toBe('http://127.0.0.1:4747/')
  expect(w.runs[0]!.slice(1)).toEqual(['dashboard', '--detach', '--port=4748'])
})

test('R6 a failing dh is reported in the command text, not thrown', async ($, on) => {
  world(on, { procFails: true })
  const r: any = await $.command.run({ command: 'dashboard', args: '', origin: { kind: 'composer' }, presentation: { layout: 'main', columns: 80 } } as any)
  expect(r.text).toContain('failed')
})

// R20 trap
const TASKS = (status: string, prdStatus: string, extra: Record<string, string> = {}) => ({
  [`${CWD}/.harness/tasks/T-1/TASK.md`]: `| Status | ${status} |\n| PRD (RF-<n>) | PRD-001 (R1) |\n`,
  [`${CWD}/docs/prd/PRD-001-x.md`]: `| Status | ${prdStatus} |\n`,
  ...extra,
})
const commit = ($: any) => $.tool.call({ tool: 'Bash', command: 'git commit -m "feat: x"' } as any) as Promise<any>

test('R20 drift (a) denies with the tickets listed', async ($, on) => {
  world(on, { files: TASKS('pronta', 'entregue em 2026-01-01') })
  const r = await commit($)
  expect(String(r.deny ?? r.text)).toContain('T-1')
  expect(String(r.deny ?? r.text)).toContain('PRD-001')
  expect(String(r.deny ?? r.text)).toContain('commit again')
})

test('R20 drift (b) denies; message follows pt-br', async ($, on) => {
  world(on, { files: TASKS('concluída', 'aprovado', { [`${CWD}/.harness/project.yaml`]: 'language: pt-br\n' }) })
  const r = await commit($)
  expect(String(r.deny ?? r.text)).toContain('PRD-001')
  expect(String(r.deny ?? r.text)).toContain('commite de novo')
})

test('R20 no drift does not deny', async ($, on) => {
  world(on, { files: TASKS('concluída', 'entregue em 2026-01-01') })
  const r = await commit($)
  expect(String(r.deny ?? r.text)).not.toContain('dev-harness')
})

test('non-commit commands are never checked', async ($, on) => {
  world(on, { files: TASKS('pronta', 'entregue em 2026-01-01') })
  const ls: any = await $.tool.call({ tool: 'Bash', command: 'ls' } as any)
  expect(ls.deny).toBe(undefined)
})

test('R20 read error does not deny (fail open)', async ($, on) => {
  world(on, { files: {} })
  expect(String((await commit($)).deny ?? '')).not.toContain('drift')
})
