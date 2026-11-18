# Adapter do Claude Code

Mapeia as fontes canônicas do kit para o layout de plugin do Claude Code. O mapeamento é uma cópia mais um manifesto gerado. Ele não reescreve corpos de prompt.

Caminho de carga documentado para a v1:

```bash
claude --plugin-dir "/absolute/path/to/clone/dist/claude-code/dev-harness" --agent coordinator
```

Quando o plugin é carregado por nome, qualifique o agente como `dh:coordinator` se o runtime exigir. Os comandos aparecem como `/dh:<name>`.

Não apresente `.agents/`, `.commands/` ou `.skills/` como diretórios que o Claude descobre sozinho.

## Fonte para pacote

| Fonte canônica | Caminho no pacote | Descoberta em runtime |
| --- | --- | --- |
| `.agents/<id>.md` | `agents/<id>.md` | Subagente. Frontmatter `name`, `description`, `model`, `tools`, `color`. |
| `.commands/<id>.md` | `commands/<id>.md` | Comando de barra `/dh:<id>`. Frontmatter `description`, `argument-hint` opcional. |
| `.skills/<id>/` | `skills/<id>/` | Skill. O `name` do `SKILL.md` deve ser igual ao nome da pasta. |
| `templates/` | `templates/` | Não é descoberto automaticamente. Referenciado por agentes, skills e pela configuração. |
| `profiles/` | `profiles/` | Não é descoberto automaticamente. Referenciado pela `project-onboarding`. |
| `cmd/dh`, `internal/` (Go) | `bin/<os>_<arch>/dh` mais o invólucro `bin/dh` | Diagnóstico mecânico, validação, build e instantâneos de sessão. `/dh:doctor` roda `bin/dh doctor` a partir do diretório do plugin. |
| `adapters/claude-code/plugin/settings.json` | `settings.json` | `subagentStatusLine` chamando `bin/dh snapshot subagents`. |
| `adapters/claude-code/plugin/hooks/hooks.json` | `hooks/hooks.json` | Declara só o mod `./panel.ts`, sem hooks de comando. O mod grava o instantâneo da sessão a partir dos eventos do engine (eventos de sessão, prompt, ferramenta, agente e turno); o plugin não chama `dh snapshot event`. Em contas Team, o engine contorna os hooks `classic.*` do plugin, por isso o mod não os usa. |
| gerado | `.claude-plugin/plugin.json` | Identidade do plugin. `name` é `dh`, `license` é `MIT`. |
| gerado | `harness-manifest.json` | Inventário derivado das fontes na hora do build. |
| gerado | `GENERATED.txt` | Marcador de que esta árvore não é uma fonte. |

Build: a partir da raiz do kit, `go run ./cmd/dh build` (Go 1.22+ com a toolchain fixada em `go.mod`; compila para cinco alvos). Ele valida as fontes, monta um pacote completo e então substitui `dist/claude-code/dev-harness`. Edite as fontes e refaça o build. Nunca altere arquivos em `dist/`.

## Mapeamento do frontmatter

| Campo | Agentes | Comandos | Skills |
| --- | --- | --- | --- |
| `name` | Obrigatório, igual ao nome do arquivo sem extensão | Não é usado; o nome do arquivo é o id do comando | Obrigatório, igual ao nome da pasta |
| `description` | Texto de acionamento que o Claude usa para escolher o agente | Subtítulo do menu | Texto de acionamento para carregar a skill |
| `model` | Família obrigatória: `haiku`, `sonnet` ou `opus`. O Coordenador ainda passa o modelo em cada chamada do Agent. | n/a | n/a |
| `tools` | Lista de permissão. Papéis somente leitura omitem Write, Edit e Bash. | n/a | n/a |
| `color` | Dica opcional de UI | n/a | n/a |
| `argument-hint` | n/a | Opcional | n/a |
| `author` | Registrado, não é uma chave de runtime | n/a | Registrado |
| `metadata.provenance` | n/a | n/a | `authorial` ou `adapted` |

"Somente leitura" em linguagem natural não é um sandbox. Papéis somente leitura omitem ferramentas de escrita no frontmatter. O shell ainda pode escrever; esses papéis não recebem Bash.

## O que este adapter não faz

- Fixar uma versão mínima do Claude Code. Esse valor continua `unverified` até a Fase 4.
- Instalar o plugin globalmente ou por um marketplace.
- Traduzir prompts para outro produto.
- Impor travas de sistema de arquivos nos arquivos de `.harness/`.
