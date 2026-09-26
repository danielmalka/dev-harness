---
name: architecture-decisions
description: Use when a choice about boundaries, integration or technology will outlive the task that raised it, when a change crosses components, data, operations or trust boundaries, or when two shapes are open and the decision will be expensive to reverse. Do not use for a local bugfix, for a choice with one viable option, or to redesign an area the current change only passes through.
author: malka
metadata:
  provenance: adapted
  sources: ["solution-architect role", "architecture ADR skill", "system-design", "code-architect", " architect-reviewer"]
---

# Architecture Decisions

## Overview

Decide a boundary once, with the cost stated and the discarded alternative named, then write down enough for a future reader to know whether the decision still holds. The effort is proportional to how expensive the decision is to reverse: a choice inside one file gets a sentence, a choice that other components will depend on gets a record. This skill writes design artifacts only: the comparison and, when durable, one ADR file inside the authorized write set. It never edits product code, tests or configuration.

## When to use

- A change crosses components, services, storage, operations or a trust boundary.
- Two shapes are still open and whichever is picked will be costly to undo.
- Other components will depend on the answer, so it must be stated once rather than re-derived.
- A technology is being introduced, replaced or removed.
- A plan cannot be sliced because the boundary underneath it is undecided.

## When not to use

- The defect is local and the current shape already holds. Fix it.
- Only one option is actually viable. Record it in a line and move on.
- The request is an interface between services. That is api-contracts.
- The change merely passes through an area you dislike. Taste is not a trigger.
- The decision is reversible in an afternoon. Choose, note why, keep going.

## Inputs

| Input | If missing |
|---|---|
| The requirement or brief driving the choice | Stop. A decision without a requirement optimizes for a preference. |
| A map of the affected area | Map it first. Options drawn without the real layout compare imaginary systems. |
| Constraints: deadline, team, operations, compliance, existing stack | Record each as unknown and say which option the answer would change. |
| Load, data volume or availability expectations | State the assumption you are using and mark it a hypothesis. |
| Prior decisions in the project | Search for them. Reversing a decision silently is worse than not making one. |

## Procedure

1. **Restate the decision that actually needs making.** One sentence, in terms of what changes if it goes the other way. Most requests that arrive as architecture questions turn out to be one narrow choice wrapped in three settled ones. Decide the narrow one.

2. **Size the decision before spending on it.** Ask how expensive it is to reverse in six months and how many components depend on the answer. Cheap and local gets one line in the plan. Expensive or widely depended on gets the full comparison and a record. Do not produce a record for a choice nobody will revisit.

3. **State the forces.** Functional requirements, then the non-functional ones that decide this specific case: expected scale, latency, availability, cost ceiling, team familiarity, operational burden, compliance. Only the forces that separate the options belong here; a force that points the same way for every option is not a criterion.

4. **List two or three proportional options, always including keeping the current shape.** Two options is a real comparison; five is a survey nobody reads. Each option must be one someone could actually build here with the team and the time available. An option listed only to be dismissed is padding, and it makes the comparison look more rigorous than it is.

5. **Trace each option through data, operations and security.** Data: shape, ownership, consistency, migration and backfill cost. Operations: what has to be deployed, watched, backed up, and what happens on failure. Security: trust boundaries crossed, new secrets, authorization surface, exposure of sensitive data. An option whose operational cost was never traced is not comparable to one whose cost was.

6. **Say how each option would be validated.** A decision that cannot be proven by a test, a measurement or an observable behavior is a bet. Name the check that would show the choice working and the signal that would show it failing. This line also tells the planner what the slices must produce.

7. **Recommend one and say why the others lose.** Name the losing reason per option, not a generic ranking. Distinguish a technology preference from a requirement out loud: if the reason is familiarity or taste, that is a legitimate cost factor and it must be labeled as such instead of being dressed up as a constraint.

8. **Write the contracts between components.** For every seam the decision creates or moves: who calls whom, what data crosses, what the error and timeout behavior is, who owns the state, and what the compatibility expectation is. Builders receive the chosen shape and these contracts, never a menu of undecided styles.

9. **Record an ADR only when the decision is durable.** Context, decision, consequences, status, and the condition that would make it worth revisiting. The revisit condition is what separates a record from a monument: it names the observation that would reopen the question, such as a load threshold, a second consumer, or a constraint expiring.

10. **Stop at the decision.** Do not implement the design, do not expand the change to make the new shape consistent everywhere, and do not reverse an unrelated prior decision in passing. Report the decision, its cost and any incident to the Coordinator, who is the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md` and `.harness/RISKS.md`.

## Output format

```
## Decision
- Question: <what changes if it goes the other way>
- Reversibility: cheap and local | expensive or widely depended on
- Forces that separate the options: <only the discriminating ones>

## Options
### Option A: <name>
| Dimension | Assessment |
|---|---|
| Complexity | |
| Cost to build | |
| Operational burden | |
| Data impact | |
| Security impact | |
| Team familiarity | |
- How it would be validated:

### Option B: <name>
<same shape>

## Choice
- Chosen: <option>
- Why: <the reason that decided it>
- Discarded: <option>, and the specific reason it loses
- Labeled preference: <any factor that is taste or familiarity, named as such> | none

## Consequences
- Becomes easier:
- Becomes harder:
- Must be revisited when:

## Component contracts
| Seam | Caller -> callee | Data crossing | Errors and timeouts | State owner | Compatibility |
|---|---|---|---|---|---|

## Effect on validation
- Check that proves it works: <command or test>, status: not-run (architecture decisions execute nothing)
- Signal that shows it failing:

## Risks
## Evidence
```

When the decision is durable, write the ADR from `templates/<lang>/ADR.md` (bundled with the kit) into `.harness/adr/ADR-<n>.md`, inside the authorized write set.

## Quick reference

| Trigger | Depth |
|---|---|
| Reversible inside one file | One line in the plan, no record |
| Reversible within the task | Short comparison in the plan, no ADR |
| Other components will depend on it | Full comparison, contracts, ADR |
| Data shape or ownership changes | Full comparison plus migration and rollback path |
| Trust boundary moves | Full comparison plus authorization and exposure trace |

| Signal | What it usually means |
|---|---|
| Every option scores the same on a dimension | That dimension is not a criterion; drop it |
| One option exists | This is not a decision; record it and move on |
| The reason is "cleaner" | Name the concrete cost it lowers, or it is taste |
| The new shape requires changing five untouched areas | The decision is bigger than the task authorizes; return the scope conflict to the Coordinator and stop |
| No check can tell the options apart | The difference may not matter yet; prefer the cheaper one |

## Common mistakes

| Mistake | Why it hurts | Do instead |
|---|---|---|
| Redesigning during a pinpoint fix | Turns a small correction into an unreviewable change | Fix the defect; record the design concern separately |
| Listing five options | Reads as rigor, delivers a survey nobody can act on | Two or three real options, including the current shape |
| Omitting "keep the current shape" | Guarantees a change even when no change is warranted | Include it and give it the same treatment |
| Presenting options without choosing | Pushes the decision back to whoever has the least context | Recommend one and name why the others lose |
| Calling a preference a requirement | Hides the real cost and blocks a legitimate cheaper option | Label familiarity and taste as cost factors |
| Skipping the operational trace | The cheap option is often the one nobody has to operate | Trace deployment, monitoring, backup and failure for each |
| Writing an ADR with no revisit condition | It becomes permanent by accident | Name the observation that reopens the question |
| Implementing the design while deciding it | Removes the review point and locks in an unreviewed shape | Deliver the decision and the contracts; let the plan slice it |
| Reversing an older decision in passing | Breaks consumers who relied on it and were never told | Supersede it explicitly with its own record |

## Example

Question: invite delivery crosses from the request path to an external mail provider. Send in process, or write an outbox row a worker consumes?

```
## Decision
- Question: does the request path own delivery, or hand it to a separate consumer?
- Reversibility: expensive or widely depended on; the admin resend path and any later notification inherits the seam.
- Forces that separate the options: provider latency inside the request, retry after a provider outage, and who holds the provider credential.

## Options
### Option A: send in process (the current shape; the handler already calls the provider)
| Dimension | Assessment |
|---|---|
| Complexity | lowest; one call inside the existing handler |
| Cost to build | one day |
| Operational burden | none new; a provider outage surfaces as a failed request |
| Data impact | none; nothing is stored about the attempt |
| Security impact | the provider credential lives in the web process, the widest trust surface we have |
| Team familiarity | high |
- How it would be validated: a test with a stubbed provider asserting the handler returns 201 and calls send once.

### Option B: outbox table consumed by a worker
| Dimension | Assessment |
|---|---|
| Complexity | moderate; one table, one worker, one retry policy |
| Cost to build | three to four days |
| Operational burden | a worker to deploy, watch and drain; a stuck-row alert |
| Data impact | new `invite_outbox` table; rows retained until delivered, then pruned |
| Security impact | the provider credential moves out of the web process into the worker |
| Team familiarity | moderate; no worker runs here yet |
- How it would be validated: a test that a failed send leaves the row pending and a second worker pass delivers it exactly once.

## Choice
- Chosen: Option B.
- Why: a provider outage must not fail invite creation, and retry has to survive the request.
- Discarded: Option A, because a provider timeout becomes a failed invite with no record that a send was ever attempted.
- Labeled preference: none.

## Consequences
- Easier: retry, audit of what was sent, and later notification types reuse the same seam.
- Harder: a second process to deploy and watch; delivery is now eventually consistent.
- Must be revisited when: a second consumer needs the outbox, or delivery latency has to drop below the worker poll interval.

## Component contracts
| Seam | Caller -> callee | Data crossing | Errors and timeouts | State owner | Compatibility |
|---|---|---|---|---|---|
| outbox enqueue | invite handler -> `invite_outbox` | invite id, recipient, template key | insert fails inside the creating transaction, so the invite is not created either | the invites service owns the table | new columns are additive; the worker ignores unknown ones |

## Effect on validation
- Check that proves it works: <test command from the project profile>, status: not-run (architecture decisions execute nothing)
- Signal that shows it failing: rows older than the poll interval still pending.
```

## Related

Roles: solution-architect, coordinator, implementation-planner, api-designer, data-engineer, security-reviewer. Used inside `/dh:plan`. Skills: repository-mapping supplies the affected area, implementation-planning consumes the chosen shape and its contracts, api-contracts for the interface itself, data-migrations when the choice changes stored data, security-review when the decision moves a trust boundary.

## Proof case

For a decision that crosses a component boundary, produce a comparison of two or three viable options that states the cost of each across data, operations and security, names the discarded alternative with the specific reason it lost, and states the check that would prove the choice works. Write the ADR only when the decision is durable, and show it carrying a revisit condition. On a local bugfix request, show that no comparison and no record were produced.
