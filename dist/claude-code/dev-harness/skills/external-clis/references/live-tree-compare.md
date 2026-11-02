# Live-tree comparison without a clone

This file is the freeze, compare and revert procedure for a CLI call that runs on the live tree, without a disposable clone: `codex` under its read-only flag, and the `agy` review. A call in a clone uses [live-tree-freeze.md](live-tree-freeze.md) instead. The comparison is not a sandbox: it only detects a change after the fact. A reviewer is read-only by contract even when its binary has no read-only flag.

## Quiet period

Nothing else may write to the tree or to the `<harness dir>` between the freeze and the comparison, the Coordinator and the delegate included. A build, a package regeneration or a test run started in parallel changes files the comparison then blames on the reviewer. Measured on 2026-09-22: a `dh build` run by the Coordinator while two CLI reviewers were working made the comparison differ on `harness-manifest.json` and `dist/`, which by the rule below discards a verdict the reviewer never threatened. Either wait, or run the parallel work in a separate checkout.

## Freeze (before the call)

- `git status --porcelain -uall` for the path list, then `git hash-object "<path>"` for every listed path (modified or untracked). A path absent on disk records the marker `deleted` in place of the hash.
- The bytes of every file under the `<harness dir>` (mode `global`: `<home>/projects/<name>/`, outside the tree, and the Coordinator writes nothing there during the call either; mode `repo`: `.harness/`, or the main checkout's `.harness/` when the project is a linked worktree), under `.harness/` when it exists, and of every file under `.git/hooks/`, of `.git/config`, and of `.claude/settings.json` and `.claude/settings.local.json` when present. Git ignores them, so the path list cannot see a change there, yet they hold the commands, reviewer list and permissions later steps trust.
- The repository-discovery items: the bytes of the gitfile (or the fact that `.git` is a real directory) and, in a linked worktree, the bytes of `<admin>/commondir` and `<admin>/gitdir`, plus the recorded absence of `<C>/commondir` and `<project>/.git/commondir`.

## Compare and revert (after the call)

Recompute all of the above. Any difference is the failure kind `tree_changed`: a changed path list, a changed hash for any path, or in the byte snapshot a changed file, a file created during the call (absent to present) or a file deleted. The verdict is discarded, the change is reverted, and the reviewer is not-run for that reason, whatever it said. This also catches an edit made inside a file that was already dirty in the frozen state, which the path list alone cannot distinguish.

Revert every differing path to its exact pre-call state:

- A path that was listed (dirty or untracked) or snapshotted before the call: the freeze keeps a copy of its pre-call bytes and mode outside the checkout, and revert writes those back. Never use `git checkout` or `git restore` for such a path: that returns it to `HEAD` and destroys changes already in the tree before the call, and git cannot restore an ignored file at all.
- A path that was clean and tracked before the call (in neither set): `git checkout -- "<path>"` is the correct restore, because `HEAD` is exactly its pre-call state.
- A path that did not exist before the call: delete it.
- Never write through the current directory entry: remove the path first (without following a symlink) and recreate it, so a symlink planted during the call is removed rather than followed.

## Linked worktree

- The `agy` live-tree review ends `isolation_unavailable` when `<project>/.git` is not a directory: the review does not run and no review round is spent, because this freeze snapshots `.git/hooks/` and `.git/config`, which are not folders in a gitfile checkout. The trigger is `.git` not being a directory, a deliberate superset of a linked worktree (it also covers a submodule and a bare repository, which lack those folders too).
- For `codex` from a linked worktree the procedure is otherwise unchanged, but the trust paths (`hooks/`, `config`, `config.worktree` when present, `info/`, `HEAD`, `packed-refs`, and the admin dir `config.worktree`, all saved and restored with non-git tools before any git command) are read from the common dir `<C>` (`linked_worktree.commondir` from `dh harness-path --json "<project>"`), not from the gitfile at `<project>/.git`. A path absent in both snapshots is not a difference, and the read-only sandbox is the guard.
- Before any git command of the comparison or the reversion, compare by bytes with non-git tools (`cat`, `cmp`, `stat`) the repository-discovery items saved at the freeze (an item recorded as absent must still be absent). Restore or remove first. If that is impossible, run no git command and record the verdict as discarded, not-run, `not revertible`. Every such byte comparison and restoration comes before the first git command.
