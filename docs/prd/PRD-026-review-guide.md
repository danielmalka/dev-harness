# PRD-026 · Review guide (REVIEW.md) and PR review mode

| Field | Value |
|---|---|
| Status | delivered on 2026-10-21 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-21 / 2026-10-21 |
| Tickets | T-2601, T-2602, T-2603, T-2604 (planned) |
| Origin | Owner chat, 2026-10-21 (discovery DISC-006; video "How Anthropic Does Code Review (What You Can Copy)", pasquadev, 2026-10-21, https://youtu.be/c6xOFjM2IbE) |
| Depends on | nothing |
| Enables | nothing |

## 1. Problem

The code review of the kit has no project-specific guide. The reviewer judges a diff with generic rules only. Fact: `.skills/code-review/SKILL.md` and `.agents/reviewer.md` name no per-project file that says what to look for.

Four gaps follow. Fact: each gap is visible in the current kit text.

- A repeated false finding costs a round each time. Nothing records that a pattern is noise, so the next review reports it again.
- The reviewer reports one instance of a bug, and the next round finds the next instance. Rounds are spent one case at a time.
- A re-entry review has no stated scope. The reviewer can re-read the whole change instead of the delta.
- A pull request cannot be reviewed as a pull request. `/dh:review` takes a scope, a ref, or `feature PRD-<n>`; it has no mode that reads a PR and its declared intent.

Also: the merge step of `/dh:auto` (`.commands/auto.md`, step 5 of the git sequence) merges on green CI plus an `approve` verdict. It has no stop for a change to a path the owner calls critical.

The video is the source of the ideas. Its numbers and cited papers are unverified. No rule in this PRD depends on them.

## 2. Solution

After delivery, each project has one review guide, `REVIEW.md`, in English, in the project repository. The path is recorded in `project.yaml` of the harness dir. The reviewer reads the guide and ties each finding to a rule.

Flow after delivery:

1. `/dh:setup` creates `REVIEW.md` from `templates/en/REVIEW.md` when the project has none, and writes its path (relative to the repository root) to `project.yaml`. This holds in mode `repo` and in mode `global`.
2. `/dh:review` reads the guide. Each finding cites the guide rule it comes from, or says no rule applies. A finding that matches an Ignore entry is reported as ignored-by-rule and does not spend a round.
3. The reviewer reports the whole class of a bug at once, labels a finding on untouched lines as "pre-existing", and on re-entry reviews only the delta plus the status of each earlier blocking finding.
4. `/dh:review pr <n>` reviews a GitHub pull request read-only. The PR body and any linked PRD or ticket are the requirement source. The report stays local. A comment on the PR is posted only with explicit owner authorization.
5. Only Critical and Major findings go to the builder. Minor and nit findings stay in the report.
6. The Coordinator records what happened to each finding (fixed, intentional, noise) in `MEMORY.md`. When the same pattern is classified noise 3 times, the Coordinator proposes an Ignore line or a lint rule to the owner. It never applies the proposal itself.
7. When a PR touches a path listed under "Critical paths" in `REVIEW.md`, the flow stops before the merge and asks the owner, even with a green gate and an approved QA verdict.

This repository is a consumer of the kit. The same batch adds a real `REVIEW.md` at the repository root.

Release 0.32.0 (author's proposal), one PR.

## 3. Rules

Marks: R1-R14 and R16-R18 carry owner decisions of 2026-10-21 (D1-D6, A1, A2). R5-R9 and the file names are the author's wording of the discovery brief. R23-R25 and the changes to R5, R12, and R13 come from the code and security review of the implementation (2026-10-21); they need no new owner decision. Each rule states how to verify it.

- R1 Template. `templates/en/REVIEW.md` exists. It has these parts: (a) a short "How to use" note; (b) a rules section where each entry has the fields Symptom, Why it matters, and Example with `file:line`; (c) an "Ignore" section where each entry reads "ignore this pattern in this kind of file"; (d) a "Critical paths" section that lists path globs. No `templates/pt-br/REVIEW.md` exists. Verifiable: the file exists with the four parts; `ls templates/pt-br/` shows no `REVIEW.md`.
- R2 Setup creates the guide. When `/dh:setup` runs and the project has no review guide file, setup creates `REVIEW.md` at the repository root from `templates/en/REVIEW.md` and writes its path to `project.yaml` in the key `review_guide` (author's proposal for the key name). The path is relative to the repository root. The file is always in English, whatever `language` says. Verifiable: the text of `.commands/setup.md` and `.skills/project-onboarding/SKILL.md` states these steps, the key name, and the English rule.
- R3 Both modes. R2 holds in mode `repo` and in mode `global`. In mode `global`, `REVIEW.md` is the only file setup writes in the repository tree. The existing sentences "nothing is written in the repository tree" (`.commands/setup.md`, `.skills/project-onboarding/SKILL.md`, `docs/tutorial.html`, `docs/en/tutorial.html`) say this exception. Verifiable: `grep -n -i 'repository tree' .commands/setup.md .skills/project-onboarding/SKILL.md docs/tutorial.html docs/en/tutorial.html` shows each hit with the exception stated.
- R4 Setup never overwrites. When the project has the file but no key, setup records the key for the existing file and leaves the file untouched. When the key points to an existing file, setup changes nothing. A dispatched specialist writes the file, not the Coordinator (the file is new and larger than the small-correction bound of PRD-025). Verifiable: the setup text states the three cases and names the writer; the setup report format lists `REVIEW.md` under Created or Left untouched.
- R5 Review is read-only and points to setup. When `/dh:review` runs and `project.yaml` has no `review_guide` key, or the file at that path does not exist at the trusted revision (R23), or the path is unsafe (R24), then the report says "no review guide", gives the reason, and points to `/dh:setup`. The review does not create the file. The review still runs with the generic rules, and the report marks the guide axis as not-run with the reason. Verifiable: `.commands/review.md` states this; `.skills/code-review/SKILL.md` has a matching step.
- R6 Findings cite the guide. When the `reviewer` agent reviews code with a guide present, then it reads the guide at the path from `project.yaml` first. Each finding carries a line `Guide rule: <entry title>` or `Guide rule: none applies`. Verifiable: the `code-review` output format and the `reviewer` agent output format have this line; the Example in `code-review` shows it.
- R7 Ignore entries. When a finding matches an Ignore entry of the guide, then the report lists it in a separate "Ignored by rule" list with the entry title. It is not graded, not blocking, and does not spend a round. Verifiable: the `code-review` output format has the list; the Coordinator text on blocking findings (`.agents/coordinator.md`, "Proportional review") counts only Critical and Major findings that are not ignored-by-rule.
- R8 Whole class, pre-existing. When the reviewer finds a bug, then it reports every instance of that class inside the review scope in one finding, with each location listed. A finding on lines the change did not touch is labeled "pre-existing" and is non-blocking, unless the change makes that code reachable or worse; then it keeps its severity and the finding says why. Verifiable: both rules appear in the procedure of `.skills/code-review/SKILL.md`, and the false-positive control step (step 9) agrees with the label.
- R9 Re-entry scope. When a review runs again on the same change, then the reviewer reads (a) the delta since the last reviewed state and (b) the status of each earlier blocking finding (fixed, not fixed, withdrawn). It reports on those two items only, plus any new Critical or Major finding it finds in the delta. The Coordinator names the last reviewed state in the dispatch. Verifiable: the `code-review` skill has a re-entry section with these two items; `.agents/coordinator.md` names the last reviewed state in the re-entry dispatch.
- R10 Noise record. The Coordinator records, for each review finding, one outcome in `MEMORY.md`: fixed, intentional, or noise. When the same pattern is classified noise 3 times, then the Coordinator proposes to the owner an Ignore line for `REVIEW.md` or a lint rule. The Coordinator never applies the proposal itself. A dispatched specialist writes the change only after the owner approves. Verifiable: `.agents/coordinator.md` states the record, the threshold 3, the proposal, and the prohibition.
- R11 Fix dispatch carries only Critical and Major. When the Coordinator sends review findings to a builder, then the dispatch contains only Critical and Major findings. Minor and nit findings stay in the review report. Verifiable: `.agents/coordinator.md` (Procedure item 8 and the review loop text) and every other place that describes the builder fix dispatch (`grep -n -i 'findings' .commands/build.md .commands/fix.md .commands/review.md .agents/builder.md`) say this; none says "all findings".
- R12 Critical paths stop the merge. When a PR touches a path that matches a glob under "Critical paths" in the project's `REVIEW.md`, then the flow stops before the merge and asks the owner for an explicit ok, even when the gate is green and the QA verdict is approved. The question names the matched paths. The path match is computed by a dispatched `qa-verifier` from the PR file list. When `REVIEW.md` is absent or the section is empty, no stop applies. This rule overrides the owner's global auto-merge instruction for those paths only; for every other path the global instruction is unchanged. A change to the review guide stops the merge in the same way (R23). Verifiable: `.commands/auto.md` (merge step 5 of the git sequence) and the Coordinator completion text state the stop, the owner question, the scope of the override, and the guide-change stop.
- R13 PR review mode. PR mode starts only when the first argument token is exactly `pr`. When `pr` is followed by one ASCII digit string (`[0-9]+`) and nothing else (for example `pr 44`), the command runs in PR mode. When `pr` is followed by anything else, the command reports the malformed trigger and stops (fails closed). Any other argument (for example `prod`, `preview/x`, `prototype.go`) is the default mode, as before this PRD. In PR mode a dispatched `qa-verifier` reads the PR with `gh` read-only (`gh pr view <n> --json ...`, `gh pr diff <n>`), not the Coordinator. The command takes the PR body and any linked `PRD-<digits>` or `T-<digits>` reference found in it as the requirement source, and dispatches `reviewer` in mode `code` as in the default mode. The report states the commit it read. It does not check out a branch and does not write to the tree. No new command exists. Verifiable: `.commands/review.md` has the mode, the exact-token rule, the fail-closed rule, the default-mode fallthrough for `prod`, `preview/x`, and `prototype.go`, the read-only `gh` calls, and the `qa-verifier` dispatch; the existing modes are unchanged.
- R14 Posting needs authorization. In `pr <n>` mode the report stays local. The command posts a comment on the GitHub PR only when the owner authorizes that post explicitly in the current conversation. The authorization covers one post. Verifiable: `.commands/review.md` states this and names the post as an external action.
- R15 Template list in `dh validate`. The slice `englishOnlyTemplates` in `internal/kit/validate.go` (line 35) lists `"REVIEW.md"` next to `"PRD.md"`. The comment above it (line 34) says that the PRD and the review guide are always written in English. A test case in `internal/kit/validate_test.go` covers `REVIEW.md` the same way it covers `PRD.md` (a copy in `templates/pt-br/` is an error; a missing copy in `templates/pt-br/` is not). Verifiable: `go test ./internal/kit/...` passes; `go run ./cmd/dh validate .` passes with `templates/en/REVIEW.md` present and no pt-br copy.
- R16 This repository is a consumer. The batch adds `REVIEW.md` at the repository root, in English. It has at least 3 real entries for this kit, each with Symptom, Why it matters, and Example with `file:line`: generated `dist/` edited by hand; a count stated in docs that drifts from the real count (agents, commands); a private path in a PRD or doc. It has an Ignore section and a Critical paths section (owner decision D7: the critical globs are exactly `adapters/claude-code/plugin/hooks/**` and `.skills/external-clis/**`). The key `review_guide: REVIEW.md` is written in the `project.yaml` of this project's harness dir, in whichever mode the project resolves to. The Coordinator writes that one line; it is a Coordinator step, not a builder write. Verifiable: the file exists with the parts above; `dh harness-path --json` gives the harness dir, and its `project.yaml` has the key.
- R17 Out of scope items stay out. This PRD adds none of the items of section 5. The round cap stays 6 (`.agents/coordinator.md`, "Proportional review and paid evaluation"). The merged verdict rule "any blocking finding from any reviewer blocks" is unchanged. Verifiable: `git diff main...HEAD -- .agents/coordinator.md` shows no change to the number 6, no majority-vote text, and no `dh review run`.
- R18 Counts. The counts stay at 13 agents and 19 commands + 1 mod command. No new agent, skill, or command. Verifiable: `go run ./cmd/dh validate .` passes, and `git diff --name-status main...HEAD -- .agents .skills .commands` lists no added or deleted file.
- R19 No trace of the guideline file. The guideline synthesis file of the video (a pt-br file in `docs/`) and its `Zone.Identifier` sibling are not in the repository at closing. The owner moves the file to the vault; the Coordinator handles that outside the kit. No file, `dist/` copy, doc, or PRD text cites its path. The video is cited by title, channel, date, and URL only. Verifiable, all with no output: `git ls-files | grep -i 'DIRETRIZES-CODE-REVIE[W]'`; `find . -name '*Zone.Identifier*' -not -path './.git/*'`; `grep -rIl 'DIRETRIZES-CODE-REVIE[W]' . --exclude-dir=.git`; `git status --short` shows neither file.
- R20 Cleanup. Every stale reference that this change creates is removed (AGENTS.md, "Premissa de limpeza"). Examples: sentences that say setup writes nothing in the repository tree (R3); sentences that say `/dh:review` takes only a scope, a ref, or `feature PRD-<n>` (R13); sentences that say "the PRD is the only English-only template" (`grep -rn -i 'only.*PRD\|PRD.*only' .skills .commands .agents templates docs README.md README.en.md`); setup `writes:` metadata in `.commands/setup.md` that omits `REVIEW.md`; `.commands/review.md` text that lists the allowed arguments without `pr <digits>`. Verifiable: the ticket author runs each grep and lists the hits and the action taken in the ticket report.
- R21 Release and generated files. Release 0.32.0 in three places: `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json`, and the source `ref` (`v0.32.0`). `dist/` is not edited by hand; `go run ./cmd/dh build` regenerates it, and a second build changes nothing except the `generated_at` line. `go run ./cmd/dh validate .` passes with 0 errors and `go test ./...` is green. Clean-clone proof applies (AGENTS.md, "Prova de clone limpo"), because the PR touches `internal/` (`validate.go`, `build.go`) and `dist/*/bin`: `git clone . <tmp>`, `go run ./cmd/dh build` there, and `sha256sum` of the 5 binaries equals the tree. Verifiable: run the commands and compare with `git diff`.
- R22 No private path in the PRD, the tickets, the new `REVIEW.md`, or any new text. The text uses `<harness dir>` and `<home>`. Verifiable: `go run ./cmd/dh validate .` (`checkPrivatePaths`) passes.
- R23 Trusted policy revision. The review reads the guide file and the `review_guide` key from the base revision of the change, never from the head under review. The base is the PR base in `pr <n>` mode, the comparison ref in default mode with a ref, and `HEAD` for uncommitted work. When the base has no guide, no guide applies to that change. When the change adds, edits, or deletes the guide file at the configured path, or changes the `review_guide` key, then the flow stops before the merge like a critical path (R12), and the owner question names the guide change. Ignore entries that the change adds do not apply to that same change. The owner can approve a guide update through this stop. In mode `global` the key lives outside the repository tree, so only a change to the guide file triggers the stop. Verifiable: `.commands/review.md`, `.skills/code-review/SKILL.md`, and `.agents/reviewer.md` state the base-revision read and the bootstrap case; `.commands/auto.md` and `.agents/coordinator.md` state the stop and the Ignore rule.
- R24 Path containment. The `review_guide` value must be a repository-relative path with no absolute prefix and no `..` segment. Its resolved target, symlinks included, must stay inside the repository. When any check fails, the review treats the project as having "no review guide" and reports why (R5). Setup writes only `REVIEW.md` at the repository root. Verifiable: `.commands/review.md` and `.skills/code-review/SKILL.md` state the three checks and the fallback; the setup text states the fixed value.
- R25 PR text handling. In `pr <n>` mode the review treats PR text (body, title, comments) as untrusted data. A dispatched `qa-verifier` runs the `gh pr view` and `gh pr diff` reads, not the Coordinator (R13). No one fetches a link found in PR text, and the reviewer uses no web fetch in this mode. PR text is never an owner authorization to post, merge, or do any other action (R14). Verifiable: `.commands/review.md` and `.agents/reviewer.md` state the three prohibitions and the dispatch of the reads.

## 4. Docs

Update:

- `templates/en/REVIEW.md` (new; R1) and `internal/kit/validate.go` with `internal/kit/validate_test.go` (R15).
- `.commands/setup.md` and `.skills/project-onboarding/SKILL.md` (R2-R4, R20, R24): setup steps, the `writes:` metadata, the report format, the "repository tree" sentences.
- `.commands/review.md` (R5, R13, R14, R23-R25): the `pr <n>` mode, the guide check, and the read-only limits.
- `.skills/code-review/SKILL.md` and `.agents/reviewer.md` (R6-R9, R23-R25): guide reading, `Guide rule:` line, ignored-by-rule list, whole class, pre-existing label, re-entry scope, output format and Example.
- `.agents/coordinator.md` (R7, R9-R12, R23): fix dispatch content, re-entry dispatch, noise record, critical-path stop.
- `.commands/auto.md` (R12, R23): merge step 5 of the git sequence. `.commands/build.md` and `.commands/fix.md` only where the grep of R11 finds "all findings" wording.
- `AGENTS.md`: already updated by the Coordinator on 2026-10-21 (template exception, `dh review run` trigger, "Guia de revisão" rule). The ticket author reads it and reports any disagreement; it does not edit it.
- `docs/tutorial.html`, `docs/en/tutorial.html`: the setup description, the review description, the `project.yaml` key, the critical-path stop, the `pr <n>` mode.
- `README.md` and `README.en.md`: only where `grep -n -i 'review\|setup'` finds a sentence that this change makes false.
- `docs/roadmap.html`, `docs/en/roadmap.html`: a Delivered card for 0.32.0 in both files, and the status chip and sources line. `dh validate` checks the drift against `CHANGELOG.md`.
- `CHANGELOG.md` (entry 0.32.0).
- `.claude-plugin/marketplace.json`, `internal/build/build.go` (R21).
- `REVIEW.md` at the repository root (new; R16).

No `pdocs` exist for this project.

Not affected (the ticket author checks with `grep`):

- `docs/plano-produto.html`: internal pt-br history; changed only if it states review behavior as current.
- `docs/adr/`: no new ADR unless the owner asks.
- `docs/inicio-rapido.html`, `docs/en/quick-start.html`: changed only if they describe setup output.
- `.skills/external-clis/`: stays in prose (D2); no change.

## 5. Out of scope

- Majority vote among reviewers. The rule "any blocking finding from any reviewer blocks" stays.
- Incident evals (`evals/incidentes/`). Paid evaluation runs only when the owner asks (`AGENTS.md`).
- A CI review job per PR. The CI runs without an API key on purpose.
- Risk-classifier agents.
- A merge queue.
- Feature flags and gradual rollout.
- Post-merge metrics.
- The `dh review run` subcommand. `external-clis` stays in prose. The trigger for it (a CLI call that needs a kill or a manual re-run inside a batch) is only in `AGENTS.md`.
- A 3-round cap, and any change to the round cap of 6.
- A new agent, skill, or command, and any change of the counts.
- A `dh doctor` check for the guide (see Appendix, suggestion S1).

## 6. Owner decisions

- D1, 2026-10-21: extend `/dh:review` with a `pr <n>` mode. No new command; the count stays 19 + 1. The PR body plus any linked PRD or ticket is the requirement source. The report stays local. A comment on the PR needs explicit authorization (R13, R14).
- D2, 2026-10-21: `external-clis` stays in prose. The `dh review run` trigger is recorded in `AGENTS.md` only. Out of scope here.
- D3, 2026-10-21: `REVIEW.md` has a critical-paths section. A PR that touches a listed path stops before the merge and asks the owner. This overrides the global auto-merge instruction for those paths only (R12).
- D4, 2026-10-21: Minor and nit findings stay in the report. Only Critical and Major go to the builder (R11).
- D5, 2026-10-21: the round cap stays 6 (R17).
- D6, 2026-10-21: the guideline file and its `Zone.Identifier` sibling leave the repository at closing (R19).
- A1, 2026-10-21: `REVIEW.md` is always added to the project. Its path is recorded in `project.yaml` of the harness dir the project resolves to, relative to the repository root. The file lives in the project repository (R2-R4, R16).
- A2, 2026-10-21: `REVIEW.md` is always in English. The template exists only in `templates/en/` (R1, R15).
- D7, 2026-10-21: the Critical paths of this repository's `REVIEW.md` are exactly `adapters/claude-code/plugin/hooks/**` (runtime guard, RISK-001) and `.skills/external-clis/**` (CLI reviewer isolation, RISK-002). A release that only bumps the version and the binaries does not stop (R16).

## Appendix

### Suggested tickets

| Ticket | Lane | Executor | Write set |
|---|---|---|---|
| T-2601 | `backend` | `harness-maintainer` with `harness-authoring` | `templates/en/REVIEW.md`; `.commands/setup.md`; `.skills/project-onboarding/SKILL.md`; `.commands/review.md`; `.skills/code-review/SKILL.md`; `.agents/reviewer.md`; `.agents/coordinator.md`; `.commands/auto.md`; `.commands/build.md` and `.commands/fix.md` if R11 grep hits. Rules R1-R14, R17, R18, R20, R23-R25 |
| T-2602 | `backend` | `builder` (Go) | `internal/kit/validate.go`; `internal/kit/validate_test.go`; `internal/build/build.go` (version); `.claude-plugin/marketplace.json`. Rules R15, R21. Clean-clone proof applies (touches `internal/` and `dist/*/bin`) |
| T-2603 | `backend` (kit source and docs) | `harness-maintainer` with `harness-authoring`, then `docs-guide` for the docs | `REVIEW.md` (root); `docs/tutorial.html`; `docs/en/tutorial.html`; `README.md`; `README.en.md`; `docs/roadmap.html`; `docs/en/roadmap.html`; `CHANGELOG.md`. Rules R3, R16, R20, R22 |
| T-2604 | `teste (unitário e integração)` | `qa-verifier` | Test of R15 in `internal/kit/validate_test.go` (with T-2602); runs of R18, R19, R21, R22; a dry read of R2-R14 and R23-R25 text against the greps listed in each rule. No production code |

Order: T-2601 and T-2602 may run in parallel only if the write sets stay disjoint (they do). T-2603 follows T-2601 (docs describe final text). T-2604 follows all. T-2602 and T-2604 write the same test file: T-2602 first.

Coordinator step (not a builder write): after `REVIEW.md` exists at the repository root, the Coordinator writes the line `review_guide: REVIEW.md` in `project.yaml` of this project's harness dir (R16). The Coordinator also moves the guideline file out of the repository (R19).

The release gate sequence is in R21: `dh validate`, `go test ./...`, `go run ./cmd/dh build` with no `dist/` diff, then the clean-clone proof, then commit candidate; push, PR, and tag only after the proof.

### Discarded alternatives

| Alternative | Why not |
|---|---|
| A new command `/dh:review-pr` | The owner chose a mode of `/dh:review` (D1). The count 19 + 1 stays. |
| `REVIEW.md` in `pt-br` too | The owner chose English only (A2). |
| `REVIEW.md` stored in the harness dir | The owner chose the project repository (A1). The reviewer in a clone or a CLI sees it with the tree. |
| Review creates the missing file | Review is read-only. Setup owns creation (R5). |
| Forward Minor findings to the builder | The owner chose to keep them in the report (D4). A builder fixes a nit and the review finds another nit. |
| Cap of 3 rounds | The owner kept 6 (D5). |

### Risks and pending decisions

- Risk: setup in mode `global` now writes in the repository tree, which the docs call forbidden · mitigation: R3 changes every sentence, and R20 greps for the old wording.
- Risk: the runtime mod warns when the main session writes a new file outside the harness dir · mitigation: R4 sends the write to a specialist.
- Risk: the critical-path match uses globs, and a glob style is not defined by the kit · assumption: the ticket author uses the same glob style as `.gitignore` and writes one example in the template. Hypothesis: this is enough for this repository's paths (R16).
- Risk: in `pr <n>` mode the local checkout may not be the PR head, so the reviewer's caller tracing reads the wrong tree · mitigation: R13 requires the report to state the commit it read; the reviewer marks traced paths it could not verify as limited.
- Assumption: the key name is `review_guide` and the release is 0.32.0 (the owner did not object, 2026-10-21).
- Suggestion S1 (not a requirement): a `dh doctor` check that warns when `project.yaml` has no `review_guide` key or the file is missing. A1 is met without it, because R5 and R2 cover the flows. It adds Go code and a test, so it is left to the owner.
- Suggestion S2 (not a requirement): a limit on comments per PR in `pr <n>` mode, to cut noise. The video suggests it; the owner has not decided.
- Pending decision: none blocks the tickets.

## Executed summary

**Delivered**: Every project gets an English `REVIEW.md` created by setup with `review_guide` in `project.yaml`; the reviewer reads it from the base revision and ties each finding to a rule; `/dh:review pr <n>` reviews a GitHub PR read-only; only Critical and Major findings go to the builder; critical paths and guide changes stop the merge for the owner.

**Rules**:

Evidence for all rules: QA rounds 2-5 on the final state (`dh validate` 0 errors and 0 warnings, `go test ./...` 7 packages ok, build idempotent with only `generated_at` changing, `claude plugin test` 235 pass and 0 fail); final-feature adherence review FEATURE APROVADA; code review `cli:codex/gpt-6-luna` approve; security review `cli:codex/gpt-6.1-sol` no demonstrated vulnerability, and `claude` pass with residual risk. The PRD was amended after the implementation review (R5, R12, R13 changed; R23-R25 added) and revalidated once (approved).

- R1: Honored. `templates/en/REVIEW.md` has the four parts; `templates/pt-br/` has no `REVIEW.md`.
- R2: Honored. `.commands/setup.md`, `.skills/project-onboarding/SKILL.md` state the steps, the key and the English rule.
- R3: Honored. The "repository tree" sentences in `.commands/setup.md`, `.skills/project-onboarding/SKILL.md`, `docs/tutorial.html`, `docs/en/tutorial.html` state the exception.
- R4: Honored. `.commands/setup.md` states the three cases and the specialist writer.
- R5: Honored. `.commands/review.md`, `.skills/code-review/SKILL.md` (no guide, guide axis not-run, points to setup).
- R6: Honored. `Guide rule:` line in `.skills/code-review/SKILL.md` and `.agents/reviewer.md`.
- R7: Honored. "Ignored by rule" list in `.skills/code-review/SKILL.md`; `.agents/coordinator.md` counts only non-ignored Critical and Major.
- R8: Honored. Whole-class and pre-existing rules in `.skills/code-review/SKILL.md`.
- R9: Honored. Re-entry section in `.skills/code-review/SKILL.md`; `.agents/coordinator.md` names the last reviewed state.
- R10: Honored. `.agents/coordinator.md` (record, threshold 3, proposal, prohibition).
- R11: Honored. `.agents/coordinator.md`, `.commands/build.md`, `.commands/fix.md`, `.commands/review.md`, `.agents/builder.md`; none says "all findings".
- R12: Honored. `.commands/auto.md` and `.agents/coordinator.md`; QA dry run touched no critical glob.
- R13: Honored. `.commands/review.md` (exact `pr` token, fail closed, default-mode fallthrough, read-only `gh`, `qa-verifier` dispatch).
- R14: Honored. `.commands/review.md` (post needs explicit authorization, one post).
- R15: Honored. `internal/kit/validate.go` (`englishOnlyTemplates`) and `internal/kit/validate_test.go`; tests mutation-checked, `go test ./...` ok.
- R16: Honored. Root `REVIEW.md` with 4 rules, 2 Ignore entries, and the two D7 globs; `review_guide: REVIEW.md` in `project.yaml`.
- R17: Honored. Round cap 6 and the merged verdict rule unchanged in `.agents/coordinator.md`; no majority vote, no `dh review run`.
- R18: Honored. `dh validate` shows 13 agents and 19 commands + 1 mod command; no added or deleted file in `.agents`, `.skills`, `.commands`.
- R19: Honored. The four greps and `git status --short` are empty.
- R20: Honored. Stale-wording greps run by T-2601 over the kit sources and by QA rounds 2 and 3 over `docs/`, the READMEs and `CHANGELOG.md`; hits fixed (see Docs).
- R21: Honored for version 0.32.0 (`internal/build/build.go`, `.claude-plugin/marketplace.json`, `ref`) and build. The clean-clone sha256 proof is pending at the time of writing; the release step runs it before push.
- R22: Honored. `checkPrivatePaths` passes in `dh validate`.
- R23: Honored. `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/reviewer.md`, `.commands/auto.md`, `.agents/coordinator.md`.
- R24: Honored. `.commands/review.md`, `.skills/code-review/SKILL.md` (three checks and fallback).
- R25: Honored. `.commands/review.md`, `.agents/reviewer.md` (three prohibitions, `qa-verifier` reads).

**Tickets**:

- T-2601 backend, `harness-maintainer`: done. Evidence: QA rounds and the three reviews above.
- T-2602 backend, `builder` (Go): done. Evidence: R15 tests (mutation-checked), version 0.32.0.
- T-2603 backend, `harness-maintainer`: done. Root `REVIEW.md` with 4 rules, 2 Ignore entries, and the two D7 globs.
- T-2604 teste, `qa-verifier`: done. Rounds 1-5, final PASS.

**Docs**: `CHANGELOG.md` (0.32.0); `docs/roadmap.html` and `docs/en/roadmap.html` (0.32.0 card, 14 templates); `docs/tutorial.html` and `docs/en/tutorial.html`; `docs/tutoriais/templates.{md,html}` and `docs/en/tutorials/templates.{md,html}`; `README.md` and `README.en.md`. No `pdocs` for this project.

**Out of scope**: The section 5 items stayed out: majority vote, incident evals, CI review job, risk-classifier agents, merge queue, feature flags, gradual rollout, post-merge metrics, `dh review run`, a 3-round cap, a new agent, skill or command, and a `dh doctor` check. Residuals: S3 - the `qa-verifier` that reads PR text inherits all tools and the controls are prose only (owner decision pending); a rename may evade critical-path glob matching (Low, left by D4); posting is gated by prose authorization (accepted).
