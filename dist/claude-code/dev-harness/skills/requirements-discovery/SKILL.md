---
name: requirements-discovery
description: Use when a request names a pain, an audience or an outcome but no testable behavior, when two readings of the same ask would produce different software, or when acceptance criteria must exist before a plan can be written. Do not use when acceptance already exists, when the user asked for a plan, a patch or a review, or when the only missing piece is knowing where the code lives.
author: malka
metadata:
  provenance: adapted
  sources: ["product-discovery role", "brainstorming", "grilling", " discovery-agent"]
---

# Requirements Discovery

## Overview

Turn a vague ask into five things a builder can act on: the current situation, the desired outcome, concrete examples, the alternatives considered, and the limits. The output is a brief short enough to read in one sitting and specific enough to plan from. Questions cost the user's attention, so only decisions that would change what gets built are worth asking about; everything findable is found, not asked. This skill writes discovery artifacts only: brief, acceptance, open decisions. Never application source, configuration or test files.

## When to use

- The request describes a pain, a user or a goal, but no observable behavior.
- Two honest readings of the same sentence would lead to different builds.
- Nobody can say what "done" looks like without inventing it mid-build.
- The ask is large enough that it may be several independent pieces rather than one.
- A previous attempt shipped the wrong thing and the cause was ambiguity, not execution.

## When not to use

- Acceptance criteria already exist and are testable. Go to implementation-planning.
- The user asked for a patch, a plan, a review or a map.
- The unknown is the codebase, not the intent. Use repository-mapping.
- A brief exists and the goal has not changed. Do not reopen a closed brief for polish.

## Inputs

| Input | If missing |
|---|---|
| The stated goal, in the user's words | Nothing to discover. In the main session, ask the owner what should become possible that is not possible now. As a dispatched specialist, return the missing input to the Coordinator and stop; never proceed on an assumed value. |
| Who is affected | Make it the first question, since it changes every later answer. |
| Known constraints: deadline, stack, out of scope | Proceed and record them as unknown; do not invent a deadline or a stack. |
| Existing product docs or a prior brief in the project | Search before asking. An answer already written down is not a question. |

## Procedure

1. **Size the ask before refining it.** If the request contains several independent subsystems, say so immediately and propose splitting it. Refining details inside a request that needs decomposition wastes every question that follows. Each piece then gets its own brief.

2. **Separate facts from decisions, and never mix the two.** Facts about the repository, existing behavior, current data shape and prior art are yours to find with reading and search. Decisions about who this serves, what happens in the awkward case, and what is out of scope belong to the user. Asking the user for a fact you could look up spends their attention on the wrong thing; answering a decision yourself is the failure mode that produces confident, wrong software.

3. **Build the question frontier.** List every open decision, then keep only the ones whose prerequisites are already settled. A question that depends on an unanswered question belongs to a later round. Order the survivors by how much they change the build.

4. **Apply the scope test to each candidate question.** Would two different answers produce different acceptance criteria, different files, or different tests? If not, drop it or decide it yourself and label the choice as a suggestion. Curiosity is not a reason to ask.

5. **Ask the smallest useful question and wait for required decisions.** Include a recommendation, but silence is never an answer or authorization. Continue independent research while a blocking scope decision is pending. Two questions may share one message only when the second is meaningless without the first; otherwise one blocking question at a time. Record the actual answers in the brief.

6. **Fill the five dimensions.** Current situation, what actually happens today, with evidence where the project can show it. Desired outcome, in behavior, not in feelings. Examples, at least one normal case and one awkward case. Alternatives, including doing nothing and the cheaper partial version. Limits, meaning explicitly out of scope, non-negotiable constraints, and the data or permissions involved. Record each claim about current behavior you could not verify as a hypothesis, with the observation that would falsify it.

7. **Write acceptance as "when X, then Y".** Each item observable by a person or a test, with X naming a trigger and Y naming a result. Mark each one `requirement` or `suggestion`. A requirement is something the user asked for or that correctness demands; anything you thought of that nobody requires is a suggestion, and shipping suggestions as requirements is scope creep with a straight face.

8. **Cover the awkward cases proportionally.** Empty, duplicate, expired, unauthorized, concurrent, and too large. Include only those the feature actually meets. An invented edge case is scope, not diligence.

9. **Record what is still open.** Each open decision gets the question, why it changes scope, and your recommendation. Do not resolve an open product decision by picking quietly; do not stall the whole brief on one unanswered item either. If a question cannot be settled by talking because it needs something to look at, say so and route it to a prototype instead of rephrasing it forever.

10. **Stop when acceptance is testable and the open decisions are explicit.** Hand the brief to the Coordinator, along with any constraint or incident worth remembering. Only the Coordinator writes `.harness/MEMORY.md`, `.harness/EPOCHAL.md` and `.harness/RISKS.md`; discovery reports facts and never edits those records.

## Output format

```
## Brief
- Problem: <what is wrong today, observable>
- Users affected: <who, and what they do now instead>
- Desired behavior: <what becomes possible>
- Out of scope: <explicitly not now>

## Examples
- Normal: given <input>, the user sees <result>
- Awkward: given <edge input>, the user sees <result>

## Alternatives considered
| Option | Cost | Why not chosen |
|---|---|---|
| Do nothing | | |
| <cheaper partial version> | | |

## Limits and constraints
- <constraint> (source: user | project instructions | observed)

## Acceptance
- When <X>, then <Y>. [requirement]
- When <X>, then <Y>. [suggestion]

## Hypotheses
- <claim>, how it would be shown wrong

## Open decisions
- <question>, why it changes scope, recommended answer

## Evidence
- <relative path> or "none yet"
```

When persisted, the brief becomes `.harness/prd/PRD-<n>.md` from the kit template `templates/<lang>/PRD.md`; keep the requirement/suggestion marking in the text.

## Quick reference

| Vague phrase | The question that makes it testable |
|---|---|
| "make it faster" | Which operation, measured how, and what number is acceptable? |
| "users get confused" | At which step do they stop, and what should happen instead? |
| "it should be secure" | Which data, which actor must be refused, and what does refusal look like? |
| "add reporting" | Which decision does the reader make from the report? |
| "it should just work" | Which case fails today, and what is the correct result there? |
| "like <other product>" | Which one behavior of it, specifically? |

| Test | Ask it | Do not ask it |
|---|---|---|
| Scope test | Two answers change the build | Two answers change only wording |
| Ownership test | It is the user's decision | It is discoverable in the repository |
| Frontier test | Its prerequisites are settled | It depends on an open question |
| Cost test | The answer is cheaper than building the wrong thing | The answer is guessable and low risk |

## Common mistakes

| Excuse | Reality |
|---|---|
| "I will ask everything upfront so I only interrupt once" | A question list gets one answer and a tired user. One at a time, ordered by impact, finishes faster. |
| "The stack question seems important" | It changes the build only if the user has a real constraint. Otherwise it is a decision for planning, not discovery. |
| "They did not mention errors, so there are none" | Unstated does not mean absent. Name the awkward case and get a decision on it. |
| "I thought of a nicer flow, so I will add it" | Your idea is a suggestion. Label it and let the user promote it. |
| "Acceptance is obvious, I will skip writing it" | Obvious acceptance is exactly what two people read differently. |
| "The scope is clear but I will run discovery anyway" | Reopening a settled brief costs a cycle and reintroduces closed questions. |
| "I could not reach the user, so I decided for them" | Record it as an open decision with a recommendation, and say the brief is provisional. |
| "It is one request, so it is one brief" | Several independent subsystems in one ask need decomposition before any refinement. |

## Example

Input: "people forget to follow up after demos".

```
## Brief
- Problem: after a demo, the next contact depends on the rep remembering; nothing in the product records that a demo happened.
- Users affected: sales reps, who currently keep notes outside the product.
- Desired behavior: a demo can be marked as held, and an unanswered one becomes visible to its owner.
- Out of scope: sending email on the user's behalf.

## Examples
- Normal: rep marks a demo held; two days with no logged reply puts it on their open list.
- Awkward: the contact replies the same day; the demo never reaches the open list.

## Alternatives considered
| Option | Cost | Why not chosen |
|---|---|---|
| Do nothing | none | the forgotten follow-up is the reported pain |

## Limits and constraints
- No outbound email on the user's behalf (source: user)

## Acceptance
- When a demo is marked held and no reply is logged within the configured window, then it appears on the owner's open list. [requirement]
- When a reply is logged, then the item leaves the list on the next load. [requirement]
- When the owner dismisses an item, then it does not return for that demo. [suggestion]

## Hypotheses
- Replies are already logged against the contact. Shown wrong by a contact with a reply that never reaches the product.

## Open decisions
- Is the window fixed or per user? It changes storage and settings. Recommended: fixed at first, configurable later.

## Evidence
- none yet
```

## Related

Roles: product-discovery, coordinator, implementation-planner, qa-verifier. Command: `/dev-harness:discover`. Skills: implementation-planning consumes this brief, architecture-decisions when a constraint here forces a durable choice, regression-testing turns the acceptance items into checks, repository-mapping when a claim about current behavior needs evidence.

## Proof case

A one-sentence vague request produces a brief with at least one normal and one awkward example, acceptance written as "when X, then Y" with requirement and suggestion marked separately, and at least one open decision stated with the reason it changes scope. The session asks no question whose two possible answers would produce the same build, and asks nothing that was already answered in the repository.
