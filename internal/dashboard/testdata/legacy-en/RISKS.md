# RISKS

Summarized history of severe incidents in this project. Written exclusively by the Coordinator.

## Usage rules

- Consult before planning or implementing a change that affects established behavior, a critical business rule, or a high-risk area, including fixes, refactors, and migrations under these conditions.
- The Coordinator evaluates this trigger from the request, file map, and area criticality. Routine changes outside these conditions do not load this file by default.
- Look for incidents by affected files, modules, and rules. Consider known renames and coupling; do not depend exclusively on the exact path.
- Carry relevant prevention rules into the plan, acceptance criteria, and the dispatch for the responsible specialist.
- Record every known severe incident with enough evidence, even if it is not yet resolved. Severity examples: data loss or corruption, significant unavailability, data exposure, or material breakage of a critical rule.
- Do not invent incidents to fill the file. A hypothetical risk belongs in the task plan; it becomes historical here only when there is a known incident.
- Subagents report occurrences and evidence to the Coordinator; they do not write to this file.
- When a severe incident occurs or is updated, the Coordinator may consult the corresponding record to maintain it. Consolidating MEMORY.md does not remove or archive this prevention history.
- Keep records concise; extensive details remain referenced in task artifacts or EPOCHAL.md. Consult this raw history only if deeper review is needed.
- Do not record secrets. Use relative project paths and dates with timezone offset when known; explicitly declare unknown information.

## Incidents

No incidents recorded. This does not prove the absence of previous incidents.

<!-- Fields for each real incident:
- Stable ID and title.
- Incident date/time with timezone offset, when known.
- Record date/time and last update.
- Severity and concrete impact.
- Location: affected file(s), module(s), symbol(s), or flow(s).
- Critical business rule and behavior that must be preserved.
- Brief description of what happened.
- Confirmed cause or explicitly identified hypothesis.
- State: open, mitigated, or resolved.
- Applied resolution or current mitigation; date when known.
- How to prevent recurrence: relevant constraint, care, and test or verification.
- Evidence and relative references to the task or historical batch.
When resolved, keep the incident and update the state and prevention. -->
