# Revisão — PRD-019 / T-1901

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

## Rodada 1 (2026-10-12)

- Status da revisão: `request changes`
- Revisores: claude (document-validator, opus). Revisão proporcional padrão: 1 revisor, primeira entrada de `reviewers.document`.

## Achados

### BLOCKING: P1 `THIRD_PARTY_NOTICES.md` é um arquivo de documentação novo, não pedido, creditado ao dono
- Reportado por: claude
- Local: PRD R1, seção 4, seção 6; T-1901 R1, seção 5, Autorização.
- Resolução: decisão do dono em 2026-10-12 — nada de arquivo de notices e nada de comentário de origem no código. Remover de R1, de docs, do conjunto de escrita e da seção 6.

### BLOCKING: P2 R5 "byte a byte" versus o gravador `JSON.parse`/`stringify` do original
- Reportado por: claude
- Resolução (padrão do Coordenador, registrado como premissa): igualdade semântica — as outras chaves lidas e iguais em profundidade, formatação normalizada para JSON com 2 espaços, como faz o original.

### BLOCKING: P3 Os controles de modelo e de esforço do Mode não têm aceite
- Reportado por: claude
- Resolução: acrescentar critério — a página Mode lista modelos, níveis de esforço e modos de nova sessão; o esforço é aplicado em `turn.step`; o modelo é trocado por `/model`; provado com a portagem dos dois testes do motor do original.

### MINOR: P4 Nome e local do arquivo de backup não declarados
- Reportado por: claude
- Resolução: `~/.claude/settings.json.dh.bak`, gravado uma vez se ausente; a origem da restauração é a chave do store, como no original.

### MINOR: P5 O texto do segmento em R4 parece saída traduzida
- Reportado por: claude
- Resolução: citar o texto atual de `dhProgress` sem traduzir (`… (k blocked)`).

### MINOR: P6 R5 pede escrita atômica além do original
- Reportado por: claude
- Resolução: remover "atômica"; manter "arquivo inválido/ilegível é erro e nunca é sobrescrito" (validação na fronteira).

## Rodada 2 (2026-10-12)

- Status da revisão: `request changes`
- Revisores: claude (document-validator, opus). Limite de rodadas atingido (2 de 2); os pontos em aberto vão para o dono.
- Assentados: P1, P2, P4, P5, P6 resolvidos; P3 com a parte do modelo resolvida e a parte do esforço em aberto.

## Achados

### BLOCKING: P3 (herdado) O esforço aplicado em `turn.step` não tem teste que possa falhar
- Reportado por: claude
- Sugestão: R9 acrescenta "depois de escolher `effort-high`, um `turn.step` com `effort` definido vai para `next` com `effort: 'high'`".

### BLOCKING: P7 A checagem "nada de arquivo de notices" de R1 não pode falhar
- Reportado por: claude
- Sugestão: trocar `git diff --stat THIRD_PARTY_NOTICES.md` por `test ! -e THIRD_PARTY_NOTICES.md`.

### MINOR: P8 R3 ainda aponta para os comentários de origem de R1
- Reportado por: claude
- Sugestão: `grep -rn "extra-mods\|/tools" adapters/claude-code/plugin/hooks` vazio.

### MINOR: P9 Não está claro se o Mode grava o backup `.dh.bak`
- Reportado por: claude
- Sugestão: acrescentar "(o mesmo `.dh.bak`, se ausente)" depois da cláusula do Mode.

## Disposição do dono (2026-10-12)

- O dono aceitou aplicar literalmente as sugestões da rodada 2 (P3, P7, P8, P9) sem uma terceira rodada. A execução de T-1901 pode começar.
