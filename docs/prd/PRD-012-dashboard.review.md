# PRD-012 · Validação

Documento: `docs/prd/PRD-012-dashboard.md`. Validador: `document-validator` (opus), skill `document-review`.

## Rodada 1 de 2 (2026-10-07)

- Review status: changes required
- Blocking findings: 6

## Findings

### BLOCKING: P1 [ambiguity] Avatar states "erro" and "concluído" have no source event; "concluído" overlaps "parado"
- Location: section 2 item 3; R13; R18
- Reported by: claude
- Problem: R18 ranks six states; R13 defines only working/idle/waiting/closed. No event produces `erro` (jevmon: `StopFailure`) or `concluído` (jevmon: `Stop`), while R13 maps `Stop` to `idle`.
- Suggested fix: extend R13 with an event→state table for every state R18 ranks (e.g. `StopFailure`→error, `Stop`→done decaying to idle after N min), or drop them; one mod test per mapping.

### BLOCKING: P2 [conflict] R13 renames `active` to `working` without saying what the mod writes, breaking the PRD-001 extension contract
- Location: R11; R13; C5
- Reported by: claude
- Problem: the extension knows active/idle/closed; R13 never says whether the mod writes `active` or `working`, nor how old readers treat `waiting`.
- Suggested fix: keep `state` in schema-1 vocabulary and add a new field (e.g. `activity`), or write new values and list the extension change in Docs.

### BLOCKING: P3 [gap] Removing `dh snapshot event` drops fields no rule tells the mod to produce
- Location: R11, R12, R13
- Reported by: claude
- Problem: today the hooks fill `session_name`, `started_at`, `cwd`, `agent`, `model.id` and the `events` ring buffer; without them R15 cannot map a session to a project.
- Suggested fix: the mod writes every field `dh snapshot event` writes today, from equivalent mod events; one mod test per event compared with `dh snapshot event`.

### BLOCKING: P4 [ambiguity] Owner's jevmon folder cannot work unchanged without a default state→pose map; sample file names undefined
- Location: R16, R17
- Reported by: claude
- Problem: jevmon poses are named `frente`, `duvida`, `puto`…; no default map without `poses.json`; fallback only works if both sets share pose names.
- Suggested fix: R16 adopts jevmon's default map; R17 names sample files with those pose names.

### BLOCKING: P5 [conflict] R21 settles pending decision 7 while the PRD lists it as open
- Location: R21; Appendix decision 7; C2
- Reported by: claude
- Problem: R21 requires option (b) of an open decision.
- Suggested fix: make R21 conditional on decision 7.

### BLOCKING: P6 [gap] The drift check cannot catch the staleness that motivates it
- Location: R19, R20; section 1; source D4
- Reported by: claude
- Problem: only "open ticket in a delivered PRD" is drift; "all tickets done / cited delivered in CHANGELOG but header not `entregue em`" passes silently.
- Suggested fix: add the reverse condition, or record the narrow check as an open decision.

### NON-BLOCKING: P7 [conflict, factual] `updated_at` already exists in schema 1
- Location: R15; section 1
- Reported by: claude
- Suggested fix: reword as refreshed on every mod event and by heartbeat if H3 holds.

### NON-BLOCKING: P8 [conflict, factual] Status counts and quotes in section 1 and C2 slightly off
- Location: section 1; C2
- Reported by: claude
- Suggested fix: "dezenas dizem 'pronta'"; "dizem algo diferente de 'entregue em'".

### NON-BLOCKING: P9 [ambiguity] "Delivered PRDs collapsed" became a count line only
- Location: section 2 item 1; R10
- Reported by: claude
- Suggested fix: state count only or expandable list.

### NON-BLOCKING: P10 [organization] Pending-decision bookkeeping
- Location: Appendix decisions 4 and 6; R23, R24
- Reported by: claude
- Suggested fix: rename decision 4 to "Animação no conjunto de exemplo"; move the risk half of R23 to Riscos.

## Rodada 2 de 2 (2026-10-07)

- Review status: changes required
- Blocking findings: 2 (P11, P12), both new; P1-P10 resolved. Cap reached: goes to the owner with P11 and P12 open.

## Findings

### BLOCKING: P11 [conflict] `done` decay measured from `updated_at`, which the heartbeat keeps refreshing, so `done` never decays under the recommended liveness option
- Location: R13; R15; decision 3 recommendation (a)
- Reported by: claude
- Suggested fix: the mod writes `activity_at` only when `activity` changes; heartbeat refreshes `updated_at` only; reader decays `done` when `now - activity_at >= N`; Go test with fresh `updated_at` and 6-min-old `activity_at` reads idle.

### BLOCKING: P12 [conflict] Option (a) of decision 7 makes R19(b) fail in this repo, so R21's check cannot pass and R20 blocks every commit
- Location: R21 option (a); R19(b); R20; decision 7
- Reported by: claude
- Suggested fix: state that consequence in 7(a) (viable only with an exemption list in R19), or drop (a) and keep (b) only.

### NON-BLOCKING: P13 [organization] Docs cites "decisão pendente 5" for the command count (it is 9); "3b" has no options
- Reported by: claude

### NON-BLOCKING: P14 [gap, hypothesis] `StopFailure` reaching the mod is not listed in H2
- Reported by: claude

### NON-BLOCKING: P15 [ambiguity] R19(b) fires when the last ticket closes, before final review/release; not stated as intended nor listed as risk
- Reported by: claude
