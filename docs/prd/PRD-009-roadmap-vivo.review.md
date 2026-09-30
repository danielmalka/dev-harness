# Review · PRD-009 · round 1

- Reviewer: document-validator (claude, opus requested; effective model unverified)
- Date: 2026-09-30
- Source: owner's request (card "Este próprio roadmap como artefato vivo"); AGENTS.md decision of 2026-09-30 (option A); Coordinator follow-up decisions on the owner's behalf.

## Verdict
- changes required
- Blocking findings: 2

## Findings
- P1 · conflict · blocking · §3 R10 (and §2). The generic kit assets (`delivery-readiness` readiness table, Coordinator docs-guide dispatch) would carry a dev-harness-specific roadmap row naming `docs/roadmap.html` and `dh validate`, while R1 says consumer projects have no roadmap; the "Met when" gate has no not-applicable wording for that case. Suggestion: conditional, path-neutral row ("when the project keeps a roadmap tied to its changelog …; not applicable when the project keeps no roadmap"), exact Coordinator sentence with the same condition, and R10's Verificável names the no-roadmap case. Reported by: claude
- P2 · ambiguity · blocking · §3 R1. Only one of the two roadmaps present is unspecified (skip vs. error). Suggestion: error naming the missing file; skip only when both are absent; Go test for it. Reported by: claude
- P3 · ambiguity · non-blocking · §3 R8 vs R5/R7. Line for a missing card and file for a pt/en mismatch unspecified. Suggestion: missing card → the `<h2 id="entregue">` line; R7 → cite `docs/en/roadmap.html`. Reported by: claude
- P4 · gap · non-blocking · §4 Docs. `docs/en/tutorials/` mirror missing from the conditional line. Reported by: claude

## Not raised
- 22 CHANGELOG headings ↔ 22 version cards per language, already newest first; two detail cards per language; R2–R4 formats match the tree; coordinator.md lines 72/76; release-manager has no enumerated checklist; R12 re-embed and no paid eval consistent with AGENTS.md; no contradiction with AGENTS.md.

---

# Review · PRD-009 · round 2 (final)

- Reviewer: document-validator (claude, opus requested; effective model unverified). Date: 2026-09-30.

## Verdict
- approved
- Blocking findings: 0

## Settlement of round 1
- P1–P4: applied and closed (R10/§2 conditional and path-neutral with exact sentences and greps; R1 skip only when both roadmaps are absent, one present is an error; R8 cites the `<h2 id="entregue">` line and `docs/en/roadmap.html`; Docs line lists `docs/en/tutorials/`).

## New findings
- P5 · ambiguity · non-blocking · §3 R10. "passes the project's check" and "one card per release" assume a check and a card-based roadmap in every consumer. Suggestion: "the project's roadmap check, when it has one"; "the new release recorded", keeping "one card per release" only inside the dev-harness parenthetical. Reported by: claude
- Coordinator: P5 applied verbatim to R10 after round 2 (wording only; the software built is the same).
