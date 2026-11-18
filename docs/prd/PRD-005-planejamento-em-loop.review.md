# Revisão: PRD-005 — rodada 1 (mesclada)

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

Estado avaliado: o primeiro rascunho de `docs/prd/PRD-005-planejamento-em-loop.md`. Fonte: `.harness/tasks/UND-003/GRILL.md` (G1–G17). Limite de rodadas: 6, decisão do dono de 2026-09-25.

| Revisor | Veredito |
|---|---|
| `claude` (document-validator, opus solicitado, efetivo não verificado) | `changes required`: 5 bloqueantes, 2 não bloqueantes |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | `changes required`: 4 bloqueantes, 1 não bloqueante |
| `cli:grok/grok-4.7` | `changes required`: 6 bloqueantes, 2 não bloqueantes |

## Veredito
- `changes required`

## Achados (mesclados). Cada item cita quem o reportou. A disposição segue as respostas do dono na rodada 4 do grill (G18–G23).

- R1: a convergência é descrita como "zero bloqueadores" e "sem ressalvas", mas o dono decidiu a aprovação unânime (G9). O PRD também não diz o que acontece quando um revisor fica `not-run`.
  - Reportado por: claude (P1), cli:codex/gpt-5.6-luna (P2).
  - Disposição: G9 + G22. O loop converge quando todo revisor que rodou devolve `approved` para o plano vencedor e pelo menos um deles é de outro fornecedor. Achados não bloqueantes não impedem a convergência.
- R2: o fallback sem CLI é descrito como o `/plan` inalterado com limite de 6, mas o texto do kit ainda diz 2.
  - Reportado por: cli:codex/gpt-5.6-luna (P1).
  - Disposição: o lote do limite de 6 em todo o kit passa a ser pré-requisito explícito da implementação deste PRD. O fallback roda o `/plan` pós-migração.
- R3: um empate em achados bloqueantes torna o AC-07 intestável.
  - Reportado por: cli:codex/gpt-5.6-luna (P3).
  - Disposição: G18. Vence o plano com menos achados no total. Se isso também empatar, vence o plano do Claude.
- R4: o PRD não diz que não autoriza nenhuma implementação, nem que a implementação é a terceira da fila.
  - Reportado por: cli:codex/gpt-5.6-luna (P4).
  - Disposição: acrescentar as duas coisas como restrição.
- R5: o conjunto de CLIs planejadoras é tratado como decidido (`reviewers.document`) e também aparece como pendente. O caso de 3 ou mais CLIs não está definido.
  - Reportado por: claude (P3), cli:grok/grok-4.7 (P6, P7), cli:codex/gpt-5.6-luna (P5).
  - Disposição: G19. Reaproveitar `reviewers.document` e percorrer as CLIs em rodízio, na ordem da lista. Remover a cláusula da "CLI só de revisão".
- R6: o PRD não diz o que acontece com uma onda cujo plano da CLI foi descartado.
  - Reportado por: claude (P4), cli:grok/grok-4.7 (P5).
  - Disposição: G20. A vez passa para a CLI seguinte na mesma onda. Se não sobrar nenhuma, a onda continua só com o plano do Claude: sem comparação, sem travar, e a onda conta para o limite.
- R7: o momento da checagem do limite de despachos é ambíguo, e também é ambíguo se os despachos que falharam contam.
  - Reportado por: claude (P2).
  - Disposição: G21. O limite é checado antes de abrir uma onda, uma onda aberta sempre termina, e os despachos que falharam contam.
- R8: o corpo do prompt do planejador não traz a instrução de somente leitura / de não gravar nem a cláusula anti-delegação (RISK-001), e a sua origem não está definida.
  - Reportado por: claude (P5).
  - Disposição: o prompt é montado a partir do corpo do `implementation-planner`, do texto da `implementation-planning` e do contexto da onda. O corpo carrega três coisas literalmente: a instrução de não gravar, uma cláusula anti-delegação redigida para planejamento (redação exata fixada no lote de implementação) e o bloco que permite leitura. Estender o AC-05 para checar que o prompt as contém.
- R9: a escrita às cegas é atribuída ao dono (G2), mas foi uma escolha de design.
  - Reportado por: claude (P6, não bloqueante), cli:grok/grok-4.7 (P4).
  - Disposição: G23 agora a torna uma decisão do dono. A onda 1 é às cegas; as ondas seguintes refinam a partir do vencedor sem ver o novo rascunho do outro. Citar G23.
- R10: "teto atingido" é aplicado de forma inconsistente. G11 trata só do caso de terminar em um limite sem convergência.
  - Reportado por: cli:grok/grok-4.7 (P3).
  - Disposição: a linha "teto atingido" aparece só quando o loop para no limite de ondas ou no limite de despachos sem convergir.
- R11: o AC-01 diz que `PLAN.md` e seus arquivos `TASK.md` são os únicos artefatos, o que conflita com a persistência por onda (G12), a estimativa em dólares (G8) e a lista de bloqueantes em aberto (G11).
  - Reportado por: cli:grok/grok-4.7 (P8).
  - Disposição: listar todos os artefatos. Por onda, sob `wave-N/`: os dois planos, as duas revisões e o registro do vencedor. No final: o `PLAN.md`/`TASK.md` vencedor e um relatório com o motivo da parada, a contagem de despachos, a estimativa em dólares e os bloqueantes em aberto.
- Não bloqueantes:
  - A meta de cobertura de avaliações no Sinal 1 vai além do grill. Reportado por: cli:grok/grok-4.7 (P1). Disposição: suavizar para uma preocupação do plano de implementação.
  - Uma linha de tabela malformada. Reportado por: cli:grok/grok-4.7 (P2). Disposição: corrigir.
  - A base da estimativa em dólares. Reportado por: claude (P7). Disposição: deixar a base para o lote de implementação e dizer isso.

# Revisão: PRD-005 — rodada 2 (mesclada)

Estado avaliado: o rascunho da correção 1 do PRD. Fonte: GRILL.md G1–G23.

| Revisor | Veredito |
|---|---|
| `claude` (document-validator, opus solicitado, efetivo não verificado) | `changes required`: R1–R11 aplicados; 3 bloqueantes (S1–S3), 5 não bloqueantes (N1–N5) |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | `changes required`: R1–R11 não reabertos; 4 bloqueantes (S1–S4) |
| `cli:grok/grok-4.7` | `not-run`: o estado foi substituído pela correção 2 antes de a sua onda rodar |

## Veredito
- `changes required`

## Achados (mesclados). As disposições do Coordenador derivam de G19–G21.

- T1: a regra de rotação não está definida em dois casos: quando uma CLI configurada não resolve nesta máquina, e quando um repasse dentro de uma onda deixa a CLI da onda seguinte pouco clara.
  - Reportado por: claude (S1), cli:codex/gpt-5.6-luna (S2).
  - Disposição:
    - A lista de rotação são as entradas `cli:` resolvíveis de `reviewers.document`, na ordem da lista, tomadas da checagem do início do loop (RF-02).
    - A CLI primária da onda N é a entrada ((N−1) mod k). Um repasse dentro de uma onda (RF-15) tenta as entradas seguintes, em ordem. Ele não move o cursor: a CLI primária da onda seguinte continua dependendo só do número da onda.
    - Uma entrada que deixa de resolver na sua vez é tratada como no RF-15.
    - Acrescentar ACs para uma entrada configurada que não resolve e para o cursor depois de um repasse.
- T2: o PRD nunca diz como decidir se uma onda "cabe" sob o limite de despachos.
  - Reportado por: claude (S2), cli:codex/gpt-5.6-luna (S1).
  - Disposição:
    - Uma onda cabe quando o limite restante ≥ 2 + 2 × R, onde R é o número de revisores configurados em `reviewers.document`: 2 planejadores mais cada revisor julgando os dois planos.
    - As tentativas de repasse além dessa reserva ainda rodam, porque uma onda aberta sempre termina. Elas contam para o limite.
    - Pôr um exemplo trabalhado no AC-11 (para R = 3: uma onda reserva 8).
- T3: artefatos. O AC-01 conflita com o AC-03, e a persistência do fallback está incompleta.
  - Reportado por: claude (S3), cli:codex/gpt-5.6-luna (S4).
  - Disposição:
    - O AC-01 vale só para o caminho do loop, isto é, pelo menos uma CLI resolvível.
    - O fallback segue o contrato completo de persistência do `/plan` atual: `PLAN.md`, `TASK.md` por fatia, ADRs quando couber, `PLAN.review.md` e atualizações de memória pelo Coordenador.
    - O modo loop também grava um `PLAN.review.md` no nível base: a revisão mesclada do plano vencedor final, para que os consumidores do `/plan` ainda o encontrem. O RF-12 e o AC-01 devem dizer isso.
- T4: a linha de status diz "rodada 1 de 2", e o §9 diz que o limite de 2 rodadas cobre "inclusive deste PRD".
  - Reportado por: cli:codex/gpt-5.6-luna (S3), claude (N2).
  - Disposição: escrever "rodada 2 de 6" e remover "inclusive deste PRD".
- Não bloqueantes:
  - N1: o §6 diz que o limite de ondas herda o limite válido em todo o kit. Disposição: declarar que `wave_cap` tem padrão próprio de 6.
  - N3: a última cláusula do AC-05 se contradiz. Disposição: reescrevê-la como "o Coordenador é o único gravador do plano".
  - N4: o PRD não dá teste para "conteúdo de plano reconhecível". Disposição: o sinal mínimo são as seções exigidas do formato da `implementation-planning` (objetivo e fatias). A checagem exata fica adiada para a implementação.
  - N5: o PRD rotula decisões como citações "textuais". Disposição: escrever "decisão do dono, Gxx (registro do grill)".
  - Todos os itens não bloqueantes foram reportados pelo claude.

## Rodada 2: correção 2 despachada ao `product-discovery`.

# Revisão: PRD-005 — rodada 3 (mesclada)

Estado avaliado: o rascunho da correção 2.

| Revisor | Veredito |
|---|---|
| `claude` (document-validator, opus solicitado, efetivo não verificado) | `changes required`: T1–T4 e os ids N aplicados, nenhum id R regrediu; 3 bloqueantes (U1–U3), 1 não bloqueante (U4) |
| `cli:grok/grok-4.7` | `changes required`: 2 bloqueantes mais o R11 parcial; 2 não bloqueantes |
| `cli:codex/gpt-5.6-luna` | `not-run`: estado substituído pela correção 3 antes da sua onda |

## Veredito
- `changes required`

## Achados (mesclados)
O Coordenador decidiu cada disposição abaixo, derivando cada uma das decisões do dono citadas.

- V1 — O repasse não diz se a varredura dá a volta depois da última entrada da rotação. Reportado por: claude (U1), cli:grok/grok-4.7 (U2).
  - Disposição (G20, "não trava"): o repasse é circular. Ele começa na entrada seguinte à primária da onda e tenta cada uma das k entradas no máximo uma vez por onda. A onda roda sem comparação só quando as k entradas falharam nela.
  - Afirmar isso no AC-17.
- V2 — "Melhor vencedor até aquele ponto" (G11) pode significar o vencedor da última onda ou o melhor vencedor entre as ondas. Reportado por: claude (U2).
  - Disposição (G11 + G18): o melhor plano até então é o vencedor com menos achados bloqueantes entre todas as ondas rodadas. Um empate vai para o de menos achados no total, e um empate que restar vai para a onda mais recente.
  - O RF-12 grava esse plano como `PLAN.md`, e a sua revisão mesclada vira o `PLAN.review.md` base.
  - Afirmar isso no AC-12.
- V3 — O caminho do loop não registra ADRs nem atualizações do MEMORY, e o RF-01 ainda descreve a saída como só `PLAN.md` e `TASK.md`. Reportado por: claude (U3), cli:grok/grok-4.7 (R11 parcial).
  - Disposição: o caminho do loop também grava ADRs quando couber, e o Coordenador atualiza `.harness/MEMORY.md`, como no `/dev-harness:plan`.
  - O RF-01 descreve o artefato final como o contrato completo do `/plan` mais os artefatos extras do loop, sem um "só" conflitante.
  - O AC-01 os lista.
- V4 — Uma onda pode ser lida como uma passada única de pontuação ou como um ciclo de correção aninhado do `/plan`. Reportado por: cli:grok/grok-4.7 (U1).
  - Disposição (G5: cada onda é uma rodada de revisão): cada candidato é revisado uma vez por onda. Nenhum planejador é redespachado dentro de uma onda, e os bloqueantes em aberto alimentam a onda seguinte (RF-08). O ciclo de correção do `/plan` vale só para o fallback sem CLI.
  - Acrescentar um critério: um veredito `changes required` não redespacha o planejador dentro da mesma onda, e uma onda faz 2×R despachos de revisão, mais qualquer repasse do RF-15.
- Não bloqueantes:
  - O item de "Não entra" sobre mais de duas CLIs é reescrito como o claude sugeriu. Reportado por: claude (U4).
  - O risco do §8 diz por engano que o `document-review` tem avaliadores `type: llm`. Corrigir ou remover. Reportado por: cli:grok/grok-4.7 (U3).
  - A referência que credita G18–G23 por ampliar o papel da CLI deve citar G6 e G15. Reportado por: cli:grok/grok-4.7 (U4).

## Rodada 3 — correção 3 despachada ao `product-discovery`.

# Revisão: PRD-005 — rodada 4 (mesclada)

Estado avaliado: o rascunho da correção 3.

| Revisor | Veredito |
|---|---|
| `claude` (document-validator, opus solicitado, efetivo não verificado) | `approved`. V1–V4 aplicados, sem regressão. 3 achados não bloqueantes (W1–W3). |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | `changes required`. 1 achado bloqueante. |
| `cli:grok/grok-4.7` | `not-run`. A sua onda nunca rodou porque a correção 4 substituiu o estado antes. |

## Veredito
- `changes required` (um revisor bloqueia)

## Achados
- X1: a redação do RF-15 deixa um leitor tentar de novo a CLI primária que falhou dentro da mesma onda, inclusive quando ela é a única CLI configurada.
  - Reportado por: cli:codex/gpt-5.6-luna (W1).
  - O claude leu o mesmo texto como correto ("conta a tentativa da própria primária"). O sentido está assentado; isto é redação ambígua, não divergência sobre o comportamento.
  - Disposição (G20): o repasse tenta cada outra CLI resolvível no máximo uma vez e nunca tenta de novo a primária que falhou da onda dentro dessa onda. Com k = 1, uma falha vai direto para uma onda de candidato único. Espelhar isso no AC-17.
- Não bloqueantes (claude):
  - W1: mover o `PLAN.review.md` base do AC-01 para a lista de arquivos que compõem o contrato do `/plan`.
  - W2: no RF-12 e no AC-14, dizer que o Coordenador atualiza `.harness/MEMORY.md` fora do diretório base.
  - W3: no RF-06, uma onda de dois candidatos gasta 2×R revisões e uma onda de candidato único gasta R, mais qualquer repasse.

## Rodada 4: correção 4 despachada ao `product-discovery`.

# Revisão: PRD-005 — rodada 5 (mesclada)

Estado avaliado: o rascunho da correção 4.

| Revisor | Veredito |
|---|---|
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | `approved`: X1 e W1–W3 aplicados |
| `cli:grok/grok-4.7` | `changes required`: 1 achado bloqueante |
| `claude` (document-validator) | `not-run`: a correção 5 substituiu o estado antes da sua onda; ele aprovou a rodada 4 |

## Veredito
- `changes required`

## Achados
- Z1: o RF-01, o RF-12 e o AC-14 põem os ADRs, e no RF-01 o `.harness/MEMORY.md`, dentro de `.harness/tasks/<id>/`. O `/dev-harness:plan` grava os ADRs em `.harness/adr/ADR-<n>.md` e a memória em `.harness/MEMORY.md`, ambos fora desse diretório. O fallback (AC-03) gravaria, portanto, em um lugar diferente do caminho do loop.
  - Reportado por: cli:grok/grok-4.7 (Y1). Confirmado contra `.commands/plan.md`, linhas 9 e 22.
  - Disposição: "sob `.harness/tasks/<id>/`" vale só para o `PLAN.md`, os arquivos `TASK.md` por fatia, o `PLAN.review.md` base, `wave-N/` e o relatório final. Os ADRs, quando couber, vão para `.harness/adr/ADR-<n>.md`, e o Coordenador atualiza `.harness/MEMORY.md`. Os dois caminhos são os que o `/dev-harness:plan` já usa.

## Rodada 5: correção 5 despachada. A rodada 6 é a última sob o limite de 6.

# Revisão: PRD-005 — rodada 6 (mesclada, rodada final: limite de 6 atingido)

Estado avaliado: o rascunho da correção 5.

| Revisor | Veredito |
|---|---|
| `claude` (document-validator, opus solicitado, efetivo não verificado) | `approved`. Z1 aplicado, sem regressão, sem achado novo. |
| `cli:grok/grok-4.7` | `approved` |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | `changes required`. Z1 aplicado; 1 novo achado bloqueante (AA1). |

## Veredito
- `changes required`. Um revisor bloqueia, e o limite de rodadas foi atingido. Pela regra do kit, o documento vai ao dono com o ponto em aberto listado abaixo.

## Ponto em aberto
- AA1: o conjunto de CLIs candidatas não exige um fornecedor diferente do Claude.
  - Reportado por: cli:codex/gpt-5.6-luna.
  - Por que importa: G6 diz que o plano concorrente vem de uma CLI externa "de outro fornecedor". A `external-clis` também aceita `claude` como binário de CLI, então uma entrada `cli:claude/<slug>` configurada deixaria o Claude competir contra si mesmo.
  - Recomendação do Coordenador: excluir as entradas `cli:claude/...` da lista de rotação dos planejadores no RF-05 e no AC-06. Se só resolverem entradas do mesmo fornecedor, usar o fallback sem loop (RF-02). É uma mudança de uma frase. Fica para o dono porque o limite de rodadas foi atingido.
