# Review: PRD-032

## Round 1

| Field | Value |
|---|---|
| Document | docs/prd/PRD-032-lean-coordinator-context.md |
| Source | tasks/PRD-032/SOURCE.md (harness dir) |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-30 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 4. Raw reply: `tasks/PRD-032/evidence/document-r1/agy-reply.txt` (harness dir).

Note (Coordinator): the reviewer reported running `go run ./cmd/dh validate .` and `wc -c`, against the anti-delegation clause; the tree was unchanged and the verdict was kept by the Coordinator. The reviewer's own reply names its report as `docs/prd/PRD-032-lean-coordinator-context.md.review.md`; this file is the path the Coordinator set.

## Comprehension
- What the document proposes: PRD-032 proposes making the Coordinator session leaner during long tasks by adopting patterns inspired by Google's Elephant-Goldfish Model. Specialist agents write reports exceeding 40 lines to files in the task folder, returning only a compact summary and verdict within a 200-word cap. When multiple reviewers evaluate a change in a stage, docs-guide synthesizes their report files into a single merged artifact rather than sending raw reports into the Coordinator chat. Single-command guidance sections move out of coordinator.md into their respective command files, chained command inlining in /dh:auto is replaced by a targeted section-reading table, and MEMORY.md is governed by entry and file size limits with doctor warnings. Finally, an opt-in repository bootstrap phase is added to /dh:setup to generate a concise, reusable REPO-MAP.md without adding new commands or agents.
- How the current system works, as the document states it: The Coordinator currently loads large instruction and state files into its context on every execution, including coordinator.md and MEMORY.md, and /dh:auto inlines six full chained command files across the entire run. Specialist dispatches lack numeric reply caps, allowing full reports such as extensive code reviews to return directly into the main conversation. The Coordinator itself merges review findings in chat when multiple reviewers report, and it reads reports in full even when changes are approved without blocking issues. MEMORY.md lacks a byte ceiling and can expand to thousands of lines, while consolidation happens only on explicit owner command. Furthermore, the kit does not persist a repository map, leaving subsequent planning and scouting runs to repeatedly rediscover layout and architecture from scratch.
- Sections not paraphrased: none

## Findings

- **P1** [pending] [conflict] `repo-scout` report-to-file instruction in R1 conflicts with no-write toolset in R2
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46) and R2 (line 50)
  - Evidence: R1 asserts in its verifiable criteria that `repo-scout.md` states the report file rule in its Output format ("builder.md, reviewer.md, qa-verifier.md, docs-guide.md, and repo-scout.md each state the same rule in their Output format", where the rule is "A specialist writes its full report to `<task folder>/reports/<role>-<n>.md` when the report is longer than 40 lines"). However, R2 explicitly contradicts this for `repo-scout`: "`repo-scout` has no Write tool. Its reply stays within the cap and lists paths, not file bodies. When a scout map is longer, the Coordinator dispatches `docs-guide` to write it (R10)."
  - Suggestion: Remove `repo-scout.md` from the list of agents stating the report-to-file rule in R1. Update R1 to specify that `repo-scout`'s Output format enforces only the reply cap, delegating any report file writing to `docs-guide` per R2 and R10.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [pending] [ambiguity] Review-specific fields mandated for non-review specialist replies in R1
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46)
  - Evidence: R1 dictates: "A specialist writes its full report to `<task folder>/reports/<role>-<n>.md` when the report is longer than 40 lines. Its reply then holds the verdict, the report path, the counts by severity, at most 5 blocking findings as `location · title`, and a summary within the cap. ... builder.md, reviewer.md, qa-verifier.md, docs-guide.md, and repo-scout.md each state the same rule in their Output format." Non-review specialists (`builder`, `qa-verifier`, `docs-guide`, `debugger`) do not produce verdicts, severity counts, or blocking findings.
  - Suggestion: Disambiguate R1 by distinguishing review specialist replies from non-review specialist replies: state that reviewer replies include verdict, report path, severity counts, and top 5 blocking findings, whereas non-review specialists (builder, qa-verifier, docs-guide, debugger) include the report path and a role-appropriate status summary (e.g., delivered slices, check results, or updated docs) within the 200-word cap.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [pending] [ambiguity] Unspecified execution source for `git rev-parse HEAD` in R10 and R12
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R10 (lines 78, 84) and R12 (line 86)
  - Evidence: R10 requires `REPO-MAP.md` to include "1. Header: date, `git rev-parse HEAD`, and the command that produced it" while simultaneously stating "The writer reads the repository and does not run project code or scripts", and `.agents/docs-guide.md` lacks a shell tool ("its stamp script needs a shell, which this role does not have"). Similarly, R12 states "When the header sha differs from the current `HEAD`, `repo-scout` reports the map as stale", but `.agents/repo-scout.md` has no shell tool ("Product read-only. No Write, Edit, or shell"). Neither rule defines who executes `git rev-parse HEAD` or how the current commit sha reaches `docs-guide` and `repo-scout`.
  - Suggestion: Explicitly state in R10 and R12 that the Coordinator (or a dispatched `qa-verifier` in lane `infra`) resolves the current `git rev-parse HEAD` sha and passes it in the dispatch payload to `docs-guide` during setup/bootstrap and to `repo-scout` during reuse.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [gap] `.skills/security-review/SKILL.md` omitted from documentation update scope in Section 4
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 4, Docs (line 124)
  - Evidence: Section 4 lists `.skills/code-review/SKILL.md`, but omits `.skills/security-review/SKILL.md`. Under R2, R4, and `coordinator.md`, the artifact-first report file rule, the 200-word reply cap, and the synthesizer merge apply to both `code` and `security` review stages. The kit reviewer executes mode `security` using `.skills/security-review/SKILL.md`.
  - Suggestion: Add `.skills/security-review/SKILL.md` to the document list in Section 4.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [pending] [organization] `debugger.md` listed in R2 and Section 4 but omitted from R1 verifiable output format list
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46) and R2 (line 48)
  - Evidence: R2 includes `debugger` ("builder, qa-verifier, docs-guide, and debugger write their own report file") and Section 4 lists `.agents/debugger.md` for update. However, R1's verifiable condition omits `debugger.md` ("builder.md, reviewer.md, qa-verifier.md, docs-guide.md, and repo-scout.md each state the same rule in their Output format").
  - Suggestion: Add `debugger.md` to R1's list of specialist definitions updated with the report-to-file rule.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [pending] [organization] Procedural and descriptive sentence length violations against PRD writing rules (ASD-STE100)
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R4 (line 54) and R5 (line 64)
  - Evidence: R4 sentence 1 ("When a review stage has more than one reviewer entry in the dispatch, a docs-guide run (model sonnet) merges the reviewer report files into one file <task folder>/reports/<stage>-merged-<round>.md.") has 30 words (limit: 20 procedural / 25 descriptive). R5 sentence 3 ("auto.md carries a table with one row for each chained stage (discover, plan, build, secure, review, release): the file, the sections to read, and the byte count budget.") has 29 words (limit: 25 words).
  - Suggestion: Break these sentences into shorter statements of at most 20 words each, following the kit's PRD writing rules (ASD-STE100).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P7** [pending] [ambiguity] Threshold precedence between 40-line report file trigger and 200-word reply cap
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46)
  - Evidence: R1 specifies that "The default cap is 200 words" and "A specialist writes its full report to `<task folder>/reports/<role>-<n>.md` when the report is longer than 40 lines. ... A report of 40 lines or fewer stays in the reply." A report with 30 lines but 300 words satisfies the condition to stay in the reply while violating the 200-word reply cap.
  - Suggestion: Reconcile the line and word thresholds in R1, for example: "A specialist writes its full report to the file when the report exceeds 40 lines or exceeds the reply cap; a report meeting both limits stays in the reply."
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised
- Lean Coordinator context motivation: Owner request for reduced Coordinator context load in long tasks is directly addressed across R1-R8.
- Elephant-Goldfish Model alignment: Goldfish read-only subagents, reports written to disk, and deduplicated synthesis (EGM findings in SOURCE.md) map cleanly to R1, R2, R3, and R4.
- Repository bootstrap mechanism: Extending `/dh:setup` via an opt-in second phase preserving the fixed 19 + 1 command count maps cleanly to R9-R12 without altering documented command limits.
- Chained command inlining reduction: R5 correctly replaces inlining of full command files in `/dh:auto` with a six-row table targeting only `Routing`, `Prerequisites`, and `Limits`.
- Coordinator single-command section reduction: R6 moves four candidate sections ("Worktree", "Consolidation on explicit request", "External CLI planners", "Drift report") to their respective commands and targets a verified >= 5500 byte reduction.
- Memory budget and consolidation: R7 establishes concrete entry (<= 1200 chars) and file (<= 40000 bytes) limits, and R8 provides observable `dh doctor` size checks without violating owner-triggered consolidation.
- PRD-031 dependency and timing: Header, Section 1, R6, R15, and Section 4 correctly state and respect the dependency on the merged state of PRD-031.
- Documented counts integrity: R13 preserves 13 agents, 22 skills, 19+1 commands, 7 profiles, 14 templates, and prohibits new `project.yaml` keys.
- Safeguards preservation: R14 retains not-run honesty, memory ownership, independent review, Coordinator conversation exclusion, the anti-delegation clause, authorization limits, and tree immutability for CLI reviewers.
- Measurement discipline: R15 defines byte baselines and results in the Executed summary without relying on paid evaluations or token-counting tools.
- Cleanup and reproducibility: R16 mandates clean diff in `dist/` after build and clean-clone proof for changes touching Go code in `internal/`.
- Out-of-scope boundaries: Appendix accurately rejects EGM hooks, six-lens brainstorm, README trees, model modifications, and automated memory rewriting.

## Round 2

| Field | Value |
|---|---|
| Document | docs/prd/PRD-032-lean-coordinator-context.md |
| Source | tasks/PRD-032/SOURCE.md (harness dir) |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-31 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 0. Raw reply: `tasks/PRD-032/evidence/document-r2/agy-reply.txt` (harness dir). The document was corrected between the rounds; the round-2 review is on the corrected document.

## Comprehension
- What the document proposes: PRD-032 proposes making the Coordinator session leaner during long tasks by adopting patterns inspired by Google's Elephant-Goldfish Model. Specialist agents write reports exceeding 40 lines to files in the task folder, returning only a compact, role-tailored summary within a 200-word cap. When multiple reviewers evaluate a change in a stage, docs-guide synthesizes their report files into a single merged artifact rather than sending raw reports into the Coordinator chat. Single-command guidance sections move out of coordinator.md into their respective command files, chained command inlining in /dh:auto is replaced by a targeted section-reading table, and MEMORY.md is governed by entry and file size limits with doctor warnings. Finally, an opt-in repository bootstrap phase is added to /dh:setup to generate a concise, reusable REPO-MAP.md without adding new commands or agents.
- How the current system works, as the document states it: The Coordinator currently loads large instruction and state files into its context on every execution, including coordinator.md and MEMORY.md, and /dh:auto inlines six full chained command files across the entire run. Specialist dispatches lack numeric reply caps, allowing full reports such as extensive code reviews to return directly into the main conversation. The Coordinator itself merges review findings in chat when multiple reviewers report, and it reads reports in full even when changes are approved without blocking issues. MEMORY.md lacks a byte ceiling and can expand to thousands of lines, while consolidation happens only on explicit owner command. Furthermore, the kit does not persist a repository map, leaving subsequent planning and scouting runs to repeatedly rediscover layout and architecture from scratch.
- Sections not paraphrased: none

## Findings

- **P1** [applied] [conflict] `repo-scout` report-to-file instruction in R1 conflicts with no-write toolset in R2
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (lines 52-53) and R2 (line 57)
  - Resolution: R1 now explicitly specifies that `repo-scout` has no Write tool and always returns inline within the cap, while R2 and R10 delegate any report file writing to `docs-guide`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] Review-specific fields mandated for non-review specialist replies in R1
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (lines 47-52)
  - Resolution: R1 now specifies distinct, role-appropriate reply formats for `reviewer`, `qa-verifier`, `builder`, `docs-guide`, `debugger`, and `repo-scout`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] Unspecified execution source for `git rev-parse HEAD` in R10 and R12
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R10 (line 86) and R12 (line 94)
  - Resolution: R10 and R12 now explicitly designate `qa-verifier` to resolve `git rev-parse HEAD` and state that the Coordinator passes the resulting sha in the dispatch to `docs-guide` and `repo-scout`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [gap] `.skills/security-review/SKILL.md` omitted from documentation update scope in Section 4
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 4, Docs (line 133)
  - Resolution: `.skills/security-review/SKILL.md` is now included in the Section 4 documentation update list.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] `debugger.md` listed in R2 and Section 4 but omitted from R1 verifiable output format list
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 53)
  - Resolution: `debugger.md` is now explicitly included in the verifiable criteria of R1 alongside the other specialists.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] Procedural and descriptive sentence length violations against PRD writing rules (ASD-STE100)
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R4 (line 61) and R5 (line 72)
  - Resolution: Long sentences in R4 and R5 were split into concise statements satisfying the ASD-STE100 length limits (<= 20 words for procedural, <= 25 words for descriptive).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P7** [applied] [ambiguity] Threshold precedence between 40-line report file trigger and 200-word reply cap
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46)
  - Resolution: R1 clarifies that the 200-word cap strictly bounds every reply, ensuring that reports staying within the reply must satisfy both the line and word thresholds.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Gaps: Checked all points from SOURCE.md (EGM goldfish pattern, file reports, synthesis, repo analysis/bootstrap, context levers, MEMORY.md growth, PRD-031 baseline dependency). All source requirements are fully mapped to R1-R16.
- Ambiguities: Checked role-specific reply shapes (R1), write permissions (R2), coordinator reading discipline (R3), synthesizer rules (R4), `/dh:auto` table usage (R5), section moves (R6), memory limits (R7), doctor warnings (R8), and git sha resolution via qa-verifier for REPO-MAP.md (R10, R12). Every rule has unambiguous operational semantics.
- Conflicts: Checked internal consistency and alignment with AGENTS.md rules. Toolsets match write sets (R2, R10), review immutability is preserved (R2, R14), command counts remain fixed at 19+1 (R9, R13), and memory ownership stays strictly with the Coordinator.
- Excesses: Checked requirement boundaries against SOURCE.md and out-of-scope declarations. No extra commands, skills, agents, or templates were introduced, and discarded EGM features (six-lens brainstorm, README trees, hooks) are explicitly excluded.
- Weak criteria: Checked verifiable criteria across all rules (R1-R16). Every condition relies on deterministic checks (exact strings, numbers, line limits, byte reduction targets, and `dh validate` assertions) rather than subjective expectations.
- Organization: Checked document structure, continuous rule numbering (R1-R16), balanced Solution-to-Rules mapping (Parts A-D), ASD-STE100 sentence lengths, and paired documentation updates in Section 4.
- PRD writing rules (ASD-STE100): Procedural and descriptive statements follow active voice, present tense, and length constraints (<= 20 words for procedural, <= 25 words for descriptive).

## Round 3 (amendment revalidation)

| Field | Value |
|---|---|
| Document | docs/prd/PRD-032-lean-coordinator-context.md |
| Source | tasks/PRD-032/SOURCE.md (harness dir) |
| Round | 3 of 3 (amendment revalidation) |
| Verdict | approved |
| Updated | 2026-10-31 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 0. Raw reply: `tasks/PRD-032/evidence/document-r3/agy-reply.txt` (harness dir).

Amendment reviewed: D1-D8 closed on the Coordinator's recommendations (owner, 2026-10-31), new rule R17 (PRD-031 carry-over fixes), R6 target from the merged 0.36.0 size.

Note (docs-guide): the raw reply file ends with a truncated trailing fragment after its Evidence section, a duplicate report body cut off mid-word and followed by "error: interrupted". The verdict and the sections above it are complete; the fragment is not reproduced here.

## Comprehension
- What the document proposes: PRD-032 establishes a leaner execution profile for the Coordinator during long tasks by adopting file-based artifact reporting and synthesizer deduplication inspired by the Elephant-Goldfish Model. Dispatched specialists must adhere to a 200-word reply cap and output reports longer than 40 lines to dedicated task-folder files rather than returning verbose text into the main conversation. Chained command inlining in `/dh:auto` is replaced by a structured table reading only stage-essential sections, single-command guidance sections move out of `coordinator.md` to achieve a 5500-byte reduction, and `MEMORY.md` is constrained by strict entry and file budgets with doctor warnings. An opt-in repository bootstrap phase is integrated into `/dh:setup` to generate and maintain a compact, reusable `REPO-MAP.md` without increasing the fixed agent or command counts. In addition to formalizing closed decisions D1 through D8, the amendment introduces rule R17 to resolve four carried-over documentation defects from PRD-031 delivery.
- How the current system works, as the document states it: The Coordinator presently loads substantial static instruction and state files into its context across all sessions, most notably `coordinator.md` and `MEMORY.md`. Execution of `/dh:auto` inlines all six chained command files across the entire run, maintaining unnecessary instructions in the primary conversation. Specialist agent dispatches lack numeric reply limits, permitting extensive review and build reports to return directly into chat and triggering full Coordinator re-reads even when work is approved. Unconstrained record accumulation allows `MEMORY.md` to exceed tens of thousands of bytes without automated size monitoring, leaving consolidation entirely dependent on manual owner initiation. Finally, the harness persists no baseline architectural layout or repository map, forcing subsequent scouting and planning passes to repeatedly re-examine code structure from scratch.
- Sections not paraphrased: none

## Findings

No new findings. P1-P7 from round 2 are applied, as the reply reports:

- **P1** [applied] [conflict] `repo-scout` report-to-file instruction in R1 conflicts with no-write toolset in R2
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (lines 52-53) and R2 (line 57)
  - Resolution: R1 now explicitly specifies that `repo-scout` has no Write tool and always returns inline within the cap, while R2 and R10 delegate any report file writing to `docs-guide`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] Review-specific fields mandated for non-review specialist replies in R1
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (lines 47-52)
  - Resolution: R1 now specifies distinct, role-appropriate reply formats for `reviewer`, `qa-verifier`, `builder`, `docs-guide`, `debugger`, and `repo-scout`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] Unspecified execution source for `git rev-parse HEAD` in R10 and R12
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R10 (line 86) and R12 (line 94)
  - Resolution: R10 and R12 now explicitly designate `qa-verifier` to resolve `git rev-parse HEAD` and state that the Coordinator passes the resulting sha in the dispatch to `docs-guide` and `repo-scout`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [gap] `.skills/security-review/SKILL.md` omitted from documentation update scope in Section 4
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 4, Docs (line 133)
  - Resolution: `.skills/security-review/SKILL.md` is now included in the Section 4 documentation update list.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] `debugger.md` listed in R2 and Section 4 but omitted from R1 verifiable output format list
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 53)
  - Resolution: `debugger.md` is now explicitly included in the verifiable criteria of R1 alongside the other specialists.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] Procedural and descriptive sentence length violations against PRD writing rules (ASD-STE100)
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R4 (line 61) and R5 (line 72)
  - Resolution: Long sentences in R4 and R5 were split into concise statements satisfying the ASD-STE100 length limits (<= 20 words for procedural, <= 25 words for descriptive).
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P7** [applied] [ambiguity] Threshold precedence between 40-line report file trigger and 200-word reply cap
  - Location: docs/prd/PRD-032-lean-coordinator-context.md section 3, R1 (line 46)
  - Resolution: R1 clarifies that the 200-word cap strictly bounds every reply, ensuring that reports staying within the reply must satisfy both the line and word thresholds.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Gaps: Checked all points from `SOURCE.md` (EGM goldfish pattern, offloaded file reports, synthesis, repo analysis/bootstrap, context levers, `MEMORY.md` growth, PRD-031 baseline dependency), owner words (closing decisions D1-D8, noting separate bootstrap command in memory), and PRD-031 Executed summary (all four carried-over defects). All source points are completely mapped to R1-R17 and D1-D8.
- Ambiguities: Checked role-specific reply shapes (R1), write permissions (R2), coordinator reading discipline (R3), synthesizer rules (R4), `/dh:auto` table usage (R5), section moves (R6), memory limits (R7), doctor warnings (R8), git sha resolution via qa-verifier for `REPO-MAP.md` (R10, R12), and the four specific defect repairs in R17. Every rule and verifiable check provides unambiguous operational semantics.
- Conflicts: Checked internal consistency and alignment with `AGENTS.md` rules. Toolsets match write sets (R2, R10), review immutability is preserved (R2, R14), command counts remain fixed at 19+1 (R9, R13, D6), and memory ownership stays strictly with the Coordinator. The updated R6 target (`<= 26342` bytes) matches the baseline calculation (`31842 - 5500`) exactly.
- Excesses: Checked requirement boundaries against `SOURCE.md`, owner instructions, and out-of-scope declarations. No unrequested commands, skills, agents, or templates were introduced, and discarded EGM features (six-lens brainstorm, README trees, hooks) remain explicitly excluded.
- Weak criteria: Checked verifiable criteria across all rules (R1-R17). Every condition relies on deterministic checks (exact string occurrences, character/byte limits, line caps, doctor unit tests in `internal/`, clean diffs, and exact grep commands) rather than subjective expectations.
- Organization: Checked document structure, continuous rule numbering (R1-R17), balanced Solution-to-Rules mapping (Parts A-D), closed decisions D1-D8, and paired documentation updates in Section 4.
- PRD writing rules (ASD-STE100): Procedural and descriptive statements follow active voice, present tense, and length constraints (<= 20 words for procedural, <= 25 words for descriptive).
