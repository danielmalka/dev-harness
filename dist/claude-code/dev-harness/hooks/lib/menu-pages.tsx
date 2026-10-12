import type { Effort, Panel, ZenColors, ZenTheme } from '../mods-types.ts'
import { EFFORTS, MODELS, MODES, modeLabel, modelFamily, prettyModel } from './modes.ts'
import { colorsFor, ZEN_THEMES } from './zen.ts'

export const MENU_NAME = 'dh mods'
const PANE_BACKGROUND = '#0b0b0b'

// What the pages show, and what their buttons do (closures built where $ is in scope).
export type MenuView = {
  isZenOn: boolean
  zenTheme: ZenTheme
  custom: ZenColors
  level: Effort | null
  model: string | null
  isStatusOn: boolean
  defaultMode: string | null
}
export type MenuActions = {
  open: (p: Panel) => unknown
  zenOn: (on: boolean) => unknown
  zenTheme: (t: ZenTheme) => unknown
  zenColor: (part: keyof ZenColors, value: string) => unknown
  statusOn: (on: boolean) => unknown
  model: (alias: string) => unknown
  effort: (l: Effort) => unknown
  mode: (m: string) => unknown
}

// The menu, one item per row so the arrow keys walk it top to bottom.
export function menuPage(
  { Box, Text, Button, Input }: { Box: any; Text: any; Button: any; Input: any },
  panel: Panel,
  size: { bodyColumns: number; bodyRows: number },
  v: MenuView,
  a: MenuActions,
) {
  // A radio row: ● for the chosen one, focusable and pressable.
  const option = (key: string, label: string, isOn: boolean, onPress: () => unknown, swatch?: string) => (
    <Box key={`row-${key}`} gap={1}>
      {swatch !== undefined && <Text color={swatch}>■</Text>}
      <Button key={key} plain label={`${isOn ? '●' : '○'} ${label}`} dimColor={!isOn} onPress={onPress} />
    </Box>
  )
  const heading = (text: string) => (
    <Box key={`h-${text}`} marginTop={1}>
      <Text bold>{text}</Text>
    </Box>
  )
  const page = (title: string, rows: unknown[]) => (
    <Box flexDirection="column" paddingX={1} backgroundColor={PANE_BACKGROUND} width={size.bodyColumns} minHeight={size.bodyRows}>
      <Box justifyContent="space-between">
        <Text bold color="claude">
          ◆ {title}
        </Text>
        {title !== MENU_NAME && <Button key="mods-back" plain dimColor hotkey="b" label="← back" onPress={() => a.open('main')} />}
      </Box>
      {rows}
      <Box marginTop={1}>
        <Text dimColor>↑↓ move · enter select · esc close</Text>
      </Box>
    </Box>
  )

  if (panel === 'zen') {
    return page('Zen mode', [
      <Text key="zen-about" dimColor>
        Hides tool calls and their output. Shows Claude's text and a progress band above the prompt.
      </Text>,
      heading('Zen mode'),
      option('zen-on', 'On', v.isZenOn, () => a.zenOn(true)),
      option('zen-off', 'Off', !v.isZenOn, () => a.zenOn(false)),
      heading('Band color'),
      ...ZEN_THEMES.map(t => option(`zen-theme-${t}`, t, t === v.zenTheme, () => a.zenTheme(t), colorsFor(t, v.custom).bar)),
      v.zenTheme === 'custom' && (
        <Box key="zen-custom" flexDirection="column">
          {(['bar', 'text', 'border'] as const).map(part => (
            <Input
              key={`zen-color-${part}`}
              label={`${part} color`}
              value={v.custom[part]}
              placeholder="#3b82f6"
              submitLabel="set"
              onSubmit={(value: string) => a.zenColor(part, value)}
            />
          ))}
        </Box>
      ),
    ])
  }

  if (panel === 'status') {
    return page('Status line', [
      <Text key="status-about" dimColor>
        Model, folder and branch, then context, 5h and 7d limits and cache, with the dh progress on the same row.
      </Text>,
      heading('Status line'),
      option('status-on', 'On', v.isStatusOn, () => a.statusOn(true)),
      option('status-off', 'Off', !v.isStatusOn, () => a.statusOn(false)),
      <Text key="status-note" dimColor>
        On takes your own statusLine command out of settings.json; Off puts it back.
      </Text>,
    ])
  }

  if (panel === 'mode') {
    const family = modelFamily(v.model)
    return page('Mode', [
      heading('Model'),
      ...MODELS.map(m => option(`model-${m.alias}`, m.label, m.alias === family, () => a.model(m.alias))),
      <Button key="model-default" plain dimColor label="↺ Claude Code's default" onPress={() => a.model('default')} />,
      heading('Effort'),
      ...EFFORTS.map(l => option(`effort-${l}`, l, l === v.level, () => a.effort(l))),
      heading('Running in'),
      // The engine gives mods no permission-mode event (BUG-001: engine events only), so the running mode is unknown.
      <Text key="mode-running">
        {modeLabel(null)}
        <Text dimColor> · shift+tab switches it</Text>
      </Text>,
      heading('New sessions start in'),
      ...MODES.map(m => (
        <Box key={`row-default-mode-${m}`} gap={1}>
          <Button key={`default-mode-${m}`} plain label={m === v.defaultMode ? '●' : '○'} onPress={() => a.mode(m)} />
          <Text bold={m === v.defaultMode}>{modeLabel(m)}</Text>
        </Box>
      )),
      v.defaultMode === 'bypassPermissions' && (
        <Text key="mode-warning" color="warning">
          Bypass skips every permission check. Claude Code asks you to confirm it when the next session starts.
        </Text>
      ),
    ])
  }

  const row = (hotkey: string, name: string, value: string, target: Panel) => (
    <Box key={`row-${target}`} flexDirection="column" marginTop={1}>
      <Button
        key={`open-${target}`}
        plain
        hotkey={hotkey}
        autoFocus={target === 'status' ? true : undefined}
        label={`${name}  ›`}
        onPress={() => a.open(target)}
      />
      <Text dimColor>   {value}</Text>
    </Box>
  )
  return page(MENU_NAME, [
    row('1', 'Status line', v.isStatusOn ? 'on' : 'off', 'status'),
    row('2', 'Zen mode', v.isZenOn ? `on · ${v.zenTheme} band` : 'off', 'zen'),
    row('3', 'Mode', `${v.model === null ? '' : `${prettyModel(v.model)} · `}effort ${v.level ?? 'default'}`, 'mode'),
  ])
}
