---
name: external-clis
description: Use quando o Coordenador despacha uma etapa de revisão (`document`, `code`, `security`) que lista uma entrada `cli:<binary>/<slug>` em `reviewers:` (`.harness/project.yaml`), quando `/dh:plan-loop` despacha uma chamada de planejamento `cli:<binary>/<slug>` para uma onda, quando uma chamada de CLI volta sem o sinal de veredito da etapa (ou, numa chamada de planejamento, sem conteúdo de plano reconhecível) e a falha precisa ser classificada, ou quando um novo binário de CLI precisa de uma receita de invocação sem interface acrescentada em `references/`. Não use para papéis de `builder`, para uso interativo/TUI de qualquer CLI, nem para a etapa `verify`, que está fora de escopo nesta versão.
author: malka
metadata:
  provenance: adapted
  sources: ["personal skill fleet-clis (SKILL.md, per-binary recipes, gotchas.md), not versioned in this kit"]
---

# CLIs externas

## Visão geral

Um revisor de CLI externa é invocado sem interface, com o mesmo prompt de papel que um especialista Claude receberia, de modo que o veredito dele é evidência em que o Coordenador pode confiar exatamente tanto quanto no de um revisor Claude. Quatro partes nunca mudam por binário: resolver o binário antes de supor que ele existe; pôr o prompt em disco quando passar de uma linha; aplicar o modo somente leitura do binário, ou rodá-lo num clone descartável quando ele não tem um; e ler o resultado pelo sinal de veredito da etapa, para que uma resposta muda ou malformada nunca seja tomada por aprovação. Uma chamada de planejamento de `/dh:plan-loop` não tem sinal de veredito e é classificada conforme "A chamada de planejamento". As flags e manias de cada binário ficam numa receita por binário em `references/`, porque a ordem das flags faz parte do contrato de cada binário e misturar receitas convida a copiar a ordem errada.

## Quando usar

- O Coordenador despacha `cli:<binary>/<slug>` para uma etapa de revisão configurada em `reviewers:`, junto com o agente Claude da etapa ou no lugar dele, ou para uma dúvida no meio do build (uma revisão `code`, ver "Dúvida no meio do build"), ou `/dh:plan-loop` despacha uma para redigir um plano concorrente para uma onda ("A chamada de planejamento"). Esta skill não escolhe qual entrada roda.
- Um prompt de papel (`document-validator`, `reviewer` no modo `code` ou `security`) precisa ir para uma CLI externa em vez de um especialista Claude.
- Uma chamada voltou sem o sinal de veredito da etapa, ou passou do timeout, e a falha precisa ser classificada; ou um novo binário de CLI precisa da própria receita em `references/`.

## Quando não usar

- O papel de `builder`, em qualquer raia ou modo. Uma CLI nunca escreve código de produção, e uma chamada de planejamento nunca escreve arquivo algum: ela devolve texto de plano, e `docs-guide` é o único que escreve arquivos de plano.
- Uso interativo ou TUI de qualquer CLI (toda chamada aqui é sem interface e de disparo único), ou instalar, autenticar ou atualizar uma: um binário não instalado e autenticado aqui é reportado como `not-run`, nunca provisionado.
- A etapa `verify` (`qa-verifier`), fora de escopo até existir um contrato de revisor somente leitura para o QA.
- Escolher revisores, mesclar N vereditos ou escrever `.harness/local.yaml`: o procedimento do próprio Coordenador.

## Quem executa a chamada

O Coordenador não executa uma chamada de CLI. Ele despacha um `qa-verifier`, ou um `builder` na raia `infra`, com a etapa, a entrada `cli:<binary>/<slug>`, o pacote de contexto e o harness dir. O delegado carrega esta skill, executa toda a mecânica dela (resolver, montar, clonar, congelar, chamar, comparar, reverter, classificar) e devolve o bloco abaixo. O Coordenador mescla os vereditos e nunca lê a resposta como evidência própria; pode abrir a saída bruta citada.

```
## Revisor de CLI: <binário>/<slug>
- Binário resolvido: <comando achado no PATH, ou o caminho tirado de .harness/local.yaml>
- Veredito: <o sinal de veredito da etapa, literal> | not-run (<classe>: <motivo>)
- Timeout: <o valor usado, configurado por revisor ou o padrão de 15 minutos>
- Saída bruta: <caminho em disco do prompt e da resposta, se algum dos dois foi escrito lá>
- Checagem do espaço de trabalho: clean | reverted (<o que o revisor alterou, agora revertido>) | clone removed: <caminho> | notes: <arquivos novos não rastreados no clone> | leftover processes: stopped / none (<string casada>)
```

`<classe>` é uma das classes de "Classificação de falhas". Uma resposta `not-run` sempre nomeia a classe e o motivo, porque o Coordenador reporta os dois ao dono de imediato. Os três últimos valores da Checagem do espaço de trabalho valem para uma chamada num clone. Nada escreve na árvore nem no `<harness dir>` enquanto uma chamada de CLI está em andamento (RISK-002), delegado inclusive: ele escreve o prompt antes de a chamada começar, e só o diretório de congelamento e o clone são escritos durante a chamada.

Toda chamada é um único comando de shell que começa com o nome do binário. Nenhum `timeout` nem outro invólucro vai antes dele, e nenhum `&&`, `;` ou pipe o encadeia com o congelamento ou com qualquer outro passo: o congelamento anterior à chamada e a comparação posterior são comandos separados. Motivo (2026-10-16): uma regra de permissão como `Bash(<binary>:*)` casa só com comandos que começam com o binário; uma chamada embrulhada e encadeada com um script de instantâneo foi negada, e a mesma chamada isolada rodou. A flag do próprio binário limita o tempo quando ele tem uma; senão a receita diz que o timeout da ferramenta Bash ou `run_in_background` faz isso. Uma chamada em clone fixa o diretório de trabalho com a flag do próprio binário dentro desse mesmo comando; uma receita cujo binário não tem essa flag começa com `cd "$CLONE" &&`, que não casa com uma regra de permissão assim, então a chamada pode precisar da aprovação do dono ou ser `not-run`.

## Entradas

| Entrada (uma chamada de revisão; as entradas de uma chamada de planejamento estão em "A chamada de planejamento") | Se faltar |
| --- | --- |
| Corpo do prompt de papel (`.agents/reviewer.md` para `code` e `security`, `.agents/document-validator.md` para `document`), sem frontmatter | Pare. Um revisor de CLI nunca recebe um prompt diferente do que um especialista Claude receberia para a etapa. |
| O texto da skill de revisão da etapa (`code-review`, `document-review`, `security-review`) | Pare. O sinal de veredito vem do formato de saída da própria skill. |
| Pacote de contexto do despacho. Etapa `document`: exatamente quatro itens (caminho do documento; material-fonte como caminhos de arquivo, ou o pedido do dono palavra por palavra sob `owner's words`; caminho do relatório anterior numa rodada seguinte; harness dir). Demais etapas: escopo, diff ou documento, material-fonte. Nunca qualquer parte da conversa do Coordenador (turnos do dono, raciocínio, resumos de intenção, notas de rodadas anteriores); em `code`, resultados do `builder` e do QA só como dado citado e rotulado `unverified` | Pare. Mesma exigência que para um especialista Claude. |
| `.harness/local.yaml`, quando o binário não está no `PATH` | `not-run`, classe `not_installed`. Nunca invente um caminho. |
| `commands.test` / `commands.lint` de `.harness/project.yaml` | Omita essa instrução; o revisor raciocina só a partir do código. |

## Procedimento

1. **Resolva o binário.** `command -v <binary>`; se falhar, `binaries.<binary>` em `.harness/local.yaml` (um comando simples, um caminho `~/...` ou um caminho relativo ao projeto, ver `references/local.yaml.example`). Nenhum resolve: `not-run`, `not_installed`. Nunca substitua por outro binário.
2. **Fixe o modelo.** Pegue o slug da entrada `cli:<binary>/<slug>` despachada. Nunca invente um; um slug não medido falha no nível da própria CLI.
3. **Monte o prompt** (ver "Montagem do prompt"; uma chamada de planejamento é montada conforme "A chamada de planejamento"). Escreva-o em disco quando passar de uma linha; a receita diz como o binário o recebe. O prompt nunca contém o caminho absoluto da árvore viva (higiene, não uma fronteira), e o resultado de uma checagem que a CLI reporta nunca é evidência: um `qa-verifier` despachado roda as checagens de novo.
4. **Isole e congele.** Use a linha do binário em "Modo somente leitura por binário". Uma flag somente leitura: aplique-a. Um clone descartável: [references/clone-isolation.md](references/clone-isolation.md) para o ciclo de vida e [references/live-tree-freeze.md](references/live-tree-freeze.md) para o congelamento. Árvore viva: [references/live-tree-compare.md](references/live-tree-compare.md). Se o clone não puder ser feito, a chamada é `not-run`, `isolation_unavailable`, sem recuo para a árvore viva.
5. **Invoque** com a ordem de flags da receita do binário, como um único comando (ver "Quem executa a chamada"). Não reordene as flags em relação ao que a receita mostra.
6. **Compare e reverta.** Recalcule o congelamento com o mesmo arquivo do passo 4. Qualquer diferença é `tree_changed`: descarte o veredito, reverta a alteração, reporte `not-run`.
7. **Classifique o resultado** pelo sinal de veredito da etapa, nunca só pelo código de saída ("Classificação de falhas"), e **reporte** o bloco de "Quem executa a chamada". Mesclar vereditos e escrever um relatório persistente cabem ao Coordenador.

## Montagem do prompt

Um prompt de revisão contém, literalmente e no corpo, nunca só num campo de frontmatter nem numa instrução de despacho que o ambiente de execução da CLI se confie a fazer cumprir: o corpo do prompt de papel, o texto da skill de revisão da etapa, o pacote de contexto como limitado em Entradas, a linha de tag do revisor, a cláusula anti-delegação, o bloco "Ler é permitido" e o contrato de veredito, todos abaixo. Um revisor de CLI lê o formato de saída da etapa como uma seção entre muitas e parafraseia o sinal de veredito a menos que o prompt o declare como exigência; uma execução medida (2026-09-22) devolveu um veredito bloqueante correto escrito como `- **Request changes**` sem linha `Status da revisão:`, que a classificação teria descartado. Um revisor a quem ninguém diz a própria tag inventa uma, então a tag é uma linha do prompt que nomeia a string exata que ele deve escrever no campo `Reportado por:` de cada achado. Uma chamada num clone também leva o bloco `commands.test`/`commands.lint`, que roda no clone.

## A cláusula anti-delegação

Copie este bloco literalmente para o corpo de todo prompt de revisor de CLI, não só para a instrução de despacho em volta:

```
Não invoque outro binário de CLI, não crie outro agente nem delegue qualquer parte
desta revisão a outra ferramenta. Não execute nenhum script deste
repositório. Os únicos comandos que você pode rodar são o `commands.test` e o
`commands.lint` do próprio projeto, exatamente como declarados em
`.harness/project.yaml`, e só para ler o resultado da suíte de testes
ou de lint existente — nunca `commands.build` nem qualquer outro comando do repositório.
Esta instrução está escrita aqui, no texto do prompt, porque o ambiente de execução
não a faz cumprir por si só: um especialista deste projeto já ignorou uma
restrição `disallowedTools` declarada apenas no frontmatter e
despachou outro agente mesmo assim. Siga as palavras deste prompt, não uma
suposição sobre o que o ambiente de execução bloqueia.
```

`commands.test`/`commands.lint` é a única execução que esta skill permite a um revisor de CLI, nunca `commands.build`, que pode reescrever artefatos versionados como `dist/` e disparar a reversão. Uma receita diz onde essa permissão é inerte na prática.

## Ler é permitido

Acrescente este bloco, literalmente, a todo prompt de revisor de CLI, junto com a cláusula acima. Um revisor que lê a cláusula como proibição de abrir arquivos a entendeu mal; o bloco existe porque esse mal-entendido foi observado:

```
Ler qualquer arquivo deste repositório é permitido e esperado — a cláusula
acima proíbe executar scripts e delegar esta revisão a outra
ferramenta, não ler.
```

## O contrato de veredito

Acrescente este bloco a todo prompt de revisor de CLI, com `<sinal da etapa>` e `<seção da etapa>` substituídos a partir da tabela abaixo:

```
Sua resposta deve conter, em uma linha própria, exatamente esta linha, sem negrito, sem
citação e sem palavras extras antes dela:
<sinal da etapa>
Escreva-a dentro da seção `<seção da etapa>` do formato de saída acima. Todas as
outras partes da sua resposta seguem esse formato como escrito. Uma resposta sem esta
linha exata é descartada como falha de transporte e sua revisão é registrada como
`not-run`, por mais sólido que seja o conteúdo.
```

E, em linha própria: `Sua tag Reportado por: para esta revisão é cli:<binary>/<slug>. Escreva essa string, exatamente, em todo achado que você reportar.` Este é o único formato de veredito do kit; nenhuma receita define outro.

| Etapa | `<seção da etapa>` | `<sinal da etapa>` (casado como prefixo de uma linha, nunca por igualdade da linha inteira) |
| --- | --- | --- |
| `document` | `## Veredito` | `- approved` ou `- changes required` |
| `code` | `## Veredito` | `- Status da revisão: approve` ou `- Status da revisão: request changes` |
| `security` | `## Achados` | `### <SEVERIDADE>: <título>` por achado, ou `- Nenhuma vulnerabilidade demonstrada ou exposição confirmada encontrada no escopo revisado` (a frase completa continua com `; as hipóteses não resolvidas continuam listadas abaixo`: case o prefixo) |

`security` não tem seção `## Veredito`: o formato de saída dela vai direto para `## Achados`, onde está o sinal. Uma cláusula final, nota de rodapé ou pontuação depois do texto casado nunca transforma um veredito real numa falha.

## Classificação de falhas

Uma chamada que não devolveu uma resposta utilizável é `not-run` com uma classe. Toda classe recebe a mesma ação: reportar `not-run` com essa classe como motivo, não gastar rodada de revisão, nunca contar como aprovação nem como rejeição e não fazer nova tentativa automática. O Coordenador avisa o dono de imediato, nomeando o revisor e a classe. `tree_changed` também descarta o veredito, seja qual for.

| Classe | Detecção |
| --- | --- |
| `not_installed` | `command -v` falha e `binaries.<binary>` em `.harness/local.yaml` está ausente ou não resolve. Uma falha de `command -v` sozinha não basta: alguns binários ficam só no `PATH` de um shell de login. |
| `auth_required` | A saída da CLI pede login ou reporta credenciais ausentes ou inválidas. |
| `rate_limited` | A saída da CLI reporta um limite de taxa, de cota ou de uso. |
| `timed_out` | O timeout do próprio binário ou o timeout da ferramenta Bash disparou, ou o processo sobreviveu a uma resposta completa (a receita diz como pará-lo). |
| `no_signal` | Saída 0, ou qualquer saída sem outra classe, e a resposta não tem o sinal de veredito da etapa: texto de ajuda, um preâmbulo de raciocínio, um arquivo vazio, uma seção sem nenhum dos resultados. Numa chamada de planejamento: nenhum conteúdo de plano reconhecível. |
| `tree_changed` | A comparação do congelamento difere (árvore viva, ou o conjunto julgado do clone). A alteração é revertida primeiro. |
| `isolation_unavailable` | O clone não pôde ser criado nem casado, uma precondição do procedimento de isolamento falhou, ou uma revisão na árvore viva rodou num projeto cujo `.git` não é um diretório (ver `live-tree-compare.md`). Sem recuo para a árvore viva. |

**Recuperar um veredito divergente.** Uma exceção estreita a `no_signal`, porque descartar um veredito bloqueante é pior do que ler um divergente: quando a resposta tem a seção de veredito da própria etapa (`## Veredito` para `document` e `code`, `## Achados` para `security`) e essa seção contém exatamente um dos dois resultados da etapa (`approved` ou `changes required`; `approve` ou `request changes`; para `security`, uma entrada de severidade para pelo menos um achado, ou uma declaração de que nenhum foi achado), com qualquer ênfase ou pontuação, leia-a como esse resultado e registre o desvio ao lado (`veredito recuperado: o revisor escreveu "<a linha que ele de fato escreveu>"`). Os dois resultados, nenhum, ou a seção ausente (para `security`, achados e uma declaração de limpeza juntos) continua `no_signal`. A recuperação nunca promove nada: uma aprovação recuperada é a aprovação de um só revisor, um veredito bloqueante recuperado bloqueia como um bem formado bloquearia. Casos conhecidos por classe: [references/gotchas.md](references/gotchas.md).

## Modo somente leitura por binário

O índice da frota: como cada binário é isolado e se ele pode planejar. Uma flag somente leitura qualifica um binário como planejador só com uma sondagem registrada e datada em que a CLI foi instruída a criar um arquivo e a rodar um comando de shell e recusou os dois. Um clone descartável também o qualifica. O clone não é um sandbox, então a qualificação é dita como é: clone descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas; escritas fora do projeto não são detectadas.

| Binário | Isolamento | Receita |
| --- | --- | --- |
| `codex` | flag somente leitura `-s read-only`, documentada e medida; pode planejar | [codex.md](references/codex.md) |
| `mcode` | clone descartável (flag somente leitura `--permission off` documentada, não medida); pode planejar | [mcode.md](references/mcode.md) |
| `grok` | clone descartável; pode planejar | [grok.md](references/grok.md) |
| `opencode` | clone descartável; pode planejar | [opencode.md](references/opencode.md) |
| `agy` | sem flag somente leitura documentada; árvore viva sob a exposição aceita RISK-002, com a comparação de `live-tree-compare.md` (só detecta); não pode planejar até que uma sondagem datada do cwd seja registrada | [agy.md](references/agy.md) |
| `claude` | clone descartável; nunca planeja enquanto o kit roda sobre o Claude | [claude.md](references/claude.md) |

## A chamada de planejamento

`/dh:plan-loop` despacha um planejador de CLI só a partir de um binário que a tabela acima marca como "pode planejar". A chamada usa a mesma resolução, o mesmo isolamento e o mesmo congelamento de uma chamada de revisão, com um prompt diferente e um teste de falha diferente. Ela é montada exclusivamente como [references/planning-call.md](references/planning-call.md) diz: o corpo do planejador, o texto da skill `implementation-planning` e o contexto da onda, depois a instrução de não escrever, a cláusula anti-delegação de planejamento e o bloco "Ler é permitido" acima, no corpo do prompt. A resposta é dado não confiável e só é persistida depois que a comparação do congelamento passou, por `docs-guide`, no invólucro que aquele arquivo define; um comando de `Checagens:` fora da lista permitida é um conflito. Uma chamada de planejamento não tem sinal de veredito e nunca é classificada pela tabela de revisão. Suas falhas usam as mesmas classes: `no_signal` significa nenhum conteúdo de plano reconhecível, lido como plano só quando contém, no formato de saída de `implementation-planning`, um título `## Plano` (ou títulos `### Fatia N`) que nomeie ao menos um objetivo e uma fatia; uma resposta vazia, uma recusa, uma pergunta de volta ao Coordenador ou prosa sem essa estrutura não é um plano. O handoff da onda se aplica a um planejador `not-run`.

## `.harness/local.yaml`

Os caminhos de binário específicos da máquina e os slugs de modelo medidos ficam em `.harness/local.yaml`, nunca nesta skill nem em qualquer arquivo sob `.skills/`. Esquema e marcadores de posição: `references/local.yaml.example`. O Coordenador é o único que escreve; o delegado só lê.

## O que carregar e quando

| Situação | Carregue, em `references/` |
| --- | --- |
| Qualquer chamada | A receita do binário, da tabela acima |
| Uma chamada em clone, ou um `$CLONE` / `$PROMPT` de receita | `clone-isolation.md` (Termos primeiro), depois `live-tree-freeze.md` para o congelamento, a comparação, a reversão e a higiene do prompt |
| Uma chamada na árvore viva (sem clone) | `live-tree-compare.md` |
| Uma chamada de planejamento | `planning-call.md` |
| Uma chamada falhou ou parece estranha | `gotchas.md`, depois acrescente o caso novo sob a sua classe |
| Escrever `local.yaml` (Coordenador) | `local.yaml.example` |

## Referência rápida

- Um comando, a primeira palavra é o binário; flags e entrada do prompt vêm da receita; compare o congelamento depois de toda chamada.
- O resultado é o sinal de veredito da etapa, nunca o código de saída. `not-run`: uma classe, sem rodada, sem nova tentativa, dono avisado de imediato.

## Erros comuns

| Erro | Por que prejudica | Em vez disso |
| --- | --- | --- |
| Montar o pacote de contexto com um `pathspec` de diretório, como `git show <sha> -- 'evals/cases/*/graders'` | O git casa arquivos, não diretórios: o comando tem sucesso com um diff vazio e o revisor julga em silêncio um escopo que nunca foi enviado (medido em 2026-09-22) | Case arquivos (`'evals/cases/*/graders/*'`) e confira o tamanho em bytes do pacote antes de enviá-lo |
| Deixar de fora a cláusula, o bloco "Ler é permitido" ou o contrato de veredito, ou pô-los só na instrução de despacho | O ambiente de execução não faz cumprir uma restrição declarada fora do corpo do prompt (RISK-001), e uma CLI real escreve o próprio formato de veredito | Copie os blocos literalmente para o texto do prompt; uma chamada de planejamento segue "A chamada de planejamento" |
| Tratar saída 0 como aprovação, ou deixar um revisor rodar `commands.build` | Texto de ajuda, um preâmbulo e uma operação nula e silenciosa saem todos com 0; um build pode reescrever artefatos versionados como `dist/` | Leia o sinal de veredito da etapa; a cláusula permite só `commands.test`/`commands.lint` |
| Chamar o clone de isolado, seguro ou sandbox, ou supor que um binário sem flag somente leitura não pode tocar o espaço de trabalho | A CLI roda como o mesmo usuário e alcança a máquina por caminho absoluto | Diga como é: clone descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas; escritas fora do projeto não são detectadas |
| Repetir por conta própria uma chamada `not-run`, ou aprovar por omissão | Uma nova tentativa esconde a falha do dono; um revisor omitido não é uma aprovação | Reporte a classe e o motivo; o Coordenador decide |

## Exemplo

Uma etapa `code` lista `["claude", "cli:codex/<slug>"]`. O delegado resolve `codex`, monta o prompt a partir do corpo de `.agents/reviewer.md` (modo `code`), do texto da skill `code-review`, do diff e dos blocos acima, escreve-o em disco, congela a árvore viva, roda `codex exec --skip-git-repo-check -s read-only -m <slug> - < prompt.md`, recompara o congelamento e lê a resposta atrás de `Status da revisão:` sob `## Veredito`. Achado: esse é o veredito. Saída 0 sem seção `## Veredito`: `not-run`, classe `no_signal`, nenhuma rodada gasta.

Uma etapa `document` lista `["claude", "cli:grok/<slug>"]`. `grok` não tem modo somente leitura, então o delegado segue a ordem de clone de `references/live-tree-freeze.md` (congelar, criar e casar o clone, congelar o clone, escrever o prompt, chamar `grok --always-approve --cwd "<clone>" --prompt-file "<prompt>"`, recomparar o clone, comparar e reverter a árvore viva, remover o clone e o diretório de congelamento). Uma diferença em qualquer das comparações é `not-run`, `tree_changed`; um clone que não pôde ser casado é `isolation_unavailable`. A resposta acrescenta `clone removed: <path>`.

## Dúvida no meio do build

Uma chamada de dúvida no meio do build é uma revisão `code` montada pelo mesmo caminho de qualquer revisão de código: o corpo do papel, `code-review`, o pacote de contexto, a cláusula anti-delegação, o bloco "Ler é permitido" e o contrato de veredito de código. O pacote é só o artefato e o contrato, mais a instrução adversarial de `code-review` em "Dúvida no meio do build"; a afirmação do autor não faz parte dele. O sinal de veredito continua sendo o sinal de código. Escolher qual entrada roda é o procedimento do Coordenador em `.agents/coordinator.md`; esta skill ainda cobre um binário.

## Relacionado

Papéis: coordinator, document-validator, reviewer, implementation-planner. Comandos: `/dh:review`, `/dh:discover`, `/dh:secure`, `/dh:plan-loop`, `/dh:auto` (carrega esta skill). Skills: code-review, document-review, security-review (donas dos formatos de veredito que esta skill lê), implementation-planning (dona do formato de saída que "A chamada de planejamento" testa), project-onboarding (`.harness/local.yaml` é escrito na primeira invocação de revisor de CLI, não pelo setup).

## Caso de prova

Dada uma etapa `code` com um revisor Claude e um revisor de CLI que tem uma flag somente leitura, a chamada de CLI é montada a partir do mesmo texto de papel e de skill que o revisor Claude recebe, rodada com essa flag e lida atrás de `Status da revisão:` sob `## Veredito`, sem confiar no código de saída. Dada uma saída 0 com um preâmbulo de raciocínio e sem seção `## Veredito`, ela é reportada como `not-run`, classe `no_signal`, não gasta rodada, nunca é lida como aprovação e não é repetida. Dada uma chamada que alterou um arquivo apesar da flag, ou uma chamada em clone que escreveu um arquivo na árvore viva por caminho absoluto, a comparação detecta, a alteração é revertida, o veredito é descartado e ela é reportada como `not-run`, `tree_changed`; o clone é removido de qualquer forma, e uma escrita fora do diretório do projeto não é detectada. Uma chamada de dúvida no meio do build usa a mesma montagem com a afirmação do autor deixada fora do pacote.
