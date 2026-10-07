import { expect, test } from 'claude-code/testing'

// Engine-level tests of the composed module (panel.ts): hooks sit above the test's own `on`
// hooks, which stand for the engine (fs, session, tools). Module state is shared across
// tests, so each test starts by running a green gate Bash to clear the dirty flag.
// AI-mention strings are built by concatenation on purpose.
const AI = 'Co-Authored' + '-By: Cla' + 'ude'
const CWD = '/proj'

type World = { editFails?: boolean; yaml?: string; tasks?: Record<string, string>; fsFails?: boolean; roots?: Record<string, string> }

function world(on: any, w: World = {}) {
  const statuses: string[] = []
  on('session.cwd', () => ({ value: CWD }) as any)
  on('session.usage', () => ({ value: { startedAt: 0, context: { window: 100, percent: 12, tokens: 12 }, rateLimits: [] } }) as any)
  on('fs.read', (_$: any, e: any) => {
    const path = e.path
    if (w.fsFails) throw new Error('boom')
    if (path === `${CWD}/.harness/project.yaml` && w.yaml !== undefined) return { value: w.yaml } as any
    const m = /tasks\/([^/]+)\/TASK.md$/.exec(path)
    if (m && w.tasks?.[m[1]] !== undefined) return { value: w.tasks[m[1]] } as any
    throw new Error('ENOENT')
  })
  on('fs.exists', (_$: any, { path }: any) => {
    const m = /tasks\/([^/]+)\/TASK.md$/.exec(path)
    return { value: !!(m && w.tasks?.[m[1]] !== undefined) } as any
  })
  on('fs.list', () => ({ value: Object.keys(w.tasks ?? {}).map(name => ({ name, kind: 'dir', size: 0, mtimeMs: 0 })) }) as any)
  // `git -C <dir> rev-parse --show-toplevel`: known dirs answer their root, others fail (the mod falls back to the dir).
  on('process.run', (_$: any, e: any) => {
    const root = w.roots?.[e.argv[2]]
    return { value: { exitCode: root ? 0 : 128, stdout: root ? `${root}\n` : '', stderr: '', isStdoutTruncated: false, isStderrTruncated: false } } as any
  })
  on('ui.status', (_$: any, e: any) => { statuses.push(e.text); return { value: undefined } as any })
  on('ui.toast', () => ({ value: undefined }) as any)
  on('session.start', () => ({ cwd: CWD }) as any)
  if (w.editFails) on('tool.call', { tool: 'Edit' }, () => ({ result: 'no match', isError: true }) as any)
  on('tool.call', () => ({ result: 'ok' }) as any)
  on('agent.spawn', () => ({ agentId: 'ag1', model: 'm' }) as any)
  return { statuses }
}

const bash = ($: any, command: string) => $.tool.call({ tool: 'Bash', command } as any) as Promise<any>
const text = (r: any) => String(r.deny ?? r.text ?? '')
const clean = async ($: any) => { await bash($, 'go test ./...'); await bash($, 'cd /other && go test ./...') }

test('1 AI mention in commit is denied by the panel; clean commit is not', async ($, on) => {
  world(on)
  await clean($)
  const bad = await bash($, `git commit -m "x\n\n${AI}"`)
  expect(text(bad)).toContain('dev-harness')
  expect(text(bad)).toContain('mentions AI')
  const ok = await bash($, 'git commit -m "feat: x"')
  expect(ok.deny).toBe(undefined)
  expect(text(ok)).not.toContain('dev-harness')
})

test('2 edit then commit is denied naming gate commands; green gate re-allows', async ($, on) => {
  world(on, { yaml: 'commands:\n  test:\n    value: go test ./...\n' })
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  const d = await bash($, 'git commit -m "feat: x"')
  expect(text(d)).toContain('go test ./...')
  expect(text(d)).toContain('after the last green gate')
  await bash($, 'go test ./...')
  const ok = await bash($, 'git commit -m "feat: x"')
  expect(text(ok)).not.toContain('green gate')
})

for (const tool of ['Edit', 'Write']) {
  test(`3 D5 ${tool} under .harness/ does not dirty the gate`, async ($, on) => {
    world(on)
    await clean($)
    await $.tool.call({ tool, file_path: '.harness/tasks/X/TASK.md', content: 'x' } as any)
    await $.tool.call({ tool, file_path: '/proj/.harness/tasks/X/TASK.md', content: 'x' } as any)
    const r = await bash($, 'git commit -m "feat: x"')
    expect(text(r)).not.toContain('green gate')
  })
}

test('4a failed Edit does not dirty the gate', async ($, on) => {
  world(on, { editFails: true })
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  expect(text(await bash($, 'git commit -m "feat: x"'))).not.toContain('green gate')
})

test('4b Edit denied by runtime-guard (subagent on MEMORY.md) does not dirty the gate', async ($, on) => {
  world(on)
  await clean($)
  const d: any = await $.tool.call({ tool: 'Edit', file_path: '.harness/MEMORY.md', agentId: 'sub' } as any)
  expect(text(d)).toContain('written only by the Coordinator')
  expect(text(await bash($, 'git commit -m "feat: x"'))).not.toContain('green gate')
})

test('5 D2 deny text follows project.yaml language', async ($, on) => {
  world(on, { yaml: 'language: pt-br\n' })
  await clean($)
  const pt = await bash($, `git commit -m "x\n\n${AI}"`)
  expect(text(pt)).toContain('menciona IA')
})
test('5 D2 deny text is English without language', async ($, on) => {
  world(on, { yaml: 'name: x\n' })
  await clean($)
  const en = await bash($, `git commit -m "x\n\n${AI}"`)
  expect(text(en)).toContain('mentions AI')
})
test('5 D2 gate deny is Portuguese with pt-br', async ($, on) => {
  world(on, { yaml: 'language: pt-br\n' })
  await clean($)
  await $.tool.call({ tool: 'Write', file_path: 'src/a.go', content: 'x' } as any)
  expect(text(await bash($, 'git commit -m "feat: x"'))).toContain('último gate verde')
})

test('6 status line has ctx, English labels, and dh progress (via agent.spawn refresh)', async ($, on) => {
  const w = world(on, { tasks: { A: '| Status | concluída |', B: '| Status | bloqueada |', C: '| Status | open |' } })
  await clean($)
  await $.agent.spawn({ prompt: 'hi', description: 'd' } as any)
  const last = w.statuses[w.statuses.length - 1] ?? ''
  expect(last).toContain('ctx')
  expect(last).toContain('agents 1')
  expect(last).toContain('dh ')
  expect(last).toContain('1/3')
  expect(last).toContain('(1 blocked)')
})

test('6 status after real session.start event', async ($, on) => {
  const w = world(on, { tasks: { A: '| Status | done |' } })
  const start = ($ as any).session?.start
  if (typeof start !== 'function') throw new Error('NOT-RUN: engine $ cannot raise session.start')
  await start({ source: 'startup', cwd: CWD })
  expect(w.statuses.join('|')).toContain('ctx')
})

test('7 runtime-guard via composed module: spawn from subagent denied, main-session write allowed', async ($, on) => {
  world(on)
  const d: any = await $.agent.spawn({ prompt: 'x', parentAgentId: 'a1' } as any)
  expect(text(d)).toContain('only the Coordinator dispatches agents')
  const w: any = await $.tool.call({ tool: 'Write', file_path: '.harness/MEMORY.md', content: 'x' } as any)
  expect(w.deny).toBe(undefined)
  expect(text(w)).not.toContain('Coordinator')
})

test('8 fail-open: fs read failure does not block a commit or Bash', async ($, on) => {
  world(on, { fsFails: true })
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  const r = await bash($, 'git commit -m "feat: x"')
  // handler's readProject catches the fs failure; the call must still be answered, not thrown
  expect(r.isError === true && /boom/.test(text(r))).toBe(false)
  const ls = await bash($, 'ls')
  expect(ls.deny).toBe(undefined)
})

// Bug (0.18.1): the gate trap judged every commit by the session's edits, so a commit in another repository
// was denied after an edit here, and `git -C <dir> commit` escaped the trap altogether.
test('9a edit here, commit in another repo (cd) is not gated', async ($, on) => {
  world(on)
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  expect(text(await bash($, 'cd /other && git add -A && git commit -m "docs: x"'))).not.toContain('green gate')
  expect(text(await bash($, 'git commit -m "feat: x"'))).toContain('green gate')
  await clean($)
})

test('9b git -C: other repo passes, edited repo is gated', async ($, on) => {
  world(on)
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  expect(text(await bash($, 'git -C /other commit -m "docs: x"'))).not.toContain('green gate')
  expect(text(await bash($, 'git -C /proj commit -m "feat: x"'))).toContain('green gate')
  await clean($)
})

test('9c edit in another repo gates only that repo; its own gate clears it', async ($, on) => {
  world(on)
  await clean($)
  await $.tool.call({ tool: 'Write', file_path: '/other/x.go', content: 'x' } as any)
  expect(text(await bash($, 'git commit -m "feat: x"'))).not.toContain('green gate')
  expect(text(await bash($, 'cd /other && git commit -m "feat: x"'))).toContain('green gate')
  await bash($, 'go test ./...')
  expect(text(await bash($, 'cd /other && git commit -m "feat: x"'))).toContain('green gate')
  await bash($, 'cd /other && go test ./...')
  expect(text(await bash($, 'cd /other && git commit -m "feat: x"'))).not.toContain('green gate')
})

test('9d commit from a subdirectory resolves the repo root through git', async ($, on) => {
  world(on, { roots: { '/proj/sub': '/proj', '/proj': '/proj' } })
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: '/proj/a.go' } as any)
  expect(text(await bash($, 'cd sub && git commit -m "feat: x"'))).toContain('green gate')
  await clean($)
})

test('9e a gate followed by cd clears the gate\'s repo, not the later one', async ($, on) => {
  world(on)
  await clean($)
  await $.tool.call({ tool: 'Edit', file_path: 'src/a.go' } as any)
  await bash($, 'go test ./... && cd /other')
  expect(text(await bash($, 'git commit -m "feat: x"'))).not.toContain('green gate')
})
