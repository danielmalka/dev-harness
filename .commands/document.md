---
name: document
description: Write the docs needed to use and operate a change
author: malka
argument-hint: "[subject or area to document]"
metadata:
  roles: [coordinator, docs-guide, document-validator]
  skills: [delivery-readiness, context-handoff, document-review]
  writes: documentation files in the authorized scope
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Subject: $ARGUMENTS. Decide first whether the change alters how someone uses, runs or continues the work; when it does not, state that and stop. Otherwise dispatch `docs-guide` with the Agent tool, model `haiku`, requiring the kit skill `delivery-readiness` for operational and runbook sections or `context-handoff` for continuity documents. For a pt-BR HTML document in the kit's standard, also require `doc-template-html`; that role has no shell, so instruct it to copy the template rather than run the stamp script. Raise the model to `sonnet` only when the subject spans several subsystems, and record the reason.
A delivery PR opens with the PRD or brief, the decisions and any ADR, and closes with the project documentation kept current; the Coordinator never opens a documentation-only PR; at the end of every batch, whichever command closes it, the Coordinator dispatches docs-guide before delivery. This `docs-guide` dispatch happens at the end of an implementation batch regardless of whether a validated-template document (PRD, STORY, TASK, ADR) is involved — a tutorial-only or README-only update still counts.
The validation cycle below applies only when the persisted document is one of these four: `templates/<lang>/PRD.md`, `STORY.md`, `TASK.md`, `ADR.md` (or their `templates/en/` mirror). Any other output of this command — a tutorial, a runbook, a standalone HTML internal document, a handoff — does not trigger it. When it applies: once `docs-guide` has persisted the document, dispatch `document-validator`, model `opus`, requiring the kit skill `document-review`, with the document as the artifact under review and its sourced material (files, acceptance criteria, QA evidence, decisions) as the source. The validator is read-only and never edits the document; dispatch stays exclusively with the Coordinator. Persist the report it returns as `<document>.review.md`. On `changes required`, re-dispatch `docs-guide` with only the listed blocking points and validate again. Cap: six correction rounds; the owner receives the document only after `approved` or the round cap, with any point left open.

## Prerequisites
The subject and the sourced material it documents: files, acceptance criteria, QA evidence, decisions. Without sourced material, report the gap rather than writing from inference. The authorized write path for the document; without it, `docs-guide` returns the text in its reply and nothing is written. Commands and paths quoted in the documentation come from `.harness/project.yaml` or the tree; an unconfirmed command is marked unverified.

## Output
The document in the format its skill defines, persisted at the authorized path under `.harness/` (for example `.harness/tasks/<id>/` for task documentation), with planned behavior and implemented behavior clearly distinguished and every check reported as passed, failed or not-run. `docs-guide` reports proposed memory updates back to you; you record them in `.harness/MEMORY.md`.

## Limits
- `docs-guide` never edits `.harness/MEMORY.md`, `EPOCHAL.md` or `RISKS.md`; it reports updates to the Coordinator.
- Never document a behavior as implemented without evidence that it exists; planned stays labeled planned.
- No product source, configuration or test file is edited by this command; use relative paths only.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:release` when the documentation completes a delivery package, or `/dh:handoff` to close the session.
