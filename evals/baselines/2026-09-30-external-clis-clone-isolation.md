# external-clis-clone-isolation + external-clis-clone-lifecycle — baseline (2026-09-30)

Cases: `evals/cases/external-clis-clone-isolation/` (T-1010, PRD-008 R13 actions a–l) and `evals/cases/external-clis-clone-lifecycle/` (T-1011, R1–R4, R9–R12). Package built from `feat/prd-008-worktree-isolation`. Local `claude plugin eval`, model and judge sonnet, `--runs 1`, `--ablation none`, `--scaffold`, `--case 'external-clis-clone-*'`. Budget authorized by the owner: USD 3.50 for the whole measurement (D1).

Limit of this measurement: it scores how the Coordinator classifies canned before/after snapshots given in fixtures. It does not show that the R5/R6 commands detect anything in a real tree; that evidence is T-1015 (executed by the Coordinator on a throwaway copy, every action in agreement after two kit-text corrections, local and git-ignored under `.harness/tasks/T-1015/evidence/`).

## Run 1

| Case | Score | Failed graders | Cost (agent + judge) |
|---|---|---|---|
| clone-isolation | 0.737 (14/18 graders, weighted) | 01 (a), 02 (b), 05 (e), 17 (f regex) | 0.689 + 0.332 |
| clone-lifecycle | 0.619 (12/19 graders, weighted) | 01, 02, 04, 05, 08, 15, 18 | 0.700 + 0.338 |

Total USD 2.058 (agent 1.389 + judge 0.670; the tool's reported `costUsd` of 1.389 excludes the judge). Reading the agent's final message (grader `evidence`): the Coordinator's answer matched the kit text in every failing scenario. Case defects: two regex graders fired on legitimate text (a negation "no not-run reason, nothing reverted" in (f); the live path quoted inside the `grep -F` check in (r4)); several llm graders demanded exact phrasing or details the scenario did not ask for; grader 15 was ambiguous about Read/Glob; the lifecycle prompt pointed to `.commands/plan-loop.md`, which does not exist in the packaged plugin. Decided by the Coordinator on the owner's behalf: correct the graders and the prompt path, repeat only these two cases.

## Run 2 (after the grader corrections)

| Case | Score | Failed graders | Cost (agent + judge) |
|---|---|---|---|
| clone-isolation | 0.882 (15/17) | 14 (clone removed at the exact stored path after R5 and R6, in every clone scenario), 15 (no real call or git) | 0.677 + 0.338 |
| clone-lifecycle | 0.950 (18/19) | 04 (r3: the answer did not state that the clone holds no ignored file and no `.harness` / `.claude/settings.local.json`) | 0.749 + 0.365 |

Total USD 2.129 (agent 1.426 + judge 0.703). Both runs: **USD 4.187, above the USD 3.50 the owner authorized**. The Coordinator read the tool's `costUsd` (agent only, 1.389 after run 1) as the whole spend and set run 2's `--max-cost-usd 2.1` from it; `--max-cost-usd` also checks agent cost only. The overrun is USD 0.687, reported to the owner.

## Outcome after run 2 (superseded by Round A below)

- The pass bar used by earlier baselines (every grader passes) is **not** met: 3 graders out of 36 fail in run 2. The CI gate threshold of `.github/workflows/evals.yml` (0.8) is met by both cases.
- The budget is already exceeded, so no further run. Decided by the Coordinator: stop here, record the three misses as open, no routine repetition.
- Lesson for later measurements: budget the judge separately; `--max-cost-usd` and the aggregate `costUsd` count the agent only (the per-run `judgeCostUsd` is separate).
- Open misses: lifecycle 04 is a real omission in the answer (the r3 report did not restate what the clone never holds); isolation 14 and 15 were judged FAIL 3/3 and their evidence was not inspected further within this budget, so they stay open, neither dismissed nor fixed.

## Round A (owner's choice 2 on 2026-09-30: "faça mais 3 rodadas")

Changes before the run, from the run-2 evidence: the clone-isolation prompt asks, for scenarios g–k, what happens to the clone and whether any file of it reaches the live tree (grader 14 unchanged); grader 15 (llm, "no real call or git", unverifiable from the last message) replaced by a deterministic `tool_used` grader on `Bash`, min 0, max 0; the clone-lifecycle r3 prompt also asks what the clone never contains (grader 04 unchanged).

| Case | Score | Failed graders | Cost (agent + judge) |
|---|---|---|---|
| clone-isolation | 1.0 (17/17) | none | 0.992 |
| clone-lifecycle | 1.0 (19/19) | none | 1.222 |

Round total USD 2.214. The all-graders bar is met; the two remaining authorized rounds were not needed. Measurement total over three runs: USD 6.40 (agent + judge).
