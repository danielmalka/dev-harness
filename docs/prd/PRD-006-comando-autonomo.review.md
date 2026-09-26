# Review: PRD-006 — comando autônomo `/dh:auto`

| Field | Value |
|---|---|
| Document | docs/prd/PRD-006-comando-autonomo.md |
| Source | .harness/tasks/UND-004/GRILL.md; CLAUDE.md decisions of 2026-09-25 |
| Reviewers | claude (document-validator, opus); cli:codex/gpt-5.6-luna; cli:grok/grok-4.7 |

## Round 1 of 6 — 2026-09-25

| Reviewer | Verdict |
|---|---|
| claude | changes required (7 blocking, 4 non-blocking) |
| cli:codex/gpt-5.6-luna | changes required (4 blocking) |
| cli:grok/grok-4.7 | changes required (12 blocking, 1 non-blocking) |
| Merged | changes required |

Workspace freeze: identical after all three calls.

- **A1** [blocking] `/dh:resume` must read `Stage:`, but `resume.md` does not and AC-16 forbids editing it. — claude P1, codex P2, grok P2. Disposition: `.commands/resume.md` (or `context-handoff`) reads `Stage:` from `AUTO-<n>/BRIEF.md` when MEMORY names an `AUTO-<n>` task; add it to the allowed edits.
- **A2** [blocking] AC-16 contradicts docs-in-PR, eval cases, the re-embedded coordinator copies in eval prompts, regenerated `dist/`, and runtime `AUTO-<n>/BRIEF.md`. — claude P2, codex P4, grok P1. Disposition: rewrite as claude suggested, and add re-embeds, runtime artifacts, and `resume.md` (A1).
- **A3** [blocking] AC-14/AC-15 and "Não entra" forbid the RF-11 reviewer reduction. — claude P3. Disposition: "with RF-11 applied; no other reduction".
- **A4** [blocking] No step executes commit, push, PR, merge or tag, in no set order; the standing git rules are missing. — claude P4, codex P3, grok P3. Disposition: new RF+AC. After release reports ready, the Coordinator does only what (c)-(f) authorized, in this order: commit → proof the project requires (clean clone for `dist/`) → push → PR → merge only with green CI + APPROVED review + non-main branch → tag on the merge commit. Never `--force` or `--no-verify`; CI off or broken is not green. An unauthorized item ends the run with state recorded. RF-13 PR/merge steps only apply when authorized.
- **A5** [blocking] How stages are chained is not stated; `release` has `disable-model-invocation`. — claude P5, grok P11 (part). Disposition: the Coordinator follows each `.commands/<stage>.md` Routing/Prerequisites/Output/Limits inline in the same session, including `release`; no slash invocation.
- **A6** [blocking] Mode (b)/(c) questions vs RF-03 "no further questions" and their mechanics. — claude P6, codex P1. Disposition: they are the only exception besides the RF-06 stops. The question is asked in chat after the state is recorded (Stage + MEMORY), so `/dh:resume` works if the session ends. AC-06 is reworded.
- **A7** [blocking] "Confirmed vulnerability" is undefined against the kit's classes. — claude P7, grok P10. Disposition: a demonstrated vulnerability or a confirmed exposure stops the run; a hypothesis does not.
- **A8** [blocking] Detection of an approved PRD or validated plan uses the wrong signal. — claude P8 (non-blocking), grok P4, P5. Merged as blocking. Disposition: an approved PRD is one whose table row `Status` reads `aprovado`/`approved`; a validated plan has a `PLAN.review.md` next to it whose last round verdict is `approved`; anything else proposes `discover`.
- **A9** [blocking] The RF-12 secure trigger paraphrases `secure.md` narrower. — grok P6. Disposition: cite `.commands/secure.md`'s own surface (authentication, data, input, execution) verbatim; no own list.
- **A10** [blocking] RF-11 (iii) drops "plano de". — grok P7. Disposition: "o plano de um lote mecânico passa por um revisor".
- **A11** [blocking] "Achado menor" is undefined. — grok P8. Disposition: a finding the stage's review skill classifies as non-blocking (Minor/Cosmetic in `code-review`, non-blocking in `document-review`/`security-review`).
- **A12** [blocking] "Missed eval bar" is undefined. — grok P9. Disposition: the pass bar the plan fixes per case before any paid run; a miss is a grader score in the result JSON below that bar.
- **A13** [blocking] The RF-14 gate claims to be existing kit practice. — grok P12. Disposition: state it as this command's rule. Estimated judge cost = the judge/agent ratio from the case's last recorded baseline × `--max-cost-usd`, or `--max-cost-usd` itself when no baseline exists.
- **A14** [non-blocking] `<base-branch>` and the PR target are undefined. — claude P9. Disposition: a feature branch created from `main` at start; per-slice PRs target `main`.
- **A15** [non-blocking] "docs-guide" named as a command; verify/secure are not invoked by stages. — claude P10, grok P11. Disposition: `/dh:secure` per RF-12; the `docs-guide` agent at the end of the batch; drop `verify`.
- **A16** [non-blocking] The questionnaire has 8 items (7 + start-stage confirmation). — claude P11. Accepted.
- **A17** [non-blocking] The AC-17 "why" sentence in MEMORY is extra. — grok P13. Accepted: drop it.

Coordinator dispositions A1, A4–A15 are decided on the owner's behalf under the session authorization of 2026-09-25.

## Round 2 of 6 — 2026-09-25

| Reviewer | Verdict |
|---|---|
| claude | changes required (1 blocking, 4 non-blocking); A1–A17 confirmed applied |
| cli:codex/gpt-5.6-luna | changes required (1 blocking) |
| cli:grok/grok-4.7 | changes required (3 blocking, 1 non-blocking) |
| Merged | changes required |

Workspace freeze: identical after all three calls.

- **B1** [blocking] `Stage:` meaning undefined. — codex P1. Disposition: `Stage:` always names the next unit (stage or slice) to execute; completing a unit advances it; a stop records the interrupted unit, not marked complete.
- **B2** [blocking] Branch-per-slice needs per-slice commits, which RF-16's order forbids; branch point undefined. — claude P1. Disposition: in branch-per-slice mode with (c) authorized, each slice is committed on its own branch once its build/QA/review gate is green (not a violation of RF-16). Slice N+1 branches from slice N (stacked). Push → PR → merge in slice order → tag run after release. Branch-per-slice without (c) is flagged as incompatible in the questionnaire, which proposes single branch.
- **B3** [blocking] RF-16 step 2 says the proof runs "before committing dist"; the candidate is the local commit. — grok A4, claude P5 (non-blocking). Merged blocking. Disposition: the step (1) commit is the candidate. When `dist/` was regenerated, run the clean-clone sha256 proof on that commit; on failure amend only while unpushed, never `--force`; push only after the proof passes. Per-slice PR steps apply only with (d), merge steps only with (e).
- **B4** [blocking] RF-09 and section 6 attribute a commit/push gate to CLAUDE.md that it does not state. — grok A18. Disposition: commit and push require questionnaire (c) plus the B3 proof when `dist/` changed. The green CI + APPROVED + non-main gate applies to merge (RF-16 step 5). Remove the misattribution.
- **B5** [blocking] The position of `/dh:secure` in the chain is undefined. — grok A19. Disposition: per slice, secure runs right after that slice's build gates and before the next slice when the slice touches the `.commands/secure.md` surface; otherwise one skip is recorded. The batch `review` runs after the slice loop and before `release`.
- **B6** [non-blocking] CI not green at merge. — claude P2. Disposition: a pending CI may be awaited; red, off or broken CI ends the run with state recorded.
- **B7** [non-blocking] Status detection: the value begins with `aprovado`/`approved`, case-insensitive. — claude P3. Accepted.
- **B8** [non-blocking] Sinal 3 points at section 6 for proportionality cases. — claude P4. Disposition: "observable in the eval cases the plan defines (section 6)".
- **B9** [non-blocking] RF-15 "ou context-handoff" vs AC-16. — grok A20. Disposition: `.commands/resume.md` only.

All dispositions decided by the Coordinator on the owner's behalf (session authorization of 2026-09-25).

## Round 3 of 6 — 2026-09-25

| Reviewer | Verdict |
|---|---|
| claude | changes required (1 blocking, 3 non-blocking); A1–A17, B1–B9 confirmed |
| cli:codex/gpt-5.6-luna | changes required (1 blocking) |
| cli:grok/grok-4.7 | changes required (7 blocking) |
| Merged | changes required |

Workspace freeze: identical after all three calls.

- **D1** [blocking] AC-01 limits the chain to `.agents/` agents; configured `reviewers.*` CLI entries are not agents. — codex P1. Disposition: `/dh:auto` adds no new agent, skill or command beyond `auto.md`, and keeps every configured reviewer, including `cli:` entries.
- **D2** [blocking] The AC-16 exception list omits the `append_system_prompt` coordinator copies in existing `evals/cases/*/prompt.md`. — grok A2. Accepted: name them.
- **D3** [blocking] AC-24 says slices branch from `main`; RF-13/AC-41 say stacked. — grok C1. Disposition: `<base-branch>` is the feature branch from `main`; slice 1 branches from it; slice N+1 from slice N; every PR targets `main`.
- **D4** [blocking] The clean-clone proof runs only on the last slice commit. — claude C1, grok C2. Disposition: every commit that regenerates `dist/` is its own candidate (the batch commit, or each slice commit). The proof runs on the tip of that branch before its push; amend only while that branch is unpushed; never `--force`.
- **D5** [blocking] End-of-batch docs (docs-guide before release) have no commit in branch-per-slice. — grok C3. Disposition: step (1) in branch-per-slice commits, on the last slice branch, whatever is still uncommitted after docs-guide and release, and does not recommit slice files.
- **D6** [blocking] Per-slice order of secure vs. slice commit. — grok C4. Disposition: build gates → `/dh:secure` or the recorded skip → slice commit. An RF-06 stop during secure happens before the commit.
- **D7** [blocking] AC-20 says "only one minor finding". — grok C5. Disposition: any number of findings the stage skill classifies as non-blocking opens no correction round.
- **D8** [blocking] Chained command Prerequisites say "ask the owner", which conflicts with mode (a). — grok C6. Disposition: a chained command's instruction to ask the owner never overrides RF-05. In mode (a) the Coordinator decides and records; in (b)/(c) RF-05/AC-29 apply. Name `.commands/build.md` Prerequisites as covered.
- **D9** [non-blocking] Release returns "not ready". — claude C2. Disposition: the run ends with the reason and `Stage: release` recorded (RF-07), without RF-16.
- **D10** [non-blocking] Stacked-PR merge mechanics. — claude C3. Disposition: per-slice PRs are merged in slice order with a merge commit, never squash, so no rebase or `--force` is needed; single-branch PRs keep the project's usual squash.
- **D11** [non-blocking] Section 9 cites only round 1. — claude C4. Accepted: rounds 1–3, A/B/D items.

All dispositions decided by the Coordinator on the owner's behalf (session authorization of 2026-09-25).

## Round 4 of 6 — 2026-09-25

| Reviewer | Verdict |
|---|---|
| claude | approved (2 non-blocking); A/B/D items confirmed |
| cli:codex/gpt-5.6-luna | approved (0 findings) |
| cli:grok/grok-4.7 | changes required (4 blocking) |
| Merged | changes required |

Workspace freeze: identical after all three calls.

- **E1** [blocking] AC-41 still commits the slice right after build gates, before secure. — grok D6. Disposition: AC-41 commits only after build gates and secure (or its recorded skip); an RF-06 stop during secure comes before the commit.
- **E2** [blocking] AC-48 lacks a branch-mode guard. — grok E1. Disposition: it opens with "in branch-per-slice mode with (c) authorized". Single-branch mode has no per-slice commit; its only commit is RF-16 step 1 after release.
- **E3** [blocking] AC-23/AC-24 do not condition commit, PR and merge on (c)/(d)/(e). — grok E2. Disposition: the single-branch commit happens only with (c) and the PR only with (d). Per-slice PRs open only with (d) and merge only with (e), in slice order with merge commits. Otherwise the run stops at that RF-16 step (AC-32).
- **E4** [blocking] Mode (a) vs "ask the owner" in a chained command's Routing (`discover.md` Routing) as well as Prerequisites. — grok E3. Disposition: any instruction to ask the owner, in the Routing or Prerequisites of a chained command, yields to RF-05. Name `discover.md` Routing and `build.md` Prerequisites. In mode (a) the Coordinator decides and records; in (b)/(c) the question follows RF-05/AC-29.
- **E5** [non-blocking] The tag target in branch-per-slice. — claude P1. Accepted: the merge commit of the last slice PR on `main`, or the squash commit of the single PR.
- **E6** [non-blocking] The RF-13 phrase about RF-16's scope. — claude P2. Accepted: "que rege o commit de fechamento e push/PR/merge/tag".

## Round 5 of 6 — 2026-09-25

| Reviewer | Verdict |
|---|---|
| claude | approved (2 non-blocking); E1–E6 confirmed |
| cli:codex/gpt-5.6-luna | approved (1 non-blocking) |
| cli:grok/grok-4.7 | changes required (4 blocking) |
| Merged | changes required |

Workspace freeze: identical after all three calls.

- **G1** [blocking] `build.md` Limits and coordinator Procedure 8 say "replan with the owner" at the 6-round cap; RF-07 says stop. — grok F1. Accepted: in RF-07, that text yields. The run records the blocked state with `Stage:` at the interrupted unit and ends; replanning happens only in a later session.
- **G2** [blocking] Coordinator Intake says "wait for approval before the first edit". — grok F2. Accepted: in RF-05, Intake's wait is not an extra gate. In mode (a) a plan `approved` by the validation panel goes to build without a question. In modes (b)/(c) the Coordinator asks only an open product/contract decision (AC-29) and does not wait otherwise.
- **G3** [blocking] "Não entra" vs AC-16. — grok F3. Accepted with grok's wording.
- **G4** [blocking as ambiguity; grok's suggested direction rejected] What "APPROVED review" means at merge. — grok F4. Coordinator disposition: the signal is the independent in-session review verdict `approve` (code-review / `.commands/review.md`, recorded in MEMORY). A GitHub PR review state is not required. Evidence: the owner's standing merge rule names the QA reviewer's APPROVED verdict, and this is a single-owner project with no PR reviewers, so requiring a PR review state would make every autonomous merge impossible. State this explicitly in RF-16 step 5 and section 6 so no builder reads it either way.
- **G5** [non-blocking] The section 8 "sempre" order covers both modes. — claude P1. Accepted: "Em branch por slice … commit; em branch única, build → secure/skip, sem commit por slice".
- **G6** [non-blocking] Section 9 cites rounds 1–3. — claude P2, codex P1. Accepted: rounds 1–5, items A–G.

## Round 6 of 6 (final) — 2026-09-26

| Reviewer | Verdict |
|---|---|
| claude | approved (1 non-blocking: G6 incomplete) |
| cli:codex/gpt-5.6-luna | approved (0 findings) |
| cli:grok/grok-4.7 | 1st call: timeout with no verdict (transport failure, not-run). Retry: `changes required` (2 blocking), read even though the retry also ended with exit 124 after writing its verdict |
| Merged | changes required. The cap is exhausted. |

Workspace freeze: identical after all calls.

- **H1** [blocking] The plan stage of `/dh:auto` would use `plan.md` Routing, which dispatches only `document-validator` and ignores `reviewers.document`. — grok P1. Coordinator disposition: accepted with grok's wording. The plan stage uses the `discover.md` `reviewers.document` procedure; RF-11 (iii) applies.
- **H2** [blocking] A stop inside RF-16 does not say what `Stage:` records. — grok P2. Coordinator disposition: accepted. `Stage:` names the interrupted RF-16 step (`commit`, `proof`, `push`, `pr`, `merge`, `tag`); `/dh:resume` continues there without rerunning release or completed steps.
- **G6** [non-blocking] Section 9 omits E1–E6. — claude, grok. Accepted.

Round cap reached with H1/H2 open. The Coordinator applies H1, H2 and G6 as dispositioned, with no seventh validation round, and hands the PRD to the owner. The report states that these three edits were not re-validated by the panel.
