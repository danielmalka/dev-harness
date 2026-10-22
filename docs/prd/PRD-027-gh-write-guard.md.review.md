# Review: PRD-027

| Field | Value |
|---|---|
| Document | docs/prd/PRD-027-gh-write-guard.md |
| Source | Owner decisions D1-D12 (2026-10-21), discovery DISC-007, measurement M2 |
| Round | 1 of 2 |
| Reviewer | cli:agy/gemini-3.8-flash-high |
| Verdict | approved |
| Updated | 2026-10-21 |

## Findings

(No open findings.)

## Not raised

- Gaps: The three defects of batch 0.32.1 (`gh` write guard in subagent Bash, opencode v2.0.32 recipe, `pgrep`/`pkill -f` self-match gotcha) map from DISC-007, the owner's chat and M2 to problem descriptions G1-G3 and rules R1-R14.
- Gaps: Decision D1 (kit subagents `dh:*` denied except `dh:release-manager`; fail-closed on an unknown `agentId`; main session and non-kit subagents unaffected) is covered in sections 2.1, D1, R4, R5 and R7. The cleanup scope of AGENTS.md ("Premissa de limpeza") is covered in R15 and section 4.
- Ambiguities: Command-segment splitting at shell separators, skipping of environment assignments and wrappers (`env`, `command`, `exec`, `sudo`), and argument handling are defined in 2.3 and R2. The pure-function contract of `hooks/lib/` (R1) fixes input, output and zero I/O. The six execution contexts of the engine test (R5) do not overlap. Fail-open on a guard exception and fail-closed on an unknown `agentId` are kept separate (R7).
- Conflicts: The ADR-005 amendment (R9) resolves the narrow exception for `gh` write verbs and `git push` against ADR-005's general rejection of command analysis. The `AGENTS.md` "Trava de `gh`" line matches 2.1, D1 and R5. Invariant counts (13 agents, 19 + 1 commands, 14 templates) match the inventory (R16). The version 0.32.1 synchronization (R17) and the clean-clone proof for Go changes agree with AGENTS.md. The critical-path merge stop (R18) matches the `REVIEW.md` entries touched by T-2701 and T-2702.
- Excesses: Out-of-scope items (OS sandbox, curl and MCP blocking, restricting `qa-verifier` tools, `dh review run`, new evals or probes, `reviewers:` changes, non-Claude Code runtimes) stay in section 5. Known bypasses are documented in 2.4 and R9 without added scope.
- Weak criteria: All rules R1-R19 state a deterministic, observable check. M1 (R8) has explicit pass and fail conditions under bypass permission mode.
- Organization: The document follows `templates/en/PRD.md`. Rules R1-R19 are numbered in sequence. Suggested tickets T-2701 to T-2704 have correct lanes, write sets and dependencies (ADR-003, PRD-007).
- PRD writing style: conforms to the ASD-STE100 rules in `.skills/requirements-discovery/SKILL.md`.

## Round log

- Round 1 of 2, 2026-10-21: approved, 0 blocking findings. Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Coordinator clarification after round 1, 2026-10-21 (not revalidated, within D3): section 2.3 and R2 also skip global options between the binary and its subcommand (`gh -R`/`--repo`, `git -C`, `-c`, `--opt=value`); Status set to approved.
- Amendment revalidation (single round), 2026-10-21: approved, 0 blocking findings, no open findings. Amended after security review round 1: sections 2.1-2.4, R1, R2, R3, R5, R6, R9, D1, D2, D3, D6 and the appendix risks (read allow-list for `gh`, two-pass matching, Bash/Monitor/PowerShell). Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Amendment 2 revalidation (single round), 2026-10-21: approved, 0 blocking findings, no open findings. Amended after security review round 2: sections 2.2, 2.3, 2.4, R2, R9(a), D2, D3 and the appendix false-positive risk (two real passes, check at every `gh`/`git` word with a search/print exception, new `git` denies, method-override header). Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Amendment 3 revalidation (single round), 2026-10-21: approved, 0 blocking findings (lexer with a context stack, recursion and size limits, narrower search exception, PowerShell no-string deny). Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean. Superseded by amendment 4 before implementation.
- Amendment 4 revalidation (single round), 2026-10-21: approved, 0 blocking findings, no open findings. Owner decision: the shell parser is replaced by a short rule (trigger words, one plain gh read or one plain grep, no shell metacharacters). Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Amendment 5 revalidation (single round), 2026-10-21: approved, 0 blocking findings, no open findings. Amended after security review round 4: gh read checked in order, `auth status` flag list, word boundaries, `git-<verb>` forms, PowerShell without exceptions, general false-positive rule. Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Amendment 6 revalidation (single round), 2026-10-22: approved, 0 blocking findings, no open findings. Amended after security review round 5: lower-case trigger text, `\` handling per shell (PowerShell path separator), `git` ends before `-` (dash list removed), `gh api` without `--hostname` and without `://` tokens, upper-case gh read as an accepted false positive. Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Coordinator correction after amendment 6, 2026-10-22 (not revalidated, wording only): the Windows path example in 2.3 and R2 uses a `<drive>` placeholder in place of a real drive letter, because `dh validate` rejects a drive-letter path as a private path; the rule and the case are unchanged.
- Amendment 7 revalidation (single round), 2026-10-22: approved, 0 blocking findings, no open findings. Amended after security review round 6 (owner choice "a": one extra round): two Bash trigger texts (every backslash removed, or `\\` mapped to `/`), triggered when either triggers; gh read exception only when no text triggered on `git`; R2 backtick cases; 2.4 records the `GH_ENTERPRISE_TOKEN` hypothesis. Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
- Amendment 8 revalidation (single round), 2026-10-22: approved, 0 blocking findings, no open findings. Amended after security review round 7 (owner choice "a"): 2.4 names the backslash-newline form inside a command name as an accepted bypass; header updated. Reviewer: cli:agy/gemini-3.8-flash-high. Workspace check: clean.
