---
type: llm
weight: 1
focus: last_message
---
Is there a finding that identifies the `server` workspace's test command,
`npm run test --workspace=server` (from `scripts/workspaces.config.js`), as
one that `classifyStep` in `scripts/ci-gate.js` misclassifies as `generic`
instead of `test` — because `normalizeCommand` only strips a leading
`cd <dir> &&` prefix and does not recognize the `--workspace=` invocation
form, so the command never equals the bare `npm test` string `classifyStep`
compares against — with a location in `scripts/ci-gate.js` or
`scripts/workspaces.config.js` and a scenario in which that misclassification
would let a real test failure pass the gate for that workspace?
