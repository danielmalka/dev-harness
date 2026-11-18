# ADR-011 · O kit é só pt-BR: fim do seletor de idioma e do espelho em inglês

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-11-17 |
| Decisor | Daniel Malka (dono) |
| Origem | Decisão do dono de 2026-11-17 e plano "dev-harness só pt-BR" v1.2 (quatro perguntas respondidas na mesma data) |
| Reversibilidade | cara: o kit inteiro muda de idioma em cinco versões (0.42.0 a 0.46.0); voltar ao regime bilíngue exige traduzir tudo de novo |
| Relação | substitui as regras de idioma que o `AGENTS.md` trazia desde a 0.1.0 ("Idioma dos arquivos" e "Idioma em execução"); não muda nenhuma regra de revisão, de memória ou de trava |

## 1. Contexto

- Fato: até a 0.41.0 o kit mantinha duas línguas em espelho. `.skills/`, `.agents/`, `.commands/` e `profiles/` eram em inglês; `templates/` tinha `en/` e `pt-br/`; `docs/` tinha o espelho `docs/en/`; havia `README.md` e `README.en.md`.
- Fato: o `/dh:setup` perguntava o idioma e gravava `language` (`en` ou `pt-br`) no `project.yaml`. O Coordenador respondia nesse idioma, e a chave escolhia a pasta de `templates/`.
- Fato: manter o espelho custava uma checagem de paridade no `dh validate`, uma pasta de template por idioma, mensagens em pares nos mods e uma tradução a cada mudança.
- Decisão do dono (2026-11-17): o dev-harness passa a ser somente pt-BR. Acabam a divisão `en`/`pt-br` e o seletor de idioma, e não fica nada em inglês além do que é identificador.

## 2. Decisão

- D1 **O kit é só pt-BR.** Prosa, headings, `description`, exemplos, mensagens do `dh` e dos mods, templates, documentos, despachos e respostas entre agentes são em pt-BR.
- D2 **Fica em inglês só o que é identificador:**
  - o que a máquina lê: chaves de YAML, JSON e frontmatter (`name`, `description`, `model`, `reviewers`, `mode`), valores de enum já existentes (`repo`, `global`, `always`, `on-request`) e campos da API do painel;
  - o que a pessoa digita como nome: comandos `/dh:*`, subcomandos e flags do `dh`, nomes de agente, slugs de skill, nomes de arquivo e de pasta (`MEMORY.md`, `TASK.md`, `REVIEW.md`, `PRD-NNN-…`);
  - identificadores de código.
  - Nome técnico citado no meio do texto vai entre crases. Os slugs das 22 skills **não mudam**: trocar quebraria a invocação, as referências cruzadas, o inventário do manifesto e quem já aponta para eles.
- D3 **A opção de idioma deixa de existir.** `language` sai de `ConfigKeys`, do `Config`, do `config_save`, do `dh config`, do schema dos perfis, do `/dh:setup` e dos docs.
  - `dh config set language <valor>` (e `unset`) cai no mesmo caminho de qualquer chave desconhecida: `chave desconhecida "language" (mode, dashboard.sprites, worktree)`, exit diferente de 0, sem mensagem especial.
  - Uma linha `language:` antiga, no `project.yaml` ou no `config.yaml` da pasta global, é **lida e ignorada**, nunca reescrita (um `SaveConfig` posterior a preserva byte a byte). O `dh doctor` avisa, sem falhar nem mudar o exit: `` `language` não existe mais e é ignorada (o kit é só pt-BR); pode apagar a linha de <arquivo> ``.
- D4 **`templates/` é uma pasta só**, sem subpasta de idioma. O PRD e o `REVIEW.md`, que só existiam em inglês, foram traduzidos. O `dh validate` reprova a volta de `templates/en/` ou de `templates/pt-br/`.
- D5 **Commits e comentários de código seguem o que o projeto define.** A fonte é a primeira que diz algo, nesta ordem: `AGENTS.md`, `CLAUDE.md`, `README.md` do projeto e, depois, a chave opcional `conventions.code_language` do `project.yaml` (`en` ou `pt-br`). Se nenhuma define, o padrão é **inglês**. O Coordenador resolve isso no início, grava a origem no `MEMORY.md` e põe o idioma em todo despacho que escreve código ou commit.
  - Não é a opção de idioma que foi removida: a conversa, os documentos e os registros do kit são sempre pt-BR; esta chave vale só para código e commit do projeto consumidor.
  - O `dh` valida a chave: `en` e `pt-br` passam, qualquer outro valor reprova no `dh doctor` (exit 1). O `/dh:setup` não pergunta: só a grava, com a origem, quando acha a definição num arquivo de diretriz.
  - Conflito entre o arquivo de diretriz e a chave: vale o arquivo de diretriz, que é o texto do dono, e o Coordenador avisa o dono uma vez.
  - Neste repositório vale inglês (`CONTRIBUTING.md`), então os commits e os comentários de código do kit continuam em inglês.
- D6 **Regras de escrita do PRD em pt-BR** no lugar do resumo do ASD-STE100: frase curta (até cerca de 25 palavras) e uma ação por frase; voz ativa com sujeito explícito; um termo por conceito; "deve", "pode" e "não deve" com o mesmo peso sempre; número com unidade e limite fechado; sem "etc.", "e afins", "entre outros", "quando possível", "adequado" ou "rápido" sem número; lista para passos e critérios. O texto normativo vai para `.skills/requirements-discovery/SKILL.md` na 0.43.0; o `document-validator` confere com a mesma severidade (ambiguidade bloqueia, estilo puro é Minor).
- D7 **O histórico em inglês é traduzido** (PRD-020 a PRD-034, as revisões `.review.md` e os ADR-009 e ADR-010), em quatro lotes dentro dos PRs da 0.43.0 à 0.46.0, porque o repositório não aceita PR só de documentação. Os nomes de arquivo ficam, e cada arquivo traduzido leva no topo a linha "Traduzido em <data>; o original em inglês está no commit `aa1e80a`". Não há par pt/en no histórico; o que era duplicado (`docs/en/`, `README.en.md`) é apagado na 0.45.0.
- D8 **O painel continua lendo os registros antigos em inglês.** Os cabeçalhos `Created / updated` e `Done on` dos tickets, o `delivered on` dos PRDs e os headings de `EPOCHAL.md` e `RISKS.md` em inglês seguem aceitos na leitura; a escrita nova é pt-BR.

## 3. Consequências

- Quem tinha `language: en` passa a ver o kit em pt-BR. É uma quebra de experiência para o usuário de língua inglesa, registrada em destaque no `CHANGELOG.md` da 0.42.0; o kit está em 0.x, então cabe numa versão minor.
- `MEMORY.md`, `EPOCHAL.md`, `RISKS.md`, tickets, ADR e PRD de um projeto consumidor em inglês ficam como estão: o kit continua lendo esses arquivos, as entradas novas saem em pt-BR e o projeto fica misto até o dono traduzir, se quiser. Um `REVIEW.md` em inglês continua valendo como guia (o kit lê o conteúdo, não o idioma).
- Um script que chama `dh config set language en` falha com "chave desconhecida", como qualquer chave inexistente.
- O plugin atualizado por cima perde `templates/en/` e `templates/pt-br/` do `dist`; nenhum registro de projeto aponta para esses caminhos.
- A extensão VS Code não quebra: os campos de `/api/state` não mudam. O texto dela é outro trabalho, fora deste repositório.
- Os mods passam a ter uma mensagem só (pt-BR) onde havia pares por idioma; nenhuma decisão de negar, perguntar ou avisar muda, só o texto.
- Transição: da 0.42.0 à 0.44.0 o `dh validate` aceitou `description` de skill que comece com "Use when" **ou** "Use quando", e aceitou citações antigas `templates/<lang>/` nas fontes (conferia o arquivo em `templates/`). Desde a 0.45.0 aceita só "Use quando" (o teste `validate_test` reprova "Use when"), e desde a 0.46.0 reprova por lista fixa, com arquivo e linha, os marcadores do formato antigo nas fontes do kit, inclusive `templates/<lang>/` (`internal/kit/legacy_markers.go`).
- Leitores duplos de rótulo: da 0.43.0 à 0.45.0, `/dh:plan`, `/dh:review`, `/dh:plan-loop`, o `coordinator` e a `external-clis` leram os rótulos de plano e de relatório em pt-BR e também em inglês, para não quebrar `PLAN.md` e relatórios escritos antes da 0.43.0. A 0.46.0 encerrou a leitura em inglês (a trava e a leitura dupla não convivem: a trava reprova a forma antiga no texto dos leitores). Quem tem um `PLAN.md` ou um relatório antigo renomeia os rótulos à mão; o painel continua lendo `Done on` e `Created / updated` de tickets antigos (leitura no código, sem texto em inglês nas fontes).

## 4. Fatias

| Versão | O que entra |
|---|---|
| 0.42.0 | Esta decisão; remoção de `language` (código, schema, setup, docs); `conventions.code_language`; `templates/` único; `dh validate`, `roadmap_check` e `doc_counts` só em pt; mensagens do `dh` e dos mods em pt-BR |
| 0.43.0 | As 22 skills em pt-BR, as regras de escrita do PRD e o lote H1 do histórico (ADR-009 e ADR-010) |
| 0.44.0 | Agentes, comandos e perfis em pt-BR, regra do Coordenador e lote H2 (PRD-020 a PRD-027) |
| 0.45.0 | Docs e restos, `docs/en/` e `README.en.md` removidos, evals traduzidas, `validate` só com "Use quando" e lote H3 (PRD-028 a PRD-034) |
| 0.46.0 | Trava contra o inglês voltar no `dh validate` (lista fixa de marcadores), fim da leitura dupla de rótulos, lote H4 (as 25 revisões `.review.md` em inglês) e as citações em inglês que restavam em PRDs e ADRs antigos; a migração termina aqui |

## 5. Alternativas recusadas

- **Manter os dois idiomas.** Custo permanente de espelho e de paridade, sem benefício para o dono, que só escreve em português.
- **Manter a opção de idioma com um único valor.** Sobraria chave, validação, pasta de template e mensagem em par para uma escolha que não existe mais.
- **Traduzir os slugs das skills.** Quebra a invocação e as referências cruzadas; a `description` em pt-BR já é o que decide quando a skill é usada.
- **Congelar o histórico em inglês.** Deixaria 41 arquivos (cerca de 90 mil palavras) fora da regra "não fica nada em inglês"; o dono pediu a tradução.
