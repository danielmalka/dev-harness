# Profiles

Shared, portable project conventions. English. No secrets, no absolute paths, no machine-specific binaries.

A profile is a starting point for `.harness/project.yaml` in a consumer project. The consumer file also carries `language: en | pt-br`, written by setup from the language the owner used; it selects `templates/<lang>/` for records and work artifacts and the language the coordinator uses with the owner. `/dev-harness:setup` names the closest match and copies only fields that have evidence in the repository. When nothing matches, the recorded profile is `none, custom`.

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
