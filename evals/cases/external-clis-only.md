id: external-clis-only
asset: .agents/coordinator.md
type: positive
dimension: routing
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. `.harness/project.yaml` for this run carries
  `reviewers: {code: ["cli:codex/gpt-5.6-sol"]}` — no `claude` entry.
  Review diff.md.
  For the `cli:codex/gpt-5.6-sol` call, do not invoke a real binary: resolve
  it as if `command -v codex` succeeded, and read its reply from
  cli-output-codex-approve.md instead of running the CLI. Return the merged
  report only.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Per .agents/coordinator.md ("External CLI reviewers"), "When a stage
  lists no `claude`, no Claude agent is dispatched for it and the merged
  report comes only from the CLI verdicts" -> the code-reviewer agent is
  never dispatched; the Dispatches section names only the codex entry, and
  the merged verdict (approve, zero blocking findings) is built solely from
  cli-output-codex-approve.md.
  Must not: dispatch code-reviewer (claude) "just to have a Claude opinion"
  even though the list does not name it; log a claude dispatch in the
  Dispatches section; treat the absence of a claude entry as an error or
  fall back to a Claude-only flow; invoke a real codex binary.
rubric: evals/rubrics/external-clis.md
