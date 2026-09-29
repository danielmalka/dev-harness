# New PRD flow oriented around subagents

Portuguese version: [docs/novo-fluxo-prd.md](../novo-fluxo-prd.md)

Proposal consolidated on 2026-09-28. This does not yet replace `templates/` or the current kit flow.

One small feature equals one PRD. The implementing subagent does not read the entire PRD. It reads the ticket, which already carries the story, the flow, and the copy of the rules that affect it.

Three barriers close the process: an unapproved PRD does not become tickets, an unreviewed ticket does not close, and a feature is not delivered without a final review.

## 1. Input

The requester may start with any of the following:

- an idea
- a request
- an imagined solution
- a bug
- an improvement
- a user complaint
- a perceived technical need

The input does not need to be complete. It is a signal for investigation, not a finished PRD.

## 2. PRD creation

The orchestrator calls the **PRD Creation Subagent**.

This subagent grills the requester before writing. The grill continues until the capability is understandable and important decisions are explicit. It is not a bureaucratic interview.

The grill must clarify:

- what problem actually exists
- who is affected
- what behavior must change
- how the user will perceive the change
- which scenarios must work
- which behaviors must not happen
- which limits cannot be crossed
- which existing flows or documents are affected

Rule amendments also belong to this subagent. The requester may ask for changes.

### PRD structure

Only these sections, in this order. If one is missing, the PRD does not become tickets.

#### Problem

Who suffers, what does not work, in which situation, what is the impact, and what must change.

One capability. Two independent capabilities become two PRDs before the first ticket.

#### Solution

What becomes true after delivery.

Prioritize behavior, user outcome, expected flow, and system response. This is not an architecture description.

#### Rules

Non-negotiable constraints, numbered as `R1`, `R2`, `R3`.

Each rule must be verifiable, observable, a required condition, or an explicit prohibition. Generic text, preference, or intent is not a rule.

```text
R1. A user without permission cannot view the report.
R2. The system must record every denied attempt.
R3. The API cannot return another user's data.
```

#### Docs

Paths of documented flows touched by the feature.

When none are affected, use this exact line:

```text
No documented flow affected.
```

An empty section does not count. If the feature changes an already documented flow, updating the documentation is part of the delivery.

## 3. Adversarial PRD review

The orchestrator calls the **Adversarial PRD Reviewer**.

The reviewer does not improve the text based on stylistic preference. It looks for failures that could cause an incorrect implementation.

It checks:

- Is the problem clear?
- Is there only one capability?
- Does the solution describe observable behavior?
- Are the rules truly non-negotiable?
- Can every rule be verified?
- Is any rule missing?
- Is there a contradiction between the solution and the rules?
- Does the flow leave an important scenario unanswered?
- Does the PRD contain an implicit decision?
- Is the Docs section correct?
- Can it be split into tickets without inventing requirements?

Result:

```text
APPROVED
```

or:

```text
REJECTED

Reasons:
- ...

Required corrections:
- ...
```

A rejected PRD returns to the PRD Creation Subagent. It advances only after approval.

## 4. Ticket decomposition

Once the PRD is approved, the orchestrator calls the **PRD Decomposition Subagent**.

This role is separate from the PRD author, so rule writing and work decomposition are not mixed.

This subagent:

- splits the work by lane
- creates backend, frontend, and test tickets when needed
- copies the relevant rules literally
- identifies the necessary tests
- identifies documentation to update
- preserves the expected flow
- adds no scope that is not in the PRD

A story that needs both an API and a screen becomes two tickets when the responsibilities are independent. “Also test it” at the end of a backend ticket is not a test ticket.

Lanes:

```text
backend
frontend
test (unit)
test (integration)
test (unit and integration)
```

One lane per ticket.

## 5. Ticket structure

### Lane

```text
Lane: backend
```

or:

```text
Lane: test (integration)
```

### Story

```text
As a <actor>,
I want <action>,
so that <outcome>.
```

### Flow

What the actor does and what happens next, in observable steps. The actor may be a user, system, or test. This block is the acceptance basis.

```text
1. The user submits the valid form.
2. The system creates the record.
3. The system returns confirmation.
4. The record appears in the list.
```

### Rules

Literal copies of only the rules this ticket must honor. No link. No summary. A rule that does not apply does not belong here.

```text
R1. A user without permission cannot view the report.
R3. The API cannot return another user's data.
```

### Docs

Paths to update in this ticket, or `none`. It inherits the PRD Docs section and is refined per ticket.

```text
Docs:
- docs/flows/report.md
```

or:

```text
Docs:
- none
```

### Done when

The flow happens and every copied rule has evidence.

```text
Done when:
- the described flow happens
- rules R1 and R3 were verified
- relevant tests are green
- the listed documentation was updated
```

### Minimal example

```text
Lane: test (integration)
Story: As a visitor, I want to be rejected when I submit an empty password, so that no session is created.
Flow:
1. The visitor submits login with an empty password.
2. No session exists.
3. The visitor sees "password required".
Rules:
R1. An empty password does not authenticate.
R2. The message is "password required" and does not reveal whether the email exists.
Docs: none
Done when: the test cites R1 and R2 and is green.
```

## 6. Ticket execution

The orchestrator calls a **specialized Development Subagent** according to the ticket lane. It checks the evidence. It does not implement the ticket itself.

The subagent receives the complete ticket, the minimum repository context, the copied rules, the expected files or areas, and the evidence criteria. It does not need to interpret the entire PRD.

It:

1. executes only the received ticket
2. respects every copied rule
3. does not invent scope
4. updates the listed documentation
5. runs the required checks
6. returns real evidence

## 7. Agentic loop

```text
Development Subagent
        |
        v
Test and Review Subagent
        |
        v
Approved or Rejected
```

### Development

Executes the ticket and returns:

- changed files
- implemented behavior
- executed tests
- commands used
- command results
- limitations or questions
- visual evidence, when applicable

A self-report is not evidence. The orchestrator verifies the result.

### Test and review

The reviewer checks the ticket against the Story, Flow, copied Rules, Done when criteria, project conventions, required tests, and listed documentation.

Approved:

```text
APPROVED

Evidence:
- ...

Tests:
- command: ...
- result: ...

Docs:
- ...
```

Rejected:

```text
REJECTED

Problems:
- ...

Affected rule or criterion:
- R2
- Done when, item 3

Required correction:
- ...
```

A rejection returns to the Development Subagent with the specific correction. After correction, the same ticket goes through Test and Review again.

The ticket does not close with a missing test, failing test, unverified rule, pending documentation, behavior different from the flow, or evidence that was only asserted.

## 8. Final feature review

Once all tickets are approved, the orchestrator announces the end of development and calls feature-review subagents.

This stage reviews the whole feature, not each isolated ticket.

It may call reviews for:

- backend/frontend integration
- functional behavior
- tests
- documentation
- security
- regressions
- PRD adherence

When needed, it dispatches test and documentation subagents.

If a screen is accessible, the PRD Creation Subagent, or another subagent with screen access, may use the app as a client. This happens only when a visual surface exists. It does not replace backend or integration tests.

Using the app as a client means:

1. access the screen
2. act as a user
3. follow the flow described in the Solution
4. verify the result
5. record evidence
6. test relevant invalid or prohibited behaviors as well

At this stage, the PRD Creator evaluates whether the result matches the original intent. It does not write new requirements.

Result:

```text
FEATURE APPROVED
```

or:

```text
FEATURE REJECTED

Reasons:
- ...

Tickets or areas to reopen:
- ...
```

A rejection reopens only the affected tickets or areas and restarts the required loop.

## 9. Responsibilities

- **Requester:** supplies an idea, request, solution, bug, or equivalent. It does not need to be complete. May request a rule change.
- **Orchestrator:** calls creation, adversarial review, decomposition, specialized development, and final review. Checks evidence. Does not implement tickets.
- **PRD Creation Subagent:** runs the grill, writes Problem, Solution, and Rules, and amends rules. At the end, when a screen exists, may test the execution as a client.
- **Adversarial PRD Reviewer:** evaluates structure, rules, and content. Rejects what is not aligned.
- **PRD Decomposition Subagent:** breaks the PRD into tickets and copies rules.
- **Specialized Development Subagent:** executes one ticket in its lane and returns evidence.
- **Test and Review Subagent:** approves or rejects the ticket. A rejection returns it to development.
- **Feature Review Subagents:** review the feature at closing. Test and documentation are included when needed.

## 10. Guardrails

1. Incomplete input may start the process. An incomplete PRD cannot start execution.
2. Each PRD is one capability. Two capabilities become two PRDs before the first ticket.
3. The PRD is created after the requester grill.
4. Rule amendments belong to the PRD Creation Subagent. The requester may request a change.
5. A PRD becomes tickets only after adversarial approval.
6. A ticket contains only requirements present in the PRD.
7. A ticket rule is a literal copy of the PRD rule.
8. Every observable rule has a test ticket that cites its ID. Without this, the PRD does not enter execution.
9. Backend and frontend do not close while the corresponding rule test is missing or red.
10. A ticket that changes a documented flow updates that document in the same delivery. Stale documentation keeps the ticket open.
11. If implementation finds a documented flow outside Docs, work stops, the PRD is corrected, and only then does work continue.
12. A newly discovered rule stops the work. The PRD is updated before continuing.
13. A changed rule rewrites tickets not yet started. A ticket in progress stops.
14. Two tickets changing the same file do not run concurrently.
15. A rejection returns to development with a specific correction. The same ticket is reviewed again.
16. A self-report is not proof. Evidence is a green test, exit code, URL, or reproduced flow step.
17. Final review checks the feature against the PRD, not only against tickets.
18. If a screen exists, the main flow is tested as a client.
19. A new desire becomes a PRD amendment or a new PRD. It does not enter a ticket as informal scope.
20. “Done” does not deliver. Delivery requires checked evidence and an Executed Summary.

## 11. Closing the PRD

The PRD closes only after all tickets are approved, relevant tests have run, documentation is updated, final review is approved, and the orchestrator has checked the evidence.

The orchestrator writes this block into the PRD.

### Executed Summary

**Delivered.** One sentence describing what became true, compared with the Solution.

**Rules.**

```text
R1. Honored.
Evidence: ...

R2. Honored.
Evidence: ...

R3. Not honored.
Reason: ...
Pending decision: ...
```

**Tickets.**

```text
- TKT-001, backend, approved
  Evidence: ...

- TKT-002, frontend, approved
  Evidence: ...

- TKT-003, integration test, approved
  Evidence: ...
```

**Docs.**

```text
- docs/flows/report.md, updated
- docs/api/report.md, updated
```

or:

```text
No documented flow was affected.
```

**Out of scope.** What the original Solution promised but the delivery did not include.

A rule deviation without an explicit requester decision means the PRD is not delivered.

## Flow

```text
Requester input
        |
        v
Requester grill
        |
        v
PRD Creation Subagent
        |
        v
Adversarial PRD Reviewer
        |
        v
PRD approved?
   |-- no --> correct PRD
   |-- yes
        |
        v
PRD Decomposition Subagent
        |
        v
Tickets
        |
        v
Specialized development
        |
        v
Ticket test and review
   |-- rejected --> development
   |-- approved
        |
        v
Are all tickets approved?
   |-- no --> next ticket
   |-- yes
        |
        v
Final feature review
        |
        v
Tests, documentation, and client-on-screen test when available
        |
        v
Executed Summary in the PRD
        |
        v
PRD delivered
```
