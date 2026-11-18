# Revisão: PRD-022

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-022-remove-visual-mods.md |
| Fonte | Decisões do dono 2026-10-14 (despacho do Coordenador, T-2201); árvore de hooks antes da mudança (dist/claude-code/dev-harness/hooks) |
| Rodada | 1 de 2 |
| Veredito | changes required |
| Atualizado | 2026-10-14 |

## Achados

- **P1** [pending] [weak criterion] O grep da R3 não consegue passar, porque `--exclude-dir` casa só nomes-base, então docs/prd e docs/adr não são excluídos e PRD-019 e PRD-022 dão hit. Ele também não pega `mods-types.ts`, `lib/status.ts`, `lib/zen.ts` e `mods.engine.test.ts`.
  - Local: §3 R3
  - Sugestão: `git ls-files adapters dist | grep -E 'mods-|lib/(status|status-line|zen|zen-band|menu-pages|modes|settings)\.'` vazio; `git grep` com os pathspecs `':!docs/adr' ':!docs/prd' ':!CHANGELOG.md'` e o padrão mais amplo, com hits só nos cards de release do roadmap
  - Reportado por: claude
- **P2** [pending] [weak criterion] Nenhuma regra prova a remoção dos avisos de 60%/80% (`$.ui.toast` em panel.ts). A R1 casa os stubs do motor `on('ui.status')` nos arquivos de teste mantidos.
  - Local: §2 item dos avisos, §3 R1
  - Sugestão: grep da R1 `'\$\.ui\.(status|toast)|PromptHint|zen_progress|ZEN_PROMPT|/mods'` com `--exclude='*.test.ts'`
  - Reportado por: claude
- **P3** [pending] [weak criterion] O `dh validate` não checa a contagem de comandos do mod "19 + N" (doc_counts.go só compara "N comandos" com os arquivos em disco), então a R4 não consegue falhar.
  - Local: §3 R4
  - Sugestão: grep de `19 \+ 2|2 comandos registrados|2 commands registered|/mods` em AGENTS.md, nos dois READMEs e nos dois tutoriais sem nenhum hit, mais `dh validate` verde
  - Reportado por: claude
- **P4** [pending] [ambiguity] (`non-blocking`) Os itens do modo de permissão e da contagem de agentes não nomeiam arquivo. O mapa `running` em panel.ts poderia ficar como estado morto.
  - Local: §2
  - Reportado por: claude
- **P5** [pending] [organization] (`non-blocking`) O status diz "aprovado em 2026-10-14" antes da validação e antes de o dono ver o PRD.
  - Local: cabeçalho
  - Reportado por: claude

## Não levantado

- Toda decisão do dono de 2026-10-14 tem lugar na §2, §3, §5 ou §6. Nenhum excesso e nenhum conflito interno.
- Os padrões da R1 e da R3 não dão hit em nenhum dos arquivos mantidos (runtime-guard.ts, snapshot-writer.ts, status-drift.ts, dashboard-command.ts, lib/gate.ts), exceto pelos stubs de teste da P2.
- R2, R5, R6 e R7 podem ser checadas; a prova de clone limpo se aplica corretamente.
- O estilo cumpre as regras de redação de PRD.

## Rodada 2 (2026-10-14)

- Status da revisão: `changes required` (1 `blocking`, novo). P1-P5 aplicados. Limite de rodadas atingido; a P6 vai para o dono.
- P1 [applied] A R3 usa `git ls-files` e `git grep` com pathspecs `:!` e o padrão mais amplo; os hits da árvore atual estão só nos roadmaps.
- P2 [applied] A R1 cobre `$.ui.toast` e exclui `*.test.ts`.
- P3 [applied] A R4 tem uma checagem por grep.
- P4 [applied] A §2 nomeia os arquivos e o mapa `running`.
- P5 [applied] O status é "rascunho".
- **P6** [assentado: o dono escolheu a opção A em 2026-10-14; `/mods` removido do padrão da R4] [weak criterion] O termo `/mods` da R4 dá hit no texto correto das notas de release em `README.md:20` e `README.en.md:20`, então a R4 não consegue passar. Sugestão: remova `/mods` do padrão da R4; a R1 o cobre no código.
  - Reportado por: claude
