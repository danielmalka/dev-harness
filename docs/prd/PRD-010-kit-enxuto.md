# PRD-010 · Kit enxuto: menos verificação, menos agentes, regras só vigentes

| Campo | Valor |
|---|---|
| Status | aprovado (2026-09-30, pelo Coordenador por delegação do dono) |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-09-30 / 2026-09-30 |
| Tickets | a derivar |

## 1. Problema

O kit cresceu por acréscimo e o dono paga o custo: "ficamos dias a dias rodando testes e testes que não sei o quanto
foram produtivos". Fatos medidos no repositório (30/09/2026):

- 23 casos de eval em `evals/cases/`, 19 com `append_system_prompt` copiado do corpo de um agente. Uma checagem do
  `dh validate` (`checkEmbeddedAgentBodies`, `internal/kit/evals.go`, 589 linhas) reprova qualquer deriva dessa cópia
  ("append_system_prompt no longer matches"), e o conserto é um script fora do repositório
  (`/tmp/claude-1000/resync.py`). Toda edição de agente vira um lote de ressincronização.
- 19 agentes de 103 a 129 linhas, cinco deles construtores com o mesmo esqueleto (`backend-builder`,
  `frontend-builder`, `data-engineer`, `devops-engineer`, `refactorer`), dois revisores (`code-reviewer`,
  `security-reviewer`) e dois de desenho (`api-designer`, `solution-architect`). Três comandos repetem o texto do
  agente e da skill que despacham (`understand`/`repo-scout`/`repository-mapping`, `release`/`release-manager`/
  `delivery-readiness`, `verify`/`qa-verifier`/`regression-testing`).
- `AGENTS.md` (49 linhas) e `.agents/coordinator.md` (223 linhas) misturam regra vigente com data, motivo e histórico.
- A prova de clone limpo (sha256) é exigida em todo PR, mesmo os que não tocam Go nem binários.
- Só existe `docs/inicio-rapido.html` (mais o tutorial 00), sem um guia completo de uso.

Decisão do dono (`AGENTS.md`, "Kit enxuto", 30/09/2026): simplificar processo e verificação sem perder capacidade;
o kit não vira um monstro. Palavras dele: "não quero que vc tire coisas, talvez um merge conforme sugerido já resolva
bastante coisa."

## 2. Solução

O kit faz o mesmo, com menos peças e menos ritual. Os 19 comandos `/dh:` continuam. Evals saem do fluxo normal e
ficam em 5 casos essenciais rodados só a pedido. 19 agentes viram 13, com a especialidade movida para as skills e
para o despacho. `AGENTS.md` e o Coordenador passam a dizer só o que vale hoje; a história vai para um ADR. A verificação
padrão passa a ser checagens determinísticas mais um revisor Claude (proposta R11, decidida em 2026-09-30). Há um tutorial completo de uso em pt-br e en. O PRD-009
fecha seu pendente. Sai como 0.14.0.

Agentes depois do merge (13): `coordinator`, `product-discovery`, `implementation-planner`, `architect`, `builder`,
`reviewer`, `debugger`, `document-validator`, `qa-verifier`, `docs-guide`, `release-manager`, `repo-scout`,
`harness-maintainer`. O `debugger` continua separado (decisão de 2026-09-30, "não quero que vc tire coisas"); fundir
mais agentes está fora de escopo.

Casos de eval mantidos (5), escolhidos por medir comportamento que o kit não pode perder e que já tem baseline em
`evals/baselines/`:

1. `security-reviewer-object-authz-miss` — o revisor de segurança acha a falha de autorização por objeto sem acusar o
   código correto (baselines `2026-09-24/25-security-reviewer-*`).
2. `code-reviewer-equivalent-form-miss` — o revisor de código pega a forma equivalente que o classificador não trata
   (baselines `2026-09-24/25-code-reviewer-*`). Junto com o 1, cobrem os dois modos do `reviewer`.
3. `external-clis-risk-001` — a cláusula anti-delegação dentro do prompt montado vincula o revisor CLI (RISK-001, o
   risco grave registrado do kit).
4. `doc-validator-gap` — o validador de documento reprova PRD com lacuna e critério não verificável
   (baseline `2026-09-24/25-document-validator-*`); o validador continua para PRD e plano.
5. `auto-hard-stop-recorded-state` — o `/dh:auto` para no ponto duro sem laço nem repetição e registra o estado
   igual em BRIEF e MEMORY (baseline `2026-09-28-auto-hard-stop-*`); é a salvaguarda do modo autônomo.

Os outros 18 casos saem do fluxo: vão para `evals/archive/`, fora do `dh build` e de qualquer checagem do `validate`, rodáveis a pedido.
Baselines antigos ficam como histórico.

## 3. Regras

- R1 A capacidade do kit não muda: os 19 comandos de `.commands/` continuam, `plan-loop`, CLIs externas com clone,
  validador de documento e docs em pt-br e en continuam, e cada capacidade dos agentes fundidos (R5 a R7) continua
  alcançável por comando ou despacho. Verificável: `ls .commands | wc -l` dá 19; `go run ./cmd/dh validate .` sai
  limpo; toda decisão do `AGENTS.md` anterior está no novo ou no ADR-003 (R9).
- R2 `evals/cases/` contém exatamente os 5 casos da Solução; os outros 18 estão em `evals/archive/`, que entra na
  exclusão de primeiro nível do `dh build` e de `evalsExcluded`. A checagem de cópia (`checkEmbeddedAgentBodies`,
  `checkFixtureAgentCopies` e helpers só usados por elas) é removida com seus testes; a rejeição de `fixtures/` como
  link simbólico, hoje dentro de `checkFixtureAgentCopies`, é preservada em `checkEvalReferences` com seu teste.
  `checkAntiDelegationClause` fica. `resync.py` deixa de ser necessário e nenhum documento vigente o cita.
  Verificável: `ls evals/cases | wc -l` dá 5; `ls dist/*/evals/archive` não existe; `grep -rn "no longer matches"
  internal/` vazio; `grep -rn resync AGENTS.md .agents .skills .commands docs --exclude-dir=prd --exclude-dir=adr`
  vazio; `go test ./...` verde.
- R3 Nenhum caso mantido depende de cópia literal de corpo de agente, inclusive em `fixtures/` (o
  `external-clis-risk-001` passa a usar `.agents/reviewer.md`; os casos que embutem o Coordenador seguem a mesma
  abordagem). Hipótese (`evals.go` diz que o runner não aceita arquivo em `append_system_prompt`; não confirmado contra
  `claude plugin eval`): o caso nomeia o agente (`dh:reviewer`), ou o `scaffold.sh` copia `.agents/<papel>.md` para a
  fixture na hora da execução, sem versionar a cópia. Nenhuma eval paga roda para confirmar; se a hipótese cair, os casos
  mantidos ficam marcados `não rodável` até nova decisão do dono, sem voltar a cópia versionada. Verificável:
  `grep -rlnE '\.agents/code-reviewer|fixtures/code-reviewer\.md' evals/cases` vazio; `evals/cases/external-clis-risk-001/fixtures/code-reviewer.md` não existe; nenhum arquivo em `evals/cases/*/fixtures/` é corpo de agente; nenhum
  `append_system_prompt` de mais de 5 linhas.
- R4 Eval paga só roda por pedido explícito do dono ou dentro de `/dh:improve`, que o dono invoca; nenhum outro comando,
  fluxo do Coordenador ou CI a dispara; a resposta do dono ao item de teto de eval do questionário do `/dh:auto` conta
  como pedido explícito (emenda de 2026-09-30, por delegação do dono); teto zero ou ausente = nenhuma eval. O job `evals` continua pulando sem chave. Verificável: `grep -rn "claude plugin
  eval" .commands .agents` só encontra `improve.md`; o texto do Coordenador diz "only when the owner asks".
- R5 Existe `.agents/builder.md`, que recebe no despacho a lane (`backend`, `frontend`, `dados`, `infra`) e o modo
  (`build` ou `refactor`). Substitui `backend-builder`, `frontend-builder`, `data-engineer`, `devops-engineer` e
  `refactorer`; `debugger` fica. Roteamento de `.commands/build.md` preservado: lane vira `builder` com essa lane; sem
  lane, o fallback atual por tipo de trabalho escolhe a lane (`infra` fora do fallback); toda lane carrega
  `incremental-implementation`, `dados` também `data-migrations`; modo `refactor` carrega `safe-refactoring`; lanes de
  teste seguem só com `qa-verifier`. Preserva: só implementa fatia autorizada; não inventa decisão de produto; não
  publica nem aprova a própria entrega; não aplica mudança destrutiva em dado real nem provisiona infraestrutura;
  refactor não muda comportamento; relata checagens não rodadas. `MinAgents` em `internal/kit/validate.go` vai para 13
  (e seu teste). `profiles/*.yaml` passam a citar `builder` (lane `dados`) e `architect`. Verificável: os 5 arquivos
  antigos não existem; `grep -rnE 'api-designer|data-engineer|solution-architect|backend-builder|frontend-builder|
  devops-engineer|refactorer' profiles/ .commands .agents` vazio; `validate` limpo.
- R6 Existe `.agents/reviewer.md`, com modo `code` ou `security` no despacho (skill `code-review` ou `security-review`),
  somente leitura (mesma lista de ferramentas dos revisores atuais). Substitui `code-reviewer` e `security-reviewer`.
  Modelo explícito no despacho do Coordenador, `/dh:secure`, `/dh:build` e `/dh:review`: `security` com `opus`, `code`
  com `sonnet`. Preserva: achados com severidade, local, cenário e impacto; não corrige; a cláusula anti-delegação
  continua só no Coordenador e em `external-clis` e é anexada ao prompt do revisor CLI na montagem (`reviewer.md` não a
  copia); o prompt do revisor CLI é montado do corpo de `reviewer.md`, da skill do estágio e do contexto do despacho;
  `/dh:secure` despacha `reviewer` em modo `security`; `reviewers:` do `project.yaml` não muda. Verificável:
  `checkAntiDelegationClause` verde; `grep -n '## External CLI reviewers' .agents/coordinator.md` encontra a seção e o
  bloco cercado; `secure.md` cita `reviewer` e `opus`.
- R7 Existe `.agents/architect.md` que substitui `api-designer` e `solution-architect`, com a skill `api-contracts`
  (`sonnet`) ou `architecture-decisions` (`opus`) conforme o pedido, modelo explícito no despacho. Preserva: contrato
  público com operações, erros, autenticação e compatibilidade; decisão de fronteira com opções proporcionais; não
  redesenha por correção local. Verificável: os dois arquivos antigos não existem; `grep` dos nomes antigos vazio fora
  de história; `validate` verde.
- R8 Sem remover comando, a sobreposição atrás de três comandos deixa uma casa só para cada texto: o procedimento fica
  na skill (`repository-mapping`, `delivery-readiness`, `regression-testing`), o agente (`repo-scout`,
  `release-manager`, `qa-verifier`) fica com papel, limites e formato de saída, e o comando (`understand`, `release`,
  `verify`) fica com o despacho e o critério de parada. Verificável: a soma de linhas dos 9 arquivos cai, e cada uma das
  frases "Do not use for a whole-repo audit", "Do not publish because you were asked to prepare" e "You may write
  tests; you do not patch production code" (ou o equivalente nos textos atuais, fixado no ticket) aparece em um só
  arquivo; os 3 comandos seguem em `.commands/`.
- R9 `AGENTS.md` fica com cerca de 15 regras vigentes, cada uma em 1 a 3 linhas, sem data nem "porquê" (datas de
  revisão pendentes, como "revisitar em 21/10/2026", ficam na regra vigente). Toda regra removida ou reescrita, com sua
  data e seu motivo, entra em `docs/adr/ADR-003-kit-enxuto.md`; o EPOCHAL, escrito pelo Coordenador, é cópia adicional,
  nunca o único destino. Nenhuma regra fixada pelo dono se perde: só a história muda de lugar. O `CLAUDE.md` continua
  importando `@AGENTS.md`. Verificável: `grep -c '^- ' AGENTS.md` menor ou igual a 18; toda decisão do `AGENTS.md`
  anterior (`git show HEAD:AGENTS.md`) aparece no novo `AGENTS.md` ou no ADR-003.
- R10 `.agents/coordinator.md` fica só com as regras em vigor, com no máximo 150 linhas (hoje 223), e cita os agentes
  fundidos e o despacho por lane e modo. A rotina de memória e o roteamento por lane continuam, e a seção
  `## External CLI reviewers` com o bloco cercado da cláusula anti-delegação fica onde está. Verificável: `wc -l` menor
  ou igual a 150; `grep -n '## External CLI reviewers' .agents/coordinator.md`; `validate` e `go test ./...` verdes.
- R11 (decidido em 2026-09-30, por delegação do dono) Padrão de verificação: uma mudança recebe as checagens
  determinísticas (`dh validate`, build, `go test`) mais 1 revisor Claude. O painel de 3 só entra quando o dono pede ou
  a mudança toca segurança. Não se criam casos de eval nem sondas novas sem pedido do dono; as sondas do PRD-008
  (`mcode`, `agy` H3) continuam autorizadas. O teto de rodadas de correção continua, mas a rodada termina quando não há
  achado bloqueante. A exceção do `/dh:plan-loop` (onda julgada pela lista completa de `reviewers.document`) continua.
  Verificável: `AGENTS.md` traz a regra; o Coordenador diz "panel of three" só com essas duas
  condições.
- R12 A prova de clone limpo com sha256 só é exigida quando o PR toca `cmd/`, `internal/`, `go.mod` ou
  `dist/*/bin`. PR só de texto do kit e docs segue com `dh validate` e `go run ./cmd/dh build` sem diferença em `dist/`.
  Verificável: a regra em `AGENTS.md` e em `implementation-planning` cita esses quatro caminhos.
- R13 Existe um tutorial completo de uso em `docs/tutorial.html` (pt-br) e `docs/en/tutorial.html` (en), com link
  cruzado, seguindo `doc-template-html`, com as seções: instalação; primeira sessão (`setup`, `doctor`); o que vive em
  `.harness/`; fluxos (feature de `discover` a `release`, `/dh:auto`, bug com `fix`, `refactor`, `verify`, sessão longa
  com `handoff`, `resume` e `consolidate-memory`, cuidado do kit com `improve`); configuração (`project.yaml`:
  `language`, `commands`, `reviewers:` por etapa; `local.yaml`); avançado (`plan-loop`, CLIs externas com isolamento por
  clone, evals a pedido); tabela de referência dos 19 comandos (quando usar, o que produz, quem aprova). As quatro
  páginas antigas (`docs/inicio-rapido.html`, `docs/en/quick-start.html`, `docs/tutoriais/00-primeira-maquina.html`,
  `docs/en/tutorials/00-first-machine.html`) não são apagadas: cada uma vira uma página curta que aponta para o
  tutorial novo. `README.md` e `README.en.md` apontam para o tutorial. Verificável: a tabela de referência
  (`<table id="comandos">`) tem 19 `<tr>` de corpo, um por arquivo de `.commands/`, conferido por nome, em cada língua;
  as 4 páginas antigas têm menos de 40 linhas e um link para o tutorial; `validate` verde.
- R14 O `docs/prd/PRD-009-roadmap-vivo.md` ganha a seção "Resumo executado" com os campos do RF-14 do PRD-007
  (Entregue, Regras, Tickets, Docs, Fora), citando a 0.13.0 e o PR #21, e o Status passa a "entregue em 2026-09-30".
  Verificável: `grep` do PRD-009 encontra a seção e o Status novo.
- R15 Release 0.14.0: versão em `.claude-plugin/marketplace.json` e nos demais pontos de versão, `CHANGELOG.md` com
  `## 0.14.0 — <data>`, card 0.14.0 no topo de Entregue nos dois roadmaps, chips, linha Fontes e "Contagem atual"
  atualizadas (13 agentes, 19 comandos), `dist/` regenerado por `go run ./cmd/dh build`. Como toca `internal/`, vale a
  prova de clone limpo do R12. Verificável: `go run ./cmd/dh validate .` sai limpo (inclui a checagem do roadmap da
  0.13.0); sha256 dos binários do clone limpo igual aos da árvore.

## 4. Docs

- AGENTS.md
- .agents/coordinator.md
- docs/adr/ADR-003-kit-enxuto.md
- docs/tutorial.html
- docs/en/tutorial.html
- docs/roadmap.html
- docs/en/roadmap.html
- README.md
- README.en.md
- CHANGELOG.md
- docs/prd/PRD-009-roadmap-vivo.md
- docs/inicio-rapido.html, docs/en/quick-start.html, docs/tutoriais/00-primeira-maquina.html, docs/en/tutorials/00-first-machine.html (viram páginas curtas que apontam para o tutorial)

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Cortar comandos ou `plan-loop` | O dono vetou: nenhuma capacidade sai. |
| Manter a checagem de cópia e só automatizar o `resync.py` | Mantém o custo de 589 linhas e um script fora do repositório; o problema é a cópia, não o conserto. |
| Apagar os 18 casos excedentes | Perde histórico medido; arquivar custa uma pasta e nenhuma checagem. |
| Fundir também `debugger`, `repo-scout`, `release-manager`, `qa-verifier` | O dono não quer tirar coisas; papéis distintos ficam. |
| Apagar os 4 docs antigos | O dono não quer tirar coisas; viram páginas curtas de link (R13). |

### Riscos e decisões pendentes

- Risco: agente fundido perde uma regra sem ninguém notar · mitigação: R5 a R7 listam o que preservar e o revisor do PR
  confere contra os corpos antigos (`git show HEAD:.agents/<nome>.md`).
- Hipótese (R3): o runner de `claude plugin eval` só aceita texto em `append_system_prompt` · falsa se um caso que nomeia o
  agente ou copia na hora do `scaffold.sh` rodar e medir o agente atual (só se o dono pedir a eval).
- Decisão pendente: nenhuma. R11 confirmado por delegação do dono em 2026-09-30; sem eval paga, o R3 cai em `não rodável` se a hipótese falhar.
