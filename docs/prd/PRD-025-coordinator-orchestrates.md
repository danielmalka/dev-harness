# PRD-025 · O Coordenador orquestra, não executa

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-16 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-16 / 2026-10-16 |
| Tickets | T-2501, T-2502, T-2503 (planejados) |
| Origem | Conversa com o dono, 2026-10-16 (card do roadmap `Coordinator orchestrates, does not execute`, PRD-024) |
| Depende de | nada |
| Habilita | nada |

## 1. Problema

O Coordenador faz trabalho braçal que um subagente pode fazer. Isso gasta contexto da sessão principal e mistura dois ofícios: orquestração e execução. Fato: observação do dono, 2026-10-16.

Inventário da autoexecução em `.agents/coordinator.md` (os números de linha são os da base; o autor do ticket reencontra cada linha pelo texto):

- Linha 44: o Coordenador confere por amostragem uma afirmação sustentadora de um relatório não revisado.
- Linha 48: o Coordenador descobre fatos do repositório `by reading`.
- Linha 50: a fase somente leitura explora o repositório.
- Linha 56: numa edição trivial, o Coordenador `does the scoped work` e as checagens.
- Linha 60: o Coordenador congela o estado relevante.
- Linha 65: a checagem do roadmap no fim do lote não nomeia quem a executa.
- Linha 69: o Coordenador muda o status do ticket e preenche `Concluído em` / `Done on`.
- Linha 70: o Coordenador escreve o "Resumo executado" do PRD.
- Linha 87: o Coordenador executa a mecânica de chamada de um revisor de CLI externo (resolver, montar, clonar, congelar, reverter, classificar) por meio de `external-clis`.
- Linha 106: o Coordenador escreve `<document>.review.md`.
- Linha 108: o Coordenador escreve `local.yaml` e roda `git check-ignore`.
- Linha 143: o Coordenador escreve os arquivos `wave-N/candidate-*.md` do `/dh:plan-loop`.
- Linhas 145-150: o Coordenador executa a consolidação.

Outros textos do kit com o mesmo padrão: `.skills/external-clis/references/live-tree-freeze.md` (em torno da linha 340, `Coordinator reruns the check`), `.skills/context-handoff/SKILL.md` (em torno da linha 64), `.commands/setup.md` (inicialização do registro), `.commands/release.md`. O autor do ticket confere cada um e muda só os que fazem o Coordenador rodar um comando, um grep ou uma edição de várias linhas.

O mod de runtime não ajuda. Fato (leitura de `adapters/claude-code/plugin/hooks/runtime-guard.ts`): `guardFile` só age quando o evento tem `agentId` (a linha `const hit = e.agentId ? protectedHit(...) : undefined`), e `denySpawn` só age quando o evento tem `parentAgentId`. Ou seja, o mod trata um evento sem `agentId` como a sessão principal. Fato: `runtime-guard.test.ts` afirma essa convenção nos casos `without agentId passes`. Hipótese: o motor real deixa `agentId` vazio para a sessão principal do mesmo jeito; os testes de motor existentes (`panel.engine.test.ts`) são o lugar para confirmar.

Local do código-fonte: o código-fonte do mod fica em `adapters/claude-code/plugin/hooks/*.ts`. A cópia em `dist/claude-code/dev-harness/hooks/` é gerada por `go run ./cmd/dh build`. O `AGENTS.md` escreve `hooks/*.ts` como forma curta do primeiro caminho.

## 2. Solução

Após a entrega, o Coordenador decide e julga. Os subagentes executam e verificam. O Coordenador escreve por conta própria só três tipos de arquivo: seus próprios registros, mudanças de status de ticket de uma linha e pequenas correções.

Divisão do trabalho:

| Trabalho | Quem faz |
|---|---|
| Ler o repositório, grep, mapear, descobrir um fato | `repo-scout` |
| Rodar um comando ou uma checagem (build, `dh validate`, testes, `git check-ignore`) | `qa-verifier` |
| Conferir por amostragem uma afirmação sustentadora | Um despacho a `repo-scout` ou `qa-verifier` |
| Chamada de revisão por CLI (resolver, montar, clonar, congelar, reverter, classificar) | Um `qa-verifier` despachado, ou `builder` na `lane` `infra`. Ele devolve a linha de veredito e o caminho da saída bruta. O Coordenador funde os vereditos. |
| Escrever o "Resumo executado" do PRD e persistir `<document>.review.md` | `docs-guide`. O Coordenador julga o conteúdo. |
| Escrever `wave-N/candidate-*.md` no `/dh:plan-loop` | `docs-guide`. O Coordenador julga o resultado. |
| Consolidação (anexar a `EPOCHAL.md`, podar `MEMORY.md`) e `local.yaml` | O Coordenador (sem mudança). Só a checagem `git check-ignore` sobre `local.yaml` vai para o `qa-verifier`. |
| Editar mais do que poucas linhas, ou qualquer arquivo novo | Um especialista despachado para a `lane` |
| Escrever `MEMORY.md`, `EPOCHAL.md`, `RISKS.md` | O Coordenador (sem mudança) |
| Mudar o status de um ticket e preencher `Concluído em` / `Done on` (uma linha) | O Coordenador |
| Pequena correção ("poucas linhas", definida na R3) | O Coordenador, direto |

"Julgar" significa: ler o relatório do subagente e a evidência que ele cita por caminho relativo. O Coordenador pode abrir um arquivo citado. O Coordenador não roda comando e não faz grep para julgar.

A regra de congelamento não muda: nada escreve na árvore nem no harness dir enquanto uma chamada de CLI está em andamento (RISK-002). Agora ela vincula também o delegado.

O mod ganha uma regra que só avisa. Quando a sessão principal chama `Edit` ou `Write` sobre um arquivo fora da lista de permissão (R5), o mod emite um aviso. Ele nunca nega. O `Bash` fica fora da regra (ADR-005). Os especialistas não são afetados.

O `AGENTS.md` ganha uma linha que declara essa regra de escopo para o repositório. Ela prevalece, neste repositório, sobre a regra global "se a tarefa for trivial, faça direto".

Release 0.31.0, uma PR.

## 3. Regras

Marcação: R1-R10 e R18 são requisitos da decisão do dono de 2026-10-16. Os números exatos na R3 e na R5 são proposta do autor. Cada regra abaixo diz como verificá-la.

- R1 O Coordenador não lê o repositório, não faz grep, não faz build, não valida nem roda checagens por conta própria. Ele despacha `repo-scout` para leituras e `qa-verifier` para comandos e checagens. Verificável: `grep -n -i -E 'by reading|spot-check|do the scoped work|freeze the relevant|reruns the check' .agents/coordinator.md .skills/external-clis/references/live-tree-freeze.md` não encontra frase que mande o Coordenador fazer isso por conta própria; cada ocorrência antiga nomeia o papel que a executa.
- R2 O Coordenador mantém seus próprios registros. Ele escreve `MEMORY.md`, `EPOCHAL.md` e `RISKS.md` por conta própria, e a mudança de status do ticket com `Concluído em` / `Done on` (uma linha). Verificável: `.agents/coordinator.md` ainda nomeia o Coordenador como único escritor dos três registros, e a mudança de status com a data ainda está no texto do Coordenador.
- R3 Pequena correção (proposta): o Coordenador pode editar direto, em qualquer arquivo, só quando as três condições valem: um arquivo existente; no máximo 5 linhas alteradas (o maior entre a contagem de linhas antigas e a de linhas novas); nenhum arquivo novo. Qualquer outra coisa vai para um subagente. Verificável: `.agents/coordinator.md` declara essas três condições num só lugar, com o número 5, e a linha `trivial edit` (linha 56) remete à R3 e não diz mais `do the scoped work`.
- R4 As chamadas de revisão por CLI rodam num `qa-verifier` ou `builder` (`lane` `infra`) despachado. O delegado devolve a linha de veredito e o caminho relativo da saída bruta. O Coordenador funde os vereditos e não roda `command -v`, nem clona, nem congela, nem reverte. Verificável: a linha 87 e os passos de entrada de `external-clis` atribuem a mecânica ao delegado; `grep -n -i 'coordinator' .skills/external-clis/SKILL.md` não mostra passo em que o Coordenador execute uma mecânica.
- R5 Regra de aviso do mod. A regra se aplica quando um evento `Edit` ou `Write` não tem `agentId` (sessão principal) e o arquivo-alvo está fora da lista de permissão. Gatilho para `Edit`: mais de 5 linhas alteradas. Gatilho para `Write`: conteúdo de mais de 5 linhas, ou criação de um arquivo novo. O hook confere em disco se o arquivo existe; se não puder saber, ele tem falha aberta (sem aviso). O mod então emite um aviso e repassa o evento. Lista de permissão: uma checagem dedicada, não `protectedHit` (que casa só os três arquivos de registro): qualquer caminho dentro do harness dir (`.harness/` no modo repo, `<home>/projects/<name>/` no modo global). Nenhum outro caminho é excluído por nome: um `Edit` pequeno de 5 linhas ou menos em qualquer arquivo não gera aviso. Verificável: testes novos em `runtime-guard.test.ts` (ou num arquivo de teste novo ao lado dele) cobrem esses casos, e cada um passa.
- R6 A regra de aviso nunca nega e nunca muda o evento. Ela tem falha aberta diante de entrada malformada, como `guardFile`. Verificável: testes mostram que o resultado nunca é `{ deny }` para nenhuma entrada, que entrada malformada passa (como já fazem os casos R5 existentes), e que os testes de negação existentes para `MEMORY.md`, `EPOCHAL.md` e `RISKS.md` continuam passando sem mudança.
- R7 A regra de aviso não age sobre especialistas e não age sobre os fluxos de setup e de registro. Verificável: testes mostram nenhum aviso quando `agentId` está presente; nenhum aviso para `.harness/project.yaml`, `.harness/local.yaml`, `.harness/tasks/<id>/TASK.md` e os três registros na sessão principal; e a ferramenta `Bash` não é registrada pela regra.
- R8 O "Resumo executado" do PRD e a persistência de `<document>.review.md` são escritos pelo `docs-guide`. Verificável: `.agents/coordinator.md` (linhas 70 e 106) e as instruções do `docs-guide` nomeiam o `docs-guide` como escritor; o texto do Coordenador diz que ele julga o conteúdo.
- R9 "Julgar" é definido uma só vez, com as palavras da seção 2, em `.agents/coordinator.md`. A regra `spot-check a load-bearing claim` vira um despacho. Verificável: `grep -n -i 'judge' .agents/coordinator.md` encontra a definição; a linha 44 diz despachar.
- R10 O `AGENTS.md` ganha exatamente uma linha nova em "Regras vigentes". A linha declara a regra de escopo do Coordenador neste repositório: o que ele escreve, o limite de 5 linhas e o mod que só avisa. Ela diz também que prevalece sobre `if trivial, do it inline`. Verificável: `git diff main...HEAD -- AGENTS.md` mostra uma linha acrescentada e nenhuma removida.
- R11 A regra de congelamento não muda e se aplica ao delegado. Verificável: a frase da linha 107 permanece, e o texto do delegado de CLI diz que nada escreve enquanto uma chamada de CLI está em andamento.
- R12 As contagens continuam em 13 agentes e 19 + 1 comandos. Nenhum agente, skill ou comando novo. Verificável: `go run ./cmd/dh validate .` passa, e `git diff --name-status main...HEAD -- .agents .skills .commands` não lista arquivo acrescentado nem removido.
- R13 Release 0.31.0 em três lugares: `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte (`v0.31.0`). Verificável: `grep -n '0.31.0' internal/build/build.go .claude-plugin/marketplace.json` mostra os três valores.
- R14 `dist/` não é editado à mão. Verificável: `git diff --name-only main...HEAD -- dist` lista só arquivos que `go run ./cmd/dh build` regenera (`plugin.json`, `harness-manifest.json`, `bin/**`, e a cópia dos arquivos do kit alterados e de `hooks/`).
- R15 `go run ./cmd/dh validate .` passa com 0 erros, `go test ./...` fica verde, os testes do mod (`runtime-guard.test.ts` e os outros em `hooks/`) passam com o comando que o `CONTRIBUTING.md` ou o job de CI existente usa, e um segundo build não muda nada em `dist/` exceto a linha `generated_at`. Verificável: rode os comandos e compare com `git diff`.
- R16 A prova de clone limpo se aplica, porque a PR toca `internal/build/build.go` (um arquivo sob `internal/`) e `dist/*/bin`. A mudança do mod sozinha não a acionaria: `hooks/` é uma cópia, não um binário. Verificável: `git clone . <tmp>`, `go run ./cmd/dh build` lá, o `sha256sum` dos 5 binários é igual ao da árvore: 5 de 5.
- R17 Nenhum caminho de máquina privada é escrito neste PRD, num ticket ou numa fixture de teste. Verificável: `go run ./cmd/dh validate .` (`checkPrivatePaths`) passa.
- R18 Toda chamada de CLI externa é um único comando de shell que começa com o nome do binário (`agy`, `codex`, `grok`, `opencode`, `mcode`, `claude`). Nenhum `timeout` nem outro invólucro vem antes dele. Nenhum `&&`, `;` ou pipe o encadeia com o congelamento ou com qualquer outro passo. A flag do próprio binário limita o tempo (`agy --print-timeout`). Para um binário sem essa flag, a receita declara o que limita o tempo (o timeout da ferramenta Bash, ou `run_in_background`). O instantâneo anterior à chamada e a comparação posterior são comandos separados. Uma chamada em clone usa a flag de diretório de trabalho do próprio binário (`--cwd`, `--dir`, `-C`), escrita dentro do único comando que começa com o nome do binário. O T-2501 confere, em cada receita em modo clone (`grok`, `opencode`, `mcode`, `claude`), se existe tal flag na ajuda ou nos docs do próprio binário. Para um binário sem essa flag, a receita de clone mantém `cd "$CLONE" && <binary> ...` como única exceção permitida. Essa receita declara que essa forma não casa com uma regra de permissão `Bash(<binary>:*)`, então a chamada pode precisar da aprovação do dono ou ficar não rodada. Fato (medido em 2026-10-16 por um `qa-verifier` despachado): `cd /tmp` e depois `pwd` numa chamada Bash separada imprimiu a raiz do repositório, então o diretório de trabalho de um subagente não persiste entre chamadas Bash. O status de modo clone do `claude` não muda. Motivo (fato, observado em 2026-10-16): a regra de permissão do dono `Bash(agy:*)` casa só comandos que começam com `agy`. Um comando de subagente `timeout 960 agy ...`, encadeado com um script de instantâneo, não casou e o classificador do modo automático o negou. A mesma chamada como comando `agy ...` isolado rodou. Verificável: cada receita em `.skills/external-clis/references/` mostra a chamada começando com o nome do binário, e `grep -n -E '(timeout|&&|;).*(agy|codex|grok|opencode|mcode|claude) ' .skills/external-clis/references/*.md` não encontra receita que encadeie ou embrulhe a chamada.

## 4. Docs

Atualizar:

- `.agents/coordinator.md` (R1-R4, R8, R9, R11, R18).
- `.skills/external-clis/SKILL.md`, `.skills/external-clis/references/live-tree-freeze.md`, e as receitas `references/agy.md`, `codex.md`, `grok.md`, `opencode.md`, `mcode.md`, `claude.md` onde mostram uma chamada embrulhada ou encadeada (R4, R11, R18). A receita `agy.md` ganha também um fato: em 2026-10-16 uma chamada `agy --print` escreveu a resposta completa em cerca de 2 minutos, mas o processo ficou vivo além de `--print-timeout 15m` e precisou de SIGKILL. A receita manda o executor tratar uma resposta completa em disco com a linha de veredito como o resultado, e então parar o processo. A seção 3 de `references/clone-isolation.md` também muda: ela troca `cd "$CLONE" &&` pela flag de cwd do próprio binário onde existir (R18).
- `.agents/docs-guide.md`: ganha a escrita do Resumo executado do PRD, de `<document>.review.md` e dos arquivos candidatos do `/dh:plan-loop` (R8).
- `.skills/document-review/SKILL.md`: seus passos de persistência (o passo 9 e o protocolo de rodadas) nomeiam o Coordenador como escritor de `<document>.review.md`; passam a nomear o `docs-guide` (R8).
- `.skills/context-handoff/SKILL.md`, `.commands/setup.md`, `.commands/release.md`, `.commands/plan-loop.md` (só onde o texto faz o Coordenador rodar um comando ou uma edição longa; o autor do ticket confere cada um com `grep`).
- `AGENTS.md` (R10).
- `docs/tutorial.html` e `docs/en/tutorial.html` (a descrição das atribuições do Coordenador; o autor do ticket encontra as passagens com `grep -n -i 'coordenador\|coordinator'`).
- `README.md` e `README.en.md` só se o `grep` encontrar uma frase que diga que o Coordenador lê, roda ou edita por conta própria. Não existe linha de versão para mudar (o PRD-024 a removeu).
- `docs/roadmap.html` e `docs/en/roadmap.html`: mover o card `Coordinator orchestrates, does not execute` / `O Coordenador orquestra, não executa` para fora de Planned; acrescentar um card Delivered para a 0.31.0 nos dois arquivos; atualizar o chip de status e a linha `releases 0.1.0 to 0.30.1` para 0.31.0. O `dh validate` confere a deriva contra o `CHANGELOG.md`.
- `CHANGELOG.md` (entrada 0.31.0).
- `.claude-plugin/marketplace.json`, `internal/build/build.go` (R13).
- `adapters/claude-code/plugin/hooks/runtime-guard.ts` e seu teste (R5-R7).

Não afetados (o autor do ticket confere com `grep`):

- `docs/plano-produto.html`: histórico interno em pt-br; muda só se declarar as atribuições do Coordenador como comportamento atual.
- `docs/adr/`: o ADR-005 fica como histórico; a regra que só avisa o estende e não precisa de ADR novo, a menos que o dono peça.
- `docs/inicio-rapido.html`, `docs/en/quick-start.html`.

## 5. Fora de escopo

- Qualquer regra de negação para a sessão principal. O mod só avisa.
- Uma trava sobre o `Bash` (o ADR-005 mantém o `Bash` de fora).
- Um agente, skill ou comando novo, e qualquer mudança na contagem de agentes.
- Uma mudança em como os especialistas escrevem, e qualquer mudança na posse dos três registros.
- Mudar a regra global do usuário fora deste repositório.
- Avaliação paga. Nenhum caso de avaliação é acrescentado (`AGENTS.md`).

## 6. Decisões do dono

- Dono, 2026-10-16, Q1: o Coordenador não executa trabalho por conta própria. Ele pode fazer pequenas correções de poucas linhas em qualquer arquivo. Ele mantém seus registros. Leituras de repositório, grep, builds, validate e checagens vão para subagentes. O limite exato é proposta do autor (R3) e precisa de aprovação do dono.
- Dono, 2026-10-16, Q2: as chamadas de revisão por CLI rodam num `qa-verifier` ou `builder` (`lane` `infra`) despachado, que devolve a linha de veredito e o caminho da saída bruta. O Coordenador funde os vereditos. A regra de congelamento se aplica ao delegado.
- Dono, 2026-10-16, Q3: redação em `.agents/coordinator.md` mais uma regra de mod que só avisa sobre `Edit` e `Write` fora de uma lista de permissão. Sem negação. O `Bash` fica de fora. Especialistas e os fluxos de setup e de registros ficam sem afetação. Testes são obrigatórios.
- Dono, 2026-10-16, Q4: o `docs-guide` escreve o Resumo executado do PRD e persiste `<document>.review.md`. O Coordenador julga o conteúdo. A mudança de status do ticket e o `Done on` ficam com o Coordenador.
- Dono, 2026-10-16, Q5: "julgar" significa ler o relatório e a evidência citada. O Coordenador pode abrir arquivos citados e não roda comandos nem faz grep. A conferência por amostragem vira um despacho.
- Decisões pendentes: os itens do Apêndice.

## Apêndice

### Tickets sugeridos

| Ticket | `lane` | Executor | Escopo |
|---|---|---|---|
| T-2501 | `backend` | `harness-maintainer` com `harness-authoring` | Texto do kit: R1-R4, R8-R12, R18 (receitas de CLI); docs da seção 4 |
| T-2502 | `backend` | `builder` (TypeScript sob `adapters/`) | Primeiro passo: confirmar o canal do aviso (ver Apêndice). Depois a regra de aviso em `runtime-guard.ts`: R5-R7 |
| T-2503 | `teste (unitário e integração)` | `qa-verifier` | Testes para R5-R7 em `runtime-guard.test.ts`; execução de R15-R16 |

T-2502 e T-2503 escrevem o mesmo arquivo de teste. Ordem: T-2502 primeiro, depois T-2503. Os arquivos de release (R13) vão com o T-2501.

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Negar `Edit` e `Write` para a sessão principal | O dono escolheu só avisar. Uma negação quebra o setup e as pequenas correções. |
| Travar o `Bash` também | O ADR-005 mantém o `Bash` de fora; um comando de shell não pode ser classificado de forma confiável. |
| Só redação, sem mod | O dono pediu também um sinal em runtime. |
| Contar linhas com um diff exato no mod | Caro demais. A contagem usa o maior entre a contagem de linhas antigas e a de novas do evento `Edit`. |

### Riscos e decisões pendentes

- Risco: o motor real define `agentId` para a sessão principal, então a regra de aviso nunca dispara ou sempre dispara · mitigação: o T-2503 confirma com um teste de motor (como `panel.engine.test.ts`) antes de a regra ser dada como pronta.
- Risco: o aviso inunda a sessão principal durante um setup legítimo · mitigação: a lista de permissão (R5) e os testes da R7.
- Risco: mais despachos custam mais tokens e tempo · mitigação: o limite de 5 linhas mantém as correções triviais diretas; o dono aceitou a troca.
- Decidido: o limite de 5 linhas alteradas (R3, R5) · aprovado pelo dono, 2026-10-16.
- Item pendente (só a ligação ao motor): o canal do aviso. Fato: nenhum mod em `adapters/claude-code/plugin/hooks/` chama `$.ui`; os handlers de `tool.call` ali devolvem `{ deny }` ou o resultado de `next(e)`. Fato: cinco mods do repositório vizinho `malka-claude-mods` (`cleanscreen`, `menu`, `mode`, `permissions` e outros) chamam `$.ui.toast(text)` para um aviso que não bloqueia. Hipótese: um handler de `tool.call` pode chamar `$.ui.toast(text)` e depois `return next(e)`. Não confirmei isso numa declaração de tipos do motor. O T-2502 começa por essa leitura (os tipos do motor, ou um teste de motor no estilo de `panel.engine.test.ts`). O contrato de teste não depende disso: R5-R7 são afirmadas sobre uma função pura em `hooks/lib/` que recebe o evento, o home do harness, o diretório de trabalho e um callback de existência de arquivo, e devolve o texto do aviso ou `undefined`. A ligação (`toast`, ou outro canal que a leitura achar) é o único item em aberto. Se o motor não tiver canal, o dono escolhe · respondido por: `builder` no T-2502 · bloqueia: só a ligação do T-2502.
- Decisão pendente: se as fontes de texto do kit fora de `coordinator.md` (R1, Seção 4) precisam de mudança, conforme o `grep` do autor do ticket · respondido por: autor do T-2501, reportado ao Coordenador · bloqueia: o escopo do T-2501.

A R18 mudou depois da rodada 2 de validação (aprovada); o dono a aceitou sem nova rodada, 2026-10-16.

## Resumo executado

**Entregue**: O Coordenador agora orquestra e julga enquanto os especialistas leem, rodam comandos e escrevem documentos, com um `toast` do mod que só avisa para edições diretas além do limite de 5 linhas e chamadas de CLI externa de um comando só com um passo de limpeza.

**Regras**:

- R1: Cumprida. `.agents/coordinator.md:47` manda as leituras para o repo-scout e os comandos para o qa-verifier; o grep do PRD deixa uma ocorrência, que nomeia o papel (`.skills/external-clis/references/live-tree-freeze.md:353`).
- R2: Cumprida. `.agents/coordinator.md:38` e `:49`.
- R3: Cumprida. `.agents/coordinator.md:50`.
- R4: Cumprida. `.skills/external-clis/SKILL.md:35` e `:37`.
- R5: Cumprida. Testes em `adapters/claude-code/plugin/hooks/runtime-guard.test.ts` (235 passam, 0 falham); a correção de `inHarness` foi conferida pelo revisor opencode.
- R6: Cumprida. Mesmos testes e mesma conferência do revisor da R5.
- R7: Cumprida. Mesmos testes e mesma conferência do revisor da R5.
- R8: Cumprida. `.agents/coordinator.md:51` e `:113`, `.agents/docs-guide.md:56`, `.skills/document-review/SKILL.md`.
- R9: Cumprida. `.agents/coordinator.md:47-48`.
- R10: Cumprida. O diff de `AGENTS.md` é +1 linha, -0.
- R11: Cumprida. `.skills/external-clis/SKILL.md:37`, `.agents/coordinator.md:114`.
- R12: Cumprida. `dh validate` mostra 13 agentes e 19 comandos; nenhum arquivo do kit criado ou removido.
- R13: Cumprida: versão 0.31.0 em `internal/build/build.go` (`defaultVersion`) e `.claude-plugin/marketplace.json` (`version`, `ref` v0.31.0); o `dh build` regenerou `dist/` e `harness-manifest.json`.
- R14: Cumprida. Só os arquivos gerados diferem.
- R15: Cumprida. `dh validate` 0 erros, `go test` ok, teste do plugin 235/0, build idempotente.
- R16: Cumprida: clone limpo mais o estado exato da árvore de trabalho, `go run ./cmd/dh build`; o sha256 dos 5 binários de plataforma (e de `dh`, `dh.cmd`) é igual ao da árvore: 5 de 5.
- R17: Cumprida. `checkPrivatePaths` passa no `dh validate`.
- R18: Cumprida. As receitas começam com o nome do binário (`grok --cwd`, `opencode run --dir`, `mcode exec --cwd`); o `claude` mantém a exceção documentada do `cd`.
- Revisão de código: entrada 1 `cli:codex/gpt-6-luna` `request changes` (1 `Major`, corrigido: `warn.ts` `inHarness`); entrada 2 `cli:opencode/opencode-go/muse-spark-1.3-contributor` `approve`.

**Tickets**:

- T-2501 backend, `harness-maintainer`: aceito após 2 retrabalhos.
- T-2502 backend, `builder`: aceito após a correção da revisão (`warn.ts` `inHarness`).
- T-2503 teste (unitário e integração), `qa-verifier`: 9 testes, 235 passam.

**Docs**: `CHANGELOG.md` (entrada 0.31.0); `docs/roadmap.html` e `docs/en/roadmap.html` (card Planned movido para um card Delivered da 0.31.0, chip de status e linha de fontes); `docs/tutorial.html` e `docs/en/tutorial.html` (atribuições do Coordenador e o mod que só avisa); este Resumo executado. `README.md` e `README.en.md`: nenhuma frase dizia que o Coordenador lê, roda ou edita por conta própria, e não existe linha de versão, então sem mudança.

**Fora de escopo**: qualquer regra de negação para a sessão principal; uma trava sobre o `Bash`; um agente, skill ou comando novo e qualquer mudança na contagem de agentes; uma mudança em como os especialistas escrevem ou na posse dos três registros; mudar a regra global do usuário fora deste repositório; avaliação paga (nenhum caso de avaliação acrescentado).
