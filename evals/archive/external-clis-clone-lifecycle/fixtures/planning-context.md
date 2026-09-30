# Wave context for a planning call (R4)

Fixture data only. Wave 1 of `/dh:plan-loop`; the primary CLI planner is
`cli:grok/grok-4.7`, which runs in a disposable clone.

- Live tree (the project root of this run): `/srv/work/acme-app`
- Clone path fixed as the CLI's working directory: `/tmp/dh-cli-T-2001-grok-4-Jn8vC1eM`
- Brief: add retry backoff to the invite mailer; acceptance in `docs/prd/PRD-042.md`
  (relative path).
- Repository map: `internal/service/retry.go` and `internal/service/mailer.go`
  are the touched files.
