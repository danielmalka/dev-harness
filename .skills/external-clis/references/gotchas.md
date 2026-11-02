# Gotchas (append-only)

Generalized failure patterns observed invoking these CLIs as reviewers, grouped by the failure kinds of `SKILL.md`, "Failure classification". One entry per miss: what happened, what to do instead. Append a new entry under its kind; never rewrite or reorder an old one. No machine path and no consumer-project name belongs here; see `SKILL.md` metadata for provenance.

## not_installed

- `command -v <binary>` misses a binary that is installed when its directory is only on a login-shell `PATH`. · Look up `binaries.<binary>` in `.harness/local.yaml` before classifying; a `command -v` miss alone is not `not_installed`.

## auth_required

- None recorded.

## rate_limited

- None recorded.

## timed_out

- An `agy --print` call wrote its full reply in about 2 minutes, but the process stayed alive past `--print-timeout 15m` and needed SIGKILL (measured 2026-10-16). · A complete reply on disk that carries the stage's verdict signal is the result; stop the process by the pid recorded at the start of the call; see [agy.md](agy.md).
- A large diff sent with `--prompt-file` at high reasoning effort timed out. · It is a failure, never a lenient pass.

## no_signal

- A CLI review call exited 0 after printing only a reasoning preamble, with no verdict block, on a large diff run at high reasoning effort. · Exit 0 with no verdict signal is a failure exactly like a non-zero exit or a timeout; it is never read as an implicit pass.
- A real reviewer run returned a correct blocking verdict written as `- **Request changes**` and no `Review status:` line (2026-09-22). · State the verdict signal in the prompt as a requirement, verbatim, as the verdict contract in `SKILL.md` says.
- The same reviewer runs signed their findings with tags that were not the dispatched slug. · Tell the reviewer its exact `Reported by:` string in the prompt.
- `agy` printed its `--add-dir` help and ignored the packet when `--add-dir PATH` came before the prompt. · Inline the files in the prompt; see [agy.md](agy.md).
- `codex exec` outside a trusted git repo exited 0 with empty stdout. · Pass `--skip-git-repo-check`; see [codex.md](codex.md).
- Two concurrent calls to the same CLI binary, from the same account, failed on the second with a database-lock error. · Serialize calls to a binary known to reject concurrent runs into their own batch; see the binary's recipe for which binaries this applies to.

## tree_changed

- A CLI given an explicitly read-only, validation-only brief ran a repository build script on its own and rewrote a versioned output directory (the rebuilt result happened to match the sources, so no damage occurred, but the run was not read-only in practice). · A read-only brief must forbid running repository scripts, not only editing files; check the modification times of any versioned output directory after the call, not just `git status` on tracked files.
- A `dh build` run by the Coordinator while two CLI reviewers were working made the post-call comparison differ on the manifest and `dist/` (2026-09-22). · Nothing else writes to the tree while a call is in flight; see [live-tree-compare.md](live-tree-compare.md).

## isolation_unavailable

- None recorded.

## No kind (the call ran; the prompt or the scheduling was wrong)

- A CLI running as an implementer under a full-access sandbox flag invoked a second CLI binary on its own, inside the same checkout, delegating part of the work without being asked. · The anti-delegation clause must be copied verbatim into the body of every reviewer prompt, never left only in the dispatch instruction around it (RISK-001); when delegation is suspected, check running processes during the call.
- `codex exec -s read-only` blocks a Go-toolchain `commands.test` from actually compiling. · Treat that reviewer's pass as reading and reasoning about the code, not as behavioral validation; see [codex.md](codex.md).
- `pgrep -f` or `pkill -f` with a pattern that contains the delegate's own command or path matched the delegate's own shell (exit 144, endless waits; seen 2026-10-21). · Kill by the pid recorded when the call started, or use a bracket pattern that does not match its own text (`'[o]pencode run'`).
