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
Run this yourself; dispatch no specialist, because these records have a single writer. Follow the "Consolidation on explicit request" section below, with the kit skill `context-handoff` loaded for the surrounding rules.

## Prerequisites
An explicit statement from the owner in this session that no other session is writing these records; the absence of evidence of another writer is not confirmation, so ask and do not consolidate when in doubt. `.harness/MEMORY.md` must exist and hold operational content; skip the archival and say so when it does not. If a previous run was interrupted after archival, locate and verify that batch before continuing.

## Output
A batch appended to `.harness/EPOCHAL.md`: a clearly delimited, verbatim copy of `.harness/MEMORY.md` with batch id `<YYYY-MM-DD>-<n>` and an ISO 8601 timestamp with timezone offset, original dates and text preserved, any dated comment kept outside the raw block. Read the stored batch back and compare it with the original. Only after that verification, shorten `.harness/MEMORY.md` to the current work, decisions, blockers, needed evidence, next step and brief recent history, and record the batch id and consolidation time. Report the batch and the remaining active state.

## Limits
- On any failure or mismatch, leave `.harness/MEMORY.md` intact and report the operation as incomplete.
- Never re-append a batch already confirmed stored, and never remove information that was not preserved.
- Leave `.harness/RISKS.md` untouched.
- Do not claim atomicity across the two files, and never consolidate silently just because the memory grew. When `MEMORY.md` is over 40000 bytes, the Coordinator offers this command once per session; the offer is the only automatic step, and consolidation runs only when the owner asks.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Consolidation on explicit request
Run `consolidate-memory` through `/dh:consolidate-memory` only when the owner asks to archive.
1. Ensure exclusive ownership; read `MEMORY.md`, skip archival when there is no operational content. Pick a stable batch ID with timestamp and timezone and append a delimited verbatim copy of `MEMORY.md` to `EPOCHAL.md` (dated comments outside the raw block).
2. Read the batch back, compare it with the original and confirm `MEMORY.md` did not change meanwhile; on mismatch preserve `MEMORY.md` and report. Only then shorten `MEMORY.md`, keeping current work, decisions, blockers, needed evidence, next step and brief recent history, and record the batch ID and time.
3. After an interruption verify that batch before continuing; never append a confirmed batch twice or drop unpreserved information. Leave `RISKS.md` untouched, report the batch and remaining state, claim no atomicity across the two files.

## Next
`/dh:resume` to continue the remaining authorized work from the trimmed memory.
