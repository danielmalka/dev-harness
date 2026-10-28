# Live-tree freeze, clone comparison and prompt hygiene

This file is the comparison procedure for a CLI call that runs in a disposable
clone (lifecycle in [clone-isolation.md](clone-isolation.md)). It replaces
step 7 of `SKILL.md` for such a call; step 7 as written stays in force for a
call without a clone (`codex`, and the `agy` review on the live tree). The `agy`
live-tree review of a project whose `.git` is not a directory (a linked
worktree; a deliberate superset that also covers submodule and bare) ends
`isolation unavailable` and does not run: its step-7 freeze reads `.git/hooks/`
and `.git/config` of `<project>/.git`, which are not folders there (`SKILL.md` step 7).

The delegate runs every step of this file: a dispatched `qa-verifier`, or
`builder` in lane `infra` (`SKILL.md`, "Who runs the call"). It returns the
verdict line and the relative path of the raw output; the Coordinator runs
none of these steps.

The qualification for every mention of this mechanism, stated in full:

disposable clone; writes outside it are detected on the live tree (R5), never prevented; writes outside the project are not detected

Two comparisons run, on two trees. The live-tree comparison (R5) detects and
reverts what the CLI wrote where it should not. The clone comparison (R6)
judges whether the CLI changed the content it was given. Neither prevents
anything.

## Order

1. Preconditions and freeze (section 1). The precondition is step 1 of
   clone-isolation.md: `.git` a real directory, or a gitfile with a non-null
   `linked_worktree` and the `test -f`/`test ! -L` check right before the clone;
   anything else is `isolation unavailable`. When clone calls overlap in time,
   take one freeze before the first of them. Save the copies outside the
   project tree and outside the clone: create a freeze directory with
   `mktemp -d "<tmp root>/dh-freeze-<call-id>-XXXXXXXX"` (`<tmp root>` and
   `<call-id>` as in clone-isolation.md, Terms; the call-id is that of the first
   call; mode 0700) and write the marker `<freeze dir>/dh-freeze.lock` with the
   same two lines as the clone marker (`pid=` and `started=`,
   clone-isolation.md step 1). Then compute the digest of the freeze (end of
   section 1). If any step of this freeze fails (a missing tool included, see the note on
   GNU tools in Notation), the call is not-run, `isolation unavailable`: no
   CLI starts and no comparison runs. Remove the partial freeze directory under
   the step 8 guard.
2. Create, populate and match every overlapping clone (clone-isolation.md steps
   1 and 2), and freeze each clone (section 3). Overlapping calls share the one
   freeze directory, with one clone and one prompt file each.
3. Write every call's prompt (section 5), all of them before any CLI starts.
4. Start each CLI as one shell command that begins with the binary name, in its
   own tool call, with the binary's own working-directory flag where its recipe
   shows one (`grok --cwd`, `mcode exec --cwd`); `claude` and `opencode` have no
   such flag and keep `cd "$CLONE" && <binary> ...` as the two named exceptions
   (`references/claude.md`, `references/opencode.md`; the form starts with `cd`,
   does not match a `Bash(<binary>:*)` allow rule and may need the owner's
   approval; a separate `cd` in its own call does not persist). The opencode
   form is `cd "<clone>" && <opencode binary> run --auto -m
   <provider/model>[#<variant>] --file=<prompt file> -- "Follow the attached
   prompt."`. No `timeout`, `setsid`, subshell,
   `&&`, `;` or pipe wraps or chains the call, and the freeze and the later
   comparisons are separate commands (a permission rule `Bash(<binary>:*)`
   matches only a command that starts with the binary). The binary's own time
   flag, or the Bash tool timeout, bounds the call. Then, as one separate
   command run once the reply is on disk or the timeout has passed (before any
   comparison), stop this call's leftover processes by pid: `kill <pid>`. The
   pid comes from the call itself: the call runs in the background
   (`run_in_background`) and its pid is recorded when it starts, or the delegate
   lists `pgrep -x <binary>` before the call and after it and kills only the new
   pid. Never `pkill -f` with the clone path or the command text: the pattern
   matches the delegate's own shell (`references/gotchas.md`); if a pattern is
   unavoidable, use a bracket pattern (`'[o]pencode run'`). The pid is unique to
   this call, so it never matches the owner's own sessions of the same binary; a
   pid already gone is ignored. This narrows residual (b) of clone-isolation.md
   section 7 (a process that outlives the call); it does not close it, since a
   child that has a different pid is missed. Report it as `leftover processes` in the
   Workspace check.
5. When a call ends, recompare its clone (section 3, R6). R6 stays per clone.
6. After the last overlapping call has ended, and only then, on the live tree:
   verify the freeze digest, extract the copies, compare and restore the trust
   paths with non-git tools (section 2, part A), and, only when part A reports no
   remaining difference, run the git-based comparison (section 2, part B).
7. Revert what the comparison found (section 4).
8. Remove each clone (clone-isolation.md step 3), then the freeze directory,
   prompt files included, with the same guard (parent equals `<tmp root>`,
   basename starts with `dh-freeze-`, `test -d`, `test ! -L`). Nothing is ever
   swept by age on its own.

Clone calls are dispatched in closed groups: a new clone call never starts (its
clone is never created) while a CLI from an earlier group is still running. A
group is the set of calls that overlap in time under one freeze.

Nothing else may write to the live tree or to the `<harness dir>` between step 1 and step 6, the delegate and the Coordinator included (the freeze directory is shared by overlapping calls, so this covers all of them): a build, a
package regeneration or a test run started in parallel changes files the
comparison then blames on the CLI (the measured case is in step 7 of
`SKILL.md`). Overlapping clone calls are governed by the same rule: R5 and the
reversion run only after the last one has ended, because a sibling's reversion
would otherwise hide the writer. Any R5 difference marks every overlapping call
not-run, `change detected`, as `coordinator.md` already does for calls without
a clone. The cap of two concurrent specialists is unchanged. The reversion in
section 4 is only sound when no process started by the call is still running; a
process that outlives the call can write again after the reversion, and that
stays undetected (clone-isolation.md section 7, item (b)). Either wait, or run
the parallel work in a separate checkout.

Not-run reasons for a clone call: `change detected` ("mudança detectada"),
from this file, and `isolation unavailable` ("isolamento indisponivel"), from
clone-isolation.md step 4. Both discard the verdict, neither counts as an
approval or a rejection, and neither spends a correction round.

## Notation

`R5GIT` stands for
`GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false -c core.hooksPath=/dev/null`,
that is, the prefix `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 -c core.fsmonitor=false -c core.untrackedCache=false -c core.hooksPath=/dev/null`
placed on `git`. Every git command of the live-tree comparison, the reversion
and the clone comparison runs as `R5GIT <subcommand>`, so a global or system
configuration can neither change the result nor execute a command, and no hook
runs. The `core.hooksPath` override is needed because a git command that
updates a ref or the index runs the live repository's own hooks otherwise
(measured: `update-ref` ran a planted `reference-transaction` hook twice; with
the override it ran none). `<project>` is the live tree.

This file assumes GNU tools in a few places: `sha256sum` (macOS: `shasum -a 256`),
`stat -c` (BSD: `stat -f`) and `tar --null -T`.
A tool that is missing during the freeze, with no equivalent, is
`isolation unavailable`.

Hashing rule for one path, used everywhere below: a regular file is hashed
with `R5GIT hash-object --no-filters -- <path>`; a symlink is recorded as its
`readlink` text and never followed; a path absent on disk records the marker
`deleted`. A directory is never hashed, only expanded to the files and
symlinks under it (`find <dir> \( -type f -o -type l \)`, which does not
follow symlinks). Every recorded path also records its mode (`stat -c %a` on
GNU, `stat -f %Lp` on BSD). Entries that are neither a regular file nor a
symlink (directories, fifos, sockets) are handled by the freeze list of
section 1, not by this rule.

## 1. Freeze the live tree (R5)

RISK-002 prevention: nothing writes to the tree or to the `<harness dir>` between the freeze and the comparison, the Coordinator included. Detection and revert below cover the resolved `<harness dir>` in mode `repo`, wherever it is (inside the project, or the main checkout's `.harness/` when the project is a linked worktree); in mode `global` it is outside the freeze of this procedure (RISK-002 residual).

Run from `<project>`. Record all of these, before the call:

- The trust paths, as bytes and mode, saved with non-git tools (`cp -p`, or
  `tar` from an explicit NUL list, below): everything under the resolved `<harness dir>` in mode `repo`, wherever it is (inside the project, or the main checkout's `.harness/` when the project is a linked worktree), and under `.harness/` when it exists. In mode `global` the harness dir (`<home>/projects/<name>/`) is outside the freeze of this procedure, a stated RISK-002 residual; only the step 7 freeze of `SKILL.md`, for a call without a clone, saves it by bytes in both modes;
  `.claude/settings.json` and `.claude/settings.local.json` (each when
  present, and the list of those absent); from the common dir `<C>`, `hooks/`,
  `config`, `config.worktree` (when present), `info/`, `HEAD` and `packed-refs`; for a linked worktree, also the admin dir `config.worktree` (read from the `gitdir` field of `linked_worktree`; part (b) below); and every `.gitattributes` file
  in the working tree (found with `find . -path ./.git -prune -o -name .gitattributes -print`).
  `<C>` is `<project>/.git` for a normal checkout and `linked_worktree.commondir`
  (from `dh harness-path --json "<project>"`) for a linked worktree, where
  `<project>/.git` is only a gitfile and holds none of these paths; an item
  absent in both snapshots is not a difference and never an error. In a linked worktree the harness dir is the one `dh harness-path --json` gives
  in `dir` (the main checkout's `.harness/` in mode `repo`), not the ignored copy
  of `.harness/` inside the worktree; in mode `global` it stays outside the freeze of this procedure.
  Also trust items, compared by bytes with non-git tools only (`cat`, `cmp`,
  `stat`; section 2, part A): in a normal checkout, that `<project>/.git` is
  still a real directory (not a file or a symlink); for a linked worktree, the
  bytes of the gitfile `<project>/.git`, and `<admin>/commondir` and
  `<admin>/gitdir`; and that `<C>/commondir` and `<project>/.git/commondir` are
  absent when they were absent at the freeze (a planted one makes git
  re-discover another repository, with that repository's config and
  `info/attributes`). They are saved before the call like the other trust paths.
  The `config.worktree` files are trust paths like `<C>/config`: with
  `extensions.worktreeConfig` on, git reads them (`core.attributesFile`,
  `filter.<x>.smudge`, `include.path`, `core.worktree`), so a planted one could run
  a command in the delegate's own git. They are saved before the call and
  restored with non-git tools in part A of section 2, before any git command.
- `R5GIT status --porcelain -z -uall --ignored`, expanded to files (an ignored
  directory such as `!! build/` is walked, without following symlinks), with
  the hash of each path listed, ignored or not. Every ignored path is hashed
  (the owner's decision D2); a cost limit is set only by an amendment after a
  measurement.
- The hash of every tracked file (`R5GIT ls-files -z`), not only the listed
  ones.
- `R5GIT ls-files -v --stage`, saved as text. Its first letter is the index
  flag: `H` neither flag, `S` skip-worktree only, `h` assume-unchanged only,
  `s` both skip-worktree and assume-unchanged. Each entry also holds the mode
  and blob id.
- `R5GIT for-each-ref --format='%(refname) %(objectname)'`, saved as text.
- Git dirs, two parts; `config.worktree` (both dirs) is a restored trust path
  (above), and so are the gitfile, `<admin>/commondir` and `<admin>/gitdir`;
  the rest is detected only. Part (a), every project: in the common dir
  `<C>`, besides the trust paths above, the files under `refs/` (bytes and
  mode) and the entries under `worktrees/` when they exist, with the files of
  each entry that part (b) lists. Part (b), only when `<project>/.git` is a
  gitfile: the bytes of the gitfile and, in the worktree admin dir (read from the
  `gitdir` field of `linked_worktree`), `HEAD`, `commondir`, `gitdir`,
  `config.worktree`, the files under `refs/` and the presence of `locked`. An
  item absent in both snapshots is not a difference and never an error. The raw
  `index` files (of the common dir and of each admin dir) stay out of the freeze.
  A change in these items during the call is `change detected`, as in a normal
  checkout, and the verdict is discarded. Except the restored trust paths above, they are not
  saved for reversion: each is reported as `not revertible: <path>`, like an
  altered ignored path.
- Every entry under the project that is not a regular file or a symlink,
  excluding `.git`: `find . -path ./.git -prune -o ! -type f ! -type l -print`
  (directories included), saved as text with the mode of each entry
  (`stat -c %a` on GNU, `stat -f %Lp` on BSD). A new empty directory, FIFO or socket
  then appears as a created path, and section 4 removes it with `rm -rf --`
  after the ancestor check. A directory that disappeared is recreated with
  `mkdir -p`; a directory whose mode changed gets `chmod <saved mode> <dir>`
  during the reversion; a FIFO or socket that disappeared is reported as
  `not revertible: <path>`.

Save the bytes and the mode only of the paths a reversion can need: every
non-ignored path listed by `status` (dirty or untracked), every tracked path
whose saved `ls-files -v` letter is not `H` (a file flagged `h`, `S` or `s`
can hold an uncommitted edit that `status` hides), and the trust paths. A
tracked path is "clean" only when its letter is `H` and `status` does not list
it; a clean tracked path is not copied, it is restored from `HEAD` (section 4,
item 3). Write the NUL-separated list of
the saved paths to `<freeze dir>/paths.nul`, leaving out every path recorded as
`deleted` (tar fails with "Cannot stat" on a path that is missing, exit 2; a
deleted path has nothing to save and is recorded as `deleted` in the hash
list), then run
`tar --null -T "<freeze dir>/paths.nul" -cf "<freeze dir>/live.tar"` from
`<project>` (`tar` stores a symlink as a symlink; GNU tar accepts this, and
that bsdtar does too is a hypothesis, not measured here). An ignored path is
hashed and its mode recorded, but its bytes are not saved: an altered or
deleted ignored path cannot be reverted. It discards the verdict and is
reported as `not revertible: <path>`. Also record, for every ancestor
directory of a saved path, whether the freeze saw it as a real directory
(`test -d` and `test ! -L`).

Right after the records above are complete, compute `sha256sum` of
`<freeze dir>/live.tar` and of every saved text snapshot in `<freeze dir>` (the
manifest: the hash list, the `status` list, `ls-files -v --stage`,
`for-each-ref`, the list of non-file entries and the ancestor records). Keep
the digest in the delegate's own session output, never in a file the CLI can
reach. This detects tampering with the copies; it does not prevent it.

## 2. Compare the live tree, after the call (R5)

Before part A, verify the freeze: recompute the `sha256sum` of `live.tar` and of
the manifest files and compare with the digest the delegate kept. A mismatch
is `change detected` plus `freeze tampered`: nothing is restored from the
copies, no git command runs, and the report lists only the digest mismatch and
a `find` listing of the live trust paths for the owner to restore by hand (no
`cmp` against the untrusted copies). This is detection, not prevention: a CLI
that runs as the same user can still write to the freeze directory, and
only the tampering is noticed. When the digest matches, extract the copies:
first require `test -d "<freeze dir>" && test ! -L "<freeze dir>"` on the stored
freeze-directory path; a failed check is `freeze tampered`. Create the
extraction directory fresh, `mktemp -d "<freeze dir>/copy-XXXXXXXX"`, store its
literal path as `<copy dir>`, and extract with
`tar -xpf "<freeze dir>/live.tar" -C "<copy dir>"`. Never `mkdir -p` a fixed
name here: a planted symlink at a fixed name would send the extraction
elsewhere. Every restore below reads from `<copy dir>`.

Ancestor check, before any `find` or `rm` in part A or part B. Every path
involved (every saved, created or listed path, and the trust roots `<C>`,
`<C>/hooks`, `<C>/info`, the resolved `<harness dir>` in mode `repo` (`dir` from `dh harness-path --json "<project>"`: `<project>/.harness` in a normal checkout, `<main>/.harness` in a linked worktree), `.harness` when it exists, and `.claude`) must have each ancestor
still a real directory (`test -d` and `test ! -L`). An ancestor that the freeze
recorded as a real directory and that is now a symlink or a regular file is
planted: remove only that ancestor (`rm -- <ancestor>`), drop from the created
list every path below it (they were reached through the planted link and are not
in the project), and recreate the ancestor with `mkdir`. `rm -rf` does not
follow the final component of its argument, and the intermediate components are
the ones this check covers. Never pass a trust directory with a trailing slash
to `find` or `rm` (a trailing slash makes both follow a symlink).

Part A, trust paths first, with non-git tools only. Repo-local configuration
can execute a command during `status`, `hash-object` (without `--no-filters`),
`checkout` and similar (a `filter.<name>.clean` or `smudge` driver, selected
by a `.gitattributes` file), so a hostile `<C>/config` or `.gitattributes`
must be put back before any git command runs. Using only `tar`, `cmp`, `find`, `rm`, `mv`, `cat`, `mktemp`, `chmod`, `mkdir` and
`ln`, compare every trust path with the copy in `<copy dir>`, list the trust
directories (the resolved `<harness dir>` in mode `repo` (`dir` from `dh harness-path --json "<project>"`: `<project>/.harness` in a normal checkout, `<main>/.harness` in a linked worktree), `.harness` when it exists, `<C>/hooks`, `<C>/info`) with `find` and compare
the list, and restore every difference at once as in section 4, items 1 and 2, for the
trust paths and the `.gitattributes` files only. The same non-git byte
comparison covers the repository-discovery items of section 1: the gitfile
bytes (or, in a normal checkout, that `.git` is still a real directory, not a
file or a symlink), `<admin>/commondir` and `<admin>/gitdir`, and that
`<C>/commondir` and `<project>/.git/commondir` are still absent when they were
absent at the freeze. On a difference, restore from the freeze copy or remove
the planted file before any git command. If it cannot be restored or removed,
no git command runs at all: skip part B and the reversion, and record the
verdict as discarded, not-run, with the reason `not revertible`. No git
command runs in part A. A difference here is `change detected`.

Order: every part A byte comparison and restoration happens before the first git
command of part B; the git commands of part B and of the reversion (items 3 to
5 of section 4) run only after part A reports no remaining difference.

Part B, then the git-based comparison. Recompute every other item of section 1
and compare each with the freeze. Any difference is `change detected`: a path
created (absent, now present), deleted or altered (a different hash, mode or
link text); an index entry whose flag letter or blob id differs; a ref
created, moved or deleted (`for-each-ref` differs). The verdict is discarded
regardless of what the CLI said, the rest of the reversion of section 4 (items
1 and 2 for the other paths, then items 3 to 5) runs, and the report says
`change detected`.

## 3. Freeze and compare the clone (R6)

Right after the match of clone-isolation.md step 2, record the judged set: the
tracked paths (`R5GIT -C <clone path> ls-files -z`) and the untracked
non-ignored paths (`R5GIT -C <clone path> ls-files -z --others --exclude-standard`),
each with its hash (`deleted` when absent).

Before the clone is removed, recompute the hash of each path of that initial
set:

- A changed hash, or a removed path of the set, is `change detected`. The
  verdict is discarded and nothing is reverted: the clone is thrown away.
- Ignored paths of the clone and any new untracked file are outside the set.
  Neither reaches the live tree nor changes the judged content, so neither
  discards the verdict. Each new untracked file (a path in `ls-files` now that
  is not in the initial set) is recorded as a note in the call's report.
- The clone's `.git`, including `dh-cli.lock`, is not compared.

## 4. Revert the live tree (R5)

Revert from the freeze copies, and for a clean tracked path from `HEAD` (item
3). Never write through an existing path: what is there now is removed first,
without following a symlink, and the saved content is placed by a rename.
Items 1 and 2 first restore the trust paths and the `.gitattributes` files, in
part A of section 2, with no git command; in part B they run again for every
other path. Items 3 to 5 run only in part B, after part A. Afterwards recompute
section 1 and confirm it equals the freeze; if it does not, report `reversion
incomplete` beside `change detected`.

1. **A path created** (absent before): `rm -rf -- <path>`, which does not
   follow a symlink. A tracked path saved with letter `S` whose worktree file
   was absent before the call is recorded `deleted`, so a file the CLI creates
   there is removed by this item. All created paths are removed before any path is
   restored, so a created path never blocks a restore.
2. **A path altered or deleted, with a saved copy.**
   - First the ancestors (the ancestor check of section 2 has already run;
     repeat it here for each path), from the project root down. Git does not track
     through a symlink, so every ancestor of a tracked path was a real
     directory before the call by definition; for any other path the freeze
     recorded whether the ancestor was a real directory. An ancestor that was a
     real directory and is now a symlink (`test -L`) or a regular file is
     planted: remove it with `rm -- <ancestor>`. A `.git` (or any other
     ancestor) that was already a symlink before the call is left alone. Then
     `mkdir -p -- <parent>` for every missing ancestor.
   - A path with saved bytes and mode, written as one chained subshell so the
     temporary name is not lost:
     `( umask 077; tmp=$(mktemp <parent>/.dh-restore-XXXXXXXX) && cat -- <saved copy> > "$tmp" && chmod <saved mode> "$tmp" && rm -rf -- <path> && mv -f "$tmp" <path> )`.
     A rename replaces the destination entry itself and does not follow a final
     symlink: a symlink planted over `.harness/project.yaml` is removed and its
     target is untouched (measured: exit 0, target untouched, regular file with
     the saved mode).
   - A path that was a symlink before the call, in the same shape, after
     removing the name `mktemp` created (`ln -s` fails on an existing name):
     `( tmp=$(mktemp <parent>/.dh-restore-XXXXXXXX) && rm -f "$tmp" && ln -s <saved link text> "$tmp" && rm -rf -- <path> && mv -f "$tmp" <path> )`.
3. **A clean tracked path** (saved letter `H`, not listed by `status`, so no
   saved copy) that is now altered or deleted. Part B only, after part A has
   restored the hooks, the config and the attributes, with the
   `core.hooksPath` override in force. Handle the ancestors as in item 2, then
   remove what is there, `rm -rf -- <path>`, then clear both index flags, one
   command each (`R5GIT update-index --no-skip-worktree -- <path>`, then
   `R5GIT update-index --no-assume-unchanged -- <path>`), and only then run
   `R5GIT checkout HEAD -- <path>`. Plain `checkout -- <path>` is not used: it
   restores the index blob, and a tampered index entry gives back the tampered
   content (measured). `checkout HEAD -- <path>` restores the `HEAD` content,
   resets the index entry to `HEAD`, leaving `H`. Measured (T-1015, git 2.53.0):
   `checkout HEAD -- <path>` on a removed path whose index entry carries
   skip-worktree fails with exit 1, "pathspec did not match any file(s) known to
   git"; clearing the flags first, as above, works and leaves `H`. Item 4 then
   sets the saved letters. For such a path the index equals `HEAD` by
   construction, so `HEAD` is its pre-call state. Any filter that runs there is
   the owner's own pre-call configuration, restored in part A; a `cat-file`
   restore would break filtered repositories such as git-lfs.
4. **Index entries, then flags**, after items 2 and 3. The saved letters and
   blob ids are the final word: compare the saved `ls-files -v --stage` with
   the new one, path by path:
   - first entries: a stage entry whose mode or blob id differs is set back
     with `R5GIT update-index --cacheinfo <mode>,<blob id>,<path>` using the
     saved values; an entry that did not exist before is removed with
     `R5GIT update-index --force-remove -- <path>`;
   - then flags, one flag per invocation (measured: `--no-skip-worktree
     --no-assume-unchanged` in one invocation leaves `S` set; separate
     invocations clear it). Clear first: `R5GIT update-index --no-skip-worktree -- <path>`,
     then, as its own command, `R5GIT update-index --no-assume-unchanged -- <path>`.
     Then set what the saved letter had, each in its own command: for `S`,
     `R5GIT update-index --skip-worktree -- <path>`; for `h`,
     `R5GIT update-index --assume-unchanged -- <path>`; for `s` both, in two
     commands; for `H` nothing more.
5. **Refs**, last, from the saved `for-each-ref` list:
   - a ref created by the call: `R5GIT update-ref -d <refname> <new id>`;
   - a ref moved by the call: `R5GIT update-ref <refname> <old id> <new id>`;
   - a ref deleted by the call: `R5GIT update-ref <refname> <old id>`.
   `update-ref` can rewrite `<C>/packed-refs`. The final recompute includes
   the trust paths, so a `packed-refs` that differs after this item is
   restored again from the freeze copy (item 2, no git command).

The reversion restores what R5 saved copies for, `config.worktree` included. The
other git-dir items of section 1 are detected and reported as
`not revertible: <path>`, never reverted. It does not restore what R5 does not name:
objects a `git tag -a` wrote, reflogs and the other live git-dir paths listed
in clone-isolation.md section 7 stay as the CLI left them, and the verdict is
already discarded.

## 5. The prompt (R4, R7)

The prompt carries the context as it does today: role, skill and bundle. A clone
call writes it to a file created by `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )`,
unique per call, inside the freeze directory: it is removed and swept together
with that directory, and no prompt file is ever age-swept on its own. (A call
without a clone, the `agy` review, creates its prompt with
`( umask 077; mktemp "<tmp root>/dh-prompt-XXXXXXXX" )` and removes it in its
own cleanup.) It
keeps the `commands.test`/`commands.lint` block from the anti-delegation
clause, and those commands run in the clone.

R4 is hygiene, not a boundary. The prompt must not contain the absolute path
of the live tree; check with `grep -F -- "<absolute path of the project>" <prompt file>`,
which must find nothing. The CLI can find that path by other means (the
comparisons above detect what it does with it), so this does not prevent a
write.

R7. An artifact that a legitimate test generates in the clone (an ignored file
or a new untracked file) is outside the judged set and does not discard the
verdict; a test that rewrites a tracked file discards it like any other change.
A check result a CLI reports is never evidence: a dispatched `qa-verifier`
reruns the checks. An ignored `_test.go` planted in the clone can only mislead the
CLI's own reply, which is already untrusted input. If a command fails for a
missing ignored file (a dependency, a `.env`), the CLI reports "not validated",
and the Coordinator reads that report as neither a failure of the code nor an
approval.

## 6. What the report carries

The "Workspace check" line of the Output format in `SKILL.md` takes, for a
call in a clone: `clean` or `reverted (<what changed, now reverted>)` for the
live tree; `clone removed: <path>` once the clone is gone; and
`notes: <new untracked files in the clone>` when section 3 found any; and
`leftover processes: stopped` or `none`, with the match string. An
altered ignored path is reported as `not revertible: <path>`.

The residual of clone-isolation.md section 7 applies unchanged: disposable
clone; writes outside it are detected on the live tree (R5), never prevented;
writes outside the project are not detected. It also holds for the freeze
directory. It contains copies of the `<harness dir>`, of `.harness/`, of `.claude/settings*.json` and of
the dirty or untracked non-ignored files, it sits in `<tmp root>`, and a
CLI that runs as the same user can read and alter it (mode 0700 does not stop
that). Tampering with the copies falls under the residual (a), a write outside
the project, and reading the copies falls under (e), reading and returning
secrets reachable by the user; the owner accepted both on 2026-09-29. The
freeze directory is removed at the end of the call, and swept as an orphan if
the call is interrupted (clone-isolation.md step 5).

Residual of the freeze: configuration files pulled in by `include.path` or `includeIf` from `<C>/config` that live inside the project are not restored before the first git command of part B; the `R5GIT` overrides do not neutralise them. Accepted residual of the live-tree freeze.
