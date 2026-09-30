---
type: tool_order
before:
  tool: Edit
  input_match: '"file_path":[ \t]*"[^"]*memory\.md"'
after:
  tool: Agent
weight: 1
---

The call that persists the same summary to `fixtures/memory.md` (standing
in for `.harness/MEMORY.md`) occurs before any `Agent` call that dispatches
the first chained stage, mirroring grader 02 for the second of the two
records `.commands/auto.md` requires "before dispatching the first stage".
`fixtures/memory.md` is an *existing* record (the fixture seeds it with
the real bundled `MEMORY.md` template, not an absent file), so the correct
tool for the Coordinator to use here is `Edit`, not `Write` — this grader
checks `Edit`, not `Write`, to match that. `input_match` is anchored to
the serialized `file_path` key (`"file_path": "fixtures/memory.md"`)
rather than a bare `memory\.md` substring, so a `Write` to `brief.md`
whose *content* happens to mention `memory.md` — the false-positive this
grader previously scored on, since `tool_order`'s `input_match` matches
against the whole serialized tool input, including `content` — cannot
satisfy it by accident. Both records must land before dispatch; this
grader closes the gap grader 02 leaves if only `BRIEF.md` were checked.
