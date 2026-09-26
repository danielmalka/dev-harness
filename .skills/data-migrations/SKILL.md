---
name: data-migrations
description: Use when a change alters persisted structure or data, including new tables and columns, type or constraint changes, renames, drops, backfills, index changes, and the queries that depend on them. Also use when a schema change must ship without downtime, or when someone needs to know whether a migration can be rolled back. Do not use for a query change that leaves the schema untouched, and never use it to apply changes against a production database.
author: malka
metadata:
  provenance: adapted
  sources: [" postgres-pro", "parallel change practice"]
---

# Data Migrations

## Overview

Design a persistence change as a sequence of steps that are each safe on their own, applied in an order where the running application is correct before, during, and after every step. Structure moves through expand and contract: add the new shape, write to both, move readers, then remove the old shape once nothing reads it. Every migration ships with a rollback plan, or with an explicit statement that rollback is impossible and why. Verification happens on a disposable database with synthetic data. Never against real data, and never against a database whose address came from an environment file, a configuration default, or a previous session. Applying to a shared or real environment is an owner decision outside this procedure.

## When to use

- A feature needs a new table, column, index, constraint, enumeration value, or relationship.
- An existing column changes type, nullability, default, or meaning.
- Data must be backfilled, corrected, deduplicated, or moved between tables.
- A structure is being removed or renamed and consumers still reference it.
- A query became slow and the fix is an index, a denormalization, or a partitioning decision.

## When not to use

- The change rewrites a query against a schema that is not moving. That is ordinary implementation work.
- The task is to operate a database: backups, replication, failover, capacity. That is operations, not this procedure.
- Someone wants a migration executed against a shared or production environment. That is an owner decision with its own authorization, outside this procedure.

## Inputs

| Input | If missing |
| --- | --- |
| The behavior the change must support, with its acceptance criteria | Stop and get it. A schema designed from a guess outlives the guess. |
| The current schema: tables, columns, types, constraints, indexes | Read it from the migration history in the repository. Do not connect to an environment to discover it. |
| The access patterns: which queries read and write this data, and how often | List the call sites found in the code and mark unresolved consumers as unknown. |
| Data volume and growth, at least in orders of magnitude | Assume the table is large enough to need batching and say that you assumed it. |
| The deployment model: whether old and new application code run at the same time | Assume they overlap. It is the safe assumption and it forces expand and contract. |
| A disposable database to test against | Report the migration as untested and say what could not be verified. Do not substitute a shared environment. |

## Procedure

1. **State the target shape and the invariant it protects.** Write the intended structure and, for each constraint, the rule it enforces in business terms. A constraint with no stated rule is either missing or arbitrary. Prefer expressing integrity in the database where the database can enforce it, because application checks do not survive a second writer.

2. **Map the consumers.** Find every read and write of the affected structure in the repository: queries, ORM models, serializers, reports, background jobs, exports, and fixtures. Each one becomes a row in the change plan. A consumer you did not find is the one that breaks.

3. **Classify the change.** Additive, mutating, or destructive. Additive changes are safe when nullable or defaulted. Mutating changes need a compatibility window. Destructive changes need an explicit decision, because they are the ones that cannot be undone.

4. **Decompose into expand and contract.** For anything beyond a purely additive change, write the sequence:
   - Expand: add the new structure alongside the old, nullable or defaulted, with no consumer depending on it.
   - Dual write: the application writes both shapes. Inventory old application instances, workers, and other writers that still update only the old shape. Before completing the backfill, retire those writers or establish synchronization that keeps both shapes current for their writes.
   - Backfill: fill the new shape for existing rows, in batches, restartable, with no long-running transaction. Prevent a backfill from overwriting a newer concurrent write; state the synchronization or conditional-update rule that preserves it.
   - Migrate readers: first verify that all active writers maintain both shapes and that existing rows have converged. Then move readers to the new shape, one at a time, verifying each. A completed batch cursor alone does not establish convergence.
   - Enforce: add the constraint or non-null requirement once the data supports it.
   - Contract: remove the dual write, then the old structure, in a later deployment, once nothing reads it.
   Each phase is a separate migration and usually a separate release. Collapsing them is what produces the outage.

5. **Plan the locks.** For each step, state what it locks and for how long. Adding a column with a volatile default, changing a column type, and adding a constraint that validates existing rows are the usual offenders; they can hold a write lock over the whole table. Prefer the non-blocking form the engine offers: create the index concurrently, add the constraint as not validated and validate it in a second step, add the column without a rewrite and backfill separately. If a step must take a strong lock, say so, bound it, and set a lock timeout so it fails fast instead of queueing every other query behind it.

6. **Plan the transactions.** Schema changes go in a transaction when the engine supports transactional data definition, so a failure leaves nothing half-applied. Backfills go outside, in batches, because a single transaction over millions of rows blocks cleanup and risks exhausting the log. Say which of the two each step is. Never mix a long backfill into the same transaction as a structural change.

7. **Plan the indexes deliberately.** Add an index because a named query needs it, and say which one. Check whether an existing index already serves the predicate as a prefix. Build concurrently where the engine allows it, and know that the concurrent build can fail and leave an invalid index that must be dropped and retried. Removing an index is also a migration: confirm no query depends on it before dropping.

8. **Write the backfill as a restartable job.** Bounded batches, ordered by a stable key, with a recorded high-water mark so an interrupted run resumes instead of restarting. A pause between batches so ordinary traffic keeps its share of the database. Idempotent per row, so re-running a batch is harmless. Include the query that reports how many rows remain.

9. **Write the rollback plan, or state that there is none.** For each step, give the inverse: drop the added column, restore the previous default, recreate the dropped index. Where the inverse cannot restore the data, say so in those words. Dropping a column, truncating, deduplicating, narrowing a type, and overwriting a value during a backfill are one-way steps. For each one-way step, record what would be lost, whether a backup exists that the owner has actually verified, and what the forward recovery path is. An unverified backup is not a rollback plan.

10. **Update the affected queries with the schema.** Rewrite the reads and writes found in step 2 so they work in the compatibility window: tolerant of the new column being null, not selecting every column blindly, not depending on column order. State which queries change in which phase, and check the plan for the ones whose cost changes because of a new index or a changed type.

11. **Test on a disposable database with synthetic data.** If the project already provides a disposable database and the dispatch authorizes running it, create the schema from the migration history and seed it with generated data that covers the shapes that matter, including nulls, duplicates, boundary values, and the volume needed to see the batching work. Otherwise record every verification row as not-run and report the migration as untested. Do not install or start a database service, and do not provision one, to obtain a test target. With a target available:
    - Apply the migration forward and confirm the structure and the constraints.
    - Confirm the invariant: the rows that must exist, the rows that must be rejected.
    - Run the affected queries and confirm the results and, where relevant, the plan.
    - Apply the rollback and confirm the database returns to a usable prior state, or record that it cannot.
    - Interrupt the backfill partway and restart it, confirming it resumes without duplicating work.
    - For dual-shape migrations, exercise overlapping writers and a backfill racing with an update. Confirm stale old-only writes are prevented or synchronized and the newest value survives before readers switch.
    Never seed with a copy of real data containing personal information, and never point the test run at a database you did not create for the test.

12. **Write the application plan.** Ordered steps, each with what it does, its expected duration on the real volume, what it locks, how to tell it is progressing, and the point of no return. Name the metric to watch during each step: replication lag, lock waits, error rate, row count remaining. Name the abort condition for each step and the action it triggers.

13. **Stop at the decision.** Do not apply a migration to any real environment automatically, and do not connect to a database whose address you inferred from an environment file, a configuration default, or a previous session. Hand the plan, the evidence from the disposable run, and the residual risks to the owner. Only the Coordinator writes `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`; report a one-way step, a failed rollback, or a data-loss near miss as an incident for the Coordinator to record.

## Output format

```
## Change
- Goal: <behavior this supports>
- Current shape: <tables and columns affected>
- Target shape: <what it becomes>
- Invariants: <rule> enforced by <constraint>
- Classification: additive | mutating | destructive

## Consumers
| Consumer (relative path) | Reads | Writes | Phase it changes in |
|---|---|---|---|

## Phases
| # | Step | Kind | Locks | Transaction | Reversible |
|---|------|------|-------|-------------|------------|
| 1 | <expand: add column> | schema | none | yes | yes: drop column |
| 2 | <dual write> | code | none | n/a | yes |
| 3 | <backfill in batches of N> | data | none | per batch | no: overwrites <field> |

## Rollback
- Per step: <inverse or "none">
- One-way steps: <step> - <what is lost> - <verified backup? yes/no> - <forward recovery>

## Verification (disposable database)
- Seed: <what synthetic data was generated>
- <check> - passed | failed | not-run
- Forward: <result>
- Rollback: <result, or "not possible: reason">
- Interrupted backfill resumed: <result>
- Writer barrier and data convergence: <evidence that every active writer maintains both shapes and no stale rows remain>
- Concurrent update preserved: <result>

## Application plan
| Order | Step | Expected duration | Watch | Abort if |
|-------|------|-------------------|-------|----------|

## Residual risks
- <risk and who must decide>
```

## Quick reference

| Change | Safe form | Trap |
| --- | --- | --- |
| Add column | Nullable, or with a constant default the engine applies without a rewrite | A volatile default rewrites every row under a write lock |
| Add non-null column | Add nullable, backfill, then enforce | Adding non-null in one step blocks and fails on existing rows |
| Rename column | Add new, dual write, migrate readers, drop old | A rename breaks every deployed instance of the old code instantly |
| Change type | New column with the new type, backfill, switch readers, drop old | An in-place type change rewrites and locks the table |
| Add constraint | Add as not validated, then validate separately | Validating on creation scans the whole table under a lock |
| Add index | Build concurrently, verify it is valid | A plain build blocks writes; a failed concurrent build leaves an invalid index |
| Drop column or table | Stop reading, ship, then drop in a later release | Dropping in the same release breaks the instances still running old code |
| Backfill | Bounded restartable batches with a high-water mark | One transaction over the whole table blocks cleanup and can exhaust the log |
| Delete or deduplicate rows | Owner decision, verified backup, reversible staging table first | There is no inverse; this is the step that loses data |

| Question | Answer before writing the migration |
| --- | --- |
| Can old and new application code run at the same time? | Assume yes, design for the overlap |
| What does each step lock, and for how long? | Stated per step, with a lock timeout |
| Is there an inverse for each step? | Stated per step, including "none" |
| How large is the table in production? | At least an order of magnitude, or batch anyway |
| Which queries change plan because of this? | Listed with the phase they change in |

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Shipping expand and contract in one release | Every instance running the previous code breaks the moment the migration lands | Separate releases, with the dual-write window between them |
| Renaming a column in place | Instantly breaks deployed code and every cached query plan that named it | Add, dual write, migrate readers, drop later |
| Treating an unverified backup as a rollback plan | The recovery path is discovered to be broken during the incident | Require a restore the owner has actually tested |
| Connecting to a database found in a configuration file | Runs a migration against an environment nobody authorized | Use only a database created for this test |
| Testing with a copy of production data | Spreads personal data into development and evidence | Generate synthetic data covering the same shapes |
| Leaving the affected queries for later | The schema lands and the application breaks on the next read | Change the queries in the same plan, phase by phase |

## Example

Input: split a single `full_name` column into `given_name` and `family_name` without downtime.

```
Phase 1 (release A): add given_name and family_name, both nullable. Reversible: drop them.
Phase 2 (release A): application writes all three columns on create and update. Reversible.
  Finish rollout to every writer, including workers; confirm no old-only writer remains.
Phase 3 (job): backfill in batches of 5000 ordered by id, high-water mark stored,
  splitting on the last space, leaving rows with no space for manual review.
  Atomically fill only rows whose new columns are still unset, preserving concurrent updates.
  One-way: overwrites nothing, so reversible by clearing the new columns.
Phase 4 (release B): resolve manual-review rows and verify convergence, then switch readers
  to given_name and family_name. Reversible: revert the code.
Phase 5 (release B): enforce non-null on given_name once the remaining-rows query returns zero.
Phase 6 (release C): stop writing full_name, then drop it.
  One-way: the original unsplit text is lost. Verified backup required before this step.

## Verification (disposable database)
- Seed: 50000 synthetic rows (no space, trailing space, multi-part family name)
- forward migration applied - passed
- rollback of phases 1-4 - passed
- interrupted backfill resumed - passed
- concurrent update preserved - not-run (no concurrent harness in this fixture)
```

## Related

Roles: data-engineer, backend-builder, coordinator, qa-verifier, devops-engineer, release-manager. Commands: used inside `/dh:plan` and `/dh:build`. Skills: architecture-decisions when the shape is a durable decision, api-contracts when the schema change reaches an exposed contract, regression-testing for the checks that protect the invariant, security-review when the data is sensitive, delivery-readiness for the rollout decision.

## Proof case

Given a fixture database and a change that splits one column into two, the plan preserves every fixture row through the forward migration and the rollback of the reversible phases, and the backfill resumes correctly after being interrupted. An old-only writer blocks the reader switch until retired or synchronized, and a concurrent update survives the backfill. The one phase that cannot be reversed is stated explicitly with what would be lost and what backup the owner must verify first. No remote database is contacted and no destructive step is applied without an owner decision.
