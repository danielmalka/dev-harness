---
type: regex
weight: 1
pattern: '^-\s*changes required\b'
flags: m
match: contains
target: last_message
---
Is the verdict `changes required`?
