# grok (Grok Build)

Binary: `grok`, resolved via `.harness/local.yaml` `binaries.grok` when `command -v grok` fails. A login-shell-only `PATH` is common for this binary: `command -v` misses it even when it is installed, so do not treat that miss alone as `not_installed`. Coding agent with hands; also a reviewer and a planner (in a disposable clone).

## Invocation

```bash
grok --always-approve --cwd "$CLONE" --prompt-file "$PROMPT"
# short prompt:
grok --always-approve --cwd "$CLONE" -p "$(cat "$PROMPT")"
```

One shell command whose first word is the binary: write the resolved command there as a literal (bare `grok`, or the literal path taken from `.harness/local.yaml`; a literal path does not match a `Bash(grok:*)` allow rule, so that form may need the owner's approval or be not-run), with no `timeout`, `&&`, `;` or pipe in front of or after it. `--cwd <CWD>` is the binary's own working-directory flag, so no `cd` is needed. `-p` / `--single` prints to stdout and exits. `--always-approve` auto-approves tools, needed for an unattended call; `--permission-mode bypassPermissions` is the named equivalent.

| Flag | Why |
|---|---|
| `-p` / `--single` | one-shot prompt on argv |
| `--prompt-file` | one-shot prompt from disk |
| `-m` / `--model` | pin the slug |
| `--always-approve` | unattended tools |
| `--cwd` | workspace; `--cwd "$CLONE"` on every call |
| `--reasoning-effort` / `--effort` | `none` / `minimal` / `low` / `medium` / `high` / `xhigh` / `max` (only the levels the model advertises) |
| `--output-format` | `plain` (default) / `json` / `streaming-json` / `streaming-messages-json` |
| `--json-schema` | constrains JSON output (implies `--output-format json`) |
| `--max-turns` | cap agentic turns |
| `--no-subagents` | block spawning a subagent from inside the call |

Do not launch the interactive TUI from an orchestrator turn.

## Prompt input

`--prompt-file` is the contract for anything longer than a line; `-p` takes a short prompt on argv. `$CLONE` and `$PROMPT` are defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

No read-only flag documented: disposable clone ([clone-isolation.md](clone-isolation.md)), compared as [live-tree-freeze.md](live-tree-freeze.md) says. `--no-subagents` narrows but does not replace the anti-delegation clause: pass both.

## Model flag

`-m` / `--model <slug>`, from `models.grok` in `.harness/local.yaml`.

## Timeout

No documented wall-clock flag besides `--max-turns`. The Bash tool timeout (set to the reviewer's timeout) or `run_in_background` bounds the time.

## Output capture

Reply on stdout; `--output-format` selects the shape (default `plain`).

## Failure signals

Large diffs with `--prompt-file` at high effort have been seen to time out or to print only a reasoning preamble with exit 0.

## Known gotchas

See [gotchas.md](gotchas.md), kinds `timed_out` and `no_signal`.

## Measured on

2026-10-16 (`grok --help`: `--cwd`).
