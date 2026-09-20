---
name: understand
description: Map the code behind a question, read-only
author: malka
argument-hint: "[question or area to map]"
metadata:
  roles: [coordinator, repo-scout]
  skills: [repository-mapping]
  writes: none (read-only)
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Question or area: $ARGUMENTS. Dispatch `repo-scout` with the Agent tool, model `haiku`, and require it to load the kit skill `repository-mapping`. Pass the question phrased as behavior, the target directory, and the check commands from `.harness/project.yaml` when it exists.

## Prerequisites
A question phrased as behavior. If $ARGUMENTS names no behavior, ask the owner what should be located before dispatching; a map with no question maps everything and answers nothing. No index, database or external tool is required, and none may be requested.

## Output
The `repository-mapping` report: question, scope, file map with symbols and how each was confirmed, flow, blast radius, checks with discovery confirmed or inferred and execution always not-run, guesses, unknowns, evidence with paths and lines, and the first files to touch. The map is returned in the reply; persist it only under `.harness/tasks/<id>/` when a task already owns it. Record the mapped area and its evidence paths in `.harness/MEMORY.md`.

## Limits
- Read-only: no file in the project is created or edited, and no check is executed.
- Every line is observed in a file that was read or labeled a guess; grep-only findings say so.
- Mapping authorizes no implementation; do not start one from this command.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dev-harness:plan` when the work is understood, or `/dev-harness:discover` when the behavior itself is still undecided.
