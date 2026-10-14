# PRD-021 · Lean kit: faster hot path, dashboard parse cache, English-only PRD template and a cleanup rule

| Field | Value |
|---|---|
| Status | delivered on 2026-10-13 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-13 / 2026-10-13 |
| Tickets | T-2101, T-2102 |
| Origin | Owner chat, 2026-10-13 (the owner delegated: "follow your recommendations as long as you follow my premises") |
| Depends on | PRD-020 (English PRD), PRD-019 (version and `ref` pattern) |
| Enables | nothing |

## 1. Problem

Facts measured on 2026-10-13 with Go pprof on the 0.27.0 tree, linux/amd64:

- The `dh` binary is 7.65 MB (flags `-s -w -trimpath`). Only `internal/dashboard` imports `net/http`. The binary has no external module.
- `dh snapshot event` runs on every hook event. It takes about 3.96 ms per call. An empty process takes about 1.07 ms. 200 calls take 0.793 s.
- The work of `dh snapshot event` is 62 µs and 46 allocations. Package init adds about 0.55 ms.
- The init cost is: `internal/kit` 0.24 ms (854 allocs), `internal/dashboard` 0.20 ms (599 allocs), `internal/harness` 0.11 ms (77 allocs).
- There are 34 package-level `regexp.MustCompile` calls in the three packages: kit 20, dashboard 12, harness 2. `internal/doctor` has 1 and `internal/snapshot` has none. `internal/kit` also compiles 3 regexps inside functions (`validate.go` lines 445, 459, 482).
- Dashboard `buildState` on the owner's real data takes 55.8 ms and 12.9 MB per call, with 11,573 allocs.
- 85% of that CPU is in `projectProgress`. 56% is regexp `FindStringSubmatch` in `internal/dashboard/projects.go`. 18% is `snapshot.ReadCapped`.
- Every refresh (cache TTL 1 s) re-reads and re-parses every ticket and PRD file, even when the file is unchanged.
- `dh validate` takes about 0.25 s and allocates 81 MB. 70 MB is `io.ReadAll` and `readText`, because each check re-reads the same files. This item is a follow-up candidate. It is out of scope.
- `templates/pt-br/PRD.md` is a byte-identical copy of the English PRD template (PRD-020). It adds a file, a mirror rule and a count claim for no reader.
- The mirror check `checkTemplateRefs` in `internal/kit/validate.go` (lines about 501-521) is the only blocker to delete that copy.
- No rule says that a change must remove the leftovers it makes obsolete.

## 2. Solution

- The dashboard stays in the single `dh` binary. The whole binary is analysed with pprof to improve and optimise it, never to remove a part.
- Package-level regexps in `internal/kit`, `internal/dashboard` and `internal/harness` compile on first use. `dh snapshot event` then pays no init for them.
- The dashboard keeps a parse cache keyed by file path, size and modification time. Unchanged ticket and PRD files are not re-read on each refresh.
- `templates/pt-br/PRD.md` is deleted. `checkTemplateRefs` has one list of English-only templates with the entry `PRD.md`. Every claim that depends on the mirror is corrected.
- `.agents/coordinator.md` "Model selection" gets a short ruler for model and effort.
- `AGENTS.md` gets the cleanup premise: every change removes the leftovers it makes obsolete. Leftover junk blocks delivery.
- `AGENTS.md` states that all 5 binaries are generated and committed on each release. It states that pprof optimises the binary and does not trim it.
- Version 0.28.0.

## 3. Rules

- R1 `templates/pt-br/PRD.md` is deleted. `checkTemplateRefs` has one list of English-only templates with the single entry `PRD.md`.
  - A template in that list may exist only in `templates/en/`. The same name in another language folder is an error.
  - Any other template missing in one language stays an error.

  Verifiable: `ls templates/pt-br` shows no `PRD.md`; Go tests cover three cases (PRD only in en passes; PRD also in pt-br fails; TASK only in en fails); `go test ./internal/kit` is green.
- R2 Every claim that depends on the mirror is correct.
  - `harness-manifest.json` lists 13 templates after `go run ./cmd/dh build`. `checkDocumentedCounts` reads the count from disk. `internal/kit/doc_counts.go` gets no count change. Only its two regexps compile on first use (R3).
  - `docs/roadmap.html` and `docs/en/roadmap.html` say 13 templates (line about 685) and "7 templates in `en/`, 6 in `pt-br/`" (line about 678).
  - The first bullet of `AGENTS.md`, the templates tutorials (html and md, pt and en), the README lines on templates and `docs/dev-harness-mapa.excalidraw` are correct.
  - No leftover remains.

  Verifiable: `grep -rniE "espelhad|mirrored|cópia idêntica|identical copy" README.md README.en.md AGENTS.md docs --include=*.html --include=*.md --include=*.excalidraw` has only hits that are history (`docs/adr`, `docs/prd`, `CHANGELOG.md`, roadmap release cards) or that name only other templates; `dh validate` is green; `go run ./cmd/dh build` leaves no `pt-br/PRD.md` in `dist/`.
- R3 Package-level `regexp.MustCompile` in `internal/kit`, `internal/dashboard` and `internal/harness` compile on first use, for example with `sync.OnceValue`. Verifiable: `GODEBUG=inittrace=1 dh snapshot event < /dev/null 2>&1 | grep -E "internal/(kit|dashboard|harness)"` shows each of the three below 0.05 ms and below 100 allocs; the 200-call loop is not slower than 0.793 s on the same machine (report before and after); all existing tests are green.
- R4 The dashboard keeps a parse cache keyed by file path, size and modification time.
  - An unchanged ticket or PRD file is not re-read or re-parsed on a refresh.
  - A changed file is re-parsed.
  - The cache lookup runs after the same `os.Lstat`, regular-file and size checks that `snapshot.ReadCapped` runs. A refused file gives the same "skipped" warning with and without the cache.

  Verifiable: one Go test shows that a changed file is re-parsed, an unchanged file comes from the cache, and a ticket replaced by a symlink is refused with the same warning; `buildState` returns the same result with and without the cache on the same data; a scratch benchmark (not committed) shows at least 3x lower ns/op on repeated calls with unchanged data; the PR body reports before and after.
- R5 The PR body and `CHANGELOG.md` give the measured before and after of R3 and R4. They name the follow-up candidate: a shared read cache for `validate`. No new committed benchmark and no new eval exist. Verifiable: `grep -n "validate" CHANGELOG.md` finds the follow-up in the 0.28.0 entry; `git diff main...HEAD -- '*.go' | grep -c '^+func Benchmark'` is 0; `git diff --name-only main...HEAD -- evals/` is empty.
- R6 `.agents/coordinator.md` "Model selection" gets at most 10 new lines. The lines state these rules:
  - (a) Raise the model or the effort when ambiguity, coupling or consequence is high. The default is the cheap choice.
  - (b) Change the effort on the current model before you change the model.
  - (c) Do not delegate a dependent chain or work that fits one context.
  - (d) Do not use Haiku for security recon.
  - (e) The dispatch states the effort when it is raised. The log records it.

  The lines hold no benchmark number, no model release number, no machine slug and no absolute path. Verifiable: `grep -n -i "effort" .agents/coordinator.md` finds the lines; `git diff --stat` shows at most 10 added lines in that file; `dh validate` is green.
- R7 `AGENTS.md` gets three edits. It stays in Portuguese.
  - One bullet with the cleanup premise.
  - One bullet or edit states three things. The 5 binaries are generated and committed on each release. The dashboard stays. pprof optimises the binary and does not trim it.
  - The first bullet drops the identical-copy sentence for the PRD template and says that the PRD template exists only in `templates/en/`.

  Verifiable: `grep -n "lixo\|limpeza" AGENTS.md` finds the premise; `grep -n "pprof" AGENTS.md` finds the binary bullet; `grep -n "cópia idêntica" AGENTS.md` is empty.
- R8 Version 0.28.0. `defaultVersion` in `internal/build/build.go`, `version` in `.claude-plugin/marketplace.json` and the source `ref` (`v0.28.0`) in the same file change together. A forgotten `ref` is a delivery failure. The clean-clone proof applies, because the PR touches `internal/` and `dist/*/bin`. `CHANGELOG.md` has a 0.28.0 entry in Portuguese. Both roadmaps reflect it. Verifiable: `grep` of the three values shows `0.28.0` or `v0.28.0`; `git clone . <tmp>`, `go run ./cmd/dh build` there, and equal `sha256sum` of the binaries; `dh validate` is green.
- R9 Behavior is preserved. Verifiable: `go test ./...` is green; `claude plugin test` is green in `dist/claude-code/dev-harness`; the `dh validate` output for the repository is the same as before, except for the template count.
- R10 Scope guard. The dashboard routes and JSON output do not change. No new dependency exists. `go.mod` does not change. The other templates do not change. The old PRDs are not rewritten. Verifiable: `git diff go.mod` is empty; `git diff --name-only main...HEAD` lists no other template and no old PRD (`PRD-020` and its review file are history).

## 4. Docs

Checked with `grep` on 2026-10-13. Only files that carry an affected claim are listed.

- `AGENTS.md` (first bullet; cleanup premise; binaries and pprof bullet)
- `README.md`, `README.en.md` (line 63, "versões espelhadas" / "mirrored versions", and the version line)
- `docs/tutoriais/templates.html`, `docs/tutoriais/templates.md`, `docs/en/tutorials/templates.html`, `docs/en/tutorials/templates.md` (the mirror claim at the start; the PRD section)
- `docs/roadmap.html`, `docs/en/roadmap.html` (lines about 678 and 685; 0.28.0 entry)
- `docs/dev-harness-mapa.excalidraw` (note about the mirrored templates, line about 37)
- `CHANGELOG.md` (0.28.0)
- `.claude-plugin/marketplace.json` (`version` and `ref`)
- `harness-manifest.json` lists 13 templates after `go run ./cmd/dh build`; do not edit it by hand. `internal/kit/doc_counts.go` gets no count change. Only its two regexps compile on first use (R3).

`docs/tutorial.html` and `docs/en/tutorial.html` carry no mirror claim. They are not listed.

## 5. Out of scope

- A dashboard split, UPX, an own HTTP server, a history rewrite, and any removal of a binary or platform.
- Download of binaries on demand. All 5 binaries stay committed (ADR-001).
- A shared read cache for `dh validate`. It is a follow-up candidate.
- A committed benchmark and a new eval.
- The local plugin cache. The owner keeps only the latest version there, outside the repository.
- A change to other templates and a rewrite of old PRDs.

## 6. Owner decisions

- Owner, 2026-10-13: the dashboard stays in the single binary. pprof is for improving and optimising the whole binary, never for removing a part.
- Owner, 2026-10-13: keep only the latest plugin version in the local cache. This is done outside the repository and is not a rule here.
- Owner, 2026-10-13: all 5 platform binaries are generated every release and stay committed.
- Owner, 2026-10-13: a model-and-effort ruler goes into `.agents/coordinator.md`.
- Owner, 2026-10-13: remove `templates/pt-br/PRD.md` and fix every place that depends on it. The PRD is always English. This replaces the PRD-020 decision to keep the identical copy.
- Owner premise, 2026-10-13: every change removes the leftovers it makes obsolete. Leftover junk blocks delivery.
- Coordinator decision, 2026-10-13, under the owner's delegation: the `dh validate` read cache is a follow-up, out of scope.
- Coordinator decision, 2026-10-13, under the owner's delegation: the content of the R6 ruler, items (a)-(e).
- Analyst assumptions, not owner decisions:
  - (a) The lazy-init target is the 34 package-level regexps in the three packages of R3. `internal/doctor` keeps its 1. The 3 regexps inside functions in `validate.go` stay, because they run only in `validate`.
  - (b) The cache of R4 lives in memory in the dashboard process. It drops entries for paths not seen in the last build.
  - (c) The version is 0.28.0 because the change adds behavior with no break.
  - (d) `harness-manifest.json` changes only through `dh build`, not by hand. `internal/kit/doc_counts.go` gets no count change. Only its two regexps compile on first use (R3).
- Pending decisions: none.

## Executed summary

- **Delivered:** the PRD template exists only in `templates/en/` with a check in `dh validate`; package-level regexps compile on first use; the dashboard keeps a parse cache; the Coordinator has a model-and-effort ruler; `AGENTS.md` has the cleanup premise and the binaries rule; version 0.28.0.
- **Rules:**
  - R1 Honored: `templates/pt-br/PRD.md` is deleted; `englishOnlyTemplates` in `internal/kit/validate.go`; three Go test cases; mutation check red when the list is empty.
  - R2 Honored: manifest lists 13 templates; roadmaps say 13 and "7 in en/, 6 in pt-br/"; README, tutorials, excalidraw note and AGENTS.md first bullet corrected; mirror grep shows only history or other templates; `dist/` has no `pt-br/PRD.md`.
  - R3 Honored: inittrace kit 0.01-0.02 ms/62 allocs, dashboard 0.01 ms/51, harness 0 ms/6; 200-call loop 0.66-0.68 s against 0.75-0.80 s on `main`.
  - R4 Honored: cache keyed by path, size and mtime after the `Lstat`, regular-file and size guards; five Go tests; `buildState` 50-58 ms to 14-15 ms (about 3.8x), 11.5 MB to 0.67 MB.
  - R5 Honored: CHANGELOG 0.28.0 carries the before/after table and the `validate` read-cache follow-up; no committed benchmark; no change in `evals/`.
  - R6 Honored: four lines in `.agents/coordinator.md` "Model selection", no numbers, slugs or paths.
  - R7 Honored: cleanup premise at `AGENTS.md` line 12, binaries and pprof sentence at line 16, "cópia idêntica" gone.
  - R8 Honored: `defaultVersion`, marketplace `version` and `ref` are 0.28.0 and v0.28.0; CHANGELOG and both roadmaps; clean-clone proof run at commit time.
  - R9 Honored: `go test ./... -count=1`, `-race` on dashboard and kit, `claude plugin test` (254 pass) and `dh validate` green.
  - R10 Honored: `go.mod`, routes and JSON unchanged; no other template touched; no old PRD rewritten.
- **Tickets:** T-2101 (backend, builder) approved; T-2102 (backend, harness-maintainer) approved. QA matrix, code review by `claude` (approve), final review FEATURE APROVADA.
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutoriais/templates.html` and `.md`, `docs/en/tutorials/templates.html` and `.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/dev-harness-mapa.excalidraw`, `CHANGELOG.md`, `.claude-plugin/marketplace.json`.
- **Out of scope:** dashboard split, binary download, platform cuts, UPX, history rewrite, a shared read cache for `dh validate` (follow-up), new evals.
