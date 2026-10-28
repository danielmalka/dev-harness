# dev-harness

Português: [README.md](README.md)

## What it is

A portable AI-assisted development kit, installed as a Claude Code plugin: a Coordinator, specialist agents, skills and `/dh:<command>` commands that take a request from idea to delivery with review and verification. MIT license.

Credit: part of the review flow and the drift report are inspired by the ideas of the Elephant-Goldfish Model (EGM) from Google (https://github.com/google/elephant-goldfish-model, Apache-2.0 license, "not an officially supported Google product"); the kit reimplements the mechanisms in its own words and copies no text.

## Requirements

Claude Code, authenticated with model access; Git.

## Install

Inside a Claude Code session, in the project you will work on. The URL form downloads only `marketplace.json`, without cloning the repository:

```text
/plugin marketplace add https://raw.githubusercontent.com/danielmalka/dev-harness/main/.claude-plugin/marketplace.json
/plugin install dh@dev-harness
```

The short form `/plugin marketplace add danielmalka/dev-harness` also works, but it clones the whole repository (about 90 MB, because of the `dh` binaries of every version). `install` downloads only the plugin folder, at the tag pinned in `.claude-plugin/marketplace.json`, which sets the version. For a team, the project can declare the marketplace in `extraKnownMarketplaces` and the plugin in `enabledPlugins` in `.claude/settings.json`.

## Quick start

In the session, inside the project, run `/dh:doctor` and then `/dh:setup`. The first diagnoses the environment; the second picks the mode (`.harness/` in the repository or the global folder `~/.harness/`) and registers the project. Then use the work commands (`/dh:discover`, `/dh:plan`, `/dh:build`, `/dh:review`, `/dh:verify`, or `/dh:auto` to chain them). Step by step: [tutorial](docs/en/tutorial.html) and [quick start](docs/en/quick-start.html).

The Coordinator answers in the language you write in (English by default). `/dh:setup` writes `language: en` or `language: pt-br` to `project.yaml`; the PRD is always in English.

## Update

The plugin comes from the tag pinned in `.claude-plugin/marketplace.json`. For a new version, update the marketplace (`claude plugin marketplace update dev-harness`) and the plugin (`claude plugin update dh@dev-harness`); check with `/dh:doctor`. What changed in each version is in [`CHANGELOG.md`](CHANGELOG.md).

## Uninstall and remove data

1. In the session: `/plugin uninstall dh@dev-harness`. This removes the plugin; it does not delete your data.
2. Global folder: delete `~/.harness` (or the path in `DH_HOME`) to remove the project registry, the global-mode records, the snapshots and the dashboard token.
3. If you added `~/.harness` to `additionalDirectories` in `~/.claude/settings.json`, remove that entry.
4. In each repository in repository mode: delete the `.harness/` folder (the project's memory, tasks, PRDs and ADRs). The kit never deletes it by itself.

## What the kit writes and where

In global mode, the `~/.harness/` folder (`DH_HOME` moves it) is unique per machine and per environment; an in-repository `.harness/` wins when it exists. Details in the "Global mode" section of the [tutorial](docs/en/tutorial.html).

| Where | What | Written by |
| --- | --- | --- |
| `<repository>/.harness/` | repository mode: MEMORY, EPOCHAL, RISKS, `project.yaml`, tasks, PRDs, ADRs | Coordinator and `/dh:` commands |
| `~/.harness/config.yaml` | project registry and defaults (`language`, `mode`, `worktree`, `dashboard.sprites`) | only `dh` (`dh config`, `dh link`) |
| `~/.harness/projects/<name>/` | global mode: the same records as `.harness/` | Coordinator and `/dh:` commands |
| `<repository>/REVIEW.md` | project review guide, created by `/dh:setup` when missing, in both modes; the only file setup writes in the tree in global mode; never overwrites | `/dh:setup` (via a specialist) |
| `~/.harness/sessions/` | one snapshot per Claude Code session (`dh snapshot prune` deletes those older than 7 days) | plugin mod |
| `~/.harness/dashboard/` | dashboard stop token, per port | dashboard |
| `~/.harness/projects/<name>/pdocs/` | project HTML documentation (`/dh:document pdocs`) | docs-guide |
| `127.0.0.1:4747` | local dashboard, this machine only, when you start it with `/dashboard` | `dh` binary |

With `worktree` set to `always` (or when you ask), `/dh:build` and `/dh:auto` work in a sibling worktree `<checkout>.worktrees/<id>`. When `/dh:plan` leaves open decisions, it generates an HTML approval page in `pdocs`. Details in the "Configuration" section of the [tutorial](docs/en/tutorial.html).

## Turn things off

- Dashboard: `dh dashboard --stop` (`--port N` for another port). It never kills another program's process.
- Plugin: `claude plugin disable dh@dev-harness` (or uninstall, as above); the mod and the guards go with it.
- There is no per-hook switch for the mod or the commit trap: turning them off means disabling the whole plugin.

## If something fails

- Run `/dh:doctor` in the session or `dh doctor` in a terminal; it points at what is missing.
- Port busy: the dashboard warns; use `/dashboard --port N` and `dh dashboard --stop --port N`.
- Plugin not found: check `/plugin marketplace add ...` and `/plugin install dh@dev-harness`; when loading from a clone, do not put `~` inside quotes in `--plugin-dir` (Claude silently ignores the plugin).
- Still stuck: open an issue on [GitHub Issues](https://github.com/danielmalka/dev-harness/issues). For a vulnerability, see [`SECURITY.md`](SECURITY.md).

## Use together with

Optional; the kit needs neither of them.

- [`danielmalka/malka-claude-mods`](https://github.com/danielmalka/malka-claude-mods): visual Claude Code mods (status line, menu, Zen, Mode) that moved out of this kit in 0.29.0.
- [`akitaonrails/ai-memory`](https://github.com/akitaonrails/ai-memory): long-term memory and handoff across coding CLIs.

## Docs map

- [Quick start](docs/en/quick-start.html) and [tutorial](docs/en/tutorial.html): start here. [Roadmap](docs/en/roadmap.html) and [`CHANGELOG.md`](CHANGELOG.md): what shipped and what is next.
- `docs/prd/`, `docs/adr/` and `docs/discovery/`: design history, in Portuguese.
- `evals/`: measurement history; not needed to use the kit.
- [`docs/plano-produto.html`](docs/plano-produto.html): product plan, Portuguese and internal.
- Portuguese pages: the same topics under Portuguese names in `docs/` (quick start: `docs/inicio-rapido.html`).

## For contributors

This section is only for people who maintain the kit; users do not need it. Full rules in [`CONTRIBUTING.md`](CONTRIBUTING.md). Needs Go 1.22+ (and Python 3 only for the `slice-01` fixture).

```bash
go test ./...
go run ./cmd/dh build
go run ./cmd/dh validate .
```

| Path | Role |
| --- | --- |
| `.agents/` `.commands/` `.skills/` `templates/` `profiles/` | Canonical sources |
| `cmd/dh`, `internal/` | Go binary `dh`. `go run ./cmd/dh build` generates `dist/claude-code/dev-harness` |
| `dist/` | Generated package. Do not edit. |
| `adapters/` | Mapping for Claude Code (`claude-code/`) and manual use in another AI, without parity (`generic/`) |

Load from the clone, in the project where you will work:

```bash
claude --plugin-dir "$HOME/path/to/clone/dist/claude-code/dev-harness" --agent coordinator
```

Use `$HOME/...` or the full path, never `~` inside quotes. If the runtime requires the qualified name, use `--agent dh:coordinator`. Diagnostics without Claude: `dist/claude-code/dev-harness/bin/dh doctor dist/claude-code/dev-harness`. Work templates: [`docs/en/tutorials/templates.html`](docs/en/tutorials/templates.html).
