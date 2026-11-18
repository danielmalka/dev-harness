# PRD-029 · Regras de revisão independente e relatório de deriva

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-28 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-28 / 2026-10-28 |
| Tickets | a definir pelo planejador |
| Origem | Pedido do dono em 2026-10-28; descoberta UND-009 (pesquisa em `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Depende de | PRD-007 (fluxo de tickets e revisão final da feature, entregue), PRD-026 (guia de revisão, entregue), PRD-028 (base de versão 0.33.0, entregue) |
| Entrega | Um pull request, versão 0.34.0 |

## 1. Problema

Uma revisão só é independente na medida do contexto que recebe. Hoje existem cinco lacunas. As cinco estão no fluxo de revisão e verificação.

**A · O pacote de contexto não precisa estar livre da conversa do Coordenador.** Fato: só a dúvida no meio do build retém a afirmação do autor (`.agents/coordinator.md`, "Dúvida no meio do build", linha 143). Fato: o contrato de despacho em `.agents/coordinator.md` (Procedimento, item 6, linha 66) lista "artefatos-fonte" e não proíbe texto de chat. Fato: a fonte do `document-validator` é "notas de descoberta, o pedido do dono ou o registro da tarefa" (`.skills/document-review/SKILL.md`, Entradas, linha 34). Nada impede o Coordenador de colar o próprio resumo do chat. Um revisor que lê esse resumo herda a leitura do autor.

**B · Nenhuma checagem de que o validador entendeu o documento.** Fato: a saída do validador (`.skills/document-review/SKILL.md` linhas 80 a 101, `.agents/document-validator.md` linhas 98 a 116) tem veredito, achados e `Not raised`. Ela não tem uma declaração do que o validador entendeu. Um validador pode aprovar um texto que leu do jeito errado.

**C · Os achados de uma revisão de código não passam por triagem por evidência.** Fato: os achados Critical e Major vão para o builder. Os achados Minor ficam no relatório (`.agents/coordinator.md` item 8, linha 68). Fato: o Coordenador registra `fixed`, `intentional` ou `noise` para cada achado (mesma linha). Fato: `.skills/code-review/SKILL.md` ("Para o autor", passo 3, linha 89) pede que o autor conteste com um arquivo, um teste ou uma restrição. Faltam cinco coisas. Um achado de código não tem id estável (só os achados de documento têm `P1`, `P2`). Nenhuma regra trava a severidade. Nenhuma regra mantém assentado um achado já assentado. Nenhuma regra cobre um achado corrigido que volta. Nada detecta um loop que corrige e quebra o mesmo lugar.

**D · Os loops não têm regra de convergência além de "uma rodada limpa".** Fato: "uma rodada encerra o loop quando não tem achado bloqueante" (`.agents/coordinator.md`, linha 119; `AGENTS.md`, bullet "Verificação padrão", linha 20). Fato: uma etapa em que todo revisor ficou `not-run` nunca é aprovada (`.agents/coordinator.md`, linha 111). Falta: uma regra de que uma rodada com um revisor `not-run` não é limpa. Revisores novos amostram de forma diferente a cada rodada, então uma rodada limpa é um sinal fraco.

**E · Nenhuma checagem de um documento de design contra o código.** Fato: o kit checa só a deriva mecânica (`internal/kit/status_drift.go`, `internal/kit/roadmap_check.go`, `internal/kit/doc_counts.go`). Nada compara um PRD ou um ADR com o código que o implementa. Uma regra como "um produto vendido nunca é apagado" pode quebrar meses depois sem nenhum sinal.

A fonte das ideias é o Elephant-Goldfish Model (EGM) do Google: um plugin de skills em markdown e personas de revisor, Apache-2.0. O EGM não é um modelo de aprendizado de máquina. Fato: o EGM não publica nenhum resultado que mostre melhores desfechos (pesquisa, seção 3). Este PRD toma os mecanismos e os reimplementa nas palavras deste kit.

## 2. Solução

O lote tem cinco partes e um crédito. As quatro primeiras mudam texto em ativos existentes. A quinta acrescenta uma skill. Nenhum agente e nenhum comando é acrescentado.

**Parte A · Independência do pacote de contexto.** O prompt de um revisor tem duas partes. O pacote de contexto é o que o revisor lê sobre o trabalho. A montagem do prompt é o corpo do papel, a skill, a cláusula anti-delegação e o contrato de veredito. Este lote restringe só o pacote de contexto. Para um revisor de documento e para qualquer revisor de CLI, o pacote de contexto tem arquivos e dados citados. Ele não tem chat, nem resumo da intenção do dono, nem raciocínio do Coordenador ou do autor. A dúvida no meio do build mantém a própria regra.

**Parte B · Bloco de compreensão.** O `document-validator` escreve uma seção `## Comprehension` em todo relatório. Ela declara com as palavras do próprio validador o que o documento propõe e como o sistema funciona, conforme o documento diz. Uma seção que o validador não consegue parafrasear em uma leitura é uma ambiguidade bloqueante. O Coordenador compara a paráfrase com o pedido do dono antes de o dono ver o documento.

**Parte C · Triagem dos achados de código.** O Coordenador dá um id a cada achado Critical e Major. Ninguém muda a severidade que o revisor escreveu. O builder responde a cada achado com uma disposição e uma evidência. Um achado assentado não é levantado de novo sem evidência nova. Um achado corrigido que volta é reaberto sob o id antigo. Um achado que existe só por causa de uma correção anterior é oscilação, e a oscilação interrompe o loop.

**Parte D · Convergência.** Uma rodada é limpa só quando nenhum achado está `open`, `reopened` ou `escalated` e todo revisor despachado devolveu um veredito. O loop de build/QA e revisão de código termina após duas rodadas limpas consecutivas, no limite de seis rodadas, ou na oscilação. A expressão do dono, "regra de duas camadas mínimas", significa isto. O dono confirmou a leitura e o escopo em 2026-10-28 (decisão D1).

**Parte E · Relatório de deriva.** Uma skill nova, `drift-report`, compara um documento de design (PRD ou ADR) com o código que ele cobre. Um despacho somente leitura do `repo-scout` encontra o código. Um revisor novo compara o documento seção por seção. O relatório lista só as diferenças, em quatro classes, mais a lógica de código que nenhuma seção menciona. Cada diferença cita o documento e o código. O relatório é um arquivo na pasta da tarefa. Ele nunca corrige nada. O Coordenador o roda a pedido. Depois de `FEATURE APROVADA` em `/dh:review feature PRD-<n>`, o Coordenador o oferece em uma linha opcional.

**Crédito.** O kit credita o EGM na nota de procedência de cada skill que toma um mecanismo dele, e em uma linha de cada README (decisão D3). Nenhum texto do EGM é copiado.

## 3. Regras

Todas as regras são finais. O dono decidiu D1, D2 e D3 em 2026-10-28. As decisões estão registradas no Apêndice.

### 3.1 Parte A · Independência do pacote de contexto

- R1 O prompt de um revisor tem duas partes. O pacote de contexto é o que o revisor lê sobre o trabalho. A montagem do prompt é o corpo do papel, o texto da skill, a cláusula anti-delegação, o contrato de veredito e o bloco de leitura permitida. As regras R1 a R9 restringem só o pacote de contexto. A montagem do prompt continua como `.agents/coordinator.md` e `.skills/external-clis/SKILL.md` a definem. Este lote não a muda.
- R2 A conversa do Coordenador é qualquer turno de chat do dono, qualquer raciocínio do Coordenador, qualquer resumo de intenção e qualquer nota sobre rodadas anteriores. O Coordenador nunca coloca esses itens no pacote de contexto de um revisor da etapa de documento.
- R3 Quando o Coordenador despacha um revisor da etapa de documento, o pacote de contexto tem só cinco itens. Os itens são o caminho do documento, o material-fonte, o caminho do relatório de revisão anterior numa rodada posterior, o idioma do projeto e o harness dir. O material-fonte é uma lista de caminhos de arquivo ou o pedido do dono citado palavra por palavra.
- R4 Quando um revisor de CLI de qualquer etapa (`document`, `code`, `security`) é chamado, o pacote de contexto dele não tem nenhuma parte da conversa do Coordenador.
- R5 O pedido do dono só pode entrar num pacote de contexto em uma de duas formas. A primeira forma é uma citação literal rotulada `owner's words`. A segunda forma é o caminho de um arquivo que contém o pedido. O Coordenador não acrescenta paráfrase.
- R6 Os resultados de checagem do builder e do QA podem ficar num pacote de contexto da etapa de código como dados citados rotulados `unverified`. O `.commands/review.md` (Roteamento, linha 15) já permite isso. R2 não o proíbe.
- R7 O Coordenador lista os elementos de cada pacote de contexto na entrada de Despachos do `MEMORY.md`. Cada elemento é um caminho ou um rótulo. Um leitor posterior pode conferir R3 e R4 por essa lista.
- R8 A regra de R1 a R7 aparece com o mesmo sentido em quatro arquivos.
  - `.agents/coordinator.md` (Procedimento, item 6).
  - `.skills/external-clis/SKILL.md` (linha de Entradas do pacote de contexto, e passo 3).
  - `.skills/document-review/SKILL.md` (Entradas e protocolo de rodadas).
  - `.agents/document-validator.md` (Entradas mínimas).
- R9 A regra da dúvida no meio do build continua como está escrita. O caminho do Claude não ganha regra de parada para um host sem subagentes novos.

Checagem de R1 a R9: uma leitura dos quatro arquivos citados em R8, e uma leitura de uma entrada de Despachos do `MEMORY.md` desde o primeiro uso.

### 3.2 Parte B · Bloco de compreensão

- R10 Quando o `document-validator` revisa um PRD, um plano, uma história ou um ADR, a resposta dele e o relatório persistido têm, cada um, uma seção `## Comprehension`.
- R11 A seção tem três itens. O primeiro é `What the document proposes`, em no máximo 5 frases. O segundo é `How the current system works, as the document states it`, em no máximo 5 frases. O terceiro é `Sections not paraphrased`: uma lista de ids de seção, ou `none`.
- R12 O validador escreve as duas paráfrases com as próprias palavras. Uma frase copiada fica entre aspas.
- R13 Quando o validador não consegue escrever uma leitura única de uma seção, ele registra um achado bloqueante da categoria `ambiguity` para essa seção. Ele lista a seção em `Sections not paraphrased`. Esta regra é uma exceção. A regra de severidade diz que uma ambiguidade só bloqueia quando duas leituras dão softwares diferentes (`.skills/document-review/SKILL.md`, passo 6, linha 59). A skill declara a exceção.
- R14 Tome um termo específico do projeto. Suponha que nenhum documento, nenhuma fonte e nenhum arquivo do repositório o defina. Então o validador registra um achado de ambiguidade pela regra de severidade normal. Um termo geral de engenharia nunca é um achado.
- R15 Quando o veredito é `approved` com zero achados, o relatório ainda tem `## Comprehension`.
- R16 A seção `## Verdict` continua sendo a primeira seção da resposta, com as linhas atuais. O relatório persistido coloca `## Comprehension` depois da tabela de cabeçalho e antes de `## Findings`. Na resposta, ela vem depois de `## Verdict`.
- R17 Antes de o Coordenador enviar um documento validado ao dono, ele lê `## Comprehension`. A mensagem dele ao dono diz que a paráfrase coincide com o pedido, ou cita a frase que difere. Uma diferença vai ao dono como ponto em aberto. O Coordenador não a resolve sozinho.
- R18 O bloco vale só para a etapa de documento. Os modos `code` e `security` do `reviewer` não ganham esse bloco. O bloco não acrescenta rodada e não muda o limite de duas rodadas de correção.

Checagem de R10 a R18: uma leitura de `.skills/document-review/SKILL.md` (Formato de saída, os dois blocos, e passo 6), `.agents/document-validator.md` (Formato de saída) e `.agents/coordinator.md` (handoff de documento).

### 3.3 Parte C · Triagem dos achados de código

- R19 Quando o Coordenador junta os vereditos da etapa de código de uma rodada, ele dá um id da forma `F<round>-<n>` a cada achado Critical ou Major novo. `<round>` é a rodada em que o achado apareceu pela primeira vez. `<n>` conta a partir de 1 entre os achados novos dessa rodada, na ordem do relatório reunido. Um achado levado a uma rodada posterior mantém o id original. Um achado nunca recebe um segundo id.
- R20 O Coordenador registra cada id no `MEMORY.md`. O registro tem a severidade como o revisor a escreveu, o local, uma linha de texto, o estado e a evidência do estado. O Coordenador atualiza o registro de um id quando o estado muda. Ele nunca cria um segundo registro para o mesmo achado.
- R21 Os estados são `open`, `fixed`, `intentional`, `noise`, `escalated` e `reopened`. Os quatro primeiros reutilizam os desfechos do item 8 de `.agents/coordinator.md`, mais `open`. O lote não acrescenta nenhum outro vocabulário para o loop de fim de fatia.
  - `open`: um revisor levantou o achado e nada o fechou. Um achado que o builder reporta como corrigido continua `open` até um revisor confirmar a correção.
  - `fixed`: um revisor confirmou numa rodada posterior que a correção funciona. Isto segue a checagem de status de `.skills/code-review/SKILL.md` ("Reentrada"). O achado está fechado.
  - `intentional` e `noise`: o Coordenador aceitou a evidência de R24, ou um revisor retirou o achado (`noise`). O achado está fechado e assentado.
  - `escalated`: o achado espera o dono (R25 e R29). A resposta do dono o move para `intentional` ou `open`.
  - `reopened`: um achado `fixed` voltou (R27). Ele passa a `fixed` quando um revisor confirma a correção de novo.
- R22 O Coordenador e o builder nunca mudam a severidade de um achado. Um revisor pode mudá-la numa rodada posterior. O revisor cita a evidência nova no achado.
- R23 Para cada achado Critical ou Major que vai ao builder, o builder reporta uma disposição: `fixed`, `intentional`, `noise` ou `escalated`. Uma disposição `fixed` é uma afirmação. O estado continua `open` até um revisor confirmá-la (R21).
- R24 Uma disposição `intentional` ou `noise` cita um arquivo, um teste ou uma restrição como evidência. Sem essa evidência, o achado continua `open`. A linha 75 de `.agents/coordinator.md` continua valendo: um bloqueio não é descartado porque o builder discorda.
- R25 Uma disposição `escalated` significa que o builder não consegue julgar o achado com as informações da árvore. O Coordenador o envia ao dono como prosa numerada. O estado continua `escalated` até o dono responder. Depois da resposta, o estado passa a `intentional` quando o dono aceita o achado como está. Passa a `open` quando o dono pede uma correção.
- R26 Um achado está assentado quando o estado dele é `intentional` ou `noise`. Num despacho de reentrada, o Coordenador lista cada achado assentado com o id, o local e o caminho da evidência. A lista não tem nenhum texto do raciocínio do builder.
- R27 Quando um revisor reentra, ele não levanta de novo um achado assentado no mesmo local e cenário. Ele pode levantar o achado com evidência nova, que ele cita. Quando um achado `fixed` aparece de novo no mesmo escopo, o Coordenador define o estado dele como `reopened` e mantém o id. Um achado `reopened` é bloqueante até um revisor confirmar a correção de novo. Então o estado dele passa a `fixed`.

Seguem as regras de oscilação.

- R28 Quando um revisor reentra, ele acrescenta uma linha `Churn:` a cada achado Critical ou Major novo. O valor é `none`, ou o id ou o local de uma correção anterior, com as palavras `reverting that fix would resolve this finding`. O revisor lê o delta desde o último estado revisado para escrever essa linha.
- R29 Quando um achado de uma rodada tem um valor de `Churn:` diferente de `none`, o Coordenador interrompe o loop e não gasta mais nenhuma rodada. Ele define o estado desse achado como `escalated`. Ele faz ao dono uma pergunta numerada. As opções são: aceitar o achado, mandar corrigir, ou replanejar.
  - Quando o dono aceita o achado, o estado dele passa a `intentional`. O achado está assentado, e o loop continua.
  - Quando o dono manda corrigir, o estado dele passa a `open`. A correção pode reverter a correção anterior. O loop continua, e o limite da etapa continua valendo.
  - Quando o dono escolhe replanejar, o estado continua `open`. O Coordenador reporta `partial`, como R35 diz.

Checagem de R19 a R29: uma leitura de `.agents/coordinator.md` (itens 7 e 8), `.skills/code-review/SKILL.md` (Reentrada, Formato de saída), `.agents/reviewer.md` (Entradas mínimas, modo `code`) e `.commands/review.md` (Roteamento). Os achados de documento mantêm os ids `P<n>` e os estados `applied`, `rejected`, `pending`. A dúvida no meio do build mantém as quatro classes e o limite de 3 ciclos.

### 3.4 Parte D · Convergência

- R30 Uma rodada da etapa de código é limpa quando duas condições valem. Primeira, depois de o Coordenador registrar os estados da rodada, nenhum achado Critical ou Major do loop está `open`, `reopened` ou `escalated`. Segunda, todo revisor que a etapa despachou devolveu um veredito. Um revisor `not-run` torna a rodada não limpa. Para a etapa padrão com um revisor, um revisor `not-run` já torna a etapa `not-run` (`.agents/coordinator.md`, linha 111).
- R31 Quando duas rodadas consecutivas da etapa de código são limpas, o loop de build/QA e revisão de código termina como convergido. Isto vale para o loop de cada ticket e para a revisão final da feature.
- R32 A rodada de confirmação é a segunda rodada limpa. Ela revisa o mesmo estado avaliado que a primeira rodada limpa. O estado avaliado é a árvore do projeto mais as entradas do despacho. Os registros do Coordenador no harness dir não fazem parte dele. Esses registros são o `MEMORY.md`, o registro da tarefa e os estados dos achados. Registrar os estados da primeira rodada limpa não é uma edição. Nenhuma edição do estado avaliado pode acontecer entre as duas rodadas. Quando uma edição acontece, a contagem de rodadas limpas recomeça do zero.
- R33 Uma rodada de confirmação lê a mudança inteira, como uma primeira revisão faz. Ela não lê só um delta. Um delta estaria vazio (`.skills/code-review/SKILL.md`, passo 1, linha 46). A rodada de confirmação ainda conta como reentrada para o rodízio da etapa de código. Ela usa a próxima entrada de `reviewers.code`. Sem outra entrada, o mesmo `reviewer` roda de novo num contexto novo.
- R34 A rodada de confirmação conta no limite de seis rodadas.
- R35 Quando o limite de seis rodadas termina com menos de duas rodadas limpas consecutivas, o Coordenador reporta `partial` ou `blocked`, como `.agents/coordinator.md` já diz. Uma rodada limpa que nenhuma segunda rodada limpa confirmou nunca é reportada como verificada.

As regras R30 a R35 valem só para o loop de build/QA e revisão de código. O dono decidiu que quatro outros loops mantêm a regra atual. A etapa de documento mantém o limite de 2, e uma validação limpa a encerra. A dúvida no meio do build mantém o limite de 3. A etapa `security` e o `/dh:plan-loop` não mudam.

Checagem de R30 a R35: uma leitura de `.agents/coordinator.md` (Revisão proporcional, item 7, rodízio da etapa de código), `.skills/code-review/SKILL.md` (Reentrada) e `AGENTS.md` (bullet "Verificação padrão"). O bullet do `AGENTS.md` substitui as palavras "rodada sem achado bloqueante encerra o loop" para o loop de build/QA e revisão de código. A mesma frase em `.agents/coordinator.md` linha 119 muda do mesmo jeito.

### 3.5 Parte E · A skill `drift-report`

- R36 O lote acrescenta uma skill em `.skills/drift-report/SKILL.md`. O `name` é igual ao nome da pasta `drift-report`. A `description` começa com `Use when`, declara os gatilhos e declara os casos em que a skill não se aplica. O frontmatter tem `author: malka` e `metadata.provenance: adapted`. As seções seguem o contrato de skill de `.skills/harness-authoring/SKILL.md` (tabela na linha 110).
- R37 As entradas são um documento de design (um PRD ou um ADR), nomeado pelo dono ou pelo Coordenador, e a área de código que ele cobre. Quando o documento falta, o despacho para e devolve a entrada faltante ao Coordenador. Um despacho compara um documento.
- R38 Quando a área de código não é nomeada, o Coordenador despacha primeiro o `repo-scout`, somente leitura, para achar o código que o documento cobre. O Coordenador passa só a lista de caminhos para a comparação. Ele não passa a narrativa do `repo-scout`.
- R39 A comparação roda num despacho do `reviewer` no modo `code`, modelo `sonnet`, com a skill `drift-report` nomeada no despacho. O pacote de contexto tem o caminho do documento, a lista de caminhos de código e o harness dir (R1 a R4 valem). O `.agents/reviewer.md` ganha uma frase: quando o despacho nomeia `drift-report`, o revisor segue a saída dessa skill no lugar da saída do modo `code`. Nenhum agente é acrescentado.
- R40 O revisor divide o documento pelo título de seção. Para cada seção, ele lista as afirmações (regras, fluxos, dados, decisões), acha o código que sustenta cada afirmação e as compara. Ele reporta só as diferenças.
- R41 Uma diferença pertence a uma classe.
  - `missing feature`: o documento a declara e a área de código não tem implementação.
  - `incomplete model`: o código implementa parte dela (faltam campos, estados ou casos).
  - `logic shortcut`: o código chega ao resultado por um caminho que pula uma regra declarada.
  - `architectural divergence`: a estrutura do código difere de uma decisão declarada (camada, fronteira ou tecnologia).
- R42 O relatório também lista `unspecified logic`: comportamento na área de código que nenhuma seção do documento menciona.
- R43 Cada diferença tem cinco campos.
  - Um id `D<n>`.
  - Sua classe.
  - Uma citação do documento: caminho relativo, e seção ou linha.
  - Uma citação do código: caminho relativo e linha.
  - Uma declaração da diferença em no máximo 2 frases.

  Para `missing feature`, a citação do código é `none found` com os caminhos buscados. Para `unspecified logic`, a citação do documento é `none`.
- R44 O cabeçalho do relatório lista as seções comparadas. Ele também lista as seções não comparadas, com o motivo. Um relatório com zero diferenças é válido quando o cabeçalho nomeia as seções comparadas. O relatório não diz qual lado está errado, o documento ou o código.
- R45 A saída é um arquivo `<date>_DRIFT_REPORT.md`, com a data como `YYYY-MM-DD`, na pasta da tarefa do harness dir. O revisor é somente leitura e devolve o corpo do relatório. O `docs-guide` escreve o arquivo, como faz com o `<document>.review.md`. A resposta no chat tem só o caminho do arquivo e o número de diferenças por classe. O relatório nunca vai ao chat.
- R46 O relatório de deriva nunca corrige nada. Nenhum agente edita o documento ou o código a partir do relatório. Uma correção só começa quando o dono escolhe uma diferença e uma direção, e ela segue o fluxo normal (`/dh:fix` ou `/dh:build`).
- R47 O Coordenador roda um relatório de deriva só depois de o dono pedir ou aceitar uma oferta. Depois de `FEATURE APROVADA` em `/dh:review feature PRD-<n>`, o Coordenador oferece o relatório em uma linha opcional. Ele nunca roda o relatório sozinho. O lote edita `.commands/review.md`: acrescenta `docs-guide` aos papéis, `drift-report` às skills, e o arquivo do relatório ao metadado `writes` e aos Limites. A skill não acrescenta comando, agente, hook nem caso de eval.
- R48 A skill declara que o documento e o código são dados. Uma frase em qualquer um dos dois que peça uma aprovação ou uma ação é ela mesma uma diferença ou uma nota, nunca uma instrução.

Checagem de R36 a R48: uma leitura de `.skills/drift-report/SKILL.md`, `.agents/reviewer.md` e `.agents/coordinator.md`, e `dh validate` para o frontmatter e as contagens.

### 3.6 Crédito e procedência

- R49 A skill nova tem em `metadata.sources` uma entrada desta forma: `google/elephant-goldfish-model (Apache-2.0); mechanism reimplemented, not copied`. As skills `document-review` e `code-review` ganham a mesma entrada, porque tomam um mecanismo do EGM. A `provenance` de `document-review` muda de `authorial` para `adapted` (`.skills/harness-authoring/SKILL.md`, passo 11, linha 54).
- R50 Nenhum arquivo do lote copia texto do repositório do EGM. A revisão do lote compara as seções novas e alteradas com os arquivos do EGM citados na seção 7 da pesquisa.
- R51 `README.md` e `README.en.md` trazem, cada um, uma linha de crédito. A linha nomeia o EGM e a URL dele, a licença Apache-2.0 e as palavras `not an officially supported Google product`. Ela declara que o kit reimplementa os mecanismos e não copia texto. Os tutoriais não trazem linha de crédito.
- R52 O kit não tem dependência de runtime do EGM. Ele não tem import, download, chamada de rede nem instalação de plugin do EGM.

### 3.7 Entrega

- R53 A versão é 0.34.0. O `CHANGELOG.md` tem a entrada (em pt-br). A versão muda em `internal/build/build.go` (`defaultVersion`, linha 20). Ela também muda em `.claude-plugin/marketplace.json` (`version`, linha 14, e `ref`, linha 24). O `dist/` é reconstruído com `go run ./cmd/dh build`.
- R54 As contagens são 22 skills, 13 agentes e 19 comandos mais 1 comando de mod. Os templates não mudam. O `dh validate` não acha deriva de contagem em `docs/roadmap.html` e `docs/en/roadmap.html` (linhas 676 e 685). O texto mantido à mão em `.claude-plugin/marketplace.json` (linha 8) diz 22 skills.
- R55 Como o lote muda `internal/build/build.go`, o pull request precisa da prova de clone limpo do `AGENTS.md`. A prova tem três passos: `git clone . <tmp>`, depois `go run ./cmd/dh build` lá, depois uma comparação do sha256 dos binários com o da árvore. O Coordenador só faz o commit depois da prova.
- R56 O lote não edita nenhum arquivo chamado `references/clone-isolation.md` ou `references/live-tree-freeze.md`. Ele não edita nenhuma linha da cláusula anti-delegação. Ele não muda nenhuma parte da trava do PRD-027. O lote não precisa de etapa `security`. Ele precisa da revisão de código padrão e de uma matriz do `qa-verifier`.
- R57 O `AGENTS.md` recebe só as regras atuais, em pt-br. São a regra do pacote de contexto (R1 a R5) e a frase substituída da regra de convergência (R30 a R35). Um bullet cobre o relatório de deriva (R46, R47). O histórico fica neste PRD.
- R58 A mudança atualiza os documentos da seção 4 no mesmo pull request. Não existe pull request só de documentação.
- R59 Nenhum caso de eval é acrescentado e nenhum eval pago roda. No primeiro uso, o Coordenador registra no `MEMORY.md` o custo de uma rodada de confirmação e de um relatório de deriva. O registro é uma nota, não uma medição.
- R60 Todo arquivo novo ou alterado usa caminhos relativos ou `~/...`. O `dh validate` não acha nenhum caminho privado absoluto.

## 4. Docs

- README.md
- README.en.md
- docs/tutorial.html
- docs/en/tutorial.html
- docs/roadmap.html
- docs/en/roadmap.html
- AGENTS.md
- CHANGELOG.md
- .claude-plugin/marketplace.json
- .agents/coordinator.md
- .agents/document-validator.md
- .agents/reviewer.md
- .skills/document-review/SKILL.md
- .skills/code-review/SKILL.md
- .skills/external-clis/SKILL.md
- .skills/drift-report/SKILL.md
- .commands/review.md

Notas sobre esta lista.

- Os roadmaps ganham o card da 0.34.0, os chips (data e status `0.34.0`) e a contagem 22 no bloco de estrutura e na linha "Contagem atual".
- Os tutoriais descrevem a regra do pacote de contexto, o bloco `## Comprehension`, os ids e estados dos achados, as duas rodadas limpas e como pedir um relatório de deriva.
- Os README ganham a linha de crédito de R51. Fato: nenhum dos dois README declara uma contagem de skills (conferido por `grep`).
- O `.commands/review.md` muda pela oferta de R47. O `.agents/reviewer.md` muda pela frase de R39.
- As notas de procedência de R49 ficam em `.skills/drift-report/SKILL.md`, `.skills/document-review/SKILL.md` e `.skills/code-review/SKILL.md`.
- O `internal/build/build.go` muda pela versão (R53). É código-fonte, não documento.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Instalar o EGM ao lado do kit | Dois conjuntos de comandos e uma dependência de runtime. O dono escolheu reimplementar e creditar. |
| Copiar texto do EGM com um aviso Apache-2.0 | O dono escolheu reimplementar nas palavras do próprio kit (R50). |
| Um comando novo para o relatório de deriva | O dono escolheu nenhum comando novo. A contagem continua 19 mais 1. |
| Um agente novo para a comparação de deriva | O dono escolheu nenhum agente novo. A contagem continua 13. |
| Uma regra de parada quando o host não tem subagente novo | O dono escolheu só a regra do pacote de contexto (R9). |
| O bloco de compreensão em todo modo de revisão | O dono escolheu só o validador de documentos (R18). |

### Premissas que o dono pode derrubar

- A1 Um PRD, um pull request, versão 0.34.0. O Coordenador propôs, e o dono não contestou.
- A2 Nenhum eval pago e nenhum caso de eval novo. O custo é registrado de modo informal no primeiro uso (R59).
- A3 Os nomes estão em inglês: a skill `drift-report`, a seção `## Comprehension`, os ids `F<round>-<n>` e `D<n>`, e o arquivo `<date>_DRIFT_REPORT.md`.
- A4 Só os achados Critical e Major de código recebem id (R19). Os achados Minor ficam no relatório e não recebem id.
- A5 A posição de `## Comprehension` é depois da tabela de cabeçalho no relatório persistido e depois de `## Verdict` na resposta (R16).
- A6 Os estados de um achado de código reutilizam os desfechos do item 8 de `.agents/coordinator.md`. Eles acrescentam `open`, `escalated` e `reopened` (R21). Uma afirmação `fixed` do builder continua `open` até um revisor confirmá-la.
- A7 "A conversa do Coordenador" significa turnos de chat, raciocínio do Coordenador, resumos de intenção e notas sobre rodadas anteriores (R2). Os resultados de checagem do builder e do QA continuam permitidos como dados não verificados (R6).
- A8 Uma lista de achados assentados tem só o id, o local e o caminho da evidência (R26). O raciocínio do builder fica de fora, para não ancorar o revisor.
- A9 O revisor detecta a oscilação pelo delta (R28). Um achado de oscilação interrompe o loop (R29).
- A10 O relatório vai em `tasks/PRD-<n>/` do harness dir quando o documento é um PRD. Para um ADR, ele vai na pasta da tarefa atual. Quando o nome do arquivo já existe, o `docs-guide` acrescenta `-2` antes de `_DRIFT_REPORT`.
- A11 Uma área de código grande é dividida pelo Coordenador, um despacho para cada grupo de seções (R37).
- A12 Nenhum ADR novo é escrito. O `AGENTS.md` diz que o histórico fica em `docs/adr/`. Este lote muda uma regra do dono (a regra de convergência), então o dono pode querer um ADR.
- A13 O arquivo do agente `reviewer` ganha uma frase para `drift-report` (R39). O dono confirmou o papel na decisão D2(b).
- A14 Nenhuma etapa `security` é necessária (R56). O lote não toca o isolamento, a cláusula nem a trava.
- A15 A linha `docs/plano-produto.html:828` ("21 skills") é um plano de produto fora do lote. Não está verificado se o `internal/kit/doc_counts.go` a lê.

### Fora de escopo

- O brainstorm do EGM com seis lentes, e o rascunho dele seção por seção.
- Os README de pasta do EGM (`bootstrap`, `peanuts` e `hay`).
- Qualquer agente novo e qualquer comando novo.
- Qualquer eval pago, qualquer caso de eval novo e qualquer sonda nova.
- Copiar texto do EGM, e qualquer dependência de runtime do EGM.
- Uma regra de parada para hosts sem subagentes novos.
- Os hooks e os arquivos de plugin do EGM para o Claude Code.
- Uma mudança no isolamento dos revisores de CLI, na trava do PRD-027 ou na cláusula anti-delegação.

### Riscos e decisões pendentes

- Risco: o custo em tokens de uma rodada de confirmação e de um relatório de deriva não está medido · mitigação: R59 registra o primeiro uso.
- Risco: a rodada de confirmação acrescenta pelo menos uma revisão a todo loop que passa na primeira rodada · mitigação: o dono aceitou esse custo na decisão D1(b) e escolheu todo loop. R59 registra o custo, e o dono pode restringir o escopo depois.
- Risco: o EGM não publica evidência de que esses mecanismos melhoram os desfechos (pesquisa, seção 3) · mitigação: R59 e os primeiros usos dão ao dono um registro local.
- Risco: o revisor e o Coordenador usam o mesmo modelo base, então um contexto novo remove só o viés de contexto · mitigação: R33 faz o rodízio para a próxima entrada de `reviewers.code` quando a lista tem uma.
- Risco: um relatório de deriva sobre uma leitura parcial dá diferenças falsas · mitigação: R43 exige duas citações, e R44 lista as seções não comparadas.
- Risco: a atribuição Apache-2.0 se aplica se texto for copiado · mitigação: R50 proíbe copiar, e R49 e R51 creditam a fonte.
- Decisão pendente: nenhuma. O dono decidiu D1, D2 e D3 em 2026-10-28. O registro segue.

#### D1 · Significado e escopo de "duas camadas mínimas" (decidido)

Contexto: o dono escreveu "Prefiro adicionar a regra de 2 camadas mínimas". O Coordenador leu como a regra do EGM "duas rodadas limpas consecutivas". Antes deste lote, uma rodada sem achado bloqueante encerrava um loop (`AGENTS.md` linha 20, `.agents/coordinator.md` linha 119). A regra nova custa pelo menos mais uma revisão para cada loop que passa na primeira rodada.

Parte (a) · O significado.

1. (Escolhida) Duas rodadas limpas consecutivas, com uma rodada de confirmação sobre o mesmo estado avaliado (R30 a R35). Consequência: uma contagem clara, pelo menos duas revisões para cada loop. Com `reviewers.code` que lista duas entradas, o rodízio faz a segunda revisão vir de outro revisor.
2. Dois revisores diferentes precisam aprovar, em qualquer ordem. Consequência: precisa de duas entradas em `reviewers.code`. Com uma entrada, a regra não consegue valer ou recai numa repetição do mesmo revisor.
3. Duas camadas significa as checagens determinísticas mais uma revisão independente. Consequência: isto já é verdade hoje, então a regra só o repetiria. Nenhum custo extra, nenhuma segurança nova.

Decidido pelo dono em 2026-10-28: opção 1.

Parte (b) · Os loops que a regra cobre.

- Loop de build/QA e revisão de código (limite 6).
  1. (Escolhida) Todo loop, cada ticket e a revisão final da feature. Consequência: o maior custo, o sinal mais forte.
  2. Só a revisão final da feature. Consequência: os loops de ticket mantêm uma rodada limpa, e o custo continua baixo.
  3. Só um loop que já achou um achado bloqueante. Consequência: uma mudança que passa na primeira rodada não paga nada, mas uma primeira rodada de sorte não é confirmada.
  Decidido pelo dono em 2026-10-28: opção 1.
- Validação de documento (limite 2). Duas rodadas limpas não cabem em duas rodadas de correção.
  1. Manter o limite de 2 e exigir que a última validação seja limpa. Consequência: na prática igual a hoje, com uma frase nova.
  2. Subir o limite para 3 e exigir duas rodadas limpas consecutivas. Consequência: mais uma validação Opus para a maioria dos documentos.
  3. (Escolhida) Deixar os documentos com a regra antiga. Consequência: nenhuma mudança e nenhum custo extra.
  Decidido pelo dono em 2026-10-28: opção 3.
- Dúvida no meio do build (limite 3). Decidido pelo dono em 2026-10-28: inalterada. É uma checagem única no meio do build, e uma segunda dúvida limpa dobraria o custo dela.
- `/dh:plan-loop`. Decidido pelo dono em 2026-10-28: inalterado. Cada onda já é julgada pela lista completa de `reviewers.document`.
- A etapa `security`. Decidido pelo dono em 2026-10-28: inalterada. Ela já roda a lista completa e não faz rodízio.

Parte (c) · As palavras que substituem a frase do `AGENTS.md` "Build/QA e revisão de código: teto de 6 rodadas, e rodada sem achado bloqueante encerra o loop".

1. (Escolhida) Substituí-la por este sentido: "Build/QA e revisão de código: limite de 6 rodadas. Uma rodada é limpa quando nenhum achado bloqueante está aberto, reaberto ou escalado e todo revisor despachado devolveu um veredito. Duas rodadas limpas consecutivas encerram o loop. Um achado de oscilação ou o limite o encerra antes, e o Coordenador reporta ao dono." A frase sobre documentos continua inalterada, como a parte (b) decide. Consequência: uma regra para o loop de código, em um só lugar.
2. Manter a frase antiga para a revisão padrão e usar a nova só na revisão final da feature. Consequência: duas regras no `AGENTS.md`, e um leitor precisa saber qual vale.
3. O dono dá outras palavras. Consequência: R30 a R35 e a linha 119 de `.agents/coordinator.md` as seguem.

Decidido pelo dono em 2026-10-28: opção 1. O texto em pt-br vai no `AGENTS.md` na hora do build, com o mesmo sentido.

#### D2 · Ponto de entrada do relatório de deriva e quem roda a comparação (decidido)

Contexto: o dono escolheu uma skill sem comando novo e sem agente novo. Os pontos em aberto eram se um comando oferece o relatório, e qual papel compara.

Parte (a) · O ponto de entrada.

1. (Escolhida) Uma linha opcional que o Coordenador oferece depois de `FEATURE APROVADA` em `/dh:review feature PRD-<n>`. Ela nunca roda sozinha. Consequência: PRD e código estão ambos estáveis nesse ponto. O `.commands/review.md` muda: papéis, skills, metadado `writes` e os Limites (hoje `writes: none` e somente leitura na linha 35).
2. Um assunto do `/dh:document`. Consequência: o metadado `writes` dele já permite documentos, mas o comando escreve documentação e este relatório é uma revisão.
3. Um modo do `/dh:verify`. Consequência: o `qa-verifier` escreve testes, e uma comparação de documento e código é uma revisão. É a que se encaixa pior.
4. Só a pedido, sem texto de comando. Consequência: nenhum arquivo de comando muda e o dono precisa lembrar de pedir. Os tutoriais descrevem como.

Decidido pelo dono em 2026-10-28: opção 1. Nota: a revisão final já lê a aderência a partir da evidência do QA, e o relatório de deriva lê o código. As duas checagens diferem.

Parte (b) · Quem compara.

1. (Escolhida) `reviewer` no modo `code`, com a skill `drift-report` nomeada no despacho (R39). Consequência: Sonnet, somente leitura, um contexto novo. Uma frase é acrescentada a `.agents/reviewer.md`.
2. `document-validator` com a skill. Consequência: Opus, somente leitura, e a missão dele já é "documento contra a fonte". Custa mais, e a saída dele é um veredito.
3. Um terceiro modo `drift` do `reviewer`. Consequência: uma edição maior de `.agents/reviewer.md` (description, procedimento, saída) e nenhum agente novo.

Decidido pelo dono em 2026-10-28: opção 1.

#### D3 · Onde vai a linha de crédito (decidido)

Contexto: R49 já põe o crédito na nota de procedência de cada skill que toma um mecanismo, como o kit faz com outras fontes. O dono pediu o crédito em uma linha dos docs. Nenhum dos dois README tem hoje uma seção de créditos.

1. Só as notas de procedência. Consequência: nenhuma linha de documento, então o pedido do dono não é atendido.
2. (Escolhida) As notas de procedência e uma linha em `README.md` e `README.en.md`. Consequência: o crédito fica visível no ponto de entrada, nos dois idiomas.
3. As notas de procedência e uma linha nos tutoriais (`docs/tutorial.html`, `docs/en/tutorial.html`). Consequência: o crédito fica onde o fluxo novo é descrito, e um leitor do README não o vê.

Decidido pelo dono em 2026-10-28: opção 2.

## Resumo executado

- **Entregue:** Os revisores leem um pacote de contexto fechado em vez da conversa do Coordenador, o validador de documentos declara a própria compreensão, os achados da revisão de código passam por triagem por id e estado com detecção de oscilação e uma regra de convergência de duas rodadas limpas consecutivas, e uma skill nova somente leitura, `drift-report`, compara um documento com o código, creditada ao EGM e lançada como 0.34.0 com 22 skills.
- **Regras:**

  | Regra | Veredito | Evidência |
  |---|---|---|
  | R1 | Cumprida | T-2901, T-2904, T-2908; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R2 | Cumprida | T-2901, T-2904, T-2908; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R3 | Cumprida | T-2901, T-2904, T-2908; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R4 | Cumprida | T-2901, T-2904, T-2908; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R5 | Cumprida | T-2901, T-2904, T-2908; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R6 | Cumprida | T-2901, T-2904, T-2907; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R7 | Cumprida | T-2904; matriz do T-2911 PASS (texto em `coordinator.md`); primeiro uso: as entradas de Despachos do AUTO-2 no registro do Coordenador listam os elementos do pacote (ver nota g); reconferido no QA final (`qa-final`) |
  | R8 | Cumprida | T-2901, T-2904; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R9 | Cumprida | T-2904; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R10 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R11 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R12 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R13 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R14 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R15 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R16 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R17 | Cumprida | T-2904; matriz do T-2911 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R18 | Cumprida | T-2901; matriz do T-2911 PASS (leitura e grep) |
  | R19 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R20 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R21 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R22 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R23 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R24 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R25 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R26 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R27 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R28 | Cumprida | T-2902, T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R29 | Cumprida | T-2905; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`) |
  | R30 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a |
  | R31 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a |
  | R32 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a |
  | R33 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a; `reviewer.md:81` foi corrigido depois da matriz e reconfirmado (nota b) |
  | R34 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a |
  | R35 | Cumprida | T-2906, T-2908, T-2902; matriz do T-2912 PASS (leitura e grep); reconferido no QA final (`qa-final`); ver nota a |
  | R36 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R37 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R38 | Cumprida | T-2903; matriz do T-2913 PASS (leitura); reconferido no QA final (`qa-final`) |
  | R39 | Cumprida | T-2902, T-2903, T-2906; matriz do T-2913 PASS (uma frase em `reviewer.md`; 13 arquivos de agente); reconferido no QA final (`qa-final`) |
  | R40 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R41 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R42 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R43 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R44 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R45 | Cumprida | T-2903, T-2906, T-2907, T-2908; matriz do T-2913 PASS (leitura); reconferido no QA final (`qa-final`) |
  | R46 | Cumprida | T-2903, T-2906, T-2907, T-2908; matriz do T-2913 PASS (leitura); reconferido no QA final (`qa-final`) |
  | R47 | Cumprida | T-2903, T-2906, T-2907, T-2908; matriz do T-2913 PASS (leitura); reconferido no QA final (`qa-final`); a primeira execução do T-2913 falhou (ver nota b), reconfirmado |
  | R48 | Cumprida | T-2903; matriz do T-2913 PASS (leitura) |
  | R49 | Cumprida | T-2901, T-2902, T-2903; o frontmatter das três skills traz a string da fonte; matrizes PASS, T-2914 PASS |
  | R50 | Cumprida | T-2901, T-2902, T-2903; comparação de sequências de palavras com os arquivos públicos do EGM (ver nota c) |
  | R51 | Cumprida | T-2909; leitura dos dois README PASS; os tutoriais não trazem linha de crédito |
  | R52 | Cumprida | T-2903; `grep -rn` do nome do repositório do EGM: só entradas de `sources`, os README, o CHANGELOG (e o `dist/` gerado); os dois PASS |
  | R53 | Cumprida | T-2909, T-2910, T-2903; lugares da versão, topo do `CHANGELOG.md`, contagens 22/13/19+1; `dh validate` 0 erros; T-2914 PASS |
  | R54 | Cumprida | T-2909, T-2910, T-2903; lugares da versão, topo do `CHANGELOG.md`, contagens 22/13/19+1; `dh validate` 0 erros; T-2914 PASS |
  | R55 | Cumprida | processo (nenhum ticket de build); prova de clone limpo (ver nota d) |
  | R56 | Cumprida | restrição (nenhum ticket de build); o diff não toca as referências de isolamento, a cláusula, `hooks/` nem `adapters/`; os dois PASS (ver nota e) |
  | R57 | Cumprida | T-2908; leitura do diff do `AGENTS.md` PASS; reconferido no QA final (`qa-final`) |
  | R58 | Cumprida | T-2909; lista de Docs conferida com `git diff main..HEAD --name-only`, PASS |
  | R59 | Cumprida | T-2906; diff de `evals/` vazio, PASS; nota de custo registrada (ver nota f); reconferido no QA final (`qa-final`) |
  | R60 | Cumprida | T-2901, T-2902, T-2903; varredura de caminho privado e `dh validate` 0 erros; todos PASS |

  Notas sobre a tabela. Os caminhos de evidência ficam em `~/.harness/projects/dev-harness/tasks/PRD-029/evidence/`, salvo indicação contrária.
  - (a) R30 a R35 estão Cumpridas no texto. Os loops deste lote rodaram sob a regra anterior (uma rodada sem achado bloqueante encerra o loop), porque a regra nova vale a partir do merge. A primeira aplicação da regra de duas rodadas limpas vem depois do merge.
  - (b) Duas checagens falharam uma vez e foram corrigidas. T-2913 R47: o bullet de `AGENTS.md:32` não tinha o nome `<date>_DRIFT_REPORT.md`; o Coordenador corrigiu uma linha e a regra foi reconfirmada. T-2912 R33 (recomendado pelo QA): `reviewer.md:81` declara que uma rodada de confirmação lê a mudança inteira; o Coordenador corrigiu uma linha e a reconfirmou. O QA final (`qa-final`) reconferiu as duas.
  - (c) R50 está Cumprida. A comparação de sequências de palavras usou os arquivos públicos do EGM, obtidos em somente leitura no commit f0ba151 (`egm/MANIFEST.md`). A maior sequência em comum tem 6 palavras e o limite é 8. Evidência: `qa-T-2911/r50.txt`, `qa-T-2912/r50.txt`, `qa-T-2913/r50-output.txt`. As sequências de 6 palavras são `reviewer raises a settled finding again`, `system works as the document` (um nome de item que o PRD exige) e `reverting that fix would resolve`.
  - (d) R55 está Cumprida. O T-2914 rodou um clone limpo de `3bb762c`, rodou `go run ./cmd/dh build` lá e comparou o sha256: os cinco binários são idênticos na árvore e no clone. O resultado está registrado no registro do Coordenador (`MEMORY.md`, entrada "AUTO-2 COMMIT CANDIDATO e T-2914"); o T-2914 não salvou nenhum arquivo de saída. A prova é repetida no commit final antes da publicação remota, com um log salvo.
  - (e) R56 está Cumprida. Uma passagem de `/dh:secure` rodou no T-2901 como precaução de superfície de execução decidida pelo Coordenador (lista completa de `reviewers.security`; nenhum achado; evidência `security-T-2901/`). As referências de isolamento, a cláusula, `hooks/` e `adapters/` não foram tocados.
  - (f) R59 está Cumprida. Nenhum eval pago rodou e `evals/` não foi tocado. A nota de custo do primeiro uso (informal, não medida) está registrada no registro do Coordenador depois da revisão final.
  - (g) R7 é evidência de processo no primeiro uso: as entradas do AUTO-2 no registro do Coordenador nomeiam os elementos do pacote de cada despacho.
  - (h) Limites, ditos às claras: as matrizes de QA do T-2911 ao T-2913 são asserções de leitura e grep, não testes executáveis. A segunda revisão de vários tickets leu só o delta (dito em cada registro). O painel da revisão final (`review-final-codex`, `review-final-opencode`) aprovou sem nenhum achado bloqueante.

- **Tickets:** Todos os tickets estão `done` em 2026-10-28. Os caminhos de evidência ficam em `~/.harness/projects/dev-harness/tasks/PRD-029/evidence/`.
  - T-2901, lane `backend` (etapa de documento, pacote, `## Comprehension`), `done`. Evidência: `doubt-T-2901/` (approve), `review-T-2901/` (approve), `security-T-2901/` (nenhum achado), QA `qa-T-2911/`.
  - T-2902, lane `backend` (etapa de código, reentrada, achados assentados, `Churn:`), `done`. Evidência: `review-T-2902/` (um Major, corrigido), `review-T-2902-r2/` (approve), QA `qa-T-2912/`.
  - T-2903, lane `backend` (skill `drift-report`, contagens do roadmap), `done`. Evidência: `review-T-2903/` (um Major, corrigido), `review-T-2903-r2/` (approve), QA `qa-T-2913/`.
  - T-2904, lane `backend` (pacote do Coordenador e leitura de Comprehension), `done`. Evidência: `review-batch-C/` (nenhum achado), QA `qa-T-2911/`.
  - T-2905, lane `backend` (triagem do Coordenador), `done`. Evidência: `review-batch-C/` (nenhum achado), QA `qa-T-2912/`.
  - T-2906, lane `backend` (convergência, execução do drift-report, nota de custo), `done`. Evidência: `review-batch-C/` (um Major, corrigido), `review-T-2906-r2/` (approve), QA `qa-T-2912/`.
  - T-2907, lane `backend` (oferta do `/dh:review` e convergência), `done`. Evidência: `review-batch-C/` (nenhum achado), QA `qa-T-2913/`.
  - T-2908, lane `backend` (regras do `AGENTS.md`), `done`. Evidência: `review-batch-C/` (nenhum achado), QA `qa-T-2911/`, `qa-T-2912/`, `qa-T-2913/` (R47 corrigido, nota b).
  - T-2909, lane `backend` (docs do lote), `done`. Evidência: `review-batch-D/` (uma linha `request changes` com 0 achados, causada por resultados de checagem ausentes), `review-batch-D-r2/` (approve), T-2914.
  - T-2910, lane `infra` (versão 0.34.0, texto do marketplace, rebuild), `done`. Evidência: `review-batch-D-r2/` (approve), T-2914 (`go test ./...` 7 pacotes ok, `dh validate` 0 erros).
  - T-2911, lane `teste (unitário)` (R1 a R18, R49, R50, R60), `done`. Evidência: `qa-T-2911/` (21 PASS e 4 checagens extras PASS).
  - T-2912, lane `teste (unitário)` (R19 a R35, R49, R50, R60, escopo negativo), `done`. Evidência: `qa-T-2912/` (todos PASS, 5 linhas de escopo negativo PASS).
  - T-2913, lane `teste (unitário)` (R36 a R48, R50, R52, R56), `done`. Evidência: `qa-T-2913/` (R47 FAIL corrigido e reconfirmado, o resto PASS).
  - T-2914, lane `teste (integração)` (gate final, clone limpo, R49, R51 a R60), `done`. Evidência: o registro do Coordenador (gate verde, 5 binários idênticos, nota d); `qa-final/` (0 falhas, `validate` 0 erros com 13 agentes, 19 comandos, 22 skills).

- **Docs:** `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html` (uma frase em cada uma corrigida depois da revisão final), `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md`, `.claude-plugin/marketplace.json`, `.agents/coordinator.md`, `.agents/document-validator.md`, `.agents/reviewer.md`, `.skills/document-review/SKILL.md`, `.skills/code-review/SKILL.md`, `.skills/external-clis/SKILL.md`, `.skills/drift-report/SKILL.md` (nova), `.commands/review.md`, `docs/prd/PRD-024-public-readiness.md` (redação de contagem em uma linha, fora da lista de Docs da seção 4: limpeza de um texto que esta mudança tornou obsoleto), este PRD e o registro de revisão dele `PRD-029-independent-review-and-drift-report.review.md`. Não são documentos: `internal/build/build.go` (constante de versão), `harness-manifest.json` e `dist/` (gerados).

- **Fora de escopo:** O brainstorm do EGM com seis lentes e o rascunho dele seção por seção; os README de pasta do EGM (`bootstrap`, `peanuts` e `hay`); qualquer eval pago, caso de eval novo ou sonda nova; copiar texto do EGM e qualquer dependência de runtime do EGM; uma regra de parada para hosts sem subagentes novos; os hooks e os arquivos de plugin do EGM para o host dele; uma mudança no isolamento dos revisores de CLI, na trava do PRD-027 ou na cláusula anti-delegação. Também fora de escopo: o bloco de compreensão nos revisores de código (só o validador de documentos); uma mudança no sistema de memória; nenhum agente, comando, hook, template ou caso de eval novo.
