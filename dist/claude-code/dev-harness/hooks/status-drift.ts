// PRD-012 R19/R20: ticket/PRD status drift. Pure, mirrors internal/dashboard + internal/kit/status_drift.go.
// Shared fixture: internal/dashboard/testdata/status-cases.json.
const STATUS_ROW = /^\|\s*Status\s*\|\s*([^|]+)\|/m
const STATUS_BOLD = /\*\*Status:?\*\*:?\s*([^\n]+)/
const PRD_FIELD = /^\|\s*(?:Story \/ PRD|PRD \(RF-[^)]*\))\s*\|\s*([^|]+)\|/m
const PRD_ID = /PRD-\d+/

export function statusValue(md: string): string {
  return (md.match(STATUS_ROW)?.[1] ?? md.match(STATUS_BOLD)?.[1] ?? '').trim().toLowerCase()
}

// The PRD id a ticket links to, or '' ("fora"/"nenhum" or no PRD-n).
export function prdLink(md: string): string {
  const v = md.match(PRD_FIELD)?.[1]
  if (v === undefined || /^(fora|nenhum)/.test(v.trim().toLowerCase())) return ''
  return v.match(PRD_ID)?.[0] ?? ''
}

export type Drift = { prd: string; kind: 'open-in-delivered' | 'all-done'; tickets: string[] }

// tickets: id + TASK.md text; prds: PRD id -> lowercased header Status. classify = taskStatus (panel.ts).
export function findDrift(tickets: { id: string; md: string }[], prds: Map<string, string>, classify: (md: string) => string): Drift[] {
  const by = new Map<string, { total: number; open: string[] }>()
  for (const t of tickets) {
    const id = prdLink(t.md)
    if (!id || !prds.has(id)) continue
    const c = by.get(id) ?? { total: 0, open: [] }
    c.total++
    if (classify(t.md) !== 'done') c.open.push(t.id)
    by.set(id, c)
  }
  const out: Drift[] = []
  for (const id of [...by.keys()].sort()) {
    const c = by.get(id)!
    const delivered = (prds.get(id) ?? '').startsWith('entregue em')
    c.open.sort()
    if (delivered && c.open.length) out.push({ prd: id, kind: 'open-in-delivered', tickets: c.open })
    if (!delivered && c.total > 0 && !c.open.length) out.push({ prd: id, kind: 'all-done', tickets: [] })
  }
  return out
}

export function driftMessage(drift: Drift[], lang: 'en' | 'pt-br', cwd = ''): string {
  const pt = lang === 'pt-br'
  const lines = drift.map(d => d.kind === 'open-in-delivered'
    ? (pt ? `${d.prd}: tickets ainda abertos (${d.tickets.join(', ')}) mas o PRD está "entregue em"; conclua os tickets ou corrija o Status do PRD.`
      : `${d.prd}: tickets still open (${d.tickets.join(', ')}) but the PRD is marked "entregue em"; close the tickets or fix the PRD Status.`)
    : (pt ? `${d.prd}: todos os tickets ligados estão concluídos mas o Status do PRD não começa com "entregue em"; atualize o Status.`
      : `${d.prd}: all linked tickets are done but the PRD Status does not start with "entregue em"; update the Status.`))
  return `dev-harness: ${pt ? 'deriva de status entre tickets e PRD bloqueia este commit' : 'ticket/PRD status drift blocks this commit'}\n${lines.join('\n')}\n${pt ? `Corrija e commite de novo. (Verificado só em ${cwd || 'o diretório da sessão'}.)` : `Fix it and commit again. (Checked in the session directory only: ${cwd || 'cwd'}.)`}`
}
