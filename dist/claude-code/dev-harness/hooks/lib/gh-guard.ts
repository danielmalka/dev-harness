// PRD-027 2.3: a short rule on the command text, no parser. No I/O.
// ponytail: accepted false positives and bypasses are listed in PRD 2.3 and 2.4. Ceiling: any text that names
// `gh` or a git write word is denied unless it is a plain gh read or a plain grep; the main session runs the rest.

// Word boundaries (PRD 2.3): a word starts at the start or after a char that is not a letter, digit or `_`.
const START = '(?<=^|[^A-Za-z0-9_])'
const word = (w: string, end = '(?=$|[^A-Za-z0-9_-])') => new RegExp(`${START}${w}${end}`)
const GH = word('gh', '(?=$|\\.exe|[^A-Za-z0-9_-])') // gh also ends before `.exe`, but not before `-` (gh-pages)
const GIT = word('git', '(?![A-Za-z0-9_])') // git ends before `-` too: git-push, git-lfs push
const WRITE = ['push', 'send-pack', 'subtree', 'p4', 'svn', 'dcommit', 'alias'].map(w => ({ w, re: word(w) }))

// gh reads (PRD 2.2): group -> allowed verbs.
const READ: Record<string, string[]> = {
  pr: ['view', 'diff', 'list', 'status', 'checks'], issue: ['view', 'list', 'status'], repo: ['view', 'list'],
  run: ['view', 'list', 'watch'], release: ['view', 'list'], workflow: ['view', 'list'],
}
const SAFE = /^[A-Za-z0-9 ._/:=,@%+-]+$/
const API_FLAG = /^(--paginate|-i|--include|--silent)$/ // flags without a value
const API_VALUE = /^(--jq|-q|--cache)$/ // flags with a value, consumed with the next token
const API_EQ = /^(--jq|--cache)=/
const isR = (x: string | undefined) => x === '-R' || x === '--repo'

// Left to right (PRD 2.2): a value flag consumes the next token whatever it looks like; a value flag without one is not a read.
function isGhRead(raw: string): boolean {
  const t = raw.trim().split(/\s+/)
  if (!/(^|\/)gh(\.exe)?$/.test(t[0]!) || !SAFE.test(raw.trim())) return false
  let i = 1
  while (isR(t[i])) { if (i + 1 >= t.length) return false; i += 2 }
  const group = t[i]
  if (!group || group === '--version' || group === 'status' || group === 'search') return true
  const rest = t.slice(i + 1)
  if (group === 'api' || group === 'auth') {
    // strict: only the flags of the table; -R / --repo is not one of them
    if (group === 'auth' ? rest[0] !== 'status' : rest.some(x => x.includes('://'))) return false // a full URL can send the token away
    for (let j = group === 'auth' ? 1 : 0; j < rest.length; j++) {
      const x = rest[j]!
      const value = group === 'auth' ? /^(-h|--hostname)$/ : API_VALUE
      if (value.test(x)) { if (++j >= rest.length) return false }
      else if (!x.startsWith('-')) continue
      else if (group === 'auth' ? !/^(-a|--active|(-h|--hostname)=.*)$/.test(x) : !API_FLAG.test(x) && !API_EQ.test(x)) return false
    }
    return true
  }
  const a: string[] = []
  for (let j = 0; j < rest.length; j++) {
    if (isR(rest[j])) { if (++j >= rest.length) return false } else a.push(rest[j]!)
  }
  return !!READ[group] && (!a[0] || READ[group]!.includes(a[0]))
}

// The name for the deny message: `gh pr merge`, `git push`, or just `gh` / `git`.
function verbOf(text: string, ghHit: boolean): string {
  if (!ghHit) return `git ${WRITE.find(x => x.re.test(text))?.w ?? ''}`.trim()
  const words = text.slice(text.search(GH) + 2).split(/\s+/).filter(x => /^[a-z][\w-]*$/.test(x))
  return ['gh', ...(words[0] === 'api' ? ['api'] : words.slice(0, 2))].join(' ')
}

// The matched verb ('gh pr merge', 'git push', 'gh'), or undefined when the command is not denied.
// `powershell`: the backtick is an ignored escape there, so it is removed and no exception applies (PRD 2.3).
export function ghWriteVerb(command: string, powershell = false): string | undefined {
  const lc = command.toLowerCase()
  const cur = lc.replace(powershell ? /\\/g : /\\\\/g, '/').replace(powershell ? /['"`]/g : /[\\'"]/g, '')
  // Bash: also the text with every backslash removed, since a backtick substitution drops one backslash layer first.
  const texts = powershell ? [cur] : [lc.replace(/[\\'"]/g, ''), cur]
  const hit = (t: string) => GIT.test(t) && WRITE.some(x => x.re.test(t))
  const gitHit = texts.some(hit)
  const text = texts.find(t => (gitHit ? hit(t) : GH.test(t)))
  if (text === undefined) return undefined
  const ghHit = GH.test(text)
  if (powershell) return verbOf(text, ghHit)
  if (!gitHit && isGhRead(command)) return undefined
  if (/^\s*(grep|egrep|fgrep)\s/.test(command) && !/[;|&$()<>`\n]/.test(command)) return undefined
  return verbOf(text, ghHit)
}

// PRD-027 R6: who may run it, and no retry, no workaround. `unknown` = the agentId has no recorded type.
export const GH_MSG = {
  en: (verb: string, unknown: boolean) => `dev-harness: \`${verb}\` is a remote write.${unknown ? ' Denied because the subagent type is unknown (unknown subagent type).' : ''} Only release-manager or the main session may run it. Do not retry and do not work around it; report it to the Coordinator. For a text search, use the Grep tool.`,
  'pt-br': (verb: string, unknown: boolean) => `dev-harness: \`${verb}\` é uma escrita remota.${unknown ? ' Negado porque o tipo do subagente é desconhecido (unknown subagent type).' : ''} Só o release-manager ou a sessão principal pode rodá-lo. Não tente de novo nem contorne; avise o Coordinator. Para busca de texto, use a ferramenta Grep.`,
}

// The deny for a Bash call, or undefined. Kit subagents (dh:) other than release-manager and unknown ids are
// denied on a match; the main session (no agentId) and non-kit subagents are not.
export function ghDenied(command: string, agentId: unknown, type: string | undefined, unreadable = false, powershell = false): { verb: string; unknown: boolean } | undefined {
  if (!agentId) return undefined
  if (type !== undefined && (type === 'dh:release-manager' || !type.startsWith('dh:'))) return undefined
  const verb = unreadable ? 'unreadable command' : ghWriteVerb(command, powershell)
  return verb ? { verb, unknown: type === undefined } : undefined
}
