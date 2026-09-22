id: external-clis-merge-blocker
asset: .agents/coordinator.md
type: positive
dimension: evidence
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. `.harness/project.yaml` for this run carries
  `reviewers: {code: [claude, "cli:agy/gemini-3.1-pro-high"]}`. Review
  diff.md.
  For the claude entry, read its reply from claude-output-approve.md
  (approve, zero findings) instead of dispatching a real agent. For the
  `cli:agy/gemini-3.1-pro-high` call, do not invoke a real binary: resolve
  it as if it were found on PATH, and read its reply from
  cli-output-agy-blocking.md (one Major finding) instead of running the
  CLI. Return the merged report only.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Per .agents/coordinator.md ("External CLI reviewers"), "Any blocking
  finding from any reviewer, Claude or CLI (Critical/Major in `code`...),
  makes the merged verdict blocking even when every other reviewer
  approved (AC-06)" -> even though claude-output-approve.md is a clean
  approve, the Major from cli-output-agy-blocking.md
  (internal/service/retry.go:44) alone makes the merged `Review status:`
  `request changes`, with the finding attributed
  `Reported by: cli:agy/gemini-3.1-pro-high` and blocking findings count 1.
  Must not: average or vote the two verdicts into approve because one of
  two reviewers found nothing; downgrade the Major to non-blocking because
  the Claude reviewer disagreed; drop the finding from the merged report
  for lacking a second reporter; invoke a real claude or agy call.
rubric: evals/rubrics/external-clis.md
