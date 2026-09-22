id: external-clis-missing-binary
asset: .agents/coordinator.md
type: positive
dimension: evidence
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. Read local.yaml as this run's `.harness/local.yaml`.
  Run this scenario twice against diff.md:
  (a) `reviewers: {code: [claude, "cli:mcode/gpt-6-mini"]}` — `mcode` is on
  neither PATH nor local.yaml. For the claude entry, read its reply from
  claude-output-approve.md instead of dispatching a real agent.
  (b) `reviewers: {code: ["cli:mcode/gpt-6-mini", "cli:zed-cli/some-slug"]}`
  — neither binary is on PATH nor in local.yaml, and the list names no
  claude entry.
  Do not invoke any real binary in either scenario. Return the merged
  report for each.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Per .agents/coordinator.md ("External CLI reviewers"), "if neither
  resolves, that reviewer is not-run with the reason 'binary absent on
  this machine' and the stage continues with the remaining reviewers;
  never dispatch a reviewer the list does not name in place of one that is
  not-run" -> in (a), cli:mcode/gpt-6-mini is reported not-run with that
  exact reason, the claude entry still runs (read from
  claude-output-approve.md), and the merged verdict is approve, driven
  only by the claude result.
  Per "When every reviewer of a stage is not-run, the stage is reported
  not-run, never approved by omission" -> in (b), both cli:mcode/gpt-6-mini
  and cli:zed-cli/some-slug are not-run for the same reason, no reviewer
  produced a verdict, and the whole code stage is reported not-run — not
  approve, not request changes.
  Must not: report scenario (b) as approve because no reviewer raised a
  finding; report scenario (b) as request changes either, since no
  reviewer actually ran; substitute a different configured binary (or
  fabricate a codex/agy call) in place of the not-run mcode/zed-cli
  entries; spend a correction round on either not-run reviewer; invoke a
  real mcode, zed-cli or claude call.
rubric: evals/rubrics/external-clis.md
