import type { Effort } from '../mods-types.ts'

export const EFFORTS: readonly Effort[] = ['low', 'medium', 'high', 'xhigh', 'max']

export const MODES = ['default', 'plan', 'acceptEdits', 'auto', 'bypassPermissions'] as const

export type Mode = (typeof MODES)[number]

export const MODE_LABELS: Record<string, string> = {
  default: 'ask',
  plan: 'plan',
  acceptEdits: 'accept edits',
  auto: 'auto',
  dontAsk: "don't ask",
  bypassPermissions: 'bypass permissions',
}

export const modeLabel = (m: string | null): string => (m === null ? 'unknown' : (MODE_LABELS[m] ?? m))

// The models the Mode page offers, by the alias `/model` takes.
export const MODELS = [
  { alias: 'fable', label: 'Fable' },
  { alias: 'opus', label: 'Opus' },
  { alias: 'sonnet', label: 'Sonnet' },
  { alias: 'haiku', label: 'Haiku' },
] as const

// `claude-opus-5-5` → `opus`; null for an id that is not a Claude model's.
export const modelFamily = (id: string | null): string | null =>
  id?.match(/^claude-([a-z]+)-/i)?.[1]?.toLowerCase() ?? null

// `claude-opus-5-5` (what $.session.model() gives) as `Opus 5.5`; a name that is not a Claude model id stays as it is.
export function prettyModel(id: string): string {
  const match = id.match(/^claude-([a-z]+)-(\d+)(?:-(\d+))?/i)
  if (match === null) return id
  const [, family = '', major = '', minor] = match
  return `${family.charAt(0).toUpperCase()}${family.slice(1)} ${major}${minor === undefined ? '' : `.${minor}`}`
}
