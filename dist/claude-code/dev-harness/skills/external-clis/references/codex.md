# codex

Binário: `codex`. Revisor; planejador (modo somente leitura documentado e medido). Mãos: sim, via `codex exec`, a menos que rode sob `-s read-only`.

## Invocação

```bash
codex exec --skip-git-repo-check -s read-only -m <slug> - < "$PROMPT"
```

Um comando de shell cuja primeira palavra é `codex`, com o redirecionamento de `stdin` como parte dele: nenhum `timeout`, `&&`, `;` ou pipe antes ou depois dele. `exec` é o caminho não interativo; `codex` sozinho é a TUI e nunca é iniciado a partir de um turno de orquestrador.

| Flag | Por quê |
|---|---|
| `exec` | disparo único |
| `-` | lê o prompt do `stdin` |
| `-m` / `--model` | fixa o slug |
| `--skip-git-repo-check` | obrigatória fora de um repositório git confiável; senão, saída 0 com `stdout` vazio |
| `-C` / `--cd` | raiz do espaço de trabalho |
| `--add-dir` | diretórios extras graváveis |
| `-s` / `--sandbox` | `read-only` / `workspace-write` / `danger-full-access` |
| `--ephemeral` | não persiste a sessão |

`codex review` (também `codex exec review`) é o caminho dedicado de revisão; prefira-o quando o trabalho é só um veredito sobre a árvore atual:

```bash
codex review --uncommitted -c model="<slug>"
```

`--uncommitted` cobre alterações preparadas para commit, não preparadas e não rastreadas. `--base BRANCH` revisa contra essa branch; `--commit SHA` revisa esse commit. `-m`/`-C` não são válidos no subcomando `review`, e `--uncommitted` não se combina com um argumento `[PROMPT]` / `-`.

## Entrada do prompt

`stdin` (`-`). Se o `stdin` é canalizado e há um argumento de prompt, o `stdin` é acrescentado como um bloco `<stdin>`. `$PROMPT` é definido em [clone-isolation.md](clone-isolation.md), Termos.

## Isolamento

Flag somente leitura: `-s read-only`. Sem clone; a árvore viva é comparada como [live-tree-compare.md](live-tree-compare.md) diz. Sob `-s read-only`, `/tmp` é montado como somente leitura, então um `go test` da cadeia de ferramentas Go e checagens parecidas que compilam em tmp falham lá dentro. Trate a chamada como somente leitura e raciocínio, não como validação de comportamento: rode o gate de testes e de lint do próprio projeto no host se for preciso um resultado real. É por isso que a permissão de `commands.test`/`commands.lint` na cláusula anti-delegação é inerte na prática para este binário. `-s workspace-write` também bloqueia a rede de saída do sandbox, então uma suíte que precisa de um serviço ativo também não roda lá.

## Flag de modelo

`-m` / `--model <slug>`, o slug de `models.codex` em `.harness/local.yaml`. O Codex não tem subcomando `models`: a lista de slugs é atualizada à mão quando o operador mede um novo.

## Timeout

`codex exec` não tem flag de tempo real. O timeout da ferramenta Bash (ajustado ao timeout do revisor) ou `run_in_background` limita o tempo.

## Captura da saída

Resposta no `stdout`.

## Sinais de falha

- Saída 0 com `stdout` vazio: falta `--skip-git-repo-check` fora de um repositório git confiável.
- Um `go test` que não consegue compilar sob `-s read-only` (ver Isolamento).

## Pegadinhas conhecidas

Ver [gotchas.md](gotchas.md), classes `no_signal` e `tree_changed`.

## Medido em

2026-09-22 (execução de formatação do veredito). A data da sondagem somente leitura que qualifica `codex` como planejador: desconhecida.
