# Glossário fixo da tradução para pt-BR

Este glossário vale para todo texto do kit traduzido do inglês (ADR-011): skills, agentes, comandos, perfis, docs e o histórico. Um termo definido aqui tem **uma** tradução, sem sinônimo, em todos os arquivos. Termo novo entra aqui antes de entrar em qualquer arquivo.

## 1. Regras de tradução

1. **Fiel.** A tradução não muda regra nem comportamento. Toda frase normativa mantém a força (deve, nunca, só, no máximo N), os números, os limites, as proibições e as exceções. Nada é resumido, acrescentado ou reordenado de forma que mude o sentido.
2. **Mesma estrutura.** Mesmo número de títulos, itens de lista, linhas de tabela, blocos de código e passos numerados, na mesma ordem. Isso deixa o diff conferível lado a lado com o original.
3. **Identificadores ficam em inglês e entre crases.** Slugs de skill, nomes de agente (`coordinator`, `builder`), comandos (`/dh:build`, `dh validate`), subcomandos e flags, nomes de arquivo e de pasta (`MEMORY.md`, `TASK.md`, `.harness/`), chaves de YAML, JSON e frontmatter, variáveis de ambiente, código e saída literal de ferramentas. Nome técnico citado no meio do texto vai entre crases.
4. **Valores de enum ficam em inglês e entre crases** (lista na seção 4): outras skills, agentes e avaliações os leem literalmente.
5. **Títulos, rótulos de relatório e `description` são pt-BR** (seção 3). Os modelos de saída (blocos de código com o formato do relatório) têm os rótulos e os marcadores `<...>` traduzidos; os valores de enum dentro deles ficam como estão.
6. **Frontmatter:** só o valor de `description` é traduzido. `name`, `author` e as chaves e valores de `metadata` (inclusive `sources`, que são nomes de origem) ficam como estão.
7. **Imperativo do "você" sem o pronome** ("Rode", "Leia", "Pergunte"), igual aos docs em pt-BR do repositório. Frase curta, voz ativa, um termo por conceito.
8. **Papéis em prosa:** "o Coordenador" (o agente é `coordinator`), "o dono", "o revisor" (o agente é `reviewer`), "o especialista". Quando o texto cita o papel como identificador (um valor, uma chave, um `subagent_type`), vale o nome em inglês entre crases.
9. **Exemplos:** traduzem-se a prosa e os dados ilustrativos; ficam em inglês o código, os comandos, os nomes de tabela, campo e função, e a saída literal de ferramentas.
10. **Nada em inglês fora dos identificadores.** Se uma palavra em inglês não é identificador nem enum, ela tem tradução aqui ou entra aqui.

## 2. Termos

| Inglês | pt-BR | Observação |
|---|---|---|
| owner | dono | o dono do projeto |
| Coordinator | Coordenador | o agente é `coordinator` |
| specialist | especialista | |
| dispatch (substantivo) | despacho | "despachar" é o verbo |
| dispatched | despachado | |
| main session | sessão principal | |
| reviewer | revisor | o agente é `reviewer` |
| author | autor | |
| validator | validador | o agente é `document-validator` |
| round | rodada | |
| slice | fatia | |
| stage | etapa | |
| batch | lote | |
| wave | onda | |
| ticket | ticket | já usado em pt-BR |
| task | tarefa | |
| plan | plano | |
| brief | resumo | o resumo de requisitos da skill `requirements-discovery` |
| requirement | requisito | |
| acceptance criteria | critérios de aceite | |
| acceptance | aceite | |
| scope | escopo | |
| write set | conjunto de escrita | os arquivos que a tarefa pode escrever |
| authorization | autorização | |
| authorized | autorizado | |
| finding | achado | |
| blocking finding | achado bloqueante | |
| non-blocking | não bloqueante | |
| verdict | veredito | |
| verdict signal | sinal de veredito | a linha com o prefixo que marca o veredito |
| evidence | evidência | |
| unverified | não verificado | quando é rótulo de relatório; o valor de enum `unverified` fica |
| claim | afirmação | |
| assumption | premissa | |
| hypothesis | hipótese | |
| open decision / decision needed | decisão em aberto / decisão necessária | |
| risk | risco | |
| limitation | limitação | |
| rollback | reversão | |
| revert | reverter | |
| rollout | implantação | |
| release | release | já usado em pt-BR |
| release notes | notas de release | |
| compatibility | compatibilidade | |
| regression | regressão | |
| regression test | teste de regressão | |
| characterization test | teste de caracterização | |
| root cause | causa raiz | |
| symptom | sintoma | |
| reproduction | reprodução | |
| failure | falha | |
| fail-open | falha aberta | quando é conceito; o texto explica |
| not-run | não rodado | rótulo de relatório; o valor de enum `not-run` fica |
| call (de CLI, de ferramenta) | chamada | |
| call site | ponto de chamada | |
| caller / callee | chamador / chamado | |
| consumer | consumidor | |
| contract | contrato | |
| interface | interface | |
| entry point | ponto de entrada | |
| trust boundary | fronteira de confiança | |
| untrusted input | entrada não confiável | |
| blast radius | raio de impacto | |
| flow | fluxo | |
| map | mapa | |
| guess | palpite | |
| unknown (substantivo) | desconhecido | |
| live tree | árvore viva | a árvore de trabalho real do projeto |
| disposable clone | clone descartável | |
| clone isolation | isolamento por clone | |
| freeze (substantivo) | congelamento | "congelar" é o verbo; o diretório é `freeze directory` = diretório de congelamento |
| worktree | worktree | já usado em pt-BR |
| linked worktree | worktree vinculado | |
| working directory | diretório de trabalho | |
| harness dir | harness dir | já usado em pt-BR; em crases quando é o nome do conceito do kit |
| record (substantivo) | registro | |
| memory | memória | |
| handoff | handoff | já usado em pt-BR; o comando é `/dh:handoff` |
| resume (verbo) | retomar | o comando é `/dh:resume` |
| consolidation | consolidação | |
| archive | arquivar | |
| trim | podar | |
| write guard | trava de escrita | |
| hash guard | trava de hash | |
| lock | trava | quando é arquivo de trava, `lock` em crases |
| read-only | somente leitura | |
| read-only flag | flag somente leitura | |
| prompt | prompt | já usado em pt-BR |
| prompt body | corpo do prompt | |
| prompt assembly | montagem do prompt | |
| clause | cláusula | |
| anti-delegation clause | cláusula anti-delegação | |
| binary | binário | |
| model | modelo | |
| timeout | timeout | já usado em pt-BR |
| in-flight doubt | dúvida no meio do build | como no `AGENTS.md` |
| proof case | caso de prova | |
| quick reference | referência rápida | |
| common mistakes | erros comuns | |
| rationalizations | racionalizações | |
| red flags | sinais de alerta | |
| gotchas | pegadinhas | |
| drift | deriva | |
| drift report | relatório de deriva | |
| neighbor | vizinho | |
| criterion (de avaliação) | critério | |
| non-discriminating | que não discrimina | |
| baseline | linha de base | |
| asset | ativo | um arquivo-fonte do kit (skill, agente, comando, perfil, template) |
| provenance | procedência | |
| conventions | convenções | |
| existing instructions | instruções existentes | |
| inventory | inventário | |
| capability impact | impacto na capacidade | |
| mode | modo | a chave `mode` fica |
| profile | perfil | |
| template | template | já usado em pt-BR |
| skeleton | esqueleto | |
| shell (de HTML) | casca | |
| stamp (verbo, do `stamp.sh`) | carimbar | o script é `stamp.sh` |
| matrix | matriz | |
| acceptance matrix | matriz de aceite | |
| coverage gap | lacuna de cobertura | |
| coverage | cobertura | |
| disposable database | banco descartável | |
| migration | migração | |
| phase | fase | |
| backfill | preenchimento retroativo | |
| decision (ADR) | decisão | |
| option | opção | |
| choice | escolha | |
| consequences | consequências | |
| error catalog | catálogo de erros | |
| operation | operação | |
| idempotent | idempotente | |
| severity | severidade | os valores (`Critical`, `Major`, `Minor`) ficam |
| step | passo | |
| check (substantivo) | checagem | |
| gate | gate | já usado em pt-BR; "quality gate" também |
| quality gate | quality gate | |
| commit | commit | já usado em pt-BR |
| pull request | pull request | já usado em pt-BR |
| diff | diff | já usado em pt-BR |
| branch | branch | já usado em pt-BR |
| tag | tag | já usado em pt-BR |
| shell command | comando de shell | |
| tool | ferramenta | |
| reply cap | limite de resposta | |
| owner's words | palavras do dono | só quando é rótulo; o valor `owner's words` fica |
| context bundle | pacote de contexto | |
| source material | material-fonte | |
| document | documento | |
| section | seção | |
| heading | título | de seção |
| line prefix | prefixo de linha | |

## 3. Rótulos fixos de relatório e de skill

Estes rótulos aparecem em quase toda skill e têm uma tradução só.

| Inglês | pt-BR |
|---|---|
| Overview | Visão geral |
| When to use | Quando usar |
| When not to use | Quando não usar |
| Inputs | Entradas |
| Procedure | Procedimento |
| Output format | Formato de saída |
| Quick reference | Referência rápida |
| Common mistakes | Erros comuns |
| Rationalizations / Common rationalizations | Racionalizações |
| Red flags | Sinais de alerta |
| Example | Exemplo |
| Related | Relacionado |
| Proof case | Caso de prova |
| Verdict | Veredito |
| Findings | Achados |
| Evidence | Evidência |
| Scope | Escopo |
| Checks | Checagens |
| Limitations | Limitações |
| Next step | Próximo passo |
| Status | Status |
| Summary | Resumo |
| Documentation | Documentação |
| Comprehension | Compreensão |
| Not raised | Não levantado |
| Conclusion | Conclusão |
| Files changed | Arquivos alterados |
| Open decisions | Decisões em aberto |
| Decisions needed | Decisões necessárias |
| Assumptions | Premissas |
| Self-review | Autorrevisão |
| Re-entry | Reentrada |
| In-flight doubt | Dúvida no meio do build |
| Review status | Status da revisão |
| Reported by | Reportado por |
| Blocking findings | Achados bloqueantes |
| Guide rule | Regra do guia |
| Review guide | Guia de revisão |
| Ignored by rule | Ignorado por regra |
| Spec compliance | Aderência à especificação |
| Correctness and regressions | Correção e regressões |
| Maintainability notes | Notas de manutenibilidade |
| Reported, not fixed | Reportado, não corrigido |
| Residual risks | Riscos residuais |

## 4. Valores que ficam em inglês (enum, entre crases)

Outras skills, agentes e avaliações os casam literalmente; só mudam junto com quem os lê.

- Veredito e status de revisão: `approved`, `approve`, `request changes`, `changes required`, `incomplete`, `blocked`, `rejected`, `pending`.
- Severidade: `Critical`, `Major`, `Minor`, `Cosmetic`, `HIGH`, `MEDIUM`, `LOW`, `INFO`.
- Estado de achado: `open`, `closed`, `escalated`.
- Resultado de checagem: `passed`, `failed`, `not-run`, `reported (unverified)`, `unverified`.
- Status de conclusão de um especialista e de entrega: `ready for QA/review`, `ready for independent QA/review`, `partial`, `blocked`, `inconclusive`, `diagnosed, no fix authorized`, `fixed, pending independent QA and review`, `ready for decision`, `not ready`.
- Estado de ticket e de execução: `ready`, `in progress`, `blocked`, `done`, `running`, `completed`, `cancelled`.
- Valores de configuração: `repo`, `global`, `none`, `always`, `on-request`, `claude`, `cli:<binário>/<slug>`.
- Classes de falha de CLI externa (`not_installed`, `isolation_unavailable` etc.) e qualquer valor que a skill defina como identificador com sublinhado ou hífen.

## 5. `description` das skills

- Formato: "Use quando …" para o que dispara a skill, "Use também quando …" para o segundo gatilho e "Não use quando …" para o que a exclui. A frase é a do original, traduzida, sem acrescentar nem tirar gatilho.
- Começa sempre com "Use quando" (o `dh validate` exige a partir da 0.45.0).
- Uma linha, sem quebra, entre aspas só se o YAML exigir.

## 6. Idioma de commit e de comentário de código

Vale o que o projeto define (`AGENTS.md`, `CLAUDE.md`, `README.md`, depois `conventions.code_language` do `project.yaml`); sem definição, inglês (ADR-011). Este glossário não muda essa regra.

## 7. Rótulos lidos por outros arquivos: escrita em pt-BR, leitura nas duas formas até a 0.46.0

Os rótulos e títulos abaixo são escritos em pt-BR pelas skills e pelos agentes que produzem o formato. Quem **lê** um deles (comandos, o `coordinator`, a `external-clis`) aceita também a forma em inglês, para não quebrar `PLAN.md` e relatórios escritos antes da 0.43.0. A leitura nas duas formas vale até a 0.46.0 (fatia S5, que trava o inglês); depois dela, só pt-BR.

| pt-BR (escrita) | Inglês (também lido) | Quem lê | Quem escreve |
|---|---|---|---|
| `## Decisões necessárias` | `## Decisions needed` | `.commands/plan.md` (página de aprovação) | `implementation-planning`, `implementation-planner` |
| `### Fatia <n> (<id>)` | `### Slice <n> (<id>)` | `.commands/review.md` | `implementation-planning`, `implementation-planner` |
| `## Plano` | `## Plan` | `external-clis` (teste de falha de transporte do planejador) | `implementation-planning`, `implementation-planner` |
| `Status da revisão:` | `Review status:` | `coordinator`, `external-clis` | `code-review`, `reviewer` |
| `## Achados` | `## Findings` | `coordinator`, `external-clis` | `code-review`, `security-review`, `document-review`, `reviewer`, `document-validator` |
| `### <SEVERIDADE>: <título>` | `### <SEVERITY>: <title>` | `coordinator`, `external-clis` | `security-review`, `reviewer` |
| `Reportado por:` | `Reported by:` | `coordinator` | `code-review`, `security-review`, `document-review` |
| `## Compreensão` | `## Comprehension` | `coordinator` | `document-review`, `document-validator` |
| `## Não levantado` | `## Not raised` | `coordinator` | `document-review`, `document-validator` |
| `Regra do guia:` | `Guide rule:` | `.commands/review.md` | `code-review`, `reviewer` |
| `cobertura incompleta: <eixo>` | `coverage incomplete: <axis>` | `coordinator` | `code-review`, `security-review`, `document-review`, `reviewer`, `document-validator` |

O pacote de contexto da etapa de documento tem **quatro** itens (documento, material-fonte, relatório anterior, harness dir). O `coordinator` e o `document-validator` não citam mais o item "idioma do projeto".

O contrato desta tabela é conferido por `internal/kit/rotulos_test.go`.
