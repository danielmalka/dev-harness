# ADR-002 — review

Document: `docs/adr/ADR-002-cli-planejadora.md`. Source: PRD-005 (RF-03/04/05/13, AC-05/06/15, §6), GRILL G6/G15/G18–G23, PLAN-011 slice 1 and T-801, RISK-001, CLAUDE.md (AA1; exception of 2026-09-28), `.skills/external-clis/SKILL.md`, `.agents/coordinator.md`, `internal/kit/evals.go`.
Depth: lean rule — one Claude validator (`document-validator`, opus requested, effective unverified), two rounds.

## Round 1 (2026-09-28)
- Verdict: changes required — 1 blocking, 4 non-blocking.
- A1 — conflict — blocking — §5, line 55. The failure signal "the two clauses differ in anything beyond review → planning" fires on the correct text fixed in T-801 (lines 38-51), which says "this planning task" and wraps the first lines differently. Fix: the signal becomes "the planning clause differs from the text fixed in T-801 / PLAN-011 slice 1; compared with each other, they differ in anything beyond 'this review' → 'this planning task' (ignoring line wrapping)".
- A2 — ambiguity — non-blocking — §4 lines 48-49 and §5 line 54. "the second pair" suggests a copy of the clause in coordinator.md, but T-801 publishes the planning clause only in external-clis "The planning call". Fix: say it is published once, coordinator.md only references it, and the manual sync is review × planning.
- A3 — ambiguity — non-blocking — §3 bullet 3 (line 41). "referenced" can be read as applying to the prompt, while AC-05 requires the "Reading is allowed" block verbatim in the prompt. Fix: copied verbatim into the prompt; referenced, not rewritten, in the skill.
- A4 — organization — non-blocking — §3 bullet 5 (line 43) and §7. Missing citations: RF-03, RF-05/G19, RF-13/AC-15; RF-13 is absent from §7.
- A5 — ambiguity — non-blocking — §4 line 47. "the same signals as the CLI reviewer" contradicts §3 bullet 4 (the verdict table does not apply). Fix: the same freeze record (`clean`/`reverted`) and the same not-run classes.
- Not raised: scope, header, template shape, the §1 quotes, G6/AA1/RISK-001, the §3 rationale and bullets 1, 2 and 4, the description of `checkAntiDelegationClause`, and the present-tense claims (all verified).

### Coordinator dispositions
- A1–A5 accepted; correction 1 of 2 returned to the author (solution-architect).

## Round 2 (2026-09-28) — final
- Verdict: approved — 0 blocking.
- A1 applied: the §5 signal was checked against T-801:38-51 and SKILL.md:72-82; a correct implementation passes and real drift is flagged. A2–A5 applied.
- No regression: only §3 bullets 3 and 5, §4:47-49, §5:55 and §7:63/65 changed. No new findings.
- Validator note, not raised: §5:54 ("outside and after" the first fenced block of the reviewers section) is met automatically, since coordinator.md holds no copy of the clause.
