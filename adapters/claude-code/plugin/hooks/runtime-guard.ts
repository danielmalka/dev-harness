import type { EngineInterface, Register } from 'claude-code'

const FILES = ['MEMORY.md', 'EPOCHAL.md', 'RISKS.md']

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

export function denySpawn($: any, e: any, next: any) {
  try {
    if (e.parentAgentId) return { deny: 'only the Coordinator dispatches agents (dev-harness runtime guard)' }
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
    if (hit) return { deny: `${hit} is written only by the Coordinator (dev-harness runtime guard); return the proposed update instead.` }
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

const guard = async ($: EngineInterface, e: any, next: any) => {
  const { home, cwd, user } = await harnessHome($)
  return guardFile(home, e, next, cwd, user)
}

export const register: Register = on => {
  on('agent.spawn', denySpawn)
  on('tool.call', { tool: 'Write' }, guard)
  on('tool.call', { tool: 'Edit' }, guard)
}
