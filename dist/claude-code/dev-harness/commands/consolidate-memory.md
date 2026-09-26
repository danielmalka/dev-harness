---
name: consolidate-memory
description: Archive memory into EPOCHAL.md, then trim it
author: malka
argument-hint: ""
metadata:
  roles: [coordinator]
  skills: [context-handoff]
  writes: .harness/MEMORY.md and .harness/EPOCHAL.md only
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Run this yourself; dispatch no specialist, because these records have a single writer. Follow the consolidation procedure in the `coordinator` agent definition, with the kit skill `context-handoff` loaded for the surrounding rules.

## Prerequisites
An explicit statement from the owner in this session that no other session is writing these records; the absence of evidence of another writer is not confirmation, so ask and do not consolidate when in doubt. `.harness/MEMORY.md` must exist and hold operational content; skip the archival and say so when it does not. If a previous run was interrupted after archival, locate and verify that batch before continuing.

## Output
A batch appended to `.harness/EPOCHAL.md`: a clearly delimited, verbatim copy of `.harness/MEMORY.md` with batch id `<YYYY-MM-DD>-<n>` and an ISO 8601 timestamp with timezone offset, original dates and text preserved, any dated comment kept outside the raw block. Read the stored batch back and compare it with the original. Only after that verification, shorten `.harness/MEMORY.md` to the current work, decisions, blockers, needed evidence, next step and brief recent history, and record the batch id and consolidation time. Report the batch and the remaining active state.

## Limits
- On any failure or mismatch, leave `.harness/MEMORY.md` intact and report the operation as incomplete.
- Never re-append a batch already confirmed stored, and never remove information that was not preserved.
- Leave `.harness/RISKS.md` untouched.
- Do not claim atomicity across the two files, and never consolidate silently just because the memory grew.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:resume` to continue the remaining authorized work from the trimmed memory.
