# PRD-011 · Trava de runtime por mod: só o Coordenador despacha, só o Coordenador escreve a memória

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-05 |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-10-05 / 2026-10-05 |
| Tickets | T-1101, T-1102 |

## 1. Problema

Duas regras críticas do kit vivem só em prosa e dependem de o modelo obedecer:

- "Só o Coordenador despacha". O RISK-001 (`.harness/RISKS.md`) registra um especialista que despachou outro agente apesar
  de `disallowedTools: [Agent]`: o runtime não honra o campo em agente de plugin. O dono removeu o campo (falsa segurança)
  e o estado do risco segue "mitigado", sem trava.
- "Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`" (`AGENTS.md`). Nada impede um subagente de
  chamar `Write` ou `Edit` nesses arquivos.

Fato (verificado pelo Coordenador em 2026-10-05, Claude Code 2.1.289): `claude plugin validate` aceita um `hooks.json` com
`hooks` clássicos e `modules`; o evento `agent.spawn` traz `parentAgentId` (ausente no laço principal) e aceita
`{ deny: reason }`; o evento `tool.call` traz `agentId` (ausente no laço principal). `claude plugin test <pasta>` existe
(`claude plugin test --help`: roda todo `*.test.ts` e `*.test.tsx` sob a pasta, importa o kit de `claude-code/testing`, sai
com 1 se um teste falha) e um teste de fumaça passou numa sonda em scratchpad (verificação do Coordenador, 2026-10-05).
Logo o runtime pode negar essas duas ações de forma determinística.

## 2. Solução

O plugin do kit passa a trazer um mod (TypeScript, roda no motor do Claude Code, sem Node) com duas travas, em todo
projeto consumidor:

- Subagente que tenta despachar outro agente tem o `agent.spawn` negado, com motivo legível ("only the Coordinator
  dispatches").
- Subagente que tenta `Write`/`Edit` em `.harness/MEMORY.md`, `.harness/EPOCHAL.md` ou `.harness/RISKS.md` tem a chamada
  negada, com motivo legível. O laço principal (Coordenador) não é afetado.

O `hooks/hooks.json` ganha `"modules": ["./<arquivo>.ts"]` ao lado dos `hooks` clássicos do `dh snapshot`, que não mudam.
Lacuna aceita pelo dono: escrita por `Bash` não é bloqueada; a regra em prosa continua valendo. Adaptadores que não são
Claude Code mantêm só a regra em prosa.

## 3. Regras

- R1 O mod fica em `adapters/claude-code/plugin/hooks/` e é referenciado por `modules` em `hooks/hooks.json`; os `hooks`
  clássicos de `dh snapshot` ficam byte a byte iguais. Verificável: `git diff` de `hooks.json` só adiciona a chave
  `modules`; `claude plugin validate` no pacote construído passa.
- R2 O mod exporta `register(on, options)` e registra exatamente dois handlers: `agent.spawn` e `tool.call`. Nenhum outro
  evento, painel, status line ou UI. Verificável: leitura do arquivo; os testes de R6.
- R3 `agent.spawn`: quando `e.parentAgentId` está definido, retorna `{ deny: <motivo> }`; quando ausente, não nega.
  Verificável: teste com os dois casos.
- R4 `tool.call` das ferramentas `Write` e `Edit` (só elas): quando a chamada traz `agentId` e o caminho alvo, normalizado
  (separadores `/` e `\`), termina em `.harness/MEMORY.md`, `.harness/EPOCHAL.md` ou `.harness/RISKS.md`, qualquer que seja
  a raiz (relativo, absoluto, ou de outro checkout: negar também é aceito), retorna `{ deny: <motivo> }`. Sem `agentId`
  (Coordenador) ou com outro caminho, não nega. Verificável: testes para cada um dos 3 arquivos, caminho relativo, absoluto
  e com `\`, com e sem `agentId`, um caminho vizinho (`.harness/project.yaml`) que passa e um caminho fora do projeto
  terminando em `.harness/MEMORY.md` que nega.
- R5 O mod falha aberto: erro inesperado dentro do handler (campo ausente, caminho não string) não nega nem derruba a
  sessão. Verificável: teste com entrada malformada que não nega.
- R6 O pacote construído traz o mod e um teste `*.test.ts` com os casos de R3 a R5. `go run ./cmd/dh build` copia ambos sem
  mudança de Go (fato: `internal/build/build.go` copia `adapters/claude-code/plugin` inteiro com `copyTree` e confere com
  `compareCopiedTree`). Verificável: os arquivos existem em `dist/claude-code/dev-harness/hooks/`; `claude plugin test`
  na pasta do pacote passa.
- R7 Nenhuma mudança em Go. `dh validate` continua verde com os `.ts` presentes. Limite registrado: `iterTextFiles`
  (`internal/kit/validate.go`) só varre `.md .yaml .yml .html .py .sh .json .txt`, então os `.ts` não passam pela checagem
  de caminho privado. Verificável: `git diff --stat` sem `cmd/`, `internal/` nem `go.mod`.
- R8 Texto vigente diz a verdade. (a) `.agents/coordinator.md`, bullet "You are the sole writer of `.harness/MEMORY.md`..."
  (fora do bloco cercado da cláusula anti-delegação, que não muda nem em `coordinator.md` nem em `external-clis`): ganha a frase
  "enforced at runtime for subagents by the Claude Code plugin mod (writes through Bash are not blocked)". (b) `AGENTS.md`:
  o bullet de runtime cita a exceção do mod TypeScript no motor do Claude Code (ADR-005), e a regra de despacho/escrita da
  memória diz que no Claude Code ela é imposta pelo mod e nos demais adaptadores é prosa. (c) `harness-authoring` deixa de
  sugerir `disallowedTools` como trava. Verificável: `grep -n "enforced at runtime" .agents/coordinator.md`;
  `grep -n ADR-005 AGENTS.md`; `checkAntiDelegationClause` (`go run ./cmd/dh validate .`) verde.
- R9 Versão mínima do Claude Code e comportamento em runtime antigo: só 2.1.289 está comprovada. As versões 2.1.287 e
  2.1.288 existem nesta máquina em `~/.local/share/claude/versions/`; T-1102 as usa para ver o carregamento do plugin com
  `modules` (módulo ignorado, ou plugin falha). O achado vai no relatório do T-1102. Se ficar "não verificado", ou se o
  plugin falhar ao carregar, o release espera decisão explícita do dono, registrada neste PRD.
- R10 Fora de escopo: congelar a árvore durante chamadas de CLI externa de revisão; painéis, status lines e UI; adaptadores
  que não são Claude Code; bloquear escrita por `Bash`.

## 4. Docs

- AGENTS.md
- .agents/coordinator.md
- .skills/harness-authoring/SKILL.md
- docs/tutorial.html, docs/en/tutorial.html, docs/roadmap.html, docs/en/roadmap.html, CHANGELOG.md (fim do lote, pelo docs-guide)
- docs/adr/ADR-005-trava-de-runtime-por-mod.md (nova: API de mod em acesso antecipado como dependência do pacote)

## Apêndice

### Tickets (formato ADR-003 / PRD-007)

**T-1101 · lane `backend` · `harness-maintainer` com `harness-authoring` (fonte do kit)**
- Fluxo: criar `adapters/claude-code/plugin/hooks/<arquivo>.ts` (R2 a R5); editar `hooks/hooks.json` (R1); atualizar os
  textos de R8; escrever ADR-005.
- Pronto quando: `go run ./cmd/dh validate --source-only .`; `go run ./cmd/dh build` sem diferença inesperada em `dist/`;
  `go test ./...`; `go run ./cmd/dh validate .`; `claude plugin validate` no pacote construído.
- Documentação: AGENTS.md, coordinator.md, harness-authoring, ADR-005.

**T-1102 · lane de teste · `qa-verifier`**
- Fluxo: escrever o `*.test.ts` com os casos de R3 a R5; verificar R9 com as versões 2.1.287 e 2.1.288.
- Pronto quando: `claude plugin test` no pacote construído passa; os checks de T-1101 seguem verdes; sessão real com o
  plugin construído carregado (`claude --plugin-dir dist/claude-code/dev-harness`): um subagente que tenta `Agent` e
  `Write .harness/MEMORY.md` é negado, e a sessão principal grava `.harness/MEMORY.md` sem negação; o achado de R9 está no
  relatório como fato ou "não verificado".
- Documentação: nenhum.

Revisão final da feature: integração e aderência a este PRD.

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Voltar `disallowedTools` no frontmatter | O runtime não honra (RISK-001); falsa segurança. |
| Hook clássico `PreToolUse` em shell/Go | O dono escolheu mod; evita processo e binário por chamada. A comparação não foi medida aqui. |
| Bloquear `Bash` que escreve em `.harness/` | Fora de escopo por decisão do dono; heurística frágil. |

### Riscos e decisões pendentes

- Risco: a API de mod é de acesso antecipado e pode mudar entre versões do Claude Code · mitigação: `claude plugin test` no
  pacote em cada release; ADR-005 registra a dependência.
- Risco: runtime sem `modules` quebra o carregamento do plugin para quem não atualizou · mitigação: R9; se confirmado,
  decisão do dono abaixo.
- Risco: falso positivo nega a escrita legítima do Coordenador se o motor passar `agentId` no laço principal · mitigação:
  teste de R4 e uma sessão real de `/dh:build` antes do release.
- Proposta (RISK-001, quem decide é o Coordenador): se `claude plugin test` provar a negação de `agent.spawn` por
  `parentAgentId`, citar o teste como evidência determinística no RISK-001. A regra D4 do risco (5/5 execuções de eval)
  não é substituída por isso sem decisão do dono; sugestão: mover para "mitigado por trava de runtime" no Claude Code,
  sem "resolvido".
- Decididas em 2026-10-05 (Coordenador, padrões alinhados ao dono): `.ts` no `validate` = não (R7); ADR-005 = sim, curto;
  R4 só `Write` e `Edit`; RISK-001 segue como proposta, a decisão é do Coordenador/dono na entrega.
- Decisão pendente: runtime antigo (R9) · responde: dono, se o achado for "não verificado" ou "plugin falha" · bloqueia: release.

## Resumo executado

- **Entregue**: o plugin do Claude Code traz o mod `hooks/runtime-guard.ts`, que nega o despacho de agente feito por subagente e o `Write`/`Edit` de subagente em `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, sem afetar a sessão principal (0.16.0).
- **Regras**:
  - R1 Honrada: `hooks.json` só ganhou `modules`; `claude plugin validate` no pacote passa.
  - R2 Honrada: dois handlers (`agent.spawn`, `tool.call` em `Write`/`Edit`); revisão de código.
  - R3 Honrada: testes do mod e sessão real (spawn de subagente negado).
  - R4 Honrada: 48 casos gerados e 2 pelo motor; sessão real (subagente negado, sessão principal grava). Limite registrado: sem colapso de `//`/`./`, sensível a maiúsculas.
  - R5 Honrada: casos de entrada malformada não negam.
  - R6 Honrada: `claude plugin test dist/claude-code/dev-harness` 62 pass, 0 fail; mutação com 19 falhas.
  - R7 Honrada: nenhum Go de comportamento mudou; a única linha em `internal/` é a constante de versão do lote (0.16.0), com prova de clone limpo.
  - R8 Honrada: `grep` de "enforced at runtime" em `coordinator.md` e de ADR-005 em `AGENTS.md`; `dh validate .` verde.
  - R9 Honrada: 2.1.287 carrega e impõe; 2.1.288 carrega o plugin (imposição não demonstrada, o modelo recusou o comando de teste); nenhum runtime falhou ao carregar.
  - R10 Honrada: nada fora de escopo entrou.
- **Tickets**: T-1101 (backend, harness-maintainer) concluída, revisão de código aprovada sem bloqueante; T-1102 (teste, qa-verifier) concluída, APPROVED.
- **Docs**: AGENTS.md, `.agents/coordinator.md`, `.skills/harness-authoring/SKILL.md`, ADR-005, CHANGELOG.md, README.md, README.en.md, docs/roadmap.html, docs/en/roadmap.html, docs/tutorial.html, docs/en/tutorial.html.
- **Fora**: escrita por `Bash`; congelamento da árvore durante CLI externa; painéis; adaptadores não-Claude. Correção à parte no lote: frontmatter do `/dh:review`.
