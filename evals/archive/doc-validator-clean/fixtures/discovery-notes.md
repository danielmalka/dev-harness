# Discovery notes: invite-only signup

Session with the workspace owner, fictional product, used as a fixture only.

1. Signup is closed. A person creates an account only by redeeming an invite code; a direct signup attempt without a code is refused.
2. Only a workspace admin creates an invite, and each invite is bound to one email address. A non-admin member who tries to create one is refused.
3. An invite expires 7 days after it is created. Redeeming an expired invite is refused and the person is told the invite expired.
4. An admin revokes a pending invite before it is redeemed. A revoked invite can no longer be redeemed.
5. The invite email is sent once when the invite is created. The admin can resend it at most 3 times per invite; the fourth resend is refused.

Not discussed: pricing, seat limits, single sign-on, referral rewards.
