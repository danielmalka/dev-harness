---
type: regex
weight: 1
pattern: '^[ \t]*(?:[-*][ \t]+(?:\*\*)?(?:Status da revisão:(?:\*\*)?[ \t]*(?:\*\*)?)?|(?:\*\*)?Status da revisão:(?:\*\*)?[ \t]*(?:\*\*)?)[Rr]equest changes\b'
flags: m
match: contains
target: last_message
---
O veredito é `request changes`, seja na linha `Status da revisão: request changes` da skill (com ou sem marcador de lista, com ou sem negrito), seja na linha `- request changes` do próprio agente (também com ou sem negrito)? Um `request changes` solto no texto, sem marcador de lista e sem o rótulo, não conta.
