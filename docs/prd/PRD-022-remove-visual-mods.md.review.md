# Review: PRD-022

| Field | Value |
|---|---|
| Document | docs/prd/PRD-022-remove-visual-mods.md |
| Source | Owner decisions 2026-10-14 (Coordinator dispatch, T-2201); pre-change hooks tree (dist/claude-code/dev-harness/hooks) |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-14 |

## Findings

- **P1** [pending] [weak criterion] R3's grep cannot pass, because `--exclude-dir` matches base names only, so docs/prd and docs/adr are not excluded and PRD-019 and PRD-022 hit. It also misses `mods-types.ts`, `lib/status.ts`, `lib/zen.ts` and `mods.engine.test.ts`.
  - Location: §3 R3
  - Suggestion: `git ls-files adapters dist | grep -E 'mods-|lib/(status|status-line|zen|zen-band|menu-pages|modes|settings)\.'` empty; `git grep` with `':!docs/adr' ':!docs/prd' ':!CHANGELOG.md'` pathspecs and the wider pattern, with hits only in roadmap release cards
  - Reported by: claude
- **P2** [pending] [weak criterion] No rule proves the 60%/80% toast removal (`$.ui.toast` in panel.ts). R1 matches the `on('ui.status')` engine stubs in the kept test files.
  - Location: §2 toasts bullet, §3 R1
  - Suggestion: R1 grep `'\$\.ui\.(status|toast)|PromptHint|zen_progress|ZEN_PROMPT|/mods'` with `--exclude='*.test.ts'`
  - Reported by: claude
- **P3** [pending] [weak criterion] `dh validate` does not check the "19 + N" mod-command count (doc_counts.go only compares "N comandos" with the files on disk), so R4 cannot fail.
  - Location: §3 R4
  - Suggestion: grep for `19 \+ 2|2 comandos registrados|2 commands registered|/mods` in AGENTS.md, both READMEs and both tutorials has no hit, plus `dh validate` green
  - Reported by: claude
- **P4** [pending] [ambiguity] (non-blocking) The permission-mode and agent-count bullets name no file. The `running` map in panel.ts could stay as dead state.
  - Location: §2
  - Reported by: claude
- **P5** [pending] [organization] (non-blocking) Status says "approved on 2026-10-14" before validation and before the owner has seen the PRD.
  - Location: header
  - Reported by: claude

## Not raised

- Every owner decision of 2026-10-14 maps to §2, §3, §5 or §6. No excess and no internal conflict.
- The R1 and R3 patterns hit none of the kept files (runtime-guard.ts, snapshot-writer.ts, status-drift.ts, dashboard-command.ts, lib/gate.ts), apart from the test stubs in P2.
- R2, R5, R6 and R7 can be checked; the clean-clone proof applies correctly.
- Style meets the PRD writing rules.

## Round 2 (2026-10-14)

- Review status: changes required (1 blocking, new). P1-P5 applied. Round cap reached; P6 goes to the owner.
- P1 [applied] R3 uses `git ls-files` and `git grep` with `:!` pathspecs and the wider pattern; current-tree hits are only in the roadmaps.
- P2 [applied] R1 covers `$.ui.toast` and excludes `*.test.ts`.
- P3 [applied] R4 has a grep check.
- P4 [applied] §2 names the files and the `running` map.
- P5 [applied] Status is "draft".
- **P6** [settled: owner chose option A on 2026-10-14; `/mods` removed from R4's pattern] [weak criterion] R4's `/mods` term hits the correct release-note text in `README.md:20` and `README.en.md:20`, so R4 cannot pass. Suggestion: remove `/mods` from R4's pattern; R1 covers it in the code.
  - Reported by: claude
