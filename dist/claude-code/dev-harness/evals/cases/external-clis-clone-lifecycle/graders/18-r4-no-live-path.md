---
type: llm
weight: 1
focus: last_message
---
In (r4), does the answer say the planning prompt does not contain the live tree's absolute path (it names the live tree as `<project>` or by relative path)? Quoting the path inside the grep check on the prompt file, such as `grep -F -- "<path>" <prompt file>` expecting no match, is allowed; the path appearing as prompt content is not.
