# claude (Claude Code)

Binary: `claude`. Usable as a CLI reviewer, distinct from a Coordinator
dispatching a Claude agent in-process: this path shells out to a
separate `claude` process instead.

## Headless

```bash
claude -p --dangerously-skip-permissions --model <slug> --effort high \
  "$(cat /tmp/prompt.md)"
```

`-p` / `--print` prints and exits. `--dangerously-skip-permissions`
bypasses the permission prompt, needed for an unattended call.
`--permission-mode bypassPermissions` is the named equivalent.

| Flag | Why |
|---|---|
| `-p` / `--print` | one-shot; bare `claude` is the TUI |
| `--dangerously-skip-permissions` | unattended |
| `--model` | an alias (`haiku`, `sonnet`, `opus`, `fable`, …) or a full model id |
| `--effort` | `low` / `medium` / `high` / `xhigh` / `max` |
| `--fallback-model` | pin a fallback if the default is unavailable |
| `--add-dir` | extra workspace roots |
| `--output-format` | print-only: `text` / `json` / `stream-json` |
| `--input-format` | print-only: `text` or `stream-json` |

Prompt as the last argv, or keep it in a file and expand as above. Do not
start the interactive TUI from an orchestrator turn.

## As a reviewer

Demand a `VERDICT: APPROVED|NEEDS_WORK`, `INDEPENDENCE: EXTERNAL`,
`FINDINGS`, `SCENARIO CHECK` block when this recipe is used outside the
kit's own review skills; inside the kit, the stage's own verdict signal
(SKILL.md, "Transport failure vs. a real verdict") is what the Coordinator
reads. A long run that never emits either is a transport failure. `claude`
has no documented read-only flag as a CLI process; the guarantee that it
did not edit the workspace comes from the post-call `git status` check
the parent skill describes (SKILL.md step 7).
