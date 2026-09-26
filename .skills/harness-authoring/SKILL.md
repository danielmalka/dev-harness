---
name: harness-authoring
description: Use when an agent, skill, command or tutorial of this kit must be created or changed, when a trigger fires too widely or too narrowly, when an asset's instructions are ambiguous or contradict another asset, or when a provenance, version or compatibility note is needed for kit content. Do not use for changes to a consumer project's own code, for a one-off preference about tone, or to hand-edit the generated package instead of the sources.
author: malka
metadata:
  provenance: adapted
  sources: ["skill-development", "writing-for-agents", "agent-development", "command-development", "writing-skills", " skill-creator role", "dev-harness harness-maintainer agent"]
---

# Harness Authoring

## Overview

This procedure produces or changes one asset of this kit: an agent, a skill, a command or a tutorial. It fixes the contract each asset type owes, keeps triggering text separate from procedure text, and pushes heavy material behind a pointer so the loaded file stays legible. A new asset needs an authorized use case and testable expectations. A repair needs a reproduced failure or a concrete contract contradiction; an isolated preferred answer is not evidence of improvement.

## When to use

- A new agent, skill, command or tutorial is needed and its contract must be written correctly the first time.
- A trigger fires on work it should not touch, or fails to fire on work it owns.
- Two assets give conflicting instructions, or one asset silently duplicates another.
- An asset's description summarizes its workflow, so consumers act on the description and skip the body.
- The loaded file has grown past what it needs and reference material must move behind a pointer.
- Provenance, license or version information for kit content is missing or wrong.

## When not to use

- The change belongs to a consumer project's application code or its instructions file.
- The only evidence is a single reply someone liked better. Collect a reproducible case first.
- The request is to evaluate whether a change helped. That is harness-evaluation.
- The intent is to patch the generated package under `dist/`. Change the sources.

## Inputs

| Input | If missing |
| --- | --- |
| A concrete use case and expected behavior; for repairs, observed failure or conflicting instructions | Derive a case from the authorized request. For a new asset, current behavior may be absent. For a static contradiction, cite both instructions and the affected scenario; do not claim a model reproduction you did not run. |
| The asset that owns the behavior | Trace it: which trigger fired, which file was loaded, which instruction produced the output. Name the owner before editing anything. |
| Kit version or the identity of the harness checkout | Record the checkout state you worked from, and mark the version as unknown rather than guessing. |
| Existing baseline results for the asset | Capture the before state. If model runs are unavailable, distinguish static correction from unmeasured behavioral improvement and report that limit. |
| Provenance of any borrowed material: origin and license | Do not incorporate it. Naming a source does not make its content available. Write the procedure yourself or leave the gap recorded. |

## Procedure

1. **Establish the case in the harness checkout.** Work on sources. For a repair, reproduce the failure or cite the contradictory contract and a concrete scenario. For a new asset, record the authorized behavior currently missing. Preserve the before state and distinguish observed runs from static analysis.
2. **Locate the owning asset.** Follow the path the runtime took: which description matched, which body loaded, which instruction the output followed. A failure in routing belongs to a description; a failure in method belongs to a body; a failure in scope belongs to a limit.
3. **Choose the asset type by what the change must do.** An agent is a responsible party with its own context and limits. A skill is the reusable procedure. A command is the user's entry point and carries routing, parameters and contract only, with the method living in the skill. If the fix needs a new type, say so explicitly rather than stretching an existing one.
4. **Write and store the test cases before editing anything.** At least one positive and one negative case under `evals/cases/`, so harness-evaluation can capture the baseline on the unchanged asset. Running and grading them is harness-evaluation.
5. **Write or repair the contract for that type.** Use the contract tables in Quick reference. Every field is either filled or explicitly marked not applicable with a reason.
6. **Write the description as triggering conditions only.** Third person, starting with the conditions that should load the asset, plus the contraindications. Never summarize the procedure in the description: a description that describes the workflow becomes a shortcut the consumer follows instead of reading the body. Keep it under 500 characters.
7. **Match the form of the guidance to the form of the failure.** Use the failure-to-form table. A rule skipped under pressure needs an explicit counter and a rationalization table. Output of the wrong shape needs a positive recipe stating what the output is, not a list of prohibitions. A missing element needs a required slot in the template. Behavior that depends on a condition needs a conditional keyed to something observable.
8. **Prefer the positive target over the prohibition.** Naming the forbidden behavior makes it more available, not less. Keep a prohibition only as a hard guardrail that cannot be phrased positively, and pair it with the target behavior.
9. **Apply progressive disclosure.** The loaded file holds what every path through it needs. Material only some paths reach, and any block of reference over roughly 100 lines, moves to `references/<topic>.md` inside the skill's own folder and is reached by a folder-relative link, which resolves both in the source tree (`.skills/`) and in the generated package (`skills/`). Keep each meaning in one place; the same rule in two files is a maintenance defect and inflates its apparent importance.
10. **Keep the change minimal and the identifiers stable.** Change the smallest thing that fixes the reproduced case. Do not restyle neighboring assets in the same pass. Renaming an identifier is a breaking change and needs its own reason.
11. **Record provenance in the frontmatter.** Preserve `author: malka`. Mark the asset authorial when the procedure is your own or adapted when it synthesizes outside material into something different. Record identifiable origins and verify applicable licenses before incorporating protected text. This kit's policy forbids verbatim copied skills even when licensed; attribution alone does not make a copy an adaptation.
12. **State compatibility and propose a version.** Name what breaks for a consumer already using the kit: a renamed identifier, a changed output shape, a narrowed trigger, a removed reference. Propose the version and mark it as a proposal.
13. **Regenerate the package from the sources.** Never independently edit prompts under `dist/`. From the kit root run `go run ./cmd/dh build` (or `bin/dh build` from a built package). It validates sources, stages a complete package, then replaces `dist/claude-code/dev-harness`. If source and generated content disagree, source wins and the build must be re-run.
14. **Report, do not record.** Case, change, evidence and compatibility go back to the Coordinator, who is the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md` and `.harness/RISKS.md`.

## Output format

```
# Asset change: <asset id>

## Case
- Prompt:
- Expected:
- Observed:
- Reproduced: yes / no (<reason>)

## Owning asset
- File: <relative path> | type: agent / skill / command / tutorial
- Why this asset owns the behavior:

## Change
- What changed (smallest edit):
- What deliberately did not change:

## Contract check
| Field | State |
| --- | --- |
| <field from the contract table> | filled / not applicable (<reason>) |

## Test cases
- Positive: <relative path under evals/cases/> | must now <behavior>
- Negative: <relative path under evals/cases/> | must still <behavior>

## Provenance
- Class: authorial / adapted
- Origin:
- License verified: yes / no / not applicable

## Compatibility
- Breaks for existing consumers:
- Identifiers renamed:

## Version proposal
- <value> (proposal, needs confirmation)

## Evidence
- <relative path> - <what it establishes>
```

## Quick reference

Contract per asset type. Every field is filled or marked not applicable with a reason.

| Asset | Required contract |
| --- | --- |
| Agent | Mission; positive triggers; negative triggers; minimum inputs and missing-input behavior; procedure; tool policy (omit `tools` to inherit everything; allowlist only for read-only roles; no `disallowedTools` — the runtime did not enforce it (RISK-001, 2026-09-21), so the rule that only the Coordinator dispatches lives in the prompt body of every specialist); limits; output format; context handoff; stop condition; evaluable examples; explicit default model (never inherit), with scoped Coordinator override and effective model recorded or unverified |
| Skill | Frontmatter with name matching the folder, triggering description, author, provenance; overview with the core principle; when to use; when not to use; inputs; numbered verifiable procedure; output format; quick reference; common mistakes; one worked example; related roles, commands and sibling skills; proof case |
| Command | Routing to the owning role and skill; parameters; expected output contract; missing-prerequisite behavior; next action. The method stays in the skill |
| Tutorial | Prerequisites; the exercise; the completion condition; what the reader must be able to do unaided at the end |

Match the form of the guidance to the form of the failure.

| Baseline failure | Right form | Wrong form |
| --- | --- | --- |
| Knows the rule, skips it under pressure | Explicit counter plus a rationalization table | Soft wording like prefer or consider |
| Complies but the output has the wrong shape | Positive recipe naming the parts of the output in order | A list of prohibitions |
| Omits a required element from something already produced | A required slot in the template | A prose reminder near the template |
| Behavior should depend on a condition | Conditional keyed to an observable predicate | Unconditional rule with exemption clauses |
| Fires on work it does not own | Sharpen the contraindications in the description | Add caveats in the body, which loads too late |

Where material lives.

| Material | Location |
| --- | --- |
| Triggering conditions | Frontmatter description only |
| Procedure every path needs | The loaded file |
| Reference only some paths reach, or over roughly 100 lines | `references/<topic>.md` inside the skill's own folder, linked folder-relatively |
| Something a script can verify mechanically | A script, not prose |
| A fact the environment already states | Leave it in the environment; a restatement goes stale |

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Editing an asset because one reply read better | Kit rules get rewritten by taste and drift away from the cases they were built for | Require a reproduced case with prompt, expected and observed |
| Description that summarizes the workflow | Consumers follow the summary and never read the body, so the body becomes decoration | Triggering conditions and contraindications only |
| Hand-editing the generated package | The next build silently reverts the fix and the sources never learn about it | Change the sources and rebuild |
| Fixing routing by adding text to the body | The body loads only after the routing decision was already made | Fix the description, which is what the decision reads |
| Big rewrite when a clause was wrong | The change cannot be attributed to the case, so nobody can tell what fixed what | Smallest edit that fixes the reproduced case |
| Renaming identifiers while fixing behavior | Two changes land together and every consumer reference breaks for an unrelated reason | Keep identifiers stable unless the rename is the fix |
| Copying text from another source because it is good | Unverified license and unresolved references travel with it | Authorial or adapted, with origin and license recorded |
| Naming an external tool the kit does not ship | The asset breaks on any machine that lacks it | Phrase it as available-or-fallback, and name the fallback |
| The same rule written in two assets | They drift, and the contradiction surfaces mid-task | One source of truth, referenced by relative path |
| Absolute paths from the author's machine | The kit stops being portable the moment it moves | Paths relative to the kit root or the consumer project |
| Only a positive test case | A widened trigger passes it and starts capturing work it does not own | Always pair it with a negative case |

## Example

Input: the mapping role explores the whole tree for a question about one known file.

```
Case
  Prompt:   "What does src/auth/token.ts export?"
  Expected: read that file, answer, stop.
  Observed: full-tree map, 40 files read.

Owning asset: .agents/repo-scout.md, description field.
  The body already says "smallest sufficient map". The description
  matched before the body was ever loaded, so the body could not help.

Change (one line, description only):
  add "Do not use when the user names a single known file and the
  question is answerable from it."

Test cases
  Positive: broad question with no file named -> the role still loads.
  Negative: the prompt above -> the role does not load; the file is read directly.

Compatibility: no identifier changed, trigger narrowed. Proposal: patch version.
```

## Related

Roles: harness-maintainer, coordinator, qa-verifier. Command: `/dh:improve`. Skills: harness-evaluation, context-handoff, delivery-readiness. Templates: the `templates/` folder bundled with the kit (`templates/<lang>/MEMORY.md`).

## Proof case

Given a reproduced routing failure, the flow names the owning asset and the field responsible, makes a single minimal edit in the harness checkout, and leaves the generated package untouched until it is rebuilt from the sources. The asset ends with a description that states only triggering conditions, a complete contract for its type, a positive and a negative test case stored for replay, a provenance class with its origin, and a stated compatibility impact with a proposed version. A request to change an asset based on one isolated preferred reply is refused with the reason, and no rule is edited.
