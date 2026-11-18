---
name: discover
description: Transforme um pedido vago em um resumo com aceite
author: malka
argument-hint: "[ideia, dor ou pedido]"
metadata:
  roles: [coordinator, product-discovery, document-validator]
  skills: [requirements-discovery, document-review, external-clis]
  writes: .harness/prd/ quando autorizado
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Ideia ou dor: $ARGUMENTS. Despache `product-discovery` com a ferramenta Agent, modelo `sonnet`, e exija que ele carregue a skill do kit `requirements-discovery`. Instrua-o a buscar nos documentos de produto existentes e em qualquer resumo anterior do projeto antes de perguntar qualquer coisa, e a levantar uma pergunta bloqueante por vez por meio de você.
Quando o PRD estiver escrito, despache `document-validator`, modelo `opus`, exigindo a skill do kit `document-review`, só quando `reviewers.document` estiver ausente de `.harness/project.yaml` ou contiver `claude`; quando `reviewers.document` estiver presente e não contiver `claude`, despache em vez disso a primeira entrada listada (uma entrada `cli:...`, via `external-clis`) — ainda exatamente um revisor, nunca zero e nunca a lista inteira. Despache com o caminho do PRD, as notas de descoberta ou o pedido do dono e o relatório de revisão anterior quando existir. Entradas da forma `cli:...` em `reviewers.document` seguem o procedimento de `external-clis` descrito em `.agents/coordinator.md` ("Revisores de CLI externos"); tanto o despacho do Claude quanto qualquer entrada `cli:...` contam contra o limite de dois especialistas concorrentes. Persista o corpo do relatório que ele devolve como `<document>.review.md`; o validador é somente leitura. Em `changes required`, despache `product-discovery` de novo só com os pontos bloqueantes listados e valide outra vez. Limite: duas rodadas de correção; a lista completa só roda quando o dono pede o painel neste pedido ou o PRD toca segurança.

## Pré-requisitos
O objetivo declarado nas palavras do dono. Se $ARGUMENTS não traz objetivo, pergunte ao dono o que deve passar a ser possível e hoje não é, depois despache. Restrições, prazo e stack que não podem ter fonte são registrados como desconhecidos, nunca inventados.

## Saída
O resumo de `requirements-discovery`: problema, usuários afetados, comportamento desejado, fora de escopo, exemplos, alternativas consideradas, limites, aceite como "Quando X, então Y" marcado `requirement` ou `suggestion`, hipóteses, decisões em aberto e evidência. Quando o dono quiser mantê-lo, persista-o como `.harness/prd/PRD-<n>.md` a partir do template de PRD que acompanha este kit. Registre os itens de aceite e as decisões em aberto em `.harness/MEMORY.md`. O dono recebe o PRD só depois de `approved` ou do limite de rodadas, com o caminho do relatório `<document>.review.md` e qualquer ponto deixado em aberto.

## Limites
- A descoberta escreve só artefatos de descoberta: resumo, aceite, decisões em aberto. Nunca código-fonte da aplicação, configuração nem arquivos de teste.
- Pergunte só sobre decisões que mudam o que é construído; o que pode ser encontrado é encontrado, não perguntado.
- Um resumo não autoriza implementação, instalação de dependência, commit, push, publicação nem deploy.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:plan` quando o aceite estiver acordado e as decisões em aberto que bloqueiam o escopo forem respondidas.
