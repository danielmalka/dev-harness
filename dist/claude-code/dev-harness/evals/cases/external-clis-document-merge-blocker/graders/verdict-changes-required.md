---
type: regex
target: last_message
pattern: '## Verdict[\s\S]{0,80}changes required'
flags: ''
match: contains
weight: 2
---

The last message contains the `## Verdict` heading from
`.skills/document-review/SKILL.md`'s Output format, followed within a short
distance by a line starting with `changes required` (prefix match, not
exact equality) — the merged verdict must stay blocking because the CLI
entry's fixture double reports a blocking gap, even though the claude
entry's fixture double is a clean approve.
