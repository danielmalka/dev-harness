---
name: discover
description: Turn a vague ask into a brief with acceptance
author: malka
argument-hint: "[idea, pain or request]"
metadata:
  roles: [coordinator, product-discovery, document-validator]
  skills: [requirements-discovery, document-review, external-clis]
  writes: .harness/prd/ and .harness/stories/ when authorized
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Idea or pain: $ARGUMENTS. Dispatch `product-discovery` with the Agent tool, model `sonnet`, and require it to load the kit skill `requirements-discovery`. Instruct it to search the existing product docs and any prior brief in the project before asking anything, and to raise one blocking question at a time through you.
Once the PRD is written, dispatch `document-validator`, model `opus`, requiring the kit skill `document-review`, only when `reviewers.document` is absent from `.harness/project.yaml` or contains `claude`; with the PRD path, the discovery notes or the owner's request, and the previous review report when it exists. Entries of the form `cli:...` in `reviewers.document` follow the `external-clis` procedure described in `.agents/coordinator.md` ("External CLI reviewers"); both the Claude dispatch and any `cli:...` entry count against the two-specialist concurrency cap. Persist the report body it returns as `<document>.review.md`; the validator is read-only. On `changes required`, re-dispatch `product-discovery` with only the listed blocking points and validate again. Cap: two correction rounds.

## Prerequisites
The stated goal in the owner's words. If $ARGUMENTS carries no goal, ask the owner what should become possible that is not possible now, then dispatch. Constraints, deadline and stack that cannot be sourced are recorded as unknown, never invented.

## Output
The `requirements-discovery` brief: problem, users affected, desired behavior, out of scope, examples, alternatives considered, limits, acceptance as "When X, then Y" marked requirement or suggestion, hypotheses, open decisions and evidence. When the owner wants it kept, persist it as `.harness/prd/PRD-<n>.md` from the PRD template bundled with this kit. Record the acceptance items and open decisions in `.harness/MEMORY.md`. The owner receives the PRD only after `approved` or the round cap, with the report path `<document>.review.md` and any point left open.

## Limits
- Discovery writes discovery artifacts only: brief, acceptance, open decisions. Never application source, configuration or test files.
- Ask only about decisions that change what gets built; anything findable is found, not asked.
- A brief authorizes no implementation, dependency install, commit, push, publish or deploy.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dev-harness:plan` once acceptance is agreed and the open decisions that block scope are answered.
