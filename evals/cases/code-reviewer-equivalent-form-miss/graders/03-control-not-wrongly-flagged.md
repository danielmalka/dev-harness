---
type: llm
weight: 1
focus: last_message
---
Todo comando de lint (`.lint` de root, server e worker em
`scripts/workspaces.config.js`) e o próprio `classifyLintStep` são tratados
como funcionando na resposta, sem nenhum achado que alegue que um deles está
classificado errado, quebrado ou incompleto? Um achado sobre um comando de
teste (por exemplo `server.test` caindo em `generic`, mesmo que cite o
despacho compartilhado de `classify()` ou de `normalizeCommand`) não é sobre um
comando de lint e não muda esta resposta; uma observação cosmética de nomenclatura
que não alegue classificação errada também deixa o caminho do lint tratado como
funcionando.
