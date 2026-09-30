# Clone verification results (R2, R3)

Fixture data only. Nothing here is run; each block is the printed result of a
check the Coordinator would run for one call. `<tmp root>` resolved to `/tmp`
(`TMPDIR` unset). Live tree `<project>`. Both calls are `cli:grok/grok-4.7`
review calls with `<call-id>` = `T-2001-grok-1`.

## Call X: every check matches

```
clone path: /tmp/dh-cli-T-2001-grok-1-Q7mZ2pLa
git -C <clone> remote                          -> (empty)
git -C <clone> rev-parse --git-common-dir      -> /tmp/dh-cli-T-2001-grok-1-Q7mZ2pLa/.git
test ! -e <clone>/.git/objects/info/alternates -> exit 0
test ! -e <clone>/.git/logs                    -> exit 0
judged set (live == clone): 3 tracked paths, 1 untracked non-ignored path
per-path hash and mode, live vs clone          -> all equal
git status --porcelain -uall, live             -> " M internal/service/retry.go" / "?? notes.txt"
git status --porcelain -uall, clone            -> " M internal/service/retry.go" / "?? notes.txt"
test ! -e <clone>/.harness                     -> exit 0
test ! -e <clone>/.claude/settings.local.json  -> exit 0
```

## Call Y: the match diverges

```
clone path: /tmp/dh-cli-T-2001-grok-2-Hd4nT8wc
remote / git-common-dir / alternates / logs    -> all as in call X
per-path hash: internal/service/retry.go       -> live 1111ffff, clone 1111aaaa
git status --porcelain -uall, live             -> "M  internal/service/retry.go" (staged)
git status --porcelain -uall, clone            -> " M internal/service/retry.go" (unstaged)
```
