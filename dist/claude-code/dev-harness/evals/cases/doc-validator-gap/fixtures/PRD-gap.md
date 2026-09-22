# PRD-902 · Invite-only signup

| Field | Value |
|---|---|
| Status | in review |
| Owner | workspace owner |
| Created / updated | 2026-09-20 / 2026-09-20 |
| Stories | ST-901, ST-902 |

## 1. Problem

Signup is open, so anyone who finds the URL creates an account in a workspace they were never meant to join. The workspace owner currently removes those accounts by hand after the fact.

## 2. Expected outcome

- Signal 1: accounts created without an invite · today: unbounded · target: 0.
- Signal 2: pending invites that expire unnoticed · today: no baseline · target: measured after the first month.

## 3. Users and scenarios

| User | Scenario | Today | After |
|---|---|---|---|
| Admin | Brings a colleague into the workspace | Sends them the signup URL | Creates an invite bound to their email |
| Invited person | Creates an account | Signs up directly | Redeems the invite code |

## 4. Scope

**In scope**
- RF-01 A person creates an account only by redeeming a valid invite code; a signup attempt without a code is refused.
- RF-02 Only a workspace admin creates an invite, and the invite is bound to exactly one email address.
- RF-03 An invite expires 7 days after creation and can no longer be redeemed.
- RF-04 The invite email is sent once on creation, and the admin resends it at most 3 times per invite.

**Out of scope**
- Pricing and seat limits.
- Single sign-on.
- Referral rewards.

## 5. Acceptance criteria

- AC-01 (RF-01) When a person submits the signup form with no invite code, then the account is not created and the response states that an invite is required.
- AC-02 (RF-02) When a member who is not an admin submits an invite, then the invite is not created and the response states that only admins may invite. When an admin submits an invite for an email address, then the invite is bound to that address and no other address can redeem it.
- AC-03 (RF-03) When an invite is redeemed more than 7 days after its creation timestamp, then the redemption is refused and the response states that the invite expired.
- AC-04 (RF-04) When the admin uses the invite form, then the form must be fast.

## 6. Constraints and technical impact

- Data: new invite record with email, code, creation timestamp, state and resend count.
- Contracts: the signup endpoint gains a required invite code field.
- Security: the invite code is single use and bound to one email address.

## 7. Discarded alternatives

| Alternative | Why not |
|---|---|
| Shared workspace password | Cannot be bound to one person or revoked individually |

## 8. Risks and pending decisions

- Risk: an admin revokes the wrong invite · mitigation: the revoked invite is shown in the invite list with its state.

## 9. References

- discovery-notes.md
