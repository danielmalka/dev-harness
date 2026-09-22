id: external-clis-attribution
asset: .agents/coordinator.md
type: positive
dimension: evidence
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. `.harness/project.yaml` for this run carries
  `reviewers: {code: ["cli:codex/gpt-5.6-sol", "cli:agy/gemini-3.1-pro-high"]}`.
  Review diff.md.
  For the `cli:codex/gpt-5.6-sol` call, do not invoke a real binary: read
  its reply from cli-output-codex-major-duplicate.md. For the
  `cli:agy/gemini-3.1-pro-high` call, read its reply from
  cli-output-agy-blocking.md. Neither binary is actually invoked. Return
  the merged report only.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Both doubles flag internal/service/retry.go:44 as a Major, worded
  differently. Per .agents/coordinator.md ("External CLI reviewers"),
  "Equivalent findings from different reviewers (same file and line, or
  same document section) are listed once with every reporter attributed
  (AC-13)" and "Each finding in that record carries the tag of the
  reviewer that produced it, `claude` or `cli:<binary>/<slug>`" (AC-09) ->
  the merged report lists internal/service/retry.go:44 exactly once, under
  one Major entry, with both `Reported by: cli:codex/gpt-5.6-sol` and
  `Reported by: cli:agy/gemini-3.1-pro-high` present against that one
  entry (or a single `Reported by:` line naming both).
  Must not: list the finding twice, once per reviewer; keep only one
  reviewer's attribution and drop the other; merge the two differently
  worded findings into a new, reworded third finding with no reporter
  tags; invoke a real codex or agy call.
rubric: evals/rubrics/external-clis.md
