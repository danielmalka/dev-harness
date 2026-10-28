# Review: PRD-030

## Round 1

| Field | Value |
|---|---|
| Document | docs/prd/PRD-030-diagram-guidance.md |
| Source | docs/roadmap.html:1726-1741, docs/en/roadmap.html:1726-1741, .agents/docs-guide.md, .skills/doc-template-html/SKILL.md, .skills/doc-template-html/scripts/stamp.sh, AGENTS.md, docs/prd/PRD-029-independent-review-and-drift-report.md, owner's words |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-28 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 1. Raw reply: `tasks/AUTO-3/evidence/validation-prd-r1/reply.md` (harness dir).

## Comprehension
- What the document proposes: The document establishes concise guidance in docs-guide to write Mermaid diagrams inside `<pre class="mermaid">` blocks and link Excalidraw diagrams as separate `.excalidraw` files without embedding. It updates four hand-written HTML pages to replace floating Mermaid version imports with a pinned reference to version 11.17.2, matching the version pinned in `stamp.sh`. It ensures no untracked floating Mermaid references remain outside distribution outputs and historical records, and updates the roadmap diagram card to reflect the delivered scope. Finally, it aligns version identifiers to 0.35.0, regenerates package distribution files, and confirms that agent, command, and skill counts remain unchanged.
- How the current system works, as the document states it: The roadmap contains a planned card for specialized Excalidraw and Mermaid diagram skills, but `docs-guide` currently provides no rules on when or how to draw diagrams. While `stamp.sh` and dashboard `pdocs` already pin Mermaid to version 11.17.2, four hand-written HTML documentation pages still load Mermaid via a floating major-version tag. In addition, the roadmap card text describes Mermaid as the only permitted CDN exception, lagging behind `AGENTS.md` where Google Fonts is recognized as a second exception. An external diagram design skill was evaluated but rejected because its inline SVG emission and external runtime dependencies conflict with kit principles.
- Sections not paraphrased: none

## Findings

- **P1** [applied] [conflict] (blocking) Acceptance criterion "all say 0.35.0" contradicts ref "v0.35.0" and leaves the five version places undefined.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R5 (lines 46-47)
  - Evidence: R5 line 46 specifies `ref v0.35.0` in `.claude-plugin/marketplace.json`, while the criterion on line 47 claims all five version places say `0.35.0`. Moreover, R5 describes 7 distinct version targets across the repository.
  - Suggestion: Explicitly name each target location and its exact expected string (`0.35.0` vs `v0.35.0`) in the acceptance criteria.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] Instruction placement after Procedure step 2 leaves step numbering and the line 75 cross-reference ambiguous.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R1 (line 29)
  - Evidence: Renumbering Procedure steps desynchronizes the reference in `.agents/docs-guide.md:75` ("Procedure step 5").
  - Suggestion: Clarify that the instruction is placed as an unnumbered paragraph or sub-item within step 2 without renumbering steps 1 through 5.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [organization] Section 4 Docs omits the two tutorial files modified by R2 from its bulleted list.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 4 (lines 57-63)
  - Evidence: `docs/tutoriais/templates.html` and `docs/en/tutorials/templates.html` are modified in R2 and mentioned in prose, but missing from the list.
  - Suggestion: Add both tutorial files to the bulleted list in Section 4.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [organization] Header Status field is marked approved prior to validation completion.
  - Location: docs/prd/PRD-030-diagram-guidance.md line 5
  - Evidence: `templates/en/PRD.md` prescribes `in review` or `draft` while undergoing validation.
  - Suggestion: Change `Status` to `in review`.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Sentences in R4 and R5 exceed the 25-word limit of ASD-STE100 rules.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R4 (line 42) and R5 (line 46)
  - Evidence: Descriptive sentences in R4 and R5 each contain 28 words.
  - Suggestion: Split each sentence into two shorter sentences under 25 words.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Gaps: All points from `docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`, owner's words, and `AGENTS.md:29` mapped to R1-R5 and Out of scope.
- Excesses: No unrequested agents, skills, or commands introduced.
- Diagram rules (R1): Mermaid in `<pre class="mermaid">` as default and Excalidraw as a separate linked `.excalidraw` file with no script or embed correctly honors the owner's 2026-09-25 decision and `doc-template-html` rules.
- Pinned imports (R2): Version 11.17.2 accurately matches `stamp.sh:89` and `internal/dashboard/pdocs.go:18`, and correctly targets all four hand-written HTML pages with floating tags.
- Floating reference purge (R3): Clean tree verification via `grep` correctly excludes generated files under `dist/` and historical archives in `docs/prd/`.
- Roadmap card sync (R4): Updating status to delivered in 0.35.0, removing "única exceção de CDN" in favor of the two exceptions in `AGENTS.md:29`, and keeping pt and en in sync correctly honors roadmap maintenance rules.
- Deterministic checks (R6): Running `dh validate`, `go test ./...`, and `dh build` while preserving invariant counts (13 agents, 19+1 commands, 22 skills) accurately matches repository verification standards.

## Round 2

| Field | Value |
|---|---|
| Document | docs/prd/PRD-030-diagram-guidance.md |
| Source | docs/roadmap.html:1726-1741, docs/en/roadmap.html:1726-1741, .agents/docs-guide.md, .skills/doc-template-html/SKILL.md, .skills/doc-template-html/scripts/stamp.sh, AGENTS.md, docs/prd/PRD-029-independent-review-and-drift-report.md, owner's words |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-28 |

Reviewer: cli:agy/gemini-3.8-flash-high (document stage). Blocking findings: 0. Raw reply: `tasks/AUTO-3/evidence/validation-prd-r2/reply.md` (harness dir). PRD corrections between the rounds were made by product-discovery.

## Comprehension
- What the document proposes: The PRD adds a concise instruction to `docs-guide` establishing Mermaid inside `<pre class="mermaid">` as the default diagram format and restricting Excalidraw diagrams to linked external `.excalidraw` files. It updates four hand-written HTML documentation files by replacing floating Mermaid imports with a reference pinned to exact version 11.17.2, matching `stamp.sh`. It ensures no floating Mermaid imports remain tracked across the codebase and updates the roadmap diagram card to reflect the delivered scope under version 0.35.0. It updates repository version identifiers and changelog entries to 0.35.0, regenerates package distribution files, and verifies that agent, command, and skill counts remain strictly unchanged.
- How the current system works, as the document states it: The roadmap lists a planned card for specialized Excalidraw and Mermaid diagram skills, while `docs-guide` currently contains no instructions on when or how diagrams should be included. Although `stamp.sh` and dashboard `pdocs` pin Mermaid to version 11.17.2, four hand-written HTML documentation files continue to import Mermaid using a floating major-version tag. The roadmap card still describes Mermaid as the only allowed CDN exception, trailing `AGENTS.md` which specifies two CDN exceptions including Google Fonts. An external diagram design skill was evaluated but rejected by the owner because its inline SVG emission and external dependencies contradict kit architecture and ADR-003.
- Sections not paraphrased: none

## Findings

- **P1** [applied] [conflict] (blocking) Acceptance criterion "all say 0.35.0" contradicts ref "v0.35.0" and leaves the five version places undefined.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R5 (lines 46-47)
  - Resolution: The author updated the acceptance criterion to explicitly reference the files named in R5 and accounted for the `v0.35.0` prefix in marketplace ref ("When the files named in this rule are read, then each version field says 0.35.0 (the marketplace ref says v0.35.0)"), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] Instruction placement after Procedure step 2 leaves step numbering and the line 75 cross-reference ambiguous.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R1 (line 29)
  - Resolution: The author clarified that the instruction is added at the end of Procedure step 2 within the same list item, explicitly asserting that step numbers remain unchanged, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [organization] Section 4 Docs omits the two tutorial files modified by R2 from its bulleted list.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 4, line 62
  - Resolution: The author added both tutorial files to the Section 4 list ("docs/tutoriais/templates.html and docs/en/tutorials/templates.html: import pin only (R2)"), round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [organization] Header Status field is marked approved prior to validation completion.
  - Location: docs/prd/PRD-030-diagram-guidance.md line 5
  - Evidence: Header table declares `Status: approved (...)`, whereas `templates/en/PRD.md` prescribes `draft` or `in review` while the document is undergoing validation.
  - Suggestion: Keep as a non-blocking organization record; the status will formally transition upon coordinator acceptance. Note (Coordinator, after this round): the header Status row was set to `delivered on 2026-10-28` at the end of the PRD cycle; the reviewer's state for P4 is unchanged.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Sentences in R4 and R5 exceed the 25-word limit of ASD-STE100 rules.
  - Location: docs/prd/PRD-030-diagram-guidance.md section 3, R4 (line 42) and R5 (line 46)
  - Resolution: The author split the descriptive sentences in R4 and R5 into concise statements under 25 words each, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Gaps: All source points from `docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`, `AGENTS.md:29`, and the owner's request ("Vamos seguir com a B e pode fixar a versão do mermaid junto") are traced to R1-R6 and Out of scope.
- Ambiguities: All instructions and rules provide unambiguous single implementations. Placement of the docs-guide instruction at the end of Procedure step 2 preserves step numbering and `.agents/docs-guide.md:75` cross-references.
- Conflicts: Target version string specifications in R5 match Git tag requirements (`v0.35.0` for marketplace ref vs `0.35.0` for package versions). Pinned version `11.17.2` matches `stamp.sh:89` and `internal/dashboard/pdocs.go:18`. CDN exception statements match `AGENTS.md:29`.
- Excesses: No new agents, skills, or commands are introduced; out of scope explicitly rejects embedding, SVG export, and external diagram dependencies in accordance with owner decision B and ADR-003.
- Weak criteria: All acceptance criteria under R1-R6 use deterministic "When ..., then ..." structures verifiable by exact string matching, HTTP status, file contents, or tool execution.
- Organization: Requirements and documentation impact lists are properly structured and categorized.
- PRD writing style: All descriptive sentences conform to the 25-word limit specified by ASD-STE100 guidelines.
