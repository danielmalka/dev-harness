# PRD-032 · Lean Coordinator context and repository bootstrap

| Field | Value |
|---|---|
| Status | delivered on 2026-10-31 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-30 / 2026-10-31 |
| Tickets | to be defined by the planner |
| Origin | Owner request 2026-10-30, source `tasks/PRD-032/SOURCE.md` in the harness dir |
| Depends on | PRD-031 (lean builder flow), merged in `main` at 0.36.0 (merge commit 40538a49). This PRD builds on that state. |
| Delivery | One pull request |

## 1. Problem

In a long task, the Coordinator (the main session) fills its context. The owner loses quality late in the run. The owner wants the Coordinator to do less work itself, and wants a bootstrap command for a repository.

Most facts below come from `wc` and file reads on 2026-10-30, on the PRD-031 branch. PRD-031 is now merged (0.36.0). The merged `.agents/coordinator.md` has 174 lines and 31842 bytes. R15 re-measures the sizes.

- Fact: the harness `MEMORY.md` of this project has 176 lines and 185182 bytes. One line has 73512 characters. Every command reads this file first. It is the largest item the Coordinator loads, about 5.7 times the size of `.agents/coordinator.md`.
- Fact: `.agents/coordinator.md` had 174 lines and 32663 bytes on that branch. It has 31842 bytes at 0.36.0. The sections "Worktree", "Consolidation on explicit request", "External CLI planners", and "Drift report" total about 6100 bytes. Each section serves one or two commands. The Coordinator loads all of them in every session.
- Fact: `/dh:auto` tells the Coordinator to follow each chained command text "inline". The chained files are `discover.md`, `plan.md`, `build.md`, `secure.md`, `review.md`, and `release.md`. They total about 36600 bytes and stay in the main context for the whole run.
- Fact: `.agents/coordinator.md`, Procedure item 6, requires a "reply format" and a "budget" in each dispatch. No rule gives a number. Only `document-validator` has a numeric limit (at most 5 sentences per paraphrase).
- Fact: a specialist report returns whole to the Coordinator chat. A code review report returns to the Coordinator and the Coordinator writes the result into `MEMORY.md`.
- Fact: the Elephant-Goldfish Model (EGM, `google/elephant-goldfish-model`) keeps long context in one session (the elephant) and sends fresh read-only subagents (goldfish) to review. Reports go to files. A synthesizer merges findings before the elephant reads them. EGM states no measure of context saved (research in `tasks/UND-009/egm-research.md`).
- Fact: PRD-029 already took the closed review bundle, the comprehension paraphrase, and the drift report from EGM. It left out the six-lens brainstorm, the README bootstrap ("peanuts" and "hay"), and hooks.
- Fact: nothing persists a repository map. `/dh:setup` records the stack and the test, lint, and build commands in `project.yaml`. `/dh:understand` maps one slice. The skill `repository-mapping` forbids a full-repository tour. `/dh:doctor` checks the environment only.
- Hypothesis: `MEMORY.md` growth is the largest cost in a long run, larger than any rule text. The hypothesis is wrong if a measured session shows that the report text in the chat is larger than the `MEMORY.md` read. Parts A and B each address one of the two.

## 2. Solution

After delivery:

- A specialist report above a size limit goes to a file in the task folder. The Coordinator gets the verdict, the path, and a short summary within a numeric cap.
- When a stage has more than one reviewer, one specialist merges the verdicts into one file. The Coordinator reads the merged file.
- `/dh:auto` reads only the named sections of each chained command, at the stage that needs them.
- Sections that serve one or two commands move out of `coordinator.md` into those commands.
- `MEMORY.md` has a size budget for each entry and for the file. The Coordinator offers consolidation when the file is over budget. Consolidation stays owner-triggered.
- `/dh:setup` can write a repository map. Later sessions and `repo-scout` reuse it.

Part A is R1 to R4 (specialist replies). Part B is R5 to R8 (what the Coordinator loads). Part C is R9 to R12 (repository bootstrap). Part D is R13 to R17 (limits, proof, and PRD-031 carry-over fixes).

## 3. Rules

### Part A · Artifact-first replies

- R1 Reply cap and report file. The dispatch contract (`.agents/coordinator.md`, Procedure, the item that lists the dispatch fields) gives a numeric reply cap. The default cap is 200 words. It limits every reply. One more rule governs where a report lives: a report of more than 40 lines goes to `<task folder>/reports/<role>-<n>.md`, and a report of 40 lines or fewer stays in the reply. The reply always stays within the cap, so a short report may also be shorter than 40 lines. The Coordinator may set a different cap in the dispatch and states the reason. Each role has its own reply shape:
  - `reviewer`: verdict, counts by severity, at most 5 blocking findings as `location · title` (the rest only counted). It has no Write tool, so it always returns inline within the cap, like `repo-scout`.
  - `qa-verifier`: matrix verdict (pass, fail, or not-run per check), counts of each, failed checks as `check · reason`, report path.
  - `builder`: slice status, files changed, `Red:` and check results, report path. It has no severity counts.
  - `docs-guide`: files written with paths, and what each holds. It has no verdict.
  - `debugger`: diagnosis status (confirmed, hypothesis, or inconclusive), the cause in one sentence, report path.
  - `repo-scout`: it has no Write tool, so it always returns inline within the cap, as a path list with one line for each path (see R2).
  Verifiable: `coordinator.md` contains the phrase `reply cap` and the numbers `200` and `40`. `builder.md`, `reviewer.md`, `qa-verifier.md`, `docs-guide.md`, `debugger.md`, and `repo-scout.md` each state their own reply shape in the Output format. `document-validator` keeps its own limit.
- R2 Who writes the report file. These writers keep their current write sets and add only the report path named in the dispatch:
  - `builder`, `qa-verifier`, `docs-guide`, and `debugger` write their own report file.
  - `reviewer` stays fully read-only. It has no Write tool and writes no report file (D1). It returns its reply inline (R1).
  - `repo-scout` has no Write tool. It returns inline within the cap and lists paths, not file bodies (R1). For a map longer than the cap, the Coordinator dispatches `docs-guide` (R10).
  - A CLI reviewer never writes. The `qa-verifier` or `builder` (lane `infra`) that runs the call writes the report file, as PRD-025 sets for the call itself.
  Verifiable: the `tools` of `reviewer.md` have no Write. The statement "review is read-only" remains true without exception.
- R3 Coordinator read discipline. The Coordinator reads a report file only to judge a blocking finding, to settle a dispute, or to make a decision that the reply cannot support. It does not read a report in full when the verdict is approved and no blocking finding exists. The "In-flight doubt" section and the `plan-loop.md` Procedure items that judge wave reports follow this rule. `MEMORY.md` records the report path and the verdict, not the report text. Verifiable: `coordinator.md` contains this rule once, and no section contradicts it. A search for "read the report" in `.agents/` and `.commands/` shows only places that name a condition.
- R4 Synthesizer. A review stage can have more than one reviewer entry in the dispatch. Then a `docs-guide` run (model `sonnet`) merges the reviewer outputs (the CLI report files written by the delegate and the inline reply of a `claude` reviewer, passed in the dispatch as quoted data) into one file `<task folder>/reports/<stage>-merged-<round>.md`. The Coordinator reads this file, not the raw reports. The merge rules:
  - Input is the reviewer outputs named in the dispatch (files, or an inline reply quoted as data), never the chat.
  - Findings with the same location and title become one finding. The merge lists each reviewer that raised it.
  - The merge never removes a blocking finding.
  - The merge never lowers a severity.
  - Each reviewer verdict stays verbatim in the file.
  - The merge adds no finding and no judgment. The Coordinator judges.
  With one reviewer entry, no synthesizer runs. The document stage keeps `<document>.review.md` as the merged file, so it needs no second merge. Verifiable: `coordinator.md` ("External CLI reviewers", merge step) names `docs-guide` as the merger for the `code` and `security` stages. `docs-guide.md` lists the merge task and its rules.

### Part B · What the Coordinator loads

- R5 `/dh:auto` stops inlining chained commands. At the start of each chained stage, the Coordinator reads only three sections of the chained command file: `Routing`, `Prerequisites`, and `Limits`. It skips `Output` and every section that the stage does not use. `auto.md` carries a table with one row for each chained stage (`discover`, `plan`, `build`, `secure`, `review`, `release`). Each row gives the file, the sections to read, and a byte budget. The Coordinator reads the stage sections once per run. A fresh session that resumes the run (`/dh:resume`) reads them again. The table rule replaces the sentence "Follow each chained command's own Routing and Prerequisites text inline". The chained gates, caps, and authorization checks do not change. Verifiable: `auto.md` contains the stage table with six rows. `auto.md` no longer says that the Coordinator reads `release.md` "inline in this same session" in full.
- R6 Move single-command sections out of `coordinator.md`. Each of these sections moves to the command that uses it, with the text unchanged except for references:
  - "Worktree" moves to `.commands/build.md`. `auto.md` points to that section.
  - "Consolidation on explicit request" moves to `.commands/consolidate-memory.md`.
  - "External CLI planners" moves to `.commands/plan-loop.md`.
  - "Drift report" moves to `.commands/review.md`.
  `coordinator.md` keeps one line for each moved section that names the command. The sections "Language", "Model selection", "External CLI reviewers" (it holds the anti-delegation clause), and "Proportional review and paid evaluation" stay. Target: `coordinator.md` is smaller by at least 5500 bytes than its size at 0.36.0 (31842 bytes). Measure with `wc -c`. Verifiable: `wc -c .agents/coordinator.md` is at most 26342. `go run ./cmd/dh validate` passes, including `checkAntiDelegationClause`. A search finds no reference to a moved section by its old place.
- R7 `MEMORY.md` budget. The Coordinator keeps each `MEMORY.md` entry (one bullet, or one paragraph under a heading) at or below 1200 characters. It keeps the whole file at or below 40000 bytes. Detail beyond an entry goes to a file in the task folder. The entry keeps the path and one sentence. Existing entries stay as they are: the disk is the truth, and the Coordinator does not rewrite old records to meet the budget. The rule applies to every new or changed entry. Verifiable: `.agents/coordinator.md` ("Memory ownership") and `templates/en/MEMORY.md` and `templates/pt-br/MEMORY.md` state both numbers.
- R8 Consolidation offer. When a read or `/dh:doctor` shows `MEMORY.md` over the file budget (R7), the Coordinator tells the owner once in that session and offers `/dh:consolidate-memory`. It never runs consolidation without the owner's request. `/dh:doctor` prints the byte size of `MEMORY.md` in the harness dir and a warning above 40000 bytes. Verifiable: a test of the `doctor` check, in `internal/`, for a file below, at, and above the budget. `consolidate-memory.md` states that the offer is the only automatic step.

### Part C · Repository bootstrap

- R9 Mechanism. The bootstrap extends `/dh:setup` as an opt-in second phase. It is not a new command. After the profile, `/dh:setup` offers to write the repository map. The owner accepts or declines. A re-run with the note `bootstrap` in the arguments refreshes the map without redoing the profile. The command counts stay at 19 + 1. Verifiable: `setup.md` contains the offer, the `bootstrap` argument, and the map path. The command count in `go run ./cmd/dh validate` is unchanged.
- R10 Map file and writer. The map is `<harness dir>/REPO-MAP.md`, a markdown file in the `language` of `project.yaml`. It has no template and no `project.yaml` key. The skill `project-onboarding` gets a section "Repository bootstrap" that defines the format and the procedure. A `docs-guide` run (model `sonnet`) with the skills `repository-mapping` and `project-onboarding` reads the repository and writes the file. Its reply follows R1. The map has at most 200 lines and these parts:
  1. Header: date and the `HEAD` sha. The Coordinator gets the sha from a `qa-verifier` run (`git rev-parse HEAD`) and passes it in the dispatch. The writer copies it and does not run git.
  2. Top-level layout: one line per directory.
  3. Entry points and how to run, test, lint, and build, with the source of each command.
  4. Architecture in at most 15 lines: the main components and how they connect.
  5. Risk areas: paths where a change needs extra care, with the reason.
  6. Conventions the code shows, each with one example path.
  Every claim cites a relative path. A claim the writer did not verify is marked `hypothesis`. When the repository has more than 15 top-level directories, the writer stops and the Coordinator asks the owner to name the scope (the EGM safety gate, adapted). The writer reads the repository and does not run project code or scripts. Verifiable: a run on a fixture produces a map of at most 200 lines with all six parts and a header sha.
- R11 Reuse. `repo-scout` reads `REPO-MAP.md` first when it exists and treats it as a hint. The Coordinator adds the map path to `sources` in dispatches to `repo-scout`, `implementation-planner`, and `builder`. The Coordinator does not read the map itself. `builder` and `implementation-planner` treat the map as a hint and verify it against the tree, as `repo-scout` does. `repository-mapping` gains one line: use the map as a pointer and verify any claim before it supports a plan. Verifiable: `repo-scout.md`, `builder.md`, `implementation-planner.md`, and `repository-mapping` each contain `REPO-MAP.md`. `coordinator.md` does not tell the Coordinator to read it.
- R12 Staleness and refresh. The Coordinator passes the current `HEAD` sha (from `qa-verifier`) in each `repo-scout` dispatch. When the header sha differs from it, `repo-scout` reports the map as stale and still uses it as a pointer. The refresh is owner-triggered (R9). `/dh:setup` may overwrite `REPO-MAP.md`, the only exception to its rule "never overwrite an existing record", because the file is derived data. In `repo` mode, `.harness/REPO-MAP.md` is versioned like other internal records. Verifiable: `setup.md` Limits names the exception for this one file. The `repo-scout` report has a line `Map: fresh | stale | none`.

### Part D · Limits and proof

- R13 No new agent, skill, command, template, or `project.yaml` key. The counts stay: 13 agents, 22 skills, 19 + 1 commands. The template count stays. `go run ./cmd/dh validate` passes the documented-counts check.
- R14 These safeguards stay in force:
  - Not-run honesty: a check that did not run is `not-run`, never `passed`.
  - Memory ownership: only the Coordinator writes `MEMORY.md`, `EPOCHAL.md`, and `RISKS.md`. A report file in the task folder is not one of them.
  - Independent review: a builder claim is not review, and a reviewer changes no repository file.
  - The review bundle excludes the Coordinator conversation (the context package rule in `AGENTS.md`).
  - The anti-delegation clause, verbatim in `coordinator.md` and `external-clis`.
  - Authorization: plan approval does not authorize commit, push, publish, or deploy.
  - A CLI reviewer run changes nothing in the tree. A change in the tree discards the verdict.
  Verifiable: `go run ./cmd/dh validate` passes. The builder lists each safeguard with the file that holds it in the pull request description.
- R15 Measurement. The first ticket records a baseline, and the last ticket records the result, in the Executed summary of this PRD. The records give: `wc -c` of `.agents/coordinator.md`; the bytes the Coordinator reads from the six chained commands in `/dh:auto` (before: the six files in full; after: the table sections only); and the size of the `MEMORY.md` template. No paid eval and no token-count tool is part of this PRD, and the Executed summary says that the figures are bytes, not tokens. Verifiable: the Executed summary has the baseline and result lines.
- R16 Cleanup premise. The same pull request updates every document that states a behavior this PRD changes (section 4) and regenerates `dist/` with `go run ./cmd/dh build`. A search for each moved section name finds no stale reference outside `docs/prd/`, `docs/adr/`, `CHANGELOG.md` history, and `evals/baselines/`. Because R8 changes `internal/`, the clean-clone proof in `AGENTS.md` applies before any push. Verifiable: `go run ./cmd/dh validate`, `go test ./...`, and a clean `git diff` of `dist/` after `go run ./cmd/dh build`.
- R17 PRD-031 carry-over fixes. Four defects from the PRD-031 delivery are fixed in this pull request:
  - `.skills/code-review/SKILL.md` asks again, in the review procedure, whether a plausible regression would fail a test. Verifiable: the file contains that question once.
  - `.skills/incremental-implementation/SKILL.md` restores the three safeguards of the official-page rule: skip the fetch when only one option is viable; send a conflict with local convention to the open-decision gate; put the URL in the report, not in a code comment. Verifiable: the file contains each safeguard once.
  - `docs/prd/PRD-031-rule-inventory.md` states the total as 55830 bytes. Verifiable: `grep -n "55830\|55,830" docs/prd/PRD-031-rule-inventory.md` shows the figure and no other total.
  - `.skills/regression-testing/SKILL.md` step 13 gives one sentence for each heuristic. Verifiable: step 13 has three sentences, one for each heuristic of PRD-031 R3.

## 4. Docs

The builder confirms each document with a search and edits it in the same pull request. Pairs update together (pt-br and `docs/en/`).

- `AGENTS.md` (bullets "Escopo do Coordenador", "Agentes", "Verificação padrão", and "Pasta global do harness")
- `.agents/coordinator.md`
- `.agents/builder.md`
- `.agents/reviewer.md`
- `.agents/qa-verifier.md`
- `.agents/docs-guide.md`
- `.agents/debugger.md`
- `.agents/repo-scout.md`
- `.agents/implementation-planner.md`
- `.commands/auto.md`
- `.commands/build.md`
- `.commands/consolidate-memory.md`
- `.commands/plan-loop.md`
- `.commands/review.md`
- `.commands/setup.md`
- `.commands/doctor.md`
- `.skills/project-onboarding/SKILL.md`
- `.skills/repository-mapping/SKILL.md`
- `.skills/code-review/SKILL.md`
- `.skills/security-review/SKILL.md`
- `.skills/incremental-implementation/SKILL.md` (R17)
- `.skills/regression-testing/SKILL.md` (R17)
- `docs/prd/PRD-031-rule-inventory.md` (R17)
- `templates/en/MEMORY.md`
- `templates/pt-br/MEMORY.md`
- `internal/` (the `doctor` size check and its test)
- `docs/tutorial.html`
- `docs/en/tutorial.html`
- `docs/novo-fluxo-prd.md`
- `docs/en/new-prd-flow.md`
- `docs/roadmap.html`
- `docs/en/roadmap.html`
- `CHANGELOG.md` (new entry, version bump, counts line)
- `docs/adr/ADR-010-lean-coordinator-context.md` (new: records the artifact-first rule, the section moves, and the bootstrap choice)

The builder also searches `adapters/`, `hooks/`, and `internal/` for code that reads the moved `coordinator.md` sections or the `MEMORY.md` layout. A match is in scope.

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| New `/dh:bootstrap` command | The fixed count 19 + 1 changes in `AGENTS.md`, ADR-003, `internal/kit/doc_counts.go`, roadmap and tutorial (pt and en), `marketplace.json`, and `REVIEW.md`. `/dh:setup` already owns first-time project knowledge. |
| Bootstrap through `/dh:understand` | The command maps one slice. A whole-repository mode blurs it and contradicts `repository-mapping`. |
| README tree ("peanuts" and "hay") | It adds files to the consumer repository. A single map in the harness dir has the same use with no repository change. |
| Stage-runner subagent for each `/dh:auto` stage | A subagent cannot start subagents, so a runner could not dispatch builders and reviewers. |
| Edit `MEMORY.md` automatically at the budget | Consolidation stays owner-triggered (`AGENTS.md`). The disk is the truth for records. |
| EGM six-lens brainstorm, EGM hooks | Excluded in PRD-029. They add context, not remove it. |
| A cheaper model for the Coordinator | Out of the owner's request. Model choice is not the cause of context growth. |
| Paid evals of context saved | Only on the owner's request (`AGENTS.md`). The measure is bytes (R15). |

### Risks and pending decisions

- Risk: the byte measure is a proxy for tokens. The real saving shows only in long runs. Mitigation: R15 states the figures as bytes. The owner judges by use.
- Risk: R6 saves about 5500 bytes, small beside the `MEMORY.md` read (185182 bytes in this project). Mitigation: R7 and R8 target the larger item. The existing `MEMORY.md` is not rewritten by this PRD. The owner can run `/dh:consolidate-memory`.
- Risk: the Coordinator trusts the reply and skips the file (R3). A summary can hide a blocking finding. Mitigation: the reply must list up to 5 blocking findings, and R4 forbids dropping or lowering one.
- Risk: a `reviewer` reply with no report file can lose detail. Mitigation: the reply keeps the verdict, the counts, and up to 5 blocking findings (R1). The Coordinator can ask for a narrower re-review.
- Risk: sections moved out of `coordinator.md` (R6) can be skipped when a command is run by hand without the Coordinator reading it. Mitigation: each command already reads its own text. The one-line pointer stays in `coordinator.md`.
- Risk: a stale map misleads planning (R12). Mitigation: it is a hint, the sha header marks staleness, and claims are verified.
- Risk: sections can change before the build starts. Mitigation: this PRD cites section names, not lines, and R6 measures from the 0.36.0 size.
- Suggestion (not a rule): after each `/dh:auto` stage, the status line could say that the run can continue in a fresh session with `/dh:resume`. The `Stage:` line already supports it. The Coordinator cannot measure its own context, so the PRD does not make this a rule.

### Decisions (closed 2026-10-31, owner: follow recommendations)

Each item is decided 2026-10-31 (owner: follow recommendations).

- D1 `reviewer` gets no Write tool and stays read-only. It returns its report inline (R1, R2). Decided 2026-10-31. Reason: a limit written only in prose was a confirmed exposure. A runtime-guard rule that allows one narrow reviewer write is a future candidate, out of scope here.
- D2 `docs-guide`, model `sonnet`, runs the synthesizer (R4). Decided.
- D3 `/dh:auto` reads the named sections through a stage table (R5). Resuming in a fresh session stays an owner choice and a suggestion. Decided.
- D4 Thresholds: reply cap 200 words, report to file above 40 lines, entry at most 1200 characters, file at most 40000 bytes, move at least 5500 bytes (R1, R6, R7). Decided.
- D5 `dh doctor` prints the `MEMORY.md` size and warns above budget (R8). Decided.
- D6 The bootstrap extends `/dh:setup`. The counts stay 19 + 1 (R9). Owner note 2026-10-31: a separate bootstrap command may be needed in the future; out of scope here. Decided.
- D7 The map is `<harness dir>/REPO-MAP.md`, written by `docs-guide`, with no template and no key (R10). Decided.
- D8 Minor version bump from 0.36.0 to 0.37.0. Decided.

Amendment 2026-10-31: D1 reverted after a security review (medium finding: untrusted PR text could reach `project.yaml` or `REVIEW.md` through a `reviewer` write). R1, R2, and the related Risk follow D1. R11 now says that `builder` and `implementation-planner` verify the map.

### Out of scope

- EGM hooks, the six-lens brainstorm, and the README tree.
- A model change for the Coordinator.
- Paid evals and any token-counting tool.
- Automatic `MEMORY.md` consolidation or rewrite of existing records.
- Changes to review rules, finding states, or round caps (PRD-031 owns them).

## Executed summary

Figures are bytes, not tokens. Measured on 0.36.0 before any PRD-032 edit (`wc -c`; a section runs from its `## <name>` header up to the next `## `).

- Baseline coordinator.md: 31842 bytes (`.agents/coordinator.md`).
- Baseline auto before: 36583 bytes, the six chained commands in full (discover 3339, plan 7427, build 4734, secure 3205, review 16061, release 1817). Sections `## Routing` / `## Prerequisites` / `## Limits` in bytes: discover 1463 / 266 / 415; plan 4546 / 530 / 416; build 2266 / 393 / 914; secure 757 / 469 / 884; review 8442 / 1037 / 1784; release 494 / 275 / 227.
- Baseline MEMORY template: `templates/en/MEMORY.md` 2579 bytes; `templates/pt-br/MEMORY.md` 2762 bytes.

Result (T-3217, commit a159079, same method and files as the baseline; `wc -c`):

- Result coordinator.md: 26338 bytes, 5504 below the baseline, within the R6 limit of 26342.
- Result auto after: 25524 bytes, the `Routing` / `Prerequisites` / `Limits` sections of the six chained commands (discover 2142, plan 5482, build 3553, secure 2110, review 11241, release 996), against 36583 before. The `auto.md` table budgets (2150, 5500, 3600, 2150, 11300, 1000) cover these sizes.
- Result MEMORY template: `templates/en/MEMORY.md` 2849 bytes; `templates/pt-br/MEMORY.md` 3050 bytes (they grew because they now state the R7 budget).

**Delivered:** The Coordinator now loads less (moved sections, an `auto.md` stage table, and a `MEMORY.md` budget with a doctor warning), roles reply within a 200-word cap with long reports in files, and `/dh:setup` offers an opt-in `REPO-MAP.md` bootstrap that planning and build roles use as a hint, in one pull request at 0.37.0 with no new agent, skill, command, template, or key.

**Rules:**

- R1 Honored. QA matrix `tasks/PRD-032/reports/qa-verifier-2.md` (pass); reply shapes in the six role files.
- R2 Honored. `reviewer.md` has no Write tool (D1 reverted by the 2026-10-31 amendment); QA matrix `qa-verifier-2.md`; security review `evidence/review-r2/`.
- R3 Honored. `qa-verifier-2.md` (rule once in `coordinator.md`, other mentions name a condition).
- R4 Honored. `qa-verifier-2.md` (`coordinator.md` names `docs-guide` as merger; `docs-guide.md` lists the rules).
- R5 Honored. `.commands/auto.md` stage table with six rows; sizes above; `qa-verifier-2.md`. The R5 byte budgets were not re-measured in that report; they were re-measured here.
- R6 Honored. 26338 bytes (limit 26342); `go run ./cmd/dh validate` passes with the anti-delegation clause hash unchanged (`qa-verifier-2.md`, `qa-verifier-3.md`).
- R7 Honored. Both numbers in `coordinator.md` and both MEMORY templates (`qa-verifier-2.md`).
- R8 Honored. `printHarnessFor` in `internal/doctor/harness.go` with a test below, at, and above the budget; `go test ./...` passes (`qa-verifier-3.md`). The `/dh:doctor` command text was not run (limitation in `qa-verifier-2.md`).
- R9 Honored. `qa-verifier-2.md`; command count stays 19 + 1.
- R10 Honored. The fixture map (30 lines, six parts, header sha equal to `git rev-parse HEAD`) was checked in `reports/qa-verifier-1.md` on an earlier state. It was not repeated on the final state. Observation in that report: the writer ran `git ls-files` beyond the sha.
- R11 Honored. `qa-verifier-2.md`.
- R12 Honored. `qa-verifier-2.md`.
- R13 Honored. Counts 13 agents, 22 skills, 19 + 1 commands, 14 templates; validate passes (`qa-verifier-2.md`, `qa-verifier-3.md`).
- R14 Honored in the tree: the safeguards are intact (validate passes, anti-delegation clause sha256 equal in `coordinator.md` and `external-clis`, critical paths untouched: `qa-verifier-3.md`). Pending: the list of safeguards with their files in the pull request body is not written yet; this is not claimed as done.
- R15 Honored. Baseline and result lines above; the figures are bytes, not tokens.
- R16 Honored. `dist/` rebuilt; clean-clone proof PASS in `qa-verifier-3.md` (7 binaries with the same sha256, validate 0 errors, `go test ./...` ok, only `generated_at` differs).
- R17 Honored. `qa-verifier-2.md`.

**Tickets:**

| Id | Lane | Verdict | Evidence |
|---|---|---|---|
| T-3201 | backend | done | baseline above |
| T-3202 | backend | done | `qa-verifier-2.md` (R17); `evidence/review-r3/` |
| T-3203 | backend | done | `go test ./...` in `qa-verifier-3.md`; `evidence/review-r3/` |
| T-3204 | backend | done | coordinator.md 26338 bytes; `evidence/review-r3/` |
| T-3205 | backend | done | `qa-verifier-2.md` (R1, R2, R4); `evidence/review-r3/` |
| T-3206 | backend | done | security `evidence/review-r2/`: `cli:codex/gpt-6.1-sol` found no vulnerability; the `claude` reviewer approved after F1-1 and F1-2 were fixed; 2 Low findings, non-blocking |
| T-3207 | backend | done | `qa-verifier-2.md` (R1 to R4); `evidence/review-r4/` (codex approve) |
| T-3208 | backend | done | `qa-verifier-2.md` (R9 to R12); `evidence/review-r4/` |
| T-3209 | backend | done | `qa-verifier-2.md` (R7, R8); `evidence/review-r4/` |
| T-3210 | backend | done | coordinator.md closes at 26338 bytes; `evidence/review-r4/` |
| T-3211 | backend | done | `auto.md` table; `qa-verifier-2.md` |
| T-3212 | test | pass (13 pass; R10, R14, R15, R16 pending then, closed above) | `reports/qa-verifier-2.md` |
| T-3213 | backend | done | `AGENTS.md`, `docs/adr/ADR-010-lean-coordinator-context.md`; `evidence/review-r5/` (opencode approve) |
| T-3214 | backend | done | docs and `CHANGELOG.md` 0.37.0; `evidence/review-r5/` |
| T-3215 | infra | done | version and `dist/` rebuild; `evidence/review-r5/` |
| T-3216 | test | PASS | `reports/qa-verifier-3.md` |
| T-3217 | backend | done | this section |

Reviews of the batch: `evidence/review-r3/` (opencode approve), `evidence/review-r4/` (codex approve), `evidence/review-r5/` (opencode approve, docs and version). Final feature review: `evidence/review-r6/` (codex; finding F6-1 settled as intentional per the PLAN line 72) and `evidence/review-r7/` (opencode approve).

Document revalidation after the 2026-10-31 amendment (D1 reverted): `evidence/doc-r4/` (agy, changes required). The Coordinator fixed points P8 to P10 without a further round. Recorded for the owner.

**Docs:** `AGENTS.md` (current rules), `docs/adr/ADR-010-lean-coordinator-context.md`, the tutorials and the PRD flow, both roadmaps, `CHANGELOG.md` 0.37.0, the `dist/` rebuild, and this PRD (baseline, result, Executed summary, Status). Source: PLAN.md ownership table, T-3213 to T-3215.

**Out of scope:** unchanged from the section above: EGM hooks and the six-lens brainstorm, a Coordinator model change, paid evals and token counting, automatic `MEMORY.md` consolidation, and review-rule changes. A separate bootstrap command and a runtime-guard rule for a narrow reviewer write remain future candidates.
