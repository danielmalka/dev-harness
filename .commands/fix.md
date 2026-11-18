---
name: fix
description: Reproduza, diagnostique e corrija um defeito com evidência
author: malka
argument-hint: "[sintoma, texto do erro ou issue]"
metadata:
  roles: [coordinator, debugger, qa-verifier, reviewer]
  skills: [systematic-debugging, regression-testing, code-review]
  writes: só o conjunto de escrita autorizado
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` incluída no kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Sintoma: $ARGUMENTS. Em sequência: despache `debugger` com a ferramenta Agent, modelo `sonnet`, exigindo as skills `systematic-debugging` e `regression-testing` do kit, e diga se ele só diagnostica ou está autorizado a alterar arquivos, com o conjunto de escrita. Depois que uma correção for aplicada e o estado estiver congelado, despache `qa-verifier` (modelo `sonnet`, skill `regression-testing`). Em seguida despache `reviewer` no modo `code` (modelo `sonnet`, skill `code-review`) com escopo na correção, proporcional ao seu raio de impacto. No máximo dois especialistas rodam ao mesmo tempo, e só em trabalho independente.

## Pré-requisitos
O texto literal do erro ou o comportamento observado, e os passos de reprodução ou o comando que falha. Um erro parafraseado aponta para o componente errado: peça a saída exata. Consulte `.harness/RISKS.md` quando o sintoma tocar uma regra crítica ou uma área de incidente conhecida. Pergunte ao dono só pela autorização para alterar arquivos ou o ambiente.

## Saída
O relatório de `systematic-debugging`: status, sintoma, reprodução antes e depois, causa raiz marcada como fato ou hipótese com caminhos e linhas de evidência, pontos de chamada checados, correção ou "não aplicada", proteção contra regressão, causas descartadas e o relatório de incidente para o Coordenador. Grave a evidência em `.harness/tasks/<id>/evidence/`. Registre o incidente com severidade, causa, resolução e prevenção em `.harness/RISKS.md`, e o status atual em `.harness/MEMORY.md`.

## Limites
- Sem reprodução, reporte a lacuna e o caminho inconclusivo; nunca apresente uma correção comprovada.
- Um pedido de diagnóstico não autoriza edições de código-fonte, instrumentação, instalações nem mudanças em dados externos.
- Corrija a origem do valor ruim, não a linha onde ele apareceu; reporte cada ponto de chamada irmão checado.
- A proteção contra regressão só conta quando foi observada falhando antes da correção e passando depois; do contrário fica `not-run`, com o motivo.
- Nunca faça commit, push, publicação nem deploy a partir deste comando.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:verify` para ampliar a matriz de aceite, ou `/dh:release` quando a correção for a entrega.
