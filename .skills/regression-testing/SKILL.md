---
name: regression-testing
description: Use when a delivered change must be verified against its acceptance criteria, when tests are needed for a slice that is already stabilized, or when a report must state exactly which checks passed, failed, and were never run. Also use when deciding how much testing a change deserves. Do not use while the implementation is still moving, to diagnose an unexplained failure, or when the request is to change production code rather than verify it.
author: malka
metadata:
  provenance: adapted
  sources: ["tdd reference (good tests, mocking)", "test-driven-development and writing-good-tests", "verification-before-completion", "testing-strategy", " qa-engineer"]
---

# Regression Testing

## Overview

This procedure turns acceptance criteria into executed checks and returns a matrix where every criterion is linked to a command and to a literal result: passed, failed, or not-run. Verification depth follows risk, not habit. A passing test establishes only its assertions; observing it catch the intended defect provides additional regression evidence. A check that was not executed never becomes a pass, no matter how confident anyone is. Any required failed or not-run check leaves verification partial or blocked.

## When to use

- A change is stabilized and its acceptance criteria must be verified.
- A slice needs behavior tests that a reviewer can trust.
- A report needs an honest account of what was verified and what was not.
- Coverage of a risky path is unknown and must be established before release.

## When not to use

- The implementation is still changing. Wait for a stabilized state.
- The failure cause is unknown. Use systematic-debugging.
- The fix itself is the task. Use incremental-implementation.
- The check is a user interface interaction matrix. Use ui-verification for that surface and this skill for the rest.

## Inputs

| Input | If missing |
| --- | --- |
| Acceptance criteria for the change | Ask the Coordinator. Without criteria, testing has no target and the matrix has no left column. |
| A stabilized change and the list of files touched | Do not begin. Verifying a moving target produces results nobody can reproduce. |
| The project test commands and how to run a subset | Discover them from the project profile or the test configuration. If none exists, record every check as not-run with the reason. |
| The write scope agreed for tests | Reuse recorded authorization if it covers the files. Otherwise run permitted existing checks without edits and report any needed test or configuration changes to the Coordinator. Missing scope grants no write permission. Never production code. |
| Environment needs (fixtures, disposable database, browser) | Use synthetic data inside the authorized scope. Creating, seeding, or destroying databases and external services requires explicit authorization; without it the check is not-run with the missing dependency named. |

## Procedure

1. **Read the criteria and the change.** Map the changed files to the behaviors they can affect. Confirm the state is stable; if the implementer is still editing, stop and say so.
2. **Turn each criterion into an observable.** For every criterion write what an outside observer would see when it holds. A criterion with no observable outcome is a defect in the criterion, not a reason to skip it; report it back.
3. **Size the depth by risk.** Use the risk table in Quick reference. A low-risk change gets its behavior asserted once; a high-risk one gets boundaries, failures, and the integration seam. Depth is a decision you state, not an accident.
4. **Inventory existing coverage first.** Run the existing suite for the area and read what it already asserts. Do not re-test what is covered; do name the gap between existing coverage and the criteria.
5. **Write tests inside the agreed scope, sequentially.** QA writes test files and test configuration only. A production defect is reported, never patched here. Never write tests into files the implementer is currently editing; test authoring follows implementation, it does not run beside it.
6. **Check failure detection where needed.** Reuse documented before-fix evidence for the same test and defect, then run the test on the current state. If additional defect-detection evidence is needed, use an explicitly authorized test fixture, or an isolated disposable copy created inside the authorized test scope; never outside it, and never revert a fix or inject a defect into the production source under evaluation. Keep the evaluated state unchanged, identify the experimental state separately, and record the expected failure plus the current-state result. If no authorized method is available, report the gap; a required detection check stays not-run and leaves verification partial or blocked.
7. **Run the checks and read the output.** Full command, full output, exit status, failure count. Do not extrapolate from a partial run or present a result from an earlier state as current-state verification. Historical failure evidence remains labeled with its original state.
8. **Build the acceptance matrix.** One row per criterion with the exact command, the literal result, and the evidence path. Identify required checks from the acceptance and task context; do not silently downgrade them. Criteria you could not verify keep a not-run row with the blocking reason. Use verified only when all required checks pass on the current stable state; otherwise report partial or blocked. QA's result does not replace independent code review.
9. **Report defects, do not fix them.** For each defect give reproduction, expected, actual, the criterion it violates, and a severity. Hand it to the implementer through the Coordinator.
10. **Stop your own loop at three.** If your test code fails three times for reasons inside the test rather than the product, stop and report the blocker instead of grinding. This cap applies to your own test code only. The correction loop between implementer and QA still stops at two rounds and returns to the Coordinator.
11. **Report, do not record.** Results, defects, and any incident go to the Coordinator, the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`. Detailed evidence belongs under `.harness/tasks/<id>/`.

## Output format

```
## Verdict
<verified / partial / blocked> and the required evidence that decides it

## Acceptance matrix
| Criterion | Check (exact command) | Result | Evidence |
| --- | --- | --- | --- |
| <id or text> | <command> | passed / failed / not-run | <relative path or output line> |

## Failure detection proof
- <test name> - <prior evidence or authorized isolated fixture>, observed failure <message>, current evaluated state <result>; or <not-run and reason>

## Defects found
- <id> | severity | reproduction | expected | actual | criterion violated

## Coverage gaps
- <behavior with no check> - <why it matters>

## Not verified
- <check> - <blocking reason>

## Next step
- <the single next action, and who owns it>
```

## Quick reference

Depth by risk. Pick the highest row that applies to the change.

| Risk signal | Minimum depth |
| --- | --- |
| Money, authentication, authorization, personal data, destructive operation | Normal case, each failure branch, boundary values, and the integration seam with real collaborators |
| Persisted schema, migration, or public contract | Normal case, compatibility of old and new shape, and rollback or reversibility where possible |
| Business rule with branches | Each branch asserted once, plus the boundary between branches |
| Internal behavior with one path | One behavior test on the observable outcome |
| Rendering-only or copy change | Pertinent existing checks and direct inspection; exercise affected interaction paths when present, without requiring new tests or invented failure cases |

What a check proves.

| Check | Proves | Does not prove |
| --- | --- | --- |
| Behavior test | The asserted behavior holds at the seam | Anything not asserted |
| Typecheck | Types agree | That the feature works |
| Lint | Style and some defect classes | Behavior |
| Build | It compiles and packages | Behavior |
| Manual run | It worked once, here, now | That it is repeatable |

Properties of a test worth keeping.

| Property | Test it against |
| --- | --- |
| Behavior over implementation | Would it survive a refactor that preserves behavior? If not, rewrite it |
| One reason to fail | Does the name contain "and"? Split it |
| Independent expectation | Is the expected value computed the way the code computes it? Replace with a literal |
| Names the break it catches | Can you state the production change that makes it fail? If not, it is a change detector |
| Deterministic | Does it depend on a sleep, wall clock, ordering, or shared state? Wait on a condition, freeze time, isolate the fixture |
| Real code where it matters | Mock only at system boundaries, mirror the real data shape completely, and never assert on the mock itself |
| Isolated | Can it run alone and in any order? |

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Marking an unavailable check as passed | One false green makes the whole matrix worthless | not-run with the missing dependency named |
| Accepting typecheck as behavior verification | Well-typed code does the wrong thing all the time | Assert the observable outcome |
| Claiming defect-detection proof from a green test alone | Passing assertions do not establish that the intended defect is caught | Use relevant before-fix evidence or an authorized isolated fixture; report missing proof honestly |
| Editing production code to make a test pass | The verifier stops being independent of the work | Report the defect, let the implementer fix it |
| Asserting that a mock was called | The assertion holds as long as the mock exists, whatever the code does | Assert the real observable result |
| Sleeping to wait for async work | Flaky under load, slow when it works | Wait on the condition with a timeout |
| Testing everything at the same depth | Effort lands on trivia while the risky path stays uncovered | Use the risk table |
| Copying the implementation into the expectation | The test passes by construction | Hand-derived literal or fixture |
| Partial mock of an external payload | Passes locally, breaks in integration when an omitted field is read | Mirror the complete documented shape |
| Running a subset and reporting the suite | The claim is broader than the evidence | Report exactly what ran |

## Example

Input: criterion AC-2, "a duplicate invite is rejected with 409 already_invited".

```
## Verdict
Partial: the concurrent AC-2 case fails and required AC-3 is not-run.

## Acceptance matrix
| Criterion | Check | Result | Evidence |
| --- | --- | --- | --- |
| AC-2 sequential duplicate invite rejected | npm test -- src/invites -t 'rejects sequential duplicate invite' | passed | src/invites/service.test.ts:41 |
| AC-2 concurrent duplicate invite rejected | npm test -- src/invites -t 'rejects concurrent duplicate invite' | failed | src/invites/service.test.ts:58 |
| AC-3 invite email delivered | npm run test:integration | not-run | no mail sandbox in this environment |

## Failure detection proof
- rejects sequential duplicate invite - builder's before-fix record shows
  "expected 409, received 201" for this test; QA reran it on the current
  evaluated state and it passed. QA made no production source edits.

## Defects found
- D1 | medium | invite twice within the same second | 409 | two rows created |
  criterion AC-2 | race on concurrent insert, no unique constraint

## Coverage gaps
- Invite expiry has no check; nothing would catch a pending invite that never expires

## Not verified
- npm run test:integration - no mail sandbox in this environment

## Next step
- Implementer fixes D1; QA reruns AC-2 concurrent (owner: backend-builder)
```

AC-3 stays not-run with its reason. The race is a defect report, not a fix by QA.

## Related

Roles: qa-verifier, backend-builder, frontend-builder, code-reviewer, coordinator. Commands: `/dev-harness:verify`, and used inside `/dev-harness:build`, `/dev-harness:fix`, `/dev-harness:refactor`. Skills: incremental-implementation, ui-verification, systematic-debugging, safe-refactoring, delivery-readiness.

## Proof case

A known defect in an authorized test fixture is caught by the suite and reported with reproduction, expected, and actual, without edits to the production source under evaluation. With no test write authorization, existing permitted checks run without edits. An unavailable required check appears as not-run with the missing dependency named and leaves the verdict partial or blocked. A typecheck result is never presented as verification of a behavior criterion.
