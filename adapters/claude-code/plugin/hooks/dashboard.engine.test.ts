import { expect, test, mock } from 'claude-code/testing'
import { applyEvent } from './snapshot-writer'

// Engine-level tests of the mod's snapshot writer, /dh:dashboard command and R20 commit trap.
// The world is an in-memory filesystem under the plugin; mock.clock moves only when told.
const CWD = '/proj'
const DH = '/dh-home' // DH_HOME of every test: nothing touches a real ~/.harness
const SNAP = `${DH}/sessions`
const SID = 'sess-1'
const FILE = `${SNAP}/${SID}.json`
const T0 = Date.parse('2026-10-07T10:00:00Z')

type Opts = { files?: Record<string, string>; limits?: { kind: string; percentUsed: number }[]; writeFails?: boolean; toolGate?: boolean; procFails?: boolean; env?: Record<string, string>; harness?: string }

function world(on: any, o: Opts = {}) {
  const files = new Map(Object.entries(o.files ?? {}))
  const registered: any[] = []
  const runs: string[][] = []
  mock.env(on, o.env ?? { DH_HOME: DH })
  const clock = mock.clock(on, { now: T0 })
  const w: any = { exists: 0, sid: SID, writeGate: undefined, toolGate: undefined }
  on('session.id', () => ({ value: w.sid }) as any)
  on('session.model', () => ({ value: 'opus-x' }) as any)
  on('session.cwd', () => ({ value: CWD }) as any)
  on('session.usage', () => ({ value: { startedAt: 0, context: { window: 100, percent: 1, tokens: 1 }, rateLimits: o.limits ?? [] } }) as any)
  on('fs.read', (_$: any, e: any) => { if (!files.has(e.path)) throw new Error('ENOENT'); return { value: files.get(e.path) } as any })
  on('fs.exists', (_$: any, e: any) => { w.exists++; return { value: files.has(e.path) } as any })
  on('fs.write', async (_$: any, e: any) => { if (o.writeFails) throw new Error('EACCES'); if (w.writeGate) await w.writeGate; files.set(e.path, e.text); return { value: undefined } as any })
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
    if (o.procFails) throw new Error('spawn failed')
    // `dh harness-path --json <dir>` answers o.harness (JSON) or fails like a dh without the subcommand.
    // git (repo root lookups) answers "not a repository": the mod falls back to the dir itself.
    if (e.argv[0] === 'git') return { value: { exitCode: 128, stdout: '', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } as any
    if (e.argv[1] === 'harness-path') return { value: { exitCode: o.harness ? 0 : 2, stdout: o.harness ?? '', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } as any
    runs.push(e.argv)
    return { value: { exitCode: 0, stdout: 'http://127.0.0.1:4747/\n', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } as any
  })
  on('ui.status', () => ({ value: undefined }) as any)
  on('ui.toast', () => ({ value: undefined }) as any)
  on('session.start', () => ({ cwd: CWD }) as any)
  on('tool.call', async () => { if (o.toolGate && w.toolGate) await w.toolGate; return { result: 'ok' } as any })
  // The engine's own handling of the events the mod records: nothing but an answer.
  on('prompt.submit', (_$: any, e: any) => ({ text: e.text }) as any)
  on('turn.complete', (_$: any, e: any) => ({ text: e.answer }) as any)
  on('session.end', (_$: any, e: any) => ({ sessionId: e.sessionId }) as any)
  on('agent.spawn', (_$: any, e: any) => ({ agentId: e.description, model: 'm' }) as any)
  return Object.assign(w, { files, registered, runs, clock })
}

// Hot-path hooks (prompt, tool, subagent) do not wait for their snapshot write: give the queue real time to drain.
const settled = () => new Promise(r => setTimeout(r, 20))
const snap = (w: any) => JSON.parse(w.files.get(FILE))
// Engine events the mod maps onto the snapshot (BUG-001): no classic.* hook is involved.
const start = ($: any, extra: any = {}) => $.session.start({ source: 'startup', cwd: CWD, ...extra })
const prompt = ($: any) => $.prompt.submit({ text: 'x' })
const toolUse = ($: any) => $.tool.call({ tool: 'Read', file_path: '/a' } as any)
const TURN = { answer: '', durationMs: 1, isAborted: false, turnId: 't' }
const stop = ($: any) => $.turn.complete({ ...TURN, reason: 'answer' } as any)
const stopFail = ($: any) => $.turn.complete({ ...TURN, reason: 'error' } as any)
const subStop = ($: any, id = 'a1') => $.turn.complete({ ...TURN, reason: 'answer', agentId: id } as any)
const spawn = ($: any, d = 'a1') => $.agent.spawn({ tool_use_id: 't', prompt: 'p', description: d, subagentType: 'builder' } as any)
const end = ($: any) => $.session.end({ reason: 'other', sessionId: SID, resume: {} } as any)

// R13: one test per row of the event table.
const ROWS: [string, ($: any) => Promise<any>, string, string][] = [
  ['SessionStart', $ => start($), 'idle', 'idle'],
  ['prompt.submit', $ => prompt($), 'working', 'active'],
  ['tool.call (tool activity)', $ => toolUse($), 'working', 'active'],
  ['turn.complete reason error', $ => stopFail($), 'error', 'idle'],
  ['turn.complete reason answer', $ => stop($), 'done', 'idle'],
  ['agent.spawn', $ => spawn($), 'working', 'active'],
  ['turn.complete of a subagent', $ => subStop($), 'working', 'active'],
  ['session.end', $ => end($), 'idle', 'closed'],
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

test('BUG-001 the mod registers no classic.* hook: classic events write nothing', async ($, on) => {
  const w = world(on)
  on('classic.Stop' as any, () => ({}) as any)
  on('classic.UserPromptSubmit' as any, () => ({}) as any)
  await ($ as any).classic.Stop({ session_id: SID, stop_hook_active: false })
  await ($ as any).classic.UserPromptSubmit({ session_id: SID, prompt: 'x' })
  await settled()
  expect(w.files.has(FILE)).toBe(false)
})

test('activity_at moves only when activity changes; updated_at moves every event', async ($, on) => {
  const w = world(on)
  await prompt($); await settled()
  await w.clock.advance(5000)
  await toolUse($); await settled()
  expect(snap(w).activity_at).toBe('2026-10-07T10:00:00Z')
  expect(snap(w).updated_at).toBe('2026-10-07T10:00:05Z')
  await w.clock.advance(5000)
  await stop($)
  expect(snap(w).activity_at).toBe('2026-10-07T10:00:10Z')
})

test('preserves unknown fields (tasks) and recovers from partial JSON', async ($, on) => {
  const w = world(on, { files: { [FILE]: JSON.stringify({ schema: 1, session_id: SID, tasks: [{ id: 't1' }], future: { a: 1 }, started_at: 'T0' }) } })
  await prompt($); await settled()
  expect(snap(w).tasks).toEqual([{ id: 't1' }])
  expect(snap(w).future).toEqual({ a: 1 })
  expect(snap(w).started_at).toBe('T0')
})

test('partial JSON is left alone (write skipped), a missing file starts fresh', async ($, on) => {
  const w = world(on, { files: { [FILE]: '{"schema":2,"sess' } })
  await stop($)
  expect(w.files.get(FILE)).toBe('{"schema":2,"sess')
  w.sid = 'fresh'
  await stop($)
  expect(JSON.parse(w.files.get(`${SNAP}/fresh.json`)!).activity).toBe('done')
})

test('a file that stays unparseable is overwritten on the third skipped write', async ($, on) => {
  const w = world(on, { files: { [FILE]: 'garbage' } })
  await stop($)
  await stop($)
  expect(w.files.get(FILE)).toBe('garbage')
  await stop($)
  expect(snap(w).activity).toBe('done')
})

test('SubagentStop keeps the current activity but still records the event', async ($, on) => {
  const w = world(on)
  await stop($)
  await subStop($); await settled()
  expect(snap(w).activity).toBe('done')
  expect(snap(w).events.length).toBe(1)
})

test('a heartbeat never resurrects a closed session or reopens a missing file', async ($, on) => {
  const w = world(on)
  await start($)
  await end($)
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
  for (let i = 0; i < 60; i++) await spawn($, `a${i}`)
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
  w.sid = '../evil'
  await stop($)
  expect([...w.files.keys()]).toEqual([])
})

// PRD-014 R10: <home>/sessions only; DEV_HARNESS_SNAPSHOT_DIR and CLAUDE_CONFIG_DIR no longer pick the folder.
test('R10 with DH_HOME the snapshot goes to $DH_HOME/sessions, old variables ignored', async ($, on) => {
  const w = world(on, { env: { DH_HOME: '/x/h', HOME: '/u', DEV_HARNESS_SNAPSHOT_DIR: '/old', CLAUDE_CONFIG_DIR: '/cfg' } })
  await start($)
  expect([...w.files.keys()]).toEqual([`/x/h/sessions/${SID}.json`])
})
test('R10 without DH_HOME the snapshot goes to $HOME/.harness/sessions', async ($, on) => {
  const w = world(on, { env: { HOME: '/u', DEV_HARNESS_SNAPSHOT_DIR: '/old', CLAUDE_CONFIG_DIR: '/cfg' } })
  await start($)
  expect([...w.files.keys()]).toEqual([`/u/.harness/sessions/${SID}.json`])
})
test('R10 USERPROFILE is the last fallback', async ($, on) => {
  const w = world(on, { env: { USERPROFILE: 'C:\\Users\\u' } })
  await start($)
  // On Linux the engine reads `C:/...` as relative to the working directory; on Windows it is absolute.
  expect([...w.files.keys()].every(k => k.endsWith(`C:/Users/u/.harness/sessions/${SID}.json`))).toBe(true)
  expect(w.files.size).toBe(1)
})
test('R10 no home variable writes nothing', async ($, on) => {
  const w = world(on, { env: {} })
  w.sid = 'none'
  await stop($)
  expect([...w.files.keys()]).toEqual([])
})

// R13b parity with `dh snapshot event`. GOLDEN is the time-stripped output of the Go binary for this same
// input sequence (printf '<json>' | dh snapshot event, then with its snapshot dir in a temp folder), captured 2026-10-07.
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
  await end($)
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
  await prompt($); await settled()
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

// PRD-014 R19: global mode reads tickets, project.yaml and PRDs from the resolved folder.
// Last in the file: the mod caches a global answer for CWD for the rest of the module's life.
const G = `${DH}/projects/proj`
const GLOBAL = JSON.stringify({ mode: 'global', dir: G, defaults: {} })
test('R19 global drift: tickets under <home>/projects/<name>/tasks are checked', async ($, on) => {
  world(on, { harness: GLOBAL, files: {
    [`${G}/tasks/T-1/TASK.md`]: '| Status | pronta |\n| PRD (RF-<n>) | PRD-001 (R1) |\n',
    [`${CWD}/docs/prd/PRD-001-x.md`]: '| Status | entregue em 2026-01-01 |\n',
    [`${G}/project.yaml`]: 'language: pt-br\n',
  } })
  const r = await commit($)
  expect(String(r.deny ?? r.text)).toContain('T-1')
  expect(String(r.deny ?? r.text)).toContain('commite de novo')
})
test('R19 global drift: PRDs under <home>/projects/<name>/prd count too', async ($, on) => {
  world(on, { harness: GLOBAL, files: {
    [`${G}/tasks/T-1/TASK.md`]: '| Status | concluída |\n| PRD (RF-<n>) | PRD-001 (R1) |\n',
    [`${G}/prd/PRD-001-x.md`]: '| Status | aprovado |\n',
  } })
  expect(String((await commit($)).deny ?? '')).toContain('PRD-001')
})
test('R19 global mode ignores <repo>/.harness/tasks (the resolved folder wins)', async ($, on) => {
  world(on, { harness: GLOBAL, files: TASKS('pronta', 'entregue em 2026-01-01') })
  expect(String((await commit($)).deny ?? '')).not.toContain('drift')
})
test('R19 no tickets in either place is no drift (CI case)', async ($, on) => {
  world(on, { harness: GLOBAL, files: { [`${CWD}/docs/prd/PRD-001-x.md`]: '| Status | entregue em 2026-01-01 |\n' } })
  expect(String((await commit($)).deny ?? '')).not.toContain('drift')
})

// BUG-001: the cycle on engine events only (a Team account floors every classic.* hook).
test('BUG-001 full cycle writes the snapshot with rate_limits and the expected activity transitions', async ($, on) => {
  const w = world(on, { limits: [{ kind: 'five_hour', percentUsed: 12 }, { kind: 'seven_day', percentUsed: 34 }] })
  await start($)
  expect(snap(w)).toMatchObject({ schema: 2, session_id: SID, state: 'idle', activity: 'idle', cwd: CWD, model: { id: 'opus-x' }, rate_limits: { five_hour: 12, seven_day: 34 }, started_at: '2026-10-07T10:00:00Z' })
  await prompt($); await settled()
  expect(snap(w)).toMatchObject({ state: 'active', activity: 'working' })
  await toolUse($); await settled()
  expect(snap(w).activity).toBe('working')
  await spawn($, 'a1'); await settled()
  expect(snap(w).events).toMatchObject([{ event: 'SubagentStart', agent_type: 'builder', agent_id: 'a1' }])
  await subStop($, 'a1'); await settled()
  expect(snap(w).activity).toBe('working') // a subagent turn does not overwrite the main activity
  expect(snap(w).events.map((e: any) => e.event)).toEqual(['SubagentStart', 'SubagentStop'])
  await stop($)
  expect(snap(w)).toMatchObject({ state: 'idle', activity: 'done' })
  await subStop($, 'a2'); await settled()
  expect(snap(w).activity).toBe('done')
  await w.clock.advance(30_000)
  expect(snap(w).updated_at).toBe('2026-10-07T10:00:30Z')
  await end($)
  expect(snap(w)).toMatchObject({ state: 'closed', activity: 'idle' })
  const closed = snap(w)
  await w.clock.advance(120_000)
  expect(snap(w)).toEqual(closed)
})

test('BUG-001 guards are unchanged by the snapshot: Bash AI mention still denied, tool.call still recorded', async ($, on) => {
  const w = world(on)
  const r: any = await $.tool.call({ tool: 'Bash', command: 'git commit -m "Co-Authored' + '-By: Cla' + 'ude"' } as any)
  expect(String(r.deny)).toContain('mentions AI')
  await settled()
  expect(snap(w).activity).toBe('working')
})

test('BUG-001 a failing snapshot write never changes a guard outcome', async ($, on) => {
  world(on, { writeFails: true })
  const r: any = await $.tool.call({ tool: 'Read', file_path: '/a' } as any)
  expect(r.result).toBe('ok')
  await settled()
})

const gate = () => { let open!: () => void; const p = new Promise<void>(r => { open = r }); return { p, open } }
const fileOf = (w: any, id: string) => JSON.parse(w.files.get(`${SNAP}/${id}.json`))

// BUG-001 round 1, finding 1: /clear continues the process under a new id and no session.start fires.
test('BUG-001 after /clear the first write for the new id carries cwd, started_at and model; the old file stays closed; the heartbeat follows', async ($, on) => {
  const w = world(on)
  await start($)
  await $.session.end({ reason: 'clear', sessionId: SID, resume: {} } as any)
  expect(fileOf(w, SID).state).toBe('closed')
  await w.clock.advance(5000)
  w.sid = 'sess-2'
  await prompt($); await settled()
  expect(fileOf(w, 'sess-2')).toMatchObject({ session_id: 'sess-2', state: 'active', activity: 'working', cwd: CWD, model: { id: 'opus-x' }, started_at: '2026-10-07T10:00:05Z' })
  expect(fileOf(w, SID).state).toBe('closed')
  await w.clock.advance(30_000)
  expect(fileOf(w, 'sess-2').updated_at).toBe('2026-10-07T10:00:35Z') // heartbeat moved to the new id
  expect(fileOf(w, SID).updated_at).toBe('2026-10-07T10:00:00Z')
})

// Finding 2: id captured at event time; a closed file is never reopened by a non-start write.
test('BUG-001 a record queued behind a slow job lands on the id it was raised under', async ($, on) => {
  const w = world(on)
  await start($)
  const g = gate(); w.writeGate = g.p
  await prompt($)  // job 1 hangs in fs.write
  await toolUse($)  // queued behind it
  w.sid = 'sess-2'
  w.writeGate = undefined; g.open(); await settled()
  expect([...w.files.keys()].sort()).toEqual([FILE])
})

test('BUG-001 a tool record that drains after session.end does not reopen the closed file', async ($, on) => {
  const w = world(on, { toolGate: true })
  await start($)
  const g = gate(); w.toolGate = g.p
  const call = toolUse($) // the tool is still running when the session ends
  await end($)
  expect(snap(w).state).toBe('closed')
  g.open(); await call; await settled()
  expect(snap(w)).toMatchObject({ state: 'closed', activity: 'idle' })
})

// Finding 3: session.end has 1.5s in all.
test('BUG-001 session.end stops waiting for the chain after 1s and the close write still lands', async ($, on) => {
  const w = world(on)
  await start($)
  const g = gate(); w.writeGate = g.p
  let done = false
  const ending = end($).then(() => { done = true })
  await settled()
  expect(done).toBe(false)
  await w.clock.advance(1000)
  await ending
  expect(done).toBe(true)
  w.writeGate = undefined; g.open(); await settled()
  expect(snap(w).state).toBe('closed')
  const closed = snap(w)
  await w.clock.advance(60_000) // the heartbeat was cancelled on the capped path too
  expect(snap(w)).toEqual(closed)
})

test('BUG-001 round 2: session.end cancels the heartbeat at once, a hung close job cannot leave it running', async ($, on) => {
  const w = world(on)
  await start($)
  const g = gate(); w.writeGate = g.p
  const ending = end($)
  await w.clock.advance(1000); await ending
  const before = w.exists
  await w.clock.advance(60_000) // a live beat would queue touch jobs behind the hung close job
  w.writeGate = undefined; g.open(); await settled()
  expect(w.exists - before).toBe(0) // the close job read before it hung; a live beat would add a touch read
  expect(snap(w).state).toBe('closed')
})

test('BUG-001 round 2: a file closed by another process does not freeze a session continued in-process', async ($, on) => {
  const w = world(on, { files: { [`${SNAP}/sess-b.json`]: JSON.stringify({ schema: 2, session_id: 'sess-b', state: 'closed', activity: 'idle', started_at: 'T0', cwd: CWD, model: { id: 'm' } }) } })
  await start($)
  await end($)
  w.sid = 'sess-b' // /resume of an old id: no session.start, file closed by the previous process
  await prompt($); await settled()
  expect(fileOf(w, 'sess-b')).toMatchObject({ state: 'active', activity: 'working', started_at: 'T0' })
  await w.clock.advance(30_000)
  expect(fileOf(w, 'sess-b').updated_at).toBe('2026-10-07T10:00:30Z') // heartbeat follows B
})

test('BUG-001 round 2: a prompt after session.end of the same id (in-process resume) reopens it; a late tool record does not', async ($, on) => {
  const w = world(on)
  await start($)
  await end($)
  await toolUse($); await settled()
  expect(snap(w).state).toBe('closed') // late record of the id this instance closed
  await prompt($); await settled()
  expect(snap(w)).toMatchObject({ state: 'active', activity: 'working' })
  await toolUse($); await settled()
  expect(snap(w).state).toBe('active')
  await w.clock.advance(30_000)
  expect(snap(w).updated_at).toBe('2026-10-07T10:00:30Z')
})

test('BUG-001 session.end waits for the close write when the chain is quick', async ($, on) => {
  const w = world(on)
  await start($)
  await end($)
  expect(snap(w).state).toBe('closed') // awaited: written before end() resolved
})
