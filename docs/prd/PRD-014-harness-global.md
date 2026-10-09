# PRD-014 · `.harness` global: pasta única por máquina com config e estado dos projetos

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-08 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-08 / 2026-10-08 |
| Tickets | T-1601 (este PRD); T-1601-01 a T-1601-08 (plano em `.harness/tasks/T-1601/PLAN.md`) |

> Nota (0.24.0, PRD-017): o dashboard ganha a rota `/pdocs/<nome>/` e o campo `docs` em `/api/state`; `dh harness-path --json` ganha `pdocs`. Ver `docs/prd/PRD-017-pdocs.md` e o `CHANGELOG.md`.

## 1. Problema

Hoje o estado do dh de cada projeto vive em `<repo>/.harness/`. No trabalho os colegas não usam o dh e pediram ao dono
para não alterar o `.gitignore` dos repositórios nem deixar pasta do dh na árvore. Além disso, o `setup` e o `doctor`
recomendam commitar ou ignorar o `.harness/`, por herança de iterações anteriores; o dono quer o contrário: o `.harness/`
interno é versionado (salvo repositórios públicos que ele sinalizar), e quem não pode deixar rastro usa uma pasta global.

Fatos verificados em 2026-10-08 (leitura do repositório):

- A trava de memória casa pelo sufixo `.harness/MEMORY.md|EPOCHAL.md|RISKS.md`
  (`adapters/claude-code/plugin/hooks/runtime-guard.ts`, constante `PROTECTED`); uma pasta externa não seria protegida.
- O `doctor` olha só `<cwd>/.harness` e avisa "`.harness/` is ignored by git" (`internal/doctor/doctor.go`).
- Os snapshots de sessão ficam em `~/.claude/dev-harness/sessions` (ou `$CLAUDE_CONFIG_DIR/dev-harness/sessions`, ou
  `DEV_HARNESS_SNAPSHOT_DIR`) (`internal/snapshot/snapshot.go`, `SnapshotDir`). O token de parada do dashboard fica em
  `<UserConfigDir>/dev-harness/dashboard` (`internal/dashboard/stop.go`).
- O dashboard descobre projetos por `DH_DASHBOARD_ROOTS` e sprites por `DH_DASHBOARD_SPRITES`; `/api/state.config` expõe
  `{roots, sprites, port}` (`internal/dashboard/projects.go`, `server.go`).
- Proposta do dono, conversa de 2026-10-08, registrada em `.harness/projeto-harness-global.md`: a pasta global é criada pelo binário `dh` em Go, com `config.yaml` na
  raiz e uma pasta por projeto com a estrutura normal; ao abrir um projeto o Claude deve identificar se o `.harness` é o
  local ou o global (local tem prioridade); o nome da pasta do projeto é idêntico ao da pasta do repositório.

## 2. Solução

Existe uma pasta global por máquina e por ambiente, `~/.harness/` (`DH_HOME` sobrescreve o lugar), criada pelo `dh` no
primeiro uso e fonte da verdade do que o dh guarda fora dos repositórios:

```
<home>/                       # ~/.harness ou $DH_HOME
  config.yaml                 # padrões globais + mapa projects:
  projects/<nome>/            # MEMORY, EPOCHAL, RISKS, project.yaml, local.yaml, tasks/, prd/ (direto, sem subpasta .harness)
  sessions/                   # todos os snapshots de sessão
  dashboard/                  # token de parada do dashboard
```

- Todo projeto fica registrado em `config.yaml` (`projects:`, caminho absoluto do repositório para nome, o nome sendo a
  pasta do repositório), com `.harness/` interno ou global.
- A resolução é única (`dh harness-path`): `<repo>/.harness/` se existir; senão `<home>/projects/<nome>/` pelo mapa; senão
  nenhum, e o `setup` oferece criar. A sessão recebe o caminho resolvido ao abrir.
- Em modo global nada do dh é escrito na árvore do repositório: sem pasta, sem arquivo, sem linha de `.gitignore`.
- A trava de memória, o dashboard, o `doctor`, o `validate` e os textos de agentes e skills passam a usar a pasta resolvida.
- `.harness/` interno é versionado; `setup` e `doctor` não sugerem ignorá-lo (só o `local.yaml` mantém o aviso).
- Entrega única na 0.21.0, com quebras declaradas (R9, R12).

## 3. Regras

Nota: R17 a R20 foram acrescentadas na validação e ficam listadas depois de R13 de propósito.

Pasta global e config

- R1 `<home>` segue uma única ordem, usada pelo Go e pelos mods: `$DH_HOME`; senão `$HOME/.harness`; senão
  `%USERPROFILE%\.harness`; senão o diretório home do usuário no SO mais `.harness`. Sem nenhum deles, o home é
  desconhecido e os comandos falham com mensagem clara (nunca um caminho relativo). O `dh` cria `<home>` (e `config.yaml`,
  `projects/`) no primeiro uso, sem pedir confirmação. Cada ambiente (WSL, Windows, macOS) tem a sua pasta, nunca
  compartilhada; o kit não traduz caminhos. Verificável: teste Go com `DH_HOME` apontando para um diretório temporário
  inexistente e outro com a variável ausente (home do usuário simulado); teste com `HOME` diferente de `USERPROFILE`
  (vence `HOME`) nos dois lados, Go e mod; teste sem nenhuma fonte de home (falha com mensagem, sem caminho relativo).
- R2 O `config.yaml` é escrito e lido só pelo binário Go, em subconjunto mínimo de YAML plano, sem dependência Go nova
  (ADR-001). Chaves de padrão: `language`, `reviewers`, `mode` (modo padrão do `setup`: `repo` ou `global`) e
  `dashboard.sprites`; mais o mapa `projects:`. Os caminhos em `projects:` são gravados entre aspas; o leitor aceita um
  nível de aninhamento (`dashboard:`, `reviewers:`, `projects:`). Mods e scripts nunca leem o arquivo: perguntam ao `dh harness-path`.
  Verificável: `go.mod` sem dependência nova; teste de ida e volta (escrever, ler, comparar), incluindo uma chave de caminho Windows (com letra de unidade e barra invertida) e de arquivo inválido
  (erro legível, nada sobrescrito); `grep` nos mods sem `config.yaml`.
- R3 O `setup` copia os padrões do `config.yaml` (`language`, `reviewers`) para o `project.yaml` do projeto novo; depois
  vale o do projeto. O modo do projeto novo vem de `mode` no `config.yaml` quando definido; sem a chave, o `setup`
  pergunta ao dono uma vez, para aquele projeto (`repo` ou `global`). Verificável: teste do `setup` com `config.yaml`
  preenchido (com e sem `mode`; sem `mode` a pergunta é feita) e sem ele (cai nos padrões de hoje).

Registro e resolução

- R4 Todo projeto fica registrado em `projects:` (caminho absoluto do repositório para nome); o `setup` sempre registra e
  o `dh link` registra e corrige. O nome padrão é o nome da pasta do repositório; se já existe uma entrada com o mesmo
  nome, o `dh link` a substitui (religa, R6) quando o caminho dela já não é um diretório, e recusa exigindo nome
  explícito quando esse caminho ainda existe. O nome da pasta em `<home>/projects/` é igual ao da pasta do
  repositório, salvo escolha explícita. Verificável: teste com dois repositórios de mesma pasta-base (recusa sem nome;
  aceita com nome) e com entrada de caminho removido (substituída, religa).
- R5 `dh harness-path [dir] [--json]` é o único resolvedor e imprime o modo (`repo`, `global` ou `none`) e o diretório
  absoluto resolvido. Ordem: (1) `<repo>/.harness/` se existir (modo `repo`); (2) `<home>/projects/<nome>/` pelo mapa
  (modo `global`); (3) nenhum (`none`). Se existirem os dois, vence o interno e o `doctor` avisa. `--json` imprime um
  objeto com `mode` e `dir`. Verificável: teste de tabela com as quatro situações (só interno, só global, os dois,
  nenhum) e a saída com e sem `--json`.
- R6 `dh link [name] [--off] [dir]` registra o repositório (`dir`, padrão o diretório atual) em `projects:` e, quando
  não há `.harness/` interno, cria `<home>/projects/<name>/`. Reaproveita a pasta quando ela já existe (religa depois de
  mover ou renomear o repositório). `--off` remove a entrada do mapa e não apaga a pasta do projeto. Mover ou renomear o
  repositório deixa o mapa desatualizado: o `doctor` mostra `mode: none` e aponta `dh link`. Verificável: testes de
  criar, religar com pasta existente, `--off` preservando a pasta, e `doctor` após mover.
- R7 Rastro zero no repositório em modo global: depois de `dh link` num repositório novo sem `.harness/`, `git status
  --porcelain -uall` fica vazio, `<repo>/.harness` não existe e nenhum `.gitignore` foi tocado; o mesmo vale depois de
  uma sessão que escreve MEMORY e uma tarefa. Verificável: teste de integração com repositório git temporário.

Travas e permissões

- R8 A trava de memória (`runtime-guard.ts`) nega a escrita de subagente em MEMORY.md, EPOCHAL.md e RISKS.md nos dois
  lugares: `<repo>/.harness/` e `<home>/projects/*/`, com `<home>` pela ordem de R1 (sem home conhecido, só o padrão interno `<repo>/.harness/` é conferido, ADR-007 §4). Cada modo tem
  teste. Verificável: teste do mod nega `<home>/projects/x/MEMORY.md` e `<repo>/.harness/MEMORY.md` de um subagente e
  permite a escrita do Coordenador; arquivo com nome parecido fora dessas pastas continua permitido. Nenhum arquivo do
  `<home>` fora desse padrão ganha trava nova.
- R9 Permissões: no trabalho, a entrada `~/.harness` em `additionalDirectories` das configurações de usuário do Claude
  Code (`~/.claude/settings.json`) é o caminho documentado. O `doctor` detecta a falta (com `<home>` fora do projeto) e
  sugere a entrada; nunca grava. Verificável: teste do `doctor` com e sem a entrada; `grep` sem escrita em `settings.json`.

Sessões, dashboard e extensão

- R10 Snapshots de sessão ficam só em `<home>/sessions/` e são atribuídos ao projeto pelo mapa `projects:`
  (caminho do `cwd` da sessão dentro de um repositório registrado). `DEV_HARNESS_SNAPSHOT_DIR` e `CLAUDE_CONFIG_DIR`
  deixam de definir o destino (decidido): `DH_HOME` é o único controle, e os testes o definem. O mod
  `snapshot-writer.ts` grava em `<home>/sessions/`, com `<home>` pela ordem de R1, e
  ignora `DEV_HARNESS_SNAPSHOT_DIR` e `CLAUDE_CONFIG_DIR`. Verificável: teste de `SnapshotDir` e de atribuição por `cwd`
  (repositório raiz e subpasta); teste do mod com e sem `DH_HOME` e com `HOME` diferente de `USERPROFILE` (vence `HOME`).
- R11 O dashboard lista os projetos do `config.yaml` (modo `repo` ou `global`) e lê `dashboard.sprites` dele. O token de
  parada passa a viver em `<home>/dashboard/`. `/api/state.config` passa a ser `{home, sprites, port}` (sem `roots`).
  Verificável: teste do servidor com um projeto global (sem nada na árvore) listado e com sessão atribuída; teste do
  JSON de `config`; teste de `dh dashboard --stop` com o token no novo caminho; as barras de progresso de um projeto global vêm de
  `<home>/projects/<nome>/tasks`.
- R12 Quebras da 0.21.0 (cinco), todas no CHANGELOG com passos de migração: `DH_DASHBOARD_ROOTS` e `DH_DASHBOARD_SPRITES` são
  removidas; `/api/state.config` perde `roots` e ganha `home`; snapshots saem de `~/.claude/dev-harness/sessions` (ou
  `$CLAUDE_CONFIG_DIR/dev-harness/sessions`, ou `DEV_HARNESS_SNAPSHOT_DIR`; as duas últimas deixam de ser seletor) para `<home>/sessions/`; o token de parada sai
  de `<UserConfigDir>/dev-harness/dashboard` para `<home>/dashboard/`. Os projetos em repositório configurados antes da 0.21.0 só aparecem no dashboard depois de `dh link` (ou
  `setup`) em cada um; o `doctor` aponta isso. Um dashboard de versão anterior em execução tem de
  ser parado uma vez à mão (a versão antiga não conhece o token novo). Verificável: CHANGELOG e docs listam os cinco
  itens; `grep` sem as duas variáveis no código.
- R13 A extensão VS Code (repositório `dev-harness-vscode`) muda em versão própria com tag própria, combinada com a
  sessão `claude-dh-vscode`; fora do escopo deste PRD salvo o contrato: leitura de `/api/state` com `config`
  `{home, sprites, port}`, `projects[]` com `mode` e `harness` (R17), `dh projects --json` (R18), e fim de `dh.dashboard.roots` e das variáveis removidas em R12. Este PRD não altera nada
  naquele repositório. Verificável: nota no CHANGELOG; contrato descrito em Docs.

Contrato com a extensão

- R17 (contrato combinado com a sessão da extensão VS Code em 2026-10-08, decisão do Coordenador) `dh dashboard` relê o `config.yaml` a cada montagem do estado (editar o arquivo não exige reiniciar); só `--port` e
  o bind ficam fixos na partida. `/api/state.projects[]` lista os projetos registrados cujo caminho é diretório existente;
  caminho ausente, ou existente mas que resolve `none`, é omitido e contado no `warning` existente; `mode` é sempre `repo` ou `global`; o teto `MaxProjects` = 100 e seu aviso permanecem; caminhos
  absolutos e limpos; cada item traz `name`, `path`, `mode` (`repo|global`) e `harness` (diretório resolvido). Verificável: teste que edita o `config.yaml` com o servidor no ar e vê a mudança; teste com caminho
  inexistente e com caminho existente que resolve `none` (omitidos, `warning` os conta), 101 projetos (teto e aviso) e caminhos limpos.
- R18 `dh projects --json` (somente leitura, sem servidor) imprime `name`, `path`, `mode` (`repo|global`) e `harness`
  para o mesmo conjunto e na mesma ordem de `/api/state.projects[]` (sem `open` nem `delivered`); sem config sai com código
  0 e `[]`. `harness-path`, `link` e `projects` são só subcomandos do `dh`: nenhum comando `/dh:` novo, contagem segue
  19 + 1. Verificável: teste que compara esses quatro campos com `/api/state.projects[]`, de config ausente (`[]`, exit 0)
  e `grep` em `.commands/`.
- R19 Os leitores de código usam a pasta resolvida (R5), nos dois modos: `panel.ts` (`project.yaml`/`language`, `tasks`,
  `prd` e o `isHarnessPath` do gate de commit), `ProjectProgress` do dashboard e `checkStatusDrift` do `dh validate` (`internal/kit/status_drift.go`). Escrita em
  `<home>/projects/<nome>/` não suja o gate de commit. Verificável: testes do mod e do dashboard com projeto `repo` e
  projeto `global` (progresso e `checkStatusDrift` lidos da pasta global: projeto global com tickets em `<home>/projects/<nome>/tasks` detecta deriva, e sem `.harness/tasks` em nenhum dos dois lugares não há erro, caso da CI; escrita em `<home>/projects/<nome>/` não marca o
  gate como pendente).
- R20 Ao iniciar a sessão, o mod resolve o harness (chama `dh harness-path --json`) e injeta `mode` e `dir` no contexto da
  sessão pelo evento do mod `prompt.compose` (seção escopada à sessão; provado na sonda H2 da Fase 0 em 2026-10-08, ver
  `.harness/tasks/T-1601/TASK.md`). Com `none`, a seção aponta para `/dh:setup`. O mod resolve no primeiro `prompt.compose` da sessão e guarda o resultado enquanto o modo é `repo` ou `global`; enquanto
  for `none`, resolve de novo a cada compose (`/dh:setup` ou `dh link` valem sem reiniciar). Falha do `dh` não bloqueia
  a sessão (fail-open). Verificável: teste de motor para `repo`, `global` e `none`, para a troca de `none` para `repo`/`global`
  na mesma sessão, para o cache enquanto `repo`/`global`, e para falha do `dh`.

Doctor, setup e versionamento

- R14 `.harness/` interno é versionado: `setup` e `doctor` não sugerem ignorá-lo nem commitá-lo por obrigação; o aviso
  "`.harness/` is ignored by git" do `doctor` some para a pasta inteira e fica só para `local.yaml` (que deve ficar
  ignorado). O `doctor` mostra o modo e a pasta resolvidos (via o mesmo resolvedor de R5). Verificável: teste do `doctor`
  com `.harness/` ignorado (sem aviso), com `local.yaml` versionado (aviso) com modo `global` (mostra `mode: global` e o
  caminho), com as duas pastas presentes (modo `repo` e aviso) e com `<repo>/.harness/` presente e o caminho fora de
  `projects:` (mostra `mode: repo` e avisa que o projeto não aparece no dashboard até `dh link`).
- R15 A migração de `repo` para `global` e de volta é só documentação, em três linhas por sentido, sem comando próprio
  (nenhum `dh migrate`). Verificável: texto presente no tutorial pt-br e en.
- R16 Entrega única na 0.21.0: nenhum item de R1 a R20 sai em versão anterior. Resolução e memória em modo `repo`
  funcionam como hoje (a pasta global é criada com `config.yaml` e `projects/`, sem projetos além dos registrados). Verificável:
  suíte de regressão do modo `repo` verde; `dh validate` e `go run ./cmd/dh build` sem diferença em `dist/` além da
  versão; prova de clone limpo (o PR toca `cmd/`, `internal/` e `dist/*/bin`).

## 4. Docs

- AGENTS.md (regras: pasta global e `config.yaml`; resolvedor `dh harness-path`; trava nos dois lugares; sessões e token em
  `<home>`; `DH_DASHBOARD_*` removidas; `.harness/` interno versionado; ambientes separados; contagem de comandos)
- README.md, README.en.md (modo global, `dh link`, `dh harness-path`)
- docs/tutorial.html, docs/en/tutorial.html (uso no trabalho: instalar, `dh link`, `additionalDirectories`; migração de
  `repo` para `global` e de volta em três linhas; fim das variáveis do dashboard)
- docs/inicio-rapido.html, docs/en/quick-start.html (passo de `setup` com o modo)
- docs/roadmap.html, docs/en/roadmap.html (entrada da 0.21.0; deriva contra `CHANGELOG.md` barrada por `dh validate`)
- CHANGELOG.md (0.21.0, com as quebras de R12)
- docs/adr/ (ADR do registro único e da pasta global, se o plano o julgar necessário)
- docs/prd/PRD-012-dashboard.md e docs/prd/PRD-013-extensao-vscode.md (nota de remissão: `DH_DASHBOARD_*`, `config.roots`,
  caminho de snapshots e token mudam na 0.21.0; contrato novo da extensão em R13)
- `.skills/` e `.agents/` e `.commands/` que citam `.harness/` (`setup`, `doctor`, agentes e skills que leem MEMORY,
  tasks e prd) e `templates/en/`, `templates/pt-br/` (espelhados) quando citarem o caminho
- `.harness/projeto-harness-global.md` (nota de que a
  seção 8 virou este PRD)

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Linha no ignore global do git (`~/.config/git/ignore`) | A pasta continua na árvore, aparece no editor, morre em `git clean -fdx` e depende de cada máquina; não dá a visão única dos projetos. |
| Ponteiro no `.git/config` local (`dh.project`) | O dono fixou o mapa no `config.yaml`; o `.git/config` não é lido pelo mod e some ao reclonar. |
| Layout `projects/<nome>/.harness/` | Substituído: arquivos direto em `projects/<nome>/` (seção 8 e texto do dono). |
| Segunda pasta (`~/.dh`) ou token em `~/.config` | Uma pasta só: tudo fora dos repositórios fica em `<home>`. |
| Comando de migração `dh migrate` | Documentação em três linhas basta. |
| Compartilhar `<home>` entre WSL e Windows | Ambientes têm caminhos diferentes; o kit não traduz caminhos. |

### Riscos e decisões pendentes

- Risco: a trava de memória não cobrir a pasta global (casa hoje pelo sufixo `.harness/`) · mitigação: R8, um teste por
  modo.
- Risco: o Claude Code recusar escrever em `~/.harness` sem permissão · mitigação: R9 (`additionalDirectories`; hipótese
  não provada até a máquina do trabalho; o dono confirma no uso real).
- Risco: dashboard antigo em execução segue com o token antigo · mitigação: R12 (parar à mão uma vez).
- Risco: mover o repositório desfaz a ligação sem aviso · mitigação: R6 (`doctor` mostra `mode: none`).
- Nota de aplicação (R9 do projeto original, seção 3 de `.harness/projeto-harness-global.md`): nada é aplicado no WSL pessoal
  além do desenvolvimento e teste; o uso real é na máquina do trabalho.
- Decisões pendentes: nenhuma. Discovery fechada e aprovada pelo dono em 2026-10-08.

## Resumo executado

- **Entregue:** `~/.harness` virou a fonte da verdade da máquina. O `config.yaml` registra todos os projetos, o `.harness/` interno tem prioridade, a sessão recebe a pasta resolvida, a trava vale nos dois lugares e o dashboard e as sessões ficam num lugar só. Tudo saiu na 0.21.0.
- **Regras:**
  - R1 Honrada. A ordem DH_HOME→HOME→USERPROFILE é a mesma no Go e no mod (`TestHomeOrder`, `runtime-guard.test.ts`).
  - R2 Honrada. O leitor e gravador sem dependência nova passa no teste de ida e volta com caminho Windows. O `config.yaml` tem trava com nonce, e os mods nunca o leem.
  - R3 Honrada, com verificação estática de `setup.md` e de `project-onboarding`. O `setup` ao vivo não rodou.
  - R4 Honrada. Religar e recusar funcionam, inclusive com `EqualFold` (`TestR4_RelinkStaleRefuseLive`).
  - R5 Honrada. Os quatro casos de resolução estão em texto e em `--json` (`TestR5_ResolveFourSituations`).
  - R6 Honrada. `dh link` cria, religa e o `--off` preserva a pasta.
  - R7 Honrada. Com git real, `git status` fica vazio e não aparece `.harness` (`TestR7_ZeroTraceInRepo`). A sessão real `claude -p` deixou o repositório intocado.
  - R8 Honrada. A trava nega nos dois padrões, com `~`, sem distinguir maiúsculas e com caminhos normalizados (`runtime-guard.test.ts`).
  - R9 Honrada. O `doctor` só sugere e nunca grava. O `additionalDirectories` com `~` foi provado em `claude -p`; o uso interativo no trabalho ainda falta confirmar.
  - R10 Honrada. Os snapshots vão para `<home>/sessions`, e a sessão real gravou lá.
  - R11 Honrada. O `config` é `{home, sprites, port}` e o token fica em `<home>/dashboard` (`TestR11_GlobalProjectTokenAndConfig`).
  - R12 Honrada. As cinco quebras e a migração estão no CHANGELOG, e não sobrou nenhum `DH_DASHBOARD_*` no código.
  - R13 Honrada. O contrato foi enviado à sessão da extensão, e a extensão muda na própria release.
  - R14 Honrada. O `doctor` mostra o modo e a pasta, além dos avisos de ambos presentes, de projeto não registrado e de `local.yaml` (`TestR14_*`).
  - R15 Honrada. A migração está em três linhas por sentido no tutorial pt/en, e não existe comando `migrate`.
  - R16 Honrada. O modo repo não mudou, e a prova de clone limpo deu sha256 idêntico nos seis binários.
  - R17 Honrada. A configuração é relida a cada montagem, os projetos ausentes ou `none` são omitidos e o teto é 100 (`TestR17_*`, `-race`).
  - R18 Honrada. Há paridade entre `dh projects --json` e `/api/state` (`TestR18_ProjectsJSONParityWithAPIState`), e continuam 19 comandos.
  - R19 Honrada. O painel, o `ProjectProgress` e o `checkStatusDrift` usam a pasta resolvida (`panel.engine.test.ts`, `TestStatusDriftGlobalProject`).
  - R20 Honrada. A seção `dh:harness` cobre repo, global, none e falha (`panel.engine.test.ts`). A sessão real respondeu "mode none" e, depois do `dh link`, "mode global".
- **Tickets:**
  - T-1601-01 backend: aprovado, com dúvida em voo e revisão.
  - T-1601-02 backend: aprovado.
  - T-1601-03 backend: aprovado.
  - T-1601-04 backend: aprovado. A brecha do `~` na trava foi corrigida.
  - T-1601-05 backend (prosa): aprovado.
  - T-1601-06 backend (versão, build e clone limpo): aprovado.
  - T-1601-07 teste: verde, R1 a R18.
  - T-1601-08 teste: verde, mais a sessão real.
  - Revisões: Fable adversarial, código claude e codex, segurança com a lista completa (claude e codex).
- **Docs:** `CHANGELOG.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `README.md`, `README.en.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `AGENTS.md`, as notas em PRD-012 e PRD-013, ADR-007 e o cabeçalho de ADR-005 e ADR-006.
- **Fora:**
  - Sincronização entre máquinas e compartilhamento com colegas.
  - O comando de migração.
  - A mudança da extensão VS Code, que vai no repositório dela.
  - Limites aceitos: caminhos por symlink na trava, a dupla tomada de uma trava velha depois de um crash e a colisão NFC/NFD no macOS.
