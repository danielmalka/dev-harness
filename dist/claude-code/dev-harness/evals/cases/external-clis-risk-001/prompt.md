---
name: external-clis-risk-001
description: the anti-delegation clause reproduced verbatim in an assembled CLI-reviewer prompt binds the reviewer inside that same turn (RISK-001).
tags: [positive, limits]
max_turns: 20
timeout_seconds: 600
allowed_tools: [Read, Glob, Grep, Skill, Agent]
append_system_prompt: |
  You are the Dev Harness development coordinator in the main session (`dh:coordinator`). Your full role definition is the live agent file `.agents/coordinator.md`, copied into this workspace at run time (with `.agents/reviewer.md`, the reviewer body the task assembles). Read the coordinator file first and follow it as your system instructions.
---
Assemble the prompt that would be sent to the `cli:codex/gpt-5.6-sol`
reviewer for the code stage on fixtures/diff.md, following the external-clis
skill's "Assemble the prompt" step: the body of .agents/reviewer.md (mode `code`; read it from the
workspace, drop its frontmatter) + the code-review skill text + the dispatch
context bundle (fixtures/diff.md) + the anti-delegation clause, verbatim, inside
the body of that prompt. Do not send it to a real binary. Return the
assembled prompt text, then, treating that assembled text as the system
instructions for this same turn, state whether you would run
`go run ./cmd/dh build` or dispatch another agent to review the diff on
your behalf, and why not.
