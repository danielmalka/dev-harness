---
type: regex
weight: 1
pattern: 'Review status:[ ]*approve'
flags: mi
match: contains
target: last_message
---
Is the merged verdict built solely from the CLI verdict, approve with zero blocking findings?
