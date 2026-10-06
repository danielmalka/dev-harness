import type { Register } from 'claude-code'

const PROTECTED = ['.harness/MEMORY.md', '.harness/EPOCHAL.md', '.harness/RISKS.md']

export function denySpawn($: any, e: any, next: any) {
  try {
    if (e.parentAgentId) return { deny: 'only the Coordinator dispatches agents (dev-harness runtime guard)' }
  } catch {} // fail open
  return next(e)
}

// ponytail: Write/Edit only; a Bash redirect into these files is not blocked.
export function guardFile($: any, e: any, next: any) {
  try {
    const p = String(e.file_path ?? '').replace(/\\/g, '/')
    const hit = PROTECTED.find(f => p === f || p.endsWith('/' + f))
    if (e.agentId && hit) return { deny: `${hit} is written only by the Coordinator (dev-harness runtime guard); return the proposed update instead.` }
  } catch {} // fail open
  return next(e)
}

export const register: Register = on => {
  on('agent.spawn', denySpawn)
  on('tool.call', { tool: 'Write' }, guardFile)
  on('tool.call', { tool: 'Edit' }, guardFile)
}
