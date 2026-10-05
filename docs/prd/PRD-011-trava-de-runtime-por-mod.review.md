# Review: PRD-011

| Field | Value |
|---|---|
| Document | docs/prd/PRD-011-trava-de-runtime-por-mod.md |
| Source | Owner decision 2026-10-05 ("Pode 1 e 2, dentro do plugin") + Coordinator verified facts (dispatch); AGENTS.md; .harness/RISKS.md RISK-001 |
| Round | 2 of 2 |
| Verdict | approved |
| Updated | 2026-10-05 |

## Findings

- **P1** [applied] [conflict] `claude plugin test` was labelled a verified fact but was not in the source's verified facts; the real-session check was in no ticket.
  - Location: section 1 lines 22-24; T-1102 lines 100-103
  - Resolution: round 2. The fact is now scoped to the Coordinator's post-round-1 verification (`--help`, `claude-code/testing`, smoke test). T-1102 requires a real session with the built plugin: subagent `Agent`/`Write .harness/MEMORY.md` denied, main session allowed.
  - Coordinator note (2026-10-05): verified after round 1 — `claude plugin test` exists in 2.1.289 (`claude plugin test --help`; a smoke `*.test.ts` importing `claude-code/testing` ran 1 pass in a scratchpad probe).
  - Reported by: claude
- **P2** [applied] [weak criterion] T-1102 could close with R9 "não verificado" and bypass the release gate.
  - Location: R9 lines 72-75; line 130
  - Resolution: round 2. "Não verificado" or a plugin load failure → the release waits for an explicit owner decision recorded in the PRD. T-1102 tests 2.1.287 and 2.1.288, which were confirmed present.
  - Reported by: claude
- **P3** [applied] [conflict] The AGENTS.md runtime bullet would become false.
  - Location: R8(b) lines 67-69
  - Resolution: round 2. The runtime bullet cites the TS mod exception (ADR-004), checked by `grep -n ADR-004 AGENTS.md`.
  - Reported by: claude
- **P4** [applied] [weak criterion] R8 did not say which coordinator.md text changes; the grep named no string.
  - Location: R8(a) lines 65-67, 70-71
  - Resolution: round 2. Target is the "sole writer" bullet (coordinator.md line 37); the fenced block is unchanged in both files; check is `grep -n "enforced at runtime"` plus `checkAntiDelegationClause` green.
  - Reported by: claude
- **P5** [applied] [ambiguity] R7 left an optional `.ts` scan to the builder.
  - Location: R7 lines 62-64; T-1101
  - Resolution: round 2. No Go change; the limit is recorded; checked by `git diff --stat` showing no `cmd/`, `internal/` or `go.mod`.
  - Reported by: claude
- **P6** [applied] [excess] `NotebookEdit` and a conditional `MultiEdit`.
  - Location: R4 line 50
  - Resolution: round 2. `Write` and `Edit` only.
  - Reported by: claude
- **P7** [applied] [ambiguity] R4 path matching had no root anchor.
  - Location: R4 lines 50-55
  - Resolution: round 2. Suffix match on the normalized path (`/` and `\`), any root, with a test for an out-of-project path.
  - Reported by: claude
- **P8** [applied] [organization] The test ticket wrote docs; T-1102 had no Documentação field.
  - Location: R9; section 4 line 84; T-1102 line 104
  - Resolution: round 2. Docs at the end of the batch by the docs-guide; T-1102 "Documentação: nenhum".
  - Reported by: claude

## Not raised

- No scope beyond items 1 and 2; R10 exclusions intact; R2 allows exactly two handlers.
- R3/R4 still match the verified `agent.spawn`/`parentAgentId` and `tool.call`/`agentId` facts; R5 and R6 unchanged.
- Line 128 decisions match the Coordinator resolutions given to the author.
- Observation (non-blocking): the R8(b) dispatch/memory clause in AGENTS.md has no grep of its own; missing it leaves AGENTS.md true, not false.
- Observation: the `claude --plugin-dir` flag was not verified by the validator; the T-1102 acceptance holds whatever load method is used.
- All six categories re-checked on the revised text; identifiers unique.
