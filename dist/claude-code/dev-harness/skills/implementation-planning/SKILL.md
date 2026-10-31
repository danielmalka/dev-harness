---
name: implementation-planning
description: Use when authorized work is wider than one obvious patch, when several files or roles would collide without an order, or when another agent must be able to start the first slice with no further conversation. Do not use for a single bounded edit, while architecture or contract questions are still open, and never to implement the plan it produces.
author: malka
metadata:
  provenance: adapted
  sources: ["implementation-planner role", "writing-plans", "plan-document-reviewer", " harness-planner", "addyosmani/agent-skills (MIT); mechanism adapted, not copied"]
---

# Implementation Planning

Split authorized work into slices that can each be finished and proven alone, ordered so none depends on a later one. A plan is finished when an agent holding only the brief, the map and this document can start slice one unaided.

## When to use

Use when the brief covers more than one patch, several files or roles would collide without an order, the work goes to another agent or session, or the sequence carries risk. Do not use for one file and one behavior (say so in one line, name the file and the check, return it to the Coordinator, do not edit), while boundaries or contracts are undecided (architecture-decisions, api-contracts first), without acceptance (requirements-discovery first), or when asked to build (this skill never edits product code).

## Inputs

- Brief or acceptance: if missing, stop and run discovery.
- Repository map or file evidence: if missing, map first; never invent candidate files.
- Constraints and durable decisions: if missing, record them as unknown and flag dependent slices.
- Authorized write scope: if unconfirmed, plan anyway and assume no permission to commit or publish.
- Check commands from the project profile: if undiscoverable, mark affected slices provisional and name the missing check.

## Procedure

1. **Decide whether a plan is warranted.** A bounded edit gets one line; when in doubt, plan, because hidden complexity found mid-task means re-planning.
2. **Ask the Coordinator for relevant `RISKS.md` incidents when the work touches consolidated or high-risk behavior.** Carry each into the affected slice as a constraint and an acceptance item; only the Coordinator writes the records.
3. **Write global constraints once at the top, with exact values, then list the files to create or modify.** Split by responsibility and follow the local pattern.
4. **Cut slices by independently verifiable behavior, two to five per focused pass.** Fold setup, config, fixtures and docs into the slice that needs them, and split where a reviewer could approve one slice and reject its neighbor.
5. **Put the test that proves a PRD rule `R<n>` in its own test-lane ticket that cites the id.** Tests tied to no PRD rule fold into their slice.
6. **Write `Doubt: yes (<trust boundary | public contract | migration>)` on an implementation slice that declares one, else `Doubt: no`.** A mechanical slice or a bug whose regression test already fails is `no`.
7. **Order by dependency, list every prior slice needed, and mark parallel slices within the two-specialist ceiling.** No slice depends on a later one.
8. **State acceptance as a checkable condition and name what each slice consumes and produces.** Prescribe code only for a load-bearing signature; the builder of slice four sees only slice four.
9. **Give each slice a check line (status `not-run`) and a risk line with the earliest signal.** A slice gets at most six correction rounds before it is reported blocked and replanned.
10. **For a slice changing `internal/build`, `internal/kit` or `dist/`, order checks: code change, `go run ./cmd/dh build`, then `go test ./...` and `go run ./cmd/dh validate .`.** The reproducibility test compares committed `dist/` with a fresh rebuild.
11. **Require the clean-clone proof only when the PR touches `cmd/`, `internal/`, `go.mod` or `dist/*/bin`.** Commit locally first (a clone sees only HEAD), clone, rerun the CI sequence there and compare binary sha256; amend only while unpushed, never `--force`. Kit text and docs need only `dh validate` and a build leaving `dist/` unchanged.
12. **Self-review against the brief: coverage, readiness, interface names, rule coverage (list an uncited `R<n>` as a `gap`) and Doubt.** Every implementation slice has `yes` or `no`; invent no value to clear a placeholder.
13. **Deliver the plan and the decisions it needs; do not implement or widen scope.** Report facts, open decisions and incidents to the Coordinator.

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
### Slice <n> (<id>): <name>
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
- Doubt: yes (<trust boundary | public contract | migration>) | no
- Risks: <what could break, and the earliest signal>
## Decisions needed
- <question>, which slice it blocks, recommended answer
## Assumptions
- <default chosen>, why it did not go to Decisions needed, what changes if wrong
## Self-review
- Coverage: <acceptance item -> slice, or gap>
- Rule coverage: <R<n> -> test-lane ticket id, or gap>
- Doubt: <slice -> yes (<reason>) | no>
- Unresolved prerequisites: none | <affected slices and what must be resolved>
- Interface consistency: checked
## Evidence
```

When persisted, each slice becomes `.harness/tasks/<id>/TASK.md` from the kit template `templates/<lang>/TASK.md`. The `<id>` in the slice heading is that same ticket id: one path segment, no `/` or `..`.

## Example

Brief: expired invites must be refused. Two slices.

```
### Slice 1 (T-1): expiry rule in the domain
- Produces: Invite.isExpired(now: Timestamp) -> bool
- Done when: a past-timestamp invite reports expired, with the clock injected
### Slice 2 (T-2): refusal at the redeem entry point
- Consumes: Invite.isExpired(now: Timestamp) -> bool
- Done when: redeeming an expired token answers 410 and the record is unchanged
```

## Related

Roles: implementation-planner, coordinator, architect, builder, qa-verifier. Commands: `/dh:plan`, `/dh:plan-loop`. Skills: requirements-discovery, repository-mapping, architecture-decisions, incremental-implementation, regression-testing, context-handoff, external-clis (dispatches a CLI planner against this output format; see its "The planning call").

A CLI planner's reply follows the output format in full. The `external-clis` transport-failure test checks only a `## Plan` heading (or `### Slice N` headings) naming a goal and one slice; a reply passing that but incomplete is judged by the wave's reviewers, not discarded.

## Proof case

Given a brief and a map, produce a plan whose first slice another agent finishes with no further conversation: candidate files exist, the done condition names a check, every symbol a later slice consumes is defined earlier. On consolidated behavior, show the risk record consulted. On a one-file request, produce no plan document. A public-contract slice is `Doubt: yes`, a one-file rename `Doubt: no`, and an `R<n>` with no citing test-lane ticket leaves the plan unfinished.
