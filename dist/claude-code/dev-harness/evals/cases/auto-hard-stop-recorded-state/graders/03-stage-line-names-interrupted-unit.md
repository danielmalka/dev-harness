---
type: regex
target:
  source: file
  path: fixtures/brief.md
pattern: 'Stage:[^\n]{0,160}T-1'
flags: ''
match: contains
weight: 1
---

`fixtures/brief.md` (the fixture's writable copy of
`.harness/tasks/AUTO-1/BRIEF.md`) has a `Stage:` line that names the
interrupted unit, `T-1` (Slice 1 of `fixtures/PLAN.md`), within the same
line or the text immediately following it — read directly from disk, not
from the transcript. RF-15/PLAN-010 T-704 require `Stage:` to always name
the next (here, interrupted) unit; the placeholder this file starts from
(`Stage: not-started`) does not mention `T-1` and so cannot match this
pattern by itself.
