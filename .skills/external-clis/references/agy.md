# agy (Antigravity)

Binary: `agy`. Reviewer among other roles.

## Headless

```bash
agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat "$PROMPT")"
```

The call is one shell command whose first word is `agy` (R18): no `timeout`, `&&`, `;` or pipe in front of or after it, and the pre-call snapshot and the post-call compare are separate commands. `--print-timeout` bounds the time. `agy` has no clone recipe: it runs on the live tree and does not plan until a dated H3 probe records that `agy` treats the process cwd as its workspace and which flag, if any, sets it (not checked here). If a clone recipe is ever added, it uses that flag inside the single command, or states the `cd` exception as `claude.md` does. Here `$PROMPT` is the literal path of the review prompt, created by `( umask 077; mktemp "<tmp root>/dh-prompt-XXXXXXXX" )` because the `agy` review runs without a clone; it is removed in the call's own cleanup (an interrupted call can leave it behind, mode 0600; it is not swept and is left for the owner; `<tmp root>` is defined in [clone-isolation.md](clone-isolation.md), Terms).

Order is load-bearing:

1. `--dangerously-skip-permissions`
2. optional session flags (`--model`, `--effort`, `--print-timeout`, `--output-format`)
3. `--print` (or `-p` / `--prompt`)
4. the prompt string

`--add-dir PATH` **before** the prompt has been seen to print `--add-dir` help
and ignore the packet. Prefer inlining files in the prompt instead of
`--add-dir`.

`--print-timeout` default is 5m; a review call often needs longer (`15m`
above).

**The timeout may not stop the process.** Fact (measured 2026-10-16): an `agy --print` call wrote its full reply in about 2 minutes, but the process stayed alive past `--print-timeout 15m` and needed SIGKILL. Treat a complete reply on disk that carries the stage's verdict line as the result, then stop the process with one separate command, `pkill -f -- "<the basename of $PROMPT>"`: the prompt's first line is `Call: <that basename>`, so the string is in the process arguments and unique to this call, so the owner's own `agy` sessions are never matched. Report it as `leftover processes` in the Workspace check. A reply that is incomplete or lacks the verdict line is still a transport failure, whether or not the process is alive.

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

`--input-format stream-json` reads NDJSON turns on stdin and requires
`--output-format stream-json`. Default path is a prompt string.

Refresh slugs on this machine: `agy models`; record the result in
`.harness/local.yaml`, never in this file.

## As a reviewer

Require a `VERDICT:` / `INDEPENDENCE:` / `FINDINGS:` / `SCENARIO CHECK:`
block. Help text or a missing `VERDICT` is a transport failure — it does
not spend a review round. `agy` has no documented read-only flag and, until the H3 probe is recorded, runs on the live tree: the post-call step-7 comparison (path list, per-path hashes and byte snapshot; SKILL.md step 7) is not a sandbox and only detects a workspace edit after the fact (the RISK-002 exposure the owner accepted for review).
