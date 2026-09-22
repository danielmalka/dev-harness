# agy (Antigravity)

Binary: `agy`. Reviewer among other roles.

## Headless

```bash
agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat /tmp/prompt.md)"
```

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
not spend a review round. `agy` has no documented read-only flag; the
guarantee that it did not edit the workspace comes from the post-call
`git status` check the parent skill describes (SKILL.md step 7).
