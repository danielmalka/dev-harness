# Novo fluxo de PRD orientado a subagents

English version: [docs/en/new-prd-flow.md](en/new-prd-flow.md)

Proposta consolidada em 2026-09-28. Ainda não substitui `templates/` nem o fluxo atual do kit.

Uma feature pequena é um PRD. O subagent que implementa não lê o PRD inteiro. Lê o ticket, que já carrega a story, o fluxo e a cópia das regras que o afetam.

Três barreiras fecham o processo: PRD não aprovado não vira ticket, ticket não revisado não fecha, feature não entregue sem revisão final.

## 1. Entrada

O solicitante pode trazer qualquer ponto de partida:

- uma ideia
- um pedido
- uma solução imaginada
- um bug
- uma melhoria
- uma reclamação de usuário
- uma necessidade técnica percebida

A entrada não precisa estar completa. É sinal para investigação, não um PRD pronto.

## 2. Criação do PRD

O orquestrador chama o **Subagent de Criação de PRD**.

Esse subagent conduz um grill com o solicitante antes de escrever. O grill continua até a capability estar compreensível e as decisões importantes estarem explícitas. Não é entrevista burocrática.

O grill precisa deixar claro:

- qual problema realmente existe
- quem é afetado
- qual comportamento precisa mudar
- como o usuário percebe a mudança
- quais cenários precisam funcionar
- quais comportamentos não podem acontecer
- quais limites não podem ser ultrapassados
- quais fluxos ou documentações existentes são afetados

Emenda de regra também é deste subagent. O solicitante pode pedir alteração.

### Estrutura do PRD

Só estas seções, nesta ordem. Faltou uma, não quebra em ticket.

#### Problema

Quem sofre, o que não funciona, em qual situação, qual o impacto, o que precisa mudar.

Uma capability. Duas capabilities independentes viram dois PRDs antes do primeiro ticket.

#### Solução

O que passa a ser verdade depois da entrega.

Prioriza comportamento, resultado para o usuário, fluxo esperado e resposta do sistema. Não é descrição de arquitetura.

#### Regras

Constraints inegociáveis, numeradas (`R1`, `R2`, `R3`).

Cada regra é verificável, observável, uma condição obrigatória ou uma proibição explícita. Texto genérico, preferência ou intenção não é regra.

```text
R1. Um usuário sem permissão não pode visualizar o relatório.
R2. O sistema deve registrar toda tentativa negada.
R3. A API não pode retornar dados de outro usuário.
```

#### Docs

Paths dos fluxos já documentados que esta feature toca.

Se não toca nenhum, a linha exata:

```text
Nenhum fluxo documentado afetado.
```

Seção vazia não vale. Se a feature altera um fluxo já documentado, a atualização da documentação faz parte da entrega.

## 3. Revisão adversarial do PRD

O orquestrador chama o **Revisor Adversarial de PRD**.

O revisor não melhora o texto por preferência. Procura falha que cause implementação errada.

Verifica:

- O problema está claro?
- Existe apenas uma capability?
- A solução descreve comportamento observável?
- As regras são realmente inegociáveis?
- Cada regra pode ser verificada?
- Há regra faltando?
- Há contradição entre solução e regras?
- O fluxo deixa cenário importante sem resposta?
- O PRD contém decisão implícita?
- A seção Docs está correta?
- Dá para quebrar em tickets sem inventar requisito?

Resultado:

```text
APROVADO
```

ou:

```text
RECUSADO

Motivos:
- ...

Correções obrigatórias:
- ...
```

PRD recusado volta ao Subagent de Criação de PRD. Só avança depois da aprovação.

## 4. Decomposição em tickets

PRD aprovado: o orquestrador chama o **Subagent de Decomposição de PRD**.

Papel separado do autor do PRD, para não misturar escrita de regra com quebra de trabalho.

Esse subagent:

- divide o trabalho por lane
- cria tickets de backend, frontend e teste quando necessário
- copia literalmente as regras relevantes
- indica os testes necessários
- identifica a documentação a atualizar
- preserva o fluxo esperado
- não adiciona escopo que não está no PRD

Story que precisa de API e tela vira dois tickets quando as responsabilidades são independentes. "E testa também" no final do backend não é ticket de teste.

Lanes:

```text
backend
frontend
teste (unitário)
teste (integração)
teste (unitário e integração)
```

Uma lane por ticket.

## 5. Estrutura do ticket

### Lane

```text
Lane: backend
```

ou:

```text
Lane: teste (integração)
```

### Story

```text
Como <ator>,
quero <ação>,
para <resultado>.
```

### Fluxo

O que o ator faz e o que acontece em seguida, em passos observáveis. Ator é usuário, sistema ou teste. Este bloco é o aceite.

```text
1. O usuário envia o formulário válido.
2. O sistema cria o registro.
3. O sistema retorna confirmação.
4. O registro aparece na listagem.
```

### Regras

Cópia literal só das regras que este ticket tem que honrar. Sem link. Sem resumo. Regra que não cabe aqui não entra.

```text
R1. Um usuário sem permissão não pode visualizar o relatório.
R3. A API não pode retornar dados de outro usuário.
```

### Docs

Paths a atualizar neste ticket, ou `nenhum`. Herda a seção Docs do PRD e é refinada por ticket.

```text
Docs:
- docs/fluxos/relatorio.md
```

ou:

```text
Docs:
- nenhum
```

### Pronto quando

O fluxo acontece e cada regra copiada tem evidência.

```text
Pronto quando:
- o fluxo descrito acontece
- as regras R1 e R3 foram verificadas
- os testes relevantes estão verdes
- a documentação indicada foi atualizada
```

### Exemplo mínimo

```text
Lane: teste (integração)
Story: Como visitante, quero ser barrado se enviar senha vazia, para não criar sessão.
Fluxo:
1. Visitante envia login com senha vazia.
2. Não há sessão.
3. Vê "senha obrigatória".
Regras:
R1. Senha vazia não autentica.
R2. A mensagem é "senha obrigatória" e não revela se o email existe.
Docs: nenhum
Pronto quando: o teste cita R1 e R2 e fica verde.
```

## 6. Execução dos tickets

O orquestrador chama um **Subagent de Desenvolvimento especializado** conforme a lane. Confere a evidência. Não implementa o ticket.

O subagent recebe o texto completo do ticket, o contexto mínimo do repositório, as regras copiadas, os arquivos ou áreas esperadas e os critérios de evidência. Não precisa interpretar o PRD inteiro.

Ele:

1. executa apenas o ticket recebido
2. respeita todas as regras copiadas
3. não inventa escopo
4. atualiza a documentação indicada
5. roda as verificações necessárias
6. devolve evidência real

## 7. Loop agêntico

```text
Subagent de Desenvolvimento
        |
        v
Subagent de Teste e Revisão
        |
        v
Aprovado ou Reprovado
```

### Desenvolvimento

Executa o ticket e devolve:

- arquivos alterados
- comportamento implementado
- testes executados
- comandos utilizados
- resultado dos comandos
- limitações ou dúvidas
- evidência visual, quando aplicável

Self-report não é evidência. O orquestrador confere o resultado.

### Teste e revisão

O revisor confere o ticket contra a Story, o Fluxo, as Regras copiadas, o Pronto quando, os padrões do projeto, os testes necessários e a documentação indicada.

Aprovado:

```text
APROVADO

Evidências:
- ...

Testes:
- comando: ...
- resultado: ...

Docs:
- ...
```

Reprovado:

```text
REPROVADO

Problemas:
- ...

Regra ou critério afetado:
- R2
- Pronto quando, item 3

Correção exigida:
- ...
```

Reprovado volta ao Subagent de Desenvolvimento com a correção específica. Depois da correção, o mesmo ticket passa de novo pelo Subagent de Teste e Revisão.

O ticket não fecha com teste ausente, teste vermelho, regra não verificada, documentação pendente, comportamento diferente do fluxo, ou evidência só declarada.

## 8. Revisão final da feature

Todos os tickets aprovados: o orquestrador informa o fim do desenvolvimento e convoca subagents de revisão da feature.

Essa etapa olha a feature inteira, não cada ticket isolado.

Pode chamar revisão de:

- integração entre backend e frontend
- comportamento funcional
- testes
- documentação
- segurança
- regressões
- aderência ao PRD

Se necessário, dispara subagents de teste e de documentação.

Se houver tela acessível, pode convocar o Subagent de Criação de PRD, ou outro subagent com acesso à tela, para usar o app como cliente. Esse passo só existe quando a superfície visual está disponível. Não substitui teste de backend nem de integração.

O uso como cliente:

1. acessa a tela
2. age como usuário
3. segue o fluxo descrito na Solução
4. verifica o resultado
5. registra evidência
6. testa também os comportamentos proibidos ou inválidos relevantes

Nessa fase o Criador de PRD avalia se o resultado corresponde à intenção original. Não escreve requisito novo.

Resultado:

```text
FEATURE APROVADA
```

ou:

```text
FEATURE REPROVADA

Motivos:
- ...

Tickets ou áreas que precisam retornar:
- ...
```

Reprovação reabre só os tickets ou áreas afetadas e reinicia o loop necessário.

## 9. Quem faz o quê

- **Solicitante:** passa ideia, pedido, solução, bug ou equivalente. Não precisa estar completo. Pode pedir alteração de regra.
- **Orquestrador:** chama criação, revisão adversarial, decomposição, desenvolvimento especializado e revisão final. Confere evidência. Não implementa ticket.
- **Subagent de Criação de PRD:** faz o grill, escreve Problema, Solução e Regras, emenda regra. No fim, se houver tela, pode testar a execução como cliente.
- **Revisor Adversarial de PRD:** avalia estrutura, regras e conteúdo. Recusa o que não está aderente.
- **Subagent de Decomposição de PRD:** quebra os tickets e copia as regras.
- **Subagent de Desenvolvimento especializado:** executa um ticket da sua lane e devolve evidência.
- **Subagent de Teste e Revisão:** aprova ou reprova o ticket. Reprovação volta ao desenvolvimento.
- **Subagents de revisão da feature:** no fechamento, revisam a feature. Teste e documentação entram quando necessário.

## 10. Guardrails

1. Entrada incompleta pode iniciar o processo. PRD incompleto não inicia a execução.
2. Cada PRD é uma capability. Duas capabilities viram dois PRDs antes do primeiro ticket.
3. O PRD nasce depois do grill com o solicitante.
4. Emenda de regra é do Subagent de Criação de PRD. O solicitante pode pedir alteração.
5. PRD só quebra em ticket depois da aprovação adversarial.
6. Ticket só contém requisito presente no PRD.
7. Regra do ticket é cópia literal da regra do PRD.
8. Toda regra observável tem um ticket de teste que a cita pelo id. Sem isso o PRD não entra em execução.
9. Backend e frontend não fecham com o teste da regra ausente ou vermelho.
10. Ticket que muda fluxo documentado atualiza esse doc na mesma entrega. Doc velho deixa o ticket aberto.
11. Implementação que acha fluxo documentado fora da seção Docs para, corrige o PRD, e só então segue.
12. Regra nova descoberta no meio para o trabalho. O PRD é atualizado antes de continuar.
13. Regra alterada no meio: ticket não iniciado é reescrito. Ticket em execução para.
14. Dois tickets no mesmo arquivo não rodam juntos.
15. Reprovação volta ao desenvolvimento com correção específica. O mesmo ticket é revisado de novo.
16. Self-report não prova. Evidência é teste verde, exit code, URL ou passo do fluxo reproduzido.
17. A revisão final confere a feature contra o PRD, não só contra os tickets.
18. Se houver tela, o fluxo principal é testado como cliente.
19. Desejo novo é emenda no PRD ou PRD novo. Não entra como escopo informal do ticket.
20. "Feito" não entrega. Entrega exige evidência conferida e Resumo executado.

## 11. Fechamento

O PRD só encerra com todos os tickets aprovados, testes relevantes executados, documentação atualizada, revisão final aprovada e evidência conferida pelo orquestrador.

O orquestrador escreve este bloco no próprio PRD.

### Resumo executado

**Entregue.** Uma frase do que ficou verdadeiro, contra a Solução.

**Regras.**

```text
R1. Honrada.
Evidência: ...

R2. Honrada.
Evidência: ...

R3. Não honrada.
Motivo: ...
Decisão pendente: ...
```

**Tickets.**

```text
- TKT-001, backend, aprovado
  Evidência: ...

- TKT-002, frontend, aprovado
  Evidência: ...

- TKT-003, teste de integração, aprovado
  Evidência: ...
```

**Docs.**

```text
- docs/fluxos/relatorio.md, atualizado
- docs/api/relatorio.md, atualizado
```

ou:

```text
Nenhum fluxo documentado foi afetado.
```

**Fora.** O que a Solução prometeu e não entrou.

Desvio de regra sem decisão explícita do solicitante: o PRD não é entregue.

## Fluxo

```text
Entrada do solicitante
        |
        v
Grill com solicitante
        |
        v
Subagent de Criação de PRD
        |
        v
Revisor Adversarial de PRD
        |
        v
PRD aprovado?
   |-- não --> corrigir PRD
   |-- sim
        |
        v
Subagent de Decomposição de PRD
        |
        v
Tickets
        |
        v
Desenvolvimento especializado
        |
        v
Teste e revisão do ticket
   |-- reprovado --> desenvolvimento
   |-- aprovado
        |
        v
Todos os tickets aprovados?
   |-- não --> próximo ticket
   |-- sim
        |
        v
Revisão final da feature
        |
        v
Testes, documentação e cliente na tela, se houver
        |
        v
Resumo executado no PRD
        |
        v
PRD entregue
```
