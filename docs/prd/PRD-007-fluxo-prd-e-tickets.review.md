# Review · PRD-007 · round 1

- Reviewer: document-validator (claude, model opus requested; effective model unverified)
- Date: 2026-09-29
- Source: docs/novo-fluxo-prd.md; CLAUDE.md (owner decisions 1a–4a of 2026-09-29, lean review, docs inside the PR); current kit text; docs/prd/PRD-006-comando-autonomo.md

## Verdict
- changes required
- Blocking findings: 8

## Findings
- P1 · conflict · blocking · §4 RF-16 / AC-32 vs RF-03. Proposal guardrail 4 gives rule amendments to the PRD creator; RF-03 agrees, RF-16 gives the mid-execution amendment to the Coordinator. Suggestion: the Coordinator re-dispatches `product-discovery` to amend (Docs or Regras), then `document-validator` runs one revalidation round. Reported by: claude
- P2 · gap · blocking · §4 RF-05 / §5 AC-07. The proposal's ticket has a "Pronto quando" block (guardrail 20 depends on it); RF-05/AC-07 omit it while RF-08/AC-14 reference it. Suggestion: add Pronto quando with the proposal's minimum content (flow happens, each copied R<n> has evidence, relevant tests green, listed docs updated). Reported by: claude
- P3 · conflict · blocking · §4 RF-05 / §6 Dados. Decision 1a says the new format replaces PRD/STORY/TASK; RF-05 keeps every current TASK section and adds STORY fields, leaving two acceptance sources (Fluxo vs §7 AC-xx). Suggestion: RF-05 names the closed list of ticket sections, states which current TASK/STORY sections are removed, and makes Fluxo + Pronto quando the only acceptance source; keeping technical sections beyond 1a goes to §8 as an owner decision. Reported by: claude
- P4 · conflict · blocking · §4 RF-17, §5 AC-33, §8 second risk. (a) Wrong citation: the per-slice code-review assumption lives in PRD-006 RF-09, RF-12 and AC-14 and in `.commands/auto.md` line 19, not in PRD-006 RF-01. (b) Editing an approved contract changes the shape of delivery (CLAUDE.md "perguntar antes de escolher"; 1a "PRDs antigos ficam como estão"). Suggestion: fix the citation; move to §8 as a pending owner decision: (1) amend PRD-006, or (2) leave PRD-006 untouched with a supersession note in PRD-007 plus the auto.md edit; recommend one. Reported by: claude
- P5 · conflict · blocking · §4 RF-12 / AC-23. PRD-006 RF-12 and auto.md run `/dh:secure` per slice; PRD-007 places it in the final review while claiming "do mesmo jeito". Suggestion: state trigger = `.commands/secure.md`, position = once in the final feature review (or per ticket), and list the change in RF-17. Reported by: claude
- P6 · conflict · blocking · §4 RF-10 / AC-19; §7 row "Nova lane kit/harness". Routing kit-source work to `harness-maintainer` under lane `backend` breaks the proposal's lane→specialist rule, is an author decision presented as requirement, and puts kit-specific checks into the generic flow. Suggestion: move to §8 as an owner decision: (a) builder override scoped to this repository; (b) routing note in dev-harness CLAUDE.md only. Keep AC-18 as the only lane→builder rule meanwhile. Reported by: claude
- P7 · conflict · blocking · §4 RF-12, §6 Stack, AC-24. Decision 4a names Playwright headless; the PRD downgrades to "any headless browser", and ui-verification marks rows not-run without a browser, so guardrail 18 degrades silently. Suggestion: state Playwright headless and that its absence makes the feature verdict partial/blocked, never approved; or record tool-agnostic as an owner decision. Reported by: claude
- P8 · gap · blocking · §6 Contratos; RF-05/RF-06 claims. Omitted files: `.agents/coordinator.md` (Procedure 7 QA then code review; Flows Feature; STORY template mention), `.commands/document.md` (STORY), `.commands/discover.md` frontmatter `writes: .harness/stories/`, `.skills/implementation-planning/SKILL.md` (step 5 folds tests into the slice; step 10 does not check R<n>→test coverage), `templates/*/PRD.md` in Contratos. Suggestion: list them with the passage each must change; rephrase RF-06 as a new self-review pass in step 10. Reported by: claude
- P9 · excess · non-blocking · RF-01 / AC-01. Keeping the header table is an author decision (grounded in PRD-006 RF-02); list in §8 for the owner to confirm. Reported by: claude
- P10 · conflict · non-blocking · RF-14. The kit rule is "only the Coordinator writes the three records", not "the Coordinator writes only them". Suggestion: "O Coordenador grava a seção Resumo executado no PRD; nenhum outro trecho do PRD é escrito pelo Coordenador." Reported by: claude
- P11 · ambiguity · non-blocking · RF-14 / AC-29. No marker for "PRD entregue"; name it (e.g. `Status` → `entregue em <data>` only when no rule is Não honrada without owner decision). Reported by: claude
- P12 · excess · non-blocking · RF-01 relocation of alternatives/risks to tickets and RISKS.md does not trace; drop it and let §8 carry the appendix decision. Reported by: claude
- P13 · organization · non-blocking · §4 "Não entra" bullet 5: "named" → "nomeia". Reported by: claude

## Not raised
- Decisions 1a (no migration), 2a, 3a, 4a (product-discovery judges only), lean review for rejection, one-round revalidation, 6-round per-ticket cap: traced.
- Guardrail 9 resolved by deferring closure (RF-07, AC-11/12): grounded; uses existing `em revisão` status; does not break implementation-planning step 6.
- Guardrails 1–20 mapped; 1, 6, 16 covered by existing kit rules.
- RF-15/AC-30 deterministic literal-copy check correctly a suggestion; validate.go scope claim matches.
- Decomposable into tickets once P2, P3, P6 are resolved; RF/AC ids unique and cross-referenced.

## Evidence
- Read: PRD-007; docs/novo-fluxo-prd.md (English mirror not read); CLAUDE.md; .agents/coordinator.md; .agents/product-discovery.md (head); .commands/build.md, discover.md, plan.md, review.md, document.md, auto.md (grep); PRD-006 lines 30–49, AC-03, AC-14; templates/pt-br/TASK.md; .skills/ui-verification/SKILL.md (grep); .skills/implementation-planning/SKILL.md steps 1–11; internal/kit/validate.go checkTemplateRefs.
- Not checked: whether validate.go's templatePattern matches the bare "STORY.md" in document.md; bodies of qa-verifier.md and harness-maintainer.md.
- No scripts run, no files edited.

---

# Review · PRD-007 · round 2 (final)

- Reviewer: document-validator (claude, opus requested; effective model unverified). Date: 2026-09-29.

## Verdict
- changes required
- Blocking findings: 2 (both new; P1–P13 of round 1 all applied)
- Round cap reached (2 of 2): the document goes to the owner with P14 and P15 open.

## Settlement of round 1
- P1–P13: applied. P5 judged consistent: /dh:secure stays per ticket after the qa-verifier gate, matching PRD-006 RF-12 and `.commands/auto.md` line 19. Observation, not a finding: a security issue that exists only once tickets integrate is seen by no security reviewer.

## New findings
- P14 · conflict · blocking · §5 AC-07 "contendo exatamente a lista fechada" presupposes the answer to the open §8 decision on TASK technical sections; `.commands/build.md` line 21 needs a place for status/results. Suggestion: closed list plus only the technical sections the owner keeps in §8; no separate acceptance section; Fluxo + Pronto quando the only acceptance source. Reported by: claude
- P15 · ambiguity · blocking · §6 Contratos, `.agents/coordinator.md` bullet: "fluxo de slice legado" is undefined. Reading A: per-unit code review survives only outside a ticketed feature (/dh:fix, /dh:refactor, harness instruction changes), build.md loses it entirely. Reading B: build.md keeps a per-slice branch for pre-PRD-007 plans, contradicting AC-15. Suggestion: state Reading A explicitly (or amend AC-15 for B). Reported by: claude
- P16 · organization · non-blocking · RF-17 quote attributed to PRD-006 RF-09 is actually AC-14 (line 77). Reported by: claude
- P17 · gap · non-blocking · no file in §6 carries the Playwright client-pass rule; §9 still calls ui-verification tool-agnostic. Suggestion: list `.commands/review.md` or `.skills/ui-verification/SKILL.md` with "Playwright headless; ausência → not-run, veredito parcial/bloqueado". Reported by: claude

## Evidence
- Read: corrected PRD-007; round-1 report; grep of `.commands/` for TASK/BUG-/2B; PRD-006 RF-09, RF-12, AC-14. Not checked: bodies of qa-verifier.md and harness-maintainer.md; plan-loop.md judged covered by the implementation-planning edit. No scripts, edits or dispatches.
