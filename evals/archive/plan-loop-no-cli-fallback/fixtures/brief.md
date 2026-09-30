# Brief: `widget export --dry-run`

Fictional brief, used as fixture data only, for
`plan-loop-no-cli-fallback`. Small on purpose: this case measures the
`/dh:plan-loop` no-CLI fallback routing, not planning depth.

## Scope

- RF-01 `widget export` accepts a `--dry-run` flag.
  - AC-01 When `--dry-run` is passed, then no output file is written and
    the command reports what it would have written instead.
- RF-02 Without `--dry-run`, behavior is unchanged.
  - AC-02 When `--dry-run` is absent, then `widget export` writes the file
    exactly as it does today.

## Direct file evidence

The existing command is `fixtures/widget-export.md` (a fixture stand-in
for the real source file — read directly instead of dispatching a
repository-map task, since this one small file is the whole surface this
brief touches).
