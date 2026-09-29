# PRD-000 · <título curto>

<!-- Local sugerido no projeto consumidor: .harness/prd/PRD-000.md
Preencha só o que muda a decisão. Apague os comentários ao finalizar.
Fato, hipótese e decisão ficam marcados como tal. Sem segredos. -->

| Campo | Valor |
|---|---|
| Status | rascunho / em revisão / aprovado / entregue em <data> / cancelado |
| Dono | <quem decide o escopo> |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |
| Tickets | T-000, T-000 |

## 1. Problema

<!-- Quem sofre, o que não funciona, em qual situação, qual o impacto, o que precisa mudar. Uma evidência (log, métrica, relato) vale mais que um adjetivo. Uma capability por PRD: duas capabilities independentes viram dois PRDs antes do primeiro ticket. -->

## 2. Solução

<!-- O que passa a ser verdade depois da entrega: comportamento observável, resultado para o usuário, fluxo esperado, resposta do sistema. Não é descrição de arquitetura. -->

## 3. Regras

<!-- Constraints inegociáveis, numeradas R1, R2, R3... Cada regra precisa passar neste teste: é uma condição obrigatória ou uma proibição explícita, e é verificável. Texto genérico, de preferência ou de intenção não é regra. -->

- R1 <condição obrigatória ou proibição explícita, verificável>
- R2 <condição obrigatória ou proibição explícita, verificável>

## 4. Docs

<!-- Paths dos fluxos já documentados que esta feature toca, um por linha. Se a feature não toca nenhum fluxo documentado, apague a linha de exemplo abaixo e escreva exatamente a linha: "Nenhum fluxo documentado afetado." Seção vazia não vale — document-validator reporta gap. -->

- <caminho/do/fluxo.md>

## Apêndice (opcional)

<!-- Só quando o autor achar necessário. A ausência deste apêndice não impede a virada em ticket, do mesmo jeito que a seção Docs já funciona quando vazia de itens reais (linha "Nenhum fluxo documentado afetado."). -->

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| <opção> | <custo, risco ou limite> |

### Riscos e decisões pendentes

- Risco: <descrição> · mitigação: <ação>
- Decisão pendente: <pergunta> · responde: <quem> · bloqueia: <RF/ticket>
