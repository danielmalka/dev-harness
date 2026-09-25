---
type: regex
weight: 1
pattern: 'Blocking findings:[ \t]*[1-9]'
flags: m
match: contains
target: last_message
---
Does the reply report at least one blocking finding?
