id: doc-validator-gap
asset: .agents/document-validator.md
type: positive
dimension: evidence
fixture: evals/fixtures/prd-review
prompt: |
  Load the kit skill document-review. Source: discovery-notes.md. Document: PRD-gap.md. Return only your output format.
expected: |
  Verdict is `changes required` with at least one blocking finding.
  A gap finding cites discovery point 4 (an admin revokes a pending invite before it is redeemed) as absent from section 4 of PRD-gap.md, and suggests a concrete requirement plus an acceptance criterion for it.
  A weak-criterion finding cites AC-04 ("the form must be fast") as unverifiable and suggests an observable threshold.
  Both findings carry a category, a severity of blocking, a location in PRD-gap.md, the source evidence, and a suggestion.
  The `Not raised` section states at least one thing checked and found sound.
  Must not: edit PRD-gap.md or any fixture file; address the owner instead of the coordinator; reply in a language other than English; raise the missing resend acceptance criterion and the missing revocation requirement as one merged finding without citing either location.
rubric: evals/rubrics/document-review.md
