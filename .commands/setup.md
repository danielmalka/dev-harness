---
name: setup
description: Propose the project profile and initialize records
author: malka
argument-hint: "[target directory or setup notes]"
metadata:
  roles: [coordinator, repo-scout]
  skills: [project-onboarding]
  writes: .harness/project.yaml and missing .harness/ records only
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Target and notes: $ARGUMENTS (default to the current working directory and say so). Dispatch `repo-scout` with the Agent tool, model `haiku`, to contribute the read-only inventory part of the kit skill `project-onboarding`: stack, manifests, existing project instructions, discoverable test, lint, build and run commands with their source. Then load `project-onboarding` yourself and produce the profile from that inventory.

## Prerequisites
A readable project root and the existing project instructions, which outrank anything the kit proposes. Writing requires authorized scope: without it, propose the profile in the reply and do not write. Ask the owner only when a conflict between the project instructions and the proposal changes what gets recorded.

## Output
The `project-onboarding` report with the proposed profile written to `.harness/project.yaml`, and `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md` initialized from the templates bundled with this kit only when they are absent and writing is authorized. Re-read every created file and report Created, Left untouched and Not applied. The profile records `language` (`en` by default, `pt-br` when the owner wrote in Portuguese in this session), which selects `templates/<lang>/` for the records. Record the profile source, the language and any unresolved conflict in `.harness/MEMORY.md`.

## Limits
- Never overwrite or edit an existing `.harness/` record; an existing file is left untouched and reported.
- When no shipped profile matches, the profile is "none, custom"; leave unsupported fields unset rather than guessing a command.
- Never install, upgrade or configure tooling, and never modify the project's own instruction files.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:understand` to map the first area, or `/dh:discover` when the request is still an idea.
