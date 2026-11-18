# agy (Antigravity)

Binário: `agy`. Revisor entre outros papéis. Sem planejamento até que uma sondagem datada registre que `agy` trata o cwd do processo como o espaço de trabalho dele e qual flag, se houver, o define.

## Invocação

```bash
agy --dangerously-skip-permissions --print-timeout 15m --print "$(cat "$PROMPT")"
```

Um comando de shell cuja primeira palavra é `agy`: nenhum `timeout`, `&&`, `;` ou pipe antes ou depois dele. O congelamento anterior à chamada e a comparação posterior são comandos separados. A ordem é determinante:

1. `--dangerously-skip-permissions`
2. flags de sessão opcionais (`--model`, `--effort`, `--print-timeout`, `--output-format`)
3. `--print` (ou `-p` / `--prompt`)
4. a string do prompt

| Flag | Por quê |
|---|---|
| `--dangerously-skip-permissions` | aprovação de ferramentas sem supervisão |
| `--print` / `-p` / `--prompt` | disparo único; `agy` sozinho é a TUI |
| `--print-timeout` | limite de espera (`15m` para uma chamada de revisão) |
| `--model` | o slug de `models.agy` em `.harness/local.yaml` |
| `--effort` | `low` / `medium` / `high` (não `xhigh`) |
| `--output-format` | `text` (padrão) / `json` / `stream-json` |
| `--json-schema` | saída final estruturada |
| `--mode` | `accept-edits` ou `plan` |
| `--sandbox` | restringe o terminal |

`--input-format stream-json` lê turnos NDJSON no `stdin` e exige `--output-format stream-json`. O caminho padrão é uma string de prompt. Se algum dia se acrescentar uma receita de clone, ela usa a flag de cwd achada pela sondagem dentro do único comando, ou declara a exceção do `cd` como [claude.md](claude.md) faz.

## Entrada do prompt

Embutida na string do prompt. Prefira embutir arquivos no prompt a usar `--add-dir`: já se viu `--add-dir PATH` antes do prompt imprimir a ajuda de `--add-dir` e ignorar o pacote. `$PROMPT` é definido em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Nenhuma flag somente leitura documentada; `agy` opera no cwd do processo (não sondado). A revisão fica na árvore viva sob a exposição aceita RISK-002, com a comparação de [live-tree-compare.md](live-tree-compare.md): não é um sandbox, só detecta depois do fato uma edição no espaço de trabalho. Um projeto cujo `.git` não é um diretório termina em `isolation_unavailable` (ver aquele arquivo).

## Flag de modelo

`--model <slug>`, de `models.agy` em `.harness/local.yaml`. Atualize os slugs desta máquina com `agy models` e registre o resultado em `.harness/local.yaml`, nunca neste arquivo.

## Timeout

`--print-timeout`; o padrão é 5m, e uma chamada de revisão muitas vezes precisa de mais (`15m` acima). O timeout pode não parar o processo (ver Pegadinhas conhecidas).

## Captura da saída

Resposta no `stdout`; `--output-format` seleciona `text` (padrão), `json` ou `stream-json`.

## Sinais de falha

- Texto de ajuda em vez de uma resposta: `--add-dir` posto antes do prompt.
- Uma resposta sem o sinal de veredito da etapa é `no_signal`, com o processo vivo ou não.

## Pegadinhas conhecidas

O timeout pode não parar o processo. Uma resposta completa em disco que traga o sinal de veredito da etapa é o resultado; pare o processo com um comando separado, `kill <pid>` do pid registrado quando a chamada começou (a chamada roda em segundo plano e o pid é registrado, ou o delegado lista `pgrep -x agy` antes e depois da chamada e mata só o pid novo), para que as sessões `agy` do próprio dono nunca sejam casadas. Não use `pkill -f` com o nome-base do prompt nem com qualquer texto do comando: ele casa com o shell do próprio delegado ([gotchas.md](gotchas.md)); um padrão com colchetes (`pkill -f '[a]gy --print'`) é o recurso de reserva. Reporte isso como `leftover processes` na checagem do espaço de trabalho.

## Medido em

2026-10-16 (processo mantido vivo além de `--print-timeout 15m` depois de a resposta completa ter sido escrita).
