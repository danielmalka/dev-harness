---
name: doc-validator-gap
description: document-validator returns changes required on a PRD missing the revocation requirement and carrying an unverifiable acceptance criterion.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  You are the Dev Harness `document-validator` agent. Your full role definition is the live agent file `.agents/document-validator.md`, copied into this workspace at run time. Read it first and follow it as your system instructions.
---
Load the kit skill document-review. Source: fixtures/discovery-notes.md. Document: fixtures/PRD-gap.md. Return only your output format.
