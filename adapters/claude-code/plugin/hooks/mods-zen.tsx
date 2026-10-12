import { atom, read, update } from 'claude-code'
import type { Register } from 'claude-code'

import type { ZenColors } from './mods-types.ts'
import { addStep, colorsFor, parseZenInput, percentOf, todoSteps, updateSteps, ZEN_TOOL } from './lib/zen.ts'
import { zenBand } from './lib/zen-band.tsx'

// Zen: the band above the prompt, hidden tool rows and the steps from zen_progress, TodoWrite and the Task tools.
// The loader reads atoms declared in the file that uses them: the keys below are shared with mods-menu.tsx and panel.ts.
const zenOn = atom({ plugin: 'dh', key: 'zenOn' } as const, false)
const zenTheme = atom({ plugin: 'dh', key: 'zenTheme' } as const, 'blue')
const zenCustom = atom(
  { plugin: 'dh', key: 'zenCustom' } as const,
  { bar: '#3b82f6', text: '#e2e8f0', border: '#3b82f6' } satisfies ZenColors,
)
const zenCollapsed = atom({ plugin: 'dh', key: 'zenCollapsed' } as const, false)
const task = atom({ plugin: 'dh', key: 'task' } as const, null)
const now = atom({ plugin: 'dh', key: 'now' } as const, 0)

export const registerZen: Register = on => {
  // The band above the prompt. A failure draws nothing of ours.
  on('ui.render', { component: 'AbovePrompt' }, async ($, e, next) => {
    try {
      if (e.props.hasSurvey) return next(e)
      const t = await read($, task)
      if (!(await read($, zenOn)) || t === null) return next(e)
      const theme = await read($, zenTheme)
      const ui = $.ui.resolve(e)
      return zenBand(ui, {
        t,
        theme,
        colors: colorsFor(theme, await read($, zenCustom)),
        isCollapsed: await read($, zenCollapsed),
        clock: await read($, now),
        columns: Math.max(30, e.props.bodyColumns ?? 80),
        onToggle: () => update($, zenCollapsed, c => !c),
      })
    } catch {
      return next(e)
    }
  })

  // Tool rows draw nothing while it is on.
  for (const component of ['ToolUse', 'ToolResult', 'ToolGroup', 'ToolProgress'] as const) {
    on('ui.render', { component }, async ($, e, next) => {
      let hide = false
      try { hide = await read($, zenOn) } catch {}
      return hide ? $.ui.resolve(e).Box({ display: 'none' }) : next(e)
    })
  }

  // Claude's own progress tool. Each call sends the whole step list. Fail-open: a state error answers
  // with a plain result instead of failing the tool call.
  on('tool.call', { tool: ZEN_TOOL }, async ($, e) => {
    try {
      if (!(await read($, zenOn))) return { result: 'Zen mode is off; the progress band is not shown.' }
      const { title, steps } = parseZenInput(e)
      // ponytail: a task begun by this call has no elapsed-time ticker (panel.ts starts that one on a typed prompt).
      if ((await read($, task)) === null) {
        const startedAt = await $.clock.now()
        await update($, task, () => ({ title: 'Working', startedAt, doneAt: null, steps: [] }))
        await update($, now, () => startedAt)
      }
      await update($, task, t => (t === null ? t : { ...t, title: title || t.title, steps }))
      const t = await read($, task)
      const done = steps.filter(s => s.status === 'completed').length
      return { result: `Progress band updated: ${done} of ${steps.length} steps done (${t === null ? 0 : percentOf(t)}%).` }
    } catch {
      return { result: 'Progress band unavailable.' }
    }
  })

  // The three below only follow the tool's input or result: a state error is dropped, the call itself goes on once.
  on('tool.call', { tool: 'TodoWrite' }, async ($, e, next) => {
    try { await update($, task, t => (t === null ? t : { ...t, steps: todoSteps(e.todos) })) } catch {}
    return next(e)
  })

  on('tool.call', { tool: 'TaskCreate' }, async ($, e, next) => {
    const ran = await next(e)
    try {
      const id = (ran.result as { task?: { id?: string } } | undefined)?.task?.id
      if (id !== undefined) await update($, task, t => (t === null ? t : { ...t, steps: addStep(t.steps, id, e.subject) }))
    } catch {}
    return ran
  })

  on('tool.call', { tool: 'TaskUpdate' }, async ($, e, next) => {
    const ran = await next(e)
    if (ran.deny !== undefined || ran.isError === true) return ran
    try { await update($, task, t => (t === null ? t : { ...t, steps: updateSteps(t.steps, e) })) } catch {}
    return ran
  })
}
