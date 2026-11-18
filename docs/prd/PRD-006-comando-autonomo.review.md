# Revisão: PRD-006 — comando autônomo `/dh:auto`

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-006-comando-autonomo.md |
| Fonte | .harness/tasks/UND-004/GRILL.md; decisões do CLAUDE.md de 2026-09-25 |
| Revisores | claude (document-validator, opus); cli:codex/gpt-5.6-luna; cli:grok/grok-4.7 |

## Rodada 1 de 6 — 2026-09-25

| Revisor | Veredito |
|---|---|
| claude | `changes required` (7 bloqueantes, 4 não bloqueantes) |
| cli:codex/gpt-5.6-luna | `changes required` (4 bloqueantes) |
| cli:grok/grok-4.7 | `changes required` (12 bloqueantes, 1 não bloqueante) |
| Mesclado | `changes required` |

Congelamento do workspace: idêntico depois das três chamadas.

- **A1** [blocking] `/dh:resume` deve ler `Stage:`, mas o `resume.md` não lê e o AC-16 proíbe editá-lo. — claude P1, codex P2, grok P2. Disposição: `.commands/resume.md` (ou `context-handoff`) lê `Stage:` de `AUTO-<n>/BRIEF.md` quando o `MEMORY.md` cita uma tarefa `AUTO-<n>`; acrescente-o às edições permitidas.
- **A2** [blocking] O AC-16 contradiz a documentação no pull request, os casos de avaliação, as cópias reembutidas do coordinator nos prompts de avaliação, o `dist/` regenerado e o `AUTO-<n>/BRIEF.md` de runtime. — claude P2, codex P4, grok P1. Disposição: reescreva como o claude sugeriu e acrescente os reembutimentos, os artefatos de runtime e o `resume.md` (A1).
- **A3** [blocking] O AC-14/AC-15 e "Não entra" proíbem a redução de revisores do RF-11. — claude P3. Disposição: "com o RF-11 aplicado; nenhuma outra redução".
- **A4** [blocking] Nenhum passo executa commit, push, PR, merge ou tag, em nenhuma ordem definida; faltam as regras permanentes de git. — claude P4, codex P3, grok P3. Disposição: novo RF+AC. Depois que o release reporta pronto, o Coordenador faz só o que (c)-(f) autorizaram, nesta ordem: commit → prova que o projeto exige (clone limpo para `dist/`) → push → PR → merge só com CI verde + revisão APPROVED + branch que não seja a main → tag no commit de merge. Nunca `--force` nem `--no-verify`; CI desligado ou quebrado não é verde. Um item não autorizado encerra a execução com o estado registrado. Os passos de PR/merge do RF-13 só valem quando autorizados.
- **A5** [blocking] Não está dito como as etapas são encadeadas; o `release` tem `disable-model-invocation`. — claude P5, grok P11 (parte). Disposição: o Coordenador segue o Roteamento/Pré-requisitos/Saída/Limites de cada `.commands/<stage>.md` em linha, na mesma sessão, inclusive o `release`; sem invocação por barra.
- **A6** [blocking] As perguntas dos modos (b)/(c) contra o RF-03 "sem mais perguntas" e a mecânica delas. — claude P6, codex P1. Disposição: elas são a única exceção além das paradas do RF-06. A pergunta é feita no chat depois que o estado é registrado (Stage + `MEMORY.md`), para o `/dh:resume` funcionar se a sessão acabar. O AC-06 é reescrito.
- **A7** [blocking] "Vulnerabilidade confirmada" não está definida frente às classes do kit. — claude P7, grok P10. Disposição: uma vulnerabilidade demonstrada ou uma exposição confirmada para a execução; uma hipótese não.
- **A8** [blocking] A detecção de um PRD aprovado ou de um plano validado usa o sinal errado. — claude P8 (`non-blocking`), grok P4, P5. Mesclado como bloqueante. Disposição: um PRD aprovado é aquele cuja linha `Status` da tabela diz `aprovado`/`approved`; um plano validado tem um `PLAN.review.md` ao lado, cujo veredito da última rodada é `approved`; qualquer outro caso propõe `discover`.
- **A9** [blocking] O gatilho seguro do RF-12 parafraseia o `secure.md` de forma mais estreita. — grok P6. Disposição: cite literalmente a superfície do próprio `.commands/secure.md` (autenticação, dados, entrada, execução); sem lista própria.
- **A10** [blocking] O RF-11 (iii) perde "plano de". — grok P7. Disposição: "o plano de um lote mecânico passa por um revisor".
- **A11** [blocking] "Achado menor" não está definido. — grok P8. Disposição: um achado que a skill de revisão da etapa classifica como não bloqueante (Minor/Cosmetic no `code-review`, `non-blocking` no `document-review`/`security-review`).
- **A12** [blocking] "Nota de avaliação não atingida" não está definida. — grok P9. Disposição: a nota de corte que o plano fixa por caso antes de qualquer execução paga; uma falha é uma nota do avaliador no JSON de resultado abaixo dessa nota de corte.
- **A13** [blocking] O gate do RF-14 afirma ser prática existente do kit. — grok P12. Disposição: declare-o como regra deste comando. Custo estimado do juiz = a razão juiz/agente da última linha de base registrada do caso × `--max-cost-usd`, ou o próprio `--max-cost-usd` quando não há linha de base.
- **A14** [non-blocking] `<base-branch>` e o alvo da PR não estão definidos. — claude P9. Disposição: uma feature branch criada a partir da `main` no início; as PRs por fatia têm a `main` como alvo.
- **A15** [non-blocking] "docs-guide" citado como comando; as etapas não invocam verify/secure. — claude P10, grok P11. Disposição: `/dh:secure` conforme o RF-12; o agente `docs-guide` no fim do lote; retire o `verify`.
- **A16** [non-blocking] O questionário tem 8 itens (7 + confirmação da etapa inicial). — claude P11. Aceito.
- **A17** [non-blocking] A frase do "porquê" do AC-17 no `MEMORY.md` é supérflua. — grok P13. Aceito: retire-a.

As disposições A1, A4–A15 do Coordenador são decididas em nome do dono, sob a autorização de sessão de 2026-09-25.

## Rodada 2 de 6 — 2026-09-25

| Revisor | Veredito |
|---|---|
| claude | `changes required` (1 bloqueante, 4 não bloqueantes); A1–A17 confirmados como aplicados |
| cli:codex/gpt-5.6-luna | `changes required` (1 bloqueante) |
| cli:grok/grok-4.7 | `changes required` (3 bloqueantes, 1 não bloqueante) |
| Mesclado | `changes required` |

Congelamento do workspace: idêntico depois das três chamadas.

- **B1** [blocking] O significado de `Stage:` não está definido. — codex P1. Disposição: `Stage:` sempre nomeia a próxima unidade (etapa ou fatia) a executar; concluir uma unidade o avança; uma parada registra a unidade interrompida, não marcada como concluída.
- **B2** [blocking] Branch por fatia exige commits por fatia, que a ordem do RF-16 proíbe; o ponto de ramificação não está definido. — claude P1. Disposição: no modo branch por fatia com (c) autorizado, cada fatia recebe commit na própria branch quando o gate de build/QA/revisão fica verde (não é violação do RF-16). A fatia N+1 sai da fatia N (empilhadas). Push → PR → merge na ordem das fatias → a tag roda depois do release. Branch por fatia sem (c) é sinalizado como incompatível no questionário, que propõe branch única.
- **B3** [blocking] O passo 2 do RF-16 diz que a prova roda "antes de commitar o dist"; o candidato é o commit local. — grok A4, claude P5 (`non-blocking`). Mesclado como bloqueante. Disposição: o commit do passo (1) é o candidato. Quando o `dist/` foi regenerado, rode a prova de sha256 em clone limpo nesse commit; se falhar, faça amend só enquanto não houver push, nunca `--force`; faça push só depois que a prova passar. Os passos de PR por fatia só valem com (d), os de merge só com (e).
- **B4** [blocking] O RF-09 e a seção 6 atribuem ao CLAUDE.md um gate de commit/push que ele não declara. — grok A18. Disposição: commit e push exigem o questionário (c) mais a prova B3 quando o `dist/` mudou. O gate de CI verde + APPROVED + fora da main vale para o merge (passo 5 do RF-16). Remova a atribuição errada.
- **B5** [blocking] A posição do `/dh:secure` na cadeia não está definida. — grok A19. Disposição: por fatia, o secure roda logo depois dos gates de build dessa fatia e antes da fatia seguinte, quando a fatia toca a superfície do `.commands/secure.md`; caso contrário, registra-se um pulo. O `review` do lote roda depois do laço de fatias e antes do `release`.
- **B6** [non-blocking] CI não verde no merge. — claude P2. Disposição: um CI pendente pode ser aguardado; CI vermelho, desligado ou quebrado encerra a execução com o estado registrado.
- **B7** [non-blocking] Detecção de status: o valor começa com `aprovado`/`approved`, sem diferenciar maiúsculas de minúsculas. — claude P3. Aceito.
- **B8** [non-blocking] O Sinal 3 aponta para a seção 6 nos casos de proporcionalidade. — claude P4. Disposição: "observável nos casos de avaliação que o plano define (seção 6)".
- **B9** [non-blocking] RF-15 "ou context-handoff" contra o AC-16. — grok A20. Disposição: só o `.commands/resume.md`.

Todas as disposições decididas pelo Coordenador em nome do dono (autorização de sessão de 2026-09-25).

## Rodada 3 de 6 — 2026-09-25

| Revisor | Veredito |
|---|---|
| claude | `changes required` (1 bloqueante, 3 não bloqueantes); A1–A17, B1–B9 confirmados |
| cli:codex/gpt-5.6-luna | `changes required` (1 bloqueante) |
| cli:grok/grok-4.7 | `changes required` (7 bloqueantes) |
| Mesclado | `changes required` |

Congelamento do workspace: idêntico depois das três chamadas.

- **D1** [blocking] O AC-01 limita a cadeia aos agentes de `.agents/`; as entradas de CLI `reviewers.*` configuradas não são agentes. — codex P1. Disposição: o `/dh:auto` não acrescenta nenhum agente, skill ou comando novo além do `auto.md` e mantém todo revisor configurado, inclusive as entradas `cli:`.
- **D2** [blocking] A lista de exceções do AC-16 omite as cópias do coordinator em `append_system_prompt` nos `evals/cases/*/prompt.md` existentes. — grok A2. Aceito: nomeie-as.
- **D3** [blocking] O AC-24 diz que as fatias saem da `main`; o RF-13/AC-41 dizem empilhadas. — grok C1. Disposição: `<base-branch>` é a feature branch criada a partir da `main`; a fatia 1 sai dela; a fatia N+1 sai da fatia N; toda PR tem a `main` como alvo.
- **D4** [blocking] A prova em clone limpo roda só no último commit de fatia. — claude C1, grok C2. Disposição: todo commit que regenera o `dist/` é um candidato próprio (o commit do lote ou cada commit de fatia). A prova roda na ponta dessa branch antes do push dela; faça amend só enquanto essa branch não tiver push; nunca `--force`.
- **D5** [blocking] A documentação de fim de lote (docs-guide antes do release) não tem commit no branch por fatia. — grok C3. Disposição: o passo (1) no branch por fatia commita, na última branch de fatia, tudo o que ainda não tiver commit depois do docs-guide e do release, e não recommita arquivos de fatia.
- **D6** [blocking] Ordem por fatia entre o secure e o commit da fatia. — grok C4. Disposição: gates de build → `/dh:secure` ou o pulo registrado → commit da fatia. Uma parada do RF-06 durante o secure acontece antes do commit.
- **D7** [blocking] O AC-20 diz "só um achado menor". — grok C5. Disposição: qualquer número de achados que a skill da etapa classifica como não bloqueantes não abre rodada de correção.
- **D8** [blocking] Os Pré-requisitos dos comandos encadeados dizem "pergunte ao dono", o que conflita com o modo (a). — grok C6. Disposição: a instrução de um comando encadeado de perguntar ao dono nunca se sobrepõe ao RF-05. No modo (a) o Coordenador decide e registra; em (b)/(c) valem o RF-05/AC-29. Nomeie os Pré-requisitos do `.commands/build.md` como cobertos.
- **D9** [non-blocking] O release devolve "not ready". — claude C2. Disposição: a execução termina com o motivo e `Stage: release` registrado (RF-07), sem o RF-16.
- **D10** [non-blocking] Mecânica de merge de PRs empilhadas. — claude C3. Disposição: as PRs por fatia são mescladas na ordem das fatias com um commit de merge, nunca squash, de modo que não precisa de rebase nem de `--force`; as PRs de branch única mantêm o squash habitual do projeto.
- **D11** [non-blocking] A seção 9 cita só a rodada 1. — claude C4. Aceito: rodadas 1–3, itens A/B/D.

Todas as disposições decididas pelo Coordenador em nome do dono (autorização de sessão de 2026-09-25).

## Rodada 4 de 6 — 2026-09-25

| Revisor | Veredito |
|---|---|
| claude | `approved` (2 não bloqueantes); itens A/B/D confirmados |
| cli:codex/gpt-5.6-luna | `approved` (0 achados) |
| cli:grok/grok-4.7 | `changes required` (4 bloqueantes) |
| Mesclado | `changes required` |

Congelamento do workspace: idêntico depois das três chamadas.

- **E1** [blocking] O AC-41 ainda commita a fatia logo depois dos gates de build, antes do secure. — grok D6. Disposição: o AC-41 commita só depois dos gates de build e do secure (ou do pulo registrado); uma parada do RF-06 durante o secure vem antes do commit.
- **E2** [blocking] O AC-48 não tem trava de modo de branch. — grok E1. Disposição: ele abre com "no modo branch por fatia com (c) autorizado". O modo branch única não tem commit por fatia; seu único commit é o passo 1 do RF-16 depois do release.
- **E3** [blocking] O AC-23/AC-24 não condicionam commit, PR e merge a (c)/(d)/(e). — grok E2. Disposição: o commit de branch única acontece só com (c) e a PR só com (d). As PRs por fatia abrem só com (d) e são mescladas só com (e), na ordem das fatias, com commits de merge. Caso contrário, a execução para nesse passo do RF-16 (AC-32).
- **E4** [blocking] Modo (a) contra "pergunte ao dono" também no Roteamento de um comando encadeado (Roteamento do `discover.md`), além dos Pré-requisitos. — grok E3. Disposição: qualquer instrução de perguntar ao dono, no Roteamento ou nos Pré-requisitos de um comando encadeado, cede ao RF-05. Nomeie o Roteamento do `discover.md` e os Pré-requisitos do `build.md`. No modo (a) o Coordenador decide e registra; em (b)/(c) a pergunta segue o RF-05/AC-29.
- **E5** [non-blocking] O alvo da tag no branch por fatia. — claude P1. Aceito: o commit de merge da PR da última fatia na `main`, ou o commit de squash da PR única.
- **E6** [non-blocking] A frase do RF-13 sobre o escopo do RF-16. — claude P2. Aceito: "que rege o commit de fechamento e push/PR/merge/tag".

## Rodada 5 de 6 — 2026-09-25

| Revisor | Veredito |
|---|---|
| claude | `approved` (2 não bloqueantes); E1–E6 confirmados |
| cli:codex/gpt-5.6-luna | `approved` (1 não bloqueante) |
| cli:grok/grok-4.7 | `changes required` (4 bloqueantes) |
| Mesclado | `changes required` |

Congelamento do workspace: idêntico depois das três chamadas.

- **G1** [blocking] Os Limites do `build.md` e o Procedimento 8 do coordinator dizem "replaneje com o dono" no limite de 6 rodadas; o RF-07 diz pare. — grok F1. Aceito: no RF-07, esse texto cede. A execução registra o estado bloqueado com `Stage:` na unidade interrompida e termina; o replanejamento acontece só numa sessão posterior.
- **G2** [blocking] O Intake do coordinator diz "aguarde a aprovação antes da primeira edição". — grok F2. Aceito: no RF-05, a espera do Intake não é um gate extra. No modo (a), um plano `approved` pelo painel de validação vai para o build sem pergunta. Nos modos (b)/(c), o Coordenador pergunta só uma decisão aberta de produto/contrato (AC-29) e não espera nos demais casos.
- **G3** [blocking] "Não entra" contra o AC-16. — grok F3. Aceito com a redação do grok.
- **G4** [blocking como ambiguidade; direção sugerida pelo grok rejeitada] O que significa "revisão APPROVED" no merge. — grok F4. Disposição do Coordenador: o sinal é o veredito `approve` da revisão independente na sessão (code-review / `.commands/review.md`, registrado no `MEMORY.md`). Um estado de revisão de PR do GitHub não é exigido. Evidência: a regra permanente de merge do dono cita o veredito APPROVED do revisor de QA, e este é um projeto de um único dono, sem revisores de PR; exigir um estado de revisão de PR tornaria impossível todo merge autônomo. Declare isso explicitamente no passo 5 do RF-16 e na seção 6, para nenhum builder ler de outro jeito.
- **G5** [non-blocking] A ordem "sempre" da seção 8 cobre os dois modos. — claude P1. Aceito: "Em branch por slice … commit; em branch única, build → secure/skip, sem commit por slice".
- **G6** [non-blocking] A seção 9 cita as rodadas 1–3. — claude P2, codex P1. Aceito: rodadas 1–5, itens A–G.

## Rodada 6 de 6 (final) — 2026-09-26

| Revisor | Veredito |
|---|---|
| claude | `approved` (1 não bloqueante: G6 incompleto) |
| cli:codex/gpt-5.6-luna | `approved` (0 achados) |
| cli:grok/grok-4.7 | 1ª chamada: timeout sem veredito (falha de transporte, `not-run`). Nova tentativa: `changes required` (2 bloqueantes), lida mesmo que a nova tentativa também tenha terminado com exit 124 depois de gravar o veredito |
| Mesclado | `changes required`. O limite está esgotado. |

Congelamento do workspace: idêntico depois de todas as chamadas.

- **H1** [blocking] A etapa de plano do `/dh:auto` usaria o Roteamento do `plan.md`, que despacha só o `document-validator` e ignora `reviewers.document`. — grok P1. Disposição do Coordenador: aceito com a redação do grok. A etapa de plano usa o procedimento `reviewers.document` do `discover.md`; o RF-11 (iii) se aplica.
- **H2** [blocking] Uma parada dentro do RF-16 não diz o que `Stage:` registra. — grok P2. Disposição do Coordenador: aceito. `Stage:` nomeia o passo interrompido do RF-16 (`commit`, `proof`, `push`, `pr`, `merge`, `tag`); o `/dh:resume` continua ali sem reexecutar o release nem os passos concluídos.
- **G6** [non-blocking] A seção 9 omite E1–E6. — claude, grok. Aceito.

Limite de rodadas atingido com H1/H2 em aberto. O Coordenador aplica H1, H2 e G6 como dispostos, sem uma sétima rodada de validação, e entrega o PRD ao dono. O relatório declara que essas três edições não foram revalidadas pelo painel.

## Verificação pós-aprovação do corte do branch por fatia — 2026-09-26

O dono aprovou o PRD-006 e cortou o branch por fatia da v1. O corte passou por uma verificação com escopo restrito e um único revisor (claude, document-validator), conforme a regra de proporcionalidade: `changes required`, 1 bloqueante, 2 não bloqueantes.
- O item da seção 8 ligava a tag a `e` em vez de `f`. Corrigido: a tag acontece "só com `f`, depois do merge autorizado por `e`".
- A sobra do AC-43 "(em qualquer estratégia de branch)" foi removida.
- A data de atualização do cabeçalho agora é 2026-09-26.

O Coordenador aplicou essas três correções de uma frase diretamente, como edição trivial.
