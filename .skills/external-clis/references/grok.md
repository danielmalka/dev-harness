# grok (Grok Build)

Binário: `grok`, resolvido por `binaries.grok` em `.harness/local.yaml` quando `command -v grok` falha. Um `PATH` só de shell de login é comum para este binário: `command -v` não o encontra mesmo quando ele está instalado, então não trate essa falha sozinha como `not_installed`. Agente de código com mãos; também revisor e planejador (num clone descartável).

## Invocação

```bash
grok --always-approve --cwd "$CLONE" --prompt-file "$PROMPT"
# prompt curto:
grok --always-approve --cwd "$CLONE" -p "$(cat "$PROMPT")"
```

Um comando de shell cuja primeira palavra é o binário: escreva ali o comando resolvido como literal (`grok` simples, ou o caminho literal tirado de `.harness/local.yaml`; um caminho literal não casa com uma regra de permissão `Bash(grok:*)`, então essa forma pode precisar da aprovação do dono ou ser `not-run`), sem `timeout`, `&&`, `;` ou pipe antes ou depois dele. `--cwd <CWD>` é a flag de diretório de trabalho do próprio binário, então nenhum `cd` é necessário. `-p` / `--single` imprime no stdout e sai. `--always-approve` aprova ferramentas automaticamente, necessário para uma chamada sem supervisão; `--permission-mode bypassPermissions` é o equivalente nomeado.

| Flag | Por quê |
|---|---|
| `-p` / `--single` | prompt de disparo único no argv |
| `--prompt-file` | prompt de disparo único vindo do disco |
| `-m` / `--model` | fixa o slug |
| `--always-approve` | ferramentas sem supervisão |
| `--cwd` | espaço de trabalho; `--cwd "$CLONE"` em toda chamada |
| `--reasoning-effort` / `--effort` | `none` / `minimal` / `low` / `medium` / `high` / `xhigh` / `max` (só os níveis que o modelo anuncia) |
| `--output-format` | `plain` (padrão) / `json` / `streaming-json` / `streaming-messages-json` |
| `--json-schema` | restringe a saída JSON (implica `--output-format json`) |
| `--max-turns` | limita os turnos agênticos |
| `--no-subagents` | impede criar um subagente de dentro da chamada |

Não inicie a TUI interativa a partir de um turno de orquestrador.

## Entrada do prompt

`--prompt-file` é o contrato para qualquer coisa maior que uma linha; `-p` recebe um prompt curto no argv. `$CLONE` e `$PROMPT` são definidos em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Nenhuma flag somente leitura documentada: clone descartável ([clone-isolation.md](clone-isolation.md)), comparado como [live-tree-freeze.md](live-tree-freeze.md) diz. `--no-subagents` estreita mas não substitui a cláusula anti-delegação: passe as duas.

## Flag de modelo

`-m` / `--model <slug>`, de `models.grok` em `.harness/local.yaml`.

## Timeout

Nenhuma flag de tempo real documentada além de `--max-turns`. O timeout da ferramenta Bash (ajustado ao timeout do revisor) ou `run_in_background` limita o tempo.

## Captura da saída

Resposta no stdout; `--output-format` seleciona o formato (padrão `plain`).

## Sinais de falha

Já se viu diffs grandes com `--prompt-file` em esforço alto estourarem o timeout ou imprimirem só um preâmbulo de raciocínio com saída 0.

## Pegadinhas conhecidas

Ver [gotchas.md](gotchas.md), classes `timed_out` e `no_signal`.

## Medido em

2026-10-16 (`grok --help`: `--cwd`).
