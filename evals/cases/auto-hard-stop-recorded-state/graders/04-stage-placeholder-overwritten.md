---
type: regex
target:
  source: file
  path: fixtures/brief.md
pattern: 'Stage:[ \t]*not-started'
flags: ''
match: not_contains
weight: 1
---

`fixtures/brief.md`'s placeholder `Stage: not-started` line is gone after
the run — read directly from disk. A run that never touches this file
(leaving the placeholder untouched) fails this grader even if the final
reply narrates a stop in chat, matching the plan's own requirement that "a
chat-only mention of 'stopping' or 'paused' with no matching `Stage:`
write fails this grader."
