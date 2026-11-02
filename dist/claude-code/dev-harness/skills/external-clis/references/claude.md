# claude (Claude Code)

Binary: `claude`. Usable as a CLI reviewer, distinct from a Coordinator dispatching a Claude agent in-process: this path shells out to a separate `claude` process. Not a planner: `cli:claude/<slug>` stays excluded from `/dh:plan-loop` while the kit runs on Claude.

## Invocation

```bash
cd "$CLONE" && claude -p --dangerously-skip-permissions --model <slug> --effort high \
  "$(cat "$PROMPT")"
```

`claude --help` lists no working-directory flag (only `--add-dir`, which adds roots and does not set the cwd), so this recipe keeps `cd "$CLONE" && claude ...` as the first of two named exceptions (with `opencode`, [opencode.md](opencode.md)) to the rule that a call is one command starting with the binary name. A Bash call does not keep its cwd for the next call, so the `cd` stays inside the same command. That form does not match a `Bash(claude:*)` allow rule: the call may need the owner's approval, or is not-run if it is refused. No `timeout` wrapper goes before it. The call never passes `--add-dir` for the live tree. `-p` / `--print` prints and exits. `--dangerously-skip-permissions` bypasses the permission prompt, needed for an unattended call; `--permission-mode bypassPermissions` is the named equivalent.

| Flag | Why |
|---|---|
| `-p` / `--print` | one-shot; bare `claude` is the TUI |
| `--dangerously-skip-permissions` | unattended |
| `--model` | an alias (`haiku`, `sonnet`, `opus`, `fable`, …) or a full model id |
| `--effort` | `low` / `medium` / `high` / `xhigh` / `max` |
| `--fallback-model` | pin a fallback if the default is unavailable |
| `--add-dir` | extra workspace roots (never the live tree in a clone call) |
| `--output-format` | print-only: `text` / `json` / `stream-json` |
| `--input-format` | print-only: `text` or `stream-json` |

Do not start the interactive TUI from an orchestrator turn.

## Prompt input

Last argv, or expanded from a file as above. `$CLONE` and `$PROMPT` are defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

No documented read-only flag as a CLI process: disposable clone ([clone-isolation.md](clone-isolation.md)), compared as [live-tree-freeze.md](live-tree-freeze.md) says.

## Model flag

`--model <slug>`: an alias or a full model id.

## Timeout

No wall-clock flag. The Bash tool timeout or `run_in_background` bounds the time.

## Output capture

Reply on stdout; `--output-format` selects `text`, `json` or `stream-json`.

## Failure signals

A long run that never emits the stage's verdict signal is `no_signal`.

## Known gotchas

None recorded in [gotchas.md](gotchas.md). The verdict format is the one in `SKILL.md`, "The verdict contract"; this recipe defines no other.

## Measured on

2026-10-16 (`claude --help`: no cwd flag). The call itself: unknown (not measured on this kit).
