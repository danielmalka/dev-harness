# PRD-029 · Independent review rules and drift report

| Field | Value |
|---|---|
| Status | delivered on 2026-10-28 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-28 / 2026-10-28 |
| Tickets | to be defined by the planner |
| Origin | Owner request 2026-10-28; discovery UND-009 (research in `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Depends on | PRD-007 (ticket flow and final-feature review, delivered), PRD-026 (review guide, delivered), PRD-028 (version base 0.33.0, delivered) |
| Delivery | One pull request, version 0.34.0 |

## 1. Problem

A review is only as independent as the context it gets. Five gaps exist today. All five sit in the review and verification flow.

**A · The context bundle is not required to be free of the Coordinator's conversation.** Fact: only the in-flight doubt withholds the author's claim (`.agents/coordinator.md`, "In-flight doubt", line 143). Fact: the dispatch contract in `.agents/coordinator.md` (Procedure, item 6, line 66) lists "source artifacts" and does not forbid chat text. Fact: the `document-validator` source is "discovery notes, the owner's request, or the task record" (`.skills/document-review/SKILL.md`, Inputs, line 34). Nothing stops the Coordinator from pasting its own summary of the chat. A reviewer that reads that summary inherits the author's reading.

**B · No check that the validator understood the document.** Fact: the validator output (`.skills/document-review/SKILL.md` lines 80 to 101, `.agents/document-validator.md` lines 98 to 116) has verdict, findings and "Not raised". It has no statement of what the validator understood. A validator can approve a text it read in a wrong way.

**C · Findings of a code review are not triaged by evidence.** Fact: Critical and Major findings go to the builder. Minor findings stay in the report (`.agents/coordinator.md` item 8, line 68). Fact: the Coordinator records `fixed`, `intentional` or `noise` for each finding (same line). Fact: `.skills/code-review/SKILL.md` ("For the author", step 3, line 89) asks the author to push back with a file, a test or a constraint. Five things are missing. A code finding has no stable id (only document findings have `P1`, `P2`). No rule locks the severity. No rule keeps a settled finding settled. No rule covers a fixed finding that comes back. Nothing detects a loop that fixes and breaks the same place.

**D · The loops have no convergence rule beyond "one clean round".** Fact: "a round ends the loop when it has no blocking finding" (`.agents/coordinator.md`, line 119; `AGENTS.md`, bullet "Verificação padrão", line 20). Fact: a stage with every reviewer not-run is never approved (`.agents/coordinator.md`, line 111). Missing: a rule that a round with one not-run reviewer is not clean. Fresh reviewers sample differently each round, so one clean round is a weak signal.

**E · No check of a design document against the code.** Fact: the kit checks only mechanical drift (`internal/kit/status_drift.go`, `internal/kit/roadmap_check.go`, `internal/kit/doc_counts.go`). Nothing compares a PRD or an ADR with the code that implements it. A rule such as "a sold product is never deleted" can break months later with no signal.

The source of the ideas is the Elephant-Goldfish Model (EGM) of Google: a plugin of markdown skills and reviewer personas, Apache-2.0. EGM is not a machine-learning model. Fact: EGM publishes no result that shows better outcomes (research, section 3). This PRD takes the mechanisms and reimplements them in the words of this kit.

## 2. Solution

The batch has five parts and one credit. The first four change text in existing assets. The fifth adds one skill. No agent and no command is added.

**Part A · Independence of the context bundle.** A reviewer prompt has two parts. The context bundle is what the reviewer reads about the work. The prompt assembly is the role body, the skill, the anti-delegation clause and the verdict contract. This batch restricts the context bundle only. For a document reviewer and for any CLI reviewer, the context bundle holds files and quoted data. It holds no chat, no summary of the owner's intent, and no reasoning of the Coordinator or of the author. The in-flight doubt keeps its own rule.

**Part B · Comprehension block.** The `document-validator` writes a section `## Comprehension` in every report. It states in its own words what the document proposes and how the system works, as the document says. A section that the validator cannot paraphrase in one reading is a blocking ambiguity. The Coordinator compares the paraphrase with the owner's request before the owner sees the document.

**Part C · Triage of code findings.** The Coordinator gives each Critical and Major finding an id. Nobody changes the severity that the reviewer wrote. The builder answers each finding with a disposition and evidence. A settled finding is not raised again without new evidence. A fixed finding that comes back reopens under its old id. A finding that exists only because of an earlier fix is churn, and churn stops the loop.

**Part D · Convergence.** A round is clean only when no finding is `open`, `reopened` or `escalated` and every dispatched reviewer returned a verdict. The build/QA and code-review loop ends after two consecutive clean rounds, at the cap of six rounds, or at churn. The owner's wording "rule of two minimum layers" means this. The owner confirmed the reading and its scope on 2026-10-28 (decision D1).

**Part E · Drift report.** A new skill `drift-report` compares one design document (PRD or ADR) with the code it covers. A read-only `repo-scout` dispatch finds the code. A fresh reviewer compares the document section by section. The report lists only differences, in four classes, plus code logic that no section mentions. Each difference cites the document and the code. The report is a file in the task folder. It never fixes anything. The Coordinator runs it on request. After `FEATURE APROVADA` in `/dh:review feature PRD-<n>`, the Coordinator offers it in one optional line.

**Credit.** The kit credits EGM in the provenance note of each skill that takes a mechanism from it, and in one line of each README (decision D3). No EGM text is copied.

## 3. Rules

All rules are final. The owner decided D1, D2 and D3 on 2026-10-28. The decisions are recorded in the Appendix.

### 3.1 Part A · Independence of the context bundle

- R1 A reviewer prompt has two parts. The context bundle is what the reviewer reads about the work. The prompt assembly is the role body, the skill text, the anti-delegation clause, the verdict contract and the reading-is-allowed block. Rules R1 to R9 restrict the context bundle only. The prompt assembly stays as `.agents/coordinator.md` and `.skills/external-clis/SKILL.md` define it. This batch does not change it.
- R2 The Coordinator's conversation is any owner chat turn, any reasoning of the Coordinator, any summary of intent, and any note on earlier rounds. The Coordinator never puts these in the context bundle for a reviewer of the document stage.
- R3 When the Coordinator dispatches a document-stage reviewer, then the context bundle holds only five items. The items are the document path, the source material, the previous review report path on a later round, the project language and the harness dir. The source material is a list of file paths or the owner's request quoted word for word.
- R4 When a CLI reviewer of any stage (`document`, `code`, `security`) is called, then its context bundle holds no part of the Coordinator's conversation.
- R5 The owner's request may enter a context bundle only in one of two forms. The first form is a verbatim quote labelled "owner's words". The second form is the path of a file that holds the request. The Coordinator adds no paraphrase.
- R6 Builder and QA check results may stay in a code-stage context bundle as quoted data labelled unverified. `.commands/review.md` (Routing, line 15) already allows this. R2 does not forbid it.
- R7 The Coordinator lists the elements of each context bundle in the Dispatches entry of `MEMORY.md`. Each element is a path or a label. A later reader can check R3 and R4 from this list.
- R8 The rule of R1 to R7 appears with the same meaning in four files.
  - `.agents/coordinator.md` (Procedure, item 6).
  - `.skills/external-clis/SKILL.md` (Inputs row for the context bundle, and step 3).
  - `.skills/document-review/SKILL.md` (Inputs and round protocol).
  - `.agents/document-validator.md` (Minimum inputs).
- R9 The in-flight doubt rule stays as written. The Claude path gets no stop rule for a host without fresh subagents.

Check for R1 to R9: a read of the four files named in R8, and a read of one `MEMORY.md` Dispatches entry from the first use.

### 3.2 Part B · Comprehension block

- R10 When `document-validator` reviews a PRD, a plan, a story or an ADR, then its reply and its persisted report each have a section `## Comprehension`.
- R11 The section has three items. The first is "What the document proposes", in at most 5 sentences. The second is "How the current system works, as the document states it", in at most 5 sentences. The third is "Sections not paraphrased": a list of section ids, or `none`.
- R12 The validator writes the two paraphrases in its own words. A copied phrase is in quotation marks.
- R13 When the validator cannot write one reading of a section, then it files a blocking finding of category `ambiguity` for that section. It lists the section under "Sections not paraphrased". This rule is an exception. The severity rule says an ambiguity blocks only when two readings give different software (`.skills/document-review/SKILL.md`, step 6, line 59). The skill states the exception.
- R14 Take a term that is specific to the project. Suppose no document, no source and no file of the repository defines it. Then the validator files an ambiguity finding under the normal severity rule. A general engineering term is never a finding.
- R15 When the verdict is `approved` with zero findings, then the report still holds `## Comprehension`.
- R16 The `## Verdict` section stays the first section of the reply, with its current lines. The persisted report places `## Comprehension` after the header table and before `## Findings`. In the reply, it follows `## Verdict`.
- R17 Before the Coordinator sends a validated document to the owner, it reads `## Comprehension`. Its message to the owner says that the paraphrase matches the request, or it quotes the sentence that differs. A difference goes to the owner as an open point. The Coordinator does not settle it alone.
- R18 The block applies to the document stage only. The `reviewer` modes `code` and `security` get no such block. The block adds no round and does not change the cap of two correction rounds.

Check for R10 to R18: a read of `.skills/document-review/SKILL.md` (Output format, both blocks, and step 6), `.agents/document-validator.md` (Output format) and `.agents/coordinator.md` (document handoff).

### 3.3 Part C · Triage of code findings

- R19 When the Coordinator merges the code-stage verdicts of a round, then it gives an id of the form `F<round>-<n>` to each new Critical or Major finding. `<round>` is the round in which the finding first appeared. `<n>` counts from 1 among the new findings of that round, in the order of the merged report. A finding carried into a later round keeps its original id. A finding never gets a second id.
- R20 The Coordinator records each id in `MEMORY.md`. The record has the severity as the reviewer wrote it, the location, one line of text, the state and the evidence for the state. The Coordinator updates the record of an id when its state changes. It never creates a second record for the same finding.
- R21 The states are `open`, `fixed`, `intentional`, `noise`, `escalated` and `reopened`. The first four reuse the outcomes of `.agents/coordinator.md` item 8, plus `open`. The batch adds no other vocabulary for the end-of-slice loop.
  - `open`: a reviewer raised the finding and nothing has closed it. A finding that the builder reports as fixed stays `open` until a reviewer confirms the fix.
  - `fixed`: a reviewer confirmed in a later round that the fix works. This follows the status check of `.skills/code-review/SKILL.md` ("Re-entry"). The finding is closed.
  - `intentional` and `noise`: the Coordinator accepted the evidence of R24, or a reviewer withdrew the finding (`noise`). The finding is closed and settled.
  - `escalated`: the finding waits for the owner (R25 and R29). The owner's answer moves it to `intentional` or `open`.
  - `reopened`: a `fixed` finding came back (R27). It moves to `fixed` when a reviewer confirms the fix again.
- R22 The Coordinator and the builder never change the severity of a finding. A reviewer may change it in a later round. The reviewer cites new evidence in the finding.
- R23 For each Critical or Major finding that goes to the builder, the builder reports one disposition: `fixed`, `intentional`, `noise` or `escalated`. A disposition `fixed` is a claim. The state stays `open` until a reviewer confirms it (R21).
- R24 A disposition `intentional` or `noise` names a file, a test or a constraint as evidence. Without that evidence, the finding stays `open`. `.agents/coordinator.md` line 75 stays: a blocker is not discarded because the builder disagrees.
- R25 A disposition `escalated` means that the builder cannot judge the finding with the information in the tree. The Coordinator sends it to the owner as numbered prose. The state stays `escalated` until the owner answers. After the answer, the state becomes `intentional` when the owner accepts the finding as it is. It becomes `open` when the owner asks for a fix.
- R26 A finding is settled when its state is `intentional` or `noise`. In a re-entry dispatch, the Coordinator lists each settled finding with its id, its location and the path of its evidence. The list holds no text of the builder's reasoning.
- R27 When a reviewer re-enters, then it does not raise a settled finding again at the same location and scenario. It may raise the finding with new evidence, which it cites. When a `fixed` finding appears again at the same scope, then the Coordinator sets its state to `reopened` and keeps its id. A `reopened` finding is blocking until a reviewer confirms the fix again. Then its state becomes `fixed`.

Churn rules follow.

- R28 When a reviewer re-enters, then it adds a line `Churn:` to each new Critical or Major finding. The value is `none`, or the id or location of an earlier fix, with the words "reverting that fix would resolve this finding". The reviewer reads the delta since the last reviewed state to write this line.
- R29 When one finding of a round has a `Churn:` value other than `none`, then the Coordinator stops the loop and spends no further round. It sets the state of that finding to `escalated`. It asks the owner one numbered question. The options are: accept the finding, order a fix, or replan.
  - When the owner accepts the finding, its state becomes `intentional`. The finding is settled, and the loop continues.
  - When the owner orders a fix, its state becomes `open`. The fix may revert the earlier fix. The loop continues, and the cap of the stage still applies.
  - When the owner chooses to replan, the state stays `open`. The Coordinator reports partial, as R35 says.

Check for R19 to R29: a read of `.agents/coordinator.md` (items 7 and 8), `.skills/code-review/SKILL.md` (Re-entry, Output format), `.agents/reviewer.md` (Minimum inputs, mode `code`) and `.commands/review.md` (Routing). Document findings keep their ids `P<n>` and their states `applied`, `rejected`, `pending`. The in-flight doubt keeps its four classes and its cap of 3 cycles.

### 3.4 Part D · Convergence

- R30 A code-stage round is clean when two conditions hold. First, after the Coordinator records the states of the round, no Critical or Major finding of the loop is `open`, `reopened` or `escalated`. Second, every reviewer that the stage dispatched returned a verdict. A not-run reviewer makes the round not clean. For the default stage with one reviewer, a not-run reviewer already makes the stage not-run (`.agents/coordinator.md`, line 111).
- R31 When two consecutive code-stage rounds are clean, then the build/QA and code-review loop ends as converged. This holds for the loop of every ticket and for the final feature review.
- R32 The confirming round is the second clean round. It reviews the same evaluated state as the first clean round. The evaluated state is the project tree plus the dispatch inputs. The records of the Coordinator in the harness dir are not part of it. These records are `MEMORY.md`, the task record and the finding states. Recording the states of the first clean round is not an edit. No edit of the evaluated state may happen between the two rounds. When an edit happens, the count of clean rounds restarts at zero.
- R33 A confirming round reads the whole change, as a first review does. It does not read only a delta. A delta would be empty (`.skills/code-review/SKILL.md`, step 1, line 46). The confirming round still counts as a re-entry for the Code-stage rotation. It uses the next entry of `reviewers.code`. With no other entry, the same `reviewer` runs again in a new context.
- R34 The confirming round counts in the cap of six rounds.
- R35 When the cap of six rounds ends with fewer than two consecutive clean rounds, then the Coordinator reports partial or blocked, as `.agents/coordinator.md` already says. One clean round that no second clean round confirmed is never reported as verified.

Rules R30 to R35 apply to the build/QA and code-review loop only. The owner decided that four other loops keep their current rule. The document stage keeps its cap of 2, and one clean validation ends it. The in-flight doubt keeps its cap of 3. The `security` stage and `/dh:plan-loop` do not change.

Check for R30 to R35: a read of `.agents/coordinator.md` (Proportional review, item 7, Code-stage rotation), `.skills/code-review/SKILL.md` (Re-entry) and `AGENTS.md` (bullet "Verificação padrão"). The `AGENTS.md` bullet replaces the words "rodada sem achado bloqueante encerra o loop" for the build/QA and code-review loop. The same sentence in `.agents/coordinator.md` line 119 changes in the same way.

### 3.5 Part E · The skill `drift-report`

- R36 The batch adds one skill at `.skills/drift-report/SKILL.md`. The `name` equals the folder name `drift-report`. The `description` starts with "Use when", states the triggers and states the cases where the skill does not apply. The frontmatter has `author: malka` and `metadata.provenance: adapted`. The sections follow the skill contract of `.skills/harness-authoring/SKILL.md` (table at line 110).
- R37 The inputs are one design document (a PRD or an ADR), named by the owner or by the Coordinator, and the code area that it covers. When the document is missing, the dispatch stops and returns the missing input to the Coordinator. One dispatch compares one document.
- R38 When the code area is not named, then the Coordinator dispatches `repo-scout` first, read-only, to find the code that the document covers. The Coordinator passes only the list of paths to the comparison. It does not pass the narrative of `repo-scout`.
- R39 The comparison runs in a `reviewer` dispatch in mode `code`, model `sonnet`, with the skill `drift-report` named in the dispatch. The context bundle holds the document path, the list of code paths and the harness dir (R1 to R4 apply). `.agents/reviewer.md` gets one sentence: when the dispatch names `drift-report`, the reviewer follows the output of that skill in place of the mode `code` output. No agent is added.
- R40 The reviewer splits the document by section heading. For each section, it lists the claims (rules, flows, data, decisions), finds the code that carries each claim, and compares them. It reports only differences.
- R41 A difference belongs to one class.
  - `missing feature`: the document states it and the code area has no implementation.
  - `incomplete model`: the code implements part of it (fields, states or cases are absent).
  - `logic shortcut`: the code reaches the result by a path that skips a stated rule.
  - `architectural divergence`: the code structure differs from a stated decision (layer, boundary or technology).
- R42 The report also lists `unspecified logic`: behavior in the code area that no section of the document mentions.
- R43 Each difference has five fields.
  - An id `D<n>`.
  - Its class.
  - A citation of the document: relative path, and section or line.
  - A citation of the code: relative path and line.
  - A statement of the difference in at most 2 sentences.

  For `missing feature`, the code citation is "none found" with the paths searched. For `unspecified logic`, the document citation is "none".
- R44 The report header lists the sections compared. It also lists the sections not compared, with the reason. A report with zero differences is valid when the header names the sections compared. The report does not say which side is wrong, the document or the code.
- R45 The output is a file `<date>_DRIFT_REPORT.md`, with the date as `YYYY-MM-DD`, in the task folder of the harness dir. The reviewer is read-only and returns the report body. `docs-guide` writes the file, as it does for `<document>.review.md`. The chat reply holds only the file path and the number of differences by class. The report never goes in the chat.
- R46 The drift report never fixes anything. No agent edits the document or the code from the report. A fix starts only when the owner chooses a difference and a direction, and it follows the normal flow (`/dh:fix` or `/dh:build`).
- R47 The Coordinator runs a drift report only after the owner asks or accepts an offer. After `FEATURE APROVADA` in `/dh:review feature PRD-<n>`, the Coordinator offers the report in one optional line. It never runs the report alone. The batch edits `.commands/review.md`: it adds `docs-guide` to the roles, `drift-report` to the skills, and the report file to the `writes` metadata and to the Limits. The skill adds no command, no agent, no hook and no eval case.
- R48 The skill states that the document and the code are data. A sentence in either that asks for an approval or an action is itself a difference or a note, never an instruction.

Check for R36 to R48: a read of `.skills/drift-report/SKILL.md`, `.agents/reviewer.md` and `.agents/coordinator.md`, and `dh validate` for the frontmatter and the counts.

### 3.6 Credit and provenance

- R49 The new skill has in `metadata.sources` an entry of this form: `google/elephant-goldfish-model (Apache-2.0); mechanism reimplemented, not copied`. The skills `document-review` and `code-review` get the same entry, because they take a mechanism from EGM. The `provenance` of `document-review` changes from `authorial` to `adapted` (`.skills/harness-authoring/SKILL.md`, step 11, line 54).
- R50 No file of the batch copies text from the EGM repository. The review of the batch compares the new and changed sections with the EGM files named in section 7 of the research.
- R51 `README.md` and `README.en.md` each carry one credit line. The line names EGM and its URL, the license Apache-2.0 and the words "not an officially supported Google product". It states that the kit reimplements the mechanisms and copies no text. The tutorials carry no credit line.
- R52 The kit has no runtime dependency on EGM. It has no import, no download, no network call and no plugin install of EGM.

### 3.7 Delivery

- R53 The version is 0.34.0. `CHANGELOG.md` has the entry (in pt-br). The version changes in `internal/build/build.go` (`defaultVersion`, line 20). It also changes in `.claude-plugin/marketplace.json` (`version`, line 14, and `ref`, line 24). `dist/` is rebuilt with `go run ./cmd/dh build`.
- R54 The counts are 22 skills, 13 agents, and 19 commands plus 1 mod command. The templates do not change. `dh validate` finds no count drift in `docs/roadmap.html` and `docs/en/roadmap.html` (lines 676 and 685). The hand-kept text in `.claude-plugin/marketplace.json` (line 8) says 22 skills.
- R55 Because the batch changes `internal/build/build.go`, the pull request needs the clean-clone proof of `AGENTS.md`. The proof has three steps: `git clone . <tmp>`, then `go run ./cmd/dh build` there, then a comparison of the sha256 of the binaries with the tree. The Coordinator commits only after the proof.
- R56 The batch edits no file named `references/clone-isolation.md` or `references/live-tree-freeze.md`. It edits no line of the anti-delegation clause. It changes no part of the PRD-027 guard. The batch needs no `security` stage. It needs the default code review and a `qa-verifier` matrix.
- R57 `AGENTS.md` receives only the current rules, in pt-br. These are the context-bundle rule (R1 to R5) and the replaced sentence of the convergence rule (R30 to R35). One bullet covers the drift report (R46, R47). History stays in this PRD.
- R58 The change updates the documents of section 4 in the same pull request. No documentation-only pull request exists.
- R59 No eval case is added and no paid eval runs. At first use, the Coordinator records in `MEMORY.md` the cost of one confirming round and of one drift report. The record is a note, not a measure.
- R60 Every new or changed file uses relative paths or `~/...`. `dh validate` finds no absolute private path.

## 4. Docs

- README.md
- README.en.md
- docs/tutorial.html
- docs/en/tutorial.html
- docs/roadmap.html
- docs/en/roadmap.html
- AGENTS.md
- CHANGELOG.md
- .claude-plugin/marketplace.json
- .agents/coordinator.md
- .agents/document-validator.md
- .agents/reviewer.md
- .skills/document-review/SKILL.md
- .skills/code-review/SKILL.md
- .skills/external-clis/SKILL.md
- .skills/drift-report/SKILL.md
- .commands/review.md

Notes on this list.

- The roadmaps get the card for 0.34.0, the chips (date and status `0.34.0`), and the count 22 in the structure block and in the "Contagem atual" line.
- The tutorials describe the context-bundle rule, the `## Comprehension` block, the finding ids and states, the two clean rounds, and how to ask for a drift report.
- The README files get the credit line of R51. Fact: neither README states a skill count (checked by `grep`).
- `.commands/review.md` changes for the offer of R47. `.agents/reviewer.md` changes for the sentence of R39.
- The provenance notes of R49 sit in `.skills/drift-report/SKILL.md`, `.skills/document-review/SKILL.md` and `.skills/code-review/SKILL.md`.
- `internal/build/build.go` changes for the version (R53). It is source code, not a document.

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| Install EGM beside the kit | Two sets of commands and a runtime dependency. The owner chose to reimplement and credit. |
| Copy EGM text with an Apache-2.0 notice | The owner chose to reimplement in the kit's own words (R50). |
| A new command for the drift report | The owner chose no new command. The count stays 19 plus 1. |
| A new agent for the drift comparison | The owner chose no new agent. The count stays 13. |
| A stop rule when the host has no fresh subagent | The owner chose the context-bundle rule only (R9). |
| The comprehension block in every review mode | The owner chose the document validator only (R18). |

### Assumptions the owner may overturn

- A1 One PRD, one pull request, version 0.34.0. The Coordinator proposed it, and the owner did not contest it.
- A2 No paid eval and no new eval case. The cost is recorded informally at first use (R59).
- A3 The names are in English: the skill `drift-report`, the section `## Comprehension`, the ids `F<round>-<n>` and `D<n>`, and the file `<date>_DRIFT_REPORT.md`.
- A4 Only Critical and Major code findings get an id (R19). Minor findings stay in the report and get no id.
- A5 The placement of `## Comprehension` is after the header table in the persisted report and after `## Verdict` in the reply (R16).
- A6 The states of a code finding reuse the outcomes of `.agents/coordinator.md` item 8. They add `open`, `escalated` and `reopened` (R21). A builder claim of `fixed` stays `open` until a reviewer confirms it.
- A7 "The Coordinator's conversation" means chat turns, Coordinator reasoning, summaries of intent, and notes on earlier rounds (R2). Builder and QA check results stay allowed as unverified data (R6).
- A8 A settled-finding list holds the id, the location and the evidence path only (R26). The builder's reasoning stays out, so it does not anchor the reviewer.
- A9 The reviewer detects churn from the delta (R28). One churn finding stops the loop (R29).
- A10 The report goes in `tasks/PRD-<n>/` of the harness dir when the document is a PRD. For an ADR it goes in the current task folder. When the file name exists, `docs-guide` adds `-2` before `_DRIFT_REPORT`.
- A11 A large code area is split by the Coordinator, one dispatch for each group of sections (R37).
- A12 No new ADR is written. `AGENTS.md` says history lives in `docs/adr/`. This batch changes one rule of the owner (the convergence rule), so the owner may want an ADR.
- A13 The `reviewer` agent file gets one sentence for `drift-report` (R39). The owner confirmed the role in decision D2(b).
- A14 No `security` stage is needed (R56). The batch does not touch isolation, the clause or the guard.
- A15 The line `docs/plano-produto.html:828` ("21 skills") is a product plan outside the batch. Whether `internal/kit/doc_counts.go` reads it is not verified.

### Out of scope

- The brainstorm of EGM with six lenses, and its section-by-section draft.
- The folder READMEs of EGM ("bootstrap", "peanuts" and "hay").
- Any new agent and any new command.
- Any paid eval, any new eval case and any new probe.
- Copying EGM text, and any runtime dependency on EGM.
- A stop rule for hosts without fresh subagents.
- The hooks and the plugin files of EGM for Claude Code.
- A change to the isolation of CLI reviewers, to the PRD-027 guard or to the anti-delegation clause.

### Risks and pending decisions

- Risk: the token cost of a confirming round and of a drift report is not measured · mitigation: R59 records the first use.
- Risk: the confirming round adds at least one review to every loop that passes at the first round · mitigation: the owner accepted this cost in decision D1(b) and chose every loop. R59 records the cost, and the owner may narrow the scope later.
- Risk: EGM publishes no evidence that these mechanisms improve outcomes (research, section 3) · mitigation: R59 and the first uses give the owner a local record.
- Risk: the reviewer and the Coordinator use the same base model, so a fresh context removes only the context bias · mitigation: R33 rotates to the next entry of `reviewers.code` when the list has one.
- Risk: a drift report on a partial read gives false differences · mitigation: R43 requires two citations, and R44 lists the sections not compared.
- Risk: Apache-2.0 attribution applies if text is copied · mitigation: R50 forbids copying, and R49 and R51 credit the source.
- Pending decision: none. The owner decided D1, D2 and D3 on 2026-10-28. The record follows.

#### D1 · Meaning and scope of "two minimum layers" (decided)

Context: the owner wrote "Prefiro adicionar a regra de 2 camadas mínimas". The Coordinator read it as the EGM rule "two consecutive clean rounds". Before this batch, one round with no blocking finding ended a loop (`AGENTS.md` line 20, `.agents/coordinator.md` line 119). The new rule costs at least one more review for each loop that passes at the first round.

Part (a) · The meaning.

1. (Chosen) Two consecutive clean rounds, with a confirming round on the same evaluated state (R30 to R35). Consequence: a clear count, at least two reviews for each loop. With `reviewers.code` that lists two entries, the rotation makes the second review come from another reviewer.
2. Two different reviewers must approve, in any order. Consequence: it needs two entries in `reviewers.code`. With one entry, the rule cannot hold or falls back to a repeat of the same reviewer.
3. Two layers means the deterministic checks plus one independent review. Consequence: this is already true today, so the rule would only restate it. No extra cost, no new safety.

Decided by the owner on 2026-10-28: option 1.

Part (b) · The loops that the rule covers.

- Build/QA and code-review loop (cap 6).
  1. (Chosen) Every loop, each ticket and the final feature review. Consequence: the highest cost, the strongest signal.
  2. The final feature review only. Consequence: the ticket loops keep one clean round, and the cost stays low.
  3. Only a loop that already found a blocking finding. Consequence: a change that passes at the first round pays nothing, but a lucky first round is not confirmed.
  Decided by the owner on 2026-10-28: option 1.
- Document validation (cap 2). Two clean rounds cannot fit in two correction rounds.
  1. Keep the cap of 2 and require the last validation to be clean. Consequence: in effect the same as today, with a new sentence.
  2. Raise the cap to 3 and require two consecutive clean rounds. Consequence: one more Opus validation for most documents.
  3. (Chosen) Leave the documents on the old rule. Consequence: no change and no extra cost.
  Decided by the owner on 2026-10-28: option 3.
- In-flight doubt (cap 3). Decided by the owner on 2026-10-28: unchanged. It is a one-shot check in the middle of the build, and a second clean doubt would double its cost.
- `/dh:plan-loop`. Decided by the owner on 2026-10-28: unchanged. Each wave is already judged by the full `reviewers.document` list.
- The `security` stage. Decided by the owner on 2026-10-28: unchanged. It already runs the full list and does not rotate.

Part (c) · The words that replace the `AGENTS.md` sentence "Build/QA e revisão de código: teto de 6 rodadas, e rodada sem achado bloqueante encerra o loop".

1. (Chosen) Replace it with this meaning: "Build/QA and code review: cap of 6 rounds. A round is clean when no blocking finding is open, reopened or escalated and every dispatched reviewer returned a verdict. Two consecutive clean rounds end the loop. A churn finding or the cap ends it early, and the Coordinator reports to the owner." The sentence on documents stays unchanged, as part (b) decides. Consequence: one rule for the code loop, in one place.
2. Keep the old sentence for the default review and use the new one only in the final feature review. Consequence: two rules in `AGENTS.md`, and a reader must know which one applies.
3. The owner gives other words. Consequence: R30 to R35 and `.agents/coordinator.md` line 119 follow them.

Decided by the owner on 2026-10-28: option 1. The pt-br text goes in `AGENTS.md` at build time, with the same meaning.

#### D2 · Entry point of the drift report and who runs the comparison (decided)

Context: the owner chose a skill with no new command and no new agent. The open points were whether a command offers the report, and which role compares.

Part (a) · The entry point.

1. (Chosen) An optional line that the Coordinator offers after `FEATURE APROVADA` in `/dh:review feature PRD-<n>`. It is never run alone. Consequence: PRD and code are both stable at that point. `.commands/review.md` changes: roles, skills, `writes` metadata and the Limits (today `writes: none` and read-only in line 35).
2. A subject of `/dh:document`. Consequence: its `writes` metadata already allows documents, but the command writes documentation and this report is a review.
3. A mode of `/dh:verify`. Consequence: `qa-verifier` writes tests, and a comparison of document and code is a review. This fits worst.
4. On request only, with no command text. Consequence: no command file changes and the owner must remember to ask. The tutorials describe how.

Decided by the owner on 2026-10-28: option 1. Note: the final review already reads adherence from QA evidence, and the drift report reads the code. The two checks differ.

Part (b) · Who compares.

1. (Chosen) `reviewer` in mode `code`, with the skill `drift-report` named in the dispatch (R39). Consequence: Sonnet, read-only, a fresh context. One sentence is added to `.agents/reviewer.md`.
2. `document-validator` with the skill. Consequence: Opus, read-only, and its mission is already "document against source". It costs more, and its output is a verdict.
3. A third mode `drift` of `reviewer`. Consequence: a larger edit of `.agents/reviewer.md` (description, procedure, output) and no new agent.

Decided by the owner on 2026-10-28: option 1.

#### D3 · Where the credit line goes (decided)

Context: R49 already puts the credit in the provenance note of each skill that takes a mechanism, as the kit does for other sources. The owner asked for the credit in one line of the docs. Neither README has a credits section today.

1. The provenance notes only. Consequence: no document line, so the owner's request is not met.
2. (Chosen) The provenance notes and one line in `README.md` and `README.en.md`. Consequence: the credit is visible at the entry point, in both languages.
3. The provenance notes and one line in the tutorials (`docs/tutorial.html`, `docs/en/tutorial.html`). Consequence: the credit sits where the new flow is described, and a reader of the README does not see it.

Decided by the owner on 2026-10-28: option 2.

## Executed summary

- **Delivered:** Reviewers read a closed context bundle instead of the Coordinator's conversation, the document validator states its own comprehension, code-review findings are triaged by id and state with churn detection and a two-consecutive-clean-rounds convergence rule, and a new read-only skill `drift-report` compares a document with the code, credited to EGM and released as 0.34.0 with 22 skills.
- **Rules:**

  | Rule | Verdict | Evidence |
  |---|---|---|
  | R1 | Honored | T-2901, T-2904, T-2908; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R2 | Honored | T-2901, T-2904, T-2908; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R3 | Honored | T-2901, T-2904, T-2908; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R4 | Honored | T-2901, T-2904, T-2908; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R5 | Honored | T-2901, T-2904, T-2908; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R6 | Honored | T-2901, T-2904, T-2907; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R7 | Honored | T-2904; T-2911 matrix PASS (text in `coordinator.md`); first use: the Dispatches entries of AUTO-2 in the Coordinator's record list the bundle elements (see note g); re-checked in the final QA (`qa-final`) |
  | R8 | Honored | T-2901, T-2904; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R9 | Honored | T-2904; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R10 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R11 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R12 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R13 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R14 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R15 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R16 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R17 | Honored | T-2904; T-2911 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R18 | Honored | T-2901; T-2911 matrix PASS (read and grep) |
  | R19 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R20 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R21 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R22 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R23 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R24 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R25 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R26 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R27 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R28 | Honored | T-2902, T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R29 | Honored | T-2905; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`) |
  | R30 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a |
  | R31 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a |
  | R32 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a |
  | R33 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a; `reviewer.md:81` was corrected after the matrix and reconfirmed (note b) |
  | R34 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a |
  | R35 | Honored | T-2906, T-2908, T-2902; T-2912 matrix PASS (read and grep); re-checked in the final QA (`qa-final`); see note a |
  | R36 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R37 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R38 | Honored | T-2903; T-2913 matrix PASS (read); re-checked in the final QA (`qa-final`) |
  | R39 | Honored | T-2902, T-2903, T-2906; T-2913 matrix PASS (one sentence in `reviewer.md`; 13 agent files); re-checked in the final QA (`qa-final`) |
  | R40 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R41 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R42 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R43 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R44 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R45 | Honored | T-2903, T-2906, T-2907, T-2908; T-2913 matrix PASS (read); re-checked in the final QA (`qa-final`) |
  | R46 | Honored | T-2903, T-2906, T-2907, T-2908; T-2913 matrix PASS (read); re-checked in the final QA (`qa-final`) |
  | R47 | Honored | T-2903, T-2906, T-2907, T-2908; T-2913 matrix PASS (read); re-checked in the final QA (`qa-final`); the first T-2913 run failed (see note b), reconfirmed |
  | R48 | Honored | T-2903; T-2913 matrix PASS (read) |
  | R49 | Honored | T-2901, T-2902, T-2903; frontmatter of the three skills carries the source string; matrices PASS, T-2914 PASS |
  | R50 | Honored | T-2901, T-2902, T-2903; word-run comparison with the public EGM files (see note c) |
  | R51 | Honored | T-2909; read of both READMEs PASS; tutorials carry no credit line |
  | R52 | Honored | T-2903; `grep -rn` of the EGM repository name: only sources entries, READMEs, CHANGELOG (and generated `dist/`); both PASS |
  | R53 | Honored | T-2909, T-2910, T-2903; version places, `CHANGELOG.md` head, counts 22/13/19+1; `dh validate` 0 errors; T-2914 PASS |
  | R54 | Honored | T-2909, T-2910, T-2903; version places, `CHANGELOG.md` head, counts 22/13/19+1; `dh validate` 0 errors; T-2914 PASS |
  | R55 | Honored | process (no build ticket); clean-clone proof (see note d) |
  | R56 | Honored | constraint (no build ticket); diff does not touch the isolation references, the clause, `hooks/` or `adapters/`; both PASS (see note e) |
  | R57 | Honored | T-2908; read of the `AGENTS.md` diff PASS; re-checked in the final QA (`qa-final`) |
  | R58 | Honored | T-2909; Docs list checked against `git diff main..HEAD --name-only`, PASS |
  | R59 | Honored | T-2906; `evals/` diff empty, PASS; cost note recorded (see note f); re-checked in the final QA (`qa-final`) |
  | R60 | Honored | T-2901, T-2902, T-2903; private-path scan and `dh validate` 0 errors; all PASS |

  Notes on the table. Evidence paths are under `~/.harness/projects/dev-harness/tasks/PRD-029/evidence/` unless stated.
  - (a) R30 to R35 are Honored in the text. The loops of this batch ran under the previous rule (one round without a blocking finding ends the loop), because the new rule applies from the merge on. The first application of the two-clean-rounds rule comes after the merge.
  - (b) Two checks failed once and were fixed. T-2913 R47: the `AGENTS.md:32` bullet lacked the name `<date>_DRIFT_REPORT.md`; the Coordinator corrected one line and the rule was reconfirmed. T-2912 R33 (recommended by the QA): `reviewer.md:81` states that a confirming round reads the whole change; the Coordinator corrected one line and reconfirmed it. The final QA (`qa-final`) re-checked both.
  - (c) R50 is Honored. The word-run comparison used the public EGM files, fetched read-only at commit f0ba151 (`egm/MANIFEST.md`). The longest shared run is 6 words and the limit is 8. Evidence: `qa-T-2911/r50.txt`, `qa-T-2912/r50.txt`, `qa-T-2913/r50-output.txt`. The 6-word runs are "reviewer raises a settled finding again", "system works as the document" (an item name the PRD requires), and "reverting that fix would resolve".
  - (d) R55 is Honored. T-2914 ran a clean clone of `3bb762c`, `go run ./cmd/dh build` there, and compared sha256: the five binaries are identical in the tree and in the clone. The result is recorded in the Coordinator's record (`MEMORY.md`, entry "AUTO-2 COMMIT CANDIDATO e T-2914"); T-2914 saved no output file. The proof is repeated on the final commit before the remote publication, with a saved log.
  - (e) R56 is Honored. One `/dh:secure` pass ran on T-2901 as an execution-surface precaution decided by the Coordinator (full `reviewers.security` list; no finding; evidence `security-T-2901/`). The isolation references, the clause, `hooks/` and `adapters/` are untouched.
  - (f) R59 is Honored. No paid eval ran and `evals/` is untouched. The first-use cost note (informal, not measured) is recorded in the Coordinator's record after the final review.
  - (g) R7 is process evidence at first use: the AUTO-2 entries of the Coordinator's record name the bundle elements of each dispatch.
  - (h) Limits, stated plainly: the QA matrices of T-2911 to T-2913 are read and grep assertions, not executable tests. The second review of several tickets read only the delta (stated in each record). The final review panel (`review-final-codex`, `review-final-opencode`) approved with no blocking finding.

- **Tickets:** All tickets are done on 2026-10-28. Evidence paths are under `~/.harness/projects/dev-harness/tasks/PRD-029/evidence/`.
  - T-2901, lane backend (document stage, bundle, `## Comprehension`), done. Evidence: `doubt-T-2901/` (approve), `review-T-2901/` (approve), `security-T-2901/` (no finding), QA `qa-T-2911/`.
  - T-2902, lane backend (code stage, re-entry, settled findings, `Churn:`), done. Evidence: `review-T-2902/` (one Major, fixed), `review-T-2902-r2/` (approve), QA `qa-T-2912/`.
  - T-2903, lane backend (skill `drift-report`, roadmap counts), done. Evidence: `review-T-2903/` (one Major, fixed), `review-T-2903-r2/` (approve), QA `qa-T-2913/`.
  - T-2904, lane backend (Coordinator bundle and Comprehension read), done. Evidence: `review-batch-C/` (no finding), QA `qa-T-2911/`.
  - T-2905, lane backend (Coordinator triage), done. Evidence: `review-batch-C/` (no finding), QA `qa-T-2912/`.
  - T-2906, lane backend (convergence, drift-report run, cost note), done. Evidence: `review-batch-C/` (one Major, fixed), `review-T-2906-r2/` (approve), QA `qa-T-2912/`.
  - T-2907, lane backend (`/dh:review` offer and convergence), done. Evidence: `review-batch-C/` (no finding), QA `qa-T-2913/`.
  - T-2908, lane backend (`AGENTS.md` rules), done. Evidence: `review-batch-C/` (no finding), QA `qa-T-2911/`, `qa-T-2912/`, `qa-T-2913/` (R47 fixed, note b).
  - T-2909, lane backend (docs of the batch), done. Evidence: `review-batch-D/` (a `request changes` line with 0 findings, caused by missing check results), `review-batch-D-r2/` (approve), T-2914.
  - T-2910, lane infra (version 0.34.0, marketplace text, rebuild), done. Evidence: `review-batch-D-r2/` (approve), T-2914 (`go test ./...` 7 packages ok, `dh validate` 0 errors).
  - T-2911, lane teste (unitário) (R1 to R18, R49, R50, R60), done. Evidence: `qa-T-2911/` (21 PASS and 4 extra checks PASS).
  - T-2912, lane teste (unitário) (R19 to R35, R49, R50, R60, negative scope), done. Evidence: `qa-T-2912/` (all PASS, 5 negative-scope rows PASS).
  - T-2913, lane teste (unitário) (R36 to R48, R50, R52, R56), done. Evidence: `qa-T-2913/` (R47 FAIL fixed and reconfirmed, rest PASS).
  - T-2914, lane teste (integração) (final gate, clean clone, R49, R51 to R60), done. Evidence: the Coordinator's record (gate green, 5 binaries identical, note d); `qa-final/` (0 failures, `validate` 0 errors with 13 agents, 19 commands, 22 skills).

- **Docs:** `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html` (one sentence in each corrected after the final review), `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md`, `.claude-plugin/marketplace.json`, `.agents/coordinator.md`, `.agents/document-validator.md`, `.agents/reviewer.md`, `.skills/document-review/SKILL.md`, `.skills/code-review/SKILL.md`, `.skills/external-clis/SKILL.md`, `.skills/drift-report/SKILL.md` (new), `.commands/review.md`, `docs/prd/PRD-024-public-readiness.md` (one-line count wording, not in the Docs list of section 4: cleanup of a text this change made obsolete), this PRD and its review record `PRD-029-independent-review-and-drift-report.review.md`. Not documents: `internal/build/build.go` (version constant), `harness-manifest.json` and `dist/` (generated).

- **Out of scope:** The brainstorm of EGM with six lenses and its section-by-section draft; the folder READMEs of EGM ("bootstrap", "peanuts" and "hay"); any paid eval, new eval case or new probe; copying EGM text and any runtime dependency on EGM; a stop rule for hosts without fresh subagents; the hooks and plugin files of EGM for its host; a change to the isolation of CLI reviewers, to the PRD-027 guard or to the anti-delegation clause. Also out of scope: the comprehension block in code reviewers (document validator only); a change to the memory system; no new agent, command, hook, template or eval case.
