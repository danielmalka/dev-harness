# Case notes: external-clis-distinct-findings

Proves the negative counterpart `external-clis-attribution`'s own positive
merge coverage does not test: two genuinely unrelated findings from two CLI
reviewers, at different files and lines with no plausible overlap, stay as
two separate, correctly attributed entries rather than being merged into
one.

## Paid run 1 (2026-09-28)

`--runs 1`, model and judge sonnet. Result: graders 01-03 PASS 3/3; grader
04 FAIL 3/3. Behavior was correct — this was a case defect. Trace read
manually by the Coordinator (`/tmp/claude-eval-JyRG4o/out/trace.jsonl`):
the merged reply reproduces both findings verbatim — the retry.go:44 entry
carrying `cli:codex/gpt-5.6-sol`'s own text and `Reported by:` tag, the
lru.go:87 entry carrying `cli:agy/gemini-3.1-pro-high`'s own text and tag —
exactly the behavior this case exists to prove.

## Correction 1/6 (case directory only)

Grader 04 (`04-no-paraphrase-or-fusion.md`, `type: llm`) asked "Is neither
entry's text a paraphrase or fusion...?" — the same double-negation trap
already seen in PLAN-007: a judge answering the plain fact ("no, it is not
a paraphrase or fusion") reads its own "no" as a negative verdict and
scores the criterion FAIL, even though "no" was the correct answer to the
question as posed.

Fixed by deleting `04-no-paraphrase-or-fusion.md` and replacing it with two
deterministic graders, `type: regex`, `target: last_message`,
`match: contains`, one per reviewer's own finding:

- `04-codex-finding-verbatim.md` — matches, verbatim, the first sentence of
  `cli:codex/gpt-5.6-sol`'s fixture reply ("The final failed attempt now
  returns immediately with no backoff sleep").
- `05-agy-finding-verbatim.md` — matches, verbatim, the first sentence of
  `cli:agy/gemini-3.1-pro-high`'s fixture reply ("Set no longer evicts the
  oldest entry before inserting when the cache is already at capacity").

Both patterns contain no regex metacharacters needing escape. Validated
offline against the kept run's final assistant message
(`/tmp/claude-eval-JyRG4o/out/trace.jsonl`): both PASS. A synthetic
paraphrased negative for each sentence (codex: "the last retry attempt
skips its backoff delay entirely"; agy: "the cache never frees its oldest
key once it is full") does not match either pattern.

Pass bar is now five graders (01, 02, 03, 04, 05), all must pass. No `dh
build` and no paid eval run made for this correction; `go run ./cmd/dh
validate --source-only .` re-run clean (exit 0, `errors: 0`, with
`pipefail`).
