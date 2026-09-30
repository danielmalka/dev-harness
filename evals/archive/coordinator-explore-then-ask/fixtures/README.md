# receipts-api

A tiny internal service that stores purchase receipts so support agents can
look them up while helping a customer. Run the test suite with the script
declared in `package.json`.

## Layout

- `src/receipts.js` — the receipt lookup and creation handlers.

## Receipt record

- `customerEmail` — the customer's email address; this is where a receipt
  gets sent when a support agent emails one.
- `scanImagePath` — set when a paper receipt was scanned in: a path under
  this service's existing `scans/` storage. Can be a large image.

## Auth

Every route in this service already runs behind the app's existing
support-agent auth middleware. A new route needs no new authentication or
authorization work.

## Outbound email

Outbound email already goes through the app's existing shared mail
client; a new route reuses it, so sending a receipt needs no new mail
provider or dependency.

## Background work

This service has no job queue, worker process, or message broker of any
kind.
