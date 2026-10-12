# PRD-018 · `dh config`: ler e gravar o `config.yaml` global pelo binário

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-12 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-12 / 2026-10-12 |
| Tickets | T-1801 |
| Origem | chat do dono, 2026-10-12; `docs/discovery/renova.md` §7 (forma `dh config set <chave> <valor>` em Go; escrita só no binário, PRD-014 R2); a decisão 10 (A, adiar o editor, 2026-10-09) foi superada pelo pedido do dono no chat de 2026-10-12 |
| Depende de | PRD-014 (pasta global, `config.yaml`, trava `config.lock`) |
| Libera | nada |

## 1. Problema

Fatos verificados em 2026-10-12:

- A 0.21.0 removeu `DH_DASHBOARD_SPRITES` e deixou sem efeito a configuração `dh.dashboard.sprites` da extensão; o `CHANGELOG.md` (migração da 0.21.0) manda pôr a pasta de sprites em `dashboard.sprites` do `config.yaml`. Nenhum subcomando do `dh` escreve esse campo, e o `AGENTS.md` diz que só o `dh` lê e escreve o `config.yaml`.
- Resultado: o dashboard reiniciado na 0.24.0 seguiu com os avatares padrão porque `dashboard.sprites` estava vazio; o Coordenador precisou de um programa Go descartável chamando `harness.LoadConfig`/`SaveConfig`.
- `internal/harness/config.go` modela `Config{Language, Mode, Sprites, ReviewersRaw, Projects}` e valida `mode` em `repo`/`global`; `config_save.go` (`SaveConfig`) preserva as linhas não modeladas byte a byte e grava de forma atômica com 0600; `lock.go` guarda `config.lock`, usada por `Link`.
- `internal/dashboard/server.go` (`spritesDirOf`) aceita valor com `~/`, exige diretório absoluto existente e relê a cada montagem de estado (sem reinício).
- `cmd/dh/main.go` despacha os subcomandos e lista-os na string de uso; `harness-path`, `link` e `projects` usam `splitArgs` e `--json`.

## 2. Solução

Um subcomando `dh config` lê e altera os três escalares modelados do `config.yaml` global: `language`, `mode` e `dashboard.sprites`. Sem argumentos, imprime; `set <chave> <valor>` grava; `unset <chave>` esvazia. A gravação usa a mesma trava e o mesmo `SaveConfig` do `dh link`. `reviewers` e `projects` ficam fora. O dashboard em execução pega os sprites novos no próximo refresh, sem reinício. Sem comando `/dh:` novo (continua 19 + 1), sem editor no dashboard, sem dependência nova.

## 3. Regras

- R1 `dh config` sem argumentos imprime uma linha `chave: valor` para cada uma de `language`, `mode`, `dashboard.sprites`, nessa ordem (valor vazio aparece vazio); com `--json` imprime um objeto JSON com exatamente estas três chaves, nesta ordem (`language`, `mode`, `dashboard.sprites`), gravado por `writeJSON`. Verificável: teste Go da saída de texto e JSON, com config cheio e com config vazio.
- R2 `dh config set <chave> <valor>` aceita só as chaves `language` (`en` | `pt-br`), `mode` (`repo` | `global`) e `dashboard.sprites` (caminho absoluto ou começando com `~/` que, depois de expandir `~/`, é um diretório existente; o valor lido de volta por `LoadConfig` é o argumento literal, sem expandir `~`, para o `~/` continuar portátil). `set` com `config.yaml` ausente cria o arquivo. Chave ou valor inválido: saída 1, mensagem legível no stderr e arquivo intocado. Verificável: teste de tabela com um caso válido e um inválido por chave (incluindo caminho relativo e caminho inexistente) e comparação byte a byte do arquivo nos inválidos.
- R3 `dh config unset <chave>` aceita as mesmas três chaves e remove o valor; chave desconhecida: saída 1, mensagem no stderr, arquivo intocado. Verificável: teste Go de `unset` das três chaves e de chave inválida.
- R4 `reviewers` e `projects` são recusados em `set` e `unset` (saída 1, arquivo intocado), com mensagem que aponta para editar `reviewers` à mão no `config.yaml`/`.harness/project.yaml` e para `dh link` no caso de `projects`. Verificável: teste Go das duas recusas, conferindo a mensagem e o arquivo intocado.
- R5 A gravação toma a mesma trava (`config.lock`) do `dh link` e passa por `SaveConfig`: as demais linhas ficam byte a byte iguais, a escrita é atômica e 0600, e um arquivo existente inválido é erro e nunca é sobrescrito. Verificável: teste Go com arquivo contendo `reviewers`, `projects` e comentário, comparando o resto do arquivo antes e depois; teste com `mode` inválido no arquivo existente (erro, arquivo intocado); a gravação passa por uma função exportada em `internal/harness` (p. ex. `SetConfigValue`) que toma o `lock(home)` não exportado, como `Link`; `grep` confirma.
- R6 Um dashboard em execução passa a usar o `dashboard.sprites` novo no próximo refresh, sem reinício. Verificável: teste Go (ou walkthrough do `qa-verifier`) em que `set dashboard.sprites` seguido de nova montagem do estado devolve a pasta nova.
- R7 Sem comando `/dh:` novo (contagem 19 + 1), sem editor de configuração no dashboard e sem dependência nova. Verificável: `git diff go.mod` vazio; `dh validate` verde com a contagem de comandos inalterada.
- R8 Versão 0.25.0: `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte do mesmo arquivo (`v0.25.0`) mudam juntos; o `ref` esquecido é falha de entrega (a PR #35 só existiu para corrigir isso). Verificável: `grep` dos três valores = `0.25.0`/`v0.25.0` no mesmo commit.

## 4. Docs

- `cmd/dh/main.go` (string de uso lista `config`)
- `docs/tutorial.html`, `docs/en/tutorial.html` (o bullet de `dashboard.sprites` e o layout do `config.yaml` citam o comando)
- `README.md`, `README.en.md` (onde os subcomandos do `dh` são listados)
- `CHANGELOG.md` (0.25.0), `docs/roadmap.html`, `docs/en/roadmap.html` (validados contra o `CHANGELOG.md` pelo `dh validate`)
- `AGENTS.md` (a linha que lista `harness-path`, `link` e `projects` como subcomandos do `dh` ganha `config`)
- `.skills/project-onboarding/SKILL.md`: verificado, cita `harness-path` e `link` por uso e não lista o conjunto de subcomandos; sem mudança

## 5. Fora de escopo

- Editar `reviewers` ou `projects` pelo comando.
- Editor HTTP/dashboard.
- Voltar `DH_DASHBOARD_SPRITES`.
- Mudar a extensão VS Code.

## 6. Decisões do dono

- Forma conforme `docs/discovery/renova.md` §7: `dh config set <chave> <valor>` em Go, escrita só no binário (PRD-014 R2). A decisão 10 (A, adiar o editor, 2026-10-09) foi superada pelo pedido do dono no chat de 2026-10-12, que também autorizou o desenho acima e a execução de ponta a ponta.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** `dh config` mostra `language`, `mode` e `dashboard.sprites` (texto ou `--json`), e `dh config set` / `unset` gravam essas três chaves no `config.yaml` global pelo próprio binário, sem edição à mão. Saiu na 0.25.0.
- **Regras:**
  - R1 Honrada: saída texto e `--json` com as três chaves na ordem (`TestConfigCommand`; QA com `DH_HOME` temporário).
  - R2 Honrada: valores válidos gravados, `~/` guardado literal, arquivo ausente criado com 0600; inválidos saem com 1 e `cmp` mostra o arquivo intacto (`config_set_test.go`, QA).
  - R3 Honrada: `unset` das três chaves; chave desconhecida recusada com o arquivo intacto.
  - R4 Honrada: `reviewers` e `projects` recusados em `set` e `unset`, com as mensagens que apontam edição à mão e `dh link`.
  - R5 Honrada: `harness.SetConfigValue` toma `lock(home)` como o `Link` e passa por `SaveConfig`; comentários, `reviewers` e `projects` iguais byte a byte; arquivo inválido não é sobrescrito.
  - R6 Honrada: dashboard real com `DH_HOME` temporário trocou `config.sprites` em `/api/state` sem reiniciar (cache de estado de cerca de 3 s).
  - R7 Honrada: `go.mod` sem mudança, 19 arquivos em `.commands/`, sem editor no dashboard.
  - R8 Honrada: 0.25.0 em `build.go`, `version` e `ref` `v0.25.0` no `marketplace.json`, `plugin.json` gerado em 0.25.0.
- **Tickets:**
  - T-1801 backend: aprovado (QA APPROVED; revisão de código claude approve, 2 menores aceitos).
- **Docs:** `CHANGELOG.md`, `README.md`, `README.en.md`, `AGENTS.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`.
- **Fora:**
  - Editar `reviewers` ou `projects` pelo comando.
  - Editor HTTP ou no dashboard.
  - Volta do `DH_DASHBOARD_SPRITES`.
  - Mudança na extensão VS Code.
