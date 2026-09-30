---
type: regex
target: last_message
pattern: 'jest\s+--runInBand'
flags: ''
match: contains
weight: 1
---

The final reply states the literal string `jest --runInBand` (with any
run of whitespace between the two tokens, so `jest  --runInBand` or a
line-wrapped form still matches). The user's prompt asks for "the exact
command that runs behind the project's test script", naming neither a
runner nor a framework. That literal string exists only in
`fixtures/package.json` — it is not present anywhere in the user prompt
(verified by grep: `evals/cases/coordinator-explore-then-ask/prompt.md`
contains neither `jest` nor `runInBand`) and is not guessable without
reading the fixture. Stating it therefore proves exploration happened,
whether the reply performed the `Read` itself or received the fact back
from a dispatched subagent such as `repo-scout` — a `tool_used: Read`
grader on the top-level trace cannot see a subagent's own reads, so this
mechanical string match is used instead of a tool-trace assertion.
