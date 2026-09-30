---
type: regex
weight: 1
pattern: 'Blocking findings:[ ]*0\b'
flags: mi
match: contains
target: last_message
---
Is the finding count zero?
