---
type: regex
target: last_message
pattern: '## Findings[\s\S]{0,200}### (Critical|High|Medium|Low):'
flags: i
match: contains
weight: 2
---

The last message contains the `## Findings` heading from
`.skills/security-review/SKILL.md`'s Output format, followed by a
`### <SEVERITY>: <one-line title>` entry — the format `security` uses for a
blocking finding, distinct from `document`'s `## Verdict` /
`changes required` signal. The stage's CLI entry fixture double
(`fixtures/cli-output-codex-blocking.md`) reports one HIGH finding.
