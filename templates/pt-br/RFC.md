# RFC-000 · <título curto>

<!-- Local sugerido: .harness/rfc/RFC-000.md
Um RFC propõe uma mudança relevante e coleta objeções ANTES de construir.
Quando aprovado, a decisão durável vira ADR; o trabalho vira stories/tasks.
Mantenha curto: quem lê deve conseguir objetar em uma passada. -->

| Campo | Valor |
|---|---|
| Status | rascunho / em discussão / aprovado / rejeitado / retirado |
| Autor | <quem propõe> |
| Revisores | <quem precisa opinar> · prazo: AAAA-MM-DD |
| Origem | PRD-000 / ST-000 / incidente RISKS.md#id / dívida técnica |
| Resultado | ADR-000 / ST-000 (preenchido ao fechar) |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |

## 1. Resumo

<!-- Três frases: o problema, a proposta, o que muda para quem. -->

## 2. Motivação

<!-- Por que agora. Evidência do problema (métrica, incidente, custo). O que acontece se não fizermos nada. -->

## 3. Proposta

<!-- A mudança em termos concretos: componentes, fluxo, contratos, dados. Diagrama em texto se ajudar. Sem código de implementação, só interfaces que importam. -->

- Componentes afetados: `<caminho ou módulo>`
- Fluxo proposto: <passo → passo → passo>
- Contratos novos ou alterados: <API, evento, schema; compatibilidade>
- Dados: <entidade, migração, backfill, reversibilidade>

## 4. Alternativas

| Alternativa | Vantagem | Por que não |
|---|---|---|
| Manter como está | | |
| <opção B> | | |

## 5. Impactos

<!-- Preencha só o que se aplica. -->

- Segurança / dados sensíveis: <...>
- Operação: <deploy, rollout, feature flag, observabilidade, rollback>
- Dependências: <lib instalar / atualizar / remover; versão mínima>
- Compatibilidade: <o que quebra; plano de migração para consumidores>
- Custo: <esforço estimado ou "não estimado">

## 6. Plano de adoção

<!-- Fatias na ordem, cada uma verificável. Nomeie a primeira fatia que prova a ideia. -->

1. <fatia> → prova: <check>
2. <fatia> → prova:
3. <fatia> → prova:

## 7. Riscos e questões abertas

- Risco: <o que pode dar errado> · mitigação: <ação> · incidente relacionado: <RISKS.md#id ou nenhum>
- Questão aberta: <pergunta> · responde: <quem> · bloqueia aprovação: sim / não

## 8. Discussão

<!-- Objeções e respostas, datadas. Não apague objeção resolvida; marque como resolvida. -->

- AAAA-MM-DD · <quem> · <objeção> → <resposta> · resolvida / aberta

## 9. Decisão

<!-- Preenchido ao fechar. -->

- Resultado: aprovado / rejeitado / retirado · data: AAAA-MM-DD · por: <quem>
- Motivo: <uma frase>
- Desdobramentos: ADR-000, ST-000
