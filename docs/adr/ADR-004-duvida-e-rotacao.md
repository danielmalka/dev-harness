# ADR-004 · Dúvida no meio do build e um revisor rotacionado na etapa code

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-05 |
| Decisor | Daniel Lemos (dono) |
| Origem | pedido de encaixar mecanismos de addyosmani/agent-skills sem skill nova; autorização de 2026-10-05 |
| Reversibilidade | barata: edição de fonte do kit, sem mudança de binário |

## 1. Contexto

- Pergunta: como apertar a qualidade no meio da fatia sem acrescentar comando, agente ou skill, e sem chamar outra CLI quando o projeto não configurou nenhuma.
- Forças: o fluxo já está lento por orquestração; o vazamento é pular a regra sob pressão e só pegar isso no painel do fim; `reviewers:` em `.harness/project.yaml` já é o contrato de CLIs revisoras.
- Fato: o campo não se chama `config.yaml`. Ausente, o fluxo é só Claude.
- A verificação do ADR-003 (R11, 2026-09-30) continua como história daquela data. Esta decisão estreita o que "a mudança toca segurança" faz na etapa `code`.

## 2. Decisão

- Escolhida: quatro mecanismos adaptados nas skills que o fluxo já carrega. Proveniência `adapted`, origem addyosmani/agent-skills (MIT). Texto original do kit; nada copiado verbatim.
- Dúvida: predicado do Coordenador, um despacho em contexto limpo, artefato e contrato, sem a conclusão do autor. Dispara em fronteira de confiança, contrato público, migração, ou quando o builder devolve a opção estrutural (O2). Não dispara em fatia mecânica nem em bug cujo teste de regressão já falha antes da correção. No máximo 3 ciclos, fora do teto de 6 rodadas de correção. O painel do fim continua, uma vez, no estado congelado. O builder não despacha revisor.
- Rotação: um revisor por entrada da etapa `code`. A primeira entrada da tarefa é `claude` se `reviewers.code` o lista, senão a primeira entrada; sem o campo, só Claude e não há CLI para rotacionar. Reentrada avança uma entrada, cursor no `MEMORY.md`, por id de tarefa. Entrada not-run não move o cursor. Painel completo só quando o dono pede essa etapa.
- Segurança: a etapa `security` despacha a lista configurada inteira e não rotaciona. Mudança que toca segurança dispara `/dh:secure`. Isso não amplia a entrada da etapa `code`.
- Fora deste lote: o corte da segunda passagem do `/dh:auto` numa fatia só.
- Descartada: oferecer Gemini ou Codex como menu novo, e rodar toda CLI configurada em toda entrada. Isso é o painel, que continua só a pedido do dono ou na etapa `security`.
- Descartada: skill, comando ou agente novo. Os 13 agentes e os 19 comandos continuam.

## 3. Consequências

- Fica mais fácil: achar o pulo no meio da fatia, e não repetir o mesmo revisor quando a mesma tarefa volta à revisão.
- Fica mais difícil: uma fatia com `Doubt: yes` ganha um despacho antes de o builder seguir. O teto de 6 não absorve esses ciclos.
- Dívida assumida: o painel de 3 "porque a mudança toca segurança" deixa de valer na etapa `code`. A etapa `security` cobre essa lista. Rever se uma fatia de código com superfície sensível passar sem `/dh:secure`.
- Contratos afetados: `.agents/coordinator.md` é o procedimento. `AGENTS.md` só enuncia a regra. `reviewers:` não ganha chave nova.

## 4. Revisar quando

- O dono pedir o corte da segunda passagem do `/dh:auto` em lote de uma fatia.
- A data de revisita da prosa de `external-clis` (21/10/2026) chegar e a rotação tiver mudado o contrato de chamada.
- Uma entrada da etapa `code` em mudança que toca segurança sair sem o `/dh:secure` correspondente.
