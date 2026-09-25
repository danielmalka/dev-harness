# code-reviewer — after baseline (T-303, PLAN-006)

- Case `code-reviewer-equivalent-form-miss` (hardened fixture), edited `.agents/code-reviewer.md` + `.skills/code-review/SKILL.md`, CLI 2.1.282, same flags as before (`--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`).
- Grader result: score 0.78, 1/3 full passes. Graders: 01 request changes 3/3, 02 names the `--workspace=` miss 3/3, 03 control not wrongly flagged 1/3. This is identical to the before pass 2 (0.78, 1/3; 03 1/3).
- Manual reading of all six final replies (before pass 2 and after): no reply claims that `classifyLintStep` or any lint path is broken. Every failing grader-03 run had unanimous FAIL votes on a reply that states the lint and build paths are correct, or that names `classifyLintStep` only in a cosmetic naming note. Grader 03, phrased as a double negative ("Does no finding claim…"), is a **case defect**: the sonnet judge mis-scores correct replies. It is not changed here, because editing a grader after seeing its result is the tuning the plan forbids. It is recorded for a separate fix.
- Reading by manual evidence: detection 3/3 → 3/3, control intact 3/3 → 3/3.
- **AC-04 (code): no measured gain.** The case does not discriminate on real behavior. No regression was found.
- Behavioral change visible in the after replies: every run now opens with a `## Coverage` section that names what was traced for correctness, regression, spec compliance and security surface (D1b proof of coverage). The before replies had no such section.
- Cost: agent 0.83 + judge 0.08.
