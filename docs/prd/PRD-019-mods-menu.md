# PRD-019 · Menu de mods (`/mods`): status line, Zen e Mode no mod do dh

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-12 (as partes visuais foram removidas pelo PRD-022, 0.29.0) |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-12 / 2026-10-12 |
| Tickets | T-1901 |
| Origem | chat do dono, 2026-10-12; repositório `github.com/Jhonatan-de-Souza/ClaudeMods` (plugin `extra-mods` 0.4.1, commit `9b4b59f9b9683465e5044828f00c6a4abd1e5f63`, clonado só para leitura) |
| Depende de | PRD-012 (mod `panel.ts`, `/dashboard`), BUG-001 (0.24.1, eventos `classic.*` não disparam em conta Team) |
| Libera | nada |

## 1. Problema

Fatos verificados em 2026-10-12:

- O plugin `extra-mods` tem um menu lateral (botão de canto "◆ extra-mods ▲" ou `/tools`) com quatro ferramentas: Zen mode (esconde as linhas de ferramenta, registra a ferramenta `zen_progress`, injeta `ZEN_PROMPT` em `prompt.submit`, faixa de progresso acima do prompt, segue `TodoWrite`/`TaskCreate`/`TaskUpdate`), Theme (busca `terminalcolors.com`, grava `~/.claude/themes` e as configurações do Windows Terminal), Mode (modelo por `/model`, esforço em `turn.step`, `permissions.defaultMode` em `~/.claude/settings.json` com backup; lê o modo por `classic.UserPromptSubmit`/`classic.PostToolUse`) e Status line (modelo, pasta, branch, contexto, limites 5h/7d com reinício, cache, pílula de shells em segundo plano; ligada, remove o `statusLine` do usuário de `~/.claude/settings.json` com backup e o restaura ao desligar; ligada por padrão).
- O repositório de origem não tem arquivo de licença. O dono falou com o autor, que disse que não há licença porque qualquer pessoa pode copiar e usar sem atribuição.
- O mod do dh é `adapters/claude-code/plugin/hooks/panel.ts`, o único módulo do `hooks.json` (um módulo por plugin; `on(evento)` duplicado é recusado pelo carregador), encadeando `runtime-guard.ts`, `snapshot-writer.ts`, `dashboard-command.ts` e `status-drift.ts`. A barra atual é `$.ui.status(...)` em `refresh()`: `ctx ▰▰▱ N%`, `agents N`, `dh ▰▰▱ done/total (k blocked)` (conta `tasks/*/TASK.md` da pasta do harness por `dhProgress`), `gate pending`, limites 5h/wk, custo; mais avisos de contexto em 60/80%. A `description` do `hooks.json` diz "local files only, no network".
- BUG-001 (0.24.1): eventos `classic.*` não disparam em conta Team; o mod usa só eventos do motor.
- A contagem de comandos está documentada como "19 + 1" (19 arquivos mais `/dashboard` do mod) em `AGENTS.md` e `CHANGELOG.md`.

## 2. Solução

O mod do dh ganha um menu aberto por `/mods` (e por um botão de canto) com três ferramentas copiadas do `extra-mods`: Status line, Zen e Mode. Theme fica de fora (busca de rede, escreve configuração do Windows Terminal, não testado em Linux/WSL, quebraria a promessa "sem rede"). A status line, ligada por padrão como no original, substitui a barra `$.ui.status` atual e mantém no fim da linha os segmentos do dh (`dh ▰▰▱`, `gate pending`, `agents N`); os avisos de contexto em 60/80% ficam. Tudo entra na cadeia de `panel.ts`. O acompanhamento do modo de permissão não usa `classic.*`. Contagem de comandos vira 19 + 2. Sem dependência nova.

## 3. Regras

- R1 Código copiado de `Jhonatan-de-Souza/ClaudeMods` (plugin `extra-mods` 0.4.1, commit `9b4b59f9b9683465e5044828f00c6a4abd1e5f63`): status line, Zen e Mode, nada de Theme. Os trechos copiados ficam em `adapters/claude-code/plugin/hooks/` (módulos `lib` sob a mesma pasta), sem arquivo de avisos de terceiros e sem comentário de origem no código. Verificável: `grep -ri "terminalcolors\|Windows Terminal" adapters/claude-code/plugin/hooks` vazio; `git diff --name-only main...HEAD -- THIRD_PARTY_NOTICES.md` vazio e `git status --porcelain -- THIRD_PARTY_NOTICES.md` vazio (o arquivo existente desde a 0.18.0 não é alterado) e `grep -rn "9b4b59f\|ClaudeMods" adapters/claude-code/plugin/hooks` vazio.
- R2 `hooks.json` continua com um único módulo (`panel.ts`); o código novo entra na cadeia de `panel.ts` (junto de `runtime-guard.ts`, `snapshot-writer.ts`, `dashboard-command.ts` e `status-drift.ts`) e nenhum `on(evento)` fica duplicado. A `description` do `hooks.json` segue verdadeira ("local files only, no network"): nenhum código novo faz requisição de rede. Verificável: `hooks.json` com um só módulo; `grep -rn "fetch(\|http\.\|https\.\|XMLHttpRequest" adapters/claude-code/plugin/hooks/*.ts` sem chamada de rede (fora de `dashboard-command.ts`, que só imprime a URL local); `claude plugin test` verde.
- R3 O comando do menu é `/mods` (não `/tools`) e o botão de canto tem rótulo que identifica o dh (não "extra-mods"); `/mods` abre o menu com exatamente três itens: Status line, Zen e Mode. Verificável: teste do mod (`claude plugin test`) que registra o comando `mods`, lê o rótulo do botão e lista os três itens; `grep -rn "extra-mods\|/tools" adapters/claude-code/plugin/hooks` vazio.
- R4 Status line ligada por padrão, substituindo o `$.ui.status` atual. A linha mostra o conteúdo da original (modelo, pasta, branch, contexto, limites 5h/7d com reinício, cache, pílula de shells em segundo plano) e, no fim da linha, nesta ordem, os segmentos do dh: `dh ▰▰▱ done/total (k blocked)` (texto atual do `dhProgress`, sem tradução), `gate pending` e `agents N`, com as mesmas regras de hoje (contagem de `tasks/*/TASK.md` por `dhProgress`; `gate pending` só quando pendente; `agents` só com N > 0). Os avisos de contexto em 60% e 80% continuam. Verificável: testes de função pura da linha com dh presente, sem dh e com `gate pending`; `status-drift.test.ts` e os testes de `panel` atualizados e verdes.
- R5 Mesmo escopo de escrita do original em `~/.claude/settings.json`: com a status line ligada, o mod remove o `statusLine` do usuário depois de gravar o backup `~/.claude/settings.json.dh.bak` (escrito uma só vez, apenas se ainda não existir) e o restaura ao desligar a partir do valor guardado na chave do armazenamento do plugin (como no original); Mode grava `permissions.defaultMode` (o mesmo `.dh.bak`, se ainda não existir). Paridade com o original (padrão do Coordenador, não decisão do dono): as demais chaves ficam semanticamente iguais (lidas e comparadas por deep-equal) com a formatação normalizada para JSON de 2 espaços, e o `statusLine` restaurado é deep-equal ao valor guardado. Nenhuma outra chave é tocada; `settings.json` ilegível ou inválido é erro e nunca é sobrescrito. Verificável: testes com `settings.json` temporário (HOME em pasta temporária): ligar, desligar, `statusLine` restaurado deep-equal, `defaultMode` alterado com as demais chaves deep-equal, `.dh.bak` criado uma vez e não reescrito, arquivo inválido intocado.
- R6 Nenhum caminho novo depende de eventos `classic.*` (BUG-001): o acompanhamento do modo de permissão usa eventos do motor e, sem eles, mostra "desconhecido". Verificável: `grep -n "classic\." adapters/claude-code/plugin/hooks/*.ts` vazio nos arquivos novos e copiados; teste em que nenhum evento de permissão chega e o modo aparece como desconhecido sem erro.
- R7 Zen: a ferramenta `zen_progress` é registrada pelo mod do dh (nome efetivo `mcp__<plugin>__zen_progress`, com o prefixo que o motor der ao plugin `dh`), o `ZEN_PROMPT` é injetado em `prompt.submit` só com Zen ligado, e o Zen acompanha `TodoWrite`/`TaskCreate`/`TaskUpdate`. Verificável: testes das funções puras do original portados (progresso, prompt, acompanhamento de tarefas) verdes, mais um teste do nome registrado.
- R8 Comportamento atual preservado: `runtime-guard`, trava de commit, `snapshot-writer` (snapshots de sessão com o mesmo schema) e `/dashboard` não mudam. Verificável: `git diff` de `runtime-guard.ts`, `snapshot-writer.ts` e `dashboard-command.ts` vazio (ou só import/registro necessário), e os testes existentes (`runtime-guard.test.ts`, `dashboard.engine.test.ts`, `panel.engine.test.ts`) verdes.
- R9 Os testes do original para as funções puras de status, Zen e Mode são portados para o runner `claude plugin test`, ao lado dos testes de `panel`. A página do Mode lista modelos, níveis de esforço e modos de nova sessão; escolher um esforço o aplica em `turn.step`; escolher um modelo passa por `/model`. Verificável: `claude plugin test` verde com os dois testes de motor do original portados ("the mode page lists effort levels and new-session modes" e "the mode page switches the model through /model"), mais o teste: após escolher `effort-high`, um `turn.step` com `effort` definido segue para `next` com `effort: 'high'`; os arquivos de teste novos existem em `adapters/claude-code/plugin/hooks/`.
- R10 Contagem de comandos passa a "19 + 2" (19 arquivos em `.commands/` mais `/dashboard` e `/mods` do mod); `.commands/` continua com 19 arquivos e sem dependência Go nova. Verificável: `ls .commands | wc -l` = 19; `git diff go.mod` vazio; `dh validate` verde; `grep -rn "19 + 1" AGENTS.md CHANGELOG.md` sem ocorrência vigente (PRDs antigos e entradas históricas do `CHANGELOG.md` ficam).
- R11 Versão 0.26.0: `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte do mesmo arquivo (`v0.26.0`) mudam juntos; o `ref` esquecido é falha de entrega. Verificável: `grep` dos três valores = `0.26.0`/`v0.26.0` no mesmo commit; `go run ./cmd/dh build` sem diferença em `dist/` além da esperada; prova de clone limpo (toca `internal/` e `dist/*/bin`).

## 4. Docs

- `AGENTS.md` (linha "19 comandos em arquivo, mais 1 comando registrado pelo mod" passa a 19 + 2 com `/mods`; linha da pasta global, "contagem segue 19 + 1")
- `README.md`, `README.en.md` (parágrafo de status: status line do dh, `/mods`, `/dashboard`)
- `docs/tutorial.html`, `docs/en/tutorial.html` (onde `/dashboard` é descrito, acrescentar `/mods` e o aviso de que a status line mexe em `~/.claude/settings.json`)
- `CHANGELOG.md` (0.26.0), `docs/roadmap.html`, `docs/en/roadmap.html` (validados contra o `CHANGELOG.md` pelo `dh validate`)
- `adapters/claude-code/plugin/hooks/hooks.json` (conferir a `description`; só muda se deixar de ser verdadeira)

## 5. Fora de escopo

- Theme (busca na rede, escrita das configurações do Windows Terminal).
- Qualquer chamada de rede no mod.
- Mudar `runtime-guard`, a trava de commit, o `snapshot-writer` ou o `/dashboard`.
- Mudar a extensão VS Code.
- Segundo módulo no `hooks.json`.

## 6. Decisões do dono

- Copiar o código do `extra-mods` (opção c): o repositório não tem licença, o autor disse ao dono que qualquer um pode copiar e usar sem atribuição (2026-10-12). Decisão do dono (2026-10-12): copiar sem arquivo de avisos de terceiros e sem comentário de origem no código; o autor disse que é livre para copiar sem menção. A proveniência fica só neste PRD (commit 9b4b59f).
- Trazer Status line, Zen e Mode; Theme fica fora.
- Comando `/mods` (não `/tools`); rótulo do botão identifica o dh; contagem 19 + 2; docs na mesma PR.
- Status line ligada por padrão, substituindo a barra `$.ui.status`, com `dh ▰▰▱`, `gate pending` e `agents N` no fim da linha; avisos de contexto 60/80% ficam.
- Dono, 2026-10-12: aceitou aplicar as sugestões da rodada 2 do validador (P3, P7, P8, P9) sem terceira rodada.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** o mod do dh ganhou o menu `/mods` (botão `◆ dh mods ▲`) com Status line, Zen e Mode copiados do `extra-mods`; a status line, ligada por padrão, substitui a barra `$.ui.status` e termina com `dh ▰▰▱ done/total (k blocked)`, `gate pending` e `agents N`. Saiu na 0.26.0.
- **Regras:**
  - R1 Honrada: greps de `terminalcolors`/`Windows Terminal` e de `9b4b59f`/`ClaudeMods` em `hooks/` vazios; `THIRD_PARTY_NOTICES.md` sem diff nem status.
  - R2 Honrada: `hooks.json` com um módulo e a descrição "no network"; código dividido por registro (`mods-menu.tsx`, `mods-zen.tsx`) e funções puras em `hooks/lib/`, sem `on(evento)` duplicado.
  - R3 Honrada: comando `mods`, botão `◆ dh mods ▲`, teste conta exatamente 3 itens; grep de `extra-mods`/`/tools` vazio.
  - R4 Honrada: teste da ordem dos segmentos do dh no fim da linha e do caso sem dh; avisos 60/80% mantidos; a barra antiga só some depois que a status line desenhou.
  - R5 Honrada: testes com fs em memória (backup `.dh.bak` uma vez, restauração deep-equal, arquivo inválido intocado, `CLAUDE_CONFIG_DIR` só absoluto).
  - R6 Honrada: sem `classic.*` no código novo; modo corrente mostrado como "unknown" (nenhum evento do motor).
  - R7 Honrada: `mcp__dh__zen_progress` registrado; entrada saneada (`parseZenInput`); `ZEN_PROMPT` só com Zen ligado; fail-open testado.
  - R8 Honrada: `runtime-guard.ts`, `snapshot-writer.ts`, `dashboard-command.ts` e `status-drift.ts` sem diff; testes existentes verdes.
  - R9 Honrada: testes portados do original e o teste do esforço em `turn.step` (só no loop principal).
  - R10 Honrada: 19 arquivos em `.commands/`, `go.mod` sem diff, `AGENTS.md`, README, tutorial, CHANGELOG e roadmap em 19 + 2.
  - R11 Honrada: 0.26.0 em `build.go`, `version` e `ref` `v0.26.0` no `marketplace.json`; prova de clone limpo dos binários.
- **Tickets:**
  - T-1901 backend: aprovado (QA APPROVED; revisão de código `claude` approve, `cli:codex/gpt-5.6-luna` approve na reentrada; 2 rodadas de correção; `claude plugin test` 253/0).
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html`.
- **Fora:**
  - Theme (rede e Windows Terminal).
  - Modo de permissão da sessão corrente (sem evento do motor).
  - Escrita atômica do `settings.json` e os demais limites aceitos na revisão.
