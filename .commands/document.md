---
name: document
description: Escreva a documentação necessária para usar e operar uma mudança
author: malka
argument-hint: "[assunto ou área a documentar]"
metadata:
  roles: [coordinator, docs-guide, document-validator]
  skills: [delivery-readiness, context-handoff, document-review]
  writes: arquivos de documentação no escopo autorizado
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Assunto: $ARGUMENTS. Decida primeiro se a mudança altera como alguém usa, roda ou continua o trabalho; quando não altera, diga isso e pare. Caso contrário, despache `docs-guide` com a ferramenta Agent, modelo `haiku`, exigindo a skill do kit `delivery-readiness` para as seções operacionais e de runbook ou `context-handoff` para documentos de continuidade. Para um documento HTML em pt-BR no padrão do kit, exija também `doc-template-html`; esse papel não tem shell, então instrua-o a copiar o template em vez de rodar o script de carimbo. Assunto `pdocs` (o HTML de visão geral do projeto): (1) rode `dh harness-path --json` com o `bin/dh` que acompanha o kit (`bin\dh.cmd` no Windows), como faz `/dh:doctor`, e leia `pdocs`; (2) quando `pdocs` está vazio, pare com "rode dh link"; no modo global exija também o setup feito (`project.yaml` presente na pasta do harness), caso contrário pare e aponte para `/dh:setup`; (3) nos dois modos lembre o dono de que `~/.harness` deve estar em `permissions.additionalDirectories` de `~/.claude/settings.json` (nunca escreva esse arquivo); (4) despache `docs-guide`, modelo `haiku`, exigindo `doc-template-html`, tipo `projeto`, para criar `<pdocs>/index.html` copiando `assets/projeto-modelo.html` (sem shell) ou para atualizar o existente; páginas extras ficam planas na mesma pasta, em kebab-case, com link a partir de `index.html`; fontes: `docs/prd/` do repositório e `prd/` da pasta do harness para a tabela "Planejado e desenvolvido" (`id`, título, `Status`, data de geração), `CHANGELOG.md` para os títulos de release (copiados, nunca reescritos), pastas e entradas de primeiro nível para o diagrama macro obrigatório; uma lacuna continua "ainda não fechado", nunca texto inventado; (5) sem ciclo de validação (não é PRD, TASK nem ADR). Instrua `docs-guide` a substituir toda linha e todo fato do modelo pelos fatos do projeto de destino e a manter `ainda não fechado` onde for desconhecido.
Suba o modelo para `sonnet` só quando o assunto abranger vários subsistemas, e registre o motivo.
Um PR de entrega abre com o PRD ou resumo, as decisões e qualquer ADR, e fecha com a documentação do projeto atualizada; o Coordenador nunca abre um PR só de documentação; ao fim de todo lote, seja qual for o comando que o fecha, o Coordenador despacha docs-guide antes da entrega. Esse despacho de `docs-guide` acontece ao fim de um lote de implementação independentemente de haver um documento de template validado (PRD, TASK, ADR) envolvido — uma atualização só de tutorial ou só de README também conta.
O ciclo de validação abaixo se aplica apenas quando o documento persistido é um destes três: `templates/PRD.md`, `templates/TASK.md`, `templates/ADR.md`. Qualquer outra saída deste comando — um tutorial, um runbook, um documento interno HTML autônomo, um handoff — não o dispara. Quando se aplica: depois que `docs-guide` persistir o documento, despache `document-validator`, modelo `opus`, exigindo a skill do kit `document-review`, com o documento como artefato em revisão e o seu material-fonte (arquivos, critérios de aceite, evidência de QA, decisões) como fonte. O validador é somente leitura e nunca edita o documento; o despacho fica exclusivamente com o Coordenador. Persista o relatório que ele devolve como `<document>.review.md`. Em `changes required`, despache `docs-guide` de novo só com os pontos bloqueantes listados e valide outra vez. Limite: duas rodadas de correção; a lista completa só roda quando o dono pede o painel neste pedido ou o documento toca segurança; o dono recebe o documento só depois de `approved` ou do limite de rodadas, com qualquer ponto deixado em aberto.

## Pré-requisitos
O assunto e o material-fonte que ele documenta: arquivos, critérios de aceite, evidência de QA, decisões. Sem material-fonte, reporte a lacuna em vez de escrever por inferência. O caminho de escrita autorizado para o documento; sem ele, `docs-guide` devolve o texto na resposta e nada é escrito. Comandos e caminhos citados na documentação vêm de `.harness/project.yaml` ou da árvore; um comando não confirmado é marcado como não verificado.

## Saída
O documento no formato que a sua skill define, persistido no caminho autorizado sob `.harness/` (para `pdocs` o caminho é o absoluto de `dh harness-path --json`, nunca digitado nem inferido; por exemplo `.harness/tasks/<id>/` para a documentação de uma tarefa), com o comportamento planejado e o comportamento implementado claramente distinguidos e toda checagem reportada como `passed`, `failed` ou `not-run`. `docs-guide` reporta de volta a você as atualizações de memória propostas; você as registra em `.harness/MEMORY.md`.

## Limites
- `docs-guide` nunca edita `.harness/MEMORY.md`, `EPOCHAL.md` nem `RISKS.md`; ele reporta as atualizações ao Coordenador.
- Nunca documente um comportamento como implementado sem evidência de que ele existe; o planejado continua rotulado como planejado.
- Este comando não edita arquivo de código-fonte do produto, de configuração nem de teste; use só caminhos relativos, exceto o caminho de `pdocs` acima.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:release` quando a documentação completa um pacote de entrega, ou `/dh:handoff` para encerrar a sessão.
