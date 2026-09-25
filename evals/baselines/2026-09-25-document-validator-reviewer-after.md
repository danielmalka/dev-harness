# document-validator — after baseline (T-305, PLAN-006)

- Case `document-validator-scope-overclaim`, run against the edited `.agents/document-validator.md` and `.skills/document-review/SKILL.md`. CLI 2.1.282, same flags as the before run (`--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`).
- Graded result: score 0.67, 0/3 full passes. This is identical to before.

  | Grader | Before | After |
  |---|---|---|
  | 01 changes required | 3/3 | 3/3 |
  | 02 blocking findings ≥ 1 | 3/3 | 3/3 |
  | 03 cites the contradicted source point | 3/3 | 3/3 |
  | 04 correct claim not flagged | 1/3 | 1/3 |
  | 05 correct claim listed under Not raised | 0/3 | 0/3 |

- Manual reading of all six replies: the "control" claim is not clean.
  - The claim: fixtures are copied into the workspace by the case's own scaffold script.
  - `fixtures/source-notes.md` point 6 adds the condition "invoked only when the runner is passed `--scaffold`", and the document drops it.
  - In 5 of 6 runs, the validator reports that dropped condition as a finding. The finding is defensible: an omitted precondition is a gap.
  - So graders 04 and 05 penalize a legitimate finding. This is a **fixture defect**, in the same family as code grader 03.
  - It is not changed here, because the plan forbids editing a fixture after its before-run was observed. It is recorded for a separate fix.
- Reading from the manual evidence: the planted four-roles scope overclaim is detected 3/3 before and 3/3 after.
- **AC-04 (document): no measured gain.** The only discriminating graders measure a defective control. No regression.
- AC-03 no-regression guard: `doc-validator-clean` 6/6 graders and `doc-validator-gap` 11/11 graders pass on the edited role, 1 run each (`evals/results/2026-09-25-doc-validator-neighbours.json`). A clean document still gets `approved`, and the gap document still gets `changes required`.
- Limit: this eval runs on sonnet, but `document-validator` ships on opus (P9).
- Cost: scope-overclaim agent 0.62 + judge 0.08; doc-validator-* agent 0.88 + judge 0.53.
