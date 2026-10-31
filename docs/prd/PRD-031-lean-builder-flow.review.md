# Review: PRD-031

## Round 1

| Field | Value |
|---|---|
| Document | docs/prd/PRD-031-lean-builder-flow.md |
| Source | tasks/PRD-031/SOURCE.md (harness dir) |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-30 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 5. Raw reply: `tasks/PRD-031/evidence/document-r1/agy-output.md` (harness dir).

## Comprehension
- What the document proposes: The document establishes a leaner builder workflow by instituting test-first development as the default for code implementation slices and introducing three heuristics in regression testing to prevent empty or dishonest tests. It standardizes commit message bodies to document reason, code changes, and verification checks, while constraining harness procedure steps to single-sentence imperatives with optional one-sentence justifications. It imposes a size ceiling requiring four builder-phase files to shrink to at most 40% of their current size without dropping load-bearing safeguards tracked in a rule inventory. Finally, it outlines seven discrete review-stage simplifications, including single clean round completion and consolidated finding states, submitted for individual owner approval.
- How the current system works, as the document states it: The builder phase is governed by four files totaling 830 lines and over 107 kilobytes of procedural instructions. Test-first execution currently applies strictly to bug regression testing, leaving general feature slices without a red-first mandate. Review cycles carry substantial token and bookkeeping overhead from multi-stage machinery, including six finding states, structured finding identifiers, vendor rotation cursors, and mandatory two-consecutive-clean-round confirmations. Harness authoring documentation lacks constraints on procedure step length, and commit message bodies have no enforced structural standards. In-flight doubts permit up to three fresh-context cycles and rely on complex multi-way finding classifications.
- Sections not paraphrased: none

## Findings

- **P1** [pending] [gap] (blocking) Test-first default for implementation slices dropped from regression-testing
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R1 and R3
  - Evidence: SOURCE.md approved recommendation 1 requires landing test-first in incremental-implementation and regression-testing; R3 explicitly keeps red-first restricted to bug tests.
  - Suggestion: Update R3 to land test-first requirements and exemptions for slices in regression-testing alongside the three heuristics.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [pending] [conflict] (blocking) Five finding states in source asserted as six finding states in document
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 1, line 31 and section 3, R13
  - Evidence: SOURCE.md section "Source findings", bullet 4 states five finding states; PRD Section 1 and R13 assert six finding states and six exits.
  - Suggestion: Align Section 1 and R13 with SOURCE.md by stating that finding states shrink from five to three and correcting the exit count to five.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [pending] [conflict] (blocking) In-flight doubt remaining shape falsely attributed to R15
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R2
  - Evidence: R2 claims R15 sets the remaining shape of the In-flight doubt, but R15 defines review re-entry churn replacement.
  - Suggestion: Remove the reference to R15 from R2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [weak criterion] (blocking) Incomplete grep check for R1 exemption list and report requirement
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R1
  - Evidence: `grep -n "fail"` cannot observe whether the exemption list or slice reporting rule is implemented.
  - Suggestion: Use an explicit pattern check for exemptions and report structure.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [pending] [ambiguity] (blocking) Undefined "substantive" finding condition in two-class classification
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R14
  - Evidence: Binary classification ("actionable or not") conflicts with undefined escalation for "substantive" findings with zero actionable.
  - Suggestion: Define concrete criteria for what makes a non-actionable finding substantive.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [pending] [excess] (non-blocking) Unrequested reduction of round cap from 6 to 3
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R16
  - Evidence: SOURCE.md does not request changing the build/QA and review cap from 6 to 3.
  - Suggestion: Drop R16 or document under Decisions as an unrequested proposal requiring owner sign-off.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P7** [pending] [organization] (non-blocking) Review-bundle stage omission in Section 1 hypothesis
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 1, Hypothesis
  - Evidence: Source findings list review-bundle rules as a weight contributor, which PRD replaces with in-flight doubt despite removing bundles in R17.
  - Suggestion: Include review-bundle rules in Section 1's stage list.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- File lines and bytes measurements match SOURCE.md working tree measurements.
- Size budget formula (<= 330 lines, <= 43000 bytes) correctly mirrors 40% target.
- Three test heuristics match SOURCE.md recommendation 2 verbatim.
- Commit content requirements (why, what, how) match recommendation 3; D1 soundly justifies landing in `delivery-readiness`.
- Procedure step length and cross-reference constraints in R5 match recommendation 4.
- All seven load-bearing safeguards in R8 correctly preserve dispatch requirements.
- Agent (13) and command (19+1) counts preserved with no new `project.yaml` keys.
- Appendix accurately excludes all eight rejected source items and future planning direction.
- Cost, risk, and recommendations provided for R11, R12, R15, and R17.

## Round 2

| Field | Value |
|---|---|
| Document | docs/prd/PRD-031-lean-builder-flow.md |
| Source | tasks/PRD-031/SOURCE.md (harness dir) |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-30 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 0. Raw reply: `tasks/PRD-031/evidence/document-r2/agy-output.md` (harness dir). The document was corrected between the rounds; the round-2 review is on the corrected document.

## Comprehension
- What the document proposes: The document establishes a mandatory test-first development workflow for code implementation slices across both `incremental-implementation` and `regression-testing`, exempting only documentation, configuration, and copy slices. It equips `regression-testing` with three heuristics to prevent empty or tautological tests and requires commit message bodies to document why, what, and how changes were verified. It limits new harness procedure steps to a single imperative sentence with at most one sentence of rationale, moving cross-references to step endings. It enforces a compression ceiling restricting four core builder-phase files to at most 40% of their current size while preserving all load-bearing safeguards through a comprehensive rule inventory. Finally, it packages seven discrete review-stage simplifications—such as single clean round completion, collapsed finding states, and narrowed doubt cycles—for individual owner approval.
- How the current system works, as the document states it: Four procedural files govern the builder phase, consuming 830 lines and 107,478 bytes of instructional text. The harness currently enforces red-first testing only for bug regression verification, leaving general feature slices without a test-first requirement. Commit message bodies lack structural standardization, and skill authoring rules impose no brevity limits on procedure steps. Review cycles accumulate substantial token overhead from multi-stage state machines, including six finding states, granular finding identifiers, reviewer rotation cursors, and mandatory two-consecutive-clean-round verifications. In-flight doubts allow up to three fresh-context dispatches based on four-way classifications before reaching standard slice reviews.
- Sections not paraphrased: none

## Findings

- **P1** [applied] [gap] (blocking) Test-first default for implementation slices dropped from regression-testing
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R1 and R3
  - Resolution: R1 now explicitly lands test-first in both `incremental-implementation` and `regression-testing`, and verifies exact phrases in both files, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) Five finding states in source asserted as six finding states in document
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 1, line 27 and section 3, R13
  - Resolution: Section 1 explicitly cites the six states from `.agents/coordinator.md` item 8, aligned with SOURCE.md Correction 2026-10-30, round 2.
  - Note (Coordinator): the source analysis in `tasks/PRD-031/SOURCE.md` said five finding states, and the reviewer's round-1 P2 was raised against that count. A dated correction at the end of that file (2026-10-30, Coordinator) sets six states per `.agents/coordinator.md` item 8, which the PRD uses. The round-1 text of P2 is kept as the reviewer wrote it.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [conflict] (blocking) In-flight doubt remaining shape falsely attributed to R15
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R2
  - Resolution: Reference to R15 removed from R2, leaving only R14, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [weak criterion] (blocking) Incomplete grep check for R1 exemption list and report requirement
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R1
  - Resolution: Verifiable criteria updated to check exact phrases `expected failure`, `docs-only`, `config-only`, `copy-only`, `test-first`, and report line `Red:`, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [ambiguity] (blocking) Undefined "substantive" finding condition in two-class classification
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R14
  - Resolution: Undefined "substantive" phrasing removed; classification is strictly binary ("actionable or not") with explicit escalation if none are actionable, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [excess] (non-blocking) Unrequested reduction of round cap from 6 to 3
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 3, R16
  - Resolution: R16 explicitly labeled as an unrequested Coordinator proposal subject to owner approval, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P7** [applied] [organization] (non-blocking) Review-bundle stage omission in Section 1 hypothesis
  - Location: docs/prd/PRD-031-lean-builder-flow.md section 1, Hypothesis
  - Resolution: Section 1 hypothesis updated to include review-bundle rules, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- File lines and bytes measurements (830 lines, 107478 bytes) match SOURCE.md working tree measurements.
- Size budget ceiling (<= 330 lines, <= 43000 bytes) correctly mirrors 40% target.
- Three test heuristics match SOURCE.md recommendation 2 verbatim.
- Commit content requirements (why, what, how) match recommendation 3; D1 soundly justifies landing in `delivery-readiness`.
- Procedure step length and cross-reference constraints in R5 match recommendation 4.
- All seven load-bearing safeguards in R8 correctly preserve dispatch requirements.
- Agent (13), skill (22), and command (19+1) counts preserved with no new `project.yaml` keys.
- Appendix accurately excludes all eight rejected source items and future planning direction.
- Cost, risk, and recommendations provided for R11, R12, R13, R14, R15, R16, and R17.
