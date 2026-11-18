---
name: api-contracts
description: Use quando uma interface pública entre componentes estiver sendo criada ou alterada e os consumidores precisarem concordar com o provedor antes de o código ser escrito, incluindo APIs HTTP, RPC, de eventos, de webhooks e entre módulos, ou quando uma mudança puder quebrar clientes implantados e precisar de um caminho de migração. Não use quando a mudança for interna e sem superfície visível ao consumidor, quando o contrato já existir e só a implementação estiver em questão, ou quando a tarefa for construir o servidor em vez de especificá-lo.
author: malka
metadata:
  provenance: adapted
  sources: ["api-design reference (REST/GraphQL conventions)", " api-designer role description", "dev-harness api-designer agent", "addyosmani/agent-skills (MIT); mechanism adapted, not copied"]
---

# Contratos de API

## Visão geral

Este procedimento produz um contrato que um consumidor e um provedor leem de forma idêntica: operações, catálogo de erros, um exemplo de sucesso, um exemplo de falha e um registro explícito do que quebra. O contrato é escrito antes de o servidor existir, para que a divergência apareça no texto e não na integração. O princípio central é a aplicabilidade: cada preocupação abaixo só é incluída quando um gatilho declarado está presente, e tudo o que for incluído sem gatilho é especulação, não projeto.

## Quando usar

- Uma interface nova de HTTP, RPC, eventos, webhook ou entre módulos precisa ser definida.
- Uma interface existente muda de formato, de códigos de status, de códigos de erro ou de autorização.
- Dois papéis estão prestes a implementar a mesma interface de lados diferentes.
- Um consumidor já implantado pode quebrar e o caminho de migração não está claro.
- Os critérios de aceite dependem do que a interface devolve no sucesso e na falha.

## Quando não usar

- Refactor interno sem superfície visível ao consumidor. Use `safe-refactoring`.
- Um único endpoint conhecido cujo comportamento já está especificado e só precisa ser construído. Use `incremental-implementation`.
- Escolher o transporte ou a própria fronteira do serviço. Use `architecture-decisions`.
- Trabalho puramente de UI contra uma API inalterada.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Autorização e conjunto de escrita para o arquivo do contrato | Não escreva. Devolva o contrato na resposta e peça o caminho e o escopo ao Coordenador. |
| Os consumidores da interface e o que cada um faz com ela | Pergunte ao Coordenador. Um consumidor desconhecido muda a paginação, a autenticação e a compatibilidade. Não adivinhe. |
| Convenções existentes no repositório (caminhos, formato de erro, autenticação, nomenclatura) | Leia dois ou três contratos vizinhos e adote o que eles fazem. Registre como decisão, não como fato. |
| Critérios de aceite que dependem do contrato | Escreva o contrato para as operações que você consegue embasar numa fonte e liste os critérios que não conseguiu mapear em Decisões em aberto. |
| Versão atual e estado de implantação dos consumidores existentes | Trate todo consumidor como implantado. É a premissa segura e custa só uma nota de migração. |

## Procedimento

1. **Leia antes de escrever.** Localize contratos existentes, formatos de erro, middleware de autenticação e convenções de nomenclatura com busca local. Cite cada convenção com um caminho relativo. Um contrato que contradiz os vizinhos é um defeito mesmo quando é internamente consistente.
2. **Nomeie os consumidores.** Liste cada consumidor e a única pergunta que ele faz à interface. Uma operação que nenhum consumidor listado chama não pertence a este contrato.
3. **Modele as operações a partir do domínio, não do armazenamento.** Para cada operação registre intenção, entradas, saída no sucesso, pré-condições e efeitos observáveis. Nomeie as ações pela intenção, não pela mutação de tabela.
4. **Aplique o gate de aplicabilidade.** Percorra a tabela em Referência rápida. Para cada preocupação, ou declare o gatilho que a faz valer e especifique-a, ou registre uma linha dizendo que não se aplica e por quê. Silêncio não é decisão.
5. **Defina um catálogo de erros para a superfície inteira.** Cada erro carrega um código estável legível por máquina, um status de transporte, uma mensagem humana, detalhes opcionais por campo e um identificador de correlação. A mesma falha nunca aparece sob dois códigos.
6. **Escreva um exemplo de sucesso e um exemplo de falha por operação.** Valores completos, concretos e que se possa copiar e colar, sem marcadores como `...`. Esses dois exemplos são o artefato compartilhado: o consumidor constrói seu parser a partir deles e o provedor constrói sua resposta a partir deles.
7. **Registre a compatibilidade.** Classifique cada delta como aditivo, incompatível ou não resolvido, pendente de evidência de consumidor nomeado, usando a tabela de compatibilidade. Não finalize o trabalho que depende de compatibilidade enquanto essa evidência faltar. Para cada item incompatível registre o que quebra, qual consumidor quebra e a sequência de expandir e contrair que evita uma virada simultânea.
8. **Escreva cenários de integração.** Um dado/quando/então por critério de aceite, cada um apontando para uma operação e para um dos dois exemplos. Eles viram as entradas de `regression-testing`.
9. **Pare na especificação.** Não implemente o servidor, o cliente nem a migração neste papel. Devolva o contrato com evidência, limitações e decisões em aberto.
10. **Reporte, não registre.** Fatos, evidência, decisões em aberto e qualquer incidente voltam ao Coordenador, que é o único a escrever `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`. Escreva o contrato só dentro do conjunto de escrita atribuído pelo Coordenador, normalmente `.harness/tasks/<id>/`. Se o projeto guarda especificações em outro lugar, peça esse caminho antes de escrever.

## Formato de saída

```
# Contrato: <nome da interface>

## Consumidores
- <consumidor> -> <o que precisa desta interface>

## Convenções seguidas
- <convenção> (fonte: <caminho relativo>)

## Operações
### <id da operação> <forma de transporte, p. ex. POST /v1/invites>
- Intenção:
- Requisição: <campos, tipos, obrigatório ou opcional>
- Sucesso: <status> <formato da resposta>
- Erros: <código> -> <status> -> <quando>
- Autorização: <regra, ou "não aplicável: motivo">
- Idempotência: <regra, ou "não aplicável: motivo">
- Paginação: <regra, ou "não aplicável: motivo">
- Versionamento: <regra, ou "não aplicável: motivo">
- Limitação de taxa: <regra, ou "não aplicável: motivo">
- Falha parcial: <regra, ou "não aplicável: motivo">

## Catálogo de erros
| código | status | significado | detalhes do campo | id de correlação |
| --- | --- | --- | --- | --- |

## Exemplos
### Sucesso
<requisição e resposta completas>
### Falha
<requisição e resposta completas>

## Compatibilidade
- Aditivas: <itens>
- Incompatíveis: <item> | <consumidor afetado> | <passo de migração>
- Não resolvidas: <item> | <consumidor> | <evidência necessária antes de finalizar a compatibilidade>
- Sequência: expandir -> migrar consumidores -> contrair
- Retirada: não aplicável | <superfície> | consumidores: <contagem medida ou desconhecida> | substituto: <nome ou "sem substituto", decisão do dono>

## Cenários de integração
- Dado <estado>, quando <chamada>, então <resultado observável> (critério: <id>)

## Decisões em aberto
- <pergunta> | <quem decide> | <o que muda com cada resposta>

## Evidência
- <caminho relativo> - <o que estabelece>
```

## Referência rápida

Gate de aplicabilidade. Inclua a preocupação só quando o gatilho valer.

| Preocupação | Incluir quando | Registrar quando omitida |
| --- | --- | --- |
| Autorização | A operação acessa dados protegidos ou exige um papel, capacidade, direito ou permissão em nível de objeto, inclusive sobre os dados do próprio chamador | "Operação pública irrestrita: <motivo pelo qual nenhuma permissão é exigida>" |
| Idempotência | Uma nova tentativa é realista e o efeito não é naturalmente repetível (criação, cobrança, despacho) | "GET, naturalmente idempotente" |
| Paginação | A coleção cresce sem um limite que o chamador controle | "Conjunto pequeno e fixo, limitado por projeto" |
| Versionamento | Ao menos um consumidor já está implantado contra a forma antiga | "Primeiro release, nenhum consumidor implantado" |
| Limitação de taxa | A superfície é pública, cara ou abusável | "Interna, atrás de um chamador confiável" |
| Semântica de falha parcial | A operação toca mais de um recurso ou sistema | "Recurso único, atômico" |

Classificação de compatibilidade.

| Mudança | Classe |
| --- | --- |
| Campo de requisição opcional novo ou operação nova | Aditiva quando as chamadas existentes mantêm seu significado e comportamento |
| Campo de resposta novo | Aditiva só quando o contrato do consumidor ou os consumidores verificados toleram campos desconhecidos; caso contrário, incompatível, ou não resolvida até ser conferida |
| Código de erro novo para uma condição antes não tratada | Aditiva se os clientes têm um ramo padrão, caso contrário incompatível |
| Valor novo de enum | Incompatível, a menos que esteja especificado que os consumidores toleram valores desconhecidos |
| Remover ou renomear um campo, operação ou código de erro | Incompatível |
| Tornar obrigatório um campo de requisição opcional, ou apertar a validação | Incompatível |
| Mudar um tipo, um código de status, um padrão ou uma unidade | Incompatível |
| Mudar os padrões de paginação ou a ordenação | Incompatível |

Expandir e contrair, a única sequência que evita uma virada simultânea: adicione a forma nova ao lado da antiga, sirva as duas, migre cada consumidor com evidência, anuncie a data de remoção e então remova.

Retirada, só quando a mudança remove ou descontinua uma operação, campo ou código de erro visível ao consumidor. Meça o uso primeiro: pontos de chamada, mais qualquer consumidor implantado nomeado nas entradas. Uma contagem desconhecida continua desconhecida. Não escreva desconhecido como zero. Nomeie um substituto comprovado, ou escreva "sem substituto" como decisão em aberto para o dono. O padrão é consultivo. Só é obrigatório com uma data, um guia de migração com passos concretos e uma forma de encontrar os chamadores restantes. Não adicione comportamento a uma superfície marcada para remoção. A remoção de coluna e de esquema continua em `data-migrations`. Esta seção não repete esse procedimento.

Não rotule como aditiva uma checagem de compatibilidade não resolvida. Nomeie o consumidor afetado e a evidência necessária para decidir. Por exemplo, um consumidor que rejeita propriedades de resposta desconhecidas pode quebrar quando um campo é adicionado. Propriedade, sozinha, também não resolve a autorização: mudar a própria conta ainda pode exigir um papel ou um direito.

## Erros comuns

| Erro | Por que prejudica | Faça em vez disso |
| --- | --- | --- |
| Campos adicionados "para depois" | Todo campo especulativo vira uma obrigação de compatibilidade no instante em que um consumidor o lê | Só campos de que um consumidor nomeado precisa agora |
| Exemplos com marcadores | Consumidor e provedor preenchem as lacunas de jeitos diferentes e a divergência aparece na integração | Valores completos e concretos dos dois lados |
| Dois códigos de erro para a mesma causa | Os clientes ramificam pelo status de transporte e perdem a distinção de qualquer forma | Um código estável legível por máquina por causa |
| 5xx para um erro do chamador | Esconde bugs do cliente dentro dos alertas do provedor e anula a política de nova tentativa | 4xx para o chamador, 5xx só para falha do provedor |
| Contrato moldado como as tabelas do banco | Refatorações do armazenamento viram mudanças incompatíveis na API | Modele a intenção do domínio |
| Mudança incompatível com um aumento de versão e mais nada | Um número de versão não é uma migração; os clientes implantados continuam quebrando | Registre o consumidor afetado mais os passos de expandir e contrair |
| Escrever o servidor enquanto escreve a especificação | A especificação deixa de ser revisável e vira uma descrição do código | Só especificação, a implementação é outro papel |
| Tratar um rascunho como implantado | O trabalho a jusante planeja contra algo que não existe | Marque o status explicitamente |

## Exemplo

Entrada: "Adicione um endpoint de convites. O frontend cria convites, o worker tenta de novo em caso de timeout."

```
### create_invite  POST /v1/invites
- Intenção: convidar um endereço de e-mail para uma organização.
- Requisição: { "email": string (obrigatório), "role": "member" | "admin" (obrigatório) }
- Sucesso: 201 { "id": "inv_7Kq", "email": "sam@example.com",
                 "role": "member", "status": "pending" }
- Erros: invalid_email -> 422 | already_invited -> 409 | forbidden -> 403
- Autorização: o chamador deve ter admin na organização de destino.
- Idempotência: cabeçalho Idempotency-Key obrigatório; um reenvio devolve o corpo 201 original.
- Paginação: não aplicável, recurso único.

### Exemplo de falha
POST /v1/invites   { "email": "sam@", "role": "member" }
422 { "error": { "code": "invalid_email", "message": "O e-mail é inválido",
      "details": [{ "field": "email", "issue": "format" }],
      "request_id": "req_abc123" } }

## Compatibilidade
- Aditivas: create_invite é uma operação nova, sem consumidor implantado

## Cenários de integração
- Dada uma organização com um convite pendente para sam@example.com, quando create_invite é
  chamada de novo, então 409 already_invited (critério: AC-2)

## Decisões em aberto
- Expiração de um convite pendente | dono do produto | uma expiração adiciona um campo e um
  segundo motivo para o 409 desaparecer
```

A idempotência está presente porque o worker tenta de novo. A paginação é omitida com um motivo declarado. É a partir do exemplo de falha que o frontend constrói seu ramo de erro.

## Relacionado

Papéis: `architect`, `builder`, `qa-verifier`. Comandos: usados dentro de `/dh:plan`. Skills: `architecture-decisions`, `implementation-planning`, `incremental-implementation`, `regression-testing`.

## Caso de prova

Dado um contrato produzido por este procedimento, um papel de consumidor e um papel de provedor descrevem, de forma independente, o mesmo corpo de sucesso e o mesmo corpo de falha, incluindo status, código de erro e nomes de campo, sem mais esclarecimentos. Uma mudança que remove um campo de resposta é classificada como incompatível, nomeia o consumidor afetado e carrega uma sequência de expandir e contrair em vez de um mero aumento de versão. Adicionar um campo de resposta para um consumidor que rejeita propriedades desconhecidas também é incompatível; sem evidência de tolerância, permanece não resolvido. Uma operação sobre a conta do próprio chamador que exige um direito mantém essa regra de autorização. A remoção de uma operação visível ao consumidor nomeia uma contagem de consumidores, medida ou explicitamente desconhecida, e um substituto ou uma decisão explícita de não ter substituto. Desconhecido não é escrito como zero. A remoção de uma coluna fica com `data-migrations`.
