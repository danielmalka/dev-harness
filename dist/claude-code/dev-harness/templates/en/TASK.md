# T-000 · <short title>

<!-- Location: .harness/tasks/T-000/TASK.md
This is the single file a builder agent receives: it already carries the
story, the expected flow, and a literal copy of the rules it must honor.
It does not read the whole PRD. Fluxo/Flow (section 3) plus Pronto
quando/Done when (section 6) are the only acceptance source for this
ticket — there is no separate "Acceptance criteria" section.

For Tipo/Type: bug, there is no dedicated defect section: the Flow
section carries the reproduction (numbered steps ending in the
observed-vs-expected pair) and the Done when section gains the
regression-test item. See the comments in those two sections below.

Delete the comments when finished. -->

| Field | Value |
|---|---|
| Type | task / bug |
| Status | ready / in progress / blocked / in review / done |
| PRD (RF-<n>) | PRD-000 (RF-01) |
| Suggested role | backend-builder / frontend-builder / data-engineer / devops-engineer / qa-verifier / ... |
| Depends on | T-000, <contract, decision, migration> |
| Blocks | T-000 |
| Authorization | <what may be written; commit/push/deploy require explicit authorization> |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |

## 1. Lane

<!-- One value per ticket, closed list. A story that needs an API and a
screen with independent responsibilities becomes two tickets, one per
lane; "and test it too" at the end of this ticket does not by itself
create a test ticket. -->

Lane: backend / frontend / dados / infra / teste (unitário) / teste (integração) / teste (unitário e integração)

## 2. Story

<!-- One sentence, always in this shape. -->

As a <actor>, I want <action>, so that <result>.

## 3. Flow

<!-- Numbered, observable steps, "When X, then Y" — actor is user, system,
or test. This block, together with Done when (section 6), is the only
acceptance source for the ticket.

For Type: bug, this field carries the defect's reproduction: the same
numbered "When X, then Y" steps end in the observed-vs-expected pair —
what happens today (the defective behavior) and what should happen
instead. Do not create a separate section for this. -->

1. When <action>, then <observable result>.
2. When <action>, then <observable result>.

## 4. Rules

<!-- Literal copy, no paraphrase and no summary, only of the R<n> rules
from the PRD this ticket must honor. No link. A rule that does not fit
here does not enter. -->

- R1 <literal copy of the matching PRD rule>

## 5. Docs

<!-- Paths to update in this ticket, inherited and refined from the PRD's
Docs section, or the exact line "none". -->

- <path/to/doc.md>

## 6. Done when

<!-- The ticket's minimum content: the flow in section 3 happened; each
rule copied in section 4 has evidence; the relevant tests are green; the
documentation listed in section 5 was updated.

For Type: bug, this field gains one more item: a regression test for
this defect is green. -->

- the flow described in section 3 happened
- the rules copied in section 4 have evidence
- the relevant tests are green
- the documentation listed in section 5 was updated (or the section says "none")

## 7. Result

<!-- Filled in by the executor on delivery. Only what was done and executed. -->

- Conclusion: <what was delivered, in one sentence>
- Files changed: `<path>`, `<path>`
- Checks: `<command>` → passed / failed / not-run (<reason>)
- Limitations: <what was not verified or remains pending>
- Pending items for the Coordinator: <decision, incident, memory update>
- Next step: <concrete action> · authorization required: <which>

## 8. Risks and limits

- Risk: <what may break> · mitigation: <action>
- When encountering an open decision: stop and report; do not invent behavior.
- After 6 unsuccessful correction rounds: return to the Coordinator with evidence.
