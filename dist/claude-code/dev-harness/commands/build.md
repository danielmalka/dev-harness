---
name: build
description: Implement one authorized slice with QA and review
author: malka
argument-hint: "[task id or slice]"
metadata:
  roles: [coordinator, builder, qa-verifier, reviewer]
  skills: [incremental-implementation, data-migrations, regression-testing, ui-verification, code-review]
  writes: authorized write set only
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Slice: $ARGUMENTS. Every implementation slice dispatches `builder` with the Agent tool, model `sonnet` in every case, mode `build`, and the lane named in the dispatch. A present `Lane` value selects that lane alone: `backend`, `frontend`, `dados` (also skill `data-migrations`) or `infra`. A slice with no `Lane` field falls back, unchanged, to the original kind-of-work clauses: lane `backend` for server, service or domain work, lane `frontend` for user interface work, lane `dados` (also skill `data-migrations`) for schema work; lane `infra` is not reachable through this fallback. For lane `backend`, `frontend`, `dados`, `infra`, or the lane-less fallback, in sequence: dispatch `builder` with the selected lane, requiring the kit skill `incremental-implementation`. When the slice's `Doubt` line is `yes`, the Coordinator runs the in-flight doubt in `.agents/coordinator.md` before the builder continues. The builder does not dispatch that review. Freeze the resulting state, then dispatch `qa-verifier` (model `sonnet`, skill `regression-testing`, plus `ui-verification` when the slice has a user interface). Then dispatch `reviewer` in mode `code` (model `sonnet`, skill `code-review`) on that same stable state. For lane `teste (unitário)`, `teste (integração)` or `teste (unitário e integração)`: dispatch only `qa-verifier` (model `sonnet`, skill `regression-testing`, plus `ui-verification` when the slice has a user interface); no `reviewer` is dispatched for this lane. When the "Worktree" section below says this task runs in a worktree (`worktree: always`, or the owner asked for one), create one for the task, never per slice: branch is the lowercase `<task-id>` (`^[A-Za-z0-9][A-Za-z0-9_-]*$`, checked before any git command; the removal first requires the worktree's actual branch, from `git worktree list --porcelain`, to equal it), base is the current `HEAD`; the builder, `qa-verifier` and `reviewer` dispatches carry its absolute path, and its relative path and branch go in `TASK.md`. Cleanup follows that section, only after merge and tag. At most two specialists run concurrently and only on independent work; never run QA or review while the evaluated code, contracts, configuration or tests are still changing.

## Prerequisites
An authorized slice with its completion condition, its write set and the acceptance it must satisfy. Missing any of these: report the gap and that `/dh:plan` produces the slice, then stop. Check commands come from `.harness/project.yaml`; when undiscoverable, every check is reported not-run with the reason. Ask the owner only for a decision that changes scope or contract.

## Output
The `incremental-implementation` report (conclusion, evidence, files changed, checks table with its `Red:` line, limitations, next step), the QA acceptance matrix, and the independent review with severity-graded findings and their resolution. Evidence and screenshots go under `.harness/tasks/<id>/evidence/`; update `.harness/tasks/<id>/TASK.md` with status and results. Record status, effective models, findings and residual risk in `.harness/MEMORY.md`, and any severe incident in `.harness/RISKS.md`.

## Limits
- Build only inside the authorized write set; naming a file does not authorize editing it.
- No dependency install, upgrade, removal, commit, push, publish or deploy without authorization recorded for this slice.
- QA writes tests after the builder finishes and never edits production code; any later edit invalidates the affected evidence and review, so rerun both on the new state. For a lane `teste (unitário)`, `teste (integração)` or `teste (unitário e integração)` ticket, a correction reruns only the `qa-verifier` matrix; `reviewer` is never dispatched for that lane.
- After six unsuccessful correction rounds on the same issue, report blocked or partial and replan with the owner; changing the model does not reset that limit.
- A required failed or not-run check means partial or blocked, never verified.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Worktree
`/dh:build` and `/dh:auto` use this section; `/dh:fix` and `/dh:refactor` do not. You run every git command here in the main session, never a specialist and never the `EnterWorktree` tool.
- **Effective value.** Take the first valid `worktree` (`always` or `on-request`): the top-level line `worktree: <value>` in `local.yaml` of the harness dir (plain text, no parser), then `defaults.worktree` from `dh harness-path --json`, then `on-request`. Tell the owner about an invalid value and skip to the next source.
- **When.** `always` runs the task (`/dh:build`) or run (`/dh:auto`) in a worktree; `on-request` only when the owner asks. One worktree per task or run, never per slice.
- **Create.** End the run before any git command unless `<task-id>` (task id or `AUTO-<n>`) and `/dh:auto`'s `<base-branch>` match `^[A-Za-z0-9][A-Za-z0-9_-]*$`. Path: `<parent of the main checkout>/<checkout folder name>.worktrees/<task-id>`, never under `<home>`. Run `git worktree add -b "<branch>" "<path>" "<base>"`, quoting every path and branch argument in every command here (`git -C "<path>" ...`, `git worktree remove "<path>"`, `git branch -D -- "<branch>"`). `/dh:build`: branch is `<task-id>` in lowercase, base is the main checkout's current `HEAD`; `/dh:auto`: branch is the run's single `<base-branch>`, base is `main`. Then confirm `dh harness-path --json "<path>"` returns the same `dir` as the main checkout. On any failure (invalid id, path present that is not this task's registered worktree, branch in use elsewhere, `git worktree add` error, mismatched `dir`) report the cause and stop: no automatic suffix, no fallback to the main checkout without the owner's answer; a worktree that failed the `dir` check stays and is reported.
- **Record.** Write the path relative to the main checkout (`../<folder>.worktrees/<task-id>`) and the branch in the task record (`TASK.md` for `/dh:build`, `BRIEF.md` for `/dh:auto`) and in `MEMORY.md`. Records stay in the main checkout's harness dir; no agent writes records in a `.harness/` copy inside the worktree.
- **Dispatch.** Each `builder` dispatch carries the absolute path, computed at dispatch time and never written to a versioned file; `qa-verifier` and `reviewer` work on the same path. Each specialist dispatch says the specialist runs no `git worktree add`, `remove` or `prune` (prose only; no runtime guard blocks them today).
- **Cleanup.** Only you, only after the PR is merged and the tag exists, never with `--force`. Check the recorded path matches the computed one for the validated `<task-id>`, and that the actual branch at that path in `git worktree list --porcelain` equals the recorded and the expected branch (`/dh:build`: lowercase task id; `/dh:auto`: the recorded `<base-branch>`). `git worktree remove` still deletes IGNORED files (local builds, `.env`, the ignored `.harness/` copy) that `status --porcelain` omits, so list them with `git -C "<path>" status --porcelain --ignored` in the run report first. Remove only when `git -C "<path>" status --porcelain` prints nothing AND (the worktree `HEAD` equals `headRefOid` from `gh pr view <n> --json state,headRefOid` OR `git rev-list HEAD --not --remotes` run there prints nothing): `git worktree remove "<path>"`, then `git branch -D -- "<branch>"`. On any mismatch or failed condition remove nothing and report path and reason. A run that stops before merge and tag leaves the worktree in place; no automatic sweep of orphans.

## Next
`/dh:build` on the next ready slice, or `/dh:release` when the authorized scope is complete.
