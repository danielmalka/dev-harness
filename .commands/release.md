---
name: release
description: Assemble a release readiness package for decision
author: malka
argument-hint: "[change, slice or version]"
metadata:
  roles: [coordinator, release-manager, builder]
  skills: [delivery-readiness]
  writes: readiness record and notes only; no remote action
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Change or version: $ARGUMENTS. Dispatch `release-manager` with the Agent tool, model `sonnet`, and require it to load the kit skill `delivery-readiness`. When the package needs clean-checkout pipeline or packaging evidence, also dispatch `builder` in lane `infra`, mode `design` (model `sonnet`, same skill), with no product write, to reproduce the build from a fresh checkout; it installs nothing, provisions nothing and reaches no live server. Both work on the same frozen state.

## Prerequisites
The verified change, its QA result and its independent review. Stop and report the gap when QA is missing (`/dh:verify` produces it) or the review is missing (`/dh:review` produces it); every other missing input follows `delivery-readiness`'s Inputs table.

## Output
The package in `delivery-readiness`'s Output format, persisted under `.harness/tasks/<id>/`. Record the verdict, the residual risks and the pending decision in `.harness/MEMORY.md`.

## Limits
- No commit, push, tag, publish or deploy runs from this command; `release-manager`'s Limits govern an existing publish authorization.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:handoff` once the owner has the package, or `/dh:fix` when a blocking finding remains.
