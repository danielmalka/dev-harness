---
type: regex
target:
  source: file
  path: fixtures/memory.md
pattern: '7\.25'
flags: ''
match: contains
weight: 1
---

`fixtures/memory.md` (the fixture's writable copy of `.harness/MEMORY.md`)
also contains the literal eval spend cap `7.25` after the run, mirroring
grader 04 for the second record `.commands/auto.md` requires: "the same
summary to `.harness/MEMORY.md`". The bundled template this file starts
from ("No execution recorded.") has no digits at all, so an unmodified
fixture fails this grader; only a reply that actually persists the
questionnaire's answers into this file, not only into `brief.md`, passes
both graders 04 and 05 together.
