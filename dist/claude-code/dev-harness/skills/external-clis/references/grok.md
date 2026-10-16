# grok (Grok Build)

Binary: `grok`, resolved via `.harness/local.yaml` `binaries.grok` when
`command -v grok` fails (a login-shell-only `PATH` is common for this
binary; `command -v` misses it even when it is installed — do not treat
that miss alone as "not installed"). Coding agent with hands; also a
reviewer.

## Headless

```bash
grok --always-approve --cwd "$CLONE" --prompt-file "$PROMPT"
# short prompt:
grok --always-approve --cwd "$CLONE" -p "$(cat "$PROMPT")"
```

Each call is one shell command whose first word is the binary
(R18): write the resolved command there as a literal (bare `grok`, or the literal path taken from `.harness/local.yaml`; a literal path does not match a `Bash(grok:*)` allow rule, so that form may need the owner's approval or be not-run), with no `timeout`, `&&`, `;` or pipe in front of or after it. `--cwd <CWD>` is the binary's own working-directory flag (fact, `grok --help`, 2026-10-16), so no `cd` is needed. Time is bounded by the Bash tool timeout (set it to the reviewer's timeout) or `run_in_background`; `grok` has no documented wall-clock flag besides `--max-turns`.

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
| `--cwd` | workspace; `--cwd "$CLONE"` on every call |
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
