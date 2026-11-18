# Revisão: PRD-011

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-011-trava-de-runtime-por-mod.md |
| Fonte | Decisão do dono 2026-10-05 ("Pode 1 e 2, dentro do plugin") + fatos verificados pelo Coordenador (despacho); AGENTS.md; .harness/RISKS.md RISK-001 |
| Rodada | 2 de 2 |
| Veredito | `approved` |
| Atualizado | 2026-10-05 |

## Achados

- **P1** [applied] [`conflict`] `claude plugin test` foi rotulado como fato verificado, mas não estava nos fatos verificados da fonte; a checagem em sessão real não estava em nenhum ticket.
  - Local: seção 1 linhas 22-24; T-1102 linhas 100-103
  - Resolução: rodada 2. O fato agora se limita à verificação do Coordenador pós-rodada-1 (`--help`, `claude-code/testing`, teste de fumaça). O T-1102 exige uma sessão real com o plugin construído: `Agent`/`Write .harness/MEMORY.md` do subagente negados, sessão principal permitida.
  - Nota do Coordenador (2026-10-05): verificado após a rodada 1 — `claude plugin test` existe na 2.1.289 (`claude plugin test --help`; um `*.test.ts` de fumaça que importa `claude-code/testing` rodou com 1 aprovado numa sonda em scratchpad).
  - Reportado por: claude
- **P2** [applied] [`weak criterion`] O T-1102 poderia fechar com a R9 "não verificado" e contornar o gate de release.
  - Local: R9 linhas 72-75; linha 130
  - Resolução: rodada 2. "Não verificado" ou uma falha de carga do plugin → o release espera uma decisão explícita do dono registrada no PRD. O T-1102 testa 2.1.287 e 2.1.288, cuja presença foi confirmada.
  - Reportado por: claude
- **P3** [applied] [`conflict`] O marcador de runtime do AGENTS.md se tornaria falso.
  - Local: R8(b) linhas 67-69
  - Resolução: rodada 2. O marcador de runtime cita a exceção do mod em TS (ADR-004), conferida por `grep -n ADR-004 AGENTS.md`.
  - Reportado por: claude
- **P4** [applied] [`weak criterion`] A R8 não dizia qual texto do coordinator.md muda; o grep não nomeava nenhuma string.
  - Local: R8(a) linhas 65-67, 70-71
  - Resolução: rodada 2. O alvo é o marcador "sole writer" (coordinator.md linha 37); o bloco cercado não muda em nenhum dos dois arquivos; a checagem é `grep -n "enforced at runtime"` mais `checkAntiDelegationClause` verde.
  - Reportado por: claude
- **P5** [applied] [`ambiguity`] A R7 deixava uma varredura opcional de `.ts` a cargo do builder.
  - Local: R7 linhas 62-64; T-1101
  - Resolução: rodada 2. Nenhuma mudança em Go; o limite é registrado; conferido por `git diff --stat` sem `cmd/`, `internal/` nem `go.mod`.
  - Reportado por: claude
- **P6** [applied] [`excess`] `NotebookEdit` e um `MultiEdit` condicional.
  - Local: R4 linha 50
  - Resolução: rodada 2. Só `Write` e `Edit`.
  - Reportado por: claude
- **P7** [applied] [`ambiguity`] A correspondência de caminho da R4 não tinha âncora na raiz.
  - Local: R4 linhas 50-55
  - Resolução: rodada 2. Correspondência por sufixo no caminho normalizado (`/` e `\`), qualquer raiz, com um teste para um caminho fora do projeto.
  - Reportado por: claude
- **P8** [applied] [`organization`] O ticket de teste escrevia docs; o T-1102 não tinha o campo Documentação.
  - Local: R9; seção 4 linha 84; T-1102 linha 104
  - Resolução: rodada 2. Docs no fim do lote pelo docs-guide; T-1102 "Documentação: nenhum".
  - Reportado por: claude

## Não levantado

- Nenhum escopo além dos itens 1 e 2; exclusões da R10 intactas; a R2 permite exatamente dois handlers.
- R3/R4 ainda batem com os fatos verificados `agent.spawn`/`parentAgentId` e `tool.call`/`agentId`; R5 e R6 sem mudança.
- As decisões da linha 128 batem com as resoluções do Coordenador dadas ao autor.
- Observação (não bloqueante): a cláusula de despacho/memória da R8(b) no AGENTS.md não tem grep próprio; perdê-la deixa o AGENTS.md verdadeiro, não falso.
- Observação: a flag `claude --plugin-dir` não foi verificada pelo validador; o aceite do T-1102 vale seja qual for o método de carga.
- As seis categorias foram reconferidas no texto revisado; identificadores únicos.
