# Review: PRD-005 — round 1 (merged)

Evaluated state: the first draft of `docs/prd/PRD-005-planejamento-em-loop.md`. Source: `.harness/tasks/UND-003/GRILL.md` (G1–G17). Round cap: 6, owner decision of 2026-09-25.

| Reviewer | Verdict |
|---|---|
| `claude` (document-validator, opus requested, effective unverified) | changes required: 5 blocking, 2 non-blocking |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | changes required: 4 blocking, 1 non-blocking |
| `cli:grok/grok-4.7` | changes required: 6 blocking, 2 non-blocking |

## Verdict
- changes required

## Findings (merged). Each item cites its reporters. The disposition follows the owner's answers in grill round 4 (G18–G23).

- R1: convergence is described as "zero blockers" and "without reservations", but the owner decided unanimous approval (G9). The PRD also does not say what happens when a reviewer is not-run.
  - Reported by: claude (P1), cli:codex/gpt-5.6-luna (P2).
  - Disposition: G9 + G22. The loop converges when every reviewer that ran returns `approved` for the winning plan and at least one of them is from another vendor. Non-blocking findings do not prevent convergence.
- R2: the no-CLI fallback is described as the unchanged `/plan` with a cap of 6, but the kit text still says 2.
  - Reported by: cli:codex/gpt-5.6-luna (P1).
  - Disposition: the kit-wide 6-cap batch becomes an explicit prerequisite of this PRD's implementation. The fallback runs the post-migration `/plan`.
- R3: a tie in blocking findings makes AC-07 untestable.
  - Reported by: cli:codex/gpt-5.6-luna (P3).
  - Disposition: G18. The plan with fewer total findings wins. If that also ties, Claude's plan wins.
- R4: the PRD does not say that it authorizes no implementation, or that implementation is third in the queue.
  - Reported by: cli:codex/gpt-5.6-luna (P4).
  - Disposition: add both as a constraint.
- R5: the CLI planner pool is treated as decided (`reviewers.document`) and is also listed as pending. The case of 3 or more CLIs is undefined.
  - Reported by: claude (P3), cli:grok/grok-4.7 (P6, P7), cli:codex/gpt-5.6-luna (P5).
  - Disposition: G19. Reuse `reviewers.document`, and cycle through the CLIs in list order. Remove the "reviewer-only CLI" clause.
- R6: the PRD does not say what happens to a wave whose CLI plan was discarded.
  - Reported by: claude (P4), cli:grok/grok-4.7 (P5).
  - Disposition: G20. The turn passes to the next CLI in the same wave. If none is left, the wave continues with Claude's plan alone: no comparison, no stall, and the wave counts toward the cap.
- R7: the dispatch-cap check timing is ambiguous, and so is whether failed dispatches count.
  - Reported by: claude (P2).
  - Disposition: G21. The cap is checked before opening a wave, an opened wave always finishes, and failed dispatches count.
- R8: the planner prompt body does not carry the read-only / no-write instruction or the anti-delegation clause (RISK-001), and its source is undefined.
  - Reported by: claude (P5).
  - Disposition: the prompt is built from the `implementation-planner` body, the `implementation-planning` text and the wave context. Its body carries three things verbatim: the no-write instruction, an anti-delegation clause worded for planning (exact wording fixed in the implementation batch), and the reading-is-allowed block. Extend AC-05 to check that the prompt contains them.
- R9: blind writing is attributed to the owner (G2), but it was a design choice.
  - Reported by: claude (P6, non-blocking), cli:grok/grok-4.7 (P4).
  - Disposition: G23 now makes it an owner decision. Wave 1 is blind; later waves refine from the winner without seeing the other's new draft. Cite G23.
- R10: "teto atingido" is applied inconsistently. G11 is only about ending at a cap without convergence.
  - Reported by: cli:grok/grok-4.7 (P3).
  - Disposition: the "teto atingido" line appears only when the loop stops at the wave cap or the dispatch cap without converging.
- R11: AC-01 says `PLAN.md` and its `TASK.md` files are the only artifacts, which conflicts with the per-wave persistence (G12), the dollar estimate (G8) and the open-blockers list (G11).
  - Reported by: cli:grok/grok-4.7 (P8).
  - Disposition: list every artifact. Per wave, under `wave-N/`: both plans, both reviews and the winner record. Final: the winning `PLAN.md`/`TASK.md` and a report with the stop reason, the dispatch count, the dollar estimate and any open blockers.
- Non-blocking:
  - The eval-coverage target in Sinal 1 goes beyond the grill. Reported by: cli:grok/grok-4.7 (P1). Disposition: soften it to an implementation-plan concern.
  - A malformed table row. Reported by: cli:grok/grok-4.7 (P2). Disposition: fix it.
  - The basis of the dollar estimate. Reported by: claude (P7). Disposition: leave the basis to the implementation batch and say so.

# Review: PRD-005 — round 2 (merged)

Evaluated state: the correction 1 draft of the PRD. Source: GRILL.md G1–G23.

| Reviewer | Verdict |
|---|---|
| `claude` (document-validator, opus requested, effective unverified) | changes required: R1–R11 applied; 3 blocking (S1–S3), 5 non-blocking (N1–N5) |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | changes required: R1–R11 not reopened; 4 blocking (S1–S4) |
| `cli:grok/grok-4.7` | not-run: the state was superseded by correction 2 before its wave ran |

## Verdict
- changes required

## Findings (merged). The Coordinator's dispositions derive from G19–G21.

- T1: the rotation rule is undefined in two cases: when a configured CLI does not resolve on this machine, and when a hand-off inside a wave leaves the next wave's CLI unclear.
  - Reported by: claude (S1), cli:codex/gpt-5.6-luna (S2).
  - Disposition:
    - The rotation list is the resolvable `cli:` entries of `reviewers.document`, in list order, taken from the check at loop start (RF-02).
    - Wave N's primary CLI is entry ((N−1) mod k). A hand-off inside a wave (RF-15) tries the following entries in order. It does not move the cursor: the next wave's primary CLI still depends only on its wave number.
    - An entry that stops resolving on its turn is handled as in RF-15.
    - Add ACs for an unresolvable configured entry and for the cursor after a hand-off.
- T2: the PRD never says how to decide whether a wave "fits" under the dispatch cap.
  - Reported by: claude (S2), cli:codex/gpt-5.6-luna (S1).
  - Disposition:
    - A wave fits when the remaining cap ≥ 2 + 2 × R, where R is the number of reviewers configured in `reviewers.document`: 2 planners plus each reviewer judging both plans.
    - Hand-off attempts beyond that reservation still run, because an opened wave always finishes. They count toward the cap.
    - Put a worked example in AC-11 (for R = 3: a wave reserves 8).
- T3: artifacts. AC-01 conflicts with AC-03, and the fallback persistence is incomplete.
  - Reported by: claude (S3), cli:codex/gpt-5.6-luna (S4).
  - Disposition:
    - AC-01 applies to the loop path only, meaning at least one resolvable CLI.
    - The fallback follows the complete current `/plan` persistence contract: `PLAN.md`, per-slice `TASK.md`, ADRs when applicable, `PLAN.review.md`, and memory updates by the Coordinator.
    - Loop mode also writes a base-level `PLAN.review.md`: the merged review of the final winning plan, so `/plan` consumers still find it. RF-12 and AC-01 must say so.
- T4: the status line says "rodada 1 de 2", and §9 says the 2-round cap covers "inclusive deste PRD".
  - Reported by: cli:codex/gpt-5.6-luna (S3), claude (N2).
  - Disposition: write "rodada 2 de 6", and remove "inclusive deste PRD".
- Non-blocking:
  - N1: §6 says the wave cap inherits the kit-wide cap. Disposition: state that `wave_cap` has its own default of 6.
  - N3: AC-05's last clause contradicts itself. Disposition: reword it as "o Coordenador é o único gravador do plano".
  - N4: the PRD gives no test for "conteúdo de plano reconhecível". Disposition: the minimum signal is the required sections of the `implementation-planning` format (goal and slices). The exact check is deferred to implementation.
  - N5: the PRD labels decisions as "textual" quotes. Disposition: write "decisão do dono, Gxx (registro do grill)".
  - All non-blocking items reported by claude.

## Round 2: correction 2 dispatched to product-discovery.

# Review: PRD-005 — round 3 (merged)

Evaluated state: the correction 2 draft.

| Reviewer | Verdict |
|---|---|
| `claude` (document-validator, opus requested, effective unverified) | changes required: T1–T4 and N-ids applied, no R-id regressed; 3 blocking (U1–U3), 1 non-blocking (U4) |
| `cli:grok/grok-4.7` | changes required: 2 blocking plus R11 partial; 2 non-blocking |
| `cli:codex/gpt-5.6-luna` | not-run: state superseded by correction 3 before its wave |

## Verdict
- changes required

## Findings (merged)
The Coordinator decided every disposition below, deriving each from the named owner decisions.

- V1 — Hand-off does not say whether the scan wraps past the last rotation entry. Reported by: claude (U1), cli:grok/grok-4.7 (U2).
  - Disposition (G20, "não trava"): the hand-off is circular. It starts at the entry after the wave's primary and tries each of the k entries at most once per wave. The wave runs without a comparison only when all k entries failed in it.
  - Assert this in AC-17.
- V2 — "Melhor vencedor até aquele ponto" (G11) could mean the last wave's winner or the best winner across waves. Reported by: claude (U2).
  - Disposition (G11 + G18): the best plan so far is the winner with the fewest blocking findings across all waves run. A tie goes to fewer total findings, and a remaining tie goes to the most recent wave.
  - RF-12 writes that plan as `PLAN.md`, and its merged review becomes the base `PLAN.review.md`.
  - Assert this in AC-12.
- V3 — The loop path does not record ADRs or MEMORY updates, and RF-01 still describes the output as only `PLAN.md` and `TASK.md`. Reported by: claude (U3), cli:grok/grok-4.7 (R11 partial).
  - Disposition: the loop path also writes ADRs when applicable, and the Coordinator updates `.harness/MEMORY.md`, as in `/dev-harness:plan`.
  - RF-01 describes the final artifact as the full `/plan` contract plus the loop's extra artifacts, with no conflicting "only".
  - AC-01 lists them.
- V4 — A wave could be read as one scoring pass or as a nested `/plan` correction cycle. Reported by: cli:grok/grok-4.7 (U1).
  - Disposition (G5: each wave is one review round): each candidate is reviewed once per wave. No planner is re-dispatched inside a wave, and open blockers feed the next wave (RF-08). The `/plan` correction cycle applies only to the no-CLI fallback.
  - Add a criterion: a `changes required` verdict does not re-dispatch the planner within the same wave, and a wave makes 2×R review dispatches, plus any RF-15 hand-off.
- Non-blocking:
  - The "Não entra" bullet about more than two CLIs is reworded as claude suggested. Reported by: claude (U4).
  - The §8 risk wrongly says `document-review` has `type: llm` graders. Correct it or drop it. Reported by: cli:grok/grok-4.7 (U3).
  - The reference that credits G18–G23 with extending the CLI's role should cite G6 and G15 instead. Reported by: cli:grok/grok-4.7 (U4).

## Round 3 — correction 3 dispatched to product-discovery.

# Review: PRD-005 — round 4 (merged)

Evaluated state: the correction 3 draft.

| Reviewer | Verdict |
|---|---|
| `claude` (document-validator, opus requested, effective unverified) | approved. V1–V4 applied, no regression. 3 non-blocking findings (W1–W3). |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | changes required. 1 blocking finding. |
| `cli:grok/grok-4.7` | not-run. Its wave never ran because correction 4 superseded the state first. |

## Verdict
- changes required (one reviewer blocks)

## Findings
- X1: the RF-15 wording lets a reader retry the failed primary CLI within the same wave, including when it is the only configured CLI.
  - Reported by: cli:codex/gpt-5.6-luna (W1).
  - claude read the same text as correct ("counts the primary's own attempt"). The meaning is settled; this is ambiguous wording, not a disagreement about behavior.
  - Disposition (G20): the hand-off tries each other resolvable CLI at most once and never retries the wave's failed primary within that wave. With k = 1, a failure goes straight to a single-candidate wave. Mirror this in AC-17.
- Non-blocking (claude):
  - W1: move the base `PLAN.review.md` in AC-01 into the list of files that make up the `/plan` contract.
  - W2: in RF-12 and AC-14, say the Coordinator updates `.harness/MEMORY.md` outside the base directory.
  - W3: in RF-06, a two-candidate wave spends 2×R reviews and a single-candidate wave spends R, plus any hand-off.

## Round 4: correction 4 dispatched to product-discovery.

# Review: PRD-005 — round 5 (merged)

Evaluated state: the correction 4 draft.

| Reviewer | Verdict |
|---|---|
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | approved: X1 and W1–W3 applied |
| `cli:grok/grok-4.7` | changes required: 1 blocking finding |
| `claude` (document-validator) | not-run: correction 5 superseded the state before its wave; it approved round 4 |

## Verdict
- changes required

## Findings
- Z1: RF-01, RF-12 and AC-14 place ADRs, and in RF-01 `.harness/MEMORY.md`, inside `.harness/tasks/<id>/`. `/dev-harness:plan` writes ADRs to `.harness/adr/ADR-<n>.md` and memory to `.harness/MEMORY.md`, both outside that directory. The fallback (AC-03) would therefore write them to a different place than the loop path.
  - Reported by: cli:grok/grok-4.7 (Y1). Confirmed against `.commands/plan.md` lines 9 and 22.
  - Disposition: "sob `.harness/tasks/<id>/`" applies only to `PLAN.md`, the per-slice `TASK.md` files, the base `PLAN.review.md`, `wave-N/` and the final report. ADRs, when applicable, go to `.harness/adr/ADR-<n>.md`, and the Coordinator updates `.harness/MEMORY.md`. Both paths are the ones `/dev-harness:plan` already uses.

## Round 5: correction 5 dispatched. Round 6 is the last round under the cap of 6.

# Review: PRD-005 — round 6 (merged, final round: cap of 6 reached)

Evaluated state: the correction 5 draft.

| Reviewer | Verdict |
|---|---|
| `claude` (document-validator, opus requested, effective unverified) | approved. Z1 applied, no regression, no new finding. |
| `cli:grok/grok-4.7` | approved |
| `cli:codex/gpt-5.6-luna` (`-s read-only`) | changes required. Z1 applied; 1 new blocking finding (AA1). |

## Verdict
- changes required. One reviewer blocks, and the round cap is reached. Under the kit rule the document goes to the owner with the open point listed below.

## Open point
- AA1: the pool of candidate CLIs does not require a vendor other than Claude.
  - Reported by: cli:codex/gpt-5.6-luna.
  - Why it matters: G6 says the competing plan comes from an external CLI "from another vendor". `external-clis` also supports `claude` as a CLI binary, so a configured `cli:claude/<slug>` entry would let Claude compete against itself.
  - Coordinator recommendation: exclude `cli:claude/...` entries from the planner rotation list in RF-05 and AC-06. If only same-vendor entries resolve, use the no-loop fallback (RF-02). This is a one-sentence change. It is left to the owner because the round cap is reached.
