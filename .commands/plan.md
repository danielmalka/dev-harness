---
name: plan
description: Split authorized work into provable slices
author: malka
argument-hint: "[brief, PRD path or scope]"
metadata:
  roles: [coordinator, implementation-planner, solution-architect, api-designer, data-engineer]
  skills: [implementation-planning, architecture-decisions, api-contracts]
  writes: .harness/tasks/<id>/ and .harness/adr/ when authorized
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Brief or scope: $ARGUMENTS. Dispatch `implementation-planner` with the Agent tool, model `sonnet`, and require it to load the kit skill `implementation-planning`. Only when the plan genuinely depends on them, dispatch in addition: `solution-architect` (model `opus`, skill `architecture-decisions`) for a boundary that other components will depend on, `api-designer` (model `sonnet`, skill `api-contracts`) for a new or changed interface, `data-engineer` (model `sonnet`, skill `data-migrations`) for a schema change. At most two specialists run concurrently and only on independent work.

## Prerequisites
Acceptance criteria or a brief, and a repository map or direct file evidence. Missing acceptance: stop and report that `/dev-harness:discover` produces it. Missing map: stop and report that `/dev-harness:understand` produces it. Check commands come from `.harness/project.yaml`; when undiscoverable, affected slices stay provisional with the missing check named. Consult `.harness/RISKS.md` when the work touches established behavior, a critical rule or a high-risk area, and carry each relevant incident into a slice constraint.

## Output
The `implementation-planning` plan: goal, global constraints, risk record, file structure, and per slice the goal, inputs, outputs, candidate files, consumes, produces, dependencies, parallelizable-with, checks with status not-run, readiness, done-when and risks, plus decisions needed and the self-review. Persist each slice as `.harness/tasks/<id>/TASK.md` from the TASK template bundled with this kit; a durable boundary decision becomes `.harness/adr/ADR-<n>.md` from the ADR template. Record the slice list, dependencies and open decisions in `.harness/MEMORY.md`.

## Limits
- The planner never implements; no product source, configuration or test file is edited by this command.
- Plan approval authorizes nothing: no dependency install, commit, push, publish or deploy follows from it.
- Do not invent candidate files from directory names, and do not claim a check is runnable without its source.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dev-harness:build` on the first ready slice, once the decisions that block it are answered.
