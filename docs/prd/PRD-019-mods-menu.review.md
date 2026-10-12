# Review — PRD-019 / T-1901

## Round 1 (2026-10-12)

- Review status: request changes
- Reviewers: claude (document-validator, opus). Default proportional review: 1 reviewer, first entry of `reviewers.document`.

## Findings

### BLOCKING: P1 THIRD_PARTY_NOTICES.md is an unrequested new doc file credited to the owner
- Reported by: claude
- Location: PRD R1, section 4, section 6; T-1901 R1, section 5, Autorização.
- Resolution: owner decision 2026-10-12 — no notices file and no origin comment in code. Remove from R1, docs, write set, and section 6.

### BLOCKING: P2 R5 "byte a byte" vs the original's JSON.parse/stringify writer
- Reported by: claude
- Resolution (coordinator default, recorded as assumption): semantic equality — other keys parsed and deep-equal, formatting normalized to 2-space JSON, as the original does.

### BLOCKING: P3 Mode's model and effort controls have no acceptance
- Reported by: claude
- Resolution: add criterion — Mode page lists models, effort levels and new-session modes; effort applied on `turn.step`; model switched through `/model`; proven by porting the two engine tests of the original.

### MINOR: P4 Backup file name and location not stated
- Reported by: claude
- Resolution: `~/.claude/settings.json.dh.bak`, written once if absent; restore source is the store key, as in the original.

### MINOR: P5 Segment text in R4 reads like translated output
- Reported by: claude
- Resolution: quote the current `dhProgress` text untranslated (`… (k blocked)`).

### MINOR: P6 R5 asks atomic write beyond the original
- Reported by: claude
- Resolution: drop "atômica"; keep "invalid/unreadable file is an error and never overwritten" (boundary validation).

## Round 2 (2026-10-12)

- Review status: request changes
- Reviewers: claude (document-validator, opus). Round cap reached (2 of 2); open points go to the owner.
- Settled: P1, P2, P4, P5, P6 resolved; P3 model part resolved, effort part open.

## Findings

### BLOCKING: P3 (carried over) Effort applied on `turn.step` has no test that can fail
- Reported by: claude
- Suggestion: R9 adds "after choosing `effort-high`, a `turn.step` with `effort` set goes to `next` with `effort: 'high'`".

### BLOCKING: P7 R1's "no notices file" check cannot fail
- Reported by: claude
- Suggestion: replace `git diff --stat THIRD_PARTY_NOTICES.md` with `test ! -e THIRD_PARTY_NOTICES.md`.

### MINOR: P8 R3 still points to R1's origin comments
- Reported by: claude
- Suggestion: `grep -rn "extra-mods\|/tools" adapters/claude-code/plugin/hooks` empty.

### MINOR: P9 Whether Mode writes the `.dh.bak` backup is unclear
- Reported by: claude
- Suggestion: add "(the same `.dh.bak`, if absent)" after the Mode clause.

## Owner disposition (2026-10-12)

- Owner accepted applying the round-2 suggestions (P3, P7, P8, P9) verbatim without a third round. Execution of T-1901 may start.
