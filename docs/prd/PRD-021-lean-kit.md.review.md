# Review: PRD-021

| Field | Value |
|---|---|
| Document | docs/prd/PRD-021-lean-kit.md (+ tasks T-2101 and T-2102 in the harness dir) |
| Source | Owner decisions 2026-10-13 and pprof measurements (Coordinator dispatch) |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-13 |

## Findings

- **P1** [pending] [conflict] T-2101 cannot finish green: deleting `templates/pt-br/PRD.md` makes the real-tree roadmap test (`TestRealTreeRoadmapsClean`) fail on `docs/roadmap.html:685` and `docs/en/roadmap.html:685` ("14 templates"), and those files are in T-2102's write set only.
  - Evidence: `internal/kit/doc_counts.go:88-126`, `internal/kit/roadmap_check_repo_test.go:24-34`
  - Suggestion: move the line 685 fix (pt and en) into T-2101's authorization, flow and "Pronto quando"
  - Reported by: claude
- **P2** [pending] [conflict] `doc_counts.go` holds no template count and `dh build` does not regenerate it; the count comes from disk and `harness-manifest.json`.
  - Location: R2, §4 last bullet, §6 (d), T-2102 authorization note
  - Suggestion: "`harness-manifest.json` lists 13 templates after `go run ./cmd/dh build`. `checkDocumentedCounts` reads the count from disk. `internal/kit/doc_counts.go` does not change."
  - Reported by: claude
- **P3** [pending] [conflict] The "37 regexps in kit, dashboard, harness and snapshot" fact is wrong: `internal/snapshot` has no regexp; package level is kit 20, dashboard 12, harness 2; `internal/doctor` has 1; kit has 3 more inside functions (`validate.go:445,459,482`).
  - Location: §1, §6 (a)
  - Reported by: claude
- **P4** [pending] [gap] `docs/dev-harness-mapa.excalidraw:37` carries the mirror claim and is in no list or write set.
  - Suggestion: add it to §4, T-2102 §5 and its authorization
  - Reported by: claude
- **P5** [pending] [weak criterion] R2's `grep "pt-br/PRD"` cannot see the real mirror claims ("espelhadas", "mirrored", "cópia idêntica", "identical copy", roadmap line 678).
  - Suggestion: grep for the mirror phrases; line 678 of both roadmaps states 7 templates in `en/` and 6 in `pt-br/`
  - Reported by: claude
- **P6** [pending] [ambiguity] R4 does not say whether the cache keeps the read guards of `snapshot.ReadCapped` (no symlink, regular file, size cap, same "skipped" warning).
  - Suggestion: the lookup runs after the same `os.Lstat`, regular-file and size checks (as `scanEpochal`/`epochalCache` do); a refused file gives the same warning with and without the cache; the Go test covers a ticket replaced by a symlink
  - Reported by: claude
- **P7** [pending] [weak criterion] R5's check (`*_bench*` file names) cannot see a Go benchmark.
  - Suggestion: `git diff main...HEAD -- '*.go' | grep -c '^+func Benchmark'` is 0; `git diff --name-only main...HEAD -- evals/` is empty
  - Reported by: claude
- **P8** [pending] [excess] (non-blocking) R6 items (a)-(e) and "`dist/` holds only the current version" are not traced in §6; R7 bullet 2 has no check.
  - Reported by: claude
- **P9** [pending] [ambiguity] (non-blocking) The dashboard cache keeps entries for deleted or renamed paths.
  - Suggestion: drop entries for paths not seen in the last build
  - Reported by: claude

## Coordinator disposition (round 1)

- Spot-checked P1 (`roadmap_check_repo_test.go`, roadmap line 685 says "14 templates") and P3 (`internal/snapshot` has no regexp).
- All nine findings returned to the author with the validator's suggestions accepted. P8: the "`dist/` holds only the current version" sentence is dropped; R6 content is recorded in §6 as a Coordinator decision under the owner's delegation.

## Round 2 (2026-10-13)

- Review status: approved (0 blocking). P1-P9 applied. Three new non-blocking findings.
- P1 [applied] T-2101 owns the roadmap template-count lines (about 678 and 685); T-2102 adds only the 0.28.0 card after it.
- P2 [applied] R2, §4 and §6 (d): the count comes from disk and the manifest; `doc_counts.go` does not change.
- P3 [applied] §1 and §6 (a): 34 package-level regexps (kit 20, dashboard 12, harness 2), doctor 1, snapshot none, 3 in-function calls in kit.
- P4 [applied] `docs/dev-harness-mapa.excalidraw` is in R2, §4 and T-2102.
- P5 [applied] R2 greps the mirror phrases.
- P6 [applied] R4 keeps the `ReadCapped` read guards; the Go test covers a symlink.
- P7 [applied] R5 checks `^+func Benchmark` and `evals/`.
- P8 [applied] R6 content recorded in §6 as a Coordinator decision; the `dist/` sentence is gone.
- P9 [applied] The cache drops entries for paths not seen in the last build.
- **P10** [applied by the Coordinator's dispatch] §4 AGENTS.md entry says "binaries and pprof bullet" (not a `dist/` bullet). Non-blocking.
- **P11** [applied by the Coordinator's dispatch] T-2101 §4 lists the roadmap-count part of R2; T-2102 §5 hands only the 0.28.0 card to T-2102. Non-blocking.
- **P12** [applied by the Coordinator's dispatch] R7 bullet 2 split into short sentences. Non-blocking.

## Coordinator disposition (round 2)

- Text-only fixes for P10-P12 sent to the author; no third validation round (cap reached, findings are non-blocking and mechanical).
