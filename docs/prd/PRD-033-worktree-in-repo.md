# PRD-033 · Kit worktrees in `<checkout>/.worktrees/`

| Field | Value |
|---|---|
| Status | approved on 2026-11-16 |
| Owner | Daniel Malka |
| Created / updated | 2026-11-16 / 2026-11-16 |
| Tickets | to be defined by the planner |
| Origin | Owner decision 2026-11-16, finding A2 of the repository dossier, card 212 |
| Depends on | PRD-028 and ADR-008 (decision D1 is amended here); version base 0.39.2 |
| Delivery | One pull request, version 0.40.0 |

## 1. Problem

Facts:

- The kit creates a development worktree in a folder next to the checkout, `<parent>/<checkout folder>.worktrees/<task-id>` (`.commands/build.md`, section "Worktree", Create and Record; `docs/adr/ADR-008-worktree-de-desenvolvimento.md`, D1; `README.md`, `README.en.md`, both tutorials).
- The owner fixed another rule for the agents that work on the repositories: `<repo>/.worktrees/<task>`. PR #58 already puts `.worktrees/` in the `.gitignore` of this repository. The kit contradicts its owner's own rule.
- A sibling folder is outside the session directory. The tutorials note that nobody measured whether Claude Code asks for write permission there.
- A worktree inside the checkout breaks `dh validate`. Reproduced on `main` at 0.39.2: `go run ./cmd/dh validate .` in a checkout that has a worktree in `.worktrees/` fails with "private path" and "unresolved link" errors on the copy of the repository. The walker `iterTextFiles` (`internal/kit/validate.go`) does not skip `.worktrees`.
- The resolver (`internal/harness/worktree.go`) reads the gitfile and does not depend on where the worktree folder is. The real-git test `TestR18_R19_R26_R61_RealGitWorktree` (`cmd/dh/worktree_cmd_test.go`) only builds the old sibling path.

## 2. Solution

After delivery, the kit creates the worktree of a task or of an `/dh:auto` run in `<checkout>/.worktrees/<task-id>`. The folder is git-ignored. `/dh:setup` adds the line `.worktrees/` to the project `.gitignore`, and `/dh:build` stops before it creates a worktree if the folder is not ignored. `dh validate` skips the folder. Worktrees that exist in the old sibling folder keep working: nothing is moved. The ADR-008 decision D1 is amended, and the docs say one thing.

## 3. Rules

- R1 `.commands/build.md` (Create, Record, Cleanup), `.commands/resume.md` and `.commands/auto.md` use `.worktrees/<task-id>` under the main checkout as the path, and the same relative form in the task record (`TASK.md`, `BRIEF.md`) and in `MEMORY.md`.
  - When `/dh:build` creates a worktree for task `T-9`, then the path is `<main checkout>/.worktrees/T-9`, never under `<home>` and never a sibling folder.
  - When a record is written, then its relative path is `.worktrees/T-9`.
  - When `/dh:auto` creates the worktree of a run, then its folder is `.worktrees/AUTO-<n>`.
- R2 `/dh:build` checks the ignore before it creates a worktree, and never edits `.gitignore` itself.
  - When `git -C "<main checkout>" check-ignore -q -- ".worktrees/<task-id>"` fails, then `/dh:build` stops before `git worktree add`, tells the owner to run `/dh:setup` or to add the line `.worktrees/` to the project `.gitignore`, and changes no file.
  - When the check passes, then the creation goes on as before (id check, `git worktree add`, `dh harness-path` check).
- R3 `/dh:setup` adds the line `.worktrees/` to the root `.gitignore` of a git repository when `git check-ignore -q .worktrees/x` fails. A specialist (`builder`, lane `infra`) writes it, the check runs again, and the report lists it under Created or Left untouched. This is the only `.gitignore` edit that setup makes, and the Limits line of `setup.md` says so. A folder that is not a git repository needs nothing.
  - When the line is missing in a git repository, then setup appends exactly `.worktrees/` and no other line.
  - When the folder is already ignored, then setup changes nothing and reports Left untouched.
  - When the project is not a git repository, then setup writes no `.gitignore`.
- R4 Records made before 0.40.0 still work. `/dh:resume` and the cleanup checks of `build.md` accept the old relative path `../<checkout folder name>.worktrees/<task-id>` for that existing worktree. Every new worktree uses the new path. Nothing is moved, converted or removed.
  - When a record holds the old path and `git worktree list --porcelain` lists it with the recorded branch, then `/dh:resume` reuses it.
  - When a record holds any other path than the new one or the old one, then `/dh:resume` reports it and stops, as before.
- R5 `dh validate` skips the folder `.worktrees` when it reads the repository (`iterTextFiles`).
  - When a worktree copy of the repository sits in `<root>/.worktrees/T-1` and holds a private path and a broken link, then `Validate` reports neither.
  - When the same test runs on `main` at 0.39.2, then it fails (reproduction before the fix).
- R6 The resolver does not change. The real-git test builds the worktree at `<main>/.worktrees/T-X`, in the modes `repo` and `global`, and keeps every assertion of today.
  - When `go test ./cmd/dh/ -run RealGitWorktree -v` runs with `git` in the PATH, then both subtests pass and neither is skipped.
- R7 The docs say one thing. `AGENTS.md` has one new line with the rule. `docs/adr/ADR-008-worktree-de-desenvolvimento.md` D1 and its consequences are amended with the date and this PRD; the comparison table stays as the record. `README.md`, `README.en.md`, `docs/tutorial.html` and `docs/en/tutorial.html` name the new folder and the ignore, and say that setup also writes the `.worktrees/` line in mode `global`. The permission note is replaced by the compatibility note. `.skills/external-clis/references/clone-isolation.md` says that `.worktrees/` must be ignored and states the residual: a call whose project is the main checkout lists the worktrees in the `--ignored` status, so a write by the builder during the call discards the verdict (a false alarm, not a miss).
  - When the live docs and commands are searched for the old sibling folder, then the only matches are the compatibility sentences (R4) and the amendment text.
  - When PRD-028 is read, then a note at its top says that R63 is replaced by this PRD. The history in `CHANGELOG.md` and in the old roadmap cards keeps its words.
- R8 Version 0.40.0 is set in `internal/build/build.go`, in `.claude-plugin/marketplace.json` (version and ref `v0.40.0`) and in `harness-manifest.json`. `CHANGELOG.md` has an entry dated 2026-11-16. `docs/roadmap.html` and `docs/en/roadmap.html` have updated chips, the Sources range, and a new delivered version card. `dist/` is rebuilt, binaries included. The counts stay at 13 agents, 19 commands plus 1 mod command, and 22 skills.
  - When the files named in this rule are read, then each version field says 0.40.0.
  - When `go run ./cmd/dh build` runs, then `dist/` shows no difference except the `generated_at` line of the manifests.
  - When `go test ./...`, `go run ./cmd/dh validate .` and the clean-clone proof of `AGENTS.md` run, then all pass. This PRD touches `internal/` and `cmd/`, so the proof is required.

Out of scope: moving or deleting worktrees that exist; an automatic `.gitignore` edit by `/dh:build`; `.git/info/exclude`; `worktree` in the `gh` guard list; the `EnterWorktree` tool; the sandbox clone of the CLI reviewers; any new command, agent or skill.

## 4. Docs

- `.commands/build.md`, `.commands/resume.md`, `.commands/auto.md`, `.commands/setup.md`: the flows touched (R1 to R4).
- `docs/adr/ADR-008-worktree-de-desenvolvimento.md`, `docs/prd/PRD-028-worktree-and-approval-page.md` (note), `AGENTS.md` (R7).
- `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `.skills/external-clis/references/clone-isolation.md` (R7).
- `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html` (R8).

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| Keep the sibling folder | The owner fixed `<repo>/.worktrees/<task>` for the agents, and the sibling folder needs a write permission outside the session directory. |
| `.git/info/exclude` instead of `.gitignore` | The owner asked for `.gitignore`, and a versioned line also protects the other clones and agents of the project. |
| `/dh:build` edits `.gitignore` itself | `.gitignore` belongs to the project. Setup is the one place that may write the line, and the Coordinator writes only its own records. |
| Move the existing sibling worktrees | A worktree holds the work of the owner. The resolver and git do not care where the folder is, so nothing needs to move (R4). |
| A new ADR | The decision amends D1 of ADR-008 and keeps every other part. The amendment is dated in the ADR. |

### Risks and pending decisions

- Risk: a project that never ran `/dh:setup` after 0.40.0 has no ignore line. Mitigation: R2 stops `/dh:build` before it creates the worktree and names both fixes.
- Risk: a review call whose project is the main checkout sees the writes of the builder in `.worktrees/`. Mitigation: the false alarm is safe (it discards the verdict). The residual is written in `clone-isolation.md`, and the review points at the worktree path.
- Risk: `setup` and the tutorial said that mode `global` writes only `REVIEW.md` in the tree. Mitigation: R7 changes both sentences.

### Decisions (closed 2026-11-16, owner)

- D1 The worktree folder is `<checkout>/.worktrees/<task-id>` (dossier A2, card 212).
- D2 `/dh:setup` guarantees the `.worktrees/` line in the `.gitignore` of the consumer project (card 212).
- D3 Minor version bump from 0.39.2 to 0.40.0: a command changes behavior and `/dh:setup` writes a new file line.
