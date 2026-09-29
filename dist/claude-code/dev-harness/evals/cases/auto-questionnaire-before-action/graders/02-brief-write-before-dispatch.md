---
type: tool_order
before:
  tool: Write
  input_match: '"file_path":[ \t]*"[^"]*brief\.md"'
after:
  tool: Agent
weight: 1
---

The `Write` call that persists the questionnaire record to
`fixtures/brief.md` (standing in for `.harness/tasks/AUTO-1/BRIEF.md`)
occurs before any `Agent` call that dispatches the first chained stage.
`input_match` is anchored to the serialized `file_path` key
(`"file_path": "fixtures/brief.md"`) rather than a bare `brief\.md`
substring, so a `Write` to some other file whose *content* happens to
mention `brief.md` (for example while writing `memory.md`'s summary text)
cannot satisfy this grader by accident. `.commands/auto.md` requires the
answers to be persisted "before dispatching the first stage" and "never
write a file or dispatch a specialist before this questionnaire is
answered" — a run that dispatches a stage first and only writes the record
afterward, or never writes it at all, fails this grader on the trace
alone, independent of what the final reply claims.
