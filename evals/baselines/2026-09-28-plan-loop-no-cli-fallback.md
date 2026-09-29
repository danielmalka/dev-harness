# plan-loop-no-cli-fallback — baseline (2026-09-28)

Case: `evals/cases/plan-loop-no-cli-fallback/` (T-803, PRD-005 RF-02/AC-02/AC-03). Package built from `feat/loop-planning`. Local `claude plugin eval` (CLI 2.1.283), model and judge sonnet, `--runs 1`, `--max-cost-usd 1.00`, `--allow-tools Read Glob Grep Skill Agent Write Edit`. Pass bar fixed before any run: all five graders pass.

- Run 1: 3/5. Grader 01 (`focus: trace` cannot see assistant text) and grader 03 (the regex missed the plugin-qualified `dh:implementation-planner`) were case defects. Behavior was correct in the trace: the notice "k = 0 → no-CLI fallback" came before any dispatch, then exactly one planner dispatch and one validator dispatch, with no `wave-N/` write. US$0.85 agent + 0.25 judge.
- Correction 1 (qa-verifier): grader 03 anchored on `"subagent_type":[ \t]*"(dh:)?implementation-planner"`; grader 01 moved to `last_message`. The ordering claim was verified manually on the run-1 trace.
- **Run 2: 5/5 (bar met).** US$0.70 agent + 0.14 judge. JSON: `2026-09-28-plan-loop-no-cli-fallback.json`.
- Batch eval spend so far (US$6 cap): ≈ US$1.94.
- Case review (code-reviewer, sonnet): 1 Major — no grader checked that `PLAN.review.md` is written (AC-03). Correction 2 added grader 06 (deterministic regex on the trace, anchored on `file_path`). It was validated offline: it matches run 2's trace (a `Write` to `.harness/tasks/PLAN-1/PLAN.review.md`) and rejects a synthetic negative. Following the offline rule for deterministic graders, the run-2 result stands at **6/6** with no new paid run.
