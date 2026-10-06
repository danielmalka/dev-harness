import { expect, test } from 'claude-code/testing'
import { bar, deniesAiMention, isGateCommand, isHarnessPath, parseGateCommands, parseLanguage, taskStatus } from './panel'

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
