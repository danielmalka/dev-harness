# Revisão: PRD-020

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-020-prd-ste100.md (+ ~/.harness/projects/dev-harness/tasks/T-2001/TASK.md) |
| Fonte | Decisões do dono 2026-10-13 (despacho do Coordenador, pontos 1-7) |
| Rodada | 2 de 2 |
| Veredito | `changes required` |
| Atualizado | 2026-10-13 |

## Achados

- **P1** [pending] [conflict] Apagar `templates/pt-br/PRD.md` falha na checagem de espelho de `checkTemplateRefs` (validate.go:501-521). Aí é preciso uma mudança em Go, o que contradiz R10 e os pontos 2 e 7 da fonte.
  - Local: §2 linha 30, R4, R8, R10, §6 (c); T-2001 §3.2, R10, §8
  - Evidência: validate.go:514-519; validate_test.go:39
  - Sugestão: manter pt-br/PRD.md como cópia idêntica byte a byte do template em inglês (a contagem continua 14), ou obter uma decisão do dono para uma exceção de PRD em `validate.go` e emendar R10 e §5
  - Reportado por: claude
- **P2** [pending] [conflict] `dh config set language` grava `<home>/config.yaml` (o padrão global), não `project.yaml`.
  - Local: §2 linha 33, R7, §6 (d); T-2001 §3.4
  - Evidência: internal/harness/config_set.go:13-48; docs/tutorial.html:772
  - Sugestão: o Coordenador propõe uma edição de `language` em `<harness dir>/project.yaml` e a grava depois que o dono concordar
  - Reportado por: claude
- **P3** [pending] [gap] `.commands/document.md:18` (`templates/<lang>/PRD.md`) e `.agents/docs-guide.md:59` (PRD escrito em `language`) não estão em R6 nem na autorização de T-2001.
  - Local: R6; T-2001 Autorização
  - Evidência: o `grep` de R6 sobre `.commands` falha; ponto 1 da fonte
  - Sugestão: acrescentar os dois arquivos a R6 e à autorização
  - Reportado por: claude
- **P4** [pending] [gap] A lista de Docs não traz a linha 678 dos roadmaps ("7 templates espelhados por idioma"), a linha 740 dos tutoriais e a linha 69 de README.md/README.en.md.
  - Local: §4; T-2001 §5
  - Evidência: ponto 6 da fonte; doc_counts.go:65-87 (uma linha de um só substantivo não é conferida)
  - Sugestão: listar essas linhas explicitamente
  - Reportado por: claude
- **P5** [pending] [ambiguity] A frase da linha 79 da §6 tem cerca de 70 palavras, acima do limite de 25 palavras de R11 (só de estilo, não bloqueante).
  - Local: §6 linha 79
  - Evidência: R11
  - Sugestão: um item de lista por premissa
  - Reportado por: claude
- **P6** [pending] [weak criterion] O `grep` de R7 não cobre `.commands/setup.md` (não bloqueante).
  - Local: R7
  - Evidência: setup.md:22
  - Sugestão: acrescentar setup.md ao `grep`
  - Reportado por: claude

## Não levantado

- Cada um dos pontos 1-7 da fonte corresponde a §2, R1-R14 e §5. Nenhum excesso foi encontrado.
- As linhas citadas existem: project-onboarding:66, coordinator.md Language (73-74) e Executed summary (70), product-discovery:53, requirements-discovery:98.
- O dashboard já lê `delivered on` (R9). A contagem de templates é 14 hoje, e doc_counts confere só as linhas ancoradas com vários substantivos (R8). Nenhum código em Go nomeia `PRD.md`.
- Os ids R1-R14 são únicos. A lane e o papel de T-2001 batem com o AGENTS.md.

## Disposição do Coordenador (rodada 1)

- P1: o Coordenador escolheu o caminho sem Go: `templates/pt-br/PRD.md` continua como cópia idêntica byte a byte de `templates/en/PRD.md`. Conferência por amostragem de `checkTemplateRefs` em `internal/kit/validate.go:501-521`.
- P2-P6: devolvidos ao autor.

## Rodada 2 (2026-10-13)

- Status da revisão: `changes required` (1 bloqueante, novo)
- P1 [applied] o template de PRD pt-br continua como cópia idêntica byte a byte (§2, R4, R8, R10, §6). Sem mudança em Go.
- P2 [applied] o Coordenador edita `<harness dir>/project.yaml`; `dh config set` é só o padrão da máquina. O texto restante em R7 foi registrado como P7.
- P3 [applied] docs-guide.md:59 e document.md:18 acrescentados a R6 e à autorização de T-2001.
- P4 [applied] linhas 69-70 do README, linha 740 do tutorial e linha 678 do roadmap acrescentadas a §4 e a T-2001 §5.
- P5 [applied] premissas da §6 divididas em itens.
- P6 [applied] `setup.md` acrescentado ao `grep` de R7.
- **P7** [pending] [conflict] R7 ainda dizia que o Coordenador "propõe uma atualização de `project.yaml`", contradizendo "sem segunda confirmação" em §2, §6 e T-2001 §3.4.
  - Local: §3 R7 (linha 44)
  - Sugestão: "a pedido explícito do dono para trocar, edita `language` em `<harness dir>/project.yaml` sem segunda confirmação"; acrescentar `grep -n "switch only after agreement" .agents/coordinator.md` vazio.
  - Reportado por: claude

## Disposição do Coordenador (rodada 2, limite atingido)

- P7: o dono delegou o lote e está ausente. A correção só alinha R7 a §2 e §6, que as decisões delegadas pelo dono já fixam. O autor aplicou a sugestão exata do validador. Nenhuma terceira rodada de validação rodou. Sinalizado ao dono no relatório do lote para revisão posterior.

## Revalidação da emenda (2026-10-13, uma rodada)

Emenda: R10 permite que `internal/kit/status_drift.go` e `adapters/claude-code/plugin/hooks/status-drift.ts` (mais os testes) aceitem `delivered on`; R6 acrescenta `Solution` e o item de idioma do AGENTS.md; R1/R2/R6/R7 divididos para R11. Achada pela revisão de código (T-2001 rodada 1) e pelo QA.

- Status da revisão: `changes required` (1 bloqueante)
- **P7** [applied] [conflict] R7 agora diz editar sem segunda confirmação, verificado pelo `grep` de `switch only after agreement`.
- **P8** [pending] [conflict] §2 linha 30 ("nenhuma mudança em Go é necessária") e §6 linha 107 ("Nenhum código novo") contradizem o R10 emendado. A §6 não tem registro da decisão da emenda de status-drift.
  - Local: §2 linha 30, §6 linha 107
  - Evidência: R10 linhas 69-74; status_drift.go:82; status-drift.ts:35
  - Sugestão: restringir a linha 30 ao template, acrescentar um item na §2 para as duas formas e acrescentar uma linha de decisão do Coordenador na §6
  - Reportado por: claude

## Disposição do Coordenador (rodada da emenda)

- P8: dono ausente, lote delegado. O autor aplicou o texto exato do validador. Nenhuma outra rodada rodou. Sinalizado ao dono no relatório do lote.
