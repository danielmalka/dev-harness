---
name: release
description: Assemble a release readiness package for decision
author: malka
argument-hint: "[change, slice or version]"
metadata:
  roles: [coordinator, release-manager, devops-engineer]
  skills: [delivery-readiness]
  writes: readiness record and notes only; no remote action
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Change or version: $ARGUMENTS. Dispatch `release-manager` with the Agent tool, model `sonnet`, and require it to load the kit skill `delivery-readiness`. When the package needs clean-checkout pipeline or packaging evidence, also dispatch `devops-engineer` (model `sonnet`, same skill) to reproduce the build from a fresh checkout; it installs nothing, provisions nothing and reaches no live server. At most two specialists run concurrently, only on independent work and only on the same frozen state. A delivery PR opens with the PRD or brief, the decisions and any ADR, and closes with the project documentation kept current; the Coordinator never opens a documentation-only PR; at the end of every batch, whichever command closes it, the Coordinator dispatches docs-guide before delivery.

## Prerequisites
The verified change with its files and acceptance criteria, the QA result with each check marked required or optional and passed, failed or not-run, and the independent review with its findings and their resolution. Missing QA: report that `/dh:verify` produces it. Missing review: report that `/dh:review` produces it. Missing version policy: propose a version and mark it a proposal. Treat every consumer as deployed and every configuration as unset when they are unstated, and record the assumption.

## Output
The `delivery-readiness` package: verdict ready for decision or not ready with the blocking reason, scope, evidence table, findings with severity and state, version and compatibility classification, consumers affected, expand-and-contract sequence when breaking, release notes, and the rollout and rollback plan with its triggers and irreversible parts. Persist it under `.harness/tasks/<id>/`. Record the verdict, the residual risks and the pending decision in `.harness/MEMORY.md`. The Coordinator confirms `docs-guide` ran for this batch before the package is reported ready for decision.

## Limits
- Perform no commit, push, tag, publish or deploy here; "prepare the release" is not permission to publish.
- An explicit publish authorization that already exists is recognized and routed to a capable authorized executor, never re-asked and never expanded by this command.
- Any required failed or not-run check, or an incomplete required review, leaves the package not ready even with zero reported findings.
- Secrets and configuration are named by purpose only; a value never enters the package.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:handoff` once the owner has the package, or `/dh:fix` when a blocking finding remains.
