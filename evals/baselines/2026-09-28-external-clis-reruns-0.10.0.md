# Deferred 0.9.0 re-runs, measured in the 0.10.0 batch (2026-09-28)

Owner decision (2026-09-28): the three `external-clis-*` re-runs deferred from the 0.9.0 measurement move to this batch. Package built from `feat/loop-planning`, model and judge sonnet, `--runs 1`, `--max-cost-usd 0.75`, budget gate before each run.

| Case | Result | Cost (agent + judge) |
|---|---|---|
| external-clis-document-merge-blocker | 1.0 (all pass) | 0.16 + 0.01 |
| external-clis-parallel | 1.0 (all pass) | 0.36 + 0.03 |
| external-clis-document-only | 0.8: `cli-fixture-only` FAIL 3/3 (run 1) and 3/3 (run 2) | 0.17 + 0.01; 0.15 + 0.01 |

`document-only` diagnosis: the verdict, the absence of a claude dispatch and the approved report are correct. In both runs the final reply also adds a "My spot-check: I read both files myself…" paragraph beside the single CLI verdict; the grader requires the fixture to be the *only* verdict source. The same case against the **v0.8.0** package also spot-checks and fails 2-1 (0.15 + 0.01). **This behavior is pre-existing and was not introduced by 0.9.0 or 0.10.0.** Likely cause: the coordinator rule "Spot-check load-bearing claims from a report that skipped QA and independent review" is being applied to a single CLI review verdict, which is independent review. Recorded as debt: clarify that rule (not changed in this batch).

Batch eval spend (US$6 cap): ≈ US$3.33.
