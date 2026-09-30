# Running a CLI call in a disposable clone

This file is the exact lifecycle for a CLI call that runs in a disposable
clone: create, populate, verify, use, remove, sweep orphans. It applies to a
review or planning call to a binary with no documented and measured
read-only mode. `codex` keeps its own sandbox and never uses it.

Every mention of this mechanism carries one qualification, stated in full:

disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected

The clone is a
working directory, not a sandbox. The CLI runs as the same OS user and still
reaches the live tree and the rest of the machine by absolute path. Nothing
here prevents a write; it only puts the call's working directory somewhere
else and lets the Coordinator throw that directory away.

This file defines only the clone lifecycle. Which binaries may plan is decided
under "The planning call" in `SKILL.md`, not here.

Terms: `<project>` is the project root (the live tree). `$CLONE` is the exact
clone path. Shell variables do not survive between separate tool calls, so the
call stores the literal path that `mktemp` printed, and every later command
(verification, comparison, removal) uses that literal string, not `$CLONE`.
`<call-id>` is chosen once by the Coordinator before step 1, as
`<task id>-<binary>-<n>` with `<n>` a per-batch counter, and may contain only
`[A-Za-z0-9_-]`; anything else is `isolation unavailable`. It is used verbatim
in both names of a call: `dh-cli-<call-id>-XXXXXXXX` (the clone) and
`dh-freeze-<call-id>-XXXXXXXX` (the freeze directory, which also holds the
call's prompt file). Calls that overlap in time share one freeze directory,
named after the first call, with one clone and one prompt file each.

`<tmp root>` is `(cd "${TMPDIR:-/tmp}" && pwd -P)`, resolved once when the call (or
the batch) starts and stored as a literal path. Every `mktemp` template, every
removal guard and the orphan sweep use `<tmp root>`, never the unresolved
variable.

The project must have `.harness/` and `.claude/settings*.json` git-ignored. If
it does not, the live tree holds them as tracked or untracked files, the clone
would carry them, the absence checks in step 2 fail, and every call ends as
`isolation unavailable`. It is never fixed up by deleting them from the clone.

## 1. Create the clone (R2)

Before creating anything, `test -d <project>/.git && test ! -L <project>/.git`
must pass: a project whose `.git` is not a real directory, a linked worktree
with a gitfile for example, is `isolation unavailable`. Run these in order. Every git command in this file, in every step, runs with
the prefix `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1` (the commands
written out in this step included). `mktemp` gives the random suffix, mode 0700 and
uniqueness. `<tmp root>` is the only root the kit creates clones in and
the only root the orphan sweep (step 5) reads.

```
CLONE=$(mktemp -d "<tmp root>/dh-cli-<call-id>-XXXXXXXX")
GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git clone --no-hardlinks <project> "$CLONE"
GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" remote remove origin
rm -rf "$CLONE/.git/logs"
printf 'pid=%s\nstarted=%s\n' "$PPID" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$CLONE/.git/dh-cli.lock"
```

`$PPID` is the process that hosts the Coordinator's shell, not the short-lived
shell itself. `started` is ISO 8601 in UTC; `date -Iseconds` (GNU) also works
but the `date -u +%Y-%m-%dT%H:%M:%SZ` form is portable. The clone has its own
`.git` (hooks, config, refs, objects, `info/exclude`), holds no ignored file
of the live tree, and nothing uncommitted: it is born clean at `HEAD`. Until
the third command runs, `origin` still points at the live tree; the reflog
under `.git/logs` records the live path, which is why the fourth command
deletes it.

Verify before going on. Each check must pass:

- `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" remote` prints nothing.
- `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" rev-parse --git-common-dir` resolves to `$CLONE/.git`, the
  clone's own directory.
- `test ! -e "$CLONE/.git/objects/info/alternates"`.
- `test ! -e "$CLONE/.git/logs"`.

## 2. Populate the clone and match it (R3)

The clone must hold the state under judgment: the current commit, plus the
tracked modifications and the new non-ignored files of the live tree at the
freeze. Run the live-tree commands from `<project>`, each prefixed with
`GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false`.
Every git command run against the clone (`apply`, `status`, `hash-object`)
gets the same `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1` prefix, so a
global or system configuration cannot change or run anything there.

1. Staged part: save
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false diff --cached --binary --no-ext-diff --no-textconv --src-prefix=a/ --dst-prefix=b/`
   to a file. If the diff is empty, skip this step: `git apply` exits 128
   ("No valid patches in input") on empty input, so an empty diff is never
   piped in. (`git apply --allow-empty`, from git 2.35, is an alternative;
   the skip works on any git version and is the rule.) Otherwise apply it with
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" apply --index`.
2. Unstaged part: the same command without `--cached`, with the same skip on
   an empty diff, applied with
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" apply` (no `--index`).
3. New files: from `<project>`, list them with
   `git ls-files -z --others --exclude-standard`. Any entry ending in `/` (a
   nested repository shows as `sub/`) is `isolation unavailable` on its own,
   before anything is copied; it is not left to the match. If the list is
   empty, skip the copy. Otherwise run, under `set -o pipefail` so a failing
   `ls-files` is not masked by the `tar` stages,
   `git ls-files -z --others --exclude-standard | tar --null -T - -cf - | tar -xpf - -C "$CLONE"`.
   The NUL-separated list survives any file name, `tar` stores a symlink as a
   symlink instead of following it, and `-p` keeps the mode. GNU tar accepts
   this pipeline (an empty list also exits 0 there; the skip is kept anyway);
   that bsdtar (macOS, BSD) accepts it is a hypothesis, not measured here.
   Ignored files are never copied.

Then match. All of these must hold:

- The judged set (tracked paths plus untracked non-ignored paths) is the same
  path set in the clone and in the live tree. A path absent on disk records the
  marker `deleted` on its side, and the markers must equal too.
- Per path, by kind. A regular file: the same `git hash-object --no-filters` hash and the
  same mode (`stat -c %a` on GNU, `stat -f %Lp` on BSD), for tracked and
  untracked paths alike. A symlink: the same link text from `readlink`, never
  a hash of the target it points at. A submodule: the same commit id.
- `git status --porcelain -uall` (never `--ignored`) prints the same lines in
  the clone (run with the clone-side prefix above) and in the live tree, including the two XY status letters, so the
  clone reproduces the index (staged part and unstaged part), not only the
  content. Only the marker `.git/dh-cli.lock` is ignored, and it is inside
  `.git`, so status does not show it.
- `test ! -e "$CLONE/.harness"` and
  `test ! -e "$CLONE/.claude/settings.local.json"` pass.

A state that does not reproduce exactly (a rename, a mode change, a symlink, a
submodule, a nested repository (refused before the copy), a binary file, or an XY that the two-step
population cannot recreate) is `isolation unavailable` (step 4). The mismatch
is the gate. Never adjust the clone by hand until it matches, and never
continue on the live tree instead.

## 3. Use the clone and remove it (R8)

Fix the CLI's working directory to the stored clone path: every clone recipe
starts with `cd "$CLONE" &&` and keeps the binary's own cwd flag as
`references/<binary>.md` shows (`--cwd` or `--dir`). Measured (T-1015): with
`grok --cwd "$CLONE"` launched from the Coordinator's cwd, `/proc` showed the
`timeout` wrapper, and briefly `grok` itself, with a cwd inside the live tree;
with `cd "$CLONE" &&` first, every process of the group had cwd = the clone. The
prompt does not carry the live tree's absolute path.

When the call ends (success, transport failure, timeout or change detected),
and after the comparisons of the live tree and of the clone have run, the
Coordinator removes the clone with `rm -rf` on the exact stored path, and only
that. Guard first; all of these must hold: the parent directory of the
stored path equals `<tmp root>`; its basename starts with `dh-cli-` and
contains no `/` and no `..`; `test -d` passes; `test ! -L` passes. A path that
fails the guard is never passed to `rm`. Then confirm the directory no longer exists. If
`rm` fails, report it and leave the directory for the orphan sweep. Nothing
from the clone is copied or merged into the live tree: the verdict or the plan
text comes back as the reply, as before.

## 4. Isolation unavailable (R11)

The call is not-run, with the reason `isolation unavailable` ("isolamento
indisponivel" in the owner's language), when any of these holds: the project
is not a git repository; its `.git` is not a real directory (`test -d` and
`test ! -L`, a linked worktree for example); any step of the freeze
(`references/live-tree-freeze.md` section 1) fails (no CLI starts and no comparison runs; the partial freeze directory is removed under the guard of step 3);
`HEAD` has no commit yet; disk is short; `<call-id>`
has a character outside `[A-Za-z0-9_-]`; a check in step 1 fails;
the population or the match in step 2 fails. Remove any partial clone at its
exact stored path under the exact guard of step 3, including one that has no marker
because its creation failed, report the reason, and stop. There is no
fallback: a binary with no read-only mode is never run on the live tree in
place of the clone. The round is not spent, and the call counts as neither an
approval nor a rejection.

## 5. Orphan sweep (R9)

Interrupted runs can leave clones behind. Once per batch of dispatches, run by
the Coordinator before the first dispatch of the batch and never inside a
call, list the direct children of `<tmp root>` and remove a directory
only when all five conditions hold. The sweep covers two kinds: a clone
(`dh-cli-`, marker `.git/dh-cli.lock`) and a freeze directory (`dh-freeze-`,
marker `dh-freeze.lock` at its top, written by
`references/live-tree-freeze.md`); the same five conditions apply to both,
with the prefix and marker of its kind:

1. Its name starts with `dh-cli-` (or `dh-freeze-`), it is a direct child of `<tmp root>`,
   and it is a real directory, not a symlink (`test -d` and `test ! -L`).
2. It belongs to the current user (`stat -c %U` on GNU, or `find -user`).
3. It contains its marker: `.git/dh-cli.lock` (or `dh-freeze.lock`).
4. The `pid` in the marker no longer exists (`kill -0 <pid>` fails with "no
   such process"). A `kill -0` that fails with EPERM means the process exists,
   so the PID counts as alive and the directory stays.
5. The `started` in the marker is older than the call timeout (15 minutes by
   default) plus 5 minutes of margin.

The conditions are conjunctive. Never touch a directory that fails any one of
them: no prefix, a symlink, no marker, another user, a live PID, or a recent
`started` is left alone, including a directory of another project.

## 6. Parallel calls (R10)

Two calls in clones at the same time each get their own `<call-id>` and their
own `$CLONE` (distinct `mktemp` results). Each call removes only its own
clone, and never touches another `dh-cli-` directory; only the sweep in step 5
does, before the batch. When such calls overlap in time, the order is: one live-tree freeze taken
before the first of them; then every clone created, populated and matched;
then every prompt written; only then the first CLI started. Calls are dispatched in closed groups: a new clone call never starts while a
CLI from an earlier group is still running. They share one
freeze directory (named after the first call), each with its own clone and
prompt file. R5 and the reversion run only after the last one has ended; any R5 difference marks every overlapping call not-run,
`change detected`, and R6 stays per clone (`references/live-tree-freeze.md`,
Order). The cap of two concurrent specialists and the serialization of
`opencode` do not change.

## 7. What stays outside the clone (residual, stated in full)

The clone is a working directory, not a sandbox: disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected.

What the live-tree comparison (R5) covers: the list of
`git status --porcelain -uall --ignored` with the hash of each path; the hash
of every tracked file; the output of `git ls-files -v --stage` (index flags
and blob ids, which catches `--skip-worktree` and `--assume-unchanged`); the
bytes of `.harness/`, of `.claude/settings*.json` and, from the live `.git`,
of `hooks/`, `config`, `info/`, `HEAD` and `packed-refs`; and the output of
`git for-each-ref`.

Not covered by R5, so not detected (this list of live `.git` paths is the one
the author knows and is not exhaustive; any live `.git` path R5 does not name
is outside the coverage): `.git/index` as raw bytes (only its logical content
is read by `ls-files`), `logs/` (reflogs), loose objects and packs,
`worktrees/`, `refs/` as files (only what `for-each-ref` shows), `modules/`,
`FETCH_HEAD`, `ORIG_HEAD` and `COMMIT_EDITMSG`.

Accepted by the owner on 2026-09-29, until the OS sandbox batch, and not
detected:

- (a) a write outside the project directory: dotfiles under the home
  directory, `~/.ssh`, the global git configuration, credentials and the
  CLI's own configuration, other repositories, other folders under `/tmp`;
- (b) processes that outlive the call, network calls, and `git push` to any
  remote that is not the live repository;
- (e) reading secrets (`.env`, `~/.aws`, anything the user reaches) and
  returning them in the reply, because the clone only hides them from a
  relative path.

The freeze directory of `references/live-tree-freeze.md` is covered by (a) and
(e) above. It holds copies of `.harness/`, `.claude/settings*.json` and the
dirty or untracked non-ignored files, it sits in `<tmp root>`, and a CLI
running as the same user can read and alter it; mode 0700 does not stop that.
Tampering with the copies is a write outside the project (a), reading them is
reading secrets (e). The same holds between overlapping calls: a running CLI
can read or alter a sibling call's prompt file or the shared freeze directory
(residual (a)). Writing every prompt before any CLI starts narrows this window;
it does not close it. Starting each CLI in its own process group and killing
the group after it ends narrows residual (b); it does not close it. An
interrupted `agy` review can leave its `dh-prompt-XXXXXXXX` file (mode 0600)
under `<tmp root>`: it is not swept, and it is left for the owner to remove.

Also not detected: (c) a write to the live tree that is undone before the
second snapshot or that preserves the bytes; (f) a concurrent write by another
process, the Coordinator included, which step 7 of `SKILL.md` already says to
avoid. Only an operating-system isolation would close these; it is a later
batch.
