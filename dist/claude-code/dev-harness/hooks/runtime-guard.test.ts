import { test, expect } from 'claude-code/testing'
import { guardFile, homeFrom, protectedHit, register, warnThenNext } from './runtime-guard.ts'
import { warnWrite } from './lib/warn.ts'

// Capture the plugin's handlers with a fake `on`, so agentId / parentAgentId (which the
// engine's $.tool.call and $.agent.spawn drop) can be set on the event.
type H = (...a: any[]) => any
const handlers: Record<string, H[]> = {}
;(register as any)((ev: string, ...rest: any[]) => {
  ;(handlers[ev] ||= []).push(rest[rest.length - 1])
})
const NEXT = 'passed'
const next = () => NEXT
const spawn = (e: any) => handlers['agent.spawn'][0]({}, e, next)
// A fake $ with env and cwd; `{}` (no env) is the fail-open case: only the suffix pattern holds.
const fake = (env: Record<string, string> = {}, cwd = '/proj') => ({ env: { get: async (k: string) => env[k] }, session: { cwd: async () => cwd } })
const file = (tool: string, e: any, $: any = {}) => Promise.all(handlers['tool.call'].map(h => h($, { tool, ...e }, next)))
const FILES = ['MEMORY.md', 'EPOCHAL.md', 'RISKS.md']

test('registers agent.spawn and Write+Edit tool.call hooks', () => {
  expect(handlers['agent.spawn'].length).toBe(1)
  expect(handlers['tool.call'].length).toBe(2)
})

// R3
test('R3 spawn with parentAgentId is denied', () => {
  const r = spawn({ parentAgentId: 'a1' })
  expect(r.deny).toContain('only the Coordinator dispatches agents')
})
test('R3 spawn without parentAgentId passes', () => {
  expect(spawn({})).toBe(NEXT)
})

// R4: each file x path style x tool, in a subagent
for (const tool of ['Write', 'Edit']) {
  for (const f of FILES) {
    const rel = `.harness/${f}`
    for (const [style, p] of [
      ['relative', rel],
      ['absolute', `/home/u/proj/${rel}`],
      ['backslash', `C:\\proj\\.harness\\${f}`],
      ['out-of-project', `/tmp/x/${rel}`],
    ]) {
      test(`R4 ${tool} ${style} ${f} with agentId is denied`, async () => {
        for (const r of await file(tool, { file_path: p, agentId: 'sub1' })) expect(r.deny).toContain(rel)
      })
      test(`R4 ${tool} ${style} ${f} without agentId passes`, async () => {
        for (const r of await file(tool, { file_path: p })) expect(r).toBe(NEXT)
      })
    }
  }
  test(`R4 ${tool} .harness/project.yaml with agentId passes`, async () => {
    for (const r of await file(tool, { file_path: '.harness/project.yaml', agentId: 's' })) expect(r).toBe(NEXT)
  })
  test(`R4 ${tool} lookalike name (not-MEMORY.md) with agentId passes`, async () => {
    for (const r of await file(tool, { file_path: 'docs/.harness/MEMORY.md.bak', agentId: 's' })) expect(r).toBe(NEXT)
    for (const r of await file(tool, { file_path: 'xharness/MEMORY.md', agentId: 's' })) expect(r).toBe(NEXT)
  })
  // R5: malformed input fails open
  test(`R5 ${tool} malformed input does not deny`, async () => {
    for (const e of [{}, { agentId: 's' }, { file_path: 42, agentId: 's' }, { file_path: null, agentId: 's' }, { file_path: { a: 1 }, agentId: 's' }, { file_path: false, agentId: 's' }])
      for (const r of await file(tool, e)) expect(r).toBe(NEXT)
  })
}
// PRD-014 R8: global mode, <home>/projects/<name>/<file>, with <home> from DH_HOME, else HOME/.harness, else USERPROFILE/.harness.
for (const tool of ['Write', 'Edit']) {
  for (const f of FILES) {
    test(`R8 ${tool} <home>/projects/x/${f} (DH_HOME) with agentId is denied`, async () => {
      for (const r of await file(tool, { file_path: `/h/projects/x/${f}`, agentId: 's' }, fake({ DH_HOME: '/h', HOME: '/u' }))) expect(r.deny).toContain(`/h/projects/x/${f}`)
    })
    test(`R8 ${tool} ~/.harness/projects/x/${f} (HOME) with agentId is denied`, async () => {
      for (const r of await file(tool, { file_path: `/u/.harness/projects/x/${f}`, agentId: 's' }, fake({ HOME: '/u' }))) expect(r.deny).toContain('written only by the Coordinator')
    })
    test(`R8 ${tool} <home>/projects/x/${f} by the Coordinator passes`, async () => {
      for (const r of await file(tool, { file_path: `/h/projects/x/${f}` }, fake({ DH_HOME: '/h' }))) expect(r).toBe(NEXT)
    })
    test(`R8 ${tool} <repo>/.harness/${f} with agentId is still denied with DH_HOME set`, async () => {
      for (const r of await file(tool, { file_path: `/proj/.harness/${f}`, agentId: 's' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain(`.harness/${f}`)
    })
  }
  test(`R8 ${tool} other <home> files and lookalikes with agentId pass`, async () => {
    for (const p of ['/h/projects/x/notes/MEMORY.md', '/h/config.yaml', '/h/sessions/x.json', '/h/projects/MEMORY.md', '/h/projects/x/project.yaml', '/h/projects/x/tasks/T-1/TASK.md', '/tmp/MEMORY.md', '/hx/projects/x/MEMORY.md', '/h/projects/x/MEMORY.md.bak'])
      for (const r of await file(tool, { file_path: p, agentId: 's' }, fake({ DH_HOME: '/h' }))) expect(r).toBe(NEXT)
  })
  test(`R8 ${tool} relative path and .. are normalized against the session cwd`, async () => {
    for (const r of await file(tool, { file_path: 'x/MEMORY.md', agentId: 's' }, fake({ DH_HOME: '/h' }, '/h/projects'))) expect(r.deny).toContain('/h/projects/x/MEMORY.md')
    for (const r of await file(tool, { file_path: '/h/projects/x/notes/../RISKS.md', agentId: 's' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain('/h/projects/x/RISKS.md')
    for (const r of await file(tool, { file_path: '/p/.harness/tasks/../MEMORY.md', agentId: 's' }, fake({}))) expect(r.deny).toContain('.harness/MEMORY.md')
  })
  test(`R8 ${tool} without any home variable only the suffix pattern holds (fail open)`, async () => {
    for (const r of await file(tool, { file_path: '/h/projects/x/MEMORY.md', agentId: 's' }, fake({}))) expect(r).toBe(NEXT)
    for (const r of await file(tool, { file_path: '/h/projects/x/MEMORY.md', agentId: 's' }, {})) expect(r).toBe(NEXT)
  })
}
for (const tool of ['Write', 'Edit']) {
  test(`R8 ${tool} leading ~ expands to the user home (HOME, else USERPROFILE)`, async () => {
    for (const p of ['~/.harness/projects/x/MEMORY.md', '~/.harness/MEMORY.md', '~\\.harness\\projects\\x\\RISKS.md'])
      for (const r of await file(tool, { file_path: p, agentId: 's' }, fake({ HOME: '/u' }, '/elsewhere'))) expect(r.deny).toContain('written only by the Coordinator')
    for (const r of await file(tool, { file_path: '~/.harness/projects/x/MEMORY.md', agentId: 's' }, fake({ USERPROFILE: '/w' }, '/elsewhere'))) expect(r.deny).toContain('/w/.harness/projects/x/MEMORY.md')
    for (const r of await file(tool, { file_path: '~/.harness/projects/x/notes/MEMORY.md', agentId: 's' }, fake({ HOME: '/u' }))) expect(r).toBe(NEXT)
  })
  test(`R8 ${tool} both patterns are case-insensitive (a look-alike denied is the safe side)`, async () => {
    for (const p of ['/p/.HARNESS/memory.md', '.Harness/Risks.MD', '/h/Projects/x/Memory.md', '/H/projects/x/EPOCHAL.md'])
      for (const r of await file(tool, { file_path: p, agentId: 's' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain('written only by the Coordinator')
  })
}
test('R8 protectedHit without a user home leaves ~ literal (relative to cwd)', () => {
  expect(protectedHit('/h', '~/x/MEMORY.md', '/w')).toBe(undefined)
  expect(protectedHit('/u/.harness', '~/.harness/projects/x/MEMORY.md', '/w', '/u')).toBe('/u/.harness/projects/x/MEMORY.md')
})
test('R8 homeFrom: DH_HOME, then HOME/.harness, then USERPROFILE/.harness, then empty', () => {
  expect(homeFrom('/h', '/u', 'C:\\Users\\u')).toBe('/h')
  expect(homeFrom(undefined, '/u', 'C:\\Users\\u')).toBe('/u/.harness')
  expect(homeFrom(undefined, undefined, 'C:\\Users\\u')).toBe('/C:/Users/u/.harness')
  expect(homeFrom(undefined, undefined, undefined)).toBe('')
  expect(homeFrom('rel/h', undefined, undefined, '/w')).toBe('/w/rel/h')
})
test('R8 Windows drive paths meet whatever the slash or case of the drive', () => {
  const home = homeFrom('c:\\dev\\.harness', undefined, undefined)
  expect(protectedHit(home, 'C:\\dev\\.harness\\projects\\x\\MEMORY.md')).toBe('/C:/dev/.harness/projects/x/MEMORY.md')
  expect(protectedHit(home, 'C:/dev/.harness/projects/x/notes/MEMORY.md')).toBe(undefined)
})
test('R8 guardFile(home, e, next) is the panel signature', () => {
  expect(guardFile('/h', { file_path: '/h/projects/x/MEMORY.md', agentId: 's' }, next).deny).toContain('Coordinator')
  expect(guardFile('', { file_path: '/h/projects/x/MEMORY.md', agentId: 's' }, next)).toBe(NEXT)
})

test('R5 spawn with malformed event does not throw', () => {
  expect(spawn({})).toBe(NEXT)
  expect(handlers['agent.spawn'][0]({}, undefined, next)).toBe(NEXT)
})

// Through the engine: plugin hooks sit above the test's hooks; no agentId at the main session.
test('engine: main-session Write to .harness/MEMORY.md reaches the tool (not denied)', async ($, on) => {
  on('tool.call', { tool: 'Write' }, () => ({ result: 'ok' }) as any)
  const r: any = await $.tool.call({ tool: 'Write', file_path: '.harness/MEMORY.md', content: 'x' } as any)
  expect(r.deny).toBe(undefined)
})
// Observed: the engine keeps agentId on $.tool.call (the types say it is dropped), so the
// full chain plugin hook -> deny is reachable from a test.
for (const tool of ['Write', 'Edit']) {
  test(`engine: ${tool} on .harness/RISKS.md with agentId is denied end to end`, async ($, on) => {
    on('tool.call', { tool }, () => ({ result: 'ok' }) as any)
    const r: any = await $.tool.call({ tool, file_path: '.harness/RISKS.md', agentId: 'sub' } as any)
    expect(r.text ?? r.deny).toContain('written only by the Coordinator')
  })
}
test('engine: spawn from the test (no parent) is not denied', async ($, on) => {
  on('agent.spawn', () => ({ agentId: 'x', model: 'm' }) as any)
  const r: any = await $.agent.spawn({ prompt: 'hi' })
  expect(r.deny).toBe(undefined)
})

// PRD-017 R3: pdocs/ is writable by a subagent; the records next to it stay protected.
test('R3 pdocs: subagent may write under projects/x/pdocs/, records outside it stay denied', () => {
  const h = '/h'
  expect(protectedHit(h, '/h/projects/x/pdocs/index.html')).toBe(undefined)
  expect(protectedHit(h, '/h/projects/x/pdocs/d.svg')).toBe(undefined)
  expect(protectedHit(h, '/h/projects/x/pdocs/MEMORY.md')).toBe(undefined)
  expect(protectedHit(h, '~/.harness/projects/x/pdocs/EPOCHAL.md', '/', '/u')).toBe(undefined)
  expect(guardFile(h, { file_path: '/h/projects/x/pdocs/index.html', agentId: 's' }, next)).toBe(NEXT)
  expect(guardFile(h, { file_path: '/h/projects/x/MEMORY.md', agentId: 's' }, next).deny).toContain('Coordinator')
})

// PRD-025 R5-R7: warn-only rule for the main session.
const L = (n: number) => Array.from({ length: n }, (_, i) => `l${i}`).join('\n')
const W = (e: any, ex: any = () => true, home = '/h', cwd = '/proj') => warnWrite({ tool: 'Edit', ...e }, home, cwd, ex)
test('R5 warn: Edit with 6 changed lines warns, 5 does not (larger side counts)', async () => {
  expect(await W({ file_path: 'src/a.ts', old_string: L(6), new_string: 'x' })).toContain('6-line Edit')
  expect(await W({ file_path: 'src/a.ts', old_string: 'x', new_string: L(6) })).toContain('src/a.ts')
  expect(await W({ file_path: 'src/a.ts', old_string: L(5), new_string: L(5) })).toBe(undefined)
  expect(await W({ file_path: 'src/a.ts', old_string: L(5) + '\n', new_string: 'x' })).toBe(undefined) // trailing newline
})
test('R5 warn: Write content over 5 lines warns; new file warns; small existing file does not', async () => {
  const w = (e: any, ex: any) => warnWrite({ tool: 'Write', ...e }, '/h', '/proj', ex)
  expect(await w({ file_path: 'a.ts', content: L(6) }, () => true)).toContain('6-line Write')
  expect(await w({ file_path: 'a.ts', content: 'x' }, () => false)).toContain('new file')
  expect(await w({ file_path: 'a.ts', content: L(5) }, async () => true)).toBe(undefined)
})
test('R5 warn: exists unknown or throwing fails open', async () => {
  const w = (ex: any) => warnWrite({ tool: 'Write', file_path: 'a.ts', content: L(9) }, '/h', '/proj', ex)
  expect(await w(() => undefined)).toBe(undefined)
  expect(await w(() => { throw new Error('x') })).toBe(undefined)
  expect(await w(async () => { throw new Error('x') })).toBe(undefined)
})
test('R7 warn: agentId, harness paths (repo and global) and Bash do not warn', async () => {
  const big = { old_string: L(9), new_string: L(9), content: L(9) }
  expect(await W({ file_path: 'src/a.ts', agentId: 's', ...big })).toBe(undefined)
  for (const p of ['.harness/project.yaml', '.harness/local.yaml', '.harness/tasks/T-1/TASK.md', '.harness/MEMORY.md', '/proj/.harness/RISKS.md', '/h/projects/x/project.yaml', '/h/projects/x/tasks/T-1/TASK.md', '/h/projects/x/MEMORY.md'])
    for (const tool of ['Edit', 'Write']) expect(await W({ tool, file_path: p, ...big }, () => false)).toBe(undefined)
  expect(await W({ tool: 'Bash', file_path: 'src/a.ts', command: 'x', ...big }, () => false)).toBe(undefined)
  expect(handlers['tool.call'].length).toBe(2) // no Bash hook registered
})
test('R7 warn: only this project harness dir is exempt, a foreign .harness warns', async () => {
  const big = { old_string: L(9), new_string: L(9), content: L(9) }
  for (const tool of ['Edit', 'Write']) {
    expect(await W({ tool, file_path: '/tmp/other/.harness/payload.ts', ...big }, () => true)).toContain('/tmp/other/.harness/payload.ts')
    expect(await W({ tool, file_path: '/proj/.harness/x', ...big }, () => true)).toBe(undefined)
    expect(await W({ tool, file_path: '/h/projects/p/x', ...big }, () => true)).toBe(undefined)
  }
})
test('R7 warn: a .harness dir equal to <home> is not exempt, so <home>/config.yaml warns even when home is under an ancestor of cwd', async () => {
  const e = { file_path: '/home/u/.harness/config.yaml', old_string: L(9), new_string: L(9) }
  expect(await W(e, () => true, '/home/u/.harness', '/home/u/proj')).toContain('config.yaml')
})
test('R6 warn: malformed input never throws or warns', async () => {
  for (const e of [undefined, null, 42, {}, { tool: 'Edit' }, { tool: 'Edit', file_path: 42 }, { tool: 'Write', file_path: 'a', content: 5 }, { tool: 'Edit', file_path: 'a', old_string: 5, new_string: {} }])
    expect(await warnWrite(e, '/h', '/proj', () => true)).toBe(undefined)
})
test('R6 warnThenNext: always calls next(e), never denies, toasts only when warned', async () => {
  const toasts: string[] = []
  const $ = { ui: { toast: (t: string) => toasts.push(t) }, fs: { exists: async () => true } }
  const e1 = { tool: 'Edit', file_path: 'src/a.ts', old_string: L(7), new_string: 'x' }
  expect(await warnThenNext($, e1, '/h', '/proj', next)).toBe(NEXT)
  expect(toasts.length).toBe(1)
  expect(toasts[0]).toContain('warning only')
  expect(await warnThenNext($, { ...e1, old_string: 'x' }, '/h', '/proj', next)).toBe(NEXT)
  expect(await warnThenNext($, { ...e1, agentId: 's' }, '/h', '/proj', next)).toBe(NEXT)
  expect(toasts.length).toBe(1)
  for (const bad of [undefined, {}, { tool: 'Edit', file_path: 7 }]) expect(await warnThenNext($, bad, '/h', '/proj', next)).toBe(NEXT)
  // broken engine: no fs / throwing toast still passes
  expect(await warnThenNext({}, e1, '/h', '/proj', next)).toBe(NEXT)
  expect(await warnThenNext({ ui: { toast: () => { throw new Error('x') } }, fs: { exists: () => true } }, e1, '/h', '/proj', next)).toBe(NEXT)
})
test('R6 warn: through the registered handlers a warned main-session Edit still passes, a subagent record write is still denied', async () => {
  const $: any = { ...fake({ DH_HOME: '/h' }), ui: { toast: () => {} }, fs: { exists: async () => true } }
  const big = { file_path: 'src/a.ts', old_string: L(9), new_string: 'x' }
  for (const r of await file('Edit', big, $)) expect(r).toBe(NEXT)
  for (const r of await file('Edit', { file_path: '.harness/MEMORY.md', agentId: 's' }, $)) expect(r.deny).toContain('.harness/MEMORY.md')
})
