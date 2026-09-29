---
type: tool_used
tool: Agent
input_match: '"subagent_type":[ \t]*"(dh:)?implementation-planner"'
min: 1
max: 3
weight: 2
---

`implementation-planner` is dispatched through the `Agent` tool at least
once (the fallback still needs a plan drafted) and at most three times —
`.commands/plan.md:16`'s own two-correction-round cap (the fallback's real
source of truth, per `.commands/plan-loop.md` Routing item 1) allows up to
two additional redispatches of the planner on `changes required`, on top
of the first draft, so "exactly once" would wrongly fail a legitimate
correction round. `input_match` is anchored to the serialized
`subagent_type` key rather than a bare `implementation-planner` substring,
so a dispatch of a different role whose context merely mentions
"implementation-planner" in passing cannot inflate the count. The optional
`(dh:)?` prefix accounts for the runtime serializing the plugin-qualified
name (`dh:implementation-planner`) rather than the bare role name observed
in run 1's trace (Correction 1: `plan-loop-no-cli-fallback/README.md`).

Concurrency is not independently checkable by this grader type — this
kit's eval schema has no concurrency primitive on the trace (`tool_used`
counts and orders calls, it does not time them). This case's own prompt
instructs the assistant to stop once the first validator verdict is in
hand, without a correction redispatch, so in practice exactly one
`implementation-planner` dispatch is exercised and no two calls ever have
the chance to overlap; the 1-3 bound stays the mechanically checked
acceptance shape for a run that does redispatch.
