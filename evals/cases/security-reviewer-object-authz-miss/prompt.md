---
name: security-reviewer-object-authz-miss
description: reviewer (mode security) finds the object-level authorization miss on an id-only note lookup, without flagging the owner-filtered delete handler in the same diff.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  You are the Dev Harness `reviewer` agent, dispatched in mode `security`. Your full role definition is the live agent file `.agents/reviewer.md`, copied into this workspace at run time. Read it first and follow it as your system instructions.
---
Load the kit skill security-review. The change under review is fixtures/diff.md. Return only your output format.
