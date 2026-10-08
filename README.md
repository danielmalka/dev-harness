# dev-harness

English: [README.en.md](README.en.md)

Kit portátil de desenvolvimento assistido por IA. Clone, carregue no Claude Code, desenvolva.

Licença MIT. Plano de produto: [`docs/plano-produto.html`](docs/plano-produto.html). Tutorial: [`docs/tutorial.html`](docs/tutorial.html).

## Instalar

Pelo marketplace do próprio repositório, sem clonar, dentro de uma sessão do Claude Code no projeto em que você vai trabalhar:

```text
/plugin marketplace add danielmalka/dev-harness
/plugin install dh@dev-harness
```

A versão vem da tag fixada em `.claude-plugin/marketplace.json`. Para um time, o projeto pode declarar o marketplace em `extraKnownMarketplaces` e o plugin em `enabledPlugins` no `.claude/settings.json`, e o Claude instala para quem confiar na pasta.

**Versão atual: 0.19.1** (tag `v0.19.1`). Extensão VS Code ([`dev-harness-vscode`](https://github.com/danielmalka/dev-harness-vscode), verificada em VS Code real no Remote WSL e no Windows): view com dashboard local, status bar com limites, lançador de sessão; instale o `.vsix` da [release mais recente](https://github.com/danielmalka/dev-harness-vscode/releases/latest) via "Extensions: Install from VSIX...". Dashboard local em http://127.0.0.1:4747 com `/dashboard` (registrado pelo mod); mostra projetos de `DH_DASHBOARD_ROOTS`, uma barra de progresso por PRD aberto, sessões do Claude Code com estado, avatar com limites de 5h/wk (layout 70/30: projetos à esquerda, avatar e limites à direita sticky). Painel de status abaixo do prompt: contexto, agentes, dh slice progress, `gate pending`, limites 5h/wk (🔴 80%+), custo. Travas de commit: nega menção a IA em `git commit`, `gh pr create/edit`, `gh release create`; nega commit quando houve edição após último gate verde; agora escopo do repositório da cada commit (não mais bloqueia edições em outro repo). Validação: ticket aberto sob PRD marcado `entregue em`, ou todos os tickets prontos mas PRD não. Mensagens de painel e dh seguem idioma de `.harness/project.yaml`. Testes: 147 pass. Kit enxuto (13 agentes, 19 comandos, 21 skills) e evals fora do fluxo. Tutorial: [`docs/tutorial.html`](docs/tutorial.html). CLIs externas rodam em clone descartável. PRD com quatro seções; ticket único substitui story. `/dh:plan-loop` disputa ondas com CLIs; `/dh:auto` encadeia discover-plan-build. Roadmap: [`docs/roadmap.html`](docs/roadmap.html). Notas: [`CHANGELOG.md`](CHANGELOG.md).

**Migrando de uma instalação anterior à 0.8.0:** o id do plugin mudou de `dev-harness@dev-harness` para `dh@dev-harness` e os comandos de `/dev-harness:<comando>` para `/dh:<comando>`. Rode `/plugin uninstall dev-harness@dev-harness` e depois `/plugin install dh@dev-harness`; troque `/dev-harness:` por `/dh:` em scripts e anotações próprias. Nome do marketplace, repositório, binário Go `dh` e caminho de snapshot não mudam. Ver `CHANGELOG.md` 0.8.0.

## Carregar a partir do clone (desenvolvimento)

Pré-requisitos: Git e Claude Code autenticado com acesso a um modelo. Python 3 só para a fixture `slice-01`. Go 1.22+ só para quem mantém o kit.

No projeto em que você vai trabalhar:

```bash
claude --plugin-dir "/caminho/absoluto/do/clone/dist/claude-code/dev-harness" --agent coordinator
```

Substitua o caminho pelo resultado de `pwd` no clone, mais `/dist/claude-code/dev-harness`. Atenção: dentro de aspas o `~` não expande e o Claude ignora o plugin em silêncio (`--agent coordinator` falha com "not found"). Use `$HOME/...` ou o caminho completo. Se o runtime exigir o nome qualificado, use `--agent dh:coordinator`.

Depois, na sessão:

```text
/dh:doctor
/dh:setup
```

O diagnóstico mecânico também roda sem o Claude:

```bash
dist/claude-code/dev-harness/bin/dh doctor dist/claude-code/dev-harness
```

Tutorial completo: [`docs/tutorial.html`](docs/tutorial.html). Fixture da primeira fatia: [`evals/fixtures/slice-01`](evals/fixtures/slice-01).

## Fontes e pacote

| Caminho | Papel |
| --- | --- |
| `.agents/` `.commands/` `.skills/` `templates/` `profiles/` | Fontes canônicas |
| `cmd/dh`, `internal/` | Binário Go `dh`: `validate`, `build`, `doctor`, `snapshot`. `go run ./cmd/dh build` gera `dist/claude-code/dev-harness` |
| `dist/` | Pacote gerado. Não editar. |
| `adapters/claude-code/` | Mapeamento para o Claude Code |
| `adapters/generic/` | Uso manual em outra IA, sem paridade |

Templates de trabalho em `templates/en/` e `templates/pt-br/` (versões espelhadas): PRD, Task/Bug (o ticket), ADR e RFC, além dos três registros de memória (MEMORY, EPOCHAL, RISKS). O `setup` grava `language` em `.harness/project.yaml` e escolhe a pasta. Quando usar cada um: [`docs/tutoriais/templates.html`](docs/tutoriais/templates.html).

Perfis (`base`, `go-api`, `typescript-web`, `php`) estão em inglês em `profiles/`.

## Idioma

- O Coordenador responde no idioma em que você escreve. Inglês é o padrão quando não há sinal; escreva em português e ele responde em português.
- No primeiro `/dh:setup`, o idioma da sessão é gravado em `.harness/project.yaml` como `language: en` ou `language: pt-br`. Esse campo escolhe `templates/<lang>/` para os registros de memória e os artefatos de trabalho (PRD, task, ADR, RFC). Para trocar, edite o campo ou peça ao Coordenador; registros já escritos não são traduzidos.
- Entre agentes tudo é inglês: despachos, respostas dos especialistas, código, commits e casos de eval. `project.yaml` e perfis também.
- Documentos HTML gerados pela skill `doc-template-html` seguem o mesmo campo: `bash .skills/doc-template-html/scripts/stamp.sh --lang en ...` (padrão `pt-br`).
- Documentação do kit: pt-br em `docs/` e inglês em `docs/en/`; cada página tem link para a outra versão. Este README tem versão em [`README.en.md`](README.en.md). O plano de produto é interno e existe só em pt-br.

## Manutenção

```bash
go test ./...
go run ./cmd/dh validate --source-only .
go run ./cmd/dh build
go run ./cmd/dh validate .
```
