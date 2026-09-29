---
type: regex
weight: 1
target: last_message
match: contains
pattern: 'The final failed attempt now returns immediately with no backoff sleep'
flags: ''
---
Deterministic check, replacing a `type: llm` judge (`04-no-paraphrase-or-fusion.md`)
that failed 3/3 on a verifiably correct run — the same double-negation trap
already seen in PLAN-007: the criterion asked "is neither entry's text a
paraphrase or fusion", and a judge answering the plain fact ("no, it is not
a paraphrase") reads as the literal word "no" and gets scored as a FAIL
rather than a PASS. The merged reply's own retry.go entry is proven instead
by matching, verbatim, the first sentence of `cli:codex/gpt-5.6-sol`'s
fixture reply (`fixtures/cli-output-codex-retry-backoff.md`): "The final
failed attempt now returns immediately with no backoff sleep." Offline-
validated against the kept run's final assistant message
(`/tmp/claude-eval-JyRG4o/out/trace.jsonl`): PASSES: a synthetic negative
that paraphrases the sentence (for example "the last retry attempt skips
its backoff delay entirely") does not match.
