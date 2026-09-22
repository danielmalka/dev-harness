# codex

Binary: `codex`. Reviewer. Hands: yes, via `codex exec`, unless run under
`-s read-only` (see below).

## Headless

```bash
codex exec --skip-git-repo-check -s read-only -m <slug> - < /tmp/prompt.md
```

`exec` is the non-interactive path. Prompt on stdin when using `-`. If
stdin is piped **and** a prompt argv is set, stdin is appended as a
`<stdin>` block. `-m` / `--model` pins the slug from `models.codex` in
`.harness/local.yaml`.

| Flag | Why |
|---|---|
| `exec` | one-shot; bare `codex` is the TUI |
| `-` | read prompt from stdin |
| `-m` / `--model` | pin the slug |
| `--skip-git-repo-check` | required outside a trusted git repo; otherwise exit 0 with empty stdout |
| `-C` / `--cd` | workspace root |
| `--add-dir` | extra writable dirs |
| `-s` / `--sandbox` | `read-only` / `workspace-write` / `danger-full-access` |
| `--ephemeral` | do not persist the session |

`codex review` (also `codex exec review`) is the dedicated review path;
prefer it when the job is only a verdict on the current tree:

```bash
codex review --uncommitted -c model="<slug>"
```

`--uncommitted` covers staged, unstaged and untracked changes. `--base
BRANCH` reviews against that branch; `--commit SHA` reviews that commit.
`-m`/`-C` are not valid on the `review` subcommand, and `--uncommitted`
does not combine with a `[PROMPT]` / `-` argument.

**Sandbox note (read-only mode)**: `codex exec -s read-only` runs with
`/tmp` mounted read-only, so a Go-toolchain `go test` and similar
compile-in-tmp checks fail inside it. Treat the call as read/reasoning
only, not behavioral validation — run the project's own test/lint gate
on the host if a real result is needed. This is the reason the parent
skill's anti-delegation clause treats the `commands.test`/`commands.lint`
permission as inert in practice for this binary under `-s read-only`.
`-s workspace-write` also blocks outbound network from the sandbox, so a
suite that needs a live service will not run there either.

Do not start the interactive CLI (`codex` with no subcommand) from an
orchestrator turn. Codex has no `models` subcommand — the slug list lives
in `.harness/local.yaml`, refreshed by hand when the operator measures a
new one.
