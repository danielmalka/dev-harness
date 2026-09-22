# mcode (MiniMax Code)

Binary: `mcode`, resolved via `.harness/local.yaml` `binaries.mcode` when
`command -v mcode` fails (a login-shell-only `PATH` is common for this
binary — do not treat that miss alone as "not installed"). Coding agent
with hands.

## Headless

```bash
MCODE=<resolved path or bare "mcode">
"$MCODE" exec --permission off --cwd "$PWD" \
  --model <provider>/<slug> \
  --input - < /tmp/prompt.md
```

Read-only, unattended review call: `--permission off` disables tool
approval entirely rather than auto-approving it, which is the binary's
own read-only mode. It is used **in addition to**, not instead of, the
unconditional post-call `git status` check every reviewer gets
(SKILL.md step 7) — that check is never skipped just because a binary
has its own read-only flag. `--permission full` is the write-capable
unattended mode, for a builder role only — never pass it from this
skill, which covers review stages.

| Flag | Why |
|---|---|
| `exec` | one-shot; bare `mcode` is the TUI |
| `--permission off` | read-only, unattended (use for a reviewer) |
| `--permission full` | unattended with write access (builder roles only, out of this skill's scope) |
| `--cwd` | workspace |
| `--model provider/model` | pin the slug from `models.mcode` in `.harness/local.yaml` |
| `--input -` | read prompt from stdin (`--input-format text` / `json`) |
| `--file PATH` | attach a file (repeatable) |
| `--timeout 15m` | a review call often needs longer than the default |
| `--max-steps` | cap assistant steps |
| `--output-format` | `text` / `json` / `stream-json` |
| `--output-schema` | JSON Schema for the final answer |
| `-o` / `--output-last-message PATH` | final agent message on disk |

Slugs: `models.mcode` in `.harness/local.yaml`, refreshed by hand when
the operator measures a new one for this account.

Do not start the interactive TUI from an orchestrator turn.
