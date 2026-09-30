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
Acceptance criteria and a stabilized change. Stop and report which is missing when either is absent; `/dh:discover` or `/dh:plan` produce criteria. Ask the owner only to authorize an environment dependency such as a disposable database.

## Output
The report in `regression-testing`'s Output format. When a user interface is involved, the `ui-verification` matrix and its screenshots go under `.harness/tasks/<id>/evidence/`. Record the verdict and the remaining gaps in `.harness/MEMORY.md`.

## Limits
- Test and fixture edits happen only inside the authorized test scope; `qa-verifier`'s Limits and the skills govern the rest.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:review` for the independent review, or `/dh:fix` when the matrix shows a defect.
