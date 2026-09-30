# grok (Grok Build)

Binary: `grok`, resolved via `.harness/local.yaml` `binaries.grok` when
`command -v grok` fails (a login-shell-only `PATH` is common for this
binary; `command -v` misses it even when it is installed — do not treat
that miss alone as "not installed"). Coding agent with hands; also a
reviewer.

## Headless

```bash
GROK=<resolved path or bare "grok">
cd "$CLONE" && "$GROK" --always-approve --cwd "$CLONE" --prompt-file "$PROMPT"
# short prompt:
cd "$CLONE" && "$GROK" --always-approve --cwd "$CLONE" -p "$(cat "$PROMPT")"
```

`$CLONE` in these recipes is the literal clone path stored when the clone was created ([clone-isolation.md](clone-isolation.md), Terms), written into the command as that literal string. The prompt file `$PROMPT` is the literal path of a file created by `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )` inside the call's freeze directory ([live-tree-freeze.md](live-tree-freeze.md), section 5); it is unique per call and is removed with that directory.

`-p` / `--single` prints to stdout and exits. `--prompt-file` is the
contract for anything longer than a line. `-m` / `--model` pins the slug
from `models.grok` in `.harness/local.yaml`. `--always-approve`
auto-approves tools, needed for an unattended call.
`--permission-mode bypassPermissions` is the named equivalent.

| Flag | Why |
|---|---|
| `-p` / `--single` | one-shot prompt on argv |
| `--prompt-file` | one-shot prompt from disk |
| `-m` / `--model` | pin the slug |
| `--always-approve` | unattended tools |
| `--cwd` | workspace; `cd "$CLONE" &&` first, then `--cwd "$CLONE"`, for every call |
| `--reasoning-effort` / `--effort` | `none` / `minimal` / `low` / `medium` / `high` / `xhigh` / `max` (only the levels the model advertises) |
| `--output-format` | `plain` (default) / `json` / `streaming-json` / `streaming-messages-json` |
| `--json-schema` | constrains JSON output (implies `--output-format json`) |
| `--max-turns` | cap agentic turns |
| `--no-subagents` | block spawning a subagent from inside the call |

Large diffs with `--prompt-file` at high effort have been seen to time
out or print only a reasoning preamble with exit 0 — that is a transport
failure, not a lenient pass; see `references/gotchas.md`.

`grok` has no documented read-only flag, so it runs in a clone: disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected. The live tree is compared as `references/live-tree-freeze.md` describes (SKILL.md step 7 covers only a call without a clone). `--no-subagents` narrows but
does not replace the anti-delegation clause — pass both.

Do not launch the interactive TUI from an orchestrator turn.
