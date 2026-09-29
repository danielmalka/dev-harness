---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)start[ -]?stage[^\n]{0,30}discover'
flags: ''
---
Deterministic check, replacing a `type: llm` judge that split 2-1 FAIL on a
verifiably correct brief (`evals/baselines/2026-09-28-auto-cases-0.9.0.md`
r6: the brief `Write` records all seven items verbatim and the 2-1 split is
judge variance). Proves questionnaire item 1 only (start stage): the
persisted brief names the start stage as `discover` on the same line,
within 30 characters after the words "start stage". Offline-validated
(Python `re.search`, a faithful proxy for this runtime's regex engine —
none of this file's or its five siblings' patterns use lookaheads or
dotall) against the r6 evidence text: PASSES on the full positive text; a
synthetic negative with item 1 deleted FAILS this grader alone and does
not false-pass on the bare word "discover" that still appears in the
fixture's own "## Request" section, since that occurrence sits far outside
this pattern's 30-character, same-line window; the unmodified placeholder
fixture (`Stage: not-started`) FAILS. Items 2-6 are each proven by their
own sibling grader (`06-` through `10-questionnaire-*.md`); item 7 stays
covered by the existing, unmodified grader 04.
