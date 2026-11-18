---
name: setup
description: Proponha o perfil do projeto e inicialize os registros
author: malka
argument-hint: "[diretório de destino ou notas de setup | bootstrap]"
metadata:
  roles: [coordinator, repo-scout]
  skills: [project-onboarding]
  writes: project.yaml e os registros ausentes no harness dir; REPO-MAP.md no harness dir quando o dono aceita; REVIEW.md na raiz do repositório quando ausente; a linha `.worktrees/` no `.gitignore` do projeto quando não está ignorada; a linha `Fatos do sistema:` no `AGENTS.md` (ou no `CLAUDE.md`, ou num `AGENTS.md` mínimo novo) só com o sim do dono
disable-model-invocation: true
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `<harness dir>/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Destino e notas: $ARGUMENTS (use o diretório de trabalho atual como padrão e diga isso). Despache `repo-scout` com a ferramenta Agent, modelo `haiku`, para contribuir com a parte de inventário somente leitura da skill do kit `project-onboarding`: stack, manifestos, instruções de projeto existentes, comandos de teste, lint, build e execução que se possam descobrir, com a fonte de cada um. Depois carregue você mesmo `project-onboarding` e produza o perfil a partir desse inventário. Decida onde os registros ficam, nesta ordem: (1) rode `dh harness-path --json` (o `bin/dh` que acompanha o kit (`bin\dh.cmd` no Windows), a partir do diretório do plugin carregado, como em `.commands/doctor.md`; toda chamada `dh` aqui o usa) e leia `defaults`; (2) o modo é `defaults.mode` quando preenchido; caso contrário, pergunte ao dono uma vez para este projeto (`repo` ou `global`); com `defaults` todo vazio (sem `config.yaml` ou sem chaves), use os padrões de hoje; (3) em `repo`, crie `<repo>/.harness/` com os registros a partir dos templates; (4) rode `dh link [name]`, que sempre registra o projeto (em `repo` só registra, então `<home>/projects/<name>/` não existe depois; em `global` também cria essa pasta), depois cole o bloco `defaults.reviewers` no novo `project.yaml`; a partir daí os valores do próprio projeto prevalecem. Nunca leia `config.yaml` com a ferramenta Read. Em `global` os registros nascem em `<home>/projects/<name>/` e o único arquivo escrito na árvore do repositório é `REVIEW.md` (veja o passo do guia de revisão abaixo); lembre o dono de acrescentar `~/.harness` a `additionalDirectories` em `~/.claude/settings.json` (nunca o escreva).

Guia de revisão, nos dois modos (skill `project-onboarding`): quando o projeto não tem arquivo de guia de revisão, despache um especialista (`builder`, `lane` `infra`; não o Coordenador, já que o arquivo é novo) para criar `REVIEW.md` na raiz do repositório a partir de `templates/REVIEW.md`, depois escreva `review_guide: REVIEW.md` (um caminho relativo ao repositório, nunca absoluto nem com `..`) no `project.yaml`. O arquivo é escrito em pt-BR (ADR-011). Quando o arquivo existe mas a chave está ausente, registre a chave e deixe o arquivo intacto. Quando a chave aponta para um arquivo existente, não mude nada. Nunca sobrescreva.

Ignore do worktree, em um repositório git: os worktrees de desenvolvimento do kit ficam em `<repo>/.worktrees/<task-id>`. Rode `git check-ignore -q .worktrees/x` na raiz do repositório; quando falhar, despache um especialista (`builder`, `lane` `infra`) para acrescentar a única linha `.worktrees/` ao `.gitignore` da raiz (criando o arquivo quando ausente), depois rode a checagem de novo. Reporte isso em Criado ou Deixado intacto. Uma pasta que não é repositório git não precisa de nada.

Projeto só de documentação (sem código): uma pasta que não é repositório git é aceita; registre os campos de `commands` como null; `REVIEW.md` lista apenas os documentos ou caminhos críticos; o formato do `project.yaml` vem de `profiles/examples/project.yaml`.

Linha de referência do arquivo de fatos (skill `project-onboarding`, seção "Linha de referência do arquivo de fatos"), nos dois modos: pergunte ao dono, uma vez, se acrescenta ao `AGENTS.md` do projeto (senão ao `CLAUDE.md`; sem nenhum dos dois, pergunte se cria um `AGENTS.md` mínimo só com a linha) a linha `Fatos do sistema: .harness/FACTS.md (dica com evidência; o código vence).` no modo `repo`, ou, no modo `global`, `Fatos do sistema: mantidos pelo dev-harness fora deste repositório; caminho em "dh harness-path --json" (campo harness) mais FACTS.md.` Sem o sim do dono, não toque em arquivo nenhum. Quando o arquivo já tem uma linha que começa com `Fatos do sistema:`, não acrescente outra. Despache um especialista (`builder`, `lane` `infra`) para a escrita, com o arquivo de instruções como único conjunto de escrita, e reporte em Criado ou Deixado intacto. Nunca inclua o `FACTS.md` por `@import` e não escreva o bloco `<!-- dh:facts:inicio -->` do modo `global`. O setup não cria o `FACTS.md`.

Bootstrap do repositório (opcional, skill `project-onboarding`): depois do perfil, ofereça escrever o mapa do repositório `<harness dir>/REPO-MAP.md`; o dono aceita ou recusa. Com `bootstrap` nos argumentos, pule o perfil e atualize só o mapa. Obtenha o sha de `HEAD` de uma execução de `qa-verifier` (`git rev-parse HEAD`) e despache `docs-guide` (modelo `sonnet`, skills `repository-mapping` e `project-onboarding`) com ele e os achados do `repo-scout`; quando ele reportar mais de 15 diretórios de primeiro nível, peça ao dono que nomeie o escopo. No modo `repo` o arquivo é versionado. A parte 6 do mapa é só a linha `Convenções: ver FACTS.md`; o `docs-guide` devolve as convenções que observou no bloco `## Fatos candidatos` (no máximo 5), e você as grava no `FACTS.md` pela regra "Arquivo de fatos" do `coordinator` (`dh memory write --file FACTS.md`, `--expect absent` quando o arquivo ainda não existe, e nunca o `docs-guide`).

## Pré-requisitos
Uma raiz de projeto legível e as instruções de projeto existentes, que prevalecem sobre qualquer coisa que o kit proponha. Escrever exige escopo autorizado: sem ele, proponha o perfil na resposta e não escreva. Pergunte ao dono só quando um conflito entre as instruções do projeto e a proposta mudar o que é registrado.

## Saída
O relatório de `project-onboarding` com o modo, o registro, o perfil proposto escrito em `project.yaml` e `MEMORY.md`, `EPOCHAL.md` e `RISKS.md` inicializados (todos no harness dir) a partir dos templates que acompanham este kit, só quando estiverem ausentes e a escrita estiver autorizada. Releia todo arquivo criado e reporte Criado, Deixado intacto e Não aplicado; `REVIEW.md` aparece em Criado ou Deixado intacto. O kit é só pt-BR (ADR-011): o setup nunca pergunta nem registra um idioma. Ele escreve `conventions.code_language` (`en` ou `pt-br`, o idioma das mensagens de commit e dos comentários de código) no `project.yaml` só quando um arquivo de instruções do projeto (`AGENTS.md`, `CLAUDE.md` ou `README.md`) o define, e anota em `MEMORY.md` o arquivo de onde veio; sem definição, a chave é omitida e vale o inglês. Registre em `MEMORY.md` a fonte do perfil, o modo e qualquer conflito não resolvido.

## Limites
- Nunca sobrescreva nem edite um registro existente no harness dir; um arquivo existente é deixado intacto e reportado. A única exceção é `REPO-MAP.md`, dado derivado que uma atualização pode sobrescrever.
- Quando nenhum perfil distribuído corresponde, o perfil é `none, custom`; deixe sem valor os campos sem suporte em vez de adivinhar um comando.
- Nunca instale, atualize nem configure ferramentas, e nunca modifique os arquivos de instruções do próprio projeto, exceto para acrescentar a única linha `Fatos do sistema:` com o sim do dono (Roteamento).
- Nunca edite o `.gitignore`, exceto para acrescentar a única linha `.worktrees/` quando não está ignorada (Roteamento); só o `local.yaml` (escrito depois, na primeira chamada de revisor de CLI) deve continuar ignorado.
- Só o Coordenador escreve `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, e o primeiro `MEMORY.md` com `dh memory write --expect absent` (`context-handoff`, trava de escrita da memória). O `FACTS.md` só o Coordenador grava (candidatos do bootstrap), com `dh memory write --file FACTS.md`.

## Próximo
`/dh:understand` para mapear a primeira área, ou `/dh:discover` quando o pedido ainda é uma ideia.
