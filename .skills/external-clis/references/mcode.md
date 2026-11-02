# mcode (MiniMax Code)

Binary: `mcode`, resolved via `.harness/local.yaml` `binaries.mcode` when `command -v mcode` fails. A login-shell-only `PATH` is common for this binary: `command -v` misses it even when it is installed, so do not treat that miss alone as `not_installed`. Coding agent with hands; reviewer and planner (in a disposable clone).

## Invocation

```bash
mcode exec --permission off --cwd "$CLONE" --timeout 15m \
  --model <provider>/<slug> \
  --input - < "$PROMPT"
```

One shell command whose first word is the binary: write the resolved command there as a literal (bare `mcode`, or the literal path taken from `.harness/local.yaml`; a literal path does not match a `Bash(mcode:*)` allow rule, so that form may need the owner's approval or be not-run), with no `timeout`, `&&`, `;` or pipe in front of or after it. `--cwd <path>` is the binary's own working-directory flag ("workspace directory"). The `< "$PROMPT"` redirection is part of the same command. `--permission full` is the write-capable unattended mode, for a builder role only: never pass it from this skill, which covers review stages.

| Flag | Why |
|---|---|
| `exec` | one-shot; bare `mcode` is the TUI |
| `--permission off` | documented, unmeasured: read-only, unattended (use for a reviewer, in a clone) |
| `--permission full` | unattended with write access (builder roles only, out of this skill's scope) |
| `--cwd` | workspace; `--cwd "$CLONE"` on every call |
| `--model provider/model` | pin the slug |
| `--input -` | read prompt from stdin (`--input-format text` / `json`) |
| `--file PATH` | attach a file (repeatable) |
| `--timeout 15m` | a review call often needs longer than the default |
| `--max-steps` | cap assistant steps |
| `--output-format` | `text` / `json` / `stream-json` |
| `--output-schema` | JSON Schema for the final answer |
| `-o` / `--output-last-message PATH` | final agent message on disk |

Do not start the interactive TUI from an orchestrator turn.

## Prompt input

Stdin (`--input -`) or `--file`. `$CLONE` and `$PROMPT` are defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

Read-only flag documented, not yet measured on this kit: `--permission off` disables tool approval entirely rather than auto-approving it, which is the binary's own read-only mode. No probe is recorded confirming it refuses a file-create and a shell-command attempt, so until a dated probe is on record `mcode` runs in a disposable clone ([clone-isolation.md](clone-isolation.md)) for review and planning alike. The flag is used in addition to, never instead of, the comparison every clone call gets ([live-tree-freeze.md](live-tree-freeze.md)).

## Model flag

`--model <provider>/<slug>`, from `models.mcode` in `.harness/local.yaml`, refreshed by hand when the operator measures a new slug for this account.

## Timeout

`--timeout 15m` bounds the time.

## Output capture

Reply on stdout; `-o` / `--output-last-message PATH` writes the final agent message to disk; `--output-format` selects `text`, `json` or `stream-json`.

## Failure signals

unknown

## Known gotchas

None recorded in [gotchas.md](gotchas.md).

## Measured on

2026-10-16 (`mcode exec --help`: `--cwd`). `--permission off`: unknown (not measured).
