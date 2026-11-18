---
name: regression-testing
description: Use quando uma mudança entregue precisa ser verificada contra os critérios de aceite dela, quando são necessários testes para uma fatia já estabilizada, ou quando um relatório precisa declarar exatamente quais checagens passaram, falharam e nunca foram rodadas. Use também quando for decidir quanto teste uma mudança merece. Não use enquanto a implementação ainda está em movimento, para diagnosticar uma falha sem explicação, nem quando o pedido é mudar o código de produção em vez de verificá-lo.
author: malka
metadata:
  provenance: adapted
  sources: ["tdd reference (good tests, mocking)", "test-driven-development and writing-good-tests", "verification-before-completion", "testing-strategy", " qa-engineer", "akitaonrails/my-skills pr-audit @7b8a5e5"]
---

# Teste de regressão

## Visão geral

Este procedimento transforma critérios de aceite em checagens executadas e devolve uma matriz em que cada critério está ligado a um comando e a um resultado literal: `passed`, `failed` ou `not-run`. A profundidade da verificação segue o risco, não o hábito. Um teste que passa estabelece só as suas asserções; observá-lo pegar o defeito pretendido dá evidência adicional de regressão. Uma checagem que não foi executada nunca vira `passed`, por mais confiante que todos estejam. Qualquer checagem exigida `failed` ou `not-run` deixa a verificação `partial` ou `blocked`.

## Quando usar

- Uma mudança está estabilizada e os critérios de aceite dela precisam ser verificados.
- Uma fatia precisa de testes de comportamento em que um revisor possa confiar.
- Um relatório precisa de um relato honesto do que foi verificado e do que não foi.
- A cobertura de um caminho arriscado é desconhecida e precisa ser estabelecida antes do release.

## Quando não usar

- A implementação ainda está mudando. Espere um estado estabilizado.
- A causa da falha é desconhecida. Use `systematic-debugging`.
- A própria correção é a tarefa. Use `incremental-implementation`.
- A checagem é uma matriz de interação de interface de usuário. Use `ui-verification` para essa superfície e esta skill para o resto.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Critérios de aceite da mudança | Pergunte ao Coordenador. Sem critérios, o teste não tem alvo e a matriz não tem coluna da esquerda. |
| Uma mudança estabilizada e a lista dos arquivos tocados | Não comece. Verificar um alvo em movimento produz resultados que ninguém consegue reproduzir. |
| Os comandos de teste do projeto e como rodar um subconjunto | Descubra-os no perfil do projeto ou na configuração de testes. Se não existirem, registre toda checagem como `not-run` com o motivo. |
| O escopo de escrita combinado para os testes | Reuse a autorização registrada se ela cobrir os arquivos. Senão, rode as checagens existentes permitidas sem editar e reporte ao Coordenador as mudanças de teste ou de configuração necessárias. A falta de escopo não concede permissão de escrita. Nunca código de produção. |
| Necessidades de ambiente (fixtures, banco descartável, navegador) | Use dados sintéticos dentro do escopo autorizado. Criar, popular ou destruir bancos de dados e serviços externos exige autorização explícita; sem ela, a checagem fica `not-run` com a dependência faltante nomeada. |

## Procedimento

1. **Leia os critérios e a mudança.** Mapeie os arquivos alterados para os comportamentos que eles podem afetar. Confirme que o estado é estável; se o implementador ainda está editando, pare e diga isso.
2. **Transforme cada critério num observável.** Para cada critério, escreva o que um observador externo veria quando ele vale. Um critério sem resultado observável é um defeito do critério, não um motivo para pulá-lo; reporte-o de volta.
3. **Dimensione a profundidade pelo risco.** Use a tabela de risco na Referência rápida. Uma mudança de baixo risco tem o comportamento afirmado uma vez; uma de alto risco recebe limites, falhas e a costura de integração. A profundidade é uma decisão que você declara, não um acidente.
4. **Faça primeiro o inventário da cobertura existente.** Rode a suíte existente da área e leia o que ela já afirma. Não teste de novo o que já está coberto; nomeie a lacuna entre a cobertura existente e os critérios.
5. **Escreva testes dentro do escopo combinado, em sequência.** O QA escreve só arquivos de teste e configuração de teste. Um defeito de produção é reportado, nunca corrigido aqui. Nunca escreva testes em arquivos que o implementador está editando no momento; a escrita de testes segue a implementação, não corre ao lado dela.
6. **Confira a detecção de falha onde for preciso.** "Falha antes da correção" significa uma falha de asserção sobre o comportamento sob teste; um erro de compilação ou de import não conta. Quando o teste novo não compila contra a base, prove o vermelho revertendo só a correção numa cópia isolada, sob as regras de isolamento deste passo. Reuse evidência documentada de antes da correção para o mesmo teste e o mesmo defeito, depois rode o teste no estado atual. Se for preciso mais evidência de detecção de defeito, use uma fixture de teste explicitamente autorizada, ou uma cópia descartável isolada criada dentro do escopo de teste autorizado; nunca fora dele, e nunca reverta uma correção nem injete um defeito no código-fonte de produção sob avaliação. Mantenha o estado avaliado inalterado, identifique o estado experimental à parte e registre a falha esperada mais o resultado no estado atual. Se nenhum método autorizado estiver disponível, reporte a lacuna; uma checagem de detecção exigida continua `not-run` e deixa a verificação `partial` ou `blocked`.
7. **Rode as checagens e leia a saída.** Comando completo, saída completa, status de saída, contagem de falhas. Não extrapole a partir de uma execução parcial nem apresente um resultado de um estado anterior como verificação do estado atual. A evidência histórica de falha continua rotulada com o estado original dela.
8. **Monte a matriz de aceite.** Uma linha por critério com o comando exato, o resultado literal e o caminho da evidência. Identifique as checagens exigidas a partir do aceite e do contexto da tarefa; não as rebaixe em silêncio. Os critérios que você não conseguiu verificar mantêm uma linha `not-run` com o motivo do bloqueio. Use `verified` só quando todas as checagens exigidas passam no estado estável atual; senão, reporte `partial` ou `blocked`. O resultado do QA não substitui a revisão de código independente.
9. **Reporte os defeitos, não os corrija.** Para cada defeito, dê reprodução, esperado, obtido, o critério que ele viola e uma severidade. Entregue-o ao implementador por meio do Coordenador.
10. **Pare o seu próprio laço em três.** Se o seu código de teste falhar três vezes por motivos dentro do teste e não do produto, pare e reporte o bloqueio em vez de insistir. Este limite vale só para o seu próprio código de teste. O laço de correção entre o implementador e o QA continua parando em seis rodadas e volta ao Coordenador.
11. **Reporte, não registre.** Resultados, defeitos e qualquer incidente vão ao Coordenador, o único que escreve `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`. A evidência detalhada fica em `.harness/tasks/<id>/`. Quando um ticket é marcado como concluído, a linha de cabeçalho `Concluído em` (`Done on` em en) é preenchida com a data de hoje, `AAAA-MM-DD`.
12. **Adote o teste primeiro por padrão, de responsabilidade do `builder`.** O `builder` escreve o teste que falha e registra a falha esperada antes de mudar o código, e o QA verifica esse resultado vermelho e o resultado verde depois da mudança; a regra do vermelho primeiro do passo 6 continua valendo para testes de regressão de bug.
13. **Aplique três heurísticas contra testes vazios.** Veja o teste novo falhar pelo motivo esperado antes de ele contar. Apague um teste que continua passando quando toda função sob teste devolve nada. Quando um teste honesto é impraticável, reporte a checagem como `not-run` e não escreva um teste fraco.
14. **Confira a documentação que a mudança invalida.** Uma mudança de comportamento, de arquitetura, de operação ou de uma decisão registrada do dono deve atualizar `README`, `AGENTS.md` e `docs/DECISIONS.md` (ou os equivalentes do projeto nomeados nos arquivos de instrução dele) na mesma mudança. Trate `documentation updated` como um critério: documentação faltando é `failed`, reportada como qualquer outro defeito, com o documento e a afirmação defasada; uma mudança que não toca nenhum desses quatro declara isso explicitamente. A verificação nunca edita os documentos por conta própria, só os reporta.

## Formato de saída

```
## Veredito
<verified / partial / blocked> e a evidência exigida que o decide

## Matriz de aceite
| Critério | Checagem (comando exato) | Resultado | Evidência |
| --- | --- | --- | --- |
| <id ou texto> | <comando> | passed / failed / not-run | <caminho relativo ou linha da saída> |

## Documentação
- <documento atualizado - o que mudou nele> | faltando: <documento> ainda diz <afirmação defasada> | sem impacto em documentação, porque <motivo>

## Prova de detecção de falha
- <nome do teste> - <evidência anterior ou fixture isolada autorizada>, falha observada <mensagem>, estado atual avaliado <resultado>; ou <not-run e o motivo>

## Testes adicionados
- <caminho relativo> - <critério que ele cobre>

## Defeitos encontrados
- <id> | severidade | reprodução | esperado | obtido | critério violado

## Lacunas de cobertura
- <comportamento sem checagem> - <por que importa>

## Não verificado
- <checagem> - <motivo do bloqueio>

## Próximo passo
- <a próxima ação única e quem é o dono dela>
```

## Referência rápida

Profundidade por risco. Escolha a linha mais alta que se aplica à mudança.

| Sinal de risco | Profundidade mínima |
| --- | --- |
| Dinheiro, autenticação, autorização, dados pessoais, operação destrutiva | Caso normal, cada ramo de falha, valores de limite e a costura de integração com colaboradores reais |
| Schema persistido, migração ou contrato público | Caso normal, compatibilidade do formato antigo e do novo, e reversão ou reversibilidade quando possível |
| Regra de negócio com ramos | Cada ramo afirmado uma vez, mais o limite entre os ramos |
| Comportamento interno com um só caminho | Um teste de comportamento sobre o resultado observável |
| Mudança só de renderização ou de texto | Checagens existentes pertinentes e inspeção direta; exercite os caminhos de interação afetados quando houver, sem exigir testes novos nem casos de falha inventados |

O que uma checagem prova.

| Checagem | Prova | Não prova |
| --- | --- | --- |
| Teste de comportamento | O comportamento afirmado vale na costura | Qualquer coisa não afirmada |
| Checagem de tipos | Os tipos concordam | Que a funcionalidade funciona |
| Lint | Estilo e algumas classes de defeito | Comportamento |
| Build | Compila e empacota | Comportamento |
| Execução manual | Funcionou uma vez, aqui, agora | Que é repetível |

Propriedades de um teste que vale a pena manter.

| Propriedade | Teste-o contra |
| --- | --- |
| Comportamento acima de implementação | Ele sobreviveria a uma refatoração que preserva o comportamento? Se não, reescreva-o |
| Um só motivo para falhar | O nome contém "e"? Divida-o |
| Expectativa independente | O valor esperado é calculado do jeito que o código o calcula? Troque por um literal |
| Nomeia a quebra que pega | Você consegue declarar a mudança de produção que o faz falhar? Se não, é um detector de mudança |
| Determinístico | Depende de um sleep, do relógio real, de ordem ou de estado compartilhado? Espere por uma condição, congele o tempo, isole a fixture |
| Código real onde importa | Faça mock só nas fronteiras do sistema, espelhe por completo o formato real dos dados e nunca afirme sobre o próprio mock |
| Isolado | Roda sozinho e em qualquer ordem? |

## Erros comuns

| Erro | Por que prejudica | Em vez disso |
| --- | --- | --- |
| Marcar como `passed` uma checagem indisponível | Um falso verde torna a matriz inteira sem valor | `not-run` com a dependência faltante nomeada |
| Aceitar a checagem de tipos como verificação de comportamento | Código bem tipado faz a coisa errada o tempo todo | Afirme o resultado observável |
| Alegar prova de detecção de defeito só a partir de um teste verde | Asserções que passam não estabelecem que o defeito pretendido é pego | Use evidência pertinente de antes da correção ou uma fixture isolada autorizada; reporte com honestidade a prova faltante |
| Editar código de produção para fazer um teste passar | O verificador deixa de ser independente do trabalho | Reporte o defeito e deixe o implementador corrigi-lo |
| Afirmar que um mock foi chamado | A asserção vale enquanto o mock existir, faça o código o que fizer | Afirme o resultado real observável |
| Dormir para esperar trabalho assíncrono | Instável sob carga, lento quando funciona | Espere pela condição, com um timeout |
| Testar tudo na mesma profundidade | O esforço cai em trivialidades enquanto o caminho arriscado fica sem cobertura | Use a tabela de risco |
| Copiar a implementação para a expectativa | O teste passa por construção | Literal derivado à mão ou fixture |
| Mock parcial de um payload externo | Passa localmente, quebra na integração quando um campo omitido é lido | Espelhe o formato documentado completo |
| Rodar um subconjunto e reportar a suíte | A afirmação é mais ampla do que a evidência | Reporte exatamente o que rodou |

## Exemplo

Entrada: critério AC-2, "um convite duplicado é rejeitado com `409 already_invited`".

```
## Veredito
partial: o caso concorrente de AC-2 falha e o AC-3 exigido está not-run.

## Matriz de aceite
| Critério | Checagem | Resultado | Evidência |
| --- | --- | --- | --- |
| AC-2 convite duplicado sequencial rejeitado | npm test -- src/invites -t 'rejects sequential duplicate invite' | passed | src/invites/service.test.ts:41 |
| AC-2 convite duplicado concorrente rejeitado | npm test -- src/invites -t 'rejects concurrent duplicate invite' | failed | src/invites/service.test.ts:58 |
| AC-3 e-mail de convite entregue | npm run test:integration | not-run | sem sandbox de e-mail neste ambiente |

## Prova de detecção de falha
- rejects sequential duplicate invite - o registro de antes da correção do `builder`
  mostra "expected 409, received 201" para este teste; o QA rodou de novo no
  estado atual avaliado e ele passou. O QA não fez nenhuma edição no código-fonte de produção.

## Defeitos encontrados
- D1 | medium | convidar duas vezes no mesmo segundo | 409 | duas linhas criadas |
  critério AC-2 | corrida na inserção concorrente, sem constraint de unicidade

## Lacunas de cobertura
- A expiração do convite não tem checagem; nada pegaria um convite pendente que nunca expira

## Não verificado
- npm run test:integration - sem sandbox de e-mail neste ambiente

## Próximo passo
- O implementador corrige D1; o QA roda de novo o AC-2 concorrente (dono: builder, lane backend)
```

AC-3 continua `not-run` com o motivo. A corrida é um relatório de defeito, não uma correção do QA.

## Relacionado

Papéis: `qa-verifier`, `builder`, `reviewer`, `coordinator`. Comandos: `/dh:verify`, e usado dentro de `/dh:build`, `/dh:fix`, `/dh:refactor`. Skills: `incremental-implementation`, `ui-verification`, `systematic-debugging`, `safe-refactoring`, `delivery-readiness`.

## Caso de prova

Um defeito conhecido numa fixture de teste autorizada é pego pela suíte e reportado com reprodução, esperado e obtido, sem edições no código-fonte de produção sob avaliação. Sem autorização de escrita de testes, as checagens existentes permitidas rodam sem edições. Uma checagem exigida indisponível aparece como `not-run` com a dependência faltante nomeada e deixa o veredito `partial` ou `blocked`. Um resultado de checagem de tipos nunca é apresentado como verificação de um critério de comportamento.
