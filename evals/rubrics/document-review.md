# Rubric: document-review

Each criterion is answered yes or no from the observed run. A criterion that cannot be observed is recorded as not-run, never as yes.

## doc-validator-gap

1. Is the verdict `changes required`?
2. Is there a finding categorized as a gap that names discovery point 4, the admin revoking a pending invite?
3. Does that finding give a location inside PRD-gap.md rather than a general statement?
4. Is there a finding categorized as a weak criterion that quotes AC-04, "the form must be fast"?
5. Does the weak-criterion finding propose an observable threshold instead of restating the problem?
6. Is every finding marked blocking or non-blocking?
7. Does the reply contain a `Not raised` section with at least one item checked and found sound?
8. Did the run leave every fixture file unmodified?
9. Is the whole reply in English?

## doc-validator-clean

1. Is the verdict `approved`?
2. Is the finding count zero?
3. Does `Not raised` show the five discovery points mapped to requirements?
4. Are the out-of-scope items absent from the findings?
5. Did the run avoid demanding any requirement the discovery notes do not raise?
6. Did the run leave every fixture file unmodified?

## doc-validator-round-cap

1. Was `document-validator` dispatched after the PRD was written and before the owner saw it?
2. Was the dispatch made with model opus and the skill `document-review`?
3. Did the dispatch carry the PRD path, the source material and the previous report when one existed?
4. On `changes required`, was the correction dispatch limited to the listed blocking points?
5. Did validation run at most twice?
6. Did the PRD reach the owner after `approved` or after the second validation, with the report path?
7. Did the persistent report keep the round-1 identifiers and mark their states?
8. Did `document-validator` write no file?
