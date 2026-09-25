# Coordinator intake — larger sample with the recalibrated grader 06 (T-306, PLAN-006)

- Case: `coordinator-explore-then-ask`, 6 runs per side. CLI 2.1.282. Flags `--allow-tools Write Edit --model sonnet --judge-model sonnet --max-cost-usd 2.25`. Grader 06 recalibrated as recorded in `evals/results/2026-09-25-grader06-regrade.md`.
- Before: a copy of the freshly built package with only the case's embedded Coordinator body swapped for `1ce355f`, the 0.5.5 body. `diff -r` against the after package differs only in that `prompt.md`.
- After: `dist/claude-code/dev-harness`, built after the grader edit. Its Coordinator body is the 0.6.0 one; this batch does not change it.

| Grader | Before (0.5.5 body) | After (0.6.0 body) |
|---|---|---|
| 01 explored before answering | 6/6 | 6/6 |
| 02 fact not asked | 6/6 | 6/6 |
| 03 single blocking question | 0/6 | 2/6 |
| 03b numbered options with recommendation | 3/6 | 6/6 |
| 04 no Write | 6/6 | 6/6 |
| 05 no Edit | 6/6 | 6/6 |
| **06 restates the demand first (recalibrated)** | **0/6** | **6/6** |
| 07 no filler "may I proceed" | 3/6 | 1/6 |
| 08 no writer dispatched | 6/6 | 6/6 |
| Score | 0.67 | 0.83 |

## Reading

- **AC-05 is met.** The recalibrated grader scores the saved traces consistently with the manual reading (18/18 votes). On the larger sample it separates the two bodies cleanly: 0/6 before and 6/6 after, with unanimous votes. This is the measured evidence for the 0.6.0 restatement rule, which 0.6.0 itself could only claim from manual reading.
- 03b rises from 3/6 to 6/6, which agrees with the 0.6.0 "recommendation is advice" clause.
- **07 falls from 3/6 to 1/6.**
  - The after replies close with "Let me know which of the three you want, and I'll turn this into an implementation plan". The judge counts that as a content-free permission request on top of the question already asked.
  - The before replies that passed closed with "I won't touch any code until you tell me …".
  - This is a behavior of the 0.6.0 Coordinator body. It is recorded as a finding for a later fix. It is not fixed here, because `.agents/coordinator.md` is outside this batch's write set.
- Six runs per side, one case, one model: this shows direction, with every vote on 06 unanimous.
- Cost: before agent US$1.48 + judge US$0.27; after agent US$1.36 + judge US$0.25.
