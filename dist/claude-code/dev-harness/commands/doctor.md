---
name: doctor
description: Diagnose kit, runtime and project readiness
author: malka
argument-hint: ""
metadata:
  roles: [coordinator]
  skills: [project-onboarding]
  writes: none (read-only)
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `.harness/MEMORY.md` if present before anything else.

## Routing
Run the diagnosis yourself; do not dispatch a specialist. Load the kit skill `project-onboarding` and follow its diagnosis mode only, never its configuration mode. When a shell is available, run the bundled `bin/dh doctor` from the loaded plugin directory (the `--plugin-dir` path) or from the kit checkout `bin/dh doctor`, and incorporate its output as the mechanical inventory and environment probe. If the script is missing or cannot run, say so and continue with the observational diagnosis.

## Prerequisites
None beyond a resolvable kit package. List the agents, commands and skills actually discovered by name from the agents, skills and templates bundled with this kit; if the kit location cannot be resolved, report the kit as unverified and stop. Ask the owner nothing: this command only observes.

## Output
The `project-onboarding` diagnosis report: kit inventory (commands, roles, skills found by name, unresolved references), the environment table by layer with command or observation, status and evidence, capability impact, whether `.harness/` and `.harness/project.yaml` exist, which tools are available, limitations and next step. Include the `bin/dh doctor` status when it ran. Write no files. When `.harness/MEMORY.md` already exists, record that the diagnosis ran and what was found missing.

## Limits
- Never install, upgrade or configure anything; report only what was observed.
- Never report a layer as present without the observed output that shows it; unknown stays unknown.
- Do not create `.harness/` or any record here; `/dev-harness:setup` owns that.
- Only the Coordinator writes `.harness/MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dev-harness:setup` when `.harness/` or the project profile is missing; otherwise the work command for the task at hand.
