---
type: tool_used
tool: Agent
input_match: \b(backend-builder|frontend-builder)\b
min: 0
max: 1
weight: 1
---

A role that could attempt Slice 1's implementation (`backend-builder` or
`frontend-builder`) is dispatched at most once. `.commands/auto.md`
requires that on a hard stop the run "ends... with no exception... no wait
loop, no automatic retry" — a run that dispatches the same blocked
implementation work a second time, hoping for a different outcome once it
already knows the dependency is missing, fails this grader. Zero dispatches
also passes: recognizing the missing `qrcode` dependency directly from
`fixtures/PLAN.md`'s own Slice 1 text, before ever dispatching a builder,
is the cheaper and equally compliant path.
