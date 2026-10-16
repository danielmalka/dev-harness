# Review: PRD-024

| Field | Value |
|---|---|
| Document | docs/prd/PRD-024-public-readiness.md |
| Source | Owner-approved scope, chat 2026-10-16 (Coordinator dispatch); audit UND-006; docs/discovery/renova.md §7; AGENTS.md |
| Round | 2 of 2 |
| Verdict | changes required |
| Updated | 2026-10-16 |

## Findings

- **P1** [applied] [conflict] R8 could never pass because the PRD held the owner's path.
  - Resolution: round 2. The literal appears only in R8 (§3:97), with a `':!docs/prd/PRD-024*'` pathspec that also covers this review file. §1, §2 and §6 now use "the owner's home/checkout path".
  - Reported by: claude
- **P2** [applied] [conflict] The count rule misstated checkDocumentedCounts, and R7 held on the base.
  - Resolution: round 2. §2:62 allows at most one anchored line with two or more nouns (four anchors). The R7 grep is widened. §3:85 names R7, R11, R17, R18, R19 and R21 as regression guards.
  - Reported by: claude
- **P3** [applied] [weak criterion] The R4 greps passed on the base, and R1 counted only headings.
  - Resolution: round 2. The §2 table fixes 12 en/pt headings. R1 checks the order. R4 is scoped per section and includes per-hook/por hook. R3 requires optional/opcional.
  - Reported by: claude
- **P4** [applied] [gap] The dead-reference list was incomplete, and R14 was ambiguous on RISKS.md.
  - Resolution: round 2. §1:26-27 lists the exact lines in both files (they match the base). §2:73 keeps 760 and 839 and rewords the RISK-002 Relations. R14 has CLAUDE.md empty and RISKS.md lines named (see P9).
  - Reported by: claude
- **P5** [applied] [ambiguity] (non-blocking) Current README blocks had no target section.
  - Resolution: round 2. §2:54-60 maps Language, clone-load, sources, the 0.8.0 note, Global mode and line 20.
  - Reported by: claude
- **P6** [applied] [weak criterion] (non-blocking) Card titles (a)-(c) were not fixed.
  - Resolution: round 2. §2:67-70 gives the exact en/pt `<h3>` texts, and R12 checks them.
  - Reported by: claude
- **P7** [applied] [ambiguity] (non-blocking) The R10 `claude` grep forbade naming Claude Code.
  - Resolution: round 2. The pattern is now `co-authored-by|generated with`, and Claude Code and `.claude-plugin` are allowed.
  - Reported by: claude
- **P8** [applied in part] [organization] (non-blocking) Sentence length and passive voice.
  - Resolution: round 2. §1:31 and card (d) are split. §2:62 still has one sentence of about 27 words (Minor).
  - Reported by: claude
- **P9** [accepted by owner] [conflict] The card (d) text puts `RISKS.md` in both roadmaps, so R14's "exactly 2 lines per file" fails.
  - Location: §2:70 against §3 R14
  - Evidence: Card (d) says the Coordinator "still owns the writes to `MEMORY.md`, `EPOCHAL.md` and `RISKS.md`" (also §6:148). The grep then prints 760, 839 and the card line.
  - Suggestion: allow exactly 3 lines per file (760, 839 and the card (d) line). State that the 0.30.1 Delivered card writes neither `CLAUDE.md` nor `RISKS.md`.
  - Reported by: claude
- **P10** [accepted by owner] [ambiguity] (non-blocking) R4's "7 table rows that start with `|`" does not say whether the header and separator lines count, and a literal count gives 9.
  - Location: §3 R4, Writes bullet
  - Suggestion: "7 body rows (after the header and the separator)".
  - Reported by: claude

## Not raised

- Base RISKS.md lines per roadmap are 760, 839 and the two RISK-002 Relations lines, so exactly 2 would remain without the new card. CLAUDE.md lines are as listed in §1:26.
- The Requirements row (Claude Code; Git) matches base README.en.md:28. Go moving to For contributors is consistent with the prebuilt binaries.
- The R1, R5 and R6 heading checks hold against the fixed headings ("contribut" / "contribu").
- Card (d) still matches the owner decision, and its implementation stays out of scope (§5:138).
- All round 1 Not raised items still hold: source mapping, owner decisions, audit and release facts, out of scope, repo rules, shape.

## Owner disposition (2026-10-16)

- Validation cap reached with P9 open. The owner chose to apply the validator's correction for P9 (R14: exactly 3 `RISKS.md` lines per roadmap, the 0.30.1 Delivered card cites neither `CLAUDE.md` nor `RISKS.md`) and P10 (7 body rows), with no further round.

## In-flight amendment, 1 of 1 (2026-10-16)

- Review status: changes required (1 blocking, goes to the owner). Scope: R8, R15, R22, §4:120, §6:151-152.
- **P11** [accepted by owner] [conflict] The R15 marketplace clause assumes a stale count that is not stale.
  - Location: §3 R15, §4:120, §6:152
  - Evidence: marketplace.json:8 says "coordinator, 12 specialist agents". AGENTS.md says 13 agents including the coordinator, so the text is true. R15 forbids the string and allows a wrong "13 specialist agents". checkDocumentedCounts does not scan marketplace.json.
  - Suggestion: drop the clause, or require "13 agents in total or coordinator plus 12 specialists".
  - Reported by: claude
- **P12** [accepted by owner] [weak criterion] (non-blocking) The R8 checks do not cover HTML files, and the "scans the tracked files" note is inaccurate.
  - Location: §3 R8
  - Evidence: validate.go:579 skips `.html` and :584 skips `evals/baselines` and `evals/results`. The git grep covers only baselines. The roadmaps are edited in this PR.
  - Suggestion: `git grep -l '/ho[m]e/ma[l]ka'` over the whole tree prints nothing.
  - Reported by: claude
- **P13** [accepted by owner] [ambiguity] (non-blocking) R22 can remove the disable step that R4 requires, and it does not say whether the CLI form proves the slash form.
  - Location: §3 R22, R4
  - Suggestion: fall back to `/plugin uninstall` in Turn-off, and state that a command is verified in the form the README prints.
  - Reported by: claude
- **P14** [accepted by owner] [organization] (non-blocking) §6:152 is not in the "Owner, <date>:" decision format.
  - Location: §6:152
  - Reported by: claude

### Not raised

- PRD-024 and its review file contain no private-path substring (the `checkPrivatePaths` pattern), so checkPrivatePaths passes on them.
- R22 numbering is unique (R1-R21 are contiguous). R22 is checkable through the QA report.
- R8 keeps the parse loop and the check for 28 `<repo>` files. The pathspec exclusion is removed.
- The P9 and P10 corrections are present in R14 and R4. Nothing outside the declared amendment scope changed.

### Owner disposition (2026-10-16)

- Single round ended with P11 open. The owner chose to remove the marketplace description clause (P11) and apply the validator's corrections for P12, P13 and P14, with no further round.
