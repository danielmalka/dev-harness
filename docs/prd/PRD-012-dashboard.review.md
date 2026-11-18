# PRD-012 · Validação

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

Documento: `docs/prd/PRD-012-dashboard.md`. Validador: `document-validator` (opus), skill `document-review`.

## Rodada 1 de 2 (2026-10-07)

- Status da revisão: `changes required`
- Achados bloqueantes: 6

## Achados

### BLOCKING: P1 [ambiguity] Os estados do avatar "erro" e "concluído" não têm evento de origem; "concluído" se sobrepõe a "parado"
- Local: seção 2 item 3; R13; R18
- Reportado por: claude
- Problema: R18 ordena seis estados; R13 define só `working`/`idle`/`waiting`/`closed`. Nenhum evento produz `erro` (jevmon: `StopFailure`) nem `concluído` (jevmon: `Stop`), enquanto R13 mapeia `Stop` para `idle`.
- Correção sugerida: estender R13 com uma tabela evento→estado para cada estado que R18 ordena (por exemplo `StopFailure`→`error`, `Stop`→`done`, decaindo para `idle` após N min), ou removê-los; um teste de mod por mapeamento.

### BLOCKING: P2 [conflict] R13 renomeia `active` para `working` sem dizer o que o mod escreve, quebrando o contrato da extensão do PRD-001
- Local: R11; R13; C5
- Reportado por: claude
- Problema: a extensão conhece `active`/`idle`/`closed`; R13 nunca diz se o mod escreve `active` ou `working`, nem como os leitores antigos tratam `waiting`.
- Correção sugerida: manter `state` no vocabulário do schema 1 e acrescentar um campo novo (por exemplo `activity`), ou escrever os valores novos e listar a mudança da extensão em Docs.

### BLOCKING: P3 [gap] Remover `dh snapshot event` descarta campos que nenhuma regra manda o mod produzir
- Local: R11, R12, R13
- Reportado por: claude
- Problema: hoje os hooks preenchem `session_name`, `started_at`, `cwd`, `agent`, `model.id` e o buffer circular `events`; sem eles, R15 não consegue mapear uma sessão para um projeto.
- Correção sugerida: o mod escreve todos os campos que `dh snapshot event` escreve hoje, a partir de eventos equivalentes do mod; um teste de mod por evento, comparado com `dh snapshot event`.

### BLOCKING: P4 [ambiguity] A pasta jevmon do dono não funciona sem alteração, pois falta um mapa padrão estado→pose; nomes dos arquivos de exemplo indefinidos
- Local: R16, R17
- Reportado por: claude
- Problema: as poses do jevmon se chamam `frente`, `duvida`, `puto`…; não há mapa padrão sem `poses.json`; o fallback só funciona se os dois conjuntos compartilharem os nomes das poses.
- Correção sugerida: R16 adota o mapa padrão do jevmon; R17 nomeia os arquivos de exemplo com esses nomes de pose.

### BLOCKING: P5 [conflict] R21 assenta a decisão pendente 7 enquanto o PRD a lista como em aberto
- Local: R21; Apêndice decisão 7; C2
- Reportado por: claude
- Problema: R21 exige a opção (b) de uma decisão em aberto.
- Correção sugerida: tornar R21 condicional à decisão 7.

### BLOCKING: P6 [gap] A checagem de deriva não consegue pegar a defasagem que a motiva
- Local: R19, R20; seção 1; fonte D4
- Reportado por: claude
- Problema: só "ticket aberto num PRD entregue" é deriva; "todos os tickets `done` / citado como entregue no CHANGELOG, mas o cabeçalho não diz `entregue em`" passa em silêncio.
- Correção sugerida: acrescentar a condição inversa, ou registrar a checagem estreita como decisão em aberto.

### NON-BLOCKING: P7 [conflict, factual] `updated_at` já existe no schema 1
- Local: R15; seção 1
- Reportado por: claude
- Correção sugerida: reescrever como atualizado a cada evento do mod e por heartbeat se H3 valer.

### NON-BLOCKING: P8 [conflict, factual] Contagens de status e citações na seção 1 e em C2 estão levemente erradas
- Local: seção 1; C2
- Reportado por: claude
- Correção sugerida: "dezenas dizem 'pronta'"; "dizem algo diferente de 'entregue em'".

### NON-BLOCKING: P9 [ambiguity] "PRDs entregues recolhidos" virou só uma linha de contagem
- Local: seção 2 item 1; R10
- Reportado por: claude
- Correção sugerida: declarar só a contagem ou uma lista expansível.

### NON-BLOCKING: P10 [organization] Controle das decisões pendentes
- Local: Apêndice decisões 4 e 6; R23, R24
- Reportado por: claude
- Correção sugerida: renomear a decisão 4 para "Animação no conjunto de exemplo"; mover a metade de risco de R23 para Riscos.

## Rodada 2 de 2 (2026-10-07)

- Status da revisão: `changes required`
- Achados bloqueantes: 2 (P11, P12), ambos novos; P1-P10 resolvidos. Limite atingido: vai para o dono com P11 e P12 em aberto.

## Achados

### BLOCKING: P11 [conflict] O decaimento de `done` é medido a partir de `updated_at`, que o heartbeat não para de atualizar; então `done` nunca decai sob a opção de vivacidade recomendada
- Local: R13; R15; recomendação (a) da decisão 3
- Reportado por: claude
- Correção sugerida: o mod escreve `activity_at` só quando `activity` muda; o heartbeat atualiza só `updated_at`; o leitor faz `done` decair quando `now - activity_at >= N`; teste em Go com `updated_at` recente e `activity_at` de 6 min atrás lê `idle`.

### BLOCKING: P12 [conflict] A opção (a) da decisão 7 faz R19(b) falhar neste repositório, então a checagem de R21 não consegue passar e R20 bloqueia todo commit
- Local: opção (a) de R21; R19(b); R20; decisão 7
- Reportado por: claude
- Correção sugerida: declarar essa consequência em 7(a) (viável só com uma lista de isenção em R19), ou descartar (a) e manter só (b).

### NON-BLOCKING: P13 [organization] Docs cita "decisão pendente 5" para a contagem de comandos (é a 9); "3b" não tem opções
- Reportado por: claude

### NON-BLOCKING: P14 [gap, hypothesis] `StopFailure` chegar ao mod não está listado em H2
- Reportado por: claude

### NON-BLOCKING: P15 [ambiguity] R19(b) dispara quando o último ticket fecha, antes da revisão final e do release; não está declarado como intencional nem listado como risco
- Reportado por: claude
