---
name: safe-refactoring
description: Use when code structure must improve while contracted behavior stays identical, including duplication, tangled conditionals, oversized units, poor boundaries, and legacy code with thin coverage. Also use when deciding whether a cleanup is worth its risk, or when an earlier cleanup changed behavior by accident. Do not use when the goal is new behavior, when a defect must be corrected, or when a public contract is meant to change.
author: malka
metadata:
  provenance: adapted
  sources: ["tech-debt prioritization", "code-simplifier review agent", " refactoring-specialist"]
---

# Safe Refactoring

## Overview

This procedure changes structure while keeping contracted behavior identical, and proves it with the same inputs producing the same observable outputs before and after. Every step is small enough to revert alone, and checks run between steps rather than at the end. The core principle is that preserved behavior is the acceptance criterion: a cleanup that improves the code and changes one observable result has failed, no matter how much better the code reads.

## When to use

- A named structural problem is slowing work or causing defects: duplication, deep nesting, an oversized unit, a leaking boundary.
- Legacy code must be made changeable before a feature can land safely.
- A planned change is blocked by structure and the structural step should be separate and reviewable.

## When not to use

- The task adds or alters behavior. Use incremental-implementation.
- Something is broken and the cause is unknown. Use systematic-debugging.
- The public contract is supposed to change. Use api-contracts, then plan the migration.
- The only motivation is taste, with no observed cost. A refactor with no named problem is churn.

## Inputs

| Input | If missing |
| --- | --- |
| Authorization and write set, including any characterization tests | Reuse recorded authorization for this scope. If it does not cover needed edits, continue read-only assessment and report the missing boundary to the Coordinator; the need for coverage does not authorize test or configuration edits. |
| The named problem and the evidence of its cost | Do not start. Ask the Coordinator for the problem statement. Without one there is no way to say the work succeeded. |
| The behavior boundary: what is contracted versus incidental | Derive it from callers, tests, and published contracts, then state your reading as a decision before touching code. |
| Existing checks and their current results | Run them and record the baseline. A suite already failing is a baseline fact, not a reason to skip it. |
| Coverage at the boundary you will move | Use existing tests or recorded examples first. If preservation needs new characterization tests, add them only within the authorized test scope (step 5); otherwise report the gap before changing the structure. |
| Representative inputs and their current outputs | Capture them yourself from the current code and keep them as the comparison set. |

## Procedure

1. **State the problem and the target reduction.** Name the structure that hurts, the evidence that it hurts, and the measure you expect to move: duplicated blocks removed, branches in a function, call sites of a leaked internal, size of a unit. Vague goals produce unverifiable work.
2. **Draw the behavior boundary.** Write down what is contracted and therefore frozen: public signatures, transport shapes, persisted formats, error codes, observable side effects, ordering that callers rely on, and any performance figure that was promised. Everything else is incidental and may change. When a fact sits on the line, treat it as contracted.
3. **Find the callers and the invariants.** Search the repository for every use of the code you will move, including dynamic and configuration-driven uses. Record each one with a relative path. An invariant nobody wrote down is still an invariant; state it explicitly before it becomes a regression.
4. **Establish the baseline.** Run the existing checks and record command plus result. Capture the observable output of the representative inputs. If you claim a performance property will be preserved, measure it now, in the same way you will measure later.
5. **Add characterization tests where needed and authorized.** Use existing coverage and recorded examples when sufficient for the risk. If new tests are necessary, check the authorized test scope before writing them; if absent, report the gap rather than expanding the write set. A characterization test asserts what the code does today, including behavior that looks wrong. Do not correct anything while writing it. Report apparent defects to the Coordinator as separate items.
6. **Plan small reversible steps.** One transformation per step: extract a unit, inline a unit, rename, introduce a parameter object, replace a conditional chain, move a boundary. Each step must be revertible on its own and must leave the tree working.
7. **Execute one step, then check.** Run the pertinent checks after every step. Never chain two transformations before a check; when something breaks, you need to know which move broke it.
8. **Compare after against baseline.** Same inputs, same observable outputs, same check results. Any difference is a defect in the refactor. Revert the step and report it. A contract change is a separate task with its own authorization and migration, never a line item inside a cleanup.
9. **Verify the reduction.** Measure the problem again with the measure from step 1 and report the before and after figure. A refactor that moved code without reducing the named problem should be reverted, not defended.
10. **Hold the scope.** No new feature, no contract change, no unrelated fix, no dependency added. Improvements you notice go into the report as candidates.
11. **Report and hand over.** Deliver preservation evidence, the reduction, the remaining risks, and how to revert. A required failed or not-run preservation check leaves the work partial or blocked, even when before and after statuses match. Passing builder checks make the stable state ready for independent QA and review, not self-approved. Facts and incidents go to the Coordinator, the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`.

## Output format

```
## Status
<ready for independent QA/review / partial / blocked> - name every required check that is failed or not-run; any of them keeps the status partial or blocked

## Problem
<the structural problem and the evidence of its cost>

## Behavior boundary
- Contracted (frozen): <items, with relative paths>
- Incidental (free to change): <items>

## Baseline
| Check or input | Before | After |
| --- | --- | --- |
| <command or input case> | <result or output> | <result or output> |

## Steps applied
1. <transformation> - <files> - checks after: <result>

## Reduction achieved
- <measure>: <before> -> <after>

## Characterization tests added
- <test name> - <behavior it pins, including anything that looks wrong>

## Behavior preserved
- <evidence that the contracted list is unchanged>

## Remaining risks
- <risk> - <why it remains> - <what would detect it>

## Revert path
- <how to undo, step by step or by commit range>

## Reported, not fixed
- <defect or improvement noticed, left for a separate task>
```

## Quick reference

| Change | Allowed inside a refactor |
| --- | --- |
| Rename an internal symbol, split or merge internal units | Yes |
| Replace a conditional chain with dispatch, extract a value object | Yes |
| Move code between files or layers, adjust imports | Yes, if no consumer path changes |
| Change a public signature, route, payload, or error code | No, that is a contract change |
| Change a persisted format or migration behavior | No, that is a data task |
| Fix a defect noticed along the way | No, report it |
| Add a feature flag, option, or extension point | No, no consumer exists |
| Change log wording that a consumer parses | No, it is contracted in practice |

Order the work by cost, not by irritation: priority rises with how much the structure slows current work and with the risk it carries, and falls with the effort of the change. The first candidates are the ones that block a planned feature, carry a history of defects, and can be moved in small steps.

Step size rule: if you cannot describe the step in one sentence and revert it with one action, it is two steps.

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Refactoring with no baseline | Nothing can prove behavior survived | Record checks and outputs before touching code |
| Fixing a defect during the cleanup | The diff now mixes behavior change with structure, and review cannot separate them | Report it, schedule it |
| Writing characterization tests that assert the desired behavior | They fail immediately and tempt a silent behavior change | Assert what the code does today |
| Several transformations between checks | The breakage cannot be attributed and the revert is large | One step, one check |
| Adjusting the expected value to match the new output | The only evidence of preservation is destroyed | Treat the difference as a defect |
| Changing a public signature because it is cleaner | Consumers break under the banner of cleanup | Contract change, separate task, migration |
| Extracting to remove visual repetition | A wrong abstraction is costlier than the duplication it replaced | Extract when it is the same concept with the same reason to change |
| Declaring success without measuring the problem again | Movement gets reported as improvement | Report the before and after figure |
| Refactoring the whole area because you are already there | The change stops being reviewable and revertible | Hold the named problem |

## Example

Input: "checkout total is computed in three places, they already disagreed once".

```
## Behavior boundary
- Contracted: checkout response totals, rounding to 2 decimals,
  discount-before-tax order (src/checkout/api.ts:40, tests at :118)
- Incidental: internal helper names, file layout

## Baseline
| Check or input | Before | After |
| --- | --- | --- |
| npm test -- src/checkout | 18 passed | 18 passed |
| cart [10.00, 5.55] coupon SAVE10 | 14.00 | 14.00 |
| cart [0.01] no coupon | 0.01 | 0.01 |

## Steps applied
1. Characterize the three totals with 6 literal cases - checks: passed
2. Extract computeTotal into src/checkout/total.ts - checks: passed
3. Point all three call sites at it - checks: passed

## Reduction achieved
- Total computation sites: 3 -> 1

## Characterization tests added
- checkout totals, 6 literal cases, pins the admin rounding order as it is today

## Behavior preserved
- The three representative carts return the same totals before and after, and
  the 18 existing checkout tests pass unchanged

## Remaining risks
- No coverage of multi-currency carts; a rounding change there would go unnoticed

## Revert path
- Revert steps 3, 2, 1 in that order; each is a single commit

## Reported, not fixed
- The admin path rounds before the discount, producing a one cent difference.
  Pinned by a characterization test as current behavior. Needs a decision.
```

The one cent divergence is frozen by a test and escalated, not corrected inside the refactor.

## Related

Roles: refactorer, qa-verifier, code-reviewer, coordinator, backend-builder. Commands: `/dev-harness:refactor`. Skills: regression-testing, incremental-implementation, systematic-debugging, architecture-decisions, api-contracts.

## Proof case

An authorized fixture refactor with a named structural problem preserves the same representative inputs and observable outputs, shown side by side, while the stated measure drops by a reported figure. A behavior that looks wrong is pinned by an authorized characterization test and reported rather than corrected, and no public signature changes. Missing test write scope produces a reported gap without edits; an unavailable required preservation check leaves the result partial or blocked.
