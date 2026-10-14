# PRD-022 · Remove the visual mods that the owner now keeps in a separate repository

| Field | Value |
|---|---|
| Status | delivered on 2026-10-14 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-14 / 2026-10-14 |
| Tickets | T-2201 |
| Origin | Owner chat, 2026-10-14 (scope decided in chat) |
| Depends on | PRD-019 (added the mods menu) |
| Enables | nothing |

## 1. Problem

- The owner now keeps personal Claude Code mods in a separate repository, `malka-claude-mods`. Its plugins are statusline, menu, cleanscreen, mode, permissions and sentinel.
- The `dh` mod in `adapters/claude-code/plugin/hooks/` draws the same things. Two status bars now conflict on the screen.
- `panel.ts` calls `$.ui.status('')` in `refresh()`. The engine pins an empty plugin status, because only `undefined` clears it. The prompt then shows an empty "⚠ dh:" notice.
- The kit does not need these visual parts to work. The memory lock, the commit gate and the dashboard snapshot do not use them.

## 2. Solution

Remove the visual components from the `dh` mod. Keep every guard and the dashboard.

Remove:

- The status line under the prompt (`lib/status-line.tsx`, `lib/status.ts`, the `PromptHint` render, the `dhSegments` and `statusDrawn` state) and the fallback `$.ui.status(...)` bar in `panel.ts` `refresh()`. Copy: statusline.
- The `/mods` menu (`mods-menu.tsx`, `lib/menu-pages.tsx`, `lib/settings.ts`). Copy: menu.
- The Zen band, the `zen_progress` tool and the `ZEN_PROMPT` injection (`mods-zen.tsx`, `lib/zen.ts`, `lib/zen-band.tsx`). Copy: cleanscreen.
- The model and effort mode (`lib/modes.ts`, the `turn.step` effort override). Copy: mode.
- The default permission mode (in `mods-menu.tsx`, `lib/settings.ts` and `lib/menu-pages.tsx`). Copy: permissions.
- The running-agents count (the `running` map and the `agents N` segment in `panel.ts`). Remove the map too, unless the snapshot uses it. Copy: sentinel.
- The context toasts at 60% and 80%. No copy exists. The owner approved the removal and sent the idea to the `malka-claude-mods` session.

Keep:

- `runtime-guard.ts` (memory lock, ADR-005).
- `snapshot-writer.ts` (dashboard live sessions).
- The commit gate (dirty-after-gate, AI mention deny).
- `status-drift.ts`.
- The `dh:harness` prompt section.
- The `/dashboard` command (`dashboard-command.ts`).

The dashboard still shows ticket progress per PRD (`internal/dashboard/projects.go`, `projectProgress`). That code does not change.

Version 0.29.0. The command count is 19 file commands and 1 mod command (`/dashboard`), not 19 and 2.

## 3. Rules

- R1 No visual component remains in `adapters/claude-code/plugin/hooks/`. There is no `PromptHint` render, no `ui.status` call, no `ui.toast` call, no `/mods` command and no `zen_progress` tool. Verifiable: `grep -rnE '\$\.ui\.(status|toast)|PromptHint|zen_progress|ZEN_PROMPT|/mods' adapters/claude-code/plugin/hooks --exclude='*.test.ts'` has no hit.
- R2 The kept guards work as before. Verifiable: the mod tests for `runtime-guard`, `snapshot-writer`, the commit gate and `status-drift` are green.
- R3 No orphan remains: no file, test, import, reference, doc or generated copy of a removed part (cleanup premise in `AGENTS.md`). Verifiable: `git ls-files adapters dist | grep -E 'mods-|lib/(status|status-line|zen|zen-band|menu-pages|modes|settings)\.'` is empty; `git grep -nE 'mods-menu|mods-zen|mods-types|menu-pages|status-line|zen-band|lib/(modes|settings|status|zen)' -- . ':!docs/adr' ':!docs/prd' ':!CHANGELOG.md'` has hits only in roadmap release cards.
- R4 The command count is 19 + 1 in every place that states it. Verifiable: `grep -nE '19 \+ 2|2 comandos registrados|2 commands registered' AGENTS.md README.md README.en.md docs/tutorial.html docs/en/tutorial.html` has no hit; `dh validate` is green.
- R5 The build is complete. Verifiable: `dh validate` and `go run ./cmd/dh build` pass with `dist/` regenerated; `go test ./...` and the mod tests pass.
- R6 The change does not touch a `statusLine` setting of the user. The owner confirmed that none was saved. Verifiable: `git diff --name-only main...HEAD` lists no user settings file.
- R7 Version 0.29.0. `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json` and the source `ref` (`v0.29.0`) change together. The clean-clone proof applies, because the PR touches `dist/*/bin`. Verifiable: `grep` of the three values shows `0.29.0` or `v0.29.0`; the `sha256sum` of the binaries is equal in the clean clone.

## 4. Docs

- `AGENTS.md` (mods list; "19 + 2" becomes "19 + 1")
- `README.md`, `README.en.md`
- `docs/tutorial.html`, `docs/en/tutorial.html`
- `docs/roadmap.html`, `docs/en/roadmap.html` (0.29.0 entry)
- `CHANGELOG.md` (0.29.0)
- `docs/prd/PRD-019-mods-menu.md`: header note "superseded in part by PRD-022".
- `docs/adr/ADR-006-painel-e-travas-de-commit.md`: header note "superseded in part by PRD-022".
- `.claude-plugin/marketplace.json` (`version` and `ref`)

## 5. Out of scope

- Any change to `runtime-guard.ts`, `snapshot-writer.ts`, the commit gate, `status-drift.ts`, the `dh:harness` section or `/dashboard`.
- Any change to the dashboard (`internal/dashboard`).
- The `malka-claude-mods` repository and its plugins.
- A replacement for the removed features inside the kit.

## 6. Owner decisions

- Owner, 2026-10-14: remove from the `dh` mod only what has a copy in `malka-claude-mods`, and remove the whole status line even without a `dh` segment there.
- Owner, 2026-10-14: remove the context toasts at 60% and 80%. The idea went to the `malka-claude-mods` session.
- Owner, 2026-10-14: the guards, the commit gate, `status-drift.ts`, the `dh:harness` section and `/dashboard` stay.
- Root cause, fact: `panel.ts` calls `$.ui.status('')`. The engine pins an empty plugin status. The empty "⚠ dh:" notice disappears with the removal.
- Pending decisions: none.

## Executed summary

- **Delivered:** the mod no longer draws a status line, a `/mods` menu, a Zen band, a Mode page or context toasts; the memory lock, snapshot recording, commit gate, status drift, the `dh:harness` section and `/dashboard` stay, and the empty "⚠ dh:" notice is gone.
- **Rules:**
  - R1 Honored: the R1 `grep` has no hit.
  - R2 Honored: `claude plugin test` in `dist/claude-code/dev-harness` is 226 pass, 0 fail, including a new test for the `<cwd>/.harness` fallback (a mutation that removed the fallback made it fail).
  - R3 Honored: the `git ls-files` check is empty; the `git grep` check has hits only in roadmap release cards.
  - R4 Honored: the `grep` has no hit (`/mods` removed from the pattern by owner decision P6, option A); `dh validate` green.
  - R5 Honored: `dh validate` 0 errors, `go run ./cmd/dh build` idempotent, `go test ./...` green.
  - R6 Honored: no user settings file is in the change.
  - R7 Honored: `0.29.0` in `internal/build/build.go` and `.claude-plugin/marketplace.json`, `ref` `v0.29.0`; clean-clone proof: 7 of 7 binaries byte-identical.
- **Tickets:** T-2201 (backend, harness-maintainer) approved. QA matrix by `qa-verifier`; code review by `claude` (2 non-blocking findings, fixed) and, on re-entry, by `cli:codex/gpt-5.6-luna` (approve, 0 findings).
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `docs/prd/PRD-019-mods-menu.md` (note), `docs/adr/ADR-006-painel-e-travas-de-commit.md` (note), `.claude-plugin/marketplace.json`.
- **Out of scope:** context toasts in the `statusline` mod of malka-claude-mods (idea sent to that project), any change to the dashboard, `/dashboard` or the kept guards.
