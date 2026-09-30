---
type: tool_used
tool: Agent
min: 0
max: 0
arm: both
weight: 1
---

Neither reviewer is dispatched for real: this case substitutes a fixture for
the `claude` entry ("read its reply from fixtures/claude-output-approve.md
instead of dispatching a real agent") and for the `cli:agy` call, so no
`Agent` dispatch may appear on any scored run.
