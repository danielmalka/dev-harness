---
type: regex
target: last_message
pattern: 'Reported by: \s*cli:codex/gpt-5\.6-sol'
flags: ''
match: contains
weight: 2
---

The blocking finding in the last message carries `Reported by:
cli:codex/gpt-5.6-sol`, per `.skills/security-review/SKILL.md`'s
`Reported by: claude | cli:<binary>/<slug>` field and
`.agents/coordinator.md` ("External CLI reviewers"): "Each finding in that
record carries the tag of the reviewer that produced it."
