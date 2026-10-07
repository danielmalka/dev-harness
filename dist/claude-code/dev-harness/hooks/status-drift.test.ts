import { expect, test } from 'claude-code/testing'
import { taskStatus } from './panel'
import { driftMessage, findDrift, prdLink, statusValue } from './status-drift'
import { STATUS_CASES } from './status-cases.fixture'

// The same table internal/dashboard/dashboard_test.go (TestStatusCasesFixture) reads.
for (const c of STATUS_CASES) {
  test(`fixture status: ${c.name}`, () => expect(taskStatus(c.md)).toBe(c.want))
}

// Same cases as TestPRDLink in internal/dashboard/dashboard_test.go.
const LINKS: Record<string, string> = {
  '| Story / PRD | PRD-003 (RF-1) |': 'PRD-003',
  '| PRD (RF-<n>) | PRD-011 (R1, R2) |': 'PRD-011',
  '| PRD (RF-<n>) | fora do PRD-008 (kit only) |': '',
  '| Story / PRD | nenhum |': '',
  '| Story / PRD | story X |': '',
  '| Tipo | PRD-001 |': '',
}
for (const [md, want] of Object.entries(LINKS)) {
  test(`fixture PRD link: ${md}`, () => expect(prdLink(md + '\n')).toBe(want))
}

test('PRD status reads table then bold, lowercased', () => {
  expect(statusValue('| Status | Entregue em 2026-01-01 |\n')).toBe('entregue em 2026-01-01')
  expect(statusValue('**Status:** rascunho\n')).toBe('rascunho')
})

const T = (status: string, prd: string) => `| Status | ${status} |\n| PRD (RF-<n>) | ${prd} |\n`

test('R19 conditions (a) open ticket in a delivered PRD and (b) all done in an undelivered PRD', () => {
  const prds = new Map([['PRD-001', 'entregue em 2026-01-01'], ['PRD-002', 'aprovado'], ['PRD-003', 'aprovado']])
  const d = findDrift([
    { id: 'T-1', md: T('pronta', 'PRD-001') },
    { id: 'T-2', md: T('concluída', 'PRD-002') },
    { id: 'T-3', md: T('concluída', 'PRD-002') },
    { id: 'T-4', md: T('pronta', 'PRD-003') },
    { id: 'T-5', md: T('pronta', 'fora do PRD-003') },
    { id: 'T-6', md: T('pronta', 'PRD-099') },
  ], prds, taskStatus)
  expect(d).toEqual([
    { prd: 'PRD-001', kind: 'open-in-delivered', tickets: ['T-1'] },
    { prd: 'PRD-002', kind: 'all-done', tickets: [] },
  ])
})

test('clean state has no drift', () => {
  const prds = new Map([['PRD-001', 'entregue em x'], ['PRD-002', 'aprovado']])
  expect(findDrift([{ id: 'T-1', md: T('concluída', 'PRD-001') }, { id: 'T-2', md: T('pronta', 'PRD-002') }], prds, taskStatus)).toEqual([])
})

test('message lists tickets in both languages', () => {
  const d = [{ prd: 'PRD-001', kind: 'open-in-delivered' as const, tickets: ['T-1', 'T-9'] }]
  expect(driftMessage(d, 'en')).toContain('T-1, T-9')
  expect(driftMessage(d, 'en')).toContain('commit again')
  expect(driftMessage(d, 'pt-br')).toContain('T-1, T-9')
  expect(driftMessage(d, 'pt-br')).toContain('commite de novo')
})
