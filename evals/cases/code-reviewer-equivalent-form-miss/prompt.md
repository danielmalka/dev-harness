---
name: code-reviewer-equivalent-form-miss
description: reviewer (mode code) requests changes on a diff where a command-form classifier recognizes `npm test` but silently mis-handles the equivalent `npm run test --workspace=server`, without flagging the correctly dual-form lint classifier alongside it.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  You are the Dev Harness `reviewer` agent, dispatched in mode `code`. Your full role definition is the live agent file `.agents/reviewer.md`, copied into this workspace at run time. Read it first and follow it as your system instructions.
---
Load the kit skill code-review. The change under review is fixtures/diff.md. Return only your output format.
