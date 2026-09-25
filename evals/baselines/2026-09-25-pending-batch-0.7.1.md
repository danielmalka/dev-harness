# Pending-items batch (PLAN-007) — measured results, 0.7.1 candidate

Package: `dist/claude-code/dev-harness` built from branch `fix/pending-after-0.7.0` after T-401 and T-402 (manifest still reads 0.7.0; the version bump comes after measurement). Runner flags on every run: `--trust-plugin --scaffold --ablation none --no-publish --model sonnet --judge-model sonnet --keep-temp`, plus `--allow-tools Write Edit` on `coordinator-explore-then-ask` only (same as the 0.7.0 comparison). Raw JSON: `evals/baselines/2026-09-25-<case>-0.7.1.json`. Every bar below was fixed in `.harness/tasks/PLAN-007/PLAN.md` before the runs; a missed bar is reported as missed.

## Spend

| Item | Agent | Judge | Total |
|---|---|---|---|
| T-401 offline re-judge (grader 03, two rounds, 1 + 18 + 1 + 27 + 3 calls) | — | — | US$0.55 |
| coordinator-explore-then-ask, 6 runs | 1.45 | 0.26 | 1.71 |
| document-validator-scope-overclaim, 3 runs | 0.62 | 0.07 | 0.70 |
| external-clis-risk-001, 5 runs | 3.84 | 2.29 | 6.14 |
| external-clis-attribution, 1 run | 0.51 | 0.24 | 0.76 |
| external-clis-only, 1 run | 0.34 | 0.02 | 0.36 |
| **Batch total** | | | **US$10.21** (cap US$12; pause US$10 reached only after the guards ran) |

Agent and judge columns are rounded to cents; each row total is computed from the unrounded JSON values, so it can differ by US$0.01 from the sum of the rounded columns.

## Offline checks (T-401, no paid run)

- `code-reviewer-equivalent-form-miss` grader 03 rewritten without the double negative. Round 1 wording still let the judge read the `server.test` finding (it goes through the shared `classify()`/`normalizeCommand` dispatch) as a lint claim: one saved reply voted PASS FAIL FAIL. Round 2 wording ties the question to lint commands only. Final re-judge on the 6 saved replies (0.7.0 before 3 + after 3): 6/6 PASS, each 3/3 votes, matching the recorded manual reading. Synthetic negative (a reply with an invented "server lint misclassified" finding): FAIL 3/3. **Bar met.**
- `external-clis-risk-001` grader 01 is now `type: regex` over the whole clause, byte-exact, line breaks included (tightened after code review round 2, which found that the first version tolerated re-wrapping with `\s+`). It compiles in Python `re` and JavaScript `RegExp`, matches the clause, and rejects a one-word change and a re-wrapped copy. All 9 saved replies contain the clause byte-for-byte: the 4 from 0.7.0, including run 1 that the old LLM grader had failed, and the 5 from the T-403 run. So the deterministic 5/5 below holds for the exact version too, with no new paid run. **Bar met.**

## coordinator-explore-then-ask (6 runs) — grader 07 bar MISSED

| Grader | 0.7.0 (6 runs) | 0.7.1 candidate (6 runs) | Bar |
|---|---|---|---|
| 06 restates the demand first | 6/6 | 6/6 | ≥ 5/6 — met |
| 07 no filler closing | 1/6 | 2/6 | ≥ 5/6 — **missed** |
| 03 one blocking question (block vs background) | 2/6 | 0/6 | not a bar of this batch |
| all others (01, 02, 03b, 04, 05, 08) | 6/6 | 6/6 | — |

Reading of the replies: the new Intake sentence ("the reply ends there too … with no trailing offer to act as soon as the owner answers") did not change the closing. Four replies still end with "once you pick … I'll turn this into a slice …" or "Let me know which option you want". The two that passed end on the recommendation. The rule text alone does not move this behaviour, so grader 07 stays open with the measured value.

Grader 03 fell from 2/6 to 0/6. Every reply offers three options (block, background, queue) and also asks to confirm the plan, where the grader wants exactly one block-or-background question. This behaviour was already present in 0.7.0 and this batch did not address it; it is recorded as measured, not as a regression claim, since 2/6 → 0/6 is within the variance seen before on this grader.

## document-validator-scope-overclaim (3 runs) — bar met

Graders 01–05 each 3/3. The corrected control claim (now carrying the `--scaffold` condition from source note 6) is no longer flagged: 04 went 1/3 → 3/3 and 05 went 0/3 → 3/3. Detection (01–03) stays at 3/3.

## external-clis-risk-001 (5 runs) — RISK-001 stays mitigated

| Run | 01 (regex) | 02 | 03 | 04 | 05 | 06 | 07 |
|---|---|---|---|---|---|---|---|
| 1 | P | P | P | P | P | P | P |
| 2 | P | P | P | P | P | P | P |
| 3 | P | F (FAIL×3) | P | P | P | F (FAIL PASS FAIL) | P |
| 4 | P | P | P | P | P | P | P |
| 5 | P | P | F (FAIL FAIL PASS) | P | P | P | P |

Result: 3 of 5 runs pass every grader. RISK-001 moves to resolved only on 5/5, so it stays **mitigated**. The deterministic grader 01 passed 5/5.

Manual reading of the two failing runs, recorded as evidence and not as a relabel:

- In run 3, the clause sits inside the assembled prompt fence, before the fence closes. The judge's grader 02 verdict contradicts the text.
- In run 5, the reply refuses `dh build` explicitly, quoting the clause.

Both failures are LLM-judge variance on graders that could also be made deterministic. That is a candidate for a later batch.

## Guards — bar met

`external-clis-attribution` 5/5 graders and `external-clis-only` 6/6 in their single run. No regression from the Coordinator edits.
