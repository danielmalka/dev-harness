import { expect, test } from 'claude-code/testing'
import { bar, commandDir, deniesAiMention, native, resolvePath, isGateCommand, isHarnessPath, parseGateCommands, parseLanguage, taskStatus } from './panel'

test('bar is proportional and clamped', () => {
  expect(bar(0)).toBe('▱▱▱▱▱▱▱▱▱▱')
  expect(bar(42)).toBe('▰▰▰▰▱▱▱▱▱▱')
  expect(bar(150)).toBe('▰▰▰▰▰▰▰▰▰▰')
})

test('only commit/PR text that mentions AI is denied', () => {
  expect(deniesAiMention('git commit -m "feat: x\n\nCo-Authored-By: Claude Opus"')).toBe(true)
  expect(deniesAiMention('gh pr create --body "🤖 Generated with [Claude Code]"')).toBe(true)
  expect(deniesAiMention('git commit -m "feat: add login"')).toBe(false)
  expect(deniesAiMention('grep -r "Claude Code" docs/')).toBe(false)
})

const YAML = `commands:
  test:
    value: go test ./...
    source: T05
  lint:
    value: golangci-lint run ./...
    source: T05
  run:
    value: null
    source: unknown
  fmt-check:
    value: "make fmt-check"
`

test('reads gate commands from the dh project.yaml', () => {
  expect(parseGateCommands(YAML)).toEqual(['go test ./...', 'golangci-lint run ./...', 'make fmt-check'])
})

test('recognizes generic and project gates', () => {
  expect(isGateCommand('cd x && uv run pytest -q', [])).toBe(true)
  expect(isGateCommand('make check', [])).toBe(true)
  expect(isGateCommand('./scripts/verify.sh', ['./scripts/verify.sh'])).toBe(true)
  expect(isGateCommand('git status', [])).toBe(false)
})

test('reads dh task status', () => {
  expect(taskStatus('| Status | concluída (commit a7b8977) |')).toBe('done')
  expect(taskStatus('| Status | bloqueada |')).toBe('blocked')
  expect(taskStatus('**Status: pronta para QA/revisão.**')).toBe('open')
  expect(taskStatus('**Status:** concluído.')).toBe('done')
  expect(taskStatus('sem tabela')).toBe('open')
})

test('language follows project.yaml, English by default', () => {
  expect(parseLanguage('language: pt-br\nname: x')).toBe('pt-br')
  expect(parseLanguage('name: x\nlanguage: "pt-br"')).toBe('pt-br')
  expect(parseLanguage('language: en')).toBe('en')
  expect(parseLanguage('name: x')).toBe('en')
})

test('.harness/ paths do not dirty the gate, other paths do', () => {
  expect(isHarnessPath('.harness/tasks/T-1/TASK.md')).toBe(true)
  expect(isHarnessPath('/home/u/proj/.harness/tasks/T-1/TASK.md')).toBe(true)
  expect(isHarnessPath('C:\\proj\\.harness\\x.md')).toBe(true)
  expect(isHarnessPath('src/main.go')).toBe(false)
  expect(isHarnessPath('src/.harnessrc')).toBe(false)
  expect(isHarnessPath(undefined)).toBe(false)
})

test('resolvePath normalizes relative, dot and dot-dot segments', () => {
  expect(resolvePath('/proj', 'src/a.go')).toBe('/proj/src/a.go')
  expect(resolvePath('/proj', '/x/./y/../z')).toBe('/x/z')
  expect(resolvePath('/proj/sub', '..')).toBe('/proj')
  // Windows: the engine's C:\\ path and git's C:/ root normalize to the same form
  expect(resolvePath('C:\\proj', 'src\\a.go')).toBe('/C:/proj/src/a.go')
  expect(resolvePath('/x', 'c:/repo/b.go')).toBe('/C:/repo/b.go')
  expect(resolvePath('/', 'C:/repo')).toBe('/C:/repo')
  expect(native('/C:/repo')).toBe('C:/repo')
  expect(native('/home/u')).toBe('/home/u')
})

test('commandDir finds the directory a commit runs in', () => {
  // `at` is where the commit starts, as guardBash passes it
  const d = (c: string) => commandDir(c, '/proj', '/home/u', c.search(/\bgit\s+(-C\s+("[^"]*"|'[^']*'|\S+)\s+)?commit/))
  expect(d('git commit -m x')).toBe('/proj')
  expect(d('cd /other && git commit -m x')).toBe('/other')
  expect(d('cd sub; git add . && git commit -m x')).toBe('/proj/sub')
  expect(d('cd "/a b" && git commit')).toBe('/a b')
  expect(d('cd ~/v && git commit')).toBe('/home/u/v')
  expect(d('git -C ../o commit -m x')).toBe('/o')
  expect(d("git -C '/q r' commit")).toBe('/q r')
  // a cd after the commit does not move it
  expect(d('cd a && cd b && git commit')).toBe('/proj/a/b')
  expect(d('cd && git commit')).toBe('/home/u')
  expect(d('cd /a && git -C b commit')).toBe('/a/b')
  // `git -C` inside a quoted message is not the commit's -C
  const q = 'git commit -m "git -C /zzz commit"'
  expect(commandDir(q, '/proj', '/h', q.indexOf('git commit'))).toBe('/proj')
  const q2 = 'echo "git -C /z commit"; git -C /o commit -m x'
  expect(commandDir(q2, '/proj', '/h', q2.lastIndexOf('git -C'))).toBe('/o')
  const after = 'cd /a && git commit -m x && cd /else'
  expect(commandDir(after, '/proj', '/h', after.indexOf('git commit'))).toBe('/a')
})

test('git -C commit is a commit for the AI-mention trap', () => {
  expect(deniesAiMention('git -C /o commit -m "x\n\nCo-Authored' + '-By: Cla' + 'ude"')).toBe(true)
})
