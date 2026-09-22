# PRD: teammate invites

Fictional PRD, used as a fixture only. Written from fixtures/source.md, and
covers every point of it — this fixture is deliberately clean, so the
`document` review stage of `external-clis-document-only` produces zero
findings.

## Scope

- RF-01 An admin invites a new teammate by email (source point 1).
  - AC-01 When an admin submits an email, then a pending invite is created
    for it.
- RF-02 An invite expires after 7 days (source point 2).
  - AC-02 When 7 days pass without redemption, then the invite can no
    longer be redeemed.
- RF-03 An admin revokes a pending invite before it is redeemed, and a
  revoked invite cannot be redeemed (source point 3).
  - AC-03 When an admin revokes a pending invite, then redeeming it
    afterward is refused.
- RF-04 Only admins can view the list of pending invites (source point 4).
  - AC-04 When a non-admin requests the pending-invites list, then the
    request is refused.
- RF-05 A redeemed invite creates the teammate's account with the invited
  email pre-filled (source point 5).
  - AC-05 When an invite is redeemed, then the new account's email field is
    pre-filled with the invited address.
