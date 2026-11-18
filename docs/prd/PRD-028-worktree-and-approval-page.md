# PRD-028 · Configuração de worktree e página HTML de aprovação

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-28 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-27 / 2026-10-27 |
| Tickets | a definir pelo planejador |
| Origem | Pedido do dono em 2026-10-27 (descoberta encerrada pelo dono) |
| Depende de | PRD-014 (harness global, entregue), PRD-008 e ADR-002 (isolamento por clone dos revisores de CLI, entregues) |
| Entrega | Um pull request, versão 0.33.0 |

Nota de 2026-11-16: R63 (a pasta irmã `<checkout>.worktrees`) é substituída pelo PRD-033: a pasta do worktree é `<checkout>/.worktrees/<task-id>`. O resto deste PRD continua valendo.

## 1. Problema

Duas necessidades viajam em um só lote. As duas pertencem ao fluxo de desenvolvimento.

**W · O desenvolvimento sempre roda no checkout principal.** Fato: nenhuma configuração deixa o dono dizer "desenvolva esta feature em um worktree do git". Fato: `.commands/build.md` e `.commands/auto.md` nunca criam um worktree. Fato: o `auto.md` cria um `<base-branch>` a partir de `main` por execução, no checkout principal. Impacto: o dono não consegue manter limpo o checkout principal enquanto uma execução roda.

Três partes do kit bloqueiam ou contradizem um worktree hoje.

- Fato: `dh harness-path` resolve por `<repo>/.harness` e depois pelo caminho canônico do repositório no registro (`internal/harness/resolve.go`). Um worktree é outro caminho. No modo global, a resolução dá `none`.
- Fato: `.skills/external-clis/references/clone-isolation.md` e `references/live-tree-freeze.md` recusam um projeto cujo `.git` não é um diretório de verdade. Um worktree vinculado tem um gitfile, então toda chamada de revisão por CLI termina como `isolation unavailable`.
- Fato: `docs/roadmap.html` e `docs/en/roadmap.html` ainda prometem a um lote futuro `"isolation by disposable git worktree"` para os revisores de CLI. Isso contradiz a decisão do clone em `AGENTS.md`, no ADR-002 e no PRD-008.

**A · O dono aprova um plano sem uma página estruturada.** Fato: o dono aprova um plano a partir de texto no chat. Um plano com pontos de decisão não tem uma visão que mostre todos os passos, os trechos propostos e os pontos a aprovar. O dono também não tem um lugar fixo para deixar o motivo de uma rejeição.

## 2. Solução

Depois da entrega, o dono escolhe onde o desenvolvimento roda. Uma configuração `worktree` tem dois valores: `always` e `on-request`. O padrão é `on-request`. O dono a define globalmente com `dh config set worktree <value>`. O dono a sobrescreve por projeto em `local.yaml`.

Quando o valor efetivo exige um worktree, o Coordenador cria um worktree para a execução. O Coordenador despacha o `builder` com o caminho do worktree. Um revisor de CLI chamado a partir desse worktree se comporta conforme o tipo de chamada.

- Uma chamada baseada em clone aceita um worktree vinculado válido e mantém todas as travas do RISK-002 (R26 a R34).
- Uma chamada somente leitura sem clone, hoje a do `codex`, não muda. Ela não depende do layout do worktree.
- A revisão da árvore viva do `agy` termina como `isolation unavailable` a partir de um worktree vinculado (R64).

Na fase de planejamento, quando o plano tem pontos de decisão, o `docs-guide` escreve uma página HTML de aprovação na pasta de pdocs. A página lista os passos do trabalho, cada um com uma descrição. Ela mostra os trechos propostos. Ela mostra alguns pontos como cartões. O dono aprova ou rejeita cada cartão e pode acrescentar um comentário. A página não tem servidor. Um controle "Copiar decisões" produz um bloco de Markdown. O dono cola esse bloco no chat.

## 3. Regras

### 3.1 Configuração `worktree`

- R1 Quando `dh config set worktree always` ou `dh config set worktree on-request` roda, então o comando grava o valor em `<home>/config.yaml`.
- R2 Quando `dh config set worktree <value>` roda com qualquer outro valor, então o comando sai com um erro que lista os dois valores válidos. Ele deixa o `config.yaml` sem alteração.
- R3 Quando `dh config unset worktree` roda, então a chave sai de `<home>/config.yaml`. O valor global efetivo passa a ser `on-request`.
- R4 Quando `dh config` roda, então a saída lista `worktree` com seu valor gravado. Quando `dh config --json` roda, então o JSON tem um campo `worktree`. Nos dois casos, uma chave ausente mostra um valor vazio, como `language`, `mode` e `dashboard.sprites` fazem. O padrão `on-request` vem de R5 e R7.
- R5 Quando `dh harness-path --json` roda, então o objeto `defaults` tem o campo `worktree`. O campo guarda o valor global, ou `on-request` quando a chave está ausente. Esse campo não inclui o valor do `local.yaml`.
- R6 A chave `worktree` em `<harness dir>/local.yaml` sobrescreve o valor global para aquele projeto naquela máquina. Um valor válido é `always` ou `on-request`.
- R7 O valor efetivo é o primeiro encontrado nesta ordem: `local.yaml`, depois a configuração global, depois `on-request`.
- R8 Quando o `local.yaml` tem um valor inválido de `worktree`, então o Coordenador avisa o dono. O Coordenador ignora essa fonte e usa a próxima na ordem de R7. O mesmo vale para um valor inválido que o dono escreveu à mão no `config.yaml` global. O código em Go não confere o valor gravado de `worktree` contra `always` e `on-request`. Um valor que quebra a regra de escalar do YAML aborta a leitura, como acontece com `language`. Essa regra rejeita um valor sem aspas que tem `:` ou começa com `{`, `[`, `|`, `>`, `&`, `*` ou `!`.
- R9 O Coordenador lê o `local.yaml` como texto simples. A mudança não acrescenta analisador em Go para o `local.yaml`. O arquivo `.skills/external-clis/references/local.yaml.example` documenta a chave `worktree`.

### 3.2 Comportamento no `/dh:build` e no `/dh:auto`

- R10 Quando o valor efetivo é `always`, então `/dh:build` e `/dh:auto` rodam o desenvolvimento da tarefa ou da execução dentro de um worktree do git.
- R11 Quando o valor efetivo é `on-request`, então esses comandos criam um worktree só quando o dono pede um no pedido.
- R12 O Coordenador cria o worktree na sessão principal. Todo despacho de um builder ou de outro especialista em uma execução com worktree declara no seu texto que o especialista não deve rodar `git worktree add`, `git worktree remove` nem `git worktree prune`. Hoje não existe bloqueio de runtime para esses comandos.
- R13 O Coordenador cria um worktree para cada execução (`/dh:auto`) ou cada tarefa (`/dh:build`). Ele não cria worktree para cada fatia.
- R14 Quando um worktree existe, então o Coordenador despacha cada `builder` com o caminho absoluto do worktree. O `qa-verifier` e o `reviewer` também trabalham nesse caminho.
- R15 No `/dh:auto`, o `<base-branch>` único da execução é o branch do worktree. O comando não cria um segundo branch.
- R16 O Coordenador registra o caminho e o branch do worktree no registro da tarefa, para o `/dh:resume` os encontrar. O registro da tarefa é o `TASK.md` no `/dh:build` e o `BRIEF.md` no `/dh:auto`.
- R17 Quando o Coordenador não consegue criar o worktree, então ele reporta a causa e para. Ele não volta ao checkout principal sem a resposta do dono.

### 3.3 Resolução do harness e registros

- R18 Quando `dh harness-path --json` roda dentro de um worktree vinculado válido (R26) de um projeto vinculado, então os campos `dir` e `pdocs` são iguais aos valores que ele devolve no checkout principal. A chave dos dois campos é o caminho do checkout principal. Isso vale no modo repo e no modo global.
- R19 Quando um worktree vinculado contém um diretório `.harness/` (modo repo, versionado), então a resolução ignora essa cópia. Ela devolve o harness dir do checkout principal.
- R20 No modo repo, o Coordenador continua gravando registros no harness dir do checkout principal. O `builder` muda só os arquivos do seu conjunto de escrita no worktree. Nenhum agente grava registros na cópia de `.harness/` dentro do worktree.
- R21 Quando `dh harness-path` roda em um caminho que não é um worktree vinculado nem um projeto vinculado, então o resultado não muda em relação a hoje.

### 3.4 Limpeza do worktree

- R22 O Coordenador remove o worktree só depois que o pull request é mergeado e a tag existe. No `/dh:auto`, a remoção é o último passo de `Stage:`, `cleanup`, que roda depois de `tag`.
- R23 Quando o worktree tem mudanças não commitadas, então o Coordenador não o remove. Ele reporta o caminho e o motivo ao dono.
- R24 Quando o branch do worktree tem commits que o remoto não tem, então o Coordenador não o remove. Ele reporta o caminho e o motivo ao dono. Um merge com squash que apaga o branch remoto deixa os commits em nenhuma ref remota. Os commits contam como enviados quando o `HEAD` do worktree é igual ao `headRefOid` do pull request mergeado. A remoção nunca usa `--force`.
- R25 Quando a execução para antes do merge e da tag, então o Coordenador deixa o worktree no lugar e registra seu caminho no registro da tarefa. Quando uma execução do `/dh:auto` para durante o `cleanup`, então `Stage: cleanup` fica no `BRIEF.md` para o `/dh:resume`.

### 3.5 Isolamento do revisor de CLI a partir de um worktree

- R26 Quando `<project>/.git` é um diretório, então valem as checagens de isolamento existentes e nenhuma outra checagem. Quando não é um diretório, então a chamada precisa de um `linked_worktree` não nulo (R61). O `.git` também deve apontar para um git dir de um layout de worktree válido: `<common dir>/worktrees/<name>`, com o `commondir` e o ponteiro de volta correspondentes. Logo antes do clone, o `.git` deve ser um arquivo regular (`test -f`) e não um symlink (`test ! -L`).
- R27 Quando o `.git` é um symlink, então a chamada termina como `isolation unavailable`. O mesmo vale quando `linked_worktree` é `null` e o `.git` não é um diretório. O mesmo vale quando o `.git` aponta para fora de um layout de worktree válido. O mesmo vale quando o `.git` não é um arquivo regular na checagem antes do clone.
- R28 O sandbox da CLI é sempre um clone novo feito com `git clone --no-hardlinks` em uma pasta temporária descartável. O kit nunca usa o worktree como sandbox. O kit nunca usa `git worktree` para o sandbox.
- R29 O clone tem seu próprio diretório `.git`. Ele não tem remoto `origin`. Ele não tem hardlink para a origem. Essas checagens são as mesmas de um checkout normal.
- R30 O congelamento da árvore viva cobre duas partes. Um item ausente nos dois instantâneos não é uma diferença e nunca é um erro. Em um checkout normal, o diretório comum é o diretório `.git`.
  - Parte (a), todo projeto: a cobertura existente de R5, e o harness dir resolvido com sua cobertura em bytes. Em um worktree vinculado, o harness dir resolvido (R18) não é `<project>/.harness/`, que é só a cópia dentro do worktree. No modo global, o harness dir continua fora do congelamento. Também no diretório comum: os arquivos sob `refs/`, e as entradas sob `worktrees/` quando existem, com os arquivos de cada entrada que a parte (b) lista.
  - Parte (b), só quando `<project>/.git` é um gitfile (worktree vinculado): os bytes do gitfile e, no diretório administrativo do worktree, `HEAD`, `commondir`, `gitdir`, `config.worktree`, os arquivos sob `refs/` e a presença de `locked`.
  - Os arquivos `index` brutos ficam fora do congelamento.
- R31 Quando o congelamento detecta uma mudança na árvore ou nos git dirs de R30 durante uma chamada de CLI, então a chamada descarta o veredito. Um checkout normal se comporta da mesma forma.
- R32 O Coordenador não grava nada no harness dir nem na árvore enquanto uma chamada de CLI está em andamento. Esta é a regra existente do RISK-002. Ela continua.
- R33 A cláusula anti-delegação continua literal em todo prompt de CLI externa (RISK-001).
- R34 O RISK-002 continua "mitigated". Nenhum arquivo da mudança o chama de "resolved".

### 3.6 Texto desatualizado

- R35 Os arquivos `docs/roadmap.html` e `docs/en/roadmap.html` não prometem isolamento por `git worktree` para os revisores de CLI. O texto declara a decisão do clone.
- R36 O `AGENTS.md` nomeia a configuração `worktree` com suas três fontes e seu padrão. Ele declara uma regra: um worktree é permitido para o desenvolvimento, e o sandbox da CLI continua sendo um clone. R30 reduz os resíduos aceitos `"paths of the live .git outside the R5 snapshot"`. O `AGENTS.md` lista os resíduos que permanecem. O ADR-008 registra quais resíduos diminuem. A mudança não edita o ADR-002 nem o ADR-003. Ela mantém o resto do `AGENTS.md` como está.

### 3.7 Página de aprovação

- R37 O `docs-guide` gera a página com a skill `doc-template-html`. A skill ganha um `--type` novo para a página de aprovação. A mudança atualiza `.skills/doc-template-html/scripts/stamp.sh` e `references/types.md`. Ela acrescenta um esqueleto em `assets/skeletons/pt-br/` e um em `assets/skeletons/en/`.
- R38 O `/dh:plan` despacha o `docs-guide` para a página quando o plano tem pelo menos um ponto de decisão. Quando o plano não tem nenhum, o `/dh:plan` não gera página.
- R39 A página lista todos os passos do trabalho em ordem. Cada passo tem um título e uma descrição.
- R40 A página mostra os trechos propostos para a feature. Cada trecho fica em um bloco de código e nomeia o passo a que pertence.
- R41 A página mostra cada ponto de decisão como um cartão. Cada cartão tem um id, um título, uma descrição curta, um controle Aprovar e um controle Rejeitar.
- R42 Cada cartão tem um campo de comentário. O comentário é opcional quando o dono aprova. O comentário é obrigatório quando o dono rejeita.
- R43 Quando o dono rejeita um cartão com o comentário vazio (só caracteres em branco), então a página mostra uma mensagem visível nesse cartão. A página não inclui esse cartão como decidido no bloco copiado.
- R44 O controle "Copiar decisões" produz um bloco de Markdown. Cada cartão dá uma entrada: id do cartão, decisão (`approved` ou `rejected`) e comentário. Um cartão sem decisão aparece como `pending`.
- R45 Quando a área de transferência não está disponível ou o navegador recusa a cópia, então a página mostra o bloco de Markdown em uma área de texto. O dono seleciona e copia à mão.
- R46 A página não faz escrita no servidor. A mudança não acrescenta endpoint de escrita no dashboard. Nenhum processo lê as decisões automaticamente. O dono cola o bloco no chat.
- R47 A página guarda as escolhas em andamento no armazenamento do navegador. Toda chamada ao armazenamento fica dentro de `try/catch`. A página funciona com o armazenamento bloqueado ou vazio.
- R48 A página escapa todo texto inserido: texto do plano, trechos e o campo de comentário. Um texto como `<script>` ou `<img onerror=...>` aparece como texto simples e não executa código. A página não usa `innerHTML` com texto sem escape.
- R49 A página não carrega nenhum script externo. As únicas origens externas são o CDN fixado do Mermaid e o Google Fonts, que já valem para os pdocs.
- R50 A página funciona só com o teclado. Todo controle tem um estilo de foco visível e um rótulo de texto. Aprovar e Rejeitar não são mostrados só pela cor.
- R51 A página não tem rolagem horizontal em uma largura de celular de 360 px.
- R52 A página é um arquivo `.html` plano na pasta de pdocs. O caminho da pasta vem de `dh harness-path --json`, campo `pdocs`. Quando esse campo está vazio (um projeto sem entrada no registro), o Coordenador reporta a causa e não gera página. Dentro de um worktree vinculado, esse campo resolve a partir do checkout principal (R18). O arquivo respeita o limite de 8 MiB do dashboard.
- R53 Quando `<pdocs>/index.html` existe, então a geração o atualiza com um link para a página.
- R54 O idioma da página segue o `language` do projeto (`pt-br` ou `en`).

### 3.8 Entrega

- R55 A versão é 0.33.0. O `CHANGELOG.md` tem a entrada. O `dist/` é reconstruído com `go run ./cmd/dh build`.
- R56 Um ADR novo (ADR-008) registra a decisão do worktree. Ele declara as três fontes da configuração, quem cria o worktree e por que o sandbox da CLI continua sendo um clone.
- R57 A mudança atualiza os documentos da seção 4 dentro do mesmo pull request.
- R58 A contagem continua em 19 arquivos de comando mais 1 comando de mod. A mudança não acrescenta nenhum comando `/dh:`.
- R59 Como a mudança edita instruções executáveis do harness que sustentam um controle de segurança (R26 a R34), o lote precisa da etapa de revisão `security`. A etapa despacha a lista completa configurada em `reviewers.security`. O lote também precisa da revisão de código padrão e da matriz do `qa-verifier`.
- R60 Todo arquivo novo usa caminhos relativos ou `~/...`. O `dh validate` não acha nenhum caminho privado absoluto.

### 3.9 Acrescentado pela emenda após o ADR-008

- R61 O JSON de `dh harness-path --json` ganha um campo aditivo `linked_worktree`. Fora de um worktree vinculado válido, ele é `null`. Dentro de um, é um objeto com três campos (`main`, `gitdir` e `commondir`), cada um um caminho canônico.
  - `main`: o checkout principal.
  - `gitdir`: o diretório administrativo do worktree `<common dir>/worktrees/<name>`. Não é o gitfile.
  - `commondir`: o diretório comum.
  - Os campos existentes mantêm seu significado. O passo de isolamento (R26, R30) lê o diretório administrativo de `gitdir`.
- R62 Para o `/dh:build`, o branch do worktree é o `<task-id>` em minúsculas. Para o `/dh:auto`, vale R15. O Coordenador cria o worktree na sessão principal com `git worktree add -b`. Ele não usa a ferramenta `EnterWorktree` do Claude Code.
- R63 A pasta do worktree é `<parent of the main checkout>/<checkout folder name>.worktrees/<task-id>`. O `<task-id>` é o id da tarefa, ou `AUTO-<n>` para o `/dh:auto`. Ele tem só os caracteres `[A-Za-z0-9_-]`. Qualquer outro caractere encerra a execução como em R17. A pasta nunca fica sob `<home>`. O registro da tarefa (R16) guarda o caminho relativo ao checkout principal. O despacho ao builder (R14) guarda o caminho absoluto.
- R64 Quando o projeto é um worktree vinculado e a chamada é a revisão de árvore viva do `agy`, então a chamada termina como `isolation unavailable`. A revisão não roda, e nenhuma rodada de revisão é gasta. O motivo é que seu congelamento tira instantâneos de `.git/hooks/` e `.git/config`, que não são pastas em um checkout com gitfile.

## 4. Docs

- README.md
- README.en.md
- docs/tutorial.html
- docs/en/tutorial.html
- docs/roadmap.html
- docs/en/roadmap.html
- docs/adr/ADR-008-worktree-de-desenvolvimento.md
- AGENTS.md
- CHANGELOG.md
- .commands/build.md
- .commands/auto.md
- .commands/plan.md
- .skills/external-clis/references/clone-isolation.md
- .skills/external-clis/references/live-tree-freeze.md
- .skills/external-clis/references/local.yaml.example
- .skills/doc-template-html/references/types.md

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Um worktree para cada fatia | Mais worktrees para criar e limpar. O dono decidiu um para cada execução ou tarefa. |
| Um subagente cria o worktree | A sessão principal detém a capacidade de worktree e o registro da tarefa. Uma regra do prompt mantém os subagentes longe dos comandos de worktree (R12). |
| Usar o worktree como sandbox da CLI | Um worktree compartilha todo o `.git` do checkout principal. Uma CLI poderia escrever ali. A decisão do clone continua. |
| Continuar recusando um worktree no isolamento da CLI (opção b) | O dono escolheu a opção (a): aceitar um worktree com as travas da seção 3.5. |
| Um endpoint do dashboard que recebe as decisões | Fora de escopo por premissa do dono. Ele acrescenta uma superfície de escrita. |
| Um analisador novo em Go para o `local.yaml` | O Coordenador lê duas fontes em texto simples. Nenhuma necessidade encontrada. |

### Premissas que o dono pode derrubar

- A1 O `/dh:plan` é o único comando que gera a página. `/dh:plan-loop` e `/dh:document` estão fora de escopo.
- A2 A configuração vale para `/dh:build` e `/dh:auto`. `/dh:fix` e `/dh:refactor` estão fora de escopo.
- A3 O canal de retorno é o bloco de Markdown de "Copiar decisões" que o dono cola no chat (R44 a R46).
- A4 A página de aprovação é um `--type` novo do `doc-template-html`, não uma skill nova e não um comando novo.
- A5 O worktree fica fora da árvore do checkout principal. O ADR-008 decide a pasta exata, o nome do branch e o método de criação (R62, R63).
- A7 Este lote não muda a trava de Bash do PRD-027. Uma mudança na trava também negaria formas de leitura como `git worktree list`. Ela precisa de sua própria revisão de segurança. R12 usa só uma regra de prompt.
- A8 A atualização de `<pdocs>/index.html` (R53) segue a regra existente do Coordenador. O `docs-guide` atualiza esse arquivo quando ele existe. Só `/dh:document pdocs` o cria.
- A6 O nome do tipo é escolhido pelo planejador. Ele segue os nomes existentes no `stamp.sh`.

### Fora de escopo

- Qualquer endpoint de escrita do dashboard, leitura automática das decisões e sincronização multiusuário.
- Tratamento de worktree para `/dh:fix`, `/dh:refactor` e qualquer comando não nomeado em R10 e R11.
- Uma mudança na decisão do sandbox da CLI (o clone continua).
- Casos de eval e sondas novas.

### Riscos e decisões pendentes

- Risco: R26 a R31 ampliam um controle de segurança · mitigação: R59 (etapa completa de segurança), os casos de recusa de R27 e os casos de QA para um gitfile symlink e um gitfile que aponta para fora de um layout válido.
- Risco: um `local.yaml` editado à mão tem um valor ruim · mitigação: R8.
- Decidido no ADR-008: R18 e R19 resolvem lendo o gitfile, o `commondir` e o ponteiro de volta em Go, sem `git`. O comportamento de R18 a R21 é o requisito.
- Candidato (reforço opcional): acrescentar `worktree` à lista de escrita do git da trava, com sua própria revisão de segurança · respondido por: dono · bloqueia: nada neste lote.
- Decidido no ADR-008: pasta do worktree (A5), nome do branch e método de criação (R62, R63). Nada pende aqui.
- Decisão pendente: nome do novo `--type` (A6) · respondido por: planejador · bloqueia: o ticket da página.
- Decidido no ADR-008: a sessão principal cria o worktree com `git worktree add -b`, não com a ferramenta de worktree do harness (R62).
- Ponto em aberto: não está verificado se o Claude Code pede permissão de escrita na pasta irmã `<checkout folder name>.worktrees` fora do modo bypass · respondido por: uma sonda do `qa-verifier`, ou uma nota de configuração nos docs · bloqueia: o ticket de docs. Este PRD não afirma nenhum resultado.

## Resumo executado

- **Entregue:** A configuração `worktree` (`always` ou `on-request`, padrão `on-request`) agora define o fluxo de worktree do `/dh:build`, do `/dh:auto` e do `/dh:resume`, um resolvedor em Go lê o harness dir e o `pdocs` a partir do checkout principal de um worktree vinculado, o isolamento do revisor de CLI aceita um worktree vinculado com um congelamento mais amplo, e o `/dh:plan` escreve uma página HTML de aprovação (tipo `aprovacao`), entregue como versão 0.33.0.
- **Regras:**

  | Regra | Veredito | Evidência |
  |---|---|---|
  | R1 | Honrada | T-2801, T-2803; teste em Go `TestR1_R2_R3_ConfigSetUnsetWorktree` |
  | R2 | Honrada | T-2801, T-2803; teste em Go `TestR1_R2_R3_ConfigSetUnsetWorktree` (caso de recusa) |
  | R3 | Honrada | T-2801, T-2803; teste em Go `TestR1_R2_R3_ConfigSetUnsetWorktree` |
  | R4 | Honrada | T-2801, T-2803; teste em Go `TestR4_ConfigShowsWorktree` |
  | R5 | Honrada | T-2801, T-2803; teste em Go `TestR5_R7_DefaultsWorktree` |
  | R6 | Honrada | T-2806, T-2807; checagem de prosa (fraca, ver nota e) |
  | R7 | Honrada | T-2803, T-2807; teste em Go `TestR5_R7_DefaultsWorktree`; checagem de prosa (fraca) |
  | R8 | Honrada | T-2803, T-2807; teste em Go `TestR8_HarnessPathScalarRules`; checagem de prosa (fraca) |
  | R9 | Honrada | T-2807, T-2814; `grep` do diff, nenhuma leitura de `local.yaml` no código em Go (fraca) |
  | R10 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R11 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R12 | Honrada | T-2806, T-2807; checagem de prosa, sem bloqueio de runtime de propósito (fraca) |
  | R13 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R14 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R15 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R16 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R17 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R18 | Honrada | T-2802, T-2803; teste em Go `TestR18_R19_RepoModeResolvesMainHarness`, caso com `git worktree add` real |
  | R19 | Honrada | T-2802, T-2803; teste em Go `TestR18_R19_RepoModeResolvesMainHarness`; mutação M3 detectada |
  | R20 | Honrada | T-2806, T-2807; checagem de prosa (fraca) |
  | R21 | Honrada | T-2802, T-2803; teste em Go `TestR21_NormalCheckoutAndPlainDirAreNotWorktrees` |
  | R22 | Honrada | T-2806, T-2807; checagens de prosa e da ordem da limpeza (fraca) |
  | R23 | Honrada | T-2806, T-2807; checagens de prosa e da ordem da limpeza (fraca) |
  | R24 | Honrada | T-2806, T-2807; checagens de prosa e da ordem da limpeza (fraca) |
  | R25 | Honrada | T-2806, T-2807; checagens de prosa e da ordem da limpeza (fraca) |
  | R26 | Honrada | T-2802, T-2803, T-2805; testes em Go e matriz de shell (11 PASS) |
  | R27 | Honrada | T-2802, T-2803, T-2805; testes em Go e matriz de shell (26 PASS) |
  | R28 | Honrada | T-2805; matriz de shell, clone com `--no-hardlinks` (7 PASS) |
  | R29 | Honrada | T-2805; matriz de shell (3 PASS) |
  | R30 | Honrada | T-2805, T-2804; matriz de shell (55 PASS) e rodada 4 do T-2804 |
  | R31 | Honrada | T-2805; matriz de shell (1 PASS, linha única) |
  | R32 | Honrada | T-2805; matriz de shell (5 PASS) e checksum da cláusula |
  | R33 | Honrada | T-2805; matriz de shell (5 PASS) e checksum da cláusula |
  | R34 | Honrada | T-2805, T-2814; matriz de shell (11 PASS) |
  | R35 | Honrada | T-2811; rodada 3 da matriz e nova checagem manual |
  | R36 | Honrada | T-2811; linhas da matriz e nova checagem manual (ver nota c) |
  | R37 | Honrada | T-2808, T-2810; matriz R37 PASS (ver nota a) |
  | R38 | Honrada | T-2809, T-2810; asserção de leitura (fraca) |
  | R39 | Honrada | T-2808, T-2810; matriz do navegador headless (156 PASS) |
  | R40 | Honrada | T-2808, T-2810; matriz do navegador headless |
  | R41 | Honrada | T-2808, T-2810; matriz do navegador headless |
  | R42 | Honrada | T-2808, T-2810; matriz do navegador headless (8 PASS) |
  | R43 | Honrada | T-2808, T-2810; matriz do navegador headless (8 PASS) |
  | R44 | Honrada | T-2808, T-2810; matriz do navegador headless (16 PASS), ver nota b |
  | R45 | Honrada | T-2808, T-2810; matriz do navegador headless (12 PASS) |
  | R46 | Honrada | T-2810, T-2814; log de rede sem escrita, `git status` limpo |
  | R47 | Honrada | T-2810; matriz do navegador headless (origem opaca, armazenamento, área de transferência) |
  | R48 | Honrada | T-2808, T-2810; matriz do navegador headless (18 PASS), varredura do código-fonte com 0 pontos de risco; ver nota c para o modelo de ameaças |
  | R49 | Honrada | T-2810; matriz do navegador headless (nenhum `src` de script) |
  | R50 | Honrada | T-2810; matriz do navegador headless (28 PASS, teclado e foco) |
  | R51 | Honrada | T-2808, T-2810; matriz do navegador headless (12 PASS) e checagem do limite de saída |
  | R52 | Honrada | T-2809, T-2810; checagem de rota (200) e matriz do T-2809 |
  | R53 | Honrada | T-2809, T-2810; asserção de leitura, nunca rodado de ponta a ponta (fraca) |
  | R54 | Honrada | T-2808, T-2810; matriz do navegador headless (8 PASS, pt-BR e en) |
  | R55 | Honrada | T-2813, T-2814; checagem de prosa, `go test ./internal/kit` passa (fraca: sem pasta de evidência do T-2813) |
  | R56 | Honrada | T-2814; checagem de prosa, status do ADR-008 accepted |
  | R57 | Honrada | T-2812, T-2814; checagem de `git status` dos 16 arquivos da seção 4 |
  | R58 | Honrada | T-2814; log do validate (19 comandos, 13 agentes) |
  | R59 | Honrada | T-2814; passes de segurança (ver nota d) |
  | R60 | Honrada | T-2814; log do validate (0 erros) e checagem de prosa de caminhos privados |
  | R61 | Honrada | T-2802, T-2803; testes em Go `TestR18_R61_GlobalModeLinkedWorktree` e casos relacionados |
  | R62 | Honrada | T-2806, T-2807; checagem de prosa da ordem dos comandos citados (fraca) |
  | R63 | Honrada | T-2806, T-2807; checagem de prosa da ordem dos comandos citados (fraca) |
  | R64 | Honrada | T-2804, T-2805; matriz de shell (17 PASS), não rodada de novo na passada final |

  Notas sobre a tabela:
  - (a) R37 é Honrada por uma decisão registrada do Coordenador. O `docs-guide` não tem shell, então a página é montada em três despachos: o `docs-guide` escreve o JSON de dados, um `builder` na `lane` `infra` roda o `stamp.sh`, e o `docs-guide` acrescenta o link no índice. O texto de R37 diz que o `docs-guide` gera a página.
  - (b) O bloco copiado também guarda o título do cartão, achatado em uma linha (R44).
  - (c) O `AGENTS.md` também registra a decisão do dono de 2026-10-28 sobre o modelo de ameaças (ferramenta local pessoal). Isso vai além do trecho de R36.
  - (d) R59 é Honrada. A lista completa de `reviewers.security` (`cli:codex/gpt-6.1-sol` e `claude`) rodou em três passes, mais um passe para a página de aprovação. Os achados foram corrigidos e verificados de novo. Os relatórios do `claude` ficam literais em `~/.harness/projects/dev-harness/tasks/PRD-028/evidence/sec1` a `sec4`. Os relatórios do codex estão em `sec2` a `sec4`. A saída bruta do codex do primeiro passe não foi guardada; só o seu resumo está no registro do Coordenador. Nenhum passe de segurança cobre o estado final, porque o dono limitou o modelo de ameaças em 2026-10-28. A revisão de código padrão e a matriz do `qa-verifier` rodaram para cada ticket. A revisão final da feature rodou uma vez.
  - (e) Evidência fraca, dita com clareza: R6 a R17, R20, R22 a R25, R62 e R63 se apoiam em checagens de prosa somente leitura. O T-2807 deixou saídas brutas e nenhuma matriz de vereditos. O QA final atualizou as checagens, e elas passam.

- **Tickets:** Todos os tickets estão done. T-2801 a T-2807 ficaram done em 2026-10-27; os outros em 2026-10-28. Os caminhos de evidência abaixo ficam sob `~/.harness/projects/dev-harness/tasks/`.
  - T-2801, `lane` `backend` (chave de configuração em Go), done. Evidência: `PRD-028/evidence/final/rtests-verbose.log`.
  - T-2802, `lane` `backend` (resolvedor em Go), done. Evidência: `PRD-028/evidence/final/rtests-verbose.log`.
  - T-2803, `lane` `test` (testes em Go), done. Evidência: `T-2803/evidence/` (por exemplo `02-new-tests-verbose.txt`, execuções de mutação M2 e M3).
  - T-2804, `lane` `backend` (prosa de isolamento e congelamento), done. Evidência: `T-2804/evidence/w7/matrix.md`.
  - T-2805, `lane` `test` (isolamento com um worktree vinculado real), done. Evidência: `T-2805/evidence/w7/matrix.md`.
  - T-2806, `lane` `backend` (prosa de build, auto, coordinator e resume), done. Evidência: `T-2806/evidence/w5/matrix.md`, `w7/matrix.md`.
  - T-2807, `lane` `test` (checagens de prosa), done. Evidência: só saídas brutas (`T-2807/evidence/run1.txt` a `run3-text.txt`); nenhuma matriz de vereditos.
  - T-2808, `lane` `backend` (tipo da página de aprovação, stamp, esqueletos), done. Evidência: a matriz da rodada 3 citada em `PRD-028/evidence/final/coverage.md` (linha de R36).
  - T-2809, `lane` `backend` (gancho do `/dh:plan` e `docs-guide`), done. Evidência: `T-2809/evidence/matrix.md`.
  - T-2810, `lane` `test` (navegador headless): 156 casos principais e 49 casos extras, todos PASS. Evidência: `T-2810/evidence/matrix.md`, `T-2810/evidence/screenshots/`.
  - T-2811, `lane` `backend` (texto do roadmap, `AGENTS.md`, frase do ADR-008), done. Evidência: `T-2811/evidence/matrix.md`.
  - T-2812, `lane` `backend` (docs), done. Evidência: checagem de `git status` (R57); sem pasta de evidência.
  - T-2813, `lane` `backend` (versão 0.33.0 e `dist/`), done. Evidência: sem pasta de evidência; a checagem do `dist/` é coberta pelo T-2814.
  - T-2814, `lane` `test` (QA final): 21 PASS, 0 FAIL, 1 `NOT-RUN`, julgado pelo Coordenador. A prova de sha256 em clone limpo dos cinco binários e dos dois launchers passou. Evidência: `PRD-028/evidence/final/` (`coverage.md`, `validate.log`, `prose-recheck.txt`, `go-test.log`).

- **Docs:** `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md`, `docs/adr/ADR-008-worktree-de-desenvolvimento.md` (novo, com seu registro de revisão), este PRD e seu registro de revisão `PRD-028-worktree-and-approval-page.review.md`, `.agents/coordinator.md`, `.agents/docs-guide.md`, `.commands/build.md`, `.commands/auto.md`, `.commands/resume.md`, `.commands/plan.md`, `.skills/external-clis/` (`SKILL.md` e referências) e `.skills/doc-template-html/` (`SKILL.md`, `README.md`, `references/types.md`, `stamp.sh`, dois esqueletos).

- **Fora de escopo:** Os itens do Apêndice: qualquer endpoint de escrita do dashboard, leitura automática das decisões, sincronização multiusuário; tratamento de worktree para `/dh:fix`, `/dh:refactor` e outros comandos não nomeados em R10 e R11; uma mudança na decisão do sandbox da CLI (o clone continua); casos de eval e sondas novas. Também fora de escopo:
  - Nenhuma mudança na trava de runtime do PRD-027 para `git worktree`. A proibição para os especialistas é só prosa.
  - Nenhum endpoint de escrita do dashboard. A página devolve as decisões por um bloco de Markdown colado.
  - Nenhum agente, skill ou comando novo.
  - Lotes candidatos achados durante o trabalho, não iniciados: um falso positivo da trava em `git p4` em um padrão de busca; reforço da trava para `git worktree`; aplicar a ordenação da Parte A ao passo 7 do caminho de árvore viva do `agy`; uma passada de higiene nos marcadores `<path>` sem aspas nos comandos de reversão de `live-tree-freeze.md`; rejeitar ids de decisão duplicados em páginas mais antigas; lições de receita para o `opencode` (`--standalone`, MCP do codegraph desligado) e o `codex` (`-c mcp_servers.codegraph.enabled=false`).
