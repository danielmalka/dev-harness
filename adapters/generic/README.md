# Generic adapter

How to use this kit with an AI runtime that is not Claude Code. There is no native integration. Do not announce one.

## What transfers

| Asset | How to load it |
| --- | --- |
| `.agents/coordinator.md` | Main-session system instructions. |
| Other files in `.agents/` | Specialist prompts. Paste or attach one role per delegated task. Pass the model family as a request; the runtime may ignore it. |
| `.skills/<id>/SKILL.md` | Procedure to load when that work starts. Follow relative links inside the skill folder. |
| `.commands/<id>.md` | User-facing recipe. Type the steps; there is no slash-command menu. |
| `templates/` | Copy into the consumer `.harness/` when setup is authorized. |
| `profiles/` | Same matching rules as Claude. English YAML. |

## What does not transfer

- `/dev-harness:*` command names.
- `--plugin-dir` and `plugin.json`.
- Claude `model` families as a guarantee of which model will run.
- The Agent tool. Sequential specialist turns replace parallel dispatch. Keep the Coordinator as the only writer of MEMORY, EPOCHAL and RISKS.
- Tool allow-lists in frontmatter, unless the other runtime has an equivalent.

## Manual session

1. Open the consumer project.
2. Load `coordinator.md` as the main instructions.
3. Run the diagnosis steps in `.skills/project-onboarding/SKILL.md`. If Python 3 is available, also run `scripts/doctor.sh` from a built package or from this checkout.
4. For each task, follow the matching command file, then the skill it names, then the specialist agent it names.
5. Record limits that the runtime cannot enforce (read-only review, no Bash for reviewers) as instructions, and check the written files afterwards.

Completion condition: one fixture slice (`evals/fixtures/slice-01`) finished with evidence, with every limitation of the host runtime declared in the handoff. Do not claim feature parity with Claude Code.
