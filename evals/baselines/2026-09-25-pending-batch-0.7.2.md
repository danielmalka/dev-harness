# Grader 07 recalibration (PLAN-008): consolidated results for the 0.7.2 candidate

Sources:
- Offline re-judge and live measurement, in detail: `evals/baselines/2026-09-25-grader07-offline-rejudge.md`.
- Raw live run: `evals/baselines/2026-09-25-coordinator-explore-then-ask-0.7.2.json`.

## What changed
- The Coordinator Intake sentence and grader 07 of `coordinator-explore-then-ask` now share one definition of filler: a content-free request for permission or continuation, or a bare restatement of the open question.
- A closing that names the concrete next step after the owner's decision is not filler. Both texts include the same example: "I'll turn this into an implementation plan".

## Results
| Check | Result | Bar |
|---|---|---|
| Offline re-judge (13 items × 3 votes: 12 saved replies + 1 synthetic negative) | Round 0: 11/13 majorities matched → text corrected before any paid run. Round 1: 13/13 | every majority matches the pre-registered expectation: met |
| Saved replies under the corrected grader | 0.7.0: 6/6; 0.7.1: 5/6 | like-for-like baseline |
| Live 0.7.2, grader 07 (6 runs) | 6/6 | ≥ 5/6: met |
| Live 0.7.2, grader 06 (6 runs) | 5/6 | ≥ 5/6: met |
| Live 0.7.2, grader 03 (not in scope) | 2/6 | — |

## Reading
- Grader 07 passes because its definition changed. The saved 0.7.0 and 0.7.1 replies already score 6/6 and 5/6 under the corrected grader, so 0.7.2's 6/6 is not evidence of a behaviour change in the Coordinator.
- The old-grader figures (0.7.0 1/6, 0.7.1 2/6) are not comparable.
- No guard cases (`external-clis-*`) ran in this batch.

## Spend
- Offline re-judge: US$0.68 (80 `claude -p` calls across two rounds).
- Live run: US$1.68 (agent 1.41 + judge 0.27).
- **Total: US$2.36**, against a cap of US$3. The live run was capped at `--max-cost-usd 1.97` by the budget gate (3.00 − 0.68 − 0.35).
