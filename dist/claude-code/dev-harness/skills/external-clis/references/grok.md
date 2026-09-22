# grok (Grok Build)

Binary: `grok`, resolved via `.harness/local.yaml` `binaries.grok` when
`command -v grok` fails (a login-shell-only `PATH` is common for this
binary; `command -v` misses it even when it is installed — do not treat
that miss alone as "not installed"). Coding agent with hands; also a
reviewer.

## Headless

```bash
GROK=<resolved path or bare "grok">
"$GROK" --always-approve --prompt-file /tmp/prompt.md
# short prompt:
"$GROK" --always-approve -p "$(cat /tmp/prompt.md)"
```

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
| `--cwd` | workspace |
| `--reasoning-effort` / `--effort` | `none` / `minimal` / `low` / `medium` / `high` / `xhigh` / `max` (only the levels the model advertises) |
| `--output-format` | `plain` (default) / `json` / `streaming-json` / `streaming-messages-json` |
| `--json-schema` | constrains JSON output (implies `--output-format json`) |
| `--max-turns` | cap agentic turns |
| `--no-subagents` | block spawning a subagent from inside the call |

Large diffs with `--prompt-file` at high effort have been seen to time
out or print only a reasoning preamble with exit 0 — that is a transport
failure, not a lenient pass; see `references/gotchas.md`.

`grok` has no documented read-only flag; the guarantee that it did not
edit the workspace comes from the post-call `git status` check the
parent skill describes (SKILL.md step 7). `--no-subagents` narrows but
does not replace the anti-delegation clause — pass both.

Do not launch the interactive TUI from an orchestrator turn.
