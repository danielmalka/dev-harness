---
type: regex
target:
  source: file
  path: fixtures/brief.md
pattern: 'Stage:[^\n]{0,160}T-1'
flags: ''
match: contains
weight: 1
---

O `fixtures/brief.md` (a cópia gravável da fixture de
`.harness/tasks/AUTO-1/BRIEF.md`) tem uma linha `Stage:` que nomeia a
unidade interrompida, `T-1` (Fatia 1 do `fixtures/PLAN.md`), na mesma
linha ou no texto imediatamente seguinte — lido direto do disco, não da
transcrição. O RF-15/T-704 do PLAN-010 exigem que `Stage:` sempre nomeie a
próxima unidade (aqui, a interrompida); o placeholder de que este arquivo
parte (`Stage: not-started`) não menciona `T-1` e por isso não casa com
este padrão sozinho.
