# PRD-020 · PRD in English with ASD-STE100 writing rules, and `language` as the owner conversation language

| Field | Value |
|---|---|
| Status | delivered on 2026-10-13 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-13 / 2026-10-13 |
| Tickets | T-2001 |
| Origin | Owner chat, 2026-10-13 (the owner delegated: "follow your recommendations"; decision 1 changed the same day) |
| Depends on | PRD-007 (PRD and ticket flow), PRD-019 (version and `ref` pattern) |
| Enables | nothing |

## 1. Problem

Facts verified on 2026-10-13:

- The PRD is the contract of the kit. `product-discovery` writes it from `templates/<lang>/PRD.md`.
- The template and the skill `requirements-discovery` give no writing style. Long sentences and variable terms make a rule ambiguous. Two builders then build different software.
- `document-review` blocks an ambiguity only when the two readings give different software. It cites no writing rules.
- ASD-STE100 (Simplified Technical English) is an English standard. The official approved-word dictionary is copyrighted by ASD.
- Today `language` in `.harness/project.yaml` picks the PRD language. `.agents/coordinator.md` (Language section) tells the Coordinator to reply in the language the owner writes in. `.skills/project-onboarding/SKILL.md` (line 66) sets `language` from the owner language that setup detects.
- The owner said: write the PRD always in English with the standard. The configuration must remind the agent which language to use between the orchestrator and the operator.

## 2. Solution

- Every PRD that the kit writes is in English. It follows short STE100-inspired writing rules, whatever the project `language`.
- One short rules block holds the rules. It lives in the skill `requirements-discovery`. The PRD template and the skill `document-review` point to it. They do not repeat it.
- `document-validator` checks the style. A violation that makes a rule or an acceptance ambiguous blocks. A style deviation that keeps the meaning is Minor.
- The kit does not copy the ASD dictionary. It summarizes the rules in its own words and cites the standard by name.
- `templates/pt-br/PRD.md` stays as a byte-identical copy of `templates/en/PRD.md`. The template count stays 14. No Go change is needed for the template.
- `dh validate` and the status-drift mod accept both `entregue em` and `delivered on`.
- The PRD header, its section names and its "Executed summary" always use the English forms: `delivered on <date>`, Delivered, Rules, Tickets, Docs, Out of scope.
- The skill, the Coordinator and `product-discovery` stop saying "the project language picks the form" for the PRD.
- `language` in `.harness/project.yaml` is the language of the conversation between the Coordinator and the owner. The Coordinator reads it at execution start and replies to the owner in that language. Specialist replies and dispatches stay in English. The Coordinator translates.
- The owner switches `language` by explicit request. That request is the authorization: the Coordinator edits `language` in `<harness dir>/project.yaml` and switches, with no second confirmation. `dh config set language` only sets the machine default for new projects. The Coordinator mentions it only if the owner asks for a default. No new config key exists.
- `language` still picks `templates/<lang>/` for the other records (MEMORY, EPOCHAL, RISKS, TASK, ADR, RFC). Tickets keep following `language`.

## 3. Rules

- R1 The skill `.skills/requirements-discovery/SKILL.md` has one rules block of at most 20 lines for PRD writing. The block states these rules:
  - Procedural sentences have at most 20 words. Descriptive sentences have at most 25.
  - One instruction per sentence. Active voice.
  - One term per concept. No synonyms.
  - Simple tenses. Verbs with one sense. Articles kept.
  - No idiom or phrasal verb when one plain verb exists.
  - Lists for sequences. Short paragraphs.

  Verifiable: `grep -n "20 words" .skills/requirements-discovery/SKILL.md` and `grep -n "25 words" ...` each find one line; the block has at most 20 lines.
- R2 The block says that ASD-STE100 is an English standard. It says that the block is a summary in the kit's own words. It says that every PRD is in English whatever the `language`. The block has no pt-br adaptation. Verifiable: `grep -n "ASD-STE100" .skills/requirements-discovery/SKILL.md` finds the line; `grep -rn -i "estilo STE\|pt-br adaptation" .skills .agents templates` is empty.
- R3 The kit does not copy the approved-word dictionary or text of the standard. The block cites the standard by name, with no issue. Verifiable: `grep -rn -i "approved word\|approved-word\|dictionary" .skills/requirements-discovery/SKILL.md` finds only the sentence that says the dictionary is not copied; no list of approved words exists in `.skills/`, `templates/` or `.agents/`.
- R4 `templates/pt-br/PRD.md` is a byte-identical copy of `templates/en/PRD.md`. `templates/en/PRD.md` has one comment line that points to the block of R1 and does not repeat the rules. Verifiable: `cmp templates/en/PRD.md templates/pt-br/PRD.md` is silent; `grep -c "requirements-discovery" templates/en/PRD.md` = 1; `grep -c "20 words" templates/en/PRD.md` = 0.
- R5 `.skills/document-review/SKILL.md` has a criterion that points to the block of R1 when the document is a PRD. A violation that makes a rule or an acceptance ambiguous blocks. A style deviation that keeps the meaning is Minor and does not block. The skill repeats no numeric limit. Verifiable: `grep -n "requirements-discovery" .skills/document-review/SKILL.md` finds the criterion; `grep -c "20 words\|25 words" .skills/document-review/SKILL.md` = 0.
- R6 The PRD is always English in the sources.
  - `.skills/requirements-discovery/SKILL.md` says the brief becomes the PRD from `templates/en/PRD.md`.
  - `.agents/coordinator.md` (Executed summary), `.agents/product-discovery.md` (amendment branch and shared contract), `.agents/docs-guide.md` (line 59) and `.commands/document.md` (line 18) use the English forms for the PRD.
  - These files no longer say that the project language picks the form or the template folder.
  - `.agents/product-discovery.md` (line 53) and `.commands/review.md` (line 19) name the section "Solution", not "Solução".
  - `AGENTS.md` ("Idioma em execução") says that the Coordinator replies in the language recorded in `language`.

  Verifiable: `grep -rn "templates/<lang>/PRD" .skills .agents .commands` is empty; `grep -n "project .language. picks the form" .agents/coordinator.md` is empty; `grep -n "Solução" .agents/product-discovery.md .commands/review.md` is empty; the amendment branch names the section "Docs or Rules" only; `grep -n "idioma do dono" AGENTS.md` is empty.
- R7 `language` is the owner conversation language.
  - `.agents/coordinator.md` (Language section) says that the Coordinator reads `language` at execution start.
  - It says that the Coordinator replies to the owner in that language and keeps dispatches and specialist replies in English.
  - It says that, on an explicit owner request to switch, the Coordinator edits `language` in `<harness dir>/project.yaml` with no second confirmation. That request is the authorization.
  - `.skills/project-onboarding/SKILL.md` (line 66) and `.commands/setup.md` describe `language` as the owner conversation language and as the template choice for the other records.

  Verifiable: `grep -n "conversation language" .agents/coordinator.md .skills/project-onboarding/SKILL.md .commands/setup.md` finds all three; `grep -n "the language the owner writes in" .agents/coordinator.md` is empty; `grep -n "switch only after agreement" .agents/coordinator.md` is empty.
- R8 The template count and the mirror rule stay true. The count stays 14. `AGENTS.md` keeps "every new template enters in both" and adds that the PRD template is English in both folders. Verifiable: `dh validate` green; `cmp templates/en/PRD.md templates/pt-br/PRD.md` is silent; `ls templates/*/ | wc` is unchanged.
- R9 Scope. The TASK, ADR and RFC templates do not change. The existing PRDs in `docs/prd/` are not rewritten. The dashboard code does not change: it already parses `delivered on`. Verifiable: `git diff --name-only main...HEAD` lists no `templates/*/TASK.md`, `ADR.md`, `RFC.md`, no old PRD and no `internal/dashboard/` file.
- R10 Go changes are limited to three items.
  - The `defaultVersion` line in `internal/build/build.go`.
  - `internal/kit/status_drift.go` accepts both `entregue em` and `delivered on`, as `internal/dashboard/metrics.go` (line 91) already does. A Go test covers `delivered on`.
  - The mod `adapters/claude-code/plugin/hooks/status-drift.ts` accepts both forms. A test in `status-drift.test.ts` covers it.

  Reason: a PRD closed as `delivered on <date>` fails both status-drift checks today. `dh validate` gets no style check, and `go.mod` stays unchanged. Verifiable: `git diff main -- '*.go'` shows only `build.go`, `status_drift.go` and its test; `git diff go.mod` is empty; `go test ./...` is green; `claude plugin test` is green in `dist/claude-code/dev-harness`.
- R11 This PRD follows its own rules. Rule lines that carry a `Verifiable:` command are exempt from the length limit. Verifiable: the validator finds no other sentence above 20 words (procedure) or 25 words (description).
- R12 Version 0.27.0. `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json` and the source `ref` (`v0.27.0`) in the same file change together. A forgotten `ref` is a delivery failure. Verifiable: `grep` of the three values shows `0.27.0` or `v0.27.0` in one commit; `go run ./cmd/dh build` shows no `dist/` difference except the expected one.
- R13 Clean-clone proof applies, because the PR touches `internal/` and `dist/*/bin`. Verifiable: `git clone . <tmp>`, `go run ./cmd/dh build` there, and equal `sha256sum` of the binaries.
- R14 `CHANGELOG.md` has a 0.27.0 entry. `docs/roadmap.html` and `docs/en/roadmap.html` reflect it. Verifiable: `dh validate` green.

## 4. Docs

Checked with `grep` on 2026-10-13. Only files that describe PRD writing, templates or `language` are listed.

- `AGENTS.md`, with three edits:
  - the mirror rule: the PRD template is English in both folders;
  - the "Idioma em execução" bullet: the Coordinator replies in the `language` language;
  - one short bullet for the STE block
- `README.md`, `README.en.md` (line "Templates de trabalho", the `language` bullet on lines 69-70, the version line 0.27.0)
- `docs/tutorial.html`, `docs/en/tutorial.html` (step `/dh:discover`; `language` description, line 823; the `/dh:setup` description, line 740)
- `docs/tutoriais/templates.html`, `docs/en/tutorials/templates.html` and their `templates.md` (PRD section and template list)
- `docs/roadmap.html`, `docs/en/roadmap.html` (line 678: reword "mirrored per language" only if it is false for the PRD; 0.27.0 entry)
- `CHANGELOG.md` (0.27.0)
- `.claude-plugin/marketplace.json` (`version` and `ref`)

## 5. Out of scope

- A change to binary distribution. ADR-001 stays and the binaries stay committed.
- A Go check or `dh validate` check for PRD style.
- Other templates: TASK, ADR, RFC, story.
- A rewrite of the existing PRDs in `docs/prd/`.
- A copy of the ASD dictionary or of any text of the standard.
- A new config key for the conversation language.

## 6. Owner decisions

- Owner, 2026-10-13: the PRD is always English and follows the standard. The configuration must remind the agent which language to use between the orchestrator and the operator. This replaces the first decision (pt-br adaptation).
- Verification is prose only. `document-validator` with `document-review` checks the style. An ambiguity blocks. A style-only deviation is Minor. No new style check exists.
- Coordinator decision, 2026-10-13, under the owner's delegation: the status-drift check in Go and in the mod accepts `delivered on`, because every PRD is English. This is a fix to an existing check, not a style check.
- Scope: PRD only. One rules block, in `requirements-discovery`. Copyright: own words, standard cited by name, no dictionary.
- Version 0.27.0 with clean-clone proof. Binaries stay committed.
- Coordinator decision, 2026-10-13: `templates/pt-br/PRD.md` stays as a byte-identical copy of the English template. This avoids a Go change in `checkTemplateRefs`.
- Coordinator decision, 2026-10-13: an explicit owner request to switch `language` is the authorization to edit `project.yaml`.
- Analyst assumptions, not owner decisions:
  - (a) The block lives in the skill, because template comments are deleted when the PRD is final.
  - (b) The length limits apply to running text. `Verifiable:` lines are exempt.
  - (c) The Coordinator does not use `dh config set language` for a project switch.
  - (d) The old pt-br PRDs keep their Portuguese headers.
- Pending decisions: none.

## Executed summary

- **Delivered:** every PRD the kit writes is English with the STE writing-rules block, and `language` is the owner conversation language.
- **Rules:**
  - R1 Honored: block in `.skills/requirements-discovery/SKILL.md` lines 100-111, 20/25-word limits on line 102.
  - R2 Honored: ASD-STE100 named as an English standard on line 101; "estilo STE" grep empty.
  - R3 Honored: the only dictionary mention says the dictionary is not copied.
  - R4 Honored: `cmp templates/en/PRD.md templates/pt-br/PRD.md` silent; one pointer line.
  - R5 Honored: `document-review` line 55; no numeric limit repeated.
  - R6 Honored: `templates/<lang>/PRD` and `Solução` greps empty; `idioma do dono` grep empty in `AGENTS.md`.
  - R7 Honored: "conversation language" in coordinator, onboarding and setup; `coordinator.md` line 74.
  - R8 Honored: `dh validate` green, 14 templates.
  - R9 Honored: no TASK/ADR/RFC template, old PRD or dashboard file in the diff.
  - R10 Honored: Go diff is `build.go`, `status_drift.go` and its test; `go test ./...` and `claude plugin test` (254) green.
  - R11 Honored: QA found long sentences; the PRD was amended and revalidated.
  - R12 Honored: `defaultVersion`, marketplace `version` and `ref` are 0.27.0 / v0.27.0.
  - R13 Honored: clean-clone build gives equal sha256 for all binaries; only `generated_at` differs in the manifests.
  - R14 Honored: CHANGELOG 0.27.0 and both roadmaps; `dh validate` green.
- **Tickets:** T-2001, backend, approved (QA matrix; code review `claude` then `cli:codex/gpt-5.6-luna` approve; final review FEATURE APROVADA).
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/tutoriais/templates.html`, `docs/tutoriais/templates.md`, `docs/en/tutorials/templates.html`, `docs/en/tutorials/templates.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`.
- **Out of scope:** binary distribution (ADR-001 stays), a `dh validate` style check, other templates, a rewrite of old PRDs.
