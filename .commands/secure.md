---
name: secure
description: Bounded security review of a change, read-only
author: malka
argument-hint: "[scope, diff range or file list]"
metadata:
  roles: [coordinator, security-reviewer]
  skills: [security-review]
  writes: none (read-only)
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Scope: $ARGUMENTS. Dispatch `security-reviewer` with the Agent tool, model `opus`, and require it to load the kit skill `security-review`. Pass the stable diff or explicit file list, the boundaries the project considers sensitive, the contracts the code exposes, any existing scanner output, and the statement that no execution is authorized unless the owner granted it in this session.

## Prerequisites
The change under review as a stable diff or an explicit file list; without it, ask the owner and stop rather than sweeping the repository. Sensitive boundaries and data classes come from the project; when unstated, they are inferred and every inference is listed as a hypothesis, not a fact. Consult `.harness/RISKS.md` for prior incidents in the touched area. When the change has no authentication, data, input or execution surface, say so and stop.

## Output
The `security-review` report: scope and boundaries touched, findings each classified as demonstrated vulnerability, confirmed exposure or hypothesis requiring investigation, with source, path, sink, mitigation present and why it is insufficient, attacker and gain, scenario, smallest mitigation and its verification; plus open hypotheses, non-issues with the evidence that dismissed them, and checks. Record confirmed exposures as incidents in `.harness/RISKS.md` and the verdict in `.harness/MEMORY.md`.

## Limits
- Read-only: no fix is applied and no file is edited by this command.
- Never execute a payload, exploit or scanner without authorization recorded in this session, and never touch a live, shared or third-party system.
- Never print a secret, token or personal datum into the report; cite its location and kind instead.
- Missing evidence alone does not establish a non-issue; unproven concerns stay qualified as hypotheses, and a missing exploit chain never erases a confirmed exposure.
- Do not downgrade this review from the Opus default to save cost.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dev-harness:fix` for a demonstrated vulnerability or confirmed exposure, then re-run this command on the new state.
