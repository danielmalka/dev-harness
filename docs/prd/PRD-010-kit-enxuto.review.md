# Revisão · PRD-010 · rodada 1

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

- Revisor: document-validator (claude, opus solicitado; modelo efetivo não verificado). Data: 2026-09-30.

## Veredito
- `changes required`
- Achados bloqueantes: 13
- Relatório: docs/prd/PRD-010-kit-enxuto.review.md (o Coordenador grava o corpo abaixo neste caminho)

## Achados
- P1
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: docs/prd/PRD-010-kit-enxuto.md, seção 3, R5/R7 (e as checagens "validate limpo" de R1/R10)
  - Evidência: `internal/kit/validate.go:18` define `MinAgents = 18`, e `checkCounts` (linha 334) falha com `agents: 12 < 18`. Nenhuma regra baixa o mínimo. Como está escrito, passar de 19 agentes para 12 faz toda checagem "`validate` verde / sai limpo" falhar, e as únicas saídas que o builder enxerga são `SkipMinimumCounts` ou uma mudança de código não planejada.
  - Sugestão: acrescentar à R5: "`MinAgents` em `internal/kit/validate.go` passa a 12 (o teste dele também); `MinCommands` continua em 17 ou passa a 19."
  - Reportado por: claude
- P2
  - Categoria: `conflict`
  - Severidade: `blocking`
  - Local: seção 3, R6 ("cláusula anti-delegação idêntica nos dois modos"; checagem "`grep -c` da cláusula em `reviewer.md` e no Coordenador igual ao de hoje")
  - Evidência: hoje a cláusula está só em `.agents/coordinator.md` (bloco cercado sob `## External CLI reviewers`, linhas 112-133) e em `.skills/external-clis/SKILL.md`. `code-reviewer.md` e `security-reviewer.md` não a contêm (o grep por RISK-001/anti-delegation em `.agents/` só acha coordinator.md). `checkAntiDelegationClause` compara o coordinator com a skill, não com os revisores. Um builder que siga a R6 poria uma terceira cópia em `reviewer.md` que nenhuma checagem cobre, e "igual ao de hoje" para um arquivo que não existe hoje não tem valor de referência.
  - Sugestão: "R6 Preserva: a cláusula anti-delegação continua só no Coordenador e em `external-clis` e é anexada ao prompt do revisor CLI na montagem; `reviewer.md` não a copia."
  - Reportado por: claude
- P3
  - Categoria: `weak criterion`
  - Severidade: `blocking`
  - Local: seção 3, R6 checagem "`checkAntiDelegationClause` verde"; R10 (coordinator ≤ 150 linhas)
  - Evidência: `internal/kit/evals.go:15-36` retorna em silêncio quando o coordinator.md não tem o título `## External CLI reviewers` ou não tem bloco cercado depois dele. A R10 corta o coordinator.md de 223 linhas para 150, e o texto de CLI externo é 42% dele (segundo a auditoria). Isso torna provável mover ou renomear essa seção, e a checagem continuaria verde sem cláusula nenhuma. A checagem não pode falhar no cenário que deveria proteger.
  - Sugestão: acrescentar à R6/R10: "o bloco cercado continua sob `## External CLI reviewers` em `.agents/coordinator.md` (ou `checkAntiDelegationClause` passa a reprovar quando o bloco some); verificável: `grep -n '## External CLI reviewers' .agents/coordinator.md` encontra a seção e o bloco."
  - Reportado por: claude
- P4
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: seção 3, R6 e R7 ("Preserva: …")
  - Evidência: `security-reviewer.md:24` é `model: opus` e `code-reviewer.md:24` é `model: sonnet`. `solution-architect.md` é `opus` e `api-designer.md` é `sonnet`. Um agente fundido tem um só modelo no frontmatter. Nem a R6 nem a R7 dizem como cada modo mantém o seu modelo, então o builder escolhe um e a revisão de segurança é rebaixada em silêncio para sonnet, ou o inverso.
  - Sugestão: acrescentar à R6: "modo `security` despachado com `opus`, modo `code` com `sonnet` (modelo explícito no despacho do Coordenador, `/dh:secure`, `/dh:build`, `/dh:review`)". Acrescentar o mesmo à R7 (`architecture-decisions` → opus, `api-contracts` → sonnet).
  - Reportado por: claude
- P5
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: seção 3, R5 checagem "`profiles/` continua válido no `validate`"
  - Evidência: `profiles/go-api.yaml:15,18` nomeiam `api-designer` e `data-engineer`, e `profiles/php.yaml:22` nomeia `data-engineer`. `checkProfiles` (`validate.go:438-453`) só confere o `id`, então o `validate` passa enquanto os perfis apontam para papéis apagados. Os perfis eram um ponto de preservação nomeado.
  - Sugestão: acrescentar à R5/R7: "`profiles/*.yaml` passam a citar `builder` (lane `dados`) e `architect`; verificável: `grep -rnE 'api-designer|data-engineer|solution-architect|backend-builder' profiles/` vazio."
  - Reportado por: claude
- P6
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: seção 3, R5 ("carrega a skill da lane (`incremental-implementation`, `data-migrations`, `safe-refactoring`, `systematic-debugging`)")
  - Evidência: `.commands/build.md:15` roteia por lane e mantém um fallback sem lane por tipo de trabalho (servidor → backend, UI → frontend, schema → dados; `infra` não é alcançável pelo fallback). Toda lane carrega `incremental-implementation`, e `dados` também carrega `data-migrations`. A R5 não nomeia nem o fallback nem essa composição de skills, e a lista dela mistura skills de modo com skills de lane. Uma leitura dá à lane `dados` só `data-migrations` e descarta `incremental-implementation`. Um ticket sem lane chega ao `builder` sem lane.
  - Sugestão: "R5 Roteamento de `.commands/build.md` preservado: lane → `builder` com essa lane; sem lane, o fallback por tipo de trabalho atual escolhe a lane (`infra` fora do fallback); toda lane carrega `incremental-implementation`, `dados` também `data-migrations`; modo `refactor` → `safe-refactoring`, modo `debug` → `systematic-debugging`; lanes de teste seguem só com `qa-verifier`."
  - Reportado por: claude
- P7
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: seção 3, R3 e a checagem dela (`grep -c append_system_prompt …`); seção 2, casos mantidos 3 e 5
  - Evidência: `evals/cases/external-clis-risk-001/fixtures/code-reviewer.md` traz uma cópia do agente, e `prompt.md:208` diz ao modelo que o prompt é `"the body of .agents/code-reviewer.md, shipped with this case"`. A checagem da R3 só olha `append_system_prompt`. Depois que a R2 remove `checkFixtureAgentCopies`, o caso RISK-001 continua medindo o corpo de um agente apagado sem que ninguém perceba. Além disso, `external-clis-risk-001` e `auto-hard-stop-recorded-state` embutem o corpo do Coordenador (prompt.md:9), que a R10 reescreve. O caminho da R3 de "nomear o agente (`dh:reviewer`)" não cobre o coordinator da sessão principal, e a R3 não diz o que acontece se nenhuma das alternativas funcionar.
  - Sugestão: estender a R3: "inclui cópias em `fixtures/` (o `external-clis-risk-001` passa a usar `.agents/reviewer.md`); os casos com corpo do Coordenador seguem a mesma abordagem; se a hipótese cair, os casos mantidos ficam marcados como `não rodável` até nova decisão do dono, sem voltar a cópia versionada". Checagem: `grep -rln "code-reviewer" evals/cases` não retorna nada, e nenhum arquivo em `evals/cases/*/fixtures/` tem corpo de agente.
  - Reportado por: claude
- P8
  - Categoria: `gap`
  - Severidade: `blocking`
  - Local: seção 3, R2 ("`checkFixtureAgentCopies` e helpers só usados por elas … removida")
  - Evidência: `internal/kit/evals.go:415-420`: a única rejeição de um `fixtures/` que seja link simbólico num caso está dentro de `checkFixtureAgentCopies`. `checkCaseGraders` só cobre `graders/`. Remover a função, "com seus testes", derruba a trava de contenção dos diretórios de caso que o código documenta como intencional (linhas 191-196). O dono pediu que nada fosse tirado.
  - Sugestão: "R2 A rejeição de `fixtures/` como link simbólico (hoje em `checkFixtureAgentCopies`) é preservada, movida para `checkEvalReferences` com seu teste."
  - Reportado por: claude
- P9
  - Categoria: `ambiguity`
  - Severidade: `blocking`
  - Local: seção 2, "vão para `evals/archive/`, fora de qualquer checagem do `validate`"; R2
  - Evidência: `checkEvalsCoverage`/`evalsExcluded` (`evals.go:312-354`) e o `dh build` copiam todo o `evals/` exceto `results/`, `baselines/` e `not-run/`. `evals/archive/` iria para o `dist/` e seria comparado por hash. As duas leituras dão software diferente: (a) archive distribuído e checado, o que contradiz "fora de qualquer checagem"; (b) `archive` acrescentado às exclusões do build e do validate.
  - Sugestão: "`archive` entra na lista de exclusão de primeiro nível do `dh build` e de `evalsExcluded`; verificável: `ls dist/*/evals/archive` não existe."
  - Reportado por: claude
- P10
  - Categoria: `weak criterion`
  - Severidade: `blocking`
  - Local: seção 3, R2 checagem "`grep -rn resync AGENTS.md .agents .skills .commands docs` vazio"
  - Evidência: o próprio `docs/prd/PRD-010-kit-enxuto.md` menciona `resync.py` (linhas 18, 72, 168), e o ADR-003 (R9) registrará a regra substituída. O grep nunca pode voltar vazio, então a checagem falha até numa entrega correta.
  - Sugestão: "… `grep -rn resync AGENTS.md .agents .skills .commands docs --exclude-dir=prd --exclude-dir=adr` vazio."
  - Reportado por: claude
- P11
  - Categoria: `conflict`
  - Severidade: `blocking`
  - Local: seção 3, R4 ("nenhum fluxo do Coordenador, comando … a dispara sozinho"; checagem "`grep` … por 'eval' mostra só a frase …")
  - Evidência: `.commands/improve.md:8,16,19,22` carrega `harness-evaluation` e exige `"a baseline must be captured before the asset changes"` mais um relatório de antes e depois. O dono manteve o `improve`. O `grep eval` também casa "harness-evaluation", "evaluate" e `evals/cases/` em improve.md, auto.md, build.md, review.md, plan-loop.md e coordinator.md (22 ocorrências). Cumprir a R4 significa cortar a medição do `improve` (perda de capacidade) ou falhar a checagem.
  - Sugestão: "R4 Eval paga só roda por pedido explícito do dono ou dentro de `/dh:improve`, que o dono invoca; nenhum outro comando, fluxo do Coordenador ou CI a dispara. Verificável: `grep -rn 'claude plugin eval' .commands .agents` só em `improve.md`/`harness-evaluation`."
  - Reportado por: claude
- P12
  - Categoria: `weak criterion`
  - Severidade: `blocking`
  - Local: seção 3, R4 (frase "a pedido do dono") e R11 ("`grep` do Coordenador por 'painel de 3'")
  - Evidência: a linha 7 do AGENTS.md diz que `.agents/` e `.commands/` são escritos em inglês. Frases em português não podem aparecer em `coordinator.md` nem nos comandos sem quebrar essa regra fixa, então os dois greps são vazios (zero ocorrências numa entrega correta) ou exigem quebrar a regra de idioma.
  - Sugestão: escrever as frases dos greps em inglês, como aparecerão no arquivo, por exemplo `"only when the owner asks"` e `"panel of three"`.
  - Reportado por: claude
- P13
  - Categoria: `weak criterion`
  - Severidade: `blocking`
  - Local: seção 3, R13 checagem "a tabela tem 19 linhas … (`grep -c` das linhas `/dh:` em cada língua)"
  - Evidência: a R13 exige seções de fluxo (discover → release, auto, fix, handoff, plan-loop…) que mencionam `/dh:` muitas vezes. `grep -c '/dh:'` sobre o HTML inteiro conta essas linhas também e não dará 19 num tutorial correto.
  - Sugestão: "a tabela de referência (`<table id=\\"comandos\\">`) tem 19 `<tr>` de corpo, um por arquivo de `.commands/`, conferido por nome."
  - Reportado por: claude
- P14
  - Categoria: `ambiguity`
  - Severidade: `non-blocking`
  - Local: seção 2, linhas 36-37 ("A verificação padrão fica em checagens determinísticas mais um revisor Claude"); R11
  - Evidência: a R11 está corretamente marcada como PROPOSTA com sua decisão pendente no apêndice, mas a Solução afirma o mesmo padrão como já decidido. A regra atual (AGENTS.md "Revisão enxuta" (2)) também traz o painel de três para um contrato novo e para lógica, o que a R11 eliminaria. A R11 também proíbe "sondas novas sem pedido do dono", o que toca decisões fixas (`mcode` "até ter sonda medida", `agy` "só entra se a sonda de cwd passar"). "achado Major" não é um termo que o kit defina (o kit usa bloqueante/não bloqueante).
  - Sugestão: na Solução, "…mais um revisor Claude (proposta R11, pendente do dono)". Na R11, "achado bloqueante" em vez de "Major", e uma nota de que as sondas do PRD-008 continuam autorizadas.
  - Reportado por: claude
- P15
  - Categoria: `ambiguity`
  - Severidade: `non-blocking`
  - Local: seção 3, R9
  - Evidência: a R9 deixa "o que for memória de projeto" ir para o EPOCHAL, enquanto a checagem exige que cada decisão anterior apareça no AGENTS.md ou no ADR-003. Uma regra movida para o EPOCHAL reprova a checagem, ou o dono perde a regra de vista. "Sem data" também colide com regras cuja data é operacional ("revisitar em 21/10/2026" no PRD-003/PRD-008).
  - Sugestão: "Toda decisão anterior aparece no novo AGENTS.md ou no ADR-003 (EPOCHAL é cópia adicional, nunca o único destino); datas de revisão pendentes ficam na regra vigente."
  - Reportado por: claude
- P16
  - Categoria: `excess`
  - Severidade: `non-blocking`
  - Local: seção 2 (debugger no builder); R13 (remoção de 4 docs)
  - Evidência: o item (4) do dono nomeia "construtores", os revisores e os arquitetos. A fusão do `debugger` é um acréscimo que o próprio PRD sinaliza ("também"). A R13 apaga `inicio-rapido`, `quick-start` e o tutorial 00, mas o dono pediu um tutorial novo e disse "não quero que vc tire coisas". O PRD dá uma razão (apêndice), mas o dono não pediu isso.
  - Sugestão: listar os dois como escolha explícita para o dono confirmar, ou mantê-los como excessos registrados.
  - Reportado por: claude
- P17
  - Categoria: `ambiguity`
  - Severidade: `non-blocking`
  - Local: seção 3, R14
  - Evidência: o AC-28 do PRD-007 (RF-14) define o bloco "Resumo executado" com Entregue, Regras, Tickets, Docs e Fora. A R14 só diz "o que foi entregue, com a versão 0.13.0 e o PR #21".
  - Sugestão: "…com os campos do RF-14 do PRD-007 (Entregue, Regras, Tickets, Docs, Fora)."
  - Reportado por: claude
- P18
  - Categoria: `weak criterion`
  - Severidade: `non-blocking`
  - Local: seção 3, checagens de R8 e R12
  - Evidência: na R8, "uma frase-chave de cada procedimento" é escolhida por quem entrega, e uma contagem de linhas que diminui passa com qualquer corte. Na R12, "`git diff --stat main` … mostra se a prova foi necessária" não observa a regra. As outras metades das duas checagens são observáveis, então nenhuma bloqueia.
  - Sugestão: R8: nomear as três frases no PRD. R12: manter só "a regra em `AGENTS.md` e em `implementation-planning` cita os quatro caminhos".
  - Reportado por: claude

## Não levantado
- Capacidade: os 19 comandos existem em `.commands/` (todos contados, incluindo `plan-loop`, `improve`, `verify`, `consolidate-memory`, `understand`, `release`), e a R1 os mantém. Nenhum comando ou capacidade que o dono manteve (CLIs externos com clone, validador de documentos para PRD/plano, docs pt/en) é removido pelo PRD. O item 9 do dono (manter pt e en) é mantido pela R1 e pela R13.
- Contagem de agentes: há 19 arquivos reais. 19 − 10 (6 builders + 2 revisores + 2 arquitetos) + 3 = 12, e a lista da seção 2 bate exatamente.
- Ferramentas da R6: `code-reviewer` e `security-reviewer` têm a mesma lista de `tools:` (Skill, Read, Grep, Glob, LSP, ToolSearch, Monitor, SendMessage, WebFetch), então "mesma lista" se sustenta. `/dh:secure` (`secure.md:7`) é nomeado na R6, e o campo `reviewers:` do `project.yaml` é mantido.
- Os 5 casos mantidos existem em `evals/cases/`. As baselines citadas existem (`2026-09-24/25-security-reviewer-*`, `code-reviewer-*`, `document-validator-*`, `2026-09-28-auto-hard-stop-recorded-state-0.9.0.json`), e o RISK-001 tem `2026-09-25-risk-001-*`. Há 23 casos no total, então 18 são arquivados. `ls evals/cases` esconde `.gitkeep`, então a checagem dos "5" se sustenta.
- R2: `checkEmbeddedAgentBodies` e `checkFixtureAgentCopies` existem (evals.go:372/407, validate.go:96), e os helpers usados só por elas (agentBodies, matchingAgentBody, lineOverlap, firstNonEmptyLine, bodyAfterFrontmatter, yamlBlockScalar) cabem em "só usados por elas". `sameTrimmedLines` é usado pela checagem de anti-delegação e é mantido por essa cláusula. O evals.go tem 589 linhas, como diz o PRD. `no longer matches` só aparece em evals.go.
- Remover a checagem de cópia tem fonte (item 1 do dono: "sem a checagem de cópia … nem o `resync.py`"). O `resync.py` só é citado no AGENTS.md e no PRD-010.
- R3: corretamente escrita como hipótese, com decisão pendente e um falsificador no apêndice.
- R11: marcada como "PROPOSTA, depende de confirmação do dono", com uma decisão pendente que bloqueia a R11 e o texto do Coordenador.
- A R12 bate com o item 10 do dono ("só quando o PR mexe em Go ou binários") e mantém a prova deste PR por meio da R15.
- Consistência com a checagem de roadmap da 0.13.0: o formato `## 0.14.0 — <data>` na R15 bate com `changelogHeading` (`roadmap_check.go:32`). Cartão, chips, linha Fontes e paridade pt/en são nomeados na R15. A "Contagem atual" (os totais-alvo de agentes e comandos) é exigida pela R15 e imposta por `checkDocumentedCounts`. Apagar páginas na R13 é coberto por `checkLinks` no `validate`.
- Decisões do dono contra o AGENTS.md: cada um dos 43 marcadores atuais é coberto pela R9 (regra mantida ou ADR-003), exceto o caminho EPOCHAL (P15). O PRD não reescreve nenhuma decisão fixa por conta própria, além da R11 (uma proposta) e da R12 (o item 10 do dono).
- R14: o PRD-009 tem status `aprovado (2026-09-30)` e hoje não tem "Resumo executado". PR #21 = commit 1d0d66d (0.13.0).
- Organização: os identificadores R1-R15 são únicos e estão em ordem, e a seção Docs lista os arquivos tocados. Não há regras duplicadas.

## Evidência
- Lido: docs/prd/PRD-010-kit-enxuto.md, AGENTS.md, internal/kit/evals.go (arquivo inteiro), internal/kit/validate.go (linhas 80-110, 330-460), internal/kit/doc_counts.go (55-135), internal/kit/roadmap_check.go (regex), .commands/build.md; grep em .agents/, .commands/, .skills/, profiles/, templates/, evals/cases (os 5 casos mantidos, incluindo fixtures), evals/baselines, docs/prd/PRD-007 e PRD-009, CHANGELOG.md.
- Rodada 1: não havia relatório anterior.
- Não verificado: a lista numerada da auditoria (o que eram as opções 2, 3, 7 e 8). A fonte que recebi só dá os vereditos, então conferi os itens "não" contra o marcador "Kit enxuto", que lista o que fica. Não conferi caso a caso a afirmação "19 de 23 casos copiam um corpo de agente": os 23 têm `append_system_prompt`, mas não confirmei quais copiam um corpo inteiro. Isso não afeta nenhuma regra. Nenhum comando foi rodado (somente leitura).

## Coordenador (em nome do dono)
- Todos aceitos. P16: o `debugger` continua um agente separado (dono: "não quero que vc tire coisas"); as páginas antigas de quick-start e do tutorial 00 continuam como páginas curtas que levam ao novo tutorial, em vez de serem apagadas.

---

# Revisão · PRD-010 · rodada 2 (final)

## Veredito
- `changes required` · bloqueantes 1 (P19, introduzido pela correção do P7)

## Assentamento
- P1–P18 aplicados.
- P19 · `weak criterion` · `blocking` · o grep da R3 por "code-reviewer" nunca pode ficar vazio enquanto o caso mantido `code-reviewer-equivalent-form-miss` existir. Sugestão: fazer grep por `.agents/code-reviewer|fixtures/code-reviewer.md` e afirmar que a cópia de fixture do risk-001 sumiu. Reportado por: claude
- P20 · `organization` · as quatro páginas de link faltam na §4 Docs. Reportado por: claude
- P21 · `ambiguity` · a R11 deve manter a exceção da lista completa do `/dh:plan-loop`. Reportado por: claude

## Coordenador (delegação do dono)
- P19, P20, P21 aplicados literalmente após a rodada 2. PRD-010 marcado como aprovado por delegação (2026-09-30).
