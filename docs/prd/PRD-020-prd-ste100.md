# PRD-020 · PRD em inglês com regras de redação ASD-STE100, e `language` como idioma de conversa com o dono

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-13 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-13 / 2026-10-13 |
| Tickets | T-2001 |
| Origem | Chat com o dono, 2026-10-13 (o dono delegou: "siga suas recomendações"; a decisão 1 mudou no mesmo dia) |
| Depende de | PRD-007 (fluxo de PRD e tickets), PRD-019 (padrão de versão e `ref`) |
| Habilita | nada |

## 1. Problema

Fatos verificados em 2026-10-13:

- O PRD é o contrato do kit. O `product-discovery` o escreve a partir de `templates/<lang>/PRD.md`.
- O template e a skill `requirements-discovery` não dão nenhum estilo de redação. Frases longas e termos variáveis tornam uma regra ambígua. Dois `builder` então constroem softwares diferentes.
- O `document-review` bloqueia uma ambiguidade só quando as duas leituras dão softwares diferentes. Ele não cita nenhuma regra de redação.
- O ASD-STE100 (Simplified Technical English) é um padrão do inglês. O dicionário oficial de palavras aprovadas tem direitos autorais da ASD.
- Hoje `language` em `.harness/project.yaml` escolhe o idioma do PRD. `.agents/coordinator.md` (seção `Idioma`) manda o Coordenador responder no idioma em que o dono escreve. `.skills/project-onboarding/SKILL.md` (linha 66) define `language` a partir do idioma do dono que o setup detecta.
- O dono disse: escreva o PRD sempre em inglês, com o padrão. A configuração deve lembrar o agente de qual idioma usar entre o orquestrador e o operador.

## 2. Solução

- Todo PRD que o kit escreve é em inglês. Ele segue regras curtas de redação inspiradas no STE100, seja qual for o `language` do projeto.
- Um bloco curto de regras guarda as regras. Ele fica na skill `requirements-discovery`. O template do PRD e a skill `document-review` apontam para ele. Eles não o repetem.
- O `document-validator` confere o estilo. Uma violação que torna uma regra ou um aceite ambíguo bloqueia. Um desvio de estilo que mantém o sentido é Minor.
- O kit não copia o dicionário da ASD. Ele resume as regras com as próprias palavras e cita o padrão pelo nome.
- `templates/pt-br/PRD.md` continua como cópia idêntica byte a byte de `templates/en/PRD.md`. A contagem de templates continua 14. Nenhuma mudança em Go é necessária para o template.
- O `dh validate` e o mod status-drift aceitam `entregue em` e `delivered on`.
- O cabeçalho do PRD, os nomes de suas seções e o seu `Executed summary` usam sempre as formas em inglês: `delivered on <date>`, `Delivered`, `Rules`, `Tickets`, `Docs`, `Out of scope`.
- A skill, o Coordenador e o `product-discovery` deixam de dizer "o idioma do projeto escolhe a forma" para o PRD.
- `language` em `.harness/project.yaml` é o idioma da conversa entre o Coordenador e o dono. O Coordenador o lê no início da execução e responde ao dono nesse idioma. As respostas e os despachos dos especialistas continuam em inglês. O Coordenador traduz.
- O dono troca `language` por pedido explícito. Esse pedido é a autorização: o Coordenador edita `language` em `<harness dir>/project.yaml` e troca, sem segunda confirmação. `dh config set language` só define o padrão da máquina para projetos novos. O Coordenador só o menciona se o dono perguntar por um padrão. Nenhuma chave nova de configuração existe.
- `language` continua escolhendo `templates/<lang>/` para os outros registros (MEMORY, EPOCHAL, RISKS, TASK, ADR, RFC). Os tickets continuam seguindo `language`.

## 3. Regras

- R1 A skill `.skills/requirements-discovery/SKILL.md` tem um bloco de regras de no máximo 20 linhas para a redação do PRD. O bloco declara estas regras:
  - Frases de procedimento têm no máximo 20 palavras. Frases descritivas têm no máximo 25.
  - Uma instrução por frase. Voz ativa.
  - Um termo por conceito. Nenhum sinônimo.
  - Tempos verbais simples. Verbos com um só sentido. Artigos mantidos.
  - Nenhuma expressão idiomática nem phrasal verb quando existe um verbo simples.
  - Listas para sequências. Parágrafos curtos.

  Verificável: `grep -n "20 words" .skills/requirements-discovery/SKILL.md` e `grep -n "25 words" ...` encontram cada um uma linha; o bloco tem no máximo 20 linhas.
- R2 O bloco diz que o ASD-STE100 é um padrão do inglês. Diz que o bloco é um resumo com as palavras do próprio kit. Diz que todo PRD é em inglês, seja qual for o `language`. O bloco não tem nenhuma adaptação para pt-br. Verificável: `grep -n "ASD-STE100" .skills/requirements-discovery/SKILL.md` encontra a linha; `grep -rn -i "estilo STE\|pt-br adaptation" .skills .agents templates` é vazio.
- R3 O kit não copia o dicionário de palavras aprovadas nem texto do padrão. O bloco cita o padrão pelo nome, sem a edição. Verificável: `grep -rn -i "approved word\|approved-word\|dictionary" .skills/requirements-discovery/SKILL.md` encontra só a frase que diz que o dicionário não é copiado; nenhuma lista de palavras aprovadas existe em `.skills/`, `templates/` ou `.agents/`.
- R4 `templates/pt-br/PRD.md` é uma cópia idêntica byte a byte de `templates/en/PRD.md`. `templates/en/PRD.md` tem uma linha de comentário que aponta para o bloco da R1 e não repete as regras. Verificável: `cmp templates/en/PRD.md templates/pt-br/PRD.md` é silencioso; `grep -c "requirements-discovery" templates/en/PRD.md` = 1; `grep -c "20 words" templates/en/PRD.md` = 0.
- R5 `.skills/document-review/SKILL.md` tem um critério que aponta para o bloco da R1 quando o documento é um PRD. Uma violação que torna uma regra ou um aceite ambíguo bloqueia. Um desvio de estilo que mantém o sentido é Minor e não bloqueia. A skill não repete nenhum limite numérico. Verificável: `grep -n "requirements-discovery" .skills/document-review/SKILL.md` encontra o critério; `grep -c "20 words\|25 words" .skills/document-review/SKILL.md` = 0.
- R6 O PRD é sempre em inglês nas fontes.
  - `.skills/requirements-discovery/SKILL.md` diz que o resumo vira o PRD a partir de `templates/en/PRD.md`.
  - `.agents/coordinator.md` (Resumo executado), `.agents/product-discovery.md` (ramo de emenda e contrato compartilhado), `.agents/docs-guide.md` (linha 59) e `.commands/document.md` (linha 18) usam as formas em inglês para o PRD.
  - Esses arquivos deixam de dizer que o idioma do projeto escolhe a forma ou a pasta de templates.
  - `.agents/product-discovery.md` (linha 53) e `.commands/review.md` (linha 19) nomeiam a seção `Solution`, não `Solução`.
  - `AGENTS.md` ("Idioma em execução") diz que o Coordenador responde no idioma registrado em `language`.

  Verificável: `grep -rn "templates/<lang>/PRD" .skills .agents .commands` é vazio; `grep -n "project .language. picks the form" .agents/coordinator.md` é vazio; `grep -n "Solução" .agents/product-discovery.md .commands/review.md` é vazio; o ramo de emenda nomeia só a seção `Docs or Rules`; `grep -n "idioma do dono" AGENTS.md` é vazio.
- R7 `language` é o idioma de conversa com o dono.
  - `.agents/coordinator.md` (seção `Idioma`) diz que o Coordenador lê `language` no início da execução.
  - Diz que o Coordenador responde ao dono nesse idioma e mantém os despachos e as respostas dos especialistas em inglês.
  - Diz que, num pedido explícito do dono para trocar, o Coordenador edita `language` em `<harness dir>/project.yaml` sem segunda confirmação. Esse pedido é a autorização.
  - `.skills/project-onboarding/SKILL.md` (linha 66) e `.commands/setup.md` descrevem `language` como o idioma de conversa com o dono e como a escolha de template para os outros registros.

  Verificável: `grep -n "conversation language" .agents/coordinator.md .skills/project-onboarding/SKILL.md .commands/setup.md` encontra os três; `grep -n "the language the owner writes in" .agents/coordinator.md` é vazio; `grep -n "switch only after agreement" .agents/coordinator.md` é vazio.
- R8 A contagem de templates e a regra de espelho continuam verdadeiras. A contagem continua 14. O `AGENTS.md` mantém `every new template enters in both` e acrescenta que o template do PRD é em inglês nas duas pastas. Verificável: `dh validate` verde; `cmp templates/en/PRD.md templates/pt-br/PRD.md` é silencioso; `ls templates/*/ | wc` não muda.
- R9 Escopo. Os templates TASK, ADR e RFC não mudam. Os PRDs existentes em `docs/prd/` não são reescritos. O código do dashboard não muda: ele já lê `delivered on`. Verificável: `git diff --name-only main...HEAD` não lista nenhum `templates/*/TASK.md`, `ADR.md`, `RFC.md`, nenhum PRD antigo e nenhum arquivo de `internal/dashboard/`.
- R10 As mudanças em Go se limitam a três itens.
  - A linha `defaultVersion` em `internal/build/build.go`.
  - `internal/kit/status_drift.go` aceita `entregue em` e `delivered on`, como `internal/dashboard/metrics.go` (linha 91) já faz. Um teste em Go cobre `delivered on`.
  - O mod `adapters/claude-code/plugin/hooks/status-drift.ts` aceita as duas formas. Um teste em `status-drift.test.ts` cobre isso.

  Motivo: um PRD fechado como `delivered on <date>` falha hoje nas duas checagens de status-drift. O `dh validate` não ganha nenhuma checagem de estilo, e o `go.mod` não muda. Verificável: `git diff main -- '*.go'` mostra só `build.go`, `status_drift.go` e o seu teste; `git diff go.mod` é vazio; `go test ./...` verde; `claude plugin test` verde em `dist/claude-code/dev-harness`.
- R11 Este PRD segue as próprias regras. As linhas de regra que carregam um comando `Verifiable:` ficam isentas do limite de tamanho. Verificável: o validador não encontra nenhuma outra frase acima de 20 palavras (procedimento) ou 25 palavras (descrição).
- R12 Versão 0.27.0. `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da fonte (`v0.27.0`) no mesmo arquivo mudam juntos. Um `ref` esquecido é uma falha de entrega. Verificável: o `grep` dos três valores mostra `0.27.0` ou `v0.27.0` em um só commit; `go run ./cmd/dh build` não mostra nenhuma diferença em `dist/`, exceto a esperada.
- R13 A prova de clone limpo vale, porque o PR toca `internal/` e `dist/*/bin`. Verificável: `git clone . <tmp>`, `go run ./cmd/dh build` ali, e `sha256sum` igual dos binários.
- R14 O `CHANGELOG.md` tem uma entrada 0.27.0. `docs/roadmap.html` e `docs/en/roadmap.html` a refletem. Verificável: `dh validate` verde.

## 4. Docs

Conferidos com `grep` em 2026-10-13. Só entram os arquivos que descrevem a redação do PRD, os templates ou `language`.

- `AGENTS.md`, com três edições:
  - a regra de espelho: o template do PRD é em inglês nas duas pastas;
  - o item "Idioma em execução": o Coordenador responde no idioma de `language`;
  - um item curto para o bloco STE
- `README.md`, `README.en.md` (linha "Templates de trabalho", o item de `language` nas linhas 69-70, a linha de versão 0.27.0)
- `docs/tutorial.html`, `docs/en/tutorial.html` (passo `/dh:discover`; descrição de `language`, linha 823; descrição do `/dh:setup`, linha 740)
- `docs/tutoriais/templates.html`, `docs/en/tutorials/templates.html` e os seus `templates.md` (seção do PRD e lista de templates)
- `docs/roadmap.html`, `docs/en/roadmap.html` (linha 678: reescreva "espelhado por idioma" só se for falso para o PRD; entrada 0.27.0)
- `CHANGELOG.md` (0.27.0)
- `.claude-plugin/marketplace.json` (`version` e `ref`)

## 5. Fora de escopo

- Uma mudança na distribuição dos binários. A ADR-001 continua e os binários continuam commitados.
- Uma checagem em Go ou no `dh validate` para o estilo do PRD.
- Outros templates: TASK, ADR, RFC, story.
- Uma reescrita dos PRDs existentes em `docs/prd/`.
- Uma cópia do dicionário da ASD ou de qualquer texto do padrão.
- Uma chave nova de configuração para o idioma de conversa.

## 6. Decisões do dono

- Dono, 2026-10-13: o PRD é sempre em inglês e segue o padrão. A configuração deve lembrar o agente de qual idioma usar entre o orquestrador e o operador. Isto substitui a primeira decisão (adaptação para pt-br).
- A verificação é só em prosa. O `document-validator` com o `document-review` confere o estilo. Uma ambiguidade bloqueia. Um desvio só de estilo é Minor. Nenhuma checagem nova de estilo existe.
- Decisão do Coordenador, 2026-10-13, sob a delegação do dono: a checagem de status-drift em Go e no mod aceita `delivered on`, porque todo PRD é em inglês. Isto é uma correção de uma checagem existente, não uma checagem de estilo.
- Escopo: só o PRD. Um bloco de regras, em `requirements-discovery`. Direitos autorais: palavras próprias, padrão citado pelo nome, nenhum dicionário.
- Versão 0.27.0 com prova de clone limpo. Os binários continuam commitados.
- Decisão do Coordenador, 2026-10-13: `templates/pt-br/PRD.md` continua como cópia idêntica byte a byte do template em inglês. Isso evita uma mudança em Go no `checkTemplateRefs`.
- Decisão do Coordenador, 2026-10-13: um pedido explícito do dono para trocar `language` é a autorização para editar o `project.yaml`.
- Premissas do analista, não decisões do dono:
  - (a) O bloco fica na skill, porque os comentários do template são apagados quando o PRD fica final.
  - (b) Os limites de tamanho valem para o texto corrido. As linhas `Verifiable:` ficam isentas.
  - (c) O Coordenador não usa `dh config set language` para uma troca de projeto.
  - (d) Os PRDs antigos em pt-br mantêm os cabeçalhos em português.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** todo PRD que o kit escreve é em inglês com o bloco de regras de redação STE, e `language` é o idioma de conversa com o dono.
- **Regras:**
  - R1 Honrada: bloco em `.skills/requirements-discovery/SKILL.md` linhas 100-111, limites de 20/25 palavras na linha 102.
  - R2 Honrada: ASD-STE100 nomeado como padrão do inglês na linha 101; `grep` de "estilo STE" vazio.
  - R3 Honrada: a única menção ao dicionário diz que o dicionário não é copiado.
  - R4 Honrada: `cmp templates/en/PRD.md templates/pt-br/PRD.md` silencioso; uma linha de apontador.
  - R5 Honrada: `document-review` linha 55; nenhum limite numérico repetido.
  - R6 Honrada: os `grep` de `templates/<lang>/PRD` e `Solução` vazios; o `grep` de `idioma do dono` vazio em `AGENTS.md`.
  - R7 Honrada: "conversation language" no coordinator, no onboarding e no setup; `coordinator.md` linha 74.
  - R8 Honrada: `dh validate` verde, 14 templates.
  - R9 Honrada: nenhum template TASK/ADR/RFC, PRD antigo ou arquivo do dashboard no diff.
  - R10 Honrada: o diff em Go é `build.go`, `status_drift.go` e o seu teste; `go test ./...` e `claude plugin test` (254) verdes.
  - R11 Honrada: o QA achou frases longas; o PRD foi emendado e revalidado.
  - R12 Honrada: `defaultVersion`, `version` do marketplace e `ref` são 0.27.0 / v0.27.0.
  - R13 Honrada: o build em clone limpo dá sha256 igual para todos os binários; só `generated_at` difere nos manifestos.
  - R14 Honrada: CHANGELOG 0.27.0 e os dois roadmaps; `dh validate` verde.
- **Tickets:** T-2001, backend, aprovado (matriz de QA; revisão de código `claude` e depois `cli:codex/gpt-5.6-luna` com `approve`; revisão final FEATURE APROVADA).
- **Docs:** `AGENTS.md`, `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/tutoriais/templates.html`, `docs/tutoriais/templates.md`, `docs/en/tutorials/templates.html`, `docs/en/tutorials/templates.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`.
- **Fora de escopo:** distribuição dos binários (a ADR-001 continua), uma checagem de estilo no `dh validate`, outros templates, uma reescrita dos PRDs antigos.
