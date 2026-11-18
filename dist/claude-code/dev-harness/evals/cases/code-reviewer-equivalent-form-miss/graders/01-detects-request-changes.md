---
type: regex
weight: 1
pattern: '^-\s*(Status da revisão:[ \t]*)?[Rr]equest changes\b'
flags: m
match: contains
target: last_message
---
O veredito é `request changes`, seja como a linha `Status da revisão: request changes` da skill, seja como a linha `- request changes` do próprio agente?
