# dev-harness

Kit portátil de desenvolvimento assistido por IA. Plano de produto em `docs/plano-produto.html`.

## Decisões fixadas

- Idioma: `.skills/`, `.agents/`, `.commands/` e `profiles/` são escritos em inglês. `templates/` existe em duas versões espelhadas, `templates/en/` e `templates/pt-br/`, e todo template novo entra nas duas. `docs/` é pt-br com espelho em inglês em `docs/en/` (exceto o plano de produto, interno); toda doc nova entra nas duas línguas com link cruzado. `README.md` pt-br e `README.en.md`.
- Memória do projeto consumidor (`.harness/MEMORY.md`, `EPOCHAL.md`, `RISKS.md`): só o Coordenador escreve.
- Quando duas regras conflitarem e a resolução mudar a forma da entrega, perguntar antes de escolher. Não assumir.
- Licença: MIT (`LICENSE`). Redistribuição do kit usa este texto; o `plugin.json` gerado deve repetir `"license": "MIT"`.
- Fontes canônicas: `.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`. `dist/` é gerado por `go run ./cmd/dh build` e não se edita à mão.
- Idioma em execução: o Coordenador responde no idioma em que o dono escreve, inglês por padrão; o `setup` grava `language` (`en` ou `pt-br`) em `.harness/project.yaml`, que escolhe `templates/<lang>/` para registros e artefatos. Despachos, respostas entre agentes, código e commits são sempre em inglês.
- O disco é a verdade sobre `.harness/`: se o dono apagou ou zerou os registros, o Coordenador não reintroduz tarefas, despachos ou evidências lembrados da conversa.
- Extensão VS Code (repo `danielmalka/dev-harness-vscode`, clone em `~/projetos/dev-harness-vscode`): só visualizador e lançador. Lê snapshots em disco escritos pelos scripts de `statusLine`/`subagentStatusLine` e hooks do plugin; abre sessão com o plugin e insere `/dev-harness:<comando>`. Nunca vira runtime nem orquestra pipelines. Não é fork da `context-kit-extension` (essa fica arquivada como doadora de código por cópia).
- Ordem: etapa 1 no plugin (marketplace.json no repo, statusLine, subagentStatusLine e hooks gravando `~/.claude/dev-harness/sessions/<session_id>.json`, estado ativo/ocioso por hooks), validada no danlemos, antes de codar a extensão. Painel v1 só ao vivo, só Claude Code. OTEL e Agent SDK: decisão adiada pelo dono; não planejar em cima deles.
- Runtime dos scripts do kit: Go, um binário `dh` (validate, build, doctor, snapshot) com binários pré-compilados commitados em `dist/.../bin/<os>_<arch>/` para linux amd64/arm64, darwin amd64/arm64 e windows amd64 (ADR-001 em `docs/adr/`). Migração concluída em 20/09/2026: `scripts/` removido; Python só resta na fixture `slice-01`. A extensão VS Code é TypeScript, sem binário Go.
- Etapa 1 do plugin especificada em `docs/prd/PRD-001-etapa-1.md`: marketplace.json, snapshot de sessões e subagentes, hooks de estado. Código só depois da spec aprovada.
- Premissa (pedido do dono em 20/09/2026, PRD-002 em `docs/prd/`): todo documento que vira contrato (PRD, story, plano, ADR) passa por um revisor adversarial somente leitura antes de chegar ao dono; o Coordenador roda o ciclo criação → validação até `approved` ou duas rodadas. Primeiro alvo: validador de PRD após `discover`. Ainda não implementado.
