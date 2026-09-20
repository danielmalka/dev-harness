# dev-harness

Kit portátil de desenvolvimento assistido por IA. Plano de produto em `docs/plano-produto.html`.

## Decisões fixadas

- Idioma: `.skills/`, `.agents/`, `.commands/` e `profiles/` são escritos em inglês. `templates/` existe em duas versões espelhadas, `templates/en/` e `templates/pt-br/`, e todo template novo entra nas duas. `docs/` é pt-br com espelho em inglês em `docs/en/` (exceto o plano de produto, interno); toda doc nova entra nas duas línguas com link cruzado. `README.md` pt-br e `README.en.md`.
- Memória do projeto consumidor (`.harness/MEMORY.md`, `EPOCHAL.md`, `RISKS.md`): só o Coordenador escreve.
- Quando duas regras conflitarem e a resolução mudar a forma da entrega, perguntar antes de escolher. Não assumir.
- Licença: MIT (`LICENSE`). Redistribuição do kit usa este texto; o `plugin.json` gerado deve repetir `"license": "MIT"`.
- Fontes canônicas: `.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`. `dist/` é gerado por `scripts/build-claude-code.sh` e não se edita à mão.
- Idioma em execução: o Coordenador responde no idioma em que o dono escreve, inglês por padrão; o `setup` grava `language` (`en` ou `pt-br`) em `.harness/project.yaml`, que escolhe `templates/<lang>/` para registros e artefatos. Despachos, respostas entre agentes, código e commits são sempre em inglês.
- O disco é a verdade sobre `.harness/`: se o dono apagou ou zerou os registros, o Coordenador não reintroduz tarefas, despachos ou evidências lembrados da conversa.
