# PRD-031 · Fluxo enxuto do builder e test-first como padrão

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-31 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-30 / 2026-10-30 |
| Tickets | a definir pelo planejador |
| Origem | Pedido do dono em 2026-10-30, fonte `tasks/PRD-031/SOURCE.md` no harness dir |
| Depende de | PRD-029 (triagem de achados, base de versão 0.35.0) |
| Entrega | Uma pull request |

## 1. Problema

O fluxo da fase do builder é longo e gasta muitos tokens. O dono pede um fluxo mais direto e concreto.

- Fato: quatro arquivos carregam a fase do builder. Somam 830 linhas e 107478 bytes (`wc`, 2026-10-30).

| Arquivo | Linhas | Bytes |
|---|---|---|
| `.skills/incremental-implementation/SKILL.md` | 183 | 15466 |
| `.skills/implementation-planning/SKILL.md` | 190 | 17098 |
| `.agents/coordinator.md` | 207 | 49460 |
| `.skills/code-review/SKILL.md` | 250 | 25454 |

- Fato: o `ub-mode` (Uncle Bob) cobre a fase do builder em cerca de 70 linhas de regras imperativas de uma linha. O `SKILL.md` inteiro dele tem 152 linhas e 3958 bytes.
- Hipótese: a maior parte do peso vem das etapas e do estado por fatia, não da prosa. As etapas são os ids de achado `F<round>-<n>`, seis estados de achado (`.agents/coordinator.md`, item 8: `open`, `fixed`, `intentional`, `noise`, `escalated`, `reopened`), o protocolo de churn, duas rodadas limpas consecutivas, o cursor de rotação da etapa de código, as regras do pacote de revisão e a dúvida no meio do build. Esta hipótese é falsa se as etapas somarem menos da metade dos quatro arquivos. O inventário de regras (R7) a testa.
- Fato: uma feature não tem regra de test-first. O `.agents/coordinator.md` diz: `A feature has no red-first rule, so this single dispatch is the disproof.` O red-first existe só para os testes de regressão de bug.
- Fato: nenhuma regra diz o que o corpo de uma mensagem de commit deve conter.
- Fato: `.skills/harness-authoring/SKILL.md` não tem regra sobre o tamanho de um passo de procedimento.

## 2. Solução

Depois da entrega:

- Um `builder` escreve primeiro o teste que falha, vê o teste falhar pelo motivo esperado e só então muda o código.
- A `regression-testing` tem três heurísticas que rejeitam testes vazios.
- O corpo de um commit diz por quê, o que mudou e como a mudança foi provada.
- Um passo de procedimento novo numa skill do kit é uma frase imperativa mais, no máximo, uma frase de motivo.
- Os quatro arquivos da fase do builder miram 30% do tamanho de hoje (teto de 40%, R6), sem perder nenhuma salvaguarda que sustenta o fluxo.
- A maquinaria de revisão por fatia fica mais simples. O dono aceita ou rejeita cada simplificação (R11 a R15). O dono aceitou as cinco em 2026-10-30.

A Parte A é R1 a R5 (comportamento novo). A Parte B é R6 a R10 (compressão e consistência). A Parte C é R11 a R15 (remoções, aceitas pelo dono em 2026-10-30).

## 3. Regras

### Parte A · Comportamento novo

- R1 A `incremental-implementation` manda o `builder`, em cada fatia de implementação, escrever o teste que falha, rodá-lo, ver a falha esperada e só então mudar o código. O relatório do `builder` nomeia o resultado que falhou. Fatias só de documentação, só de configuração e só de texto ficam isentas, e o relatório diz qual isenção vale. A `regression-testing` traz o mesmo padrão em seu Procedimento, então uma fatia de feature e uma correção de bug começam com um teste que falha. Verificável: `incremental-implementation` e `regression-testing` contêm, cada uma, a frase exata `expected failure`. A `incremental-implementation` contém também as frases `docs-only`, `config-only`, `copy-only` e uma linha de Checagens `Red:` no seu Formato de saída. A `regression-testing` contém `test-first`.
- R2 A dúvida no meio do build se estreita. O test-first (R1) substitui o motivo `disproof`, então o `.agents/coordinator.md` não diz mais que uma feature não tem regra de red-first. A dúvida fica só para uma fatia com `Doubt: yes`. A regra R14 define a forma que sobra.
- R3 A `regression-testing` contém três heurísticas, cada uma numa frase:
  - Veja o teste novo falhar pelo motivo esperado antes de ele valer.
  - Apague um teste que ainda passa quando toda função sob teste não devolve nada.
  - Quando um teste honesto é impraticável, reporte a checagem como `not-run`. Não escreva um teste fraco.
  A regra de red-first existente para testes de regressão de bug continua.
- R4 Regra de conteúdo do commit. O corpo do commit diz por que a mudança existe, o que mudou e como foi provada. O "como" nomeia a checagem que rodou. O cabeçalho de Conventional Commits do `AGENTS.md` continua sem mudança. A regra fica só em `.skills/delivery-readiness/SKILL.md` (veja Decisões, D1). Verificável: `grep -n "body" .skills/delivery-readiness/SKILL.md` mostra uma regra.
- R5 `.skills/harness-authoring/SKILL.md` contém esta regra de estilo: um passo de procedimento é uma frase imperativa mais, no máximo, uma frase de motivo, e as referências cruzadas vão no fim do passo ou da seção, não dentro de uma frase. A regra vale para todo passo que este PRD escreve ou reescreve.

### Parte B · Compressão e consistência

- R6 Orçamento de tamanho, um teto flexível. Depois da entrega, os quatro arquivos da seção 1 miram 250 linhas e 30000 bytes. O teto é 330 linhas e 43000 bytes (cerca de 40% de hoje). O teto pode ser excedido em até cerca de 10% (363 linhas, 47300 bytes) só quando o inventário de regras (R7) nomeia as regras que precisam do espaço extra. Não é um gate rígido. Meça com `wc -lc` nos quatro caminhos. Nenhum texto vai para um arquivo novo para cumprir o orçamento. O limite de bytes existe porque o arquivo do Coordenador guarda 46% dos bytes. Verificável: o total de `wc -lc` está no teto ou abaixo dele, ou o inventário tem uma seção `Budget excess` que nomeia cada regra que precisa do espaço e o total fica dentro de 10%.
  Decisão do dono em 2026-10-31: o total medido (432 linhas / 55,830 bytes, −48% de bytes) é aceito acima do teto; o próximo corte (mover seções de um só comando para fora de `coordinator.md`) é a R6 do PRD-032.
- R7 Inventário de regras. Antes de qualquer reescrita, o `builder` escreve `docs/prd/PRD-031-rule-inventory.md`. Ele lista toda regra dos quatro arquivos, uma linha cada: arquivo de origem e regra, depois `kept`, `merged into <row>` ou `removed` com o motivo. Uma regra removida sob a Parte C cita o número da sua regra. Depois da entrega, toda linha tem um status final, e as contagens de `kept`, `merged` e `removed` aparecem no fim do arquivo.
- R8 Estas salvaguardas continuam em vigor, nos quatro arquivos ou num arquivo que o inventário nomeia:
  - Honestidade de `not-run`: uma checagem que não rodou é `not-run`, nunca `passed`.
  - Autorização: a aprovação do plano não autoriza commit, push, publicação nem deploy.
  - Conjunto de escrita: um especialista escreve só dentro do conjunto que recebeu.
  - Revisão independente: uma afirmação do `builder` não é revisão, e a revisão é somente leitura.
  - Posse da memória: só o Coordenador escreve `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.
  - Cláusula anti-delegação, literal em `coordinator.md` e na `external-clis`, com sua fixação por sha256.
  - RISK-001 (a cláusula fica no texto do prompt) e RISK-002 (nenhuma escrita na árvore nem no harness dir durante uma chamada de CLI externo).
  Verificável: `go run ./cmd/dh validate` passa, inclusive `checkAntiDelegationClause`, e o inventário tem uma linha `kept` ou `merged` para cada item.
- R9 Nenhum agente, skill, comando ou chave de `project.yaml` novo. As contagens ficam: 13 agentes, 22 skills, 19 + 1 comandos. Não existe a chave `commands.quality`.
- R10 Premissa de limpeza. A mesma pull request atualiza todo documento que afirma um comportamento que este PRD muda (seção 4). `dist/` é regenerado com `go run ./cmd/dh build`. Uma busca por cada termo removido (`F<round>`, `reopened`, `rotation`, `cursor`, `Churn`, `two consecutive`, e qualquer outro termo que o dono aceite remover) não acha nenhuma menção restante fora de `docs/prd/`, `docs/adr/`, do histórico do `CHANGELOG.md` e de `evals/baselines/`. Verificável: `go run ./cmd/dh validate` e um `git diff` limpo de `dist/` depois de `go run ./cmd/dh build`.

### Parte C · Remoções de comportamento (aceitas pelo dono em 2026-10-30)

Cada regra dá o custo de hoje e o risco que volta. R16 e R17 foram propostas e rejeitadas (veja Fora de escopo). O limite de seis rodadas e o registro do pacote em `MEMORY.md` continuam.

- R11 Uma rodada limpa encerra o laço de revisão de código. O laço para depois da primeira rodada limpa, não depois de duas rodadas limpas consecutivas. Uma rodada continua "limpa" pela definição de hoje: nenhum achado bloqueante está `open` ou `escalated`, e todo revisor despachado devolveu um veredito. A revisão final da feature segue a mesma regra.
  - Custo hoje: uma revisão completa extra da mudança para cada ticket e para a revisão final. O PRD-029 (0.34.0) a acrescentou.
  - Risco se removida: uma falha de um revisor permanece. As checagens determinísticas e a revisão independente continuam, e o dono pode pedir o painel completo.
  - Recomendação: aceitar.
- R12 O cursor de rotação da etapa de código é removido. Uma entrada de código e toda reentrada usam a primeira entrada de `reviewers.code` (`claude` quando listado). O dono pede o painel completo quando quiser. O Coordenador não escreve mais um cursor em `MEMORY.md`.
  - Custo hoje: cerca de 11 linhas de regra, um campo de `MEMORY.md` por tarefa e a lógica de nova tentativa da "entrada `not-run`". A rodada de confirmação (R11) é o principal motivo dela.
  - Risco se removido: uma reentrada não muda mais de fornecedor, então um ponto cego de um revisor se repete (propósito da ADR-004).
  - Recomendação: aceitar. Depende da R11. Se a R11 for rejeitada, mantenha o comportamento atual da R12.
- R13 Os estados de achado encolhem de seis para três: `open`, `closed`, `escalated`. Uma correção que o `builder` afirma ter feito fica `open` até um revisor confirmar. `closed` cobre corrigido e confirmado e resolvido com evidência (`intentional` ou `noise`). Os ids de achado `F<round>-<n>` são removidos. Um achado é identificado por seu local e seu título.
  - Custo hoje: uma máquina de estados com seis saídas, um esquema de ids e uma entrada de `MEMORY.md` por achado.
  - Risco se removido: um achado reaberto fica mais difícil de ligar ao seu primeiro reporte, e a lista de resolvidos numa reentrada fica menos exata.
  - Recomendação: aceitar.
- R14 A dúvida no meio do build permite um ciclo, não três. O Coordenador classifica os achados em duas classes, acionável ou não. Se a revisão devolve achados e o Coordenador não classifica nenhum como acionável, o Coordenador escala ao dono. O gatilho "decisão em aberto com opção estrutural O2" é removido. Só as fatias com `Doubt: yes` o disparam.
  - Custo hoje: até três revisões de contexto novo antes da revisão de fim de fatia, uma classificação em quatro vias e a linha `Doubt` em toda fatia.
  - Risco se removida: uma leitura errada de contrato ou um erro de fronteira de confiança é pego um ciclo depois, na revisão de fim de fatia.
  - Recomendação: aceitar. É o estreitamento que a R2 exige. Remover a dúvida por inteiro não é recomendado, porque o test-first não pega uma leitura errada de contrato.
- R15 A linha `Churn:` e o protocolo de churn de três opções são substituídos por uma frase. Quando um achado novo reverte uma correção aceita antes, o laço para e o Coordenador faz ao dono uma pergunta numerada.
  - Custo hoje: um revisor lê o delta em toda reentrada para escrever um valor de `Churn:` para cada achado, e um protocolo de decisão separado vem em seguida.
  - Risco se removida por inteiro: a oscilação de corrigir e reverter corre até o limite de rodadas.
  - Recomendação: substituir pela frase única, não remover por inteiro.

## 4. Docs

Documentos que afirmam um comportamento que este PRD muda. O `builder` confirma cada um com uma busca e o edita na mesma pull request. Os pares são atualizados juntos (pt-br e `docs/en/`).

- `AGENTS.md` (itens "Verificação padrão", "PRD e tickets", "Revisão por CLI externa" e "Agentes")
- `.agents/coordinator.md`
- `.agents/builder.md`
- `.agents/reviewer.md`
- `.agents/qa-verifier.md`
- `.commands/build.md`
- `.commands/plan-loop.md`
- `.skills/incremental-implementation/SKILL.md`
- `.skills/implementation-planning/SKILL.md`
- `.skills/code-review/SKILL.md`
- `.skills/regression-testing/SKILL.md`
- `.skills/delivery-readiness/SKILL.md`
- `.skills/harness-authoring/SKILL.md`
- `.skills/external-clis/SKILL.md`
- `templates/en/TASK.md`
- `templates/pt-br/TASK.md`
- `docs/tutorial.html`
- `docs/en/tutorial.html`
- `docs/novo-fluxo-prd.md`
- `docs/en/new-prd-flow.md`
- `docs/roadmap.html`
- `docs/en/roadmap.html`
- `CHANGELOG.md` (entrada nova, aumento de versão e a linha de contagens)
- `docs/adr/ADR-004-duvida-e-rotacao.md` (nota de status: substituída em parte pela nova ADR)
- `docs/adr/ADR-009-lean-builder-flow.md` (novo: registra as remoções aceitas)
- `docs/prd/PRD-031-rule-inventory.md` (novo: R7)

O `builder` deve também buscar em `adapters/` e `internal/` código que leia um campo removido de `MEMORY.md` (cursor, estados de achado, ids). Uma ocorrência está no escopo e dispara a prova de clone limpo do `AGENTS.md`.

## Apêndice

### Fora de escopo

- Registro por ação: acrescenta contabilidade, o oposto deste PRD.
- Ferramentas e limiares `crapper`, `dryer` e `mutator`: ferramental novo que o kit não roda.
- Push direto em `main`: conflita com as regras de branch e de revisão do `AGENTS.md`.
- Limpeza de worktree sem trava: as regras de Limpeza do `coordinator.md` protegem os arquivos ignorados.
- `Build the Lever`, `Redesign From First Principles`, `Attack the Premise`: escopo sem limite que nenhum aceite consegue testar.
- Uma chave `commands.quality` nova: a R9 proíbe chaves novas.
- Planejamento iterativo e artefatos de desenvolvimento visual: o dono os chama de direção futura ("norte"), e eles não bloqueiam este PRD.
- R16, limite de 3 rodadas para o laço de build/QA e de revisão de código: rejeitada pelo dono em 2026-10-30. O limite de seis rodadas continua.
- R17, nenhum registro em `MEMORY.md` dos elementos do pacote de revisão: rejeitada pelo dono em 2026-10-30. O registro de Despachos continua.
- Compressão de arquivos além dos quatro da seção 1, e remoção das regras de revisão da etapa de documento.

### Decisões (fechadas pelo dono, 2026-10-30)

- D1 A regra de conteúdo do commit fica em `.skills/delivery-readiness/SKILL.md` (R4). Fechada.
- D2 O orçamento é um teto flexível (R6): meta de 250 linhas e 30000 bytes, teto de 330 linhas e 43000 bytes, com cerca de 10% de margem quando o inventário justifica. Fechada.
- D3 R11 a R15 aceitas. R16 e R17 rejeitadas. Fechada.
- D4 Versão 0.36.0 (aumento minor a partir da 0.35.0). Fechada.

## Resumo executado

- **Entregue**: O `builder` agora escreve primeiro um teste que falha e o laço de revisão para depois de uma rodada limpa, com o cursor de rotação, os seis estados de achado e a dúvida de três ciclos removidos (R11 a R15), mas os quatro arquivos da fase do builder somam 432 linhas e 55,830 bytes, acima do teto da R6 que o dono aceitou em 2026-10-31.
- **Regras**: Os caminhos do repositório são relativos à raiz do repositório; os caminhos do harness têm o prefixo `<harness dir>/`. A matriz de QA `<harness dir>/tasks/T-3111/TASK.md` (seção 7) e `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/` cobrem as regras marcadas como Honrada.
  - R1 Honrada. `.skills/incremental-implementation/SKILL.md` (linha de Checagens `Red:`, linha 49); `.skills/regression-testing/SKILL.md` (`test-first`). Evidência: matriz de QA e `qa-T-3111/`.
  - R2 Honrada. `.agents/coordinator.md` não diz mais `disproof`; a dúvida vale só para `Doubt: yes`. Evidência: matriz de QA (`disproof` ausente).
  - R3 Honrada. `.skills/regression-testing/SKILL.md` (linhas 47, 53, 54). Evidência: matriz de QA.
  - R4 Honrada. `.skills/delivery-readiness/SKILL.md` (linha 58), a única regra de corpo. Evidência: matriz de QA.
  - R5 Honrada. `.skills/harness-authoring/SKILL.md` (linha 58). Evidência: matriz de QA.
  - R6 Não honrada. Medido: 432 linhas / 55,830 bytes contra o teto de 330 / 43,000 (os quatro arquivos, `wc -lc`). Motivo: a reescrita cortou cerca de 48% dos bytes, mas não até o teto. Decisão: o dono aceitou o excesso em 2026-10-31 (registrado sob a R6); o próximo corte (mover seções de um só comando para fora de `coordinator.md`) é a R6 do PRD-032.
  - R7 Honrada. `docs/prd/PRD-031-rule-inventory.md` (toda linha tem um status final, com as contagens no fim). Evidência: matriz de QA e `qa-T-3111/state.txt`.
  - R8 Honrada. Salvaguardas mantidas nos quatro arquivos ou no inventário; `go run ./cmd/dh validate --source-only .` passa, inclusive `checkAntiDelegationClause`. Evidência: `qa-T-3111/validate-source.txt`, matriz de QA.
  - R9 Honrada. As contagens ficam em 13 agentes, 22 skills, 19 + 1 comandos; nenhuma chave nova. Evidência: `qa-T-3111/validate-source.txt`, matriz de QA.
  - R10 Honrada. A busca pelos termos removidos não acha nenhuma menção restante fora dos caminhos excluídos. Evidência: `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/r10.txt`, matriz de QA. A regeneração de `dist/` e a prova de clone limpo rodaram em T-3113 e T-3114 (`<harness dir>/tasks/PRD-031/evidence/qa-T-3114/`).
  - R11 Honrada. Uma rodada limpa encerra o laço de revisão de código. Evidência: matriz de QA; a revisão de código convergiu em `<harness dir>/tasks/PRD-031/evidence/review-batch-r1..r3/`; a revisão final convergiu em `<harness dir>/tasks/PRD-031/evidence/review-final-r1..r3/`.
  - R12 Honrada. Nenhum cursor de rotação; as reentradas usam a primeira entrada de `reviewers.code`. Evidência: matriz de QA.
  - R13 Honrada. Três estados de achado (`open`, `closed`, `escalated`); ids `F<round>-<n>` removidos. Evidência: matriz de QA, `qa-T-3111/r10.txt`.
  - R14 Honrada. Um ciclo de dúvida; só as fatias com `Doubt: yes` o disparam. Evidência: matriz de QA.
  - R15 Honrada. A linha `Churn:` e o protocolo de três opções são substituídos por uma frase. Evidência: matriz de QA.
  - R16 Não implementada. Rejeitada pelo dono em 2026-10-30 (limite de 3 rodadas); o limite de seis rodadas continua.
  - R17 Não implementada. Rejeitada pelo dono em 2026-10-30; o registro do pacote em `MEMORY.md` continua.
- **Tickets**: Todos com veredito `done` (`concluída` em cada `TASK.md`). A evidência de cada um é a seção Resultado de `<harness dir>/tasks/T-31NN/TASK.md`.
  - T-3101 a T-3110, lane `backend` (`harness-maintainer`, skill `harness-authoring`).
  - T-3111, lane `teste` (unitário), `qa-verifier`; evidência `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/`.
  - T-3112, lane `backend` (`docs-guide`, skill `doc-template-html`).
  - T-3113, lane `infra` (`builder`, modo build); regenerou `dist/` e o manifesto.
  - T-3114, lane `teste` (integração), `qa-verifier`; evidência `<harness dir>/tasks/PRD-031/evidence/qa-T-3114/`.
- **Docs**: `AGENTS.md`; `docs/adr/ADR-009-lean-builder-flow.md`; `docs/adr/ADR-004-duvida-e-rotacao.md`; `docs/tutorial.html`; `docs/en/tutorial.html`; `docs/novo-fluxo-prd.md`; `docs/en/new-prd-flow.md`; `docs/roadmap.html`; `docs/en/roadmap.html`; `CHANGELOG.md`; `docs/prd/PRD-031-rule-inventory.md`.
- **Fora de escopo**: Achados `Minor` herdados para o PRD-032:
  - A `code-review` perdeu a pergunta explícita `would a plausible regression fail a test?`.
  - A regra de página oficial da `incremental-implementation` perdeu três salvaguardas: pular a busca quando há uma só opção viável; um conflito com a convenção local vai para o gate de decisão em aberto; a URL vai no relatório, não num comentário de código.
  - O total do inventário de regras é 55,808 bytes, enquanto outros documentos dizem 55,830.
  - O passo 13 da `regression-testing` tem quatro frases.
