import { resolvePath } from '../runtime-guard.ts'

export const WARN_LINES = 5 // PRD-025 R3/R5: more than this many changed lines is not a small correction

// This project's harness dir only (not any `.harness` on disk): repo mode = `<dir>/.harness/` for `dir` the
// working dir or one of its ancestors (the repo root is one of them); global mode = under `<home>/projects/`.
// ponytail: the project <name> is not known to a pure function, so any project under <home>/projects/ passes;
// narrow it by passing the resolved harness dir if a sibling project's folder must warn.
function inHarness(abs: string, cwd: string, home: string): boolean {
  const c = resolvePath('/', cwd).split('/').filter(Boolean)
  const h = home === '/' ? '' : home
  for (let i = c.length; i >= 0; i--) {
    const dir = `/${c.slice(0, i).join('/')}${i ? '/' : ''}.harness`
    if (abs.startsWith(`${dir}/`) && dir !== h) return true // <home> itself is dh's, not this project's
  }
  return !!home && abs.startsWith(`${home === '/' ? '' : home}/projects/`)
}

// A trailing newline does not open a new line.
const lines = (s: unknown): number => (typeof s === 'string' && s ? s.replace(/\n$/, '').split('\n').length : 0)

// PRD-025 R5-R7: warning text for a main-session Edit/Write outside the harness dir, else undefined.
// `exists` answers whether the absolute path is on disk; undefined/throw = unknown, so no warning (fail open).
// Never denies. `home` is the ADR-007 <home>; '' skips the global-mode allowlist.
export async function warnWrite(
  e: any, home: string, cwd: string, exists: (abs: string) => boolean | undefined | Promise<boolean | undefined>,
): Promise<string | undefined> {
  try {
    if (!e || e.agentId || typeof e.file_path !== 'string' || !e.file_path) return undefined
    if (e.tool !== 'Edit' && e.tool !== 'Write') return undefined
    const abs = resolvePath(cwd, e.file_path)
    if (inHarness(abs, cwd, home)) return undefined
    let why: string
    if (e.tool === 'Edit') {
      const n = Math.max(lines(e.old_string), lines(e.new_string))
      if (n <= WARN_LINES) return undefined
      why = `${n}-line Edit`
    } else {
      if (typeof e.content !== 'string') return undefined
      const n = lines(e.content)
      const there = await exists(abs)
      if (there === undefined) return undefined // unknown: fail open
      if (n <= WARN_LINES && there) return undefined
      why = there ? `${n}-line Write` : 'Write of a new file'
    }
    return `Coordinator: ${why} on ${e.file_path}. Dispatch a builder for work beyond a small correction (${WARN_LINES} lines, one existing file) (dev-harness runtime guard, warning only).`
  } catch {
    return undefined // fail open
  }
}
