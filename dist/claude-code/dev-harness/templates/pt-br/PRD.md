# PRD-000 · <título curto>

<!-- Local sugerido no projeto consumidor: .harness/prd/PRD-000.md
Preencha só o que muda a decisão. Apague os comentários ao finalizar.
Fato, hipótese e decisão ficam marcados como tal. Sem segredos. -->

| Campo | Valor |
|---|---|
| Status | rascunho / em revisão / aprovado / cancelado |
| Dono | <quem decide o escopo> |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |
| Stories | ST-000, ST-000 |

## 1. Problema

<!-- Situação atual, quem sofre com ela, com que frequência e qual o custo. Uma evidência (log, métrica, relato) vale mais que um adjetivo. -->

## 2. Resultado esperado

<!-- O que muda para o usuário quando estiver pronto. Como saberemos: 1 a 3 sinais observáveis, com baseline se existir. -->

- Sinal 1: <métrica ou comportamento> · hoje: <valor> · alvo: <valor>

## 3. Usuários e cenários

| Usuário | Cenário | Hoje | Depois |
|---|---|---|---|
| <perfil> | <o que tenta fazer> | <o que acontece> | <o que passa a acontecer> |

## 4. Escopo

**Entra**
- RF-01 <comportamento em uma frase, verificável>
- RF-02

**Não entra** (e por quê, quando não for óbvio)
- <item>

## 5. Critérios de aceite

<!-- Um por requisito. Formato: Quando <condição>, então <resultado observável>. -->

- AC-01 (RF-01) Quando <X>, então <Y>.
- AC-02 (RF-02) Quando <X>, então <Y>.

## 6. Restrições e impacto técnico

<!-- Só o que restringe a solução. Deixe em branco o que não se aplica. -->

- Dados: <entidades novas/alteradas, migração, retenção, LGPD>
- Contratos: <APIs, eventos ou integrações afetadas; compatibilidade exigida>
- Segurança: <autenticação, autorização, dados sensíveis>
- Operação: <ambientes, feature flag, rollout, observabilidade>
- Stack e dependências: <lib nova, versão mínima, remoção>
- Prazo / custo: <se existir>

## 7. Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| <opção> | <custo, risco ou limite> |

## 8. Riscos e decisões pendentes

<!-- Risco conhecido = o que pode dar errado + mitigação. Decisão pendente = pergunta cuja resposta muda o escopo + quem responde + até quando. Consultar .harness/RISKS.md se a área tem incidente registrado. -->

- Risco: <descrição> · mitigação: <ação>
- Decisão pendente: <pergunta> · responde: <quem> · bloqueia: <RF/ST>

## 9. Referências

- <caminho relativo ou link: pesquisa, protótipo, incidente, ADR>
