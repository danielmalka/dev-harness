---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)merge authorization[^\n]{0,20}\bno\b'
flags: ''
---
Deterministic check, sibling of `01-questionnaire-shown-first.md`: proves
questionnaire item 5 only (merge authorization). The persisted brief names
merge authorization `no` on the same line, within 20 characters after the
words "merge authorization". Offline-validated (Python `re.search`)
against the r6 evidence text: PASSES on the full positive text ("Merge
authorization: no."); a synthetic negative with this exact line deleted
FAILS this grader alone, and the other five item-graders still PASS on
that same negative; the unmodified placeholder fixture (`Stage:
not-started`) FAILS.
