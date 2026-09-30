---
type: regex
weight: 1
pattern: '^- approved\b'
flags: m
match: contains
target: last_message
---
Is the verdict `approved`?
