id: external-clis-parallel
asset: .agents/coordinator.md
type: positive
dimension: routing
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. `.harness/project.yaml` for this run carries
  `reviewers: {code: [claude, "cli:codex/gpt-5.6-sol"]}` (the fixture's
  project.yaml, narrowed to this one stage). Review diff.md.
  For the `cli:codex/gpt-5.6-sol` call, do not invoke a real binary: resolve
  it as if `command -v codex` succeeded, and read its reply from
  cli-output-codex-approve.md instead of running the CLI. Return the merged
  report only.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Per .agents/coordinator.md ("External CLI reviewers"), "Dispatch every
  entry of the stage's list" -> both the code-reviewer agent (claude) and
  the cli:codex/gpt-5.6-sol call are dispatched, each counted as one
  concurrent specialist against the two-specialist cap (Procedure item 5).
  Per "Reviewers of one stage run in parallel; there is no sequential
  mode. Produce the merged report only after every reviewer of the stage
  has returned or been marked not-run", the merged report is produced only
  once both the claude dispatch and the codex double have a result; the
  codex double resolves to `approve` from cli-output-codex-approve.md, so
  the merged verdict is approve with zero blocking findings.
  The Dispatches section (coordinator output format) logs the codex entry
  with its resolved binary and slug in place of model fields.
  Must not: dispatch the codex call, wait for its result, then dispatch
  claude afterward (or the reverse) — the two calls are concurrent, not
  sequential; produce a merged report before both results are in hand;
  invoke a real codex/claude binary or any command outside reading the
  fixture doubles; treat the codex entry as a Claude specialist in the
  Dispatches log.
rubric: evals/rubrics/external-clis.md
