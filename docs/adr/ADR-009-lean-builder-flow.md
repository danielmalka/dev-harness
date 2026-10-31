# ADR-009 · Fluxo de build enxuto: uma rodada limpa, sem rotação, três estados, uma dúvida

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-30 |
| Decisor | Daniel Malka (dono) |
| Origem | PRD-031 (R10 a R17); inventário em `docs/prd/PRD-031-rule-inventory.md` |
| Reversibilidade | barata: edição de fonte do kit, sem mudança de binário |
| Relação | supera em parte ADR-004 (rotação, cursor, ciclos de dúvida); a segunda rodada limpa vinha do PRD-029 (0.34.0) |

## 1. Contexto

- Pergunta: que regras do fluxo de build custam mais do que pegam.
- Fato: o fluxo cobra uma segunda revisão completa por ticket, um cursor de rotação no `MEMORY.md`, seis estados de achado com ids, até 3 ciclos de dúvida e um protocolo de churn.

## 2. Decisão

Aceitas pelo dono em 2026-10-30:

- R11 Uma rodada limpa encerra o loop de código (ticket e revisão final). Rodada limpa: nenhum achado bloqueante `open` ou `escalated` e todo revisor devolveu veredito.
  - Custo removido: uma revisão completa extra por ticket. Risco: um erro de um revisor passa; ficam as checagens determinísticas, a revisão independente e o painel a pedido do dono.
- R12 Sem cursor de rotação: entrada e reentrada da etapa `code` usam a primeira entrada de `reviewers.code`.
  - Custo removido: ~11 linhas de regra, um campo no `MEMORY.md` por tarefa, a lógica de entrada not-run. Risco: o ponto cego de um revisor se repete na reentrada (propósito do ADR-004).
- R13 Estados `open`, `closed`, `escalated`; sem ids `F<rodada>-<n>`. Achado = local e título; `fixed` do builder segue `open` até um revisor confirmar.
  - Risco: reabertura mais difícil de ligar ao primeiro relato.
- R14 Dúvida no meio do build: um ciclo, só com `Doubt: yes`; duas classes (acionável ou não); sem acionável, escala ao dono. O gatilho da opção estrutural O2 sai.
  - Risco: erro de contrato ou de fronteira de confiança é pego um ciclo depois, na revisão do fim da fatia.
- R15 `Churn:` e o protocolo de três opções viram uma frase: achado novo que reverte correção aceita para o loop e o Coordenador faz ao dono uma pergunta numerada.
  - Risco: sem a frase, a oscilação corrige e reverte até o teto.
- Regra do conteúdo de commit fica na skill `delivery-readiness`.

Rejeitadas pelo dono: R16 (o teto de 6 rodadas continua) e R17 (o registro de Dispatches no `MEMORY.md` continua).

## 3. Consequências

- Fica mais fácil: menos despachos por ticket e menos estado a registrar.
- Fica mais difícil: nada novo; o risco por item está acima.
- Contratos afetados: `.agents/coordinator.md` e `.skills/code-review/SKILL.md` são o procedimento; `AGENTS.md` só enuncia a regra vigente. Sem espelho em `docs/en/` (não há `docs/en/adr/`).

## 4. Revisar quando

- Um defeito passar por revisão de rodada única que a segunda rodada teria pego.
- Um blind spot de revisor se repetir em reentrada.
