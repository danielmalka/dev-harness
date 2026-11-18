# Contribuindo

## Fontes e pacote

As fontes canônicas são `.agents/`, `.skills/`, `.commands/`, `templates/` e `profiles/`. Edite essas pastas, nunca o pacote gerado.

`dist/` vem de `go run ./cmd/dh build` e nunca é editado à mão. Rode o build depois de mudar uma fonte e faça o commit do `dist/` regenerado junto com ela.

## Checagens

Rode as três antes de abrir um pull request:

```bash
go test ./...
go run ./cmd/dh validate .
go run ./cmd/dh build
```

Depois do build, `git diff dist/` deve mostrar só os arquivos que o build regenera.

## Commits

Use Conventional Commits em inglês (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`). Mantenha a mensagem curta e não mencione IA em commits, pull requests nem notas de release.

## Arquivo de regras

`AGENTS.md` é o arquivo de regras canônico do repositório. Ele é escrito em português (pt-br). `CLAUDE.md` apenas o importa.

## Ferramentas

O kit é um plugin para o Claude Code. O arquivo de marketplace dele é `.claude-plugin/marketplace.json`.
