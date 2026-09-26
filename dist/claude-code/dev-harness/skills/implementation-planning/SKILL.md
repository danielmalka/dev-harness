---
name: implementation-planning
description: Use when authorized work is wider than one obvious patch, when several files or roles would collide without an order, or when another agent must be able to start the first slice with no further conversation. Do not use for a single bounded edit, while architecture or contract questions are still open, and never to implement the plan it produces.
author: malka
metadata:
  provenance: adapted
  sources: ["implementation-planner role", "writing-plans", "plan-document-reviewer", " harness-planner"]
---

# Implementation Planning

## Overview

Split authorized work into slices that can each be finished and proven on their own, ordered so that no slice depends on a later one. The plan says what each slice must be true of when it is done, names the files it touches and the checks that decide it, and leaves the implementation shape to the builder. A plan is finished when an agent holding only the brief, the map and this document can start slice one without asking anything.

## When to use

- The brief covers more behavior than one patch can carry.
- Several files, layers or roles are involved and the order is not obvious.
- Work will be handed to another agent or another session.
- A change touches consolidated behavior and the sequence itself carries risk.
- A previous attempt sprawled because nobody drew the boundaries first.

## When not to use

- The write set is one file and one behavior. Say so in one line, name the file and the check, and return it to the Coordinator as a bounded edit. Do not make the edit yourself.
- Boundaries, contracts or technology are still undecided. Send those to architecture-decisions or api-contracts first, then plan.
- Acceptance does not exist yet. Run requirements-discovery first.
- The user asked you to build. Planning is not a stalling tactic; if the work is bounded, say so in one line and hand it back; this skill never edits product code.

## Inputs

| Input | If missing |
|---|---|
| Brief or acceptance criteria | Stop and run discovery. A plan against guessed acceptance plans the wrong thing. |
| Repository map or direct file evidence | Map the slice first. Candidate files invented from directory names are the usual cause of a wrong plan. |
| Constraints and durable decisions | Record them as unknown and flag any slice that depends on one. |
| Authorized scope: what may be written | Plan anyway, mark the scope as unconfirmed, and do not assume permission to commit or publish. |
| Check commands from the project profile | Discover commands in manifests, CI or project instructions. If still unknown, mark affected slices provisional and name the missing check; do not claim they are ready to execute. |

## Procedure

1. **Decide whether a plan is warranted.** A bounded edit gets one line naming the file and the check, not a document. Say which one you chose so the user can override it. When in doubt between bounded and planned, plan; the ratchet only goes up, and hidden complexity found mid-task means stopping and re-planning, not pushing through.

2. **Check the risk record before decomposing, through the Coordinator.** When the work touches consolidated behavior, a critical business rule or a known high-risk area, ask the Coordinator for the relevant incidents in `.harness/RISKS.md`. Carry each one into the plan as a constraint on the affected slice and as an acceptance item, not as a footnote. Routine work outside those triggers does not load the record at all. The Coordinator is the only writer of those records; specialists request the relevant excerpt instead of loading them.

3. **Write the global constraints once, at the top.** Version floors, naming rules, platform requirements, data handling limits, anything that applies to every slice. Copy exact values from the source rather than paraphrasing them. Every slice inherits this section implicitly, which is what keeps it from being repeated wrongly eight times.

4. **Draw the file structure before drawing the slices.** List what gets created and what gets modified, and what each file is responsible for. Files that change together belong together. Split by responsibility, not by technical layer for its own sake. In an existing codebase, follow the local pattern instead of importing a preferred one.

5. **Cut slices by independently verifiable behavior.** A slice is the smallest unit that carries its own check cycle and could be accepted or rejected on its own merits. Fold setup, configuration, fixtures and documentation into the slice whose deliverable needs them. Split only where a reviewer could sensibly approve one slice and reject its neighbor. Two to five behaviors per slice is a working range; a slice nobody can finish in one focused pass is too big.

6. **Order by dependency and list dependencies completely.** Foundation before the things that stand on it. Slice five naming slice four is not enough when it also needs something from slice two. No slice may depend on a later one; if two slices need each other, the cut is wrong. Mark which slices could run side by side, bounded by the ceiling of two specialists working concurrently.

7. **Specify what, precisely; leave how open.** Acceptance per slice must be objectively checkable by a test or by reading the diff. "Clean interface" is not a done condition; "the redeem endpoint answers 410 for an expired token and leaves the record unchanged" is. Prescribe code only where an exact shape is load-bearing, such as a signature another slice consumes.

8. **Name the interface each slice publishes and consumes.** The builder of slice four sees only slice four. Exact names and types of what earlier slices produced, and what later slices will rely on, are the only way the seams line up. Contract drift between slices is the most common failure of multi-slice plans.

9. **Give each slice its own check line and its own risk line.** Which check proves this slice, how it is invoked, and what result counts. Status stays `not-run` in a plan, because planning runs nothing. The risk line names what could go wrong here and the cheapest early signal that it did, remembering that a slice gets at most six correction rounds for the same issue before it is reported blocked or partial and replanned. For a slice that changes `internal/build`, `internal/kit`, or any `dist/` content, the check line's order is fixed: apply the code change, then `go run ./cmd/dh build`, then `go test ./...` / `go run ./cmd/dh validate .` — never the reverse, because the reproducibility test compares the committed `dist/` against a fresh rebuild and fails by construction if the rebuild has not run yet.

10. **Self-review the finished plan against the brief.** Three passes, inline, fixing as you go. Coverage: point at the slice that satisfies each acceptance item and list any gap. Readiness: executable slices have concrete checks and no unresolved prerequisite; unknown commands or interfaces leave only the affected slices provisional. Do not invent values to remove placeholders. Consistency: the name a later slice consumes matches the name the earlier slice produces, exactly.

11. **Deliver the plan and the decisions it needs.** Do not implement. Do not widen scope to make a slice tidier. Report facts, open decisions and any incident to the Coordinator for the memory records.

## Output format

```
## Plan: <feature>
- Goal: <one sentence>
- Source brief: <relative path or "in conversation">
- Map: <relative path or "in conversation">
- Authorization: <what may be written; commit and publish status>

## Global constraints
- <exact value, copied, one per line>

## Risk record
- Consulted: yes (via Coordinator) | not applicable (<trigger absent>)
- Incidents carried into slices: <id> -> slice <n>, as <constraint or acceptance item>

## File structure
| File (relative) | Create or modify | Responsibility |
|---|---|---|

### Slice <n>: <name>
- Goal:
- Inputs: <artifacts, prior slices, data>
- Outputs: <observable result>
- Candidate files: <relative paths>
- Consumes: <exact names and types from earlier slices>
- Produces: <exact names and types later slices rely on>
- Dependencies: <all prior slices needed, complete>
- Parallelizable with: <slice numbers or "none">
- Checks: <discovered command or explicit unresolved check>, expected result, status: not-run
- Readiness: ready | provisional (<missing prerequisite>)
- Done when: <objectively verifiable condition>
- Risks: <what could break, and the earliest signal>

## Decisions needed
- <question>, which slice it blocks, recommended answer

## Assumptions
- <default chosen>, why it did not go to Decisions needed, what changes if wrong

## Self-review
- Coverage: <acceptance item -> slice, or gap>
- Unresolved prerequisites: none | <affected slices and what must be resolved>
- Interface consistency: checked

## Evidence
```

When persisted, each slice becomes `.harness/tasks/<id>/TASK.md` from the kit template `templates/<lang>/TASK.md`.

## Quick reference

| Symptom in a draft plan | The cut is wrong because | Fix |
|---|---|---|
| A slice cannot be checked on its own | It is half a behavior | Merge it with its other half |
| Two slices must land together or nothing works | They are one slice | Merge them |
| A slice takes more than one focused pass | It is several behaviors | Split by behavior, not by file |
| A slice depends on a later slice | The order is circular | Re-cut around the shared dependency |
| A slice is "set up the project" | Setup is not a deliverable | Fold it into the first slice that needs it |
| Done when says "works correctly" | Nothing observable was named | State the input and the expected result |
| A later slice calls a name no slice defines | Interface drift | Add it to the producing slice's Produces line |

| Plan size | When it fits |
|---|---|
| One line | Single file, single behavior, check already known |
| Short plan, two or three slices | One layer, familiar area, no open decisions |
| Full plan | Several layers or roles, handoff expected, or a risk trigger fired |

## Common mistakes

| Mistake | Why it hurts | Do instead |
|---|---|---|
| Writing a long plan for a one-file edit | Costs more than the change and delays it | Say it is bounded and name the file and the check |
| Slicing by layer: all models, then all handlers | No slice proves anything until the last one | Slice by behavior that can be checked end to end |
| Leaving "handle errors appropriately" in a step | The builder invents behavior nobody agreed to | Name the case and the expected result |
| Listing only the immediately previous dependency | The builder starts a slice whose foundation is missing | List every prior slice the work needs |
| Prescribing exact code everywhere | Freezes decisions the builder is better placed to make | Fix only the load-bearing signatures |
| Marking checks as passing | Planning executes nothing; the claim is false | Status is not-run until a builder runs it |
| Skipping the risk record on a consolidated-behavior change | Repeats an incident the project already paid for | Ask the Coordinator for relevant incidents and carry them into slices |
| Implementing slice one while writing the plan | Removes the review point the plan exists to create | Deliver the plan, then let the build command start it |
| Adding a nice extra slice nobody asked for | Scope creep dressed as thoroughness | Keep it under decisions needed, as a suggestion |
| Cloning before committing to "prove" a dist rebuild | `git clone` copies HEAD; uncommitted changes are invisible in the clone, so the comparison proves nothing | Commit locally after the gate is green and review approves, then clone that commit, rerun the CI sequence there (`go test ./...`, `go run ./cmd/dh validate --source-only .`, `go run ./cmd/dh build`, `git diff --exit-code --stat -- dist` with the CI exclusions, `go run ./cmd/dh validate .`), and compare binary sha256 between clone and tree; amend only while unpushed, never with `--force`, then push |

## Example

Brief: expired invites must be refused. Map: three files and one test file. Bounded check says two slices.

```
## Global constraints
- Response for an expired invite: 410, body unchanged from the existing error shape.

## Slice 1: expiry rule in the domain
- Outputs: Invite exposes isExpired against an injected clock
- Candidate files: src/domain/invite.<ext>, tests/invite_expiry.<ext>
- Produces: Invite.isExpired(now: Timestamp) -> bool
- Dependencies: none
- Checks: <test command from the project profile>, the two new cases fail before and pass after, status: not-run
- Done when: an invite past its timestamp reports expired and one before it does not, with the clock injected, not read globally
- Risks: a global clock read makes the test flaky near midnight; the injected clock is the signal

## Slice 2: refusal at the redeem entry point
- Consumes: Invite.isExpired(now: Timestamp) -> bool
- Candidate files: src/http/invites.<ext>, tests/invite_redeem.<ext>
- Dependencies: slice 1
- Done when: redeeming an expired token answers 410 and the stored record is unchanged
```

## Related

Roles: implementation-planner, coordinator, solution-architect, backend-builder, qa-verifier. Command: `/dh:plan`. Skills: requirements-discovery supplies the acceptance, repository-mapping supplies the candidate files, architecture-decisions settles a boundary before slicing, incremental-implementation executes a slice, regression-testing turns the check lines into real checks, context-handoff carries an unfinished plan to the next session.

## Proof case

Given a brief and a map, produce a plan whose first slice another agent starts and finishes with no further conversation: candidate files exist in the tree, the done condition is verifiable by a named check, and every symbol a later slice consumes is defined by an earlier slice. On a request that touches consolidated behavior, show the risk record consulted through the Coordinator and at least one prevention carried into a slice's acceptance. On a one-file request, show that no plan document was produced.
