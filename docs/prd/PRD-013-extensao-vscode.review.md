# PRD-013 · Validação

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

Documento: `docs/prd/PRD-013-extensao-vscode.md`. Validador: `document-validator` (opus), skill `document-review`.

## Rodada 1 de 2 (2026-10-07)

- Status da revisão: `changes required`
- Achados bloqueantes: 4

## Achados

### BLOCKING: P1 [conflict] R3 fixa a origem `127.0.0.1:<port>` enquanto H2/R7 dependem de `asExternalUri`/do encaminhamento do Remote WSL, que pode render outra origem
- Local: §3 R3; Apêndice H2; risco Remote WSL; §1
- Reportado por: claude
- Problema: R3 proíbe qualquer origem que não seja `127.0.0.1:<port>`; `asExternalUri` pode devolver `localhost:<other port>` ou um host encaminhado. A trava de `Host` do dashboard (`internal/dashboard/server.go:104-111`) confere só o hostname, então portas localhost encaminhadas passam, mas um host encaminhado que não seja loopback recebe 403. O `frame-src` do CSP do webview não é nomeado.
- Correção sugerida: a URL embutida vem de `asExternalUri(http://127.0.0.1:<port>)`, construída num lugar só; o `frame-src` do CSP do webview lista só essa origem; as requisições do lado da extensão (sonda, status, proxy) vão para `http://127.0.0.1:<port>`; o "Como cairia" de H2 nomeia o caminho do 403 até R7.

### BLOCKING: P2 [gap] Não está dito onde a extensão roda sob Remote WSL; R4, R5, R8, R12 dependem disso
- Local: §3 R4, R5, R8, R12; §1
- Reportado por: claude
- Problema: `dh`, `claude` e o plugin ficam no lado WSL; uma extensão de UI (lado Windows) procuraria no home do Windows, iniciaria o dashboard errado e abriria um terminal do Windows.
- Correção sugerida: regra que declare `extensionKind: ["workspace"]`; localização, sonda, status e terminal rodam no lado remoto; verificado no `package.json` e no ticket da sonda.

### BLOCKING: P3 [conflict] R9 proíbe ler instantâneos; o AGENTS.md diz que a extensão os lê; §4 não remove essa cláusula
- Local: §3 R9, R14; §4; decisão pendente 1
- Reportado por: claude
- Problema: R9 está escrita como assentada enquanto a decisão 1 está em aberto; a atualização do AGENTS.md não remove "Lê snapshots".
- Correção sugerida: tornar R9 condicional à decisão 1 (recomendação assumida) e fazer §4 dizer que a regra do AGENTS.md perde a cláusula dos instantâneos.

### BLOCKING: P4 [ambiguity] R12 não decide se o Enter é pressionado
- Local: §3 R12; H4; decisão pendente 5
- Reportado por: claude
- Problema: `Terminal.sendText` acrescenta uma quebra de linha por padrão; a sequência exata (e como o `claude` inicia) está indefinida, então a checagem não pode ser escrita.
- Correção sugerida: sequência exata: `createTerminal({cwd})`, `claude` + Enter, depois, após o mecanismo de H4, `/dh:<nome>` com `addNewLine: false` (recomendação da decisão 5; `true` se o dono escolher Enter automático); o mock registra as duas chamadas.

### NON-BLOCKING: P5 [factual error] §1 descreve a trava de `Host` como presa à porta
- Reportado por: claude
- Correção sugerida: "valida o hostname do `Host` (`127.0.0.1` ou `localhost`, qualquer porta; `[::1]` recusado)".

### NON-BLOCKING: P6 [ambiguity] R8 não mapeia os valores de `/api/state` para texto de exibição
- Reportado por: claude
- Correção sugerida: tabela dos seis valores de `avatar.state` para texto, arredondamento para inteiro, idioma, se a idade do limite é mostrada.

### NON-BLOCKING: P7 [weak criterion] As checagens de R1/R2 nomeiam atores que não conseguem rodá-las
- Reportado por: claude
- Correção sugerida: checagem manual pelo dono no VS Code Linux e no Remote WSL (capturas de tela no ticket da sonda); `@vscode/test-electron` automatizado sob `xvfb-run` para o registro da view e a URL do iframe.

### NON-BLOCKING: P8 [organization] As decisões pendentes 4 e 6 bloqueiam R11, escrita como decidida
- Reportado por: claude
- Correção sugerida: marcar R11 "decisões pendentes 4 e 6; recomendação assumida" e acrescentar o caso sem workspace/sem raízes como checagem de R11.

## Rodada 2 de 2 (2026-10-07)

- Status da revisão: `changes required`
- Achados bloqueantes: 1 (P9, novo); P1–P8 resolvidos. Limite atingido.

## Achados

### BLOCKING: P9 [conflict] R8 mapeia chaves em inglês, mas `/api/state` devolve `avatar.state` como ids do jevmon em português
- Local: §3 R8
- Reportado por: claude
- Problema: `avatar.state` traz `esperando`, `erro`, `trabalhando`, `concluido`, `atencao`, `parado` (`internal/dashboard/sessions.go:73-81`; tabela da página `AVLABEL`, `internal/dashboard/web/index.html:46`). As chaves em inglês são os valores de `activity` da sessão. P6 da rodada 1 causou isso.
- Correção sugerida: mapear os ids do jevmon (como `AVLABEL`); valor desconhecido mostra o id cru; o teste usa como fixture uma resposta real de `/api/state` capturada de `dh dashboard`.

### NON-BLOCKING: P10 [ambiguity] R12 e H4 posicionam a espera pela integração de shell de forma diferente
- Reportado por: claude
- Correção sugerida: esperar a integração de shell (ou ~1 s sem ela), `sendText("claude", true)`, espera fixa de 3 s (H4), `sendText("/dh:<nome>", false)`.

## Disposição do Coordenador (2026-10-07)

- P9 é erro de fato com correção única, conferida pelo Coordenador em `internal/dashboard/sessions.go:73-81`. Com o teto de 2 rodadas atingido e o dono ausente sob autorização de levar a demanda até o fim, a correção foi aplicada pelo autor sem terceira rodada; o ponto fica listado para o dono no relatório final. P10 aplicado junto.
