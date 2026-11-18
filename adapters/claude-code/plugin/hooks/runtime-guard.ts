import type { EngineInterface, Register } from 'claude-code'
import { warnWrite } from './lib/warn.ts'

const FILES = ['MEMORY.md', 'EPOCHAL.md', 'RISKS.md', 'FACTS.md']

// Normalizes `p` against `base`: absolute result, no `.`/`..`/empty segments. A Windows drive path
// (`C:\\x`, `c:/x`) is absolute too and reads as `/C:/x`, so git's `C:/repo` and the engine's `C:\\repo` meet.
export function resolvePath(base: string, p: string): string {
  const abs = (x: string) => x.replace(/\\/g, '/').replace(/^\/?([A-Za-z]):(?=\/|$)/, (_m, d: string) => `/${d.toUpperCase()}:`)
  const n = abs(p)
  const out: string[] = []
  for (const seg of (n.startsWith('/') ? n : `${abs(base)}/${n}`).split('/')) {
    if (seg === '..') out.pop()
    else if (seg && seg !== '.') out.push(seg)
  }
  return `/${out.join('/')}`
}

// ADR-007: <home> = DH_HOME, else HOME/.harness, else USERPROFILE/.harness; '' when none is set.
// Pure, so each module reads the three variables itself (literal names; $ does not cross an import).
export function homeFrom(dhHome: unknown, home: unknown, userProfile: unknown, cwd = '/'): string {
  const base = dhHome ? String(dhHome) : home ? `${home}/.harness` : userProfile ? `${userProfile}/.harness` : ''
  return base ? resolvePath(cwd, base) : ''
}

// PRD-035 R12/D7: the one subagent allowed to write FACTS.md is the curator, a `docs-guide` whose dispatch names FACTS.md in its write set.
// ponytail: the write set is read from the dispatch text (a line `Arquivos permitidos:`, `Conjunto de escrita:` or `Write set:` that names FACTS.md), because the
// engine gives the mod no structured write set; only the main session can spawn (denySpawn), so the text is the Coordinator's own. State is per mod load,
// like the agentTypes map of panel.ts.
const curators = new Set<string>()
const WRITE_SET_LINE = /^[ \t>*-]*(?:arquivos permitidos|conjunto de escrita|write set)[ \t]*:[^\n]*(?<![\w.-])FACTS\.md(?![\w-]|\.\w)/im

export function isCuratorSpawn(subagentType: unknown, prompt: unknown): boolean {
  if (typeof subagentType !== 'string' || typeof prompt !== 'string') return false
  return (subagentType === 'dh:docs-guide' || subagentType === 'docs-guide') && WRITE_SET_LINE.test(prompt)
}

// Called with the spawn event and the id the engine gave the new agent; remembers the curators of this session.
export function noteSpawn(e: any, agentId: unknown): void {
  try {
    if (typeof agentId !== 'string' || !agentId) return
    if (isCuratorSpawn(e?.subagentType, e?.prompt)) curators.add(agentId)
    else curators.delete(agentId) // an id that is reused by a non-curator loses the exception
  } catch {} // fail open
}

export function denySpawn($: any, e: any, next: any) {
  try {
    if (e.parentAgentId) return { deny: 'só o Coordinator despacha agentes (trava de runtime do dev-harness)' }
  } catch {} // fail open
  return next(e)
}

// The protected record `file_path` names, or undefined: `<x>/.harness/<file>` (repo mode, by suffix) or
// exactly `<home>/projects/<one segment>/<file>` (global mode), after expanding a leading `~` to `userHome`
// and resolving against `cwd`. Case-insensitive on every platform: a look-alike denied is the safe side.
// `home` '' turns the global pattern off (fail open, ADR-007 §4).
// ponytail: accepted limits (ADR-005): a symlinked spelling of <home> is not resolved, and only Write/Edit
// are guarded; a Bash redirect or any other write tool is outside the guard.
export function protectedHit(home: string, filePath: unknown, cwd = '/', userHome = ''): string | undefined {
  if (typeof filePath !== 'string' || !filePath) return undefined
  const expanded = userHome ? filePath.replace(/^~(?=$|[\\/])/, userHome) : filePath
  const abs = resolvePath(cwd, expanded) // `..` cannot walk around either pattern
  const low = abs.toLowerCase()
  const repo = FILES.find(f => low.endsWith(`/.harness/${f.toLowerCase()}`))
  if (repo) return `.harness/${repo}`
  if (!home) return undefined
  const h = resolvePath('/', home)
  const prefix = `${h === '/' ? '' : h}/projects/`
  if (!low.startsWith(prefix.toLowerCase())) return undefined
  const rest = abs.slice(prefix.length).split('/')
  const file = FILES.find(f => f.toLowerCase() === rest[1]?.toLowerCase())
  return rest.length === 2 && rest[0] && file ? abs : undefined
}

export function guardFile(home: string, e: any, next: any, cwd = '/', userHome = '') {
  try {
    const hit = e.agentId ? protectedHit(home, e.file_path, cwd, userHome) : undefined
    if (hit && /(^|\/)facts\.md$/i.test(hit) && curators.has(e.agentId)) return next(e) // PRD-035 R12: the curator may write FACTS.md
    if (hit) return { deny: `${hit} só é escrito pelo Coordinator (trava de runtime do dev-harness); devolva a atualização proposta.` }
  } catch {} // fail open
  return next(e)
}

async function harnessHome($: EngineInterface): Promise<{ home: string; cwd: string; user: string }> {
  try {
    const cwd = await $.session.cwd().catch(() => '/')
    const dh = await $.env.get('DH_HOME'), hm = await $.env.get('HOME'), up = await $.env.get('USERPROFILE')
    return { home: homeFrom(dh, hm, up, cwd), cwd, user: hm || up || '' }
  } catch {
    return { home: '', cwd: '/', user: '' } // fail open: only the suffix pattern holds
  }
}

// PRD-025 R5-R7: warn-only toast for the main session; the event always passes on unchanged.
export async function warnThenNext($: any, e: any, home: string, cwd: string, next: any) {
  try {
    const text = await warnWrite(e, home, cwd, p => $.fs.exists(p))
    if (text) $.ui.toast(text)
  } catch {} // fail open
  return next(e)
}

const guard = async ($: EngineInterface, e: any, next: any) => {
  const { home, cwd, user } = await harnessHome($)
  return guardFile(home, e, (x: any) => warnThenNext($, x, home, cwd, next), cwd, user)
}

export const register: Register = on => {
  on('agent.spawn', ($, e, next) => denySpawn($, e, (x: any) => {
    const r = next(x)
    // sync result or promise, unchanged for the caller; the curator is noted when the engine answers with the new agent id
    if (r && typeof (r as any).then === 'function') (r as any).then((v: any) => noteSpawn(x, v?.agentId), () => undefined)
    else noteSpawn(x, (r as any)?.agentId)
    return r
  }))
  on('tool.call', { tool: 'Write' }, guard)
  on('tool.call', { tool: 'Edit' }, guard)
}
