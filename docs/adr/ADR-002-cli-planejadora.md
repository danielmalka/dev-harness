# ADR-002 · CLI externa também como planejadora, pela mesma mecânica somente-leitura de `external-clis`

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-09-28 (emendado em 2026-09-29) |
| Decisor | Daniel Malka (dono) |
| Origem | PRD-005 (RF-04); PLAN-011 Decisions needed #1; emenda de 2026-09-29: PRD-008 (R14) |
| Reversibilidade | cara ou com muitos dependentes: `.commands/plan-loop.md`, `.skills/external-clis/SKILL.md` e `.agents/coordinator.md` passam a citar esta fronteira |

## 1. Contexto

- Pergunta: o binário de uma CLI externa, hoje despachado só para julgar conteúdo existente (revisão), pode também ser despachado para escrever conteúdo novo — um plano candidato numa wave do `/dh:plan-loop` — sem abrir mão da garantia de somente-leitura? Se a resposta for não, o loop vira Claude contra Claude, que o dono já rejeitou.
- Forças:
  - Fato: `.skills/external-clis/SKILL.md` limita o escopo a etapas de revisão ("This version covers review stages only; a CLI never writes production code through this skill"). Um plano não é código de produção, mas é conteúdo novo, não um veredito sobre conteúdo existente (PRD-005 §6 Contratos).
  - Fato: o PRD-005 RF-04 exige que a chamada da CLI planejadora reaproveite a mecânica já definida para revisão — resolução de binário, modo somente-leitura do binário quando existir, checagem de congelamento do workspace (`git status --porcelain -uall` + `git hash-object` por caminho, antes e depois) — e que a CLI devolva o plano como texto; quem grava em disco é o Coordenador (decisão do dono G15, `.harness/tasks/UND-003/GRILL.md`).
  - Fato: RISK-001 (`.harness/RISKS.md`) mostrou que o runtime não honra restrições declaradas só no frontmatter (`disallowedTools`); a proibição de delegar precisa estar no corpo do prompt, verbatim. A cláusula de revisão existente fala em "this review", e o RF-04 pede uma redação ajustada para planejamento.
  - Fato: sem CLI de outro fornecedor não há loop, porque duas gerações da mesma IA não valem como comparação (decisão do dono G6). Pelo mesmo motivo, `cli:claude/<slug>` fica fora do grupo de CLIs planejadoras enquanto o kit for plugin do Claude Code (decisão do dono AA1, 25/09/2026, `CLAUDE.md`).

## 2. Opções consideradas

| Dimensão | A · estender `external-clis` com um tipo de chamada de planejamento limitado | B · manter `external-clis` só de revisão e criar mecanismo separado para a CLI planejadora | C · nenhuma CLI planeja; RF-04 cai e o planejador é só Claude |
|---|---|---|---|
| Complexidade | uma seção nova na skill e uma no Coordenador; mesma resolução de binário, mesmo congelamento | segunda skill ou segundo fluxo com resolução, congelamento e cláusulas duplicados | nenhuma nova |
| Custo de construir | prosa: instrução de não escrever, cláusula anti-delegação de planejamento, teste de falha de transporte próprio | reescrever e manter em paralelo tudo que a revisão já provou | zero, mas o loop perde o propósito |
| Custo de operar | duas cláusulas anti-delegação para manter em sincronia à mão | dois mecanismos que divergem com o tempo | nenhum |
| Impacto em dados | nenhum: a CLI nunca grava; o Coordenador é o único gravador de `wave-N/candidate-*.md` | igual a A, se o novo mecanismo repetir a regra | nenhum |
| Impacto em segurança | a mesma garantia já governada por RISK-001: cláusula no corpo do prompt, checagem de congelamento descarta qualquer alteração | superfície nova, sem o histórico de evals da revisão | nenhum |

## 3. Decisão

- Escolhida: A.
- Motivo decisivo: reaproveitar a mecânica. Resolução de binário, checagem de congelamento e a disciplina anti-delegação no corpo do prompt já estão provadas em chamadas de revisão (casos `external-clis-*`, RISK-001); a chamada de planejamento muda só o prompt e o teste de falha de transporte.
- Descartadas: B, porque duplica a resolução de binário, o congelamento e a cláusula sem ganho de isolamento — a garantia real já é a checagem pós-chamada, não o nome da skill; C, porque reduz o loop a Claude contra Claude, que o dono rejeitou (G6), e derruba o RF-04.
- Preferência declarada: nenhuma.

Forma decidida (o que os construtores recebem, sem menu):

- A CLI planejadora nunca escreve no disco (qualificado pela emenda de 2026-09-29 abaixo). O plano volta como texto da resposta; o Coordenador é o único gravador, como já é do veredito persistido de um revisor CLI.
- A mesma checagem de congelamento da revisão vale para a chamada de planejamento: qualquer mudança na árvore descarta o plano (qualificado pela emenda de 2026-09-29 abaixo), reverte a mudança e marca a CLI como not-run por falha de transporte (PRD-005 AC-05).
- O prompt de planejamento carrega, verbatim no corpo, três blocos: a instrução de não escrever nenhum arquivo, uma cláusula anti-delegação própria de planejamento (separada da de revisão, sob cabeçalho próprio) e o bloco "Reading is allowed" que `external-clis` já publica. Esse último é copiado verbatim no prompt (PRD-005 AC-05); na skill, a seção de planejamento só o referencia, sem reescrevê-lo.
- Falha de transporte de uma chamada de planejamento: exit não-zero, timeout, divergência no congelamento, ou resposta sem conteúdo de plano reconhecível no formato de `implementation-planning` (objetivo e ao menos uma fatia). Não há linha de veredito; a tabela de sinais de veredito da revisão não se aplica.
- Só o `/dh:plan-loop` despacha uma CLI planejadora, uma por wave (RF-03), da rotação sobre `reviewers.document` sem `claude` e, enquanto o kit for plugin do Claude Code, sem `cli:claude/<slug>` (RF-05, G19, AA1). A CLI planejadora conta como um especialista concorrente no teto do Coordenador, como o revisor CLI (RF-13, AC-15).
- **Emenda de 2026-09-28 (decisão do dono na revisão da T-801). SUBSTITUÍDA pela emenda de 2026-09-29, logo abaixo; o texto fica como histórico e não vale mais.** só planeja uma CLI cujo binário tem modo somente-leitura documentado na tabela "Read-only mode by binary" de `external-clis` (hoje `codex -s read-only` e `mcode --permission off`). Um binário sem esse modo (`grok`, `agy`, `opencode`) não entra na rotação de planejadores. Motivo: da rodada 5 à 9 da revisão, a checagem de congelamento numa árvore viva vazou por um caminho novo a cada rodada (arquivos ignorados, reversão via git, symlink, `.gitignore` que se esconde). O congelamento detecta, mas não isola. A allowlist dos comandos de `Checks:` só governa o que o plano pode listar para um builder autorizado depois; ela não dá à CLI nenhum direito de execução durante a chamada.
- **Emenda de 2026-09-29 (decisão do dono, PRD-008; substitui a de 2026-09-28).** Enunciado obrigatório, que acompanha qualquer frase sobre este mecanismo: a chamada roda num clone descartável; escritas fora dele são detectadas na árvore viva (R5), nunca impedidas; escritas fora do projeto não são detectadas. O clone é um diretório de trabalho, não uma caixa de areia: a CLI roda como o mesmo usuário do SO e alcança a árvore viva e o resto da máquina por caminho absoluto.
  - Quem planeja: um binário qualifica-se como planejador se tem modo somente-leitura documentado e medido com sonda datada (hoje `codex -s read-only`, que continua no modo próprio e fora do clone) OU roda no clone descartável do PRD-008 (`grok`, `opencode` e `mcode`, este até ter sonda datada). O `agy` fica fora até a sonda de cwd (H3) ser registrada, datada e confirmar que ele opera no cwd do processo (se a sonda falhar, ele continua fora); até lá ele não planeja e sua revisão segue na árvore viva sob o risco de revisão aceito no RISK-002. `cli:claude/<slug>` continua excluído (AA1). A justificativa antiga da emenda de 2026-09-28 (o congelamento só detecta; o planejador tem de ser incapaz de escrever) sai. O que a substitui é o enunciado obrigatório acima.
  - Qualificação do §3, primeira frase ("A CLI planejadora nunca escreve no disco"): passa a "nunca escreve a árvore viva nem o arquivo do plano; pode escrever dentro do clone". Escrita na árvore viva é detectada (R5), nunca impedida. O Coordenador continua o único que grava o plano (`wave-N/candidate-*.md`, `PLAN.md`) e nada do clone é copiado para a árvore viva.
  - Qualificação do §3, segunda frase ("qualquer mudança na árvore descarta o plano"): passa a duas checagens. Mudança na árvore viva detectada por R5 (caminho criado, apagado ou alterado; ref criada ou movida) descarta o veredito ou o plano, reverte o que for revertível e marca a chamada como not-run por "mudança detectada". Mudança no conjunto de R6 dentro do clone (caminhos rastreados e não rastreados-não ignorados do estado julgado) descarta o veredito, sem reversão, porque o clone é descartado. Um arquivo novo não rastreado no clone não descarta: vira nota no relatório da chamada.
  - Não detectado, e aceito pelo dono em 2026-09-29 até o lote de sandbox do SO: (a) escrita fora do diretório do projeto; (b) processos que sobrevivem à chamada, rede e `git push` para remoto que não seja o repositório vivo; (e) leitura e devolução de segredos; e os caminhos do `.git` vivo que R5 não nomeia (`logs/`, `objects/` soltos e pacotes, `worktrees/`, `refs/` como arquivos (só o que `git for-each-ref` enxerga), `modules/`, `FETCH_HEAD`, `ORIG_HEAD`, `COMMIT_EDITMSG`, bytes crus do `.git/index`; a lista é a que o PRD-008 conhece, sem pretensão de ser exaustiva). Também não detectados: escrita na árvore viva desfeita antes da segunda fotografia ou que preserve os bytes, e escrita concorrente de outro processo, o Coordenador incluso.
  - Fato de origem da troca: um `git worktree` compartilha o `.git` inteiro com o repositório vivo (PRD-008 F2), por isso o mecanismo é um clone independente com `.git` próprio. Isso elimina, para o que a CLI escreve por caminho relativo, os vetores do `.git` compartilhado; não elimina o acesso por caminho absoluto.

## 4. Consequências

- Fica mais fácil: acrescentar um segundo tipo de chamada limitada no futuro, se um dia for preciso, sobre a mesma base; auditar a CLI planejadora com o mesmo registro de congelamento (`clean`/`reverted`) e as mesmas classes de not-run (exit não-zero, timeout, divergência no congelamento) já usados para o revisor CLI.
- Fica mais difícil: manter duas cláusulas anti-delegação (revisão e planejamento) em sincronia à mão. A checagem automática `checkAntiDelegationClause` (`internal/kit/evals.go`) compara só a cláusula de revisão entre `.agents/coordinator.md` ("## External CLI reviewers") e `.skills/external-clis/SKILL.md` ("## The anti-delegation clause"). A cláusula de planejamento é publicada uma única vez, em "The planning call" de `.skills/external-clis/SKILL.md`; `.agents/coordinator.md` só a referencia, sem cópia. A sincronia manual é entre as duas cláusulas, revisão × planejamento, e depende de diff manual.
- Dívida assumida: a sincronia revisão × planejamento é disciplina manual até alguém estender a checagem para cobrir a cláusula de planejamento, ou até as duas divergirem na prática (ver §6).
- Contratos afetados: `.skills/external-clis/SKILL.md` ganha "The planning call" e deixa de dizer "review stages only"; `.agents/coordinator.md` ganha "External CLI planners" (incluindo a exceção de julgamento de wave, decisão do dono de 28/09/2026: cada wave é julgada pela lista completa de `reviewers.document`); `.skills/implementation-planning/SKILL.md` passa a ser o formato exigido da resposta da CLI; `.commands/plan-loop.md` despacha contra esse contrato. Sem impacto em dados nem em segurança além do que RISK-001 já governa.

## 5. Validação

- Check que prova que funciona: `go run ./cmd/dh validate --source-only .` sem erros, em especial sem `anti-delegation clause differs between .agents/coordinator.md and .skills/external-clis/SKILL.md` (a cláusula de revisão continua intacta); diff manual dos blocos cercados que a checagem compara; leitura direta confirmando que a cláusula de planejamento está sob cabeçalho próprio, fora e depois do primeiro bloco cercado de "## The anti-delegation clause" e de "## External CLI reviewers" — senão `checkAntiDelegationClause` passa a comparar o bloco errado. Em operação: numa wave real, o `wave-N/candidate-*.md` da CLI é gravado pelo Coordenador e o registro de congelamento sai `clean`.
- Sinal de que está falhando: a validação acusa a cláusula divergente; uma CLI planejadora altera a árvore (congelamento `reverted`); a cláusula de planejamento difere do texto fixado em T-801 / PLAN-011 fatia 1; comparadas entre si, as duas cláusulas diferem em algo além de "this review" → "this planning task" (ignorando a quebra de linha); waves de candidato único recorrentes por "no recognizable plan content".

## 6. Revisar quando

- Surgir proposta de um terceiro tipo de chamada para CLI externa; ou as cláusulas de revisão e de planejamento divergirem na prática; ou o kit passar a rodar num host que não seja o Claude Code (AA1 devolve `cli:claude` à rotação). Ou existir caixa de areia do sistema operacional (bubblewrap, só Linux, lote posterior opcional), que fecharia os residuais aceitos em 2026-09-29 (RISK-002): aí a emenda de 2026-09-29 pode ser revista. Ou a sonda H3 do `agy` ser registrada, ou a revisão de 21/10/2026 dos passos em prosa da skill `external-clis`.

## 7. Referências

- `docs/prd/PRD-005-planejamento-em-loop.md` (RF-02, RF-03, RF-04, RF-05, RF-13, AC-05, AC-15, §6 Contratos e Segurança)
- `.harness/tasks/PRD-005/PLAN.md` (PLAN-011: Decisions needed #1, fatia 1)
- `.harness/tasks/UND-003/GRILL.md` (G6, G15, G19)
- `.harness/RISKS.md` RISK-001, RISK-002
- `.skills/external-clis/SKILL.md` ("The anti-delegation clause", "Reading is allowed", passo 7 do Procedure)
- `.agents/coordinator.md` ("External CLI reviewers")
- `internal/kit/evals.go` (`checkAntiDelegationClause`)
- `CLAUDE.md` do repositório (AA1; exceção de julgamento de wave de 28/09/2026)
- `docs/prd/PRD-008-isolamento-por-clone.md` (R1 a R15, §2 Residual, Decisões do dono; base da emenda de 2026-09-29)
- `AGENTS.md` do repositório (decisões do dono de 29/09/2026: clone em vez de `git worktree`; residuais aceitos)
- `docs/prd/PRD-003-revisao-em-pares-por-cli.md` (origem das CLIs como revisoras somente-leitura)
