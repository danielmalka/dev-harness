# T-000 · <título curto>

<!-- Local: .harness/tasks/T-000/TASK.md
Este é o único arquivo que um agente builder recebe: já carrega a story, o
fluxo esperado e a cópia literal das regras que ele precisa honrar. Não lê
o PRD inteiro. Fluxo (seção 3) + Pronto quando (seção 6) são a única fonte
de aceite deste ticket — não existe seção "Critérios de aceite" separada.

Para Tipo: bug, não há seção própria de defeito: a seção Fluxo carrega a
reprodução (passos numerados terminando no par observado-vs-esperado) e a
seção Pronto quando ganha o item do teste de regressão verde. Ver os
comentários dessas duas seções abaixo.

Apague os comentários ao finalizar. -->

| Campo | Valor |
|---|---|
| Tipo | task / bug |
| Status | pronta / em andamento / bloqueada / em revisão / concluída |
| PRD (RF-<n>) | PRD-000 (RF-01) |
| Papel sugerido | builder (lane, modo) / qa-verifier / reviewer (modo) / architect / ... |
| Depende de | T-000, <contrato, decisão, migração> |
| Bloqueia | T-000 |
| Autorização | <o que pode ser escrito; commit/push/deploy exigem autorização explícita> |
| Criado / atualizado | AAAA-MM-DD / AAAA-MM-DD |

## 1. Lane

<!-- Um valor por ticket, lista fechada. Uma story que precisa de API e de
tela com responsabilidades independentes vira dois tickets, um por lane;
"e testa também" ao final deste ticket não gera sozinho um ticket de teste. -->

Lane: backend / frontend / dados / infra / teste (unitário) / teste (integração) / teste (unitário e integração)

<!-- O planejador preenche. `yes` nomeia o motivo. O builder não despacha
a revisão: `yes`, ou uma decisão em aberto cuja opção estrutural seguiria,
volta ao Coordenador, que roda a dúvida no meio do build. O rótulo do
campo fica `Doubt` nos dois idiomas, porque o procedimento casa essa linha. -->

Doubt: yes (trust boundary | public contract | migration) | no

## 2. Story

<!-- Uma frase, sempre neste formato. -->

Como <ator>, quero <ação>, para <resultado>.

## 3. Fluxo

<!-- Passos numerados, observáveis, "Quando X, então Y" — ator usuário,
sistema ou teste. Este bloco, junto com Pronto quando (seção 6), é a única
fonte de aceite do ticket.

Para Tipo: bug, este campo carrega a reprodução do defeito: os mesmos
passos numerados "Quando X, então Y" terminam no par observado-vs-esperado
— o que acontece hoje (o comportamento com defeito) e o que deveria
acontecer em vez disso. Não crie uma seção separada para isso. -->

1. Quando <ação>, então <resultado observável>.
2. Quando <ação>, então <resultado observável>.

## 4. Regras

<!-- Cópia literal, sem paráfrase e sem resumo, só das regras R<n> do PRD
que este ticket precisa honrar. Sem link. Regra que não cabe aqui não entra. -->

- R1 <cópia literal da regra correspondente do PRD>

## 5. Docs

<!-- Paths a atualizar neste ticket, herdados e refinados da seção Docs do
PRD, ou a linha exata "nenhum". -->

- <caminho/do/doc.md>

## 6. Pronto quando

<!-- O conteúdo mínimo do ticket: o fluxo da seção 3 aconteceu; cada regra
copiada na seção 4 tem evidência; os testes relevantes estão verdes; a
documentação listada na seção 5 foi atualizada.

Para Tipo: bug, este campo ganha mais um item: o teste de regressão deste
defeito está verde. -->

- o fluxo descrito na seção 3 aconteceu
- as regras copiadas na seção 4 têm evidência
- os testes relevantes estão verdes
- a documentação listada na seção 5 foi atualizada (ou a seção diz "nenhum")

## 7. Resultado

<!-- Preenchido pelo executor ao entregar. Só o que foi feito e executado. -->

- Conclusão: <o que foi entregue, em uma frase>
- Arquivos alterados: `<caminho>`, `<caminho>`
- Checks: `<comando>` → passou / falhou / não executado (<motivo>)
- Limitações: <o que não foi verificado ou ficou pendente>
- Pendências para o Coordenador: <decisão, incidente, atualização de memória>
- Próximo passo: <ação concreta> · autorização necessária: <qual>

## 8. Riscos e limites

- Risco: <o que pode quebrar> · mitigação: <ação>
- Ao encontrar decisão aberta: parar e reportar; não inventar comportamento.
- Após 6 rodadas de correção sem sucesso: devolver ao Coordenador com evidência.
