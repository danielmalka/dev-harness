---
name: verify
description: Run acceptance checks and return an evidence matrix
author: malka
argument-hint: "[task id or acceptance criteria]"
metadata:
  roles: [coordinator, qa-verifier]
  skills: [regression-testing, ui-verification]
  writes: tests in the agreed scope only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Task or criteria: $ARGUMENTS. Dispatch `qa-verifier` with the Agent tool, model `sonnet`, and require it to load the kit skill `regression-testing`, plus `ui-verification` when the change has a user interface. Pass the acceptance criteria, the list of files touched, the check commands from `.harness/project.yaml`, the authorized test write scope and, for a user interface, how to run the application locally.

## Prerequisites
Acceptance criteria and a stabilized change; without criteria the matrix has no left column, and verifying a moving target produces results nobody can reproduce, so stop and report which is missing. Missing criteria are produced by `/dh:discover` or `/dh:plan`. Undiscoverable check commands mean every check is recorded not-run with the reason. Ask the owner only to authorize an environment dependency such as a disposable database.

## Output
The `regression-testing` verdict and acceptance matrix mapping each criterion to an exact command and a literal result of passed, failed or not-run with evidence, plus failure-detection proof, defects found, coverage gaps, what was not verified and the next step. When a user interface is involved, the `ui-verification` matrix and its screenshots go under `.harness/tasks/<id>/evidence/`. Record the verdict and the remaining gaps in `.harness/MEMORY.md`.

## Limits
- Never edit production code; test and fixture edits happen only inside the authorized test scope.
- A check that was not executed is never a pass, however confident anyone is; a required failed or not-run check means partial or blocked.
- Never create, seed or destroy a database or external service without recorded authorization; without it the check is not-run with the missing dependency named.
- Never install a browser, driver or runtime to obtain evidence; record the interaction criteria as not-run and visual verification as pending.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:review` for the independent review, or `/dh:fix` when the matrix shows a defect.
