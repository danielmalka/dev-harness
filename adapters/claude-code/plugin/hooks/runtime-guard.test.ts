import { test, expect } from 'claude-code/testing'
import { denyCuratorShell, forceCuratorModel, guardFile, homeFrom, isCuratorSpawn, noteSpawn, protectedHit, register, strongModel, warnThenNext } from './runtime-guard.ts'
import { warnWrite } from './lib/warn.ts'

// Capture the plugin's handlers with a fake `on`, so agentId / parentAgentId (which the
// engine's $.tool.call and $.agent.spawn drop) can be set on the event.
type H = (...a: any[]) => any
const handlers: Record<string, H[]> = {}
const matchers: Record<string, (string | undefined)[]> = {} // the `{ tool }` each handler was registered for, undefined when it has none
;(register as any)((ev: string, ...rest: any[]) => {
  ;(handlers[ev] ||= []).push(rest[rest.length - 1])
  ;(matchers[ev] ||= []).push(rest.length > 1 ? rest[0]?.tool : undefined)
})
const NEXT = 'passed'
const next = () => NEXT
const spawn = (e: any) => handlers['agent.spawn'][0]({}, e, next)
// A fake $ with env and cwd; `{}` (no env) is the fail-open case: only the suffix pattern holds.
const fake = (env: Record<string, string> = {}, cwd = '/proj') => ({ env: { get: async (k: string) => env[k] }, session: { cwd: async () => cwd } })
// Like the engine, a handler registered for a tool only sees the calls of that tool.
const file = (tool: string, e: any, $: any = {}) => Promise.all(handlers['tool.call'].filter((_h, i) => matchers['tool.call'][i] === undefined || matchers['tool.call'][i] === tool).map(h => h($, { tool, ...e }, next)))
const FILES = ['MEMORY.md', 'EPOCHAL.md', 'RISKS.md', 'FACTS.md']

test('registers agent.spawn, Write+Edit and the curator shell-deny tool.call hooks', () => {
  expect(handlers['agent.spawn'].length).toBe(1)
  expect(matchers['tool.call']).toEqual(['Write', 'Edit', 'Bash', 'PowerShell', 'Monitor'])
})

// R3
test('R3 spawn with parentAgentId is denied', () => {
  const r = spawn({ parentAgentId: 'a1' })
  expect(r.deny).toContain('só o Coordinator despacha agentes')
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
      for (const r of await file(tool, { file_path: `/u/.harness/projects/x/${f}`, agentId: 's' }, fake({ HOME: '/u' }))) expect(r.deny).toContain('só é escrito pelo Coordinator')
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
      for (const r of await file(tool, { file_path: p, agentId: 's' }, fake({ HOME: '/u' }, '/elsewhere'))) expect(r.deny).toContain('só é escrito pelo Coordinator')
    for (const r of await file(tool, { file_path: '~/.harness/projects/x/MEMORY.md', agentId: 's' }, fake({ USERPROFILE: '/w' }, '/elsewhere'))) expect(r.deny).toContain('/w/.harness/projects/x/MEMORY.md')
    for (const r of await file(tool, { file_path: '~/.harness/projects/x/notes/MEMORY.md', agentId: 's' }, fake({ HOME: '/u' }))) expect(r).toBe(NEXT)
  })
  test(`R8 ${tool} both patterns are case-insensitive (a look-alike denied is the safe side)`, async () => {
    for (const p of ['/p/.HARNESS/memory.md', '.Harness/Risks.MD', '/h/Projects/x/Memory.md', '/H/projects/x/EPOCHAL.md'])
      for (const r of await file(tool, { file_path: p, agentId: 's' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain('só é escrito pelo Coordinator')
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
    expect(r.text ?? r.deny).toContain('só é escrito pelo Coordinator')
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
  expect(await W({ file_path: 'src/a.ts', old_string: L(6), new_string: 'x' })).toContain('Edit de 6 linhas')
  expect(await W({ file_path: 'src/a.ts', old_string: 'x', new_string: L(6) })).toContain('src/a.ts')
  expect(await W({ file_path: 'src/a.ts', old_string: L(5), new_string: L(5) })).toBe(undefined)
  expect(await W({ file_path: 'src/a.ts', old_string: L(5) + '\n', new_string: 'x' })).toBe(undefined) // trailing newline
})
test('R5 warn: Write content over 5 lines warns; new file warns; small existing file does not', async () => {
  const w = (e: any, ex: any) => warnWrite({ tool: 'Write', ...e }, '/h', '/proj', ex)
  expect(await w({ file_path: 'a.ts', content: L(6) }, () => true)).toContain('Write de 6 linhas')
  expect(await w({ file_path: 'a.ts', content: 'x' }, () => false)).toContain('arquivo novo')
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
  expect(matchers['tool.call'].slice(0, 2)).toEqual(['Write', 'Edit']) // the warn path only hooks Write and Edit; the shell hooks only deny the curator
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
  expect(toasts[0]).toContain('só aviso')
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

// PRD-035 R12 (D7): FACTS.md is a protected record like the other three; the only subagent allowed to write it is the curator, a `docs-guide`
// dispatched with FACTS.md in its write set. The matrices above already run FACTS.md through the repo and global patterns, both tools and every path style.
const CURATOR_PROMPT = 'Objetivo: curar os fatos.\nArquivos permitidos: .harness/FACTS.md\nModelo: sonnet'

test('R12 isCuratorSpawn: docs-guide with FACTS.md in the write set, and nothing else', () => {
  expect(isCuratorSpawn('dh:docs-guide', CURATOR_PROMPT)).toBe(true)
  expect(isCuratorSpawn('docs-guide', CURATOR_PROMPT)).toBe(true)
  expect(isCuratorSpawn('dh:docs-guide', '- Conjunto de escrita: <harness dir>/FACTS.md')).toBe(true)
  expect(isCuratorSpawn('dh:docs-guide', 'Write set: /h/projects/x/FACTS.md')).toBe(true)
  expect(isCuratorSpawn('dh:docs-guide', '> arquivos permitidos : FACTS.md, .harness/MEMORY.md')).toBe(true)
  // another role, even with the same text
  for (const t of ['dh:builder', 'dh:qa-verifier', 'dh:repo-scout', 'dh:reviewer', 'builder', 'general-purpose', undefined, 7]) expect(isCuratorSpawn(t, CURATOR_PROMPT)).toBe(false)
  // the docs-guide without FACTS.md in the write set, or with it only in other lines of the text
  for (const p of ['Arquivos permitidos: .harness/MEMORY.md', 'Arquivos permitidos: docs/x.md\nLeia o FACTS.md antes', 'Resumo: atualizar FACTS.md', 'Arquivos permitidos: .harness/FACTS.md.bak', 'Arquivos permitidos: .harness/NOFACTS.md', '', undefined, null, 5]) expect(isCuratorSpawn('dh:docs-guide', p)).toBe(false)
})

for (const tool of ['Write', 'Edit']) {
  test(`R12 ${tool}: the curator writes FACTS.md (repo and global), a plain docs-guide and a builder do not`, async () => {
    noteSpawn({ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }, 'cur-1')
    noteSpawn({ subagentType: 'dh:docs-guide', prompt: 'Arquivos permitidos: docs/x.md' }, 'plain-1')
    noteSpawn({ subagentType: 'dh:builder', prompt: CURATOR_PROMPT }, 'build-1')
    for (const [path, env] of [['.harness/FACTS.md', {}], ['/proj/.harness/FACTS.md', {}], ['/h/projects/x/FACTS.md', { DH_HOME: '/h' }], ['~/.harness/projects/x/FACTS.md', { HOME: '/u' }]] as const) {
      for (const r of await file(tool, { file_path: path, agentId: 'cur-1' }, fake(env as any))) expect(r).toBe(NEXT)
      for (const id of ['plain-1', 'build-1', 'unknown-1'])
        for (const r of await file(tool, { file_path: path, agentId: id }, fake(env as any))) expect(r.deny).toContain('só é escrito pelo Coordinator')
    }
  })
  test(`R12 ${tool}: the curator exception covers FACTS.md only, not the other records`, async () => {
    noteSpawn({ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }, 'cur-2')
    for (const f of ['MEMORY.md', 'EPOCHAL.md', 'RISKS.md'])
      for (const r of await file(tool, { file_path: `.harness/${f}`, agentId: 'cur-2' })) expect(r.deny).toContain(`.harness/${f}`)
    for (const r of await file(tool, { file_path: '/h/projects/x/MEMORY.md', agentId: 'cur-2' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain('/h/projects/x/MEMORY.md')
  })
}

test('R12 an id reused by a non-curator loses the exception; malformed spawn input is ignored', async () => {
  noteSpawn({ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }, 'cur-3')
  for (const r of await file('Write', { file_path: '.harness/FACTS.md', agentId: 'cur-3' })) expect(r).toBe(NEXT)
  noteSpawn({ subagentType: 'dh:builder', prompt: 'x' }, 'cur-3')
  for (const r of await file('Write', { file_path: '.harness/FACTS.md', agentId: 'cur-3' })) expect(r.deny).toContain('.harness/FACTS.md')
  for (const bad of [[undefined, undefined], [{}, undefined], [{}, ''], [{}, 7], [null, 'x'], [{ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT }, null]] as const) noteSpawn(bad[0] as any, bad[1] as any)
})

test('R12 through the registered spawn handler: the engine answer carries the agent id, sync or promise, and the result is returned unchanged', async () => {
  const spawnAs = (e: any, answer: any) => handlers['agent.spawn'][0]({}, e, () => answer)
  const ev = { subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }
  const sync = { agentId: 'cur-4', model: 'claude-sonnet-5-5' }
  expect(spawnAs(ev, sync)).toBe(sync)
  for (const r of await file('Edit', { file_path: '.harness/FACTS.md', agentId: 'cur-4' })) expect(r).toBe(NEXT)
  const pending = Promise.resolve({ agentId: 'cur-5', model: 'claude-sonnet-5-5' })
  expect(spawnAs(ev, pending)).toBe(pending)
  await pending
  for (const r of await file('Edit', { file_path: '.harness/FACTS.md', agentId: 'cur-5' })) expect(r).toBe(NEXT)
  // a denied spawn (from a subagent) never reaches the engine and notes nothing
  expect(spawnAs({ ...ev, parentAgentId: 'a1' }, { agentId: 'cur-6' }).deny).toContain('só o Coordinator despacha')
  for (const r of await file('Edit', { file_path: '.harness/FACTS.md', agentId: 'cur-6' })) expect(r.deny).toContain('.harness/FACTS.md')
  // a rejected promise is not an unhandled rejection
  const rejected = Promise.reject(new Error('boom'))
  rejected.catch(() => undefined)
  expect(spawnAs(ev, rejected)).toBe(rejected)
})

// Through the engine, composed with the panel: the main session writes FACTS.md, a subagent is denied, and a subagent spawned as the curator is not.
test('engine: Write on .harness/FACTS.md with agentId is denied end to end', async ($, on) => {
  on('tool.call', { tool: 'Write' }, () => ({ result: 'ok' }) as any)
  const r: any = await $.tool.call({ tool: 'Write', file_path: '.harness/FACTS.md', agentId: 'sub-f' } as any)
  expect(r.text ?? r.deny).toContain('.harness/FACTS.md só é escrito pelo Coordinator')
  const m: any = await $.tool.call({ tool: 'Write', file_path: '.harness/FACTS.md', content: 'x' } as any)
  expect(m.deny).toBe(undefined)
})

// PRD-035 F3 contract with the Coordinator: the dispatch line that .agents/coordinator.md tells the Coordinator to write is the line this mod recognizes. A hooks
// module may import nothing but its own files and "claude-code", so this test cannot read coordinator.md; the Go test internal/kit/facts_curation_contract_test.go
// closes the chain: it checks that coordinator.md and THIS file both spell CURATOR_DISPATCH_LINE exactly, and this test checks that the mod accepts it.
const CURATOR_DISPATCH_LINE = 'Arquivos permitidos: <harness dir>/FACTS.md'

test('R12 F3 the curator dispatch line of coordinator.md is the one the mod recognizes', async () => {
  const line = CURATOR_DISPATCH_LINE.replace('<harness dir>', '/h/projects/proj')
  const dispatch = `Objetivo: curar o FACTS.md.\nsha do HEAD: 4fc0ac3\n${line}\nModelo: sonnet`
  expect(isCuratorSpawn('dh:docs-guide', dispatch)).toBe(true)
  expect(isCuratorSpawn('dh:docs-guide', CURATOR_DISPATCH_LINE)).toBe(true) // the placeholder form, as written in the agent file
  // the same line from any other role is not a curator
  for (const type of ['dh:builder', 'dh:qa-verifier', 'dh:repo-scout', 'dh:reviewer', 'dh:coordinator']) expect(isCuratorSpawn(type, dispatch)).toBe(false)
  // and the writing really flows: the id the engine answers with may write FACTS.md and nothing else
  noteSpawn({ subagentType: 'dh:docs-guide', prompt: dispatch, model: 'sonnet' }, 'cur-f3')
  for (const r of await file('Write', { file_path: '/h/projects/proj/FACTS.md', agentId: 'cur-f3' }, fake({ DH_HOME: '/h' }))) expect(r).toBe(NEXT)
  for (const r of await file('Write', { file_path: '/h/projects/proj/MEMORY.md', agentId: 'cur-f3' }, fake({ DH_HOME: '/h' }))) expect(r.deny).toContain('/h/projects/proj/MEMORY.md')
  // without the line the same docs-guide is not a curator
  expect(isCuratorSpawn('dh:docs-guide', dispatch.replace(line, 'Arquivos permitidos: docs/x.md'))).toBe(false)
})

// Verification of F3, model: the curation needs sonnet or above. The hook sees the Agent tool's `model` and may rewrite it, and the answer carries the model the
// agent runs on, so the mod forces `sonnet` on a curator dispatch that does not ask for it and recognizes the curator only when the model that came back is strong.
test('F3 strongModel: sonnet, opus and fable (alias or full id) are strong; haiku, empty and unknown are not', () => {
  for (const m of ['sonnet', 'claude-sonnet-5-5', 'Sonnet', 'opus', 'claude-opus-5-5', 'fable', 'claude-fable-5-2']) expect(strongModel(m)).toBe(true)
  for (const m of ['haiku', 'claude-haiku-5-5', '', 'inherit', 'm', undefined, null, 7, {}]) expect(strongModel(m)).toBe(false)
})

test('F3 forceCuratorModel: a curator dispatch without a strong model is rewritten to sonnet, everything else passes untouched', () => {
  const curator = { subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, description: 'curar', tool_use_id: 't' }
  expect(forceCuratorModel(curator)).toEqual({ ...curator, model: 'sonnet' }) // no model: the agent's own model (haiku) would run
  expect(forceCuratorModel({ ...curator, model: 'haiku' })).toEqual({ ...curator, model: 'sonnet' })
  expect(forceCuratorModel({ ...curator, model: 'claude-haiku-5-5' })).toEqual({ ...curator, model: 'sonnet' })
  expect(forceCuratorModel({ ...curator, model: 'inherit', parentModel: 'haiku' })).toEqual({ ...curator, model: 'sonnet', parentModel: 'haiku' })
  for (const model of ['sonnet', 'opus', 'claude-opus-5-5', 'fable']) {
    const e = { ...curator, model }
    expect(forceCuratorModel(e)).toBe(e) // asked for a strong model: the same object
  }
  const inherit = { ...curator, model: 'inherit', parentModel: 'claude-opus-5-5' }
  expect(forceCuratorModel(inherit)).toBe(inherit) // inherit under a strong parent is strong
  // a docs-guide without the line, any other role and malformed input are not touched, and no model is added
  const plain = { subagentType: 'dh:docs-guide', prompt: 'Arquivos permitidos: docs/x.md', model: 'haiku' }
  expect(forceCuratorModel(plain)).toBe(plain)
  const builder = { subagentType: 'dh:builder', prompt: CURATOR_PROMPT }
  expect(forceCuratorModel(builder)).toBe(builder)
  for (const bad of [undefined, null, {}, 7, 'x']) expect(forceCuratorModel(bad as any)).toBe(bad)
})

test('F3 through the registered spawn handler: the engine gets sonnet, and the curator it answers with may write FACTS.md', async () => {
  let seen: any
  const ev = { subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, description: 'curar', tool_use_id: 't' }
  const answer = handlers['agent.spawn'][0]({}, ev, (x: any) => { seen = x; return { agentId: 'cur-m1', model: 'claude-sonnet-5-5' } })
  expect(seen.model).toBe('sonnet')
  expect(answer).toEqual({ agentId: 'cur-m1', model: 'claude-sonnet-5-5' })
  for (const r of await file('Write', { file_path: '.harness/FACTS.md', agentId: 'cur-m1' })) expect(r).toBe(NEXT)
  // a plain docs-guide keeps whatever model it asked for
  let plainSeen: any
  handlers['agent.spawn'][0]({}, { ...ev, prompt: 'Arquivos permitidos: docs/x.md', model: 'haiku' }, (x: any) => { plainSeen = x; return { agentId: 'plain-m1', model: 'claude-haiku-5-5' } })
  expect(plainSeen.model).toBe('haiku')
})

test('F3 a curator-shaped agent that came back on a weak model is not the curator: the write is denied and the message names the model', async () => {
  // e.g. the engine ignored the rewrite, or a caller did not go through the hook
  const ev = { subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }
  noteSpawn(ev, 'weak-1', 'claude-haiku-5-5')
  noteSpawn({ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT }, 'weak-2') // no model anywhere
  for (const [id, shown] of [['weak-1', 'claude-haiku-5-5'], ['weak-2', 'o padrão do agente (haiku)']] as const) {
    for (const tool of ['Write', 'Edit']) {
      for (const r of await file(tool, { file_path: '.harness/FACTS.md', agentId: id })) {
        expect(r.deny).toContain('modelo sonnet ou mais forte')
        expect(r.deny).toContain(shown)
        expect(r.deny).toContain('redespacha a curadoria com o modelo sonnet')
      }
    }
  }
  // and a later spawn of the same id on a strong model clears it
  noteSpawn(ev, 'weak-1', 'claude-sonnet-5-5')
  for (const r of await file('Write', { file_path: '.harness/FACTS.md', agentId: 'weak-1' })) expect(r).toBe(NEXT)
})

// Verification of F3, shell: the recognized curator has no Bash, PowerShell or Monitor for as long as it lives; Read and Grep are not touched; a docs-guide without
// the line keeps its shell.
test('F3 denyCuratorShell: Bash, PowerShell and Monitor are denied to the curator with a message that points to Read and Grep', async () => {
  noteSpawn({ subagentType: 'dh:docs-guide', prompt: CURATOR_PROMPT, model: 'sonnet' }, 'cur-sh1')
  noteSpawn({ subagentType: 'dh:docs-guide', prompt: 'Arquivos permitidos: docs/x.md', model: 'sonnet' }, 'plain-sh1')
  noteSpawn({ subagentType: 'dh:builder', prompt: CURATOR_PROMPT, model: 'sonnet' }, 'build-sh1')
  for (const tool of ['Bash', 'PowerShell', 'Monitor']) {
    for (const r of await file(tool, { command: 'cat -n src/log.go', agentId: 'cur-sh1' })) {
      expect(r.deny).toContain('não roda comandos')
      expect(r.deny).toContain('Read e Grep')
    }
    for (const id of ['plain-sh1', 'build-sh1', 'unknown-sh1']) for (const r of await file(tool, { command: 'cat -n src/log.go', agentId: id })) expect(r).toBe(NEXT)
    for (const r of await file(tool, { command: 'git status' })) expect(r).toBe(NEXT) // the main session has no agentId
  }
  // Read and Grep are not shell tools: no handler intercepts them
  for (const tool of ['Read', 'Grep', 'Glob', 'Skill']) expect((await file(tool, { file_path: 'src/log.go', agentId: 'cur-sh1' })).length).toBe(0)
  // the curator still writes FACTS.md
  for (const r of await file('Write', { file_path: '.harness/FACTS.md', agentId: 'cur-sh1' })) expect(r).toBe(NEXT)
  expect(denyCuratorShell({ agentId: 'cur-sh1' })?.deny).toContain('Read e Grep')
  for (const bad of [undefined, null, {}, 7, { agentId: 7 }]) expect(denyCuratorShell(bad as any)).toBe(undefined)
  // a reused id loses it
  noteSpawn({ subagentType: 'dh:builder', prompt: 'x' }, 'cur-sh1')
  for (const r of await file('Bash', { command: 'ls', agentId: 'cur-sh1' })) expect(r).toBe(NEXT)
})
