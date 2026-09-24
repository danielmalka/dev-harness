---
type: tool_used
tool: Write
min: 0
max: 0
arm: both
weight: 1
---
No `Write` tool call occurs on any scored run, even though `Write` is available in `allowed_tools`: this is a plan-and-ask turn, not an implementation turn, so no file is created or overwritten before the owner answers the open preference.
