---
name: refactor
description: Mude a estrutura provando que o comportamento foi preservado
author: malka
argument-hint: "[alvo e motivação]"
metadata:
  roles: [coordinator, builder, qa-verifier, reviewer]
  skills: [safe-refactoring, regression-testing, code-review]
  writes: somente o conjunto de escrita autorizado
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Alvo e motivação: $ARGUMENTS. Em sequência: despache `builder` no modo `refactor` com a ferramenta Agent, modelo `sonnet`, exigindo a skill do kit `safe-refactoring`, com a `lane` do código que está sendo reestruturado, com o conjunto de escrita e a fronteira de comportamento. Congele o resultado, depois despache `qa-verifier` (modelo `sonnet`, skill `regression-testing`) para comparar com a linha de base registrada. Em seguida despache `reviewer` no modo `code` (modelo `sonnet`, skill `code-review`) sobre esse estado estável. No máximo dois especialistas rodam ao mesmo tempo, e só em trabalho independente.

## Pré-requisitos
O problema estrutural nomeado e a evidência do custo dele; sem um, não há como dizer que o trabalho teve sucesso, então pergunte ao dono e pare. As checagens existentes e os resultados atuais delas, capturados como linha de base antes de qualquer coisa se mover, incluindo uma suíte que já falha. Cobertura na fronteira que está sendo movida: se a preservação precisa de testes de caracterização, eles são acrescentados só dentro do escopo de testes autorizado; caso contrário, reporte a lacuna antes de mudar a estrutura.

## Saída
O relatório de `safe-refactoring`: status, problema, fronteira de comportamento dividida em contratado e incidental, a tabela da linha de base de antes contra depois, passos aplicados com as checagens após cada um, redução alcançada, testes de caracterização acrescentados, evidência de que o comportamento contratado não mudou, riscos restantes e o caminho de reversão por passo. Registre a decisão de fronteira, a evidência da linha de base e o caminho de reversão em `.harness/MEMORY.md`.

## Limites
- Nenhuma mudança de contrato e nenhuma feature nova; este comando só move estrutura, e um pedido de comportamento novo vai para `/dh:plan`.
- Qualquer diferença observável entre antes e depois é um defeito, não uma melhoria, por melhor que o código fique de ler.
- Cada passo deve ser pequeno o bastante para ser revertido sozinho, com checagens entre os passos e não só no fim.
- Uma checagem exigida que falhou ou não rodou deixa o trabalho parcial ou bloqueado, nunca verificado.
- Nunca faça commit, push, publicação nem deploy a partir deste comando.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:verify` para uma matriz de aceite mais ampla, ou `/dh:build` para retomar o trabalho de feature que o refactor destravou.
