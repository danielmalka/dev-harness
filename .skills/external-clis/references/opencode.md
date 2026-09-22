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
OPENCODE=<resolved path or bare "opencode">
"$OPENCODE" run --dir "$PWD" --auto -m '<provider>/<slug>' \
  --variant high "$(cat /tmp/prompt.md)"
```

Long prompt as an attached file (do not stuff a multi-kilobyte prompt
into argv):

```bash
"$OPENCODE" run --dir "$PWD" --auto -m '<provider>/<slug>' \
  --file=/tmp/prompt.md -- "Follow the attached prompt."
```

`-f`/`--file` is a yargs array: a bare `-f file.md "msg"` swallows `msg`
as a second file, and `--file=` alone with no `--` message is refused.
Always `--file=path -- "message"`.

| Flag | Why |
|---|---|
| `run` | one-shot; bare `opencode` is the TUI |
| `--auto` | approve tools (required for an unattended call) |
| `-m provider/model` | pin; do not rely on the interactive picker |
| `--dir` | workspace; default is cwd |
| `-f` / `--file` | attach file(s) to the message (repeatable) |
| `--format json` | raw events if parsing is required |
| `--variant` | reasoning effort (only the levels that slug advertises) |
| `--thinking` | print thinking blocks |
| `--title` | session title |
| `-c` / `--continue`, `-s` / `--session`, `--fork` | resume / fork; default is a new session |

Refuse from an orchestrator turn: bare `opencode`, `opencode --mini`,
`run --interactive` / `-i`. `--agent` on `run` switches the **primary**
agent, not an in-process subagent (`task` inside opencode) — the
anti-delegation clause still applies to what that primary agent may do.

`opencode` has no documented read-only flag; the guarantee that it did
not edit the workspace comes from the post-call `git status` check the
parent skill describes (SKILL.md step 7).

**Do not run two `opencode run` calls concurrently.** Two simultaneous
runs against the same account have been seen to fail with a database
lock error; see `references/gotchas.md`. When a review stage lists more
than one `opencode` entry, or the concurrency ceiling would otherwise be
reached, the Coordinator serializes `opencode` calls into their own
batch rather than running them alongside each other.
