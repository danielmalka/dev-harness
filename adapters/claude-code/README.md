# Claude Code adapter

Maps canonical kit sources to the Claude Code plugin layout. The mapping is a copy plus a generated manifest. It does not rewrite prompt bodies.

Load path documented for v1:

```bash
claude --plugin-dir "/absolute/path/to/clone/dist/claude-code/dev-harness" --agent coordinator
```

When the plugin is loaded by name, qualify the agent as `dev-harness:coordinator` if the runtime requires it. Commands appear as `/dev-harness:<name>`.

Do not present `.agents/`, `.commands/` or `.skills/` as directories Claude discovers by itself.

## Source to package

| Canonical source | Package path | Runtime discovery |
| --- | --- | --- |
| `.agents/<id>.md` | `agents/<id>.md` | Subagent. Frontmatter `name`, `description`, `model`, `tools`, `color`. |
| `.commands/<id>.md` | `commands/<id>.md` | Slash command `/dev-harness:<id>`. Frontmatter `description`, optional `argument-hint`. |
| `.skills/<id>/` | `skills/<id>/` | Skill. `SKILL.md` `name` must equal the folder name. |
| `templates/` | `templates/` | Not auto-discovered. Referenced by agents, skills and setup. |
| `profiles/` | `profiles/` | Not auto-discovered. Referenced by project-onboarding. |
| `cmd/dh`, `internal/` (Go) | `bin/<os>_<arch>/dh` plus `bin/dh` wrapper | Mechanical diagnosis, validation, build and session snapshots. `/dev-harness:doctor` runs `bin/dh doctor` from the plugin directory. |
| `adapters/claude-code/plugin/settings.json` | `settings.json` | `subagentStatusLine` calling `bin/dh snapshot subagents`. |
| `adapters/claude-code/plugin/hooks/hooks.json` | `hooks/hooks.json` | Session and subagent hooks calling `bin/dh snapshot event`. |
| generated | `.claude-plugin/plugin.json` | Plugin identity. `name` is `dev-harness`, `license` is `MIT`. |
| generated | `harness-manifest.json` | Inventory derived from sources at build time. |
| generated | `GENERATED.txt` | Marker that this tree is not a source. |

Build: from the kit root, `go run ./cmd/dh build` (Go 1.22+ with the toolchain pinned in `go.mod`; cross-compiles five targets). It validates sources, stages a complete package, then replaces `dist/claude-code/dev-harness`. Edit sources and rebuild. Never patch files under `dist/`.

## Frontmatter mapping

| Field | Agents | Commands | Skills |
| --- | --- | --- | --- |
| `name` | Required, equals filename stem | Not used; filename is the command id | Required, equals folder name |
| `description` | Triggering text Claude uses to pick the agent | Menu subtitle | Triggering text for skill load |
| `model` | Required family: `haiku`, `sonnet` or `opus`. Coordinator still passes the model on each Agent call. | n/a | n/a |
| `tools` | Allow-list. Read-only roles omit Write, Edit and Bash. | n/a | n/a |
| `color` | Optional UI hint | n/a | n/a |
| `argument-hint` | n/a | Optional | n/a |
| `author` | Recorded, not a runtime switch | n/a | Recorded |
| `metadata.provenance` | n/a | n/a | `authorial` or `adapted` |

Natural-language "read only" is not a sandbox. Read-only roles omit write tools in frontmatter. Shell can still write; those roles do not receive Bash.

## What this adapter does not do

- Pin a Claude Code minimum version. That value stays `unverified` until Fase 4.
- Install the plugin globally or through a marketplace.
- Translate prompts for another product.
- Enforce filesystem locks on `.harness/` files.
