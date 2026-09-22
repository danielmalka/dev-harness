# Gotchas (append-only)

Generalized transport-failure patterns observed invoking these CLIs as
reviewers. One entry per miss: what happened, what to do instead. No
machine path and no consumer-project name belongs here — that is exactly
what made the source these were ported from unportable; see
`SKILL.md` metadata for provenance.

- A CLI running as an implementer under a full-access sandbox flag
  invoked a second CLI binary on its own, inside the same checkout,
  delegating part of the work without being asked. · The anti-delegation
  clause must be copied verbatim into the body of every reviewer prompt,
  never left only in the dispatch instruction around it (RISK-001); when
  delegation is suspected, check running processes during the call.
- A CLI given an explicitly read-only, validation-only brief ran a
  repository build script on its own and rewrote a versioned output
  directory (the rebuilt result happened to match the sources this time,
  so no damage occurred, but the run was not read-only in practice). · A
  read-only brief must forbid running repository scripts, not only
  forbid editing files; check the modification times of any versioned
  output directory after the call, not just `git status` on tracked
  files.
- `codex exec -s read-only` blocks a Go-toolchain `commands.test` from
  actually compiling. · Treat that reviewer's pass as reading and
  reasoning about the code, not as behavioral validation; see
  [codex.md](codex.md) for the mechanism.
- A CLI review call exited 0 after printing only a reasoning preamble,
  with no verdict block, on a large diff run at high reasoning effort. ·
  Zero-exit with no verdict signal is a transport failure exactly like a
  non-zero exit or a timeout; it is never read as an implicit pass.
- Two concurrent calls to the same CLI binary, from the same account,
  have been seen to fail with a database-lock error on the second one. ·
  Serialize calls to a binary known to reject concurrent runs into their
  own batch instead of dispatching them alongside each other; see the
  per-binary reference for which binaries this applies to.
