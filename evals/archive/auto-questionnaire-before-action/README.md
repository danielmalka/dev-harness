# Case notes: auto-questionnaire-before-action

Proves PRD-006 RF-03/RF-04 (AC-06/AC-07) and PLAN-010 T-704 item 1: the
seven-item `/dh:auto` questionnaire is shown before any file write or
specialist dispatch, and its answers are persisted to
`.harness/tasks/AUTO-<n>/BRIEF.md` and `.harness/MEMORY.md` before the
first stage is dispatched.

## Design choices, recorded for a later maintainer

- **Single-turn, pre-answered fixture, not a scripted two-turn
  conversation.** AC-07 ("quando o questionário for respondido...") is
  naturally a two-turn scenario: turn 1 the model asks, turn 2 the owner
  answers, only then are the records written. The `claude plugin eval`
  binary's case schema does expose a `context.history_file` +
  `execution.prompt` mechanism for resuming a prior transcript with a next
  turn (confirmed by reading strings out of the installed CLI binary,
  `strings <claude binary> | grep history_file`: `"context.history_file
  requires execution.prompt (the resumed session needs a next turn)"`).
  No case in this repository uses that mechanism yet, the harness-evaluation
  skill documents the field but shows no worked example, and its exact
  history-file content format (a session JSONL transcript? something
  simpler?) could not be confirmed without an actual `claude plugin eval`
  run — which this task was not authorized to make. Building a case on an
  unverified transcript format risked a silently-broken case that
  `dh validate` cannot catch (it only checks that `context.history_file`
  resolves to an existing path, never its content).
  Instead, the fixture's single user message gives the free-text idea and
  pre-supplies literal answers to all seven items, explicitly telling the
  assistant it need not wait for a further reply. This keeps the case a
  single, ordinary `prompt.md` turn (matching the structural precedent of
  `coordinator-explore-then-ask` and `external-clis-only`) while still
  requiring the seven-item confirmation to reach the persisted brief
  (grader 01) and the answers to be persisted before any stage dispatch
  (graders 02/03), which is the behavior RF-03/RF-04 actually gate. If a
  future maintainer
  confirms the `history_file` format from a real run, a true two-turn
  version of this case (owner does not pre-answer) would be a stronger,
  additional test — not a replacement, since the pre-answered turn still
  proves the write-before-dispatch ordering that a bare "asked and
  stopped" single turn cannot.
- **`fixtures/brief.md` / `fixtures/memory.md` stand in for
  `.harness/tasks/AUTO-1/BRIEF.md` / `.harness/MEMORY.md`.** This mirrors
  the exact pattern already used by `external-clis-document-only`'s prompt
  ("`.harness/project.yaml` for this run is fixtures/project.yaml"): the
  prompt tells the model directly which fixture path stands in for which
  real path, rather than nesting a `.harness/` tree inside `fixtures/`.
  Graders 04/05 use a `regex` grader with `target: {source: file, path:
  fixtures/brief.md|memory.md}` — per `.skills/harness-evaluation/SKILL.md`
  and the CLI's own case schema (`internal/kit/evals.go`'s
  `checkCaseGraders`/`checkFieldPath` requires this path to already resolve
  to a real file inside the case directory, which is why both fixtures
  start as placeholder stubs rather than being absent). This assumes the
  grading step reads that same relative path back from the run's own
  scaffolded workspace after the run completes (the skill's own wording:
  "a grader reading a file the run produced") — the most natural reading,
  but likewise unverified by an actual run under this task's authorization.
- **The eval spend cap answer, `US$7.25`, is deliberately odd and unique.**
  It appears nowhere else in the fixture or prompt text, so its presence in
  the persisted files (graders 04/05) is strong, hard-to-fake evidence that
  the specific owner-given answer — not an invented placeholder value —
  reached disk.
- `AUTO-1` is assumed as the task ID: the fixture has no pre-existing
  `AUTO-*` task directory, so it is the first available `AUTO-<n>` under
  RF-15's numbering.
- **Correction 1/6 (paid run r2, `--allow-tools Write`):** grader 01
  originally used `focus: last_message`. In an actual run the coordinator
  confirmed and persisted the seven answers first, then legitimately kept
  going within the same turn and ended on a mode-(b) product question ("no
  receipts service exists in this checkout... (a)/(b)/(c)") — a fully
  compliant continuation this prompt's own mode-(b) pre-answer invites, but
  one whose *last* message no longer carries the seven-item confirmation
  text. Graders 02-05 passed 3/3; only grader 01 failed, on the case, not
  the behavior. Fixed by pointing grader 01 at the persisted file instead
  of the last message (`focus: {source: file, path: fixtures/brief.md}`),
  matching graders 04/05: it now asks whether the persisted brief records
  all seven items, which is stable regardless of how many further turns
  the reply takes afterward.
- **`claude plugin eval` needs `--allow-tools Write Edit` for this case.**
  Without it, the `Write`/`Edit` calls to `fixtures/brief.md`/
  `fixtures/memory.md` are not permitted at run time, which would make
  every content-checking grader (01, 03, 04, 05) fail for a reason
  unrelated to what they measure. (Originally noted as `--allow-tools
  Write` only, before Correction 3/6 added `Edit` to both `allowed_tools`
  and the expected trace.)
- **Correction 2/6 (paid run r3, grader 01 fixed):** with grader 01 fixed,
  run r3 passed 01/04/05 3/3 but failed graders 02/03 with `"after" tool
  Agent never called`. Cause: the fixture's idea named "the receipts
  service" without that service existing anywhere in the checkout, so the
  coordinator correctly persisted the seven answers, then correctly
  stopped on a genuine mode-(b) product question ("no receipts service
  exists in this checkout... (a)/(b)/(c)") instead of dispatching
  `product-discovery` — compliant behavior the case's own fixture invited,
  and one the installed CLI's `tool_order` grader cannot pass around: it
  has no allow-missing option and fails outright when the `after` tool
  never occurs. Fixed at the fixture, not the grader (graders 01-05
  unchanged): added `fixtures/receipts/README.md`, a tiny stub describing
  the checkout's existing receipts service (`send`/`bounce`/`attachment`),
  copied into the workspace by the existing unmodified `scaffold.sh`
  alongside the other fixtures; and reworded the prompt's idea to point at
  that stub and explicitly frame `discover` as producing the PRD for one
  new capability on it, closing the "does this service even exist" gap
  that was the actual open question — not a new product question of its
  own. This should let the discover stage's first dispatch happen as the
  natural next step, with no open product question left for mode (b) to
  raise.
- **Correction 3/6 (paid run r5, kept trace):** the real trace order was
  #11 `Write fixtures/brief.md` → #12 `Edit fixtures/memory.md` → #13
  `Agent product-discovery` — correct behavior, since `memory.md` is an
  *existing* record (seeded from the bundled `MEMORY.md` template) and the
  Coordinator edits an existing record rather than overwriting it. Two
  defects, fixed here:
  1. `tool_order`'s `input_match` matches against the whole serialized
     tool input, including `content`, not just the target path. Grader 03
     (`before: Write, input_match: memory\.md`) passed in r5 only because
     the brief's `Write` *content* happened to mention `memory.md` — a
     false positive (it had correctly failed in r4, when memory really
     was written with `Edit`). Fixed by anchoring both graders 02 and 03
     to the serialized `file_path` key specifically
     (`"file_path":[ \t]*"[^"]*brief\.md"` /
     `..."[^"]*memory\.md"`, matching the trace's exact
     `"file_path": "fixtures/memory.md"` spacing, `[ \t]*` rather than
     `\s` per the private-path check), and by changing grader 03's
     `before` tool from `Write` to `Edit` to match the correct, observed
     behavior. `Edit` was added to `allowed_tools` in `prompt.md`'s
     frontmatter to make this available.
  2. Run r5 cost US$1.64 against `--max-cost-usd 0.75` (the ceiling
     doesn't stop a run mid-flight) because the Coordinator ran the whole
     `discover` loop — PRD, validator, correction, re-validation — instead
     of stopping after the first stage's specialist returned; grader 01
     was skipped with reason "cost ceiling". Fixed by adding an explicit,
     naturally-worded owner instruction to the prompt telling the
     assistant to stop and report the stage status line once
     `product-discovery` returns, without continuing into document
     validation or any later stage, and by lowering `max_turns` from 15
     to 10 (the run no longer needs headroom for the validation loop).
