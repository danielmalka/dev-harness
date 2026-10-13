# Review: PRD-020

| Field | Value |
|---|---|
| Document | docs/prd/PRD-020-prd-ste100.md (+ ~/.harness/projects/dev-harness/tasks/T-2001/TASK.md) |
| Source | Owner decisions 2026-10-13 (Coordinator dispatch, points 1-7) |
| Round | 2 of 2 |
| Verdict | changes required |
| Updated | 2026-10-13 |

## Findings

- **P1** [pending] [conflict] Deleting `templates/pt-br/PRD.md` fails the mirror check in `checkTemplateRefs` (validate.go:501-521). A Go change is then needed, which contradicts R10 and source points 2 and 7.
  - Location: §2 line 30, R4, R8, R10, §6 (c); T-2001 §3.2, R10, §8
  - Evidence: validate.go:514-519; validate_test.go:39
  - Suggestion: keep pt-br/PRD.md as a byte-identical copy of the English template (count stays 14), or get an owner decision for a PRD exception in `validate.go` and amend R10 and §5
  - Reported by: claude
- **P2** [pending] [conflict] `dh config set language` writes `<home>/config.yaml` (the global default), not `project.yaml`.
  - Location: §2 line 33, R7, §6 (d); T-2001 §3.4
  - Evidence: internal/harness/config_set.go:13-48; docs/tutorial.html:772
  - Suggestion: the Coordinator proposes an edit of `language` in `<harness dir>/project.yaml` and writes it after the owner agrees
  - Reported by: claude
- **P3** [pending] [gap] `.commands/document.md:18` (`templates/<lang>/PRD.md`) and `.agents/docs-guide.md:59` (PRD written in `language`) are not in R6 or in T-2001's authorization.
  - Location: R6; T-2001 Autorização
  - Evidence: R6's grep over `.commands` fails; source point 1
  - Suggestion: add both files to R6 and to the authorization
  - Reported by: claude
- **P4** [pending] [gap] The Docs list misses the roadmaps' line 678 ("7 templates espelhados por idioma"), the tutorials' line 740 and README.md/README.en.md line 69.
  - Location: §4; T-2001 §5
  - Evidence: source point 6; doc_counts.go:65-87 (a one-noun line is not checked)
  - Suggestion: list these lines explicitly
  - Reported by: claude
- **P5** [pending] [ambiguity] The §6 line 79 sentence has about 70 words, above R11's 25-word limit (style only, non-blocking).
  - Location: §6 line 79
  - Evidence: R11
  - Suggestion: one list item per assumption
  - Reported by: claude
- **P6** [pending] [weak criterion] The R7 grep does not cover `.commands/setup.md` (non-blocking).
  - Location: R7
  - Evidence: setup.md:22
  - Suggestion: add setup.md to the grep
  - Reported by: claude

## Not raised

- Source points 1-7 each map to §2, R1-R14 and §5. No excess was found.
- The named lines exist: project-onboarding:66, coordinator.md Language (73-74) and Executed summary (70), product-discovery:53, requirements-discovery:98.
- The dashboard already parses `delivered on` (R9). The template count is 14 today, and doc_counts checks only anchored multi-noun lines (R8). No Go code names `PRD.md`.
- R1-R14 IDs are unique. T-2001's lane and role match AGENTS.md.

## Coordinator disposition (round 1)

- P1: Coordinator picked the no-Go path: `templates/pt-br/PRD.md` stays as a byte-identical copy of `templates/en/PRD.md`. Spot-checked `checkTemplateRefs` in `internal/kit/validate.go:501-521`.
- P2-P6: returned to the author.

## Round 2 (2026-10-13)

- Review status: changes required (1 blocking, new)
- P1 [applied] pt-br PRD template stays a byte-identical copy (§2, R4, R8, R10, §6). No Go change.
- P2 [applied] Coordinator edits `<harness dir>/project.yaml`; `dh config set` is only the machine default. Leftover text in R7 filed as P7.
- P3 [applied] docs-guide.md:59 and document.md:18 added to R6 and to the T-2001 authorization.
- P4 [applied] README lines 69-70, tutorial line 740 and roadmap line 678 added to §4 and T-2001 §5.
- P5 [applied] §6 assumptions split into items.
- P6 [applied] `setup.md` added to the R7 grep.
- **P7** [pending] [conflict] R7 still said the Coordinator "proposes a `project.yaml` update", contradicting "no second confirmation" in §2, §6 and T-2001 §3.4.
  - Location: §3 R7 (line 44)
  - Suggestion: "on an explicit owner request to switch, edits `language` in `<harness dir>/project.yaml` with no second confirmation"; add `grep -n "switch only after agreement" .agents/coordinator.md` is empty.
  - Reported by: claude

## Coordinator disposition (round 2, cap reached)

- P7: the owner delegated the batch and is absent. The fix only aligns R7 with §2 and §6, which the owner-delegated decisions already set. The author applied the validator's exact suggestion. No third validation round ran. Flagged to the owner in the batch report for later review.

## Amendment revalidation (2026-10-13, one round)

Amendment: R10 allows `internal/kit/status_drift.go` and `adapters/claude-code/plugin/hooks/status-drift.ts` (plus tests) to accept `delivered on`; R6 adds "Solution" and the AGENTS.md language bullet; R1/R2/R6/R7 split for R11. Found by code review (T-2001 round 1) and QA.

- Review status: changes required (1 blocking)
- **P7** [applied] [conflict] R7 now says edit with no second confirmation, verified by the `switch only after agreement` grep.
- **P8** [pending] [conflict] §2 line 30 ("no Go change is needed") and §6 line 107 ("No new code") contradict the amended R10. §6 has no record of the status-drift amendment decision.
  - Location: §2 line 30, §6 line 107
  - Evidence: R10 lines 69-74; status_drift.go:82; status-drift.ts:35
  - Suggestion: scope line 30 to the template, add a §2 bullet for both forms, and add a §6 Coordinator decision line
  - Reported by: claude

## Coordinator disposition (amendment round)

- P8: owner absent, batch delegated. The author applied the validator's exact text. No further round ran. Flagged to the owner in the batch report.
