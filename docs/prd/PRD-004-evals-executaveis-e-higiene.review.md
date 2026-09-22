# Review: PRD-004

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-004-evals-executaveis-e-higiene.md |
| Fonte | pedido do dono em 2026-09-21 (itens 2, 3, 4, 5, 6, 10, 13, 14), decisões fechadas do Coordenador em 2026-09-21/22, fatos do runner nativo (`claude plugin eval --help` 2.1.278 e documentação oficial) |
| Rodada | 2 de 2 |
| Veredito | changes required (rodada 2, P16–P17) → fechados por disposição do Coordenador por delegação do dono; não revalidados (teto de 2 rodadas) |
| Revisor | claude (document-validator, opus) |
| Estado avaliado | rodada 2: sha256 6dcb78a22cbf97bf (rodada 1: 779bebd39f3ee304); pós-disposição: 424ac66893b393ba |
| Atualizado | 2026-09-22T03:55-03:00 |

## Achados

- **P1** [applied] [lacuna] A decisão (c) sobre como dirigir os casos comportamentais do Coordenador (prompt com `/dev-harness:<cmd>` mais `append_system_prompt` com o papel) não virou requisito; só o fallback not-run está escrito.
  - Local: seção 4, RF-02 (linha 37); seção 8, terceiro risco (linha 114)
  - Evidência: RF-02 lista os campos de frontmatter e omite `append_system_prompt`; nenhuma frase diz que um caso comportamental é dirigido pelo comando de barra. Dois builders divergem: um tenta dirigir; outro marca tudo not-run sem tentar.
  - Sugestão: RF-02: "Um caso comportamental do Coordenador é dirigido pelo `prompt.md` com o comando `/dev-harness:<cmd> <argumentos>` como prompt e com `append_system_prompt` contendo o corpo de `.agents/coordinator.md` (sem frontmatter); `allowed_tools` lista `Read, Glob, Grep, Skill, Agent`. Só quando essa forma não produzir o fluxo o caso é declarado not-run com o motivo, nunca antes de tentar." Incluir `append_system_prompt` na lista de campos.
  - Reported by: claude
- **P2** [applied] [critério fraco] AC-03 e AC-13 exigem "not-run com motivo" no JSON do runner, mas o formato nativo não tem esse estado: caso presente em `evals/cases/` sempre roda e pontua; caso ausente não aparece; não há flag de exclusão.
  - Local: seção 5, AC-03 (linha 66) e AC-13 (linha 76); seção 4, RF-02 (linha 37) e RF-12 (linha 47)
  - Evidência: JSON do runner tem `cases[].aggregates.score`, `arms.with[].error`, `aborted`, `skippedPaidGraders`; nenhum campo de não executado. O critério só se cumpriria editando a saída à mão.
  - Sugestão: RF-02: "Um caso que não pode ser dirigido fica fora de `evals/cases/` (em `evals/not-run/<id>/` com o mesmo layout e um `REASON.md` de uma linha)." AC-03: "o JSON contém em `cases[]` cada diretório de `evals/cases/`; `evals/baselines/<data>-0.5.0.md` lista cada id de `evals/not-run/` com o motivo; a soma é 10 (+ os 4 novos)." AC-13 no mesmo molde.
  - Reported by: claude
- **P3** [applied] [conflito] Sinais 1 e 5 prometem "10/10 rodam" e risk-001 "rodando e passando", enquanto RF-02/RF-12 aceitam not-run; a regra da decisão (e) (RISK-001 só muda de estado se o caso rodar e passar) não aparece.
  - Local: seção 2, Sinais 1 e 5 (linhas 18 e 22); seção 4, RF-02 (linha 37) e RF-12 (linha 47)
  - Sugestão: Sinal 1: "10/10 migrados; N executados sob o runner, os demais listados not-run com motivo". Sinal 5 e RF-12: "RISK-001 passa a `resolvido` em `.harness/RISKS.md` apenas se `external-clis-risk-001` rodar e passar no baseline; failed ou not-run mantêm `mitigado`". AC-13 acompanha.
  - Reported by: claude
- **P4** [applied] [conflito] A decisão (d) (job de CI de evals só com o secret, threshold 0.8) é vinculante, mas o PRD a cita como "(se houver)", sem RF/AC, e RF-16/AC-17 congelam `ci.yml`; "skip explícito" não está definido.
  - Local: seção 6, Operação (linha 88); seção 4, RF-16 (linha 51); seção 5, AC-17 (linha 80)
  - Sugestão: RF-18: workflow separado `.github/workflows/evals.yml` (ci.yml intocado) rodando `claude plugin eval dist/claude-code/dev-harness --trust-plugin --runs 1 --ablation none --no-publish --max-cost-usd 15 --threshold 0.8 --json` quando o secret existe; sem secret, job verde com a linha `skipped: ANTHROPIC_API_KEY absent`. AC correspondente.
  - Reported by: claude
- **P5** [applied] [conflito] O runner só aceita em `context.add_dirs` diretórios dentro do diretório do caso, mas o PRD mantém as fixtures compartilhadas `evals/fixtures/external-clis-review` e `evals/fixtures/prd-review` sem dizer para onde vão.
  - Local: seção 4, RF-02 (linha 37) e RF-06 (linha 41); seção 6, item 3 (linha 94)
  - Evidência: documentação oficial: "`context.add_dirs` — Directories inside the case directory"; cada execução começa em workspace vazio.
  - Sugestão: RF-02: "As fixtures de cada caso vivem em `evals/cases/<id>/fixtures/` (cópia por caso, duplicação aceita; `context.add_dirs: [fixtures]`); `evals/fixtures/` fica só com `slice-01` (passo Fixture de `ci.yml`)." RF-06: "falha também quando `add_dirs`/`scaffold_script`/`history_file` apontam para fora do diretório do caso." AC-02 ganha "cada `add_dirs` resolve dentro do próprio caso".
  - Reported by: claude
- **P6** [applied] [conflito] AC-04 julga um `<documento>.review.md` gravado em disco, mas o comando de RF-03 não concede `Write`/`Edit`; sob esse comando o arquivo não pode existir. AC-05 não diz se é arquivo ou resposta.
  - Local: seção 5, AC-04 (linha 67) e AC-05 (linha 68); seção 4, RF-03 (linha 38)
  - Sugestão (recomendada, sem grant): "então a última mensagem do run contém `## Verdict` / `changes required` e a linha `Reported by: cli:codex/<slug>`, graduada por `target: last_message`".
  - Reported by: claude
- **P7** [applied] [ambiguidade] RF-08 oferece dois mecanismos de exemplo com comportamento diferente; `git diff --stat` não detecta reescrita de linha já modificada (falha AC-09); `git write-tree` exige staging e altera o índice.
  - Local: seção 4, RF-08 (linha 43); seção 5, AC-09 (linha 72)
  - Sugestão: fixar: "Antes da chamada, registrar `git status --porcelain` e, para cada caminho listado (modificado ou não rastreado), `git hash-object <caminho>`. Depois, recalcular e comparar lista e hashes; qualquer diferença é falha de transporte." Remover o "por exemplo".
  - Reported by: claude
- **P8** [applied] [ambiguidade] RF-10/AC-11 validam "o plano", mas `.commands/plan.md` persiste só `.harness/tasks/<id>/TASK.md`; não há caminho para o plano nem para o relatório; N TASK.md ou um documento?
  - Local: seção 4, RF-10 (linha 45); seção 5, AC-11 (linha 74)
  - Evidência: precedente em disco `.harness/tasks/PRD-003/PLAN.md` + `PLAN.review.md`, improvisado fora do texto do comando.
  - Sugestão: RF-10: "O plano é persistido como `.harness/tasks/<PRD-id>/PLAN.md` (TASK.md por fatia continua); `document-validator` julga esse documento contra o PRD e o brief, com os TASK.md como anexo; relatório em `.harness/tasks/<PRD-id>/PLAN.review.md`; em `changes required` o Coordenador redespacha `implementation-planner` só com os pontos bloqueantes; cap de duas rodadas." AC-11 nomeia os caminhos.
  - Reported by: claude
- **P9** [applied] [ambiguidade] RF-11 diz que "a rotina distingue" documento-contrato de tutorial sem dizer a regra.
  - Local: seção 4, RF-11 (linha 46); seção 5, AC-12 (linha 75)
  - Sugestão: "O ciclo roda quando o documento é persistido a partir de `templates/<lang>/PRD.md`, `STORY.md`, `TASK.md` ou `ADR.md`; qualquer outro produto de `/dev-harness:document` não aciona o ciclo."
  - Reported by: claude
- **P10** [applied] [ambiguidade] Com o pacote como alvo, o runner grava `results/<timestamp>/` dentro de `dist/.../evals/`; sem `.gitignore`, sem `--output-dir`; e fica "possivelmente" se `dh validate` compara `evals/` fonte vs pacote (se comparar, `results/` reprova como extra).
  - Local: seção 4, RF-01 (linha 36) e RF-03 (linha 38); seção 6, item 2 (linha 93)
  - Sugestão: RF-03 acrescenta `--output-dir evals/results/<data>` (raiz do repositório); RF-01 acrescenta "`.gitignore` ganha `evals/results/`"; seção 6 decide: `evals` entra nos pares de `checkCoverage` com `results/` e `baselines/` ignorados.
  - Reported by: claude
- **P11** [applied] [ambiguidade] O JSON do baseline carrega a configuração da suíte (caminhos); `iterTextFiles` varre `.json` sob `evals/`, então um baseline com home expandido reprova `dh validate .`.
  - Local: seção 4, RF-03 (linha 38) e RF-06 (linha 41); seção 5, AC-03 (linha 66)
  - Sugestão: AC-03 ganha "e `go run ./cmd/dh validate .` continua passando com o baseline commitado"; seção 6: se o JSON trouxer caminhos absolutos, `checkPrivatePaths` ignora `evals/baselines/` e `evals/results/`; nunca editar a saída do runner.
  - Reported by: claude
- **P12** [applied] [ambiguidade] O comando de RF-03/AC-03 não traz `--trust-plugin`; sob `--json` com stdin não interativo o runner recusa com exit 1.
  - Local: seção 4, RF-03 (linha 38); seção 5, AC-03 (linha 66)
  - Sugestão: incluir `--trust-plugin` no comando de RF-03 e AC-03.
  - Reported by: claude
- **P13** [applied] [organização] RF-15/AC-16 ancoram a remissão em uma "sequência equivalente de CI" que `.skills/delivery-readiness/SKILL.md` não tem.
  - Local: seção 4, RF-15 (linha 50); seção 5, AC-16 (linha 79)
  - Sugestão: "no item 9 do procedimento (pipeline em ambiente descartável)".
  - Reported by: claude
- **P14** [applied] [organização] A tabela de cenários diz "contra o pacote ou o checkout", mas a raiz do checkout não é alvo válido (só `marketplace.json`).
  - Local: seção 3, primeira linha da tabela (linha 28)
  - Sugestão: "contra o pacote reconstruído (`dist/claude-code/dev-harness`)".
  - Reported by: claude
- **P15** [applied] [organização] O cabeçalho mantém `Stories | ST-000`.
  - Local: cabeçalho (linha 8)
  - Sugestão: "a derivar".
  - Reported by: claude

- **P16** [applied] [conflito] `evals/not-run/<id>/` mantinha o layout de caso e ia para o pacote por RF-01; o glob `**` do runner o executaria mesmo assim.
  - Local (rodada 2): RF-01, RF-02, AC-01, AC-03, seção 6, Sinal 4
  - Evidência: `claude plugin eval --help`: descoberta recursiva de `prompt.md`/`case.yaml` sob o eval dir.
  - Resolução (Coordenador por delegação do dono, 2026-09-22T03:55-03:00): RF-01 exclui `not-run/` do pacote; AC-01 exige ausência de `evals/not-run/` no pacote; `checkCoverage` ignora `not-run/` dos dois lados; RF-02: arquivo `prompt.not-run.md` sem `graders/`, nunca reconhecido pelo glob; Sinal 4 alinhado.
  - Reported by: claude
- **P17** [applied] [ambiguidade] RF-10 introduzia `PLAN.md` sem dizer se nasce de template (regra: todo template entra em en + pt-br) nem o nome do diretório quando o plano não deriva de PRD.
  - Local (rodada 2): RF-10, AC-11, seção 6 item 10
  - Resolução (Coordenador por delegação, 2026-09-22T03:55-03:00): `PLAN.md` é a saída da skill `implementation-planning` persistida tal qual, sem template novo; diretório `.harness/tasks/<PRD-id>/` quando deriva de PRD, senão `.harness/tasks/PLAN-<n>/`; AC-11 espelha.
  - Reported by: claude
- **P18** [applied] [ambiguidade] RF-18 não dizia como a CLI chega ao runner de CI nem a exportação do secret.
  - Resolução (Coordenador, 2026-09-22T03:55-03:00): RF-18 instala a Claude Code CLI em versão fixada (>= 2.1.278) e exporta `ANTHROPIC_API_KEY` só para o passo de eval.
  - Reported by: claude
- **P19** [applied] [ambiguidade] Mecanismo de RF-08 falhava em diretório não rastreado e arquivo apagado.
  - Resolução (Coordenador, 2026-09-22T03:55-03:00): `git status --porcelain -uall`; caminho ausente registra `deleted` no lugar do hash.
  - Reported by: claude
- **P20** [applied] [organização] Decisão da seção 8 citava âncora inexistente em delivery-readiness.
  - Resolução (Coordenador, 2026-09-22T03:55-03:00): alinhada com RF-15 ("item 9 do procedimento").
  - Reported by: claude

## Não levantado (rodada 2)

- O texto corrigido de P1–P15 confere com os fatos do runner e com os arquivos do kit: `append_system_prompt` é campo documentado; `Read, Glob, Grep, Skill, Agent` estão no conjunto somente leitura; `--trust-plugin`, `--output-dir`, `--threshold`, `--json [path]` existem; `target: last_message` é observável.
- RF-18 mirar `dist/claude-code/dev-harness` em CI é válido (ci.yml já reprova dist fora de sincronia; AC-01 põe `evals/` lá).
- `--output-dir` na raiz + `.gitignore` + `checkCoverage` ignorando `results/`/`baselines/`/`not-run/` fecham P10; `checkPrivatePaths` ignorando `evals/baselines/` e `evals/results/` fecha P11.
- RF-02 fixtures por caso e RF-06/AC-07 espelham a regra do runner; `evals/fixtures/slice-01` permanece para `ci.yml`.
- IDs RF-01..18 e AC-01..19 únicos; cada RF tem AC; os oito itens do dono e as decisões seguem mapeados.
- Não verificado pelo validador: conteúdo real da saída `--json`; se `--output-dir` suprime `<plugin>/evals/results/`; se um diretório com `prompt.md` sem `graders/` é tratado como caso (a disposição de P16 usa `prompt.not-run.md`, que o glob não casa, tornando o ponto irrelevante).

## Não levantado (rodada 1)

- Os oito itens do dono têm requisito: 2 em RF-01/02/03/17; 3 em RF-04; 4 em RF-05; 5 em RF-06; 6 em RF-07/08/09; 10 em RF-10/11; 13 em RF-12/13; 14 em RF-14/15/16. Sem escopo sem rótulo; versão 0.5.0 e CHANGELOG decorrem de RF-01.
- Decisões (a), (b), (1), (2), (3), (4) do Coordenador estão em RF-02, RF-05/06, RF-04, RF-01, RF-14/15 e RF-17, repetidas na seção 8.
- RF-04: sinais conferem com `.skills/external-clis/SKILL.md` (linhas 104 e 106) e `.skills/security-review/SKILL.md` (75-76).
- RF-05: os dois blocos cercados existem (coordinator.md a partir da linha 100; SKILL.md linha 69), idênticos hoje; trim por linha basta para delimitar.
- RF-06: `linkPattern` só reconhece sintaxe Markdown; `privatePattern` já cobre `.md`/`.yaml` por `iterTextFiles`.
- RF-07: Output format (~156) não tem a linha de revisor CLI, embora a seção "External CLI reviewers" (~118) já a peça.
- RF-09: os dois roadmaps existem e usam Mermaid. RF-14/RF-16: leitura de `ci.yml` correta.
- `--ablation none` coerente com o runner; o conjunto somente leitura inclui `Skill` e `Agent`.
- IDs RF-01..17 e AC-01..18 únicos; cada RF tem AC; "Não entra" coerente com o pedido (uma branch, PR aberta, sem merge).

## Evidência

- Lidos: o PRD (sha conferido), internal/build/build.go 100-160, internal/kit/validate.go 20-40, 450-560, 630-760, `.github/workflows/ci.yml`, `.gitignore`, `.agents/coordinator.md` 90-165, `.skills/external-clis/SKILL.md`, `.skills/security-review/SKILL.md`, `.commands/{discover,plan,document}.md`, `.skills/{harness-evaluation,implementation-planning,delivery-readiness}/SKILL.md`, árvore `evals/` e três casos, `.harness/RISKS.md` RISK-001, layout `.harness/tasks`, `templates/pt-br/PRD.md`, PRD-003, `.claude-plugin/marketplace.json`, `claude plugin eval --help` e a página oficial de plugin evals.
- Não verificado: conteúdo real da saída `--json` (nenhuma execução); se agentes de plugin carregam dentro da sessão filha do eval (P1; o fallback cobre).
- Rodada 1: sem relatório anterior.

## Desvios registrados na execução (Coordenador, 2026-09-22)

Registrados aqui em vez de reabrir o PRD já aprovado com ressalvas; cada um é um fato do runner descoberto ao rodar RF-03/AC-03.

- RF-03/AC-03: o comando real leva `--scaffold` e `--model sonnet`. `--model sonnet` é controle de custo (decisão do Coordenador). `--scaffold` é necessário: `context.add_dirs` (RF-02, AC-07) não materializa `fixtures/` no cwd do workspace — só concede acesso no caminho original — e o primeiro baseline (`evals/results/2026-09-22/failed-no-scaffold.json`) bloqueou 13/13 casos por cwd vazio. Cada caso passou a declarar `context.scaffold_script: scaffold.sh`, que copia `fixtures/` para o cwd. AC-07 continua válido: `dh validate` resolve `add_dirs`, `scaffold_script` e `history_file`.
- RF-02: graders `regex` não aceitam flags inline `(?m)`/`(?mi)` (motor JavaScript); as flags foram movidas para a chave `flags:`. Seis graders corrigidos.
- `--json` recebe o caminho do arquivo (não é booleano) — já observado pelo validador em "Não verificado".
