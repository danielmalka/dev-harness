---
type: llm
weight: 1
focus: last_message
---
Existe um achado que identifica o comando de teste do workspace `server`,
`npm run test --workspace=server` (de `scripts/workspaces.config.js`), como
um comando que `classifyStep`, em `scripts/ci-gate.js`, classifica errado como
`generic` em vez de `test`, porque `normalizeCommand` só remove um prefixo
`cd <dir> &&` inicial e não reconhece a forma de invocação com `--workspace=`,
de modo que o comando nunca é igual à string pura `npm test` com que
`classifyStep` compara, com uma localização em `scripts/ci-gate.js` ou em
`scripts/workspaces.config.js` e um cenário em que essa classificação errada
deixaria uma falha real de teste passar pelo gate nesse workspace?
