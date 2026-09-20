# T-000 · <short title>

<!-- Location: .harness/tasks/T-000/TASK.md (or BUG-000 for a defect).
This is the dispatch an agent receives. Everything it needs to start
without asking is here or linked. Delete the comments when finished.
For a bug, fill in section 2B and ignore 2A. -->

| Field | Value |
|---|---|
| Type | task / bug |
| Status | ready / in progress / blocked / in review / done |
| Story / PRD | ST-000 (AC-01) / PRD-000 |
| Suggested role | backend-builder / frontend-builder / data-engineer / debugger / ... |
| Depends on | T-000, <contract, decision, migration> |
| Blocks | T-000 |
| Authorization | <what may be written; commit/push/deploy require explicit authorization> |
| Created / updated | YYYY-MM-DD / YYYY-MM-DD |

## 1. Objective

<!-- One sentence with the verifiable result. Then why, in one line. -->

## 2A. Scope (task)

**In scope**
- <concrete change>

**Out of scope**
- <what the agent must NOT touch, even if it seems useful>

## 2B. Defect (bug)

- Symptom: <literal error text, sanitized log, or observed behavior>
- Expected: <what should happen>
- Reproduction: <steps or command; if not reproducible, record attempts>
- Environment: <version, platform, config, data>
- Since when / recent change: <commit, release, or unknown>
- Cause: <hypothesis> · confirmed: yes / no · evidence: <...>
- Severity: <actual impact> · incident in RISKS.md: <id or none>

## 3. Technical context

<!-- The minimum needed to avoid reinventing: where the code is, which pattern to follow, which contract to respect. Relative paths. -->

- Candidate files / modules: `<path>`
- Local pattern to follow: <existing example in `<path>`>
- Contracts to respect: <API, event, schema, interface>
- Decisions already made: <ADR or PRD/story decision>
- Critical rules / previous incidents: <RISKS.md#id or none>

## 4. Environment changes

<!-- Everything that changes the environment beyond code. Empty = none. -->

| Type | Item | Action | Reason |
|---|---|---|---|
| lib | <name@version> | install / update / remove | <why> |
| env var | <NAME> | create / change | <purpose; value does not go here> |
| migration | <file> | create / apply | <reversible? yes / no; rollback plan> |
| config / infra | <file or service> | change | <...> |

## 5. Execution plan

<!-- Small, verifiable steps. Each step ends with a check. -->

1. <step> → check: <command or observation>
2. <step> → check:
3. <step> → check:

## 6. Tests

<!-- Each acceptance criterion has at least one test. Real project command. Typecheck does not prove behavior. -->

| Scenario | Type | File / command | Covers |
|---|---|---|---|
| <happy path> | unit / integration / e2e | `<command>` | AC-01 |
| <relevant error> | | | AC-02 |
| <bug regression> | | | 2B |

Manual / browser verification: <flow> or "not applicable"

## 7. Acceptance criteria

- AC-01 When <X>, then <Y>.
- AC-02 When <X>, then <Y>.

## 8. Risks and limits

- Risk: <what may break> · mitigation: <action>
- When encountering an open decision: stop and report; do not invent behavior.
- After 2 unsuccessful correction rounds: return to the Coordinator with evidence.

## 9. Result

<!-- Filled in by the executor on delivery. Only what was done and executed. -->

- Conclusion: <what was delivered, in one sentence>
- Files changed: `<path>`, `<path>`
- Checks: `<command>` → passed / failed / not-run (<reason>)
- Limitations: <what was not verified or remains pending>
- Pending items for the Coordinator: <decision, incident, memory update>
- Next step: <concrete action> · authorization required: <which>
