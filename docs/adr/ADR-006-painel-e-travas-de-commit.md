# ADR-006 · Painel de status e travas de commit por mod do plugin

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-06 |
| Decisor | Daniel Lemos (dono) |
| Origem | pedido do dono 2026-10-06 (T-1201); mod pessoal "painel" como especificação executável |
| Reversibilidade | barata: remover `panel.ts`, `panel.test.ts` e `panel.engine.test.ts` e voltar `modules` do `hooks.json` para `runtime-guard.ts` |

## 1. Contexto

- Pergunta: como entregar a todo consumidor a barra de status e as travas de commit que hoje só existem num mod pessoal do dono.
- Fato: o mod pessoal não é portátil; só roda na máquina do dono.
- Fato: o ADR-005 já estabeleceu mods TypeScript no plugin, com falha aberta.
- Fato: o `hooks.json` aceita uma só entrada em `modules` (`claude plugin validate`: "names one hooks module per plugin; a second entry is refused"); por isso `panel.ts` é a entrada única e encadeia as travas do `runtime-guard.ts`. `runtime-guard.ts` continua dono da lógica das travas e do teste unitário delas; o encadeamento em `panel.ts` é coberto por `panel.engine.test.ts` (spawn com `parentAgentId` negado; `Write`/`Edit` em `.harness/MEMORY.md` com `agentId` negado).
- Observado na implementação: um segundo `on("agent.spawn")` sem matcher no mesmo plugin também é recusado ("on(\"agent.spawn\") is registered twice without a matcher"), então o encadeamento é feito dentro do mesmo handler.

## 2. Decisão

- Escolhida: `hooks/panel.ts` no plugin, com três peças. (1) Barra de status abaixo do prompt (`$.ui.status`): contexto, agentes em execução, progresso das tarefas do dh, `gate pending`, limites `5h` e `wk` (🔴 a partir de 80%) e custo. (2) Trava de `git commit`, `gh pr create/edit` e `gh release create` cujo texto menciona IA ou Claude; a mensagem cita a política do dev-harness, nunca uma pessoa. (3) Trava de commit quando houve edição depois do último gate verde da sessão.
- Falha aberta: erro no mod libera a ferramenta; as travas são conveniência, não segurança.
- Desvios do mod pessoal, aprovados pelo dono: D1 rótulos curtos em inglês; D2 as negativas próprias do painel (menção a IA e gate) e os avisos de 60%/80% seguem `language` do `.harness/project.yaml` (`pt-br`, senão inglês); as negativas encadeadas do `runtime-guard.ts` continuam em inglês; D3 a negativa de menção a IA cita a política do kit; D4 o aviso de 80% sugere `/compact` ou `/dh:handoff`; D5 `Edit`/`Write`/`NotebookEdit` sob `.harness/` não marcam o gate como pendente.
- Altera o ADR-005 §2 e a reversibilidade: `runtime-guard.ts` deixa de ser declarado em `modules`. Ele exporta `denySpawn` e `guardFile`, que `panel.ts` importa e encadeia antes da própria lógica, e o `register` dele não é carregado pelo plugin. Reverter o ADR-005 agora exige tirar esse import e esse encadeamento de `panel.ts`.
- Descartadas: manter como mod pessoal (não portátil); posição `AbovePrompt` (a barra fica em `$.ui.status`).

## 3. Consequências

- Só vale no Claude Code. A API de mods é de acesso antecipado: rodar `claude plugin test` a cada release.
- `runtime-guard.ts` foi alterado só para exportar `denySpawn` e `guardFile`.
- O estado do gate é por sessão e zera quando o mod recarrega.
- Lacuna conhecida: escrita por `Bash` (redirect, `sed -i`) não marca o gate como pendente.
- Repositório sem gate reconhecível mantém o commit negado depois de uma edição até que algum comando de gate conhecido rode com sucesso.
- Segundo artefato de runtime que não é Go, no mesmo regime do ADR-005; o ADR-001 continua valendo para os scripts do kit.

## 4. Revisar quando

- A API de mods mudar ou sair do acesso antecipado.
- A trava de gate negar commit sem motivo (falso positivo) em projeto real.
