# ST-000 · <short title>

<!-- Suggested location: .harness/stories/ST-000.md
One story = one end-to-end deliverable and testable behavior.
If it does not fit in a few tasks, split it. Delete the comments when finished. -->

| Field | Value |
|---|---|
| Status | draft / ready / in progress / blocked / done |
| PRD | PRD-000 (RF-01, RF-02) |
| Depends on | ST-000, <decision or contract> |
| Tasks | T-000, T-000, BUG-000 |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |

## 1. Behavior

<!-- One sentence: who does what to obtain what. Then the minimum context an agent needs to avoid inventing: business rule, real example, limit. -->

As a <user>, when <situation>, I want <action> to get <result>.

Context:
- <business rule or concrete example>

## 2. Usage scenarios

<!-- Given / When / Then. Cover: happy path, relevant error, empty or limit. Each scenario becomes a test. -->

### CN-01 · <happy path>
- Given <initial state>
- When <action>
- Then <observable result>

### CN-02 · <error or limit>
- Given
- When
- Then

## 3. Acceptance criteria

<!-- Inherit or refine the PRD ACs. Each AC references the scenario that proves it. -->

- AC-01 (CN-01) When <X>, then <Y>.
- AC-02 (CN-02) When <X>, then <Y>.

## 4. Out of scope

- <what belongs to another story, with the ID if it already exists>

## 5. Technical impact

<!-- Fill in only what applies. Used to break down tasks and choose roles. -->

| Area | Impact |
|---|---|
| Data / schema | <entity, migration, backfill> |
| API / contract | <new or changed operation; compatibility> |
| Interface | <screens, loading/error/empty states, accessibility> |
| Integrations | <external service, event, queue> |
| Dependencies | <lib to install / update / remove> |
| Security / sensitive data | <authorization, PII> |

## 6. Verification strategy

<!-- How the story will be proven ready, beyond task tests. -->

- Tests: <unit / integration / e2e; what each level proves>
- Manual or browser verification: <flow, keyboard, two widths> or "not applicable"
- Test data: <fixture or synthetic data>

## 7. Risks and pending decisions

- Risk: <description> · mitigation: <action> · related incident: <RISKS.md#id or none>
- Pending decision: <question> · answered by: <who> · blocks: <task>

## 8. Completion evidence

<!-- Filled in when closing. Only what was actually executed. -->

| AC | Evidence (relative path, command, result) | State |
|---|---|---|
| AC-01 | | passed / failed / not-run |
