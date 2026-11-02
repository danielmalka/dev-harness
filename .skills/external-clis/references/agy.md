# agy (Antigravity)

Binary: `agy`. Reviewer among other roles. No planning until a dated probe records that `agy` treats the process cwd as its workspace and which flag, if any, sets it.

## Invocation

```bash
agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat "$PROMPT")"
```

One shell command whose first word is `agy`: no `timeout`, `&&`, `;` or pipe in front of or after it. The pre-call freeze and the post-call comparison are separate commands. Order is load-bearing:

1. `--dangerously-skip-permissions`
2. optional session flags (`--model`, `--effort`, `--print-timeout`, `--output-format`)
3. `--print` (or `-p` / `--prompt`)
4. the prompt string

| Flag | Why |
|---|---|
| `--dangerously-skip-permissions` | unattended tool approval |
| `--print` / `-p` / `--prompt` | one-shot; bare `agy` is the TUI |
| `--print-timeout` | wait cap (`15m` for a review call) |
| `--model` | the slug from `models.agy` in `.harness/local.yaml` |
| `--effort` | `low` / `medium` / `high` (not `xhigh`) |
| `--output-format` | `text` (default) / `json` / `stream-json` |
| `--json-schema` | structured final output |
| `--mode` | `accept-edits` or `plan` |
| `--sandbox` | restrict the terminal |

`--input-format stream-json` reads NDJSON turns on stdin and requires `--output-format stream-json`. The default path is a prompt string. If a clone recipe is ever added, it uses the cwd flag found by the probe inside the single command, or states the `cd` exception as [claude.md](claude.md) does.

## Prompt input

Inlined in the prompt string. Prefer inlining files in the prompt to `--add-dir`: `--add-dir PATH` before the prompt has been seen to print `--add-dir` help and ignore the packet. `$PROMPT` is defined in [clone-isolation.md](clone-isolation.md), Terms.

## Isolation

No read-only flag documented; `agy` operates on the process cwd (unprobed). The review stays on the live tree under the accepted RISK-002 exposure, with the comparison of [live-tree-compare.md](live-tree-compare.md): not a sandbox, it only detects a workspace edit after the fact. A project whose `.git` is not a directory ends `isolation_unavailable` (see that file).

## Model flag

`--model <slug>`, from `models.agy` in `.harness/local.yaml`. Refresh the slugs on this machine with `agy models` and record the result in `.harness/local.yaml`, never in this file.

## Timeout

`--print-timeout`; the default is 5m, and a review call often needs longer (`15m` above). The timeout may not stop the process (see Known gotchas).

## Output capture

Reply on stdout; `--output-format` selects `text` (default), `json` or `stream-json`.

## Failure signals

- Help text instead of a reply: `--add-dir` placed before the prompt.
- A reply without the stage's verdict signal is `no_signal`, whether or not the process is alive.

## Known gotchas

The timeout may not stop the process. A complete reply on disk that carries the stage's verdict signal is the result; stop the process with one separate command, `kill <pid>` of the pid recorded when the call started (the call runs in the background and its pid is recorded, or the delegate lists `pgrep -x agy` before and after the call and kills only the new pid), so the owner's own `agy` sessions are never matched. Do not use `pkill -f` with the prompt basename or any text of the command: it matches the delegate's own shell ([gotchas.md](gotchas.md)); a bracket pattern (`pkill -f '[a]gy --print'`) is the fallback. Report it as `leftover processes` in the workspace check.

## Measured on

2026-10-16 (process kept alive past `--print-timeout 15m` after the full reply was written).
