# Review: PRD-003

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-003-revisao-em-pares-por-cli.md |
| Fonte | pedido do dono de 2026-09-21 (verbatim no despacho DISC-003) + brief de discovery do Coordenador + skill pessoal `~/.claude/skills/fleet-clis/` + restrições do kit (`profiles/README.md`, `.agents/`, `.skills/`, `.commands/`, `CLAUDE.md`, `templates/pt-br/PRD.md`) |
| Rodada | 2 de 2 |
| Veredito | approved |
| Revisor | claude (document-validator, opus) |
| Estado avaliado | sha256 b125a95417adfa0970dd92339cd77024e176814d1c2e25dd6a6675ba6d7f52c2 (rodada 1: 9062b376c1502e584de4a24e109e1d5c42da3ef1581254a692bf6d535b31b25c) |
| Atualizado | 2026-09-21T17:05-03:00 — P10–P12 e as quatro decisões pendentes fechados pelo dono e aplicados pelo Coordenador (edição cirúrgica; não revalidada pelo validador, teto de rodadas atingido) |

## Achados

- **P1** [applied] [lacuna] Nada obrigava a invocação de uma CLI revisora a ser somente-leitura nem definia o que acontece se ela editar o workspace.
  - Local (rodada 1): seção 4, RF-01 e "Não entra"; seção 6 Segurança/Operação
  - Resolução: RF-10 (linha 43) e AC-12 (linha 65) acrescentados como sugerido; seção 6 Segurança (linha 72) repete a regra, rodada 2.
- **P2** [applied] [conflito] A entrada de `.harness/local.yaml` no `.gitignore` era atribuída ao `setup`, que não escreve fora de `.harness/`.
  - Local (rodada 1): seção 4, RF-02; seção 6 Dados; AC-10
  - Resolução: RF-02 (linha 35) passa a rodar `git check-ignore` e reportar ao dono em vez de editar o `.gitignore`, cita o write set correto do `setup` e declara a premissa sobre `.harness/` versionado; AC-10 (linha 63) perdeu a cláusula de `.gitignore`; seção 6 Dados (linha 70) coerente, rodada 2. Ver P12 para um ponto residual não bloqueante.
- **P3** [applied] [ambiguidade] O "bloco de veredito" e o veredito mesclado não estavam definidos para `security`, e `VERDICT` era o marcador da fleet-clis.
  - Local (rodada 1): seção 4, RF-05 e RF-08; AC-05 e AC-06; seção 6 Contratos
  - Resolução: RF-08 (linha 41) define o sinal observável por etapa exatamente como sugerido; AC-05 (linha 58) e AC-06 (linha 59, com `blocking` para `security`) o usam; seção 6 Contratos (linha 71) afirma que nenhuma seção nova é criada e que o sinal já existe em cada formato, rodada 2.
- **P4** [applied] [ambiguidade] Para `code` e `security`, o "relatório persistente" de RF-06 era o corpo de saída, não persistido.
  - Local (rodada 1): seção 4, RF-06; AC-09
  - Resolução: RF-06 (linha 39) nomeia `<document>.review.md` para `document` e as entradas em `.harness/MEMORY.md`/`RISKS.md` para `code`/`security`; AC-09 (linha 62) nomeia esses arquivos, rodada 2.
- **P5** [applied] [conflito] A decisão `reviewers.yaml` vs `project.yaml` era justificada com um fato falso sobre o `setup` e fechada sem o dono.
  - Local (rodada 1): seção 8, primeira "Decisão pendente"; RF-03
  - Resolução: seção 8 (linha 93) reabre a decisão com `responde: o dono · bloqueia: RF-03`, corrige o fato e marca a preferência pelo arquivo irmão como hipótese; RF-03 (linha 36) diz que o local está pendente, rodada 2.
- **P6** [applied] [ambiguidade] "Sempre paralelo" não dizia o que acontece acima do teto de concorrência ou com CLI que não aceita execuções simultâneas.
  - Local (rodada 1): seção 4, RF-03; seção 8
  - Resolução: RF-03 (linha 36) ganhou a cláusula "em levas"; nova linha de Risco (linha 90) cita gotchas 2026-08-31 e 2026-09-08, rodada 2.
- **P7** [applied] [ambiguidade] "Timeout configurado" não tinha onde ser configurado nem padrão, e nenhum AC cobria o timeout.
  - Local (rodada 1): seção 4, RF-08; seção 5
  - Resolução: RF-08 (linha 41) define campo opcional por revisor em `.harness/reviewers.yaml`, padrão 15 minutos; AC-05 cobre o caminho do timeout, rodada 2.
- **P8** [applied] [critério fraco] A cláusula de deduplicação de RF-05 não tinha AC.
  - Local (rodada 1): seção 4, RF-05; seção 5
  - Resolução: AC-13 (linha 66) acrescentado, rodada 2.
- **P9** [applied] [organização] A seção 6 Contratos não listava as regras de "não rebaixar do Opus" e "não introduzir provedor pago".
  - Local (rodada 1): seção 6 Contratos; seção 3, linha 2
  - Resolução: seção 6 Contratos (linha 71) lista a emenda a `.commands/secure.md` e `.agents/coordinator.md`, rodada 2.
- **P10** [applied] [organização] A contagem e o intervalo de datas de `gotchas.md` estão desatualizados.
  - Local: seção 1 (linha 14); seção 8, primeiro risco (linha 88)
  - Evidência: "37 entradas de 2026-08-18 a 2026-09-17" — o arquivo tem hoje 45 entradas, a última de 2026-09-20 (linhas 43–50, incluindo a delegação codex→opencode de 2026-09-17 e dois incidentes de brief somente-leitura do opencode em 2026-09-20). O número veio da linha de "Não levantado" da rodada 1 do próprio validador, baseada numa leitura truncada do arquivo; o intervalo original do PRD (até 2026-09-20) estava certo. Cosmético: não muda requisito.
  - Sugestão: "45 entradas de 2026-08-18 a 2026-09-20".
  - Resolução: contagem corrigida para "45 entradas de 2026-08-18 a 2026-09-20" na seção 1 e no primeiro risco (Coordenador, 2026-09-21T17:05-03:00).
- **P11** [applied] [ambiguidade] Não está dito se "somente-leitura" (RF-10) e a proibição no prompt (RF-01) também cobrem executar scripts e comandos do repositório.
  - Local: seção 4, RF-01 (linha 34); RF-10 (linha 43)
  - Evidência: gotchas.md 2026-09-20 (linha 49): um `opencode run` com brief READ-ONLY executou `bash scripts/build-claude-code.sh --help` por conta própria e reconstruiu `dist/`; a regra registrada é "brief read-only deve proibir explicitamente rodar scripts do repositório, não só editar". RF-01 proíbe invocar outra CLI; RF-10 proíbe editar e reverte. Não bloqueia: onde o binário tem modo somente-leitura técnico nada persiste, e onde não tem, o `git status` e a reversão de RF-10 limitam o dano; só efeitos fora da árvore (rede, caches) escapam.
  - Sugestão: estender a proibição de RF-01 para "…invoque outro binário de CLI ou execute scripts e comandos de build/teste do repositório por conta própria (gotchas.md 2026-09-17 e 2026-09-20)".
  - Resolução: decisão do dono — revisor CLI não executa scripts do repositório; pode rodar apenas `commands.test`/`commands.lint` de `.harness/project.yaml` (ex.: `make test`). RF-01 estendido e AC-15 criado (Coordenador, 2026-09-21T17:05-03:00).
- **P12** [applied] [ambiguidade / critério fraco] RF-02 não diz se `.harness/local.yaml` é gravado mesmo quando o caminho não está ignorado, e o comportamento de `check-ignore` não tem AC.
  - Local: seção 4, RF-02 (linha 35); seção 6 Dados (linha 70); seção 5
  - Evidência: "se o caminho não estiver ignorado, reporta ao dono em vez de editar o `.gitignore` sozinho" contrasta reportar com editar o `.gitignore`, mas não diz se a gravação de `local.yaml` acontece. Leitura A: grava e avisa (caminhos de máquina podem ser commitados depois pelo dono; sem segredos, e caminho errado noutra máquina vira not-run por RF-07). Leitura B: não grava até o dono ignorar, deixando os revisores CLI not-run enquanto isso. As duas degradam explicitamente, por isso não bloqueia. AC-10 cobre só o grep na árvore do kit.
  - Sugestão: completar RF-02 com "…e não grava `.harness/local.yaml` até o dono ignorar o caminho; até lá os revisores CLI ficam not-run com o motivo" (ou "…grava o arquivo mesmo assim e registra o aviso em `.harness/MEMORY.md`") e acrescentar "AC-14 (RF-02) Quando `.harness/local.yaml` não estiver ignorado pelo git do projeto consumidor, então o Coordenador reporta isso ao dono antes da primeira gravação e não edita o `.gitignore`".

  - Resolução: decisão do dono — `.harness/local.yaml` é gravado esteja `.harness/` ignorado ou não; o Coordenador avisa uma vez e registra em MEMORY.md, sem editar `.gitignore`. RF-02 e seção 6 Dados reescritos; AC-14 criado (Coordenador, 2026-09-21T17:05-03:00).
## Correção pós-aprovação

- 2026-09-21T20:05-03:00 (Coordenador, a partir do achado Major da revisão de T-001): RF-08, AC-05 e a seção 6 diziam "linha literal" para o sinal de `security`; o texto real em `.skills/security-review/SKILL.md:88` é mais longo. O PRD passa a definir o sinal como prefixo de linha em todas as etapas. Sem mudança de requisito; evita que T-003 implemente igualdade exata.

## Não levantado

- A decisão reaberta sobre o local do arquivo (seção 8, linha 93) é referenciada de forma consistente: RF-03 sinaliza o local como pendente e as demais menções a `.harness/reviewers.yaml` (RF-07, RF-08, AC-01, seção 6, riscos) leem-se como nome de trabalho, não como escolha fechada.
- A retirada de "CLI que delega a outra" da lista de falhas de transporte em RF-01 é coerente com a nova proibição no prompt e com a conferência pós-execução de RF-10; o incidente codex→opencode de 2026-09-17 citado existe (gotchas.md linha 43).
- O mapeamento de severidade em RF-05/AC-06 continua fiel a `code-review` (Critical/Major bloqueiam), `document-review` (lacuna/conflito/critério fraco bloqueiam) e `security-review` (vulnerabilidade demonstrada/exposição confirmada).
- "A alteração é revertida" em AC-12 é bem definido porque a revisão roda sobre estado congelado (`.agents/coordinator.md` linhas 59 e 61).
- "Quando ele existe" em RF-10 delimita corretamente o modo somente-leitura técnico aos binários que o oferecem (`codex -s read-only`, `mcode --permission off`) e cai na conferência pós-execução para os demais.
- Todos os pontos do pedido do dono e do brief seguem mapeados (ver rodada 1); RF-01..RF-10 e AC-01..AC-13 únicos e sequenciais; todo AC ligado a um RF; nove seções do template presentes; sem excesso — cada frase nova rastreia a um achado da rodada 1, ao pedido do dono ou a uma entrada da fleet-clis.

## Evidência

- Rodada 2: relido o PRD corrigido por inteiro (hash e tamanho conferidos) e o relatório persistido da rodada 1 (P1–P9, todos `pending` antes desta rodada).
- Relido `~/.claude/skills/fleet-clis/gotchas.md` por inteiro (50 linhas, 45 entradas) para verificar a nova afirmação da seção 1; encontrada na linha 43. Isso corrige a nota de Evidência da rodada 1, que dizia não existir entrada de delegação — a leitura anterior estava truncada na linha 42.
- Restrições do kit reconferidas contra o texto alterado: `.commands/setup.md` (write set, linhas 9 e 25), `.commands/review.md` e `.commands/secure.md` (registro persistente, linha 21 em cada; limite do padrão Opus, linha 28), `.skills/security-review/SKILL.md` (formato de saída sem seção de veredito, linhas 67–101), `.agents/coordinator.md` (linhas 59, 61, 89, 91, 110).
- Rodada 1 (2026-09-21T16:30-03:00): lidos o PRD original, templates/pt-br/PRD.md, docs/prd/PRD-002-validador-de-documentos.md, profiles/README.md, .agents/{coordinator,code-reviewer,document-validator,security-reviewer}.md, .skills/{code-review,document-review,security-review}/SKILL.md, .skills/project-onboarding/SKILL.md, .commands/{review,discover,secure,setup,verify}.md, CLAUDE.md, .harness/project.yaml, .gitignore, evals/cases/doc-validator-gap.md e ~/.claude/skills/fleet-clis/{SKILL,agy,codex,grok,opencode,mcode,claude,gotchas}.md.
- Limite de rodadas atingido: o documento seguiu ao dono com este relatório; o dono fechou P10–P12 e as quatro decisões pendentes em 2026-09-21, e o Coordenador aplicou as respostas no PRD (sha256 pós-edição registrado em `.harness/MEMORY.md`).
