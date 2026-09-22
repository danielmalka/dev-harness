# PRD: teammate invites

Fictional PRD, used as a fixture only. Written from fixtures/source.md, but
deliberately drops source point 3 (revoke): this is the gap the CLI
reviewer's fixture double reports, while the claude reviewer's fixture
double misses it, to exercise the merge rule in
`.agents/coordinator.md` ("External CLI reviewers").

## Scope

- RF-01 An admin invites a new teammate by email (source point 1).
  - AC-01 When an admin submits an email, then a pending invite is created
    for it.
- RF-02 An invite expires after 7 days (source point 2).
  - AC-02 When 7 days pass without redemption, then the invite can no
    longer be redeemed.
- RF-03 Only admins can view the list of pending invites (source point 4).
  - AC-03 When a non-admin requests the pending-invites list, then the
    request is refused.
- RF-04 A redeemed invite creates the teammate's account with the invited
  email pre-filled (source point 5).
  - AC-04 When an invite is redeemed, then the new account's email field is
    pre-filled with the invited address.
