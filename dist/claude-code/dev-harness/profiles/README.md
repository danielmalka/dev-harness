# Profiles

Shared, portable project conventions. English. No secrets, no absolute paths, no machine-specific binaries.

A profile is a starting point for `.harness/project.yaml` in a consumer project. The consumer file also carries `language: en | pt-br`, written by setup from the language the owner used; it selects `templates/<lang>/` for records and work artifacts and the language the coordinator uses with the owner. `/dh:setup` names the closest match and copies only fields that have evidence in the repository. When nothing matches, the recorded profile is `none, custom`.

## Files

| File | When it applies |
| --- | --- |
| `base.yaml` | Every project. Shared harness rules. |
| `go-api.yaml` | Go modules that expose an HTTP or RPC API. |
| `typescript-web.yaml` | TypeScript apps with a rendered user interface. |
| `php.yaml` | PHP projects managed by Composer, with or without a framework or CMS. |
| `kotlin.yaml` | Kotlin services on the JVM (Gradle or Maven), with Spring Boot, Quarkus or Micronaut sections gated by build evidence. |
| `python.yaml` | Python services and packages, with Django and FastAPI sections gated by dependency evidence. |
| `typescript-api.yaml` | TypeScript or Node services without a rendered UI. |
| `examples/project.yaml` | Shape of a consumer `.harness/project.yaml`. Not a stack profile. |

Electron and other stacks are out until a real consumer needs them. `php` entered on 2026-09-20 for the first PHP consumer (a Craft CMS site). `kotlin`, `python` and `typescript-api` entered on 2026-10-09 for the owner's work projects, the first Kotlin, Python and TypeScript backend consumers.

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
2. Pick `go-api`, `typescript-web`, `php`, `kotlin`, `python` or `typescript-api` when the stack evidence is direct.
   - `package.json` with a rendered UI present (UI framework dependency such as react, vue, svelte, next, nuxt or angular, `index.html`, or a `public/` folder): `typescript-web`. A service without a UI: `typescript-api`.
   - Monorepo with a UI and an API: one rule. `setup` run at the root picks `typescript-web` and records in `notes:` of `project.yaml` the API folder and the profile `typescript-api`; `setup` run inside the API folder picks `typescript-api`.
3. Otherwise use `base` and record `none, custom` for stack-specific commands.
4. Existing project instructions outrank the profile.

## Consumer-only fields in `.harness/project.yaml`

These fields exist only in a consumer's `.harness/project.yaml`, never in a `profiles/*.yaml` stack profile — `base.yaml`, `go-api.yaml`, `typescript-web.yaml`, `php.yaml`, `kotlin.yaml`, `python.yaml` and `typescript-api.yaml` never carry them.

`reviewers:` is optional. Absent, the flow is Claude-only and there is no other CLI to rotate. When present, it maps a stage to a list of reviewers for that stage. On the `code` stage the first entry of a task is `claude` when that list includes it, otherwise the first entry; each later entry of the same task advances one entry. The `security` stage still runs its full list. The full `code` list runs only when the owner asks for that stage.

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
