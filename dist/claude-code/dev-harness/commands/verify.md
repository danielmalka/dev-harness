---
name: verify
description: Rode checagens de aceite e devolva uma matriz de evidência
author: malka
argument-hint: "[id da tarefa ou critérios de aceite]"
metadata:
  roles: [coordinator, qa-verifier]
  skills: [regression-testing, ui-verification]
  writes: somente testes no escopo combinado
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer coisa.

## Roteamento
Tarefa ou critérios: $ARGUMENTS. Despache `qa-verifier` com a ferramenta Agent, modelo `sonnet`, e exija que ele carregue a skill `regression-testing` do kit, mais `ui-verification` quando a mudança tem interface de usuário. Passe os critérios de aceite, a lista de arquivos tocados, os comandos de checagem de `.harness/project.yaml`, o escopo de escrita de testes autorizado e, para uma interface de usuário, como rodar a aplicação localmente.

## Pré-requisitos
Critérios de aceite e uma mudança estabilizada. Pare e reporte qual dos dois falta quando um estiver ausente; `/dh:discover` ou `/dh:plan` produzem critérios. Pergunte ao dono só para autorizar uma dependência de ambiente, como um banco descartável.

## Saída
O relatório no Formato de saída do `regression-testing`. Quando há interface de usuário, a matriz do `ui-verification` e as capturas de tela dela vão em `.harness/tasks/<id>/evidence/`. Registre o veredito e as lacunas restantes em `.harness/MEMORY.md`.

## Limites
- Edições de teste e de fixture acontecem só dentro do escopo de testes autorizado; os Limites do `qa-verifier` e as skills regem o resto.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:review` para a revisão independente, ou `/dh:fix` quando a matriz mostrar um defeito.
