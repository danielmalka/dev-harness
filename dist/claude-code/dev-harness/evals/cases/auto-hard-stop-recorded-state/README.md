# Case notes: auto-hard-stop-recorded-state

Proves PRD-006 RF-06/RF-07 (AC-11/AC-12) and PLAN-010 T-704 item 2: on
reaching a hard stop, `/dh:auto` ends without a wait loop or automatic
retry, and records the exact interrupted state — named identically — in
both `.harness/tasks/AUTO-<n>/BRIEF.md`'s `Stage:` line and
`.harness/MEMORY.md`.

## RF-06 condition chosen, and why

RF-06 lists seven hard-stop conditions. Per the plan
(`.harness/tasks/PLAN-010/PLAN.md`, Slice 4), "the cheapest to fixture
deterministically is the six-round-cap exhaustion or an
unauthorized-dependency-install attempt." This case uses
**unauthorized dependency install**, not the six-round cap, because:

- The six-round cap needs six real correction rounds (build, QA, and
  independent review disagreeing six times over) to exhaust deterministically
  — a long, expensive, multi-dispatch scenario to fixture and grade in one
  cheap `--runs 1 --max-cost-usd 0.75` eval run.
- An unauthorized dependency install is reachable in a single coordinator
  turn: the fixture's one-slice `PLAN.md` states, as a plain project fact,
  that implementing Slice 1 (`T-1`) needs the `qrcode` npm package, and
  `fixtures/package.json` independently confirms it is not currently a
  dependency. The `/dh:auto` seven-item questionnaire (`.commands/auto.md`
  Routing) has no item that could authorize a dependency install in the
  first place, so this condition is reachable the moment the coordinator
  reads the plan, with no subagent dispatch strictly required (grader 02
  passes on zero or one builder dispatch, never more than one).
- `Bash` is left out of `allowed_tools` entirely, and the fixture ships no
  real package manager or registry — so even a model that tried to install
  the dependency anyway could not actually reach the network or mutate
  anything outside the scored transcript. This directly satisfies the
  task's own constraint that "fixtures must not trigger anything
  destructive or network."

## Design choices shared with the sibling case

- `fixtures/project.yaml`, `fixtures/brief.md`, `fixtures/memory.md` stand
  in for `.harness/project.yaml`, `.harness/tasks/AUTO-1/BRIEF.md`, and
  `.harness/MEMORY.md`, following the same "`.harness/project.yaml` for
  this run is fixtures/project.yaml" pattern already used by
  `external-clis-document-only`'s prompt, and the same assumption recorded
  in `auto-questionnaire-before-action/README.md` about how a `regex`
  grader with `target: {source: file, path: ...}` reads the run's own
  workspace copy of that path after the run completes.
- The seven-item questionnaire is pre-answered in the single user turn
  (mode (a), "decide and record everything," chosen here so the run heads
  straight for the build stage without a second product/contract
  decision), for the same reason recorded in the sibling case's README:
  no verified `context.history_file` transcript format was available to
  script a genuine second turn under this task's authorization (no
  `claude plugin eval` run permitted). This case is about the RF-06/RF-07
  stop behavior, not the questionnaire itself (that is
  `auto-questionnaire-before-action`'s job), so a pre-answered single turn
  does not weaken what this case measures.
- Grader 05 (`llm`, `focus: {source: file, path: fixtures/memory.md}`)
  reads content from a file rather than the transcript, matching the
  plan's requirement that this be "a mechanical check on the file, not on
  the transcript" for existence/identity, with the adequacy-of-the-recorded-
  reason left to the one genuinely semantic judgment in this case. Note for
  a later maintainer: `internal/kit/evals.go`'s `checkCaseGraders` only
  statically resolves a `target.path` (used by `regex`/graders 03/04/06),
  not a `focus.path` (used by `llm` graders); `go run ./cmd/dh validate`
  will not catch a typo in grader 05's path the way it would for 03/04/06.
  The path was checked by hand against the fixture tree above.
