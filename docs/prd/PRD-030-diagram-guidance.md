# PRD-030 · Diagram guidance for docs-guide and pinned Mermaid imports

| Field | Value |
|---|---|
| Status | delivered on 2026-10-28 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-28 / 2026-10-28 |
| Tickets | to be defined by the planner |
| Origin | Owner decision 2026-10-28 (option B + pin), task AUTO-3 |
| Depends on | PRD-029 (version base 0.34.0, delivered) |
| Delivery | One pull request, version 0.35.0 |

## 1. Problem

Three facts, one capability: how the kit guides diagrams in documents.

- Fact: the roadmap card "Skills de diagrama especializadas (Excalidraw e Mermaid)" (`docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`) is planned and has no PRD.
- Fact: `.agents/docs-guide.md` says nothing about diagrams. A docs-guide run has no rule on when to draw or in which format.
- Fact: four hand-written HTML pages import Mermaid with the floating tag `mermaid@11`. The kit rule fixes an exact version (`.skills/doc-template-html/scripts/stamp.sh:89` uses `11.17.2`). The roadmap card text (pt line 1736, en line 1737) still says Mermaid is the only CDN exception. `AGENTS.md` (line 29) says there are two: Mermaid and Google Fonts.

Decision (owner, 2026-10-28): the kit does not import the external `cathrynlavery/diagram-design` skill. Fact: it emits inline SVG, needs Python and Playwright, and has 45 diagram types. This conflicts with the `doc-template-html` rule "Never embed pre-rendered SVG" and with ADR-003 (lean kit). The owner may install it outside the kit.

## 2. Solution

After delivery, `docs-guide` has one short rule on diagrams. Mermaid is the default format, written in `<pre class="mermaid">`. Excalidraw is allowed only as a separate `.excalidraw` file that the document links. The four pages import one exact Mermaid version. The roadmap shows the card as delivered in 0.35.0. No skill, agent or command is added.

## 3. Rules

- R1 `.agents/docs-guide.md` carries one short instruction. It is one sentence added at the end of Procedure step 2 (same list item), so no step is renumbered. It says: use a Mermaid diagram (`<pre class="mermaid">`, as `doc-template-html` defines) where it helps the reader follow a flow or a structure. Use an Excalidraw diagram only as a separate `.excalidraw` file that the document links. Never embed it. No script, no export (owner decision 2026-09-25). No diagram skill is added.
  - When docs-guide reads the Procedure, then the sentence ends step 2 and the step numbers are unchanged.
  - When a document needs a flow or structure picture, then the instruction names Mermaid first.
  - When Excalidraw is used, then the instruction allows a link to a separate file only.
- R2 The four imports use the exact version `mermaid@11.17.2` with the file name `mermaid.esm.min.mjs`. The files are `docs/roadmap.html:1897`, `docs/en/roadmap.html:1901`, `docs/tutoriais/templates.html:983` and `docs/en/tutorials/templates.html:983`. The version is the same as in `stamp.sh:89`. The URL answers HTTP 200.
  - When one of the four files is read, then its import names `mermaid@11.17.2/dist/mermaid.esm.min.mjs`.
  - When the URL is requested, then the answer is HTTP 200.
  - When the four versions are compared with `stamp.sh:89`, then they are equal.
  - Accepted limit (decision): these imports have no Subresource Integrity. A static ES `import` cannot carry it. The pin to an exact version is the mitigation.
- R3 No `mermaid@11/` (floating) reference remains in tracked files, except under `dist/` and the history in `docs/prd/`.
  - When `grep -rn "mermaid@11/"` runs on the tree without `dist/` and `docs/prd/`, then it returns no line.
  - When `go run ./cmd/dh build` runs, then `dist/` holds no floating tag that the sources no longer hold.
  - When a new page imports Mermaid, then it uses the exact version of R2.
- R4 The roadmap card, in pt and en, is kept in sync by hand (`dh validate` ignores detail cards). Its status is delivered in version 0.35.0. Its description states what was delivered. It names the docs-guide instruction, Excalidraw as a linked separate file, no new skill, and the external `diagram-design` left to the owner. The words "única exceção de CDN" (pt 1736) and "the only allowed CDN exception" (en 1737) are replaced. The new text uses the two-exception wording of `AGENTS.md:29` (Mermaid and Google Fonts).
  - When the card is read in pt and en, then both show delivered, 0.35.0 and the same facts.
  - When the card is searched for the old exception words, then none remain.
  - When the card is read, then it names two CDN exceptions.
- R5 Version 0.35.0 is set in `internal/build/build.go:20`, in `.claude-plugin/marketplace.json` (version and ref `v0.35.0`) and in `harness-manifest.json`. `CHANGELOG.md` has an entry dated 2026-10-28. `docs/roadmap.html` and `docs/en/roadmap.html` have updated chips and a new Entregue / Delivered version card. The card is built like the 0.34.0 card (lines 703-709). Its Sources range goes up to 0.35.0. `dist/` is regenerated only by `go run ./cmd/dh build`, never by hand.
  - When the files named in this rule are read, then each version field says 0.35.0 (the marketplace ref says v0.35.0).
  - When the roadmap is read in pt and en, then the new card exists with the same structure as 0.34.0.
  - When `dist/` is compared with the build output, then there is no difference.
- R6 The checks pass. Counts stay at 13 agents, 19 commands plus 1 mod command, and 22 skills.
  - When `dh validate` runs, then it passes.
  - When `go test ./...` runs, then it passes.
  - When `go run ./cmd/dh build` runs after the regeneration, then `dist/` shows no difference and the counts are unchanged.

Out of scope: any new skill or agent; SVG generation; Excalidraw export or embed; SRI for ESM imports; a change to the history text of ADR-003 or PRD-017; a drift-report offer; any change to `stamp.sh` or `internal/dashboard/pdocs.go`.

## 4. Docs

- `.agents/docs-guide.md`: the documented flow touched (docs-guide behavior, R1).
- `docs/roadmap.html` and `docs/en/roadmap.html`: card and version card (R4, R5).
- `CHANGELOG.md`: entry for 0.35.0 (R5).
- `docs/tutoriais/templates.html` and `docs/en/tutorials/templates.html`: import pin only (R2).

The `doc-template-html` diagram rule is unchanged. The tutorials and the README need no new text beyond the pinned imports of R2.

## Executed summary

- **Delivered:** The docs-guide now carries one diagram instruction (Mermaid first; Excalidraw only as a separate linked `.excalidraw` file, never embedded), the four hand-written Mermaid imports are pinned to `mermaid@11.17.2`, the roadmap card is delivered and corrected, and the kit is at 0.35.0; no new skill or agent was added.
- **Rules:**

  | Rule | Verdict | Evidence |
  |---|---|---|
  | R1 | Honored | `.agents/docs-guide.md`: the sentence sits at the end of Procedure step 2, in the same list item; the `## Limits` line still cites Procedure step 5. Checked in `tasks/T-3002/evidence/` (harness dir). |
  | R2 | Honored | `docs/roadmap.html:1895`, `docs/en/roadmap.html:1898`, `docs/tutoriais/templates.html:983`, `docs/en/tutorials/templates.html:983` import `mermaid@11.17.2/dist/mermaid.esm.min.mjs`, the same version as `.skills/doc-template-html/scripts/stamp.sh:89`. QA observed HTTP 200 for the URL. Accepted limit: a static ES import has no Subresource Integrity; the exact pin is the mitigation. |
  | R3 | Honored | `git grep 'mermaid@11/'` outside `dist/` and `docs/prd/` has no match (checked again at documentation time). |
  | R4 | Honored | The roadmap card is delivered in pt and en with the four facts and the two CDN exceptions (Mermaid and Google Fonts). The old exception wording is gone from both roadmaps. |
  | R5 | Honored | Version 0.35.0 in `internal/build/build.go:20`, `.claude-plugin/marketplace.json` (`version` 0.35.0, `ref` v0.35.0), `harness-manifest.json`, the `CHANGELOG.md` entry dated 2026-10-28, the chips, the Sources range and a new Delivered version card in pt and en. A rebuild in a copy reproduces `dist/` byte for byte except the `generated_at` timestamp line of `harness-manifest.json`, which is expected build output (`tasks/T-3002/evidence/host-checks/repro-diff.txt`, harness dir). |
  | R6 | Honored | `go test ./...` exit 0 and `dh validate .` exit 0 (errors 0, 22 skills), `tasks/T-3002/evidence/host-checks/` (harness dir). Counts: 13 agents, 19 commands plus 1 mod command, 22 skills. The same `dist/` note as R5 applies. |
- **Tickets:** T-3001 (backend, harness-maintainer): done. Code review round 1 by `cli:opencode/opencode-go/muse-spark-1.3-contributor`: `approve`, 0 findings. Confirming round 2 by `cli:codex/gpt-6-luna`: `request changes` with one finding F2-1 (Major, "R6 check evidence not observed in the reviewer's read-only sandbox"); the Coordinator classified it as noise and resolved it with host evidence on the same evaluated state. T-3002 (teste, qa-verifier): green, matrix R1 to R6. Evidence: `tasks/T-3001/evidence/review-r1b/`, `tasks/T-3001/evidence/review-r2/`, `tasks/T-3002/evidence/` (harness dir).
- **Docs:** `.agents/docs-guide.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/tutoriais/templates.html`, `docs/en/tutorials/templates.html`, `CHANGELOG.md`.
- **Out of scope:** any new skill or agent; SVG generation; Excalidraw export or embed; SRI for ESM imports; a change to the history text of ADR-003 or PRD-017; a drift-report offer; any change to `stamp.sh` or `internal/dashboard/pdocs.go`. Also: the delivered detail card stays in the roadmap section "Planejado" (noted by QA as a possible Minor, owner decision pending).
