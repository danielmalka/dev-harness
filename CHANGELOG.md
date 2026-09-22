# Changelog

## 0.5.1 — 2026-09-22

- Segunda rodada do ciclo: duas regras novas no Coordenador — só se registra como revisor quem foi de fato despachado (uma entrada não chamada vira not-run com motivo, e a leitura do próprio Coordenador nunca ocupa o lugar de um revisor), e o relatório mesclado registra o veredito de cada revisor ao lado do mesclado, para que uma aprovação derrubada por um achado bloqueador continue visível. `dh validate` passa a cobrir também fixtures que carregam cópia de corpo de agente. Sete graders viraram mecânicos (`tool_used`, `regex`) ou foram reescritos para forma observável; dois checks negativos foram restaurados depois que a revisão independente mostrou que a mecanização os tinha perdido.
- Três regras novas em `.agents/coordinator.md` ("External CLI reviewers"), cada uma vinda de um caso de eval que falhou no baseline 0.5.0: o relatório declara que um revisor not-run não gasta rodada de correção; a deduplicação de achados equivalentes preserva a redação de cada revisor — suspeitar do mecanismo relatado vira divergência registrada para o dono, nunca substituição silenciosa do texto; e o relatório mesclado reproduz verbatim os títulos e rótulos da skill da etapa (`## Findings`, `### <SEVERITY>: <título>`, `- Reported by: <tag>`, sem negrito nem aspas) mesmo aninhado no formato de saída do Coordenador.
- `dh validate` passa a comparar o `append_system_prompt` de cada `evals/cases/<id>/prompt.md` com o corpo do agente que ele copia e falha com "append_system_prompt no longer matches `.agents/<papel>.md`". O runner não aceita referência a arquivo nesse campo, então a cópia é o único texto que o caso realmente executa: quando ela defasa, todo baseline daquele caso mede as instruções antigas em silêncio — foi o que aconteceu entre a 0.5.0 e esta correção. As onze cópias do corpo do Coordenador foram re-sincronizadas.
- Oito correções do lado dos casos, onde o critério era inobservável por construção: `external-clis-only` e `external-clis-parallel` pediam uma seção `## Dispatches` que o próprio prompt mandava omitir; dois graders de `external-clis-parallel` liam concorrência num trace com um único despacho real; `external-clis-transport-failure` cobrava um revert num caso sem ferramenta de escrita; e o grader de audiência de `doc-validator-gap` virou forma observável. Nenhum grader passou a aceitar promessa mais fraca do que o asset garante.
- Medição antes/depois em `evals/baselines/2026-09-22-0.5.0.md` → `2026-09-22-0.5.1.md`. Com três rodadas nos quatro casos decisivos, o defeito real do ciclo (`external-clis-merge-blocker`, veredito individual de cada revisor dissolvido na mescla) vai de FFF para PPP, e dois casos que pareciam ter regredido numa rodada única passam PPP — era variância. Rodada única do conjunto inteiro oscila entre 0,84 e 0,91 com os mesmos assets, então ela é o retrato que acompanha o pacote, não o critério de aceitação. Continuam em aberto, com evidência no companion: sobre-despacho em `external-clis-parallel` (`Agent called 2x`, pego pelo grader novo na primeira execução) e uma rodada em três que renomeou e colocou em negrito o rótulo do veredito apesar da regra de reprodução verbatim.

## 0.5.0 — 2026-09-22

- Runner de evals executável (PRD-004): `evals/` migra do formato arquivo único para o layout nativo `claude plugin eval` (Claude Code CLI 2.1.278) — `evals/cases/<id>/{prompt.md,case.yaml,scaffold.sh,fixtures/,graders/}`. Os 10 casos antigos migraram e ganharam 4 novos cobrindo `reviewers.document`/`reviewers.security` (`external-clis-document-only`, `external-clis-document-merge-blocker`, `external-clis-security-only`, `external-clis-security-merge-blocker`); `doc-validator-round-cap` não é dirigível pelo runner e vive em `evals/not-run/doc-validator-round-cap/` com `REASON.md`. `evals/rubrics/` e os fixtures compartilhados `evals/fixtures/external-clis-review`/`prd-review` foram removidos — cada caso leva seus próprios graders e fixtures agora; `evals/fixtures/slice-01` continua.
- Primeiro baseline real gravado, `evals/baselines/2026-09-22-0.5.0.{json,md}`: 13 casos pontuados, 4 passaram (`doc-validator-clean`, `external-clis-document-merge-blocker`, `external-clis-document-only`, `external-clis-security-only`), score geral 0,74, custo 4,81 USD; os 9 casos falhos ficam registrados como achados comportamentais para um próximo ciclo de `/dev-harness:improve`, não como defeito do runner. `external-clis-risk-001`: os graders anti-delegação passaram, só `06-no-paraphrasing` falhou — RISK-001 continua `mitigado`.
- Cobertura de revisor CLI nas três etapas (`document`, `security`, `code`, RF-01 a RF-04) fecha o que faltava do PRD-003.
- `dh build` passa a copiar `evals/` para o pacote gerado, excluindo `results/`, `baselines/` e `not-run/` no nível raiz e `__pycache__`/`*.pyc` em qualquer profundidade; `.gitignore` ganha `evals/results/`.
- `dh validate` ganha duas checagens: paridade da cláusula anti-delegação entre `.agents/coordinator.md` e `.skills/external-clis/SKILL.md` (RF-05), e resolução das referências relativas dentro de `evals/` — `context.add_dirs`/`context.scaffold_script`/`context.history_file` de `case.yaml` e `target.path` de grader — além de cobertura fonte-vs-pacote para `evals/` (RF-06).
- `.github/workflows/evals.yml` (novo, `.github/workflows/ci.yml` intocado): roda `claude plugin eval` com `--scaffold --threshold 0.8` na CLI fixada em 2.1.278 só quando o secret `ANTHROPIC_API_KEY` existe; sem o secret, termina verde com `skipped: ANTHROPIC_API_KEY absent`.
- Correções de prosa: o formato de saída do Coordenador (`.agents/coordinator.md`) passa a nomear o revisor CLI (binário/slug ou motivo de not-run) no bloco de Dispatches (RF-07); o mecanismo de congelamento do workspace antes de uma chamada de CLI reviewer fica concreto em `.agents/coordinator.md` e `.skills/external-clis/SKILL.md` — `git status --porcelain -uall` mais `git hash-object` por caminho, recalculado depois da chamada (RF-08); o ciclo de validação de documento (criação → `document-validator` → até duas rodadas) passa a cobrir `PLAN.md` em `.commands/plan.md` (persistido em `.harness/tasks/<PRD-id>/PLAN.md`, relatório em `PLAN.review.md`) e os templates PRD/STORY/TASK/ADR em `.commands/document.md` (RF-10, RF-11); a ordem build-antes-de-test para uma fatia que altera `dist/` fica explícita em `.skills/implementation-planning/SKILL.md`, com remissão de uma linha em `.skills/delivery-readiness/SKILL.md` (RF-14, RF-15); `.skills/harness-evaluation/SKILL.md` reescrito para o layout nativo (RF-17).
- Duas lições do runner aprendidas ao rodar o baseline, incorporadas em `.skills/harness-evaluation/SKILL.md`: `context.add_dirs` não materializa fixtures no cwd do workspace — cada caso precisou de `scaffold.sh` próprio rodado via `--scaffold`; graders `regex` não aceitam flags inline (`(?m)`/`(?mi)`) — passaram para a chave `flags:`.

## 0.4.0 — 2026-09-21

- Revisão em pares por CLI externa (PRD-003, RF-01 a RF-10): skill `.skills/external-clis/` (`SKILL.md` e as receitas `references/{agy,codex,grok,opencode,mcode,claude,gotchas}.md` mais `references/local.yaml.example`), portada e generalizada de uma skill pessoal, ensina o Coordenador e os papéis de revisão a resolver um binário de CLI (`command -v` ou `.harness/local.yaml`), montar o prompt em disco, aplicar o modo somente-leitura do binário e distinguir falha de transporte de veredito real (RF-01, RF-02).
- Campo opcional `reviewers:` em `.harness/project.yaml` (RF-03), documentado numa subseção própria de `profiles/README.md` ("Consumer-only fields") e ilustrado em `profiles/examples/project.yaml`: por etapa (`document`, `code`, `security`), lista revisores como `claude`, `"cli:<binário>/<slug>"` ou `{reviewer: "cli:<binário>/<slug>", timeout_minutes: <n>}`; ausente, o fluxo continua só Claude, como hoje.
- `.agents/coordinator.md` ganha a seção "External CLI reviewers" (RF-02, RF-04, RF-05, RF-07, RF-08, RF-10): resolução de binário, despacho em paralelo ou em levas, montagem do prompt a partir de `.agents/<role>.md` (sem frontmatter) + skill de revisão + contexto do despacho, a cláusula anti-delegação embutida verbatim no corpo do prompt (RISK-001), mescla de N vereditos (achado bloqueador de qualquer revisor vence, achados equivalentes deduplicados, divergência registra os dois vereditos sem o Coordenador decidir sozinho), degradação not-run por revisor ou por etapa inteira quando todos ficam not-run, somente-leitura com reversão por `git status` pós-chamada, e escrita de `.harness/local.yaml` com aviso único quando o caminho não está no `.gitignore`; revisor CLI conta contra o teto de dois especialistas concorrentes.
- Campo `Reported by: claude | cli:<binário>/<slug>` acrescentado ao formato de achado de `.skills/code-review/SKILL.md`, `.skills/document-review/SKILL.md` (na resposta e no corpo persistido em `<documento>.review.md`) e `.skills/security-review/SKILL.md` (RF-06, AC-09); `.commands/review.md`, `.commands/discover.md` e `.commands/secure.md` passam a despachar o agente Claude da etapa só quando `reviewers.<etapa>` está ausente ou contém `claude`, citando `external-clis` para entradas `cli:...`; `.commands/secure.md` ganha a frase de autorização de provedor/modelo; `.skills/project-onboarding/SKILL.md` registra que `.harness/local.yaml` é escrito pelo Coordenador na primeira chamada de um revisor CLI, não pelo `setup`.
- Sete casos de avaliação em `evals/cases/external-clis-*.md` (`parallel`, `only`, `missing-binary`, `merge-blocker`, `attribution`, `risk-001`, `transport-failure`), rubrica `evals/rubrics/external-clis.md` e fixture com "double" de transporte em `evals/fixtures/external-clis-review/` (RF-09) — cobrem paralelo Claude+CLI, revisão só-CLI, binário ausente, mescla com bloqueador vindo só da CLI, atribuição, a regressão de RISK-001 (a cláusula anti-delegação precisa estar no corpo do prompt, nunca só em frontmatter) e falha de transporte sem consumir rodada de correção; execução comportamental não registrada — `evals/baselines/` continua só com `.gitkeep`, mesmo status do PRD-002.
- `disallowedTools: [Agent]` removido do frontmatter dos 18 arquivos de especialista em `.agents/` (o `coordinator.md` nunca declarou o campo) e `.skills/harness-authoring/SKILL.md` atualizado (RISK-001): o runtime não impunha o campo — um especialista despachou outro agente pela ferramenta `Agent` apesar da restrição declarada só no frontmatter —, então a regra "só o Coordenador despacha" passa a viver no corpo do prompt de cada especialista.
- `model: sonnet` → `model: opus` em `.agents/coordinator.md` (decisão do dono).

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
