---
type: regex
pattern: '## Dispatches[\s\S]{0,800}cli:codex/gpt-5\.6-sol'
flags: ''
match: contains
target: last_message
weight: 1
---

The reply carries a `## Dispatches` section and logs the codex entry there with
its binary and slug rather than as a second Claude dispatch, per
`.agents/coordinator.md` ("Log a CLI reviewer in the Dispatches output with its
resolved binary and slug in place of the model fields").
