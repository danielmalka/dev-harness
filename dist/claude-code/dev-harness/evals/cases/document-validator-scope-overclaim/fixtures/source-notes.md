# Session notes: reviewer eval case batch

Fictional batch, used as a fixture only.

1. Added `evals/cases/code-reviewer-form-check/` — one seeded-defect case for the `code` reviewer role.
2. Added `evals/cases/security-reviewer-authz-check/` — one seeded-defect case for the `security` reviewer role.
3. Added `evals/cases/document-validator-overclaim-check/` — one seeded-defect case for the `document` reviewer role.
4. The `coordinator` role's existing case (`coordinator-explore-then-ask`) was left untouched — no new case was added for it in this batch, and it does not pair a planted defect with a control section the way the three new cases do.
5. Each of the three new cases pairs the planted defect with a correct control section in the same fixture, so a keyword-only reviewer risks flagging the control too.
6. Every new case's fixtures are copied into the run workspace by the case's own `scaffold.sh`, invoked only when the runner is passed `--scaffold`.
