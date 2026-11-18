---
name: data-migrations
description: Use quando uma mudança alterar estrutura ou dados persistidos, incluindo tabelas e colunas novas, mudanças de tipo ou de restrição, renomeações, remoções, preenchimentos retroativos, mudanças de índice e as consultas que dependem deles. Use também quando uma mudança de esquema precisar ir ao ar sem indisponibilidade, ou quando alguém precisar saber se uma migração pode ser revertida. Não use quando a mudança for numa consulta que deixa o esquema intacto, e nunca use para aplicar mudanças num banco de produção.
author: malka
metadata:
  provenance: adapted
  sources: [" postgres-pro", "parallel change practice"]
---

# Migrações de dados

## Visão geral

Projete uma mudança de persistência como uma sequência de passos, cada um seguro por si só, aplicados numa ordem em que a aplicação em execução fica correta antes, durante e depois de cada passo. A estrutura avança por expandir e contrair: adicione a forma nova, escreva nas duas, mova os leitores e só então remova a forma antiga, quando nada mais a ler. Toda migração sai com um plano de reversão, ou com uma declaração explícita de que a reversão é impossível e por quê. A verificação acontece num banco descartável com dados sintéticos. Nunca contra dados reais, e nunca contra um banco cujo endereço veio de um arquivo de ambiente, de um padrão de configuração ou de uma sessão anterior. Aplicar num ambiente compartilhado ou real é decisão do dono, fora deste procedimento.

## Quando usar

- Uma funcionalidade precisa de uma tabela, coluna, índice, restrição, valor de enumeração ou relacionamento novo.
- Uma coluna existente muda de tipo, de nulidade, de valor padrão ou de significado.
- Dados precisam de preenchimento retroativo, correção, deduplicação ou movimentação entre tabelas.
- Uma estrutura está sendo removida ou renomeada e há consumidores que ainda a referenciam.
- Uma consulta ficou lenta e a correção é um índice, uma desnormalização ou uma decisão de particionamento.

## Quando não usar

- A mudança reescreve uma consulta contra um esquema que não está mudando. Isso é trabalho comum de implementação.
- A tarefa é operar um banco de dados: backups, replicação, failover, capacidade. Isso é operação, não este procedimento.
- Alguém quer uma migração executada num ambiente compartilhado ou de produção. Isso é decisão do dono, com autorização própria, fora deste procedimento.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| O comportamento que a mudança deve suportar, com seus critérios de aceite | Pare e obtenha. Um esquema projetado a partir de um palpite sobrevive ao palpite. |
| O esquema atual: tabelas, colunas, tipos, restrições, índices | Leia do histórico de migrações do repositório. Não conecte a um ambiente para descobrir. |
| Os padrões de acesso: quais consultas leem e escrevem esses dados, e com que frequência | Liste os pontos de chamada encontrados no código e marque os consumidores não resolvidos como desconhecidos. |
| Volume e crescimento dos dados, ao menos em ordens de grandeza | Presuma que a tabela é grande o bastante para exigir lotes e diga que presumiu. |
| O modelo de implantação: se o código antigo e o novo da aplicação rodam ao mesmo tempo | Presuma que se sobrepõem. É a premissa segura e força o expandir e contrair. |
| Um banco descartável para testar | Reporte a migração como não testada e diga o que não pôde ser verificado. Não substitua por um ambiente compartilhado. |

## Procedimento

1. **Declare a forma-alvo e a invariante que ela protege.** Escreva a estrutura pretendida e, para cada restrição, a regra que ela impõe em termos de negócio. Uma restrição sem regra declarada ou está faltando ou é arbitrária. Prefira expressar a integridade no banco onde o banco puder impô-la, porque as checagens da aplicação não sobrevivem a um segundo escritor.

2. **Mapeie os consumidores.** Encontre toda leitura e escrita da estrutura afetada no repositório: consultas, modelos de ORM, serializadores, relatórios, jobs em segundo plano, exportações e fixtures. Cada uma vira uma linha no plano da mudança. O consumidor que você não encontrou é o que quebra.

3. **Classifique a mudança.** Aditiva, mutante ou destrutiva. Mudanças aditivas são seguras quando aceitam nulo ou têm valor padrão. Mudanças mutantes precisam de uma janela de compatibilidade. Mudanças destrutivas precisam de uma decisão explícita, porque são as que não podem ser desfeitas.

4. **Decomponha em expandir e contrair.** Para qualquer coisa além de uma mudança puramente aditiva, escreva a sequência:
   - Expandir: adicione a estrutura nova ao lado da antiga, aceitando nulo ou com valor padrão, sem nenhum consumidor dependendo dela.
   - Escrita dupla: a aplicação escreve as duas formas. Inventarie as instâncias antigas da aplicação, os workers e os outros escritores que ainda atualizam só a forma antiga. Antes de concluir o preenchimento retroativo, aposente esses escritores ou estabeleça uma sincronização que mantenha as duas formas atualizadas para as escritas deles.
   - Preenchimento retroativo: preencha a forma nova para as linhas existentes, em lotes, reiniciável, sem transação de longa duração. Impeça que um preenchimento retroativo sobrescreva uma escrita concorrente mais nova; declare a regra de sincronização ou de atualização condicional que a preserva.
   - Migrar leitores: primeiro verifique que todos os escritores ativos mantêm as duas formas e que as linhas existentes convergiram. Depois mova os leitores para a forma nova, um de cada vez, verificando cada um. Um cursor de lote concluído, sozinho, não estabelece convergência.
   - Impor: adicione a restrição ou a exigência de não nulo quando os dados a sustentarem.
   - Contrair: remova a escrita dupla e depois a estrutura antiga, num deploy posterior, quando nada mais a ler.
   Cada fase é uma migração separada e, em geral, um release separado. Juntá-las é o que produz a indisponibilidade.

5. **Planeje as travas.** Para cada passo, declare o que ele trava e por quanto tempo. Adicionar uma coluna com valor padrão volátil, mudar o tipo de uma coluna e adicionar uma restrição que valida as linhas existentes são os ofensores de sempre; eles podem segurar uma trava de escrita na tabela inteira. Prefira a forma sem bloqueio que o motor oferece: crie o índice de forma concorrente, adicione a restrição como não validada e valide num segundo passo, adicione a coluna sem reescrita e preencha à parte. Se um passo precisar tomar uma trava forte, diga, limite-a e defina um timeout de trava para que ele falhe rápido em vez de enfileirar todas as outras consultas atrás dele.

6. **Planeje as transações.** Mudanças de esquema vão numa transação quando o motor suporta definição de dados transacional, para que uma falha não deixe nada aplicado pela metade. Preenchimentos retroativos vão fora, em lotes, porque uma única transação sobre milhões de linhas bloqueia a limpeza e arrisca esgotar o log. Diga qual dos dois cada passo é. Nunca misture um preenchimento retroativo longo na mesma transação de uma mudança estrutural.

7. **Planeje os índices com critério.** Adicione um índice porque uma consulta nomeada precisa dele, e diga qual. Verifique se um índice existente já atende ao predicado como prefixo. Construa de forma concorrente onde o motor permitir, e saiba que a construção concorrente pode falhar e deixar um índice inválido que precisa ser removido e refeito. Remover um índice também é uma migração: confirme que nenhuma consulta depende dele antes de removê-lo.

8. **Escreva o preenchimento retroativo como um job reiniciável.** Lotes limitados, ordenados por uma chave estável, com uma marca d'água registrada para que uma execução interrompida retome em vez de recomeçar. Uma pausa entre os lotes para que o tráfego comum mantenha sua parte do banco. Idempotente por linha, para que rodar um lote de novo seja inofensivo. Inclua a consulta que informa quantas linhas faltam.

9. **Escreva o plano de reversão, ou declare que não há.** Para cada passo, dê o inverso: remover a coluna adicionada, restaurar o valor padrão anterior, recriar o índice removido. Onde o inverso não puder restaurar os dados, diga isso com essas palavras. Remover uma coluna, truncar, deduplicar, estreitar um tipo e sobrescrever um valor durante um preenchimento retroativo são passos sem volta. Para cada passo sem volta, registre o que seria perdido, se existe um backup que o dono realmente verificou e qual é o caminho de recuperação para a frente. Um backup não verificado não é plano de reversão.

10. **Atualize as consultas afetadas junto com o esquema.** Reescreva as leituras e escritas encontradas no passo 2 para que funcionem na janela de compatibilidade: tolerantes a que a coluna nova seja nula, sem selecionar todas as colunas às cegas, sem depender da ordem das colunas. Diga quais consultas mudam em qual fase e verifique o plano de execução das que mudam de custo por causa de um índice novo ou de um tipo alterado.

11. **Teste num banco descartável com dados sintéticos.** Se o projeto já fornece um banco descartável e o despacho autoriza rodá-lo, crie o esquema a partir do histórico de migrações e popule-o com dados gerados que cubram as formas que importam, incluindo nulos, duplicatas, valores de fronteira e o volume necessário para ver a divisão em lotes funcionar. Caso contrário, registre cada linha de verificação como não rodada e reporte a migração como não testada. Não instale nem inicie um serviço de banco de dados, e não provisione um, para obter um alvo de teste. Com um alvo disponível:
    - Aplique a migração para a frente e confirme a estrutura e as restrições.
    - Confirme a invariante: as linhas que devem existir, as linhas que devem ser rejeitadas.
    - Rode as consultas afetadas e confirme os resultados e, onde for relevante, o plano de execução.
    - Aplique a reversão e confirme que o banco volta a um estado anterior utilizável, ou registre que não volta.
    - Interrompa o preenchimento retroativo no meio e reinicie-o, confirmando que retoma sem duplicar trabalho.
    - Para migrações de forma dupla, exercite escritores sobrepostos e um preenchimento retroativo em corrida com uma atualização. Confirme que escritas antigas desatualizadas, só na forma antiga, são impedidas ou sincronizadas e que o valor mais novo sobrevive antes de os leitores trocarem.
    Nunca popule com uma cópia de dados reais que contenha informação pessoal, e nunca aponte a execução de teste para um banco que você não criou para o teste.

12. **Escreva o plano de aplicação.** Passos ordenados, cada um com o que faz, sua duração esperada no volume real, o que trava, como saber que está progredindo e o ponto sem volta. Nomeie a métrica a observar durante cada passo: atraso de replicação, esperas por trava, taxa de erro, quantidade de linhas restantes. Nomeie a condição de aborto de cada passo e a ação que ela dispara.

13. **Pare na decisão.** Não aplique uma migração em nenhum ambiente real automaticamente e não conecte a um banco cujo endereço você inferiu de um arquivo de ambiente, de um padrão de configuração ou de uma sessão anterior. Entregue ao dono o plano, a evidência da execução no banco descartável e os riscos residuais. Só o Coordenador escreve `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`; reporte um passo sem volta, uma reversão que falhou ou um quase-acidente de perda de dados como incidente para o Coordenador registrar.

## Formato de saída

```
## Mudança
- Objetivo: <comportamento que isto suporta>
- Forma atual: <tabelas e colunas afetadas>
- Forma-alvo: <no que se transforma>
- Invariantes: <regra> imposta por <restrição>
- Classificação: aditiva | mutante | destrutiva

## Consumidores
| Consumidor (caminho relativo) | Lê | Escreve | Fase em que muda |
|---|---|---|---|

## Fases
| # | Passo | Tipo | Travas | Transação | Reversível |
|---|-------|------|--------|-----------|------------|
| 1 | <expandir: adicionar coluna> | esquema | nenhuma | sim | sim: remover a coluna |
| 2 | <escrita dupla> | código | nenhuma | n/d | sim |
| 3 | <preenchimento retroativo em lotes de N> | dados | nenhuma | por lote | não: sobrescreve <campo> |

## Reversão
- Por passo: <inverso ou "nenhum">
- Passos sem volta: <passo> - <o que se perde> - <backup verificado? sim/não> - <recuperação para a frente>

## Verificação (banco descartável)
- Semente: <quais dados sintéticos foram gerados>
- <checagem> - `passed` | `failed` | `not-run`
- Para a frente: <resultado>
- Reversão: <resultado, ou "impossível: motivo">
- Preenchimento retroativo interrompido retomou: <resultado>
- Barreira de escritores e convergência dos dados: <evidência de que todo escritor ativo mantém as duas formas e que não restam linhas desatualizadas>
- Atualização concorrente preservada: <resultado>

## Plano de aplicação
| Ordem | Passo | Duração esperada | Observar | Abortar se |
|-------|-------|------------------|----------|------------|

## Riscos residuais
- <risco e quem deve decidir>
```

## Referência rápida

| Mudança | Forma segura | Armadilha |
| --- | --- | --- |
| Adicionar coluna | Aceitando nulo, ou com um valor padrão constante que o motor aplica sem reescrita | Um valor padrão volátil reescreve toda linha sob uma trava de escrita |
| Adicionar coluna não nula | Adicione aceitando nulo, preencha retroativamente e só então imponha | Adicionar não nula num passo só bloqueia e falha nas linhas existentes |
| Renomear coluna | Adicione a nova, escrita dupla, migre os leitores, remova a antiga | Uma renomeação quebra na hora toda instância implantada do código antigo |
| Mudar tipo | Coluna nova com o tipo novo, preenchimento retroativo, troque os leitores, remova a antiga | Uma mudança de tipo no lugar reescreve e trava a tabela |
| Adicionar restrição | Adicione como não validada e valide à parte | Validar na criação varre a tabela inteira sob uma trava |
| Adicionar índice | Construa de forma concorrente, verifique que é válido | Uma construção simples bloqueia escritas; uma construção concorrente que falha deixa um índice inválido |
| Remover coluna ou tabela | Pare de ler, publique e só depois remova num release posterior | Remover no mesmo release quebra as instâncias que ainda rodam o código antigo |
| Preenchimento retroativo | Lotes limitados e reiniciáveis com uma marca d'água | Uma transação sobre a tabela inteira bloqueia a limpeza e pode esgotar o log |
| Excluir ou deduplicar linhas | Decisão do dono, backup verificado, antes uma tabela intermediária reversível | Não há inverso; este é o passo que perde dados |

| Pergunta | Responda antes de escrever a migração |
| --- | --- |
| O código antigo e o novo da aplicação podem rodar ao mesmo tempo? | Presuma que sim, projete para a sobreposição |
| O que cada passo trava, e por quanto tempo? | Declarado por passo, com um timeout de trava |
| Existe um inverso para cada passo? | Declarado por passo, inclusive "nenhum" |
| Qual o tamanho da tabela em produção? | Ao menos uma ordem de grandeza, ou use lotes de qualquer jeito |
| Quais consultas mudam de plano de execução por causa disto? | Listadas com a fase em que mudam |

## Erros comuns

| Erro | Por que prejudica | Faça em vez disso |
| --- | --- | --- |
| Publicar expandir e contrair num release só | Toda instância que roda o código anterior quebra no instante em que a migração entra | Releases separados, com a janela de escrita dupla entre eles |
| Renomear uma coluna no lugar | Quebra na hora o código implantado e todo plano de consulta em cache que a nomeava | Adicione, escrita dupla, migre os leitores, remova depois |
| Tratar um backup não verificado como plano de reversão | O caminho de recuperação se revela quebrado durante o incidente | Exija uma restauração que o dono tenha de fato testado |
| Conectar a um banco achado num arquivo de configuração | Roda uma migração num ambiente que ninguém autorizou | Use só um banco criado para este teste |
| Testar com uma cópia dos dados de produção | Espalha dados pessoais para o desenvolvimento e para a evidência | Gere dados sintéticos que cubram as mesmas formas |
| Deixar as consultas afetadas para depois | O esquema entra e a aplicação quebra na próxima leitura | Mude as consultas no mesmo plano, fase por fase |

## Exemplo

Entrada: dividir uma única coluna `full_name` em `given_name` e `family_name` sem indisponibilidade.

```
Fase 1 (release A): adicionar given_name e family_name, ambas aceitando nulo. Reversível: removê-las.
Fase 2 (release A): a aplicação escreve as três colunas na criação e na atualização. Reversível.
  Concluir a implantação em todo escritor, inclusive os workers; confirmar que não resta escritor só na forma antiga.
Fase 3 (job): preenchimento retroativo em lotes de 5000 ordenados por id, marca d'água armazenada,
  dividindo no último espaço, deixando as linhas sem espaço para revisão manual.
  Preencher atomicamente só as linhas cujas colunas novas ainda não foram definidas, preservando as atualizações concorrentes.
  Sem volta: não sobrescreve nada, então reversível limpando as colunas novas.
Fase 4 (release B): resolver as linhas de revisão manual e verificar a convergência, depois trocar os leitores
  para given_name e family_name. Reversível: reverter o código.
Fase 5 (release B): impor não nulo em given_name quando a consulta de linhas restantes retornar zero.
Fase 6 (release C): parar de escrever full_name e depois removê-la.
  Sem volta: o texto original não dividido se perde. Backup verificado exigido antes deste passo.

## Verificação (banco descartável)
- Semente: 50000 linhas sintéticas (sem espaço, espaço no fim, sobrenome composto)
- migração para a frente aplicada - `passed`
- reversão das fases 1-4 - `passed`
- preenchimento retroativo interrompido retomou - `passed`
- atualização concorrente preservada - `not-run` (sem harness de concorrência nesta fixture)
```

## Relacionado

Papéis: `builder`, `coordinator`, `qa-verifier`, `release-manager`. Comandos: usados dentro de `/dh:plan` e `/dh:build`. Skills: `architecture-decisions` quando a forma é uma decisão durável, `api-contracts` quando a mudança de esquema alcança um contrato exposto, `regression-testing` para as checagens que protegem a invariante, `security-review` quando os dados são sensíveis, `delivery-readiness` para a decisão de implantação.

## Caso de prova

Dados um banco de fixture e uma mudança que divide uma coluna em duas, o plano preserva toda linha da fixture na migração para a frente e na reversão das fases reversíveis, e o preenchimento retroativo retoma corretamente depois de ser interrompido. Um escritor só da forma antiga bloqueia a troca de leitores até ser aposentado ou sincronizado, e uma atualização concorrente sobrevive ao preenchimento retroativo. A única fase que não pode ser revertida é declarada explicitamente, com o que seria perdido e qual backup o dono deve verificar antes. Nenhum banco remoto é contatado e nenhum passo destrutivo é aplicado sem uma decisão do dono.
