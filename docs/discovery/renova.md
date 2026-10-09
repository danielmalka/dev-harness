# Discovery · Projeto Renova

| Campo | Valor |
|---|---|
| Status | aprovado em 2026-10-09 |
| Dono | Daniel Malka |
| Origem | `docs/projeto-renova.md` (texto do dono, 2026-10-08) |
| Criado | 2026-10-08 |
| Autor | product-discovery (Fable), despacho UND-006 |
| Saída | este documento + `docs/prd/PRD-015-linguagens.md`, `PRD-016-metricas-dashboard.md`, `PRD-017-pdocs.md` |

Cada afirmação está marcada como fato (lido no repositório, caminho citado), hipótese (não verificada, com o que a derrubaria) ou decisão (do dono; opções e recomendação). Nada aqui autoriza implementação.

## 1. O pedido em uma frase

Preparar o kit para os próximos projetos de trabalho do dono (Kotlin, Python, TypeScript no backend), fazer o dashboard contar a história de cada projeto com números e gráficos, e dar a cada projeto uma documentação HTML viva fora do repositório. Tudo com a régua de "vai ser público": simples, didático, sem dependência escondida.

## 2. Quem usa

- **O dono, no trabalho** (usuário real hoje). Projetos em modo `global` (ADR-007: colegas não querem rastro do kit no repositório). Stacks: Kotlin com Quarkus, Micronaut ou Spring Boot; Python com Django ou FastAPI; TypeScript. Quer que o builder não "saia dos padrões mais usados" e prefira o simples.
- **Usuário público programador** (futuro). Clona, instala o plugin, roda `/dh:setup`. Precisa que os profiles novos acertem a stack sem que ele edite YAML.
- **Vibe coder / não-programador** (futuro, citado na premissa). Usa o dashboard como painel e a documentação como leitura. Não edita YAML no terminal. É o único público para quem o editor de `config.yaml` no dashboard faz diferença (seção 7).

Hoje nada do Renova existe. O dono usa o perfil `go-api` nos projetos Go e, para Kotlin/Python, cairia em `none, custom` (`profiles/README.md:50`, `.skills/project-onboarding/SKILL.md:62`).

## 3. Fase 1 · Linguagens

### O que existe (fatos)

- Profiles: `base`, `go-api`, `typescript-web`, `php` (`profiles/README.md:9-15`). `typescript-web` é só UI ("TypeScript applications with a rendered user interface", `profiles/typescript-web.yaml:3`); não há profile de TypeScript backend.
- Schema fixo: `id`, `name`, `summary`, `extends`, `stack`, `commands`, `conventions`, `limits` (`profiles/README.md:19-37`). Não há campo de framework. `extends` só funde `conventions` e `limits`.
- Comando só entra em `project.yaml` quando observado no repositório consumidor (`profiles/README.md:42`); `php.yaml:14-15` já usa o padrão "valores típicos, gravados só com evidência".
- `profiles/README.md:17` diz que Python e outras stacks ficam fora da v1 "até um consumidor real precisar". Os projetos de trabalho do dono são esse consumidor.
- O `setup` lista os profiles em tempo de execução pelos `stack` hints e não depende de lista decorada (`.skills/project-onboarding/SKILL.md:62`). Profile novo não exige editar skill nem comando.
- `dh validate` confere só `id` == nome do arquivo (`internal/kit/validate.go:438-453`) e exige `base`, `go-api`, `typescript-web` (`:343`).

### Lacunas

- Nenhuma convenção de Kotlin, Python ou TypeScript backend. O builder herda só `base` e improvisa.
- Nenhum lugar diz, por framework, qual é o layout canônico, o runner de teste e o comando de dev (Quarkus `quarkus:dev`, Spring `bootRun`, Django `manage.py`, FastAPI `uvicorn`).

### Melhoria sobre a ideia do dono

O dono pensou em "regras para Kotlin (Quarkus, Micronaut, Spring Boot) e Python (Django, FastAPI)". Lido como profiles, isso dá 5 a 7 arquivos para 3 linguagens. Proposta: **um profile por linguagem, com o framework como seção de convenções condicionada à evidência** (mesmo mecanismo do `php.yaml`). Motivos: o schema não tem framework; o matching é por `stack` hint e um repositório Kotlin é detectado antes de saber o framework; manter 7 arquivos espelhando 3 linguagens triplica a manutenção pública; a regra "comando só com evidência" já faz a distinção por framework em tempo de `setup`.

Nomes propostos, seguindo o par existente `go-api` / `typescript-web`: `kotlin.yaml`, `python.yaml`, `typescript-api.yaml`.

### Risco específico

- Não há consumidor Kotlin/Python neste repositório nem fixture. A prova de que o profile acerta só acontece no `/dh:setup` da máquina do trabalho. O PRD declara isso como limitação e pede um `setup` em modo proposta (sem escrita) como aceite final, feito pelo dono.

## 4. Fase 2 · Métricas e gráficos no dashboard, leitura da memória

### O que existe (fatos)

- Dashboard em Go stdlib, só `127.0.0.1`, rotas `/`, `/api/state`, `/sprite/*`, `/api/stop` (`internal/dashboard/server.go:54-101`). CSP `default-src 'self' 'unsafe-inline'` em toda resposta (`server.go:110`). Página `web/index.html` com 108 linhas, sem biblioteca, sem CDN (PRD-012 R5).
- O estado lê só o texto de `Status` de tickets e PRDs (`internal/dashboard/projects.go`); não lê `MEMORY.md`, `RISKS.md` nem `EPOCHAL.md`.
- Datas disponíveis hoje: `Criado / atualizado` no cabeçalho do TASK.md (`templates/pt-br/TASK.md:25`), `entregue em <data>` no `Status` do PRD, incidentes datados em RISKS.md, lotes datados em EPOCHAL.md, `started_at`/`updated_at` por sessão. Não existe histórico de transição de status.
- Custo e tokens só chegam ao snapshot pela `statusLine`, que o plugin não configura (PRD-012 §1); o mod `snapshot-writer.ts` não grava custo. Métrica de custo não tem fonte de dados hoje.
- `go.mod` sem dependência; ADR-001 fixa um binário por alvo, commitado em `dist/`, com gatilho de revisão em 50 MB; o clone já pesa ~90 MB pelos binários (`CHANGELOG.md`, 0.19.2). PRD-012 R22 e R24 recusaram SQLite explicitamente.
- `/api/state` é contrato com a extensão VS Code (PRD-014 R13/R17): mudança tem de ser aditiva.

### Melhoria sobre a ideia do dono

O dono pediu "um sqlite com métricas de uso". Três problemas: (1) um driver SQLite puro-Go seria a primeira dependência e engorda os 5 binários commitados (ADR-001; CGO quebra o cross-compile); (2) alguém precisaria escrever nele a cada transição, e hoje ninguém observa transições (os status são editados à mão em Markdown); (3) todas as métricas que o dono listou são deriváveis dos arquivos que já existem, com uma única adição de campo.

Proposta: **derivar no momento da leitura, sem armazenamento novo**. Uma linha nova no cabeçalho do ticket, `Concluído em`, fecha a lacuna da data de fim. Com isso:

| Métrica pedida | Fonte |
|---|---|
| Tarefas por PRD | já existe (`ProjectProgress`) |
| Data início e fim de cada desenvolvimento | PRD: `Criado` → `entregue em`; ticket: `Criado` → `Concluído em` |
| Tempo médio em cada task dentro de um PRD | média de (`Concluído em` − `Criado`) dos tickets do PRD, em dias |
| Tempo médio para entregar um PRD | média de (`entregue em` − `Criado`) dos PRDs entregues |
| Incidentes (RISKS) por projeto | contagem de itens com ID em `## Incidentes` do RISKS.md |
| Memória e consolidação | linhas/bytes do MEMORY.md, número de lotes e data do último lote no EPOCHAL.md |

Granularidade: **dia**. As datas são escritas à mão pelos agentes; o que o dono quer ver é "quantos dias", não "quantas horas". Ticket aberto e fechado no mesmo dia aparece como "< 1 dia". Se um dia a granularidade de hora for necessária, o passo seguinte é um `events.jsonl` append-only por projeto gravado pelo Go ao detectar mudança de status, não SQLite. Fica registrado como caminho de evolução, não como escopo.

Gráficos: sem biblioteca e sem CDN (PRD-012 R5 continua valendo). Barras e linha simples em SVG gerado por JavaScript puro na página, a partir dos números já prontos em `/api/state`. Para manter `index.html` abaixo de 500 linhas, o JavaScript dos gráficos vai num segundo arquivo embutido (`web/metrics.js`), servido por rota própria.

Leitura da memória: o dashboard passa a servir o texto de `MEMORY.md` (e, por sugestão, `RISKS.md`) do projeto por rota GET própria, identificada só pelo nome registrado (nunca por caminho), com o mesmo teto de 1 MiB da leitura atual, escapado como texto. Link `file://` não funciona a partir de página `http` nem em webview, por isso a rota.

### Riscos específicos

- Tickets antigos sem `Concluído em` ficam fora das médias; o dashboard mostra "n sem data" para não mentir.
- Médias de poucos itens oscilam; o painel mostra o `n` ao lado de toda média.
- `index.html` e `metrics.js` crescem o binário. Medir como em PRD-012 R23.

## 5. Fase 3 · Documentação HTML por projeto (pdocs)

### O que existe (fatos)

- Skill `doc-template-html`: um HTML autônomo, paleta ink-paper, Google Fonts no `<head>` (`SKILL.md:56`, `scripts/stamp.sh:209`), Mermaid em `<pre class="mermaid">` renderizado por CDN "única exceção de CDN" (`SKILL.md:58`). O `stamp.sh` não emite a tag do Mermaid (grep sem resultado); a origem do CDN não está fixada em lugar nenhum.
- `AGENTS.md` diz "Mermaid via CDN é a única exceção de CDN", mas a skill também carrega Google Fonts. Inconsistência real, a resolver no PRD-017.
- Tipos da skill: `catalogo`, `plano`, `feature`, `melhoria`, ... (`references/types.md`). Não há tipo "projeto" com motivação, como rodar, planejado vs feito, changelog e diagramas.
- `docs-guide` já é despachado ao fim de todo lote (`AGENTS.md`, `.commands/document.md`). O gancho de atualização existe; falta o artefato.
- Em modo `repo`, `<home>/projects/<nome>/` não existe (`.commands/setup.md:16`); se os dois lugares existirem, o `doctor` avisa (PRD-014 R14).
- O `guard()` do dashboard aplica o mesmo CSP a toda rota (`server.go:106-112`). Um HTML com Mermaid CDN e Google Fonts servido pelo dashboard teria script e fonte bloqueados.

### Melhoria sobre a ideia do dono

1. **Mermaid CDN ou HTML puro animado?** O dono sugeriu Mermaid CDN e, na premissa, HTML/CSS animado como "muito mais atrativo para versão pública". Leitura crítica: diagrama animado à mão custa um slice de builder por diagrama e os agentes desenham pior do que descrevem; Mermaid é texto, entra em diff e qualquer agente gera. Já o CDN falha sem rede (proxy corporativo, avião) e quebra o diagrama. Recomendação: **Mermaid com fonte visível como fallback** (o `<pre>` mostra o texto quando o script não carrega) nos pdocs; HTML animado fica para a vitrine pública do próprio kit, como item posterior.
2. **Onde ficam.** O dono quer `<home>/projects/<nome>/pdocs/`, sempre global. Custo: em modo `repo` isso cria a pasta do projeto no global só com `pdocs/` e dispara o aviso "os dois existem" do `doctor`. Alternativa `<home>/pdocs/<nome>/` (irmã de `projects/`, `sessions/`, `dashboard/`) não interfere com resolvedor, trava nem `doctor`. Decisão do dono (seção 8).
3. **Como o dashboard serve.** Rota GET `/pdocs/<nome>/<arquivo>` lendo a pasta do projeto registrado, pasta plana, só basename, extensões `.html`, `.svg`, `.png`, teto por arquivo. Para essa rota o CSP admite exatamente as duas origens (Mermaid e Google Fonts); todas as outras rotas seguem como hoje.
4. **Conteúdo mínimo** (tipo novo `projeto` na skill): motivação, como rodar local, planejado vs desenvolvido (lido dos PRDs, mesma leitura do dashboard), changelog (do `CHANGELOG.md` do repositório quando existe), diagrama macro, diagramas por domínio e bounded contexts quando houver. Buraco fica "ainda não fechado", nunca inventado.
5. **Quem escreve.** `docs-guide` com a skill, em prosa, como hoje. O Go não gera HTML; só serve. Evita um segundo renderizador.
6. **Gatilho.** O dono escreveu "como requisito sempre". Duas leituras: (a) todo projeto tem pdocs; (b) diagramas são sempre obrigatórios. Lido como (a) com diagramas quando o projeto tem domínios. O momento de criação é decisão do dono (seção 8).

### Risco específico

- Documentação fora do repositório não é versionada nem compartilhada com colegas. Aceito pelo dono para o trabalho; para o usuário público em modo `repo` fica uma sugestão de cópia para `docs/` sob pedido, fora de escopo agora.

## 6. Quebra em PRDs e ordem

| PRD | Capability | Depende de | Pode andar em paralelo com |
|---|---|---|---|
| PRD-015 · Linguagens | profiles `kotlin`, `python`, `typescript-api` | nada | PRD-016 |
| PRD-016 · Métricas, gráficos e memória no dashboard | campo `Concluído em`, métricas derivadas em `/api/state`, gráficos SVG, rota de memória | nada (toca `server.go`, `index.html`) | PRD-015 |
| PRD-017 · pdocs | tipo `projeto` na skill, pasta de pdocs, rota estática com CSP próprio, link no dashboard | PRD-016 (padrão de rota por nome de projeto e leitura de PRDs; mesmos arquivos do dashboard) | nada |

Ordem recomendada: 015 e 016 em paralelo (arquivos disjuntos), 017 depois de 016. Uma release por PRD (0.22, 0.23, 0.24), tag no merge como manda `AGENTS.md`. Sugestão, não requisito.

## 7. Depois do Renova (não entra agora)

**Editor de `config.yaml` no dashboard (ideia do dono para vibe coders).** Não entra como PRD-018 neste lote. Motivos: (1) o público do Renova, declarado na primeira linha do texto do dono, são os projetos de trabalho; o editor serve ao não-programador, público futuro; (2) hoje o `config.yaml` tem quatro chaves (`language`, `mode`, `dashboard.sprites`, `reviewers`) mais o mapa `projects:` (`internal/harness/config.go:14-20`); (3) qualquer escrita pelo navegador desmonta a postura de segurança entregue: o servidor só aceita GET (PRD-012 R3), a única rota POST recusa qualquer requisição com `Origin` (`server.go:48-52`), e um formulário do navegador sempre manda `Origin`. Abrir escrita exige CSRF, token de sessão e revisão de segurança com a lista completa; é um PRD próprio, não um apêndice. Quando vier, a forma mais barata é `dh config set <chave> <valor>` em Go (a escrita continua só no binário, PRD-014 R2) e o dashboard mostrando a configuração em modo leitura com a dica do comando; o formulário HTTP é a segunda etapa, se ainda fizer falta.

**Diagramas animados em HTML/CSS** para a vitrine pública do kit (`docs/plano-produto.html` e afins): depois dos pdocs provarem o conteúdo.

**`events.jsonl` por projeto** se a granularidade de dia não bastar (seção 4).

## 8. Decisões do dono (consolidadas)

Cada PRD repete as suas. Aqui a lista inteira, numerada. Todas decididas pelo dono em 2026-10-09, aceitando a recomendação em cada uma (1A, 2A, 3A, 4A, 5A, 6B com o risco residual aceito explicitamente, 7A, 8A, 9A, 10A).

1. **Um profile por linguagem ou um por framework?** (PRD-015) A: por linguagem, framework por evidência (3 arquivos). B: por framework (7 arquivos, `extends: kotlin`). C: só por linguagem, sem seção de framework (mais barato, não atende "padrões mais usados"). Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).
2. **TypeScript backend: profile próprio ou dentro de `typescript-web`?** (PRD-015) A: `typescript-api.yaml` novo, espelhando `go-api`. B: ampliar `typescript-web` com seção de backend (um arquivo tenta cobrir UI e API; o matching por "UI presente" fica ambíguo). Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).
3. **Os três profiles novos entram na lista obrigatória do `dh validate`?** (PRD-015) A: sim, mesma garantia dos atuais. B: não, continuam opcionais. Recomendo A (uma linha, `validate.go:343`). Decidido conforme a recomendação (dono, 2026-10-09).
4. **Armazenamento de métricas.** (PRD-016) A: derivar dos arquivos, sem armazenamento novo, granularidade de dia. B: `events.jsonl` append-only por projeto gravado pelo Go (hora, mas alguém tem de observar a transição). C: SQLite (dependência nova, engorda 5 binários, contraria ADR-001 e PRD-012 R22). Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).
5. **Tickets antigos sem `Concluído em`.** (PRD-016) A: ficam fora das médias com contador "sem data". B: preencher à mão neste repositório, como a limpeza de PRD-012 R21. C: usar `atualizado` como aproximação, marcada. Recomendo A (B pode ser feito depois pelo dono, sem mudar código). Decidido conforme a recomendação (dono, 2026-10-09).
6. **Quais arquivos de memória o dashboard mostra.** (PRD-016) A: só `MEMORY.md`. B: `MEMORY.md` e `RISKS.md`. C: os três, EPOCHAL truncado a 1 MiB. Recomendo B. Decidido conforme a recomendação (dono, 2026-10-09).
7. **Onde ficam os pdocs.** (PRD-017) A: `<home>/projects/<nome>/pdocs/` (texto do dono; exige o `doctor` ignorar pasta de projeto que só tem `pdocs/`). B: `<home>/pdocs/<nome>/` (sem tocar resolvedor nem `doctor`). Recomendo A, por ser a forma pedida e o custo ser uma condição no `doctor`. Decidido conforme a recomendação (dono, 2026-10-09).
8. **Política de CDN nos pdocs e no kit.** (PRD-017) A: Mermaid e Google Fonts, as duas origens exatas no CSP da rota `/pdocs/*`, fonte do diagrama visível como fallback; `AGENTS.md` passa a dizer "duas exceções". B: só Mermaid; fontes do sistema (muda a skill e os docs já gerados). C: nenhum CDN; diagramas em HTML/CSS à mão (ideia da premissa). Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).
9. **Quando os pdocs nascem.** (PRD-017) A: sob pedido (`/dh:document pdocs`) e, uma vez existentes, atualizados pelo `docs-guide` ao fim de todo lote. B: criados pelo `setup` para todo projeto. C: só sob pedido, sem atualização automática. Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).
10. **Editor de `config.yaml` no dashboard.** A: fica para depois do Renova (seção 7). B: entra como PRD-018 agora, começando por `dh config set` e tela de leitura. C: formulário HTTP já neste lote (revisão de segurança completa). Recomendo A. Decidido conforme a recomendação (dono, 2026-10-09).

## 9. Hipóteses

- Os projetos de trabalho usam Gradle com Kotlin DSL ou Maven, e Python com `pyproject.toml`. Derrubada por um repositório com `build.gradle` Groovy ou só `requirements.txt`; os profiles tratam isso como evidência, não como premissa.
- Datas de `Criado` nos tickets estão preenchidas de fato nos projetos do dono. Derrubada por tickets com `AAAA-MM-DD` literal; o painel conta "sem data".
- A extensão VS Code tolera campos novos em `/api/state`. Derrubada por parser estrito lá; o PRD-016 exige mudança só aditiva.
- O dono aceita granularidade de dia. Derrubada se ele quiser "tempo de sessão por task"; nesse caso a fonte é `events.jsonl` ou os snapshots, não as datas dos cabeçalhos.

## 10. Limitações desta discovery

- Nenhum repositório Kotlin ou Python real foi lido; as convenções dos profiles vêm do conhecimento geral das ferramentas e precisam do `setup` em proposta no trabalho.
- Não medi o tamanho do binário com os arquivos novos embutidos; o PRD-016 pede a medida.
- Não rodei comando algum além de leitura (`cat`, `grep`, `ls`, `wc`).
- Discovery e PRDs em pt-br, como os PRDs existentes em `docs/prd/` (sem espelho em `docs/en/prd/`); se o Coordenador quiser espelho, é despacho à parte.

## 11. Evidências

`docs/projeto-renova.md`; `profiles/README.md`; `profiles/go-api.yaml`; `profiles/typescript-web.yaml`; `profiles/php.yaml`; `profiles/base.yaml`; `internal/kit/validate.go`; `internal/dashboard/server.go`; `internal/dashboard/projects.go`; `internal/dashboard/sessions.go`; `internal/dashboard/web/index.html`; `internal/snapshot/snapshot.go`; `internal/harness/config.go`; `adapters/claude-code/plugin/hooks/snapshot-writer.ts`; `.skills/doc-template-html/SKILL.md`; `.skills/doc-template-html/scripts/stamp.sh`; `.skills/doc-template-html/references/types.md`; `.skills/project-onboarding/SKILL.md`; `.commands/setup.md`; `.commands/document.md`; `templates/pt-br/TASK.md`; `templates/pt-br/MEMORY.md`; `templates/pt-br/RISKS.md`; `templates/pt-br/EPOCHAL.md`; `docs/prd/PRD-012-dashboard.md`; `docs/prd/PRD-014-harness-global.md`; `docs/adr/ADR-001-runtime-go.md`; `docs/adr/ADR-007-harness-global.md`; `CHANGELOG.md`; `AGENTS.md`.
