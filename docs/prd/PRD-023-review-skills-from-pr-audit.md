# PRD-023 · Add review and test checks from two external audit skills

| Field | Value |
|---|---|
| Status | delivered on 2026-10-16 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-16 / 2026-10-16 |
| Tickets | T-2301 (planned) |
| Origin | Owner chat, 2026-10-16 (scope decided in chat); amended the same day with the `security-audit` items |
| Depends on | nothing |
| Enables | nothing |

## 1. Problem

- Two external skills, `pr-audit` and `security-audit` (akitaonrails/my-skills, commit 7b8a5e5), find defects that the kit's review skills miss. Fact: the Coordinator compared them with the kit.
- Gaps in `code-review`:
  - It does not tell the reviewer to look for sibling call sites of a fixed bug. The idea exists only in `.skills/systematic-debugging/SKILL.md`, step 8.
  - It does not tell the reviewer to find every consumer of a changed data shape, including mirrors that people maintain by hand.
  - It checks whether an added file is safe. It does not check whether the file belongs in the tree.
  - It does not require a new entry point to prove again the invariants of the old ones.
- Gaps in `security-review`:
  - It has no gate for a hostile diff or for the supply chain.
  - It has no rule for a security fix or for the proof of a security report.
  - It accepts a clean scanner result even when the diff changed the scanner.
  - It has no row for a disabled control, a model boundary, a leaked secret, or a fail-closed identity check.
- Gap in `regression-testing`: step 6 accepts "fails before the fix" even when the failure is a compile or import error.

## 2. Solution

Add each idea to an existing skill as a few short lines. Add no skill, no agent and no command. Examples in the shipped skill text are generic. File names of this repository stay in this PRD only.

`code-review`:

1. Steps 4 and 5: when a change fixes a symptom, the reviewer searches for sibling call sites of the same helper, guard or normalization. A point fix that leaves a twin site broken is Major, unless the scope is justified.
2. "Contracts" bullet: when a shared data shape changes, the reviewer lists every consumer. This includes mirrors that people maintain by hand, such as translated docs, a fixture and its parity test, and stated counts. A mirror that the change should have updated is a finding. In this repository the examples are the pt/en doc mirrors, `status-cases.fixture.ts` with `fixture_parity_test.go`, and the stated counts.
3. The reviewer checks that each added file belongs in the tree. A stray screenshot, scratch output, reviewer scaffolding or a stray binary is a finding. This follows the cleanup premise in `AGENTS.md`.
4. A new entry point (command, mod, adapter, route, CLI surface) must prove again the cross-cutting invariants of the old ones. Examples: the write lock on the memory files, harness-dir resolution, path handling. A new entry point does not inherit old fixes.

`security-review` (`SKILL.md` and `references/checklist.md`):

5. A diff-inventory gate for a hostile change and for the supply chain.
   - File and content signals (a new checklist row): executable bits, symlinks, submodules, Unicode bidi characters and homoglyphs, unexplained encoded or binary blobs. CI workflow changes: `pull_request_target`, `permissions`, unpinned actions, secrets that reach untrusted code.
   - Dependency signals extend the existing checklist rows "Dependencies" and "Configuration" (`checklist.md`, lines 42 and 43). They add no duplicate row: lockfile drift, new registries or git dependencies, package lifecycle hooks.
   - A security fix is the highest-risk diff. The reviewer looks for changes that ride along and for guards that the diff weakens elsewhere.
   - A security report is proven or disproved before a fix. The proof reproduces the problem with trusted code. It never runs the payload of the reporter.
   - An independent refuter of a report or a fix is a Coordinator dispatch: a second reviewer run, or a configured CLI reviewer. It receives only the finding and the code evidence, never the reasoning of the author. The reviewer never spawns agents.
7. A scanner result is not a pass when the diff changed or narrowed the scanner config, its ignores, its baseline or its threshold. It is also not a pass when the scanner does not gate the merge. A suppression or ignore entry defers the finding. It does not resolve it. To dismiss an alert, the reviewer records a reason that names the boundary that makes it safe. Place: `SKILL.md` step 1, next to "scanner output is evidence", and the Dependencies and Configuration rows.
8. A new checklist row: the diff disables or weakens a security control (tests, logging, TLS or signature checks, auth, limits, scanners); tests or mocks that hide real side effects; delayed or conditional activation (time, host, CI, region).
9. Quick reference: a gap that another enforced layer already blocks is hardening (Low), unless the coverage of that layer is itself incomplete.
10. A new row in "Input and interpretation": untrusted text that reaches a model, or a tool with authority, is a trust boundary. Tool policy is enforced outside the model. Retrieved text and tool output are data.
11. The row "Secrets in the tree": a leaked secret is rotated or revoked outside git before the reviewer dismisses or baselines it. A baseline holds only exact, reviewed fingerprints. It never holds a broad exclusion by path, rule or pattern.
12. `SKILL.md` step 10: a security regression test also has a legitimate control case. A blanket deny must not pass as correct authorization.
13. The row "Authorization": the check fails closed on a missing or partial identity. Cache and index keys are fully scoped.

`regression-testing`:

6. Step 6: "fails before the fix" means an assertion failure on the behavior under test. A compile or import error does not count. When the new test cannot compile against the base, the tester proves red by reverting only the fix in an isolated copy. This uses the isolation rules that step 6 already has.

Each changed skill adds the origin to `metadata.sources` in its frontmatter. The value of `provenance` stays `adapted`.

Release: version 0.30.0.

## 3. Rules

Greps below return no hit in the tree at the time of writing. The ticket author checks again on the base before the first edit.

- R1 Item 1 is in `.skills/code-review/SKILL.md`, in steps 4 or 5. Verifiable: `grep -n -i 'sibling call' .skills/code-review/SKILL.md` has a hit, and the line says that a twin site left broken is Major.
- R2 Item 2 is in the "Contracts" bullet of `.skills/code-review/SKILL.md`, with generic examples. Verifiable: `grep -n -i 'hand-maintained' .skills/code-review/SKILL.md` has a hit in that bullet, and `grep -n -E 'status-cases|fixture_parity' .skills` has no hit.
- R3 Item 3 is in `.skills/code-review/SKILL.md`. Verifiable: `grep -n -i 'belongs in the tree' .skills/code-review/SKILL.md` has a hit, and the line names stray artifacts.
- R4 Item 4 is in `.skills/code-review/SKILL.md`. Verifiable: `grep -n -i 'new entry point' .skills/code-review/SKILL.md` has a hit, and the line names the write lock on the memory files.
- R5 Item 5 is in `.skills/security-review/SKILL.md` and `references/checklist.md`, in the places that section 2 names. Verifiable: `grep -n -E 'pull_request_target|lifecycle|bidi|symlink' .skills/security-review/references/checklist.md` finds each term; the checklist has one new row for file and CI signals and no new row for dependencies; `grep -n -i 'reporter' .skills/security-review/SKILL.md` and `grep -n -i "author's reasoning" .skills/security-review/SKILL.md` each have a hit that names the proof rule and the refuter rule. Note: "author's reasoning" has one hit today in another context; the new hit must name the refuter.
- R6 Item 6 is in `.skills/regression-testing/SKILL.md`, step 6. Verifiable: `grep -n -i 'assertion failure' .skills/regression-testing/SKILL.md` has a hit in step 6, and the step says that a compile or import error does not count.
- R7 Item 7 is in `.skills/security-review/SKILL.md` (step 1) and in the Dependencies and Configuration rows. Verifiable: `grep -n -i 'narrowed' .skills/security-review/SKILL.md .skills/security-review/references/checklist.md` has a hit, and `grep -n -i 'baseline' .skills/security-review/SKILL.md` has a hit that says a suppression defers.
- R8 Item 8 is a checklist row. Verifiable: `grep -n -i 'weaken' .skills/security-review/references/checklist.md` has a hit in a row that names conditional activation.
- R9 Item 9 is in the Quick reference of `.skills/security-review/SKILL.md`. Verifiable: `grep -n -i 'hardening' .skills/security-review/SKILL.md` has a hit in the Quick reference.
- R10 Item 10 is a row in "Input and interpretation". Verifiable: `grep -n -i 'tool policy' .skills/security-review/references/checklist.md` has a hit in that section.
- R11 Item 11 is in the "Secrets in the tree" row. Verifiable: `grep -n -i 'fingerprint' .skills/security-review/references/checklist.md` has a hit in that row, and the row says "rotated or revoked" before baseline.
- R12 Item 12 is in step 10 of `.skills/security-review/SKILL.md`. Verifiable: `grep -n -i 'control case' .skills/security-review/SKILL.md` has a hit in step 10.
- R13 Item 13 is in the "Authorization" row. Verifiable: `grep -n -i 'fail closed' .skills/security-review/references/checklist.md` has a hit in that row.
- R14 Provenance. `metadata.sources` of `code-review` names `akitaonrails/my-skills pr-audit @7b8a5e5`. The same field of `regression-testing` names it too. The same field of `security-review` names `akitaonrails/my-skills pr-audit @7b8a5e5` and `akitaonrails/my-skills security-audit @7b8a5e5`. `provenance` stays `adapted` in all three. Verifiable: `grep -n -E 'my-skills (pr|security)-audit @7b8a5e5' ` on the three `SKILL.md` frontmatters, and `grep -n 'provenance: adapted'` on the same files.
- R15 The change adds no skill, agent or command. The counts stay at 13 agents and 19 + 1 commands. Verifiable: `git diff --name-status main...HEAD` shows no added file under `.agents/` or `.commands/`, and no new `SKILL.md` under `.skills/`; `dh validate` is green.
- R16 Hand edits are limited to: the three skills and `references/checklist.md` in `.skills/`, this PRD and its review file, the ticket records, and the release files in section 4. No file in `dist/` is edited by hand. Verifiable: `go run ./cmd/dh build` regenerates `dist/`, and `git diff --name-only main...HEAD -- dist` lists only `dist/claude-code/dev-harness/skills/{code-review,security-review,regression-testing}/**`, `plugin.json`, `harness-manifest.json` and `dist/*/bin/**`.
- R17 `go run ./cmd/dh validate .` passes with 0 errors. The build runs twice, and the second run changes nothing in `dist/` besides the `generated_at` line of `harness-manifest.json`. Verifiable: run both commands and compare with `git diff`.
- R18 Clean-clone proof, because the PR touches `internal/build/build.go` and `dist/*/bin`: run `git clone . <tmp>`, run `go run ./cmd/dh build` there, and compare `sha256sum` of the binaries with the tree. Verifiable: all sums are equal.
- R19 Length caps, counted as added lines (`git diff --numstat main...HEAD -- <file>`, first column): `code-review/SKILL.md` at most 14; `regression-testing/SKILL.md` at most 4; `security-review/SKILL.md` at most 18; `references/checklist.md` at most 24. Verifiable: the numstat values stay under the caps. The frontmatter line of `metadata.sources` counts as a changed line.
- R20 No new heading. Verifiable: `git diff -U0 main...HEAD -- .skills | grep -c '^+#'` prints 0 for the four touched files.
- R21 Version 0.30.0. `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json`, and the source `ref` (`v0.30.0`) change together. Verifiable: `grep` of the three values shows `0.30.0` or `v0.30.0`.
- R22 The scope in section 5 stays true. Verifiable: the diff has no `/dh:pr-audit` command, no script, no new `dh` subcommand, no `references/ecosystem-checks.md`, and no rule for merge, retarget, attribution or batching.

## 4. Docs

Evidence: `grep -rlE 'code-review|security-review|regression-testing'` over `docs`, the READMEs, `CHANGELOG.md` and `AGENTS.md`, without PRDs and ADRs, lists these files only: `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/plano-produto.html`, `CHANGELOG.md`.

Update:

- `CHANGELOG.md` (release entry 0.30.0). The roadmap is checked against it by `dh validate`.
- `docs/roadmap.html`, `docs/en/roadmap.html` (0.30.0 card; the pt and en cards must agree).
- `.claude-plugin/marketplace.json` (`version` and `ref`).
- `internal/build/build.go` (`defaultVersion`).
- `README.md`, `README.en.md` (line 20: the current-version line becomes 0.30.0 with tag v0.30.0, and the one-sentence summary of the latest release names this PRD). Every earlier release commit (43a280b, 00c0fb4, e67dba1) changed this line.

Not affected:

- `docs/tutorial.html`, `docs/en/tutorial.html` (lines 859 and 860). The command table names `code-review` and `security-review` as the output of `review` and `secure`. It does not describe procedure steps.
- `docs/plano-produto.html`. It names the skills only in the catalog.
- `docs/inicio-rapido.html`, `docs/en/quick-start.html`, `docs/tutoriais/`, `docs/en/tutorials/`. No hit for the three skills.

## 5. Out of scope

From `pr-audit`:

- A `/dh:pr-audit` command for PRs from third-party contributors. It waits for a real case.
- Rules for carry-through by the maintainer, merge, retarget and attribution.
- Batching of several PRs.
- Sandboxed execution of untrusted PRs.
- The "Unreleased" item in the changelog.
- The parts of `pr-audit` that belong to its own project.

From `security-audit`:

- Its 6-phase full audit process.
- Triage and PATCH of GitHub alerts through `gh api`.
- Running scanners or the test gate from the reviewer.
- `references/ecosystem-checks.md`.
- `scripts/security-surface.sh`. There is no script and no new `dh` subcommand.
- Its output template sections.
- Its policy on disclosure and privacy.

Also out: any other change to the three skills.

## 6. Owner decisions

- Owner, 2026-10-16: absorb the ideas into existing skills and import nothing whole. Decision: "follow your recommendations" on the Coordinator's investigation.
- Owner, 2026-10-16: no new skill, agent or command in this change.
- Owner, 2026-10-16: the files in akitaonrails/my-skills (`pr-audit` and `security-audit`) are released for use. Literal copying of their text is allowed. Each changed skill records the origin in `metadata.sources` (R14), and `provenance` stays `adapted`.
- Owner, 2026-10-16: the `security-audit` items (7 to 13, and the refuter rule in item 5) join this PRD.
- Decision (AGENTS.md rule on tags): the change ships as version 0.30.0 (R21).
- Pending decisions: none.

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| Import `pr-audit` and `security-audit` whole | Most of their text belongs to their own project. The kit must stay lean. |
| Add a `/dh:pr-audit` command now | It serves PRs from third-party contributors. No real case exists. It would raise the command count (19 + 1). Deferred until a real case. |
| Write one new "audit" section in each skill | The owner wants a lean kit. Short lines inside the existing steps and rows carry the same checks. |
| Keep the two sources in two PRDs | The `security-audit` items land in the same files as item 5 and share one release. Two PRDs would conflict on `checklist.md`. |

### Risks and pending decisions

- Risk: the skills grow and their trigger text gets noisy · mitigation: R19 caps lines per file and R20 forbids new headings.
- Risk: added text repeats a rule that the skill already has (for example the "do not execute the payload" rule in step 8) · mitigation: the reviewer of the ticket checks for duplicates, and the new line refers to the existing rule.
- Risk: the roadmap pt and en cards disagree · mitigation: `dh validate` checks the drift (R17).
- Risk: the clean-clone proof differs from the tree · mitigation: R18 blocks the commit until the sums are equal.
- Pending decisions: none.

## Executed summary

- **Delivered:** `code-review`, `regression-testing`, `security-review` and its checklist carry the 13 checks from `pr-audit` and `security-audit` as short lines inside existing steps and rows, with no new skill, agent, command or heading, released as 0.30.0.
- **Rules:**
  - R1-R13 Honored: each item's grep hits in the named step or row (QA matrix, `qa-verifier`), and a read of the added text confirms the full meaning of items 7, 8, 10 and 13; adherence read by `product-discovery`: all 13 items adhere.
  - R14 Honored: `metadata.sources` of the three skills names `akitaonrails/my-skills pr-audit @7b8a5e5` (and `security-audit @7b8a5e5` in `security-review`); `provenance: adapted` unchanged.
  - R15 Honored: `git diff --name-status` shows only modified files; 13 agents, 19 + 1 commands, 21 skills (`dh validate`).
  - R16 Honored: hand edits limited to the four skill files, the release files of section 4 (README line amended in, revalidated); the `dist/` diff is limited to the three skills, `plugin.json`, `harness-manifest.json` and the five binaries.
  - R17 Honored: `dh validate` 0 errors, `go test ./...` green; the second build changes only `generated_at`.
  - R18 Honored: clean-clone proof, 5 of 5 binaries with equal sha256.
  - R19 Honored: added lines 5 / 2 / 7 / 7 against caps 14 / 4 / 18 / 24.
  - R20 Honored: no new `#` line in the four files.
  - R21 Honored: `defaultVersion` 0.30.0, marketplace `version` 0.30.0 and `ref` `v0.30.0`.
  - R22 Honored: no command, script, `dh` subcommand, `ecosystem-checks.md` or merge/retarget/attribution/batching rule.
- **Tickets:** T-2301 (backend, harness-maintainer) approved. QA by `qa-verifier`: verified, one minor note (item 7's scanner wording sits in the Configuration row and step 1, not in the Dependencies row; accepted). Code review by `cli:codex/gpt-6-luna`: no finding on the text; its one Major (R15-R18 evidence absent from the review bundle) was resolved with the QA matrix, with no edit.
- **Docs:** `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `README.md`, `README.en.md`, `.claude-plugin/marketplace.json`, `internal/build/build.go`.
- **Out of scope:** a `/dh:pr-audit` command for third-party PRs, maintainer merge/retarget/attribution rules, multi-PR batching, sandboxed execution of untrusted PRs, the full six-phase audit, GitHub alert triage, running scanners from the reviewer, `ecosystem-checks.md` and `security-surface.sh`.
