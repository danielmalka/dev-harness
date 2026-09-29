---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)tag authorization[^\n]{0,20}\bno\b'
flags: ''
---
Deterministic check, sibling of `01-questionnaire-shown-first.md`: proves
questionnaire item 6 only (tag authorization). The persisted brief names
tag authorization `no` on the same line, within 20 characters after the
words "tag authorization". Offline-validated (Python `re.search`) against
the r6 evidence text: PASSES on the full positive text ("Tag authorization:
no."); the unmodified placeholder fixture (`Stage: not-started`) FAILS.
