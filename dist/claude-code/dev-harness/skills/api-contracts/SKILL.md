---
name: api-contracts
description: Use when a public interface between components is being created or changed and consumers must agree with the provider before code is written, including HTTP, RPC, event, webhook, and cross-module APIs, or when a change may break deployed clients and needs a migration path. Do not use when the change is internal with no consumer-visible surface, when the contract already exists and only its implementation is in question, or when the task is to build the server rather than specify it.
author: malka
metadata:
  provenance: adapted
  sources: ["api-design reference (REST/GraphQL conventions)", " api-designer role description", "dev-harness api-designer agent"]
---

# API Contracts

## Overview

This procedure produces a contract that a consumer and a provider read identically: operations, error catalog, one success example, one failure example, and an explicit register of what breaks. The contract is written before the server exists, so disagreement surfaces in text rather than in integration. The core principle is applicability: every concern below is included only when a stated trigger is present, and anything included without a trigger is speculation, not design.

## When to use

- A new HTTP, RPC, event, webhook, or cross-module interface must be defined.
- An existing interface changes shape, status codes, error codes, or authorization.
- Two roles are about to implement the same interface from different sides.
- A consumer already deployed may break, and the migration path is unclear.
- Acceptance criteria depend on what the interface returns in success and in failure.

## When not to use

- Internal refactor with no consumer-visible surface. Use safe-refactoring.
- A single known endpoint whose behavior is already specified and only needs building. Use incremental-implementation.
- Picking the transport or the service boundary itself. Use architecture-decisions.
- Pure UI work against an unchanged API.

## Inputs

| Input | If missing |
| --- | --- |
| Authorization and write set for the contract file | Do not write. Return the contract in the reply and request the path and scope from the Coordinator. |
| The consumers of the interface and what each one does with it | Ask the Coordinator. An unknown consumer changes pagination, auth, and compatibility. Do not guess. |
| Existing conventions in the repository (paths, error shape, auth, naming) | Read two or three neighboring contracts and adopt what they do. Record it as a decision, not as a fact. |
| Acceptance criteria that depend on the contract | Write the contract for the operations you can source, and list the criteria you could not map under Open decisions. |
| Current version and deployment status of existing consumers | Treat every consumer as deployed. That is the safe assumption and it only costs a migration note. |

## Procedure

1. **Read before writing.** Locate existing contracts, error shapes, auth middleware, and naming conventions with local search. Cite each convention with a relative path. A contract that contradicts its neighbors is a defect even when it is internally consistent.
2. **Name the consumers.** List each consumer and the single question it asks the interface. An operation that no listed consumer calls does not belong in this contract.
3. **Model operations from the domain, not from the storage.** For each operation record intent, inputs, output on success, preconditions, and observable effects. Name actions by intent, not by table mutation.
4. **Apply the applicability gate.** Walk the table in Quick reference. For each concern, either state the trigger that makes it apply and specify it, or record one line saying it does not apply and why. Silence is not a decision.
5. **Define one error catalog for the whole surface.** Each error carries a stable machine code, a transport status, a human message, optional per-field details, and a correlation identifier. The same failure never appears under two codes.
6. **Write one success example and one failure example per operation.** Complete, concrete, copy-pasteable values, no placeholders like `...`. These two examples are the shared artifact: the consumer builds its parser from them and the provider builds its response from them.
7. **Register compatibility.** Classify every delta as additive, breaking, or unresolved pending named consumer evidence using the compatibility table. Do not finalize compatibility-dependent work while that evidence is missing. For each breaking item record what breaks, which consumer breaks, and the expand and contract sequence that avoids a flag day.
8. **Write integration scenarios.** One given/when/then per acceptance criterion, each pointing at an operation and at one of the two examples. These become the inputs for regression-testing.
9. **Stop at the specification.** Do not implement the server, the client, or the migration in this role. Hand the contract back with evidence, limitations, and open decisions.
10. **Report, do not record.** Facts, evidence, open decisions, and any incident go back to the Coordinator, who is the only writer of `.harness/MEMORY.md`, `.harness/EPOCHAL.md`, and `.harness/RISKS.md`. Write the contract only inside the write set assigned by the Coordinator, normally `.harness/tasks/<id>/`. If the project keeps specifications elsewhere, request that path before writing.

## Output format

```
# Contract: <interface name>

## Consumers
- <consumer> -> <what it needs from this interface>

## Conventions followed
- <convention> (source: <relative path>)

## Operations
### <operation id> <transport form, e.g. POST /v1/invites>
- Intent:
- Request: <fields, types, required or optional>
- Success: <status> <response shape>
- Errors: <code> -> <status> -> <when>
- Authorization: <rule, or "not applicable: reason">
- Idempotency: <rule, or "not applicable: reason">
- Pagination: <rule, or "not applicable: reason">
- Versioning: <rule, or "not applicable: reason">
- Rate limiting: <rule, or "not applicable: reason">
- Partial failure: <rule, or "not applicable: reason">

## Error catalog
| code | status | meaning | field details | correlation id |
| --- | --- | --- | --- | --- |

## Examples
### Success
<complete request and response>
### Failure
<complete request and response>

## Compatibility
- Additive: <items>
- Breaking: <item> | <consumer affected> | <migration step>
- Unresolved: <item> | <consumer> | <evidence needed before finalizing compatibility>
- Sequence: expand -> migrate consumers -> contract

## Integration scenarios
- Given <state>, when <call>, then <observable result> (criterion: <id>)

## Open decisions
- <question> | <who decides> | <what changes with each answer>

## Evidence
- <relative path> - <what it establishes>
```

## Quick reference

Applicability gate. Include the concern only when the trigger holds.

| Concern | Include when | Record when skipped |
| --- | --- | --- |
| Authorization | The operation accesses protected data or requires a role, capability, entitlement, or object-level permission, including on the caller's own data | "Unrestricted public operation: <reason no permission is required>" |
| Idempotency | A retry is realistic and the effect is not naturally repeatable (creation, charge, dispatch) | "GET, naturally idempotent" |
| Pagination | The collection grows without a bound the caller controls | "Fixed small set, bounded by design" |
| Versioning | At least one consumer is already deployed against the old shape | "First release, no deployed consumer" |
| Rate limiting | The surface is public, expensive, or abusable | "Internal, behind a trusted caller" |
| Partial failure semantics | The operation touches more than one resource or system | "Single resource, atomic" |

Compatibility classification.

| Change | Class |
| --- | --- |
| New optional request field or new operation | Additive when existing calls keep their meaning and behavior |
| New response field | Additive only when the consumer contract or verified consumers tolerate unknown fields; otherwise breaking, or unresolved until checked |
| New error code for a previously unhandled condition | Additive if clients have a default branch, otherwise breaking |
| New enum value | Breaking unless consumers are specified to tolerate unknown values |
| Remove or rename a field, operation, or error code | Breaking |
| Make an optional request field required, or tighten validation | Breaking |
| Change a type, a status code, a default, or a unit | Breaking |
| Change pagination defaults or ordering | Breaking |

Expand and contract, the only sequence that avoids a flag day: add the new shape beside the old one, serve both, migrate each consumer with evidence, announce the removal date, then remove.

Do not label an unresolved compatibility check additive. Name the affected consumer and the evidence needed to decide. For example, a consumer that rejects unknown response properties can break when a field is added. Ownership alone also does not settle authorization: changing one's own account can still require a role or entitlement.

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Fields added "for later" | Every speculative field becomes a compatibility obligation the moment one consumer reads it | Only fields a named consumer needs now |
| Examples with placeholders | Consumer and provider fill the gaps differently and the disagreement appears in integration | Complete concrete values on both sides |
| Two error codes for the same cause | Clients branch on transport status and lose the distinction anyway | One stable machine code per cause |
| 5xx for a caller mistake | Hides client bugs inside provider alerting and defeats retry policy | 4xx for caller, 5xx only for provider fault |
| Contract shaped like the database tables | Storage refactors then become breaking API changes | Model the domain intent |
| Breaking change with a version bump and nothing else | A version number is not a migration; deployed clients still break | Register affected consumer plus expand and contract steps |
| Writing the server while writing the spec | The spec stops being reviewable and becomes a description of the code | Specification only, implementation is another role |
| Treating a draft as deployed | Downstream work plans against something that does not exist | Mark status explicitly |

## Example

Input: "Add an invite endpoint. Frontend creates invites, the worker retries on timeout."

```
### create_invite  POST /v1/invites
- Intent: invite an email address to an organization.
- Request: { "email": string (required), "role": "member" | "admin" (required) }
- Success: 201 { "id": "inv_7Kq", "email": "sam@example.com",
                 "role": "member", "status": "pending" }
- Errors: invalid_email -> 422 | already_invited -> 409 | forbidden -> 403
- Authorization: caller must hold admin on the target organization.
- Idempotency: Idempotency-Key header required; a replay returns the original 201 body.
- Pagination: not applicable, single resource.

### Failure example
POST /v1/invites   { "email": "sam@", "role": "member" }
422 { "error": { "code": "invalid_email", "message": "Email is invalid",
      "details": [{ "field": "email", "issue": "format" }],
      "request_id": "req_abc123" } }

## Compatibility
- Additive: create_invite is a new operation, no deployed consumer

## Integration scenarios
- Given an org with a pending invite for sam@example.com, when create_invite is
  called again, then 409 already_invited (criterion: AC-2)

## Open decisions
- Expiry of a pending invite | product owner | an expiry adds a field and a
  second reason for 409 to disappear
```

Idempotency is present because the worker retries. Pagination is skipped with a stated reason. The failure example is what the frontend builds its error branch from.

## Related

Roles: api-designer, solution-architect, backend-builder, frontend-builder, qa-verifier. Commands: used inside `/dev-harness:plan`. Skills: architecture-decisions, implementation-planning, incremental-implementation, regression-testing.

## Proof case

Given a contract produced by this procedure, a consumer role and a provider role independently describe the same success body and the same failure body, including status, error code, and field names, without further clarification. A change that removes a response field is classified as breaking, names the affected consumer, and carries an expand and contract sequence rather than a bare version bump. Adding a response field for a consumer that rejects unknown properties is also breaking; absent tolerance evidence remains unresolved. An operation on the caller's own account that requires an entitlement retains that authorization rule.
