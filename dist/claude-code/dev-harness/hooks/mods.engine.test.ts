import { expect, mock, test } from 'claude-code/testing'
import type { On } from 'claude-code'

import { backgroundPill, barFill, cacheState, CACHE_TTL_MS, colorFor, fmtDuration, fmtTokens } from './lib/status.ts'
import { barCells, parseZenInput, percentOf, STEPS_MAX, stepLabel, TEXT_MAX, visibleSteps, ZEN_PROMPT, ZEN_TOOL, ZEN_TOOL_NAME } from './lib/zen.ts'
import { configDir } from './lib/settings.ts'
import { modeLabel } from './lib/modes.ts'

// Engine-level tests of the mods menu (mods.tsx, merged into panel.ts) and pure tests of lib/.
// The engine names a plugin tested from its folder after the folder (`plugin`); in the package it is `dh` (plugin.json),
// which is what mcp__dh__zen_progress and the atoms' key assume.
const PLUGIN = 'dh'
const SURFACES = ['terminal', 'desktop'] as const

const band = (bodyColumns = 80) => ({
  component: 'AbovePrompt' as const,
  props: { hasSurvey: false, isWorking: false, maxRows: 20, bodyColumns, scroll: { offset: 0, bodyRows: 20 }, view: {} },
})
const pane = {
  component: 'Pane' as const,
  requestId: 'dh-mods',
  props: { title: 'dh mods', isFocused: true, bodyColumns: 56, placement: 'dock' as const, scroll: { offset: 0, bodyRows: 26 }, view: {} },
}
const hint = {
  component: 'PromptHint' as const,
  props: { isDraft: false, isWorking: false, hint: '? for shortcuts' },
}

// Stands in for the engine's own drawing where the mod passes; the store is in memory.
const engineDraws = (on: On, opened: string[] = []) => {
  mock.store(on)
  on('ui.render', ($, e) => $.ui.resolve(e).Text({ children: 'engine' }))
  on('ui.open', (_$, e) => { opened.push((e as any).id); return { value: { isPlaced: true as const } } })
  on('ui.close', () => ({ value: undefined }))
}

const until = async (ok: () => boolean | Promise<boolean>) => {
  for (let i = 0; i < 200 && !(await ok()); i++) await new Promise(r => setTimeout(r, 10))
}

// An in-memory ~/.claude/settings.json under HOME=/home/t, read and written through the engine's fs events.
function homeWorld(on: On, settings?: string, tasks: Record<string, string> = {}) {
  const files = new Map<string, string>()
  const path = '/home/t/.claude/settings.json'
  if (settings !== undefined) files.set(path, settings)
  engineDraws(on)
  mock.clock(on)
  mock.env(on, { HOME: '/home/t' })
  const task = (path: string) => /tasks\/([^/]+)\/TASK.md$/.exec(path)?.[1]
  on('fs.exists', (_$, e) => ({ value: files.has(e.path) || (task(e.path) !== undefined && tasks[task(e.path)!] !== undefined) }) as never)
  on('fs.read', (_$, e) => ({ value: task(e.path) !== undefined ? tasks[task(e.path)!] : files.get(e.path) }) as never)
  on('fs.list', (_$, e) => ({ value: String(e.path).endsWith('/tasks') ? Object.keys(tasks).map(name => ({ name, kind: 'dir', size: 0, mtimeMs: 0 })) : [] }) as never)
  on('fs.write', (_$, e) => { files.set(e.path, e.text); return { value: undefined } as never })
  on('session.model', () => ({ value: 'claude-opus-5-5' }))
  on('session.cwd', () => ({ value: '/work/proj' }))
  on('session.usage', () => ({ value: { context: { percent: 12, tokens: 24_000, window: 200_000 }, rateLimits: [] } }) as never)
  on('process.run', () => ({ value: { exitCode: 0, stdout: 'main\n', stderr: '' } }) as never)
  const bak = `${path}.dh.bak`
  return { files, read: () => JSON.parse(files.get(path)!), text: () => files.get(path), bak, hasBak: () => files.has(bak), bakText: () => files.get(bak) }
}

const USER = { model: 'opus', statusLine: { type: 'command', command: 'echo hi', padding: 1 }, permissions: { allow: ['Bash(ls:*)'] }, env: { A: '1' } }
const { statusLine: _kept, ...REST } = USER

test('status line formats and cache state', () => {
  expect(fmtTokens(306_000)).toBe('306k')
  expect(fmtTokens(1_000_000)).toBe('1M')
  expect(fmtDuration(110 * 60_000)).toBe('1h50m')
  expect(fmtDuration(24 * 3_600_000)).toBe('1d00h')
  expect(barFill(34)).toBe(3)
  expect(colorFor(95)).toBe('#ef4444')
  expect(cacheState(null, 1000)).toBe(null)
  expect(cacheState(0, CACHE_TTL_MS - 1)).toBe('ok')
  expect(cacheState(0, CACHE_TTL_MS)).toBe('over')
})

test("the engine's shell and monitor pill is read out of its hint", () => {
  expect(backgroundPill('? for shortcuts')).toBe(null)
  expect(backgroundPill('1 shell · ? for shortcuts')).toBe('1 shell')
  expect(backgroundPill('esc to interrupt 2 shells, 1 monitor')).toBe('2 shells, 1 monitor')
  expect(backgroundPill('1 monitor')).toBe('1 monitor')
  expect(backgroundPill('3 background tasks · ? for shortcuts')).toBe('3 background tasks')
  expect(backgroundPill('ran 2 shell commands')).toBe(null)
})

const steps = (...s: ('pending' | 'in_progress' | 'completed')[]) => s.map((status, i) => ({ id: String(i), text: `step ${i}`, status }))

test('zen percent and step label follow the task list', () => {
  const t = { title: 'x', startedAt: 0, doneAt: null, steps: steps('completed', 'in_progress', 'pending', 'pending') }
  expect(percentOf(t)).toBe(25)
  expect(stepLabel(t)).toBe('Step 2 of 4')
  expect(percentOf({ ...t, steps: [] })).toBe(0)
  expect(stepLabel({ ...t, steps: [] })).toBe('Working…')
  expect(percentOf({ ...t, doneAt: 5 })).toBe(100)
})

test('zen bar fills its width, in one color or a rainbow', () => {
  expect(barCells('blue', '#3b82f6', 40, 50)).toEqual({ filled: 20, runs: [{ color: '#3b82f6', cells: 20 }] })
  const rainbow = barCells('rainbow', '', 40, 100)
  expect(rainbow.runs).toHaveLength(8)
  expect(rainbow.runs.reduce((n, r) => n + r.cells, 0)).toBe(40)
})

test('zen long lists keep the current step in view', () => {
  const t = { title: 'x', startedAt: 0, doneAt: null, steps: steps(...Array(5).fill('completed'), 'in_progress', 'pending', 'pending', 'pending', 'pending') }
  const shown = visibleSteps(t)
  expect(shown).toHaveLength(6)
  expect(shown.some(s => s.status === 'in_progress')).toBe(true)
})

test('parseZenInput drops null steps, control characters and oversized input', () => {
  const parsed = parseZenInput({
    title: `  T\u001b[31mi\u2028tle${'x'.repeat(500)}  `,
    steps: [null, 7, { text: 5 }, { text: '\u001b]0;pwn\u0007  ok ', status: 'bogus' }, { text: 'a'.repeat(1000), status: 'completed' }],
  })
  expect(parsed.title.length).toBe(TEXT_MAX)
  expect(parsed.title).not.toMatch(/[\u0000-\u001f\u2028]/)
  expect(parsed.steps.map(x => x.text.length)).toEqual([10, TEXT_MAX])
  expect(parsed.steps[0]).toEqual({ id: 'zen-3', text: ']0;pwn  ok', status: 'pending' })
  expect(parseZenInput(null)).toEqual({ title: '', steps: [] })
  expect(parseZenInput({ steps: Array.from({ length: 10_000 }, () => ({ text: 'x', status: 'pending' })) }).steps).toHaveLength(STEPS_MAX)
})

test('R7 the tool is zen_progress and the model calls it mcp__dh__zen_progress', () => {
  expect(ZEN_TOOL_NAME).toBe('zen_progress')
  expect(ZEN_TOOL).toBe('mcp__dh__zen_progress')
  expect(ZEN_PROMPT).toContain('mcp__dh__zen_progress')
})

test('R6 the running permission mode is unknown without an engine event', () => {
  expect(modeLabel(null)).toBe('unknown')
})

test('R3 the corner button names dh and opens a menu with exactly three items', async ($, on) => {
  engineDraws(on)
  for (const surface of SURFACES) {
    const corner = await $.ui.mount({ plugin: PLUGIN, surface, ...hint })
    const menu = await $.ui.mount({ plugin: PLUGIN, surface, ...pane })
    expect((await corner.find({ key: 'mods-toggle' }))?.text).toMatch(/dh mods ▲/)
    expect(await menu.find({ key: 'open-zen' })).toBeUndefined()

    await corner.press({ key: 'mods-toggle' })
    expect((await corner.find({ key: 'mods-toggle' }))?.text).toMatch(/▼/)
    for (const key of ['open-status', 'open-zen', 'open-mode']) expect(await menu.find({ key })).toBeDefined()
    expect(await menu.find({ key: 'open-themes' })).toBeUndefined()
    expect((await menu.findAll({ type: 'Button', text: '›' })).length).toBe(3)
    await corner.press({ key: 'mods-toggle' })
    await corner.unmount()
    await menu.unmount()
  }
})

test('R3 /mods registers on the first draw and opens the menu', async ($, on) => {
  const opened: string[] = []
  engineDraws(on, opened)
  const registered: string[] = []
  on('command.register', (_$, e) => { registered.push((e as any).name); return { value: undefined } as never })
  on('tool.register', (_$, e) => { registered.push((e as any).name); return { value: undefined } as never })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await until(() => registered.includes('mods'))
  expect(registered).toContain('mods')
  expect(registered).toContain(ZEN_TOOL_NAME)
  const r = await $.command.run({ command: 'mods', args: '' })
  expect(String(r.text)).toMatch(/open/)
  expect(opened).toEqual(['dh-mods'])
  await corner.unmount()
})

test('Zen mode and its band color are set from the menu', async ($, on) => {
  engineDraws(on)
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-zen' })
  await menu.press({ key: 'zen-on' })
  await menu.press({ key: 'zen-theme-rainbow' })
  expect((await menu.find({ key: 'zen-on' }))?.text).toMatch(/● On/)
  expect((await menu.find({ key: 'zen-theme-rainbow' }))?.text).toMatch(/● rainbow/)
  await menu.press({ key: 'mods-back' })
  expect(await menu.find({ text: /on · rainbow band/ })).toBeDefined()
  await menu.press({ key: 'open-zen' })
  await menu.press({ key: 'zen-off' })
  expect((await menu.find({ key: 'zen-off' }))?.text).toMatch(/● Off/)
})

test('the mode page lists effort levels and new-session modes', async ($, on) => {
  engineDraws(on)
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  for (const level of ['low', 'medium', 'high', 'xhigh', 'max']) expect(await menu.find({ key: `effort-${level}` })).toBeDefined()
  expect(await menu.find({ key: 'default-mode-bypassPermissions' })).toBeDefined()
  // R6: no permission-mode event reaches the mod, and the page says so instead of failing.
  expect(await menu.find({ text: /unknown/ })).toBeDefined()
  await menu.press({ key: 'effort-high' })
  await menu.press({ key: 'mods-back' })
  expect(await menu.find({ text: /effort high/ })).toBeDefined()
})

test('the mode page switches the model through /model', async ($, on) => {
  engineDraws(on)
  let current = 'claude-opus-5-5'
  const ran: string[] = []
  on('command.run', { command: 'model' }, async (_$, e) => {
    ran.push(e.args)
    current = e.args === 'sonnet' ? 'claude-sonnet-5-5' : current
    return { text: `Set model to ${e.args}` }
  })
  on('session.model', () => ({ value: current }))
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  for (const alias of ['fable', 'opus', 'sonnet', 'haiku', 'default']) expect(await menu.find({ key: `model-${alias}` })).toBeDefined()
  await menu.press({ key: 'model-sonnet' })
  expect(ran).toEqual(['sonnet'])
  expect((await menu.find({ key: 'model-sonnet' }))?.text).toMatch(/● Sonnet/)
  await menu.press({ key: 'mods-back' })
  expect(await menu.find({ text: /Sonnet 5\.5 · / })).toBeDefined()
})

test('R9 after choosing effort-high a turn.step with effort goes on with effort high', async ($, on) => {
  engineDraws(on)
  const seen: unknown[] = []
  on('turn.step', async function* (_$, e) {
    seen.push(e.effort)
    return { answer: 'ok', toolUses: [], usage: null } as never
  })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  await menu.press({ key: 'effort-high' })
  const step = { model: 'claude-opus-5-5', effort: 'low' } as never
  for await (const _ of $.turn.step(step)) { /* drain */ }
  expect(seen).toEqual(['high'])
})

test('a running shell or monitor gets a second row under the status line', async ($, on) => {
  const w = homeWorld(on, '{}')
  const ran: string[] = []
  on('command.run', { command: 'tasks' }, async (_$, e) => { ran.push(e.command); return { text: '' } })
  const busy = { ...hint, props: { ...hint.props, hint: '1 shell, 1 monitor · ? for shortcuts' } }
  const first = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await until(() => w.hasBak())
  await first.unmount()
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...busy })
  expect((await corner.find({ key: 'background-tasks' }))?.text).toMatch(/1 shell, 1 monitor/)
  expect(await corner.find({ text: /Opus 5\.5 · proj/ })).toBeDefined()
  await corner.press({ key: 'background-tasks' })
  expect(ran).toEqual(['tasks'])
})

test("Claude's zen_progress calls fill the band, and Zen is injected into a typed prompt only when on", async ($, on) => {
  engineDraws(on)
  mock.clock(on)
  const prompts: any[] = []
  on('prompt.submit', (_$, e) => { prompts.push(e); return { text: 'ok' } as never })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const progress = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...band() })
  await $.prompt.submit({ text: 'before', origin: { kind: 'composer' } } as never)
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-zen' })
  await menu.press({ key: 'zen-on' })
  await $.prompt.submit({ text: 'Build the dashboard', origin: { kind: 'composer' } } as never)
  expect(prompts[0].context ?? []).not.toContain(ZEN_PROMPT)
  expect(prompts[1].context).toContain(ZEN_PROMPT)

  await $.tool.call({
    tool: ZEN_TOOL,
    title: 'Build the dashboard',
    steps: [
      { text: 'Plan', status: 'completed' },
      { text: 'Fetch data', status: 'in_progress' },
      { text: 'Draw charts', status: 'pending' },
      { text: 'Test', status: 'completed' },
    ],
  } as never)
  expect(await progress.find({ text: /Build the dashboard/ })).toBeDefined()
  expect(await progress.find({ text: /50%/ })).toBeDefined()
  expect(await progress.find({ text: /Step 2 of 4/ })).toBeDefined()
  expect(await progress.find({ text: /Fetch data/ })).toBeDefined()
})

test('R4 the status line ends with the dh progress, gate pending and agents N; without dh it ends at the cache', async ($, on) => {
  homeWorld(on, '{}', { A: '| Status | concluída |', B: '| Status | open |' })
  on('tool.call', () => ({ result: 'ok' }) as never)
  on('agent.spawn', () => ({ agentId: 'ag1', model: 'm' }) as never)
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint, viewport: { columns: 200, rows: 40 } } as never)
  await until(async () => (await corner.find({ text: /ctx/ })) !== undefined)
  const line = async () => String((await corner.find({ text: /ctx/ }))?.text)
  expect(await line()).toMatch(/ctx/)
  expect(await line()).not.toMatch(/dh ▰/)

  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as never)
  await $.agent.spawn({ prompt: 'x', description: 'd' } as never)
  await until(async () => /agents 1/.test(await line()))
  expect(await line()).toMatch(/dh ▰▰▰▱▱ 1\/2  │  gate pending  │  agents 1◆/)
})

test('R4 without tickets and without a pending gate the line has no dh segments', async ($, on) => {
  homeWorld(on, '{}')
  on('tool.call', () => ({ result: 'ok' }) as never)
  on('agent.spawn', () => ({ agentId: 'ag1', model: 'm' }) as never)
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint, viewport: { columns: 200, rows: 40 } } as never)
  await until(async () => (await corner.find({ text: /ctx/ })) !== undefined)
  await $.agent.spawn({ prompt: 'x', description: 'd' } as never)
  await until(async () => /agents 1/.test(String((await corner.find({ text: /ctx/ }))?.text)))
  const line = String((await corner.find({ text: /ctx/ }))?.text)
  expect(line).toMatch(/agents 1◆/)
  expect(line).not.toMatch(/dh ▰|gate pending/)
})

test('R5 on by default: the statusLine moves to the store, the backup is written once, off restores it deep-equal', async ($, on) => {
  const w = homeWorld(on, JSON.stringify(USER))
  const original = w.text()
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  await until(() => w.hasBak())
  await until(() => w.read().statusLine === undefined)
  expect(w.read().statusLine).toBeUndefined()
  expect(w.read()).toEqual(REST)
  expect(w.bakText()).toBe(original)

  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-status' })
  await menu.press({ key: 'status-off' })
  await until(() => w.read().statusLine !== undefined)
  expect(w.read()).toEqual(USER)
  expect(w.bakText()).toBe(original)
  expect([...w.files.keys()].sort()).toEqual(['/home/t/.claude/settings.json', '/home/t/.claude/settings.json.dh.bak'])
})

test('R5 Mode writes permissions.defaultMode, leaves the other keys deep-equal and keeps the first backup', async ($, on) => {
  const w = homeWorld(on, JSON.stringify(REST))
  const original = w.text()
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  await menu.press({ key: 'default-mode-plan' })
  await until(() => w.read().permissions?.defaultMode === 'plan')
  expect(w.read()).toEqual({ ...REST, permissions: { ...REST.permissions, defaultMode: 'plan' } })
  expect(w.bakText()).toBe(original)
  await menu.press({ key: 'default-mode-auto' })
  await until(() => w.read().permissions?.defaultMode === 'auto')
  expect(w.read()).toEqual({ ...REST, permissions: { ...REST.permissions, defaultMode: 'auto' } })
  expect(w.bakText()).toBe(original)
})

test('R5 CLAUDE_CONFIG_DIR is honored, and with no home at all settings.json is not touched', async ($, on) => {
  const w = homeWorld(on, JSON.stringify(REST))
  mock.env(on, { HOME: '/home/t', CLAUDE_CONFIG_DIR: '/cfg' })
  expect(configDir('/cfg/', '/home/t')).toBe('/cfg')
  expect(configDir(undefined, '/home/t')).toBe('/home/t/.claude')
  expect(configDir('rel/cfg', '/home/t')).toBe('/home/t/.claude')
  expect(configDir('C:\\cfg', undefined)).toBe('C:/cfg')
  expect(() => configDir('rel/cfg', undefined)).toThrow()
  expect(() => configDir(undefined, 'relhome')).toThrow()
  expect(() => configDir(undefined, undefined)).toThrow()
  expect(() => configDir('', '')).toThrow()
  expect(w.files.size).toBe(1)
})

test('R5 an invalid settings.json is never overwritten and gets no backup', async ($, on) => {
  const w = homeWorld(on, '{ not json')
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  await new Promise(r => setTimeout(r, 100))
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  await menu.press({ key: 'default-mode-plan' })
  await menu.press({ key: 'mods-back' })
  await menu.press({ key: 'open-status' })
  await menu.press({ key: 'status-off' })
  await new Promise(r => setTimeout(r, 100))
  expect(w.text()).toBe('{ not json')
  expect(w.hasBak()).toBe(false)
})

test('zen_progress is ignored while Zen is off; with Zen on a hostile call is cleaned', async ($, on) => {
  engineDraws(on)
  mock.clock(on)
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const progress = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...band() })
  const off = await $.tool.call({ tool: ZEN_TOOL, steps: [{ text: 'hidden', status: 'pending' }] } as never)
  expect(String(off.result)).toMatch(/off/)
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-zen' })
  await menu.press({ key: 'zen-on' })
  const steps = [null, { text: '\u001b[2Jevil', status: 'pending' }, ...Array.from({ length: 10_000 }, () => ({ text: 's', status: 'pending' }))]
  await $.tool.call({ tool: ZEN_TOOL, title: 'T\u001b[0m', steps } as never)
  expect(await progress.find({ text: /evil/ })).toBeDefined()
  expect(await progress.find({ text: /\u001b/ })).toBeUndefined()
  expect(await progress.find({ text: /more/ })).toBeDefined()
})

test('turn.step applies the chosen effort to the main loop only, never to a subagent', async ($, on) => {
  engineDraws(on)
  const seen: unknown[] = []
  on('turn.step', async function* (_$, e) {
    seen.push([e.agentId ?? 'main', e.effort])
    return { answer: 'ok', toolUses: [], usage: null } as never
  })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-mode' })
  await menu.press({ key: 'effort-max' })
  for await (const _ of $.turn.step({ model: 'claude-opus-5-5', effort: 'low', agentId: 'sub1' } as never)) { /* drain */ }
  for await (const _ of $.turn.step({ model: 'claude-opus-5-5', effort: 'low' } as never)) { /* drain */ }
  expect(seen).toEqual([['sub1', 'low'], ['main', 'max']])
})

test('zen_progress with a failing state answers with a result, not an error', async ($, on) => {
  engineDraws(on)
  let down = false
  on('state.get' as never, (_$: any, e: any, next: any) => { if (down) throw new Error('state down'); return next(e) })
  on('state.set' as never, (_$: any, e: any, next: any) => { if (down) throw new Error('state down'); return next(e) })
  const corner = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...hint })
  const menu = await $.ui.mount({ plugin: PLUGIN, surface: 'terminal', ...pane })
  await corner.press({ key: 'mods-toggle' })
  await menu.press({ key: 'open-zen' })
  await menu.press({ key: 'zen-on' })
  down = true
  const r = await $.tool.call({ tool: ZEN_TOOL, steps: [{ text: 'a', status: 'pending' }] } as never)
  expect(String(r.result)).toBe('Progress band unavailable.')
  expect(r.isError).not.toBe(true)
})
