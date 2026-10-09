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
Subject: $ARGUMENTS. Decide first whether the change alters how someone uses, runs or continues the work; when it does not, state that and stop. Otherwise dispatch `docs-guide` with the Agent tool, model `haiku`, requiring the kit skill `delivery-readiness` for operational and runbook sections or `context-handoff` for continuity documents. For a pt-BR HTML document in the kit's standard, also require `doc-template-html`; that role has no shell, so instruct it to copy the template rather than run the stamp script. Subject `pdocs` (the project overview HTML): (1) run `dh harness-path --json` with the bundled `bin/dh` (`bin\dh.cmd` on Windows) as `/dh:doctor` does, and read `pdocs`; (2) when `pdocs` is empty, stop with "run dh link"; in global mode also require setup done (`project.yaml` present in the harness folder), otherwise stop and point to `/dh:setup`; (3) in both modes remind the owner that `~/.harness` must be in `permissions.additionalDirectories` of `~/.claude/settings.json` (never write that file); (4) dispatch `docs-guide`, model `haiku`, requiring `doc-template-html`, type `projeto`, language from `language` in `project.yaml`, to create `<pdocs>/index.html` by copying `assets/projeto-modelo.<lang>.html` (no shell) or to update the existing one; extra pages go flat in the same folder, kebab-case, linked from `index.html`; sources: `docs/prd/` of the repository and `prd/` of the harness folder for the "Planejado e desenvolvido" table (`id`, title, `Status`, generation date), `CHANGELOG.md` for release titles (copied, never rewritten), first-level folders and entries for the mandatory macro diagram; a hole stays "ainda não fechado", never invented text; (5) no validation cycle (not a PRD, TASK or ADR). Instruct `docs-guide` to replace every row and fact of the model with the target project's facts and keep "ainda não fechado" / "not yet settled" where unknown.
Raise the model to `sonnet` only when the subject spans several subsystems, and record the reason.
A delivery PR opens with the PRD or brief, the decisions and any ADR, and closes with the project documentation kept current; the Coordinator never opens a documentation-only PR; at the end of every batch, whichever command closes it, the Coordinator dispatches docs-guide before delivery. This `docs-guide` dispatch happens at the end of an implementation batch regardless of whether a validated-template document (PRD, TASK, ADR) is involved — a tutorial-only or README-only update still counts.
The validation cycle below applies only when the persisted document is one of these three: `templates/<lang>/PRD.md`, `TASK.md`, `ADR.md` (or their `templates/en/` mirror). Any other output of this command — a tutorial, a runbook, a standalone HTML internal document, a handoff — does not trigger it. When it applies: once `docs-guide` has persisted the document, dispatch `document-validator`, model `opus`, requiring the kit skill `document-review`, with the document as the artifact under review and its sourced material (files, acceptance criteria, QA evidence, decisions) as the source. The validator is read-only and never edits the document; dispatch stays exclusively with the Coordinator. Persist the report it returns as `<document>.review.md`. On `changes required`, re-dispatch `docs-guide` with only the listed blocking points and validate again. Cap: two correction rounds; the full list runs only when the owner asks for the panel in this request or the document touches security; the owner receives the document only after `approved` or the round cap, with any point left open.

## Prerequisites
The subject and the sourced material it documents: files, acceptance criteria, QA evidence, decisions. Without sourced material, report the gap rather than writing from inference. The authorized write path for the document; without it, `docs-guide` returns the text in its reply and nothing is written. Commands and paths quoted in the documentation come from `.harness/project.yaml` or the tree; an unconfirmed command is marked unverified.

## Output
The document in the format its skill defines, persisted at the authorized path under `.harness/` (for `pdocs` the path is the absolute one from `dh harness-path --json`, never typed or inferred; for example `.harness/tasks/<id>/` for task documentation), with planned behavior and implemented behavior clearly distinguished and every check reported as passed, failed or not-run. `docs-guide` reports proposed memory updates back to you; you record them in `.harness/MEMORY.md`.

## Limits
- `docs-guide` never edits `.harness/MEMORY.md`, `EPOCHAL.md` or `RISKS.md`; it reports updates to the Coordinator.
- Never document a behavior as implemented without evidence that it exists; planned stays labeled planned.
- No product source, configuration or test file is edited by this command; use relative paths only, except the `pdocs` path above.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:release` when the documentation completes a delivery package, or `/dh:handoff` to close the session.
