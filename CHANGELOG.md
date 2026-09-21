# Changelog

## 0.3.1 — 2026-09-21

- Política de ferramentas dos agentes invertida: `coordinator` não declara `tools` (como sessão principal via `--agent`, a allowlist cortava `ToolSearch` e, com ele, toda ferramenta deferida — `SendMessage`, `Monitor`, `TaskStop`, `WebFetch`, MCP). Os especialistas que escrevem também herdam tudo; os quatro somente-leitura (`code-reviewer`, `security-reviewer`, `document-validator`, `repo-scout`) mantêm allowlist explícita, agora com `LSP`, `ToolSearch`, `Monitor`, `SendMessage` e `WebFetch`. Todo especialista leva `disallowedTools: [Agent]` para que só o Coordenador despache.
- `dh validate`: `tools` passa a ser opcional; presente e vazio continua erro.

## 0.3.0 — 2026-09-20

- Todos os agentes ganham a ferramenta `Skill` na lista `tools`: sem ela, um especialista com ferramentas restritas não consegue carregar a skill do kit que o próprio prompt exige (defeito encontrado na prova do `document-validator`; confirmado empiricamente com `--agent` em modo headless).
- Validador adversarial de documentos (PRD-002): papel `document-validator` (somente leitura, Opus), skill `document-review` com as seis categorias, a regra de severidade e o relatório persistente em `<documento>.review.md`, e o ciclo criação → validação no `/dev-harness:discover` com limite de duas rodadas antes de o PRD chegar ao dono. Casos de avaliação e fixture de PRD em `evals/`.

## 0.2.1 — 2026-09-20

- `dh build` compila com `-buildvcs=false`: binários idênticos entre a máquina do mantenedor e a CI (a 0.2.0 embutia revisão e estado sujo do git, e a CI reprovou a reprodutibilidade).

## 0.2.0 — 2026-09-20

Etapa 1: distribuição pelo marketplace e observabilidade de sessões e agentes no plugin (PRD-001). Binário Go substitui os scripts.

- Binário Go `dh` (`cmd/dh`, `internal/`): `snapshot` (event, subagents, statusline, prune), `validate` (com comparação de conteúdo do pacote), `doctor` (binário, hooks, settings, diretório de snapshots, `.harness/` ignorado) e `build` (cinco alvos cross-compilados, checksums no manifesto). `scripts/` em Python e bash removidos (ADR-001).
- Plugin: `settings.json` com `subagentStatusLine` e `hooks/hooks.json` com seis hooks gravando `~/.claude/dev-harness/sessions/<session_id>.json`; wrapper `bin/dh` por plataforma.
- `.claude-plugin/marketplace.json` na raiz: instalação por `/plugin marketplace add danielmalka/dev-harness`, versão fixada por tag.
- Spec da etapa 1 em `docs/prd/PRD-001-etapa-1.md`; PRD-002 (revisor adversarial de documentos) em rascunho.
- Skill `doc-template-html`: `--lang pt-br|en` no stamp, rótulos e esqueletos por idioma em `assets/skeletons/<lang>/`.
- Docs em inglês em `docs/en/` (início rápido, tutorial 00, tutorial de templates) e `README.en.md`, com link de troca de idioma em cada página; seção Idioma nos dois READMEs.
- Templates em duas versões espelhadas, `templates/en/` e `templates/pt-br/`; validador exige paridade entre idiomas e resolve `templates/<lang>/`.
- Coordenador: regra de idioma (responde no idioma do dono, inglês por padrão, `language` gravado em `.harness/project.yaml` pelo setup; inglês entre agentes) e regra de que os registros em disco prevalecem sobre a memória da conversa. Especialistas respondem em inglês.
- Skill `project-onboarding`: lista os perfis do kit em tempo de execução em vez de uma lista fixa; o perfil `php` passa a ser encontrado.
- Perfil de stack `profiles/php.yaml` (Composer, Pest ou PHPUnit, phpstan, php-cs-fixer; comandos só com evidência no repositório).
- Templates de trabalho em `templates/`: `PRD.md`, `STORY.md`, `TASK.md` (task ou bug), `ADR.md` e `RFC.md` (opcional), em pt-br.
- Tutorial de templates em `docs/tutoriais/templates.md` e `docs/tutoriais/templates.html`, com fluxo de engenharia e relação entre arquivos em Mermaid.
- Comandos em `.commands/` com frontmatter completo: `name`, `author`, `argument-hint`, `metadata.roles`, `metadata.skills`, `metadata.writes` e `disable-model-invocation` em `setup`, `release`, `consolidate-memory` e `improve`.
- Skill `doc-template-html`: diagramas passam a ser código Mermaid via CDN, única exceção de CDN da skill; SVG pré-renderizado deixa de ser aceito.
- `CLAUDE.md` do repositório com as decisões fixadas: idioma por pasta, autor, portabilidade, escrita de memória, Mermaid via CDN e regra de perguntar antes de resolver conflito entre regras.
- Plano de produto: critério AC-15 e texto de tutoriais passam a declarar diagramas Mermaid e fontes web como dependências remotas aceitas.

## 0.1.0 — 2026-09-20

Phase 1 packaging.

- 18 agents, 17 commands and 19 skills as canonical sources.
- Claude Code plugin package generated at `dist/claude-code/dev-harness`.
- Stack profiles `base`, `go-api` and `typescript-web`.
- Static validator and doctor script, including a copy inside the plugin package.
- Guided fixture `evals/fixtures/slice-01`.
- MIT license.
- Tutorial 00 and `docs/inicio-rapido.html`.

Not in this release: runtime proof on a clean Claude Code session, eval case runs, Linux/macOS/WSL2 qualification, native adapters for other IAs.
