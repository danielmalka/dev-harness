---
name: improve
description: Diagnostique e avalie uma melhoria deste kit
author: malka
argument-hint: "[caso reproduzível: prompt, esperado, observado]"
metadata:
  roles: [coordinator, harness-maintainer, qa-verifier]
  skills: [harness-authoring, harness-evaluation]
  writes: só as fontes do checkout do harness, nunca dist/
disable-model-invocation: true
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` incluída no kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Caso: $ARGUMENTS. Despache `harness-maintainer` com a ferramenta Agent, modelo `sonnet`, e exija que ele carregue as skills `harness-authoring` do kit para o diagnóstico e a edição e `harness-evaluation` para a medição de antes e depois. Instrua-o a nomear o ativo responsável antes de editar qualquer coisa: qual gatilho disparou, qual arquivo foi carregado, qual instrução produziu a saída.

## Pré-requisitos
Um caso concreto com prompt, comportamento esperado e comportamento observado; para um reparo, uma falha reproduzida ou uma contradição de contrato citada entre duas instruções. Uma resposta preferida isolada não é evidência, então peça a reprodução ao dono e pare. Um checkout do harness para trabalhar, com o estado registrado, e um escopo autorizado claro antes de qualquer edição. Os casos ficam em `evals/cases/` com suas rubricas e fixtures; uma linha de base deve ser capturada antes de o ativo mudar.

## Saída
O relatório de mudança de ativo de `harness-authoring` (caso, ativo responsável, menor edição, o que deliberadamente não mudou, checagem de contrato, casos de teste positivos e negativos, procedência) junto com o relatório de antes e depois de `harness-evaluation`, com suas condições, checagens estáticas, resultados avaliados por modelo e variância. Registre a mudança, o efeito medido e seus limites em `.harness/MEMORY.md`.

## Limites
- Aplique mudanças só com um escopo autorizado claro no checkout do harness; sem ele, proponha a edição e pare.
- Nunca edite `dist/` à mão; ele é gerado a partir dos ativos-fonte.
- Sem um executor disponível, os critérios avaliados por modelo ficam `not-run` e o resultado é uma correção estática, nunca uma melhoria medida; um ativo alterado nunca é a sua própria linha de base de antes.
- Não incorpore material emprestado; nomear uma fonte não torna o conteúdo dela disponível.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:doctor` para confirmar que o kit ainda resolve, depois rode de novo o caso que falhava para confirmar a correção.
