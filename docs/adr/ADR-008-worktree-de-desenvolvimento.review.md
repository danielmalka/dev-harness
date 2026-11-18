# Revisão: ADR-008

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

## Rodada 1

| Campo | Valor |
|---|---|
| Documento | docs/adr/ADR-008-worktree-de-desenvolvimento.md |
| Fonte | docs/prd/PRD-028-worktree-and-approval-page.md |
| Rodada | 1 de 2 |
| Veredito | `changes required` |
| Atualizado | 2026-10-28 |

Correções na rodada 1: o ADR pelo arquiteto; as regras R30, R61, R26 e R27 do PRD-028 pelo product-discovery. As categorias de P1 a P5 foram atribuídas a partir das descrições de uma linha no registro do Coordenador.

## Achados

- **P1** [applied] [conflict] (blocking) O ADR dizia que `mode` não é validado na leitura, mas o código o valida.
  - Local: ADR-008, regra sobre a leitura de `mode`
  - Evidência: `cmd/dh/config.go` valida `mode`.
  - Resolução: ADR corrigido pelo arquiteto, rodada 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) R30 e D4 listavam caminhos de confiança que não existem num checkout normal.
  - Local: ADR-008 D4; PRD-028 R30
  - Resolução: ADR e PRD corrigidos, rodada 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] (blocking) A condição do passo 1 do clone era ambígua entre um diretório `.git` e um gitfile.
  - Local: ADR-008, passo 1 do clone
  - Resolução: ADR corrigido pelo arquiteto, rodada 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [ambiguity] (blocking) `linked_worktree.gitdir` não estava definido; significa o diretório administrativo `<C>/worktrees/<name>`.
  - Local: ADR-008, definição de `linked_worktree`
  - Resolução: ADR corrigido pelo arquiteto, rodada 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [ambiguity] (blocking) O ADR não dizia que um ponteiro de volta relativo é resolvido a partir do diretório administrativo.
  - Local: ADR-008, resolução do ponteiro de volta
  - Resolução: ADR corrigido pelo arquiteto, rodada 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- Um achado não bloqueante foi levantado na rodada 1. Seu texto não está no registro do Coordenador, por isso não é reproduzido aqui. Registrado como corrigido na rodada 1.

## Não levantado

- Checagens da rodada 1 registradas pelo Coordenador: o congelamento ficou limpo e `dh validate` reportou 0 erros. O revisor rodou como `qa-verifier` na árvore viva.

## Rodada 2

| Campo | Valor |
|---|---|
| Documento | docs/adr/ADR-008-worktree-de-desenvolvimento.md |
| Fonte | docs/prd/PRD-028-worktree-and-approval-page.md |
| Rodada | 2 de 2 (limite atingido) |
| Veredito | `approved` |
| Atualizado | 2026-10-28 |

## Achados

- P1 a P5 trazidos da rodada 1, todos `applied`; fechados na rodada 2.
- **P6** [applied] [gap] (non-blocking) R61 não nomeava os campos de `linked_worktree`.
  - Local: PRD-028 R61
  - Resolução: corrigido pelo Coordenador em uma linha, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- Checagens da rodada 2 registradas pelo Coordenador: o congelamento ficou limpo, nenhum processo do agy ou do CodeGraph ficou rodando depois da rodada e `dh validate` reportou 0 erros.
- Nenhum achado bloqueante está em aberto.
- O status do ADR foi definido como `aceito` depois da rodada 2 (registro do Coordenador).
