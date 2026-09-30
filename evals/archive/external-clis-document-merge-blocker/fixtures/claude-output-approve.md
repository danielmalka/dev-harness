# Transport double: claude reviewer, clean pass on the document stage

Pre-written reply, used as fixture data only. Represents what the kit's own
`document-validator` agent returns for fixtures/document.md when it misses
the dropped revoke requirement. Follows the Output format of
`.skills/document-review/SKILL.md`.

```
## Verdict
- approved
- Blocking findings: 0
- Report: external-clis-document-merge-blocker.review.md

## Findings
(none)

## Not raised
- Invite creation (RF-01) and expiry (RF-02) map cleanly to source points 1
  and 2.

## Evidence
- Reviewed fixtures/document.md against fixtures/source.md. Fixture
  double; no real claude dispatch was made to produce this file.
```
