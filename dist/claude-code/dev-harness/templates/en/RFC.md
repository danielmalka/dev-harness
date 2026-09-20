# RFC-000 · <short title>

<!-- Suggested location: .harness/rfc/RFC-000.md
An RFC proposes a significant change and collects objections BEFORE building.
When approved, the durable decision becomes an ADR; the work becomes stories/tasks.
Keep it short: a reader should be able to object in one pass. -->

| Field | Value |
|---|---|
| Status | draft / under discussion / approved / rejected / withdrawn |
| Author | <who proposes> |
| Reviewers | <who needs to weigh in> · deadline: YYYY-MM-DD |
| Origin | PRD-000 / ST-000 / incident RISKS.md#id / technical debt |
| Outcome | ADR-000 / ST-000 (filled in when closing) |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |

## 1. Summary

<!-- Three sentences: the problem, the proposal, what changes for whom. -->

## 2. Motivation

<!-- Why now. Evidence of the problem (metric, incident, cost). What happens if we do nothing. -->

## 3. Proposal

<!-- The change in concrete terms: components, flow, contracts, data. Text diagram if helpful. No implementation code, only interfaces that matter. -->

- Affected components: `<path or module>`
- Proposed flow: <step → step → step>
- New or changed contracts: <API, event, schema; compatibility>
- Data: <entity, migration, backfill, reversibility>

## 4. Alternatives

| Alternative | Advantage | Why not |
|---|---|---|
| Keep as is | | |
| <option B> | | |

## 5. Impacts

<!-- Fill in only what applies. -->

- Security / sensitive data: <...>
- Operations: <deploy, rollout, feature flag, observability, rollback>
- Dependencies: <lib to install / update / remove; minimum version>
- Compatibility: <what breaks; migration plan for consumers>
- Cost: <estimated effort or "not estimated">

## 6. Adoption plan

<!-- Slices in order, each verifiable. Name the first slice that proves the idea. -->

1. <slice> → proof: <check>
2. <slice> → proof:
3. <slice> → proof:

## 7. Risks and open questions

- Risk: <what may go wrong> · mitigation: <action> · related incident: <RISKS.md#id or none>
- Open question: <question> · answered by: <who> · blocks approval: yes / no

## 8. Discussion

<!-- Objections and responses, dated. Do not delete a resolved objection; mark it as resolved. -->

- YYYY-MM-DD · <who> · <objection> → <response> · resolved / open

## 9. Decision

<!-- Filled in when closing. -->

- Outcome: approved / rejected / withdrawn · date: YYYY-MM-DD · by: <who>
- Reason: <one sentence>
- Follow-ups: ADR-000, ST-000
