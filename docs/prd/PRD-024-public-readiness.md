# PRD-024 · Make the public repository ready for strangers

| Field | Value |
|---|---|
| Status | delivered on 2026-10-16 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-16 / 2026-10-16 |
| Tickets | T-2401 (planned) |
| Origin | Owner chat, 2026-10-16 (scope and decisions in chat) |
| Depends on | nothing |
| Enables | nothing |

## 1. Problem

- The GitHub repository is public. A stranger who opens it cannot tell what the kit is, how to remove it, or what it writes on the machine. Fact: audit UND-006 by `repo-scout`.
- `README.md` and `README.en.md` have these gaps:
  - Line 20 is a release-notes paragraph of about 2,500 characters. It holds the version, the VS Code extension, the dashboard, test counts ("226 pass") and kit counts. `CHANGELOG.md` already holds this content.
  - No section says what the kit is, what it needs, how to start, how to update, how to uninstall and remove data, what the kit writes and where, how to turn things off, or what to do when something fails.
  - Control information exists only in `docs/tutorial.html` and `docs/en/tutorial.html`, in scattered places (deactivate plugin, prune of mod or snapshot, stop of the dashboard).
  - Line 7 of both READMEs links `docs/plano-produto.html` as the product page. That page is in pt-br and is internal. The English README must not lead with it.
  - The maintenance block is for contributors only, but nothing says so.
  - Neither README names the two companion projects: `danielmalka/malka-claude-mods` (visual mods that moved out of this kit in 0.29.0) and `akitaonrails/ai-memory` (long-term memory and handoff across coding CLIs). The kit needs neither.
- Personal data in the tree:
  - 28 of 52 tracked `evals/baselines/*.json` files contain the absolute path of the owner's home (the repository checkout path). No other tracked file has the owner's home path. No secret or token was found.
  - "danlemos", the name of another project of the owner, appears in `docs/roadmap.html` and `docs/en/roadmap.html` (near line 1631). It means nothing to a stranger. It also appears in old PRDs and ADRs, which are history.
  - Both roadmaps cite `CLAUDE.md`, which only imports `AGENTS.md`. Base lines: pt 1543, 1632, 1698, 1780; en 1541, 1632, 1699, 1783.
- Both roadmaps cite `RISKS.md` in the Relations of RISK-002 (pt 921 and 1591; en 922 and 1589). That record lives in the maintainer's harness record, not in this repository. Lines 760 and 839 of both files also cite `RISKS.md`. Their prose is correct (consumer-project records) and stays.
- The line numbers are on the base. The builder finds each line again by its text.
- Public-repository files: only `LICENSE` and `CHANGELOG.md` exist. There is no `SECURITY.md` and no `CONTRIBUTING.md`.
- The roadmap has no card for three deferred items from `docs/discovery/renova.md` (section 7): a form for `config.yaml` in the dashboard, animated diagrams for a public showcase, and a per-project `events.jsonl`.
- The roadmap has no card for a request of the owner. The Coordinator does too much hands-on work itself. It reads files, runs build, validate and grep, edits records, assembles CLI review prompts and spot-checks by shell. It should dispatch subagents instead. Fact: owner observation, 2026-10-16.

## 2. Solution

One release, 0.30.1, in one PR. After delivery, a stranger can read what the kit is, install it, remove it, find what it wrote, and report a vulnerability. The tree holds no path of the owner's machine. Add no kit behavior. Edit no file in `.agents/`, `.skills/` or `.commands/`.

README (pt-br in `README.md`, English in `README.en.md`). Both files have the same 12 `## ` sections in this order. The heading texts are fixed:

| # | English heading | pt-br heading | Content |
|---|---|---|---|
| 1 | What it is | O que é | At most 3 lines. |
| 2 | Requirements | Requisitos | Claude Code; Git. |
| 3 | Install | Instalar | The marketplace URL form that exists today. |
| 4 | Quick start | Início rápido | The real command names from `docs/tutorial.html` and `.commands/` (for example `/dh:doctor`, `/dh:setup`). It ends with at most 2 lines on language. |
| 5 | Update | Atualizar | How to update the plugin. |
| 6 | Uninstall and remove data | Desinstalar e remover dados | `/plugin uninstall`; the contents of `~/.harness`; the `additionalDirectories` entry; the `.harness/` folder of a repository. |
| 7 | What the kit writes and where | O que o kit grava e onde | One table with 7 rows: repository `.harness/`; `~/.harness/config.yaml`; `projects/<name>/`; `sessions/`; `dashboard/`; `pdocs/`; the dashboard on `127.0.0.1:4747`. |
| 8 | Turn things off | Desligar | `dh dashboard --stop`; disable the plugin. One sentence says there is no per-hook switch for the mod or the commit trap (pt: "por hook"). |
| 9 | If something fails | Se algo falhar | `dh doctor` or `/dh:doctor`; port busy; plugin not found; link to GitHub Issues. |
| 10 | Use together with | Usar junto com | The two companion projects, with links and one line each on what they add. The text says they are optional (pt: "opcional"). |
| 11 | Docs map | Mapa da documentação | Quick start, tutorial, roadmap, `CHANGELOG.md`. `docs/prd`, `docs/adr` and `docs/discovery` are design history in pt-br. `evals/` is measurement history and is not needed to use the kit. `docs/plano-produto.html` is labeled pt-br and internal, and the English README does not link it first. |
| 12 | For contributors | Para quem contribui | The maintenance block (`go test ./...` and the build commands), the sources-and-package table, and the clone-load instructions. |

Where the current README blocks go:

- Language block: its facts shrink to at most 2 lines at the end of Quick start.
- "Load from the clone" and "Sources and package": move into section 12.
- The 0.8.0 migration note: removed. `CHANGELOG.md` keeps it.
- Global mode paragraph: one line before the table in section 7. The tutorial holds the details.
- Release-notes paragraph (line 20): removed. `CHANGELOG.md` is its source.

The README holds no test count. A kit count (13 agents, 19 commands, 21 skills) appears in at most one place: one line that has the word "atual", "esperado", "current" or "expected" and at least two counted nouns. `checkDocumentedCounts` in `internal/kit/doc_counts.go` checks only such a line against the disk. A count on any other line is not checked, so the README does not state one.

Roadmap (`docs/roadmap.html` and `docs/en/roadmap.html`, mirrored by hand):

- Four cards in the Planned section (`<h2 id="planejado">`). No `<h3>` starts with a version number. Exact `<h3>` text (en / pt-br):
  - (a) "Edit config.yaml from the dashboard" / "Editar o config.yaml pelo dashboard". It needs a CSRF or session token and a full security review. PRD-018 delivered the `dh config` write side.
  - (b) "Animated diagrams for the public showcase" / "Diagramas animados para a vitrine pública". It starts after `pdocs` prove the content.
  - (c) "Per-project events.jsonl for finer metrics" / "events.jsonl por projeto para métricas mais finas". This is PRD-016 option B. The granularity is one day today.
  - (d) "Coordinator orchestrates, does not execute" / "O Coordenador orquestra, não executa". The Coordinator dispatches subagents for execution and verification. Execution covers reading, build, validate, grep, record edits, CLI review prompts and spot-checks. The Coordinator only judges whether the delegated work was done with skill or needs rework. It still owns the writes to `MEMORY.md`, `EPOCHAL.md` and `RISKS.md`. Status: planned, no PRD yet. Levers to evaluate later: the rules and spot-check wording in `.agents/coordinator.md`, and maybe a mod guard. Relations: `.agents/coordinator.md`.
  - Each card has Description, Acceptance criteria, Status and Relations, as the existing cards do.
- "danlemos" becomes generic wording in the two roadmaps only. The history documents stay unchanged.
- Dead references: every `CLAUDE.md` hit (lines in section 1) becomes `AGENTS.md`. In the RISK-002 Relations, the text says that RISK-002 lives in the maintainer's harness record, not in this repository. The `RISKS.md` lines 760 and 839 stay unchanged.
- A Delivered card for 0.30.1 in both files.

Other files:

- `evals/baselines/*.json`: the prefix of the owner's checkout path becomes `<repo>`. The git history is not rewritten.
- `SECURITY.md` in English: report a vulnerability through GitHub private vulnerability reporting; a short policy.
- `CONTRIBUTING.md` in English: canonical sources (`.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`); `dist/` comes from `go run ./cmd/dh build` and is never edited by hand; checks `go test ./...`, `go run ./cmd/dh validate .` and `go run ./cmd/dh build`; Conventional Commits in English with no mention of AI; `AGENTS.md` is the canonical rules file and is in pt-br.
- Release 0.30.1: `CHANGELOG.md` entry, the Delivered card, `version` and `ref` in `.claude-plugin/marketplace.json`, `defaultVersion` in `internal/build/build.go`.

## 3. Rules

Each check below fails on the tree at the time of writing, except R7, R11, R17, R18, R19 and R21. Those six guard against a regression and must hold at the end. The ticket author checks again on the base before the first edit.

- R1 `README.md` and `README.en.md` have the 12 sections of section 2 in the same order, and each file has at most 110 lines. Verifiable: `grep '^## ' README.en.md` prints the 12 English headings of section 2 in that order, `grep '^## ' README.md` prints the 12 pt-br headings in that order; `wc -l README.md README.en.md` shows 110 or less for each file.
- R2 No README holds the release-notes paragraph or a test count. Verifiable: `grep -n -E 'Versão atual|Current version|226|Testes:|Tests:' README.md README.en.md` has no hit, and no line of the two files is longer than 600 characters (`awk 'length>600' README.md README.en.md` prints nothing).
- R3 Both READMEs name the two companion projects with links and say they are optional. Verifiable: `grep -c 'github.com/danielmalka/malka-claude-mods' README.md README.en.md` and `grep -c 'github.com/akitaonrails/ai-memory' README.md README.en.md` print 1 or more for each file, and the "Use together with" / "Usar junto com" section holds "optional" / "opcional".
- R4 Both READMEs state what the kit writes and where, how to uninstall and remove data, and how to turn things off. Verifiable, per section (read the text between the section heading and the next `## `):
  - Uninstall section: has `/plugin uninstall`, `~/.harness`, `.harness/` and `additionalDirectories`.
  - Writes section: has 7 body rows (after the header and the separator) that start with `|`, naming `.harness/`, `config.yaml`, `projects/`, `sessions/`, `dashboard/`, `pdocs/` and `4747`.
  - Turn-off section: has `dh dashboard --stop`, a plugin step (the disable form that QA verified, or `/plugin uninstall` when `claude plugin --help` shows no disable form), and the sentence with `per-hook` (pt: `por hook`).
- R5 The Docs map labels `docs/plano-produto.html` as pt-br and internal, and `README.en.md` does not link it before the quick start. Verifiable: `grep -n 'plano-produto' README.en.md` shows a hit only inside the Docs map section, with the label; no hit on line 7.
- R6 The maintenance block sits under a heading or lead line that says it is for contributors only. Verifiable: the section 12 heading contains "contribut" (pt: "contribu") and holds the `go test ./...` command block.
- R7 No README holds an unchecked or wrong kit count. Verifiable: `go run ./cmd/dh validate .` passes, and `grep -n -i -E '[0-9]+ (agents|agentes|commands|comandos|skills|profiles|perfis|templates)' README.md README.en.md` has no hit, or its hits are on one line per file that has "atual", "esperado", "current" or "expected" and at least two counted nouns.
- R8 No tracked file holds the path of the owner's machine. Verifiable: `git grep -l '/ho[m]e/ma[l]ka'` over the whole tree prints nothing (the bracket classes match the path, and the PRD text does not hold the literal). Each of the 28 changed JSON files still parses: `for f in evals/baselines/*.json; do python3 -m json.tool "$f" >/dev/null || echo bad $f; done` prints nothing. The replacement is exactly `<repo>`: `git grep -l '<repo>' -- evals/baselines` lists 28 files.
- R9 `SECURITY.md` exists in English and names GitHub private vulnerability reporting. Verifiable: `grep -n -i 'private vulnerability reporting' SECURITY.md` has a hit.
- R10 `CONTRIBUTING.md` exists in English and has the 5 items of section 2. Verifiable: `grep -n -E '\.agents/|go run \./cmd/dh build|go run \./cmd/dh validate|go test \./\.\.\.|Conventional Commits|AGENTS\.md' CONTRIBUTING.md` finds each term, and `grep -n -i -E 'co-authored-by|generated with' CONTRIBUTING.md` has no hit except the rule line that forbids mention of AI. The names Claude Code and `.claude-plugin` are allowed.
- R11 No `CODE_OF_CONDUCT.md`, no `.github/ISSUE_TEMPLATE` and no pull request template are added. Verifiable: `git diff --name-status main...HEAD` lists none of them.
- R12 Both roadmaps hold the four Planned cards, mirrored. Verifiable: in each file, the lines between `id="planejado"` and the next `<h2` contain the four exact `<h3>` texts of section 2 (en in `docs/en/roadmap.html`, pt-br in `docs/roadmap.html`), the Relations of card (d) name `.agents/coordinator.md`, and the two files have the same number of `<div class="card">` in that section. No `<h3>` of the four cards starts with a version number.
- R13 "danlemos" is gone from the two roadmaps and stays in the history documents. Verifiable: `grep -c danlemos docs/roadmap.html docs/en/roadmap.html` prints 0 for both; `git diff --name-only main...HEAD -- docs/adr docs/prd docs/discovery` lists only this PRD and its review file.
- R14 The two roadmaps have no dead reference. Verifiable: `grep -n 'CLAUDE\.md' docs/roadmap.html docs/en/roadmap.html` prints nothing; `grep -n 'RISKS\.md' docs/roadmap.html docs/en/roadmap.html` prints exactly 3 lines per file: base lines 760 and 839 (found again by text) and the line of card (d). No RISK-002 Relations line cites `RISKS.md`. The 0.30.1 Delivered card writes neither `CLAUDE.md` nor `RISKS.md`.
- R15 Version 0.30.1 in three places: `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json`, and the source `ref` (`v0.30.1`). Verifiable: `grep` of the three values shows `0.30.1` or `v0.30.1`.
- R16 `CHANGELOG.md` has an entry for 0.30.1, and both roadmaps have a Delivered card for 0.30.1 that agree with it. Verifiable: `grep -n '0.30.1' CHANGELOG.md docs/roadmap.html docs/en/roadmap.html` has a hit in each file, and `dh validate` reports no roadmap drift.
- R17 The change edits no file in `.agents/`, `.skills/` or `.commands/`, and adds no agent, skill or command. Counts stay at 13 agents and 19 + 1 commands. Verifiable: `git diff --name-only main...HEAD -- .agents .skills .commands` prints nothing.
- R18 Hand edits are limited to the files in section 4, this PRD with its review file, and the ticket records. No file in `dist/` is edited by hand. Verifiable: `git diff --name-only main...HEAD -- dist` lists only files that `go run ./cmd/dh build` regenerates (`plugin.json`, `harness-manifest.json`, `dist/*/bin/**`).
- R19 `go run ./cmd/dh validate .` passes with 0 errors, and `go test ./...` is green. The build runs twice, and the second run changes nothing in `dist/` besides the `generated_at` line of `harness-manifest.json`. Verifiable: run the commands and compare with `git diff`.
- R20 Clean-clone proof, because the PR touches `internal/build/build.go` and `dist/*/bin`: run `git clone . <tmp>`, run `go run ./cmd/dh build` there, and compare `sha256sum` of the binaries with the tree. Verifiable: 5 of 5 sums are equal.
- R21 The change does not rewrite git history. Verifiable: no `--force` push, and `git merge-base --is-ancestor main HEAD` succeeds.
- R22 Every command in the README sections Update, Uninstall and Turn things off counts as verified in the form that the README prints (slash or CLI). A command is verified when the docs of this repository show it, or when QA ran it in this PR. The CLI forms already run by the Coordinator are `claude plugin marketplace update dev-harness` and `claude plugin update dh@dev-harness`. QA checks `claude plugin --help` for a disable form. If the help shows none, the Turn things off section uses `/plugin uninstall` as its plugin step, and R4 accepts it. Verifiable: the QA report lists each README command with its source (doc path or the run output).

## 4. Docs

Update:

- `README.md`, `README.en.md` (rewritten as in section 2).
- `docs/roadmap.html`, `docs/en/roadmap.html` (four Planned cards, 0.30.1 Delivered card, "danlemos", dead references).
- `CHANGELOG.md` (entry 0.30.1).
- `.claude-plugin/marketplace.json` (`version` and `ref`).
- `internal/build/build.go` (`defaultVersion`).
- `evals/baselines/*.json` (28 files, path prefix).
- `SECURITY.md`, `CONTRIBUTING.md` (new).

Not affected (the ticket author checks with `grep`):

- `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/inicio-rapido.html`, `docs/en/quick-start.html`. The README links to them. It does not change them.
- `docs/plano-produto.html`. It stays as is. Only the README labels it.
- `docs/adr/`, `docs/prd/` (other PRDs) and `docs/discovery/`. History. They keep "danlemos" on purpose.

## 5. Out of scope

- Rewrite of git history to remove the path from old commits. Accepted by the owner.
- Translation of `AGENTS.md`, PRDs or ADRs.
- `CODE_OF_CONDUCT.md`, issue templates, pull request templates.
- Move of `evals/` or of `*.review.md` files.
- Any change of kit behavior, and any edit in `.agents/`, `.skills/` or `.commands/`.
- A new per-hook switch for the mod or the commit trap.
- Implementation of the four Planned items, including any edit of `.agents/coordinator.md`. Only their roadmap cards are in scope.
- GitHub description and topics. The Coordinator set them with `gh repo edit`, outside the PR.
- Change of the holder name in `LICENSE`.

## 6. Owner decisions

- Owner, 2026-10-16: release 0.30.1 for all of this, in one PR. Reason: docs never ship in a docs-only PR, and every merge is tagged `v<version>` (`AGENTS.md`). The release changes the `CHANGELOG.md` entry, the Delivered card in both roadmaps, the marketplace `version` and `ref`, and `defaultVersion`. Clean-clone proof applies (R20).
- Owner, 2026-10-16: the holder in `LICENSE` stays "Daniel Lemos". No change.
- Owner, 2026-10-16: in `evals/baselines`, the prefix of the owner's checkout path becomes `<repo>`. The git history is not rewritten. The owner accepts that old commits keep the path.
- Owner, 2026-10-16: create `SECURITY.md` (GitHub private vulnerability reporting, short policy) and a short `CONTRIBUTING.md`, both in English. No `CODE_OF_CONDUCT.md`, no issue or pull request templates.
- Owner, 2026-10-16: add a fourth Planned card, "Coordinator orchestrates, does not execute". The owner observed that the Coordinator does too much hands-on work. Desired: it orchestrates and only judges whether delegated work was done with skill or needs rework; execution and verification go to subagents; it keeps the writes to `MEMORY.md`, `EPOCHAL.md` and `RISKS.md`. Status planned, no PRD yet.
- Owner, 2026-10-16: GitHub description and topics are set outside the PR (done).
- Owner, 2026-10-16: validation cap reached with P9 open; owner accepted the validator's correction for P9 and P10 without another round.
- Owner, 2026-10-16: the build found R8 conflicting with checkPrivatePaths; the owner accepted the amendment of R8 and R22 and the validator's corrections P11-P14 without another round.
- Pending decisions: none.

## Appendix

### Discarded alternatives

| Alternative | Why not |
|---|---|
| Rewrite git history to remove the path | The owner accepts the old path. A rewrite breaks clones and tags, and the path is not a secret. |
| Keep the release-notes paragraph in the README | `CHANGELOG.md` is its source. The paragraph goes stale at each release. |
| Docs-only PR for the README | `AGENTS.md`: documentation enters inside the PR of the feature. |
| Add a switch for each hook now | It is a behavior change. It needs its own PRD. |
| Move `evals/` and `*.review.md` out of the tree | It is a large change with no gain for a reader. The README labels them instead. |

### Risks and pending decisions

- Risk: the pt and en READMEs or roadmaps drift apart · mitigation: R1 and R12 check the parity; `dh validate` checks the Delivered cards.
- Risk: the README states a count that goes stale · mitigation: R2 and R7 forbid counts outside the anchored lines that `checkDocumentedCounts` checks.
- Risk: the path replacement breaks a baseline JSON or a test that reads it · mitigation: R8 parses every file, and R19 runs `go test ./...`.
- Risk: the README grows back into a changelog · mitigation: R1 caps lines and R2 caps line length.
- Risk: the clean-clone proof differs from the tree · mitigation: R20 blocks the commit until the sums are equal.
- Pending decisions: none.

## Executed summary

- **Delivered:** the public repository has a newcomer README in pt-br and English (12 mirrored sections with install, update, uninstall, what the kit writes, how to turn things off, troubleshooting and the two optional companions), `SECURITY.md` and `CONTRIBUTING.md`, four Planned roadmap cards, no owner path in the tracked files, and release 0.30.1.
- **Rules:**
  - R1-R7 Honored: QA matrix (`qa-verifier`): 12 headings in order in both READMEs (106 lines each after the docs pass), no release paragraph or test count, both companion links marked optional, uninstall/writes/turn-off sections scoped as required, plano-produto only in the Docs map, contributor-only block, no unanchored count.
  - R8 Honored: the bracket-class grep over the whole tree prints nothing; 28 baselines hold `<repo>` and parse.
  - R9-R11 Honored: `SECURITY.md` names private vulnerability reporting; `CONTRIBUTING.md` covers the five items with no AI attribution phrase; no code of conduct or templates.
  - R12-R14 Honored: four Planned cards with the exact titles in both roadmaps; no `danlemos`; no `CLAUDE.md`; exactly three `RISKS.md` lines per roadmap.
  - R15-R16 Honored: 0.30.1 in `build.go`, `marketplace.json` version and ref; CHANGELOG entry and Delivered cards with no drift.
  - R17-R19 Honored: no edit in `.agents/`, `.skills/`, `.commands/`; `dist/` only from the build; `dh validate` 0 errors, `go test ./...` green, the second build changes only `generated_at`.
  - R20 Honored: clean-clone proof, 5 of 5 binaries with equal sha256.
  - R21 Honored: no history rewrite.
  - R22 Honored: each README command in Update, Uninstall and Turn things off has a source (repo docs or `claude plugin ... --help`); none was run against the installed plugin.
- **Tickets:** T-2401 (backend, harness-maintainer) approved. QA by `qa-verifier`: verified (minor notes settled by the docs pass or accepted). Code review by `cli:codex/gpt-6-luna`: approve, 0 findings. Adherence by `product-discovery`: every item adheres.
- **Docs:** `README.md`, `README.en.md`, `SECURITY.md`, `CONTRIBUTING.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `.claude-plugin/marketplace.json`, `internal/build/build.go`.
- **Out of scope:** history rewrite, translation of `AGENTS.md`, PRDs and ADRs, code of conduct and templates, moving evals or review files, kit behavior changes, a per-hook switch, and the implementation of the four Planned items.
