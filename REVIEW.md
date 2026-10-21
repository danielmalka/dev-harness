# Review guide

## How to use

This file tells the reviewer what to look for in this project. The owner fills it. The reviewer reads it first and ties each finding to one rule below, or says no rule applies. Keep entries few and concrete. Always write it in English.

## Rules

One entry per recurring problem. Copy the block for each new rule.

### Generated dist edited by hand

- Symptom: the diff changes a file under `dist/` without a matching change in a canonical source (`.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`).
- Why it matters: `dist/` comes from `go run ./cmd/dh build`. A hand edit is lost at the next build, or makes CI fail on the `dist` diff check. The rule is in `AGENTS.md` ("Fontes canônicas").
- Example: `dist/claude-code/dev-harness/agents/coordinator.md:2`

### Count in docs drifts from the real count

- Symptom: a doc, README or tutorial states a count (13 agents; 19 commands in files plus 1 registered by the mod) that no longer matches the sources, or one language states a different number than its mirror.
- Why it matters: the owner trusts the stated counts; `AGENTS.md` fixes them ("Agentes (13)"). A stale number says the kit lost or gained capability when it did not.
- Example: `docs/tutorial.html:781`

### Private path in a PRD, review report or doc

- Symptom: text contains an expanded home path or another concrete machine path, instead of `<harness dir>` or `<home>`.
- Why it matters: it leaks the author's machine into a shared repository. `dh validate` (`checkPrivatePaths`) rejects it and blocks the delivery.
- Example: `internal/kit/validate.go:577`

### Change leaves stale references

- Symptom: a file, count, link or copy is removed or renamed, but references to it remain in docs, templates or the pt-br/en mirror.
- Why it matters: `AGENTS.md` ("Premissa de limpeza") blocks delivery when a change leaves obsolete leftovers.
- Example: `AGENTS.md:12`

## Ignore

One entry per pattern that is noise. Each entry reads "ignore this pattern in this kind of file". A finding that matches an entry is reported as ignored-by-rule and does not spend a round.

- Ignore the `generated_at` line in `dist/claude-code/dev-harness/harness-manifest.json`. It differs between builds, and CI tolerates it.
- Ignore a finding inside a `dist/` copy when the same finding exists in its canonical source. Review the source, not the copy (for example `dist/claude-code/dev-harness/agents/coordinator.md` against `.agents/coordinator.md`).

## Critical paths

Path globs in `.gitignore` style. A pull request that touches a match stops before the merge, and the owner decides. Leave the list empty for no stop.

- `adapters/claude-code/plugin/hooks/**`: runtime guard of the memory lock and the dispatch rule.
- `.skills/external-clis/**`: isolation of external CLI reviewers.
