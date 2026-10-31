---
name: setup
description: Propose the project profile and initialize records
author: malka
argument-hint: "[target directory or setup notes | bootstrap]"
metadata:
  roles: [coordinator, repo-scout]
  skills: [project-onboarding]
  writes: project.yaml and missing records in the harness dir; REPO-MAP.md in the harness dir when the owner accepts; REVIEW.md at the repository root when absent
disable-model-invocation: true
---
## Role
Act as the kit `coordinator` in this session (read the bundled `coordinator` agent definition if this session was not started with it). Read `<harness dir>/MEMORY.md` if present before anything else.

## Routing
Target and notes: $ARGUMENTS (default to the current working directory and say so). Dispatch `repo-scout` with the Agent tool, model `haiku`, to contribute the read-only inventory part of the kit skill `project-onboarding`: stack, manifests, existing project instructions, discoverable test, lint, build and run commands with their source. Then load `project-onboarding` yourself and produce the profile from that inventory. Decide where the records live, in this order: (1) run `dh harness-path --json` (the bundled `bin/dh` (`bin\dh.cmd` on Windows) from the loaded plugin directory, as in `.commands/doctor.md`; every `dh` call here uses it) and read `defaults`; (2) the mode is `defaults.mode` when filled, otherwise ask the owner once for this project (`repo` or `global`); with `defaults` all empty (no `config.yaml` or no keys) use today's defaults; (3) in `repo`, create `<repo>/.harness/` with the records from the templates; (4) run `dh link [name]`, which always registers the project (in `repo` it only registers, so `<home>/projects/<name>/` does not exist afterwards; in `global` it also creates that folder), then paste the `defaults.reviewers` block into the new `project.yaml`; from then on the project's own values win. Never Read `config.yaml`. In `global` the records are born in `<home>/projects/<name>/` and the only file written in the repository tree is `REVIEW.md` (see the review guide step below); remind the owner to add `~/.harness` to `additionalDirectories` in `~/.claude/settings.json` (never write it).

Review guide, in both modes (skill `project-onboarding`): when the project has no review guide file, dispatch a specialist (`builder`, lane `infra`; not the Coordinator, since the file is new) to create `REVIEW.md` at the repository root from `templates/en/REVIEW.md`, then write `review_guide: REVIEW.md` (a repository-relative path, never absolute or with `..`) in `project.yaml`. The file is always English, whatever `language` says. When the file exists but the key is absent, record the key and leave the file untouched. When the key points to an existing file, change nothing. Never overwrite.

Repository bootstrap (opt-in, skill `project-onboarding`): after the profile, offer to write the repository map `<harness dir>/REPO-MAP.md`; the owner accepts or declines. With `bootstrap` in the arguments, skip the profile and refresh only the map. Get the `HEAD` sha from a `qa-verifier` run (`git rev-parse HEAD`) and dispatch `docs-guide` (model `sonnet`, skills `repository-mapping` and `project-onboarding`) with it and the `repo-scout` findings; when it reports more than 15 top-level directories, ask the owner to name the scope. In `repo` mode the file is versioned.

## Prerequisites
A readable project root and the existing project instructions, which outrank anything the kit proposes. Writing requires authorized scope: without it, propose the profile in the reply and do not write. Ask the owner only when a conflict between the project instructions and the proposal changes what gets recorded.

## Output
The `project-onboarding` report with the mode, the registration, the proposed profile written to `project.yaml`, and `MEMORY.md`, `EPOCHAL.md` and `RISKS.md` initialized (all in the harness dir) from the templates bundled with this kit only when they are absent and writing is authorized. Re-read every created file and report Created, Left untouched and Not applied; `REVIEW.md` appears under Created or Left untouched. The profile records `language`: `defaults.language` from `dh harness-path --json` when set; otherwise the owner's language (pt-br when the owner writes Portuguese, else en). It is the owner conversation language and selects `templates/<lang>/` for the other records (the PRD and `REVIEW.md` are always English). Record the profile source, the mode, the language and any unresolved conflict in `MEMORY.md`.

## Limits
- Never overwrite or edit an existing record in the harness dir; an existing file is left untouched and reported. The one exception is `REPO-MAP.md`, derived data that a refresh may overwrite.
- When no shipped profile matches, the profile is "none, custom"; leave unsupported fields unset rather than guessing a command.
- Never install, upgrade or configure tooling, and never modify the project's own instruction files.
- Never edit `.gitignore`; only `local.yaml` (written later, on the first CLI reviewer call) should stay ignored.
- Only the Coordinator writes `MEMORY.md`, `EPOCHAL.md` and `RISKS.md`.

## Next
`/dh:understand` to map the first area, or `/dh:discover` when the request is still an idea.
