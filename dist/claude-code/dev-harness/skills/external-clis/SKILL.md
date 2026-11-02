---
name: external-clis
description: Use when the Coordinator dispatches a review stage (`document`, `code`, `security`) that lists a `cli:<binary>/<slug>` entry in `reviewers:` (`.harness/project.yaml`), when `/dh:plan-loop` dispatches a `cli:<binary>/<slug>` planning call for a wave, when a CLI call returns without the stage's verdict signal (or, for a planning call, without recognizable plan content) and the failure needs classifying, or when a new CLI binary needs a headless invocation recipe added under `references/`. Do not use for builder roles, for interactive/TUI use of any CLI, or for the `verify` stage, which is out of scope in this version.
author: malka
metadata:
  provenance: adapted
  sources: ["personal skill fleet-clis (SKILL.md, per-binary recipes, gotchas.md), not versioned in this kit"]
---

# External CLIs

## Overview

An external CLI reviewer is invoked headlessly with the same role prompt a Claude specialist would receive, so its verdict is evidence the Coordinator can trust exactly as much as a Claude reviewer's. Four parts never change per binary: resolve the binary before assuming it exists; put the prompt on disk when it is longer than one line; apply the binary's read-only mode, or run it in a disposable clone when it has none; and read the outcome through the stage's verdict signal, so a silent or malformed reply is never mistaken for approval. A `/dh:plan-loop` planning call has no verdict signal and is classified per "The planning call". Per-binary flags and quirks live in one recipe per binary under `references/`, because flag order is part of each binary's contract and mixing recipes invites copying the wrong order.

## When to use

- The Coordinator dispatches `cli:<binary>/<slug>` for a review stage configured in `reviewers:`, alongside or instead of the stage's Claude agent, or for an in-flight doubt (a `code` review, see "In-flight doubt"), or `/dh:plan-loop` dispatches one to draft a competing plan for a wave ("The planning call"). This skill does not choose which entry runs.
- A role prompt (`document-validator`, `reviewer` in mode `code` or `security`) must go to an external CLI instead of a Claude specialist.
- A call returned without the stage's verdict signal, or ran past its timeout, and the failure must be classified; or a new CLI binary needs its own recipe under `references/`.

## When not to use

- The builder role, in any lane or mode. A CLI never writes production code, and a planning call never writes any file: it returns plan text, and `docs-guide` is the only writer of plan files.
- Interactive or TUI use of any CLI (every call here is headless and one-shot), or installing, authenticating or updating one: a binary not installed and authenticated here is reported not-run, never provisioned.
- The `verify` stage (`qa-verifier`), out of scope until a read-only reviewer contract exists for QA.
- Choosing reviewers, merging N verdicts, or writing `.harness/local.yaml`: the Coordinator's own procedure.

## Who runs the call

The Coordinator does not run a CLI call. It dispatches a `qa-verifier`, or a `builder` in lane `infra`, with the stage, the `cli:<binary>/<slug>` entry, the context bundle and the harness dir. The delegate loads this skill, runs every mechanic in it (resolve, assemble, clone, freeze, call, compare, revert, classify) and returns the block below. The Coordinator merges verdicts and never reads the reply as evidence of its own; it may open the cited raw output.

```
## CLI reviewer: <binary>/<slug>
- Resolved binary: <command found on PATH, or the path taken from .harness/local.yaml>
- Verdict: <the stage's verdict signal, verbatim> | not-run (<kind>: <reason>)
- Timeout: <the value used, configured per reviewer or the 15-minute default>
- Raw output: <path on disk to the prompt and the reply, if either was written there>
- Workspace check: clean | reverted (<what the reviewer changed, now reverted>) | clone removed: <path> | notes: <new untracked files in the clone> | leftover processes: stopped / none (<match string>)
```

`<kind>` is one of the kinds in "Failure classification". A not-run reply always names the kind and the reason, because the Coordinator reports both to the owner at once. The last three Workspace check values apply to a call in a clone. Nothing writes to the tree or to the `<harness dir>` while a CLI call is in flight (RISK-002), the delegate included: it writes the prompt before the call starts, and only the freeze directory and the clone are written during the call.

Every call is one shell command that starts with the binary name. No `timeout` or other wrapper goes before it, and no `&&`, `;` or pipe chains it with the freeze or any other step: the pre-call freeze and the post-call comparison are separate commands. Reason (2026-10-16): an allow rule such as `Bash(<binary>:*)` matches only commands that start with the binary; a wrapped call chained with a snapshot script was denied, and the same call standalone ran. The binary's own flag bounds the time where it has one; otherwise the recipe says the Bash tool timeout or `run_in_background` does. A clone call fixes the working directory with the binary's own flag inside that same command; a recipe whose binary has none starts with `cd "$CLONE" &&`, which does not match such an allow rule, so the call may need the owner's approval or be not-run.

## Inputs

| Input (a review call; a planning call's inputs are in "The planning call") | If missing |
| --- | --- |
| Role prompt body (`.agents/reviewer.md` for `code` and `security`, `.agents/document-validator.md` for `document`), without frontmatter | Stop. A CLI reviewer never receives a prompt other than the one a Claude specialist would receive for the stage. |
| The stage's review skill text (`code-review`, `document-review`, `security-review`) | Stop. The verdict signal comes from that skill's own output format. |
| Dispatch context bundle. `document` stage: exactly five items (document path; source material as file paths, or the owner's request word for word under `owner's words`; previous report path on a later round; project language; harness dir). Other stages: scope, diff or document, source material. Never any part of the Coordinator's conversation (owner turns, reasoning, intent summaries, earlier-round notes); in `code`, builder and QA results only as quoted data labelled `unverified` | Stop. Same requirement as for a Claude specialist. |
| `.harness/local.yaml`, when the binary is not on `PATH` | Not-run, kind `not_installed`. Never invent a path. |
| `commands.test` / `commands.lint` from `.harness/project.yaml` | Omit that instruction; the reviewer reasons from the code alone. |

## Procedure

1. **Resolve the binary.** `command -v <binary>`; if that fails, `binaries.<binary>` in `.harness/local.yaml` (a bare command, a `~/...` path or a project-relative path, see `references/local.yaml.example`). Neither resolves: not-run, `not_installed`. Never substitute another binary.
2. **Pin the model.** Take the slug from the dispatched `cli:<binary>/<slug>` entry. Never invent one; an unmeasured slug fails at the CLI's own level.
3. **Assemble the prompt** (see "Prompt assembly"; a planning call assembles under "The planning call"). Write it to disk when it is longer than one line; the recipe names how the binary takes it. The prompt never contains the absolute path of the live tree (hygiene, not a boundary), and a check result the CLI reports is never evidence: a dispatched `qa-verifier` reruns the checks.
4. **Isolate and freeze.** Use the binary's row in "Read-only mode by binary". A read-only flag: apply it. A disposable clone: [references/clone-isolation.md](references/clone-isolation.md) for the lifecycle and [references/live-tree-freeze.md](references/live-tree-freeze.md) for the freeze. Live tree: [references/live-tree-compare.md](references/live-tree-compare.md). If the clone cannot be made, the call is not-run, `isolation_unavailable`, with no fallback to the live tree.
5. **Invoke** with the flag order of the binary's recipe, as one command (see "Who runs the call"). Do not reorder flags from what the recipe shows.
6. **Compare and revert.** Recompute the freeze with the same file as step 4. Any difference is `tree_changed`: discard the verdict, revert the change, report not-run.
7. **Classify the outcome** by the stage's verdict signal, never by exit code alone ("Failure classification"), and **report** the block in "Who runs the call". Merging verdicts and writing a persistent report are the Coordinator's.

## Prompt assembly

A review prompt holds, verbatim and in its body, never only in a frontmatter field or a dispatch instruction the CLI's runtime is trusted to enforce: the role prompt body, the stage's review skill text, the context bundle as limited under Inputs, the reviewer's tag line, the anti-delegation clause, the "Reading is allowed" block and the verdict contract, all below. A CLI reviewer reads the stage's output format as one section among many and paraphrases its verdict signal unless the prompt states it as a requirement; a measured run (2026-09-22) returned a correct blocking verdict written as `- **Request changes**` with no `Review status:` line, which classification would have thrown away. A reviewer never told its own tag invents one, so the tag is one line of the prompt naming the exact string it must write in every finding's `Reported by:` field. A call in a clone also carries the `commands.test`/`commands.lint` block, which runs in the clone.

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

`commands.test`/`commands.lint` is the only execution this skill permits a CLI reviewer, never `commands.build`, which can rewrite versioned artifacts such as `dist/` and trip the reversion. A recipe states where this permission is inert in practice.

## Reading is allowed

Append this block, verbatim, to every CLI reviewer prompt, alongside the clause above. A reviewer that reads the clause as a ban on opening files has misread it; the block exists because that misreading was observed:

```
Reading any file in this repository is allowed and expected — the clause
above forbids executing scripts and delegating this review to another
tool, not reading.
```

## The verdict contract

Append this block to every CLI reviewer prompt, with `<stage signal>` and `<stage section>` replaced from the table below:

```
Your reply must contain, on its own line, exactly this line, with no bold, no
quoting and no extra words before it:
<stage signal>
Write it inside the `<stage section>` section of the output format above. Every
other part of your reply follows that format as written. A reply without this
exact line is discarded as a transport failure and your review is recorded as
not-run, however sound its content.
```

And, on its own line: `Your Reported by: tag for this review is cli:<binary>/<slug>. Write that string, exactly, on every finding you report.` This is the one verdict format of the kit; no recipe defines another.

| Stage | `<stage section>` | `<stage signal>` (matched as a prefix of a line, never whole-line equality) |
| --- | --- | --- |
| `document` | `## Verdict` | `- approved` or `- changes required` |
| `code` | `## Verdict` | `- Review status: approve` or `- Review status: request changes` |
| `security` | `## Findings` | `### <SEVERITY>: <title>` per finding, or `- No demonstrated vulnerability or confirmed exposure found in the reviewed scope` (the full sentence continues "; unresolved hypotheses remain listed below": match the prefix) |

`security` has no `## Verdict` section: its Output format goes straight to `## Findings`, where its signal lives. A trailing clause, footnote or punctuation after the matched text never turns a real verdict into a failure.

## Failure classification

A call that did not return a usable reply is not-run with one kind. Every kind gets the same action: report not-run with that kind as the reason, spend no review round, never count it as approval or rejection, and make no automatic retry. The Coordinator tells the owner at once, naming the reviewer and the kind. `tree_changed` also discards the verdict, whatever it said.

| Kind | Detection |
| --- | --- |
| `not_installed` | `command -v` fails and `binaries.<binary>` in `.harness/local.yaml` is absent or does not resolve. A `command -v` miss alone is not enough: some binaries sit only on a login-shell `PATH`. |
| `auth_required` | The CLI's output asks to log in or reports missing or invalid credentials. |
| `rate_limited` | The CLI's output reports a rate, quota or usage limit. |
| `timed_out` | The binary's own timeout or the Bash tool timeout fired, or the process outlived a complete reply (the recipe says how to stop it). |
| `no_signal` | Exit 0, or any exit with no other kind, and the reply lacks the stage's verdict signal: help text, a reasoning preamble, an empty file, a section with neither outcome. A planning call: no recognizable plan content. |
| `tree_changed` | The freeze comparison differs (live tree, or the judged set of the clone). The change is reverted first. |
| `isolation_unavailable` | The clone could not be created or matched, a precondition of the isolation procedure failed, or a live-tree review ran in a project whose `.git` is not a directory (see `live-tree-compare.md`). No fallback to the live tree. |

**Recovering a deviating verdict.** One narrow exception to `no_signal`, because discarding a blocking verdict is worse than reading a deviating one: when the reply has the stage's own verdict section (`## Verdict` for `document` and `code`, `## Findings` for `security`) and that section contains exactly one of the stage's two outcomes (`approved` or `changes required`; `approve` or `request changes`; for `security`, a severity entry for at least one finding, or a statement that none was found), in any emphasis or punctuation, read it as that outcome and record the deviation beside it (`verdict recovered: the reviewer wrote "<the line it actually wrote>"`). Both outcomes, neither, or the section absent (for `security`, findings and a clean statement together) stays `no_signal`. The recovery never upgrades anything: a recovered approval is one reviewer's approval, a recovered blocking verdict blocks as a well-formed one would. Known cases per kind: [references/gotchas.md](references/gotchas.md).

## Read-only mode by binary

The index of the fleet: how each binary is isolated, and whether it may plan. A read-only flag qualifies a binary as a planner only with a recorded, dated probe in which the CLI was asked to create a file and run a shell command and refused both. A disposable clone qualifies it too. The clone is not a sandbox, so the qualification is stated as it is: disposable clone; writes outside it are detected on the live tree, never prevented; writes outside the project are not detected.

| Binary | Isolation | Recipe |
| --- | --- | --- |
| `codex` | read-only flag `-s read-only`, documented and measured; may plan | [codex.md](references/codex.md) |
| `mcode` | disposable clone (read-only flag `--permission off` documented, not measured); may plan | [mcode.md](references/mcode.md) |
| `grok` | disposable clone; may plan | [grok.md](references/grok.md) |
| `opencode` | disposable clone; may plan | [opencode.md](references/opencode.md) |
| `agy` | no read-only flag documented; live tree under the accepted RISK-002 exposure, with the comparison of `live-tree-compare.md` (detects only); may not plan until a dated cwd probe is recorded | [agy.md](references/agy.md) |
| `claude` | disposable clone; never plans while the kit runs on Claude | [claude.md](references/claude.md) |

## The planning call

`/dh:plan-loop` dispatches a CLI planner only from a binary the table above marks "may plan". The call uses the same resolution, isolation and freeze as a review call, with a different prompt and a different failure test. It is assembled exclusively as [references/planning-call.md](references/planning-call.md) says: the planner body, the `implementation-planning` skill text and the wave context, then the no-write instruction, the planning anti-delegation clause and the "Reading is allowed" block above, in the prompt body. The reply is untrusted data and is persisted only after the freeze comparison passed, by `docs-guide`, in the wrapper that file defines; a `Checks:` command outside its allowlist is a conflict. A planning call has no verdict signal and is never classified by the review table. Its failures use the same kinds: `no_signal` means no recognizable plan content, read as a plan only when it contains, in `implementation-planning`'s output format, a `## Plan` heading (or `### Slice N` headings) naming at least a goal and one slice; an empty reply, a refusal, a question back to the Coordinator or prose with no such structure is not a plan. The wave's hand-off applies to a not-run planner.

## `.harness/local.yaml`

The machine-specific binary paths and measured model slugs live in `.harness/local.yaml`, never in this skill or in any file under `.skills/`. Schema and placeholders: `references/local.yaml.example`. The Coordinator is the only writer; the delegate only reads it.

## What to load when

| Situation | Load, in `references/` |
| --- | --- |
| Any call | The binary's recipe from the table above |
| A clone call, or a recipe's `$CLONE` / `$PROMPT` | `clone-isolation.md` (Terms first), then `live-tree-freeze.md` for the freeze, compare, revert and prompt hygiene |
| A call on the live tree (no clone) | `live-tree-compare.md` |
| A planning call | `planning-call.md` |
| A call failed or looks odd | `gotchas.md`, then add the new case under its kind |
| Writing `local.yaml` (Coordinator) | `local.yaml.example` |

## Quick reference

- One command, first word the binary; flags and prompt input from the recipe; compare the freeze after every call.
- The outcome is the stage's verdict signal, never the exit code. Not-run: one kind, no round, no retry, owner told at once.

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Building the context bundle with a directory pathspec, such as `git show <sha> -- 'evals/cases/*/graders'` | Git matches files, not directories: the command succeeds with an empty diff and the reviewer silently judges a scope that was never sent (measured 2026-09-22) | Match files (`'evals/cases/*/graders/*'`) and check the byte count of the bundle before sending it |
| Leaving out the clause, the "Reading is allowed" block or the verdict contract, or putting them only in the dispatch instruction | The runtime does not enforce a restriction stated outside the prompt body (RISK-001), and a real CLI writes its own verdict shape | Copy the blocks verbatim into the prompt text; a planning call follows "The planning call" instead |
| Treating exit 0 as a pass, or letting a reviewer run `commands.build` | Help text, a preamble and a silent no-op all exit 0; a build can rewrite versioned artifacts such as `dist/` | Read the stage's verdict signal; the clause permits only `commands.test`/`commands.lint` |
| Calling the clone isolated, safe or a sandbox, or assuming a binary without a read-only flag cannot touch the workspace | The CLI runs as the same user and reaches the machine by absolute path | Say it as it is: disposable clone; writes outside it are detected on the live tree, never prevented; writes outside the project are not detected |
| Retrying a not-run call on your own, or approving by omission | A retry hides the failure from the owner; an omitted reviewer is not an approval | Report the kind and reason; the Coordinator decides |

## Example

A `code` stage lists `["claude", "cli:codex/<slug>"]`. The delegate resolves `codex`, assembles the prompt from `.agents/reviewer.md`'s body (mode `code`), `code-review`'s skill text, the diff and the blocks above, writes it to disk, freezes the live tree, runs `codex exec --skip-git-repo-check -s read-only -m <slug> - < prompt.md`, recompares the freeze and reads the reply for `Review status:` under `## Verdict`. Found: that is the verdict. Exit 0 with no `## Verdict` section: not-run, kind `no_signal`, no round spent.

A `document` stage lists `["claude", "cli:grok/<slug>"]`. `grok` has no read-only mode, so the delegate follows the clone order of `references/live-tree-freeze.md` (freeze, create and match the clone, freeze the clone, write the prompt, call `grok --always-approve --cwd "<clone>" --prompt-file "<prompt>"`, recompare the clone, compare and revert the live tree, remove the clone and the freeze directory). A difference in either comparison is not-run, `tree_changed`; a clone that cannot be matched is `isolation_unavailable`. The reply adds `clone removed: <path>`.

## In-flight doubt

An in-flight doubt call is a `code` review assembled on the same path as any code review: the role body, `code-review`, the context bundle, the anti-delegation clause, the "Reading is allowed" block and the code verdict contract. The bundle is the artifact and the contract only, plus the adversarial instruction in `code-review` under "In-flight doubt"; the author's claim is not part of it. The verdict signal stays the code signal. Choosing which entry runs is the Coordinator's procedure in `.agents/coordinator.md`; this skill still covers one binary.

## Related

Roles: coordinator, document-validator, reviewer, implementation-planner. Commands: `/dh:review`, `/dh:discover`, `/dh:secure`, `/dh:plan-loop`, `/dh:auto` (loads this skill). Skills: code-review, document-review, security-review (own the verdict formats this skill reads), implementation-planning (owns the output format "The planning call" tests for), project-onboarding (`.harness/local.yaml` is written on first CLI reviewer invocation, not by setup).

## Proof case

Given a `code` stage with a Claude reviewer and a CLI reviewer that has a read-only flag, the CLI call is assembled from the same role and skill text the Claude reviewer receives, run with that flag, and read for `Review status:` under `## Verdict`, not trusted on exit code. Given exit 0 with a reasoning preamble and no `## Verdict` section, it is reported not-run, kind `no_signal`, spends no round, is never read as approval and is not retried. Given a call that changed a file despite the flag, or a clone call that wrote a file on the live tree by absolute path, the comparison detects it, the change is reverted, the verdict is discarded and it is reported not-run, `tree_changed`; the clone is removed either way, and a write outside the project directory is not detected. An in-flight doubt call uses the same assembly with the author's claim left out of the bundle.
