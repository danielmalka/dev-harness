# PRD-000 · <short title>

<!-- Suggested location in the consumer project: .harness/prd/PRD-000.md
Fill in only what changes the decision. Delete the comments when finished.
Fact, hypothesis, and decision are marked as such. No secrets. -->

| Field | Value |
|---|---|
| Status | draft / in review / approved / canceled |
| Owner | <who decides the scope> |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |
| Stories | ST-000, ST-000 |

## 1. Problem

<!-- Current situation, who suffers from it, how often, and what it costs. One piece of evidence (log, metric, account) is worth more than an adjective. -->

## 2. Expected outcome

<!-- What changes for the user when it is ready. How we will know: 1 to 3 observable signals, with a baseline if one exists. -->

- Signal 1: <metric or behavior> · today: <value> · target: <value>

## 3. Users and scenarios

| User | Scenario | Today | After |
|---|---|---|---|
| <profile> | <what they try to do> | <what happens> | <what will happen> |

## 4. Scope

**In scope**
- RF-01 <verifiable behavior in one sentence>
- RF-02

**Out of scope** (and why, when it is not obvious)
- <item>

## 5. Acceptance criteria

<!-- One per requirement. Format: When <condition>, then <observable result>. -->

- AC-01 (RF-01) When <X>, then <Y>.
- AC-02 (RF-02) When <X>, then <Y>.

## 6. Constraints and technical impact

<!-- Only what constrains the solution. Leave blank what does not apply. -->

- Data: <new/changed entities, migration, retention, LGPD>
- Contracts: <affected APIs, events, or integrations; required compatibility>
- Security: <authentication, authorization, sensitive data>
- Operations: <environments, feature flag, rollout, observability>
- Stack and dependencies: <new lib, minimum version, removal>
- Timeline / cost: <if any>

## 7. Discarded alternatives

| Alternative | Why not |
|---|---|
| <option> | <cost, risk, or limitation> |

## 8. Risks and pending decisions

<!-- Known risk = what may go wrong + mitigation. Pending decision = question whose answer changes the scope + who answers + by when. Consult .harness/RISKS.md if the area has a recorded incident. -->

- Risk: <description> · mitigation: <action>
- Pending decision: <question> · answered by: <who> · blocks: <RF/ST>

## 9. References

- <relative path or link: research, prototype, incident, ADR>
