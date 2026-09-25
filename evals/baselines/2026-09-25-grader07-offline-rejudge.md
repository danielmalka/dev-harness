# Grader 07 offline re-judge (PLAN-008, T-501)

Method: each saved final reply was judged against the grader-07 text with a headless `claude -p --setting-sources "" --model sonnet --tools "" --strict-mcp-config --disable-slash-commands --max-budget-usd 0.04` call, 3 votes per item. The expected outcome for each item was pre-registered in `.harness/tasks/PLAN-008/PLAN.md` before any vote. Replies: `evals/results/2026-09-25-coordinator-0.7.1/traces/` (0.7.1) and `evals/results/2026-09-25-restate-after-larger/traces/` (0.7.0).

## Round 0 (grader text as first approved)

| Item | Expected | Votes | Majority | Match |
|---|---|---|---|---|
| 0.7.1 run-0 | FAIL | FAIL FAIL PASS | FAIL | yes |
| 0.7.1 run-1 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-2 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-3 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-4 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-5 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-0 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-1 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-2 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-3 | PASS | FAIL PASS FAIL | FAIL | **no** |
| 0.7.0 run-4 | PASS | FAIL FAIL PASS | FAIL | **no** |
| 0.7.0 run-5 | PASS | PASS PASS PASS | PASS | yes |
| synthetic negative (0.7.1 run-1, closing replaced by "Should I proceed?") | FAIL | FAIL FAIL FAIL | FAIL | yes |

Round 0 spend: US$0.35.

## Round 1 (corrected text: "a short one such as \"I'll turn this into an implementation plan\" counts")

| Item | Expected | Votes | Majority | Match |
|---|---|---|---|---|
| 0.7.1 run-0 | FAIL | FAIL PASS FAIL | FAIL | yes |
| 0.7.1 run-1 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-2 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-3 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-4 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.1 run-5 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-0 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-1 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-2 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-3 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-4 | PASS | PASS PASS PASS | PASS | yes |
| 0.7.0 run-5 | PASS | PASS PASS PASS | PASS | yes |
| synthetic negative (0.7.1 run-1, closing replaced by "Should I proceed?") | FAIL | FAIL FAIL FAIL | FAIL | yes |

Round 1 spend: US$0.33.

## Outcome

- Round 0: two rows (0.7.0 run-3, run-4) disagreed with the expected PASS. In both, the judges split on whether "…and I'll turn this into an implementation plan" names a concrete next step. Under the plan's stop rule, the paid run did not start.
- The Coordinator decided that naming the artifact that follows the decision counts. The grader-07 parenthetical and the Coordinator Intake parenthetical both gained that example; see PLAN-008 Deviations.
- Round 1: all 13 majorities match the expected outcome. **Bar met.**
- Scores under the corrected grader: 0.7.0 is 6/6 and 0.7.1 is 5/6. These are the like-for-like baseline for the 0.7.2 run. The grader's definition changed, so the old-grader figures (0.7.0 1/6, 0.7.1 2/6) are not comparable.
- Total offline spend: US$0.68. It counts toward the batch cap of US$3.

## Live measurement (PLAN-008 T-503)

Package: this branch after T-501/T-502, built with `go run ./cmd/dh build`.
Command flags: `coordinator-explore-then-ask`, 6 runs, `--trust-plugin --scaffold --ablation none --no-publish --model sonnet --judge-model sonnet --keep-temp --allow-tools Write Edit --max-cost-usd 1.97`. The cap comes from the budget gate: 3.00 − 0.68 − 0.35.
Raw JSON: `evals/baselines/2026-09-25-coordinator-explore-then-ask-0.7.2.json`.

| Grader | Under the corrected grader: 0.7.0 / 0.7.1 (offline re-judge above) | 0.7.2 candidate (live) | Bar |
|---|---|---|---|
| 07 closing | 6/6 / 5/6 | 6/6 | ≥ 5/6: met |
| 06 restates demand first | 6/6 (0.7.0 live) / 6/6 (0.7.1 live) | 5/6 | ≥ 5/6: met |
| 03 one blocking question (not a bar of this batch) | 2/6 / 0/6 (live, old runs) | 2/6 | not in scope |

Reading:
- Grader 07 passes because its definition changed. The saved 0.7.0 and 0.7.1 replies already score 6/6 and 5/6 under the corrected grader, so the 0.7.2 score of 6/6 is not evidence of a behaviour change in the Coordinator.
- What changed is the rule: a closing that names the next artifact or action after the owner's decision is no longer counted as filler. The Coordinator Intake sentence and grader 07 now say the same thing.
- No guard cases were run in this batch; the change is one Intake sentence.

Spend: agent US$1.41, judge US$0.27, live total US$1.68. Batch total with the offline re-judge: **US$2.36**, against a cap of US$3.
