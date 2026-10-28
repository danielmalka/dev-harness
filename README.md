# dev-harness

English: [README.en.md](README.en.md)

## O que é

Kit portátil de desenvolvimento assistido por IA, instalado como plugin do Claude Code: um Coordenador, agentes especialistas, skills e comandos `/dh:<comando>` que levam um pedido de ideia a entrega com revisão e verificação. Licença MIT.

Crédito: parte do fluxo de revisão e o relatório de deriva se inspiram nas ideias do Elephant-Goldfish Model (EGM) do Google (https://github.com/google/elephant-goldfish-model, licença Apache-2.0, "not an officially supported Google product"); o kit reimplementa os mecanismos com palavras próprias e não copia texto.

## Requisitos

Claude Code autenticado, com acesso a um modelo; Git.

## Instalar

Dentro de uma sessão do Claude Code, no projeto em que você vai trabalhar. A forma por URL baixa só o `marketplace.json`, sem clonar o repositório:

```text
/plugin marketplace add https://raw.githubusercontent.com/danielmalka/dev-harness/main/.claude-plugin/marketplace.json
/plugin install dh@dev-harness
```

A forma curta `/plugin marketplace add danielmalka/dev-harness` também funciona, mas clona o repositório inteiro (cerca de 90 MB, por causa dos binários `dh` de cada versão). O `install` baixa só a pasta do plugin, na tag fixada em `.claude-plugin/marketplace.json`, de onde vem a versão. Para um time, o projeto pode declarar o marketplace em `extraKnownMarketplaces` e o plugin em `enabledPlugins` no `.claude/settings.json`.

## Início rápido

Na sessão, dentro do projeto, rode `/dh:doctor` e depois `/dh:setup`. O primeiro diagnostica o ambiente; o segundo escolhe o modo (`.harness/` no repositório ou pasta global `~/.harness/`) e registra o projeto. Depois use os comandos de trabalho (`/dh:discover`, `/dh:plan`, `/dh:build`, `/dh:review`, `/dh:verify`, ou `/dh:auto` para encadear). Passo a passo: [tutorial](docs/tutorial.html) e [início rápido](docs/inicio-rapido.html).

O Coordenador responde no idioma em que você escreve (inglês por padrão). O `/dh:setup` grava `language: en` ou `language: pt-br` em `project.yaml`; o PRD é sempre em inglês.

## Atualizar

O plugin vem da tag fixada em `.claude-plugin/marketplace.json`. Para uma versão nova, atualize o marketplace (`claude plugin marketplace update dev-harness`) e o plugin (`claude plugin update dh@dev-harness`); confira com `/dh:doctor`. As mudanças de cada versão estão em [`CHANGELOG.md`](CHANGELOG.md).

## Desinstalar e remover dados

1. Na sessão: `/plugin uninstall dh@dev-harness`. Isso remove o plugin; não apaga seus dados.
2. Pasta global: apague `~/.harness` (ou o caminho de `DH_HOME`) para remover o registro de projetos, os registros do modo global, os snapshots e o token do dashboard.
3. Se você adicionou `~/.harness` em `additionalDirectories` no `~/.claude/settings.json`, remova essa entrada.
4. Em cada repositório em modo repositório: apague a pasta `.harness/` (memória, tarefas, PRDs e ADRs do projeto). O kit nunca a apaga sozinho.

## O que o kit grava e onde

No modo global, a pasta `~/.harness/` (`DH_HOME` muda o lugar) é única por máquina e por ambiente; `.harness/` no repositório vence quando existe. Detalhes na seção "Modo global" do [tutorial](docs/tutorial.html).

| Onde | O que | Quem grava |
| --- | --- | --- |
| `<repositório>/.harness/` | modo repositório: MEMORY, EPOCHAL, RISKS, `project.yaml`, tarefas, PRDs, ADRs | Coordenador e comandos `/dh:` |
| `~/.harness/config.yaml` | registro de projetos e padrões (`language`, `mode`, `worktree`, `dashboard.sprites`) | só o `dh` (`dh config`, `dh link`) |
| `~/.harness/projects/<nome>/` | modo global: os mesmos registros do `.harness/` | Coordenador e comandos `/dh:` |
| `<repositório>/REVIEW.md` | guia de revisão do projeto, criado pelo `/dh:setup` quando falta, nos dois modos; único arquivo que o setup grava na árvore em modo global; nunca sobrescreve | `/dh:setup` (via especialista) |
| `~/.harness/sessions/` | snapshot por sessão do Claude Code (`dh snapshot prune` apaga os com mais de 7 dias) | mod do plugin |
| `~/.harness/dashboard/` | token de parada do dashboard, por porta | dashboard |
| `~/.harness/projects/<nome>/pdocs/` | documentação HTML do projeto (`/dh:document pdocs`) | docs-guide |
| `127.0.0.1:4747` | dashboard local, só nesta máquina, quando você o inicia com `/dashboard` | binário `dh` |

Com `worktree` em `always` (ou quando você pede), `/dh:build` e `/dh:auto` trabalham num worktree irmão `<checkout>.worktrees/<id>`. Quando o `/dh:plan` deixa decisões em aberto, ele gera uma página HTML de aprovação em `pdocs`. Detalhes na seção "Configuração" do [tutorial](docs/tutorial.html).

## Desligar

- Dashboard: `dh dashboard --stop` (`--port N` para outra porta). Ele não mata processo de outro programa.
- Plugin: `claude plugin disable dh@dev-harness` (ou desinstale, como acima); o mod e as travas saem junto.
- Não há chave por hook para o mod nem para a trava de commit: desligar é desabilitar o plugin inteiro.

## Se algo falhar

- Rode `/dh:doctor` na sessão ou `dh doctor` no terminal; ele aponta o que falta.
- Porta ocupada: o dashboard avisa; use `/dashboard --port N` e `dh dashboard --stop --port N`.
- Plugin não encontrado: confira `/plugin marketplace add ...` e `/plugin install dh@dev-harness`; ao carregar de um clone, não use `~` entre aspas em `--plugin-dir` (o Claude ignora o plugin em silêncio).
- Não resolveu: abra uma issue em [GitHub Issues](https://github.com/danielmalka/dev-harness/issues). Vulnerabilidade: ver [`SECURITY.md`](SECURITY.md).

## Usar junto com

Opcional; o kit não precisa de nenhum dos dois.

- [`danielmalka/malka-claude-mods`](https://github.com/danielmalka/malka-claude-mods): mods visuais do Claude Code (status line, menu, Zen, Mode) que saíram deste kit na 0.29.0.
- [`akitaonrails/ai-memory`](https://github.com/akitaonrails/ai-memory): memória de longo prazo e handoff entre CLIs de código.

## Mapa da documentação

- [Início rápido](docs/inicio-rapido.html) e [tutorial](docs/tutorial.html): comece por aqui. [Roadmap](docs/roadmap.html) e [`CHANGELOG.md`](CHANGELOG.md): o que foi entregue e o que vem.
- `docs/prd/`, `docs/adr/` e `docs/discovery/`: histórico de design, em pt-br.
- `evals/`: histórico de medição; não é preciso para usar o kit.
- [`docs/plano-produto.html`](docs/plano-produto.html): plano de produto, em pt-br e interno.
- Páginas em inglês: os mesmos temas dentro de `docs/en/` (início rápido: `docs/en/quick-start.html`).

## Para quem contribui

Esta seção é só para quem mantém o kit; quem usa não precisa dela. Regras completas em [`CONTRIBUTING.md`](CONTRIBUTING.md). Requer Go 1.22+ (e Python 3 só para a fixture `slice-01`).

```bash
go test ./...
go run ./cmd/dh build
go run ./cmd/dh validate .
```

| Caminho | Papel |
| --- | --- |
| `.agents/` `.commands/` `.skills/` `templates/` `profiles/` | Fontes canônicas |
| `cmd/dh`, `internal/` | Binário Go `dh`. `go run ./cmd/dh build` gera `dist/claude-code/dev-harness` |
| `dist/` | Pacote gerado. Não editar. |
| `adapters/` | Mapeamento para o Claude Code (`claude-code/`) e uso manual em outra IA, sem paridade (`generic/`) |

Carregar a partir do clone, no projeto em que você vai trabalhar:

```bash
claude --plugin-dir "$HOME/caminho/do/clone/dist/claude-code/dev-harness" --agent coordinator
```

Use `$HOME/...` ou o caminho completo, nunca `~` entre aspas. Se o runtime exigir o nome qualificado, use `--agent dh:coordinator`. Diagnóstico sem o Claude: `dist/claude-code/dev-harness/bin/dh doctor dist/claude-code/dev-harness`. Templates de trabalho: [`docs/tutoriais/templates.html`](docs/tutoriais/templates.html).
