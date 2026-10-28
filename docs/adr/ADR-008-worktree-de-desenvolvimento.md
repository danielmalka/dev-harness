# ADR-008 · Worktree de desenvolvimento criado pelo Coordenador, com o sandbox de CLI ainda em clone

| Campo | Valor |
|---|---|
| Status | aceito (validação de documento aprovada em 2026-10-27, rodada 2) |
| Data | 2026-10-27 |
| Decisor | Daniel Malka (dono), PRD-028 aprovado em 2026-10-27 com as suposições A1–A8 |
| Origem | PRD-028 (R1–R36, R56, R61–R63); pendências do PRD: pasta do worktree (A5) e ferramenta do motor × `git worktree add` |
| Reversibilidade | cara: o resolvedor (`dh harness-path`), o isolamento de CLI (controle de segurança do RISK-002), os registros de tarefa e os comandos `build`/`auto` passam a depender da forma escolhida |
| Relação | estende ADR-007 (resolvedor único) e ADR-002/PRD-008 (clone descartável, que fica); não altera ADR-005 nem a trava de `gh`/`git` do PRD-027 (A7) |

## 1. Contexto

- Pergunta: onde e por quem o worktree de uma execução é criado, como o kit reconhece um worktree ligado (para achar o harness dir e para aceitar o projeto no isolamento de CLI) e quando ele some. Cada escolha muda o que o builder, o resolvedor e o revisor CLI enxergam.
- Fato: o resolvedor procura `<repo>/.harness` e depois o caminho canônico no mapa `projects:`; um worktree é outro caminho e resolve `none` em modo global (`internal/harness/resolve.go`, `resolveIn`).
- Fato: `dh harness-path --json` calcula `pdocs` com a chave do próprio diretório (`cmd/dh/main.go`, `runHarnessPath`, `harness.PdocsDir(home, c.Projects[key])`); de dentro de um worktree o campo sai vazio.
- Fato: o isolamento exige `test -d <project>/.git && test ! -L <project>/.git` e recusa um gitfile (`.skills/external-clis/references/clone-isolation.md`, passo 1). A foto R5 cobre, do `.git` vivo, só `hooks/`, `config`, `info/`, `HEAD`, `packed-refs` e `git for-each-ref`; `worktrees/` e `refs/` como arquivos ficam de fora (mesmo arquivo, §7; ADR-002).
- Fato: a lista `WRITE` da trava tem `push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`; não tem `worktree` (`adapters/claude-code/plugin/hooks/lib/gh-guard.ts`, linha 10). A R12 depende só de prosa.
- Fato: `gate.ts` trata como escrita de registro (não código) todo caminho que contém `/.harness/` (`adapters/claude-code/plugin/hooks/lib/gate.ts`, linha 40). Um worktree dentro de `<home>` (padrão `~/.harness`) teria todo o código do builder classificado como registro, e a trava de commit do ADR-006 deixaria de marcar.
- Fato: a ferramenta de worktree do Claude Code (`EnterWorktree`) cria o worktree em `.claude/worktrees/` dentro do repositório, troca o diretório da sessão e, por padrão, parte de `origin/<branch padrão>` (descrição da ferramenta na sessão de 2026-10-27).
- Fato: `dh config` aceita só `language`, `mode`, `dashboard.sprites` (`internal/harness/config_set.go`, `ConfigKeys` e `SetConfigValue`); `Config` não tem campo para `worktree` (`internal/harness/config.go`).
- Forças: o checkout principal fica limpo; nenhuma dependência Go nova e nenhum `git` chamado pelo resolvedor (ADR-001); nenhum caminho absoluto privado em arquivo versionado; o RISK-002 continua "mitigado", sem perder cobertura; o kit não deve depender de recurso exclusivo do Claude Code onde `git` resolve.

## 2. Opções consideradas

### D1 · Pasta do worktree

| Dimensão | A · irmã do checkout: `<pai>/<nome-do-repo>.worktrees/<task-id>` | B · `<home>/worktrees/<projeto>/<task-id>` | C · `~/.dh-worktrees/<projeto>/<task-id>` | D · dentro do repo (`.claude/worktrees/`, ignorado) |
|---|---|---|---|---|
| Fora da árvore principal (A5) | sim | sim | sim | não |
| Precisa de ignore no git | não | não | não | sim |
| Trava de commit (`gate.ts`) | intacta | quebra: caminho contém `/.harness/` | intacta | intacta |
| Mesmo sistema de arquivos do repo | sim (repo em `/mnt/c` fica lá) | não necessariamente | não necessariamente | sim |
| Caminho no registro versionado | relativo ao checkout, sem nada privado | exige `~/` ou `<home>` | exige `~/` | relativo |
| Permissão de escrita no Claude Code | uma pasta extra por projeto | já coberta por `additionalDirectories` (`~/.harness`) | uma entrada fixa | já coberta |

### D2 · Método de criação

| Dimensão | A · `git worktree add` pelo Coordenador na sessão principal | B · `EnterWorktree` do Claude Code |
|---|---|---|
| Local | o de D1 | `.claude/worktrees/` dentro do repo (viola A5) |
| Diretório da sessão | não muda; mod, trava e `dh harness-path` seguem no checkout principal | troca; tudo passa a resolver a partir do worktree |
| Ponto de partida | o que o comando escolhe (`main` no `auto`) | `origin/<padrão>` ou `HEAD`, por configuração do motor |
| Portabilidade | qualquer motor com `git` | só Claude Code |

### D3 · Reconhecer um worktree no resolvedor

| Dimensão | A · Go lê o gitfile, `commondir` e o ponteiro de volta | B · Go chama `git rev-parse --git-common-dir` | C · manter como está |
|---|---|---|---|
| Executa código do repo | nada | o `git` lê a config do repo-alvo antes de validar | nada |
| Dependência | nenhuma | `git` no PATH do `dh` | nenhuma |
| Cumpre R18–R19 | sim | sim | não (`none` em global) |

### D4 · Isolamento de CLI a partir de um worktree

| Dimensão | A · validar o layout (Go, via `dh harness-path --json`) e clonar como hoje | B · validar em prosa de shell no passo 1 | C · manter a recusa |
|---|---|---|---|
| Uma regra só para resolvedor e isolamento | sim | não, duas implementações que derivam | n/a |
| Cumpre R26–R31 | sim | sim | não (opção b, rejeitada pelo dono) |

D5 e D6 têm uma opção viável cada (R9; R22–R25); a seção 3 registra só a forma e os casos-limite.

## 3. Decisão

### D1 · Pasta: A, irmã do checkout principal

- Caminho: `<pai do checkout principal>/<nome da pasta do checkout>.worktrees/<task-id>`. O `<task-id>` é o da tarefa (`/dh:build`) ou o `AUTO-<n>` (`/dh:auto`) e só aceita `[A-Za-z0-9_-]`; outro caractere encerra pela R17.
- Colisão: dois repositórios não dividem nome na mesma pasta pai, e um `task-id` só tem um worktree. Caminho existente que não seja o worktree registrado da mesma tarefa (ver D6, retomada) encerra pela R17; nunca há sufixo automático.
- Ignore: fora da árvore, o `git status` do checkout principal não o vê; nenhuma mudança em `.gitignore`.
- Registro (R16): o caminho entra no registro da tarefa relativo ao checkout principal, `../<nome>.worktrees/<task-id>`, mais a branch. O despacho ao builder (R14) leva o caminho absoluto, calculado na hora e nunca gravado em arquivo versionado. Assim, nem o registro em modo `repo` (versionado) nem este ADR carregam caminho absoluto privado.
- Motivo decisivo: B quebra a trava de commit (`gate.ts` linha 40) sem mudar o mod; A mantém o worktree no mesmo sistema de arquivos do repo (no WSL, um repo em `/mnt/c` não ganha worktree em ext4) e não depende de `<home>`.
- Descartadas: B, pela trava de commit e por misturar código em `<home>`, que é registro; C, por cruzar sistema de arquivos e exigir mais uma pasta global sem dono; D, por violar A5 e precisar de ignore.
- Preferência declarada: nenhuma.

### D2 · Criação: A, `git worktree add` pelo Coordenador

- Forma: na sessão principal, `git worktree add -b <branch> <caminho> <base>`. No `/dh:auto`, `<branch>` é o `<base-branch>` da execução e `<base>` é `main` (R15). No `/dh:build`, o comando não tem convenção de branch, então a branch é o `<task-id>` em minúsculas (R62); `<base>` é o `HEAD` atual do checkout principal.
- Depois de criar, o Coordenador confere que `dh harness-path --json <caminho>` devolve o mesmo `dir` do checkout principal. Falha aqui também é R17.
- Especialistas: proibidos em prosa de `git worktree add|remove|prune` (R12); a trava não muda (A7, confirmado em `gh-guard.ts` linha 10).
- Motivo decisivo: `EnterWorktree` cria dentro do repo e troca o diretório da sessão, o que move a resolução do harness dir e o alvo da trava para o worktree; `git` funciona em qualquer motor.
- Descartada: B, pelos três motivos acima e por escolher a base por configuração do motor, não pela regra do `auto`.

### D3 · Resolução: A, leitura em Go sem chamar `git`

- Uma função em `internal/harness` (nome a cargo do planejador) recebe o diretório e devolve o checkout principal só quando todas as condições valem; senão o diretório não é worktree e a resolução segue igual a hoje (R21):
  1. `<dir>/.git` é arquivo regular (`Lstat`; symlink recusa), de até 4 KiB, sem NUL, com uma linha `gitdir: <caminho>`.
  2. `<caminho>` (relativo ao `<dir>` quando não é absoluto) canonizado é `<C>/worktrees/<nome>`, com `<nome>` passando por `validName`, e nem `<C>/worktrees/<nome>` nem `<C>` são symlinks.
  3. `<C>/worktrees/<nome>/commondir`, resolvido relativo à pasta de administração, canoniza para o mesmo `<C>`.
  4. `<C>/worktrees/<nome>/gitdir` (ponteiro de volta), resolvido relativo à pasta de administração quando é relativo (nunca contra o diretório corrente do processo), canoniza para o `<dir>/.git` canônico.
  5. `<C>` se chama `.git` e tem `HEAD`, `objects/` e `refs/`; o checkout principal é o pai de `<C>`.
- Com worktree válido, `Resolve` e `PdocsDir` usam o checkout principal como chave: modo `repo` devolve `<principal>/.harness` e ignora a cópia versionada dentro do worktree (R19); modo global acha o nome pelo mapa como no principal (R18). O campo `pdocs` passa a usar a mesma chave (hoje usa o diretório recebido).
- `dh harness-path --json` ganha o campo aditivo `linked_worktree`: `null` fora de worktree, senão `{"main", "gitdir", "commondir"}`, todos canônicos: `main` é o checkout principal (pai de `<C>`); `gitdir` é a pasta de administração `<C>/worktrees/<nome>` (não o gitfile); `commondir` é `<C>`. Consumidor atual: o isolamento de CLI (D4), que lê a pasta de administração de `gitdir` para a R26 e para a foto da R30. Leitores antigos ignoram o campo.
- Casos-limite: repositório bare (`<C>` não se chama `.git`): não é worktree, resolução igual a hoje e isolamento recusa; submódulo (gitfile para `.git/modules/<x>`): não casa o layout, igual a hoje; gitfile symlink: recusa; worktree de worktree: o `git` grava a administração plana em `<C>/worktrees/`, então casa; subpasta de worktree: o resolvedor não sobe diretórios hoje e continua sem subir; gitfile inválido num worktree em modo `repo` com cópia de `.harness/`: cai no resultado de hoje (a cópia); o `git` também falha nesse worktree, e o Coordenador detecta pelo teste de D2.
- Motivo decisivo: ler quatro arquivos não executa nada do repositório; chamar `git` antes de validar já lê a config apontada por um gitfile qualquer.
- Descartadas: B, pela leitura de config antes da validação e pela dependência de `git` no `dh`; C, por não cumprir R18–R19.

### D4 · Isolamento: A, validação única e clone como hoje

- Passo 1 de `clone-isolation.md`, por ramo:
  - `<project>/.git` é diretório: as checagens de hoje (`test -d <project>/.git && test ! -L <project>/.git`) e nada mais.
  - `<project>/.git` não é diretório: exige `linked_worktree` não nulo em `dh harness-path --json <project>` e, no mesmo shell, logo antes do clone, `test -f <project>/.git && test ! -L <project>/.git`.
  - Qualquer outro caso é `isolation unavailable` (R27): `dh` ausente ou falho, `linked_worktree` nulo, `.git` symlink, bare, submódulo, layout inválido, conteúdo fora do formato.
- O sandbox continua `git clone --no-hardlinks <project>` em pasta descartável, sem `origin`, sem `alternates`, sem `logs/` e com as checagens de hoje (R28–R29). Nunca `git worktree` para o sandbox.
- Hipótese a provar: o clone de um worktree ligado parte do `HEAD` do worktree (não do principal). Se não partir, o passo 2 (casar o estado) já falha e a chamada termina `isolation unavailable`; a hipótese não abre brecha, só pode impedir o uso.
- Foto R5 ampliada (R30), em duas partes:
  - (a) Todo projeto: a cobertura R5 de hoje e, no diretório comum `<C>` (num checkout normal, o próprio `.git`), `refs/` como arquivos (além do `for-each-ref`) e, quando existe, a lista de entradas de `worktrees/` com os arquivos de (b) de cada entrada. Item inexistente fica gravado como "ausente" nas duas fotos; ausente igual a ausente não é diferença, nunca é erro.
  - (b) Só quando `<project>/.git` é gitfile (worktree ligado): os bytes do gitfile e, da pasta de administração `<C>/worktrees/<nome>` (lida do campo `gitdir` de `linked_worktree`), `HEAD`, `commondir`, `gitdir`, `config.worktree`, `refs/` (arquivos) e a presença de `locked`.
  - Fora da foto, como hoje: `index` cru (de `<C>` e de cada `worktrees/<nome>`), `logs/`, objetos, `FETCH_HEAD`, `ORIG_HEAD`, `COMMIT_EDITMSG`, `modules/`. O índice fica fora porque o próprio `git status` da foto o reescreve. O residual aceito "caminhos do `.git` vivo fora da foto do R5" encolhe: `worktrees/` e `refs/` como arquivos passam a ser cobertos, e ficam fora só `index` cru, `logs/`, objetos e packs, `modules/`, `FETCH_HEAD`, `ORIG_HEAD` e `COMMIT_EDITMSG` (RISK-002 segue "mitigado").
- Harness dir: hoje a foto cobre `<project>/.harness/`. De um worktree em modo `repo`, essa pasta é a cópia ignorada; a foto passa a cobrir o harness dir resolvido (o `.harness/` do checkout principal), o que mantém a cobertura de hoje. Em modo global ele segue fora da foto, como já registrado no RISK-002.
- Razão de segurança: o gitfile decide qual repositório o `git` lê. Um gitfile forjado ou corrompido pode apontar (por `..` ou caminho absoluto) para outro repositório; então o clone copiaria outro histórico, os comandos `git` do passo 2, rodados em `<project>`, carregariam a config dele (`core.hooksPath`, `filter.*.clean` acionado por `.gitattributes`, `core.sshCommand`) e a foto vigiaria o `.git` errado, deixando passar a escrita no verdadeiro. O layout exigido amarra o gitfile a um par administração/comum coerente (ponteiro de volta e `commondir` batem, nada é symlink, o comum é um `.git` de checkout) antes de qualquer `git` rodar. O clone não copia `hooks/` nem `config` da origem; o prefixo `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1` e `-c core.fsmonitor=false` continuam.
- Residuais honestos (RISK-002 fica "mitigado"): quem escreve como o mesmo usuário consegue forjar um layout inteiro coerente, e a validação não distingue isso de um worktree legítimo (é a mesma confiança que se dá ao `.git/config` de um checkout normal); a árvore do checkout principal e de outros worktrees fica fora do `<project>` da chamada e entra em "escrita fora do projeto não detectada"; janela de troca do gitfile entre a validação e o clone, detectada só se a troca persistir até a segunda foto; o harness dir global segue fora da foto.
- Descartadas: B, por duplicar a regra em shell e em Go; C, rejeitada pelo dono (PRD-028, apêndice).

### D5 · Configuração

- Chave `worktree`, valores `always` e `on-request`. `dh config set|unset worktree` grava em `<home>/config.yaml` sob a trava `config.lock` de hoje; valor inválido sai com erro listando os dois e não toca o arquivo (R1–R3). Entram: o campo em `Config`, a gravação em `SaveConfig`, `worktree` em `ConfigKeys`, o caso em `SetConfigValue`/`validConfigValue` e a mensagem de chave desconhecida.
- `dh config` e `dh config --json` mostram `worktree`, vazio quando ausente, como `language`, `mode` e `dashboard.sprites` (R4). `dh harness-path --json` leva `defaults.worktree` só com o valor global, ou `on-request` (R5). O Go repete o valor gravado de `worktree` sem validar na leitura, como já faz com `language`. `mode` é a exceção (a leitura recusa valor fora de `repo`/`global` e aborta o parse, `internal/harness/config.go` linhas 203–205), e `worktree` não a segue: um valor inválido editado à mão não pode derrubar `dh harness-path` (R8).
- Precedência (R6–R7), aplicada só pelo Coordenador: `local.yaml` do harness dir, lido como texto (linha de topo `worktree: <valor>`, R9), depois `defaults.worktree`, depois `on-request`. Valor inválido em qualquer das duas fontes é avisado ao dono e pulado (R8 estendida ao valor global editado à mão).
- Sem parser Go de `local.yaml`; `references/local.yaml.example` documenta a chave.

### D6 · Limpeza e falha

- Remoção (R22–R24), só pelo Coordenador e só depois do merge e da tag: `git -C <wt> status --porcelain` vazio; e o `HEAD` do worktree é o `headRefOid` do PR mergeado (`gh pr view <n> --json state,headRefOid`) ou `git rev-list HEAD --not --remotes` sai vazio. A primeira forma cobre o merge squash com a branch remota apagada, em que os commits não estão em nenhuma ref remota mas o GitHub guarda `refs/pull/<n>/head`. Passando, `git worktree remove <caminho>` (sem `--force`) e `git branch -D <branch>` (o squash impede `-d`; a condição acima prova que o remoto tem os commits). Falhando, nada se remove e o Coordenador relata caminho e motivo.
- Falha na criação (R17): caminho existente, branch já em uso em outro worktree, `task-id` inválido, `git worktree add` com erro ou a conferência de D2 divergente: o Coordenador relata a causa e para, sem voltar ao checkout principal sem a resposta do dono. Worktree criado que falhou na conferência fica onde está e é relatado.
- Interrupção (R25): o worktree fica; o caminho relativo e a branch já estão no registro (R16). No `/dh:resume`, o Coordenador confere que `git worktree list --porcelain` tem esse caminho com essa branch e reaproveita; se o caminho sumiu do disco mas segue listado, relata ao dono antes de qualquer `git worktree prune`. Não existe varredura automática de worktrees órfãos, ao contrário dos clones de CLI: um worktree guarda trabalho do dono.

## 4. Consequências

- Fica mais fácil: rodar um lote sem sujar o checkout principal; chamar revisor CLI de dentro de um worktree; achar harness dir e pdocs de qualquer worktree ligado.
- Fica mais difícil: o resolvedor ganha um caminho de validação de arquivos do git; a foto R5 cresce (gitfile, administração do worktree, `refs/` como arquivos); no Claude Code, o builder escreve fora do diretório da sessão e pode pedir permissão para a pasta `<nome>.worktrees` irmã (medir; configurar a pasta em `additionalDirectories` se pedir).
- Dívida assumida: R12 só em prosa até o dono decidir pôr `worktree` na lista da trava (candidato do PRD, revisão de segurança própria); repositório bare com worktrees fica sem suporte; `dh link` rodado dentro de um worktree registraria o caminho do worktree (fora do PRD-028).
- Contratos afetados: `dh harness-path --json` ganha `linked_worktree` e `defaults.worktree` (aditivos); `dh config` ganha `worktree`; `clone-isolation.md` e `live-tree-freeze.md` trocam o passo 1 e a lista da foto; `.commands/build.md` e `.commands/auto.md` ganham criação, despacho com caminho absoluto, registro e limpeza; o registro de tarefa ganha caminho relativo e branch.

## 5. Validação

- Check que prova que funciona:
  - testes Go do resolvedor com fixtures em disco: checkout normal sem mudança; worktree válido em modo `repo` com cópia de `.harness/` devolve o principal; worktree em modo global devolve o mesmo `dir` e o mesmo `pdocs`; worktree de worktree; e recusas: `.git` symlink, gitfile com `..` para outro repositório, ponteiro de volta divergente, `commondir` divergente, pasta de administração symlink, comum bare, submódulo, gitfile grande ou malformado;
  - um teste de integração com `git worktree add` real, quando `git` existir na máquina;
  - testes de `dh config` para `set`, valor inválido sem tocar o arquivo, `unset`, e a saída com `worktree` vazio quando ausente; `defaults.worktree` sem o valor do `local.yaml`;
  - matriz do `qa-verifier` no isolamento: worktree válido passa e o clone tem `.git` próprio, sem `origin` e sem `alternates`; gitfile symlink e gitfile fora do layout terminam `isolation unavailable`; escrita simulada durante a chamada em `<C>/hooks/`, em `worktrees/<nome>/config.worktree`, num arquivo de `refs/` e no gitfile descarta o veredito; prova da hipótese do `HEAD` do clone;
  - leitura do diff: cláusula anti-delegação idêntica nos prompts; nenhum arquivo chama o RISK-002 de "resolvido"; `dh validate` sem caminho privado.
  - Estado: não rodado (decisão de arquitetura não executa nada).
- Sinal de que está falhando: `mode: none` ou `pdocs` vazio dentro de um worktree de projeto ligado; revisor CLI terminando `isolation unavailable` num worktree criado pelo kit; registro gravado na cópia `.harness/` do worktree; trava de commit sem marcar depois de escrita do builder no worktree; worktree removido com mudança não enviada.

## 6. Revisar quando

- O dono puser `worktree` na lista de escrita da trava (R12 deixa de ser só prosa).
- Surgir um segundo adaptador de runtime, ou o Claude Code permitir escolher a pasta do `EnterWorktree` sem trocar o diretório da sessão.
- Alguém pedir suporte a repositório bare com worktrees.
- A caixa de areia do SO chegar e tornar a foto R5 desnecessária.

## 7. Referências

- `docs/prd/PRD-028-worktree-and-approval-page.md`
- `docs/adr/ADR-001-runtime-go.md`, `docs/adr/ADR-002-cli-planejadora.md`, `docs/adr/ADR-006-painel-e-travas-de-commit.md`, `docs/adr/ADR-007-harness-global.md`
- `internal/harness/resolve.go`, `internal/harness/config.go`, `internal/harness/config_set.go`, `cmd/dh/main.go`
- `adapters/claude-code/plugin/hooks/lib/gh-guard.ts`, `adapters/claude-code/plugin/hooks/lib/gate.ts`
- `.skills/external-clis/references/clone-isolation.md`, `.skills/external-clis/references/live-tree-freeze.md`, `.skills/external-clis/references/local.yaml.example`
