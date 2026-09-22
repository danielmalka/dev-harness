---
type: tool_used
tool: Agent
input_match: document-validator
min: 0
max: 0
arm: both
weight: 2
---

No `document-validator` (this stage's Claude reviewer) is ever dispatched
through the `Agent` tool. `reviewers.document` in
`fixtures/project.yaml` names no `claude` entry, so per
`.agents/coordinator.md` ("External CLI reviewers"): "When a stage lists no
`claude`, no Claude agent is dispatched for it and the merged report comes
only from the CLI verdicts."
