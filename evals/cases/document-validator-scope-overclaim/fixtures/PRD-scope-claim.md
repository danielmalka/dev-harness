# PLAN-910 · Reviewer eval case batch

| Field | Value |
|---|---|
| Status | in review |
| Owner | harness maintainer |
| Created / updated | 2026-09-24 / 2026-09-24 |

## 1. Goal

Add seeded-defect eval coverage for the reviewer roles.

## 2. Scope

- This batch adds one seeded-defect eval case for each of the four reviewer
  roles — code, security, document and coordinator — pairing each planted
  defect with a correct control section in the same fixture.
- Each case's fixtures are copied into the run workspace before the run by
  the case's own scaffold script, so no case depends on a shared fixture
  directory.

## 3. Acceptance criteria

- AC-01 Each of the four cases has a positive grader that detects the
  planted defect.
- AC-02 Each of the four cases has a negative grader confirming its control
  section is not wrongly flagged.

## 4. References

- source-notes.md
