# codex

Binary: `codex`. Reviewer; planner (documented, measured read-only mode). Hands: yes, via `codex exec`, unless run under `-s read-only`.

## Invocation

```bash
codex exec --skip-git-repo-check -s read-only -m <slug> - < "$PROMPT"
```

One shell command whose first word is `codex`, with the stdin redirection part of it: no `timeout`, `&&`, `;` or pipe in front of or after it. `exec` is the non-interactive path; bare `codex` is the TUI and is never started from an orchestrator turn.

| Flag | Why |
|---|---|
| `exec` | one-shot |
| `-` | read prompt from stdin |
| `-m` / `--model` | pin the slug |
| `--skip-git-repo-check` | required outside a trusted git repo; otherwise exit 0 with empty stdout |
| `-C` / `--cd` | workspace root |
| `--add-dir` | extra writable dirs |
| `-s` / `--sandbox` | `read-only` / `workspace-write` / `danger-full-access` |
| `--ephemeral` | do not persist the session |

`codex review` (also `codex exec review`) is the dedicated review path; prefer it when the job is only a verdict on the current tree:

```bash
codex review --uncommitted -c model="<slug>"
```

`--uncommitted` covers staged, unstaged and untracked changes. `--base BRANCH` reviews against that branch; `--commit SHA` reviews that commit. `-m`/`-C` are not valid on the `review` subcommand, and `--uncommitted` does not combine with a `[PROMPT]` / `-` argument.

## Prompt input

Stdin (`-`). If stdin is piped and a prompt argv is set, stdin is appended as a `<stdin>` block. `$PROMPT` is defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

Read-only flag: `-s read-only`. No clone; the live tree is compared as [live-tree-compare.md](live-tree-compare.md) says. Under `-s read-only`, `/tmp` is mounted read-only, so a Go-toolchain `go test` and similar compile-in-tmp checks fail inside it. Treat the call as read and reasoning only, not behavioral validation: run the project's own test and lint gate on the host if a real result is needed. This is why the `commands.test`/`commands.lint` permission in the anti-delegation clause is inert in practice for this binary. `-s workspace-write` also blocks outbound network from the sandbox, so a suite that needs a live service will not run there either.

## Model flag

`-m` / `--model <slug>`, the slug from `models.codex` in `.harness/local.yaml`. Codex has no `models` subcommand: the slug list is refreshed by hand when the operator measures a new one.

## Timeout

`codex exec` has no wall-clock flag. The Bash tool timeout (set to the reviewer's timeout) or `run_in_background` bounds the time.

## Output capture

Reply on stdout.

## Failure signals

- Exit 0 with empty stdout: `--skip-git-repo-check` missing outside a trusted git repo.
- A `go test` that cannot compile under `-s read-only` (see Isolation).

## Known gotchas

See [gotchas.md](gotchas.md), kinds `no_signal` and `tree_changed`.

## Measured on

2026-09-22 (verdict formatting run). The date of the read-only probe that qualifies `codex` as a planner: unknown.
