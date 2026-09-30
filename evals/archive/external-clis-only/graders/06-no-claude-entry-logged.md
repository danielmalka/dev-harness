---
type: llm
weight: 1
focus: last_message
---
Is the Dispatches section free of any `claude` reviewer entry for this stage — no row and no bullet presenting `claude` as a reviewer that ran? Prose explaining that the stage's list names no `claude` entry, so none was dispatched, is not an entry and does not count against this.

Judged rather than matched: the criterion is about what the log presents, and
the Dispatches format is prose (a bullet list in `.agents/coordinator.md`'s
Output format, a table in some runs), so a pattern anchored to either shape
passes vacuously on the other. The hard guarantee that no Claude reviewer ran
is mechanical and lives in `01-claude-left-undispatched` (`Agent` called 0
times); this criterion covers the separate failure of logging a reviewer that
was never dispatched.
