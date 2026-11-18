# PRD-033 · Worktrees do kit em `<checkout>/.worktrees/`

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | aprovado em 2026-11-16 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-11-16 / 2026-11-16 |
| Tickets | a definir pelo planejador |
| Origem | Decisão do dono em 2026-11-16, achado A2 do dossiê do repositório, card 212 |
| Depende de | PRD-028 e ADR-008 (a decisão D1 é emendada aqui); versão base 0.39.2 |
| Entrega | Um pull request, versão 0.40.0 |

## 1. Problema

Fatos:

- O kit cria um worktree de desenvolvimento em uma pasta ao lado do checkout, `<parent>/<checkout folder>.worktrees/<task-id>` (`.commands/build.md`, seção "Worktree", Criar e Registrar; `docs/adr/ADR-008-worktree-de-desenvolvimento.md`, D1; `README.md`, `README.en.md`, os dois tutoriais).
- O dono fixou outra regra para os agentes que trabalham nos repositórios: `<repo>/.worktrees/<task>`. O PR #58 já põe `.worktrees/` no `.gitignore` deste repositório. O kit contradiz a regra do próprio dono.
- Uma pasta vizinha fica fora do diretório da sessão. Os tutoriais anotam que ninguém mediu se o Claude Code pede permissão de escrita ali.
- Um worktree dentro do checkout quebra o `dh validate`. Reproduzido em `main` na 0.39.2: `go run ./cmd/dh validate .` em um checkout que tem um worktree em `.worktrees/` falha com os erros "private path" e "unresolved link" na cópia do repositório. O percorredor `iterTextFiles` (`internal/kit/validate.go`) não pula `.worktrees`.
- O resolvedor (`internal/harness/worktree.go`) lê o gitfile e não depende de onde fica a pasta do worktree. O teste com git real `TestR18_R19_R26_R61_RealGitWorktree` (`cmd/dh/worktree_cmd_test.go`) só monta o caminho vizinho antigo.

## 2. Solução

Depois da entrega, o kit cria o worktree de uma tarefa ou de uma execução de `/dh:auto` em `<checkout>/.worktrees/<task-id>`. A pasta é ignorada pelo git. O `/dh:setup` acrescenta a linha `.worktrees/` ao `.gitignore` do projeto, e o `/dh:build` para antes de criar um worktree se a pasta não for ignorada. O `dh validate` pula a pasta. Os worktrees que já existem na pasta vizinha antiga continuam funcionando: nada é movido. A decisão D1 do ADR-008 é emendada, e os docs dizem uma coisa só.

## 3. Regras

- R1 `.commands/build.md` (Criar, Registrar, Limpeza), `.commands/resume.md` e `.commands/auto.md` usam `.worktrees/<task-id>` sob o checkout principal como caminho, e a mesma forma relativa no registro da tarefa (`TASK.md`, `BRIEF.md`) e em `MEMORY.md`.
  - Quando `/dh:build` cria um worktree para a tarefa `T-9`, então o caminho é `<main checkout>/.worktrees/T-9`, nunca sob `<home>` e nunca uma pasta vizinha.
  - Quando um registro é escrito, então o caminho relativo dele é `.worktrees/T-9`.
  - Quando `/dh:auto` cria o worktree de uma execução, então a pasta dele é `.worktrees/AUTO-<n>`.
- R2 `/dh:build` checa o ignore antes de criar um worktree e nunca edita o `.gitignore` por conta própria.
  - Quando `git -C "<main checkout>" check-ignore -q -- ".worktrees/<task-id>"` falha, então `/dh:build` para antes de `git worktree add`, diz ao dono para rodar `/dh:setup` ou para acrescentar a linha `.worktrees/` ao `.gitignore` do projeto, e não altera nenhum arquivo.
  - Quando a checagem passa, então a criação segue como antes (checagem do id, `git worktree add`, checagem de `dh harness-path`).
- R3 `/dh:setup` acrescenta a linha `.worktrees/` ao `.gitignore` da raiz de um repositório git quando `git check-ignore -q .worktrees/x` falha. Um especialista (`builder`, `lane` `infra`) a escreve, a checagem roda de novo, e o relatório a lista em Criado ou Deixado intacto. Esta é a única edição de `.gitignore` que o setup faz, e a linha Limites de `setup.md` diz isso. Uma pasta que não é repositório git não precisa de nada.
  - Quando a linha falta em um repositório git, então o setup acrescenta exatamente `.worktrees/` e nenhuma outra linha.
  - Quando a pasta já é ignorada, então o setup não altera nada e reporta Deixado intacto.
  - Quando o projeto não é um repositório git, então o setup não escreve nenhum `.gitignore`.
- R4 Os registros feitos antes da 0.40.0 continuam funcionando. `/dh:resume` e as checagens de limpeza de `build.md` aceitam o caminho relativo antigo `../<checkout folder name>.worktrees/<task-id>` para esse worktree existente. Todo worktree novo usa o caminho novo. Nada é movido, convertido ou removido.
  - Quando um registro guarda o caminho antigo e `git worktree list --porcelain` o lista com o branch registrado, então `/dh:resume` o reutiliza.
  - Quando um registro guarda qualquer outro caminho que não o novo nem o antigo, então `/dh:resume` o reporta e para, como antes.
- R5 `dh validate` pula a pasta `.worktrees` quando lê o repositório (`iterTextFiles`).
  - Quando uma cópia de worktree do repositório fica em `<root>/.worktrees/T-1` e guarda um caminho privado e um link quebrado, então `Validate` não reporta nenhum dos dois.
  - Quando o mesmo teste roda em `main` na 0.39.2, então ele falha (reprodução antes da correção).
- R6 O resolvedor não muda. O teste com git real monta o worktree em `<main>/.worktrees/T-X`, nos modos `repo` e `global`, e mantém todas as asserções de hoje.
  - Quando `go test ./cmd/dh/ -run RealGitWorktree -v` roda com `git` no PATH, então os dois subtestes passam e nenhum é pulado.
- R7 Os docs dizem uma coisa só. `AGENTS.md` ganha uma linha nova com a regra. `docs/adr/ADR-008-worktree-de-desenvolvimento.md` D1 e suas consequências são emendadas com a data e este PRD; a tabela de comparação fica como registro. `README.md`, `README.en.md`, `docs/tutorial.html` e `docs/en/tutorial.html` nomeiam a pasta nova e o ignore, e dizem que o setup também escreve a linha `.worktrees/` no modo `global`. A nota de permissão é substituída pela nota de compatibilidade. `.skills/external-clis/references/clone-isolation.md` diz que `.worktrees/` precisa ser ignorada e declara o resíduo: uma chamada cujo projeto é o checkout principal lista os worktrees no status `--ignored`, então uma escrita do `builder` durante a chamada descarta o veredito (um alarme falso, não uma falha em detectar).
  - Quando os docs vivos e os comandos são pesquisados atrás da pasta vizinha antiga, então os únicos resultados são as frases de compatibilidade (R4) e o texto da emenda.
  - Quando o PRD-028 é lido, então uma nota no topo dele diz que a R63 é substituída por este PRD. O histórico em `CHANGELOG.md` e nos cards antigos do roadmap mantém as palavras dele.
- R8 A versão 0.40.0 é definida em `internal/build/build.go`, em `.claude-plugin/marketplace.json` (versão e ref `v0.40.0`) e em `harness-manifest.json`. `CHANGELOG.md` tem uma entrada datada de 2026-11-16. `docs/roadmap.html` e `docs/en/roadmap.html` têm chips atualizados, a faixa de fontes e um card novo de versão entregue. `dist/` é reconstruído, binários incluídos. As contagens continuam em 13 agentes, 19 comandos mais 1 comando de mod e 22 skills.
  - Quando os arquivos nomeados nesta regra são lidos, então cada campo de versão diz 0.40.0.
  - Quando `go run ./cmd/dh build` roda, então `dist/` não mostra diferença, exceto a linha `generated_at` dos manifestos.
  - Quando `go test ./...`, `go run ./cmd/dh validate .` e a prova de clone limpo do `AGENTS.md` rodam, então todos passam. Este PRD toca `internal/` e `cmd/`, então a prova é obrigatória.

Fora de escopo: mover ou apagar worktrees que existem; uma edição automática de `.gitignore` pelo `/dh:build`; `.git/info/exclude`; `worktree` na lista da trava do `gh`; a ferramenta `EnterWorktree`; o clone de sandbox dos revisores de CLI; qualquer comando, agente ou skill novo.

## 4. Docs

- `.commands/build.md`, `.commands/resume.md`, `.commands/auto.md`, `.commands/setup.md`: os fluxos tocados (R1 a R4).
- `docs/adr/ADR-008-worktree-de-desenvolvimento.md`, `docs/prd/PRD-028-worktree-and-approval-page.md` (nota), `AGENTS.md` (R7).
- `README.md`, `README.en.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `.skills/external-clis/references/clone-isolation.md` (R7).
- `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html` (R8).

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Manter a pasta vizinha | O dono fixou `<repo>/.worktrees/<task>` para os agentes, e a pasta vizinha exige uma permissão de escrita fora do diretório da sessão. |
| `.git/info/exclude` em vez de `.gitignore` | O dono pediu `.gitignore`, e uma linha versionada também protege os outros clones e agentes do projeto. |
| `/dh:build` edita o `.gitignore` por conta própria | O `.gitignore` pertence ao projeto. O setup é o único lugar que pode escrever a linha, e o Coordenador escreve só os próprios registros. |
| Mover os worktrees vizinhos existentes | Um worktree guarda o trabalho do dono. O resolvedor e o git não se importam com onde a pasta fica, então nada precisa mudar de lugar (R4). |
| Um ADR novo | A decisão emenda a D1 do ADR-008 e mantém todas as outras partes. A emenda é datada no ADR. |

### Riscos e decisões pendentes

- Risco: um projeto que nunca rodou `/dh:setup` depois da 0.40.0 não tem a linha de ignore. Mitigação: a R2 faz o `/dh:build` parar antes de criar o worktree e nomeia as duas correções.
- Risco: uma chamada de revisão cujo projeto é o checkout principal vê as escritas do `builder` em `.worktrees/`. Mitigação: o alarme falso é seguro (descarta o veredito). O resíduo está escrito em `clone-isolation.md`, e a revisão aponta para o caminho do worktree.
- Risco: o `setup` e o tutorial diziam que o modo `global` escreve só `REVIEW.md` na árvore. Mitigação: a R7 muda as duas frases.

### Decisões (fechadas em 2026-11-16, dono)

- D1 A pasta do worktree é `<checkout>/.worktrees/<task-id>` (dossiê A2, card 212).
- D2 `/dh:setup` garante a linha `.worktrees/` no `.gitignore` do projeto consumidor (card 212).
- D3 Aumento de versão menor, de 0.39.2 para 0.40.0: um comando muda de comportamento e o `/dh:setup` escreve uma linha nova de arquivo.
