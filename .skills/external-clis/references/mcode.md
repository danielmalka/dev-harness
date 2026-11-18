# mcode (MiniMax Code)

Binário: `mcode`, resolvido por `binaries.mcode` em `.harness/local.yaml` quando `command -v mcode` falha. Um `PATH` só de shell de login é comum para este binário: `command -v` não o encontra mesmo quando ele está instalado, então não trate essa falha sozinha como `not_installed`. Agente de código com mãos; revisor e planejador (num clone descartável).

## Invocação

```bash
mcode exec --permission off --cwd "$CLONE" --timeout 15m \
  --model <provider>/<slug> \
  --input - < "$PROMPT"
```

Um comando de shell cuja primeira palavra é o binário: escreva ali o comando resolvido como literal (`mcode` simples, ou o caminho literal tirado de `.harness/local.yaml`; um caminho literal não casa com uma regra de permissão `Bash(mcode:*)`, então essa forma pode precisar da aprovação do dono ou ser `not-run`), sem `timeout`, `&&`, `;` ou pipe antes ou depois dele. `--cwd <path>` é a flag de diretório de trabalho do próprio binário (`"workspace directory"`). O redirecionamento `< "$PROMPT"` faz parte do mesmo comando. `--permission full` é o modo sem supervisão com capacidade de escrita, só para o papel de `builder`: nunca o passe a partir desta skill, que cobre etapas de revisão.

| Flag | Por quê |
|---|---|
| `exec` | disparo único; `mcode` sozinho é a TUI |
| `--permission off` | documentada, não medida: somente leitura, sem supervisão (use para um revisor, num clone) |
| `--permission full` | sem supervisão e com acesso de escrita (só papéis de `builder`, fora do escopo desta skill) |
| `--cwd` | espaço de trabalho; `--cwd "$CLONE"` em toda chamada |
| `--model provider/model` | fixa o slug |
| `--input -` | lê o prompt do stdin (`--input-format text` / `json`) |
| `--file PATH` | anexa um arquivo (repetível) |
| `--timeout 15m` | uma chamada de revisão muitas vezes precisa de mais que o padrão |
| `--max-steps` | limita os passos do assistente |
| `--output-format` | `text` / `json` / `stream-json` |
| `--output-schema` | JSON Schema para a resposta final |
| `-o` / `--output-last-message PATH` | mensagem final do agente em disco |

Não inicie a TUI interativa a partir de um turno de orquestrador.

## Entrada do prompt

Stdin (`--input -`) ou `--file`. `$CLONE` e `$PROMPT` são definidos em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Flag somente leitura documentada, ainda não medida neste kit: `--permission off` desliga a aprovação de ferramentas por completo em vez de aprová-las automaticamente, que é o modo somente leitura do próprio binário. Nenhuma sondagem registrada confirma que ele recusa uma tentativa de criar arquivo e de rodar comando de shell, então, até haver uma sondagem datada registrada, `mcode` roda num clone descartável ([clone-isolation.md](clone-isolation.md)) tanto para revisão quanto para planejamento. A flag é usada além da comparação que toda chamada em clone recebe ([live-tree-freeze.md](live-tree-freeze.md)), nunca no lugar dela.

## Flag de modelo

`--model <provider>/<slug>`, de `models.mcode` em `.harness/local.yaml`, atualizado à mão quando o operador mede um novo slug para esta conta.

## Timeout

`--timeout 15m` limita o tempo.

## Captura da saída

Resposta no stdout; `-o` / `--output-last-message PATH` escreve a mensagem final do agente em disco; `--output-format` seleciona `text`, `json` ou `stream-json`.

## Sinais de falha

desconhecido

## Pegadinhas conhecidas

Nenhuma registrada em [gotchas.md](gotchas.md).

## Medido em

2026-10-16 (`mcode exec --help`: `--cwd`). `--permission off`: desconhecido (não medido).
