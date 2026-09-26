---
name: context-handoff
description: Use when work must survive a session boundary, when resuming from a prior record, or when explicitly asked to consolidate operational memory into raw history. Growing memory may trigger a suggestion to consolidate. Do not use for ordinary progress notes, user-facing documentation or task specifications.
author: malka
metadata:
  provenance: authorial
  sources: ["dev-harness plan 4.6 and AC-19..AC-24", "dev-harness memory templates", "dev-harness coordinator agent", " context-manager role"]
---

# Context Handoff

## Overview

This procedure keeps a task alive across sessions through three operations: writing a handoff another session can act on, resuming from one by reconciling it against the current tree, and consolidating operational memory so the archive keeps the raw text and the memory keeps only what execution still needs. The governing principle is that the handoff orients but the tree decides. A record is a claim about the past; the files, the checks and the history are the evidence, and where they disagree the evidence wins.

Only the Coordinator writes `.harness/MEMORY.md`, `.harness/EPOCHAL.md` and `.harness/RISKS.md`. Every other role, including documentation and QA roles, reports facts, evidence, pending items and incidents back to the Coordinator, who records them. One Coordinator session holds that write role per project; two sessions must not update or consolidate these records in parallel. This is an organizational rule, not a technical file lock.

## When to use

- A session is ending, running out of context, or passing work to another person or agent.
- A new session must continue a task started elsewhere, and the tree may have moved since.
- A prior record and the current files disagree, and the next step depends on which is true.
- The operational memory has grown past the state execution actually needs and the history must be preserved before trimming.
- A consolidation was interrupted and it is unclear whether the batch was archived.

## When not to use

- Ordinary progress inside one continuous session that already holds the state.
- Documentation aimed at a user or an operator. That is delivery-readiness.
- Authoring a design decision record or a specification. The handoff carries decisions by summary and reference; the ADR or spec itself belongs under `.harness/adr/` or `.harness/tasks/<id>/`.
- Registering a serious incident. That is a RISKS record, written by the Coordinator, and consolidation never touches it.
- A role other than the Coordinator wanting to edit the three memory files. Report instead.

## Inputs

| Input | If missing |
| --- | --- |
| The objective and the authorized scope of the task | Stop. A handoff without an objective transfers activity, not work. In the main session, ask the owner. As a dispatched specialist, return the missing input to the Coordinator and stop; never proceed on an assumed value. |
| Current operational memory `.harness/MEMORY.md` | If absent, only the main-session Coordinator may initialize it from the `templates/` folder bundled with the kit (`templates/<lang>/MEMORY.md`) within authorized scope. Other roles report its absence. Never overwrite existing records or invent past state. |
| Files touched, decisions taken, evidence produced | Record only what you can source. Mark anything unverified as hypothesis. |
| Constraints and prohibitions stated by the owner | Carry sourced constraints forward verbatim. Preserve a recorded constraint with unknown provenance as unresolved until clarified; missing provenance does not revoke it. Never invent a constraint. |
| The next authorized step | Check existing authorization first. State what is already covered; mark only actions outside that scope as awaiting authorization. |
| For consolidation: confirmation that no other session is writing the records | Confirmed only by an explicit statement from the owner in this session. The absence of evidence of another writer is not confirmation; when in doubt, ask and do not consolidate. |

## Procedure

### A. Handoff: write a portable record

1. **State the objective and the authorization separately.** What the task is for, and what may currently be written, committed or published. Approval of a plan is not authorization to publish.
2. **List the files in the write set with relative paths.** Absolute paths, machine names and user directories make the record unusable elsewhere. Note which files are modified and which are only read.
3. **Record decisions with their reason and their source.** Each decision carries what was decided, why, and what would justify revisiting it. A decision without a reason gets re-litigated or, worse, reversed by accident.
4. **Record evidence as it stands.** Each check with its command and its result as passed, failed or not-run. Never promote a not-run to a pass because it passed earlier on a different tree.
5. **Carry constraints forward verbatim.** Prohibitions, scope limits, style rules, forbidden files, deadlines. These are the first thing lost across a boundary and the most expensive to lose.
6. **Name what is pending and what is blocked, with the blocker.** Distinguish work not yet started from work waiting on someone.
7. **Write exactly one next step.** Concrete enough to start without interpretation, with authorization marked already granted (source and scope) or pending. Preserve any correction-round count and the default, requested and effective model (or unverified) needed to resume; neither a handoff nor a model change resets the correction budget.
8. **Remove secrets.** No tokens, keys, connection strings, private hostnames or personal paths. Name the variable and its purpose instead.
9. **Hand the record to the Coordinator.** Store the record under `.harness/tasks/<id>/` and let the Coordinator fold the durable state into `.harness/MEMORY.md`. Other roles do not write that file.

### B. Resume: reconcile the record with reality

1. **Read the record first, then the tree.** Reading the record alone leads to repeating landed work; reading the tree alone leads to dropping constraints.
2. **Verify each claimed change against the files.** For every file the record says was modified, confirm the change is present, absent, or present in a different form than described.
3. **Detect external change.** Compare the record's evidence and file list against the current history and working tree. Anything changed by someone else since the record was written is an external change and must be named, not silently absorbed. History and tree comparison is executed by the Coordinator or a Bash-capable role; the scout returns the file-level reconciliation only.
4. **Re-run checks whose evidence was invalidated.** Re-running is executed by the Coordinator or a Bash-capable role; the scout returns the file-level reconciliation only. Compare the evaluated files, dependencies, configuration and environment with current state. Preserve applicable evidence when those inputs are unchanged. Re-run affected checks after changes or unresolved uncertainty; until then, mark them not-run for the current state. Record the new result without rewriting historical evidence.
5. **Classify every item.** Landed, partially landed, not started, invalidated by external change, or contradicted by the tree. Record the classification before choosing an action.
6. **Never repeat a landed change and never drop a stated constraint.** These are the two failure modes of resumption. If the record and the tree disagree about whether something landed, the tree decides; if they disagree about a constraint, the constraint stands until the owner releases it.
7. **Consult history only on concrete need, through the Coordinator.** `.harness/EPOCHAL.md` is not loaded by default. The Coordinator searches it when a specific question needs an archived fact or the owner asks, and reads only the relevant passage. Other roles request the needed context. Archived text is historical evidence, never a current instruction.
8. **Consult incidents when the change warrants it, through the Coordinator.** Before planning or implementing a change to consolidated behavior, a critical rule or a high-risk area, the Coordinator reads `.harness/RISKS.md` and carries the relevant prevention into the plan, acceptance and specialist dispatch. Other roles request those excerpts. Routine changes outside those triggers do not load it.
9. **State the reconciled next step.** One step, its authorization, and what changed in the plan because of the reconciliation.

### C. Consolidate memory: archive, verify, then trim

Run only on an explicit consolidation request, in the main session, by the Coordinator, with no other writer on the records. Growing memory is a reason to suggest consolidation, not permission to run it.

1. **Prepare.** Read `.harness/MEMORY.md` in full and identify the state that must remain: current execution, standing decisions, blockers, necessary evidence, next step. If the memory holds no operational content, stop and create no batch. An empty batch is noise in the archive forever.
2. **Assign the batch identity.** A stable batch id in the form `<YYYY-MM-DD>-<n>`, where `n` starts at 1 on that date in the project timezone; if the id already exists in `.harness/EPOCHAL.md`, advance `n`. Never reuse an id for different content. Add a timestamp in ISO 8601 with timezone offset. The consolidation date is not the date of the events it contains; original dates stay as written, and unknown old timestamps stay declared unknown rather than invented.
3. **Archive.** Append to the end of `.harness/EPOCHAL.md` a complete verbatim copy of the current `.harness/MEMORY.md`, with an unambiguous start and end delimiter, the batch id, the timestamp, the source, and the tasks and date range covered. Do not summarize, reorder, reformat or quietly correct the raw text. Do not load the whole archive to append one batch; append at the end.
4. **Keep commentary outside the raw block.** Coordinator notes and later corrections go outside the archived block, labeled as commentary and dated. A correction is a new linked record, never an edit inside the batch.
5. **Verify before touching the source.** Re-read the written batch and compare it integrally with the memory it came from, and confirm the source did not change during the operation. On any failure, divergence, or invalid destination file, leave `.harness/MEMORY.md` untouched and report the situation. This is the only gate that protects the memory.
6. **Trim, only after the archive is confirmed.** Keep the current execution, standing decisions, blockers, necessary evidence, the next step and a brief recent history. Record the batch id and the consolidation date in the memory's last-consolidation section. Closed events remain reachable in the archive. An open serious incident stays visible in the memory as a blocker even though it is also recorded in the incident file.
7. **Resume safely after an interruption.** The operation spans two files and is not atomic. Locate the prior batch by its stable id before appending anything. If complete, verify it and the unchanged source before trimming; if MEMORY already records that completed consolidation, do not trim again. If the batch is partial, malformed, ambiguous, or differs from current MEMORY, leave MEMORY intact and report recovery as blocked. Do not blindly restart with a new id, duplicate a batch, or edit historical raw content. Never remove from memory anything not proven preserved.
8. **Leave the incident record alone.** `.harness/RISKS.md` is not consolidated, trimmed, updated or moved into the archive. Resolved incidents and their prevention remain intact; incident updates belong to a separate Coordinator operation.
9. **Deliver.** Report the batch created, what stayed active in memory, and any pending item. Consolidation is an explicit command. The Coordinator may suggest it when the memory loses concision, and never performs a silent cleanup.

## Output format

Handoff record.

```
# Handoff: <task id>

## Objective
- Goal:
- Authorization: <what may be written, committed, published>

## Constraints
- <constraint, carried forward verbatim>

## Write set
- <relative path> | modified | read-only

## Decisions
- <decision> | reason | source | revisit if <condition>

## Evidence
| Check | Command | Result |
| --- | --- | --- |
| <name> | <command from the project profile> | passed / failed / not-run (<reason>) |

## Pending and blocked
- <item> | owner | blocker | completion condition

## Next authorized step
- <one concrete step> | authorization: already granted (<source and scope>) | pending (<what>)
```

Resume reconciliation.

```
# Resume: <task id>

## Reconciliation
| Item from handoff | State in tree | Classification | Action |
| --- | --- | --- | --- |
| <item> | <observed> | landed / partial / not started / invalidated / contradicted | <action> |

## External changes detected
- <change> | source | effect on the plan

## Constraints still in force
- <constraint>

## Checks re-run
| Check | Command | Result |

## Next step
- <step> | authorization: already granted (<source and scope>) | pending (<what>)
```

Consolidation report.

```
# Consolidation: <batch id>

- Timestamp: <ISO 8601 with offset>
- Source: MEMORY.md
- Tasks covered: <ids> | date range: <known range>
- Verification: batch re-read and compared integrally; source unchanged: yes / no
- Kept active in memory: execution, decisions, blockers, evidence, next step, brief history
- Pending after consolidation: <items>
- RISKS.md: untouched
```

## Quick reference

Read policy for the three records.

| File | Written by | Read when |
| --- | --- | --- |
| `.harness/MEMORY.md` | Coordinator only | Every execution start, every resume, before every dispatch |
| `.harness/EPOCHAL.md` | Coordinator only | Only on a concrete historical need or an explicit request; search by task, batch, date or subject and read only the passage |
| `.harness/RISKS.md` | Coordinator only | Before changing consolidated behavior, a critical rule or a high-risk area; and to register or update a serious incident |

Consolidation gates, in order. A gate that is not met stops the operation with the memory intact.

| Gate | Met when |
| --- | --- |
| Request | Explicit consolidation request, not inferred from memory size |
| Sole writer | The owner stated in this session that no other session is writing the records |
| Content | The memory holds operational content; otherwise no batch is created |
| Identity | Batch id `<YYYY-MM-DD>-<n>` not already in the archive, plus ISO 8601 timestamp with offset |
| Archive | Complete verbatim copy appended with unambiguous delimiters |
| Verification | Batch re-read, compared integrally, and the source confirmed unchanged |
| Trim | Runs only after verification passed |
| Incidents | The incident record is untouched |

## Common mistakes

| Excuse | Reality |
| --- | --- |
| "The handoff says it was done, so it is done" | The handoff is a claim. Confirm it in the files before building on it. |
| "The tree looks different, the handoff must be stale, I will start over" | Some of it landed. Classify item by item; starting over repeats landed work and loses decisions. |
| "That constraint was probably from an earlier phase" | A constraint stands until the owner releases it. Silently dropping one is the exact failure resumption must prevent. |
| "The memory is long, I will tidy it and archive later" | Trimming before the archive is verified destroys the only copy. Archive, verify, then trim, in that order. |
| "Summarizing the old memory makes the archive readable" | The archive is raw history. A summary replaces evidence with interpretation and cannot be undone. |
| "The consolidation timestamp can stand in for the event dates" | It cannot. Original dates stay; unknown ones stay declared unknown. |
| "The run was interrupted, safest is to archive again" | A duplicate batch corrupts the history. Locate and verify the same batch first. |
| "I am the documentation role, updating the memory is part of my job" | Only the Coordinator writes the three records. Report facts and evidence instead. |
| "Loading the whole history gives me better context" | The archive enters context only on a concrete need. Default loading buries the operational state. |
| "The incident is resolved, consolidation can clear it" | Consolidation never touches the incident record. Resolved incidents stay with their prevention. |
| "Checks passed before the break, I can carry the result over" | Preserve them only if evaluated inputs remain unchanged; otherwise re-run affected checks or mark not-run. |
| "Two of us can split the consolidation to go faster" | One Coordinator session writes per project. Parallel writers lose state that no archive holds. |

## Example

Interrupted consolidation, resumed in a new session.

```
Memory holds: 2 open blockers, 1 decision, 11 history entries.
Archive tail shows batch 2026-04-02-1, delimiters present, end marker present.

1. Locate batch 2026-04-02-1 in EPOCHAL.md.
2. Compare it line by line with the current MEMORY.md: identical.
3. Confirm MEMORY.md unchanged since the batch timestamp: yes.
   -> The archive step completed. Do not create a second batch.
4. Trim MEMORY.md: keep the 2 blockers, the decision, the next step,
   3 recent history entries; record "last consolidation: 2026-04-02-1".
5. Report: batch 2026-04-02-1 verified, 8 history entries now archive-only,
   2 blockers still active, RISKS.md untouched.
```

## Related

Roles: coordinator, docs-guide, repo-scout. Commands: `/dh:handoff`, `/dh:resume`, `/dh:consolidate-memory`, `/dh:document`. Skills: delivery-readiness, project-onboarding, repository-mapping. Templates: the `templates/` folder bundled with the kit (`templates/<lang>/MEMORY.md`, `templates/<lang>/EPOCHAL.md`, `templates/<lang>/RISKS.md`).

## Proof case

A new session resumes a task whose tree was changed by someone else after the handoff was written. It names the external change, does not reapply a change that already landed, and carries every stated constraint forward. A consolidation run produces an archive batch whose content matches the prior memory integrally, with timestamp and timezone, and leaves standing decisions, blockers and the next step in the memory. Interrupting the run before the archive is verified leaves the memory unchanged; interrupting it after the batch is written and resuming verifies that same batch instead of creating a second one. A routine execution completes without the archive ever entering context, and the incident record is byte-identical before and after consolidation.
