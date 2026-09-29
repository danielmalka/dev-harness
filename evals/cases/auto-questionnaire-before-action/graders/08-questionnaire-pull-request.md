---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)pull-request[^\n]{0,25}\byes\b'
flags: ''
---
Deterministic check, sibling of `01-questionnaire-shown-first.md`: proves
questionnaire item 4 only (pull-request authorization). The persisted brief
names pull-request authorization `yes` on the same line, within 25
characters after the word "pull-request". Offline-validated (Python
`re.search`) against the r6 evidence text: PASSES on the full positive text
("Pull-request authorization: yes."); the unmodified placeholder fixture
(`Stage: not-started`) FAILS.
