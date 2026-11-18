---
name: handoff
description: Registre o estado transportável para a próxima sessão
author: malka
argument-hint: "[nota opcional para a próxima sessão]"
metadata:
  roles: [coordinator, docs-guide]
  skills: [context-handoff]
  writes: registro em .harness/tasks/<id>/; MEMORY.md só pelo Coordenador
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` incluída no kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Nota opcional: $ARGUMENTS. Carregue você mesmo a skill `context-handoff` do kit e seja o responsável pelo registro. Despache `docs-guide` com a ferramenta Agent, modelo `haiku`, exigindo a mesma skill, só para redigir o documento de handoff a partir do material que você fornecer; ele devolve a você o texto e as atualizações de memória propostas.

## Pré-requisitos
O objetivo e o escopo autorizado da tarefa; um handoff sem objetivo transfere atividade, não trabalho, então pergunte ao dono e pare. Os arquivos tocados, as decisões tomadas e a evidência produzida, registrados só onde puderem ser rastreados até a fonte. `.harness/MEMORY.md` deve existir; se faltar, inicialize-o a partir do template MEMORY incluído neste kit quando a configuração estiver autorizada; caso contrário, reporte o estado ausente.

## Saída
O registro de `context-handoff`: objetivo e autorização, restrições levadas adiante literalmente, conjunto de escrita, decisões com motivo, fonte e condição para rever, tabela de evidência com cada checagem `passed`, `failed` ou `not-run`, itens pendentes e bloqueados com responsável e condição de conclusão, e o único próximo passo autorizado, marcado como já concedido ou pendente. Grave-o em `.harness/tasks/<id>/`; você atualiza `.harness/MEMORY.md` com o mesmo estado.

## Limites
- Nenhuma mudança de código: este comando registra estado e não escreve nada além do documento de handoff e do registro de memória.
- Leve adiante literalmente uma restrição com fonte; uma restrição registrada de procedência desconhecida continua em aberto em vez de ser descartada.
- Nunca invente estado passado, autorização concedida ou resultado de checagem; itens não verificados são marcados como hipótese.
- `docs-guide` redige, mas nunca edita os arquivos de memória.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, e `MEMORY.md` só pela trava de escrita da memória (`context-handoff`).

## Próximo
`/dh:resume` na próxima sessão, ou `/dh:consolidate-memory` quando o registro de memória tiver crescido além do que a execução precisa.
