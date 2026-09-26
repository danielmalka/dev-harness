---
name: improve
description: Diagnose and evaluate an improvement to this kit
author: malka
argument-hint: "[reproducible case: prompt, expected, observed]"
metadata:
  roles: [coordinator, harness-maintainer, qa-verifier]
  skills: [harness-authoring, harness-evaluation]
  writes: harness checkout sources only, never dist/
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Case: $ARGUMENTS. Dispatch `harness-maintainer` with the Agent tool, model `sonnet`, and require it to load the kit skills `harness-authoring` for the diagnosis and the edit and `harness-evaluation` for the before-and-after measurement. Instruct it to name the owning asset before editing anything: which trigger fired, which file was loaded, which instruction produced the output.

## Prerequisites
A concrete case with prompt, expected behavior and observed behavior; for a repair, a reproduced failure or a cited contract contradiction between two instructions. An isolated preferred answer is not evidence, so ask the owner for the reproduction and stop. A harness checkout to work in, with its state recorded, and a clear authorized scope before any edit. The cases live under `evals/cases/` with their rubrics and fixtures; a baseline must be captured before the asset changes.

## Output
The `harness-authoring` asset-change report (case, owning asset, smallest edit, what deliberately did not change, contract check, positive and negative test cases, provenance) together with the `harness-evaluation` before-and-after report with its conditions, static checks, model-graded results and variance. Record the change, the measured effect and its limits in `.harness/MEMORY.md`.

## Limits
- Apply changes only with a clear authorized scope in the harness checkout; without one, propose the edit and stop.
- Never edit `dist/` by hand; it is generated from the source assets.
- With no runner available, model-graded criteria are not-run and the result is a static correction, never a measured improvement; a changed asset is never its own before baseline.
- Do not incorporate borrowed material; naming a source does not make its content available.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:doctor` to confirm the kit still resolves, then re-run the failing case to confirm the fix.
