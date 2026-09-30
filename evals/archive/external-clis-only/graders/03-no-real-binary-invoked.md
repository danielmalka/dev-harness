---
type: regex
pattern: 'cli-output-codex-approve\.md'
flags: ''
match: contains
target: trace
weight: 1
---

The codex verdict comes from reading
`fixtures/cli-output-codex-approve.md`, which appears in the trace. The case
grants no execution tool, so no real `codex` binary can have been invoked.
