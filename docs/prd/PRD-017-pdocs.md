# PRD-017 · Documentação HTML por projeto (pdocs) servida pelo dashboard

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-09 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-08 / 2026-10-08 |
| Tickets | T-1703-01 a T-1703-06 (plano em `.harness/tasks/T-1703/PLAN.md`) |
| Origem | `docs/projeto-renova.md` (Fase 3 e Premissa); `docs/discovery/renova.md` §5 |
| Depende de | PRD-016 (rota por nome de projeto registrado; leitura de status de PRD reutilizada; mesmos arquivos `server.go` e `index.html`) |
| Libera | nada |

## 1. Problema

Cada projeto que usa o kit tem PRDs, tickets e memória, mas nenhuma leitura de conjunto: por que o projeto existe, como rodar, o que está planejado e o que já foi feito, o changelog e os diagramas de macro e de domínio. O dono quer isso em HTML, no padrão visual do kit, fora do repositório (os colegas não querem rastro do kit, ADR-007), e acessível pelo dashboard.

Fatos verificados em 2026-10-08:

- A skill `doc-template-html` gera um HTML autônomo com Google Fonts no `<head>` (`SKILL.md:56`, `scripts/stamp.sh:209`) e Mermaid em `<pre class="mermaid">` por CDN, "a única exceção de CDN" (`SKILL.md:58`); o `stamp.sh` não emite a tag do Mermaid e a origem do CDN não está fixada em lugar nenhum. `AGENTS.md` repete "Mermaid via CDN é a única exceção": Google Fonts contradiz a regra.
- Tipos da skill (`references/types.md`): `catalogo`, `plano`, `feature`, `melhoria`, entre outros; nenhum tipo "projeto".
- `docs-guide` já é despachado ao fim de todo lote (`AGENTS.md`; `.commands/document.md`).
- Em modo `repo`, `<home>/projects/<nome>/` não existe (`.commands/setup.md:16`); com as duas pastas presentes o `doctor` avisa (PRD-014 R14).
- O `guard()` do dashboard aplica `Content-Security-Policy: default-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'` a toda rota (`internal/dashboard/server.go:119`): um HTML com Mermaid CDN e Google Fonts servido por ele teria script e fonte bloqueados.
- Um documento estampado tem cerca de 21 KB (`assets/modelo.html`).

## 2. Solução

Todo projeto registrado pode ter uma pasta de documentação, `<home>/projects/<nome>/pdocs/` (decisão 7: A, dono 2026-10-09), sempre na pasta global, mesmo quando o projeto usa `.harness/` no repositório. Dentro dela, um `index.html` no padrão da skill, tipo novo `projeto`, com motivação, como rodar local, planejado vs desenvolvido (lido dos PRDs), changelog, diagrama macro, diagramas por domínio e bounded contexts quando existirem; páginas extras ao lado quando um domínio pedir. Quem escreve é o `docs-guide` com a skill; o Go só serve.

O dashboard serve essa pasta em `GET /pdocs/<nome>/<arquivo>` e mostra no card do projeto o link "documentação" quando `index.html` existe. Nessa rota o CSP admite exatamente as duas origens externas que os documentos usam, Mermaid e Google Fonts (decisão 8: A, dono 2026-10-09); todas as outras rotas seguem como hoje. Quando o script do Mermaid não carrega (sem rede), o diagrama mostra a própria fonte em texto, legível.

Os pdocs nascem sob pedido (`/dh:document pdocs`) e, uma vez existentes, são atualizados pelo `docs-guide` ao fim de todo lote, como já acontece com a documentação do repositório (decisão 9: A, dono 2026-10-09).

## 3. Regras

Pasta e resolução

- R1 A pasta é `<home>/projects/<nome>/pdocs/`, com `<home>` pela ordem de PRD-014 R1 e `<nome>` o nome registrado em `config.yaml`; vale nos dois modos (`repo` e `global`) e nunca em `<repo>/`. `dh harness-path --json` ganha o campo aditivo `pdocs` com esse caminho absoluto (vazio quando o projeto não está registrado); `mode` e `dir` não mudam. Verificável: teste Go do campo nos dois modos e com projeto não registrado (vazio); prova em fixture de modo `repo`: criar `pdocs/index.html` pelo caminho devolvido e ver `git status --porcelain -uall` vazio no repositório.
- R1b `/dh:document pdocs` obtém a pasta por `dh harness-path --json` (`pdocs`), para com a mensagem "run dh link" quando o campo vem vazio e lembra, nos dois modos, que `~/.harness` precisa estar em `additionalDirectories` do `~/.claude/settings.json` (PRD-014 R9; nunca grava). O limite "Use relative paths. Do not depend on a personal home directory" de `.agents/docs-guide.md:55` e, em `.commands/document.md`, a regra "persisted at the authorized path under `.harness/`" (Output) e o limite "use relative paths only" (Limits, linha 28) ganham a exceção explícita para o assunto `pdocs`: o caminho vem do `dh`, nunca é digitado nem inferido pelo agente. Verificável: texto nos dois arquivos; prova do lote: `/dh:document pdocs` neste repositório (modo `repo`) escreve em `<home>/projects/dev-harness/pdocs/`.
- R2 Pasta "só pdocs" tem uma definição: `<home>/projects/<nome>/` sem nenhum de `project.yaml`, `MEMORY.md`, `EPOCHAL.md`, `RISKS.md`, `tasks/`, `prd/`. Ela não conta como harness global: o resolvedor devolve `repo` quando `<repo>/.harness/` existe e `none` quando não existe, e o `doctor` não emite o aviso "os dois existem". Verificável: teste de tabela do resolvedor (com `<repo>/.harness/` → `repo`; sem → `none`) e do `doctor` com `projects/<nome>/pdocs/` presente e nenhum dos seis itens.
- R3 A trava de memória (`runtime-guard.ts`) não muda: `pdocs/` não contém `MEMORY.md`, `EPOCHAL.md` nem `RISKS.md`, e escrever em `pdocs/` por subagente é permitido. Verificável: teste do mod permitindo `<home>/projects/x/pdocs/index.html` de um subagente.

Conteúdo

- R4 A skill `doc-template-html` ganha o tipo `projeto` em `references/types.md` e `assets/skeletons/<lang>/`, nos dois idiomas, com as seções fixas: Motivação; Como rodar local; Planejado e desenvolvido; Changelog; Diagrama macro; Domínios e bounded contexts; Referências. "Diagrama macro" é obrigatório: pelo menos um bloco Mermaid derivado do layout do repositório (pastas de primeiro nível e entradas), qualquer que seja a decisão 9 (leitura (b) de "como requisito sempre", §6). Os outros buracos ficam "ainda não fechado" / "not yet settled"; nunca texto inventado. A skill passa a trazer os modelos preenchidos `assets/projeto-modelo.pt-br.html` e `assets/projeto-modelo.en.html`, com as sete seções e a tag do Mermaid de R7, e o `docs-guide` (sem shell) copia a partir deles. Verificável: `stamp.sh --type projeto --lang pt-br` e `--lang en` geram as sete seções; os dois modelos têm os mesmos sete `h2 id` e a mesma tag de script que a saída do stamp; cada um tem ≥ 1 `<pre class="mermaid">` na seção "Diagrama macro"; `dh validate` verde.
- R5 "Planejado e desenvolvido" é uma tabela dos PRDs do projeto com `id`, título e `Status`, lida pelo `docs-guide` dos mesmos arquivos que o dashboard lê (`docs/prd/` do repositório e `prd/` da pasta do harness), com a data de geração; nada é calculado pelo Go para o HTML. Verificável: leitura do skeleton e de um pdocs gerado neste repositório (prova do lote).
- R6 "Changelog" copia os títulos de release do `CHANGELOG.md` do repositório quando existe (sem reescrever) e fica "ainda não fechado" quando não existe. Verificável: pdocs deste repositório lista as releases de `CHANGELOG.md`.
- R7 Diagramas são Mermaid em `<pre class="mermaid">`, com a tag de script emitida pelo `stamp.sh` apontando para uma origem única e uma versão exata (`x.y.z`, nunca só a maior), fixadas em um só lugar, o `stamp.sh`, que `references/types.md` cita sem repetir o número; os modelos de R4 trazem a mesma tag. Sem script carregado, o `<pre>` fica visível e legível (fallback de texto, sem `display:none` antes da renderização). Verificável: `grep -c "mermaid@[0-9]*\.[0-9]*\.[0-9]*/" scripts/stamp.sh` = 1 e a mesma string nos dois modelos; `types.md` sem número de versão; abrir um modelo com rede bloqueada (Playwright, rota do CDN abortada) mostra a fonte do diagrama; com rede, o diagrama.
- R8 Páginas extras ficam na mesma pasta, planas (sem subpasta), nomeadas em kebab-case e ligadas a partir do `index.html`. Verificável: `find pdocs -mindepth 2` vazio no pdocs deste repositório.

Servidor e página

- R9 Rota nova `GET /pdocs/<nome>/<arquivo>`: `<nome>` registrado em `config.yaml`, `<arquivo>` só basename (qualquer `/`, `\`, `..` ou vazio → 404), extensões `.html`, `.svg`, `.png` (demais → 404), arquivo regular até 8 MiB, `Content-Type` pela extensão, `Cache-Control: no-cache`. A rota só responde 200 ou 404, nunca escreve em `warnings`; o `Lstat` recusa symlink só no arquivo final (um `<home>` que é symlink é aceito, como hoje). `/pdocs/<nome>/` sem arquivo serve `index.html`. Verificável: testes 200 para `index.html`, `.svg`, `.png`; 404 para `../x`, `x.js`, nome não registrado, arquivo ausente, symlink no arquivo e arquivo de 8 MiB + 1 byte; `<home>` symlink → 200.
- R9b `buildState` é quem avisa: `index.html` acima de 8 MiB, ou symlink, põe `docs: false` e uma linha em `warnings` com o nome do projeto e o motivo. Verificável: teste do JSON com `index.html` de 8 MiB + 1 byte (`docs: false`, aviso) e com symlink.
- R10 Só para `/pdocs/*`, o CSP é `default-src 'self' 'unsafe-inline'; img-src 'self' data:; script-src 'self' 'unsafe-inline' <origem do Mermaid>; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; connect-src 'self'; sandbox allow-scripts allow-popups`, com a origem e a versão exata de R7. O `sandbox` sem `allow-same-origin` dá à página uma origem opaca: um script vindo do CDN não lê `/api/state` nem `/api/memory` (risco novo, não coberto pelo residual de PRD-012). `/`, `/api/*`, `/sprite/*`, `/metrics.js` mantêm o CSP de hoje. Verificável: teste que lê o cabeçalho nas duas classes de rota; teste com Playwright headless em que um `index.html` de fixture tenta `fetch('/api/state')` e recebe erro de rede ou bloqueio de CORS, nunca o JSON; Playwright headless com rede abre `/pdocs/<fixture>/index.html` pelo dashboard e vê um `svg` renderizado dentro de `pre.mermaid` e nenhuma violação de CSP no console; `grep` de que nenhum outro domínio aparece no código.
- R10b A mudança passa por `/dh:secure` (lista completa da etapa `security`) antes do merge, com R9, R10 e R12 no escopo. Verificável: veredito registrado no PR.
- R11 Só GET, loopback, `nosniff` e marcador do `guard()` valem para a rota; nada é escrito. Verificável: `POST /pdocs/x/index.html` → 405; `Host: exemplo.com` → 403.
- R12 `/api/state.projects[]` ganha o campo aditivo `docs` (`true` quando `pdocs/index.html` existe, é regular e cabe em 8 MiB, R9b); o card mostra o link "documentação" para `/pdocs/<nome>/` só quando `true`. Nenhum campo existente muda; `dh projects --json` não ganha o campo. Verificável: teste do JSON com e sem a pasta; walkthrough do `qa-verifier` abrindo o link em janela nova.

Processo e regra de CDN

- R13 O comando `/dh:document` aceita o assunto `pdocs` e despacha o `docs-guide` com `doc-template-html` para criar ou atualizar `pdocs/index.html` do projeto atual; o despacho de fim de lote do `docs-guide` (`AGENTS.md`) passa a incluir "atualizar os pdocs quando existirem", nunca criar sem pedido. Verificável: texto em `.commands/document.md`, `AGENTS.md` (bullet do `docs-guide` ao fim de cada lote) e nos agentes `coordinator` e `docs-guide`; prova no lote: pdocs deste repositório criados por `/dh:document pdocs` e atualizados no fechamento.
- R14 A regra de CDN fica única e com a mesma frase em cinco lugares, `AGENTS.md`, `.skills/doc-template-html/SKILL.md` (as duas frases de hoje, linha 58 "the only CDN exception" e linha 137 "beyond the declared Mermaid exception", reescritas para as duas exceções), `.skills/doc-template-html/references/types.md`, `docs/tutorial.html` e `docs/en/tutorial.html`: "duas exceções de CDN, Mermaid (origem e versão exata fixadas no `stamp.sh`) e Google Fonts; tudo o mais inline" (em inglês nos arquivos em inglês). Verificável: `grep -ci 'única exceção\|only CDN exception' AGENTS.md .skills/doc-template-html/SKILL.md .skills/doc-template-html/references/types.md` = 0; os cinco arquivos citam Mermaid e Google Fonts como as duas exceções.
- R15 Fora de escopo: gerar HTML pelo binário Go; diagramas animados em HTML/CSS; cópia dos pdocs para `docs/` do repositório; busca, índice global ou página com todos os projetos; acesso remoto; edição pela página; mudanças na extensão VS Code (que só ganha `docs` para ler). Verificável: nada disso no diff; `git diff go.mod` vazio.

## 4. Docs

- `.skills/doc-template-html/SKILL.md`, `references/types.md`, `assets/skeletons/pt-br/projeto.html`, `assets/skeletons/en/projeto.html`, `scripts/stamp.sh` (tipo novo, tag do Mermaid, fallback)
- `.commands/document.md` (assunto `pdocs`; exceção à regra "under `.harness/`" para o caminho vindo de `dh harness-path`), `.agents/docs-guide.md` (exceção ao limite de caminho relativo, linha 55, só para `pdocs`), `.agents/coordinator.md` (atualização de fim de lote)
- `internal/harness` e `cmd/dh` (`dh harness-path --json` com `pdocs`, R1)
- `docs/tutorial.html`, `docs/en/tutorial.html` (pdocs: onde ficam, como criar, link no dashboard, regra de CDN, o campo `pdocs` do `dh harness-path --json`)
- `README.md`, `README.en.md` (uma linha)
- `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md` (0.24.0 proposta)
- `docs/prd/PRD-012-dashboard.md` e `PRD-014-harness-global.md` (nota de remissão: rota `/pdocs`, CSP por rota, pasta `pdocs/` dentro de `projects/<nome>/` e a exceção de R2)
- `AGENTS.md` (regra de CDN de R14; layout de `<home>` com `pdocs/`; bullet do `docs-guide` ao fim de cada lote, R13; exceção de caminho fora de `.harness/` para `pdocs`)
- `<home>/projects/dev-harness/pdocs/index.html` deste repositório, gerado como prova do lote (não versionado)

## 5. Fora de escopo

- Ver R15.

## 6. Decisões do dono (decididas em 2026-10-09)

7. **Onde ficam os pdocs.** A: `<home>/projects/<nome>/pdocs/`, como no texto do dono (decidida; custo: R2 no resolvedor e no `doctor`). B: `<home>/pdocs/<nome>/`, irmã de `projects/` (sem tocar resolvedor, trava nem `doctor`; muda o texto do dono). Recomendação: A, porque o custo é uma condição testada e o layout fica o que ele pediu. Decidido: A (dono, 2026-10-09).
8. **CDN nos pdocs.** A: Mermaid e Google Fonts com origens exatas no CSP de `/pdocs/*` e fallback de texto; regra escrita como "duas exceções" (decidida). B: só Mermaid, fontes do sistema (muda a skill e todo HTML já gerado em `docs/`). C: nenhum CDN, diagramas em HTML/CSS feitos à mão (ideia da premissa; um slice de builder por diagrama, agentes desenham pior do que descrevem; fica para a vitrine do kit). Recomendação: A. Decidido: A (dono, 2026-10-09).
9. **Quando nascem.** O texto do dono diz "como requisito sempre"; há duas leituras: (a) todo projeto tem pdocs; (b) diagramas são obrigatórios em todo pdocs. A leitura (b) está honrada em R4 (Diagrama macro obrigatório) qualquer que seja a opção aqui; esta decisão é só sobre (a). A: sob pedido e atualizados ao fim de todo lote quando existem (decidida; não cumpre (a) para projeto que nunca pediu). B: criados pelo `setup` em todo projeto (cumpre (a); peso para projeto pequeno e o `setup` passa a escrever HTML). C: só sob pedido, sem atualização automática (não cumpre (a); envelhece). Recomendação: A; se o dono quis (a) ao pé da letra, B. Decidido: A (dono, 2026-10-09).

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Go renderizar os pdocs a partir dos PRDs | Segundo renderizador de HTML para manter; o `docs-guide` já escreve HTML com a skill. |
| Link `file://` em vez de rota | Bloqueado a partir de `http` e em webview. |
| Um CSP relaxado para todo o dashboard | Abriria script externo em `/` e `/api/*` sem necessidade; o relaxamento fica só onde o documento precisa (R10). |
| Vendorizar o `mermaid.min.js` em `<home>` | ~2,5 MB a distribuir ou a baixar em execução; o plano proíbe download de conteúdo em execução, e commitar no kit pesa em `dist/` (ADR-001). |

### Riscos e decisões pendentes

- Risco: sem rede, diagramas não renderizam · mitigação: R7 (fonte visível).
- Risco: `doctor` avisar "os dois existem" em projeto `repo` com pdocs · mitigação: R2.
- Risco: documentação fora do repositório não é versionada nem compartilhada com colegas · aceito pelo dono para o trabalho; cópia para `docs/` fica fora (R15).
- Risco (novo, não coberto pelo residual de PRD-012): a rota `/pdocs/*` serve HTML que carrega script de CDN; sem isolamento, esse script, ou um HTML colocado em `pdocs/` por outro processo do mesmo usuário, leria `/api/state` e `/api/memory` pela mesma origem · mitigação: `sandbox allow-scripts allow-popups` (origem opaca, R10), versão exata do Mermaid (R7), só loopback, só leitura, só nome registrado e extensões fixas (R9, R11), `/dh:secure` antes do merge (R10b).
- Decisões pendentes: nenhuma. 7, 8 e 9 decididas pelo dono em 2026-10-09 (7A, 8A, 9A).

## Resumo executado

- **Entregue:** cada projeto pode ter uma documentação HTML, o `pdocs`, em `<home>/projects/<nome>/pdocs/`, sempre na pasta global. Ela é criada por `/dh:document pdocs` com o novo tipo `projeto` (sete seções, com diagrama macro obrigatório) e atualizada no fim de cada lote quando existe. O dashboard serve essas páginas isoladas em sandbox e põe um link no card. Saiu na 0.24.0.
- **Regras:**
  - R1 Honrada: campo `pdocs` em `dh harness-path --json` nos dois modos, vazio quando o projeto não está registrado; o git continua vazio (`TestR1_*`).
  - R1b Honrada: fluxo do `/dh:document pdocs` com as exceções de caminho e a exigência de setup no modo global; prova feita neste repositório com `DH_HOME` temporário.
  - R2 Honrada: uma pasta que só tem `pdocs/` não conta como harness; o doctor não emite o aviso de "ambos" (`TestR2_*`).
  - R3 Honrada: a trava continua igual e a escrita em `pdocs/` é permitida (`claude plugin test` 220/220).
  - R4 Honrada: sete ids na ordem, `pre.mermaid` no diagrama macro e modelos pt-br/en.
  - R5, R6 e R8 Honradas: o pdocs real lista PRD-001 a PRD-017 a partir dos mesmos arquivos que o dashboard lê, o changelog vem do `CHANGELOG.md`, e a pasta é plana.
  - R7 Honrada: Mermaid 11.17.2 fixado só no `stamp.sh`, com SRI sha384; o texto do diagrama fica visível sem CDN (Playwright).
  - R9 e R9b Honradas: só basename e as extensões `.html`, `.svg`, `.png`; teto de 8 MiB; symlink e `:` recusados; 404 sem eco; `docs:false` com aviso.
  - R10 Honrada: CSP com sandbox e `script-src` restrito à URL exata do Mermaid; a página isolada não lê `/api/*` (CORS no navegador e 403 no servidor para `Origin: null`).
  - R10b Honrada: `/dh:secure` com claude e codex; os endurecimentos foram aplicados e o resultado ficou registrado no PR.
  - R11 Honrada: só GET e só loopback.
  - R12 Honrada: `docs` aditivo e link no card só quando existe pdocs.
  - R13 Honrada: cria só sob pedido e atualiza no fim do lote.
  - R14 Honrada: a frase das duas exceções de CDN está em AGENTS.md, SKILL.md, types.md e nos dois tutoriais.
  - R15 Honrada: `go.mod` sem mudança; nada fora do escopo.
- **Tickets:**
  - T-1703-01 backend: aprovado, com dúvida em voo (claude).
  - T-1703-02 backend (skill): aprovado, com SRI.
  - T-1703-03 backend (rota): aprovado, com dúvida em voo (codex).
  - T-1703-04 backend (prosa): aprovado.
  - T-1703-05 teste: verde.
  - T-1703-06 teste Playwright: 20/20.
- **Docs:** `.skills/doc-template-html/**`, `.commands/document.md`, `.agents/docs-guide.md`, `.agents/coordinator.md`, `AGENTS.md`, `CHANGELOG.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html` e as notas em PRD-012 e PRD-014.
- **Fora:**
  - HTML gerado em Go.
  - Diagramas animados.
  - Cópia para o `docs/` do repositório.
  - Índice global e acesso remoto.
  - Edição pelo dashboard.
  - Mudanças na extensão VS Code.
  - A janela de troca de diretório no Windows, já documentada.
