---
name: ui-verification
description: Use when a user-facing flow has been built or changed and someone must confirm it actually works when exercised, covering keyboard operation, two viewport widths, loading, error and empty states, and accessibility basics. Also use when a change was declared done on the strength of a typecheck or a passing unit test alone. Do not use for code with no rendered surface, for pure visual art direction, or as a substitute for diagnosing a failure that already has a reproduction.
author: malka
metadata:
  provenance: adapted
  sources: ["accessibility-review WCAG AA pass", "interface-design review criteria", "browser automation practice", " frontend-developer"]
---

# UI Verification

## Overview

Exercise the flow the way a person would, then report every criterion as passed, failed, or not-run. A type check, a build, and a green unit test prove that the code compiles and that the units behave; none of them proves that a human can complete the flow. When a browser is available, the verification produces observed evidence. When one is not, the interaction criteria are recorded as not-run and visual verification is declared pending. A criterion is never marked passed because it probably works.

## When to use

- A screen, component, form, or navigation path was added or changed.
- Acceptance criteria describe something a user does, not something a function returns.
- A change was reported as complete with only a compiler or linter as evidence.
- Loading, error, or empty behavior was added and nobody has seen it render.
- Before handoff or release of anything with a rendered surface.

## When not to use

- The change has no user-facing surface. Use regression-testing instead.
- The question is which visual direction to take. That is a design decision, not a verification.
- A defect is already reproducible and the job is to find its cause. That is systematic-debugging.
- The flow is still being built. Verify a stabilized state.

## Inputs

| Input | If missing |
| --- | --- |
| The acceptance criteria, as user-visible behavior | Derive a criteria list from the change and return it to the Coordinator for confirmation before running; do not verify against criteria you invented and confirmed to yourself. |
| How to run the application locally: command, port, entry route | Ask. Record every interaction criterion as not-run until it is answered. |
| Seed or fixture data that produces the loading, error, and empty states | Verify the states you can reach and record the others as not-run with the reason. |
| A browser, headless or not | Run the static checks, record the interaction criteria as not-run, and declare visual verification pending. Never install a browser, driver, or package to obtain one, and never download a runtime at verification time. |
| Required visual assets: fonts, icons, images, tokens | Do not fetch them at verification time. A missing asset is a failed criterion, because the project must ship it. |

## Procedure

1. **Turn the criteria into a checklist.** One row per user-visible behavior, phrased as an observable outcome. Add the standing rows from Quick reference that apply to the changed surface and identify the required checks from acceptance and task context. Do not invent loading or failure states for a surface that has none. This checklist is the output skeleton; every applicable row will end as passed, failed, or not-run.

2. **Establish what evidence is possible.** Determine whether the application runs locally and whether a browser is available. Headless is fine and is the default in an environment with no display. Record the answer now, because it decides which rows can ever reach passed. Never install a browser, driver, or package to obtain one, and never download a runtime at verification time. No browser means the interaction rows are not-run.

3. **Run the project's own checks first.** Build, type check, lint, unit and component tests, using the commands from the project profile. Record each as a command plus a result. These rows are evidence of compilation and unit behavior, and of nothing else. Do not let a green result close an interaction row.

4. **Reach the flow.** Start the application, navigate to the entry point, and confirm the surface renders. If it does not render, stop and report a failure here; every downstream row is not-run, not failed.

5. **Exercise the happy path with the pointer.** Complete the flow end to end with the input a real user would give. Record what you did and what appeared at each step.

6. **Exercise the same path with the keyboard only.** Reach every interactive element by tab, activate it with the expected key, and confirm the focus indicator is visible at each stop. Check that focus order follows the visual order, that nothing is reachable but invisible, and that a dialog traps focus while open and returns it to the trigger on close. A control that the pointer can use and the keyboard cannot is a failed row, not a note.

7. **Exercise the applicable states.** Drive the surface into its loading, error, and empty states using authorized fixtures or browser response controls, without editing production source. Confirm each applicable state is distinguishable, says something useful, and offers a way forward. An error state that renders a blank region is a failure. A required state that cannot be reached remains not-run with the reason.

8. **Check two viewport widths.** One narrow, around a phone, and one wide, around a desktop. Confirm no horizontal page scroll at the narrow width, that content is not clipped or overlapped, that interactive targets remain comfortably tappable, and that the flow stays completable at both. Record the two widths you used.

9. **Run the accessibility basics.** Every input has a programmatic label, not just adjacent text. Images that carry meaning have alternative text and decorative ones do not announce. Structure uses real semantics: headings in order, lists as lists, buttons as buttons, links as links, landmarks present. Interactive elements expose a name, a role, and a current value. Text contrast meets at least 4.5 to 1, large text and interface components at least 3 to 1; compute the ratio from the resolved foreground and background colors (computed style, not the token name), and if you cannot resolve both colors the row is not-run with that reason, never estimated by eye. Errors are identified in text, not by color alone.

10. **Collect the evidence.** Capture a screenshot per state that changed, at both widths where layout matters, saved under `.harness/tasks/<task-id>/evidence/`, referenced by relative path. Read the browser console and record errors and warnings; an uncaught error during the flow is a failed row even when the screen looks correct.

11. **Confirm the assets belong to the project.** Any font, icon, image, or style the surface depends on must be present in the repository or in a declared dependency. A surface that only looks right because an external resource happened to load is not verified.

12. **Fill the matrix honestly.** Passed means observed. Failed means observed to be wrong. Not-run means it was not exercised, and it always carries a reason. Accept only when every required check passed on the current stable state. A required failed row means rejected; a required not-run row means incomplete, never accepted. Visual verification is complete only when all required visual and interaction checks were exercised, even if some failed; otherwise it is pending. Browser availability alone establishes neither completion nor acceptance. Accepted here means the UI criteria were observed to hold on this state; it is not delivery approval. Marking the work verified and closing the task belongs to the Coordinator.

13. **Hand back the result.** This procedure verifies; it does not fix. Report failures to the author or the Coordinator with enough detail to reproduce. Only the Coordinator writes `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`; report facts, evidence paths, and incidents for the Coordinator to record.

## Output format

```
## UI verification
- Surface: <route or component>
- Environment: <how it was run> | browser: headless | headed | none
- Viewports: <narrow px> and <wide px>

## Matrix
| # | Criterion | Method | Result | Evidence |
|---|-----------|--------|--------|----------|
| 1 | <user-visible behavior> | pointer | passed | <relative/path/shot.png> |
| 2 | <behavior> | keyboard | failed | <relative/path/shot.png> |
| 3 | <behavior> | viewport 390 | not-run | no browser available |

## Failures
- [<criterion>] <what happened> vs <what was expected>
  - Steps: <how to reproduce>
  - Evidence: <relative path or console excerpt>

## Console
- Errors: <count> - <first error text>
- Warnings: <count>

## Checks
- <command> - passed | failed | not-run

## Summary
- Passed: <n> | Failed: <n> | Not-run: <n>
- Visual verification: complete | pending (<reason>)
- Verdict: accepted | rejected | incomplete
```

## Quick reference

| Standing row | What proves it | Never accept instead |
| --- | --- | --- |
| Flow completes with a pointer | The end state observed after real input | A passing unit test of the submit handler |
| Flow completes with the keyboard alone | Every control reached by tab and activated | An ARIA attribute present in the source |
| Focus is always visible | Observed indicator at each stop | A focus style defined in the stylesheet |
| Loading state renders | The state seen while the request is pending | The spinner component existing |
| Error state renders and explains | The state seen after a forced failure | A catch block in the code |
| Empty state renders and offers a next step | The state seen with no data | An empty array check |
| No horizontal scroll at narrow width | Observed at the stated width | A responsive utility class |
| Content usable at wide width | Observed at the stated width | A max-width rule |
| Inputs have programmatic labels | Name announced or label association observed | Visible placeholder text |
| Meaningful images have alternative text | Attribute present and accurate | A filename that describes the image |
| Text contrast at least 4.5 to 1 | Ratio computed from the resolved foreground and background colors | A palette said to be accessible, or a ratio estimated by eye |
| No uncaught console errors during the flow | Console read after the run | A clean screenshot |
| Required assets ship with the project | Asset located in the repository or a dependency | A resource that loaded from somewhere |

| Evidence available | What may reach passed | What the summary must say |
| --- | --- | --- |
| Browser, headed or headless | Only rows exercised and observed to hold | Complete only after all required visual and interaction rows were exercised; otherwise pending with the missing checks named |
| Application runs, no browser | Build and unit rows only | Visual verification pending, no browser |
| Application does not run | Static rows only | Verification incomplete, flow unreachable |

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Reporting a pass because the type check succeeded | Compilation says nothing about whether a person can finish the flow | Keep interaction rows open until they are exercised |
| Skipping the keyboard pass because the pointer path worked | Ships a surface that a keyboard user cannot operate at all | Repeat the whole flow with the keyboard alone |
| Verifying only the happy path | Applicable loading, error, and empty states can remain broken | Exercise each required state using authorized fixtures or browser controls |
| Checking one viewport width | Layout defects appear exactly at the width nobody opened | Record two widths and state the numbers |
| Treating a missing browser as a pass | Converts an unknown into a false claim of done | Mark the rows not-run and declare visual verification pending |
| Fetching a missing asset during verification | Hides a defect that will reappear for every real user | Report the missing asset as a failed row |
| Ignoring console output because the screen looked right | Silent errors become the defect reported by users next week | Read the console and record what it says |
| Fixing the defect you found | The verification loses its independence and the matrix stops reflecting the change under test | Report the failure to the author |
| Writing "accessibility OK" with no rows | An unverifiable claim that hides every real barrier | List each accessibility row with its own result |

## Example

Input: an invite form with the criteria "an invalid address shows an inline error" and "the list shows an empty state before the first invite".

```
## UI verification
- Surface: /settings/invites
- Environment: npm run dev | browser: headless
- Viewports: 390 and 1440

## Matrix
| # | Criterion                                   | Method         | Result  | Evidence                                            |
|---|---------------------------------------------|----------------|---------|-----------------------------------------------------|
| 1 | Invalid address shows an inline error       | keyboard       | passed  | .harness/tasks/T-88/evidence/invite-error-390.png   |
| 2 | Error is announced, not color only          | semantics      | failed  | error text rendered as a span                       |
| 3 | Empty list shows a next step                | pointer        | passed  | .harness/tasks/T-88/evidence/invite-empty.png       |
| 4 | No horizontal scroll at 390 px              | viewport 390   | passed  | .harness/tasks/T-88/evidence/invite-390.png         |
| 5 | Form usable without clipping at 1440 px     | viewport 1440  | passed  | .harness/tasks/T-88/evidence/invite-1440.png        |
| 6 | No uncaught console errors during the flow  | console        | failed  | console excerpt below                               |

## Failures
- [Error is announced, not color only] The message renders in a plain span with no association to the input.
  - Steps: focus the address field, enter "a@", tab out.
  - Evidence: .harness/tasks/T-88/evidence/invite-error-390.png
- [No uncaught console errors during the flow] An unhandled rejection fires on submit while the screen still looks correct.
  - Steps: submit a valid address with the network panel open.
  - Evidence: console excerpt below

## Console
- Errors: 1 - Uncaught (in promise) TypeError: cannot read properties of undefined (reading 'id')
- Warnings: 2

## Summary
- Passed: 4 | Failed: 2 | Not-run: 0
- Visual verification: complete
- Verdict: rejected
```

## Related

Roles: frontend-builder, qa-verifier, coordinator, docs-guide. Command: `/dh:verify`, and the verification step inside `/dh:build`. Skills: regression-testing for turning a failed row into an automated check, code-review for the diff behind the surface, api-contracts when the failure is in the data the surface receives, delivery-readiness for what a release needs before the decision.

## Proof case

Given a fixture flow with one keyboard trap and one unlabelled input, the verification reports both as failed rows with reproduction steps and evidence, and passes the criteria that genuinely hold. Run again in an environment with no browser, the same procedure reports the interaction rows as not-run with the reason and declares visual verification pending, never an overall pass. A browser-enabled run with an unreachable required state also stays incomplete and pending. A surface without loading or error states does not acquire invented criteria or source edits merely to exercise them.
