# PRD-000 · <short title>

<!-- Suggested location in the consumer project: .harness/prd/PRD-000.md
Fill in only what changes the decision. Delete the comments when finished.
Fact, hypothesis, and decision are marked as such. No secrets. -->

| Field | Value |
|---|---|
| Status | draft / in review / approved / delivered on <date> / canceled |
| Owner | <who decides the scope> |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |
| Tickets | T-000, T-000 |

## 1. Problem

<!-- Who suffers, what does not work, in what situation, what the impact is, what needs to change. One piece of evidence (log, metric, account) is worth more than an adjective. One capability per PRD: two independent capabilities become two PRDs before the first ticket. -->

## 2. Solution

<!-- What becomes true after delivery: observable behavior, outcome for the user, expected flow, system response. Not an architecture description. -->

## 3. Rules

<!-- Non-negotiable constraints, numbered R1, R2, R3... Each rule must pass this test: it is a mandatory condition or an explicit prohibition, and it is verifiable. Generic text, a preference, or an intention is not a rule. -->

- R1 <mandatory condition or explicit prohibition, verifiable>
- R2 <mandatory condition or explicit prohibition, verifiable>

## 4. Docs

<!-- Paths of already-documented flows this feature touches, one per line. If the feature touches no documented flow, delete the example line below and write exactly the line: "No documented flow affected." An empty section does not count — document-validator reports it as a gap. -->

- <path/to/flow.md>

## Appendix (optional)

<!-- Only when the author judges it necessary. The absence of this appendix does not block turning the PRD into tickets, the same way the Docs section already works when it has no real items (the "No documented flow affected." line). -->

### Discarded alternatives

| Alternative | Why not |
|---|---|
| <option> | <cost, risk, or limitation> |

### Risks and pending decisions

- Risk: <description> · mitigation: <action>
- Pending decision: <question> · answered by: <who> · blocks: <RF/ticket>
