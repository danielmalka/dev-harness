---
name: incremental-implementation
description: Use when authorized code must be written for a planned slice of work, on the backend or the frontend, and the delivery has to come back with evidence a reviewer can check. Also use when a builder needs to know how far a change may spread, what counts as a finished slice, and how to report checks that were not run. Do not use for diagnosing an unexplained failure, for structural change with no behavior change, or when no plan or contract exists yet.
author: malka
metadata:
  provenance: adapted
  sources: ["executing-plans", "subagent-driven-development", "execution contract (YAGNI, evidence, operator gate)", " backend-developer", " frontend-developer"]
---

# Incremental Implementation

## Overview

This procedure turns one planned slice into code, tests, and evidence, without spreading outside the authorized write set. A slice is finished only when its required criteria have passing evidence on the current stable state and independent QA and review have no unresolved blocker. Report every check as passed, failed, or not-run; a required failed or not-run check leaves the slice partial or blocked. The core principle is that the smallest sufficient delta plus honest evidence beats a large confident change: a reviewer can verify the first and can only trust the second.

## When to use

- A plan, a contract, or an authorized scope names the change to build.
- A feature slice, an endpoint, a component, or a bug correction must be implemented.
- Work must land in a state that QA and review can evaluate.

## When not to use

- The cause of a failure is unknown. Use systematic-debugging first.
- The goal is structure with no behavior change. Use safe-refactoring.
- No plan, contract, or acceptance exists. Return to implementation-planning or requirements-discovery.
- The change is a specification of a public interface. Use api-contracts.

## Inputs

| Input | If missing |
| --- | --- |
| The slice to build, with its completion condition | Do not start. Ask the Coordinator for the slice boundary. A builder inventing scope is the most expensive failure in this flow. |
| The authorized write set (files or directories you may change) | Reuse recorded authorization that covers this slice. Naming a file does not itself authorize editing it; if scope remains missing, continue read-only and request the boundary from the Coordinator. |
| The contract or acceptance the slice must satisfy | Build only what you can source, and list the rest under Limitations. |
| Local conventions of the area (patterns, error handling, test style) | Read two or three files that carry the same responsibility, not merely the neighboring ones, and follow what they do. |
| Authorization to install a dependency, commit, push, publish, or deploy | **Never widen authorization.** Do not install, upgrade, or remove a dependency, and do not commit, push, publish, or deploy, without authorization recorded for this slice. Technical necessity and existing project usage are not authorization; report the need and stop. |
| The project check commands (test, lint, typecheck, build) | Discover them from the project profile or its configuration. If none is discoverable, report every check as not-run with the reason. |

## Procedure

1. **Read before writing.** Read the slice, the contract, the task record, and the files you will touch. Identify the references that carry the same responsibility as the code you are about to add, and follow their naming, typing, layering, and error handling.
2. **Fix the write set.** Write down the files you are authorized to change before the first edit. An improvement you notice outside that set is a line in the report, never an edit. Peripheral cleanup inside a build is how an unreviewable diff is born.
3. **Take exactly one slice.** One slice means one observable behavior with one completion condition. If the assigned work contains two independent behaviors, build the first, report, and take the second next.
4. **Choose the smallest sufficient delta.** No flag, option, prop, configuration, helper, or extension point without a consumer that exists today. Do not extract an abstraction to remove visual repetition; extract only when the fragments are the same concept with the same reason to change.
5. **Apply the open-decision gate.** If a decision would materially change the result, the risk, or the cost, and you do not hold the authority for it, stop and escalate in the form below rather than inventing behavior. For low-risk ambiguity, take the simplest reversible reading consistent with the project and state the assumption in the report.
6. **Implement the behavior.** Validate input at trust boundaries. Handle errors explicitly and never swallow them. Keep secrets out of code, logs, and artifacts. If the local convention is itself the cause of the defect, state the divergence in the report before replacing it.
7. **Cover the slice proportionally.** Use existing coverage first and add behavior tests only for relevant gaps within the authorized test scope. Cover the normal case and relevant failure paths when the changed behavior has them; do not invent a failure case for a copy, styling, or single-path change. Assert observable behavior, not internal structure. See regression-testing for depth by risk; presentation-only changes may use existing checks and direct inspection.
8. **Run the pertinent checks.** Record the exact command and the literal result for each. A check you did not run is not-run, never passed. Typecheck, lint, and build prove only their own property, and none of them proves the feature works.
9. **Verify the surface you changed.** For a user interface, exercise the real flow when a browser is available, including keyboard traversal and at least two widths, and cover loading, error, and empty states when the flow can reach them. Without a browser, mark visual and interaction verification as not-run and say what is pending. For a server interface, exercise the operation end to end against the contract examples.
10. **Stabilize before handing over.** No half-applied edit, no debug output, no unexplained failing check. QA and review evaluate a stabilized state; handing over a moving target wastes both.
11. **Report in the fixed shape** below and hand facts, evidence, limitations, and any incident to the Coordinator, who is the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`. Detailed artifacts belong under `.harness/tasks/<id>/` when that path is inside your write set; otherwise return them in the reply to the Coordinator.
12. **Cap the correction loop.** After six correction rounds that fail to reach a passing state, stop. Return to planning with the evidence of what was tried. A seventh blind attempt on the same approach is churn, not progress.

## Output format

```
## Conclusion
<ready for QA/review / partial / blocked; what works and what remains unverified>

## Evidence
- <relative path>:<line or symbol> - <what it establishes>

## Files changed
- <relative path> - <nature of the change>

## Checks
| Command | Result | What it proves |
| --- | --- | --- |
| <exact command> | passed / failed / not-run | <property> |

## Limitations
- <what was not verified, what was assumed, what stays open>

## Next step
- <the single next authorized action>
```

Escalation form, used when the open-decision gate trips:

```
## Open decision: <target>
Reason: <the material decision you do not have authority for>
O1. Minimal: <delta with no contract or scope growth>
O2. Structural: <scope, benefit, risk>
Blocked until O1 or O2 is chosen.
```

## Quick reference

| Check status | Means | Never write it when |
| --- | --- | --- |
| passed | The command was executed and exited successfully in this session | You inferred the result or ran a different command |
| failed | The command was executed and did not succeed | You are about to hide it inside positive prose |
| not-run | The command was not executed, or the environment cannot run it | You want the report to look complete |

Slice completion checklist, all of it or the slice is not done.

| Item | Done when |
| --- | --- |
| Behavior | The completion condition of the slice is observable |
| Scope | Every changed file is inside the authorized write set |
| Tests | Existing or added coverage addresses the changed behavior at the agreed risk depth; relevant failure paths are covered when present |
| Checks | Every required check passed on the current state; all other results and gaps are recorded |
| Surface | Required interaction or contract checks were exercised and passed; not-run keeps the slice partial or blocked |
| Independent verification (closed by the Coordinator, not by you) | The slice is only complete after QA and review evaluate the stable state with no unresolved blocker; your handover status is ready for QA/review |
| Report | Conclusion, evidence, files, checks, limitations, next step |

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Fixing unrelated things along the way | The diff stops being reviewable and the slice cannot be reverted alone | Report the observation, edit only the write set |
| Inventing behavior for an unanswered question | The wrong answer gets built, tested, and defended | Trip the open-decision gate with two options |
| Reporting a check as passed because it usually passes | One false green destroys the value of every other line in the report | Run it or write not-run |
| Treating typecheck or build as proof of the feature | Compilable code that does the wrong thing passes both | Behavior test plus exercised flow |
| Adding options nobody consumes yet | Every unused parameter becomes a maintenance and compatibility obligation | Build for the consumer that exists |
| Handing an unstable state to QA | Review effort gets spent on a state that no longer exists | Stabilize, then hand over |
| Looping on the same failing approach | Rounds seven and beyond rarely converge; the failure is structural | Stop at six rounds, return to planning with evidence |
| Approving your own delivery | The builder cannot be the independent check on the builder | Review and QA are separate roles |
| Installing a dependency without recorded authorization | A new dependency is a permanent cost decided inside a local task | Use what the project already has; report the need and stop |

## Example

Input: slice "reject an invite to an email already invited", write set `src/invites/`, contract from api-contracts.

```
## Conclusion
Partial: create_invite rejects sequential pending duplicates with 409 already_invited.
Concurrent duplicates remain possible; the slice is not complete.

## Evidence
- src/invites/service.ts:74 - duplicate check before insert

## Files changed
- src/invites/service.ts - duplicate check before insert
- src/invites/service.test.ts - two behavior cases

## Checks
| Command | Result | What it proves |
| --- | --- | --- |
| npm test -- src/invites | passed | first invite created, duplicate returns 409 |
| npm run typecheck | passed | types only, not behavior |
| npm run lint | not-run | no lint script in this project |

## Limitations
- Race between two simultaneous invites not covered; needs a unique index, outside this write set.

## Next step
- Ask the Coordinator whether the unique index belongs to this task or a data slice.
```

The race condition is reported rather than fixed, because the migration lives outside the authorized write set.

## Related

Roles: backend-builder, frontend-builder, coordinator, qa-verifier, code-reviewer. Commands: `/dh:build`. Skills: implementation-planning, api-contracts, regression-testing, ui-verification, code-review, safe-refactoring.

## Proof case

Given a planned slice with a relevant failure path, the delivery satisfies the normal and failure cases and changes no file outside the authorized write set. An unavailable required check is marked not-run with a reason and leaves the slice partial or blocked. A presentation-only fixture receives proportional checks without an invented failure case or unnecessary test infrastructure. An unanswered product decision inside the slice produces the escalation form instead of an invented behavior.
