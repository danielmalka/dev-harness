id: doc-validator-round-cap
asset: .commands/discover.md
type: positive
dimension: routing
fixture: evals/fixtures/prd-review
prompt: |
  /dev-harness:discover invite-only signup, from the notes in discovery-notes.md
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the kit loaded and is not graded from a single agent reply.
  After `product-discovery` writes the PRD, the coordinator dispatches `document-validator` with model opus and the kit skill `document-review`, passing the PRD path, the discovery notes and the previous review report when one exists.
  On `changes required`, the coordinator re-dispatches `product-discovery` with only the listed blocking points, then validates again.
  Validation runs at most twice: after the second validation the PRD goes to the owner whatever the verdict, with the open points and the path `<document>.review.md`.
  The persistent report keeps the identifiers from round 1 and marks each one `applied`, `rejected` or `pending`.
  Must not: deliver the PRD to the owner before the first validation; open a third correction round; reset the cap by changing the model; let `document-validator` edit the PRD or any of the three memory records.
rubric: evals/rubrics/document-review.md
