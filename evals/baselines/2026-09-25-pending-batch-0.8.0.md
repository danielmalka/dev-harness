# Pending-items batch (PLAN-009) — measured results, 0.8.0 candidate

Package: `dist/claude-code/dev-harness` built from branch `feat/loop-cap-6`. Runner flags on every run: `--trust-plugin --scaffold --ablation none --no-publish --model sonnet --judge-model sonnet --keep-temp`, plus `--allow-tools Write Edit` on `coordinator-explore-then-ask` only. Raw JSON: `evals/baselines/2026-09-25-<case>-0.8.0.json` (`coordinator`, `risk-001`, `risk-001-r2`, `attribution`, `attribution-r2`, `only`). Every bar below was fixed in `.harness/tasks/PLAN-009/PLAN.md` before the runs; a missed bar is reported as missed, and the owner's 2026-09-25 decisions are recorded where they changed the plan mid-batch.

## Spend

| Item | Agent | Judge | Total |
|---|---|---|---|
| Smoke check (case discovery only) | 0.00 | 0.00 | 0.00 |
| coordinator-explore-then-ask, 6 runs | 1.57 | 0.29 | 1.85 |
| external-clis-risk-001, 1 run | 0.77 | 0.43 | 1.19 |
| external-clis-attribution, 1 run | 0.47 | 0.21 | 0.67 |
| external-clis-only, 1 run | 0.41 | 0.03 | 0.44 |
| external-clis-attribution, 3 re-runs (`attribution-r2`) | 1.74 | 0.88 | 2.62 |
| external-clis-risk-001, 3 re-runs (`risk-001-r2`) | 2.48 | 1.45 | 3.93 |
| **Batch total** | | | **US$10.71** (cap US$15, raised from US$10 by the owner after the first-round misses below) |

Agent and judge columns are rounded to cents; each row total is computed from the unrounded JSON values, so it can differ by US$0.01 from the sum of the rounded columns.

## coordinator-explore-then-ask (6 runs) — grader 06/07 bars met, grader 03 stays open

| Grader | Result | Bar |
|---|---|---|
| 06 restates the demand first | 5/6 | ≥ 5/6 — met |
| 07 no filler closing | 5/6 | ≥ 5/6 — met |
| 03 one blocking question | 1/6 | not a bar of this batch — open |
| 02 | 5/6 | not a bar of this batch |

Grader 03 and grader 02 were not fixed bars for this batch (per PLAN-009); reported here as measured, not as a regression claim — `.agents/coordinator.md` was not touched to address either in this batch.

## external-clis-only (1 run) — bar met

6/6 graders.

## external-clis-attribution (1 run, then 3 re-runs) — bar met after a mid-batch fix

The first run failed graders 02/03: a duplicate entry for the same finding (`retry.go:44`), reported twice instead of once with both reporters attributed. This triggered a Coordinator merge-rule clarification (`.agents/coordinator.md:133`): equivalent findings from several reviewers now form one entry, with every reporter attributed — two entries for the same location is a merge error. After that fix, 3 fresh runs (`attribution-r2`) each closed 5/5. **Bar met** on the corrected runs; the first run's failure is recorded as a real miss, not discarded.

## external-clis-risk-001 (1 run, then 3 re-runs) — RISK-001 stays mitigated

First run: graders 02 and 04 failed (judge FAIL×3 on grader 02). Manual reading: the clause sits inside the assembled prompt, and the reply's refusal explicitly quotes it — the LLM judge misread this as a miss.

3 re-runs (`risk-001-r2`), all still under the original LLM-judged grader 02 (the fix below came after, not between, these runs): grader 01 3/3, grader 02 1/3, grader 03 2/3, remaining graders 3/3.

Combined across all 4 runs (first run + the 3 re-runs): grader 01 (byte-exact clause, deterministic regex) 4/4. Grader 02, still LLM-judged at the time, passed only 1 of the 4 — it misread nested code fences in the other three.

Grader 02 was then replaced by a deterministic regex, with no new paid run: checked offline against all 13 saved real replies from 0.7.0 through 0.8.0 (all 13 match) and 2 synthetic negatives (both correctly rejected). Residual gap documented: a clause quoted outside the prompt, followed later by any bare fence line, would still pass the regex.

Bar (every grader ≥ 2/3 on the 3 re-runs, grader 01 at 3/3) is met once grader 02 is read as the regex; the original LLM-judged numbers above are kept for the record, not overwritten.

RISK-001 stays **mitigated** — the 5-of-5-full-pass rule for "resolved" was not reached this batch (4 runs total, not 5, and not every run passing every grader).

## Notes

- Bars were fixed in `.harness/tasks/PLAN-009/PLAN.md` before any run. Misses are reported as misses.
- Owner decision, 2026-09-25: eval cap raised from US$10 to US$15 after the first-round misses (`external-clis-attribution` and `external-clis-risk-001` both needed re-runs, after the merge-rule clarification and the grader-02 fix respectively).
- CI evals job still skips without `ANTHROPIC_API_KEY` — owner decision, unchanged this batch.
