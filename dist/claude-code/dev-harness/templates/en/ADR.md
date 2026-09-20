# ADR-000 · <decision in one sentence>

<!-- Suggested location: .harness/adr/ADR-000.md
An ADR records ONE durable decision and the reason. It is not a design doc: to
propose and discuss, use RFC. Write only enough for someone a year from now to
understand why it was done this way and when it should be reviewed. -->

| Field | Value |
|---|---|
| Status | proposed / accepted / superseded by ADR-000 / revoked |
| Date | YYYY-MM-DD |
| Decision maker | <who approved> |
| Origin | PRD-000 / ST-000 / RFC-000 / incident RISKS.md#id |
| Reversibility | cheap and local / expensive or many dependents |

## 1. Context

<!-- The situation requiring the decision and the forces that really separate the options (do not list everything). Fact marked as fact, hypothesis as hypothesis. -->

- Question: <what changes if we take the other path>
- Forces: <technical constraint, data, operations, security, timeline>

## 2. Options considered

<!-- Include "keep as is" when viable. Fill in only the dimensions that discriminate. -->

| Dimension | A · <name> | B · <name> | C · keep as is |
|---|---|---|---|
| Complexity | | | |
| Build cost | | | |
| Operating cost | | | |
| Data impact | | | |
| Security impact | | | |
| Team familiarity | | | |

## 3. Decision

- Chosen: <option>
- Decisive reason: <the reason that broke the tie>
- Discarded: <option> because <specific reason>
- Declared preference: <factor that is a preference or habit, named as such> or none

## 4. Consequences

- Becomes easier: <...>
- Becomes harder: <...>
- Debt assumed: <what we accept living with and until when>
- Affected contracts: <component A → B; data that crosses; compatibility>

## 5. Validation

- Check that proves it works: <test, metric, or command>
- Signal that it is failing: <what to observe in operation>

## 6. Review when

<!-- Observable condition, not a vague date. Without this, the ADR becomes permanent by accident. -->

- <e.g. volume exceeds N, lib X releases version Y, incident in the area>

## 7. References

- <relative path: RFC, PRD, incident, benchmark>
