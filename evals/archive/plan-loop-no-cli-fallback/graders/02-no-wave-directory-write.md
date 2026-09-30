---
type: regex
target: trace
pattern: '"file_path":[ \t]*"[^"]*wave-[0-9]+'
flags: ''
match: not_contains
weight: 2
---

No `Write` or `Edit` tool call in the trace targets a path containing a
`wave-N` segment (`wave-1`, `wave-2`, ...), anywhere in the run.
`.commands/plan-loop.md` Routing item 1 states "No `wave-N/` directory is
created on this path" for exactly the condition this case's fixture
produces (zero resolvable planner entries after the AA1 exclusion). The
pattern is anchored to the serialized `file_path` key both tools share,
not a bare `wave-` substring, so a mention of the word "wave" inside a
file's own written content (for example, a plan slice that happens to
discuss wave-based loops in prose) cannot satisfy the check by accident.
