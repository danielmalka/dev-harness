# PRD-031 · Lean builder flow and test-first default

| Field | Value |
|---|---|
| Status | delivered on 2026-10-31 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-30 / 2026-10-30 |
| Tickets | to be defined by the planner |
| Origin | Owner request 2026-10-30, source `tasks/PRD-031/SOURCE.md` in the harness dir |
| Depends on | PRD-029 (finding triage, version base 0.35.0) |
| Delivery | One pull request |

## 1. Problem

The builder-phase flow is long and costs many tokens. The owner asks for a more direct and concrete flow.

- Fact: four files carry the builder phase. They total 830 lines and 107478 bytes (`wc`, 2026-10-30).

| File | Lines | Bytes |
|---|---|---|
| `.skills/incremental-implementation/SKILL.md` | 183 | 15466 |
| `.skills/implementation-planning/SKILL.md` | 190 | 17098 |
| `.agents/coordinator.md` | 207 | 49460 |
| `.skills/code-review/SKILL.md` | 250 | 25454 |

- Fact: `ub-mode` (Uncle Bob) covers its builder phase in about 70 lines of one-line imperative rules. Its whole `SKILL.md` has 152 lines and 3958 bytes.
- Hypothesis: most of the weight comes from per-slice stages and state, not from prose. The stages are finding ids `F<round>-<n>`, six finding states (`.agents/coordinator.md`, item 8: open, fixed, intentional, noise, escalated, reopened), the churn protocol, two consecutive clean rounds, the code-stage rotation cursor, and the review-bundle rules, and the in-flight doubt. This hypothesis is wrong if the stages total less than half of the four files. The rule inventory (R7) tests it.
- Fact: a feature has no test-first rule. `.agents/coordinator.md` says: "A feature has no red-first rule, so this single dispatch is the disproof." Red-first exists only for bug regression tests.
- Fact: no rule says what a commit message body must contain.
- Fact: `.skills/harness-authoring/SKILL.md` has no rule on the length of a procedure step.

## 2. Solution

After delivery:

- A builder writes the failing test first, sees it fail for the expected reason, and then changes the code.
- `regression-testing` has three heuristics that reject empty tests.
- A commit body states why, what changed, and how the change was proved.
- A new procedure step in a kit skill is one imperative sentence plus at most one sentence of reason.
- The four builder-phase files aim at 30% of today's size (ceiling 40%, R6), with no load-bearing safeguard lost.
- The per-slice review machinery is simpler. The owner accepts or rejects each simplification (R11 to R15). The owner accepted all five on 2026-10-30.

Part A is R1 to R5 (new behavior). Part B is R6 to R10 (compression and consistency). Part C is R11 to R15 (removals, accepted by the owner on 2026-10-30).

## 3. Rules

### Part A · New behavior

- R1 `incremental-implementation` tells the builder, for each implementation slice, to write the failing test, run it, see the expected failure, and only then change the code. The builder report names the failing result. Docs-only, config-only, and copy-only slices are exempt, and the report says which exemption applies. `regression-testing` carries the same default in its Procedure, so a feature slice and a bug fix both start with a failing test. Verifiable: `incremental-implementation` and `regression-testing` each contain the exact phrase `expected failure`. `incremental-implementation` also contains the phrases `docs-only`, `config-only`, `copy-only` and a Checks line `Red:` in its Output format. `regression-testing` contains `test-first`.
- R2 The In-flight doubt narrows. Test-first (R1) replaces the "disproof" reason, so `.agents/coordinator.md` no longer says that a feature has no red-first rule. The doubt stays only for a slice with `Doubt: yes`. Rule R14 sets the remaining shape.
- R3 `regression-testing` contains three heuristics, each as one sentence:
  - See the new test fail for the expected reason before it counts.
  - Delete a test that still passes when every function under test returns nothing.
  - When an honest test is impractical, report the check as not-run. Do not write a weak test.
  The existing red-first rule for bug regression tests stays.
- R4 Commit content rule. The commit body states why the change exists, what changed, and how it was proved. The "how" names the check that ran. The Conventional Commits header in `AGENTS.md` stays unchanged. The rule lives in `.skills/delivery-readiness/SKILL.md` only (see Decisions, D1). Verifiable: `grep -n "body" .skills/delivery-readiness/SKILL.md` shows one rule.
- R5 `.skills/harness-authoring/SKILL.md` contains this style rule: a procedure step is one imperative sentence plus at most one sentence of reason, and cross-references go at the end of the step or section, not inside a sentence. The rule applies to every step this PRD writes or rewrites.

### Part B · Compression and consistency

- R6 Size budget, a soft ceiling. After delivery, the four files from section 1 aim at 250 lines and 30000 bytes. The ceiling is 330 lines and 43000 bytes (about 40% of today). The ceiling can be exceeded by up to about 10% (363 lines, 47300 bytes) only when the rule inventory (R7) names the rules that need the extra space. It is not a hard gate. Measure with `wc -lc` on the four paths. No text moves to a new file to meet the budget. The byte limit exists because the Coordinator file holds 46% of the bytes. Verifiable: the `wc -lc` total is at or below the ceiling, or the inventory has a section "Budget excess" that names each rule that needs the space and the total stays within 10%.
  Owner decision 2026-10-31: the measured total (432 lines / 55,830 bytes, −48% bytes) is accepted over the ceiling; the next cut (moving single-command sections out of coordinator.md) is PRD-032 R6.
- R7 Rule inventory. Before any rewrite, the builder writes `docs/prd/PRD-031-rule-inventory.md`. It lists every rule in the four files, one row each: source file and rule, then `kept`, `merged into <row>`, or `removed` with the reason. A rule removed under Part C cites its rule number. After delivery, every row has a final status, and the counts of `kept`, `merged`, and `removed` appear at the end of the file.
- R8 These safeguards stay in force, in the four files or in a file the inventory names:
  - Not-run honesty: a check that did not run is `not-run`, never `passed`.
  - Authorization: plan approval does not authorize commit, push, publish, or deploy.
  - Write set: a specialist writes only inside its assigned set.
  - Independent review: a builder claim is not review, and review is read-only.
  - Memory ownership: only the Coordinator writes `MEMORY.md`, `EPOCHAL.md`, and `RISKS.md`.
  - Anti-delegation clause, verbatim in `coordinator.md` and `external-clis`, with its sha256 pin.
  - RISK-001 (the clause sits in the prompt text) and RISK-002 (no write to the tree or harness dir during an external CLI call).
  Verifiable: `go run ./cmd/dh validate` passes, including `checkAntiDelegationClause`, and the inventory has a `kept` or `merged` row for each item.
- R9 No new agent, skill, command, or `project.yaml` key. The counts stay: 13 agents, 22 skills, 19 + 1 commands. There is no `commands.quality` key.
- R10 Cleanup premise. The same pull request updates every document that asserts a behavior this PRD changes (section 4). `dist/` is regenerated with `go run ./cmd/dh build`. A search for each removed term (`F<round>`, `reopened`, `rotation`, `cursor`, `Churn`, `two consecutive`, and any other term the owner accepts to remove) finds no remaining mention outside `docs/prd/`, `docs/adr/`, `CHANGELOG.md` history, and `evals/baselines/`. Verifiable: `go run ./cmd/dh validate` and a clean `git diff` of `dist/` after `go run ./cmd/dh build`.

### Part C · Behavior removals (accepted by the owner 2026-10-30)

Each rule gives the cost today and the risk that returns. R16 and R17 were proposed and rejected (see Out of scope). The six-round cap and the `MEMORY.md` bundle logging stay.

- R11 One clean round ends the code-review loop. The loop stops after the first clean round, not after two consecutive clean rounds. A round stays "clean" under today's definition: no blocking finding is `open` or `escalated`, and every dispatched reviewer returned a verdict. The final feature review follows the same rule.
  - Cost today: one extra full-change review for every ticket and for the final review. PRD-029 (0.34.0) added it.
  - Risk if removed: one reviewer miss stands. Deterministic checks and the independent review remain, and the owner can ask for the full panel.
  - Recommendation: accept.
- R12 The code-stage rotation cursor is removed. A code entry and every re-entry use the first entry of `reviewers.code` (`claude` when listed). The owner asks for the full panel when wanted. The Coordinator no longer writes a cursor in `MEMORY.md`.
  - Cost today: about 11 lines of rule, one `MEMORY.md` field per task, and the "not-run entry" retry logic. The confirming round (R11) is the main reason for it.
  - Risk if removed: a re-entry no longer changes vendor, so a blind spot of one reviewer repeats (ADR-004 purpose).
  - Recommendation: accept. It depends on R11. If R11 is rejected, keep R12's current behavior.
- R13 Finding states shrink from six to three: `open`, `closed`, `escalated`. A fix claimed by the builder stays `open` until a reviewer confirms it. `closed` covers confirmed-fixed and settled with evidence (`intentional` or `noise`). Finding ids `F<round>-<n>` are removed. A finding is identified by its location and title.
  - Cost today: a state machine with six exits, an id scheme, and one `MEMORY.md` entry per finding.
  - Risk if removed: a reopened finding is harder to link to its first report, and the settled list in a re-entry is less exact.
  - Recommendation: accept.
- R14 The In-flight doubt allows one cycle, not three. The Coordinator classifies findings in two classes, actionable or not. If the review returns findings and the Coordinator classifies none as actionable, the Coordinator escalates to the owner. The trigger "open decision with structural option O2" is removed. Only `Doubt: yes` slices trigger it.
  - Cost today: up to three fresh-context reviews before the end-of-slice review, a four-way classification, and the `Doubt` line in every slice.
  - Risk if removed: a contract misread or a trust-boundary error is caught one cycle later, at the end-of-slice review.
  - Recommendation: accept. This is the narrowing R2 requires. Removing the doubt entirely is not recommended, because test-first does not catch a misread contract.
- R15 The `Churn:` line and the three-option churn protocol are replaced by one sentence. When a new finding reverses an earlier accepted fix, the loop stops and the Coordinator asks the owner one numbered question.
  - Cost today: a reviewer reads the delta on every re-entry to write a `Churn:` value for every finding, and a separate decision protocol follows.
  - Risk if removed in full: fix and revert oscillation runs until the round cap.
  - Recommendation: replace with the one sentence, not full removal.

## 4. Docs

Documents that assert behavior this PRD changes. The builder confirms each with a search and edits it in the same pull request. Pairs update together (pt-br and `docs/en/`).

- `AGENTS.md` (bullets "Verificação padrão", "PRD e tickets", "Revisão por CLI externa", and "Agentes")
- `.agents/coordinator.md`
- `.agents/builder.md`
- `.agents/reviewer.md`
- `.agents/qa-verifier.md`
- `.commands/build.md`
- `.commands/plan-loop.md`
- `.skills/incremental-implementation/SKILL.md`
- `.skills/implementation-planning/SKILL.md`
- `.skills/code-review/SKILL.md`
- `.skills/regression-testing/SKILL.md`
- `.skills/delivery-readiness/SKILL.md`
- `.skills/harness-authoring/SKILL.md`
- `.skills/external-clis/SKILL.md`
- `templates/en/TASK.md`
- `templates/pt-br/TASK.md`
- `docs/tutorial.html`
- `docs/en/tutorial.html`
- `docs/novo-fluxo-prd.md`
- `docs/en/new-prd-flow.md`
- `docs/roadmap.html`
- `docs/en/roadmap.html`
- `CHANGELOG.md` (new entry, version bump, and the counts line)
- `docs/adr/ADR-004-duvida-e-rotacao.md` (status note: superseded in part by the new ADR)
- `docs/adr/ADR-009-lean-builder-flow.md` (new: records the accepted removals)
- `docs/prd/PRD-031-rule-inventory.md` (new: R7)

The builder must also search `adapters/` and `internal/` for code that reads a removed `MEMORY.md` field (cursor, finding states, ids). A match is in scope and triggers the clean-clone proof in `AGENTS.md`.

## Appendix

### Out of scope

- Per-action log: adds bookkeeping, the opposite of this PRD.
- `crapper`, `dryer`, and `mutator` tools and thresholds: new tooling the kit does not run.
- Pushing to `main`: conflicts with the branch and review rules in `AGENTS.md`.
- Unguarded worktree cleanup: the Cleanup rules in `coordinator.md` protect ignored files.
- "Build the Lever", "Redesign From First Principles", "Attack the Premise": unbounded scope that no acceptance can test.
- A new `commands.quality` key: R9 forbids new keys.
- Iterative planning and visual development artifacts: the owner calls them a future direction ("norte"), and they do not block this PRD.
- R16, round cap of 3 for the build/QA and code-review loop: rejected by owner 2026-10-30. The cap of six rounds stays.
- R17, no `MEMORY.md` logging of review-bundle elements: rejected by owner 2026-10-30. The Dispatches logging stays.
- Compression of files other than the four in section 1, and removal of document-stage review rules.

### Decisions (closed by the owner, 2026-10-30)

- D1 The commit-content rule lands in `.skills/delivery-readiness/SKILL.md` (R4). Closed.
- D2 The budget is a soft ceiling (R6): target 250 lines and 30000 bytes, ceiling 330 lines and 43000 bytes, with about 10% margin when the inventory justifies it. Closed.
- D3 R11 to R15 accepted. R16 and R17 rejected. Closed.
- D4 Version 0.36.0 (minor bump from 0.35.0). Closed.

## Executed summary

- **Delivered**: The builder now writes a failing test first and the review loop stops after one clean round, with the rotation cursor, the six finding states and the three-cycle doubt removed (R11 to R15), but the four builder-phase files total 432 lines and 55,830 bytes, above the R6 ceiling that the owner accepted on 2026-10-31.
- **Rules**: Repository paths are relative to the repository root; harness paths are prefixed `<harness dir>/`. QA matrix `<harness dir>/tasks/T-3111/TASK.md` (section 7) and `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/` cover the rules marked Honored.
  - R1 Honored. `.skills/incremental-implementation/SKILL.md` (Checks line `Red:`, line 49); `.skills/regression-testing/SKILL.md` (`test-first`). Evidence: QA matrix and `qa-T-3111/`.
  - R2 Honored. `.agents/coordinator.md` no longer says "disproof"; the doubt applies only to `Doubt: yes`. Evidence: QA matrix (`disproof` absent).
  - R3 Honored. `.skills/regression-testing/SKILL.md` (lines 47, 53, 54). Evidence: QA matrix.
  - R4 Honored. `.skills/delivery-readiness/SKILL.md` (line 58), the only body rule. Evidence: QA matrix.
  - R5 Honored. `.skills/harness-authoring/SKILL.md` (line 58). Evidence: QA matrix.
  - R6 Not honored. Measured 432 lines / 55,830 bytes against the ceiling of 330 / 43,000 (the four files, `wc -lc`). Reason: the rewrite cut about 48% of the bytes but not to the ceiling. Decision: owner accepted the excess on 2026-10-31 (recorded under R6); the next cut (moving single-command sections out of `coordinator.md`) is PRD-032 R6.
  - R7 Honored. `docs/prd/PRD-031-rule-inventory.md` (every row has a final status, with counts at the end). Evidence: QA matrix and `qa-T-3111/state.txt`.
  - R8 Honored. Safeguards kept in the four files or in the inventory; `go run ./cmd/dh validate --source-only .` passes, including `checkAntiDelegationClause`. Evidence: `qa-T-3111/validate-source.txt`, QA matrix.
  - R9 Honored. Counts stay at 13 agents, 22 skills, 19 + 1 commands; no new key. Evidence: `qa-T-3111/validate-source.txt`, QA matrix.
  - R10 Honored. Search for removed terms finds no remaining mention outside the excluded paths. Evidence: `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/r10.txt`, QA matrix. The `dist/` regeneration and clean-clone proof ran in T-3113 and T-3114 (`<harness dir>/tasks/PRD-031/evidence/qa-T-3114/`).
  - R11 Honored. One clean round ends the code-review loop. Evidence: QA matrix; code review converged in `<harness dir>/tasks/PRD-031/evidence/review-batch-r1..r3/`; final review converged in `<harness dir>/tasks/PRD-031/evidence/review-final-r1..r3/`.
  - R12 Honored. No rotation cursor; re-entries use the first entry of `reviewers.code`. Evidence: QA matrix.
  - R13 Honored. Three finding states (`open`, `closed`, `escalated`); ids `F<round>-<n>` removed. Evidence: QA matrix, `qa-T-3111/r10.txt`.
  - R14 Honored. One doubt cycle; only `Doubt: yes` slices trigger it. Evidence: QA matrix.
  - R15 Honored. The `Churn:` line and the three-option protocol are replaced by one sentence. Evidence: QA matrix.
  - R16 Not implemented. Rejected by owner 2026-10-30 (round cap of 3); the cap of six rounds stays.
  - R17 Not implemented. Rejected by owner 2026-10-30; `MEMORY.md` bundle logging stays.
- **Tickets**: All verdict `done` (`concluída` in each `TASK.md`). Evidence for each is the Resultado section of `<harness dir>/tasks/T-31NN/TASK.md`.
  - T-3101 to T-3110, lane `backend` (`harness-maintainer`, skill `harness-authoring`).
  - T-3111, lane `teste` (unit), `qa-verifier`; evidence `<harness dir>/tasks/PRD-031/evidence/qa-T-3111/`.
  - T-3112, lane `backend` (`docs-guide`, skill `doc-template-html`).
  - T-3113, lane `infra` (`builder`, mode build); regenerated `dist/` and the manifest.
  - T-3114, lane `teste` (integration), `qa-verifier`; evidence `<harness dir>/tasks/PRD-031/evidence/qa-T-3114/`.
- **Docs**: `AGENTS.md`; `docs/adr/ADR-009-lean-builder-flow.md`; `docs/adr/ADR-004-duvida-e-rotacao.md`; `docs/tutorial.html`; `docs/en/tutorial.html`; `docs/novo-fluxo-prd.md`; `docs/en/new-prd-flow.md`; `docs/roadmap.html`; `docs/en/roadmap.html`; `CHANGELOG.md`; `docs/prd/PRD-031-rule-inventory.md`.
- **Out of scope**: Carried-over Minor findings for PRD-032:
  - `code-review` lost the explicit question "would a plausible regression fail a test?".
  - `incremental-implementation` official-page rule lost three safeguards: skip the fetch when there is one viable option; a local-convention conflict goes to the open-decision gate; the URL goes in the report, not in a code comment.
  - The rule inventory total is 55,808 bytes, while other documents say 55,830.
  - `regression-testing` step 13 has four sentences.
