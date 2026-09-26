---
name: handoff
description: Record transportable state for the next session
author: malka
argument-hint: "[optional note for the next session]"
metadata:
  roles: [coordinator, docs-guide]
  skills: [context-handoff]
  writes: .harness/tasks/<id>/ record; MEMORY.md by the coordinator only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Optional note: $ARGUMENTS. Load the kit skill `context-handoff` yourself and own the record. Dispatch `docs-guide` with the Agent tool, model `haiku`, requiring the same skill, only to draft the handoff document from the material you supply; it returns the text and the proposed memory updates to you.

## Prerequisites
The objective and the authorized scope of the task; a handoff without an objective transfers activity, not work, so ask the owner and stop. The files touched, the decisions taken and the evidence produced, recorded only where they can be sourced. `.harness/MEMORY.md` must exist; if it is absent, initialize it from the MEMORY template bundled with this kit when setup is authorized, otherwise report the missing state.

## Output
The `context-handoff` record: objective and authorization, constraints carried forward verbatim, write set, decisions with reason, source and revisit condition, evidence table with each check passed, failed or not-run, pending and blocked items with owner and completion condition, and the single next authorized step marked already granted or pending. Persist it under `.harness/tasks/<id>/`; you update `.harness/MEMORY.md` with the same state.

## Limits
- No code change: this command records state and writes nothing outside the handoff document and the memory record.
- Carry a sourced constraint forward verbatim; a recorded constraint with unknown provenance stays unresolved rather than being dropped.
- Never invent past state, a granted authorization or a check result; unverified items are marked hypothesis.
- `docs-guide` drafts but never edits the memory files.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:resume` in the next session, or `/dh:consolidate-memory` when the memory record has grown past what execution needs.
