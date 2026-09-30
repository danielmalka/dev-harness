# Direct children of `<tmp root>` (`/tmp`) before the first dispatch of a batch (R9)

Fixture data only. Current time is 2026-09-30T11:00:00Z; the call timeout is the
15-minute default. `kill -0` results are given, not run. Current user: `malka`.

| Directory | Owner | Marker `.git/dh-cli.lock` | `kill -0 <pid>` | `started` |
|---|---|---|---|---|
| `/tmp/dh-cli-T-1900-grok-1-Ab3dE9xQ` | malka | yes, `pid=4242` | fails: no such process | 2026-09-30T09:00:00Z |
| `/tmp/dh-cli-T-1901-grok-1-Pq8sV2mN` | malka | yes, `pid=1877` | succeeds (process alive) | 2026-09-30T09:00:00Z |
| `/tmp/scratch-build-Zk29aB1c` | malka | yes, `pid=4243` | fails: no such process | 2026-09-30T09:00:00Z |
| `/tmp/dh-cli-T-1902-opencode-1-Lm5tR7yU` | malka | no marker | not applicable | not applicable |
| `/tmp/dh-cli-T-1903-grok-1-Xc6wE4jH` | root | yes, `pid=4244` | fails: no such process | 2026-09-30T09:00:00Z |

All five are real directories, not symlinks.
