# Revisão · ADR-006

- Documento: `docs/adr/ADR-006-painel-e-travas-de-commit.md`
- Fontes: `.harness/tasks/T-1201/TASK.md`, mod de especificação do dono, `adapters/claude-code/plugin/hooks/{panel.ts,runtime-guard.ts,hooks.json}`, `docs/adr/ADR-005-trava-de-runtime-por-mod.md`, `AGENTS.md`
- Validador: `document-validator` (opus), 2 rodadas (teto atingido)
- Review status: approve

## Findings

### MAJOR: ADR-005 continuava dizendo que `runtime-guard.ts` é declarado em `modules` (P1)
- Reported by: claude
- Rodada 1: bloqueante. Rodada 2: resolvido. ADR-006 §2 registra que altera o ADR-005 §2 e a reversibilidade; ADR-005 ganhou a linha "Alterado por | ADR-006" e a nota de reversibilidade.

### MAJOR: cobertura do encadeamento afirmada sem teste (P2)
- Reported by: claude
- Rodada 1: bloqueante. Rodada 2: texto corrigido, alegação pendente de conferência do Coordenador.
- Conferência do Coordenador (2026-10-06): os casos existem em `adapters/claude-code/plugin/hooks/panel.engine.test.ts` (spawn de subagente negado; Edit de subagente em `.harness/MEMORY.md` negado) e passam em `claude plugin test`; o ADR passou a citar `panel.engine.test.ts`.

### MAJOR: D2 ambíguo sobre as negativas do `runtime-guard.ts` (P3)
- Reported by: claude
- Resolvido: só as negativas do painel e os avisos de 60%/80% seguem `language`; as do `runtime-guard.ts` ficam em inglês.

### MINOR: recusa de `on` duplicado rotulada como fato sem evidência (P4)
- Reported by: claude
- Resolvido: rotulado "Observado na implementação", com a mensagem de erro citada.

### MINOR: escopo da trava omitia `gh release create` (P5)
- Reported by: claude
- Resolvido.
