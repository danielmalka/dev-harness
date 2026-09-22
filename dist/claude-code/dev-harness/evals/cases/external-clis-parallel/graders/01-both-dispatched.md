---
type: tool_used
tool: Agent
input_match: code-reviewer
min: 1
weight: 1
---

The `claude` entry of the `code` stage is dispatched as a real
`code-reviewer` agent call, per `.agents/coordinator.md` ("External CLI
reviewers"): "Dispatch every entry of the stage's list. `claude` entries are
dispatched as today." Only the `cli:codex/gpt-5.6-sol` entry is
fixture-substituted by this case's prompt, so the claude entry must appear in
the trace as an actual `Agent` call and not merely in the Dispatches log.
