---
type: regex
target:
  source: file
  path: fixtures/brief.md
pattern: '7\.25'
flags: ''
match: contains
weight: 1
---

`fixtures/brief.md` (the fixture's writable copy of
`.harness/tasks/AUTO-1/BRIEF.md`) contains the literal eval spend cap
`7.25` after the run — the specific figure the user's message gave as
answer (g), present nowhere else in the fixture or the prompt except as
that one answer. Its presence in the persisted file, read directly from
disk rather than from the transcript, is what proves the questionnaire's
actual answers reached the record — the placeholder stub this file starts
from (`Stage: not-started`, no digits at all) cannot match this pattern by
itself, so an unmodified fixture fails this grader.
