# PRD-012 · Dashboard local: projetos, sessões do Claude Code e limites numa página

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-07 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-07 / 2026-10-07 |
| Tickets | (ainda sem tickets; o plano vem depois) |


> Nota (0.21.0, PRD-014): `DH_DASHBOARD_ROOTS` e `DH_DASHBOARD_SPRITES`, `config.roots`, o caminho dos snapshots de sessão e o do token de parada mudaram. Os projetos vêm do `config.yaml` da pasta global (`~/.harness`) e `/api/state.config` é `{home, sprites, port}`. Ver `docs/prd/PRD-014-harness-global.md` e o `CHANGELOG.md`.
>
> Nota (0.23.0, PRD-016): `/api/state` ganha o campo `metrics` por projeto e o dashboard ganha as rotas `/metrics.js` e `GET /api/memory`. Ver `docs/prd/PRD-016-metricas-dashboard.md` e o `CHANGELOG.md`.

## 1. Problema

O dono trabalha com vários projetos que usam o kit e vários terminais do Claude Code abertos ao mesmo tempo. Hoje não há um
lugar único para ver (1) em que pé está cada projeto, (2) qual terminal está trabalhando, parado ou esperando uma resposta
dele, e (3) quanto resta dos limites de uso do Claude. O painel embutido do Claude Code é por sessão e some ao fechar; a
barra de status do mod (`adapters/claude-code/plugin/hooks/panel.ts`) mostra só o projeto da sessão atual.

Fatos verificados em 2026-10-07 (leitura do repositório):

- O snapshot por sessão (`internal/snapshot/snapshot.go`, `schema: 1`, estados `active`, `idle`, `closed`) é escrito por
  hooks de comando (`adapters/claude-code/plugin/hooks/hooks.json`, `dh snapshot event`). Não existe estado "esperando o
  dono", não existe marca de vida (uma sessão que morreu sem `SessionEnd` fica `active` ou `idle` para sempre) e os limites
  de 5h/semana só chegam ao snapshot pela `statusLine`, que o plugin não configura (o `settings.json` do plugin só traz
  `subagentStatusLine`).
- O mod `panel.ts` já lê `.harness/tasks/*/TASK.md` (função `taskStatus`), já lê `$.session.usage()` (limites 5h e semana)
  e já tem uma trava em `git commit` (gate pendente). Nenhum dos dois grava arquivo.
- O progresso em disco está sujo: dos `Status` de ticket em `.harness/tasks/*/TASK.md`, dezenas dizem "pronta" (e há "provisória"/"ready"),
  mesmo para trabalho entregue. Dos PRDs, só o PRD-009 e o PRD-011 dizem "entregue em"; os demais dizem outra coisa (por
  exemplo "aprovado"), mesmo os já entregues. Barras de progresso calculadas hoje mostrariam quase tudo aberto.
- Hipótese (não verificada): o dono tem o jevmon (`~/projetos/jevmon`, Go, Windows) com um avatar por estado e uma pasta de
  sprites; essa pasta deve servir sem alteração.

## 2. Solução

O comando `dh dashboard` sobe um servidor local e serve uma página única. A mesma página abre no navegador e dentro de uma
view de webview da extensão VS Code (a extensão em si é o próximo PRD). A página mostra:

1. **Projetos.** Cada subpasta com `.harness/` dentro das pastas monitoradas, com uma barra de progresso por PRD ainda
   aberto (tickets concluídos sobre tickets do PRD) e uma linha de contagem dos PRDs entregues, expansível para a lista dos seus ids.
2. **Sessões abertas.** Cada projeto com terminal do Claude Code aberto mostra o estado: trabalhando, parado ou esperando a
   resposta do dono. Sessão que morreu sem fechar não aparece como aberta.
3. **Avatar no canto superior direito.** Pose escolhida pelo pior estado entre as sessões (esperando > erro > trabalhando >
   concluído > atenção > parado), com os limites de uso do Claude (5h e semana) em barras de progresso logo abaixo.

Os dados de sessão passam a ser gravados pelo mod (e não mais por hooks de comando), e o progresso é protegido contra
status velho por uma checagem no `dh validate` e por uma trava de mod no `git commit`. Tudo é só leitura na página; nada se
edita pelo dashboard.

## 3. Regras

Superfície e servidor

- R1 A entrega é um subcomando novo `dh dashboard` no binário `dh` existente, em Go, só com a biblioteca padrão (`net/http`)
  e a página embutida por `go:embed`. Nenhuma dependência nova em `go.mod`. Verificável: `git diff go.mod go.sum` vazio;
  `dh dashboard --help` existe; a página é servida sem arquivo externo ao binário.
- R2 O servidor escuta só em `127.0.0.1` e recusa requisição cujo `Host` não seja loopback (`127.0.0.1`, `localhost`).
  Verificável: teste Go com `Host: exemplo.com` recebe 4xx; `ss -ltn` (ou equivalente) mostra só `127.0.0.1`.
- R3 O servidor só aceita `GET`. Nenhum endpoint altera arquivo, projeto, sessão ou configuração. Verificável: teste Go com
  `POST`/`PUT`/`DELETE` recebe 405 em todas as rotas.
- R4 A página lê um único endpoint JSON por `fetch` a cada cerca de 2 s. Sem websocket e sem SSE. Verificável: leitura do
  código; teste do endpoint devolve JSON válido com os três blocos (projetos, sessões, limites).
- R5 A página não depende de origem externa (sem CDN, fonte ou imagem remota) e funciona dentro de `iframe`/webview: o
  servidor não envia `X-Frame-Options: DENY|SAMEORIGIN` nem `Content-Security-Policy: frame-ancestors` restritivo, e a
  página usa só caminhos relativos. Verificável: teste Go sobre os cabeçalhos; a página abre num `iframe` de uma página
  local de teste (verificação do `qa-verifier` com Playwright headless) e renderiza os três blocos.
- R6 `/dh:dashboard` só pode existir como comando registrado pelo mod (`$.command.register` mais `$.process.run`), que sobe
  `dh dashboard` desacoplado e devolve a URL. Não entra arquivo em `.commands/`. Verificável: `ls .commands` continua com
  19 arquivos; o mod registra o comando; `dh validate` verde.

Projetos e progresso

- R7 As pastas monitoradas vêm da variável de ambiente `DH_DASHBOARD_ROOTS` (nome proposto, ver decisão 8), caminhos
  separados por `;`. Um projeto é uma subpasta direta de uma raiz que contém `.harness/`. Sem a variável ou com valor vazio,
  o dashboard mostra aviso legível e zero projetos, sem erro. Verificável: teste Go com duas raízes separadas por `;`,
  uma raiz inexistente (ignorada) e uma subpasta sem `.harness/` (não listada).
- R8 O estado de cada ticket vem de `<projeto>/.harness/tasks/*/TASK.md` (campo `Status` da tabela ou `**Status:**`) com a
  mesma classificação de `taskStatus` em `panel.ts` (concluído, bloqueado, aberto). Verificável: tabela de casos comum
  (fixture) lida pelos testes Go e pelo teste `.ts` do mod, com o mesmo resultado nos dois.
- R9 Cada ticket liga a um PRD pela primeira ocorrência de `PRD-\d+` no campo de PRD do ticket, lido de qualquer um dos dois
  nomes: `Story / PRD` (formato antigo) ou `PRD (RF-<n>)` (formato atual, `templates/pt-br/TASK.md`); ver contradição C1. Valor
  que começa com "fora" ou "nenhum" não liga. Ticket sem PRD ligável não entra em barra. Verificável: fixtures com
  `| Story / PRD | PRD-003 ... |` e com `| PRD (RF-<n>) | PRD-011 (R1, R2) |` ligam ao PRD certo; `fora do PRD-008 (...)` não liga.
- R10 Uma barra por PRD aberto: tickets concluídos sobre tickets ligados, mais a contagem de bloqueados. PRD cujo cabeçalho
  `Status` começa com `entregue em` mostra 100% e vai para a linha de contagem de entregues (com lista expansível dos ids, só leitura), qualquer que seja o status dos
  tickets. Os PRDs são procurados em `docs/prd/PRD-*.md` e em `.harness/prd/PRD-*.md`. PRD aberto sem ticket mostra "sem
  tickets" e não 0%. Verificável: fixtures com os três casos e testes Go.

Sessões e limites

- R11 O mod passa a gravar `~/.claude/dev-harness/sessions/<session_id>.json` por `$.fs.write` (mesmo diretório e mesma
  convenção de `snapshot.SnapshotDir`, incluindo `DEV_HARNESS_SNAPSHOT_DIR` e `CLAUDE_CONFIG_DIR`). O arquivo continua
  legível pela extensão do PRD-001 sem mudança dela (decidido pelo dono em 2026-10-07): campos existentes mantidos e `state` segue o vocabulário do schema 1 (`active`, `idle`, `closed`); os estados do dashboard vão num campo novo `activity` (`idle`, `working`, `waiting`, `error`, `done`) e o campo `activity_at` (RFC 3339) guarda quando `activity` mudou pela última vez; campo novo sobe `schema` para 2; o leitor Go aceita 1 e 2 (sem `activity`, deriva de `state`: `active` = working).
  A escrita preserva campos que o mod não conhece (por exemplo `tasks`, gravado por `dh snapshot subagents`) e o leitor
  tolera JSON parcial (ignora o arquivo naquele ciclo). Verificável: teste do mod com arquivo pré-existente contendo `tasks`
  e campo desconhecido; teste Go com JSON truncado.
- R12 Os hooks de comando `dh snapshot event` saem de `hooks.json`; a chave `modules` fica. O subcomando `dh snapshot event`
  permanece no binário (compatibilidade e outros adaptadores). `subagentStatusLine` fica em `settings.json`, fora de
  escopo. Verificável: `hooks.json` sem `type: command`; `claude plugin validate` no pacote construído passa.
- R13 O mod grava `activity` (e `state`) pela tabela de eventos abaixo, no estilo do jevmon. Verificável: um teste do mod por
  linha da tabela, mais uma sessão real que dispare um pedido de permissão.

  | Evento do mod | `activity` | `state` (schema 1) |
  |---|---|---|
  | SessionStart | idle | idle |
  | UserPromptSubmit, atividade de ferramenta | working | active |
  | classic.Notification com `permission_prompt`, classic.PermissionRequest | waiting | active |
  | StopFailure | error | idle |
  | Stop | done, decai para idle após 5 min (fixo; sobrescrevível por flag de `dh dashboard`, decisão 10) | idle |
  | SessionEnd | idle | closed |

  `waiting` sai na próxima atividade (prompt, ferramenta, `Stop`). O mod grava `activity_at` só quando `activity` muda; o
  batimento (R15) renova só `updated_at`. O leitor decai `done` para `idle` quando agora - `activity_at` >= 5 min. Teste Go:
  `updated_at` recente com `activity_at` de 6 min atrás e `activity: done` lê `idle`.
- R13b O mod produz todo campo que `dh snapshot event` escreve hoje: `session_name`, `started_at` (só na primeira vez), `cwd`,
  `agent`, `model.id` e o anel `events` de 50 entradas em início/fim de agente. Verificável: um teste do mod por evento
  comparando o snapshot gerado com a saída de `dh snapshot event` para a mesma entrada (campos de tempo excluídos).
- R14 Os limites 5h e semana vêm de `$.session.usage()` (`rateLimits`, tipos `five_hour` e `seven_day`, como em
  `panel.ts`) e são gravados no snapshot (`rate_limits`) a cada atualização do mod. O limite é da conta, não da sessão: o
  dashboard mostra o valor do snapshot mais recente e a idade dele ("há N min"). Sem nenhum snapshot com limite, a área
  mostra "sem dado". Verificável: teste Go com dois snapshots de idades diferentes escolhe o mais novo.
- R15 Vivacidade: `updated_at` (já existe no schema 1) é renovado a cada evento do mod e, se H3 for verdadeira, por batimento
  periódico. Sessão cujo `updated_at` é mais velho que o limite (decisão 3) não aparece como aberta;
  `closed` nunca aparece. A sessão é atribuída ao projeto cujo caminho contém o `cwd` do snapshot (igual ou ancestral).
  Verificável: teste Go com snapshot velho, snapshot `closed`, e `cwd` em subpasta do projeto.

Avatar

- R16 A pasta de sprites usa o formato do jevmon: `<pose>.png` ou quadros `<pose>_00.png`, `<pose>_01.png`..., e
  `poses.json` opcional mapeando estado para pose (mesmo formato do jevmon, `fps` de 0 a 60, inválido ignorado com aviso).
  Sem `poses.json`, vale o mapa padrão do jevmon (README dele, linha 136): esperando = `duvida`, erro = `puto`, trabalhando =
  `celular`, concluído = `mostrando`, atenção = `apontando`, parado = `frente`.
  A pasta vem de `DH_DASHBOARD_SPRITES` (nome proposto). Verificável: a pasta de sprites real do dono, sem alteração,
  renderiza o avatar; teste Go com `poses.json` inválido cai no padrão.
- R17 O kit embute um conjunto mínimo de exemplo, só com os estados que o dashboard usa (parado, trabalhando, esperando o
  dono, erro, concluído, atenção), com os nomes de pose do mapa de R16 (`duvida`, `puto`, `celular`, `mostrando`, `apontando`, `frente`) e licença compatível com a MIT do kit e registrada em `THIRD_PARTY_NOTICES`. A queda é por nome de pose: pose
  ausente na pasta do dono cai no exemplo do kit de mesmo nome; sem pasta configurada, usa só o exemplo. O servidor serve sprite só por
  nome de pose conhecido (sem caminho livre). Verificável: teste Go com pasta vazia, pasta parcial e nome `../x` (404).
- R18 Prioridade de estado (copiada do jevmon): esperando > erro > trabalhando > concluído > atenção > parado. O estado do
  avatar é o de maior prioridade entre as sessões abertas e o alerta de limite (atenção quando 5h ou semana >= 80%, mesmo
  corte do `panel.ts`). Verificável: teste Go da função de prioridade com cada par vizinho.

Confiabilidade do progresso

- R19 `dh validate` ganha uma checagem de deriva com duas condições: (a) ticket ainda aberto cujo PRD tem `Status` começando
  com `entregue em`; (b) PRD cujos tickets ligados estão todos concluídos (e são ao menos um) mas cujo `Status` não começa
  com `entregue em`. Ambas são erro, seguindo o precedente da deriva roadmap contra CHANGELOG em `internal/kit/doc_counts.go`.
  Verificável: uma fixture falha por condição, uma limpa passa; `dh validate .` verde neste repositório depois da limpeza (R21).
- R20 O mod ganha uma trava em `git commit` (no estilo da trava de gate em `panel.ts`): se o repositório sendo commitado
  tem a deriva de R19 (qualquer das duas condições), nega o commit com mensagem acionável (lista os tickets e diz o que corrigir), no idioma de
  `project.yaml`. Falha aberta se o mod não conseguir ler (como as demais travas). Verificável: teste `.ts` com deriva
  (nega, mensagem lista os tickets), sem deriva (não nega), erro de leitura (não nega).
- R21 Limpeza única, neste repositório, antes de R19/R20 valerem aqui (decisão 7, decidido pelo dono em 2026-10-07): tickets de PRD
  entregue passam a `concluída` e o `Status` dos PRDs já entregues passa a `entregue em <data>`, com datas tiradas do
  `CHANGELOG.md`. A lista de PRDs e datas é mostrada ao dono antes de gravar. Verificável: `dh validate .` verde e a barra
  de cada PRD aberto no dashboard bate com `grep` dos `Status`.

Entrega

- R22 Sem SQLite. A única fonte de verdade são os arquivos: snapshots em `~/.claude/dev-harness/sessions/` e
  `.harness/` dos projetos. Verificável: `go.mod` sem dependência de banco.
- R23 Linguagem Go para o servidor (ADR-001); TypeScript só nos mods. Na primeira construção, medir o tamanho de cada `dh`
  (hoje cerca de 3,5 MB) e registrar o valor no relatório.
- R24 Fora de escopo: SQLite, OTEL e Agent SDK, runtimes que não são Claude Code, a implementação da extensão VS Code,
  acesso remoto e autenticação, qualquer edição de dados pelo dashboard.

## 4. Docs

- README.md, README.en.md (versão, comando novo, contagens)
- docs/tutorial.html, docs/en/tutorial.html (como subir o dashboard, variáveis de ambiente, sprites)
- docs/roadmap.html, docs/en/roadmap.html, CHANGELOG.md (fim do lote, pelo docs-guide; a deriva roadmap contra CHANGELOG é barrada por `dh validate`)
- AGENTS.md (contagem "os 19 comandos" e a lista de `dh` "validate, build, doctor, snapshot", que ganha `dashboard`; a regra
  da extensão VS Code, que passa a ter o contrato de embutir a página; a regra de integração "mod antes de hook de comando",
  já presente, deve ser citada na ADR)
- docs/adr/ (ADR nova: servidor `dh dashboard`, snapshot escrito pelo mod, schema 2; cita ADR-001, ADR-005, ADR-006)
- docs/prd/PRD-001-etapa-1.md (RF-03 e RF-04: o snapshot deixa de vir de hooks de comando; só nota de remissão)
- THIRD_PARTY_NOTICES (sprites de exemplo e sua licença)
- Contagem de comandos: um comando registrado pelo mod não entra em `.commands/`, então `harness-manifest.json` e a checagem
  `checkDocumentedCounts` seguem em 19; o comando visível ao usuário passa a ser 20. A Docs deve dizer isso nos dois idiomas
  (ver decisão 9).

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| SQLite para sessões e progresso | Decisão do dono: confiabilidade vem de gravação automática pelo mod mais checagens de deriva; banco é outra peça para sincronizar. |
| Websocket ou SSE | Decisão do dono: JSON por `fetch` a cada cerca de 2 s basta para uso local. |
| TypeScript para o servidor | ADR-001: scripts do kit são Go; TypeScript só em mod que roda no motor. |
| Manter `dh snapshot event` em hooks de comando | Regra vigente do `AGENTS.md`: mod sempre que o motor oferecer o evento. O mod passa a gravar o snapshot. |
| Comando `/dh:dashboard` em `.commands/` | Decisão do dono: só como comando de mod; evita mudar a contagem de 19 comandos de arquivo. |
| Condição extra da deriva: PRD citado como entregue no `CHANGELOG.md` mas sem `entregue em` | Acopla o `validate` ao texto do CHANGELOG; a condição (b) de R19 já cobre o caso comum. |
| Reabrir o painel de status do mod para o dashboard | A barra de status é por sessão; o dashboard cobre várias sessões e projetos. |

### Hipóteses

- H1 (não verificada) `$.fs.write` consegue gravar fora do projeto, em `~/.claude/dev-harness/sessions/`. O plano deve
  testar isto primeiro (sonda em scratchpad). Como cairia: a API recusa caminho fora do projeto. Plano B na decisão 1.
- H2 (não verificada) O mod recebe `classic.Notification`, `classic.PermissionRequest` e `classic.StopFailure`, com `notification_type`, como o
  jevmon trata nos hooks. Como cairia: o motor não expõe esses eventos ao mod; sem `Notification`/`PermissionRequest` não há `waiting`, sem `StopFailure` não há `error` (o avatar nunca chega a esses estados; o plano decide o recuo com o dono).
- H3 (não verificada) O mod tem como disparar escrita periódica (batimento) ou ler o `pid` do processo. Como cairia: não há
  temporizador no mod; só eventos atualizam `updated_at`.
- H4 (não verificada) `$.command.register` e `$.process.run` permitem subir um processo desacoplado que sobrevive à
  sessão. Como cairia: o processo morre com a sessão; o comando só imprime o que rodar.
- H5 (não verificada) O webview do VS Code aceita um `iframe` para `http://127.0.0.1:<porta>` (ou o `asExternalUri` resolve).
  Como cairia: o contrato de R5 não basta e a extensão precisa de proxy; isso é problema do PRD da extensão.

### Riscos e decisões

- Risco: o `snapshot subagents` (`subagentStatusLine`) e o mod gravam o mesmo arquivo; leitura-modificação-escrita
  concorrente pode perder campos · mitigação: R11 (preservar campos desconhecidos, gravar o arquivo inteiro de uma vez,
  leitor tolera arquivo parcial); alternativa se aparecer perda: o mod grava em arquivo próprio e o leitor une.
- Risco: o `.ts` da trava (R20) duplica a lógica Go de R19 e as duas podem divergir · mitigação: fixture compartilhada
  lida pelos dois testes (R8).
- Risco: a trava de commit (R20) bloqueia qualquer commit enquanto a deriva existir, inclusive commits que nada têm a ver
  com ela · mitigação: R21 limpa antes; a mensagem lista o que corrigir; falha aberta se a leitura falhar.
- Risco: tamanho do binário cresce com HTML e sprites embutidos; se o `dh` dobrar ou mais, reabrir a divisão do binário · mitigação: R23 mede na primeira construção; sprites de
  exemplo pequenos e sem animação.
- Risco: API de mod em acesso antecipado muda entre versões do Claude Code (herdado do PRD-011) · mitigação:
  `claude plugin test` no pacote em cada release.
- Risco: a condição (b) de R19 exige `entregue em` no mesmo commit que fecha o último ticket; a revisão final da feature é um ticket próprio, mantido aberto até a entrega · mitigação: o Coordenador fecha esse ticket e grava `entregue em` no mesmo commit.
- Risco: a página é só leitura, mas lê `.harness/` de vários projetos e mostra caminhos locais · mitigação: R2 e R3
  (loopback, `Host` validado, só `GET`).

Decisões (todas decidido pelo dono em 2026-10-07, aceitando as recomendações do Coordenador; numeração original)

1. Plano B se `$.fs.write` não escrever fora do projeto (H1): o mod grava em `<projeto>/.harness/sessions/` e o dashboard lê de cada projeto monitorado. Consequência aceita: sessão em pasta sem `.harness/` não aparece; precisa de `.gitignore`; a extensão do PRD-001 muda de caminho. Só vale se a sonda de H1 falhar.
2. Porta: padrão fixo, com `--port` para trocar e erro claro se ocupada.
3. Vivacidade: (a) batimento do mod a cada 30 s e limite de 2 min se H3 for verdadeira; senão (b) sem batimento, limite de 15 min, com o limite como flag de `dh dashboard`.
4. Animação no conjunto de exemplo: só PNGs estáticos.
5. Imagens de exemplo: ícones vetoriais simples feitos no kit (SVG convertido para PNG), com os nomes de pose de R17; o dono pode trocar os arquivos depois sem mudar código.
6. (removida: o leitor aceita os dois nomes do campo, ver R9 e C1; a numeração das demais fica.)
7. Limpeza: tickets mais `Status` dos PRDs entregues, datas do `CHANGELOG.md`, lista mostrada ao dono antes de gravar (R21).
8. Variáveis `DH_DASHBOARD_ROOTS` e `DH_DASHBOARD_SPRITES`; busca só em filhos diretos da raiz.
9. Contagem de comandos: manter "19 comandos" (arquivos de `.commands/`) e acrescentar "mais 1 comando registrado pelo mod" nos textos, nos dois idiomas.
10. Decaimento de `done` para `idle`: fixo em 5 min, sobrescrevível por flag de `dh dashboard` (R13).

Contradições encontradas entre as decisões e o repositório

- C1 Os tickets em `.harness/tasks` usam dois nomes para o campo de PRD: `Story / PRD` (67 arquivos, formato antigo, por
  exemplo T-001) e `PRD (RF-<n>)` (28 arquivos, formato atual, `templates/pt-br/TASK.md`; contagem do Coordenador). O leitor
  aceita os dois, sem renomear nada. O regex `PRD-\d+` serve, mas `.harness/tasks/T-1001/TASK.md` tem "fora do PRD-008" e
  ligaria ao PRD errado sem a exceção da R9.
- C2 Só PRD-009 e PRD-011 dizem "entregue em"; PRD-001 a PRD-008 e PRD-010 dizem algo diferente de "entregue em" (por exemplo "aprovado"). Como D3 manda 100% só para
  "entregue em", a limpeza só de tickets (opção a+b da conversa) deixaria esses PRDs como abertos.
- C3 O `settings.json` do plugin só tem `subagentStatusLine`; o snapshot de limites 5h/semana hoje só chega por uma
  `statusLine` que o plugin não instala. A D2 resolve ao gravar de `$.session.usage()`, mas o schema 1 já tem o campo
  `rate_limits` que nunca é preenchido pelo plugin.
- C4 `AGENTS.md` lista "validate, build, doctor, snapshot" como a superfície de `dh` e "os 19 comandos" como contagem do kit;
  ambos mudam (R6, decisão 9).
- C5 O `PRD-001-etapa-1.md` (RF-03, RF-04) e a regra da extensão em `AGENTS.md` descrevem o snapshot como produzido por
  hooks do plugin; D2 muda a origem, e o contrato com `danielmalka/dev-harness-vscode` (caminho e schema) precisa seguir
  igual: por decisão do dono em 2026-10-07, `state` mantém o vocabulário do schema 1 e o campo novo `activity`
  carrega os estados do dashboard (schema 2 aditivo), então a extensão segue sem mudança.

Tickets: nenhum nesta etapa; o plano começa pela sonda de H1 a H4.

## Resumo executado

- **Entregue**: `dh dashboard` serve uma página só leitura em `127.0.0.1:4747` com os projetos de `DH_DASHBOARD_ROOTS`, uma barra por PRD aberto, as sessões abertas com estado e o avatar com os limites 5h e semana; o mod do plugin grava os snapshots no schema 2, registra `/dashboard` e nega commit com deriva de status; `dh validate` barra a mesma deriva.
- **Regras**:
  - R1 Honrada: stdlib `net/http` + `go:embed`, `go.mod` sem diff (QA T-1308).
  - R2 Honrada: `tcp4 127.0.0.1`, Host fora de `127.0.0.1`/`localhost` → 403 (`TestHostGuard`, `ss -ltn`, curl).
  - R3 Honrada: não-GET → 405 com `Allow: GET` (`TestMethodsAreGetOnly`).
  - R4 Honrada: JSON consultado a cada 2 s, sem websocket/SSE (Playwright).
  - R5 Honrada: sem `X-Frame-Options`/`frame-ancestors`, zero requisição externa, página renderiza dentro de iframe (Playwright).
  - R6 Honrada: `/dashboard` registrado pelo mod (o host não expõe `/dh:dashboard`, verificado com `claude -p`); `.commands/` segue com 19.
  - R7 Honrada: `TestProjects`, `TestProjectsAbsDedupeAndAllFail`; 6 projetos reais listados.
  - R8 Honrada: fixture compartilhado Go↔TS e `TestStatusCasesTSCopyMatches`.
  - R9 Honrada: `TestPRDLink` e casos TS, os dois nomes de campo e "fora do PRD-008".
  - R10 Honrada: `TestProjectProgress`, `TestReviewFilesDoNotEatThePRDCap`; entregues recolhidos.
  - R11 Honrada: schema 2 aditivo, campos desconhecidos preservados em Go e TS; escrita do mod não atômica aceita como residual (`$.fs` sem rename).
  - R12 Honrada: `hooks.json` só com `modules`; `subagentStatusLine` mantido.
  - R13 Honrada: um teste TS por linha da tabela; `waiting` via `PermissionRequest` provado, `Notification` assinado sem prova em sessão real.
  - R13b Honrada: teste de paridade com a saída real de `dh snapshot event`.
  - R14 Honrada: limites por campo com idade; o batimento renova `rate_limits` (correção da revisão).
  - R15 Honrada: batimento de 30 s (H3 provada), corte de 2 min, `--stale`.
  - R16 Honrada: pasta real de sprites do jevmon renderizou sem alteração (QA).
  - R17 Honrada: 6 ícones do kit (MIT), fallback por pose, `../x` → 404, leitura guardada (só arquivo regular, ≤ 1 MiB, sem seguir link).
  - R18 Honrada: `TestAvatarPriorityPairs`; alerta só com limite de até 5 h.
  - R19 Honrada: `TestStatusDrift`; `dh validate .` verde.
  - R20 Honrada: testes TS e prova em clone temporário com sessão real (limpo passa, deriva nega).
  - R21 Honrada: 9 PRDs marcados `entregue em` e 51 tickets `concluída`, lista aprovada pelo Coordenador por delegação do dono.
  - R22 Honrada: sem banco de dados.
  - R23 Honrada: Go no servidor, TS só no mod; binários de ~3,5 para ~7,3–8,1 MB, mantidos num binário só (decisão do Coordenador).
  - R24 Honrada: nada fora de escopo entrou.
- **Tickets**: T-1301 (backend, harness-maintainer) concluída, sondas H1/H3/H4 verdadeiras, H2 parcial, H5 not-run; T-1302 (backend) concluída; T-1303, T-1304, T-1306, T-1307 (backend, harness-maintainer) concluídas, dúvida no meio do build em cada uma; T-1305 (frontend, builder) concluída; T-1308 (teste, qa-verifier) APPROVED R1–R23; T-1309 (docs-guide) concluída; T-1310 revisão final: código aprovado na 3ª rodada (claude; codex na 2ª), segurança sem vulnerabilidade demonstrada (claude e codex; grok not-run por falha de transporte).
- **Docs**: README.md, README.en.md, docs/tutorial.html, docs/en/tutorial.html, docs/roadmap.html, docs/en/roadmap.html, CHANGELOG.md, AGENTS.md, docs/prd/PRD-001-etapa-1.md (nota de remissão), THIRD_PARTY_NOTICES.md.
- **Fora**: implementação na extensão VS Code (próximo PRD; H5 not-run); SQLite; OTEL; runtimes não-Claude; acesso remoto e autenticação. Residuais aceitos: leitura local de `/api/state` sem token, porta 4747 ocupável por outro programa, trava olha só o diretório da sessão.
