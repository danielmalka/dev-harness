# dev-harness

Português: [README.md](README.md)

Portable AI-assisted development kit. Install it in Claude Code and build.

MIT license. Product plan: [`docs/plano-produto.html`](docs/plano-produto.html). Tutorial: [`docs/en/tutorial.html`](docs/en/tutorial.html).

## Install

From the repository's own marketplace, inside a Claude Code session in the project you will work on. The URL form downloads only `marketplace.json`, without cloning the repository:

```text
/plugin marketplace add https://raw.githubusercontent.com/danielmalka/dev-harness/main/.claude-plugin/marketplace.json
/plugin install dh@dev-harness
```

The short form `/plugin marketplace add danielmalka/dev-harness` also works, but it clones the whole repository (about 90 MB, because of the `dh` binaries of every version) and failed once with `fetch-pack: invalid index-pack output` (not reproduced). With a local clone, `/plugin marketplace add /path/to/clone` works too. In every form, `install` downloads only the plugin folder, at the tag pinned in `.claude-plugin/marketplace.json`, which sets the version. For a team, the project can declare the marketplace in `extraKnownMarketplaces` and the plugin in `enabledPlugins` in `.claude/settings.json`, and Claude installs it for whoever trusts the folder.

**Current version: 0.24.1** (tag `v0.24.1`). Install from the `marketplace.json` URL, without cloning the repository. VS Code extension ([`dev-harness-vscode`](https://github.com/danielmalka/dev-harness-vscode), verified in a real VS Code on Remote WSL and Windows): view with local dashboard, status bar with limits, session launcher; install the `.vsix` from the [latest release](https://github.com/danielmalka/dev-harness-vscode/releases/latest) via "Extensions: Install from VSIX...". Local dashboard on http://127.0.0.1:4747 with `/dashboard` (registered by mod), stopped with `dh dashboard --stop [--port N]`; shows the projects registered in `~/.harness/config.yaml`, one progress bar per open PRD, Claude Code sessions with state, per-project metrics on the card (ticket and PRD durations with n and "no date", tickets per PRD, incidents, memory, last consolidation, two charts) and on-demand MEMORY/RISKS reading, avatar with 5h/wk rate limits (70/30 layout: projects on left, avatar and limits on right sticky). Status panel below prompt: context, agents, dh slice progress, `gate pending`, rate limits 5h/wk (🔴 80%+), session cost. Commit traps: deny AI/Claude mention in `git commit`, `gh pr create/edit`, `gh release create`; deny commit when edit after last gate success; now scoped to each commit's repository (no longer blocks edits in another repo). Validation: ticket open under PRD marked `entregue em`, or all tickets done but PRD not. Panel and dh messages follow `.harness/project.yaml` language. Tests: 228 pass. Lean kit (13 agents, 19 commands, 21 skills) and evals out of flow. Tutorial: [`docs/en/tutorial.html`](docs/en/tutorial.html). External CLIs run in disposable clone. PRD with four sections; single ticket replaces story. `/dh:plan-loop` runs waves with CLIs; `/dh:auto` chains discover-plan-build. Roadmap: [`docs/en/roadmap.html`](docs/en/roadmap.html). Per-project documentation (`pdocs`): `/dh:document pdocs` creates `~/.harness/projects/<name>/pdocs/index.html`, the dashboard serves it at `/pdocs/<name>/` and links it from the card. Notes: [`CHANGELOG.md`](CHANGELOG.md).

**Global mode (0.21.0):** if you cannot leave a dh folder in the repository, use `~/.harness/` (`DH_HOME` moves it), created by `dh` and unique per machine and per environment: `config.yaml` with the registry of all projects, `projects/<name>/` with MEMORY, tasks and PRDs, `sessions/` and `dashboard/`. An internal `.harness/` wins when it exists; the memory guard applies in both places. `/dh:setup` picks the mode and registers the project; `dh link` registers or relinks (after moving the repository), `dh harness-path` shows the resolved mode and folder, and `dh projects --json` lists the projects. At work, add `~/.harness` to `additionalDirectories` in `~/.claude/settings.json` (`doctor` only suggests it). The old dashboard environment variables (roots and sprites) were removed; see `CHANGELOG.md` 0.21.0 and the "Global mode" section of the [tutorial](docs/en/tutorial.html).

**Migrating from a pre-0.8.0 install:** the plugin id changed from `dev-harness@dev-harness` to `dh@dev-harness`, and commands from `/dev-harness:<command>` to `/dh:<command>`. Run `/plugin uninstall dev-harness@dev-harness` then `/plugin install dh@dev-harness`; replace `/dev-harness:` with `/dh:` in your own scripts and notes. Marketplace name, repository, Go binary `dh` and snapshot path are unchanged. See `CHANGELOG.md` 0.8.0.

## Load from the clone (development)

Prerequisites: Git and Claude Code authenticated with model access. Python 3 only for the `slice-01` fixture. Go 1.22+ only for kit maintainers.

In the project where you will work:

```bash
claude --plugin-dir "/absolute/path/to/clone/dist/claude-code/dev-harness" --agent coordinator
```

Replace the path with the output of `pwd` in the clone, plus `/dist/claude-code/dev-harness`. Note: inside quotes, `~` does not expand and Claude silently ignores the plugin (`--agent coordinator` fails with "not found"). Use `$HOME/...` or the full path. If the runtime requires the qualified name, use `--agent dh:coordinator`.

Then, in the session:

```text
/dh:doctor
/dh:setup
```

Mechanical diagnostics also run without Claude:

```bash
dist/claude-code/dev-harness/bin/dh doctor dist/claude-code/dev-harness
```

Full tutorial: [`docs/en/tutorial.html`](docs/en/tutorial.html). First-slice fixture: [`evals/fixtures/slice-01`](evals/fixtures/slice-01).

## Sources and package

| Path | Role |
| --- | --- |
| `.agents/` `.commands/` `.skills/` `templates/` `profiles/` | Canonical sources |
| `cmd/dh`, `internal/` | Go binary `dh`: `validate`, `build`, `doctor`, `snapshot`. `go run ./cmd/dh build` generates `dist/claude-code/dev-harness` |
| `dist/` | Generated package. Do not edit. |
| `adapters/claude-code/` | Mapping for Claude Code |
| `adapters/generic/` | Manual use in another AI, without parity |

Work templates in `templates/en/` and `templates/pt-br/` (mirrored versions): PRD, Task/Bug (the ticket), ADR, and RFC, plus the three memory records (MEMORY, EPOCHAL, RISKS). Setup writes `language` to `.harness/project.yaml` and selects the folder. When to use each: [`docs/en/tutorials/templates.html`](docs/en/tutorials/templates.html).

Profiles (`base`, `go-api`, `typescript-web`, `typescript-api`, `php`, `kotlin`, `python`) are in English in `profiles/`.

## Language

- The Coordinator answers in the language you write in. English is the default when there is no signal; write in Portuguese and it answers in Portuguese.
- On the first `/dh:setup`, the session language is written to `.harness/project.yaml` as `language: en` or `language: pt-br`. That field selects `templates/<lang>/` for the memory records and the work artifacts (PRD, task, ADR, RFC). To switch, edit the field or ask the Coordinator; records already written are not translated.
- Everything between agents is English: dispatches, specialist replies, code, commits and eval cases. `project.yaml` and the profiles too.
- HTML documents generated by the `doc-template-html` skill follow the same field: `bash .skills/doc-template-html/scripts/stamp.sh --lang en ...` (default `pt-br`).
- Kit documentation: Portuguese in `docs/`, English in `docs/en/`; every page links to the other version. This README has a Portuguese version at [`README.md`](README.md). The product plan is internal and exists in Portuguese only.

## Maintenance

```bash
go test ./...
go run ./cmd/dh validate --source-only .
go run ./cmd/dh build
go run ./cmd/dh validate .
```
