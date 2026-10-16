# Review: PRD-025

| Field | Value |
|---|---|
| Document | docs/prd/PRD-025-coordinator-orchestrates.md |
| Source | Owner decisions (2026-10-16) and roadmap card "O Coordenador orquestra, não executa" (docs/roadmap.html ~L1726) |
| Round | 2 of 2 |
| Reviewer | cli:agy/gemini-3.8-flash-high |
| Verdict | approved |
| Updated | 2026-10-16 |

## Findings

- **P1** [applied] [conflict] Delegation of consolidation output and local.yaml in Solution contradicts memory ownership rules, R2, line 54, and R7
  - Location: section 2, line 52
  - Resolution: Author updated Section 2 Table, assigning candidate-*.md to docs-guide, while keeping consolidation and local.yaml writing with the Coordinator (dispatching only git check-ignore to qa-verifier per line 48), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] R5 allowlist reusing protectedHit excludes other harness files and contradicts R7
  - Location: section 3, line 76 (Rule R5)
  - Resolution: Author replaced protectedHit reuse with dedicated check covering any path inside the harness dir (.harness/ or <home>/projects/<name>/), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] Unspecified warning mechanism and Write new-file semantics produce incompatible implementations
  - Location: section 3, line 76 (Rule R5) & Appendix line 156
  - Resolution: Author defined Write triggers (>5 lines or creation of new file) and clarified the test contract as an observable pure function in hooks/lib/ returning warning text or undefined, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [gap] Section 4 Docs omits .agents/docs-guide.md and .skills/document-review/SKILL.md
  - Location: section 4, lines 90-104 & section 3, line 79 (Rule R8)
  - Resolution: Author added .agents/docs-guide.md and .skills/document-review/SKILL.md to the Section 4 update list, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Suggested ticket T-2503 uses placeholder lane teste (...)
  - Location: section Appendix, line 137
  - Resolution: Author replaced placeholder ellipsis with concrete lane teste (unitário e integração), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] Descriptive sentences in R5 and R10 exceed ASD-STE100 word limits
  - Location: section 3, lines 76 and 81
  - Resolution: Author broke the descriptive sentences in R5 and R10 into concise sentences under 25 words, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- R1: Coordinator does not read repository, grep, build, validate or run checks itself; dispatches repo-scout and qa-verifier.
- R2: Coordinator keeps own records (MEMORY.md, EPOCHAL.md, RISKS.md) and ticket status flip with Concluído em / Done on.
- R3: 5 changed lines limit proposal for direct edits in existing files with no new file.
- R4 & R11: CLI review mechanics dispatched to qa-verifier or builder lane infra, with freeze rule binding the delegate and Coordinator merging verdicts.
- R5, R6 & R7: Mod warn rule specifications: warn-only, never denying, failing open, ignoring specialists (agentId present), excluding Bash, and preserving memory protection deny rules.
- R8: docs-guide writes PRD Executed summary and persists <document>.review.md.
- R9: Definition of "Judge" restricted to reading reports and cited evidence, with spot-checks converted to dispatches.
- R10: One-line scope rule addition to AGENTS.md prevailing over inline trivial work.
- R12: Invariant counts of 13 agents and 19 + 1 commands unchanged.
- R13: Version 0.31.0 in internal/build/build.go, marketplace.json, and git ref.
- R14: Invariant that dist/ is regenerated exclusively by go run ./cmd/dh build.
- R15 & R16: Deterministic verification commands (dh validate, go test, clean-clone proof due to internal/build/build.go).
- R17: Prohibition of private machine paths verified against checkPrivatePaths.
- R18: External CLI calls executed as single standalone commands starting with binary name without wrappers or chaining.

## Post-validation amendment (2026-10-16)

- R18 changed after round 2 (approved). A dispatched qa-verifier measured that a subagent's working directory does not persist between Bash calls, so the three-command `cd` sequence for clone calls was replaced by the binary's own working-directory flag, with `cd "$CLONE" && <binary>` as the only exception for a binary without one. The agy recipe gained the rule to treat a complete reply on disk as the result and stop a process that outlives `--print-timeout`.
- Not revalidated: the two-round cap is used. Owner disposition requested.
