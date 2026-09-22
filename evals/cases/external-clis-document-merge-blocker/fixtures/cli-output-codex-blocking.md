# Transport double: cli:codex/gpt-5.6-sol, one blocking gap

Pre-written reply, used as fixture data only. No real `codex` binary is
invoked to produce this file or to grade any case that reads it. Represents
a reply that carries a real verdict signal (`## Verdict`) with one blocking
gap finding. Follows the Output format of
`.skills/document-review/SKILL.md`.

```
## Verdict
- changes required
- Blocking findings: 1
- Report: external-clis-document-merge-blocker.review.md

## Findings
- P1
  - Category: gap
  - Severity: blocking
  - Location: fixtures/document.md, Scope
  - Evidence: fixtures/source.md point 3 states an admin can revoke a
    pending invite before it is redeemed, and a revoked invite cannot be
    redeemed; no requirement in fixtures/document.md covers revocation and
    no acceptance criterion mentions it
  - Suggestion: add "RF-05 An admin revokes a pending invite, and a
    revoked invite cannot be redeemed" with an acceptance criterion
    asserting the redemption is refused after revocation
  - Reported by: cli:codex/gpt-5.6-sol

## Not raised
- Invite creation (RF-01), expiry (RF-02), the admin-only view of pending
  invites (RF-03) and the pre-filled email on redemption (RF-04) all map to
  source points 1, 2, 4 and 5.

## Evidence
- Reviewed fixtures/document.md against fixtures/source.md. No real codex
  binary invoked; this is a pre-written fixture double.
```
