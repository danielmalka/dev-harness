# Changelog

## Não lançado

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
