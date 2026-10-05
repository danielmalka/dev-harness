import { test, expect } from 'claude-code/testing'
import { register } from './runtime-guard.ts'

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
const file = (tool: string, e: any) => handlers['tool.call'].map(h => h({}, { tool, ...e }, next))
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
      test(`R4 ${tool} ${style} ${f} with agentId is denied`, () => {
        for (const r of file(tool, { file_path: p, agentId: 'sub1' })) expect(r.deny).toContain(rel)
      })
      test(`R4 ${tool} ${style} ${f} without agentId passes`, () => {
        for (const r of file(tool, { file_path: p })) expect(r).toBe(NEXT)
      })
    }
  }
  test(`R4 ${tool} .harness/project.yaml with agentId passes`, () => {
    for (const r of file(tool, { file_path: '.harness/project.yaml', agentId: 's' })) expect(r).toBe(NEXT)
  })
  test(`R4 ${tool} lookalike name (not-MEMORY.md) with agentId passes`, () => {
    for (const r of file(tool, { file_path: 'docs/.harness/MEMORY.md.bak', agentId: 's' })) expect(r).toBe(NEXT)
    for (const r of file(tool, { file_path: 'xharness/MEMORY.md', agentId: 's' })) expect(r).toBe(NEXT)
  })
  // R5: malformed input fails open
  test(`R5 ${tool} malformed input does not deny`, () => {
    for (const e of [{}, { agentId: 's' }, { file_path: 42, agentId: 's' }, { file_path: null, agentId: 's' }, { file_path: { a: 1 }, agentId: 's' }, { file_path: false, agentId: 's' }])
      for (const r of file(tool, e)) expect(r).toBe(NEXT)
  })
}
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
