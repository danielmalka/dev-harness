# Profiles

Shared, portable project conventions. English. No secrets, no absolute paths, no machine-specific binaries.

A profile is a starting point for `.harness/project.yaml` in a consumer project. The consumer file also carries `language: en | pt-br`, written by setup from the language the owner used; it selects `templates/<lang>/` for records and work artifacts and the language the coordinator uses with the owner. `/dh:setup` names the closest match and copies only fields that have evidence in the repository. When nothing matches, the recorded profile is `none, custom`.

## Files

| File | When it applies |
| --- | --- |
| `base.yaml` | Every project. Shared harness rules. |
| `go-api.yaml` | Go modules that expose an HTTP or RPC API. |
| `typescript-web.yaml` | TypeScript apps with a user interface. |
| `php.yaml` | PHP projects managed by Composer, with or without a framework or CMS. |
| `examples/project.yaml` | Shape of a consumer `.harness/project.yaml`. Not a stack profile. |

Python, Electron and other stacks are out of v1 until a real consumer needs them. `php` entered on 2026-09-20 for the first PHP consumer (a Craft CMS site).

## Schema

```yaml
id: kebab-case-id          # must match the filename without .yaml
name: Human label
summary: One sentence
extends: base              # optional; omitted on base
stack:                     # hints used to match a repo; not install instructions
  - go
commands:                  # null means "do not guess"
  test: go test ./...
  lint: null
  build: go build ./...
  run: go run .
conventions:               # bullets the coordinator and specialists must follow
  - ...
limits:                    # hard stops
  - ...
```

Rules:

- `id` equals the filename stem.
- Command values are copied into `.harness/project.yaml` only when the same command is observed in the consumer repo (Makefile, `package.json`, `go.mod` scripts). Otherwise leave the field unset and label it unknown.
- `extends` merges conventions and limits; a child field replaces a parent field of the same name. It does not install tools.
- Profiles never list author-machine paths, tokens, or MCP servers.

## Matching

1. Read the consumer tree (manifests, languages, UI presence).
2. Pick `go-api`, `typescript-web` or `php` when the stack evidence is direct.
3. Otherwise use `base` and record `none, custom` for stack-specific commands.
4. Existing project instructions outrank the profile.

## Consumer-only fields in `.harness/project.yaml`

These fields exist only in a consumer's `.harness/project.yaml`, never in a `profiles/*.yaml` stack profile — `base.yaml`, `go-api.yaml`, `typescript-web.yaml` and `php.yaml` never carry them.

`reviewers:` is optional. Absent, the flow is Claude-only, as it is today. When present, it maps a stage to a list of reviewers for that stage:

```yaml
reviewers:                       # optional; absent = flow today, Claude only
  document: [claude, "cli:agy/gemini-3.1-pro-high"]
  code: [claude, "cli:codex/gpt-5.6-sol"]
  security: [claude]
```

Valid stages: `document`, `code`, `security`. `verify` is out of v1.

Each entry in a stage's list is one of exactly three forms — no fourth form:

- `claude` — dispatches the kit's own agent for the stage; the model is already fixed in `.agents/<role>.md`.
- `"cli:<binary>/<slug>"` — a string; an external CLI reviewer with the default timeout of 15 minutes.
- `{reviewer: "cli:<binary>/<slug>", timeout_minutes: <n>}` — a map; the same CLI reviewer with a custom timeout.

`internal/kit/validate.go` does not change for this field: it validates the kit's own source tree (`.agents/`, `.skills/`, `.commands/`, `profiles/`), not a consumer's `.harness/project.yaml`, so it has nothing to check here.
