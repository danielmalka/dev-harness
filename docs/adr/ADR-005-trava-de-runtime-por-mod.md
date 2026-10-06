# ADR-005 · Trava de runtime por mod do plugin

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-05 |
| Decisor | Daniel Lemos (dono) |
| Origem | PRD-011; RISK-001 (subagente ignorou `disallowedTools` e despachou outro agente) |
| Reversibilidade | barata, mas ver ADR-006: o `runtime-guard.ts` é encadeado por `panel.ts`, então reverter exige tirar o import e o encadeamento de lá |
| Alterado por | ADR-006 (composição em `panel.ts`; `runtime-guard.ts` não é mais declarado em `modules`) |

## 1. Contexto

- Pergunta: como impedir que um subagente despache agentes ou escreva `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, se até aqui a regra só existe em prosa.
- Fato: `disallowedTools` no frontmatter não foi imposto pelo runtime (RISK-001, 2026-09-21).
- Fato: o Claude Code oferece mods em plugins (hooks de função em TypeScript, API de acesso antecipado) que interceptam `agent.spawn` e `tool.call`.

## 2. Decisão

- Escolhida: um mod no plugin (`hooks/runtime-guard.ts`, declarado em `modules` no `hooks.json`) com duas travas. `agent.spawn` é negado quando há `parentAgentId` (só o Coordenador despacha). `tool.call` de `Write` e `Edit` é negado quando há `agentId` e o alvo termina em `.harness/MEMORY.md`, `EPOCHAL.md` ou `RISKS.md`.
- Falha aberta: erro no mod libera a chamada; a prosa dos agentes continua como segunda camada.
- Os hooks clássicos (`dh snapshot event`) ficam intactos.
- Descartada: voltar a `disallowedTools` (não imposto) e travar `Bash` por análise de comando (frágil, falso positivo).

## 3. Consequências

- Só vale no Claude Code. Em outros runtimes continua a prosa.
- Lacuna conhecida: escrita por `Bash` (redirect, `sed -i`) não é bloqueada.
- A API de mods é de acesso antecipado e pode mudar entre versões do Claude Code: rodar `claude plugin test` a cada release.
- Primeiro artefato de runtime que não é Go. Ele roda no motor do Claude Code, não como script nosso; o ADR-001 continua valendo para os scripts do kit (`dh`).

## 4. Revisar quando

- A API de mods mudar ou sair do acesso antecipado.
- Aparecer escrita em `.harness/` por `Bash` de subagente (então avaliar travar `tool.call` de `Bash`).
- O kit ganhar adaptador para outro runtime.
