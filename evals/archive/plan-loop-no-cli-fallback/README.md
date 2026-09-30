# Case notes: plan-loop-no-cli-fallback

Proves PRD-005 RF-02/AC-02/AC-03: with zero resolvable planner entries in
`reviewers.document` (`claude` and a `cli:claude/<slug>` entry, both
excluded — `claude` as itself, `cli:claude/<slug>` under AA1),
`/dh:plan-loop` states plainly that no loop is possible, runs `/dh:plan`'s
own flow inline, and creates no `wave-N/` directory.

Pass bar: six graders (`graders/01`-`graders/06`), every grader must pass,
no partial credit.

## Paid run 1 (2026-09-28)

Cost: US$0.85 agent + US$0.25 judge. Result: graders 02, 04, 05 PASS;
graders 01 and 03 FAIL. Both failures were case defects — the underlying
behavior was correct. Trace, read manually by the Coordinator:

- First assistant text: "fixture reviewer list gives k = 0. That means the
  no-CLI fallback…" — stated before any dispatch.
- Exactly one `Agent` call, `subagent_type: "dh:implementation-planner"`.
- Then one `Agent` call, `subagent_type: "dh:document-validator"`.
- No `wave-N` write anywhere in the trace.

## Correction 1/6 (run 1, both failures fixed in the case directory only)

1. **Grader 03 (`03-planner-dispatch-bounds.md`).** The runtime serializes
   the plugin-qualified name, `dh:implementation-planner`, not the bare
   role name `implementation-planner` the grader's `input_match` was
   anchored to. Fixed by widening the anchor to
   `"subagent_type":[ \t]*"(dh:)?implementation-planner"` — `min: 1`/
   `max: 3` unchanged. Offline-checked (Python `re`, the documented proxy
   for the runner's JS-flavored engine) against
   `"subagent_type": "dh:implementation-planner"` (matches, positive) and
   `"subagent_type": "dh:document-validator"` (does not match, negative),
   plus the bare (unprefixed) forms of both, and an unrelated role — six
   checks, all as expected.

2. **Grader 01 (renamed `01-fallback-stated-before-plan-content.md` ->
   `01-fallback-stated-explicitly.md`).** `focus: trace` gives an `llm`
   grader only the tool-call trace, not the assistant's own reply text —
   so a grader asking whether a *prose statement* precedes a *tool call*
   cannot be judged from `focus: trace` at all; the judge simply has no
   access to the statement it needs to place in that order. Run 1's own
   trace shows the correct behavior (the fallback notice was in fact the
   first assistant text, before any dispatch), but grader 01 could not see
   that text under `focus: trace` and failed for a reason unrelated to
   what actually happened.

   Fixed by rewriting the grader to `focus: last_message`, judging a
   narrower, last-message-answerable claim: does the final reply
   explicitly state that no loop was possible, name the actual reason
   (k = 0: `claude` and `cli:claude/<slug>` both excluded), and say that
   `/dh:plan`'s flow ran instead. **The "stated before any plan content"
   ordering claim is dropped from this grader — it is not judgeable by
   any `llm` grader in this kit's eval schema**, since no `focus` value
   exposes both the assistant's prose and a verifiable position relative
   to a tool call in the same judged text. That ordering was verified
   manually on run 1's trace (see above) rather than by an automated
   grader, and is recorded here rather than re-asserted by a check this
   schema cannot express.

   Checked grader 04 (`04-validator-dispatched-rounds-capped.md`) for the
   same limitation, since it also uses `focus: trace` and passed on run 1:
   its question is answerable entirely from tool calls (which role was
   dispatched, how many times) with no reference to the reply's own prose
   accounting, so `focus: trace` is the right scope for it. Tightened its
   wording anyway, dropping the earlier phrase "whether stated in the
   reply or shown by the trace" (which implied it might read prose it
   cannot see) in favor of a count taken only from the visible
   `implementation-planner` `Agent` dispatches in the trace.

No other file in this case directory changed. No `dh build` or paid eval
run was authorized or made for this correction; `go run ./cmd/dh validate
--source-only .` was re-run clean (exit 0, `errors: 0`).

## Correction 2/6 (T-803 code review, 1 Major)

`code-reviewer` found that no grader in the case checked that
`PLAN.review.md` is actually written to disk, although AC-03 (and T-803
§7) names `PLAN.md`, its `TASK.md` file(s), and `PLAN.review.md` as the
exact artifact set the fallback persists — a run that drafted and
validated the plan but skipped persisting the validator's review would
have passed all five prior graders. Grader 05 only reads the final
reply's own description of what it wrote (`focus: last_message`), so it
cannot catch a run that fails to write `PLAN.review.md` while still
*claiming*, or simply not contradicting, that it did.

Fixed by adding `graders/06-plan-review-written.md`, `type: regex`,
`target: trace`, `match: contains`, `pattern:
'"file_path":[ \t]*"[^"]*PLAN\.review\.md"'` — anchored on the serialized
`file_path` key, the same convention as grader 02, so a bare mention of
the filename in a reply's own prose cannot satisfy it. Pass bar is now
six graders, all must pass.

Validated offline, no paid run made for this correction:

- **Positive**: the pattern matched the real `Write` tool call in the kept
  run-2 trace (a `trace.jsonl` kept outside this repository, in the
  runner's own throwaway output directory), line 128, JSON `type:
  assistant` -> `tool_use` `Write` whose `file_path` ends in
  `.harness/tasks/PLAN-1/PLAN.review.md` (the scaffolded workspace's own
  absolute prefix precedes it, as it does for every path in that trace).
- **Negative**: a synthetic trace line with a `Write` `tool_use` whose
  `file_path` ends in `PLAN.md` only (no `PLAN.review.md`) did not match.

`go run ./cmd/dh validate --source-only .` re-run clean after adding the
new grader file (exit 0, `errors: 0`, with `pipefail`). No `dh build`, no
paid eval run made for this correction.
