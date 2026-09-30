# Transport double: cli:codex/gpt-5.6-sol, one HIGH finding

Pre-written reply, used as fixture data only. No real `codex` binary is
invoked to produce this file or to grade any case that reads it. Represents
a reply that carries a real verdict signal (a `### <SEVERITY>` entry under
`## Findings`) with one blocking finding. Follows the Output format of
`.skills/security-review/SKILL.md`.

```
## Scope
- Change reviewed: src/api/invoices.ts, src/repo/invoices.ts
- Boundaries touched: authenticated HTTP route, invoice records

## Findings
### HIGH: Invoice lookup missing object-level authorization
- Category: demonstrated vulnerability
- Confidence: traced directly in the diff; no mitigating check found on
  this path
- Source: src/api/invoices.ts:20 - req.params.id from the request path,
  any authenticated caller
- Path: router handler -> repository.findById -> response serializer
- Sink: src/repo/invoices.ts:58 - select by id with no account filter
- Mitigation present: none
- Attacker and gain: any registered user obtains invoices belonging to
  other accounts
- Scenario: authenticate as account A, request an invoice id owned by
  account B, receive it
- Mitigation: filter by the session account in the repository query, not
  in the handler
- Verification: test asserting a 404 when the invoice belongs to another
  account
- Reported by: cli:codex/gpt-5.6-sol

## Hypotheses requiring investigation
(none)

## Non-issues
(none)

## Checks
- go test ./... - not-run (sandbox)

## Limitations
- Fixture double; no real codex binary invoked to produce this file.
```
