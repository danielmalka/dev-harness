import type { StatusSnapshot } from '../mods-types.ts'
import { prettyModel } from './modes.ts'
import {
  backgroundPill, BAR_WIDTH, barFill, cacheState, colorFor, type Detail, fmtDuration, fmtTokens, STATUS_CYAN, STATUS_GREEN,
  STATUS_RED, statsWidth,
} from './status.ts'

// Room the engine's own mode label takes left of the hint row (⏵⏵ auto mode on · ).
export const MODE_LABEL_COLUMNS = 22

// The engine's hint, on one row: it can carry a second line holding only a separator, which drew as a stray dot.
export function hintText(hint: string) {
  return hint
    .split(/\r?\n/)
    .map(line => line.trim())
    .filter(line => /[\p{L}\p{N}]/u.test(line))
    .join(' · ')
}

type Ui = { Box: any; Text: any; Button: any }

// The hint and the menu button, when the status line is off.
export function hintRow({ Box, Text }: Ui, hint: string, button: unknown) {
  return (
    <Box justifyContent="space-between" flexGrow={1} gap={1}>
      <Text dimColor wrap="truncate">
        {hintText(hint)}
      </Text>
      {button}
    </Box>
  )
}

// The status line: model, folder and branch, context, 5h/7d limits and cache, then the dh segments, then the
// menu button; under it, while Claude has a background shell or monitor running, the engine's pill for them.
export function statusRow(
  { Box, Text, Button }: Ui,
  p: {
    snapshot: StatusSnapshot
    segments: string[]
    clock: number
    lastResponseAt: number | null
    columns: number
    label: string
    hint: string
    button: unknown
    onPill: () => unknown
  },
) {
  const { snapshot, segments, clock } = p
  const cache = cacheState(p.lastResponseAt, clock)
  const sep = <Text dimColor>{'  │  '}</Text>
  const bar = (percent: number, color: string) => (
    <Text>
      <Text color={color}>{'█'.repeat(barFill(percent))}</Text>
      <Text dimColor>{'░'.repeat(BAR_WIDTH - barFill(percent))}</Text>{' '}
    </Text>
  )

  const ctxPercent = Math.round(snapshot.ctxPercent ?? 0)
  const ctxTokens =
    snapshot.ctxTokens === null ? '' : ` ${fmtTokens(snapshot.ctxTokens)}${snapshot.ctxWindow ? `/${fmtTokens(snapshot.ctxWindow)}` : ''}`
  const limits = snapshot.limits.map(l => ({
    ...l,
    percentText: `${l.percent}%`,
    resetText: l.resetsAt === null ? '' : ` (resets ${fmtDuration(l.resetsAt - clock)})`,
  }))
  const cacheText = cache === null ? '' : `● cache ${cache}`
  const model = prettyModel(snapshot.model)
  const head = `${model} · ${snapshot.dir}${snapshot.branch === '' ? '' : ` · ${snapshot.branch}`}`
  const segmentsWidth = segments.reduce((n, s) => n + 5 + s.length, 0)
  // One row beside the engine's mode label: the most detailed level that leaves room for the button; past the
  // least detailed one, the model and folder go too.
  const room = p.columns - MODE_LABEL_COLUMNS - p.label.length - 4
  const fits = (d: Detail, withHead: boolean) =>
    statsWidth(`${ctxPercent}%${ctxTokens}`, limits, cacheText, d) + segmentsWidth + (withHead ? head.length + 5 : 0) <= room
  const showHead = fits(2, true)
  const detail: Detail = ([0, 1, 2] as const).find(d => fits(d, showHead)) ?? 2

  const stats = [
    showHead && (
      <Text key="head">
        <Text color={STATUS_CYAN}>{model}</Text>
        <Text dimColor> · </Text>
        <Text>{snapshot.dir}</Text>
        {snapshot.branch !== '' && <Text dimColor> · </Text>}
        {snapshot.branch !== '' && <Text color={STATUS_GREEN}>{snapshot.branch}</Text>}
        {sep}
      </Text>
    ),
    <Text key="ctx">
      <Text>ctx </Text>
      {detail === 0 && bar(ctxPercent, colorFor(ctxPercent))}
      <Text bold color={colorFor(ctxPercent)}>{`${ctxPercent}%`}</Text>
      <Text dimColor>{ctxTokens}</Text>
    </Text>,
    ...limits.map(l => (
      <Text key={`limit-${l.label}`}>
        {sep}
        <Text>{l.label} </Text>
        {detail === 0 && bar(l.percent, colorFor(l.percent))}
        <Text bold color={colorFor(l.percent)}>{l.percentText}</Text>
        {detail < 2 && <Text dimColor>{l.resetText}</Text>}
      </Text>
    )),
    cache !== null && (
      <Text key="cache">
        {sep}
        <Text color={cache === 'ok' ? STATUS_GREEN : STATUS_RED}>●</Text>
        <Text dimColor> cache {cache}</Text>
      </Text>
    ),
    ...segments.map(s => (
      <Text key={`dh-${s}`}>
        {sep}
        <Text>{s}</Text>
      </Text>
    )),
  ]

  const statsRow = (
    <Box justifyContent="space-between" flexGrow={1} gap={1}>
      <Text wrap="truncate">{stats}</Text>
      {p.button}
    </Box>
  )
  const pill = backgroundPill(p.hint)
  if (pill === null) return statsRow
  return (
    <Box flexDirection="column" flexGrow={1}>
      {statsRow}
      <Box gap={1}>
        <Button key="background-tasks" plain label={`⎿ ${pill}`} onPress={p.onPill} />
        <Text dimColor wrap="truncate">
          · /tasks to manage
        </Text>
      </Box>
    </Box>
  )
}
