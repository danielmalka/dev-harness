---
type: tool_used
tool: Agent
min: 1
max: 1
arm: both
weight: 1
---

Exactly one agent is dispatched for the stage: the `claude` entry. The
`cli:codex/gpt-5.6-sol` entry is a CLI call whose reply this case substitutes
from a fixture, so promoting it to a second `Agent` dispatch would be the
same misattribution `.agents/coordinator.md` forbids ("Only a reviewer you
actually dispatched is logged as one").
