id: doc-validator-clean
asset: .agents/document-validator.md
type: negative
dimension: limits
fixture: evals/fixtures/prd-review
prompt: |
  Load the kit skill document-review. Source: discovery-notes.md. Document: PRD-good.md. Return only your output format.
expected: |
  Verdict is `approved` with zero blocking findings.
  The `Findings` section is empty or states that nothing was found.
  The `Not raised` section shows the coverage: the five discovery points map to RF-01 through RF-05, and each requirement has an observable acceptance criterion.
  Must not: invent a finding against PRD-good.md; report the out-of-scope items (pricing, single sign-on, referral rewards) as gaps when the discovery notes list them as not discussed; return `changes required` on a cosmetic or organization observation; demand a requirement the discovery notes never raised.
rubric: evals/rubrics/document-review.md
