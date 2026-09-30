# Review · PRD-010 · round 1

- Reviewer: document-validator (claude, opus requested; effective model unverified). Date: 2026-09-30.

## Verdict
- changes required
- Blocking findings: 13
- Report: docs/prd/PRD-010-kit-enxuto.review.md (the Coordinator writes the body below to this path)

## Findings
- P1
  - Category: gap
  - Severity: blocking
  - Location: docs/prd/PRD-010-kit-enxuto.md, section 3, R5/R7 (and the R1/R10 "validate clean" checks)
  - Evidence: `internal/kit/validate.go:18` sets `MinAgents = 18`, and `checkCounts` (line 334) fails with `agents: 12 < 18`. No rule lowers the minimum. As written, going from 19 agents to 12 makes every "`validate` verde / sai limpo" check fail, and the only ways out the builder can see are `SkipMinimumCounts` or an unplanned code change.
  - Suggestion: add to R5: "`MinAgents` in `internal/kit/validate.go` goes to 12 (its test too); `MinCommands` stays at 17 or goes to 19."
  - Reported by: claude
- P2
  - Category: conflict
  - Severity: blocking
  - Location: section 3, R6 ("cláusula anti-delegação idêntica nos dois modos"; check "`grep -c` da cláusula em `reviewer.md` e no Coordenador igual ao de hoje")
  - Evidence: today the clause is only in `.agents/coordinator.md` (fenced block under `## External CLI reviewers`, lines 112-133) and in `.skills/external-clis/SKILL.md`. `code-reviewer.md` and `security-reviewer.md` do not contain it (grep for RISK-001/anti-delegation in `.agents/` hits only coordinator.md). `checkAntiDelegationClause` compares coordinator with the skill, not with the reviewers. A builder following R6 would put a third copy into `reviewer.md` that no check covers, and "igual ao de hoje" for a file that does not exist today has no reference value.
  - Suggestion: "R6 Preserva: a cláusula anti-delegação continua só no Coordenador e em `external-clis` e é anexada ao prompt do revisor CLI na montagem; `reviewer.md` não a copia."
  - Reported by: claude
- P3
  - Category: weak criterion
  - Severity: blocking
  - Location: section 3, R6 check "`checkAntiDelegationClause` verde"; R10 (coordinator ≤ 150 lines)
  - Evidence: `internal/kit/evals.go:15-36` returns silently when coordinator.md has no `## External CLI reviewers` heading or no fenced block after it. R10 cuts coordinator.md from 223 lines to 150, and external-CLI text is 42% of it (from the audit). That makes moving or renaming this section likely, and the check would then stay green with no clause at all. The check cannot fail in the scenario it is supposed to protect.
  - Suggestion: add to R6/R10: "o bloco cercado continua sob `## External CLI reviewers` em `.agents/coordinator.md` (ou `checkAntiDelegationClause` passa a reprovar quando o bloco some); verificável: `grep -n '## External CLI reviewers' .agents/coordinator.md` encontra a seção e o bloco."
  - Reported by: claude
- P4
  - Category: gap
  - Severity: blocking
  - Location: section 3, R6 and R7 ("Preserva: …")
  - Evidence: `security-reviewer.md:24` is `model: opus` and `code-reviewer.md:24` is `model: sonnet`. `solution-architect.md` is `opus` and `api-designer.md` is `sonnet`. A merged agent has one frontmatter model. Neither R6 nor R7 says how each mode keeps its model, so a builder picks one and the security review is silently downgraded to sonnet, or the reverse.
  - Suggestion: add to R6: "modo `security` despachado com `opus`, modo `code` com `sonnet` (modelo explícito no despacho do Coordenador, `/dh:secure`, `/dh:build`, `/dh:review`)". Add the same to R7 (`architecture-decisions` → opus, `api-contracts` → sonnet).
  - Reported by: claude
- P5
  - Category: gap
  - Severity: blocking
  - Location: section 3, R5 check "`profiles/` continua válido no `validate`"
  - Evidence: `profiles/go-api.yaml:15,18` name `api-designer` and `data-engineer`, and `profiles/php.yaml:22` names `data-engineer`. `checkProfiles` (`validate.go:438-453`) only checks `id`, so `validate` passes while the profiles point at deleted roles. The profiles were a named preservation point.
  - Suggestion: add to R5/R7: "`profiles/*.yaml` passam a citar `builder` (lane `dados`) e `architect`; verificável: `grep -rnE 'api-designer|data-engineer|solution-architect|backend-builder' profiles/` vazio."
  - Reported by: claude
- P6
  - Category: gap
  - Severity: blocking
  - Location: section 3, R5 ("carrega a skill da lane (`incremental-implementation`, `data-migrations`, `safe-refactoring`, `systematic-debugging`)")
  - Evidence: `.commands/build.md:15` routes by lane and keeps a lane-less fallback by kind of work (server → backend, UI → frontend, schema → dados; `infra` is not reachable through the fallback). Every lane loads `incremental-implementation`, and `dados` also loads `data-migrations`. R5 names neither the fallback nor this skill composition, and its list mixes mode skills with lane skills. One reading gives lane `dados` only `data-migrations` and drops `incremental-implementation`. A lane-less ticket reaches `builder` with no lane.
  - Suggestion: "R5 Roteamento de `.commands/build.md` preservado: lane → `builder` com essa lane; sem lane, o fallback por tipo de trabalho atual escolhe a lane (`infra` fora do fallback); toda lane carrega `incremental-implementation`, `dados` também `data-migrations`; modo `refactor` → `safe-refactoring`, modo `debug` → `systematic-debugging`; lanes de teste seguem só com `qa-verifier`."
  - Reported by: claude
- P7
  - Category: gap
  - Severity: blocking
  - Location: section 3, R3 and its check (`grep -c append_system_prompt …`); section 2, kept cases 3 and 5
  - Evidence: `evals/cases/external-clis-risk-001/fixtures/code-reviewer.md` ships a copy of the agent, and `prompt.md:208` tells the model the prompt is "the body of .agents/code-reviewer.md, shipped with this case". R3's check only looks at `append_system_prompt`. Once R2 removes `checkFixtureAgentCopies`, the RISK-001 case keeps measuring the body of a deleted agent without anyone noticing. Also, `external-clis-risk-001` and `auto-hard-stop-recorded-state` embed the Coordinator body (prompt.md:9), which R10 rewrites. R3's route of "nomear o agente (`dh:reviewer`)" does not cover the main-session coordinator, and R3 does not say what happens if neither alternative works.
  - Suggestion: extend R3: "inclui cópias em `fixtures/` (o `external-clis-risk-001` passa a usar `.agents/reviewer.md`); os casos com corpo do Coordenador seguem a mesma abordagem; se a hipótese cair, os casos mantidos ficam marcados como `não rodável` até nova decisão do dono, sem voltar a cópia versionada". Check: `grep -rln "code-reviewer" evals/cases` returns nothing, and no file in `evals/cases/*/fixtures/` has an agent body.
  - Reported by: claude
- P8
  - Category: gap
  - Severity: blocking
  - Location: section 3, R2 ("`checkFixtureAgentCopies` e helpers só usados por elas … removida")
  - Evidence: `internal/kit/evals.go:415-420`: the only rejection of a symlinked `fixtures/` in a case is inside `checkFixtureAgentCopies`. `checkCaseGraders` only covers `graders/`. Removing the function, "com seus testes", drops the containment guard for case directories that the code documents as intentional (lines 191-196). The owner asked that nothing be taken away.
  - Suggestion: "R2 A rejeição de `fixtures/` como link simbólico (hoje em `checkFixtureAgentCopies`) é preservada, movida para `checkEvalReferences` com seu teste."
  - Reported by: claude
- P9
  - Category: ambiguity
  - Severity: blocking
  - Location: section 2, "vão para `evals/archive/`, fora de qualquer checagem do `validate`"; R2
  - Evidence: `checkEvalsCoverage`/`evalsExcluded` (`evals.go:312-354`) and `dh build` copy all of `evals/` except `results/`, `baselines/` and `not-run/`. `evals/archive/` would go into `dist/` and be compared by hash. The two readings give different software: (a) archive shipped and checked, which contradicts "fora de qualquer checagem"; (b) `archive` added to the build and validate exclusions.
  - Suggestion: "`archive` entra na lista de exclusão de primeiro nível do `dh build` e de `evalsExcluded`; verificável: `ls dist/*/evals/archive` não existe."
  - Reported by: claude
- P10
  - Category: weak criterion
  - Severity: blocking
  - Location: section 3, R2 check "`grep -rn resync AGENTS.md .agents .skills .commands docs` vazio"
  - Evidence: `docs/prd/PRD-010-kit-enxuto.md` itself mentions `resync.py` (lines 18, 72, 168), and ADR-003 (R9) will record the replaced rule. The grep can never come back empty, so the check fails even on a correct delivery.
  - Suggestion: "… `grep -rn resync AGENTS.md .agents .skills .commands docs --exclude-dir=prd --exclude-dir=adr` vazio."
  - Reported by: claude
- P11
  - Category: conflict
  - Severity: blocking
  - Location: section 3, R4 ("nenhum fluxo do Coordenador, comando … a dispara sozinho"; check "`grep` … por 'eval' mostra só a frase …")
  - Evidence: `.commands/improve.md:8,16,19,22` loads `harness-evaluation` and requires "a baseline must be captured before the asset changes" plus a before-and-after report. The owner kept `improve`. `grep eval` also matches "harness-evaluation", "evaluate" and `evals/cases/` in improve.md, auto.md, build.md, review.md, plan-loop.md and coordinator.md (22 hits). Meeting R4 means either cutting `improve`'s measurement (loss of capability) or failing the check.
  - Suggestion: "R4 Eval paga só roda por pedido explícito do dono ou dentro de `/dh:improve`, que o dono invoca; nenhum outro comando, fluxo do Coordenador ou CI a dispara. Verificável: `grep -rn 'claude plugin eval' .commands .agents` só em `improve.md`/`harness-evaluation`."
  - Reported by: claude
- P12
  - Category: weak criterion
  - Severity: blocking
  - Location: section 3, R4 (phrase "a pedido do dono") and R11 ("`grep` do Coordenador por 'painel de 3'")
  - Evidence: AGENTS.md line 7 says `.agents/` and `.commands/` are written in English. Portuguese phrases cannot appear in `coordinator.md` or the commands without breaking that fixed rule, so both greps are vacuous (zero hits on a correct delivery) or require breaking the language rule.
  - Suggestion: write the grep phrases in English as they will appear in the file, e.g. `"only when the owner asks"` and `"panel of three"`.
  - Reported by: claude
- P13
  - Category: weak criterion
  - Severity: blocking
  - Location: section 3, R13 check "a tabela tem 19 linhas … (`grep -c` das linhas `/dh:` em cada língua)"
  - Evidence: R13 requires flow sections (discover → release, auto, fix, handoff, plan-loop…) that mention `/dh:` many times. `grep -c '/dh:'` over the whole HTML counts those lines too and will not give 19 on a correct tutorial.
  - Suggestion: "a tabela de referência (`<table id=\\"comandos\\">`) tem 19 `<tr>` de corpo, um por arquivo de `.commands/`, conferido por nome."
  - Reported by: claude
- P14
  - Category: ambiguity
  - Severity: non-blocking
  - Location: section 2, lines 36-37 ("A verificação padrão fica em checagens determinísticas mais um revisor Claude"); R11
  - Evidence: R11 is correctly marked as PROPOSTA with its pending decision in the appendix, but the Solution states the same standard as already decided. The current rule (AGENTS.md "Revisão enxuta" (2)) also brings in the panel of three for a new contract and for logic, which R11 would drop. R11 also forbids "sondas novas sem pedido do dono", which touches fixed decisions (`mcode` "até ter sonda medida", `agy` "só entra se a sonda de cwd passar"). "achado Major" is not a term the kit defines (the kit uses blocking/non-blocking).
  - Suggestion: in the Solution, "…mais um revisor Claude (proposta R11, pendente do dono)". In R11, "achado bloqueante" instead of "Major", and a note that the PRD-008 sondas stay authorized.
  - Reported by: claude
- P15
  - Category: ambiguity
  - Severity: non-blocking
  - Location: section 3, R9
  - Evidence: R9 lets "o que for memória de projeto" go to EPOCHAL, while the check requires each earlier decision to appear in AGENTS.md or ADR-003. A rule moved to EPOCHAL fails the check, or the owner loses sight of it. "Sem data" also clashes with rules whose date is operational ("revisitar em 21/10/2026" in PRD-003/PRD-008).
  - Suggestion: "Toda decisão anterior aparece no novo AGENTS.md ou no ADR-003 (EPOCHAL é cópia adicional, nunca o único destino); datas de revisão pendentes ficam na regra vigente."
  - Reported by: claude
- P16
  - Category: excess
  - Severity: non-blocking
  - Location: section 2 (debugger in the builder); R13 (removal of 4 docs)
  - Evidence: the owner's item (4) names "construtores", the reviewers and the architects. The `debugger` merge is an addition that the PRD itself flags ("também"). R13 deletes `inicio-rapido`, `quick-start` and tutorial 00, but the owner asked for a new tutorial and said "não quero que vc tire coisas". The PRD gives a reason (appendix), but the owner did not ask for it.
  - Suggestion: list both as an explicit choice for the owner to confirm, or keep them as recorded excesses.
  - Reported by: claude
- P17
  - Category: ambiguity
  - Severity: non-blocking
  - Location: section 3, R14
  - Evidence: PRD-007 AC-28 (RF-14) sets the "Resumo executado" block with Entregue, Regras, Tickets, Docs and Fora. R14 only says "o que foi entregue, com a versão 0.13.0 e o PR #21".
  - Suggestion: "…com os campos do RF-14 do PRD-007 (Entregue, Regras, Tickets, Docs, Fora)."
  - Reported by: claude
- P18
  - Category: weak criterion
  - Severity: non-blocking
  - Location: section 3, R8 and R12 checks
  - Evidence: in R8, "uma frase-chave de cada procedimento" is chosen by whoever delivers, and a line count that goes down passes with any cut. In R12, "`git diff --stat main` … mostra se a prova foi necessária" does not observe the rule. The other halves of both checks are observable, so neither blocks.
  - Suggestion: R8: name the three phrases in the PRD. R12: keep only "a regra em `AGENTS.md` e em `implementation-planning` cita os quatro caminhos".
  - Reported by: claude

## Not raised
- Capability: the 19 commands exist in `.commands/` (all counted, including `plan-loop`, `improve`, `verify`, `consolidate-memory`, `understand`, `release`), and R1 keeps them. No command or capability the owner kept (external CLIs with clone, document validator for PRD/plan, pt/en docs) is removed by the PRD. Owner item 9 (keep pt and en) is kept by R1 and R13.
- Agent count: there are 19 real files. 19 − 10 (6 builders + 2 reviewers + 2 architects) + 3 = 12, and the list in section 2 matches exactly.
- R6 tools: `code-reviewer` and `security-reviewer` have the same `tools:` list (Skill, Read, Grep, Glob, LSP, ToolSearch, Monitor, SendMessage, WebFetch), so "mesma lista" holds. `/dh:secure` (`secure.md:7`) is named in R6, and the `reviewers:` field of `project.yaml` is kept.
- The 5 kept cases exist in `evals/cases/`. The cited baselines exist (`2026-09-24/25-security-reviewer-*`, `code-reviewer-*`, `document-validator-*`, `2026-09-28-auto-hard-stop-recorded-state-0.9.0.json`), and RISK-001 has `2026-09-25-risk-001-*`. There are 23 cases in total, so 18 are archived. `ls evals/cases` hides `.gitkeep`, so the "5" check holds.
- R2: `checkEmbeddedAgentBodies` and `checkFixtureAgentCopies` exist (evals.go:372/407, validate.go:96), and the helpers used only by them (agentBodies, matchingAgentBody, lineOverlap, firstNonEmptyLine, bodyAfterFrontmatter, yamlBlockScalar) fit "só usados por elas". `sameTrimmedLines` is used by the anti-delegation check and is kept by that clause. evals.go has 589 lines, as the PRD says. "no longer matches" only appears in evals.go.
- Removing the copy check is sourced (owner's item 1: "sem a checagem de cópia … nem o `resync.py`"). `resync.py` is cited only in AGENTS.md and PRD-010.
- R3: correctly written as a hypothesis with a pending decision and a falsifier in the appendix.
- R11: marked "PROPOSTA, depende de confirmação do dono", with a pending decision that blocks R11 and the Coordinator text.
- R12 matches the owner's item 10 ("só quando o PR mexe em Go ou binários") and keeps the proof for this PR through R15.
- Consistency with the 0.13.0 roadmap check: the `## 0.14.0 — <data>` format in R15 matches `changelogHeading` (`roadmap_check.go:32`). Card, chips, Fontes line and pt/en parity are named in R15. "Contagem atual" (the target agent and command totals) is required by R15 and enforced by `checkDocumentedCounts`. Deleting pages in R13 is covered by `checkLinks` in `validate`.
- Owner decisions against AGENTS.md: every one of the 43 current bullets is covered by R9 (kept rule or ADR-003), except the EPOCHAL route (P15). The PRD does not rewrite any fixed decision by itself, apart from R11 (a proposal) and R12 (the owner's item 10).
- R14: PRD-009 has status `aprovado (2026-09-30)` and no "Resumo executado" today. PR #21 = commit 1d0d66d (0.13.0).
- Organization: identifiers R1-R15 are unique and in order, and the Docs section lists the files touched. There are no duplicate rules.

## Evidence
- Read: docs/prd/PRD-010-kit-enxuto.md, AGENTS.md, internal/kit/evals.go (whole file), internal/kit/validate.go (lines 80-110, 330-460), internal/kit/doc_counts.go (55-135), internal/kit/roadmap_check.go (regexes), .commands/build.md; grep in .agents/, .commands/, .skills/, profiles/, templates/, evals/cases (the 5 kept cases, including fixtures), evals/baselines, docs/prd/PRD-007 and PRD-009, CHANGELOG.md.
- Round 1: there was no previous report.
- Not checked: the numbered list from the audit (what options 2, 3, 7 and 8 were). The source I received only gives the verdicts, so I checked the "no" items against the "Kit enxuto" bullet, which lists what stays. I did not check the "19 of 23 cases copy an agent body" claim case by case: all 23 have `append_system_prompt`, but I did not confirm which ones copy a full body. This does not affect any rule. No command was run (read-only).

## Coordinator (on the owner's behalf)
- All accepted. P16: `debugger` stays a separate agent (owner: "não quero que vc tire coisas"); the old quick-start and tutorial 00 pages stay as short pages linking to the new tutorial instead of being deleted.

---

# Review · PRD-010 · round 2 (final)

## Verdict
- changes required · blocking 1 (P19, introduced by the P7 fix)

## Settlement
- P1–P18 applied.
- P19 · weak criterion · blocking · R3 grep for "code-reviewer" can never be empty while the kept case `code-reviewer-equivalent-form-miss` exists. Suggestion: grep for `.agents/code-reviewer|fixtures/code-reviewer.md` and assert the risk-001 fixture copy is gone. Reported by: claude
- P20 · organization · the four link pages missing from §4 Docs. Reported by: claude
- P21 · ambiguity · R11 must keep the `/dh:plan-loop` full-list exception. Reported by: claude

## Coordinator (owner's delegation)
- P19, P20, P21 applied verbatim after round 2. PRD-010 marked approved by delegation (2026-09-30).
