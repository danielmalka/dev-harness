---
type: regex
weight: 1
target: last_message
match: contains
pattern: 'Set no longer evicts the oldest entry before inserting when the cache is already at capacity'
flags: ''
---
Deterministic check, sibling of `04-codex-finding-verbatim.md`: proves the
merged reply's own lru.go entry by matching, verbatim, the first sentence
of `cli:agy/gemini-3.1-pro-high`'s fixture reply
(`fixtures/cli-output-agy-cache-eviction.md`): "Set no longer evicts the
oldest entry before inserting when the cache is already at capacity."
Offline-validated against the kept run's final assistant message
(`/tmp/claude-eval-JyRG4o/out/trace.jsonl`): PASSES: a synthetic negative
that paraphrases the sentence (for example "the cache never frees its
oldest key once it is full") does not match.
