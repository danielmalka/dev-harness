# agy (Antigravity)

Binary: `agy`. Reviewer among other roles.

## Headless

```bash
agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat "$PROMPT")"
```

A clone recipe exists, `cd "$CLONE" && agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat "$PROMPT")"`, **not used until the H3 probe is recorded** (that `agy` treats the process cwd as its workspace, dated). Until then `agy` reviews on the live tree and does not plan. `$CLONE` in this recipe is the literal clone path stored when the clone was created ([clone-isolation.md](clone-isolation.md), Terms), written into the command as that literal string. Here `$PROMPT` is the literal path of the review prompt, created by `( umask 077; mktemp "<tmp root>/dh-prompt-XXXXXXXX" )` because the `agy` review runs without a clone; it is removed in the call's own cleanup (an interrupted call can leave it behind, mode 0600; it is not swept and is left for the owner; `<tmp root>` is defined in [clone-isolation.md](clone-isolation.md), Terms). If the clone recipe is ever used, the clone-call rule applies instead ([live-tree-freeze.md](live-tree-freeze.md), section 5).

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
not spend a review round. `agy` has no documented read-only flag and, until the H3 probe is recorded, runs on the live tree: the post-call step-7 comparison (path list, per-path hashes and byte snapshot; SKILL.md step 7) is not a sandbox and only detects a workspace edit after the fact (the RISK-002 exposure the owner accepted for review).
