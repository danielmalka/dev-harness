# PRD-024 · Deixar o repositório público pronto para pessoas de fora

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-16 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-16 / 2026-10-16 |
| Tickets | T-2401 (planejado) |
| Origem | Conversa com o dono, 2026-10-16 (escopo e decisões na conversa) |
| Depende de | nada |
| Habilita | nada |

## 1. Problema

- O repositório do GitHub é público. Uma pessoa de fora que o abre não consegue saber o que é o kit, como removê-lo nem o que ele grava na máquina. Fato: auditoria UND-006 do `repo-scout`.
- `README.md` e `README.en.md` têm estas lacunas:
  - A linha 20 é um parágrafo de notas de release de cerca de 2.500 caracteres. Ela guarda a versão, a extensão do VS Code, o dashboard, contagens de teste (`"226 pass"`) e contagens do kit. O `CHANGELOG.md` já guarda esse conteúdo.
  - Nenhuma seção diz o que o kit é, o que ele exige, como começar, como atualizar, como desinstalar e remover dados, o que o kit grava e onde, como desligar as coisas, nem o que fazer quando algo falha.
  - A informação de controle existe só em `docs/tutorial.html` e `docs/en/tutorial.html`, em lugares espalhados (desativar o plugin, poda de mod ou de snapshot, parada do dashboard).
  - A linha 7 dos dois READMEs liga `docs/plano-produto.html` como a página do produto. Essa página está em pt-br e é interna. O README em inglês não deve abrir com ela.
  - O bloco de manutenção é só para quem contribui, mas nada diz isso.
  - Nenhum dos READMEs nomeia os dois projetos companheiros: `danielmalka/malka-claude-mods` (mods visuais que saíram deste kit na 0.29.0) e `akitaonrails/ai-memory` (memória de longo prazo e handoff entre CLIs de código). O kit não precisa de nenhum dos dois.
- Dados pessoais na árvore:
  - 28 dos 52 arquivos `evals/baselines/*.json` rastreados contêm o caminho absoluto da home do dono (o caminho do checkout do repositório). Nenhum outro arquivo rastreado tem o caminho da home do dono. Nenhum segredo ou token foi encontrado.
  - `"danlemos"`, o nome de outro projeto do dono, aparece em `docs/roadmap.html` e `docs/en/roadmap.html` (perto da linha 1631). Não significa nada para uma pessoa de fora. Também aparece em PRDs e ADRs antigos, que são histórico.
  - Os dois roadmaps citam `CLAUDE.md`, que só importa o `AGENTS.md`. Linhas da base: pt 1543, 1632, 1698, 1780; en 1541, 1632, 1699, 1783.
- Os dois roadmaps citam `RISKS.md` nas Relações do RISK-002 (pt 921 e 1591; en 922 e 1589). Esse registro vive no registro do harness do mantenedor, não neste repositório. As linhas 760 e 839 dos dois arquivos também citam `RISKS.md`. A prosa delas está correta (registros de projetos consumidores) e fica.
- Os números de linha são da base. O construtor acha cada linha de novo pelo texto.
- Arquivos do repositório público: existem só `LICENSE` e `CHANGELOG.md`. Não há `SECURITY.md` nem `CONTRIBUTING.md`.
- O roadmap não tem card para três itens adiados de `docs/discovery/renova.md` (seção 7): um formulário para o `config.yaml` no dashboard, diagramas animados para uma vitrine pública e um `events.jsonl` por projeto.
- O roadmap não tem card para um pedido do dono. O Coordenador faz trabalho operacional demais por conta própria. Ele lê arquivos, roda build, validate e grep, edita registros, monta prompts de revisão por CLI e faz checagens por amostragem pelo shell. Ele deveria despachar subagentes. Fato: observação do dono, 2026-10-16.

## 2. Solução

Um release, 0.30.1, em uma PR. Depois da entrega, uma pessoa de fora consegue ler o que o kit é, instalá-lo, removê-lo, achar o que ele gravou e reportar uma vulnerabilidade. A árvore não guarda nenhum caminho da máquina do dono. Não adicione comportamento ao kit. Não edite nenhum arquivo em `.agents/`, `.skills/` ou `.commands/`.

README (pt-br em `README.md`, inglês em `README.en.md`). Os dois arquivos têm as mesmas 12 seções `## ` na mesma ordem. Os textos dos títulos são fixos:

| # | Título em inglês | Título em pt-br | Conteúdo |
|---|---|---|---|
| 1 | `What it is` | `O que é` | No máximo 3 linhas. |
| 2 | `Requirements` | `Requisitos` | Claude Code; Git. |
| 3 | `Install` | `Instalar` | A forma da URL do marketplace que existe hoje. |
| 4 | `Quick start` | `Início rápido` | Os nomes reais de comando de `docs/tutorial.html` e de `.commands/` (por exemplo `/dh:doctor`, `/dh:setup`). Termina com no máximo 2 linhas sobre idioma. |
| 5 | `Update` | `Atualizar` | Como atualizar o plugin. |
| 6 | `Uninstall and remove data` | `Desinstalar e remover dados` | `/plugin uninstall`; o conteúdo de `~/.harness`; a entrada de `additionalDirectories`; a pasta `.harness/` de um repositório. |
| 7 | `What the kit writes and where` | `O que o kit grava e onde` | Uma tabela com 7 linhas: `.harness/` do repositório; `~/.harness/config.yaml`; `projects/<name>/`; `sessions/`; `dashboard/`; `pdocs/`; o dashboard em `127.0.0.1:4747`. |
| 8 | `Turn things off` | `Desligar` | `dh dashboard --stop`; desabilitar o plugin. Uma frase diz que não existe chave por hook para o mod nem para a armadilha de commit (pt: `"por hook"`). |
| 9 | `If something fails` | `Se algo falhar` | `dh doctor` ou `/dh:doctor`; porta ocupada; plugin não encontrado; link para o GitHub Issues. |
| 10 | `Use together with` | `Usar junto com` | Os dois projetos companheiros, com links e uma linha cada sobre o que acrescentam. O texto diz que eles são opcionais (pt: `"opcional"`). |
| 11 | `Docs map` | `Mapa da documentação` | Início rápido, tutorial, roadmap, `CHANGELOG.md`. `docs/prd`, `docs/adr` e `docs/discovery` são histórico de design em pt-br. `evals/` é histórico de medição e não é necessário para usar o kit. `docs/plano-produto.html` é rotulado como pt-br e interno, e o README em inglês não o liga primeiro. |
| 12 | `For contributors` | `Para quem contribui` | O bloco de manutenção (`go test ./...` e os comandos de build), a tabela de fontes e pacote e as instruções de carga pelo clone. |

Para onde vão os blocos atuais do README:

- Bloco de idioma: seus fatos encolhem para no máximo 2 linhas no fim de Início rápido.
- `"Load from the clone"` e `"Sources and package"`: movem para a seção 12.
- A nota de migração da 0.8.0: removida. O `CHANGELOG.md` a mantém.
- Parágrafo do modo global: uma linha antes da tabela da seção 7. O tutorial guarda os detalhes.
- Parágrafo de notas de release (linha 20): removido. O `CHANGELOG.md` é a fonte dele.

O README não guarda contagem de teste. Uma contagem do kit (de agentes, comandos e skills) aparece em no máximo um lugar: uma linha que tenha a palavra `"atual"`, `"esperado"`, `"current"` ou `"expected"` e pelo menos dois substantivos contados. `checkDocumentedCounts` em `internal/kit/doc_counts.go` checa contra o disco só uma linha assim. Uma contagem em qualquer outra linha não é checada, então o README não declara nenhuma.

Roadmap (`docs/roadmap.html` e `docs/en/roadmap.html`, espelhados à mão):

- Quatro cards na seção Planejado (`<h2 id="planejado">`). Nenhum `<h3>` começa com número de versão. Texto exato dos `<h3>` (en / pt-br):
  - (a) `"Edit config.yaml from the dashboard"` / `"Editar o config.yaml pelo dashboard"`. Precisa de um token CSRF ou de sessão e de uma revisão de segurança completa. O PRD-018 entregou o lado de escrita do `dh config`.
  - (b) `"Animated diagrams for the public showcase"` / `"Diagramas animados para a vitrine pública"`. Começa depois que os `pdocs` provarem o conteúdo.
  - (c) `"Per-project events.jsonl for finer metrics"` / `"events.jsonl por projeto para métricas mais finas"`. É a opção B do PRD-016. A granularidade hoje é de um dia.
  - (d) `"Coordinator orchestrates, does not execute"` / `"O Coordenador orquestra, não executa"`. O Coordenador despacha subagentes para execução e verificação. Execução cobre leitura, build, validate, grep, edição de registros, prompts de revisão por CLI e checagens por amostragem. O Coordenador só julga se o trabalho delegado foi feito com perícia ou precisa de retrabalho. Ele continua dono das escritas em `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`. Status: planejado, ainda sem PRD. Alavancas para avaliar depois: as regras e a redação da checagem por amostragem em `.agents/coordinator.md`, e talvez uma guarda de mod. Relações: `.agents/coordinator.md`.
  - Cada card tem Descrição, Critérios de aceite, Status e Relações, como os cards existentes.
- `"danlemos"` vira texto genérico só nos dois roadmaps. Os documentos de histórico ficam sem mudança.
- Referências mortas: cada ocorrência de `CLAUDE.md` (linhas da seção 1) vira `AGENTS.md`. Nas Relações do RISK-002, o texto diz que o RISK-002 vive no registro do harness do mantenedor, não neste repositório. As linhas 760 e 839 com `RISKS.md` ficam sem mudança.
- Um card Entregue para a 0.30.1 nos dois arquivos.

Outros arquivos:

- `evals/baselines/*.json`: o prefixo do caminho do checkout do dono vira `<repo>`. O histórico do git não é reescrito.
- `SECURITY.md` em inglês: reportar uma vulnerabilidade pelo relato privado de vulnerabilidade do GitHub; uma política curta.
- `CONTRIBUTING.md` em inglês: fontes canônicas (`.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`); `dist/` vem de `go run ./cmd/dh build` e nunca é editado à mão; checagens `go test ./...`, `go run ./cmd/dh validate .` e `go run ./cmd/dh build`; Conventional Commits em inglês, sem menção a IA; `AGENTS.md` é o arquivo canônico de regras e está em pt-br.
- Release 0.30.1: entrada no `CHANGELOG.md`, o card Entregue, `version` e `ref` em `.claude-plugin/marketplace.json`, `defaultVersion` em `internal/build/build.go`.

## 3. Regras

Cada checagem abaixo falha na árvore no momento da escrita, exceto R7, R11, R17, R18, R19 e R21. Essas seis protegem contra uma regressão e devem valer no fim. O autor do ticket checa de novo na base antes da primeira edição.

- R1 `README.md` e `README.en.md` têm as 12 seções da seção 2 na mesma ordem, e cada arquivo tem no máximo 110 linhas. Verificável: `grep '^## ' README.en.md` imprime os 12 títulos em inglês da seção 2 nessa ordem, `grep '^## ' README.md` imprime os 12 títulos em pt-br nessa ordem; `wc -l README.md README.en.md` mostra 110 ou menos para cada arquivo.
- R2 Nenhum README guarda o parágrafo de notas de release nem contagem de teste. Verificável: `grep -n -E 'Versão atual|Current version|226|Testes:|Tests:' README.md README.en.md` não tem resultado, e nenhuma linha dos dois arquivos passa de 600 caracteres (`awk 'length>600' README.md README.en.md` não imprime nada).
- R3 Os dois READMEs nomeiam os dois projetos companheiros com links e dizem que são opcionais. Verificável: `grep -c 'github.com/danielmalka/malka-claude-mods' README.md README.en.md` e `grep -c 'github.com/akitaonrails/ai-memory' README.md README.en.md` imprimem 1 ou mais para cada arquivo, e a seção `"Use together with"` / `"Usar junto com"` guarda `"optional"` / `"opcional"`.
- R4 Os dois READMEs dizem o que o kit grava e onde, como desinstalar e remover dados, e como desligar as coisas. Verificável, por seção (leia o texto entre o título da seção e o próximo `## `):
  - Seção de desinstalar: tem `/plugin uninstall`, `~/.harness`, `.harness/` e `additionalDirectories`.
  - Seção do que grava: tem 7 linhas de corpo (depois do cabeçalho e do separador) que começam com `|`, nomeando `.harness/`, `config.yaml`, `projects/`, `sessions/`, `dashboard/`, `pdocs/` e `4747`.
  - Seção de desligar: tem `dh dashboard --stop`, um passo do plugin (a forma de desabilitar que o QA verificou, ou `/plugin uninstall` quando `claude plugin --help` não mostra forma de desabilitar) e a frase com `per-hook` (pt: `por hook`).
- R5 O Mapa da documentação rotula `docs/plano-produto.html` como pt-br e interno, e `README.en.md` não o liga antes do início rápido. Verificável: `grep -n 'plano-produto' README.en.md` mostra um resultado só dentro da seção Mapa da documentação, com o rótulo; nenhum resultado na linha 7.
- R6 O bloco de manutenção fica sob um título ou uma linha de abertura que diz que é só para quem contribui. Verificável: o título da seção 12 contém `"contribut"` (pt: `"contribu"`) e guarda o bloco de comando `go test ./...`.
- R7 Nenhum README guarda uma contagem do kit não checada ou errada. Verificável: `go run ./cmd/dh validate .` passa, e `grep -n -i -E '[0-9]+ (agents|agentes|commands|comandos|skills|profiles|perfis|templates)' README.md README.en.md` não tem resultado, ou seus resultados estão em uma linha por arquivo que tem `"atual"`, `"esperado"`, `"current"` ou `"expected"` e pelo menos dois substantivos contados.
- R8 Nenhum arquivo rastreado guarda o caminho da máquina do dono. Verificável: `git grep -l '/ho[m]e/ma[l]ka'` sobre a árvore inteira não imprime nada (as classes de colchetes casam o caminho, e o texto do PRD não guarda o literal). Cada um dos 28 arquivos JSON alterados ainda é lido pelo parser: `for f in evals/baselines/*.json; do python3 -m json.tool "$f" >/dev/null || echo bad $f; done` não imprime nada. A substituição é exatamente `<repo>`: `git grep -l '<repo>' -- evals/baselines` lista 28 arquivos.
- R9 `SECURITY.md` existe em inglês e nomeia o relato privado de vulnerabilidade do GitHub. Verificável: `grep -n -i 'private vulnerability reporting' SECURITY.md` tem um resultado.
- R10 `CONTRIBUTING.md` existe em inglês e tem os 5 itens da seção 2. Verificável: `grep -n -E '\.agents/|go run \./cmd/dh build|go run \./cmd/dh validate|go test \./\.\.\.|Conventional Commits|AGENTS\.md' CONTRIBUTING.md` encontra cada termo, e `grep -n -i -E 'co-authored-by|generated with' CONTRIBUTING.md` não tem resultado, exceto a linha de regra que proíbe menção a IA. Os nomes Claude Code e `.claude-plugin` são permitidos.
- R11 Nenhum `CODE_OF_CONDUCT.md`, nenhum `.github/ISSUE_TEMPLATE` e nenhum template de pull request é adicionado. Verificável: `git diff --name-status main...HEAD` não lista nenhum deles.
- R12 Os dois roadmaps guardam os quatro cards Planejados, espelhados. Verificável: em cada arquivo, as linhas entre `id="planejado"` e o próximo `<h2` contêm os quatro textos exatos de `<h3>` da seção 2 (en em `docs/en/roadmap.html`, pt-br em `docs/roadmap.html`), as Relações do card (d) nomeiam `.agents/coordinator.md`, e os dois arquivos têm o mesmo número de `<div class="card">` nessa seção. Nenhum `<h3>` dos quatro cards começa com número de versão.
- R13 `"danlemos"` some dos dois roadmaps e fica nos documentos de histórico. Verificável: `grep -c danlemos docs/roadmap.html docs/en/roadmap.html` imprime 0 para os dois; `git diff --name-only main...HEAD -- docs/adr docs/prd docs/discovery` lista só este PRD e seu arquivo de revisão.
- R14 Os dois roadmaps não têm referência morta. Verificável: `grep -n 'CLAUDE\.md' docs/roadmap.html docs/en/roadmap.html` não imprime nada; `grep -n 'RISKS\.md' docs/roadmap.html docs/en/roadmap.html` imprime exatamente 3 linhas por arquivo: as linhas da base 760 e 839 (achadas de novo pelo texto) e a linha do card (d). Nenhuma linha de Relações do RISK-002 cita `RISKS.md`. O card Entregue da 0.30.1 não escreve `CLAUDE.md` nem `RISKS.md`.
- R15 Versão 0.30.1 em três lugares: `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da origem (`v0.30.1`). Verificável: `grep` dos três valores mostra `0.30.1` ou `v0.30.1`.
- R16 `CHANGELOG.md` tem uma entrada para a 0.30.1, e os dois roadmaps têm um card Entregue para a 0.30.1 que concorda com ela. Verificável: `grep -n '0.30.1' CHANGELOG.md docs/roadmap.html docs/en/roadmap.html` tem um resultado em cada arquivo, e `dh validate` não reporta deriva do roadmap.
- R17 A mudança não edita nenhum arquivo em `.agents/`, `.skills/` ou `.commands/`, e não adiciona agente, skill nem comando. As contagens continuam em 13 agentes e 19 + 1 comandos. Verificável: `git diff --name-only main...HEAD -- .agents .skills .commands` não imprime nada.
- R18 As edições à mão se limitam aos arquivos da seção 4, a este PRD com seu arquivo de revisão e aos registros do ticket. Nenhum arquivo em `dist/` é editado à mão. Verificável: `git diff --name-only main...HEAD -- dist` lista só arquivos que `go run ./cmd/dh build` regenera (`plugin.json`, `harness-manifest.json`, `dist/*/bin/**`).
- R19 `go run ./cmd/dh validate .` passa com 0 erros, e `go test ./...` está verde. O build roda duas vezes, e a segunda execução não muda nada em `dist/` além da linha `generated_at` de `harness-manifest.json`. Verificável: rode os comandos e compare com `git diff`.
- R20 Prova de clone limpo, porque a PR toca `internal/build/build.go` e `dist/*/bin`: rode `git clone . <tmp>`, rode `go run ./cmd/dh build` lá e compare o `sha256sum` dos binários com a árvore. Verificável: 5 de 5 somas iguais.
- R21 A mudança não reescreve o histórico do git. Verificável: nenhum push com `--force`, e `git merge-base --is-ancestor main HEAD` tem sucesso.
- R22 Todo comando nas seções do README Atualizar, Desinstalar e Desligar conta como verificado na forma que o README imprime (barra ou CLI). Um comando é verificado quando os docs deste repositório o mostram, ou quando o QA o rodou nesta PR. As formas de CLI que o Coordenador já rodou são `claude plugin marketplace update dev-harness` e `claude plugin update dh@dev-harness`. O QA checa `claude plugin --help` em busca de uma forma de desabilitar. Se a ajuda não mostrar nenhuma, a seção Desligar usa `/plugin uninstall` como passo do plugin, e o R4 aceita isso. Verificável: o relatório do QA lista cada comando do README com sua fonte (caminho do doc ou a saída da execução).

## 4. Docs

Atualize:

- `README.md`, `README.en.md` (reescritos como na seção 2).
- `docs/roadmap.html`, `docs/en/roadmap.html` (quatro cards Planejados, card Entregue da 0.30.1, `"danlemos"`, referências mortas).
- `CHANGELOG.md` (entrada 0.30.1).
- `.claude-plugin/marketplace.json` (`version` e `ref`).
- `internal/build/build.go` (`defaultVersion`).
- `evals/baselines/*.json` (28 arquivos, prefixo do caminho).
- `SECURITY.md`, `CONTRIBUTING.md` (novos).

Não afetados (o autor do ticket checa com `grep`):

- `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/inicio-rapido.html`, `docs/en/quick-start.html`. O README liga para eles. Não os muda.
- `docs/plano-produto.html`. Fica como está. Só o README o rotula.
- `docs/adr/`, `docs/prd/` (outros PRDs) e `docs/discovery/`. Histórico. Eles mantêm `"danlemos"` de propósito.

## 5. Fora do escopo

- Reescrita do histórico do git para remover o caminho dos commits antigos. Aceito pelo dono.
- Tradução de `AGENTS.md`, PRDs ou ADRs.
- `CODE_OF_CONDUCT.md`, templates de issue, templates de pull request.
- Mover `evals/` ou os arquivos `*.review.md`.
- Qualquer mudança de comportamento do kit, e qualquer edição em `.agents/`, `.skills/` ou `.commands/`.
- Uma nova chave por hook para o mod ou para a armadilha de commit.
- Implementação dos quatro itens Planejados, inclusive qualquer edição de `.agents/coordinator.md`. Só os cards deles no roadmap estão no escopo.
- Descrição e tópicos do GitHub. O Coordenador os definiu com `gh repo edit`, fora da PR.
- Mudança do nome do titular em `LICENSE`.

## 6. Decisões do dono

- Dono, 2026-10-16: release 0.30.1 para tudo isto, em uma PR. Motivo: docs nunca saem numa PR só de docs, e todo merge recebe a tag `v<version>` (`AGENTS.md`). O release muda a entrada do `CHANGELOG.md`, o card Entregue nos dois roadmaps, `version` e `ref` do marketplace, e `defaultVersion`. A prova de clone limpo se aplica (R20).
- Dono, 2026-10-16: o titular em `LICENSE` continua `"Daniel Lemos"`. Sem mudança.
- Dono, 2026-10-16: em `evals/baselines`, o prefixo do caminho do checkout do dono vira `<repo>`. O histórico do git não é reescrito. O dono aceita que os commits antigos mantenham o caminho.
- Dono, 2026-10-16: criar `SECURITY.md` (relato privado de vulnerabilidade do GitHub, política curta) e um `CONTRIBUTING.md` curto, ambos em inglês. Nenhum `CODE_OF_CONDUCT.md`, nenhum template de issue ou de pull request.
- Dono, 2026-10-16: adicionar um quarto card Planejado, `"Coordinator orchestrates, does not execute"`. O dono observou que o Coordenador faz trabalho operacional demais. Desejado: ele orquestra e só julga se o trabalho delegado foi feito com perícia ou precisa de retrabalho; execução e verificação vão para subagentes; ele mantém as escritas em `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`. Status planejado, ainda sem PRD.
- Dono, 2026-10-16: descrição e tópicos do GitHub são definidos fora da PR (feito).
- Dono, 2026-10-16: teto de validação atingido com P9 em aberto; o dono aceitou a correção do validador para P9 e P10 sem outra rodada.
- Dono, 2026-10-16: o build achou o R8 em conflito com checkPrivatePaths; o dono aceitou a emenda de R8 e R22 e as correções P11-P14 do validador sem outra rodada.
- Decisões pendentes: nenhuma.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Reescrever o histórico do git para remover o caminho | O dono aceita o caminho antigo. Uma reescrita quebra clones e tags, e o caminho não é segredo. |
| Manter o parágrafo de notas de release no README | O `CHANGELOG.md` é a fonte dele. O parágrafo fica velho a cada release. |
| PR só de docs para o README | `AGENTS.md`: a documentação entra dentro da PR da funcionalidade. |
| Adicionar agora uma chave para cada hook | É uma mudança de comportamento. Precisa de um PRD próprio. |
| Mover `evals/` e os `*.review.md` para fora da árvore | É uma mudança grande sem ganho para quem lê. O README os rotula no lugar. |

### Riscos e decisões pendentes

- Risco: os READMEs ou os roadmaps pt e en divergem · mitigação: R1 e R12 checam a paridade; `dh validate` checa os cards Entregues.
- Risco: o README declara uma contagem que fica velha · mitigação: R2 e R7 proíbem contagens fora das linhas ancoradas que `checkDocumentedCounts` checa.
- Risco: a troca do caminho quebra um JSON de baseline ou um teste que o lê · mitigação: R8 faz o parse de cada arquivo, e R19 roda `go test ./...`.
- Risco: o README volta a crescer até virar changelog · mitigação: R1 limita as linhas e R2 limita o tamanho de linha.
- Risco: a prova de clone limpo difere da árvore · mitigação: R20 bloqueia o commit até as somas serem iguais.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** o repositório público tem um README para quem chega, em pt-br e em inglês (12 seções espelhadas com instalar, atualizar, desinstalar, o que o kit grava, como desligar, solução de problemas e os dois companheiros opcionais), `SECURITY.md` e `CONTRIBUTING.md`, quatro cards Planejados no roadmap, nenhum caminho do dono nos arquivos rastreados e o release 0.30.1.
- **Regras:**
  - R1-R7 Cumpridas: matriz de QA (`qa-verifier`): 12 títulos em ordem nos dois READMEs (106 linhas cada um depois da passada de docs), nenhum parágrafo de release nem contagem de teste, os dois links dos companheiros marcados como opcionais, seções de desinstalar, do que grava e de desligar delimitadas como exigido, plano-produto só no Mapa da documentação, bloco só para quem contribui, nenhuma contagem sem âncora.
  - R8 Cumprida: o grep com classes de colchetes sobre a árvore inteira não imprime nada; 28 baselines guardam `<repo>` e passam no parse.
  - R9-R11 Cumpridas: `SECURITY.md` nomeia o relato privado de vulnerabilidade; `CONTRIBUTING.md` cobre os cinco itens sem frase de atribuição a IA; nenhum código de conduta nem templates.
  - R12-R14 Cumpridas: quatro cards Planejados com os títulos exatos nos dois roadmaps; nenhum `danlemos`; nenhum `CLAUDE.md`; exatamente três linhas `RISKS.md` por roadmap.
  - R15-R16 Cumpridas: 0.30.1 em `build.go`, `version` e `ref` do `marketplace.json`; entrada do CHANGELOG e cards Entregues sem deriva.
  - R17-R19 Cumpridas: nenhuma edição em `.agents/`, `.skills/`, `.commands/`; `dist/` só pelo build; `dh validate` 0 erros, `go test ./...` verde, o segundo build muda só `generated_at`.
  - R20 Cumprida: prova de clone limpo, 5 de 5 binários com sha256 igual.
  - R21 Cumprida: nenhuma reescrita de histórico.
  - R22 Cumprida: cada comando do README em Atualizar, Desinstalar e Desligar tem uma fonte (docs do repositório ou `claude plugin ... --help`); nenhum foi rodado contra o plugin instalado.
- **Tickets:** T-2401 (backend, harness-maintainer) `approved`. QA por `qa-verifier`: verificado (notas menores resolvidas pela passada de docs ou aceitas). Revisão de código por `cli:codex/gpt-6-luna`: `approve`, 0 achados. Aderência por `product-discovery`: todo item adere.
- **Docs:** `README.md`, `README.en.md`, `SECURITY.md`, `CONTRIBUTING.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `.claude-plugin/marketplace.json`, `internal/build/build.go`.
- **Fora do escopo:** reescrita de histórico, tradução de `AGENTS.md`, PRDs e ADRs, código de conduta e templates, mover evals ou arquivos de revisão, mudanças de comportamento do kit, uma chave por hook e a implementação dos quatro itens Planejados.
