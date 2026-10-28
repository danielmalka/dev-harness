# Review: PRD-029

## Round 1

| Field | Value |
|---|---|
| Document | docs/prd/PRD-029-independent-review-and-drift-report.md |
| Source | Owner decisions (2026-10-28); discovery UND-009 (research in `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-28 |

## Findings

- **P1** [applied] [conflict] (blocking) Bundle definition in R1 and R3 contradicts existing prompt assembly and dispatch contracts.
  - Location: section 3.1, rules R1 and R3
  - Evidence: R1 and R3 define the bundle as the full text passed to a reviewer and restrict it to five path/metadata items, which would strip role instructions, review skill, anti-delegation clause, and verdict contract required by `.agents/coordinator.md` lines 66 and 94.
  - Suggestion: Clarify that R1 and R3 apply strictly to the context bundle payload within prompt assembly.
  - Resolution: applied by the author, round 1 (Section 2 Part A, R1 and R3).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) R26 excludes `fixed` from settled findings, making clean rounds impossible under R30 after a bug fix.
  - Location: section 3.3 (R26) and section 3.4 (R30)
  - Evidence: R26 defines settled findings as only `intentional`, `noise`, or `escalated`. Under R30 ("No blocking finding is unsettled"), a fixed defect remains technically unsettled, permanently blocking convergence.
  - Suggestion: Define R30's clean condition as "No blocking finding is `open` or `reopened`", or specify that confirmed `fixed` findings are settled.
  - Resolution: applied by the author, round 1 (R21, R26, R30).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] (blocking) Finding ID assignment in R19 does not explicitly preserve stable IDs for unresolved findings across rounds.
  - Location: section 3.3, rule R19
  - Evidence: R19 states findings are assigned `F<round>-<n>` in each round, without stating whether open findings retain earlier IDs or are renumbered, contradicting Problem 1.C.
  - Suggestion: Explicitly state that `F<round>-<n>` applies to new findings and carried-over findings retain their original IDs.
  - Resolution: applied by the author, round 1 (R19).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [conflict] (blocking) Rule R39 is not marked as provisional despite Decision D2(b) presenting the comparison role as an open owner decision.
  - Location: section 3 preamble, section 3.5 (R39), and Appendix (D2)
  - Evidence: Appendix D2(b) presents three options for who runs the drift comparison, but Section 3 preamble asserts R39 holds in full and omits R39 from the provisional rules.
  - Suggestion: List R39 as provisional and waiting for Decision D2(b) in Section 3 preamble and Appendix pending decisions.
  - Resolution: applied by the author, round 1 (Section 3 preamble and Appendix pending decisions).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] (Minor) Procedural and descriptive sentence length limit violations under ASD-STE100 PRD writing rules.
  - Location: section 1 (Problem C), section 3.1 (R3, R8), section 3.7 (R55)
  - Evidence: Sentences in Problem C (50 words), R3 (34 words), R8 (37 words), and R55 (34 words) exceed the 20-word procedural / 25-word descriptive limits.
  - Suggestion: Break long sentences into single-instruction statements under the word limits.
  - Resolution: applied by the author, round 1 (Problem 1.C, R3, R8, R55).
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Faithfulness to owner decisions Q1 through Q6 and Coordinator assumptions.
- Coordinator reading of Q3 kept visible as open decision D1 and not attributed to owner.
- Structuring of open decisions D1, D2, and D3 with options, consequences, and recommendations.
- Provisional markers on R30-R35, R47, and R51.
- Rule testability across R1 through R60.
- All citations of dh files and line numbers verified against `source3-cited-dh-passages.md` and found 100% accurate.
- Absence of scope excess (no new agents, commands, or runtime dependencies).
- Comprehension block correctly scoped to `document-validator` only.
- Mandatory coverage of all six categories verified.

## Round 2

| Field | Value |
|---|---|
| Document | docs/prd/PRD-029-independent-review-and-drift-report.md |
| Source | Owner decisions (2026-10-28); discovery UND-009 (research in `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-28 |

Round 2 reported 0 blocking findings. Raw replies of both rounds: `~/.harness/projects/dev-harness/tasks/UND-009/evidence/validation-r1/agy-reply.raw.md` and `validation-r2/agy-reply.raw.md`. The `## Comprehension` block is not in this report: the block is proposed by this PRD and was not yet in force when the rounds ran.

## Findings

- **P1** [applied] [conflict] Bundle definition in R1 and R3 contradicts existing prompt assembly and dispatch contracts.
  - Location: section 3.1, rules R1 and R3
  - Resolution: The author amended Section 2 Part A, R1, and R3 to separate the context bundle payload (restricted by R1-R9) from prompt assembly (retaining role body, skill text, anti-delegation clause, verdict contract, and reading-is-allowed block per coordinator line 94), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] R26 excludes `fixed` from settled findings, making clean rounds impossible under R30 after a bug fix.
  - Location: section 3.3 (R26) and section 3.4 (R30)
  - Resolution: The author updated R30 to define a clean round as having no Critical or Major finding that is `open`, `reopened` or `escalated`, allowing confirmed `fixed` findings to cleanly converge, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] Finding ID assignment in R19 does not explicitly preserve stable IDs for unresolved findings across rounds.
  - Location: section 3.3, rule R19
  - Resolution: The author amended R19 to state: "A finding carried into a later round keeps its original id. A finding never gets a second id", round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [conflict] Rule R39 is not marked as provisional despite Decision D2(b) presenting the comparison role as an open owner decision.
  - Location: section 3 preamble, section 3.5 (R39), and Appendix (D2)
  - Resolution: The author added R39 to the list of provisional rules in Section 3 preamble and to the pending decisions block in the Appendix, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Procedural and descriptive sentence length limit violations under ASD-STE100 PRD writing rules.
  - Location: section 1 (Problem C), section 3.1 (R3, R8), section 3.7 (R55)
  - Resolution: The author broke the overlong sentences in Problem 1.C, R3, R8, and R55 into shorter statements under 20 and 25 words respectively, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] (Minor) Procedural and descriptive sentence length limit violations under ASD-STE100 PRD writing rules in newly edited rules.
  - Location: section 3.2 (R14), section 3.5 (R43), section 3.7 (R57, R59)
  - Evidence: Several descriptive and procedural sentences exceed the ASD-STE100 length limits defined in the PRD writing rules (procedural <= 20 words, descriptive <= 25 words): R14 sentence 1 has 30 words; R43 sentence 1 has 38 words; R57 sentence 1 has 32 words; R59 sentence 2 has 26 words.
  - Suggestion: Break these sentences into shorter statements under 20 and 25 words respectively.
  - Resolution: reported as `pending` in the round 2 reply; fixed afterwards on the Coordinator's record (style only, no round used).
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Faithfulness to owner decisions Q1 through Q6 and Coordinator assumptions.
- Coordinator reading of Q3 kept visible as open decision D1 and not attributed to owner.
- Structuring of open decisions D1, D2, and D3 with options, consequences, and recommendations.
- Provisional markers on R30-R35, R39, R47, and R51.
- Rule testability across R1 through R60.
- Settled finding states, churn handling, clean round conditions, and evaluated state isolation.
- All citations of dh files and line numbers verified against `source3-cited-dh-passages.md` and confirmed 100% accurate.
- Absence of scope excess (no new agents, commands, or runtime dependencies).
- Comprehension block correctly scoped to `document-validator` only.
- Mandatory coverage of all six categories verified.
