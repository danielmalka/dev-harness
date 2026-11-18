# Revisão: PRD-021

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-021-lean-kit.md (+ tarefas T-2101 e T-2102 no harness dir) |
| Fonte | Decisões do dono 2026-10-13 e medições do pprof (despacho do Coordenador) |
| Rodada | 2 de 2 |
| Veredito | approved |
| Atualizado | 2026-10-13 |

## Achados

- **P1** [pending] [conflict] A T-2101 não consegue terminar verde: apagar `templates/pt-br/PRD.md` faz o teste de roadmap da árvore real (`TestRealTreeRoadmapsClean`) falhar em `docs/roadmap.html:685` e `docs/en/roadmap.html:685` ("14 templates"), e esses arquivos estão só no conjunto de escrita da T-2102.
  - Evidência: `internal/kit/doc_counts.go:88-126`, `internal/kit/roadmap_check_repo_test.go:24-34`
  - Sugestão: mova a correção da linha 685 (pt e en) para a autorização, o fluxo e o "Pronto quando" da T-2101
  - Reportado por: claude
- **P2** [pending] [conflict] `doc_counts.go` não guarda contagem de templates e o `dh build` não o regenera; a contagem vem do disco e de `harness-manifest.json`.
  - Local: R2, último item da §4, §6 (d), nota de autorização da T-2102
  - Sugestão: "`harness-manifest.json` lista 13 templates depois de `go run ./cmd/dh build`. `checkDocumentedCounts` lê a contagem do disco. `internal/kit/doc_counts.go` não muda."
  - Reportado por: claude
- **P3** [pending] [conflict] O fato "37 regexps em kit, dashboard, harness e snapshot" está errado: `internal/snapshot` não tem regexp; em nível de pacote são kit 20, dashboard 12, harness 2; `internal/doctor` tem 1; o kit tem mais 3 dentro de funções (`validate.go:445,459,482`).
  - Local: §1, §6 (a)
  - Reportado por: claude
- **P4** [pending] [gap] `docs/dev-harness-mapa.excalidraw:37` traz a afirmação de espelho e não está em nenhuma lista nem conjunto de escrita.
  - Sugestão: inclua-o na §4, na §5 da T-2102 e na autorização dela
  - Reportado por: claude
- **P5** [pending] [weak criterion] O `grep "pt-br/PRD"` da R2 não enxerga as afirmações reais de espelho (`"espelhadas"`, `"mirrored"`, `"cópia idêntica"`, `"identical copy"`, linha 678 do roadmap).
  - Sugestão: faça grep das frases de espelho; a linha 678 dos dois roadmaps diz 7 templates em `en/` e 6 em `pt-br/`
  - Reportado por: claude
- **P6** [pending] [ambiguity] A R4 não diz se o cache mantém as travas de leitura de `snapshot.ReadCapped` (sem symlink, arquivo regular, limite de tamanho, o mesmo aviso `skipped`).
  - Sugestão: a consulta roda depois das mesmas checagens de `os.Lstat`, de arquivo regular e de tamanho (como fazem `scanEpochal`/`epochalCache`); um arquivo recusado dá o mesmo aviso com e sem o cache; o teste em Go cobre um ticket trocado por um symlink
  - Reportado por: claude
- **P7** [pending] [weak criterion] A checagem da R5 (nomes de arquivo `*_bench*`) não enxerga um benchmark de Go.
  - Sugestão: `git diff main...HEAD -- '*.go' | grep -c '^+func Benchmark'` é 0; `git diff --name-only main...HEAD -- evals/` é vazio
  - Reportado por: claude
- **P8** [pending] [excess] (`non-blocking`) Os itens (a)-(e) da R6 e "`dist/` guarda só a versão atual" não têm rastro na §6; o item 2 da R7 não tem checagem.
  - Reportado por: claude
- **P9** [pending] [ambiguity] (`non-blocking`) O cache do dashboard guarda entradas de caminhos apagados ou renomeados.
  - Sugestão: descarte as entradas dos caminhos não vistos no último build
  - Reportado por: claude

## Disposição do Coordenador (rodada 1)

- Conferidos por amostragem P1 (`roadmap_check_repo_test.go`, a linha 685 do roadmap diz "14 templates") e P3 (`internal/snapshot` não tem regexp).
- Os nove achados voltaram ao autor com as sugestões do validador aceitas. P8: a frase "`dist/` guarda só a versão atual" cai; o conteúdo da R6 fica registrado na §6 como decisão do Coordenador sob a delegação do dono.

## Rodada 2 (2026-10-13)

- Status da revisão: `approved` (0 `blocking`). P1-P9 aplicados. Três achados novos, não bloqueantes.
- P1 [applied] A T-2101 é dona das linhas de contagem de templates do roadmap (cerca de 678 e 685); a T-2102 só acrescenta depois dela o card 0.28.0.
- P2 [applied] R2, §4 e §6 (d): a contagem vem do disco e do manifesto; `doc_counts.go` não muda.
- P3 [applied] §1 e §6 (a): 34 regexps em nível de pacote (kit 20, dashboard 12, harness 2), doctor 1, snapshot nenhuma, 3 chamadas dentro de funções no kit.
- P4 [applied] `docs/dev-harness-mapa.excalidraw` está na R2, na §4 e na T-2102.
- P5 [applied] A R2 faz grep das frases de espelho.
- P6 [applied] A R4 mantém as travas de leitura de `ReadCapped`; o teste em Go cobre um symlink.
- P7 [applied] A R5 checa `^+func Benchmark` e `evals/`.
- P8 [applied] O conteúdo da R6 fica registrado na §6 como decisão do Coordenador; a frase de `dist/` saiu.
- P9 [applied] O cache descarta as entradas dos caminhos não vistos no último build.
- **P10** [aplicado pelo despacho do Coordenador] A entrada de `AGENTS.md` na §4 diz "item dos binários e do pprof" (não um item de `dist/`). Não bloqueante.
- **P11** [aplicado pelo despacho do Coordenador] A §4 da T-2101 lista a parte da R2 sobre a contagem do roadmap; a §5 da T-2102 entrega à T-2102 só o card 0.28.0. Não bloqueante.
- **P12** [aplicado pelo despacho do Coordenador] O item 2 da R7 foi dividido em frases curtas. Não bloqueante.

## Disposição do Coordenador (rodada 2)

- Correções só de texto para P10-P12 enviadas ao autor; sem terceira rodada de validação (limite atingido, os achados são não bloqueantes e mecânicos).
