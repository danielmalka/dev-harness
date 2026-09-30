---
type: tool_used
tool: Agent
min: 0
max: 0
arm: both
weight: 1
---

No agent is dispatched for the stage: `reviewers.code` names only
`cli:codex/gpt-5.6-sol`, and per `.agents/coordinator.md` ("External CLI
reviewers") "when a stage lists no `claude`, no Claude agent is dispatched for
it and the merged report comes only from the CLI verdicts". Asserted on the
trace, because a reply saying it dispatched nothing is not proof that it did
not.
