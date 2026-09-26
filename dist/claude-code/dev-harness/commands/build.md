---
name: build
description: Implement one authorized slice with QA and review
author: malka
argument-hint: "[task id or slice]"
metadata:
  roles: [coordinator, backend-builder, frontend-builder, data-engineer, qa-verifier, code-reviewer]
  skills: [incremental-implementation, data-migrations, regression-testing, ui-verification, code-review]
  writes: authorized write set only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Slice: $ARGUMENTS. In sequence: dispatch the builder the slice requires with the Agent tool, model `sonnet` in every case, requiring the kit skill `incremental-implementation`: `backend-builder` for server, service or domain work, `frontend-builder` for user interface work, `data-engineer` (also skill `data-migrations`) for schema work. Freeze the resulting state, then dispatch `qa-verifier` (model `sonnet`, skill `regression-testing`, plus `ui-verification` when the slice has a user interface). Then dispatch `code-reviewer` (model `sonnet`, skill `code-review`) on that same stable state. At most two specialists run concurrently and only on independent work; never run QA or review while the evaluated code, contracts, configuration or tests are still changing.

## Prerequisites
An authorized slice with its completion condition, its write set and the acceptance it must satisfy. Missing any of these: report the gap and that `/dh:plan` produces the slice, then stop. Check commands come from `.harness/project.yaml`; when undiscoverable, every check is reported not-run with the reason. Ask the owner only for a decision that changes scope or contract.

## Output
The `incremental-implementation` report (conclusion, evidence, files changed, checks table, limitations, next step), the QA acceptance matrix, and the independent review with severity-graded findings and their resolution. Evidence and screenshots go under `.harness/tasks/<id>/evidence/`; update `.harness/tasks/<id>/TASK.md` with status and results. Record status, effective models, findings and residual risk in `.harness/MEMORY.md`, and any severe incident in `.harness/RISKS.md`.

## Limits
- Build only inside the authorized write set; naming a file does not authorize editing it.
- No dependency install, upgrade, removal, commit, push, publish or deploy without authorization recorded for this slice.
- QA writes tests after the builder finishes and never edits production code; any later edit invalidates the affected evidence and review, so rerun both on the new state.
- After six unsuccessful correction rounds on the same issue, report blocked or partial and replan with the owner; changing the model does not reset that limit.
- A required failed or not-run check means partial or blocked, never verified.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:build` on the next ready slice, or `/dh:release` when the authorized scope is complete.
