# PRD-016 · Métricas, gráficos e leitura da memória no dashboard

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-09 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-08 / 2026-10-08 |
| Tickets | (ainda sem tickets; o plano vem depois da aprovação) |
| Origem | `docs/projeto-renova.md` (Fase 2); `docs/discovery/renova.md` §4 |
| Depende de | nada (PRD-012 e PRD-014 entregues) |
| Libera | PRD-017 (rota por nome de projeto e leitura de PRDs reutilizadas) |

## 1. Problema

O dashboard mostra, por projeto, barras de tickets concluídos por PRD aberto e a contagem de PRDs entregues (`internal/dashboard/projects.go`). Não responde quanto tempo um PRD levou, quanto tempo um ticket fica aberto, quantos incidentes o projeto acumulou nem como anda a memória. E a memória de cada projeto só é lida abrindo o arquivo à mão: o painel não aponta para ela.

Fatos verificados em 2026-10-08:

- O estado lê só o texto de `Status` de tickets e PRDs; não lê `MEMORY.md`, `RISKS.md` nem `EPOCHAL.md` (`projects.go`, `server.go:129`).
- Datas que existem: `Criado / atualizado` no cabeçalho do ticket (`templates/pt-br/TASK.md:25`), `entregue em <data>` no `Status` do PRD, incidentes datados em `RISKS.md`, lotes datados em `EPOCHAL.md`, `started_at`/`updated_at` por sessão. Não existe histórico de transição de status; `atualizado` é a última edição, não a data de conclusão.
- Custo e tokens só chegam ao snapshot pela `statusLine`, que o plugin não configura (PRD-012 §1); `snapshot-writer.ts` não grava custo. Não há fonte para métrica de custo.
- `go.mod` sem dependência (ADR-001); PRD-012 R22 e R24 recusaram SQLite; a página não depende de origem externa e roda em webview (PRD-012 R5); `/api/state` é contrato com a extensão VS Code, mudança só aditiva (PRD-014 R13, R17).
- `web/index.html` tem 108 linhas; o teto do repositório é 500 por arquivo (`AGENTS.md`).

## 2. Solução

O dashboard passa a calcular, na leitura e sem armazenamento novo (decisão 4: A, dono 2026-10-09), um bloco de métricas por projeto e a desenhar gráficos simples em SVG sem biblioteca. O cabeçalho do ticket ganha a linha `Concluído em`, única adição de dado. Cada card de projeto mostra: tarefas por PRD, dias de cada PRD (início → entrega), média de dias por ticket dentro do PRD, incidentes registrados, tamanho da memória e data da última consolidação, e um link que abre a memória do projeto na própria página. Tudo continua só leitura.

Granularidade é o dia (datas escritas à mão pelos agentes); itens sem data ficam fora das médias e são contados como "sem data". Custo e tokens ficam fora até existir fonte.

## 3. Regras

Dado novo

- R1 `templates/pt-br/TASK.md` e `templates/en/TASK.md` ganham a linha de cabeçalho `| Concluído em | AAAA-MM-DD |` (`Done on` em en), logo após `Criado / atualizado`, preenchida pelo papel que marca o ticket como concluído; `Status` continua como está. O `panel.ts` e o `dh validate` não mudam de vocabulário de status. Verificável: diff dos dois templates; fixture `status-cases.json` inalterada; `dh validate` verde.
- R2 Leitura de datas, só no Go do dashboard e nos dois idiomas de template: ticket, primeira data de `Criado / atualizado` ou `Created / updated` e a data de `Concluído em` ou `Done on`; PRD, primeira data de `Criado / atualizado` ou `Created / updated` e a data em `Status` começando por `entregue em` ou `delivered on`. Formato aceito `AAAA-MM-DD`; qualquer outro texto (inclusive os placeholders `AAAA-MM-DD` e `YYYY-MM-DD`) conta como ausente. Os arquivos vêm da pasta resolvida do projeto (`harness` de `/api/state.projects[]`, modo `repo` ou `global`) e de `docs/prd/` do repositório, como `ProjectProgress` já faz. Verificável: teste de tabela com data válida, placeholder, campo ausente e data invertida (fim antes do início → ausente e um aviso em `warnings`), com uma fixture pt-br e uma en por caso, incluindo um PRD en com `delivered on <data>`; teste com projeto `global`.

Métricas

- R3 `/api/state.projects[]` ganha o campo aditivo `metrics` com: `prds[]` (`id`, `status`, `tickets`, `done`, `started`, `delivered`, `days` ou `null`, `ticket_days_avg` ou `null`, `ticket_days_n`, `ticket_undated`), `prd_days_avg` e `prd_days_n` (só PRDs entregues com as duas datas), `incidents` (R3b), `memory` (`lines`, `bytes`, `missing`), `epochal` (`batches`, `last_at` ou `null`, R3c). PRD entregue, no bloco `metrics` e nos dois gráficos de R7, é o de `Status` começando por `entregue em` ou `delivered on`; os campos existentes `Open` e `Delivered` de `ProjectProgress` passam a seguir a mesma regra (hoje só `entregue em`, o que deixa projeto en sempre aberto). Arquivo recusado por `readCapped` (acima de 1 MiB ou não regular): `missing: false` e `lines`, `bytes` e `incidents` em `null`; o card mostra "ilegível (ver avisos)". Nenhum campo existente muda de nome ou tipo. Verificável: teste Go comparando o JSON antes e depois (campos antigos idênticos); PRD en com `delivered on` aparece em `Delivered` e nas métricas; teste com `MEMORY.md` de 1 MiB + 1 byte (`missing: false`, `null`, aviso) e o texto no DOM pelo walkthrough de R6b; `dh projects --json` segue sem `metrics` (PRD-014 R18 não muda).
- R3b Incidente é todo título `###` ou item de lista, fora de comentário HTML, cujo texto começa por um ID `[A-Z]+-\d+`, dentro da seção `## Incidentes` ou `## Incidents` do `RISKS.md`; comentários `<!-- ... -->` são removidos antes da contagem. Verificável: teste de tabela: template vazio (pt-br e en) = 0; arquivo com dois incidentes (um `###`, um item de lista) = 2; ID só dentro de comentário = 0.
- R3c Lote do EPOCHAL é todo título de nível 3 (`### `) abaixo de `## Lotes arquivados` ou `## Archived batches`. A cópia bruta vai de uma linha que começa por `<!-- INICIO-BRUTO` (ou `<!-- BEGIN-RAW`) até a próxima linha que começa por `<!-- FIM-BRUTO` (ou `<!-- END-RAW`); dentro dela nenhum título conta nem encerra a seção. `batches` é a contagem; `last_at` é a primeira data ISO 8601 nas linhas depois do último título de lote, antes do próximo `### ` ou do início da cópia bruta, ou `null`. Acima de 1 MiB o arquivo não é carregado inteiro: a leitura em fluxo percorre só as linhas de título, de data e de marcador, sem teto de tamanho e sem aviso. Verificável: teste de tabela com template vazio (pt-br e en) = 0/`null`; dois lotes com um `##` e um `###` dentro da cópia bruta = 2 e `last_at` do segundo; arquivo sintético de 1 MiB + 1 byte com três lotes = 3; o `.harness/EPOCHAL.md` deste repositório dá `batches` = 2 e `last_at` = `2026-09-30T14:08:31-03:00` (Consolidação do lote 2026-09-30-1).
- R4 `ticket_days_avg`, `ticket_days_n` e `ticket_undated` contam só tickets concluídos (`TaskStatus` = `done`): `ticket_undated` é o número de concluídos sem as duas datas válidas; ticket aberto ou bloqueado não entra em nenhum dos três. Toda média vem acompanhada do `n` e do número de itens sem data; a página mostra `n` ao lado da média e "n sem data" quando houver (decisão 5: A, dono 2026-10-09). Média com `n` = 0 é `null` e a página mostra "sem dado". Verificável: teste com dois concluídos datados, um concluído sem data e um aberto (`n` = 2, `undated` = 1, o aberto em lugar nenhum); teste com nenhum datado (`null`).
- R5 Ticket ou PRD aberto e fechado no mesmo dia conta 0 e a página escreve "< 1 dia". Verificável: teste com `Criado` igual a `Concluído em`.
- R6 `MEMORY.md` e `RISKS.md` são lidos da pasta resolvida do projeto por `readCapped` (1 MiB, arquivo regular); `EPOCHAL.md` pela leitura em fluxo de R3c (só arquivo regular). Arquivo ausente conta como `missing: true` sem aviso (RISKS ausente não prova ausência de incidentes, `profiles/base.yaml`); arquivo recusado segue R3 (`missing: false`, valores `null`, aviso). Verificável: teste com `MEMORY.md` de 1 MiB + 1 byte (`null`, aviso), com symlink nos três (`null`, aviso) e ausente (`missing: true`, sem aviso); teste em projeto `global` lendo de `<home>/projects/<nome>/`.
- R6b O card de cada projeto mostra em texto, além dos gráficos: incidentes (número), memória (linhas e KB, ou "ausente"), última consolidação (data ou "nunca") e número de lotes, `prd_days_avg` com o `n`, e, por PRD entregue, as datas de início e de entrega. Verificável: walkthrough do `qa-verifier` com Playwright headless contra um `<home>` de fixture com valores conhecidos, conferindo cada texto no DOM.

Gráficos e página

- R7 Gráficos em SVG gerado por JavaScript puro, sem biblioteca, sem CDN, sem fonte remota (PRD-012 R5 continua). Dois gráficos por projeto: barras de dias por PRD entregue e barras de tickets por PRD aberto; nada de linha do tempo por sessão nesta versão. Verificável: `grep -c "<script src=\"http" web/*` = 0; a página renderiza com `connect-src 'self'` do CSP atual; walkthrough do `qa-verifier` com Playwright headless contra um `<home>` de fixture: o número de barras e o valor de cada `<rect>` (atributo `data-value`) iguais ao JSON de `/api/state`.
- R8 O JavaScript dos gráficos vai em `internal/dashboard/web/metrics.js`, embutido por `go:embed` e servido em `GET /metrics.js`; `index.html` e `metrics.js` ficam cada um abaixo de 500 linhas. Verificável: `wc -l`; teste da rota (200, `text/javascript`, `nosniff`).
- R9 Cada gráfico é um `<figure>` com `<figcaption>` que diz em texto o que as barras mostram (por exemplo "dias por PRD entregue: PRD-012 3, PRD-014 1") e cada `<svg>` tem `aria-label`; a página funciona em 640 px e em 1280 px, como hoje. Verificável: walkthrough do `qa-verifier` com Playwright headless nas duas larguras; `figcaption` e `aria-label` presentes em cada gráfico.

Memória na página

- R10 Rota nova `GET /api/memory?project=<nome>&file=memory|risks` devolve o texto do arquivo em `text/plain; charset=utf-8`, `Cache-Control: no-store`, lido por `readCapped`, só para nome registrado em `config.yaml`, lido da pasta resolvida desse projeto (`repo` ou `global`), e só para os dois arquivos (decisão 6: B, dono 2026-10-09); nome desconhecido, arquivo fora da lista ou qualquer `/`, `\` ou `..` no parâmetro respondem 404 sem corpo que ecoe o parâmetro. Esta rota expõe a leitores locais o texto da memória e dos incidentes, exposição nova em relação ao PRD-012: o PR passa por `/dh:secure` antes do merge e a aceitação do residual é do dono (decisão 6). Verificável: testes 200 nos dois arquivos, nos dois modos; 404 para `epochal`, para `../x`, para nome não registrado e para arquivo ausente; nenhum caminho absoluto no corpo da resposta de erro; veredito de `/dh:secure` registrado no PR.
- R11 No card do projeto, "memória" abre um `<details>` que busca `/api/memory` só ao abrir, mostra o texto escapado em `<pre>` e não rebusca a cada tick de 2 s; `/api/state` não embute o texto da memória. Verificável: walkthrough do `qa-verifier` com Playwright headless contando as requisições a `/api/memory`: 0 antes de abrir, 1 logo após abrir, ainda 1 depois de três ticks de 2 s; o JSON de `/api/state` não contém o conteúdo do `MEMORY.md`.

Segurança e contrato

- R12 Só GET, loopback e os cabeçalhos do `guard()` valem para as rotas novas; nenhuma rota escreve. Verificável: `POST /api/memory` e `POST /metrics.js` → 405; `Host: exemplo.com` → 403; `GET /api/memory` com cabeçalho `Origin` de loopback → 200 (a página o envia; a rota não escreve, logo `Origin` não é critério de recusa aqui, ao contrário de `/api/stop`).
- R13 Tamanho do binário medido antes e depois em cada alvo, como em PRD-012 R23; crescimento acima de 15 % em qualquer alvo reabre a decisão de embutir o JavaScript em arquivo separado comprimido. Verificável: os dois números no PR e no CHANGELOG.
- R14 Fora de escopo: SQLite, `events.jsonl`, custo e tokens, métricas por sessão, exportação, gráficos de tempo real, qualquer edição pela página, mudanças na extensão VS Code (que só ganha campos novos para ler quando quiser). Verificável: nada disso no diff; `git diff go.mod` vazio.

## 4. Docs

- `docs/tutorial.html`, `docs/en/tutorial.html` (seção do dashboard: métricas, o campo `Concluído em`, a leitura da memória)
- `README.md`, `README.en.md` (uma linha nas capacidades do dashboard)
- `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md` (0.23.0 proposta; deriva barrada por `dh validate`)
- `docs/prd/PRD-012-dashboard.md` (nota de remissão: campo `metrics`, rotas `/metrics.js` e `/api/memory`)
- `templates/pt-br/TASK.md`, `templates/en/TASK.md` (R1)
- `.agents/` e `.skills/` que fecham ticket (`qa-verifier`, `coordinator`, `regression-testing`) ganham a frase "ao marcar concluída, preencher `Concluído em`"
- `.skills/context-handoff/SKILL.md` (passo 3 da consolidação, linha 77) nomeia os marcadores da cópia bruta de R3c (`<!-- INICIO-BRUTO` / `<!-- FIM-BRUTO` em pt-br, `<!-- BEGIN-RAW` / `<!-- END-RAW` em en) e `templates/pt-br/EPOCHAL.md`, `templates/en/EPOCHAL.md` os citam no comentário de estrutura do lote
- `AGENTS.md` (regra de PRD e tickets: a linha `Concluído em`)

## 5. Fora de escopo

- Ver R14. Também: preencher `Concluído em` nos tickets antigos deste repositório (o dono pode fazer depois, sem código; decisão 5).

## 6. Decisões do dono (decididas em 2026-10-09)

4. **Armazenamento.** A: derivar dos arquivos, granularidade de dia, sem armazenamento (decidida). B: `events.jsonl` append-only por projeto, gravado pelo Go quando detectar mudança de status (granularidade de hora, mas só enquanto algo observa; segunda fonte de verdade). C: SQLite (primeira dependência Go, cinco binários maiores em `dist/`, contraria ADR-001 e PRD-012 R22). Recomendação: A; B fica como evolução se o dia não bastar. Decidido: A (dono, 2026-10-09).
5. **Tickets antigos sem data de conclusão.** A: fora das médias, contador "sem data" (decidida). B: limpeza à mão neste repositório antes da entrega, como PRD-012 R21. C: usar `atualizado` como aproximação marcada (engana: é a última edição). Recomendação: A. Decidido: A (dono, 2026-10-09).
6. **Arquivos visíveis na página.** A: só `MEMORY.md`. B: `MEMORY.md` e `RISKS.md` (decidida; o contador de incidentes aponta para o arquivo). C: os três, `EPOCHAL.md` truncado a 1 MiB com aviso. Recomendação: B. **Residual novo a aceitar explicitamente:** qualquer processo local do mesmo usuário que alcance `127.0.0.1:<porta>` passa a ler o texto da memória operacional e dos incidentes de todo projeto registrado (hoje só lê status, caminhos e sessões). Sem token, como o resto do dashboard. Mitigação: só loopback, só GET, só nome registrado, `/dh:secure` antes do merge (R10). Se o dono não aceitar, cai para A ou a rota sai. Decidido: B (dono, 2026-10-09); risco residual aceito explicitamente pelo dono em 2026-10-09.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Datas por `git log` dos arquivos de ticket | Modo `global` não tem git na pasta do projeto (ADR-007); a fonte tem de valer nos dois modos. |
| Biblioteca de gráficos embutida (Chart.js etc.) | Dezenas de KB por binário e um terceiro nome para manter; barras simples cabem em SVG puro. |
| Link `file://` para a memória | Navegador e webview bloqueiam `file://` a partir de `http`; só a rota funciona. |
| Texto da memória dentro de `/api/state` | Cresce o JSON de 2 em 2 s sem necessidade; busca sob demanda. |

### Riscos e decisões pendentes

- Risco: datas escritas à mão erradas ou no placeholder · mitigação: R2 (inválido = ausente, aviso), R4 (`n` e "sem data" visíveis).
- Risco: `index.html` passar de 500 linhas · mitigação: R8.
- Risco: binário crescer · mitigação: R13.
- Risco (novo, não coberto pelo residual de PRD-012/0.20.0): `/api/memory` expõe a leitores locais o texto de `MEMORY.md` e `RISKS.md`, não só status e caminhos · mitigação: R10 (nome registrado, dois arquivos, 404 sem eco), R12 (loopback, só GET), `/dh:secure` antes do merge e aceite explícito do dono na decisão 6.
- Decisões pendentes: nenhuma. 4, 5 e 6 decididas pelo dono em 2026-10-09 (4A, 5A, 6B, com o risco residual de `/api/memory` aceito explicitamente).

## Resumo executado

- **Entregue:** o cartão de cada projeto no dashboard ganhou duração de tickets e PRDs (com n e "sem data"), tickets por PRD, datas de início e entrega, incidentes, tamanho da memória, última consolidação e dois gráficos em SVG próprio. Também ganhou a leitura sob demanda do `MEMORY.md` e do `RISKS.md` via `GET /api/memory`. Tudo é calculado dos arquivos, sem dependência nova, e saiu na 0.23.0.
- **Regras:**
  - R1 Honrada: `Concluído em`/`Done on` nos templates; o vocabulário de status não mudou (`TestR1_*`).
  - R2 Honrada: leitura de datas em pt-br e en, placeholders e inversão com aviso (`TestR2_*`).
  - R3 Honrada: `metrics` aditivo e `delivered on` conta como entregue. Os arquivos recusados aparecem como `null` e "ilegível (ver avisos)" (`TestR3_*` e o walkthrough).
  - R3b Honrada: incidentes contados fora dos comentários (`TestR3b_*`).
  - R3c Honrada: marcadores da cópia bruta respeitados e leitura em fluxo acima de 1 MiB. O EPOCHAL real deu 2 lotes e `last_at` 2026-09-30T14:08:31-03:00.
  - R4 e R5 Honradas: média só dos tickets concluídos; "< 1 dia" para entrega no mesmo dia (`TestR4_*`, `TestR5_*`, walkthrough).
  - R6 Honrada: `readCapped` e leitura por descritor; arquivo ausente, grande demais ou symlink tratado, inclusive no modo global (`TestR6_*`).
  - R6b Honrada: textos do cartão conferidos no Playwright, 80/80.
  - R7 Honrada: as barras e os `data-value` batem com o `/api/state`, sem nenhuma requisição externa.
  - R8 Honrada: rota `/metrics.js` embutida; arquivos com menos de 500 linhas.
  - R9 Honrada: `figcaption` e `aria-label` presentes; layout conferido em 640 e 1280 px.
  - R10 Honrada: só os dois arquivos fixos, só por nome registrado, 404 sem eco. A pasta harness simbólica é recusada e a leitura é feita por descritor (`openat`). O `/dh:secure` rodou com claude e codex, e as correções foram aplicadas.
  - R11 Honrada: `/api/memory` é chamada 0 vez antes de abrir o detalhe, 1 depois de abrir, e continua em 1 após três ciclos de atualização.
  - R12 Honrada: 405 para outros métodos, 403 para Host fora do loopback, 200 com Origin do loopback.
  - R13 Honrada: os binários cresceram de +0,4 % a +0,8 %.
  - R14 Honrada: `go.mod` sem mudança; nada de SQLite nem JSONL.
- **Tickets:**
  - T-1702-01 backend (templates e prosa): aprovado.
  - T-1702-02 backend (métricas): aprovado, com dúvida em voo (claude).
  - T-1702-03 backend (`/api/memory`): aprovado, com dúvida em voo (codex).
  - T-1702-04 frontend (gráficos e cartão): aprovado, depois da correção de escala.
  - T-1702-05 teste Go: verde.
  - T-1702-06 teste Playwright: 80/80.
- **Docs:** `templates/{pt-br,en}/{TASK,EPOCHAL}.md`, `.agents/coordinator.md`, `.agents/qa-verifier.md`, `.skills/context-handoff/SKILL.md`, `.skills/regression-testing/SKILL.md`, `AGENTS.md`, `CHANGELOG.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html` e a nota em `docs/prd/PRD-012-dashboard.md`.
- **Fora:**
  - SQLite, JSONL, métricas de custo e de sessões.
  - Exportação e edição pelo dashboard.
  - Mudanças na extensão VS Code.
  - Risco residual aceito pelo dono em 2026-10-09: um processo local do mesmo usuário lê MEMORY e RISKS pelo loopback.
  - Janela mínima de troca e restauração do diretório no macOS e no Windows, documentada.
