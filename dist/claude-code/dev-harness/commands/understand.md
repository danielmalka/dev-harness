---
name: understand
description: Mapeie o código por trás de uma pergunta, somente leitura
author: malka
argument-hint: "[pergunta ou área a mapear]"
metadata:
  roles: [coordinator, repo-scout]
  skills: [repository-mapping]
  writes: nenhum (somente leitura)
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` incluída no kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Pergunta ou área: $ARGUMENTS. Despache `repo-scout` com a ferramenta Agent, modelo `haiku`, e exija que ele carregue a skill `repository-mapping` do kit. Passe a pergunta formulada como comportamento, o diretório-alvo e os comandos de checagem de `.harness/project.yaml` quando ele existir.

## Pré-requisitos
Uma pergunta formulada como comportamento. Pare e pergunte ao dono o que deve ser localizado quando $ARGUMENTS não nomear nenhum comportamento.

## Saída
O relatório no Formato de saída de `repository-mapping`, devolvido na resposta; grave-o em `.harness/tasks/<id>/` só quando uma tarefa já for dona dele. Registre a área mapeada e seus caminhos de evidência em `.harness/MEMORY.md`.

## Limites
- Somente leitura: nenhum arquivo do projeto é criado nem editado.
- Mapear não autoriza implementação; não inicie uma a partir deste comando.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:plan` quando o trabalho estiver entendido, ou `/dh:discover` quando o comportamento em si ainda não estiver decidido.
