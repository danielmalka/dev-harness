import { atom, read, update } from 'claude-code'
import type { EngineInterface, Register, Timer } from 'claude-code'

import type { Effort, Panel, RateLimit, StatusSnapshot, ZenColors, ZenTheme } from './mods-types.ts'
import { EFFORTS, type Mode, modeLabel, prettyModel } from './lib/modes.ts'
import { LIMIT_LABELS } from './lib/status.ts'
import { HEX, ZEN_THEMES, ZEN_TOOL_DESCRIPTION, ZEN_TOOL_NAME, ZEN_TOOL_SCHEMA } from './lib/zen.ts'
import { configDir, moveStatusLine, parseSettings, type Settings, withDefaultMode } from './lib/settings.ts'
import { hintRow, statusRow } from './lib/status-line.tsx'
import { MENU_NAME, menuPage } from './lib/menu-pages.tsx'

// The mods menu (/mods) and the state behind it: settings.json writes (R5 of PRD-019), Mode, the status line.
// The loader follows $ only into functions of this file and reads atoms declared here: the keys below are
// shared with mods-zen.tsx and panel.ts.
const menu = atom({ plugin: 'dh', key: 'menu' } as const, null)
const zenOn = atom({ plugin: 'dh', key: 'zenOn' } as const, false)
const zenTheme = atom({ plugin: 'dh', key: 'zenTheme' } as const, 'blue')
const zenCustom = atom(
  { plugin: 'dh', key: 'zenCustom' } as const,
  { bar: '#3b82f6', text: '#e2e8f0', border: '#3b82f6' } satisfies ZenColors,
)
const task = atom({ plugin: 'dh', key: 'task' } as const, null)
const effort = atom({ plugin: 'dh', key: 'effort' } as const, null)
const defaultMode = atom({ plugin: 'dh', key: 'defaultMode' } as const, null)
const sessionEffort = atom({ plugin: 'dh', key: 'sessionEffort' } as const, null)
const sessionModel = atom({ plugin: 'dh', key: 'sessionModel' } as const, null)
const statusOn = atom({ plugin: 'dh', key: 'statusOn' } as const, false)
const status = atom({ plugin: 'dh', key: 'status' } as const, null)
const dhSegments = atom({ plugin: 'dh', key: 'dhSegments' } as const, [])
const lastResponseAt = atom({ plugin: 'dh', key: 'lastResponseAt' } as const, null)
const statusClock = atom({ plugin: 'dh', key: 'statusClock' } as const, 0)
// True once the status line is on and a PromptHint exists (start-up runs from its first draw; drawing itself
// cannot write state): panel.ts keeps its own bar until then, and on surfaces without a PromptHint.
const statusDrawn = atom({ plugin: 'dh', key: 'statusDrawn' } as const, false)

const PANE = 'dh-mods'
const PANE_SIZE = { rows: 26, columns: 56 }
const BRANCH_MS = 5_000

// ~/.claude/settings.json ----------------------------------------------------------------
// Parsed first, so an invalid file is never overwritten; the first change keeps the original text in
// settings.json.dh.bak, once.

async function settingsPath($: EngineInterface) {
  const home = (await $.env.get('HOME')) || (await $.env.get('USERPROFILE'))
  return `${configDir(await $.env.get('CLAUDE_CONFIG_DIR'), home)}/settings.json`
}

async function readSettings($: EngineInterface, path: string): Promise<{ text: string; settings: Settings }> {
  if (!(await $.fs.exists(path))) return { text: '{}', settings: {} }
  const text = (await $.fs.read(path)) as string
  return { text, settings: parseSettings(text) }
}

async function writeSettings($: EngineInterface, path: string, text: string, settings: Settings) {
  const backup = `${path}.dh.bak`
  if (!(await $.fs.exists(backup))) await $.fs.write(backup, text)
  await $.fs.write(path, `${JSON.stringify(settings, null, 2)}\n`)
}

// On: the user's own statusLine is kept in the plugin store and taken out of settings.json; off puts it back.
async function moveUserStatusLine($: EngineInterface, isOn: boolean) {
  const path = await settingsPath($)
  const { text, settings } = await readSettings($, path)
  const moved = moveStatusLine(settings, isOn, (await $.store.get('savedStatusLine')) as { statusLine: unknown } | undefined)
  if (moved === undefined) return
  if (moved.save !== undefined) await $.store.set('savedStatusLine', moved.save)
  await writeSettings($, path, text, moved.next)
  if (moved.save === undefined) await $.store.delete('savedStatusLine')
}

// permissions.defaultMode: the mode new sessions start in. A running session's mode is shift+tab's; a mod cannot set it.
async function setDefaultMode($: EngineInterface, next: Mode) {
  const path = await settingsPath($)
  const { text, settings } = await readSettings($, path)
  await writeSettings($, path, text, withDefaultMode(settings, next))
  await update($, defaultMode, () => next)
  $.ui.toast(`New sessions start in ${modeLabel(next)} mode`)
}

// Zen settings ---------------------------------------------------------------------------

async function setZenOn($: EngineInterface, isOn: boolean) {
  await update($, zenOn, () => isOn)
  await $.store.set('zenOn', isOn)
  if (!isOn) await update($, task, () => null)
}

async function setZenTheme($: EngineInterface, theme: ZenTheme) {
  await update($, zenTheme, () => theme)
  await $.store.set('zenTheme', theme)
}

async function setZenColor($: EngineInterface, part: keyof ZenColors, value: string) {
  if (!HEX.test(value.trim())) {
    $.ui.toast(`${value} is not a hex color like #3b82f6`)
    return
  }
  const colors = await update($, zenCustom, c => ({ ...c, [part]: value.trim() }))
  await $.store.set('zenCustom', colors)
}

// Mode and effort ------------------------------------------------------------------------

// The effort every request of the main loop is sent with. When /config has an effort row, it is changed too.
async function setEffort($: EngineInterface, level: Effort) {
  await update($, effort, () => level)
  await $.store.set('effort', level)
  try {
    const row = (await $.config.list()).find(r => /effort/i.test(r.key) && r.options?.includes(level))
    if (row !== undefined) await $.config.set({ key: row.key, value: level })
  } catch {}
  $.ui.toast(`Effort: ${level}`)
}

// Switches the session's model through Claude Code's own `/model`, so /model, the status line and the requests agree.
async function setModel($: EngineInterface, alias: string) {
  await $.command.run({ command: 'model', args: alias })
  const id = await $.session.model()
  await update($, sessionModel, () => id)
  if (await read($, statusOn)) void refreshStatus($, true).catch(() => {})
  $.ui.toast(`Model: ${prettyModel(id)}`)
}

async function trackModel($: EngineInterface, id: string) {
  if (id !== (await read($, sessionModel))) await update($, sessionModel, () => id)
}

// The session's own effort level before any request. First match wins: a /config row, the effortLevel setting, CLAUDE_EFFORT, the last level seen.
async function loadSessionEffort($: EngineInterface) {
  const found: unknown[] = []
  try {
    found.push((await $.config.list()).find(r => /effort/i.test(r.key))?.value)
  } catch {}
  try {
    found.push(((await $.settings.read()) as { effortLevel?: unknown }).effortLevel)
  } catch {}
  found.push(await $.env.get('CLAUDE_EFFORT'))
  found.push(await $.store.get('lastSessionEffort'))
  const level = found.find((v): v is Effort => typeof v === 'string' && EFFORTS.includes(v as Effort))
  if (level !== undefined) await update($, sessionEffort, () => level)
}

// Records the level a request actually ran at, as the engine reports it.
async function trackEffort($: EngineInterface, level: unknown) {
  if (typeof level !== 'string' || !EFFORTS.includes(level as Effort)) return
  if (level === (await read($, sessionEffort))) return
  await update($, sessionEffort, () => level as Effort)
  await $.store.set('lastSessionEffort', level)
}

async function loadSettings($: EngineInterface) {
  const saved = {
    zenOn: (await $.store.get('zenOn')) as boolean | undefined,
    zenTheme: (await $.store.get('zenTheme')) as ZenTheme | undefined,
    zenCustom: (await $.store.get('zenCustom')) as ZenColors | undefined,
    effort: (await $.store.get('effort')) as Effort | undefined,
  }
  if (saved.zenOn !== undefined) await update($, zenOn, () => saved.zenOn!)
  if (saved.zenTheme !== undefined && ZEN_THEMES.includes(saved.zenTheme)) await update($, zenTheme, () => saved.zenTheme!)
  if (saved.zenCustom !== undefined) await update($, zenCustom, () => saved.zenCustom!)
  if (saved.effort !== undefined && EFFORTS.includes(saved.effort)) await update($, effort, () => saved.effort!)
  try {
    const { permissions } = await $.settings.read({ source: 'user' })
    await update($, defaultMode, () => (permissions as { defaultMode?: string } | undefined)?.defaultMode ?? 'default')
  } catch {}
}

// Status line ----------------------------------------------------------------------------

// The git branch is looked up at most once per BRANCH_MS (and when the folder changes).
let branchAt = 0
let branchCwd = ''
let branch = ''

// Gathers what the status line draws: the model, the folder and branch, the context window and the rate limits.
async function refreshStatus($: EngineInterface, force = false) {
  const [model, cwd, usage] = await Promise.all([$.session.model(), $.session.cwd(), $.session.usage()])
  if (force || cwd !== branchCwd || Date.now() - branchAt >= BRANCH_MS) {
    branchAt = Date.now()
    branchCwd = cwd
    branch = ''
    try {
      const ran = await $.process.run(['git', '-C', cwd, 'branch', '--show-current'])
      if (ran.exitCode === 0) branch = ran.stdout.trim()
    } catch {}
  }
  const limits: RateLimit[] = usage.rateLimits
    .filter(l => LIMIT_LABELS[l.kind] !== undefined)
    .map(l => ({
      label: LIMIT_LABELS[l.kind]!,
      percent: Math.round(l.percentUsed),
      resetsAt: l.resetsAt === undefined ? null : Date.parse(l.resetsAt),
    }))
  const snapshot: StatusSnapshot = {
    model,
    dir: cwd.replace(/[\\/]+$/, '').split(/[\\/]/).pop() || cwd,
    branch,
    ctxPercent: usage.context.percent ?? null,
    ctxTokens: usage.context.tokens ?? null,
    ctxWindow: usage.context.window,
    limits,
  }
  await update($, status, () => snapshot)
}

async function setStatusOn($: EngineInterface, isOn: boolean) {
  await moveUserStatusLine($, isOn)
  await update($, statusOn, () => isOn)
  await $.store.set('statusOn', isOn)
  if (isOn) await refreshStatus($, true)
  await update($, statusDrawn, () => isOn)
}

async function loadStatus($: EngineInterface) {
  const isOn = (await $.store.get('statusOn')) as boolean | undefined
  if (isOn === true) await update($, statusOn, () => true)
  // Never chosen: on by default (the user's own statusLine, if any, is kept and comes back with Off). When
  // that fails (settings.json unreadable) it stays off, and the failure is remembered so panel.ts keeps its bar.
  if (isOn === undefined) {
    await setStatusOn($, true).catch(async (error: unknown) => {
      await $.store.set('statusOn', false)
      $.ui.toast(`dh status line stays off: ${String(error)}`)
    })
  }
  await update($, statusClock, () => Date.now())
  await refreshStatus($, true).catch(() => undefined)
  if ((await read($, statusOn)) && (await read($, status)) !== null) await update($, statusDrawn, () => true)
}

// Start-up ------------------------------------------------------------------------------

// The status line's clock: reset countdowns and the cache state move with it.
let statusTicker: Timer | undefined
let started = false
let tries = 0

// Everything the mods need once per session. Run from the first PromptHint draw: session.start is panel.ts's.
// A failed tool or command registration allows a retry on a later draw (three at most).
async function start($: EngineInterface) {
  if (started || tries >= 3) return
  started = true
  tries++
  await loadSettings($).catch(() => undefined)
  await loadStatus($).catch(() => undefined)
  await loadSessionEffort($).catch(() => undefined)
  await $.session.model().then(id => trackModel($, id)).catch(() => undefined)
  statusTicker?.cancel()
  try {
    statusTicker = $.clock.every(15_000, () => {
      void update($, statusClock, () => Date.now())
    })
  } catch {}
  const tool = await $.tool.register({ name: ZEN_TOOL_NAME, description: ZEN_TOOL_DESCRIPTION, inputSchema: ZEN_TOOL_SCHEMA }).then(() => true, () => false)
  const command = await $.command.register({
    name: 'mods',
    description: 'Open the dh mods menu: status line, Zen mode and Mode',
    argumentHint: '[status|zen|mode]',
    immediate: true,
  }).then(() => true, () => false)
  if (!tool || !command) started = false
}

// Menu ----------------------------------------------------------------------------------

const goTo = ($: EngineInterface, panel: Panel) => update($, menu, () => panel)

// Opens the menu in its pane with the keyboard, or closes it.
async function toggleMenu($: EngineInterface, panel: Panel = 'main') {
  if ((await read($, menu)) !== null) {
    // A plugin's own close skips its own ui.close hook, so reset here too.
    await $.ui.close({ id: PANE })
    await update($, menu, () => null)
    return
  }
  await update($, menu, () => panel)
  await $.ui.open({ id: PANE, title: MENU_NAME, focus: true, closeOnEscape: true, ...PANE_SIZE })
}

async function openTasks($: EngineInterface) {
  try {
    await $.command.run({ command: 'tasks' })
  } catch (error) {
    $.ui.toast(`Could not open /tasks: ${String(error)}`)
  }
}

// Watching hooks end in .catch(... next(e)): when one fails, the event goes on.
const failOpen = ($: unknown, e: any, next: any) => next(e)

export const registerMods: Register = on => {
  on('command.run', { command: 'mods' }, async ($, e) => {
    try {
      const arg = e.args.trim().toLowerCase()
      const panel: Panel = arg === 'status' || arg === 'zen' || arg === 'mode' ? arg : 'main'
      if ((await read($, menu)) !== null) await goTo($, panel)
      else await toggleMenu($, panel)
      return { text: 'The dh mods menu is open.' }
    } catch (error) {
      return { text: `The dh mods menu could not open: ${String(error)}` }
    }
  })

  on('ui.render', { component: 'Pane', requestId: PANE }, async ($, e) => {
    const panel = await read($, menu)
    if (panel === null || (e.surface !== 'terminal' && e.surface !== 'desktop')) return $.ui.resolve(e).Text({ children: '' })
    const tell = (what: string) => (error: unknown) => $.ui.toast(`${what}: ${String(error)}`)
    return menuPage(
      $.ui.resolve(e),
      panel,
      { bodyColumns: e.props.bodyColumns, bodyRows: e.props.scroll.bodyRows },
      {
        isZenOn: await read($, zenOn),
        zenTheme: await read($, zenTheme),
        custom: await read($, zenCustom),
        level: (await read($, effort)) ?? (await read($, sessionEffort)),
        model: await read($, sessionModel),
        isStatusOn: await read($, statusOn),
        defaultMode: await read($, defaultMode),
      },
      {
        open: p => goTo($, p),
        zenOn: v => setZenOn($, v),
        zenTheme: t => setZenTheme($, t),
        zenColor: (part, value) => setZenColor($, part, value),
        statusOn: v => setStatusOn($, v).catch(tell(`Could not turn it ${v ? 'on' : 'off'}`)),
        model: alias => setModel($, alias).catch(tell('Could not switch the model')),
        effort: l => setEffort($, l),
        mode: m => setDefaultMode($, m as Mode).catch(tell('Could not save the mode')),
      },
    )
  })

  // However the pane closes (esc, its close mark, the corner button), the menu resets.
  on('ui.close', { id: PANE }, async ($, e, next) => {
    const closed = await next(e)
    await update($, menu, () => null)
    return closed
  }).catch(failOpen)

  // The corner button and, when on, the status line, under the prompt.
  on('ui.render', { component: 'PromptHint' }, async ($, e) => {
    void start($).catch(() => undefined)
    const ui = $.ui.resolve(e)
    const label = `◆ ${MENU_NAME} ${(await read($, menu)) !== null ? '▼' : '▲'}`
    const { Button } = ui
    const button = <Button key="mods-toggle" plain label={label} onPress={() => toggleMenu($)} />
    const snapshot = (await read($, statusOn)) ? await read($, status) : null
    if (snapshot === null) return hintRow(ui, e.props.hint, button)
    return statusRow(ui, {
      snapshot,
      segments: await read($, dhSegments),
      clock: await read($, statusClock),
      lastResponseAt: await read($, lastResponseAt),
      columns: e.viewport?.columns ?? 120,
      label,
      hint: e.props.hint,
      button,
      onPill: () => openTasks($),
    })
  })

  // Mode: the chosen effort rides on every model request of the main loop; its model is noted for the Mode page.
  on('turn.step', async function* ($, e, next) {
    const isMain = e.agentId === undefined
    let level: Effort | null = null
    try {
      level = await read($, effort)
    } catch {}
    if (isMain) {
      // The session's own level, as the engine would send it, before ours replaces it.
      await trackEffort($, e.effort).catch(() => undefined)
      await trackModel($, e.model).catch(() => undefined)
    }
    const answered = yield* next(isMain && level !== null && e.effort !== undefined ? { ...e, effort: level } : e)
    // The main loop's responses keep the prompt cache warm; subagents' do not.
    if (isMain && answered.usage !== null) {
      const at = Date.now()
      await update($, lastResponseAt, () => at).catch(() => undefined)
      await update($, statusClock, () => at).catch(() => undefined)
      if (await read($, statusOn).catch(() => false)) void refreshStatus($).catch(() => {})
    }
    return answered
  })
}
