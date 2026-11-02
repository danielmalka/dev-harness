# opencode

Binary: `opencode`, resolved via `.harness/local.yaml` `binaries.opencode` when `command -v opencode` fails. A login-shell-only `PATH` is common for this binary: `command -v` misses it even when it is installed, so do not treat that miss alone as `not_installed`. Reviewer or builder, and planner (in a disposable clone); hands: yes (`opencode run`).

## Invocation

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>#high' \
  "$(cat "$PROMPT")"
```

Long prompt as an attached file (do not stuff a multi-kilobyte prompt into argv):

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>[#<variant>]' \
  --file="$PROMPT" -- "Follow the attached prompt."
```

`opencode run` has no working-directory flag and no `--variant`. The call therefore starts with `cd "$CLONE" &&`, the second named exception (beside `claude`) to the rule that a call is one command whose first word is the binary. A separate `cd` in its own Bash call does not persist (the runtime resets the cwd), so the `cd` stays in the same command. Because the command starts with `cd`, it does not match a `Bash(opencode:*)` allow rule: the call may need the owner's approval or be not-run. Write the resolved binary as a literal (bare `opencode`, or the literal path taken from `.harness/local.yaml`). No `timeout` goes in front of the command and no `;` or pipe chains it; the freeze and the later comparisons are separate commands.

| Flag | Why |
|---|---|
| `run` | one-shot; bare `opencode` is the TUI |
| `--auto` | approve tools (required for an unattended call) |
| `-m provider/model[#variant]` | pin; do not rely on the interactive picker; `#variant` sets the reasoning effort (only levels that slug advertises); there is no `--variant` flag |
| `cd "$CLONE" &&` | workspace; there is no `--dir`, so every call starts with it |
| `-f` / `--file` | attach file(s) to the message (repeatable) |
| `--format json` | raw events if parsing is required |
| `--thinking` | print thinking blocks |
| `--title` | session title |
| `-c` / `--continue`, `-s` / `--session`, `--fork` | resume / fork; default is a new session |

Refuse from an orchestrator turn: bare `opencode`, `opencode --mini`, `run --interactive` / `-i`. `--agent` on `run` switches the primary agent, not an in-process subagent (`task` inside opencode): the anti-delegation clause still applies to what that primary agent may do.

## Prompt input

An attached file: `--file="$PROMPT" -- "message"`. `-f`/`--file` is a yargs array: a bare `-f file.md "msg"` swallows `msg` as a second file, and `--file=` alone with no `--` message is refused. Always `--file=path -- "message"`. `$CLONE` and `$PROMPT` are defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

No read-only flag documented: disposable clone ([clone-isolation.md](clone-isolation.md)), compared as [live-tree-freeze.md](live-tree-freeze.md) says.

## Model flag

`-m '<provider>/<slug>[#<variant>]'`, from `models.opencode` in `.harness/local.yaml`. The variant (reasoning effort) is part of the model string, only a level that slug advertises; without `#variant` the model default applies. The model catalog is provider-specific and changes over time: do not invent a slug here; refresh it by running the provider's own model-listing command on this machine and record the result in `.harness/local.yaml`, never in this file.

## Timeout

`opencode run` has no wall-clock flag. The Bash tool timeout or `run_in_background` bounds the time.

## Output capture

Reply on stdout; `--format json` gives raw events if parsing is required.

## Failure signals

- A database lock error on the second of two concurrent runs against the same account.
- A refused `--file=` call with no `--` message.

## Known gotchas

Do not run two `opencode run` calls concurrently. When a review stage lists more than one `opencode` entry, or the concurrency ceiling would otherwise be reached, the Coordinator serializes `opencode` calls into their own batch rather than running them alongside each other. See also [gotchas.md](gotchas.md).

## Measured on

2026-10-21, v2.0.32 (no cwd flag, no `--variant`; the runtime resets the cwd between Bash calls).
