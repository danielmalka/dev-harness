---
name: requirements-discovery
description: Use quando um pedido nomeia uma dor, um público ou um resultado, mas nenhum comportamento testável, quando duas leituras do mesmo pedido produziriam softwares diferentes, ou quando critérios de aceite precisam existir antes que um plano possa ser escrito. Não use quando o aceite já existe, quando o usuário pediu um plano, um patch ou uma revisão, ou quando a única peça que falta é saber onde o código está.
author: malka
metadata:
  provenance: adapted
  sources: ["product-discovery role", "brainstorming", "grilling", " discovery-agent", "addyosmani/agent-skills (MIT); mechanism adapted, not copied"]
---

# Descoberta de requisitos

## Visão geral

Transforme um pedido vago em cinco coisas em que um construtor possa agir: a situação atual, o resultado desejado, exemplos concretos, as alternativas consideradas e os limites. A saída é um resumo curto o bastante para ler de uma vez e específico o bastante para planejar a partir dele. Perguntas custam a atenção do usuário, então só vale perguntar sobre decisões que mudariam o que será construído; tudo o que dá para encontrar é encontrado, não perguntado. Esta skill escreve só artefatos de descoberta: resumo, aceite, decisões em aberto. Nunca código-fonte da aplicação, configuração ou arquivos de teste.

## Quando usar

- O pedido descreve uma dor, um usuário ou um objetivo, mas nenhum comportamento observável.
- Duas leituras honestas da mesma frase levariam a construções diferentes.
- Ninguém sabe dizer como é o "pronto" sem inventá-lo no meio do build.
- O pedido é grande o bastante para ser várias peças independentes em vez de uma só.
- Uma tentativa anterior entregou a coisa errada e a causa foi ambiguidade, não execução.

## Quando não usar

- Os critérios de aceite já existem e são testáveis. Vá para `implementation-planning`.
- O usuário pediu um patch, um plano, uma revisão ou um mapa.
- O desconhecido é a base de código, não a intenção. Use `repository-mapping`.
- Um resumo existe e o objetivo não mudou. Não reabra um resumo fechado para polimento.

## Entradas

| Entrada | Se faltar |
|---|---|
| O objetivo declarado, nas palavras do usuário | Não há nada a descobrir. Na sessão principal, pergunte ao dono o que deveria passar a ser possível e hoje não é. Como especialista despachado, devolva a entrada que falta ao Coordenador e pare; nunca prossiga com um valor presumido. |
| Quem é afetado | Faça desta a primeira pergunta, pois ela muda todas as respostas seguintes. |
| Restrições conhecidas: prazo, stack, fora de escopo | Prossiga e registre-as como desconhecidas; não invente prazo nem stack. |
| Docs de produto existentes ou um resumo anterior no projeto | Busque antes de perguntar. Uma resposta que já está escrita não é uma pergunta. |

## Procedimento

1. **Dimensione o pedido antes de refiná-lo.** Se o pedido contém vários subsistemas independentes, diga isso de imediato e proponha dividi-lo. Refinar detalhes dentro de um pedido que precisa ser decomposto desperdiça todas as perguntas seguintes. Cada peça recebe então o seu próprio resumo.

2. **Separe fatos de decisões e nunca misture os dois.** Fatos sobre o repositório, o comportamento existente, o formato atual dos dados e o estado da arte cabem a você, com leitura e busca. Decisões sobre a quem isto serve, o que acontece no caso incômodo e o que fica fora de escopo pertencem ao usuário. Perguntar ao usuário um fato que você poderia consultar gasta a atenção dele na coisa errada; responder uma decisão por conta própria é o modo de falha que produz software confiante e errado.

3. **Monte a fronteira de perguntas.** Liste toda decisão em aberto e mantenha só as que já têm os pré-requisitos resolvidos. Uma pergunta que depende de uma pergunta sem resposta pertence a uma rodada posterior. Ordene as sobreviventes pelo quanto mudam a construção.

4. **Aplique o teste de escopo a cada pergunta candidata.** Duas respostas diferentes produziriam critérios de aceite, arquivos ou testes diferentes? Se não, descarte a pergunta ou decida você mesmo e rotule a escolha como sugestão. Curiosidade não é motivo para perguntar.

5. **Faça a menor pergunta útil e espere as decisões necessárias.** Inclua uma recomendação, mas o silêncio nunca é uma resposta nem uma autorização. Continue a pesquisa independente enquanto uma decisão de escopo bloqueante estiver pendente. Duas perguntas só podem dividir uma mensagem quando a segunda não faz sentido sem a primeira; fora isso, uma pergunta bloqueante por vez. Registre no resumo as respostas reais.

6. **Preencha as cinco dimensões.** Situação atual: o que de fato acontece hoje, com evidência onde o projeto puder mostrá-la. Resultado desejado: em comportamento, não em sentimentos. Exemplos: pelo menos um caso normal e um caso incômodo. Alternativas: incluindo não fazer nada e a versão parcial mais barata. Limites: o que está explicitamente fora de escopo, as restrições inegociáveis e os dados ou permissões envolvidos. Registre como hipótese cada afirmação sobre o comportamento atual que você não pôde verificar, com a observação que a refutaria.

7. **Escreva o aceite como "quando X, então Y".** Cada item observável por uma pessoa ou por um teste, com X nomeando um gatilho e Y nomeando um resultado. Marque cada um como `requirement` ou `suggestion`. Um `requirement` é algo que o usuário pediu ou que a correção exige; tudo o que você pensou e ninguém exige é uma `suggestion`, e entregar sugestões como requisitos é aumento de escopo com cara séria.

8. **Cubra os casos incômodos de forma proporcional.** Vazio, duplicado, expirado, não autorizado, concorrente e grande demais. Inclua só os que a funcionalidade de fato encontra. Um caso-limite inventado é escopo, não diligência.

9. **Registre o que ainda está em aberto.** Cada decisão em aberto recebe a pergunta, por que ela muda o escopo e a sua recomendação. Não resolva uma decisão de produto em aberto escolhendo em silêncio; também não trave o resumo inteiro por causa de um item sem resposta. Se uma pergunta não pode ser resolvida conversando porque precisa de algo para olhar, diga isso e encaminhe-a a um protótipo em vez de reformulá-la para sempre.

10. **Pare quando o aceite for testável e as decisões em aberto estiverem explícitas.** Entregue o resumo ao Coordenador, junto com qualquer restrição ou incidente que valha lembrar. Só o Coordenador escreve `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`; a descoberta reporta fatos e nunca edita esses registros.

## Formato de saída

```
## Resumo
- Problema: <o que está errado hoje, observável>
- Usuários afetados: <quem, e o que fazem hoje em vez disso>
- Comportamento desejado: <o que passa a ser possível>
- Fora de escopo: <explicitamente não agora>

## Exemplos
- Normal: dado <entrada>, o usuário vê <resultado>
- Incômodo: dada <entrada-limite>, o usuário vê <resultado>

## Alternativas consideradas
| Opção | Custo | Por que não foi escolhida |
|---|---|---|
| Não fazer nada | | |
| <versão parcial mais barata> | | |

## Limites e restrições
- <restrição> (fonte: usuário | instruções do projeto | observado)

## Aceite
- Quando <X>, então <Y>. [requirement]
- Quando <X>, então <Y>. [suggestion]

## Hipóteses
- <afirmação>, como seria mostrada errada

## Decisões em aberto
- <pergunta>, por que muda o escopo, resposta recomendada

## Evidência
- <caminho relativo> ou "ainda nenhuma"
```

Quando persistido, o resumo vira `.harness/prd/PRD-<n>.md` a partir do template do kit `templates/PRD.md`; mantenha a marcação `requirement`/`suggestion` no texto.

## Regras de escrita do PRD
O PRD é escrito em pt-BR (ADR-011). Estas sete regras valem para o texto do PRD e são próprias do kit, em português; o `document-validator` as confere com a mesma severidade de sempre: ambiguidade bloqueia, estilo puro é Minor.
1. Frase curta, de até cerca de 25 palavras, com uma ação por frase.
2. Voz ativa e sujeito explícito ("o `dh` grava", não "é gravado").
3. Um termo por conceito, sem sinônimo: o termo definido uma vez é o termo do documento inteiro.
4. Verbo de obrigação com o mesmo peso sempre: "deve" (obrigatório), "pode" (opcional), "não deve" (proibido).
5. Número com unidade e limite fechado ("até 200 linhas", "no máximo 2 rodadas").
6. Sem "etc.", "e afins", "entre outros", "quando possível", "adequado" ou "rápido" sem número.
7. Lista para passos e critérios; prosa só para contexto.

Uma linha de comando `Verifiable:` fica isenta do limite de tamanho da regra 1.

## Referência rápida

| Frase vaga | A pergunta que a torna testável |
|---|---|
| "deixe mais rápido" | Qual operação, medida como, e qual número é aceitável? |
| "os usuários ficam confusos" | Em que passo eles param, e o que deveria acontecer em vez disso? |
| "tem que ser seguro" | Quais dados, qual ator deve ser recusado, e como é a recusa? |
| "adicione relatórios" | Que decisão o leitor toma a partir do relatório? |
| "tem que simplesmente funcionar" | Qual caso falha hoje, e qual é o resultado correto nele? |
| "como <outro produto>" | Qual comportamento específico dele, exatamente? |

| Teste | Faça a pergunta | Não faça a pergunta |
|---|---|---|
| Teste de escopo | Duas respostas mudam a construção | Duas respostas mudam só a redação |
| Teste de propriedade | A decisão é do usuário | Dá para descobrir no repositório |
| Teste de fronteira | Os pré-requisitos estão resolvidos | Depende de uma pergunta em aberto |
| Teste de custo | A resposta é mais barata que construir a coisa errada | A resposta é adivinhável e de baixo risco |

## Erros comuns

| Erro | Por que prejudica | Faça em vez disso |
|---|---|---|
| Perguntar tudo de uma vez para interromper só uma vez | Uma lista de perguntas recebe uma resposta e um usuário cansado | Uma pergunta bloqueante por vez, ordenada por impacto |
| Tratar a pergunta sobre a stack como descoberta | Ela só muda a construção quando o usuário tem uma restrição real | Deixe uma stack sem restrição para o planejamento |
| Acrescentar um fluxo mais bonito em que você pensou | Sua ideia vira escopo que ninguém pediu | Rotule-a como sugestão e deixe o usuário promovê-la |
| Rodar a descoberta quando o escopo já está claro | Um resumo resolvido é reaberto e perguntas fechadas voltam | Vá para o planejamento |
| Tratar um pedido como um resumo | Subsistemas independentes são refinados antes de serem divididos | Decomponha primeiro, depois um resumo por peça |

## Racionalizações

| O que você vai ficar tentado a pensar | Por que está errado | O que fazer |
| --- | --- | --- |
| "Eu já entendi, vou escrever o PRD sem a pergunta que muda o software" | Duas leituras dessa frase produzem softwares diferentes, e o resumo congela a errada | Faça essa pergunta e espere |
| "O aceite é óbvio, vou pular a escrita" | O aceite óbvio é o que duas pessoas leem de modo diferente | Escreva cada item como "quando X, então Y" |
| "Eles não mencionaram erros, então não há nenhum" | Não declarado não quer dizer ausente | Nomeie o caso incômodo e obtenha uma decisão sobre ele |
| "Não consegui falar com o usuário, então decidi por ele" | Uma decisão silenciosa vira contrato | Registre uma decisão em aberto com uma recomendação e marque o resumo como provisório |

## Exemplo

Entrada: "as pessoas esquecem de dar retorno depois das demos".

```
## Resumo
- Problema: depois de uma demo, o próximo contato depende de o vendedor lembrar; nada no produto registra que uma demo aconteceu.
- Usuários afetados: vendedores, que hoje guardam anotações fora do produto.
- Comportamento desejado: uma demo pode ser marcada como realizada, e uma sem resposta fica visível para o seu dono.
- Fora de escopo: enviar e-mail em nome do usuário.

## Exemplos
- Normal: o vendedor marca uma demo como realizada; dois dias sem resposta registrada a colocam na lista de pendências dele.
- Incômodo: o contato responde no mesmo dia; a demo nunca chega à lista de pendências.

## Alternativas consideradas
| Opção | Custo | Por que não foi escolhida |
|---|---|---|
| Não fazer nada | nenhum | o retorno esquecido é a dor reportada |

## Limites e restrições
- Nenhum e-mail enviado em nome do usuário (fonte: usuário)

## Aceite
- Quando uma demo é marcada como realizada e nenhuma resposta é registrada dentro da janela configurada, então ela aparece na lista de pendências do dono. [requirement]
- Quando uma resposta é registrada, então o item sai da lista no próximo carregamento. [requirement]
- Quando o dono descarta um item, então ele não volta para aquela demo. [suggestion]

## Hipóteses
- As respostas já são registradas junto ao contato. Seria mostrada errada por um contato com uma resposta que nunca chega ao produto.

## Decisões em aberto
- A janela é fixa ou por usuário? Muda o armazenamento e as configurações. Recomendado: fixa no começo, configurável depois.

## Evidência
- ainda nenhuma
```

## Relacionado

Papéis: `product-discovery`, `coordinator`, `implementation-planner`, `qa-verifier`. Comando: `/dh:discover`. Skills: `implementation-planning` consome este resumo, `architecture-decisions` quando uma restrição daqui força uma escolha durável, `regression-testing` transforma os itens de aceite em checagens, `repository-mapping` quando uma afirmação sobre o comportamento atual precisa de evidência.

## Caso de prova

Um pedido vago de uma frase produz um resumo com pelo menos um exemplo normal e um incômodo, aceite escrito como "quando X, então Y" com `requirement` e `suggestion` marcados separadamente, e pelo menos uma decisão em aberto enunciada com o motivo de ela mudar o escopo. A sessão não faz nenhuma pergunta cujas duas respostas possíveis produziriam a mesma construção, e não pergunta nada que já tenha sido respondido no repositório. O resumo não é escrito enquanto houver uma pergunta que mudaria o software ainda não feita, e o aceite é escrito mesmo quando parece óbvio.
