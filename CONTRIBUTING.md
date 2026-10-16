# Contributing

## Sources and package

The canonical sources are `.agents/`, `.skills/`, `.commands/`, `templates/` and `profiles/`. Edit these, never the generated package.

`dist/` comes from `go run ./cmd/dh build` and is never edited by hand. Run the build after a source change and commit the regenerated `dist/` with it.

## Checks

Run all three before you open a pull request:

```bash
go test ./...
go run ./cmd/dh validate .
go run ./cmd/dh build
```

After the build, `git diff dist/` must show only the files the build regenerates.

## Commits

Use Conventional Commits in English (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`). Keep the message short and do not mention AI in commits, pull requests or release notes.

## Rules file

`AGENTS.md` is the canonical rules file of the repository. It is written in Portuguese (pt-br). `CLAUDE.md` only imports it.

## Tools

The kit is a plugin for Claude Code. Its marketplace file is `.claude-plugin/marketplace.json`.
