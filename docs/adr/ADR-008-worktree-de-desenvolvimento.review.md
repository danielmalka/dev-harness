# Review: ADR-008

## Round 1

| Field | Value |
|---|---|
| Document | docs/adr/ADR-008-worktree-de-desenvolvimento.md |
| Source | docs/prd/PRD-028-worktree-and-approval-page.md |
| Round | 1 of 2 |
| Verdict | changes required |
| Updated | 2026-10-28 |

Corrections in round 1: the ADR by the architect; PRD-028 rules R30, R61, R26 and R27 by product-discovery. The categories of P1 to P5 are assigned from the one-line descriptions in the Coordinator's record.

## Findings

- **P1** [applied] [conflict] (blocking) The ADR said `mode` is not validated on read, but the code validates it.
  - Location: ADR-008, rule on reading `mode`
  - Evidence: `cmd/dh/config.go` validates `mode`.
  - Resolution: ADR corrected by the architect, round 1.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) R30 and D4 listed trust paths that do not exist in a normal checkout.
  - Location: ADR-008 D4; PRD-028 R30
  - Resolution: ADR and PRD corrected, round 1.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] (blocking) The condition of clone step 1 was ambiguous between a `.git` directory and a gitfile.
  - Location: ADR-008, clone step 1
  - Resolution: ADR corrected by the architect, round 1.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [ambiguity] (blocking) `linked_worktree.gitdir` was undefined; it means the administrative directory `<C>/worktrees/<name>`.
  - Location: ADR-008, `linked_worktree` definition
  - Resolution: ADR corrected by the architect, round 1.
  - Reported by: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [ambiguity] (blocking) The ADR did not say that a relative back-pointer resolves against the administrative directory.
  - Location: ADR-008, back-pointer resolution
  - Resolution: ADR corrected by the architect, round 1.
  - Reported by: cli:agy/gemini-3.8-flash-high
- One non-blocking finding was raised in round 1. Its text is not in the Coordinator's record, so it is not reproduced here. Recorded as corrected in round 1.

## Not raised

- Round 1 checks recorded by the Coordinator: the freeze was clean and `dh validate` reported 0 errors. The reviewer ran as a qa-verifier on the live tree.

## Round 2

| Field | Value |
|---|---|
| Document | docs/adr/ADR-008-worktree-de-desenvolvimento.md |
| Source | docs/prd/PRD-028-worktree-and-approval-page.md |
| Round | 2 of 2 (cap reached) |
| Verdict | approved |
| Updated | 2026-10-28 |

## Findings

- P1 to P5 carried forward from round 1, all `applied`; closed in round 2.
- **P6** [applied] [gap] (non-blocking) R61 did not name the fields of `linked_worktree`.
  - Location: PRD-028 R61
  - Resolution: fixed by the Coordinator in one line, round 2.
  - Reported by: cli:agy/gemini-3.8-flash-high

## Not raised

- Round 2 checks recorded by the Coordinator: the freeze was clean, no agy or CodeGraph process was left running after the round, and `dh validate` reported 0 errors.
- No blocking finding is open.
- The ADR status was set to `aceito` after round 2 (Coordinator record).
