---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)decision mode[^\n]{0,20}\bb\b'
flags: ''
---
Deterministic check, sibling of `01-questionnaire-shown-first.md`: proves
questionnaire item 2 only (decision mode). The persisted brief names
decision mode `b` on the same line, within 20 characters after the words
"decision mode". Offline-validated (Python `re.search`) against the r6
evidence text: PASSES on the full positive text ("Decision mode: (b) — ask
only product or contract decisions..."); a synthetic negative with item 5
deleted still PASSES this grader (the two items are independent lines);
the unmodified placeholder fixture (`Stage: not-started`) FAILS.
