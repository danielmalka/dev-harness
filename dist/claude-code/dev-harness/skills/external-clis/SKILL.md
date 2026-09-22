---
name: external-clis
description: Use when the Coordinator dispatches a review stage (`document`, `code`, `security`) that lists a `cli:<binary>/<slug>` entry in `reviewers:` (`.harness/project.yaml`), when a CLI call returns without the stage's verdict signal and the failure needs classifying, or when a new CLI binary needs a headless invocation recipe added under `references/`. Do not use for builder roles, for interactive/TUI use of any CLI, or for the `verify` stage, which is out of scope in this version.
author: malka
metadata:
  provenance: adapted
  sources: ["personal skill ~/.claude/skills/fleet-clis (SKILL.md, agy.md, codex.md, grok.md, opencode.md, mcode.md, claude.md, gotchas.md), not versioned in this kit"]
---

# External CLIs

## Overview

An external CLI reviewer (`agy`, `codex`, `grok`, `opencode`, `mcode`, `claude`) is invoked headlessly with the same role prompt a Claude specialist would receive, so its verdict is evidence the Coordinator can trust exactly as much as a Claude reviewer's. The procedure has four parts that never change per binary: resolve the binary before assuming it exists, put the prompt on disk when it is longer than one line, apply the binary's read-only mode when it has one, and read the call's outcome through the stage's own verdict signal so a silent or malformed reply is never mistaken for approval. Per-binary flag order and quirks live in `references/<binary>.md`, one file per binary, because flag order is part of each binary's contract and mixing them in one file invites copying the wrong order.

## When to use

- The Coordinator dispatches `cli:<binary>/<slug>` for a review stage configured in `reviewers:` (`.harness/project.yaml`), alongside or instead of the stage's Claude agent.
- A role prompt (`document-validator`, `code-reviewer`, `security-reviewer`) must be sent to an external CLI binary instead of dispatched as a Claude specialist.
- A CLI call returned exit 0 with no recognizable verdict, or ran past its timeout, and the failure must be classified as transport failure rather than a verdict.
- A new CLI binary is added to the fleet and needs its own `references/<binary>.md` recipe.

## When not to use

- Builder roles (`backend-builder`, `frontend-builder`, `data-engineer`, `devops-engineer`). This version covers review stages only; a CLI never writes production code through this skill.
- Interactive or TUI use of any CLI from an orchestrator turn. Every call here is headless and one-shot.
- Installing, authenticating, or updating a CLI. A binary that is not already installed and authenticated on this machine is reported not-run; it is never provisioned.
- The `verify` stage (`qa-verifier`). Out of scope for this version; candidate for a future one once a read-only reviewer contract exists for QA.
- Choosing which reviewers a stage runs, merging N verdicts, or writing `.harness/local.yaml`. That is the Coordinator's own procedure; this skill only covers a single binary's invocation.

## Inputs

| Input | If missing |
| --- | --- |
| The role prompt body (`.agents/<role>.md`, without its frontmatter) | Stop. A CLI reviewer never receives a prompt other than the one a Claude specialist would receive for the same stage. |
| The stage's review skill text (`code-review`, `document-review` or `security-review`) | Stop. The verdict signal this skill checks for comes from that skill's own output format. |
| The dispatch context bundle (scope, diff or document, source material) | Stop. Same requirement as dispatching a Claude specialist for the stage. |
| `.harness/local.yaml`, when the binary is not on `PATH` | Treat the reviewer as not-run, reason "binary absent on this machine". Never invent a path. |
| `commands.test` / `commands.lint` from `.harness/project.yaml`, if the reviewer needs the suite's result | Omit that instruction from the prompt; the reviewer reasons from the code alone. |

## Procedure

1. **Resolve the binary.** `command -v <binary>` first. If that fails, look up `binaries.<binary>` in `.harness/local.yaml` (a bare command on `PATH`, a `~/...` path, or a path relative to the project — see `references/local.yaml.example`). If neither resolves, the reviewer is not-run, reason "binary absent on this machine". Never substitute a different binary for the one configured.
2. **Pin the model.** Take the slug from the `cli:<binary>/<slug>` entry the Coordinator dispatched. Do not invent a slug; an unmeasured slug is a transport failure at the CLI's own level (unknown model), not something this skill guesses around.
3. **Assemble the prompt**: the role prompt body, the stage's review skill text, the dispatch context bundle, and the anti-delegation clause below, verbatim, in the body of the prompt itself — never only in a frontmatter field or a dispatch instruction the CLI's runtime is trusted to enforce on its own.
4. **Write the prompt to disk** when it is longer than one line. Do not stuff a multi-kilobyte prompt into argv when the binary offers `--prompt-file`, stdin, or an attached-file flag; `references/<binary>.md` names which one that binary takes.
5. **Invoke with the binary's flag order**, from `references/<binary>.md`. Flag order is part of the contract for some binaries (a misplaced flag is read as the prompt itself); do not reorder from what the reference shows. Apply the binary's read-only flag when it has one — see the table below.
6. **Classify the outcome by the stage's verdict signal**, not by exit code alone — see "Transport failure vs. a real verdict" below. A non-zero exit, a timeout, or exit 0 with no recognizable verdict signal is a transport failure: it is reported not-run with the reason, and it does not count as either an approval or a rejection.
7. **Recompute the frozen state and compare.** Frozen before the call: `git status --porcelain -uall` for the path list, then `git hash-object <path>` for every listed path (modified or untracked; a path absent on disk records the marker `deleted` in place of the hash). After the call, recompute both. Any difference — a changed path list, or a changed hash for any path — is a transport failure: the verdict is discarded, the change is reverted, and the reviewer is not-run for that reason, regardless of what it said. This also catches an edit made inside a file that was already dirty in the frozen state, which the path list alone cannot distinguish. A reviewer is read-only by contract even when its binary has no read-only flag.
8. **Report to the Coordinator**: the resolved binary and slug, the verdict or the not-run reason, and the raw output location if one was written to disk. This skill does not merge verdicts across reviewers or write a persistent report; both are the Coordinator's procedure.

## Output format

What step 8 hands back to the Coordinator for one CLI reviewer call:

```
## CLI reviewer: <binary>/<slug>
- Resolved binary: <command found on PATH, or the path taken from .harness/local.yaml>
- Verdict: <the stage's verdict signal, verbatim> | not-run (<reason: binary absent, transport failure, timeout, or reverted change>)
- Timeout: <the value used, configured per reviewer or the 15-minute default>
- Raw output: <path on disk to the prompt and the reply, if either was written there>
- Workspace check (`git status --porcelain -uall` + `git hash-object` per path): clean | reverted (<what the reviewer changed, now reverted>)
```

This skill does not merge verdicts across reviewers or write a persistent report; the Coordinator folds this line into the stage's own report format (`code-review`, `document-review`, `security-review`), tagged with `cli:<binary>/<slug>` the same way it tags a Claude reviewer's findings with `claude`.

## The anti-delegation clause

Copy this block verbatim into the body of every CLI reviewer prompt, not just into the dispatch instruction around it:

```
Do not invoke another CLI binary, spawn another agent, or delegate any part
of this review to another tool. Do not execute any script in this
repository. The only commands you may run are the project's own
`commands.test` and `commands.lint`, exactly as declared in
`.harness/project.yaml`, and only to read the result of the existing test
or lint suite — never `commands.build` or any other repository command.
This instruction is written here, in the prompt text, because the runtime
does not enforce it by itself: a specialist in this project once ignored a
`disallowedTools` restriction declared only in its frontmatter and
dispatched another agent anyway. Follow the words in this prompt, not an
assumption about what the runtime blocks.
```

`commands.test`/`commands.lint` is the only execution this skill permits a CLI reviewer (owner decision, PRD-003 RF-01) — never `commands.build`, which can rewrite versioned artifacts such as `dist/` and trip the reversion in step 7 above. For `codex` under its read-only sandbox, this permission is inert in practice: [references/codex.md](references/codex.md) explains why.

## Read-only mode by binary

| Binary | Read-only flag | When none exists |
| --- | --- | --- |
| `codex` | `codex exec -s read-only` | — |
| `mcode` | `mcode exec --permission off` | — |
| `agy` | none documented | Guarantee comes from the post-call porcelain+hash-object check (step 7) |
| `grok` | none documented | Guarantee comes from the post-call porcelain+hash-object check (step 7) |
| `opencode` | none documented | Guarantee comes from the post-call porcelain+hash-object check (step 7) |
| `claude` | none documented, as a CLI reviewer distinct from a Coordinator-dispatched agent | Guarantee comes from the post-call porcelain+hash-object check (step 7) |

## Transport failure vs. a real verdict

The verdict signal is defined per stage, read from the stage's own output format. Every signal below is matched as a **prefix of a line**, never as exact whole-line equality — a trailing clause, a footnote, or punctuation after the matched text must never turn a real verdict into a transport failure:

| Stage | Verdict signal | Transport failure looks like |
| --- | --- | --- |
| `document` | A line under `## Verdict` starting with `approved` or with `changes required` | Exit 0 with no `## Verdict` section, or the section present with neither word |
| `code` | A line under `## Verdict` starting with `Review status:` | Exit 0 with no `## Verdict` section, or the line missing under it |
| `security` | A `### <SEVERITY>` entry under `## Findings`, or a line under `## Findings` starting with `No demonstrated vulnerability or confirmed exposure found` (the full sentence continues "in the reviewed scope; unresolved hypotheses remain listed below" — match the prefix, not the whole line) | Exit 0 with a `## Findings` section that has neither |

A non-zero exit, a timeout, or exit 0 without the stage's signal is a transport failure: not-run with the reason, never counted as an approval, and never spending a correction round. Help text, a reasoning preamble with no verdict, and an empty file on disk are all transport failures, not a lenient pass.

## `.harness/local.yaml`

The machine-specific binary paths and measured model slugs a CLI reviewer needs live in `.harness/local.yaml`, never in this skill or in any file under `.skills/`. Schema and placeholder values: `references/local.yaml.example`. The Coordinator is the only writer of that file; this skill only reads it.

## Quick reference

| Binary | One-shot flag | Prompt input | Reference |
| --- | --- | --- | --- |
| `agy` | `--print` / `-p` / `--prompt` | inlined in the prompt string (avoid `--add-dir` before the prompt) | [references/agy.md](references/agy.md) |
| `codex` | `exec` (or `exec review` / `review`) | stdin (`-`) | [references/codex.md](references/codex.md) |
| `grok` | `-p` / `--single` | `--prompt-file` for anything longer than a line | [references/grok.md](references/grok.md) |
| `opencode` | `run` | `-f` / `--file=path -- "message"` | [references/opencode.md](references/opencode.md) |
| `mcode` | `exec` | `--input -` (stdin) or `--file` | [references/mcode.md](references/mcode.md) |
| `claude` | `-p` / `--print` | last argv, or expanded from a file | [references/claude.md](references/claude.md) |

Recorded transport-failure patterns, generalized across binaries: [references/gotchas.md](references/gotchas.md).

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Treating exit 0 as a pass | Help text, a reasoning preamble, and a silent no-op all exit 0 | Read the stage's own verdict signal, never the exit code alone |
| Putting the anti-delegation clause only in the dispatch instruction | The runtime does not enforce a restriction stated outside the prompt body (RISK-001) | Copy the clause verbatim into the prompt text itself |
| Letting a CLI reviewer run `commands.build` | Can rewrite versioned artifacts such as `dist/`, which then reads as an unrelated change | The clause permits only `commands.test`/`commands.lint` |
| Assuming a binary without a read-only flag cannot touch the workspace | Several binaries in this fleet have no read-only mode at all | The post-call `git status --porcelain -uall` + `git hash-object` comparison (step 7) is the guarantee, not the binary's own flags |
| Stuffing a multi-kilobyte prompt into argv | Several binaries truncate, misparse, or treat part of it as a second argument | Write the prompt to disk and pass it the way `references/<binary>.md` shows |
| Reordering a binary's flags because another binary takes them in a different order | For some binaries a misplaced flag is read as the prompt itself | Follow the exact order in `references/<binary>.md` |
| Believing `commands.test` under `codex exec -s read-only` proves behavior | That sandbox blocks the check from actually running — see [references/codex.md](references/codex.md) | Treat that reviewer's pass as reading and reasoning only, not behavioral validation |

## Example

A `code` stage lists `["claude", "cli:codex/<slug>"]`. For the `codex` entry: resolve `codex` on `PATH`; assemble the prompt from `.agents/code-reviewer.md`'s body, `code-review`'s skill text, the diff, and the anti-delegation clause; write it to disk; invoke `codex exec --skip-git-repo-check -s read-only -m <slug> - < prompt.md`; read the reply for `Review status:` under `## Verdict`. If found, that is the verdict. If the process exits 0 with no `## Verdict` section, the reviewer is not-run, reason "no verdict signal", and the round is not spent. Either way, the frozen `git status --porcelain -uall` path list and per-path `git hash-object` are recomputed and compared against the pre-call state before the result is reported.

## Related

Roles: coordinator, document-validator, code-reviewer, security-reviewer. Commands: `/dev-harness:review`, `/dev-harness:discover`, `/dev-harness:secure`. Skills: code-review, document-review, security-review (own the verdict formats this skill reads), project-onboarding (`.harness/local.yaml` is written on first CLI reviewer invocation, not by setup).

## Proof case

Given a `code` stage configured with one Claude reviewer and one `cli:codex/<slug>` reviewer, the `codex` call is assembled from the same role and skill text the Claude reviewer receives, run with `-s read-only`, and its reply is read for `Review status:` under `## Verdict` rather than trusted on exit code alone. Given the same call returning exit 0 with a reasoning preamble and no `## Verdict` section, it is reported not-run with that reason, spends no correction round, and is never read as approval. Given a call that edited a file in the workspace despite `-s read-only`, its verdict is discarded, the edit is reverted, and it is reported not-run for that reason instead.
