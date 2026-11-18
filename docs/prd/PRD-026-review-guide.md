# PRD-026 · Guia de revisão (REVIEW.md) e modo de revisão de PR

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-21 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-21 / 2026-10-21 |
| Tickets | T-2601, T-2602, T-2603, T-2604 (planejados) |
| Origem | Conversa com o dono, 2026-10-21 (descoberta DISC-006; vídeo `"How Anthropic Does Code Review (What You Can Copy)"`, pasquadev, 2026-10-21, https://youtu.be/c6xOFjM2IbE) |
| Depende de | nada |
| Habilita | nada |

## 1. Problema

A revisão de código do kit não tem guia específico do projeto. O revisor julga um diff só com regras genéricas. Fato: `.skills/code-review/SKILL.md` e `.agents/reviewer.md` não nomeiam nenhum arquivo por projeto que diga o que procurar.

Disso decorrem quatro lacunas. Fato: cada lacuna é visível no texto atual do kit.

- Um falso achado repetido custa uma rodada a cada vez. Nada registra que um padrão é ruído, então a revisão seguinte o reporta de novo.
- O revisor reporta uma instância de um bug, e a rodada seguinte encontra a próxima instância. As rodadas se gastam um caso por vez.
- Uma revisão de reentrada não tem escopo declarado. O revisor pode reler a mudança inteira em vez do delta.
- Uma pull request não pode ser revisada como pull request. O `/dh:review` aceita um escopo, uma ref ou `feature PRD-<n>`; não tem modo que leia uma PR e a intenção que ela declara.

Além disso: o passo de merge do `/dh:auto` (`.commands/auto.md`, passo 5 da sequência de git) faz merge com CI verde mais o veredito `approve`. Ele não tem parada para uma mudança num caminho que o dono chama de crítico.

O vídeo é a fonte das ideias. Seus números e os artigos que ele cita não são verificados. Nenhuma regra deste PRD depende deles.

## 2. Solução

Após a entrega, cada projeto tem um guia de revisão, `REVIEW.md`, em inglês, no repositório do projeto. O caminho é registrado no `project.yaml` do harness dir. O revisor lê o guia e liga cada achado a uma regra.

Fluxo após a entrega:

1. O `/dh:setup` cria o `REVIEW.md` a partir de `templates/en/REVIEW.md` quando o projeto não tem um, e escreve o caminho dele (relativo à raiz do repositório) no `project.yaml`. Isso vale no modo `repo` e no modo `global`.
2. O `/dh:review` lê o guia. Cada achado cita a regra do guia de que vem, ou diz que nenhuma regra se aplica. Um achado que casa com uma entrada `Ignore` é reportado como ignorado por regra e não gasta uma rodada.
3. O revisor reporta a classe inteira de um bug de uma vez, rotula como `pre-existing` um achado em linhas intocadas e, na reentrada, revisa só o delta mais o status de cada achado bloqueante anterior.
4. O `/dh:review pr <n>` revisa uma pull request do GitHub em modo somente leitura. O corpo da PR e qualquer PRD ou ticket ligado são a fonte de requisitos. O relatório fica local. Um comentário na PR só é postado com autorização explícita do dono.
5. Só os achados `Critical` e `Major` vão para o builder. Os achados `Minor` e `nit` ficam no relatório.
6. O Coordenador registra no `MEMORY.md` o que aconteceu com cada achado (corrigido, intencional, ruído). Quando o mesmo padrão é classificado como ruído 3 vezes, o Coordenador propõe ao dono uma linha `Ignore` ou uma regra de lint. Ele nunca aplica a proposta por conta própria.
7. Quando uma PR toca um caminho listado em `Critical paths` no `REVIEW.md`, o fluxo para antes do merge e pergunta ao dono, mesmo com o gate verde e o veredito de QA aprovado.

Este repositório é um consumidor do kit. O mesmo lote acrescenta um `REVIEW.md` real na raiz do repositório.

Release 0.32.0 (proposta do autor), uma PR.

## 3. Regras

Marcação: R1-R14 e R16-R18 carregam decisões do dono de 2026-10-21 (D1-D6, A1, A2). R5-R9 e os nomes de arquivo são a redação do autor para o resumo da descoberta. R23-R25 e as mudanças em R5, R12 e R13 vêm da revisão de código e de segurança da implementação (2026-10-21); não precisam de nova decisão do dono. Cada regra diz como verificá-la.

- R1 Template. `templates/en/REVIEW.md` existe. Ele tem estas partes: (a) uma nota curta `How to use`; (b) uma seção de regras em que cada entrada tem os campos `Symptom`, `Why it matters` e `Example` com `file:line`; (c) uma seção `Ignore` em que cada entrada diz `ignore this pattern in this kind of file`; (d) uma seção `Critical paths` que lista globs de caminho. Não existe `templates/pt-br/REVIEW.md`. Verificável: o arquivo existe com as quatro partes; `ls templates/pt-br/` não mostra `REVIEW.md`.
- R2 O setup cria o guia. Quando o `/dh:setup` roda e o projeto não tem arquivo de guia de revisão, o setup cria `REVIEW.md` na raiz do repositório a partir de `templates/en/REVIEW.md` e escreve o caminho dele no `project.yaml`, na chave `review_guide` (proposta do autor para o nome da chave). O caminho é relativo à raiz do repositório. O arquivo é sempre em inglês, seja qual for o valor de `language`. Verificável: o texto de `.commands/setup.md` e de `.skills/project-onboarding/SKILL.md` declara esses passos, o nome da chave e a regra do inglês.
- R3 Os dois modos. A R2 vale no modo `repo` e no modo `global`. No modo `global`, `REVIEW.md` é o único arquivo que o setup escreve na árvore do repositório. As frases existentes `nothing is written in the repository tree` (`.commands/setup.md`, `.skills/project-onboarding/SKILL.md`, `docs/tutorial.html`, `docs/en/tutorial.html`) dizem essa exceção. Verificável: `grep -n -i 'repository tree' .commands/setup.md .skills/project-onboarding/SKILL.md docs/tutorial.html docs/en/tutorial.html` mostra cada ocorrência com a exceção declarada.
- R4 O setup nunca sobrescreve. Quando o projeto tem o arquivo mas não tem a chave, o setup registra a chave para o arquivo existente e deixa o arquivo intocado. Quando a chave aponta para um arquivo existente, o setup não muda nada. Um especialista despachado escreve o arquivo, não o Coordenador (o arquivo é novo e maior que o limite de pequena correção do PRD-025). Verificável: o texto do setup declara os três casos e nomeia o escritor; o formato de relatório do setup lista `REVIEW.md` em `Created` ou `Left untouched`.
- R5 A revisão é somente leitura e aponta para o setup. Quando o `/dh:review` roda e o `project.yaml` não tem a chave `review_guide`, ou o arquivo nesse caminho não existe na revisão confiável (R23), ou o caminho é inseguro (R24), então o relatório diz `no review guide`, dá o motivo e aponta para o `/dh:setup`. A revisão não cria o arquivo. A revisão ainda roda com as regras genéricas, e o relatório marca o eixo do guia como `not-run` com o motivo. Verificável: `.commands/review.md` declara isso; `.skills/code-review/SKILL.md` tem um passo correspondente.
- R6 Os achados citam o guia. Quando o agente `reviewer` revisa código com um guia presente, então ele lê primeiro o guia no caminho vindo do `project.yaml`. Cada achado traz uma linha `Guide rule: <entry title>` ou `Guide rule: none applies`. Verificável: o formato de saída de `code-review` e o formato de saída do agente `reviewer` têm essa linha; o exemplo de `code-review` a mostra.
- R7 Entradas `Ignore`. Quando um achado casa com uma entrada `Ignore` do guia, então o relatório o lista numa lista separada `Ignored by rule` com o título da entrada. Ele não recebe severidade, não bloqueia e não gasta uma rodada. Verificável: o formato de saída de `code-review` tem a lista; o texto do Coordenador sobre achados bloqueantes (`.agents/coordinator.md`, `Proportional review`) conta só achados `Critical` e `Major` que não são ignorados por regra.
- R8 Classe inteira, `pre-existing`. Quando o revisor encontra um bug, então ele reporta num só achado todas as instâncias dessa classe dentro do escopo da revisão, com cada local listado. Um achado em linhas que a mudança não tocou é rotulado `pre-existing` e não bloqueia, a menos que a mudança torne esse código alcançável ou pior; então ele mantém a severidade e o achado diz por quê. Verificável: as duas regras aparecem no procedimento de `.skills/code-review/SKILL.md`, e o passo de controle de falso positivo (passo 9) concorda com o rótulo.
- R9 Escopo de reentrada. Quando uma revisão roda de novo sobre a mesma mudança, então o revisor lê (a) o delta desde o último estado revisado e (b) o status de cada achado bloqueante anterior (corrigido, não corrigido, retirado). Ele reporta só sobre esses dois itens, mais qualquer achado `Critical` ou `Major` novo que encontre no delta. O Coordenador nomeia o último estado revisado no despacho. Verificável: a skill `code-review` tem uma seção de reentrada com esses dois itens; `.agents/coordinator.md` nomeia o último estado revisado no despacho de reentrada.
- R10 Registro de ruído. O Coordenador registra, para cada achado de revisão, um resultado no `MEMORY.md`: corrigido, intencional ou ruído. Quando o mesmo padrão é classificado como ruído 3 vezes, então o Coordenador propõe ao dono uma linha `Ignore` para o `REVIEW.md` ou uma regra de lint. O Coordenador nunca aplica a proposta por conta própria. Um especialista despachado escreve a mudança só depois que o dono aprova. Verificável: `.agents/coordinator.md` declara o registro, o limiar 3, a proposta e a proibição.
- R11 O despacho de correção leva só `Critical` e `Major`. Quando o Coordenador envia achados de revisão a um builder, então o despacho contém só achados `Critical` e `Major`. Os achados `Minor` e `nit` ficam no relatório de revisão. Verificável: `.agents/coordinator.md` (item 8 de `Procedure` e o texto do ciclo de revisão) e todo outro lugar que descreva o despacho de correção ao builder (`grep -n -i 'findings' .commands/build.md .commands/fix.md .commands/review.md .agents/builder.md`) dizem isso; nenhum diz `all findings`.
- R12 Caminhos críticos param o merge. Quando uma PR toca um caminho que casa com um glob em `Critical paths` no `REVIEW.md` do projeto, então o fluxo para antes do merge e pede ao dono um ok explícito, mesmo com o gate verde e o veredito de QA aprovado. A pergunta nomeia os caminhos casados. A correspondência de caminho é calculada por um `qa-verifier` despachado, a partir da lista de arquivos da PR. Quando o `REVIEW.md` está ausente ou a seção está vazia, nenhuma parada se aplica. Esta regra se sobrepõe à instrução global de auto-merge do dono só para esses caminhos; para todo outro caminho a instrução global não muda. Uma mudança no guia de revisão para o merge do mesmo modo (R23). Verificável: `.commands/auto.md` (passo 5 do merge na sequência de git) e o texto de conclusão do Coordenador declaram a parada, a pergunta ao dono, o escopo da sobreposição e a parada por mudança no guia.
- R13 Modo de revisão de PR. O modo PR só começa quando o primeiro token do argumento é exatamente `pr`. Quando `pr` é seguido de uma sequência de dígitos ASCII (`[0-9]+`) e nada mais (por exemplo `pr 44`), o comando roda no modo PR. Quando `pr` é seguido de qualquer outra coisa, o comando reporta o gatilho malformado e para (falha fechada). Qualquer outro argumento (por exemplo `prod`, `preview/x`, `prototype.go`) é o modo padrão, como antes deste PRD. No modo PR, um `qa-verifier` despachado lê a PR com `gh` em somente leitura (`gh pr view <n> --json ...`, `gh pr diff <n>`), e não o Coordenador. O comando toma o corpo da PR e qualquer referência `PRD-<digits>` ou `T-<digits>` ligada, encontrada nele, como fonte de requisitos, e despacha o `reviewer` no modo `code` como no modo padrão. O relatório declara o commit que leu. Ele não faz checkout de branch e não escreve na árvore. Não existe comando novo. Verificável: `.commands/review.md` tem o modo, a regra do token exato, a regra da falha fechada, a queda para o modo padrão de `prod`, `preview/x` e `prototype.go`, as chamadas `gh` somente leitura e o despacho do `qa-verifier`; os modos existentes não mudam.
- R14 A postagem precisa de autorização. No modo `pr <n>` o relatório fica local. O comando posta um comentário na PR do GitHub só quando o dono autoriza essa postagem de forma explícita na conversa atual. A autorização cobre uma postagem. Verificável: `.commands/review.md` declara isso e nomeia a postagem como ação externa.
- R15 Lista de templates no `dh validate`. A fatia `englishOnlyTemplates` em `internal/kit/validate.go` (linha 35) lista `"REVIEW.md"` ao lado de `"PRD.md"`. O comentário acima dela (linha 34) diz que o PRD e o guia de revisão são sempre escritos em inglês. Um caso de teste em `internal/kit/validate_test.go` cobre `REVIEW.md` do mesmo modo que cobre `PRD.md` (uma cópia em `templates/pt-br/` é erro; a falta de cópia em `templates/pt-br/` não é). Verificável: `go test ./internal/kit/...` passa; `go run ./cmd/dh validate .` passa com `templates/en/REVIEW.md` presente e sem cópia em pt-br.
- R16 Este repositório é um consumidor. O lote acrescenta `REVIEW.md` na raiz do repositório, em inglês. Ele tem pelo menos 3 entradas reais para este kit, cada uma com `Symptom`, `Why it matters` e `Example` com `file:line`: `dist/` gerado editado à mão; uma contagem declarada nos docs que se afasta da contagem real (agentes, comandos); um caminho privado num PRD ou doc. Ele tem uma seção `Ignore` e uma seção `Critical paths` (decisão do dono D7: os globs críticos são exatamente `adapters/claude-code/plugin/hooks/**` e `.skills/external-clis/**`). A chave `review_guide: REVIEW.md` é escrita no `project.yaml` do harness dir deste projeto, seja qual for o modo em que o projeto resolve. O Coordenador escreve essa única linha; é um passo do Coordenador, não uma escrita do builder. Verificável: o arquivo existe com as partes acima; `dh harness-path --json` dá o harness dir, e o `project.yaml` dele tem a chave.
- R17 Itens fora de escopo continuam fora. Este PRD não acrescenta nenhum dos itens da seção 5. O limite de rodadas continua 6 (`.agents/coordinator.md`, `Proportional review and paid evaluation`). A regra do veredito fundido "qualquer achado bloqueante de qualquer revisor bloqueia" não muda. Verificável: `git diff main...HEAD -- .agents/coordinator.md` não mostra mudança no número 6, nem texto de voto majoritário, nem `dh review run`.
- R18 Contagens. As contagens continuam em 13 agentes e 19 comandos + 1 comando de mod. Nenhum agente, skill ou comando novo. Verificável: `go run ./cmd/dh validate .` passa, e `git diff --name-status main...HEAD -- .agents .skills .commands` não lista arquivo acrescentado nem removido.
- R19 Nenhum rastro do arquivo de diretrizes. O arquivo de síntese de diretrizes do vídeo (um arquivo pt-br em `docs/`) e o irmão `Zone.Identifier` dele não estão no repositório no fechamento. O dono move o arquivo para o cofre; o Coordenador cuida disso fora do kit. Nenhum arquivo, cópia em `dist/`, doc ou texto de PRD cita o caminho dele. O vídeo é citado só por título, canal, data e URL. Verificável, todos sem saída: `git ls-files | grep -i 'DIRETRIZES-CODE-REVIE[W]'`; `find . -name '*Zone.Identifier*' -not -path './.git/*'`; `grep -rIl 'DIRETRIZES-CODE-REVIE[W]' . --exclude-dir=.git`; `git status --short` não mostra nenhum dos dois arquivos.
- R20 Limpeza. Toda referência obsoleta que esta mudança cria é removida (AGENTS.md, "Premissa de limpeza"). Exemplos: frases que dizem que o setup não escreve nada na árvore do repositório (R3); frases que dizem que o `/dh:review` aceita só um escopo, uma ref ou `feature PRD-<n>` (R13); frases que dizem `the PRD is the only English-only template` (`grep -rn -i 'only.*PRD\|PRD.*only' .skills .commands .agents templates docs README.md README.en.md`); metadado `writes:` do setup em `.commands/setup.md` que omite `REVIEW.md`; texto de `.commands/review.md` que lista os argumentos permitidos sem `pr <digits>`. Verificável: o autor do ticket roda cada grep e lista as ocorrências e a ação tomada no relatório do ticket.
- R21 Release e arquivos gerados. Release 0.32.0 em três lugares: `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte (`v0.32.0`). `dist/` não é editado à mão; `go run ./cmd/dh build` o regenera, e um segundo build não muda nada exceto a linha `generated_at`. `go run ./cmd/dh validate .` passa com 0 erros e `go test ./...` fica verde. A prova de clone limpo se aplica (AGENTS.md, "Prova de clone limpo"), porque a PR toca `internal/` (`validate.go`, `build.go`) e `dist/*/bin`: `git clone . <tmp>`, `go run ./cmd/dh build` lá, e o `sha256sum` dos 5 binários é igual ao da árvore. Verificável: rode os comandos e compare com `git diff`.
- R22 Nenhum caminho privado no PRD, nos tickets, no novo `REVIEW.md` ou em qualquer texto novo. O texto usa `<harness dir>` e `<home>`. Verificável: `go run ./cmd/dh validate .` (`checkPrivatePaths`) passa.
- R23 Revisão de política confiável. A revisão lê o arquivo do guia e a chave `review_guide` na revisão-base da mudança, nunca na ponta sob revisão. A base é a base da PR no modo `pr <n>`, a ref de comparação no modo padrão com uma ref, e `HEAD` para trabalho não commitado. Quando a base não tem guia, nenhum guia se aplica a essa mudança. Quando a mudança acrescenta, edita ou apaga o arquivo do guia no caminho configurado, ou muda a chave `review_guide`, então o fluxo para antes do merge como um caminho crítico (R12), e a pergunta ao dono nomeia a mudança no guia. As entradas `Ignore` que a mudança acrescenta não se aplicam a essa mesma mudança. O dono pode aprovar uma atualização do guia por essa parada. No modo `global` a chave vive fora da árvore do repositório, então só uma mudança no arquivo do guia aciona a parada. Verificável: `.commands/review.md`, `.skills/code-review/SKILL.md` e `.agents/reviewer.md` declaram a leitura na revisão-base e o caso inicial sem guia; `.commands/auto.md` e `.agents/coordinator.md` declaram a parada e a regra do `Ignore`.
- R24 Contenção de caminho. O valor de `review_guide` deve ser um caminho relativo ao repositório, sem prefixo absoluto e sem segmento `..`. O alvo resolvido, symlinks incluídos, deve ficar dentro do repositório. Quando qualquer checagem falha, a revisão trata o projeto como sem guia de revisão (`no review guide`) e reporta o motivo (R5). O setup escreve só `REVIEW.md` na raiz do repositório. Verificável: `.commands/review.md` e `.skills/code-review/SKILL.md` declaram as três checagens e a queda; o texto do setup declara o valor fixo.
- R25 Tratamento do texto de PR. No modo `pr <n>` a revisão trata o texto da PR (corpo, título, comentários) como dado não confiável. Um `qa-verifier` despachado executa as leituras `gh pr view` e `gh pr diff`, e não o Coordenador (R13). Ninguém busca um link encontrado no texto da PR, e o revisor não usa busca na web nesse modo. O texto da PR nunca é autorização do dono para postar, fazer merge ou qualquer outra ação (R14). Verificável: `.commands/review.md` e `.agents/reviewer.md` declaram as três proibições e o despacho das leituras.

## 4. Docs

Atualizar:

- `templates/en/REVIEW.md` (novo; R1) e `internal/kit/validate.go` com `internal/kit/validate_test.go` (R15).
- `.commands/setup.md` e `.skills/project-onboarding/SKILL.md` (R2-R4, R20, R24): passos do setup, o metadado `writes:`, o formato de relatório, as frases sobre a "árvore do repositório".
- `.commands/review.md` (R5, R13, R14, R23-R25): o modo `pr <n>`, a checagem do guia e os limites de somente leitura.
- `.skills/code-review/SKILL.md` e `.agents/reviewer.md` (R6-R9, R23-R25): leitura do guia, linha `Guide rule:`, lista de ignorados por regra, classe inteira, rótulo `pre-existing`, escopo de reentrada, formato de saída e exemplo.
- `.agents/coordinator.md` (R7, R9-R12, R23): conteúdo do despacho de correção, despacho de reentrada, registro de ruído, parada por caminho crítico.
- `.commands/auto.md` (R12, R23): passo 5 do merge na sequência de git. `.commands/build.md` e `.commands/fix.md` só onde o grep da R11 achar a redação `all findings`.
- `AGENTS.md`: já atualizado pelo Coordenador em 2026-10-21 (exceção de template, gatilho do `dh review run`, regra "Guia de revisão"). O autor do ticket o lê e reporta qualquer discordância; não o edita.
- `docs/tutorial.html`, `docs/en/tutorial.html`: a descrição do setup, a descrição da revisão, a chave do `project.yaml`, a parada por caminho crítico, o modo `pr <n>`.
- `README.md` e `README.en.md`: só onde `grep -n -i 'review\|setup'` achar uma frase que esta mudança torne falsa.
- `docs/roadmap.html`, `docs/en/roadmap.html`: um card Delivered para a 0.32.0 nos dois arquivos, e o chip de status e a linha de fontes. O `dh validate` confere a deriva contra o `CHANGELOG.md`.
- `CHANGELOG.md` (entrada 0.32.0).
- `.claude-plugin/marketplace.json`, `internal/build/build.go` (R21).
- `REVIEW.md` na raiz do repositório (novo; R16).

Não existem `pdocs` para este projeto.

Não afetados (o autor do ticket confere com `grep`):

- `docs/plano-produto.html`: histórico interno em pt-br; muda só se declarar o comportamento de revisão como atual.
- `docs/adr/`: nenhum ADR novo, a menos que o dono peça.
- `docs/inicio-rapido.html`, `docs/en/quick-start.html`: mudam só se descreverem a saída do setup.
- `.skills/external-clis/`: continua em prosa (D2); sem mudança.

## 5. Fora de escopo

- Voto majoritário entre revisores. A regra "qualquer achado bloqueante de qualquer revisor bloqueia" continua.
- Avaliações de incidentes (`evals/incidentes/`). A avaliação paga roda só quando o dono pede (`AGENTS.md`).
- Um job de revisão no CI por PR. O CI roda sem chave de API de propósito.
- Agentes classificadores de risco.
- Uma fila de merge.
- Feature flags e implantação gradual.
- Métricas pós-merge.
- O subcomando `dh review run`. O `external-clis` continua em prosa. O gatilho para ele (uma chamada de CLI que precisa de um kill ou de uma nova execução manual dentro de um lote) está só no `AGENTS.md`.
- Um limite de 3 rodadas, e qualquer mudança no limite de rodadas de 6.
- Um agente, skill ou comando novo, e qualquer mudança nas contagens.
- Uma checagem do `dh doctor` para o guia (ver Apêndice, sugestão S1).

## 6. Decisões do dono

- D1, 2026-10-21: estender o `/dh:review` com um modo `pr <n>`. Nenhum comando novo; a contagem continua 19 + 1. O corpo da PR mais qualquer PRD ou ticket ligado é a fonte de requisitos. O relatório fica local. Um comentário na PR precisa de autorização explícita (R13, R14).
- D2, 2026-10-21: o `external-clis` continua em prosa. O gatilho do `dh review run` fica registrado só no `AGENTS.md`. Fora de escopo aqui.
- D3, 2026-10-21: o `REVIEW.md` tem uma seção de caminhos críticos. Uma PR que toca um caminho listado para antes do merge e pergunta ao dono. Isso se sobrepõe à instrução global de auto-merge só para esses caminhos (R12).
- D4, 2026-10-21: os achados `Minor` e `nit` ficam no relatório. Só `Critical` e `Major` vão para o builder (R11).
- D5, 2026-10-21: o limite de rodadas continua 6 (R17).
- D6, 2026-10-21: o arquivo de diretrizes e o irmão `Zone.Identifier` saem do repositório no fechamento (R19).
- A1, 2026-10-21: o `REVIEW.md` é sempre adicionado ao projeto. O caminho dele é registrado no `project.yaml` do harness dir em que o projeto resolve, relativo à raiz do repositório. O arquivo vive no repositório do projeto (R2-R4, R16).
- A2, 2026-10-21: o `REVIEW.md` é sempre em inglês. O template existe só em `templates/en/` (R1, R15).
- D7, 2026-10-21: os `Critical paths` do `REVIEW.md` deste repositório são exatamente `adapters/claude-code/plugin/hooks/**` (trava de runtime, RISK-001) e `.skills/external-clis/**` (isolamento do revisor de CLI, RISK-002). Uma release que só sobe a versão e os binários não para (R16).

## Apêndice

### Tickets sugeridos

| Ticket | `lane` | Executor | Conjunto de escrita |
|---|---|---|---|
| T-2601 | `backend` | `harness-maintainer` com `harness-authoring` | `templates/en/REVIEW.md`; `.commands/setup.md`; `.skills/project-onboarding/SKILL.md`; `.commands/review.md`; `.skills/code-review/SKILL.md`; `.agents/reviewer.md`; `.agents/coordinator.md`; `.commands/auto.md`; `.commands/build.md` e `.commands/fix.md` se o grep da R11 achar ocorrência. Regras R1-R14, R17, R18, R20, R23-R25 |
| T-2602 | `backend` | `builder` (Go) | `internal/kit/validate.go`; `internal/kit/validate_test.go`; `internal/build/build.go` (versão); `.claude-plugin/marketplace.json`. Regras R15, R21. A prova de clone limpo se aplica (toca `internal/` e `dist/*/bin`) |
| T-2603 | `backend` (fonte do kit e docs) | `harness-maintainer` com `harness-authoring`, depois `docs-guide` para os docs | `REVIEW.md` (raiz); `docs/tutorial.html`; `docs/en/tutorial.html`; `README.md`; `README.en.md`; `docs/roadmap.html`; `docs/en/roadmap.html`; `CHANGELOG.md`. Regras R3, R16, R20, R22 |
| T-2604 | `teste (unitário e integração)` | `qa-verifier` | Teste da R15 em `internal/kit/validate_test.go` (com o T-2602); execuções de R18, R19, R21, R22; uma leitura a seco do texto de R2-R14 e R23-R25 contra os greps listados em cada regra. Nenhum código de produção |

Ordem: T-2601 e T-2602 podem rodar em paralelo só se os conjuntos de escrita continuarem disjuntos (continuam). O T-2603 segue o T-2601 (os docs descrevem o texto final). O T-2604 segue todos. T-2602 e T-2604 escrevem o mesmo arquivo de teste: T-2602 primeiro.

Passo do Coordenador (não é escrita do builder): depois que o `REVIEW.md` existir na raiz do repositório, o Coordenador escreve a linha `review_guide: REVIEW.md` no `project.yaml` do harness dir deste projeto (R16). O Coordenador também move o arquivo de diretrizes para fora do repositório (R19).

A sequência do gate de release está na R21: `dh validate`, `go test ./...`, `go run ./cmd/dh build` sem diff em `dist/`, depois a prova de clone limpo, depois o candidato a commit; push, PR e tag só depois da prova.

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Um comando novo `/dh:review-pr` | O dono escolheu um modo do `/dh:review` (D1). A contagem 19 + 1 continua. |
| `REVIEW.md` também em `pt-br` | O dono escolheu só inglês (A2). |
| `REVIEW.md` guardado no harness dir | O dono escolheu o repositório do projeto (A1). O revisor num clone ou numa CLI o vê junto com a árvore. |
| A revisão cria o arquivo ausente | A revisão é somente leitura. O setup é dono da criação (R5). |
| Encaminhar achados `Minor` ao builder | O dono escolheu mantê-los no relatório (D4). Um builder corrige um `nit` e a revisão acha outro `nit`. |
| Limite de 3 rodadas | O dono manteve 6 (D5). |

### Riscos e decisões pendentes

- Risco: o setup no modo `global` agora escreve na árvore do repositório, o que os docs chamam de proibido · mitigação: a R3 muda cada frase, e a R20 faz grep da redação antiga.
- Risco: o mod de runtime avisa quando a sessão principal escreve um arquivo novo fora do harness dir · mitigação: a R4 manda a escrita para um especialista.
- Risco: a correspondência de caminho crítico usa globs, e o kit não define um estilo de glob · premissa: o autor do ticket usa o mesmo estilo de glob do `.gitignore` e escreve um exemplo no template. Hipótese: isso basta para os caminhos deste repositório (R16).
- Risco: no modo `pr <n>` o checkout local pode não ser o head da PR, então o rastreamento de chamadores do revisor lê a árvore errada · mitigação: a R13 exige que o relatório declare o commit que leu; o revisor marca como limitados os caminhos rastreados que não conseguiu verificar.
- Premissa: o nome da chave é `review_guide` e a release é a 0.32.0 (o dono não objetou, 2026-10-21).
- Sugestão S1 (não é requisito): uma checagem do `dh doctor` que avise quando o `project.yaml` não tem a chave `review_guide` ou o arquivo está ausente. A A1 é atendida sem ela, porque a R5 e a R2 cobrem os fluxos. Ela acrescenta código Go e um teste, então fica para o dono.
- Sugestão S2 (não é requisito): um limite de comentários por PR no modo `pr <n>`, para reduzir o ruído. O vídeo a sugere; o dono não decidiu.
- Decisão pendente: nenhuma bloqueia os tickets.

## Resumo executado

**Entregue**: Todo projeto recebe um `REVIEW.md` em inglês criado pelo setup, com `review_guide` no `project.yaml`; o revisor o lê na revisão-base e liga cada achado a uma regra; o `/dh:review pr <n>` revisa uma PR do GitHub em somente leitura; só os achados `Critical` e `Major` vão para o builder; caminhos críticos e mudanças no guia param o merge para o dono.

**Regras**:

Evidência para todas as regras: rodadas de QA 2-5 sobre o estado final (`dh validate` 0 erros e 0 avisos, `go test ./...` 7 pacotes ok, build idempotente com só `generated_at` mudando, `claude plugin test` 235 passam e 0 falham); revisão de aderência da feature final FEATURE APROVADA; revisão de código `cli:codex/gpt-6-luna` `approve`; revisão de segurança `cli:codex/gpt-6.1-sol` nenhuma vulnerabilidade demonstrada, e `claude` passa com risco residual. O PRD foi emendado depois da revisão da implementação (R5, R12, R13 mudaram; R23-R25 acrescentadas) e revalidado uma vez (aprovado).

- R1: Cumprida. `templates/en/REVIEW.md` tem as quatro partes; `templates/pt-br/` não tem `REVIEW.md`.
- R2: Cumprida. `.commands/setup.md`, `.skills/project-onboarding/SKILL.md` declaram os passos, a chave e a regra do inglês.
- R3: Cumprida. As frases `repository tree` em `.commands/setup.md`, `.skills/project-onboarding/SKILL.md`, `docs/tutorial.html`, `docs/en/tutorial.html` declaram a exceção.
- R4: Cumprida. `.commands/setup.md` declara os três casos e o especialista escritor.
- R5: Cumprida. `.commands/review.md`, `.skills/code-review/SKILL.md` (sem guia, eixo do guia `not-run`, aponta para o setup).
- R6: Cumprida. Linha `Guide rule:` em `.skills/code-review/SKILL.md` e `.agents/reviewer.md`.
- R7: Cumprida. Lista `Ignored by rule` em `.skills/code-review/SKILL.md`; `.agents/coordinator.md` conta só `Critical` e `Major` não ignorados.
- R8: Cumprida. Regras de classe inteira e `pre-existing` em `.skills/code-review/SKILL.md`.
- R9: Cumprida. Seção de reentrada em `.skills/code-review/SKILL.md`; `.agents/coordinator.md` nomeia o último estado revisado.
- R10: Cumprida. `.agents/coordinator.md` (registro, limiar 3, proposta, proibição).
- R11: Cumprida. `.agents/coordinator.md`, `.commands/build.md`, `.commands/fix.md`, `.commands/review.md`, `.agents/builder.md`; nenhum diz `all findings`.
- R12: Cumprida. `.commands/auto.md` e `.agents/coordinator.md`; o ensaio do QA não tocou nenhum glob crítico.
- R13: Cumprida. `.commands/review.md` (token `pr` exato, falha fechada, queda para o modo padrão, `gh` somente leitura, despacho do `qa-verifier`).
- R14: Cumprida. `.commands/review.md` (a postagem precisa de autorização explícita, uma postagem).
- R15: Cumprida. `internal/kit/validate.go` (`englishOnlyTemplates`) e `internal/kit/validate_test.go`; testes conferidos por mutação, `go test ./...` ok.
- R16: Cumprida. `REVIEW.md` da raiz com 4 regras, 2 entradas `Ignore` e os dois globs da D7; `review_guide: REVIEW.md` no `project.yaml`.
- R17: Cumprida. Limite de rodadas 6 e regra do veredito fundido sem mudança em `.agents/coordinator.md`; sem voto majoritário, sem `dh review run`.
- R18: Cumprida. `dh validate` mostra 13 agentes e 19 comandos + 1 comando de mod; nenhum arquivo acrescentado ou removido em `.agents`, `.skills`, `.commands`.
- R19: Cumprida. Os quatro greps e `git status --short` estão vazios.
- R20: Cumprida. Greps de redação obsoleta rodados pelo T-2601 sobre as fontes do kit e pelas rodadas de QA 2 e 3 sobre `docs/`, os READMEs e o `CHANGELOG.md`; ocorrências corrigidas (ver Docs).
- R21: Cumprida para a versão 0.32.0 (`internal/build/build.go`, `.claude-plugin/marketplace.json`, `ref`) e o build. A prova sha256 do clone limpo está pendente na hora da escrita; o passo de release a roda antes do push.
- R22: Cumprida. `checkPrivatePaths` passa no `dh validate`.
- R23: Cumprida. `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/reviewer.md`, `.commands/auto.md`, `.agents/coordinator.md`.
- R24: Cumprida. `.commands/review.md`, `.skills/code-review/SKILL.md` (três checagens e queda).
- R25: Cumprida. `.commands/review.md`, `.agents/reviewer.md` (três proibições, leituras do `qa-verifier`).

**Tickets**:

- T-2601 backend, `harness-maintainer`: concluído. Evidência: as rodadas de QA e as três revisões acima.
- T-2602 backend, `builder` (Go): concluído. Evidência: testes da R15 (conferidos por mutação), versão 0.32.0.
- T-2603 backend, `harness-maintainer`: concluído. `REVIEW.md` da raiz com 4 regras, 2 entradas `Ignore` e os dois globs da D7.
- T-2604 teste, `qa-verifier`: concluído. Rodadas 1-5, `PASS` final.

**Docs**: `CHANGELOG.md` (0.32.0); `docs/roadmap.html` e `docs/en/roadmap.html` (card da 0.32.0, 14 templates); `docs/tutorial.html` e `docs/en/tutorial.html`; `docs/tutoriais/templates.{md,html}` e `docs/en/tutorials/templates.{md,html}`; `README.md` e `README.en.md`. Nenhum `pdocs` para este projeto.

**Fora de escopo**: Os itens da seção 5 ficaram de fora: voto majoritário, avaliações de incidentes, job de revisão no CI, agentes classificadores de risco, fila de merge, feature flags, implantação gradual, métricas pós-merge, `dh review run`, um limite de 3 rodadas, um agente, skill ou comando novo, e uma checagem do `dh doctor`. Resíduos: S3 - o `qa-verifier` que lê o texto da PR herda todas as ferramentas e os controles são só prosa (decisão do dono pendente); um rename pode escapar da correspondência de glob de caminho crítico (`Low`, deixado pela D4); a postagem é controlada por autorização em prosa (aceito).
