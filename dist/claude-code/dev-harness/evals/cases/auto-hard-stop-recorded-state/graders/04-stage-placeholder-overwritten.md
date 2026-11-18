---
type: regex
target:
  source: file
  path: fixtures/brief.md
pattern: 'Stage:[ \t]*not-started'
flags: ''
match: not_contains
weight: 1
---

A linha de placeholder `Stage: not-started` do `fixtures/brief.md` não
existe mais depois da execução — lido direto do disco. Uma execução que
nunca toca neste arquivo (deixando o placeholder intocado) reprova neste
grader mesmo que a resposta final narre uma parada no chat, conforme a
exigência do próprio plano de que "uma menção só no chat a 'parando' ou
'pausado', sem a gravação correspondente de `Stage:`, reprova neste
grader".
