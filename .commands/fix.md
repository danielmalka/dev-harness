---
name: fix
description: Reproduce, diagnose and fix a defect with evidence
author: malka
argument-hint: "[symptom, error text or issue]"
metadata:
  roles: [coordinator, debugger, qa-verifier, code-reviewer]
  skills: [systematic-debugging, regression-testing, code-review]
  writes: authorized write set only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Symptom: $ARGUMENTS. In sequence: dispatch `debugger` with the Agent tool, model `sonnet`, requiring the kit skills `systematic-debugging` and `regression-testing`, and state whether it is diagnose-only or authorized to change files, with the write set. Once a fix is applied and the state is frozen, dispatch `qa-verifier` (model `sonnet`, skill `regression-testing`). Then dispatch `code-reviewer` (model `sonnet`, skill `code-review`) scoped to the fix, proportional to its blast radius. At most two specialists run concurrently and only on independent work.

## Prerequisites
The literal error text or observed behavior, and reproduction steps or the failing command. A paraphrased error points at the wrong component: ask for the exact output. Consult `.harness/RISKS.md` when the symptom touches a critical rule or a known incident area. Ask the owner only for the authorization to change files or the environment.

## Output
The `systematic-debugging` report: status, symptom, reproduction before and after, root cause marked fact or hypothesis with evidence paths and lines, call sites checked, fix or "not applied", regression protection, ruled-out causes and the incident report for the Coordinator. Persist evidence under `.harness/tasks/<id>/evidence/`. Record the incident with severity, cause, resolution and prevention in `.harness/RISKS.md`, and the current status in `.harness/MEMORY.md`.

## Limits
- Without a reproduction, report the gap and the inconclusive path; never present a proven fix.
- A diagnosis request does not authorize source edits, instrumentation, installs or changes to external data.
- Fix the origin of the bad value, not the line where it surfaced; report every sibling call site checked.
- The regression guard counts only when it was observed failing before the fix and passing after; otherwise it is not-run with the reason.
- Never commit, push, publish or deploy from this command.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:verify` to widen the acceptance matrix, or `/dh:release` when the fix is the deliverable.
