# PRD-022 · Remover os mods visuais que o dono agora mantém em um repositório separado

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-14 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-14 / 2026-10-14 |
| Tickets | T-2201 |
| Origem | Chat com o dono, 2026-10-14 (escopo decidido no chat) |
| Depende de | PRD-019 (acrescentou o menu de mods) |
| Habilita | nada |

## 1. Problema

- O dono agora mantém mods pessoais do Claude Code em um repositório separado, `malka-claude-mods`. Os plugins dele são `statusline`, `menu`, `cleanscreen`, `mode`, `permissions` e `sentinel`.
- O mod `dh` em `adapters/claude-code/plugin/hooks/` desenha as mesmas coisas. Duas barras de status agora entram em conflito na tela.
- `panel.ts` chama `$.ui.status('')` em `refresh()`. O motor fixa um status de plugin vazio, porque só `undefined` o limpa. O prompt então mostra um aviso "⚠ dh:" vazio.
- O kit não precisa dessas partes visuais para funcionar. A trava de memória, o gate de commit e o snapshot do dashboard não as usam.

## 2. Solução

Remover os componentes visuais do mod `dh`. Manter todas as travas e o dashboard.

Remover:

- A linha de status sob o prompt (`lib/status-line.tsx`, `lib/status.ts`, o render de `PromptHint`, o estado `dhSegments` e `statusDrawn`) e a barra alternativa `$.ui.status(...)` em `refresh()` de `panel.ts`. Cópia: `statusline`.
- O menu `/mods` (`mods-menu.tsx`, `lib/menu-pages.tsx`, `lib/settings.ts`). Cópia: `menu`.
- A faixa Zen, a ferramenta `zen_progress` e a injeção de `ZEN_PROMPT` (`mods-zen.tsx`, `lib/zen.ts`, `lib/zen-band.tsx`). Cópia: `cleanscreen`.
- O modo de modelo e esforço (`lib/modes.ts`, a sobrescrita de esforço em `turn.step`). Cópia: `mode`.
- O modo de permissão padrão (em `mods-menu.tsx`, `lib/settings.ts` e `lib/menu-pages.tsx`). Cópia: `permissions`.
- A contagem de agentes em execução (o mapa `running` e o segmento `agents N` em `panel.ts`). Remover o mapa também, a menos que o snapshot o use. Cópia: `sentinel`.
- Os avisos de contexto a 60% e 80%. Nenhuma cópia existe. O dono aprovou a remoção e enviou a ideia para a sessão do `malka-claude-mods`.

Manter:

- `runtime-guard.ts` (trava de memória, ADR-005).
- `snapshot-writer.ts` (sessões ao vivo do dashboard).
- O gate de commit (`dirty-after-gate`, recusa de menção a IA).
- `status-drift.ts`.
- A seção de prompt `dh:harness`.
- O comando `/dashboard` (`dashboard-command.ts`).

O dashboard continua mostrando o progresso dos tickets por PRD (`internal/dashboard/projects.go`, `projectProgress`). Esse código não muda.

Versão 0.29.0. A contagem de comandos é 19 comandos de arquivo e 1 comando de mod (`/dashboard`), não 19 e 2.

## 3. Regras

- R1 Nenhum componente visual permanece em `adapters/claude-code/plugin/hooks/`. Não há nenhum render de `PromptHint`, nenhuma chamada `ui.status`, nenhuma chamada `ui.toast`, nenhum comando `/mods` e nenhuma ferramenta `zen_progress`. Verificável: `grep -rnE '\$\.ui\.(status|toast)|PromptHint|zen_progress|ZEN_PROMPT|/mods' adapters/claude-code/plugin/hooks --exclude='*.test.ts'` não tem nenhum hit.
- R2 As travas mantidas funcionam como antes. Verificável: os testes do mod para `runtime-guard`, `snapshot-writer`, o gate de commit e `status-drift` estão verdes.
- R3 Nenhum órfão permanece: nenhum arquivo, teste, import, referência, doc ou cópia gerada de uma parte removida (premissa de limpeza no `AGENTS.md`). Verificável: `git ls-files adapters dist | grep -E 'mods-|lib/(status|status-line|zen|zen-band|menu-pages|modes|settings)\.'` é vazio; `git grep -nE 'mods-menu|mods-zen|mods-types|menu-pages|status-line|zen-band|lib/(modes|settings|status|zen)' -- . ':!docs/adr' ':!docs/prd' ':!CHANGELOG.md'` só tem hits em cards de release do roadmap.
- R4 A contagem de comandos é 19 + 1 em todo lugar que a declara. Verificável: `grep -nE '19 \+ 2|2 comandos registrados|2 commands registered' AGENTS.md README.md README.en.md docs/tutorial.html docs/en/tutorial.html` não tem nenhum hit; `dh validate` verde.
- R5 O build está completo. Verificável: `dh validate` e `go run ./cmd/dh build` passam com `dist/` regenerado; `go test ./...` e os testes do mod passam.
- R6 A mudança não toca uma configuração `statusLine` do usuário. O dono confirmou que nenhuma foi salva. Verificável: `git diff --name-only main...HEAD` não lista nenhum arquivo de configuração do usuário.
- R7 Versão 0.29.0. `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte (`v0.29.0`) mudam juntos. A prova de clone limpo vale, porque o PR toca `dist/*/bin`. Verificável: o `grep` dos três valores mostra `0.29.0` ou `v0.29.0`; o `sha256sum` dos binários é igual no clone limpo.

## 4. Docs

- `AGENTS.md` (lista de mods; "19 + 2" vira "19 + 1")
- `README.md`, `README.en.md`
- `docs/tutorial.html`, `docs/en/tutorial.html`
- `docs/roadmap.html`, `docs/en/roadmap.html` (entrada 0.29.0)
- `CHANGELOG.md` (0.29.0)
- `docs/prd/PRD-019-mods-menu.md`: nota de cabeçalho "substituído em parte pelo PRD-022".
- `docs/adr/ADR-006-painel-e-travas-de-commit.md`: nota de cabeçalho "substituído em parte pelo PRD-022".
- `.claude-plugin/marketplace.json` (`version` e `ref`)

## 5. Fora de escopo

- Qualquer mudança em `runtime-guard.ts`, `snapshot-writer.ts`, no gate de commit, em `status-drift.ts`, na seção `dh:harness` ou em `/dashboard`.
- Qualquer mudança no dashboard (`internal/dashboard`).
- O repositório `malka-claude-mods` e os seus plugins.
- Uma substituição, dentro do kit, das funções removidas.

## 6. Decisões do dono

- Dono, 2026-10-14: remover do mod `dh` só o que tem cópia no `malka-claude-mods`, e remover a linha de status inteira mesmo sem um segmento `dh` ali.
- Dono, 2026-10-14: remover os avisos de contexto a 60% e 80%. A ideia foi para a sessão do `malka-claude-mods`.
- Dono, 2026-10-14: as travas, o gate de commit, `status-drift.ts`, a seção `dh:harness` e `/dashboard` ficam.
- Causa raiz, fato: `panel.ts` chama `$.ui.status('')`. O motor fixa um status de plugin vazio. O aviso "⚠ dh:" vazio some com a remoção.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** o mod não desenha mais uma linha de status, um menu `/mods`, uma faixa Zen, uma página Mode nem avisos de contexto; a trava de memória, o registro de snapshot, o gate de commit, o status drift, a seção `dh:harness` e `/dashboard` ficam, e o aviso "⚠ dh:" vazio sumiu.
- **Regras:**
  - R1 Honrada: o `grep` da R1 não tem nenhum hit.
  - R2 Honrada: `claude plugin test` em `dist/claude-code/dev-harness` dá 226 passam, 0 falham, incluindo um teste novo para o caminho alternativo `<cwd>/.harness` (uma mutação que removeu o caminho alternativo o fez falhar).
  - R3 Honrada: a checagem de `git ls-files` é vazia; a checagem de `git grep` só tem hits em cards de release do roadmap.
  - R4 Honrada: o `grep` não tem nenhum hit (`/mods` removido do padrão pela decisão P6 do dono, opção A); `dh validate` verde.
  - R5 Honrada: `dh validate` com 0 erros, `go run ./cmd/dh build` idempotente, `go test ./...` verde.
  - R6 Honrada: nenhum arquivo de configuração do usuário está na mudança.
  - R7 Honrada: `0.29.0` em `internal/build/build.go` e em `.claude-plugin/marketplace.json`, `ref` `v0.29.0`; prova de clone limpo: 7 de 7 binários idênticos byte a byte.
- **Tickets:** T-2201 (backend, harness-maintainer) aprovado. Matriz de QA por `qa-verifier`; revisão de código por `claude` (2 achados não bloqueantes, corrigidos) e, na reentrada, por `cli:codex/gpt-5.6-luna` (`approve`, 0 achados).
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `docs/prd/PRD-019-mods-menu.md` (nota), `docs/adr/ADR-006-painel-e-travas-de-commit.md` (nota), `.claude-plugin/marketplace.json`.
- **Fora de escopo:** avisos de contexto no mod `statusline` do malka-claude-mods (ideia enviada a esse projeto), qualquer mudança no dashboard, em `/dashboard` ou nas travas mantidas.
