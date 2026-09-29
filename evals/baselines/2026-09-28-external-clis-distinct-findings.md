# external-clis-distinct-findings — baseline (2026-09-28)

Case: `evals/cases/external-clis-distinct-findings/` (T-804 item b: two reviewers, two DIFFERENT findings, which must stay two entries). Model and judge sonnet, `--runs 1`, `--max-cost-usd 1.00`. Pass bar: all graders pass.

- Paid run: graders 01–03 PASS 3/3. Old 04 FAIL 3/3, a case defect: the "Is neither…" wording is the double-negation trap. The final reply reproduces both findings verbatim, each with its own `Reported by:`, and records a disagreement for the owner about codex's mechanism, as the merge rule requires. US$0.27 agent + 0.04 judge.
- Correction 1: old 04 was replaced by two deterministic regex graders (04 codex text verbatim, 05 agy text verbatim), validated offline against the run's final message (both match) and against paraphrased negatives (no match). Following the offline rule for deterministic graders: **5/5**.
- Batch eval spend (US$6 cap): ≈ US$2.26.
