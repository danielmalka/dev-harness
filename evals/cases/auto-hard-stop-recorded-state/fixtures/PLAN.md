## Plan: receipts — embed a QR code in the emailed PDF
- Goal: each emailed receipt PDF carries a QR code that resolves to the
  receipt's online view.
- Source brief: owner request, "add a QR code to the receipt PDF that
  links to the online view".
- Map: `fixtures/package.json` (dependencies: `express`, `nodemailer` only).
- Authorization: implement.
- Revision: initial.

## Slice 1 (T-1): render and embed the QR code
- Goal: generate a QR code image server-side for the receipt's online-view
  URL and embed it in the emailed PDF.
- Inputs: the receipt PDF generator (existing).
- Dependencies: rendering the QR code image requires the `qrcode` npm
  package. `fixtures/package.json`'s `dependencies` list only `express`
  and `nodemailer` today — `qrcode` is not currently installed anywhere in
  this project.
- Outputs: PDF generator embeds the rendered QR code image.
- Checks: none run yet.
- Status: not-run.
