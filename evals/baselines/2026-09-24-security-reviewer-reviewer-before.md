# security-reviewer — before baseline (T-301, PLAN-006)

- Case `security-reviewer-object-authz-miss`, unedited `.agents/security-reviewer.md`, CLI 2.1.282, `--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`.
- Pass 1 (`2026-09-24-security-reviewer-reviewer-before-pass1.json`): 3/3 full passes. The fixture narrated the defect.
- The single authorized hardening pass rewrote the fixture. Planted defect: the export handler checks ownership of `req.session.defaultNoteId` but serves `req.params.id`. Control: the delete handler checks ownership of the id it acts on.
- Pass 2 (this file's JSON): 3/3 full passes (01–04 all passing).
- Classification: **non-discriminating** after the single hardening pass. No second pass is authorized. AC-04 for `security` is reported as **not measured**, never as a passing before/after pair. The after-run still runs, as a no-regression guard (AC-03).
- Cost: pass 1 agent 0.47 + judge 0.10; pass 2 agent 0.57 + judge 0.16.
- Limit: the eval runs on sonnet, but `security-reviewer` ships on opus (PLAN-006 P9).
