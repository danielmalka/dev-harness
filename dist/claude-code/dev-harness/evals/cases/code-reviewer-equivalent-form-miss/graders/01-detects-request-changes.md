---
type: regex
weight: 1
pattern: '^-\s*(Review status:[ \t]*)?[Rr]equest changes\b'
flags: m
match: contains
target: last_message
---
Is the verdict `request changes`, either as the skill's `Review status: request changes` line or the agent's own `- request changes` line?
