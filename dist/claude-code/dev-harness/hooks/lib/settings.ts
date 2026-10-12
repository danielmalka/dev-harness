// ~/.claude/settings.json edits (R5 of PRD-019), on a parsed object: other keys are kept as they are and
// written back as 2-space JSON; an unreadable or invalid file throws before anything is written.
export type Settings = Record<string, any>

export function parseSettings(text: string): Settings {
  const s = JSON.parse(text)
  if (!s || typeof s !== 'object' || Array.isArray(s)) throw new Error('settings.json is not a JSON object')
  return s
}

const ABSOLUTE = /^(\/|[A-Za-z]:[\\/]|\\\\)/

// CLAUDE_CONFIG_DIR when it is an absolute path (POSIX, drive or UNC), else <home>/.claude; never a relative path.
export function configDir(configDirEnv: string | undefined, home: string | undefined): string {
  const dir = configDirEnv && ABSOLUTE.test(configDirEnv) ? configDirEnv : home ? `${home}/.claude` : ''
  if (!ABSOLUTE.test(dir)) throw new Error('no absolute CLAUDE_CONFIG_DIR, HOME or USERPROFILE: settings.json not touched')
  return dir.replaceAll('\\', '/').replace(/\/+$/, '')
}

// On: the user's statusLine leaves the file (and is kept by the caller); off: the kept one comes back.
// undefined when nothing changes.
export function moveStatusLine(s: Settings, isOn: boolean, saved: { statusLine: unknown } | undefined) {
  if (isOn && s.statusLine !== undefined) {
    const { statusLine, ...next } = s
    return { next, save: { statusLine } }
  }
  if (!isOn && saved !== undefined && s.statusLine === undefined) return { next: { ...s, statusLine: saved.statusLine }, save: undefined }
  return undefined
}

export const withDefaultMode = (s: Settings, mode: string): Settings => ({ ...s, permissions: { ...s.permissions, defaultMode: mode } })
