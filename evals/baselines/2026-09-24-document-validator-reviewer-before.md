# document-validator — before baseline (T-301, PLAN-006)

- Case `document-validator-scope-overclaim`, unedited `.agents/document-validator.md`, CLI 2.1.282, `--runs 3 --model sonnet --judge-model sonnet --max-cost-usd 1.25`.
- Result: score 0.67, 0/3 full passes.
  - 01 changes required: 3/3.
  - 02 blocking ≥1: 3/3.
  - 03 cites the contradicted source point: 3/3.
  - 04 correct claim not flagged: 1/3.
  - 05 correct claim listed under `## Not raised`: 0/3.
- Classification: **reproduced**. The overclaim is detected every time, but the unedited validator fails the controls. It flags the correct scaffold claim or leaves it out of `## Not raised`, so it gives no proof of what was checked. That is D1b's target.
- No hardening pass applied.
- Cost: agent 0.59 + judge 0.07.
- Limit: the eval runs on sonnet, but `document-validator` ships on opus (PLAN-006 P9).
