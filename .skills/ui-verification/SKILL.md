---
name: ui-verification
description: Use quando um fluxo voltado ao usuário foi construído ou alterado e alguém precisa confirmar que ele de fato funciona quando exercitado, cobrindo operação por teclado, duas larguras de viewport, estados de carregamento, de erro e vazio, e noções básicas de acessibilidade. Use também quando uma mudança foi declarada concluída apenas com base em uma checagem de tipos ou em um teste unitário que passou. Não use quando o código não tem superfície renderizada, para direção de arte visual pura, nem como substituto do diagnóstico de uma falha que já tem reprodução.
author: malka
metadata:
  provenance: adapted
  sources: ["accessibility-review WCAG AA pass", "interface-design review criteria", "browser automation practice", " frontend-developer"]
---

# Verificação de UI

## Visão geral

Exercite o fluxo como uma pessoa o faria e depois reporte cada critério como `passed`, `failed` ou `not-run`. Uma checagem de tipos, um build e um teste unitário verde provam que o código compila e que as unidades se comportam; nenhum deles prova que um humano consegue completar o fluxo. Quando há um navegador disponível, a verificação produz evidência observada. Quando não há, os critérios de interação são registrados como `not-run` e a verificação visual é declarada `pending`. Um critério nunca é marcado `passed` porque provavelmente funciona.

## Quando usar

- Uma tela, um componente, um formulário ou um caminho de navegação foi acrescentado ou alterado.
- Os critérios de aceite descrevem algo que o usuário faz, não algo que uma função devolve.
- Uma mudança foi reportada como completa tendo só um compilador ou um linter como evidência.
- Um comportamento de carregamento, de erro ou de vazio foi acrescentado e ninguém o viu renderizar.
- Antes do handoff ou do release de qualquer coisa com superfície renderizada.

## Quando não usar

- A mudança não tem superfície voltada ao usuário. Use `regression-testing`.
- A pergunta é qual direção visual seguir. Isso é uma decisão de design, não uma verificação.
- Um defeito já é reproduzível e o trabalho é achar a causa dele. Isso é `systematic-debugging`.
- O fluxo ainda está sendo construído. Verifique um estado estabilizado.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Os critérios de aceite, como comportamento visível ao usuário | Derive uma lista de critérios da mudança e devolva-a ao Coordenador para confirmação antes de rodar; não verifique contra critérios que você inventou e confirmou para si mesmo. |
| Como rodar a aplicação localmente: comando, porta, rota de entrada | Pergunte. Registre todo critério de interação como `not-run` até que seja respondido. |
| Dados de seed ou de fixture que produzem os estados de carregamento, de erro e vazio | Verifique os estados que você consegue alcançar e registre os demais como `not-run` com o motivo. |
| Um navegador, headless ou não | Rode as checagens estáticas, registre os critérios de interação como `not-run` e declare a verificação visual `pending`. Nunca instale um navegador, um driver ou um pacote para obtê-lo, e nunca baixe um runtime na hora da verificação. |
| Ativos visuais exigidos: fontes, ícones, imagens, tokens | Não os busque na hora da verificação. Um ativo ausente é um critério `failed`, porque o projeto deve entregá-lo. |

## Procedimento

1. **Transforme os critérios numa checklist.** Uma linha por comportamento visível ao usuário, redigida como um resultado observável. Acrescente as linhas fixas da Referência rápida que se aplicam à superfície alterada e identifique as checagens exigidas pelo aceite e pelo contexto da tarefa. Não invente estados de carregamento ou de falha para uma superfície que não os tem. Esta checklist é o esqueleto da saída; toda linha aplicável vai terminar como `passed`, `failed` ou `not-run`.

2. **Estabeleça que evidência é possível.** Determine se a aplicação roda localmente e se há um navegador disponível. Headless serve e é o padrão num ambiente sem display. Registre a resposta agora, porque ela decide quais linhas podem chegar a `passed`. Nunca instale um navegador, um driver ou um pacote para obtê-lo, e nunca baixe um runtime na hora da verificação. Sem navegador, as linhas de interação ficam `not-run`.

3. **Rode primeiro as checagens do próprio projeto.** Build, checagem de tipos, lint, testes unitários e de componente, com os comandos do perfil do projeto. Registre cada uma como um comando mais um resultado. Essas linhas são evidência de compilação e de comportamento unitário, e de nada mais. Não deixe um resultado verde fechar uma linha de interação.

4. **Alcance o fluxo.** Inicie a aplicação, navegue até o ponto de entrada e confirme que a superfície renderiza. Se não renderizar, pare e reporte uma falha aqui; toda linha seguinte fica `not-run`, não `failed`.

5. **Exercite o caminho feliz com o ponteiro.** Complete o fluxo de ponta a ponta com a entrada que um usuário real daria. Registre o que você fez e o que apareceu em cada passo.

6. **Exercite o mesmo caminho só com o teclado.** Alcance todo elemento interativo com tab, ative-o com a tecla esperada e confirme que o indicador de foco é visível em cada parada. Confira que a ordem do foco segue a ordem visual, que nada é alcançável mas invisível, e que um diálogo prende o foco enquanto está aberto e o devolve ao gatilho ao fechar. Um controle que o ponteiro usa e o teclado não é uma linha `failed`, não uma nota.

7. **Exercite os estados aplicáveis.** Leve a superfície aos estados de carregamento, de erro e vazio usando fixtures autorizadas ou controles de resposta do navegador, sem editar o código-fonte de produção. Confirme que cada estado aplicável é distinguível, diz algo útil e oferece um caminho adiante. Um estado de erro que renderiza uma região em branco é uma falha. Um estado exigido que não pode ser alcançado continua `not-run`, com o motivo.

8. **Confira duas larguras de viewport.** Uma estreita, em torno de um celular, e uma larga, em torno de um desktop. Confirme que não há rolagem horizontal da página na largura estreita, que o conteúdo não é cortado nem sobreposto, que os alvos interativos continuam confortáveis de tocar e que o fluxo continua completável nas duas. Registre as duas larguras que você usou.

9. **Rode o básico de acessibilidade.** Todo input tem um rótulo programático, não só um texto ao lado. Imagens que carregam significado têm texto alternativo e as decorativas não são anunciadas. A estrutura usa semântica real: títulos em ordem, listas como listas, botões como botões, links como links, landmarks presentes. Elementos interativos expõem nome, papel e valor atual. O contraste do texto atinge pelo menos 4.5 para 1, e o de texto grande e de componentes de interface pelo menos 3 para 1; calcule a razão a partir das cores resolvidas de primeiro plano e de fundo (estilo computado, não o nome do token) e, se você não conseguir resolver as duas cores, a linha fica `not-run` com esse motivo, nunca estimada a olho. Os erros são identificados em texto, não só por cor.

10. **Colete a evidência.** Capture uma captura de tela por estado que mudou, nas duas larguras onde o layout importa, salva em `.harness/tasks/<task-id>/evidence/`, referenciada por caminho relativo. Leia o console do navegador e registre erros e avisos; um erro não capturado durante o fluxo é uma linha `failed` mesmo quando a tela parece correta.

11. **Confirme que os ativos pertencem ao projeto.** Toda fonte, ícone, imagem ou estilo de que a superfície depende deve estar presente no repositório ou numa dependência declarada. Uma superfície que só parece certa porque um recurso externo por acaso carregou não está verificada.

12. **Preencha a matriz com honestidade.** `passed` significa observado. `failed` significa observado como errado. `not-run` significa que não foi exercitado, e sempre traz um motivo. Aceite só quando toda checagem exigida estiver `passed` no estado estável atual. Uma linha exigida `failed` significa `rejected`; uma linha exigida `not-run` significa `incomplete`, nunca `accepted`. A verificação visual só está `complete` quando todas as checagens visuais e de interação exigidas foram exercitadas, mesmo que algumas tenham falhado; senão, está `pending`. A disponibilidade de um navegador, sozinha, não estabelece nem a conclusão nem o aceite. `accepted` aqui significa que se observou que os critérios de UI valem neste estado; não é aprovação de entrega. Marcar o trabalho como verificado e fechar a tarefa cabe ao Coordenador.

13. **Devolva o resultado.** Este procedimento verifica; não corrige. Reporte as falhas ao autor ou ao Coordenador com detalhe suficiente para reproduzir. Só o Coordenador escreve `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`; reporte fatos, caminhos de evidência e incidentes para o Coordenador registrar.

## Formato de saída

```
## Verificação de UI
- Superfície: <rota ou componente>
- Ambiente: <como foi rodada> | navegador: headless | headed | nenhum
- Viewports: <px estreito> e <px largo>

## Matriz
| # | Critério | Método | Resultado | Evidência |
|---|----------|--------|-----------|-----------|
| 1 | <comportamento visível ao usuário> | ponteiro | passed | <caminho/relativo/captura.png> |
| 2 | <comportamento> | teclado | failed | <caminho/relativo/captura.png> |
| 3 | <comportamento> | viewport 390 | not-run | sem navegador disponível |

## Falhas
- [<critério>] <o que aconteceu> vs <o que era esperado>
  - Passos: <como reproduzir>
  - Evidência: <caminho relativo ou trecho do console>

## Console
- Erros: <contagem> - <texto do primeiro erro>
- Avisos: <contagem>

## Checagens
- <comando> - passed | failed | not-run

## Resumo
- Passaram: <n> | Falharam: <n> | Não rodados: <n>
- Verificação visual: complete | pending (<motivo>)
- Veredito: accepted | rejected | incomplete
```

## Referência rápida

| Linha fixa | O que prova | Nunca aceite no lugar |
| --- | --- | --- |
| O fluxo se completa com o ponteiro | O estado final observado depois de uma entrada real | Um teste unitário aprovado do handler de envio |
| O fluxo se completa só com o teclado | Todo controle alcançado por tab e ativado | Um atributo ARIA presente no código-fonte |
| O foco é sempre visível | Indicador observado em cada parada | Um estilo de foco definido na folha de estilos |
| O estado de carregamento renderiza | O estado visto enquanto a requisição está pendente | O componente de spinner existir |
| O estado de erro renderiza e explica | O estado visto depois de uma falha forçada | Um bloco catch no código |
| O estado vazio renderiza e oferece um próximo passo | O estado visto sem dados | Uma checagem de array vazio |
| Sem rolagem horizontal na largura estreita | Observado na largura declarada | Uma classe utilitária responsiva |
| Conteúdo utilizável na largura larga | Observado na largura declarada | Uma regra de max-width |
| Os inputs têm rótulos programáticos | Nome anunciado ou associação de rótulo observada | Texto de placeholder visível |
| Imagens com significado têm texto alternativo | Atributo presente e preciso | Um nome de arquivo que descreve a imagem |
| Contraste do texto de pelo menos 4.5 para 1 | Razão calculada a partir das cores resolvidas de primeiro plano e de fundo | Uma paleta dita acessível, ou uma razão estimada a olho |
| Nenhum erro não capturado no console durante o fluxo | Console lido depois da execução | Uma captura de tela limpa |
| Os ativos exigidos são entregues com o projeto | Ativo localizado no repositório ou numa dependência | Um recurso que carregou de algum lugar |

| Evidência disponível | O que pode chegar a `passed` | O que o resumo deve dizer |
| --- | --- | --- |
| Navegador, headed ou headless | Só as linhas exercitadas e observadas como válidas | `complete` só depois que todas as linhas visuais e de interação exigidas foram exercitadas; senão, `pending` com as checagens faltantes nomeadas |
| A aplicação roda, sem navegador | Só as linhas de build e unitárias | Verificação visual `pending`, sem navegador |
| A aplicação não roda | Só as linhas estáticas | Verificação `incomplete`, fluxo inalcançável |

## Erros comuns

| Erro | Por que prejudica | Em vez disso |
| --- | --- | --- |
| Reportar `passed` porque a checagem de tipos teve sucesso | A compilação não diz nada sobre se uma pessoa consegue terminar o fluxo | Mantenha as linhas de interação abertas até serem exercitadas |
| Pular a passagem pelo teclado porque o caminho do ponteiro funcionou | Entrega uma superfície que um usuário de teclado não consegue operar de jeito nenhum | Repita o fluxo inteiro só com o teclado |
| Verificar só o caminho feliz | Os estados aplicáveis de carregamento, de erro e vazio podem continuar quebrados | Exercite cada estado exigido com fixtures autorizadas ou controles do navegador |
| Conferir uma só largura de viewport | Defeitos de layout aparecem justamente na largura que ninguém abriu | Registre duas larguras e declare os números |
| Tratar um navegador ausente como `passed` | Converte um desconhecido numa afirmação falsa de concluído | Marque as linhas `not-run` e declare a verificação visual `pending` |
| Buscar um ativo ausente durante a verificação | Esconde um defeito que vai reaparecer para todo usuário real | Reporte o ativo ausente como uma linha `failed` |
| Ignorar a saída do console porque a tela parecia certa | Erros silenciosos viram o defeito reportado pelos usuários na semana seguinte | Leia o console e registre o que ele diz |
| Corrigir o defeito que você achou | A verificação perde a independência e a matriz deixa de refletir a mudança sob teste | Reporte a falha ao autor |
| Escrever "acessibilidade OK" sem linhas | Uma afirmação inverificável que esconde toda barreira real | Liste cada linha de acessibilidade com o próprio resultado |

## Exemplo

Entrada: um formulário de convite com os critérios "um endereço inválido mostra um erro inline" e "a lista mostra um estado vazio antes do primeiro convite".

```
## Verificação de UI
- Superfície: /settings/invites
- Ambiente: npm run dev | navegador: headless
- Viewports: 390 e 1440

## Matriz
| # | Critério                                        | Método         | Resultado | Evidência                                           |
|---|-------------------------------------------------|----------------|-----------|-----------------------------------------------------|
| 1 | Endereço inválido mostra um erro inline         | teclado        | passed    | .harness/tasks/T-88/evidence/invite-error-390.png   |
| 2 | O erro é anunciado, não só por cor              | semântica      | failed    | texto do erro renderizado como um span              |
| 3 | Lista vazia mostra um próximo passo             | ponteiro       | passed    | .harness/tasks/T-88/evidence/invite-empty.png       |
| 4 | Sem rolagem horizontal em 390 px                | viewport 390   | passed    | .harness/tasks/T-88/evidence/invite-390.png         |
| 5 | Formulário utilizável sem cortes em 1440 px     | viewport 1440  | passed    | .harness/tasks/T-88/evidence/invite-1440.png        |
| 6 | Nenhum erro não capturado no console no fluxo   | console        | failed    | trecho do console abaixo                            |

## Falhas
- [O erro é anunciado, não só por cor] A mensagem renderiza num span simples, sem associação com o input.
  - Passos: focar o campo de endereço, digitar "a@", sair com tab.
  - Evidência: .harness/tasks/T-88/evidence/invite-error-390.png
- [Nenhum erro não capturado no console no fluxo] Uma rejeição não tratada dispara no envio enquanto a tela ainda parece correta.
  - Passos: enviar um endereço válido com o painel de rede aberto.
  - Evidência: trecho do console abaixo

## Console
- Erros: 1 - `Uncaught (in promise) TypeError: cannot read properties of undefined (reading 'id')`
- Avisos: 2

## Resumo
- Passaram: 4 | Falharam: 2 | Não rodados: 0
- Verificação visual: complete
- Veredito: rejected
```

## Relacionado

Papéis: `builder`, `qa-verifier`, `coordinator`, `docs-guide`. Comando: `/dh:verify`, e o passo de verificação dentro de `/dh:build`. Skills: `regression-testing` para transformar uma linha `failed` numa checagem automatizada, `code-review` para o diff por trás da superfície, `api-contracts` quando a falha está nos dados que a superfície recebe, `delivery-readiness` para o que um release precisa antes da decisão.

## Caso de prova

Dado um fluxo de fixture com uma armadilha de teclado e um input sem rótulo, a verificação reporta os dois como linhas `failed` com passos de reprodução e evidência, e aprova os critérios que de fato valem. Rodado de novo num ambiente sem navegador, o mesmo procedimento reporta as linhas de interação como `not-run` com o motivo e declara a verificação visual `pending`, nunca uma aprovação geral. Uma execução com navegador e um estado exigido inalcançável também continua `incomplete` e `pending`. Uma superfície sem estados de carregamento ou de erro não ganha critérios inventados nem edições de código-fonte só para exercitá-los.
