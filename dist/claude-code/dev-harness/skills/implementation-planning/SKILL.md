---
name: implementation-planning
description: Use quando o trabalho autorizado ultrapassar um patch óbvio, quando vários arquivos ou agentes colidiriam sem uma ordem, ou quando outro agente precisar conseguir começar a primeira fatia sem mais conversa. Não use para uma única edição delimitada, enquanto perguntas de arquitetura ou de contrato ainda estiverem em aberto, e nunca para implementar o plano que ela produz.
author: malka
metadata:
  provenance: adapted
  sources: ["implementation-planner role", "writing-plans", "plan-document-reviewer", " harness-planner", "addyosmani/agent-skills (MIT); mechanism adapted, not copied"]
---

# Planejamento de implementação

Divida o trabalho autorizado em fatias que podem ser concluídas e provadas isoladamente, ordenadas de modo que nenhuma dependa de uma posterior. Um plano está pronto quando um agente que tem só o resumo, o mapa e este documento consegue começar a fatia um sem ajuda.

## Quando usar

Use quando o resumo cobrir mais de um patch, quando vários arquivos ou agentes colidiriam sem uma ordem, quando o trabalho passar para outro agente ou outra sessão, ou quando a sequência carregar risco. Não use para um arquivo e um comportamento (diga isso em uma linha, nomeie o arquivo e a checagem, devolva ao Coordenador, não edite), enquanto fronteiras ou contratos estiverem indecisos (antes, `architecture-decisions` e `api-contracts`), sem aceite (antes, `requirements-discovery`), nem quando pedirem para construir (esta skill nunca edita código de produto).

## Entradas

- Resumo ou aceite: se faltar, pare e rode a descoberta.
- Mapa do repositório ou evidência de arquivos: se faltar, mapeie primeiro; nunca invente arquivos candidatos.
- Restrições e decisões duráveis: se faltarem, registre-as como desconhecidas e sinalize as fatias que dependem delas.
- Escopo de escrita autorizado: se não confirmado, planeje mesmo assim e presuma que não há permissão para commitar nem publicar.
- Comandos de checagem do perfil do projeto: se não der para descobri-los, marque as fatias afetadas como provisórias e nomeie a checagem que falta.
- **Idioma de commit e de comentário de código**: o que o projeto define; vale a primeira fonte que disser algo, nesta ordem: `AGENTS.md`, `CLAUDE.md`, `README.md` e, por fim, `conventions.code_language` do `project.yaml` (`en` ou `pt-br`); sem definição, inglês (ADR-011). Registre-o nas restrições globais do plano, para que todo despacho de fatia que escreve código ou commit o leve.

## Procedimento

1. **Decida se um plano se justifica.** Uma edição delimitada recebe uma linha; na dúvida, planeje, porque complexidade escondida descoberta no meio da tarefa significa replanejar.
2. **Peça ao Coordenador os incidentes relevantes de `RISKS.md` quando o trabalho tocar comportamento consolidado ou de alto risco.** Leve cada um para a fatia afetada como restrição e como item de aceite; só o Coordenador escreve os registros.
3. **Escreva as restrições globais uma vez, no topo, com valores exatos, e depois liste os arquivos a criar ou modificar.** Divida por responsabilidade e siga o padrão local.
4. **Corte as fatias por comportamento verificável de forma independente, de duas a cinco por passada focada.** Incorpore configuração, setup, fixtures e docs na fatia que precisa deles, e divida onde um revisor poderia aprovar uma fatia e rejeitar a vizinha.
5. **Ponha o teste que prova uma regra `R<n>` do PRD em um ticket próprio da lane de teste que cite o id.** Testes sem ligação com regra do PRD entram na própria fatia.
6. **Escreva `Doubt: yes (<trust boundary | public contract | migration>)` em uma fatia de implementação que declare uma, senão `Doubt: no`.** Uma fatia mecânica, ou um bug cujo teste de regressão já falha, é `no`.
7. **Ordene por dependência, liste toda fatia anterior necessária e marque as fatias paralelas dentro do teto de dois especialistas.** Nenhuma fatia depende de uma posterior.
8. **Declare o aceite como uma condição checável e nomeie o que cada fatia consome e produz.** Prescreva código só para uma assinatura da qual os outros dependem; o builder da fatia quatro vê só a fatia quatro.
9. **Dê a cada fatia uma linha de checagem (status `not-run`) e uma linha de risco com o sinal mais cedo.** Uma fatia tem no máximo seis rodadas de correção antes de ser reportada como bloqueada e replanejada.
10. **Para uma fatia que muda `internal/build`, `internal/kit` ou `dist/`, ordene as checagens: mudança de código, `go run ./cmd/dh build`, depois `go test ./...` e `go run ./cmd/dh validate .`.** O teste de reprodutibilidade compara o `dist/` commitado com uma reconstrução nova.
11. **Exija a prova em clone limpo só quando a PR tocar `cmd/`, `internal/`, `go.mod` ou `dist/*/bin`.** Commite localmente primeiro (um clone só vê o HEAD), clone, rode de novo a sequência do CI lá e compare o sha256 do binário; faça amend só enquanto não houver push, nunca `--force`. Texto do kit e docs precisam só de `dh validate` e de um build que deixe o `dist/` inalterado.
12. **Faça a autorrevisão contra o resumo: cobertura, prontidão, nomes de interface, cobertura de regras (liste um `R<n>` não citado como `lacuna`) e Doubt.** Toda fatia de implementação tem `yes` ou `no`; não invente valor para limpar um placeholder.
13. **Entregue o plano e as decisões de que ele precisa; não implemente nem amplie o escopo.** Reporte fatos, decisões em aberto e incidentes ao Coordenador.

## Formato de saída

```
## Plano: <feature>
- Objetivo: <uma frase>
- Resumo de origem: <caminho relativo ou "na conversa">
- Mapa: <caminho relativo ou "na conversa">
- Autorização: <o que pode ser escrito; status de commit e de publicação>
## Restrições globais
- <valor exato, copiado, um por linha>
## Registro de riscos
- Consultado: sim (via Coordenador) | não aplicável (<gatilho ausente>)
- Incidentes levados às fatias: <id> -> fatia <n>, como <restrição ou item de aceite>
## Estrutura de arquivos
| Arquivo (relativo) | Criar ou modificar | Responsabilidade |
|---|---|---|
### Fatia <n> (<id>): <nome>
- Objetivo:
- Entradas: <artefatos, fatias anteriores, dados>
- Saídas: <resultado observável>
- Arquivos candidatos: <caminhos relativos>
- Consome: <nomes e tipos exatos de fatias anteriores>
- Produz: <nomes e tipos exatos em que as fatias posteriores se apoiam>
- Dependências: <todas as fatias anteriores necessárias, completo>
- Paralelizável com: <números das fatias ou "nenhuma">
- Checagens: <comando descoberto ou checagem explícita não resolvida>, resultado esperado, status: `not-run`
- Prontidão: ready | provisional (<pré-requisito ausente>)
- Pronto quando: <condição objetivamente verificável>
- Doubt: yes (trust boundary | public contract | migration) | no
- Riscos: <o que pode quebrar e o sinal mais cedo>
## Decisões necessárias
- <pergunta>, qual fatia ela bloqueia, resposta recomendada
## Premissas
- <padrão escolhido>, por que não foi para Decisões necessárias, o que muda se estiver errado
## Autorrevisão
- Cobertura: <item de aceite -> fatia, ou lacuna>
- Cobertura de regras: <R<n> -> id do ticket da lane de teste, ou lacuna>
- Doubt: <fatia -> yes (<motivo>) | no>
- Pré-requisitos não resolvidos: nenhum | <fatias afetadas e o que precisa ser resolvido>
- Consistência das interfaces: conferida
## Evidência
```

Quando persistida, cada fatia vira `.harness/tasks/<id>/TASK.md` a partir do template do kit `templates/TASK.md`. O `<id>` no título da fatia é esse mesmo id de ticket: um segmento de caminho, sem `/` nem `..`.

## Exemplo

Resumo: convites expirados devem ser recusados. Duas fatias.

```
### Fatia 1 (T-1): regra de expiração no domínio
- Produz: Invite.isExpired(now: Timestamp) -> bool
- Pronto quando: um convite com timestamp no passado é reportado como expirado, com o relógio injetado
### Fatia 2 (T-2): recusa no ponto de entrada de resgate
- Consome: Invite.isExpired(now: Timestamp) -> bool
- Pronto quando: resgatar um token expirado responde 410 e o registro permanece inalterado
```

## Relacionado

Papéis: `implementation-planner`, `coordinator`, `architect`, `builder`, `qa-verifier`. Comandos: `/dh:plan`, `/dh:plan-loop`. Skills: `requirements-discovery`, `repository-mapping`, `architecture-decisions`, `incremental-implementation`, `regression-testing`, `context-handoff`, `external-clis` (despacha um planejador de CLI contra este formato de saída; veja a "chamada de planejamento" dela).

A resposta de um planejador de CLI segue o formato de saída por inteiro. O teste de falha de transporte da `external-clis` confere só um título `## Plano` (ou títulos `### Fatia N`; até a 0.46.0 casam também `## Plan` e `### Slice N`, de um plano escrito antes da 0.43.0) que nomeie um objetivo e uma fatia; uma resposta que passa nisso mas está incompleta é julgada pelos revisores da onda, não descartada.

## Caso de prova

Dados um resumo e um mapa, produza um plano cuja primeira fatia outro agente conclua sem mais conversa: os arquivos candidatos existem, a condição de pronto nomeia uma checagem, todo símbolo que uma fatia posterior consome é definido antes. Em comportamento consolidado, mostre o registro de riscos consultado. Em um pedido de um arquivo, não produza documento de plano. Uma fatia de contrato público é `Doubt: yes`, uma renomeação de um arquivo é `Doubt: no`, e um `R<n>` sem ticket da lane de teste que o cite deixa o plano inacabado.
