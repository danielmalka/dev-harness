---
type: regex
target: last_message
pattern: '## Verdict[\s\S]{0,80}approved'
flags: ''
match: contains
weight: 2
---

The last message contains the `## Verdict` heading from
`.skills/document-review/SKILL.md`'s Output format, followed within a short
distance by a line starting with `approved` (prefix match, not exact
equality) — the stage has exactly one reviewer
(`cli:codex/gpt-5.6-sol`) and its fixture double
(`fixtures/cli-output-codex-approve.md`) reports zero findings.
