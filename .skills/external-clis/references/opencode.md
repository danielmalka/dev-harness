# opencode

Binary: `opencode`, resolved via `.harness/local.yaml` `binaries.opencode`
when `command -v opencode` fails (a login-shell-only `PATH` is common for
this binary — do not treat that miss alone as "not installed"). Reviewer
or builder; hands: yes (`opencode run`).

The model catalog is provider-specific and changes over time; do not
invent a slug here. Refresh and pin from `.harness/local.yaml`
`models.opencode`, populated by running the provider's own model-listing
command on this machine and recording the result there, never in this
file.

## Headless

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>#high' \
  "$(cat "$PROMPT")"
```

`opencode run` has no working-directory flag and no `--variant` (fact, measured 2026-10-21 on v2.0.32, M2 of PRD-027). The call therefore starts with `cd "$CLONE" &&`, the second named exception (beside `claude`) to the rule that a call is one command whose first word is the binary; `grok` and `mcode` keep their own working-directory flag. A separate `cd` in its own Bash call does not persist (measured 2026-10-21: the runtime resets the cwd), so the `cd` stays in the same command. Because the command starts with `cd`, it does not match a `Bash(opencode:*)` allow rule: the call may need the owner's approval or be not-run. Write the resolved binary as a literal (bare `opencode`, or the literal path taken from `.harness/local.yaml`). No `timeout` goes in front of the command and no `;` or pipe chains it; the freeze and the later comparisons are separate commands. Time is bounded by the Bash tool timeout or `run_in_background`; `opencode run` has no wall-clock flag.

The variant (reasoning effort) is part of the model string: `-m provider/model#variant`, only a level that slug advertises. Without `#variant` the model default applies.

`$CLONE` in these recipes is the literal clone path stored when the clone was created ([clone-isolation.md](clone-isolation.md), Terms), written into the command as that literal string. The prompt file `$PROMPT` is the literal path of a file created by `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )` inside the call's freeze directory ([live-tree-freeze.md](live-tree-freeze.md), section 5); it is unique per call and is removed with that directory.

Long prompt as an attached file (do not stuff a multi-kilobyte prompt
into argv):

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>[#<variant>]' \
  --file="$PROMPT" -- "Follow the attached prompt."
```

`-f`/`--file` is a yargs array: a bare `-f file.md "msg"` swallows `msg`
as a second file, and `--file=` alone with no `--` message is refused.
Always `--file=path -- "message"`.

| Flag | Why |
|---|---|
| `run` | one-shot; bare `opencode` is the TUI |
| `--auto` | approve tools (required for an unattended call) |
| `-m provider/model[#variant]` | pin; do not rely on the interactive picker; `#variant` sets the reasoning effort (only levels that slug advertises); there is no `--variant` flag |
| `cd "$CLONE" &&` | workspace; there is no `--dir`, so every call starts with it |
| `-f` / `--file` | attach file(s) to the message (repeatable) |
| `--format json` | raw events if parsing is required |
| `--thinking` | print thinking blocks |
| `--title` | session title |
| `-c` / `--continue`, `-s` / `--session`, `--fork` | resume / fork; default is a new session |

Refuse from an orchestrator turn: bare `opencode`, `opencode --mini`,
`run --interactive` / `-i`. `--agent` on `run` switches the **primary**
agent, not an in-process subagent (`task` inside opencode) — the
anti-delegation clause still applies to what that primary agent may do.

`opencode` has no documented read-only flag, so it runs in a clone: disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected. The live tree is compared as `references/live-tree-freeze.md` describes (SKILL.md step 7 covers only a call without a clone).

**Do not run two `opencode run` calls concurrently.** Two simultaneous
runs against the same account have been seen to fail with a database
lock error; see `references/gotchas.md`. When a review stage lists more
than one `opencode` entry, or the concurrency ceiling would otherwise be
reached, the Coordinator serializes `opencode` calls into their own
batch rather than running them alongside each other.
