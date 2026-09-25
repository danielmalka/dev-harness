# code-reviewer — before baseline (T-301, PLAN-006)

- Case `code-reviewer-equivalent-form-miss`, unedited `.agents/code-reviewer.md`, CLI 2.1.282, `--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`.
- Pass 1 (`2026-09-24-code-reviewer-reviewer-before-pass1.json`): 3/3 full passes (01, 02, 03 all passing). Classified **non-discriminating**. The fixture narrated the defect in prose, so the reviewer could paraphrase it without tracing the code.
- The single authorized hardening pass (PLAN.md P4) rewrote the fixture. Planted defect: `normalizeCommand` handles `cd <dir> &&` but not the `npm run test --workspace=server` form that `workspaces.config.js` actually passes. Control: every lint/build command uses a handled form.
- Pass 2 (this file's JSON): score 0.78, 1/3 full passes. Graders: 01 request changes 3/3, 02 names the miss 3/3, 03 control not wrongly flagged 1/3.
- Classification: **reproduced**. The planted defect is detected, but the unedited reviewer falsely flags the correct control in 2/3 runs, so the full bar the after-run targets is not met.
- Cost: pass 1 agent 0.56 + judge 0.05; pass 2 agent 0.75 + judge 0.06.
