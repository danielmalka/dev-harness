# opencode

Binário: `opencode`, resolvido por `binaries.opencode` em `.harness/local.yaml` quando `command -v opencode` falha. Um `PATH` só de shell de login é comum para este binário: `command -v` não o encontra mesmo quando ele está instalado, então não trate essa falha sozinha como `not_installed`. Revisor ou `builder`, e planejador (num clone descartável); mãos: sim (`opencode run`).

## Invocação

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>#high' \
  "$(cat "$PROMPT")"
```

Prompt longo como arquivo anexado (não enfie um prompt de vários kilobytes no argv):

```bash
cd "$CLONE" && opencode run --auto -m '<provider>/<slug>[#<variant>]' \
  --file="$PROMPT" -- "Siga o prompt anexado."
```

`opencode run` não tem flag de diretório de trabalho nem `--variant`. A chamada por isso começa com `cd "$CLONE" &&`, a segunda exceção nomeada (ao lado de `claude`) à regra de que uma chamada é um comando cuja primeira palavra é o binário. Um `cd` separado numa chamada Bash própria não persiste (o ambiente de execução zera o cwd), então o `cd` fica no mesmo comando. Como o comando começa com `cd`, ele não casa com uma regra de permissão `Bash(opencode:*)`: a chamada pode precisar da aprovação do dono ou ser `not-run`. Escreva o binário resolvido como literal (`opencode` simples, ou o caminho literal tirado de `.harness/local.yaml`). Nenhum `timeout` vai na frente do comando e nenhum `;` ou pipe o encadeia; o congelamento e as comparações posteriores são comandos separados.

| Flag | Por quê |
|---|---|
| `run` | disparo único; `opencode` sozinho é a TUI |
| `--auto` | aprova ferramentas (obrigatória para uma chamada sem supervisão) |
| `-m provider/model[#variant]` | fixa; não dependa do seletor interativo; `#variant` define o esforço de raciocínio (só os níveis que esse slug anuncia); não existe flag `--variant` |
| `cd "$CLONE" &&` | espaço de trabalho; não existe `--dir`, então toda chamada começa com isso |
| `-f` / `--file` | anexa arquivo(s) à mensagem (repetível) |
| `--format json` | eventos brutos se for preciso analisar |
| `--thinking` | imprime os blocos de raciocínio |
| `--title` | título da sessão |
| `-c` / `--continue`, `-s` / `--session`, `--fork` | retomar / bifurcar; o padrão é uma sessão nova |

Recuse, a partir de um turno de orquestrador: `opencode` sozinho, `opencode --mini`, `run --interactive` / `-i`. `--agent` em `run` troca o agente primário, não um subagente no mesmo processo (`task` dentro do opencode): a cláusula anti-delegação continua valendo para o que esse agente primário pode fazer.

## Entrada do prompt

Um arquivo anexado: `--file="$PROMPT" -- "message"`. `-f`/`--file` é um array do yargs: um `-f file.md "msg"` sem mais nada engole `msg` como segundo arquivo, e `--file=` sozinho, sem mensagem após `--`, é recusado. Sempre `--file=path -- "message"`. `$CLONE` e `$PROMPT` são definidos em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Nenhuma flag somente leitura documentada: clone descartável ([clone-isolation.md](clone-isolation.md)), comparado como [live-tree-freeze.md](live-tree-freeze.md) diz.

## Flag de modelo

`-m '<provider>/<slug>[#<variant>]'`, de `models.opencode` em `.harness/local.yaml`. A variante (esforço de raciocínio) faz parte da string do modelo, só um nível que esse slug anuncia; sem `#variant` vale o padrão do modelo. O catálogo de modelos é específico de cada provedor e muda com o tempo: não invente um slug aqui; atualize-o rodando nesta máquina o próprio comando de listagem de modelos do provedor e registre o resultado em `.harness/local.yaml`, nunca neste arquivo.

## Timeout

`opencode run` não tem flag de tempo real. O timeout da ferramenta Bash ou `run_in_background` limita o tempo.

## Captura da saída

Resposta no `stdout`; `--format json` dá eventos brutos se for preciso analisar.

## Sinais de falha

- Um erro de trava de banco de dados na segunda de duas execuções simultâneas contra a mesma conta.
- Uma chamada `--file=` recusada, sem mensagem após `--`.

## Pegadinhas conhecidas

Não rode duas chamadas `opencode run` ao mesmo tempo. Quando uma etapa de revisão lista mais de uma entrada `opencode`, ou o teto de concorrência seria atingido de outro modo, o Coordenador serializa as chamadas `opencode` em seu próprio lote em vez de rodá-las lado a lado. Ver também [gotchas.md](gotchas.md).

## Medido em

2026-10-21, v2.0.32 (sem flag de cwd, sem `--variant`; o ambiente de execução zera o cwd entre chamadas Bash).
