# Two clone calls dispatched in one group (R10)

Fixture data only. Task `T-2001`, `<tmp root>` = `/tmp`. Two review calls for
the same `code` stage overlap in time.

| Call | Binary/slug | `<call-id>` | Clone path (stored literal) |
|---|---|---|---|
| 1 | `cli:grok/grok-4.7` | `T-2001-grok-1` | `/tmp/dh-cli-T-2001-grok-1-Q7mZ2pLa` |
| 2 | `cli:mcode/mcode-2` | `T-2001-mcode-2` | `/tmp/dh-cli-T-2001-mcode-2-Rt9kW3sB` |

Call 1 ends first (exit 0, verdict present); call 2 ends 40 seconds later
(exit 0, verdict present). Neither clone's set changed (R6 clean for both). The
live-tree freeze for the group was taken once before call 1 started.
