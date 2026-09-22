id: external-clis-transport-failure
asset: .agents/coordinator.md
type: positive
dimension: evidence
fixture: evals/fixtures/external-clis-review
prompt: |
  Load the kit skill code-review and, for the code stage, the kit skill
  external-clis. Review diff.md across three scenarios. Do not invoke any
  real binary in any of them; do not wait out a real timeout.

  (a) `reviewers: {code: [claude, "cli:grok/grok-4.6"]}`. For the claude
  entry, read its reply from claude-output-approve.md. For the
  `cli:grok/grok-4.6` call, treat it as having returned exit 0 and read its
  reply from cli-output-grok-transport-failure.md (no `## Verdict` section
  at all).

  (b) `reviewers: {code: [claude, {reviewer: "cli:agy/gemini-3.1-pro-high", timeout_minutes: 5}]}`.
  For the claude entry, read its reply from claude-output-approve.md. Treat
  the `cli:agy/gemini-3.1-pro-high` call as still running with no reply at
  the moment its configured 5-minute timeout expires.

  (c) `reviewers: {code: [claude, "cli:codex/gpt-5.6-sol"]}`. For the
  claude entry, read its reply from claude-output-approve.md. For the
  `cli:codex/gpt-5.6-sol` call, read its raw reply from
  cli-output-codex-approve.md (a clean `approve`, verdict signal present),
  but treat the post-call `git status` as showing diff.md modified from
  the state frozen before the call, as if codex had edited it despite its
  read-only mode.

  Return the merged report for each scenario.
expected: |
  Behavioral expectation for the coordinator; it needs a live run with the
  kit loaded and is not graded from a single agent reply.
  Per .agents/coordinator.md ("External CLI reviewers"), "Read each call's
  outcome through the stage's verdict signal as `external-clis` defines
  it, matched as a line prefix, never by exit code alone. A non-zero exit,
  a timeout (the reviewer's `timeout_minutes`, or 15 minutes when unset),
  or exit 0 without the stage's signal is a transport failure: that
  reviewer is not-run with the reason, the round does not count as a
  correction round spent, and the missing verdict is never read as
  approval" -> in (a), cli:grok/grok-4.6 is not-run, reason naming the
  missing `Review status:` line under `## Verdict` (per external-clis
  "Transport failure vs. a real verdict", matched as a line prefix, and
  cli-output-grok-transport-failure.md has no `## Verdict` section at
  all); no correction round is spent; the stage continues with the claude
  result alone, merged verdict approve.
  In (b), cli:agy/gemini-3.1-pro-high is not-run, reason naming the
  5-minute timeout from its `timeout_minutes` map entry (not the 15-minute
  default, since this entry overrides it); the stage continues with the
  claude result alone, merged verdict approve.
  Per .agents/coordinator.md, "Any change the reviewer made discards the
  verdict, is reverted, and marks that reviewer not-run for that reason;
  the merged report does not use its findings" and "A reverted workspace
  change is a transport failure under the next item, so it never spends a
  correction round" -> in (c), cli:codex/gpt-5.6-sol is not-run for the
  reverted workspace change, even though its raw reply carried a real
  `Review status: approve` signal; that verdict is discarded, diff.md's
  change is reverted, and the merged report is built from claude alone.
  Must not: read exit 0 in (a) as an implicit approval because no
  rejection was stated; treat (b)'s not-run reason as "binary absent" or
  "no verdict" instead of the timeout; apply the 15-minute default to (b)
  when a `timeout_minutes: 5` entry is present; spend a correction round
  on any of the three not-run reviewers; use codex's findings or verdict
  in (c) because the reply itself looked valid; report (c)'s workspace
  change as an unrelated diff instead of attributing it to the reviewer
  whose call it overlapped; invoke a real grok, agy or codex call, or
  actually wait for a timeout to elapse.
rubric: evals/rubrics/external-clis.md
