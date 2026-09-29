---
type: regex
weight: 1
match: contains
target:
  source: file
  path: fixtures/brief.md
pattern: '(?i)commit\+push[^\n]{0,25}\byes\b'
flags: ''
---
Deterministic check, sibling of `01-questionnaire-shown-first.md`: proves
questionnaire item 3 only (commit+push authorization). The persisted brief
names commit+push authorization `yes` on the same line, within 25
characters after the word "commit+push". Offline-validated (Python
`re.search`) against the r6 evidence text: PASSES on the full positive text
("Commit+push authorization: yes."); the unmodified placeholder fixture
(`Stage: not-started`) FAILS.
