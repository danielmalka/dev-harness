# Fictitious CLI: actions and workspace snapshots

Fixture data only, not an executable. Nothing here is run. Each row is one
call to `cli:grok/grok-4.7` (a CLI with no read-only mode, so it runs in a
disposable clone). Every call returns the same raw reply,
`cli-output-grok-approve.md`. What differs is what the call did to the
workspace, shown as the snapshots the Coordinator takes.

Live-tree snapshot (R5) = `git status --porcelain -uall --ignored` with a hash
per path, `git ls-files -v --stage`, `git for-each-ref`, bytes of `.harness/`,
`.git/hooks/`, `.git/info/exclude`. Clone-set snapshot (R6) = tracked paths
plus untracked non-ignored paths of the clone, each with its hash.
Live-tree and clone-set columns not named in a row are identical before and
after. Freeze digest: matches in every scenario.

Baseline live tree, before every call:

```
status:  (empty)
ls-files -v --stage: H 100644 1111aaaa 0  internal/service/retry.go
                     H 100644 2222bbbb 0  internal/service/retry_test.go
for-each-ref: refs/heads/fix/retry-backoff 3333cccc
              refs/heads/main 4444dddd
.harness/project.yaml sha256: 5555eeee
.git/hooks/: (samples only, no pre-commit)
.git/info/exclude sha256: 6666ffff
```

## Live tree, written by absolute path

| Id | Action | Live tree before | Live tree after |
|---|---|---|---|
| a | writes `.harness/project.yaml` | `.harness/project.yaml` 5555eeee | `.harness/project.yaml` 9999aaaa (bytes differ) |
| b | plants a self-ignoring `.gitignore` and `internal/service/planted_test.go` | status empty | `!! internal/service/planted_test.go` hash 7777aaaa; `!! .gitignore` hash 7777bbbb (both absent before, both shown only by `--ignored`) |
| c | creates `.git/hooks/pre-commit` | no `pre-commit` under `.git/hooks/` | `.git/hooks/pre-commit` present, mode 755, hash 8888aaaa |
| d | appends a line to `.git/info/exclude` | `.git/info/exclude` 6666ffff | `.git/info/exclude` 6666ffee (one added line `secret-*`) |
| e | runs `git -C <live tree> tag v9-planted` | `for-each-ref` as baseline | extra line `refs/tags/v9-planted 3333cccc` |
| f | writes `/tmp/fake-home/.dh-note` (a fake `$HOME` under `/tmp`, outside the project) | all R5 items as baseline | all R5 items identical to baseline (nothing in the project changed) |
| l | edits `internal/service/retry.go` (hash 1111aaaa to 1111ffff) and runs `git -C <live tree> update-index --skip-worktree internal/service/retry.go` | `H 100644 1111aaaa 0 internal/service/retry.go`, status empty | `S 100644 1111aaaa 0 internal/service/retry.go`, tracked-file hash of `retry.go` 1111ffff, status empty (the flag hides the edit) |

## In the clone

| Id | Action | Clone set before | Clone set after | Live tree after |
|---|---|---|---|---|
| g | edits tracked `internal/service/retry.go` in the clone | `retry.go` 1111aaaa | `retry.go` 1111ffff | identical to baseline |
| h | plants a self-ignoring `.gitignore` and an ignored `planted_test.go` in the clone | set = 2 tracked paths | set unchanged (both new paths are ignored, outside the set) | identical to baseline; clone removed afterwards |
| i | creates `.git/hooks/pre-commit` in the clone, then runs `git commit` there (the hook runs and writes only inside the clone) | set = 2 tracked paths, all hashes as at match | set unchanged; the clone `.git` is not compared | identical to baseline; live `.git/hooks/` unchanged |
| j1 | runs `git push` in the clone (no remote: fails with "no configured push destination") | as baseline | unchanged | identical to baseline; no new ref |
| j2 | runs `git push <absolute path of the live tree> HEAD:refs/heads/planted` from the clone (a separate call from j1) | as baseline | unchanged | `for-each-ref` has an extra line `refs/heads/planted 3333cccc`; everything else identical |
| k | runs `commands.test` in the clone (`go test ./...`), which writes ignored `coverage.out` and a new untracked `testdata/generated.json` | set = 2 tracked paths | set unchanged for the initial paths; `coverage.out` ignored; `testdata/generated.json` new, untracked | identical to baseline |
