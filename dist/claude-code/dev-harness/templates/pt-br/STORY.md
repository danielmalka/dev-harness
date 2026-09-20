# ST-000 · <título curto>

<!-- Local sugerido: .harness/stories/ST-000.md
Uma story = um comportamento entregável e testável de ponta a ponta.
Se não cabe em poucas tasks, divida. Apague os comentários ao finalizar. -->

| Campo | Valor |
|---|---|
| Status | rascunho / pronta / em andamento / bloqueada / concluída |
| PRD | PRD-000 (RF-01, RF-02) |
| Depende de | ST-000, <decisão ou contrato> |
| Tasks | T-000, T-000, BUG-000 |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |

## 1. Comportamento

<!-- Uma frase: quem, faz o quê, para obter o quê. Depois o contexto mínimo que um agente precisa para não inventar: regra de negócio, exemplo real, limite. -->

Como <usuário>, quando <situação>, quero <ação> para <resultado>.

Contexto:
- <regra de negócio ou exemplo concreto>

## 2. Cenários de uso

<!-- Dado / Quando / Então. Cubra: caminho feliz, erro relevante, vazio ou limite. Cada cenário vira teste. -->

### CN-01 · <caminho feliz>
- Dado <estado inicial>
- Quando <ação>
- Então <resultado observável>

### CN-02 · <erro ou limite>
- Dado
- Quando
- Então

## 3. Critérios de aceite

<!-- Herde ou refine os AC do PRD. Cada AC referencia o cenário que o prova. -->

- AC-01 (CN-01) Quando <X>, então <Y>.
- AC-02 (CN-02) Quando <X>, então <Y>.

## 4. Fora de escopo

- <o que fica para outra story, com o ID se já existir>

## 5. Impacto técnico

<!-- Preencha só o que se aplica. Serve para decompor em tasks e escolher os papéis. -->

| Área | Impacto |
|---|---|
| Dados / schema | <entidade, migração, backfill> |
| API / contrato | <operação nova ou alterada; compatibilidade> |
| Interface | <telas, estados carregando/erro/vazio, acessibilidade> |
| Integrações | <serviço externo, evento, fila> |
| Dependências | <lib instalar / atualizar / remover> |
| Segurança / dados sensíveis | <autorização, PII> |

## 6. Estratégia de verificação

<!-- Como a story será provada pronta, além dos testes das tasks. -->

- Testes: <unitário / integração / e2e; o que cada nível prova>
- Verificação manual ou em navegador: <fluxo, teclado, duas larguras> ou "não se aplica"
- Dados de teste: <fixture ou dado sintético>

## 7. Riscos e decisões pendentes

- Risco: <descrição> · mitigação: <ação> · incidente relacionado: <RISKS.md#id ou nenhum>
- Decisão pendente: <pergunta> · responde: <quem> · bloqueia: <task>

## 8. Evidência de conclusão

<!-- Preenchido ao fechar. Só o que foi executado de fato. -->

| AC | Evidência (caminho relativo, comando, resultado) | Estado |
|---|---|---|
| AC-01 | | passou / falhou / não executado |
