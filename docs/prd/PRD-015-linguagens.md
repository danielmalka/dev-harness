# PRD-015 · Profiles de Kotlin, Python e TypeScript backend

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-09 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-08 / 2026-10-08 |
| Tickets | (ainda sem tickets; o plano vem depois da aprovação) |
| Origem | `docs/projeto-renova.md` (Fase 1); `docs/discovery/renova.md` §3 |
| Depende de | nada |
| Libera | nada (PRD-016 e PRD-017 não dependem deste) |

## 1. Problema

Os próximos projetos de trabalho do dono são em Kotlin (Quarkus, Micronaut ou Spring Boot) e Python (Django ou FastAPI); que haja também TypeScript no backend é hipótese (o texto do dono diz só "Typescript" e "não sei se já temos"; derrubada se os projetos TypeScript do trabalho tiverem interface, caso em que `typescript-web` já basta, ver decisão 2, opção C). O kit só conhece Go (`go-api`), TypeScript com interface (`typescript-web`, "applications with a rendered user interface", `profiles/typescript-web.yaml:3`) e PHP. Num repositório Kotlin ou Python o `setup` grava `none, custom` (`profiles/README.md:50`) e o builder trabalha só com as convenções de `base`: nenhuma regra de layout, de runner de teste ou de "não sair do padrão" da linguagem. `profiles/README.md:17` adiou Python e outras stacks "até um consumidor real precisar"; esse consumidor chegou.

Fatos verificados em 2026-10-08:

- Schema de profile: `id`, `name`, `summary`, `extends`, `stack`, `commands`, `conventions`, `limits`; sem campo de framework (`profiles/README.md:19-37`).
- Comando só é copiado para `project.yaml` quando observado no repositório; `php.yaml:14-15` já escreve "valores típicos, gravados só com evidência".
- O `setup` lista os profiles em tempo de execução e casa por `stack` hint (`.skills/project-onboarding/SKILL.md:62`); profile novo não exige editar skill ou comando.
- `dh validate` confere `id` == nome do arquivo (`internal/kit/validate.go:438-453`) e exige `base`, `go-api`, `typescript-web` (`:343`).

## 2. Solução

Três profiles novos em `profiles/`, um por linguagem, com o framework tratado como seção de convenções condicionada à evidência encontrada no repositório (decisão 1: A, dono 2026-10-09):

- `kotlin.yaml` (`stack: [kotlin, jvm]`): Gradle (Kotlin DSL ou Groovy) ou Maven; seções para Quarkus, Micronaut e Spring Boot, cada uma com o layout canônico, o comando de dev e o runner de teste típicos, gravados só com evidência.
- `python.yaml` (`stack: [python]`): `pyproject.toml`, `requirements*.txt`, `uv.lock`, `poetry.lock`; seções para Django e FastAPI; `pytest` ou `manage.py test` conforme o que o repositório mostra.
- `typescript-api.yaml` (`stack: [typescript, node]`, sem UI): espelho de `go-api` para serviços HTTP em TypeScript (decisão 2: A, dono 2026-10-09). O matching desempata com `typescript-web` pela presença de interface renderizada (framework de UI, `index.html`, pasta `public/`).

Depois da entrega, `/dh:setup` num repositório Kotlin, Python ou TypeScript backend nomeia o profile certo, copia só os comandos observados e o builder recebe, no `project.yaml`, as convenções da linguagem e do framework detectado. O estilo pedido pelo dono ("clean code, simples antes de complexo") entra como convenção verificável, não como adjetivo.

## 3. Regras

Arquivos e schema

- R1 Entram exatamente três arquivos novos: `profiles/kotlin.yaml`, `profiles/python.yaml`, `profiles/typescript-api.yaml`, no schema de `profiles/README.md` sem campo novo, todos com `extends: base`, em inglês, sem caminho absoluto, segredo ou binário de máquina. Verificável: `dh validate` verde; `ls profiles/*.yaml` mostra os três; `grep -c "^extends: base"` = 1 em cada; nenhum campo fora dos oito do schema.
- R2 `commands:` dos três profiles nasce com `test`, `lint`, `build`, `run` em `null`; os valores típicos por framework ficam em `conventions:` com a frase "recorded only with evidence", no padrão de `php.yaml`. Verificável: `grep -A4 "^commands:"` mostra só `null` nos três; cada seção de framework em `conventions:` traz a marca de evidência.
- R3 Cada profile cobre os frameworks pedidos como seções de convenção com gatilho de evidência explícito: Kotlin, Quarkus (`io.quarkus` no build), Micronaut (`io.micronaut`), Spring Boot (`org.springframework.boot`); Python, Django (`django` nas dependências ou `manage.py`), FastAPI (`fastapi` nas dependências); TypeScript backend, sem framework fixo, lendo `package.json` scripts. Verificável: `grep -i` de cada nome de framework no profile correspondente encontra a seção e o gatilho.
- R4 Convenções de estilo são verificáveis, nunca adjetivo: cada profile nomeia a fonte do padrão (Kotlin Coding Conventions; PEP 8 e o layout `src/` ou de app do framework; o guia de estilo que o `package.json` do repositório já aponta via lint) e proíbe camada nova sem necessidade mostrada com a frase exata "Do not introduce a new module, package or layer for a one-file change" (no espírito de `go-api.yaml:14`). Verificável: `grep -c "Kotlin Coding Conventions" profiles/kotlin.yaml` ≥ 1; `grep -c "PEP 8" profiles/python.yaml` ≥ 1; `grep -c "package.json" profiles/typescript-api.yaml` ≥ 1; `grep -c "Do not introduce a new module, package or layer for a one-file change"` = 1 em cada um dos três arquivos; nenhum bullet de `conventions:` sem condição ou fonte (leitura).
- R5 `limits:` de cada profile proíbe: instalar ou trocar versão de toolchain, SDK ou gerenciador (Gradle wrapper, JDK, `uv`, `poetry`, `pnpm`) sem autorização registrada no slice; supor ferramenta não observada (`ktlint`, `detekt`, `ruff`, `mypy`, `eslint`) e aplicar migração contra banco não descartável, espelhando `go-api.yaml:20-22` e `php.yaml`. Verificável: leitura dos três `limits:`.

Matching e validação

- R6 O desempate entre `typescript-web` e `typescript-api` é escrito em `profiles/README.md` (seção Matching): UI renderizada presente → `typescript-web`; serviço sem UI → `typescript-api`; monorepo com UI e API, uma regra só: `setup` rodado na raiz → `typescript-web`, com `notes:` no `project.yaml` nomeando a pasta da API e o profile `typescript-api`; `setup` rodado dentro da pasta da API → `typescript-api`. Verificável: texto presente; `/dh:setup` em modo proposta (sem escrita) num diretório temporário fora deste repositório, com um `package.json` sem UI, nomeia `typescript-api`; nada é commitado.
- R7 `profiles/README.md` perde a frase "Python, Electron and other stacks are out of v1", ganha os três arquivos na tabela com a data de entrada e o motivo (primeiro consumidor Kotlin/Python/TS backend: projetos de trabalho do dono), reescreve o passo 2 do Matching listando os sete profiles: `base`, `go-api`, `typescript-web`, `php`, `kotlin`, `python`, `typescript-api`, e inclui os três na frase da linha 55 ("never carry them"). Verificável: `grep -c "out of v1" profiles/README.md` = 0; tabela com oito linhas (sete profiles `.yaml` mais `examples/project.yaml`); `grep -c kotlin profiles/README.md` ≥ 3 (tabela, Matching, linha 55).
- R8 `dh validate` passa a exigir também `kotlin`, `python` e `typescript-api` na lista de profiles obrigatórios (`internal/kit/validate.go:343`; decisão 3: A, dono 2026-10-09), com as listas de fixture de `internal/kit/validate_test.go:328` e `:407` atualizadas. Como o PR toca `internal/`, exige a prova de clone limpo e `dist/*/bin` reconstruído (`AGENTS.md`). Verificável: teste Go que remove um dos três do inventário e vê o erro `missing profile`; `git clone . <tmp>` + `go run ./cmd/dh build` lá com sha256 dos binários igual ao da árvore.

Prova e limites

- R9 Nenhum agente, skill ou comando novo; nenhuma mudança em `.skills/project-onboarding/SKILL.md` além de, se necessário, um exemplo. O `setup` continua listando os profiles em tempo de execução. Verificável: `git diff --stat` só toca `profiles/`, `docs/`, `README*`, `CHANGELOG.md`, `AGENTS.md`, `internal/kit/validate.go`, `internal/kit/validate_test.go` e `dist/` regenerado por `go run ./cmd/dh build`; contagem de comandos segue 19 + 1.
- R10 Aceite final é do dono: `/dh:setup` em modo proposta (sem escrita) em um repositório real de cada stack na máquina do trabalho, com o profile certo nomeado e nenhum comando inventado. Até isso acontecer, o PRD fica "entregue (prova de campo pendente)" e a limitação consta no CHANGELOG. Verificável: relato do dono registrado no `MEMORY.md` pelo Coordenador; CHANGELOG com a frase.
- R11 Fora de escopo: Kotlin Android e Multiplatform; Electron; frameworks além dos cinco nomeados; fixture nova no kit; eval paga. Verificável: nada disso aparece no diff.

## 4. Docs

- `profiles/README.md` (tabela, nota de entrada, Matching com o desempate de R6)
- `README.md`, `README.en.md` (lista de profiles)
- `docs/tutorial.html`, `docs/en/tutorial.html` (seção de profiles, se citar a lista)
- `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md` (0.22.0 proposta; deriva barrada por `dh validate`)
- `AGENTS.md` (não lista profiles por nome; sem mudança prevista, confirmar no fim do lote)

## 5. Fora de escopo

- Fixture de teste por linguagem (`slice-01` continua a única).
- Profile por framework (`quarkus.yaml` etc.); só se a decisão 1 escolher B.
- Mudança no mod ou no dashboard; no binário, só a lista obrigatória de R8 (`validate.go`, `validate_test.go`) e o `dist/` regenerado.

## 6. Decisões do dono (decididas em 2026-10-09)

1. **Um profile por linguagem ou por framework?** A: por linguagem, framework por evidência (3 arquivos; decidida). B: por framework com `extends: kotlin`/`python` (7 arquivos; mais nítido por projeto, mais manutenção pública e matching em dois passos). C: por linguagem sem seção de framework (barato; não cumpre "padrões mais usados" por framework). Recomendação: A. Decidido: A (dono, 2026-10-09).
2. **TypeScript backend.** A: `typescript-api.yaml` próprio (decidida). B: seção de backend dentro de `typescript-web` (um arquivo, matching ambíguo entre UI e API). C: nenhum profile TypeScript novo; `typescript-web` basta (R3, R6 e R8 perdem a parte TS; vale se os projetos TypeScript do trabalho tiverem interface). Recomendação: A, salvo se a hipótese de §1 cair. Decidido: A (dono, 2026-10-09).
3. **Obrigatórios no `dh validate`.** A: sim (decidida; R8). B: não, ficam opcionais e a lista obrigatória segue com três. Recomendação: A. Decidido: A (dono, 2026-10-09).

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Regras de linguagem dentro de `base.yaml` | `base` vale para todo projeto; Kotlin não deve pesar num projeto Go. |
| Skill nova "kotlin-conventions" | O schema de profile já carrega convenções e limites; skill duplicaria o mecanismo. |
| Fixture por linguagem no kit | Custo de manutenção pública para provar YAML; a prova real é o `setup` no trabalho (R10). |

### Riscos e decisões pendentes

- Risco: convenções escritas sem repositório real na frente · mitigação: R2 (tudo `null` até haver evidência), R10 (prova de campo pelo dono).
- Risco: monorepo com UI e API em TypeScript cai no profile errado · mitigação: R6 (regra escrita e nota no `project.yaml`).
- Decisões pendentes: nenhuma. As três da seção 6 foram decididas pelo dono em 2026-10-09 (1A, 2A, 3A).

## Resumo executado

- **Entregue:** três profiles novos, `kotlin` (Spring Boot, Quarkus e Micronaut por evidência), `python` (Django e FastAPI) e `typescript-api`. O `dh validate` passa a exigi-los, e o README de profiles sai da regra "Python fora da v1". Saiu na 0.22.0.
- **Regras:**
  - R1 Honrada: `extends: base` e o esquema de 8 campos (QA T-1701-02).
  - R2 Honrada: `commands` null e "recorded only with evidence".
  - R3 Honrada: seções de framework condicionadas a pistas de build.
  - R4 Honrada: fonte de estilo e frase exata, uma vez por arquivo.
  - R5 Honrada: limites de toolchain, ferramentas e migrações.
  - R6 Honrada no texto (Matching com os sinais de UI e a regra de monorepo). A proposta do `/dh:setup` em pasta temporária ficou inconclusiva (haiku em plan mode) e entra com a R10 como prova de campo.
  - R7 Honrada. O critério de grep foi estreitado à frase "Python, Electron and other stacks are out of v1", porque a linha "`verify` is out of v1" é legítima e fica.
  - R8 Honrada: `TestCheckCountsRequiresNewProfiles`, mais a prova de clone limpo.
  - R9 Honrada: o diff antes da troca de versão ficou só em `profiles/` e `internal/kit/validate*.go`.
  - R10 Não honrada ainda. A prova de campo depende do dono num projeto do trabalho e ficou registrada como pendente na CHANGELOG, por decisão do dono de seguir até o fim (2026-10-09).
  - R11 Honrada: nada de Android, Multiplatform, Electron, fixture ou eval.
- **Tickets:**
  - T-1701-01 backend (harness-maintainer): aprovado (revisão claude, polimento aplicado).
  - T-1701-02 teste: verde, com a R6 parcial.
- **Docs:** `profiles/README.md`, `CHANGELOG.md`, `README.md`, `README.en.md`, `docs/roadmap.html`, `docs/en/roadmap.html`.
- **Fora:** Android e Multiplatform, Electron, fixtures, evals pagas.
