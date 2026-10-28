# PRD-028 · Worktree setting and HTML approval page

| Field | Value |
|---|---|
| Status | delivered on 2026-10-28 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-27 / 2026-10-27 |
| Tickets | to be defined by the planner |
| Origin | Owner request 2026-10-27 (discovery closed by the owner) |
| Depends on | PRD-014 (global harness, delivered), PRD-008 and ADR-002 (clone isolation of CLI reviewers, delivered) |
| Delivery | One pull request, version 0.33.0 |

## 1. Problem

Two needs ride in one batch. Both belong to the development flow.

**W · Development always runs in the main checkout.** Fact: no setting lets the owner say "develop this feature in a git worktree". Fact: `.commands/build.md` and `.commands/auto.md` never create a worktree. Fact: `auto.md` creates one `<base-branch>` from `main` per run, in the main checkout. Impact: the owner cannot keep the main checkout clean while a run executes.

Three parts of the kit block or contradict a worktree today.

- Fact: `dh harness-path` resolves by `<repo>/.harness` and then by the canonical repo path in the registry (`internal/harness/resolve.go`). A worktree is another path. In global mode it resolves to `none`.
- Fact: `.skills/external-clis/references/clone-isolation.md` and `references/live-tree-freeze.md` refuse a project whose `.git` is not a real directory. A linked worktree has a gitfile, so every CLI review call ends as "isolation unavailable".
- Fact: `docs/roadmap.html` and `docs/en/roadmap.html` still promise a future batch "isolation by disposable git worktree" for CLI reviewers. This contradicts the clone decision in `AGENTS.md`, ADR-002 and PRD-008.

**A · The owner approves a plan without a structured page.** Fact: the owner approves a plan from chat text. A plan with decision points has no view that shows all steps, the proposed snippets and the points to approve. The owner also has no fixed place to leave a reason for a rejection.

## 2. Solution

After delivery, the owner chooses where development runs. A setting `worktree` has two values: `always` and `on-request`. The default is `on-request`. The owner sets it globally with `dh config set worktree <value>`. The owner overrides it per project in `local.yaml`.

When the effective value requires a worktree, the Coordinator creates one worktree for the run. The Coordinator dispatches the `builder` with the worktree path. A CLI reviewer called from that worktree behaves by kind of call.

- A clone-based call accepts a valid linked worktree and keeps every RISK-002 guard (R26 to R34).
- A cloneless read-only call, today `codex`, is unchanged. It does not depend on the worktree layout.
- The live-tree review of `agy` ends as `isolation unavailable` from a linked worktree (R64).

At the planning phase, when the plan has decision points, `docs-guide` writes one HTML approval page in the pdocs folder. The page lists the steps of the work, each with a description. It shows proposed snippets. It shows some points as cards. The owner approves or rejects each card and may add a comment. The page has no server. A "Copy decisions" control produces a Markdown block. The owner pastes this block in the chat.

## 3. Rules

### 3.1 Setting `worktree`

- R1 When `dh config set worktree always` or `dh config set worktree on-request` runs, then the command stores the value in `<home>/config.yaml`.
- R2 When `dh config set worktree <value>` runs with any other value, then the command exits with an error that lists the two valid values. It leaves `config.yaml` unchanged.
- R3 When `dh config unset worktree` runs, then the key leaves `<home>/config.yaml`. The effective global value becomes `on-request`.
- R4 When `dh config` runs, then the output lists `worktree` with its stored value. When `dh config --json` runs, then the JSON has a `worktree` field. In both cases, an absent key shows an empty value, as `language`, `mode` and `dashboard.sprites` do. The default `on-request` comes from R5 and R7.
- R5 When `dh harness-path --json` runs, then the `defaults` object has the field `worktree`. The field holds the global value, or `on-request` when the key is absent. This field does not include the `local.yaml` value.
- R6 The key `worktree` in `<harness dir>/local.yaml` overrides the global value for that project on that machine. A valid value is `always` or `on-request`.
- R7 The effective value is the first one found in this order: `local.yaml`, then global config, then `on-request`.
- R8 When `local.yaml` has an invalid `worktree` value, then the Coordinator tells the owner. The Coordinator skips that source and uses the next one in the order of R7. The same holds for an invalid value that the owner wrote by hand in the global `config.yaml`. The Go code does not check the stored `worktree` value against `always` and `on-request`. A value that breaks the YAML scalar rule aborts the parse, as it does for `language`. That rule rejects an unquoted value that has `:` or starts with `{`, `[`, `|`, `>`, `&`, `*` or `!`.
- R9 The Coordinator reads `local.yaml` as plain text. The change adds no Go parser for `local.yaml`. The file `.skills/external-clis/references/local.yaml.example` documents the key `worktree`.

### 3.2 Behavior in `/dh:build` and `/dh:auto`

- R10 When the effective value is `always`, then `/dh:build` and `/dh:auto` run the development of the task or the run inside a git worktree.
- R11 When the effective value is `on-request`, then these commands create a worktree only when the owner asks for one in the request.
- R12 The Coordinator creates the worktree in the main session. Every dispatch of a builder or other specialist in a worktree run states in its text that the specialist must not run `git worktree add`, `git worktree remove` or `git worktree prune`. No runtime block for these commands exists today.
- R13 The Coordinator creates one worktree for each run (`/dh:auto`) or each task (`/dh:build`). It creates no worktree for each slice.
- R14 When a worktree exists, then the Coordinator dispatches each `builder` with the absolute path of the worktree. The `qa-verifier` and the `reviewer` also work on that path.
- R15 In `/dh:auto`, the single `<base-branch>` of the run is the branch of the worktree. The command creates no second branch.
- R16 The Coordinator records the worktree path and branch in the task record, so `/dh:resume` finds them. The task record is `TASK.md` for `/dh:build` and `BRIEF.md` for `/dh:auto`.
- R17 When the Coordinator cannot create the worktree, then it reports the cause and stops. It does not fall back to the main checkout without the owner's answer.

### 3.3 Harness resolution and records

- R18 When `dh harness-path --json` runs inside a valid linked worktree (R26) of a linked project, then the fields `dir` and `pdocs` equal the values that it returns in the main checkout. The key for both fields is the path of the main checkout. This holds in repo mode and in global mode.
- R19 When a linked worktree contains a `.harness/` directory (repo mode, versioned), then resolution ignores that copy. It returns the harness dir of the main checkout.
- R20 In repo mode, the Coordinator keeps writing records in the harness dir of the main checkout. The `builder` changes only the files of its write set in the worktree. No agent writes records in the `.harness/` copy inside the worktree.
- R21 When `dh harness-path` runs in a path that is not a linked worktree and not a linked project, then the result is unchanged from today.

### 3.4 Worktree cleanup

- R22 The Coordinator removes the worktree only after the pull request is merged and the tag exists. In `/dh:auto`, the removal is the last `Stage:` step, `cleanup`, which runs after `tag`.
- R23 When the worktree has uncommitted changes, then the Coordinator does not remove it. It reports the path and the reason to the owner.
- R24 When the worktree branch has commits that the remote does not have, then the Coordinator does not remove it. It reports the path and the reason to the owner. A squash merge that deletes the remote branch leaves the commits on no remote ref. The commits count as pushed when the worktree `HEAD` equals the `headRefOid` of the merged pull request. Removal never uses `--force`.
- R25 When the run stops before merge and tag, then the Coordinator leaves the worktree in place and records its path in the task record. When an `/dh:auto` run stops during `cleanup`, then `Stage: cleanup` stays in `BRIEF.md` for `/dh:resume`.

### 3.5 CLI reviewer isolation from a worktree

- R26 When `<project>/.git` is a directory, then the existing isolation checks apply and no other check applies. When it is not a directory, then the call needs a non-null `linked_worktree` (R61). The `.git` must also point to a git dir of a valid worktree layout: `<common dir>/worktrees/<name>`, with the matching `commondir` and back-pointer. Right before the clone, `.git` must be a regular file (`test -f`) and not a symlink (`test ! -L`).
- R27 When `.git` is a symlink, then the call ends as "isolation unavailable". The same holds when `linked_worktree` is `null` and `.git` is not a directory. The same holds when `.git` points outside a valid worktree layout. The same holds when `.git` is not a regular file at the check before the clone.
- R28 The CLI sandbox is always a new clone made with `git clone --no-hardlinks` into a disposable temporary folder. The kit never uses the worktree as the sandbox. The kit never uses `git worktree` for the sandbox.
- R29 The clone has its own `.git` directory. It has no `origin` remote. It has no hardlink to the source. These checks are the same as for a normal checkout.
- R30 The live-tree freeze covers two parts. An item that is absent in both snapshots is not a difference and is never an error. In a normal checkout, the common dir is the `.git` directory.
  - Part (a), every project: the existing R5 coverage, and the resolved harness dir with its byte coverage. In a linked worktree, the resolved harness dir (R18) is not `<project>/.harness/`, which is only the copy inside the worktree. In global mode, the harness dir stays outside the freeze. Also in the common dir: the files under `refs/`, and the entries under `worktrees/` when they exist, with the files of each entry that part (b) lists.
  - Part (b), only when `<project>/.git` is a gitfile (linked worktree): the bytes of the gitfile, and, in the worktree admin dir, `HEAD`, `commondir`, `gitdir`, `config.worktree`, the files under `refs/`, and the presence of `locked`.
  - The raw `index` files stay out of the freeze.
- R31 When the freeze detects a change in the tree or in the git dirs of R30 during a CLI call, then the call discards the verdict. A normal checkout behaves the same way.
- R32 The Coordinator writes nothing to the harness dir or to the tree while a CLI call is in flight. This is the existing rule for RISK-002. It stays.
- R33 The anti-delegation clause stays verbatim in every external CLI prompt (RISK-001).
- R34 RISK-002 stays "mitigated". No file of the change calls it "resolved".

### 3.6 Stale text

- R35 The files `docs/roadmap.html` and `docs/en/roadmap.html` do not promise isolation by `git worktree` for CLI reviewers. The text states the clone decision.
- R36 `AGENTS.md` names the setting `worktree` with its three sources and its default. It states one rule: a worktree is allowed for development, and the CLI sandbox stays a clone. R30 shrinks the accepted residuals "paths of the live `.git` outside the R5 snapshot". `AGENTS.md` lists the residuals that remain. ADR-008 records which residuals shrink. The change does not edit ADR-002 or ADR-003. It keeps the rest of `AGENTS.md` as it is.

### 3.7 Approval page

- R37 `docs-guide` generates the page with the skill `doc-template-html`. The skill gets one new `--type` for the approval page. The change updates `.skills/doc-template-html/scripts/stamp.sh` and `references/types.md`. It adds one skeleton in `assets/skeletons/pt-br/` and one in `assets/skeletons/en/`.
- R38 `/dh:plan` dispatches `docs-guide` for the page when the plan has at least one decision point. When the plan has none, `/dh:plan` generates no page.
- R39 The page lists every step of the work in order. Each step has a title and a description.
- R40 The page shows the proposed snippets for the feature. Each snippet is in a code block and names the step it belongs to.
- R41 The page shows each decision point as a card. Each card has an id, a title, a short description, an Approve control and a Reject control.
- R42 Each card has a comment field. The comment is optional when the owner approves. The comment is mandatory when the owner rejects.
- R43 When the owner rejects a card with an empty comment (blank characters only), then the page shows a visible message on that card. The page does not include that card as decided in the copied block.
- R44 The "Copy decisions" control produces a Markdown block. Each card gives one entry: card id, decision (`approved` or `rejected`) and comment. A card without a decision shows as `pending`.
- R45 When the clipboard is not available or the browser refuses the copy, then the page shows the Markdown block in a text area. The owner selects and copies it by hand.
- R46 The page has no server write. The change adds no dashboard write endpoint. No process reads the decisions automatically. The owner pastes the block in the chat.
- R47 The page keeps the in-progress choices in browser storage. Every storage call is inside try/catch. The page works with storage blocked or empty.
- R48 The page escapes all inserted text: plan text, snippets and the comment field. Text such as `<script>` or `<img onerror=...>` shows as plain text and runs no code. The page uses no `innerHTML` with unescaped text.
- R49 The page loads no external script. The only external origins are the pinned Mermaid CDN and Google Fonts, which already hold for pdocs.
- R50 The page works by keyboard only. Every control has a visible focus style and a text label. Approve and Reject are not shown by color alone.
- R51 The page has no horizontal scroll at a phone width of 360 px.
- R52 The page is a flat `.html` file in the pdocs folder. The folder path comes from `dh harness-path --json`, field `pdocs`. When this field is empty (a project with no registry entry), the Coordinator reports the cause and generates no page. Inside a linked worktree, this field resolves from the main checkout (R18). The file is within the 8 MiB cap of the dashboard.
- R53 When `<pdocs>/index.html` exists, then the generation updates it with a link to the page.
- R54 The page language follows `language` of the project (`pt-br` or `en`).

### 3.8 Delivery

- R55 The version is 0.33.0. `CHANGELOG.md` has the entry. `dist/` is rebuilt with `go run ./cmd/dh build`.
- R56 A new ADR (ADR-008) records the worktree decision. It states the three sources of the setting, who creates the worktree, and why the CLI sandbox stays a clone.
- R57 The change updates the documents of section 4 inside the same pull request.
- R58 The count stays 19 command files plus 1 mod command. The change adds no `/dh:` command.
- R59 Because the change edits executable harness instructions that gate a security control (R26 to R34), the batch needs the `security` review stage. The stage dispatches the full configured `reviewers.security` list. The batch also needs the default code review and the `qa-verifier` matrix.
- R60 Every new file uses relative paths or `~/...`. `dh validate` finds no absolute private path.

### 3.9 Added by the amendment after ADR-008

- R61 The JSON of `dh harness-path --json` gets one additive field `linked_worktree`. Outside a valid linked worktree, it is `null`. Inside one, it is an object with three fields (`main`, `gitdir` and `commondir`), each a canonical path.
  - `main`: the main checkout.
  - `gitdir`: the worktree admin dir `<common dir>/worktrees/<name>`. It is not the gitfile.
  - `commondir`: the common dir.
  - The existing fields keep their meaning. The isolation step (R26, R30) reads the admin dir from `gitdir`.
- R62 For `/dh:build`, the branch of the worktree is the `<task-id>` in lowercase. For `/dh:auto`, R15 applies. The Coordinator creates the worktree in the main session with `git worktree add -b`. It does not use the `EnterWorktree` tool of Claude Code.
- R63 The worktree folder is `<parent of the main checkout>/<checkout folder name>.worktrees/<task-id>`. The `<task-id>` is the task id, or `AUTO-<n>` for `/dh:auto`. It has only the characters `[A-Za-z0-9_-]`. Any other character ends the run as in R17. The folder is never under `<home>`. The task record (R16) holds the path relative to the main checkout. The dispatch to the builder (R14) holds the absolute path.
- R64 When the project is a linked worktree and the call is the `agy` live-tree review, then the call ends as `isolation unavailable`. The review is not run, and no review round is spent. The reason is that its freeze snapshots `.git/hooks/` and `.git/config`, which are not folders in a gitfile checkout.

## 4. Docs

- README.md
- README.en.md
- docs/tutorial.html
- docs/en/tutorial.html
- docs/roadmap.html
- docs/en/roadmap.html
- docs/adr/ADR-008-worktree-de-desenvolvimento.md
- AGENTS.md
- CHANGELOG.md
- .commands/build.md
- .commands/auto.md
- .commands/plan.md
- .skills/external-clis/references/clone-isolation.md
- .skills/external-clis/references/live-tree-freeze.md
- .skills/external-clis/references/local.yaml.example
- .skills/doc-template-html/references/types.md

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| A worktree for each slice | More worktrees to create and clean. The owner decided one for each run or task. |
| A subagent creates the worktree | The main session holds the worktree capability and the task record. A prompt rule keeps subagents away from worktree commands (R12). |
| Use the worktree as the CLI sandbox | A worktree shares the whole `.git` of the main checkout. A CLI could write there. The clone decision stays. |
| Keep refusing a worktree in the CLI isolation (option b) | The owner chose option (a): accept a worktree with the guards of section 3.5. |
| A dashboard endpoint that receives the decisions | Out of scope by owner assumption. It adds a write surface. |
| A new Go parser for `local.yaml` | The Coordinator reads two plain-text sources. No need found. |

### Assumptions the owner may overturn

- A1 `/dh:plan` is the only command that generates the page. `/dh:plan-loop` and `/dh:document` are out of scope.
- A2 The setting applies to `/dh:build` and `/dh:auto`. `/dh:fix` and `/dh:refactor` are out of scope.
- A3 The return channel is the "Copy decisions" Markdown block that the owner pastes in the chat (R44 to R46).
- A4 The approval page is a new `--type` of `doc-template-html`, not a new skill and not a new command.
- A5 The worktree sits outside the main checkout tree. ADR-008 decides the exact folder, the branch name and the creation method (R62, R63).
- A7 This batch does not change the Bash guard of PRD-027. A guard change would also deny read forms such as `git worktree list`. It needs its own security review. R12 uses a prompt rule only.
- A8 The update of `<pdocs>/index.html` (R53) follows the existing Coordinator rule. `docs-guide` updates that file when it exists. Only `/dh:document pdocs` creates it.
- A6 The type name is chosen by the planner. It follows the existing names in `stamp.sh`.

### Out of scope

- Any dashboard write endpoint, automatic reading of decisions, and multi-user sync.
- Worktree handling for `/dh:fix`, `/dh:refactor`, and any command not named in R10 and R11.
- A change to the CLI sandbox decision (clone stays).
- Eval cases and new probes.

### Risks and pending decisions

- Risk: R26 to R31 widen a security control · mitigation: R59 (full security stage), R27 refusal cases, and QA cases for a symlink gitfile and a gitfile that points outside a valid layout.
- Risk: a hand-edited `local.yaml` holds a bad value · mitigation: R8.
- Decided in ADR-008: R18 and R19 resolve by reading the gitfile, `commondir` and the back-pointer in Go, without `git`. The behavior in R18 to R21 is the requirement.
- Candidate (optional hardening): add `worktree` to the git write list of the guard, with its own security review · answered by: owner · blocks: nothing in this batch.
- Decided in ADR-008: folder of the worktree (A5), branch name and creation method (R62, R63). Nothing pends here.
- Pending decision: name of the new `--type` (A6) · answered by: planner · blocks: the page ticket.
- Decided in ADR-008: the main session creates the worktree with `git worktree add -b`, not with the harness worktree tool (R62).
- Open point: it is not verified whether Claude Code asks for write permission on the sibling folder `<checkout folder name>.worktrees` outside bypass mode · answered by: a `qa-verifier` probe, or a setup note in the docs · blocks: the docs ticket. This PRD asserts no result.

## Executed summary

- **Delivered:** The `worktree` setting (`always` or `on-request`, default `on-request`) now sets the worktree flow of `/dh:build`, `/dh:auto` and `/dh:resume`, a Go resolver reads the harness dir and `pdocs` from the main checkout of a linked worktree, CLI-reviewer isolation accepts a linked worktree with a wider freeze, and `/dh:plan` writes an HTML approval page (type `aprovacao`), shipped as version 0.33.0.
- **Rules:**

  | Rule | Verdict | Evidence |
  |---|---|---|
  | R1 | Honored | T-2801, T-2803; Go test `TestR1_R2_R3_ConfigSetUnsetWorktree` |
  | R2 | Honored | T-2801, T-2803; Go test `TestR1_R2_R3_ConfigSetUnsetWorktree` (refusal case) |
  | R3 | Honored | T-2801, T-2803; Go test `TestR1_R2_R3_ConfigSetUnsetWorktree` |
  | R4 | Honored | T-2801, T-2803; Go test `TestR4_ConfigShowsWorktree` |
  | R5 | Honored | T-2801, T-2803; Go test `TestR5_R7_DefaultsWorktree` |
  | R6 | Honored | T-2806, T-2807; prose check (weak, see note e) |
  | R7 | Honored | T-2803, T-2807; Go test `TestR5_R7_DefaultsWorktree`; prose check (weak) |
  | R8 | Honored | T-2803, T-2807; Go test `TestR8_HarnessPathScalarRules`; prose check (weak) |
  | R9 | Honored | T-2807, T-2814; diff grep, no `local.yaml` read in Go code (weak) |
  | R10 | Honored | T-2806, T-2807; prose check (weak) |
  | R11 | Honored | T-2806, T-2807; prose check (weak) |
  | R12 | Honored | T-2806, T-2807; prose check, no runtime block by design (weak) |
  | R13 | Honored | T-2806, T-2807; prose check (weak) |
  | R14 | Honored | T-2806, T-2807; prose check (weak) |
  | R15 | Honored | T-2806, T-2807; prose check (weak) |
  | R16 | Honored | T-2806, T-2807; prose check (weak) |
  | R17 | Honored | T-2806, T-2807; prose check (weak) |
  | R18 | Honored | T-2802, T-2803; Go test `TestR18_R19_RepoModeResolvesMainHarness`, real `git worktree add` case |
  | R19 | Honored | T-2802, T-2803; Go test `TestR18_R19_RepoModeResolvesMainHarness`; mutation M3 detected |
  | R20 | Honored | T-2806, T-2807; prose check (weak) |
  | R21 | Honored | T-2802, T-2803; Go test `TestR21_NormalCheckoutAndPlainDirAreNotWorktrees` |
  | R22 | Honored | T-2806, T-2807; prose and cleanup-order checks (weak) |
  | R23 | Honored | T-2806, T-2807; prose and cleanup-order checks (weak) |
  | R24 | Honored | T-2806, T-2807; prose and cleanup-order checks (weak) |
  | R25 | Honored | T-2806, T-2807; prose and cleanup-order checks (weak) |
  | R26 | Honored | T-2802, T-2803, T-2805; Go tests and shell matrix (11 PASS) |
  | R27 | Honored | T-2802, T-2803, T-2805; Go tests and shell matrix (26 PASS) |
  | R28 | Honored | T-2805; shell matrix, clone with `--no-hardlinks` (7 PASS) |
  | R29 | Honored | T-2805; shell matrix (3 PASS) |
  | R30 | Honored | T-2805, T-2804; shell matrix (55 PASS) and T-2804 round 4 |
  | R31 | Honored | T-2805; shell matrix (1 PASS, single row) |
  | R32 | Honored | T-2805; shell matrix (5 PASS) and clause checksum |
  | R33 | Honored | T-2805; shell matrix (5 PASS) and clause checksum |
  | R34 | Honored | T-2805, T-2814; shell matrix (11 PASS) |
  | R35 | Honored | T-2811; matrix round 3 and manual re-check |
  | R36 | Honored | T-2811; matrix rows and manual re-check (see note c) |
  | R37 | Honored | T-2808, T-2810; matrix R37 PASS (see note a) |
  | R38 | Honored | T-2809, T-2810; read assertion (weak) |
  | R39 | Honored | T-2808, T-2810; headless browser matrix (156 PASS) |
  | R40 | Honored | T-2808, T-2810; headless browser matrix |
  | R41 | Honored | T-2808, T-2810; headless browser matrix |
  | R42 | Honored | T-2808, T-2810; headless browser matrix (8 PASS) |
  | R43 | Honored | T-2808, T-2810; headless browser matrix (8 PASS) |
  | R44 | Honored | T-2808, T-2810; headless browser matrix (16 PASS), see note b |
  | R45 | Honored | T-2808, T-2810; headless browser matrix (12 PASS) |
  | R46 | Honored | T-2810, T-2814; network log with no write, `git status` clean |
  | R47 | Honored | T-2810; headless browser matrix (opaque origin, storage, clipboard) |
  | R48 | Honored | T-2808, T-2810; headless browser matrix (18 PASS), source scan with 0 sinks; see note c for threat model |
  | R49 | Honored | T-2810; headless browser matrix (no script `src`) |
  | R50 | Honored | T-2810; headless browser matrix (28 PASS, keyboard and focus) |
  | R51 | Honored | T-2808, T-2810; headless browser matrix (12 PASS) and output threshold check |
  | R52 | Honored | T-2809, T-2810; route check (200) and T-2809 matrix |
  | R53 | Honored | T-2809, T-2810; read assertion, never run end to end (weak) |
  | R54 | Honored | T-2808, T-2810; headless browser matrix (8 PASS, pt-BR and en) |
  | R55 | Honored | T-2813, T-2814; prose check, `go test ./internal/kit` pass (weak: no T-2813 evidence folder) |
  | R56 | Honored | T-2814; prose check, ADR-008 status accepted |
  | R57 | Honored | T-2812, T-2814; `git status` check of the 16 files of section 4 |
  | R58 | Honored | T-2814; validate log (19 commands, 13 agents) |
  | R59 | Honored | T-2814; security passes (see note d) |
  | R60 | Honored | T-2814; validate log (0 errors) and prose check of private paths |
  | R61 | Honored | T-2802, T-2803; Go tests `TestR18_R61_GlobalModeLinkedWorktree` and related cases |
  | R62 | Honored | T-2806, T-2807; prose check of quoted command order (weak) |
  | R63 | Honored | T-2806, T-2807; prose check of quoted command order (weak) |
  | R64 | Honored | T-2804, T-2805; shell matrix (17 PASS), not re-run in the final pass |

  Notes on the table:
  - (a) R37 is Honored by a recorded decision of the Coordinator. `docs-guide` has no shell, so the page is built in three dispatches: `docs-guide` writes the data JSON, a `builder` in lane `infra` runs `stamp.sh`, and `docs-guide` adds the index link. The R37 wording says that `docs-guide` generates the page.
  - (b) The copied block also holds the card title, flattened to one line (R44).
  - (c) `AGENTS.md` also records the owner's threat-model decision of 2026-10-28 (personal local tool). This goes beyond the R36 hunk.
  - (d) R59 is Honored. The full `reviewers.security` list (`cli:codex/gpt-6.1-sol` and `claude`) ran in three passes, plus one pass for the approval page. The findings were fixed and re-verified. The `claude` reports are kept verbatim in `~/.harness/projects/dev-harness/tasks/PRD-028/evidence/sec1` to `sec4`. The codex reports are in `sec2` to `sec4`. The raw codex output of the first pass was not kept; only its summary is in the Coordinator's record. No security pass covers the final state, because the owner limited the threat model on 2026-10-28. The default code review and the `qa-verifier` matrix ran for every ticket. The final feature review ran once.
  - (e) Weak evidence, stated plainly: R6 to R17, R20, R22 to R25, R62 and R63 rest on read-only prose checks. T-2807 left raw outputs and no verdict matrix. The final QA refreshed the checks, and they pass.

- **Tickets:** All tickets are done. T-2801 to T-2807 were done on 2026-10-27; the others on 2026-10-28. Evidence paths below are under `~/.harness/projects/dev-harness/tasks/`.
  - T-2801, lane backend (Go config key), done. Evidence: `PRD-028/evidence/final/rtests-verbose.log`.
  - T-2802, lane backend (Go resolver), done. Evidence: `PRD-028/evidence/final/rtests-verbose.log`.
  - T-2803, lane test (Go tests), done. Evidence: `T-2803/evidence/` (for example `02-new-tests-verbose.txt`, mutation runs M2 and M3).
  - T-2804, lane backend (isolation and freeze prose), done. Evidence: `T-2804/evidence/w7/matrix.md`.
  - T-2805, lane test (isolation with a real linked worktree), done. Evidence: `T-2805/evidence/w7/matrix.md`.
  - T-2806, lane backend (build, auto, coordinator and resume prose), done. Evidence: `T-2806/evidence/w5/matrix.md`, `w7/matrix.md`.
  - T-2807, lane test (prose checks), done. Evidence: raw outputs only (`T-2807/evidence/run1.txt` to `run3-text.txt`); no verdict matrix.
  - T-2808, lane backend (approval page type, stamp, skeletons), done. Evidence: the round 3 matrix cited in `PRD-028/evidence/final/coverage.md` (R36 row).
  - T-2809, lane backend (`/dh:plan` hook and `docs-guide`), done. Evidence: `T-2809/evidence/matrix.md`.
  - T-2810, lane test (headless browser): 156 main cases and 49 extra cases, all PASS. Evidence: `T-2810/evidence/matrix.md`, `T-2810/evidence/screenshots/`.
  - T-2811, lane backend (roadmap text, `AGENTS.md`, ADR-008 sentence), done. Evidence: `T-2811/evidence/matrix.md`.
  - T-2812, lane backend (docs), done. Evidence: `git status` check (R57); no evidence folder.
  - T-2813, lane backend (version 0.33.0 and `dist/`), done. Evidence: no evidence folder; the `dist/` check is covered by T-2814.
  - T-2814, lane test (final QA): 21 PASS, 0 FAIL, 1 NOT-RUN, judged by the Coordinator. The clean-clone sha256 proof of the five binaries and the two launchers passed. Evidence: `PRD-028/evidence/final/` (`coverage.md`, `validate.log`, `prose-recheck.txt`, `go-test.log`).

- **Docs:** `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md`, `docs/adr/ADR-008-worktree-de-desenvolvimento.md` (new, with its review record), this PRD and its review record `PRD-028-worktree-and-approval-page.review.md`, `.agents/coordinator.md`, `.agents/docs-guide.md`, `.commands/build.md`, `.commands/auto.md`, `.commands/resume.md`, `.commands/plan.md`, `.skills/external-clis/` (`SKILL.md` and references), and `.skills/doc-template-html/` (`SKILL.md`, `README.md`, `references/types.md`, `stamp.sh`, two skeletons).

- **Out of scope:** The items of the Appendix: any dashboard write endpoint, automatic reading of decisions, multi-user sync; worktree handling for `/dh:fix`, `/dh:refactor` and other commands not named in R10 and R11; a change to the CLI sandbox decision (the clone stays); eval cases and new probes. Also out of scope:
  - No change to the PRD-027 runtime guard for `git worktree`. The ban on specialists is prose only.
  - No dashboard write endpoint. The page returns decisions by a pasted Markdown block.
  - No new agent, skill or command.
  - Candidate batches found during the work, not started: a guard false positive on `git p4` in a search pattern; hardening the guard for `git worktree`; applying the Part A ordering to step 7 of the `agy` live-tree path; a hygiene pass for unquoted `<path>` placeholders in the reversion commands of `live-tree-freeze.md`; rejecting duplicate decision ids in older pages; recipe lessons for `opencode` (`--standalone`, codegraph MCP disabled) and `codex` (`-c mcp_servers.codegraph.enabled=false`).
