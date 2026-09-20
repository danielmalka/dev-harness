# T-000 · <título curto>

<!-- Local: .harness/tasks/T-000/TASK.md (ou BUG-000 para defeito).
Este é o despacho que um agente recebe. Tudo que ele precisa para começar
sem perguntar está aqui ou está linkado. Apague os comentários ao finalizar.
Para bug, preencha a seção 2B e ignore a 2A. -->

| Campo | Valor |
|---|---|
| Tipo | task / bug |
| Status | pronta / em andamento / bloqueada / em revisão / concluída |
| Story / PRD | ST-000 (AC-01) / PRD-000 |
| Papel sugerido | backend-builder / frontend-builder / data-engineer / debugger / ... |
| Depende de | T-000, <contrato, decisão, migração> |
| Bloqueia | T-000 |
| Autorização | <o que pode ser escrito; commit/push/deploy exigem autorização explícita> |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |

## 1. Objetivo

<!-- Uma frase com o resultado verificável. Depois o porquê, em uma linha. -->

## 2A. Escopo (task)

**Entra**
- <mudança concreta>

**Não entra**
- <o que o agente NÃO deve tocar, mesmo que pareça útil>

## 2B. Defeito (bug)

- Sintoma: <texto literal do erro, log sanitizado ou comportamento observado>
- Esperado: <o que deveria acontecer>
- Reprodução: <passos ou comando; se não reproduz, registrar tentativas>
- Ambiente: <versão, plataforma, config, dados>
- Desde quando / mudança recente: <commit, release ou desconhecido>
- Causa: <hipótese> · confirmada: sim / não · evidência: <...>
- Gravidade: <impacto real> · incidente em RISKS.md: <id ou nenhum>

## 3. Contexto técnico

<!-- O mínimo para não reinventar: onde está o código, que padrão seguir, que contrato respeitar. Caminhos relativos. -->

- Arquivos / módulos candidatos: `<caminho>`
- Padrão local a seguir: <exemplo existente em `<caminho>`>
- Contratos a respeitar: <API, evento, schema, interface>
- Decisões já tomadas: <ADR ou decisão do PRD/story>
- Regras críticas / incidentes anteriores: <RISKS.md#id ou nenhum>

## 4. Mudanças de ambiente

<!-- Tudo que altera o ambiente além do código. Vazio = nenhuma. -->

| Tipo | Item | Ação | Motivo |
|---|---|---|---|
| lib | <nome@versão> | instalar / atualizar / remover | <por quê> |
| env var | <NOME> | criar / alterar | <finalidade; valor não vai aqui> |
| migração | <arquivo> | criar / aplicar | <reversível? sim / não; plano de retorno> |
| config / infra | <arquivo ou serviço> | alterar | <...> |

## 5. Plano de execução

<!-- Passos pequenos e verificáveis. Cada passo termina com um check. -->

1. <passo> → check: <comando ou observação>
2. <passo> → check:
3. <passo> → check:

## 6. Testes

<!-- Cada critério de aceite tem pelo menos um teste. Comando real do projeto. Typecheck não prova comportamento. -->

| Cenário | Tipo | Arquivo / comando | Cobre |
|---|---|---|---|
| <caminho feliz> | unit / integração / e2e | `<comando>` | AC-01 |
| <erro relevante> | | | AC-02 |
| <regressão do bug> | | | 2B |

Verificação manual / navegador: <fluxo> ou "não se aplica"

## 7. Critérios de aceite

- AC-01 Quando <X>, então <Y>.
- AC-02 Quando <X>, então <Y>.

## 8. Riscos e limites

- Risco: <o que pode quebrar> · mitigação: <ação>
- Ao encontrar decisão aberta: parar e reportar; não inventar comportamento.
- Após 2 rodadas de correção sem sucesso: devolver ao Coordenador com evidência.

## 9. Resultado

<!-- Preenchido pelo executor ao entregar. Só o que foi feito e executado. -->

- Conclusão: <o que foi entregue, em uma frase>
- Arquivos alterados: `<caminho>`, `<caminho>`
- Checks: `<comando>` → passou / falhou / não executado (<motivo>)
- Limitações: <o que não foi verificado ou ficou pendente>
- Pendências para o Coordenador: <decisão, incidente, atualização de memória>
- Próximo passo: <ação concreta> · autorização necessária: <qual>
