# Receipts service (fixture stub)

A tiny stub of this checkout's existing receipts service, for fixture
purposes only — not a real implementation.

- `send`: emails a receipt PDF to the customer after a purchase.
- `bounce`: records a bounce reported back by the mail provider for a
  previously sent receipt.
- `attachment`: a large receipt (one with scanned images) is attached to
  the email as a separate PDF, alongside the receipt itself.

Today, resending a bounced receipt always resends the full email exactly
as it was first sent, including the large scanned attachment when one was
present.
