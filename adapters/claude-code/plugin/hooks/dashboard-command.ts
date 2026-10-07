import type { On } from 'claude-code'

// PRD-012 R6: `/dh:dashboard` exists only as a command the mod registers; it runs the plugin's own
// `dh dashboard --detach` (idempotent) and answers with the URL it prints.
export const DASHBOARD_COMMAND = { name: 'dashboard', description: 'Start the Dev Harness dashboard (local, detached) and print its URL' }

export function dhPath(root: string): string {
  return /^[A-Za-z]:[\\/]|^\\\\/.test(root) ? `${root}\\bin\\dh.cmd` : `${root}/bin/dh`
}

export function registerDashboard(on: On): void {
  on('command.run', { command: 'dashboard' }, async ($, e) => {
    try {
      const extra = e.args.trim().split(/\s+/).filter(Boolean) // dh flags only: argv, no shell
      const r = await $.process.run([dhPath($.plugin.root), 'dashboard', '--detach', ...extra], { timeoutMs: 15_000 })
      const url = r.stdout.trim()
      if (r.exitCode === 0 && !url) return { text: 'dh dashboard exited 0 but printed no URL' }
      return { text: r.exitCode === 0 ? url : `dh dashboard failed (exit ${r.exitCode}): ${(r.stderr || r.stdout).trim()}` }
    } catch (err) {
      return { text: `dh dashboard failed: ${String(err)}` }
    }
  }).catch((_$: unknown, e: any, next: any) => next(e))
}
