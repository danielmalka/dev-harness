export type Panel = 'main' | 'zen' | 'mode' | 'status'

export type ZenTheme = 'blue' | 'dark' | 'light' | 'green' | 'rainbow' | 'custom'

export type ZenColors = { bar: string; text: string; border: string }

export type StepStatus = 'pending' | 'in_progress' | 'completed'

export type ZenStep = { id: string; text: string; status: StepStatus }

export type ZenTask = {
  title: string
  startedAt: number
  doneAt: number | null
  steps: ZenStep[]
}

export type Effort = 'low' | 'medium' | 'high' | 'xhigh' | 'max'

export type RateLimit = { label: string; percent: number; resetsAt: number | null }

// What the status line draws, refreshed after each model response and on a timer.
export type StatusSnapshot = {
  model: string
  dir: string
  branch: string
  ctxPercent: number | null
  ctxTokens: number | null
  ctxWindow: number | null
  limits: RateLimit[]
}

declare module 'claude-code' {
  interface PluginState {
    dh: {
      menu: Panel | null
      zenOn: boolean
      zenTheme: ZenTheme
      zenCustom: ZenColors
      zenCollapsed: boolean
      task: ZenTask | null
      now: number
      effort: Effort | null
      defaultMode: string | null
      statusOn: boolean
      status: StatusSnapshot | null
      dhSegments: string[]
      lastResponseAt: number | null
      statusClock: number
      sessionEffort: Effort | null
      sessionModel: string | null
    }
  }
}
