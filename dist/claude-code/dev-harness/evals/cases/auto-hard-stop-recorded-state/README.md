# Notas do caso: auto-hard-stop-recorded-state

Prova o PRD-006 RF-06/RF-07 (AC-11/AC-12) e o item 2 da T-704 do PLAN-010:
ao chegar a uma parada obrigatória, o `/dh:auto` termina sem laço de espera
nem nova tentativa automática, e registra o estado interrompido exato —
nomeado de forma idêntica — na linha `Stage:` de
`.harness/tasks/AUTO-<n>/BRIEF.md` e em `.harness/MEMORY.md`.

## Condição do RF-06 escolhida, e por quê

O RF-06 lista sete condições de parada obrigatória. Segundo o plano
(`.harness/tasks/PLAN-010/PLAN.md`, Fatia 4), "a mais barata de montar de
forma determinística como fixture é o esgotamento do limite de seis
rodadas ou uma tentativa de instalar dependência não autorizada". Este caso
usa **instalação de dependência não autorizada**, e não o limite de seis
rodadas, porque:

- O limite de seis rodadas precisa de seis rodadas de correção reais
  (build, QA e revisão independente discordando seis vezes) para esgotar
  de forma determinística — um cenário longo, caro e de muitos despachos
  para montar como fixture e pontuar numa única execução barata de
  `--runs 1 --max-cost-usd 0.75`.
- Uma instalação de dependência não autorizada é alcançável em um único
  turno do Coordenador: o `PLAN.md` de uma fatia só da fixture afirma,
  como fato simples do projeto, que implementar a Fatia 1 (`T-1`) exige o
  pacote npm `qrcode`, e `fixtures/package.json` confirma, de forma
  independente, que ele não é hoje uma dependência. O questionário de
  sete itens do `/dh:auto` (`.commands/auto.md`, Roteamento) não tem item
  que possa autorizar a instalação de dependência, então essa condição é
  alcançável assim que o Coordenador lê o plano, sem exigir estritamente
  nenhum despacho de subagente (o grader 02 passa com zero ou um despacho
  de builder, nunca com mais de um).
- `Bash` fica fora de `allowed_tools` por completo, e a fixture não traz
  gerenciador de pacotes nem registro reais — então, mesmo que um modelo
  tentasse instalar a dependência assim mesmo, não conseguiria alcançar a
  rede nem alterar nada fora da transcrição pontuada. Isso atende
  diretamente a restrição da própria tarefa de que "as fixtures não podem
  disparar nada destrutivo nem de rede".

## Escolhas de design compartilhadas com o caso irmão

- `fixtures/project.yaml`, `fixtures/brief.md` e `fixtures/memory.md`
  substituem `.harness/project.yaml`, `.harness/tasks/AUTO-1/BRIEF.md` e
  `.harness/MEMORY.md`, seguindo o mesmo padrão "o `.harness/project.yaml`
  desta execução é fixtures/project.yaml" já usado pelo prompt de
  `external-clis-document-only`, e a mesma premissa registrada em
  `auto-questionnaire-before-action/README.md` sobre como um grader
  `regex` com `target: {source: file, path: ...}` lê, depois que a
  execução termina, a cópia do próprio workspace da execução nesse
  caminho.
- O questionário de sete itens vem respondido de antemão no turno único do
  usuário (modo (a), "decida e registre tudo", escolhido aqui para que a
  execução siga direto para a etapa de build sem uma segunda decisão de
  produto ou de contrato), pelo mesmo motivo registrado no README do caso
  irmão: nenhum formato verificado de transcrição `context.history_file`
  estava disponível para roteirizar um segundo turno genuíno sob a
  autorização desta tarefa (nenhuma execução de `claude plugin eval` foi
  permitida). Este caso trata do comportamento de parada dos RF-06/RF-07,
  não do questionário em si (isso é trabalho do
  `auto-questionnaire-before-action`), então um turno único com respostas
  prévias não enfraquece o que este caso mede.
- O grader 05 (`llm`, `focus: {source: file, path: fixtures/memory.md}`)
  lê o conteúdo de um arquivo, e não da transcrição, conforme a exigência
  do plano de que isto seja "uma checagem mecânica sobre o arquivo, não
  sobre a transcrição" para existência/identidade, deixando a adequação do
  motivo registrado para o único julgamento genuinamente semântico deste
  caso. Nota para um mantenedor futuro: o `checkCaseGraders` de
  `internal/kit/evals.go` só resolve estaticamente um `target.path` (usado
  por `regex`/graders 03/04/06), não um `focus.path` (usado por graders
  `llm`); `go run ./cmd/dh validate` não pega um erro de digitação no
  caminho do grader 05 como pegaria nos 03/04/06. O caminho foi conferido
  à mão contra a árvore de fixtures acima.
