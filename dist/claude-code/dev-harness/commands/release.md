---
name: release
description: Monte um pacote de prontidão de release para decisão
author: malka
argument-hint: "[mudança, fatia ou versão]"
metadata:
  roles: [coordinator, release-manager, builder]
  skills: [delivery-readiness]
  writes: somente o registro de prontidão e as notas; nenhuma ação remota
disable-model-invocation: true
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer coisa.

## Roteamento
Mudança ou versão: $ARGUMENTS. Despache `release-manager` com a ferramenta Agent, modelo `sonnet`, e exija que ele carregue a skill `delivery-readiness` do kit. Quando o pacote precisar de evidência de pipeline ou de empacotamento a partir de um checkout limpo, despache também `builder` na lane `infra`, modo `design` (modelo `sonnet`, mesma skill), sem escrita de produto, para reproduzir o build a partir de um checkout novo; ele não instala nada, não provisiona nada e não alcança nenhum servidor ativo. Os dois trabalham sobre o mesmo estado congelado.

## Pré-requisitos
A mudança verificada, o resultado do QA dela e a revisão independente dela. Pare e reporte a lacuna quando o QA faltar (`/dh:verify` o produz) ou a revisão faltar (`/dh:review` a produz); toda outra entrada ausente segue a tabela de Entradas do `delivery-readiness`.

## Saída
O pacote no Formato de saída do `delivery-readiness`, persistido em `.harness/tasks/<id>/`. Registre o veredito, os riscos residuais e a decisão pendente em `.harness/MEMORY.md`.

## Limites
- Nenhum commit, push, tag, publicação ou deploy roda a partir deste comando; os Limites do `release-manager` regem uma autorização de publicação existente.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:handoff` quando o dono tiver o pacote, ou `/dh:fix` quando restar um achado bloqueante.
