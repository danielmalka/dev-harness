---
type: regex
target:
  source: file
  path: fixtures/memory.md
pattern: 'T-1'
flags: ''
match: contains
weight: 1
---

`fixtures/memory.md` also names the interrupted unit, `T-1`, read directly
from disk — mirroring grader 03's check on `fixtures/brief.md`'s `Stage:`
line. The plan requires the interrupted unit to be "named identically in
both `Stage:` and `.harness/MEMORY.md`, not silently dropped from either";
this grader and grader 03 each check one side of that pair mechanically
(a true cross-file identity check is outside what a single-target grader
can do, so both files are required to name the same fixed identifier, `T-1`,
rather than compared to each other directly).
