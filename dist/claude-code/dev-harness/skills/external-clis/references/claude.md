# claude (Claude Code)

Binário: `claude`. Utilizável como revisor de CLI, distinto de um Coordenador que despacha um agente Claude no mesmo processo: este caminho chama um processo `claude` separado. Não é planejador: `cli:claude/<slug>` fica excluído de `/dh:plan-loop` enquanto o kit roda sobre o Claude.

## Invocação

```bash
cd "$CLONE" && claude -p --dangerously-skip-permissions --model <slug> --effort high \
  "$(cat "$PROMPT")"
```

`claude --help` não lista nenhuma flag de diretório de trabalho (só `--add-dir`, que acrescenta raízes e não define o cwd), então esta receita mantém `cd "$CLONE" && claude ...` como a primeira de duas exceções nomeadas (com `opencode`, [opencode.md](opencode.md)) à regra de que uma chamada é um comando que começa com o nome do binário. Uma chamada Bash não conserva o cwd para a chamada seguinte, então o `cd` fica dentro do mesmo comando. Essa forma não casa com uma regra de permissão `Bash(claude:*)`: a chamada pode precisar da aprovação do dono, ou é `not-run` se for recusada. Nenhum invólucro `timeout` vai antes dele. A chamada nunca passa `--add-dir` para a árvore viva. `-p` / `--print` imprime e sai. `--dangerously-skip-permissions` contorna o prompt de permissão, necessário para uma chamada sem supervisão; `--permission-mode bypassPermissions` é o equivalente nomeado.

| Flag | Por quê |
|---|---|
| `-p` / `--print` | disparo único; `claude` sozinho é a TUI |
| `--dangerously-skip-permissions` | sem supervisão |
| `--model` | um alias (`haiku`, `sonnet`, `opus`, `fable`, …) ou um id de modelo completo |
| `--effort` | `low` / `medium` / `high` / `xhigh` / `max` |
| `--fallback-model` | fixa um modelo de reserva se o padrão estiver indisponível |
| `--add-dir` | raízes extras do espaço de trabalho (nunca a árvore viva numa chamada em clone) |
| `--output-format` | só com print: `text` / `json` / `stream-json` |
| `--input-format` | só com print: `text` ou `stream-json` |

Não inicie a TUI interativa a partir de um turno de orquestrador.

## Entrada do prompt

Último argumento, ou expandido a partir de um arquivo como acima. `$CLONE` e `$PROMPT` são definidos em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Nenhuma flag somente leitura documentada como processo de CLI: clone descartável ([clone-isolation.md](clone-isolation.md)), comparado como [live-tree-freeze.md](live-tree-freeze.md) diz.

## Flag de modelo

`--model <slug>`: um alias ou um id de modelo completo.

## Timeout

Nenhuma flag de tempo real. O timeout da ferramenta Bash ou `run_in_background` limita o tempo.

## Captura da saída

Resposta no `stdout`; `--output-format` seleciona `text`, `json` ou `stream-json`.

## Sinais de falha

Uma execução longa que nunca emite o sinal de veredito da etapa é `no_signal`.

## Pegadinhas conhecidas

Nenhuma registrada em [gotchas.md](gotchas.md). O formato de veredito é o de `SKILL.md`, "O contrato de veredito"; esta receita não define outro.

## Medido em

2026-10-16 (`claude --help`: sem flag de cwd). A chamada em si: desconhecida (não medida neste kit).
