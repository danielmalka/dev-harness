# security-reviewer — after baseline (T-304, PLAN-006)

- Case: `security-reviewer-object-authz-miss` (hardened fixture), run against the edited `.agents/security-reviewer.md` and `.skills/security-review/SKILL.md`.
- Conditions: CLI 2.1.282, same flags as the before run (`--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`).
- Result: score 1.00, 3/3 full passes (graders 01–04 pass in every run), the same as before pass 2.
- AC-04 for `security`: **not measured**. The case was non-discriminating after the single hardening pass, so a matching before/after pair proves nothing about the edit.
- AC-03 no-regression guard: **holds**. The edited role still detects the mismatched-id object-authorization defect, fills the required finding fields, and leaves the owner-filtered delete handler unflagged and listed as clean.
- Limit: the eval runs on sonnet, while `security-reviewer` ships on opus (P9).
- Cost: agent US$0.51 + judge US$0.13.
