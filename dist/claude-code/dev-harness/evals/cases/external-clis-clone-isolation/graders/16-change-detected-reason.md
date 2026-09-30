---
type: regex
weight: 1
target: last_message
match: contains
pattern: '(change detected|mudança detectada)'
flags: ''
---
The not-run reason for the discarded scenarios is named with the kit's
wording, `change detected` (or `mudança detectada` in the owner's
language), as `live-tree-freeze.md` fixes it. Not `isolation unavailable`,
which is a different reason for a call that never started.
