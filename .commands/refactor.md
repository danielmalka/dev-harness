---
name: refactor
description: Change structure while proving behavior is preserved
author: malka
argument-hint: "[target and motivation]"
metadata:
  roles: [coordinator, refactorer, qa-verifier, code-reviewer]
  skills: [safe-refactoring, regression-testing, code-review]
  writes: authorized write set only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Target and motivation: $ARGUMENTS. In sequence: dispatch `refactorer` with the Agent tool, model `sonnet`, requiring the kit skill `safe-refactoring`, with the write set and the behavior boundary. Freeze the result, then dispatch `qa-verifier` (model `sonnet`, skill `regression-testing`) to compare against the recorded baseline. Then dispatch `code-reviewer` (model `sonnet`, skill `code-review`) on that stable state. At most two specialists run concurrently and only on independent work.

## Prerequisites
The named structural problem and the evidence of its cost; without one there is no way to say the work succeeded, so ask the owner and stop. Existing checks and their current results, captured as the baseline before anything moves, including a suite already failing. Coverage at the boundary being moved: if preservation needs characterization tests, they are added only inside the authorized test scope, otherwise report the gap before changing structure.

## Output
The `safe-refactoring` report: status, problem, behavior boundary split into contracted and incidental, the baseline table of before against after, steps applied with the checks after each, reduction achieved, characterization tests added, evidence that contracted behavior is unchanged, remaining risks and the revert path per step. Record the boundary decision, the baseline evidence and the revert path in `.harness/MEMORY.md`.

## Limits
- No contract change and no new feature; this command only moves structure, and a request for new behavior goes to `/dh:plan`.
- Any observable difference between before and after is a defect, not an improvement, however much better the code reads.
- Each step must be small enough to revert on its own, with checks between steps rather than only at the end.
- A required failed or not-run check leaves the work partial or blocked, never verified.
- Never commit, push, publish or deploy from this command.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:verify` for a wider acceptance matrix, or `/dh:build` to resume the feature work the refactor unblocked.
