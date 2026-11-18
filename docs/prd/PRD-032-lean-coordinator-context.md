# PRD-032 · Contexto enxuto do Coordenador e bootstrap do repositório

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-31 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-30 / 2026-10-31 |
| Tickets | a definir pelo planejador |
| Origem | Pedido do dono em 2026-10-30, fonte `tasks/PRD-032/SOURCE.md` no harness dir |
| Depende de | PRD-031 (fluxo enxuto do builder), mesclado em `main` na 0.36.0 (commit de merge 40538a49). Este PRD parte desse estado. |
| Entrega | Um pull request |

## 1. Problema

Em uma tarefa longa, o Coordenador (a sessão principal) enche o seu contexto. O dono perde qualidade no fim da execução. O dono quer que o Coordenador faça menos trabalho por conta própria e quer um comando de bootstrap para um repositório.

A maioria dos fatos abaixo vem de `wc` e de leituras de arquivo em 2026-10-30, na branch do PRD-031. O PRD-031 agora está mesclado (0.36.0). O `.agents/coordinator.md` mesclado tem 174 linhas e 31842 bytes. A R15 mede os tamanhos de novo.

- Fato: o `MEMORY.md` do harness deste projeto tem 176 linhas e 185182 bytes. Uma linha tem 73512 caracteres. Todo comando lê este arquivo primeiro. É o maior item que o Coordenador carrega, cerca de 5,7 vezes o tamanho de `.agents/coordinator.md`.
- Fato: o `.agents/coordinator.md` tinha 174 linhas e 32663 bytes naquela branch. Tem 31842 bytes na 0.36.0. As seções "Worktree", "Consolidação a pedido explícito", "Planejadores de CLI externos" e "Relatório de deriva" somam cerca de 6100 bytes. Cada seção serve a um ou dois comandos. O Coordenador carrega todas em toda sessão.
- Fato: o `/dh:auto` manda o Coordenador seguir o texto de cada comando encadeado "no próprio texto". Os arquivos encadeados são `discover.md`, `plan.md`, `build.md`, `secure.md`, `review.md` e `release.md`. Somam cerca de 36600 bytes e ficam no contexto principal durante toda a execução.
- Fato: o `.agents/coordinator.md`, item 6 do Procedimento, exige um "formato da resposta" e um "orçamento" em cada despacho. Nenhuma regra dá um número. Só o `document-validator` tem um limite numérico (no máximo 5 frases por paráfrase).
- Fato: o relatório de um especialista volta inteiro para o chat do Coordenador. O relatório de uma revisão de código volta para o Coordenador e o Coordenador escreve o resultado no `MEMORY.md`.
- Fato: o Elephant-Goldfish Model (EGM, `google/elephant-goldfish-model`) mantém o contexto longo em uma sessão (o elefante) e envia subagentes novos, somente leitura (peixes-dourados), para revisar. Os relatórios vão para arquivos. Um sintetizador junta os achados antes de o elefante lê-los. O EGM não declara nenhuma medida de contexto poupado (pesquisa em `tasks/UND-009/egm-research.md`).
- Fato: o PRD-029 já tirou do EGM o pacote de revisão fechado, a paráfrase de compreensão e o relatório de deriva. Deixou de fora o brainstorm de seis lentes, o bootstrap do README ("amendoins" e "feno") e os hooks.
- Fato: nada persiste um mapa do repositório. O `/dh:setup` registra a stack e os comandos de teste, lint e build em `project.yaml`. O `/dh:understand` mapeia uma fatia. A skill `repository-mapping` proíbe um passeio pelo repositório inteiro. O `/dh:doctor` só checa o ambiente.
- Hipótese: o crescimento do `MEMORY.md` é o maior custo em uma execução longa, maior que qualquer texto de regra. A hipótese está errada se uma sessão medida mostrar que o texto de relatório no chat é maior que a leitura do `MEMORY.md`. As partes A e B tratam cada uma um dos dois.

## 2. Solução

Depois da entrega:

- Um relatório de especialista acima de um limite de tamanho vai para um arquivo na pasta da tarefa. O Coordenador recebe o veredito, o caminho e um resumo curto dentro de um limite numérico.
- Quando uma etapa tem mais de um revisor, um especialista junta os vereditos em um arquivo. O Coordenador lê o arquivo mesclado.
- O `/dh:auto` lê só as seções nomeadas de cada comando encadeado, na etapa que precisa delas.
- As seções que servem a um ou dois comandos saem do `coordinator.md` e vão para esses comandos.
- O `MEMORY.md` tem um orçamento de tamanho para cada entrada e para o arquivo. O Coordenador oferece a consolidação quando o arquivo passa do orçamento. A consolidação continua disparada pelo dono.
- O `/dh:setup` pode escrever um mapa do repositório. Sessões posteriores e o `repo-scout` o reutilizam.

A parte A é R1 a R4 (respostas dos especialistas). A parte B é R5 a R8 (o que o Coordenador carrega). A parte C é R9 a R12 (bootstrap do repositório). A parte D é R13 a R17 (limites, prova e correções herdadas do PRD-031).

## 3. Regras

### Parte A · Respostas com artefato primeiro

- R1 Limite de resposta e arquivo de relatório. O contrato de despacho (`.agents/coordinator.md`, Procedimento, o item que lista os campos do despacho) dá um limite de resposta numérico. O limite padrão é 200 palavras. Ele limita toda resposta. Mais uma regra governa onde o relatório fica: um relatório de mais de 40 linhas vai para `<task folder>/reports/<role>-<n>.md`, e um relatório de 40 linhas ou menos fica na resposta. A resposta sempre fica dentro do limite, então um relatório curto também pode ter menos de 40 linhas. O Coordenador pode fixar outro limite no despacho e declara o motivo. Cada papel tem o seu formato de resposta:
  - `reviewer`: veredito, contagens por severidade, no máximo 5 achados bloqueantes como `location · title` (o resto só contado). Não tem a ferramenta Write, então sempre devolve no próprio texto dentro do limite, como o `repo-scout`.
  - `qa-verifier`: veredito da matriz (pass, fail ou not-run por checagem), contagem de cada um, checagens que falharam como `check · reason`, caminho do relatório.
  - `builder`: status da fatia, arquivos alterados, `Red:` e resultados das checagens, caminho do relatório. Não tem contagem por severidade.
  - `docs-guide`: arquivos escritos com os caminhos, e o que cada um contém. Não tem veredito.
  - `debugger`: status do diagnóstico (confirmed, hypothesis ou inconclusive), a causa em uma frase, caminho do relatório.
  - `repo-scout`: não tem a ferramenta Write, então sempre devolve no próprio texto dentro do limite, como uma lista de caminhos com uma linha para cada caminho (veja a R2).
  Verificável: o `coordinator.md` contém a expressão `limite de resposta` e os números `200` e `40`. `builder.md`, `reviewer.md`, `qa-verifier.md`, `docs-guide.md`, `debugger.md` e `repo-scout.md` declaram cada um o seu formato de resposta no Formato de saída. O `document-validator` mantém o seu próprio limite.
- R2 Quem escreve o arquivo de relatório. Estes escritores mantêm os seus conjuntos de escrita atuais e acrescentam só o caminho de relatório nomeado no despacho:
  - `builder`, `qa-verifier`, `docs-guide` e `debugger` escrevem o próprio arquivo de relatório.
  - O `reviewer` continua totalmente somente leitura. Não tem a ferramenta Write e não escreve arquivo de relatório (D1). Devolve a resposta no próprio texto (R1).
  - O `repo-scout` não tem a ferramenta Write. Devolve no próprio texto dentro do limite e lista caminhos, não corpos de arquivo (R1). Para um mapa maior que o limite, o Coordenador despacha o `docs-guide` (R10).
  - Um revisor de CLI nunca escreve. O `qa-verifier` ou o `builder` (`lane` `infra`) que roda a chamada escreve o arquivo de relatório, como o PRD-025 define para a própria chamada.
  Verificável: as `tools` de `reviewer.md` não têm Write. A afirmação "a revisão é somente leitura" continua verdadeira sem exceção.
- R3 Disciplina de leitura do Coordenador. O Coordenador lê um arquivo de relatório só para julgar um achado bloqueante, para resolver uma disputa ou para tomar uma decisão que a resposta não sustenta. Não lê um relatório inteiro quando o veredito é `approved` e não existe achado bloqueante. A seção "Dúvida no meio do build" e os itens do Procedimento de `plan-loop.md` que julgam relatórios de onda seguem esta regra. O `MEMORY.md` registra o caminho do relatório e o veredito, não o texto do relatório. Verificável: o `coordinator.md` contém esta regra uma vez e nenhuma seção a contradiz. Uma busca por "leia o relatório" em `.agents/` e `.commands/` mostra só lugares que nomeiam uma condição.
- R4 Sintetizador. Uma etapa de revisão pode ter mais de uma entrada de revisor no despacho. Então uma execução do `docs-guide` (modelo `sonnet`) junta as saídas dos revisores (os arquivos de relatório da CLI escritos pelo delegado e a resposta no próprio texto de um revisor `claude`, passada no despacho como dado citado) em um arquivo `<task folder>/reports/<stage>-merged-<round>.md`. O Coordenador lê este arquivo, não os relatórios brutos. As regras da junção:
  - A entrada são as saídas dos revisores nomeadas no despacho (arquivos, ou uma resposta no próprio texto citada como dado), nunca o chat.
  - Achados com o mesmo local e o mesmo título viram um achado. A junção lista cada revisor que o levantou.
  - A junção nunca remove um achado bloqueante.
  - A junção nunca reduz uma severidade.
  - O veredito de cada revisor continua literal no arquivo.
  - A junção não acrescenta achado nem julgamento. O Coordenador julga.
  Com uma entrada de revisor, nenhum sintetizador roda. A etapa de documento mantém `<document>.review.md` como o arquivo mesclado, então não precisa de uma segunda junção. Verificável: o `coordinator.md` ("Revisores de CLI externos", passo de junção) nomeia o `docs-guide` como o responsável pela junção nas etapas `code` e `security`. O `docs-guide.md` lista a tarefa de junção e as suas regras.

### Parte B · O que o Coordenador carrega

- R5 O `/dh:auto` deixa de embutir os comandos encadeados. No início de cada etapa encadeada, o Coordenador lê só três seções do arquivo do comando encadeado: `Roteamento`, `Pré-requisitos` e `Limites`. Pula `Saída` e toda seção que a etapa não usa. O `auto.md` traz uma tabela com uma linha para cada etapa encadeada (`discover`, `plan`, `build`, `secure`, `review`, `release`). Cada linha dá o arquivo, as seções a ler e um orçamento em bytes. O Coordenador lê as seções da etapa uma vez por execução. Uma sessão nova que retoma a execução (`/dh:resume`) lê de novo. A regra da tabela substitui a frase "Siga, no próprio texto, o Roteamento e os Pré-requisitos de cada comando encadeado". Os gates, os limites e as checagens de autorização encadeados não mudam. Verificável: o `auto.md` contém a tabela de etapas com seis linhas. O `auto.md` não diz mais que o Coordenador lê o `release.md` inteiro "no próprio texto, nesta mesma sessão".
- R6 Mover seções de um só comando para fora do `coordinator.md`. Cada uma destas seções vai para o comando que a usa, com o texto inalterado, exceto pelas referências:
  - "Worktree" vai para `.commands/build.md`. O `auto.md` aponta para essa seção.
  - "Consolidação a pedido explícito" vai para `.commands/consolidate-memory.md`.
  - "Planejadores de CLI externos" vai para `.commands/plan-loop.md`.
  - "Relatório de deriva" vai para `.commands/review.md`.
  O `coordinator.md` mantém uma linha para cada seção movida que nomeia o comando. As seções "Idioma", "Seleção de modelo", "Revisores de CLI externos" (ela guarda a cláusula anti-delegação) e "Revisão proporcional e avaliação paga" ficam. Meta: o `coordinator.md` fica menor em pelo menos 5500 bytes que o seu tamanho na 0.36.0 (31842 bytes). Meça com `wc -c`. Verificável: `wc -c .agents/coordinator.md` dá no máximo 26342. `go run ./cmd/dh validate` passa, inclusive `checkAntiDelegationClause`. Uma busca não acha referência a uma seção movida pelo lugar antigo.
- R7 Orçamento do `MEMORY.md`. O Coordenador mantém cada entrada do `MEMORY.md` (um marcador, ou um parágrafo sob um título) em no máximo 1200 caracteres. Mantém o arquivo inteiro em no máximo 40000 bytes. O detalhe além de uma entrada vai para um arquivo na pasta da tarefa. A entrada guarda o caminho e uma frase. As entradas existentes ficam como estão: o disco é a verdade e o Coordenador não reescreve registros antigos para caber no orçamento. A regra vale para toda entrada nova ou alterada. Verificável: `.agents/coordinator.md` ("Posse da memória") e `templates/en/MEMORY.md` e `templates/pt-br/MEMORY.md` declaram os dois números.
- R8 Oferta de consolidação. Quando uma leitura ou o `/dh:doctor` mostra o `MEMORY.md` acima do orçamento do arquivo (R7), o Coordenador avisa o dono uma vez naquela sessão e oferece o `/dh:consolidate-memory`. Nunca roda a consolidação sem o pedido do dono. O `/dh:doctor` imprime o tamanho em bytes do `MEMORY.md` no harness dir e um aviso acima de 40000 bytes. Verificável: um teste da checagem do `doctor`, em `internal/`, para um arquivo abaixo, no limite e acima do orçamento. O `consolidate-memory.md` declara que a oferta é o único passo automático.

### Parte C · Bootstrap do repositório

- R9 Mecanismo. O bootstrap estende o `/dh:setup` como uma segunda fase opcional. Não é um comando novo. Depois do perfil, o `/dh:setup` oferece escrever o mapa do repositório. O dono aceita ou recusa. Uma nova execução com a nota `bootstrap` nos argumentos atualiza o mapa sem refazer o perfil. As contagens de comandos continuam em 19 + 1. Verificável: o `setup.md` contém a oferta, o argumento `bootstrap` e o caminho do mapa. A contagem de comandos em `go run ./cmd/dh validate` não muda.
- R10 Arquivo do mapa e escritor. O mapa é `<harness dir>/REPO-MAP.md`, um arquivo markdown no `language` do `project.yaml`. Não tem template nem chave em `project.yaml`. A skill `project-onboarding` ganha uma seção "Bootstrap do repositório" que define o formato e o procedimento. Uma execução do `docs-guide` (modelo `sonnet`) com as skills `repository-mapping` e `project-onboarding` lê o repositório e escreve o arquivo. A resposta dele segue a R1. O mapa tem no máximo 200 linhas e estas partes:
  1. Cabeçalho: data e o sha do `HEAD`. O Coordenador obtém o sha de uma execução do `qa-verifier` (`git rev-parse HEAD`) e o passa no despacho. O escritor o copia e não roda o git.
  2. Estrutura de nível superior: uma linha por diretório.
  3. Pontos de entrada e como rodar, testar, fazer lint e fazer build, com a fonte de cada comando.
  4. Arquitetura em no máximo 15 linhas: os componentes principais e como se ligam.
  5. Áreas de risco: caminhos em que uma mudança exige cuidado extra, com o motivo.
  6. Convenções que o código mostra, cada uma com um caminho de exemplo.
  Toda afirmação cita um caminho relativo. Uma afirmação que o escritor não verificou é marcada como `hypothesis`. Quando o repositório tem mais de 15 diretórios de nível superior, o escritor para e o Coordenador pede ao dono que nomeie o escopo (o gate de segurança do EGM, adaptado). O escritor lê o repositório e não roda código nem scripts do projeto. Verificável: uma execução em uma fixture produz um mapa de no máximo 200 linhas com as seis partes e um sha no cabeçalho.
- R11 Reutilização. O `repo-scout` lê o `REPO-MAP.md` primeiro quando ele existe e o trata como uma dica. O Coordenador acrescenta o caminho do mapa em `sources` nos despachos para o `repo-scout`, o `implementation-planner` e o `builder`. O Coordenador não lê o mapa por conta própria. O `builder` e o `implementation-planner` tratam o mapa como uma dica e o verificam contra a árvore, como o `repo-scout` faz. A `repository-mapping` ganha uma linha: use o mapa como ponteiro e verifique qualquer afirmação antes de ela sustentar um plano. Verificável: `repo-scout.md`, `builder.md`, `implementation-planner.md` e `repository-mapping` contêm cada um `REPO-MAP.md`. O `coordinator.md` não manda o Coordenador lê-lo.
- R12 Defasagem e atualização. O Coordenador passa o sha atual do `HEAD` (do `qa-verifier`) em cada despacho do `repo-scout`. Quando o sha do cabeçalho difere dele, o `repo-scout` reporta o mapa como desatualizado e ainda o usa como ponteiro. A atualização é disparada pelo dono (R9). O `/dh:setup` pode sobrescrever o `REPO-MAP.md`, a única exceção à sua regra "nunca sobrescreva um registro existente", porque o arquivo é dado derivado. No modo `repo`, o `.harness/REPO-MAP.md` é versionado como os outros registros internos. Verificável: os Limites do `setup.md` nomeiam a exceção para este único arquivo. O relatório do `repo-scout` tem uma linha `Map: fresh | stale | none`.

### Parte D · Limites e prova

- R13 Nenhum agente, skill, comando, template nem chave de `project.yaml` novos. As contagens continuam: 13 agentes, 22 skills, 19 + 1 comandos. A contagem de templates continua. `go run ./cmd/dh validate` passa na checagem das contagens documentadas.
- R14 Estas salvaguardas continuam em vigor:
  - Honestidade do not-run: uma checagem que não rodou é `not-run`, nunca `passed`.
  - Posse da memória: só o Coordenador escreve `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`. Um arquivo de relatório na pasta da tarefa não é um deles.
  - Revisão independente: uma afirmação do builder não é revisão, e um revisor não altera nenhum arquivo do repositório.
  - O pacote de revisão exclui a conversa do Coordenador (a regra do pacote de contexto no `AGENTS.md`).
  - A cláusula anti-delegação, literal no `coordinator.md` e na `external-clis`.
  - Autorização: a aprovação do plano não autoriza commit, push, publicação nem deploy.
  - Uma execução de revisor de CLI não altera nada na árvore. Uma mudança na árvore descarta o veredito.
  Verificável: `go run ./cmd/dh validate` passa. O builder lista cada salvaguarda com o arquivo que a guarda na descrição do pull request.
- R15 Medição. O primeiro ticket registra uma linha de base, e o último ticket registra o resultado, no Resumo executado deste PRD. Os registros dão: `wc -c` de `.agents/coordinator.md`; os bytes que o Coordenador lê dos seis comandos encadeados no `/dh:auto` (antes: os seis arquivos inteiros; depois: só as seções da tabela); e o tamanho do template de `MEMORY.md`. Nenhuma avaliação paga e nenhuma ferramenta de contagem de tokens faz parte deste PRD, e o Resumo executado diz que os números são bytes, não tokens. Verificável: o Resumo executado tem as linhas da linha de base e do resultado.
- R16 Premissa de limpeza. O mesmo pull request atualiza todo documento que declara um comportamento que este PRD muda (seção 4) e regenera o `dist/` com `go run ./cmd/dh build`. Uma busca por cada nome de seção movida não acha referência velha fora de `docs/prd/`, `docs/adr/`, do histórico do `CHANGELOG.md` e de `evals/baselines/`. Como a R8 muda o `internal/`, a prova de clone limpo do `AGENTS.md` vale antes de qualquer push. Verificável: `go run ./cmd/dh validate`, `go test ./...` e um `git diff` limpo de `dist/` depois de `go run ./cmd/dh build`.
- R17 Correções herdadas do PRD-031. Quatro defeitos da entrega do PRD-031 são corrigidos neste pull request:
  - `.skills/code-review/SKILL.md` volta a perguntar, no procedimento de revisão, se uma regressão plausível faria um teste falhar. Verificável: o arquivo contém essa pergunta uma vez.
  - `.skills/incremental-implementation/SKILL.md` restaura as três salvaguardas da regra da página oficial: pular a busca quando só uma opção é viável; enviar um conflito com a convenção local ao gate de decisão em aberto; pôr a URL no relatório, não em um comentário de código. Verificável: o arquivo contém cada salvaguarda uma vez.
  - `docs/prd/PRD-031-rule-inventory.md` declara o total como 55830 bytes. Verificável: `grep -n "55830\|55,830" docs/prd/PRD-031-rule-inventory.md` mostra o número e nenhum outro total.
  - `.skills/regression-testing/SKILL.md` passo 13 dá uma frase para cada heurística. Verificável: o passo 13 tem três frases, uma para cada heurística da R3 do PRD-031.

## 4. Docs

O builder confirma cada documento com uma busca e o edita no mesmo pull request. Os pares atualizam juntos (pt-br e `docs/en/`).

- `AGENTS.md` (marcadores "Escopo do Coordenador", "Agentes", "Verificação padrão" e "Pasta global do harness")
- `.agents/coordinator.md`
- `.agents/builder.md`
- `.agents/reviewer.md`
- `.agents/qa-verifier.md`
- `.agents/docs-guide.md`
- `.agents/debugger.md`
- `.agents/repo-scout.md`
- `.agents/implementation-planner.md`
- `.commands/auto.md`
- `.commands/build.md`
- `.commands/consolidate-memory.md`
- `.commands/plan-loop.md`
- `.commands/review.md`
- `.commands/setup.md`
- `.commands/doctor.md`
- `.skills/project-onboarding/SKILL.md`
- `.skills/repository-mapping/SKILL.md`
- `.skills/code-review/SKILL.md`
- `.skills/security-review/SKILL.md`
- `.skills/incremental-implementation/SKILL.md` (R17)
- `.skills/regression-testing/SKILL.md` (R17)
- `docs/prd/PRD-031-rule-inventory.md` (R17)
- `templates/en/MEMORY.md`
- `templates/pt-br/MEMORY.md`
- `internal/` (a checagem de tamanho do `doctor` e o seu teste)
- `docs/tutorial.html`
- `docs/en/tutorial.html`
- `docs/novo-fluxo-prd.md`
- `docs/en/new-prd-flow.md`
- `docs/roadmap.html`
- `docs/en/roadmap.html`
- `CHANGELOG.md` (entrada nova, aumento de versão, linha das contagens)
- `docs/adr/ADR-010-lean-coordinator-context.md` (novo: registra a regra do artefato primeiro, as mudanças de seção e a escolha do bootstrap)

O builder também busca em `adapters/`, `hooks/` e `internal/` código que leia as seções movidas do `coordinator.md` ou o layout do `MEMORY.md`. Uma ocorrência está no escopo.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Novo comando `/dh:bootstrap` | A contagem fixa 19 + 1 muda no `AGENTS.md`, no ADR-003, em `internal/kit/doc_counts.go`, no roadmap e no tutorial (pt e en), no `marketplace.json` e no `REVIEW.md`. O `/dh:setup` já é dono do conhecimento do projeto na primeira vez. |
| Bootstrap pelo `/dh:understand` | O comando mapeia uma fatia. Um modo de repositório inteiro o confunde e contradiz a `repository-mapping`. |
| Árvore no README ("amendoins" e "feno") | Acrescenta arquivos ao repositório do consumidor. Um único mapa no harness dir tem o mesmo uso sem mudança no repositório. |
| Subagente executor de etapa para cada etapa do `/dh:auto` | Um subagente não pode iniciar subagentes, então um executor não poderia despachar builders e revisores. |
| Editar o `MEMORY.md` automaticamente no orçamento | A consolidação continua disparada pelo dono (`AGENTS.md`). O disco é a verdade para os registros. |
| Brainstorm de seis lentes do EGM, hooks do EGM | Excluídos no PRD-029. Acrescentam contexto, não tiram. |
| Um modelo mais barato para o Coordenador | Fora do pedido do dono. A escolha do modelo não é a causa do crescimento do contexto. |
| Avaliações pagas do contexto poupado | Só a pedido do dono (`AGENTS.md`). A medida é em bytes (R15). |

### Riscos e decisões pendentes

- Risco: a medida em bytes é um substituto dos tokens. A economia real só aparece em execuções longas. Mitigação: a R15 declara os números como bytes. O dono julga pelo uso.
- Risco: a R6 poupa cerca de 5500 bytes, pouco ao lado da leitura do `MEMORY.md` (185182 bytes neste projeto). Mitigação: a R7 e a R8 miram o item maior. O `MEMORY.md` existente não é reescrito por este PRD. O dono pode rodar `/dh:consolidate-memory`.
- Risco: o Coordenador confia na resposta e pula o arquivo (R3). Um resumo pode esconder um achado bloqueante. Mitigação: a resposta deve listar até 5 achados bloqueantes, e a R4 proíbe descartar ou reduzir um.
- Risco: uma resposta de `reviewer` sem arquivo de relatório pode perder detalhe. Mitigação: a resposta guarda o veredito, as contagens e até 5 achados bloqueantes (R1). O Coordenador pode pedir uma nova revisão mais estreita.
- Risco: seções movidas para fora do `coordinator.md` (R6) podem ser puladas quando um comando é rodado à mão sem o Coordenador lê-lo. Mitigação: cada comando já lê o próprio texto. O ponteiro de uma linha continua no `coordinator.md`.
- Risco: um mapa desatualizado engana o planejamento (R12). Mitigação: é uma dica, o cabeçalho com o sha marca a defasagem e as afirmações são verificadas.
- Risco: as seções podem mudar antes de o build começar. Mitigação: este PRD cita nomes de seção, não linhas, e a R6 mede a partir do tamanho da 0.36.0.
- Sugestão (não é regra): depois de cada etapa do `/dh:auto`, a linha de status poderia dizer que a execução pode continuar em uma sessão nova com `/dh:resume`. A linha `Stage:` já permite isso. O Coordenador não consegue medir o próprio contexto, então o PRD não faz disso uma regra.

### Decisões (fechadas em 2026-10-31, dono: seguir as recomendações)

Cada item foi decidido em 2026-10-31 (dono: seguir as recomendações).

- D1 O `reviewer` não ganha a ferramenta Write e continua somente leitura. Devolve o relatório no próprio texto (R1, R2). Decidido em 2026-10-31. Motivo: um limite escrito só em prosa era uma exposição confirmada. Uma regra de trava de runtime que permita uma escrita estreita do revisor é um candidato futuro, fora de escopo aqui.
- D2 O `docs-guide`, modelo `sonnet`, roda o sintetizador (R4). Decidido.
- D3 O `/dh:auto` lê as seções nomeadas por uma tabela de etapas (R5). Retomar em uma sessão nova continua sendo escolha do dono e uma sugestão. Decidido.
- D4 Limiares: limite de resposta de 200 palavras, relatório em arquivo acima de 40 linhas, entrada com no máximo 1200 caracteres, arquivo com no máximo 40000 bytes, mover pelo menos 5500 bytes (R1, R6, R7). Decidido.
- D5 O `dh doctor` imprime o tamanho do `MEMORY.md` e avisa acima do orçamento (R8). Decidido.
- D6 O bootstrap estende o `/dh:setup`. As contagens continuam 19 + 1 (R9). Nota do dono em 2026-10-31: um comando de bootstrap separado pode ser preciso no futuro; fora de escopo aqui. Decidido.
- D7 O mapa é `<harness dir>/REPO-MAP.md`, escrito pelo `docs-guide`, sem template e sem chave (R10). Decidido.
- D8 Aumento de versão menor, de 0.36.0 para 0.37.0. Decidido.

Emenda de 2026-10-31: a D1 foi revertida depois de uma revisão de segurança (achado médio: texto não confiável de um PR poderia chegar a `project.yaml` ou `REVIEW.md` por uma escrita do `reviewer`). A R1, a R2 e o Risco relacionado seguem a D1. A R11 agora diz que o `builder` e o `implementation-planner` verificam o mapa.

### Fora de escopo

- Hooks do EGM, o brainstorm de seis lentes e a árvore no README.
- Uma mudança de modelo para o Coordenador.
- Avaliações pagas e qualquer ferramenta de contagem de tokens.
- Consolidação automática do `MEMORY.md` ou reescrita dos registros existentes.
- Mudanças nas regras de revisão, nos estados de achado ou nos limites de rodada (o PRD-031 é o dono delas).

## Resumo executado

Os números são bytes, não tokens. Medidos na 0.36.0 antes de qualquer edição do PRD-032 (`wc -c`; uma seção vai do seu cabeçalho `## <name>` até o próximo `## `).

- Linha de base do coordinator.md: 31842 bytes (`.agents/coordinator.md`).
- Linha de base do auto antes: 36583 bytes, os seis comandos encadeados inteiros (discover 3339, plan 7427, build 4734, secure 3205, review 16061, release 1817). Seções `## Routing` / `## Prerequisites` / `## Limits` em bytes: discover 1463 / 266 / 415; plan 4546 / 530 / 416; build 2266 / 393 / 914; secure 757 / 469 / 884; review 8442 / 1037 / 1784; release 494 / 275 / 227.
- Linha de base do template de MEMORY: `templates/en/MEMORY.md` 2579 bytes; `templates/pt-br/MEMORY.md` 2762 bytes.

Resultado (T-3217, commit a159079, mesmo método e mesmos arquivos da linha de base; `wc -c`):

- Resultado do coordinator.md: 26338 bytes, 5504 abaixo da linha de base, dentro do limite da R6 de 26342.
- Resultado do auto depois: 25524 bytes, as seções `Routing` / `Prerequisites` / `Limits` dos seis comandos encadeados (discover 2142, plan 5482, build 3553, secure 2110, review 11241, release 996), contra 36583 antes. Os orçamentos da tabela do `auto.md` (2150, 5500, 3600, 2150, 11300, 1000) cobrem esses tamanhos.
- Resultado do template de MEMORY: `templates/en/MEMORY.md` 2849 bytes; `templates/pt-br/MEMORY.md` 3050 bytes (cresceram porque agora declaram o orçamento da R7).

**Entregue:** O Coordenador agora carrega menos (seções movidas, uma tabela de etapas no `auto.md` e um orçamento do `MEMORY.md` com aviso do doctor), os papéis respondem dentro de um limite de 200 palavras com relatórios longos em arquivos, e o `/dh:setup` oferece um bootstrap opcional do `REPO-MAP.md` que os papéis de planejamento e build usam como dica, em um pull request na 0.37.0, sem agente, skill, comando, template nem chave novos.

**Regras:**

- R1 Honrada. Matriz de QA `tasks/PRD-032/reports/qa-verifier-2.md` (pass); formatos de resposta nos seis arquivos de papel.
- R2 Honrada. O `reviewer.md` não tem a ferramenta Write (D1 revertida pela emenda de 2026-10-31); matriz de QA `qa-verifier-2.md`; revisão de segurança `evidence/review-r2/`.
- R3 Honrada. `qa-verifier-2.md` (regra uma vez no `coordinator.md`, as outras menções nomeiam uma condição).
- R4 Honrada. `qa-verifier-2.md` (o `coordinator.md` nomeia o `docs-guide` como responsável pela junção; o `docs-guide.md` lista as regras).
- R5 Honrada. Tabela de etapas de `.commands/auto.md` com seis linhas; tamanhos acima; `qa-verifier-2.md`. Os orçamentos em bytes da R5 não foram medidos de novo nesse relatório; foram medidos de novo aqui.
- R6 Honrada. 26338 bytes (limite 26342); `go run ./cmd/dh validate` passa com o hash da cláusula anti-delegação inalterado (`qa-verifier-2.md`, `qa-verifier-3.md`).
- R7 Honrada. Os dois números no `coordinator.md` e nos dois templates de MEMORY (`qa-verifier-2.md`).
- R8 Honrada. `printHarnessFor` em `internal/doctor/harness.go` com um teste abaixo, no limite e acima do orçamento; `go test ./...` passa (`qa-verifier-3.md`). O texto do comando `/dh:doctor` não foi rodado (limitação em `qa-verifier-2.md`).
- R9 Honrada. `qa-verifier-2.md`; a contagem de comandos continua 19 + 1.
- R10 Honrada. O mapa da fixture (30 linhas, seis partes, sha do cabeçalho igual a `git rev-parse HEAD`) foi conferido em `reports/qa-verifier-1.md` em um estado anterior. Não foi repetido no estado final. Observação nesse relatório: o escritor rodou `git ls-files` além do sha.
- R11 Honrada. `qa-verifier-2.md`.
- R12 Honrada. `qa-verifier-2.md`.
- R13 Honrada. Contagens 13 agentes, 22 skills, 19 + 1 comandos, 14 templates; validate passa (`qa-verifier-2.md`, `qa-verifier-3.md`).
- R14 Honrada na árvore: as salvaguardas estão intactas (validate passa, sha256 da cláusula anti-delegação igual no `coordinator.md` e na `external-clis`, caminhos críticos intocados: `qa-verifier-3.md`). Pendente: a lista de salvaguardas com os seus arquivos no corpo do pull request ainda não foi escrita; isto não é dado como feito.
- R15 Honrada. Linhas da linha de base e do resultado acima; os números são bytes, não tokens.
- R16 Honrada. `dist/` reconstruído; prova de clone limpo PASS em `qa-verifier-3.md` (7 binários com o mesmo sha256, validate 0 erros, `go test ./...` ok, só `generated_at` difere).
- R17 Honrada. `qa-verifier-2.md`.

**Tickets:**

| Id | Lane | Veredito | Evidência |
|---|---|---|---|
| T-3201 | backend | done | linha de base acima |
| T-3202 | backend | done | `qa-verifier-2.md` (R17); `evidence/review-r3/` |
| T-3203 | backend | done | `go test ./...` em `qa-verifier-3.md`; `evidence/review-r3/` |
| T-3204 | backend | done | coordinator.md 26338 bytes; `evidence/review-r3/` |
| T-3205 | backend | done | `qa-verifier-2.md` (R1, R2, R4); `evidence/review-r3/` |
| T-3206 | backend | done | segurança `evidence/review-r2/`: `cli:codex/gpt-6.1-sol` não achou vulnerabilidade; o revisor `claude` deu `approved` depois que F1-1 e F1-2 foram corrigidos; 2 achados Low, não bloqueantes |
| T-3207 | backend | done | `qa-verifier-2.md` (R1 a R4); `evidence/review-r4/` (codex approve) |
| T-3208 | backend | done | `qa-verifier-2.md` (R9 a R12); `evidence/review-r4/` |
| T-3209 | backend | done | `qa-verifier-2.md` (R7, R8); `evidence/review-r4/` |
| T-3210 | backend | done | coordinator.md fecha em 26338 bytes; `evidence/review-r4/` |
| T-3211 | backend | done | tabela do `auto.md`; `qa-verifier-2.md` |
| T-3212 | test | pass (13 pass; R10, R14, R15, R16 pendentes na época, fechadas acima) | `reports/qa-verifier-2.md` |
| T-3213 | backend | done | `AGENTS.md`, `docs/adr/ADR-010-lean-coordinator-context.md`; `evidence/review-r5/` (opencode approve) |
| T-3214 | backend | done | docs e `CHANGELOG.md` 0.37.0; `evidence/review-r5/` |
| T-3215 | infra | done | versão e reconstrução do `dist/`; `evidence/review-r5/` |
| T-3216 | test | PASS | `reports/qa-verifier-3.md` |
| T-3217 | backend | done | esta seção |

Revisões do lote: `evidence/review-r3/` (opencode approve), `evidence/review-r4/` (codex approve), `evidence/review-r5/` (opencode approve, docs e versão). Revisão final da funcionalidade: `evidence/review-r6/` (codex; achado F6-1 resolvido como intencional pela linha 72 do PLAN) e `evidence/review-r7/` (opencode approve).

Revalidação do documento depois da emenda de 2026-10-31 (D1 revertida): `evidence/doc-r4/` (agy, changes required). O Coordenador corrigiu os pontos P8 a P10 sem outra rodada. Registrado para o dono.

**Docs:** `AGENTS.md` (regras atuais), `docs/adr/ADR-010-lean-coordinator-context.md`, os tutoriais e o fluxo de PRD, os dois roadmaps, `CHANGELOG.md` 0.37.0, a reconstrução do `dist/` e este PRD (linha de base, resultado, Resumo executado, Status). Fonte: tabela de posse do PLAN.md, T-3213 a T-3215.

**Fora de escopo:** inalterado em relação à seção acima: hooks do EGM e o brainstorm de seis lentes, uma mudança de modelo do Coordenador, avaliações pagas e contagem de tokens, consolidação automática do `MEMORY.md` e mudanças nas regras de revisão. Um comando de bootstrap separado e uma regra de trava de runtime para uma escrita estreita do revisor continuam candidatos futuros.
