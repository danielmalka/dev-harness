---
type: regex
target:
  source: file
  path: fixtures/memory.md
pattern: 'T-1'
flags: ''
match: contains
weight: 1
---

O `fixtures/memory.md` também nomeia a unidade interrompida, `T-1`, lido
direto do disco — espelhando a checagem do grader 03 sobre a linha
`Stage:` do `fixtures/brief.md`. O plano exige que a unidade interrompida
seja "nomeada de forma idêntica em `Stage:` e em `.harness/MEMORY.md`, sem
ser omitida em silêncio de nenhum dos dois"; este grader e o grader 03
checam, cada um, um lado desse par de forma mecânica (uma checagem de
identidade entre os arquivos de verdade está além do que um grader de alvo
único consegue fazer, então os dois arquivos precisam nomear o mesmo
identificador fixo, `T-1`, em vez de serem comparados diretamente entre si).
