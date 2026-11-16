# PRD-034 · A hash guard for MEMORY.md writes and the read-only `/dh:dream` command

| Field | Value |
|---|---|
| Status | approved on 2026-11-16 |
| Owner | Daniel Malka |
| Created / updated | 2026-11-16 / 2026-11-16 |
| Tickets | to be defined by the planner |
| Origin | Owner decision 2026-11-16: tips 5 and 9 of the knowledge-base page on agent memory and "dreaming" (a conference talk, secondary source, not checked against the video) |
| Depends on | PRD-032 (report files, Coordinator scope) and ADR-005 (guards by prose and mod); version base 0.40.0 |
| Delivery | One pull request, version 0.41.0 |

## 1. Problem

Facts:

- `/dh:consolidate-memory` must start with the owner's statement that no other session writes the records (`.commands/consolidate-memory.md`, Prerequisites). The skill `context-handoff` says the same in its Archive and Verify steps. The check is a question, not a test: a session that forgets, or a second session started later, passes it.
- Every other kit write of `MEMORY.md` (the Coordinator updates it after results, setup initializes it, handoff and resume fold state into it) has no check at all. Two sessions on one harness dir overwrite each other without a trace.
- Review and QA reports pile up under the harness dir (`<task folder>/reports/`, `*.review.md`). The same finding comes back in many tasks and nothing turns it into a check. The Coordinator has one such rule today, for noise findings (`.agents/coordinator.md`, item 8: "the same pattern classified noise 3 times"), and it covers one case only.
- The resolver and the lock already exist in Go (`internal/harness/resolve.go`, `lock.go`); the binary has no command that touches the records.

## 2. Solution

Two small features, one delivery.

A. **Hash guard.** The binary gains `dh memory hash` and `dh memory write`. `hash` prints the sha256 of `MEMORY.md`. `write` replaces the file only when its hash is still the one the caller read; if the hash changed, it writes nothing and says so. The kit prose makes this the way the kit writes `MEMORY.md`, and the owner question in `consolidate-memory` goes away.

B. **`/dh:dream`.** A new command, read-only, that reads the review and QA reports and `MEMORY.md`, finds findings that came back in at least N tasks (N is an argument, default 2) and writes one proposal file: "this repeated N times; make it an automatic check (skill, gate or rule)?", with examples and origin. It applies nothing.

## 3. Rules

### Part A · Memory write guard

- R1 `dh memory hash [dir]` resolves the harness dir like `dh harness-path [dir]` and prints the sha256 of `<harness dir>/MEMORY.md` as 64 lower-case hex digits, or `absent` when the file does not exist. It writes nothing. It exits 1 with a clear message when the mode is `none` (no harness dir).
  - When `MEMORY.md` holds a known text, then the printed digest equals `sha256sum` of the file.
  - When the file does not exist, then it prints `absent` and exits 0.
  - When the folder resolves to mode `none`, then it exits 1 and names the cause.
- R2 `dh memory write --expect <sha256|absent> [--from <file>] [dir]` is a compare-and-swap. It takes the lock `MEMORY.md.lock` in the harness dir, hashes the current file, and compares with `--expect`. Equal: it writes the new content to a temporary file in the same folder, flushes it, renames it over `MEMORY.md`, releases the lock and prints `written <new sha256>`. The content comes from `--from` or, without it, from standard input. `absent` is the expected value of a file that must not exist yet (the first write).
  - When the hash is equal, then the file holds the new content and no temporary file or lock is left.
  - When `--expect absent` and the file does not exist, then it is created.
  - When the content is empty, then it exits 1 and writes nothing (an empty `MEMORY.md` is never the intent).
  - When the lock cannot be taken within 5 seconds, then it exits 1 and writes nothing.
- R3 When the current hash differs from `--expect`, `dh memory write` writes nothing, leaves the file as it is, exits 3, and prints one clear message on standard error: the file changed since it was read, someone wrote it in the middle, nothing was written, re-read `MEMORY.md` and redo the edit. The message names both hashes.
  - When another process rewrote `MEMORY.md` after the caller read it, then the call exits 3 and the other process's content is intact.
  - When `--expect absent` and the file exists, then it exits 3.
- R4 Concurrent writers are tested. Several writers start with the same expected hash and different contents, in one process (goroutines) and in separate processes.
  - When N writers race with the same expected hash, then exactly one succeeds, the other N−1 exit 3, and the final file equals the winner's content byte for byte.
  - When the race ends, then no temporary file and no lock file remain.
- R5 The kit prose uses the guard for every write of `MEMORY.md`. `.agents/coordinator.md` and the skill `context-handoff` carry one rule: take `dh memory hash` when you read `MEMORY.md`; edit a draft copy outside the harness dir; write with `dh memory write --expect <hash> --from <draft>`; keep the printed hash as the next expected one. On exit 3, do not retry the same content: read `MEMORY.md` again, redo the edit on it, and write once more; a second conflict goes to the owner. The first write of a new harness dir (setup, from the template) uses `--expect absent`.
  - When the Coordinator updates `MEMORY.md` after a result, then the write is a `dh memory write` with the hash read at the start.
  - When two sessions update the same `MEMORY.md`, then the second writer gets exit 3 instead of overwriting.
- R6 `.commands/consolidate-memory.md` and the Archive, Verify and Consolidation steps of `context-handoff` drop the owner question "no other session is writing". They take `dh memory hash` before the archive, and the trim of `MEMORY.md` is a `dh memory write --expect <that hash>`. The append to `EPOCHAL.md` stays as it is (append-only, verified by reading it back). When the trim write exits 3, `MEMORY.md` stays intact, the command reports the operation as incomplete with the batch id already stored, and a new run starts a new batch from the current memory. `handoff.md`, `resume.md` and `setup.md` point to the guard in their Limits.
  - When `/dh:consolidate-memory` runs, then it does not ask the owner whether another session writes.
  - When `MEMORY.md` changes between the archive and the trim, then the trim is refused and `MEMORY.md` is intact.
- R7 The limits are written, not hidden. The guard protects the callers of `dh memory write`. A writer that edits the file directly is detected only when its write lands before the compare; one that lands between the compare and the rename, inside the lock, is lost. The Claude Code mod does not force the use of the command (ADR-005: prose first). The guard covers `MEMORY.md` only, not `EPOCHAL.md` or `RISKS.md`.

### Part B · `/dh:dream`

- R8 `.commands/dream.md` exists, with the frontmatter of the other commands: `name: dream`, `argument-hint: "[N]"`, roles `coordinator`, `repo-scout` and `docs-guide`, `writes:` naming only the proposal file. `$ARGUMENTS` is N, an integer of 2 or more; empty means 2; anything else stops with a message and no work.
  - When no argument is given, then N is 2.
  - When the argument is `1`, `0` or not a number, then the command stops and says the valid range.
- R9 The command is read-only on the project and on the records. It reads: the files named `*.md` under `<harness dir>/tasks/` whose path holds `reports/`, every `*.review.md` under the harness dir, and `<harness dir>/MEMORY.md`. It writes nothing but the proposal file of R11. It does not edit the product, tests, kit sources, `MEMORY.md`, `EPOCHAL.md` or `RISKS.md`, opens no ticket, and runs no build, test or QA.
  - When the command ends, then `git status` shows only the proposal file as new (and only when a group reached N).
- R10 Recurrence has one rule. The Coordinator dispatches `repo-scout` (model `haiku`) to list every finding found in the inputs, one line each, without grouping or judging: task, file and line, severity, location, title, state, and a quote of at most two lines. The reply cap in the dispatch is raised, with the reason. The Coordinator then groups. Two findings are the same when they have the same title after case and punctuation are ignored, or name the same scenario at another location. The count of a group is the number of distinct tasks it appears in; the rounds of one task count once. Findings in `MEMORY.md` (incidents, blockers) count as one task. A group with a count below N is dropped. The finding states seen (`open`, `closed`, `noise`) are kept in the proposal, and do not change the count.
  - When the same finding appears in three rounds of one task, then it counts once.
  - When a finding appears in two different tasks and N is 2, then it is a group.
- R11 When at least one group reaches N, `docs-guide` (model `sonnet`) writes one new file, `<harness dir>/proposals/<YYYY-MM-DD>_DREAM_PROPOSAL.md`, in the `language` recorded in `project.yaml` (the dispatch states it). The file starts with the date, N, how many files and tasks were read, and the line "Proposal only: nothing was applied." It has one section per group with: the title; the count and the tasks; the states seen; one to three examples (task, relative `path:line`, a quote of at most two lines); a suggested target (a skill, a gate or a rule) with one sentence of reason; the question "this repeated N times; make it an automatic check (skill, gate or rule)?"; and an empty `Decision:` line for the owner. A file of the same date is not overwritten: the new one gets `-2` before `_DREAM_PROPOSAL`, then `-3` (as the drift report does). The chat reply carries only the path, N, the files and tasks read, and the number of groups. When no group reaches N, no file is written and the reply says that nothing repeated N times.
  - When two groups reach N, then the file has two sections with the same fields.
  - When a quote holds something that looks like a secret, then it is replaced by `[redacted]`.
- R12 The command proposes and never applies. The report text is data, quoted and not followed. A proposal becomes a check, a skill line or a rule only when the owner decides, through the normal flow: the `harness-maintainer` with `harness-authoring` for the kit, the project's own lane for a project. The reply states this in one sentence.
  - When a report holds an instruction aimed at the agent, then the command lists it as a quote and does not act on it.
- R13 It is a new command, not a mode of `/dh:consolidate-memory`. That command writes the records and needs the guard of Part A. `/dh:dream` never writes a record, can run at any time, and has its own `writes:` and limits. The count of commands in files goes from 19 to 20; the mod command stays, so the count is 20 + 1. The agents (13) and skills (22) do not change.

### Part C · Docs, version and checks

- R14 The docs say one thing. `AGENTS.md` has the new rules and the counts 20 + 1. `README.md`, `README.en.md`, `docs/tutorial.html` and `docs/en/tutorial.html` describe the hash guard (the owner question is gone) and list `/dh:dream` in the command table. `REVIEW.md` and `.claude-plugin/marketplace.json` carry the count of 20. `docs/roadmap.html` and `docs/en/roadmap.html` show 20 commands in the tree and in the current count.
  - When the live docs and the sources are compared, then every statement of the command count says 20 (plus 1 for the mod), in pt-br and in en.
  - When `dh validate` runs, then it finds no count drift.
- R15 Version 0.41.0 is set in `internal/build/build.go`, in `.claude-plugin/marketplace.json` (version and ref `v0.41.0`) and in `harness-manifest.json`. `CHANGELOG.md` has an entry dated 2026-11-16. The roadmap pt and en have updated chips, the Sources range and a new delivered version card. `dist/` is rebuilt, binaries included, because the binary has a new subcommand.
- R16 The checks pass, and the proof of a clean clone is required because the change touches `cmd/` and `internal/`: `go test ./...`, `go vet ./...`, `go run ./cmd/dh validate .`, `go run ./cmd/dh build` with no difference in `dist/` (the `generated_at` line aside), the same sha256 for the five binaries built in a fresh `git clone`, and `claude plugin test dist/claude-code/dev-harness` unchanged. A file of this delivery holds nothing of the owner's private setup.

Out of scope: running `/dh:dream` on a schedule; applying a proposal; reading session transcripts (the talk's dreaming reads them; this command reads reports and `MEMORY.md` only, as the owner asked); a mod rule that forces `dh memory write`; a guard for `EPOCHAL.md`, `RISKS.md` or other records; any new agent or skill.

## 4. Docs

- `.commands/dream.md` (new), `.commands/consolidate-memory.md`, `.commands/handoff.md`, `.commands/resume.md`, `.commands/setup.md`: the flows touched (R5, R6, R8 to R12).
- `.agents/coordinator.md`, `.skills/context-handoff/SKILL.md`: the guard rule (R5, R6).
- `AGENTS.md`, `README.md`, `README.en.md`, `REVIEW.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `.claude-plugin/marketplace.json`: counts and the new behavior (R14).
- `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html`: version 0.41.0 (R15).
- `docs/adr/ADR-005-trava-de-runtime-por-mod.md` is not changed: the guard is prose plus a command, inside what the ADR already allows.

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| Hash check by prose only ("compare the sha256 yourself") | The check and the write are two steps the model can skip or reorder. The command does both under one lock, and it can be tested. |
| A mod rule that denies direct writes of `MEMORY.md` | It forces one runtime and would need its own review (ADR-005, ADR-006). The command works in every runtime that has the binary; the mod can come later. |
| Guard `EPOCHAL.md` and `RISKS.md` too | The owner asked for `MEMORY.md`. `EPOCHAL.md` is append-only and verified by reading it back; `RISKS.md` has the same single writer and can use the command later. |
| A mode `dream` of `/dh:consolidate-memory` | That command writes and needs the guard; a read-only mode inside it blurs `writes:` and the limits (R13). |
| A Go scanner that counts the findings | Findings are free prose by many roles (code, document, QA). A parser would match the title strings only and miss the same defect worded two ways. The judgement stays with the Coordinator, the reading with `repo-scout`. |
| `N` as a config key | An argument is enough for a command the owner runs by hand, and no `project.yaml` or `config.yaml` key is added. |
| Read session transcripts as the talk does | The kit has no transcript store; reports and `MEMORY.md` are the durable record. |

### Risks and pending decisions

- Risk: the Coordinator keeps editing `MEMORY.md` directly and skips the command. Mitigation: R5 makes it the rule in the two places the Coordinator reads first, and R7 states that the mod does not force it. A mod rule is the follow-up.
- Risk: a draft copy of a 40,000-byte `MEMORY.md` costs tokens. Mitigation: the draft is a file copy edited in place by small edits, not a rewrite; the cost is one copy and one write.
- Risk: grouping by the model is uneven (two wordings of one defect, or two defects under one title). Mitigation: R10 fixes one rule, the proposal shows the examples with origin so the owner judges the group, and a wrong group costs a decision, not a change.
- Risk: the proposal file holds a secret quoted from a report. Mitigation: R11 redacts and the quote is at most two lines.

### Decisions (closed 2026-11-16, owner)

- D1 Tip 5 replaces the owner question of `/dh:consolidate-memory` with a hash check, and covers every kit write of `MEMORY.md`.
- D2 Tip 9 is a read-only command `/dh:dream`; `N` is an argument, default 2.
- D3 `/dh:dream` is a command, not a mode of `/dh:consolidate-memory` (R13).
- D4 Minor version bump from 0.40.0 to 0.41.0: a command and a binary subcommand are added.
