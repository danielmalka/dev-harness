# Review · PRD-008 · round 1

- Panel (reviewers.document; the PRD touches security): claude (document-validator, opus requested, effective unverified) — changes required. cli:grok/grok-4.7 and cli:codex/gpt-5.6-luna — not-run: the Coordinator wrote to the tree (AGENTS.md migration and .harness records) while their calls were in flight, which invalidates the freeze comparison; grok was stopped mid-call and codex never started. Post-stop freeze diff shows only the Coordinator's four known writes. No correction round is spent for them; both run in round 2.
- Date: 2026-09-29.

## Verdict
- changes required
- Blocking findings: 5

## Findings
- P1 · conflict · blocking · §2, Apêndice (D1, first risk). The residual is understated as ".git/hooks and .git/config". The shared common dir also holds refs/, packed-refs, objects/, info/exclude (writable from inside the worktree without absolute paths, e.g. `git update-ref`, appending to info/exclude which blinds the live-tree path list); and the live tree stays reachable by absolute path (the worktree's .git file and `git worktree list` reveal it), so the RISK-002 round-9 vector (self-ignored .gitignore + _test.go in the live tree) remains outside R8. Suggestion: state the full residual, say which items R8 detects and which not, and either extend R8 (for-each-ref, packed-refs, info/ of the common dir, `--ignored` on the live-tree list) or record each undetected item as an accepted residual under D1. Reported by: claude
- P2 · conflict · blocking · §3 R14, R15. The eval writes inside the worktree, where R3 guarantees .harness/ does not exist — it passes by construction and never exercises the RISK-002 threat; the ".git/hooks" action is unspecified (relative write fails with ENOTDIR); no case writes the live tree by absolute path. Suggestion: add fixture actions (i) absolute write to the live .harness/project.yaml, (ii) write to $(git rev-parse --git-common-dir)/hooks/<name>, (iii) plant .gitignore + _test.go in the live tree, with the expected verdict for each. Reported by: claude
- P3 · ambiguity · blocking · §3 R7 (with R12, R14). Whether the in-worktree comparison includes ignored paths is undefined: reading A accepts planted ignored files, reading B discards every call that ran commands.test (R12). Suggestion: name the exact comparison and how test artifacts are treated; state verdict and not-run status per fixture action. Reported by: claude
- P4 · conflict · blocking · §2 ("isolamento verificado"), R10, R4. Prevention is claimed where only detection exists; the texts R10 amends justify planners by "must be unable to write" and "the freeze is a detector, not a sandbox", and R10 does not say how that justification changes; R4 is hygiene, not a control. Suggestion: reword §2; have R10 require the three texts to state the detection-only distinction; mark R4 as hygiene. Reported by: claude
- P5 · ambiguity · blocking · §3 R2, R5, R6, R13, R8. "Kit path" and "orphan" are undefined; parallel dispatch B could prune A's live worktree or an owner's worktree; `git worktree lock` unused; R8 may flag the kit's own worktree. Suggestion: name prefix `dh-cli-<call-id>` plus `git worktree lock --reason <call-id>`; orphan = prefixed and unlocked by a live call; never touch unprefixed worktrees; R8 excludes the call's own worktree or runs after R5. Reported by: claude
- P6 · organization · non-blocking · D2 (codex half), D6 settled by R16 and .commands/plan-loop.md; D1 vs R15 provisional; D5 is a probe result. Reported by: claude
- P7 · organization · non-blocking · H1 confirmed by the Coordinator's check — move to Fatos. Reported by: claude

## Not raised
- Problem statement, goal traceability, per-binary read-only status, F1/F4/F5, R10 count, R11 contract preservation, R13 cap, R14 cost rule, R3/R9 fail-closed, R16 boundaries, Docs list and format: all traced and consistent.

---

# Review · PRD-008 · round 2 (final)

- Panel: claude — changes required (P1–P7 applied; new P8, P9). cli:grok/grok-4.7 (detached, freeze clean; verdict recovered, line without "- ") — changes required. cli:codex/gpt-5.6-luna (freeze clean; verdict recovered) — changes required. Raw CLI outputs: .harness/tasks/PRD-008/evidence/{grok,codex}-r2.out. Round cap reached (2 of 2).

## Verdict
- changes required

## Findings (merged; equivalent points share one entry)
### Live .git metadata outside the snapshot (.git/index, logs/, worktrees/) — skip-worktree and global fsmonitor can hide tracked-file edits
- Reported by: claude — "the CLI edits a clean tracked file ... then runs `git -C <live> update-index --skip-worktree` ... A `core.fsmonitor = <command>` in `~/.gitconfig` is executed by the post-call `git status`."
- Reported by: cli:codex/gpt-5.6-luna — "Writable live metadata such as `.git/index`, `.git/logs/`, `.git/worktrees/` ... is not listed as detected or accepted residual."
### Owner acceptance of residuals (a), (b), (e) attributed without a record
- Reported by: claude
### R2 "system temp dir" vs R9 "kit temp dir" undefined as the same root
- Reported by: cli:grok/grok-4.7; cli:codex/gpt-5.6-luna
### Lock marker at the clone root breaks R3's match (should live in `<clone>/.git/`)
- Reported by: cli:grok/grok-4.7
### R3's match compares against R5's --ignored photograph (ambiguous; must use only the judged set)
- Reported by: cli:grok/grok-4.7
### agy without the cwd probe: stays on the live tree for review, does not plan (R1 vs §2 wording)
- Reported by: cli:grok/grok-4.7
### §1 claims detection of anything outside the clone; only the live tree is detected
- Reported by: cli:grok/grok-4.7
### Residual (b) "git push undetected" vs R5 detecting live-ref moves
- Reported by: cli:grok/grok-4.7
### ADR-002 §3 "never writes to disk / any change discards" not covered by R12/R14 amendments
- Reported by: cli:grok/grok-4.7
### A planted ignored _test.go in the clone can influence the reviewer's own `commands.test` run
- Reported by: cli:codex/gpt-5.6-luna
### (non-blocking) D1 marked "bloqueia R6" though R6 already fixes "note only"
- Reported by: cli:grok/grok-4.7

## Owner disposition (2026-09-29)
- The owner accepted residuals (a) writes outside the project, (b) surviving processes / network / push to non-live remotes, (e) reading and returning secrets, until the OS-sandbox batch; and the Coordinator's recommendation to close the index/fsmonitor hole. All other findings are mechanical and are applied in one final correction, checked by the Coordinator (no third validation round).
