# PRD-021 · Kit enxuto: caminho quente mais rápido, cache de parse do dashboard, template de PRD só em inglês e uma regra de limpeza

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-13 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-13 / 2026-10-13 |
| Tickets | T-2101, T-2102 |
| Origem | Chat com o dono, 2026-10-13 (o dono delegou: "siga suas recomendações desde que siga as minhas premissas") |
| Depende de | PRD-020 (PRD em inglês), PRD-019 (padrão de versão e `ref`) |
| Habilita | nada |

## 1. Problema

Fatos medidos em 2026-10-13 com o pprof do Go na árvore 0.27.0, linux/amd64:

- O binário `dh` tem 7.65 MB (flags `-s -w -trimpath`). Só `internal/dashboard` importa `net/http`. O binário não tem nenhum módulo externo.
- `dh snapshot event` roda em todo evento de hook. Ele leva cerca de 3.96 ms por chamada. Um processo vazio leva cerca de 1.07 ms. 200 chamadas levam 0.793 s.
- O trabalho de `dh snapshot event` é 62 µs e 46 alocações. A inicialização dos pacotes soma cerca de 0.55 ms.
- O custo da inicialização é: `internal/kit` 0.24 ms (854 alocações), `internal/dashboard` 0.20 ms (599 alocações), `internal/harness` 0.11 ms (77 alocações).
- Há 34 chamadas `regexp.MustCompile` em nível de pacote nos três pacotes: kit 20, dashboard 12, harness 2. `internal/doctor` tem 1 e `internal/snapshot` não tem nenhuma. `internal/kit` também compila 3 regexps dentro de funções (`validate.go` linhas 445, 459, 482).
- O `buildState` do dashboard nos dados reais do dono leva 55.8 ms e 12.9 MB por chamada, com 11,573 alocações.
- 85% dessa CPU está em `projectProgress`. 56% é o `FindStringSubmatch` de regexp em `internal/dashboard/projects.go`. 18% é `snapshot.ReadCapped`.
- Todo refresh (TTL do cache de 1 s) relê e reparseia todo arquivo de ticket e de PRD, mesmo quando o arquivo não mudou.
- O `dh validate` leva cerca de 0.25 s e aloca 81 MB. 70 MB são `io.ReadAll` e `readText`, porque cada checagem relê os mesmos arquivos. Este item é um candidato a continuação. Está fora de escopo.
- `templates/pt-br/PRD.md` é uma cópia idêntica byte a byte do template de PRD em inglês (PRD-020). Ele acrescenta um arquivo, uma regra de espelho e uma afirmação de contagem, sem nenhum leitor.
- A checagem de espelho `checkTemplateRefs` em `internal/kit/validate.go` (linhas cerca de 501-521) é o único bloqueio para apagar essa cópia.
- Nenhuma regra diz que uma mudança deve remover os restos que ela torna obsoletos.

## 2. Solução

- O dashboard continua no binário único `dh`. O binário inteiro é analisado com o pprof para melhorá-lo e otimizá-lo, nunca para remover uma parte.
- As regexps em nível de pacote em `internal/kit`, `internal/dashboard` e `internal/harness` compilam no primeiro uso. `dh snapshot event` então não paga a inicialização delas.
- O dashboard mantém um cache de parse com chave por caminho, tamanho e data de modificação do arquivo. Os arquivos de ticket e de PRD que não mudaram não são relidos a cada refresh.
- `templates/pt-br/PRD.md` é apagado. `checkTemplateRefs` tem uma lista de templates só em inglês com a entrada `PRD.md`. Toda afirmação que depende do espelho é corrigida.
- `.agents/coordinator.md` "Seleção de modelo" ganha uma régua curta de modelo e esforço.
- `AGENTS.md` ganha a premissa de limpeza: toda mudança remove os restos que ela torna obsoletos. Lixo remanescente bloqueia a entrega.
- `AGENTS.md` declara que os 5 binários são gerados e commitados a cada release. Declara que o pprof otimiza o binário e não o poda.
- Versão 0.28.0.

## 3. Regras

- R1 `templates/pt-br/PRD.md` é apagado. `checkTemplateRefs` tem uma lista de templates só em inglês com a entrada única `PRD.md`.
  - Um template dessa lista só pode existir em `templates/en/`. O mesmo nome em outra pasta de idioma é um erro.
  - Qualquer outro template ausente em um idioma continua sendo um erro.

  Verificável: `ls templates/pt-br` não mostra `PRD.md`; testes em Go cobrem três casos (PRD só em `en` passa; PRD também em `pt-br` falha; TASK só em `en` falha); `go test ./internal/kit` verde.
- R2 Toda afirmação que depende do espelho está correta.
  - `harness-manifest.json` lista 13 templates depois de `go run ./cmd/dh build`. `checkDocumentedCounts` lê a contagem do disco. `internal/kit/doc_counts.go` não ganha nenhuma mudança de contagem. Só as suas duas regexps compilam no primeiro uso (R3).
  - `docs/roadmap.html` e `docs/en/roadmap.html` dizem 13 templates (linha cerca de 685) e "7 templates em `en/`, 6 em `pt-br/`" (linha cerca de 678).
  - O primeiro item do `AGENTS.md`, os tutoriais de templates (html e md, pt e en), as linhas do README sobre templates e `docs/dev-harness-mapa.excalidraw` estão corretos.
  - Nenhum resto permanece.

  Verificável: `grep -rniE "espelhad|mirrored|cópia idêntica|identical copy" README.md README.en.md AGENTS.md docs --include=*.html --include=*.md --include=*.excalidraw` só tem hits que são histórico (`docs/adr`, `docs/prd`, `CHANGELOG.md`, cards de release do roadmap) ou que nomeiam só outros templates; `dh validate` verde; `go run ./cmd/dh build` não deixa nenhum `pt-br/PRD.md` em `dist/`.
- R3 As `regexp.MustCompile` em nível de pacote em `internal/kit`, `internal/dashboard` e `internal/harness` compilam no primeiro uso, por exemplo com `sync.OnceValue`. Verificável: `GODEBUG=inittrace=1 dh snapshot event < /dev/null 2>&1 | grep -E "internal/(kit|dashboard|harness)"` mostra cada um dos três abaixo de 0.05 ms e abaixo de 100 alocações; o laço de 200 chamadas não fica mais lento que 0.793 s na mesma máquina (reporte antes e depois); todos os testes existentes verdes.
- R4 O dashboard mantém um cache de parse com chave por caminho, tamanho e data de modificação do arquivo.
  - Um arquivo de ticket ou de PRD que não mudou não é relido nem reparseado em um refresh.
  - Um arquivo que mudou é reparseado.
  - A consulta ao cache roda depois das mesmas checagens de `os.Lstat`, de arquivo regular e de tamanho que `snapshot.ReadCapped` roda. Um arquivo recusado dá o mesmo aviso `skipped` com e sem o cache.

  Verificável: um teste em Go mostra que um arquivo que mudou é reparseado, um arquivo que não mudou vem do cache, e um ticket trocado por um symlink é recusado com o mesmo aviso; `buildState` devolve o mesmo resultado com e sem o cache nos mesmos dados; um benchmark descartável (não commitado) mostra ns/op pelo menos 3x menor em chamadas repetidas com dados que não mudaram; o corpo do PR reporta antes e depois.
- R5 O corpo do PR e o `CHANGELOG.md` dão o antes e o depois medidos da R3 e da R4. Eles nomeiam o candidato a continuação: um cache de leitura compartilhado para o `validate`. Nenhum benchmark novo commitado e nenhum eval novo existem. Verificável: `grep -n "validate" CHANGELOG.md` encontra a continuação na entrada 0.28.0; `git diff main...HEAD -- '*.go' | grep -c '^+func Benchmark'` é 0; `git diff --name-only main...HEAD -- evals/` é vazio.
- R6 `.agents/coordinator.md` "Seleção de modelo" ganha no máximo 10 linhas novas. As linhas declaram estas regras:
  - (a) Eleve o modelo ou o esforço quando a ambiguidade, o acoplamento ou a consequência é alta. O padrão é a escolha barata.
  - (b) Mude o esforço no modelo atual antes de mudar o modelo.
  - (c) Não delegue uma cadeia dependente nem um trabalho que cabe em um contexto.
  - (d) Não use o Haiku para reconhecimento de segurança.
  - (e) O despacho informa o esforço quando ele é elevado. O log o registra.

  As linhas não têm nenhum número de benchmark, nenhum número de release de modelo, nenhum slug de máquina e nenhum caminho absoluto. Verificável: `grep -n -i "effort" .agents/coordinator.md` encontra as linhas; `git diff --stat` mostra no máximo 10 linhas acrescentadas nesse arquivo; `dh validate` verde.
- R7 `AGENTS.md` ganha três edições. Ele continua em português.
  - Um item com a premissa de limpeza.
  - Um item ou edição declara três coisas. Os 5 binários são gerados e commitados a cada release. O dashboard continua. O pprof otimiza o binário e não o poda.
  - O primeiro item perde a frase da cópia idêntica para o template de PRD e diz que o template de PRD existe só em `templates/en/`.

  Verificável: `grep -n "lixo\|limpeza" AGENTS.md` encontra a premissa; `grep -n "pprof" AGENTS.md` encontra o item dos binários; `grep -n "cópia idêntica" AGENTS.md` é vazio.
- R8 Versão 0.28.0. `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte (`v0.28.0`) no mesmo arquivo mudam juntos. Um `ref` esquecido é uma falha de entrega. A prova de clone limpo vale, porque o PR toca `internal/` e `dist/*/bin`. O `CHANGELOG.md` tem uma entrada 0.28.0 em português. Os dois roadmaps a refletem. Verificável: o `grep` dos três valores mostra `0.28.0` ou `v0.28.0`; `git clone . <tmp>`, `go run ./cmd/dh build` ali, e `sha256sum` igual dos binários; `dh validate` verde.
- R9 O comportamento é preservado. Verificável: `go test ./...` verde; `claude plugin test` verde em `dist/claude-code/dev-harness`; a saída do `dh validate` para o repositório é a mesma de antes, exceto pela contagem de templates.
- R10 Trava de escopo. As rotas e a saída JSON do dashboard não mudam. Nenhuma dependência nova existe. O `go.mod` não muda. Os outros templates não mudam. Os PRDs antigos não são reescritos. Verificável: `git diff go.mod` é vazio; `git diff --name-only main...HEAD` não lista nenhum outro template nem PRD antigo (`PRD-020` e o seu arquivo de revisão são histórico).

## 4. Docs

Conferidos com `grep` em 2026-10-13. Só entram os arquivos que carregam uma afirmação afetada.

- `AGENTS.md` (primeiro item; premissa de limpeza; item dos binários e do pprof)
- `README.md`, `README.en.md` (linha 63, `versões espelhadas` / `mirrored versions`, e a linha de versão)
- `docs/tutoriais/templates.html`, `docs/tutoriais/templates.md`, `docs/en/tutorials/templates.html`, `docs/en/tutorials/templates.md` (a afirmação de espelho no início; a seção do PRD)
- `docs/roadmap.html`, `docs/en/roadmap.html` (linhas cerca de 678 e 685; entrada 0.28.0)
- `docs/dev-harness-mapa.excalidraw` (nota sobre os templates espelhados, linha cerca de 37)
- `CHANGELOG.md` (0.28.0)
- `.claude-plugin/marketplace.json` (`version` e `ref`)
- `harness-manifest.json` lista 13 templates depois de `go run ./cmd/dh build`; não o edite à mão. `internal/kit/doc_counts.go` não ganha nenhuma mudança de contagem. Só as suas duas regexps compilam no primeiro uso (R3).

`docs/tutorial.html` e `docs/en/tutorial.html` não carregam nenhuma afirmação de espelho. Eles não entram na lista.

## 5. Fora de escopo

- Uma divisão do dashboard, UPX, um servidor HTTP próprio, uma reescrita do histórico, e qualquer remoção de um binário ou de uma plataforma.
- Download de binários sob demanda. Todos os 5 binários continuam commitados (ADR-001).
- Um cache de leitura compartilhado para o `dh validate`. É um candidato a continuação.
- Um benchmark commitado e um eval novo.
- O cache local do plugin. O dono mantém ali só a versão mais recente, fora do repositório.
- Uma mudança em outros templates e uma reescrita dos PRDs antigos.

## 6. Decisões do dono

- Dono, 2026-10-13: o dashboard continua no binário único. O pprof serve para melhorar e otimizar o binário inteiro, nunca para remover uma parte.
- Dono, 2026-10-13: manter só a versão mais recente do plugin no cache local. Isso é feito fora do repositório e não é uma regra aqui.
- Dono, 2026-10-13: os binários das 5 plataformas são gerados a cada release e continuam commitados.
- Dono, 2026-10-13: uma régua de modelo e esforço entra em `.agents/coordinator.md`.
- Dono, 2026-10-13: remover `templates/pt-br/PRD.md` e corrigir todo lugar que depende dele. O PRD é sempre em inglês. Isto substitui a decisão do PRD-020 de manter a cópia idêntica.
- Premissa do dono, 2026-10-13: toda mudança remove os restos que ela torna obsoletos. Lixo remanescente bloqueia a entrega.
- Decisão do Coordenador, 2026-10-13, sob a delegação do dono: o cache de leitura do `dh validate` é uma continuação, fora de escopo.
- Decisão do Coordenador, 2026-10-13, sob a delegação do dono: o conteúdo da régua da R6, itens (a)-(e).
- Premissas do analista, não decisões do dono:
  - (a) O alvo da inicialização preguiçosa são as 34 regexps em nível de pacote nos três pacotes da R3. `internal/doctor` mantém a sua 1. As 3 regexps dentro de funções em `validate.go` ficam, porque rodam só no `validate`.
  - (b) O cache da R4 vive em memória no processo do dashboard. Ele descarta as entradas dos caminhos que não apareceram no último build.
  - (c) A versão é 0.28.0 porque a mudança acrescenta comportamento sem quebra.
  - (d) `harness-manifest.json` muda só por `dh build`, não à mão. `internal/kit/doc_counts.go` não ganha nenhuma mudança de contagem. Só as suas duas regexps compilam no primeiro uso (R3).
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** o template de PRD existe só em `templates/en/`, com uma checagem no `dh validate`; as regexps em nível de pacote compilam no primeiro uso; o dashboard mantém um cache de parse; o Coordenador tem uma régua de modelo e esforço; `AGENTS.md` tem a premissa de limpeza e a regra dos binários; versão 0.28.0.
- **Regras:**
  - R1 Honrada: `templates/pt-br/PRD.md` apagado; `englishOnlyTemplates` em `internal/kit/validate.go`; três casos de teste em Go; checagem de mutação fica vermelha quando a lista fica vazia.
  - R2 Honrada: o manifesto lista 13 templates; os roadmaps dizem 13 e "7 em en/, 6 em pt-br/"; README, tutoriais, nota do excalidraw e primeiro item do AGENTS.md corrigidos; o `grep` de espelho só mostra histórico ou outros templates; `dist/` não tem `pt-br/PRD.md`.
  - R3 Honrada: inittrace kit 0.01-0.02 ms/62 alocações, dashboard 0.01 ms/51, harness 0 ms/6; laço de 200 chamadas 0.66-0.68 s contra 0.75-0.80 s em `main`.
  - R4 Honrada: cache com chave por caminho, tamanho e mtime depois das travas de `Lstat`, de arquivo regular e de tamanho; cinco testes em Go; `buildState` de 50-58 ms para 14-15 ms (cerca de 3.8x), de 11.5 MB para 0.67 MB.
  - R5 Honrada: o CHANGELOG 0.28.0 traz a tabela de antes e depois e a continuação do cache de leitura do `validate`; nenhum benchmark commitado; nenhuma mudança em `evals/`.
  - R6 Honrada: quatro linhas em `.agents/coordinator.md` "Seleção de modelo", sem números, slugs nem caminhos.
  - R7 Honrada: premissa de limpeza na linha 12 do `AGENTS.md`, frase dos binários e do pprof na linha 16, "cópia idêntica" removida.
  - R8 Honrada: `defaultVersion`, `version` do marketplace e `ref` são 0.28.0 e v0.28.0; CHANGELOG e os dois roadmaps; prova de clone limpo rodada na hora do commit.
  - R9 Honrada: `go test ./... -count=1`, `-race` em dashboard e kit, `claude plugin test` (254 passam) e `dh validate` verdes.
  - R10 Honrada: `go.mod`, rotas e JSON sem mudança; nenhum outro template tocado; nenhum PRD antigo reescrito.
- **Tickets:** T-2101 (backend, builder) aprovado; T-2102 (backend, harness-maintainer) aprovado. Matriz de QA, revisão de código por `claude` (`approve`), revisão final FEATURE APROVADA.
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutoriais/templates.html` e `.md`, `docs/en/tutorials/templates.html` e `.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/dev-harness-mapa.excalidraw`, `CHANGELOG.md`, `.claude-plugin/marketplace.json`.
- **Fora de escopo:** divisão do dashboard, download de binários, cortes de plataforma, UPX, reescrita do histórico, um cache de leitura compartilhado para o `dh validate` (continuação), evals novos.
