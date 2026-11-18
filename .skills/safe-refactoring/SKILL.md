---
name: safe-refactoring
description: Use quando a estrutura do código precisar melhorar enquanto o comportamento contratado permanece idêntico, incluindo duplicação, condicionais emaranhadas, unidades grandes demais, fronteiras ruins e código legado com pouca cobertura. Use também quando for preciso decidir se uma limpeza vale o seu risco, ou quando uma limpeza anterior mudou o comportamento por acidente. Não use quando o objetivo for um comportamento novo, quando um defeito precisar ser corrigido, nem quando um contrato público for feito para mudar.
author: malka
metadata:
  provenance: adapted
  sources: ["tech-debt prioritization", "code-simplifier review agent", " refactoring-specialist"]
---

# Refatoração segura

## Visão geral

Este procedimento muda a estrutura mantendo idêntico o comportamento contratado, e prova isso com as mesmas entradas produzindo as mesmas saídas observáveis antes e depois. Cada passo é pequeno o bastante para ser revertido sozinho, e as checagens rodam entre os passos, não no fim. O princípio central é que o comportamento preservado é o critério de aceite: uma limpeza que melhora o código e muda um único resultado observável falhou, por mais que o código fique melhor de ler.

## Quando usar

- Um problema estrutural nomeado está atrasando o trabalho ou causando defeitos: duplicação, aninhamento profundo, uma unidade grande demais, uma fronteira que vaza.
- Um código legado precisa ficar alterável antes que uma funcionalidade possa entrar com segurança.
- Uma mudança planejada está bloqueada pela estrutura e o passo estrutural deve ser separado e revisável.

## Quando não usar

- A tarefa adiciona ou altera comportamento. Use `incremental-implementation`.
- Algo está quebrado e a causa é desconhecida. Use `systematic-debugging`.
- O contrato público deve mudar. Use `api-contracts` e depois planeje a migração.
- A única motivação é gosto, sem custo observado. Um refactor sem problema nomeado é agitação sem ganho.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Autorização e conjunto de escrita, incluindo quaisquer testes de caracterização | Reutilize a autorização registrada para este escopo. Se ela não cobrir as edições necessárias, continue a avaliação somente leitura e reporte ao Coordenador a fronteira que falta; a necessidade de cobertura não autoriza edições de teste nem de configuração. |
| O problema nomeado e a evidência do seu custo | Não comece. Peça ao Coordenador o enunciado do problema. Sem ele não há como dizer que o trabalho deu certo. |
| A fronteira de comportamento: o que é contratado e o que é incidental | Derive-a dos chamadores, dos testes e dos contratos publicados e depois declare sua leitura como uma decisão antes de tocar no código. |
| Checagens existentes e seus resultados atuais | Rode-as e registre a linha de base. Uma suíte que já falha é um fato da linha de base, não motivo para pulá-la. |
| Cobertura na fronteira que você vai mover | Use primeiro os testes existentes ou os exemplos registrados. Se a preservação exigir novos testes de caracterização, adicione-os só dentro do escopo de teste autorizado (passo 5); caso contrário, reporte a lacuna antes de mudar a estrutura. |
| Entradas representativas e suas saídas atuais | Capture-as você mesmo a partir do código atual e guarde-as como o conjunto de comparação. |
| Idioma de commit e de comentário de código | O do despacho; sem ele, o que o projeto define; vale a primeira fonte que disser algo, nesta ordem: `AGENTS.md`, `CLAUDE.md`, `README.md` e, por fim, `conventions.code_language` do `project.yaml` (`en` ou `pt-br`); sem definição, inglês (ADR-011). Os comentários novos e as mensagens de commit autorizadas saem nesse idioma; o comportamento não muda. |

## Procedimento

1. **Declare o problema e a redução-alvo.** Nomeie a estrutura que dói, a evidência de que dói e a medida que você espera mover: blocos duplicados removidos, ramificações numa função, pontos de chamada de um interno que vazou, tamanho de uma unidade. Metas vagas produzem trabalho não verificável.
2. **Trace a fronteira de comportamento.** Escreva o que é contratado e, portanto, congelado: assinaturas públicas, formatos de transporte, formatos persistidos, códigos de erro, efeitos colaterais observáveis, ordenação de que os chamadores dependem e qualquer número de desempenho que foi prometido. Todo o resto é incidental e pode mudar. Quando um fato estiver na linha divisória, trate-o como contratado.
3. **Encontre os chamadores e as invariantes.** Busque no repositório todo uso do código que você vai mover, inclusive usos dinâmicos e dirigidos por configuração. Registre cada um com um caminho relativo. Uma invariante que ninguém anotou continua sendo uma invariante; declare-a explicitamente antes que vire uma regressão.
4. **Estabeleça a linha de base.** Rode as checagens existentes e registre comando mais resultado. Capture a saída observável das entradas representativas. Se você afirma que uma propriedade de desempenho será preservada, meça agora, do mesmo jeito que vai medir depois.
5. **Adicione testes de caracterização onde for necessário e autorizado.** Use a cobertura existente e os exemplos registrados quando bastarem para o risco. Se novos testes forem necessários, confira o escopo de teste autorizado antes de escrevê-los; se ele não existir, reporte a lacuna em vez de ampliar o conjunto de escrita. Um teste de caracterização afirma o que o código faz hoje, inclusive o comportamento que parece errado. Não corrija nada enquanto o escreve. Reporte os defeitos aparentes ao Coordenador como itens separados.
6. **Planeje passos pequenos e reversíveis.** Uma transformação por passo: extrair uma unidade, embutir uma unidade, renomear, introduzir um objeto de parâmetros, substituir uma cadeia de condicionais, mover uma fronteira. Cada passo deve ser reversível sozinho e deixar a árvore funcionando.
7. **Execute um passo e depois cheque.** Rode as checagens pertinentes depois de cada passo. Nunca encadeie duas transformações antes de uma checagem; quando algo quebra, você precisa saber qual movimento quebrou.
8. **Compare o depois com a linha de base.** Mesmas entradas, mesmas saídas observáveis, mesmos resultados de checagem. Qualquer diferença é um defeito do refactor. Reverta o passo e reporte. Uma mudança de contrato é uma tarefa separada, com autorização e migração próprias, nunca um item dentro de uma limpeza.
9. **Verifique a redução.** Meça o problema de novo com a medida do passo 1 e reporte o número de antes e de depois. Um refactor que moveu código sem reduzir o problema nomeado deve ser revertido, não defendido.
10. **Mantenha o escopo.** Nenhuma funcionalidade nova, nenhuma mudança de contrato, nenhuma correção não relacionada, nenhuma dependência adicionada. As melhorias que você notar vão para o relatório como candidatas.
11. **Reporte e entregue.** Entregue a evidência de preservação, a redução, os riscos restantes e como reverter. Uma checagem de preservação obrigatória `failed` ou `not-run` deixa o trabalho `partial` ou `blocked`, mesmo quando os status de antes e depois coincidem. As checagens do builder passando deixam o estado estável pronto para a QA e a revisão independentes, não autoaprovado. Fatos e incidentes vão ao Coordenador, o único a escrever `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`.

## Formato de saída

```
## Status
<ready for independent QA/review / partial / blocked> - nomeie toda checagem obrigatória `failed` ou `not-run`; qualquer uma delas mantém o status `partial` ou `blocked`

## Problema
<o problema estrutural e a evidência do seu custo>

## Fronteira de comportamento
- Contratado (congelado): <itens, com caminhos relativos>
- Incidental (livre para mudar): <itens>

## Linha de base
| Checagem ou entrada | Antes | Depois |
| --- | --- | --- |
| <comando ou caso de entrada> | <resultado ou saída> | <resultado ou saída> |

## Passos aplicados
1. <transformação> - <arquivos> - checagens depois: <resultado>

## Redução obtida
- <medida>: <antes> -> <depois>

## Testes de caracterização adicionados
- <nome do teste> - <comportamento que fixa, inclusive qualquer coisa que pareça errada>

## Comportamento preservado
- <evidência de que a lista contratada não mudou>

## Riscos residuais
- <risco> - <por que permanece> - <o que o detectaria>

## Caminho de reversão
- <como desfazer, passo a passo ou por intervalo de commits>

## Reportado, não corrigido
- <defeito ou melhoria notada, deixada para uma tarefa separada>
```

## Referência rápida

| Mudança | Permitida dentro de um refactor |
| --- | --- |
| Renomear um símbolo interno, dividir ou fundir unidades internas | Sim |
| Substituir uma cadeia de condicionais por despacho, extrair um objeto de valor | Sim |
| Mover código entre arquivos ou camadas, ajustar imports | Sim, se nenhum caminho de consumidor mudar |
| Mudar uma assinatura pública, rota, payload ou código de erro | Não, isso é uma mudança de contrato |
| Mudar um formato persistido ou o comportamento de uma migração | Não, isso é uma tarefa de dados |
| Corrigir um defeito notado pelo caminho | Não, reporte |
| Adicionar uma feature flag, opção ou ponto de extensão | Não, não existe consumidor |
| Mudar a redação de um log que um consumidor interpreta | Não, na prática ela é contratada |

Ordene o trabalho por custo, não por irritação: a prioridade sobe com o quanto a estrutura atrasa o trabalho atual e com o risco que ela carrega, e desce com o esforço da mudança. Os primeiros candidatos são os que bloqueiam uma funcionalidade planejada, carregam um histórico de defeitos e podem ser movidos em passos pequenos.

Regra do tamanho do passo: se você não consegue descrever o passo numa frase e revertê-lo com uma ação, são dois passos.

## Erros comuns

| Erro | Por que prejudica | Faça em vez disso |
| --- | --- | --- |
| Refatorar sem linha de base | Nada pode provar que o comportamento sobreviveu | Registre checagens e saídas antes de tocar no código |
| Corrigir um defeito durante a limpeza | O diff passa a misturar mudança de comportamento com estrutura, e a revisão não consegue separá-las | Reporte, agende |
| Escrever testes de caracterização que afirmam o comportamento desejado | Falham de imediato e tentam uma mudança de comportamento silenciosa | Afirme o que o código faz hoje |
| Várias transformações entre as checagens | A quebra não pode ser atribuída e a reversão é grande | Um passo, uma checagem |
| Ajustar o valor esperado para casar com a saída nova | A única evidência de preservação é destruída | Trate a diferença como defeito |
| Mudar uma assinatura pública porque fica mais limpo | Os consumidores quebram sob a bandeira da limpeza | Mudança de contrato, tarefa separada, migração |
| Extrair para eliminar uma repetição visual | Uma abstração errada custa mais que a duplicação que substituiu | Extraia quando for o mesmo conceito com o mesmo motivo para mudar |
| Declarar sucesso sem medir o problema de novo | Movimento é reportado como melhoria | Reporte o número de antes e de depois |
| Refatorar a área inteira porque você já está lá | A mudança deixa de ser revisável e reversível | Mantenha-se no problema nomeado |

## Exemplo

Entrada: "o total do checkout é calculado em três lugares, e eles já divergiram uma vez".

```
## Fronteira de comportamento
- Contratado: totais da resposta do checkout, arredondamento em 2 casas decimais,
  ordem desconto antes do imposto (src/checkout/api.ts:40, testes em :118)
- Incidental: nomes dos helpers internos, organização dos arquivos

## Linha de base
| Checagem ou entrada | Antes | Depois |
| --- | --- | --- |
| npm test -- src/checkout | 18 passed | 18 passed |
| carrinho [10.00, 5.55] cupom SAVE10 | 14.00 | 14.00 |
| carrinho [0.01] sem cupom | 0.01 | 0.01 |

## Passos aplicados
1. Caracterizar os três totais com 6 casos literais - checagens: `passed`
2. Extrair computeTotal para src/checkout/total.ts - checagens: `passed`
3. Apontar os três pontos de chamada para ela - checagens: `passed`

## Redução obtida
- Pontos de cálculo do total: 3 -> 1

## Testes de caracterização adicionados
- totais do checkout, 6 casos literais, fixa a ordem de arredondamento do admin como ela é hoje

## Comportamento preservado
- Os três carrinhos representativos retornam os mesmos totais antes e depois, e
  os 18 testes de checkout existentes passam sem alteração

## Riscos residuais
- Sem cobertura de carrinhos multimoeda; uma mudança de arredondamento ali passaria despercebida

## Caminho de reversão
- Reverter os passos 3, 2, 1 nessa ordem; cada um é um único commit

## Reportado, não corrigido
- O caminho do admin arredonda antes do desconto, produzindo uma diferença de um centavo.
  Fixada por um teste de caracterização como comportamento atual. Precisa de uma decisão.
```

A divergência de um centavo é congelada por um teste e escalada, não corrigida dentro do refactor.

## Relacionado

Papéis: `builder`, `qa-verifier`, `reviewer`, `coordinator`. Comandos: `/dh:refactor`. Skills: `regression-testing`, `incremental-implementation`, `systematic-debugging`, `architecture-decisions`, `api-contracts`.

## Caso de prova

Um refactor autorizado de fixture, com um problema estrutural nomeado, preserva as mesmas entradas representativas e saídas observáveis, mostradas lado a lado, enquanto a medida declarada cai por um número reportado. Um comportamento que parece errado é fixado por um teste de caracterização autorizado e reportado em vez de corrigido, e nenhuma assinatura pública muda. A falta de conjunto de escrita de teste produz uma lacuna reportada sem edições; uma checagem de preservação obrigatória indisponível deixa o resultado parcial ou bloqueado.
