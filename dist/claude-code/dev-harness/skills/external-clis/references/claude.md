# claude (Claude Code)

Binary: `claude`. Usable as a CLI reviewer, distinct from a Coordinator
dispatching a Claude agent in-process: this path shells out to a
separate `claude` process instead.

## Headless

```bash
cd "$CLONE" && claude -p --dangerously-skip-permissions --model <slug> --effort high \
  "$(cat "$PROMPT")"
```

`$CLONE` in these recipes is the literal clone path stored when the clone was created ([clone-isolation.md](clone-isolation.md), Terms), written into the command as that literal string. The prompt file `$PROMPT` is the literal path of a file created by `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )` inside the call's freeze directory ([live-tree-freeze.md](live-tree-freeze.md), section 5); it is unique per call and is removed with that directory. The working directory is fixed by the `cd` before the call, because no cwd flag is documented for this CLI. That call has not been measured on this kit, and it never passes `--add-dir` for the live tree.

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
| `--add-dir` | extra workspace roots (never the live tree in a clone call) |
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
has no documented read-only flag as a CLI process, so it runs in a clone: disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected. The live tree is compared as `references/live-tree-freeze.md` describes (SKILL.md step 7 covers only a call without a clone).
