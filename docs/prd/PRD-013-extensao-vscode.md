# PRD-013 · Extensão VS Code: dashboard embutido, status na barra e lançador de sessão

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-08 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-07 / 2026-10-07 |
| Tickets | T-1401..T-1407 (`.harness/tasks/`) |


> Nota (0.21.0, PRD-014): `DH_DASHBOARD_ROOTS` e `DH_DASHBOARD_SPRITES`, `config.roots`, o caminho dos snapshots de sessão e o do token de parada mudaram. Os projetos vêm do `config.yaml` da pasta global (`~/.harness`) e `/api/state.config` é `{home, sprites, port}`. O contrato novo da extensão (R13 do PRD-014) é `config` `{home, sprites, port}`, `projects[]` com `mode` e `harness` e `dh projects --json`; ela muda em versão própria. Ver `docs/prd/PRD-014-harness-global.md` e o `CHANGELOG.md`.

## 1. Problema

O dono trabalha no VS Code com vários projetos do kit e vários terminais do Claude Code. O dashboard local (PRD-012,
entregue em 0.18.0/0.18.1) mostra projetos, sessões e limites, mas vive numa aba do navegador, fora do editor. Para abrir uma
sessão do Claude Code num projeto e rodar `/dh:<comando>`, o dono ainda abre terminal, navega até a pasta, digita `claude` e
digita o comando.

Fatos verificados em 2026-10-07:

- O repositório `danielmalka/dev-harness-vscode` (clone em `~/projetos/dev-harness-vscode`) está vazio: commit inicial e
  `LICENSE`. A extensão nasce do zero, em TypeScript, sem binário Go próprio.
- `AGENTS.md` do dev-harness fixa o papel: a extensão é só visualizador e lançador, nunca runtime nem orquestradora;
  painel v1 só ao vivo e só Claude Code; OTEL e Agent SDK adiados.
- O dashboard (PRD-012) responde em `http://127.0.0.1:4747` (padrão, `--port`), só `GET`, valida o hostname do `Host` (`127.0.0.1` ou
  `localhost`, qualquer porta; `[::1]` é recusado), não envia `frame-ancestors` restritivo, não usa origem externa e tem layout de uma coluna abaixo de 640 px
  (cabe numa view estreita). Expõe `/api/state` em JSON. Sobe por `/dashboard` no Claude Code ou por `dh dashboard --detach`.
  O binário `dh` fica no plugin instalado, em `~/.claude/plugins/cache/dev-harness/dh/<versão>/bin/<os>_<arch>/dh`
  (`dh.exe` no Windows). Variáveis: `DH_DASHBOARD_ROOTS` (separadas por `;`) e `DH_DASHBOARD_SPRITES`.
- O snapshot (schema 2) traz `state` (`active`/`idle`/`closed`), `activity` (`idle|working|waiting|error|done`),
  `activity_at`, `rate_limits` (`five_hour`, `seven_day`), `cwd`, `session_name`, `agent`, `model`.
- Hipótese (não verificada, H5 do PRD-012): o webview do VS Code aceita um `iframe` para `http://127.0.0.1:<porta>`.
  O ambiente principal do dono é Linux e WSL no Windows (caso Remote WSL: o servidor roda no WSL, o webview no host).

## 2. Solução

Uma extensão VS Code, distribuída como `.vsix`, que faz três coisas e nada além:

1. **View lateral com o dashboard.** Um ícone próprio na barra de atividades abre uma view de webview que embute a página do
   dashboard (`http://127.0.0.1:<porta>`, porta configurável, padrão 4747). O dono pode arrastar a view para a barra
   lateral secundária ou para o painel; o VS Code lembra. Se o dashboard não responde, a view mostra o botão "Iniciar
   dashboard", que roda `dh dashboard --detach` do plugin instalado (com as raízes e a pasta de sprites das configurações
   passadas como variáveis de ambiente) e carrega a página quando ela responder.
2. **Item na barra de status.** Texto curto com o estado do avatar e os limites 5h/semana, lidos de `/api/state`. Clicar
   revela a view. Dashboard parado: o item diz "dashboard parado".
3. **Lançador mínimo.** O comando de paleta "dh: abrir sessão" pergunta a pasta do projeto (pastas do workspace e projetos
   sob as raízes) e um comando `/dh:<nome>` da lista dos 19, abre um terminal integrado nessa pasta rodando `claude` e
   digita o comando. Nada além disso é enviado e o terminal nunca é lido.

Não há sobreposição flutuante (a API do VS Code não tem), nem chat, nem controle de sessão.

## 3. Regras

View e dashboard

- R1 A extensão declara um contêiner de barra de atividades com ícone próprio e uma view de webview. Verificável: o
  `package.json` tem `contributes.viewsContainers.activitybar` e `contributes.views`. Automatizado: `@vscode/test-electron`
  sob `xvfb-run` confere o registro da view. Manual pelo dono (VS Code Linux e Remote WSL): a view abre pelo ícone e pode
  ser movida para a barra secundária e para o painel, com capturas de tela no ticket da sonda.
- R2 Com o dashboard respondendo, a view mostra a página dele (os blocos de projetos, sessões e avatar com limites) dentro do
  webview, sem reimplementar nada em TypeScript. Verificável: automatizado, `@vscode/test-electron` sob `xvfb-run` com
  dashboard de fixture confere a URL do `iframe` (R3); manual pelo dono (VS Code Linux e Remote WSL) confirma os três
  blocos visíveis, com capturas de tela no ticket da sonda.
- R3 A porta vem da configuração `dh.dashboard.port` (inteiro 1024 a 65535, padrão 4747). Valor inválido cai no padrão com
  aviso. A URL embutida no `iframe` vem de `vscode.env.asExternalUri(http://127.0.0.1:<porta>)`, construída num único
  lugar; o CSP do webview lista em `frame-src` só essa origem resolvida. As requisições feitas pela extensão (sondagem,
  status, proxy) vão sempre a `http://127.0.0.1:<porta>`. Nenhuma outra origem é carregada. Verificável: teste unitário da
  validação da porta; teste do CSP gerado; leitura do código mostra uma única construção de URL.
- R4 Dashboard inacessível (sem resposta em até 2 s numa sondagem `GET /api/state`) faz a view mostrar mensagem legível e o
  botão "Iniciar dashboard"; nunca uma página em branco nem erro cru. Verificável: teste com porta sem servidor.
- R5 "Iniciar dashboard" executa o `dh` do plugin instalado (a versão mais nova em
  `~/.claude/plugins/cache/dev-harness/dh/<versão>/bin/<os>_<arch>/dh`, `dh.exe` no Windows) com `dashboard --detach`
  e `--port` igual à configuração, passando `DH_DASHBOARD_ROOTS` e `DH_DASHBOARD_SPRITES` a partir de `dh.dashboard.roots`
  (lista, juntada com `;`) e `dh.dashboard.sprites`. Plugin não encontrado: mensagem acionável ("instale o plugin
  `dh@dev-harness`"); a extensão não baixa nem embute binário. Verificável: teste unitário da localização (versão mais nova,
  ausência) e do ambiente montado, com `child_process` simulado.
- R6 Depois de iniciar, a view tenta carregar a página até 10 s (sondagem a cada 500 ms) e, passado o prazo, volta à mensagem
  de R4 com o motivo. Verificável: teste com servidor que sobe depois de 1 s e com servidor que nunca sobe.
- R7 Se a hipótese H5 falhar (webview não carrega o `iframe` local), a extensão usa o proxy pelo host da extensão (decisão
  pendente 3 define a forma), mantendo R2 a R4. Verificável: o ticket de sonda registra o resultado de H5; se falha, R2 é
  provada pelo caminho do proxy.

Barra de status

- R8 O item da barra de status mostra, em pt-br, o estado do avatar e os limites 5h e semana (por exemplo
  `dh · trabalhando · 5h 42% · sem 18%`), atualizado a cada cerca de 5 s a partir de `GET /api/state`. Mapeamento dos ids de `avatar.state` devolvidos por `/api/state` (ids do jevmon, `internal/dashboard/sessions.go`; rótulos da página, `AVLABEL`): `esperando` = `esperando você`, `erro` = `erro`, `trabalhando` = `trabalhando`, `concluido` = `concluído`, `atencao` = `atenção`, `parado` = `parado`; valor desconhecido mostra o id cru. Percentuais arredondados para inteiro. Limite
  sem dado mostra `sem dado`. A idade do limite não aparece no texto (pode aparecer no tooltip). Verificável: teste
  unitário da formatação, com uma resposta real de `/api/state` capturada de `dh dashboard` como fixture, para cada um dos seis ids, valor desconhecido, arredondamento e limite ausente; captura de tela.
- R9 (decisão pendente 1; recomendação assumida: só `/api/state`) Dashboard parado: o item mostra `dh · dashboard parado`,
  sem erro e sem notificação. A extensão não lê snapshots em `~/.claude/dev-harness/sessions/` diretamente (fonte única é
  `/api/state`). Se o dono escolher ler snapshots, R9 e R15 mudam. Verificável: teste com servidor ausente; `grep`
  no código sem acesso a esse diretório.
- R10 Clicar no item revela a view do dashboard (foco nela). Verificável: teste de integração com o comando do item.

Lançador

- R11 (decisões pendentes 4 e 6; recomendação assumida) O comando de paleta `dh.openSession` ("dh: abrir sessão") pergunta primeiro a pasta (pastas do workspace aberto e
  subpastas diretas com `.harness/` das raízes de `dh.dashboard.roots`) e depois o comando, de uma lista fixa de 19
  entradas `/dh:<nome>` (auto, build, consolidate-memory, discover, doctor, document, fix, handoff, improve, plan-loop,
  plan, refactor, release, resume, review, secure, setup, understand, verify). Sem workspace aberto e sem raízes
  configuradas, mostra aviso com atalho para as configurações e não abre terminal. Cancelar em qualquer pergunta não abre
  terminal. Verificável: teste com seletor simulado: cancelar na pasta, cancelar no comando, escolha completa, e o caso sem
  workspace e sem raízes (aviso, nenhum terminal criado).
- R12 (decisão pendente 5; recomendação assumida: sem Enter automático) Escolha completa executa exatamente esta sequência:
  `createTerminal({ cwd: <pasta> })`; espera pela integração de shell (`onDidChangeTerminalShellIntegration`, ou cerca de 1 s se ela não vier); `sendText("claude", true)`; espera fixa de 3 s (H4); `sendText("/dh:<nome>", false)`. Se o dono escolher Enter automático, o
  segundo envio usa `true`. A extensão não envia nenhum outro texto e não lê a saída do terminal. Verificável: o mock de
  `vscode.window.createTerminal` registra exatamente essas duas chamadas de `sendText` com seus valores de `addNewLine`,
  e o `cwd`; `grep` sem `onDidWriteTerminalData` nem equivalente.
- R13 O nome do projeto na lista e o `cwd` passam por validação: só caminhos que são pasta existente e vêm do workspace ou de
  uma subpasta direta das raízes; nome de comando só da lista fixa, nunca texto livre. Verificável: teste com caminho
  inexistente, `../x` e comando fora da lista (recusados).

Entrega e limites

- R14 O `package.json` declara `extensionKind: ["workspace"]`: a localização do `dh`, a sondagem, o status e o terminal
  rodam no lado do workspace (no Remote WSL, dentro do WSL, onde vivem `dh`, `claude` e o plugin). Fato, não decisão:
  `dh`, `claude` e o plugin estão no lado WSL. Verificável: `package.json`; o ticket da sonda confirma, em Remote WSL, que o
  terminal abre no WSL e que o `dh` achado é o do `~` do WSL.
- R15 A extensão não contém binário, servidor, lógica de snapshot nem reimplementação do dashboard; TypeScript apenas, sem
  dependência de runtime além de `vscode` e da biblioteca padrão do Node. Verificável: `package.json` sem `dependencies` de
  runtime; `.vsix` sem arquivos `.exe`/ELF.
- R16 CI do `dev-harness-vscode` roda lint, checagem de tipos, testes unitários e empacota o `.vsix` a cada push; a cada tag
  anotada `v*` o CI anexa o `.vsix` ao release do GitHub. Verificável: workflow verde em um push; release de teste com o
  `.vsix` anexado e instalável por "Install from VSIX" numa janela limpa.
- R17 Fora de escopo: publicação no Marketplace, OTEL, Agent SDK, interface de chat, controle ou envio de entrada a sessões
  (além de R12), runtimes que não são Claude Code, reimplementar o dashboard em TypeScript, notificações e toasts.
  Verificável: revisão final do PRD contra o diff.

## 4. Docs

- AGENTS.md (a regra da extensão VS Code passa a dizer: view que embute o dashboard, item de status lido de `/api/state`,
  lançador de sessão e distribuição por `.vsix` em release do GitHub; Marketplace adiado. A cláusula "Lê snapshots em
  `~/.claude/dev-harness/sessions/<session_id>.json`" sai, condicionada à decisão pendente 1 com a recomendação assumida;
  se o dono escolher ler snapshots, ela fica)
- README.md, README.en.md (menção da extensão e do link para o repositório)
- docs/roadmap.html, docs/en/roadmap.html (menção; a deriva contra `CHANGELOG.md` é barrada por `dh validate`)
- docs/tutorial.html, docs/en/tutorial.html (seção curta: instalar o `.vsix`, configurações, abrir o dashboard e uma sessão)
- CHANGELOG.md (dev-harness) e CHANGELOG.md do dev-harness-vscode (novo)
- dev-harness-vscode/README.md (novo: instalação por `.vsix`, configurações `dh.dashboard.port|roots|sprites`, limites)
- docs/prd/PRD-001-etapa-1.md (nota de remissão: o contrato da extensão passa a este PRD)

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Reimplementar o dashboard em TypeScript dentro da extensão | Duplica a página; AGENTS.md pede só visualizador e a página já existe (PRD-012). |
| Sobreposição flutuante do avatar | A API do VS Code não tem; decisão do dono. |
| Binário Go próprio na extensão | Decisão do dono: usa o `dh` do plugin; evita segundo binário para versionar e assinar. |
| Ler snapshots direto para o item de status quando o dashboard cai | Duas fontes para o mesmo dado; o item diz "dashboard parado" (recomendação, decisão pendente 1). |
| Marketplace já na v1 | Decisão do dono: adiado; `.vsix` no release basta. |

### Hipóteses

- H1 (não verificada) O webview aceita `iframe` para `http://127.0.0.1:<porta>` (H5 do PRD-012). Como cairia: o VS Code bloqueia
  a origem por CSP do webview ou o servidor fica inalcançável no host (Remote WSL); R7 troca para proxy. Primeiro ticket: sonda no
  VS Code do dono, em Linux e em Remote WSL (precisa do dono por alguns minutos; o agente não tem VS Code).
- H2 (não verificada) `vscode.env.asExternalUri` resolve a porta local em Remote WSL e em Remote SSH. Como cairia: devolve
  um host encaminhado não loopback, e o Host guard do dashboard (valida o hostname) responde 403; ou devolve URI
  inacessível. Nos dois casos cai no proxy de R7.
- H3 (não verificada) O plugin instalado em escopo de usuário basta para `claude` reconhecer `/dh:<comando>` em qualquer pasta,
  sem flag extra. Como cairia: o comando não existe na sessão; a extensão passaria `--plugin-dir` (precisaria de nova
  configuração e decisão do dono).
- H4 (não verificada) Uma espera fixa de 3 s depois de `sendText("claude", true)` deixa o `claude` pronto para receber o segundo `sendText` sem perder caracteres (a espera anterior pela integração de shell, ou cerca de 1 s sem ela, só garante o prompt do shell). Como cairia: o texto chega antes do prompt do `claude`; o atraso vira configuração ou o lançador deixa só `claude` aberto.
- H5 (não verificada) A localização da versão mais nova do `dh` pelo nome da pasta (ordenação semântica) basta. Como cairia:
  o plugin muda a estrutura de cache; a extensão falha com mensagem de R5.

### Riscos e decisões pendentes

- Risco: Remote WSL/SSH põe o servidor num host e o webview em outro · mitigação: H2 e a sonda; R7 como plano B.
- Risco: o cache do plugin muda de caminho entre versões do Claude Code · mitigação: R5 trata ausência com mensagem; caminho num
  único módulo; configuração futura `dh.binaryPath` só se aparecer (YAGNI agora).
- Risco: iniciar o dashboard com variáveis diferentes das do `/dashboard` do Claude Code gera dois servidores concorrentes na mesma
  porta · mitigação: o erro de porta ocupada do `dh dashboard` é mostrado; a sondagem de R4 evita iniciar se já responde.
- Risco: o texto digitado no terminal segue a sessão `claude` com permissões do dono · mitigação: R12 e R13 (lista fixa, sem
  texto livre, sem leitura de saída).
- Risco: a extensão e o `dh` evoluem em repositórios diferentes; mudança de `/api/state` quebra o item de status · mitigação:
  campos lidos listados no README da extensão; falha de parse mostra `dashboard parado` e nunca lança.

Decisões pendentes (todas com recomendação; nenhuma bloqueia a sonda H1). Decisões 1, 4, 5 e 6: recomendação assumida pelo Coordenador na ausência do dono, 2026-10-07 (assumida pelo Coordenador na ausência do dono, 2026-10-07); o dono pode reverter depois. Decisões 2 e 3 seguem abertas (2 fora de escopo; 3 só se a sonda falhar).

1. Item de status com o dashboard parado: mostrar "dashboard parado" (só `/api/state`) ou ler snapshots direto? Recomendação: só
   `/api/state`; mantém uma fonte e a extensão fina. Responde: dono. Bloqueia: R9.
2. Notificação (toast) quando uma sessão entra em `waiting`: agora ou depois? Recomendação: depois, em PRD próprio, depois de
   usar a v1. Responde: dono. Bloqueia: nada (R17 já exclui).
3. Forma do proxy se H1 falhar: (a) `fetch` no host da extensão e HTML do dashboard reescrito para o webview; (b) servidor HTTP
   local da extensão encaminhando para o dashboard. Recomendação: decidir só se a sonda falhar; inclinação a (a), sem segundo
   servidor. Responde: dono, depois da sonda. Bloqueia: R7.
4. Origem da lista de comandos do lançador: lista fixa de 19 no código (R11) ou lida do plugin instalado? Recomendação: fixa, com
   teste que compara com `.commands/` do dev-harness em CI para detectar deriva. Responde: dono. Bloqueia: R11.
5. Envio do comando: digitar e pressionar Enter automaticamente, ou deixar o texto digitado sem Enter para o dono conferir?
   Recomendação: digitar sem Enter (dono confirma); o dono sempre decide executar. Responde: dono. Bloqueia: R12.
6. Seleção de projeto sem workspace aberto e sem raízes configuradas: mostrar aviso ou pedir pasta por diálogo do sistema?
   Recomendação: aviso com atalho para as configurações; diálogo só se pedido. Responde: dono. Bloqueia: R11.

Tickets: nenhum nesta etapa. O plano começa pela sonda de H1/H2 no VS Code do dono (Linux e Remote WSL), antes de qualquer
código de view. A extensão tem CI próprio e tag anotada própria por trabalho concluído.

## Resumo executado

- **Entregue**: a extensão `dev-harness-vscode` (v0.1.0, corrigida na v0.1.1), distribuída como `.vsix` na release do GitHub, embute o dashboard numa view lateral, mostra estado e limites na barra de status a partir de `/api/state` e abre sessão com "dh: abrir sessão", que digita `/dh:<comando>` sem Enter; confirmada pelo dono num VS Code real (Remote WSL e Windows local).
- **Regras**:
  - R1 Honrada: `viewsContainers`/`views` no `package.json`, registro conferido sob `xvfb-run`; view aberta e movida pelo dono (T-1406).
  - R2 Honrada: projetos, sessões e avatar com limites visíveis na view, em captura de tela enviada pelo dono no Remote WSL; no Windows local, confirmado pelo dono sem captura. Na v0.1.0 a página ficava em "carregando…": com `enableScripts: false` o VS Code tira `allow-scripts` do frame de conteúdo e o `iframe` herda; corrigido na v0.1.1 (`enableScripts: true`, CSP da view sem `script-src`).
  - R3 Honrada: teste da porta e do CSP; URL montada só em `src/dashboardUrl.ts`.
  - R4 Honrada: teste com porta sem servidor; mensagem e botão "Iniciar dashboard".
  - R5 Honrada: testes da localização do `dh` (versão mais nova, ausência) e do ambiente montado.
  - R6 Honrada: testes com servidor que sobe depois de 1 s e que nunca sobe (prazo único de 10 s).
  - R7 Honrada: H1 verdadeira depois da v0.1.1 (o `iframe` local carrega no webview); o proxy não foi necessário.
  - R8 Honrada: teste de formatação com resposta real de `/api/state` para os seis ids, id desconhecido, arredondamento e limite ausente; item da barra de status visível nas capturas do dono nas duas janelas (Remote WSL e Windows local; nesta, a extensão ainda estava na v0.1.0, que já lia `/api/state` corretamente).
  - R9 Honrada: teste com servidor ausente; o código não acessa `~/.claude/dev-harness/sessions`.
  - R10 Honrada: teste de integração do comando do item.
  - R11 Honrada: testes de cancelamento, escolha completa e caso sem pasta nem raízes; lista fixa dos 19 com teste de deriva contra `.commands/`.
  - R12 Honrada: mock registra as duas chamadas de `sendText` e o `cwd`; o dono rodou o lançador e o comando apareceu digitado depois de o `claude` abrir, com o atraso fixo de 3 s do código (H4).
  - R13 Honrada: testes com caminho inexistente, `../x`, symlink fora da raiz e comando fora da lista.
  - R14 Honrada: `extensionKind: ["workspace"]`; no Remote WSL o terminal abriu no WSL e o dashboard foi lido do lado WSL (dono).
  - R15 Honrada: `package.json` sem `dependencies`; `.vsix` com 19 arquivos, sem binário.
  - R16 Honrada: CI verde em push e PR; tags `v0.1.0` e `v0.1.1` anexaram a `.vsix` às releases; o dono instalou a `.vsix` da v0.1.0, e a v0.1.1 entrou por `code --install-extension` e por atualização manual do dono.
  - R17 Honrada: revisão final contra o diff; sem notificação proativa (mensagens só em resposta a uma ação do dono).
- **Tickets**: T-1401 (infra, builder) concluída; T-1402, T-1403, T-1404 (frontend, builder) concluídas, revisão de código por ticket; T-1405 (teste, qa-verifier) APPROVED R1–R17, 60 testes unitários estáveis e integração sob xvfb; T-1406 (sonda do dono) concluída: H1 a H4 verdadeiras depois da v0.1.1 (H2 no Remote WSL; Remote SSH não testado); H5 coberta só por teste unitário, porque o botão "Iniciar dashboard" não foi exercido na sonda (o dashboard já estava no ar); T-1407 (docs, docs-guide) concluída; revisão final FEATURE APROVADA; segurança sem vulnerabilidade demonstrada (claude e codex) na v0.1.0 e na v0.1.1.
- **Docs**: AGENTS.md, README.md, README.en.md, docs/roadmap.html, docs/en/roadmap.html, docs/tutorial.html, docs/en/tutorial.html, CHANGELOG.md, docs/prd/PRD-001-etapa-1.md (nota de remissão), e no dev-harness-vscode README.md e CHANGELOG.md.
- **Fora**: Marketplace, OTEL, Agent SDK, chat, controle de sessão e runtimes não-Claude. Residuais aceitos: a view só percebe um dashboard subido por fora ao abrir, mudar configuração ou clicar no botão; leitura de `/api/state` sem teto de tamanho (timeout de 2 s); um processo que ocupe a porta antes do dashboard roda JavaScript dentro do painel, com o alcance da própria origem; numa janela local do Windows a view lê o dashboard do WSL pelo loopback espelhado (`networkingMode=mirrored`), e o botão "Iniciar dashboard" não deve conseguir subir o `dh` do WSL (esperado por R5, não exercido).
