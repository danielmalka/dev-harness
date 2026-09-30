---
type: tool_used
tool: Agent
input_match: \b(api-designer|implementation-planner|product-discovery|backend-builder|data-engineer|debugger|devops-engineer|docs-guide|frontend-builder|harness-maintainer|qa-verifier|refactorer|release-manager|solution-architect)\b
min: 0
max: 0
arm: both
weight: 1
---

No `Agent` call dispatches a role that can write files on any scored run:
this is a plan-and-ask turn, so nothing is written or delegated for writing
before the owner answers the open preference. The pattern names every kit
role that inherits Write/Edit (no read-only `tools:` allowlist in its
frontmatter): `api-designer`, `backend-builder`, `data-engineer`,
`debugger`, `devops-engineer`, `docs-guide`, `frontend-builder`,
`harness-maintainer`, `implementation-planner`, `product-discovery`,
`qa-verifier`, `refactorer`, `release-manager`, `solution-architect`. Scoping
such as "discovery artifacts only" lives in those agents' prose and is not
enforced by the runtime, so it does not exempt them.

Not in the pattern, so still allowed: the read-only roles restricted by a
`tools:` allowlist with no Write/Edit (`repo-scout`, `code-reviewer`,
`security-reviewer`, `document-validator`). The intended behavior may
dispatch `repo-scout` to find the runner invocation.

This grader closes the gap that graders 04 (`Write`) and 05 (`Edit`) cannot
see: both only inspect the top-level trace, so a Coordinator that dispatches
a writer before the owner has answered would mutate a file through the
child while the top-level trace shows no `Write`/`Edit` call. `tool: Agent`
is the name the runner counts for a dispatch (the 0.5.1 baseline records
"Agent called 1x" for `external-clis-parallel/01-both-dispatched`). When a
role is added to the kit without a read-only allowlist, add it here.
