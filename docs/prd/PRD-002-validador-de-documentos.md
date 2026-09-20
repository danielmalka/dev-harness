# PRD-002 · Revisor adversarial somente leitura para documentos importantes (validador de PRD primeiro)

| Campo | Valor |
|---|---|
| Status | implementado no kit em 20/09/2026; avaliação comportamental pendente |
| Dono | Daniel Lemos |
| Criado / atualizado | 2026-09-20 / 2026-09-20 |
| Stories | a derivar |

## 1. Problema

O `discover` entrega um PRD direto ao dono. Lendo um PRD gerado no danlemos, o dono viu lacunas, ambiguidades e critérios fracos que um revisor pegaria antes. O plano já exige revisão independente para código (`code-review`, `security-review`), mas não para documentos. Um PRD ruim vira stories ruins e tasks ruins; o custo de corrigir cresce a cada etapa.

## 2. Resultado esperado

- Sinal 1: PRDs que chegam ao dono já passaram por um ciclo criação → validação até aprovação · hoje: zero · alvo: 100% dos PRDs gerados por `discover`.
- Sinal 2: número de correções que o dono pede num PRD depois de aprovado pelo validador · hoje: sem baseline · alvo: medir nos pilotos.

## 3. Usuários e cenários

| Usuário | Cenário | Hoje | Depois |
|---|---|---|---|
| Dono | Recebe um PRD do `discover` | Lê e encontra problemas de estrutura sozinho | Recebe um PRD que passou pelo validador, com o relatório de problemas resolvidos anexado |
| Coordenador | Termina `discover` | Entrega ao dono | Despacha o validador; se reprovar, devolve ao `product-discovery` com a lista; repete até aprovar ou esgotar o limite de rodadas; só então entrega |

## 4. Escopo

**Entra**
- RF-01 Novo papel `document-validator`: somente leitura, adversarial, responde ao Coordenador com veredito `approved` ou `changes required` e a lista numerada de problemas (P1, P2...) com categoria e sugestão concreta.
- RF-02 Categorias de problema, herdadas do prompt de referência abaixo: lacunas (o que o discovery discutiu e não virou requisito), ambiguidades, conflitos, excessos (scope creep), critérios fracos (aceite subjetivo ou não verificável), organização (duplicidade, IDs inconsistentes, domínio errado).
- RF-03 Ciclo no Coordenador: `discover` produz PRD → validador → se `changes required`, `product-discovery` corrige só os pontos listados → validador de novo. Limite de duas rodadas de correção, como o resto do kit; depois, entrega parcial ao dono com os pontos abertos.
- RF-04 Relatório persistente ao lado do documento (`<documento>.review.md`, por exemplo `.harness/prd/PRD-001.review.md`) com IDs estáveis e estado `pending` / `applied` / `rejected`, para que o ciclo não perca o que já foi discutido.
- RF-05 Skill `document-review` com a rubrica, o formato do relatório e o controle de falso positivo (um PRD sem problemas deve sair aprovado sem achados inventados).
- RF-06 Premissa geral, registrada no plano: todo documento que vira contrato (PRD, story, plano, ADR) passa por um revisor adversarial somente leitura antes de chegar ao dono. PRD é o primeiro; story e plano vêm depois.

**Não entra**
- Validador que edita o documento. No kit, quem edita é quem escreveu; o revisor só reporta. O prompt de referência edita com aprovação do usuário; aqui o ciclo é entre agentes e o dono só vê o resultado.
- Discussão interativa com o dono a cada ponto. O dono recebe o PRD pronto e o relatório; se quiser discutir, abre uma rodada.

## 5. Critérios de aceite

- AC-01 (RF-01) Quando o validador receber um PRD com uma lacuna plantada em relação ao discovery, então reporta a lacuna com categoria e sugestão, e o veredito é `changes required`.
- AC-02 (RF-05) Quando o validador receber um PRD sem problemas, então o veredito é `approved` sem achados.
- AC-03 (RF-03) Quando o validador reprovar, então o Coordenador despacha a correção só dos pontos listados e o PRD volta ao validador; o dono só recebe depois de `approved` ou do limite de rodadas.
- AC-04 (RF-04) Quando o ciclo for retomado noutra sessão, então o relatório persistente mostra o que já foi aplicado e rejeitado.

## 6. Restrições e impacto técnico

- Papel novo em `.agents/`, skill nova em `.skills/`, comando `discover` alterado para incluir o ciclo. Modelo sugerido: Opus, como os outros revisores de profundidade.
- Idioma: agente e skill em inglês; o relatório segue o `language` do projeto.
- Eval: caso positivo (lacuna plantada), caso controle (PRD limpo), caso de limite de rodadas.

## 7. Riscos e decisões pendentes

- Risco: validador gerar achados cosméticos e travar o ciclo · mitigação: severidade, só `changes required` quando há lacuna, conflito ou critério não verificável; cosmético não reprova.
- Decisão fechada em 20/09/2026: o papel se chama `document-validator`, genérico desde a v1; PRD é o primeiro alvo.
- Decisão fechada em 20/09/2026: o relatório persistente fica ao lado do documento, como `<documento>.review.md` (para `.harness/prd/PRD-001.md`, é `.harness/prd/PRD-001.review.md`).

## 8. Referência: prompt fornecido pelo dono em 20/09/2026

Base de inspiração. Difere do kit em dois pontos: edita o documento com aprovação do usuário e discute ponto a ponto no chat. No kit o revisor é somente leitura e o ciclo é entre agentes.

```text
Voce e o PRD Validator, um analista de produto critico especializado em revisar e refinar user stories e requisitos de software.

## Seu papel

Voce recebe user stories e requisitos gerados automaticamente a partir de um discovery-notes.md. Seu trabalho e revisar com olho critico, identificar problemas, discutir com o usuario e EDITAR o arquivo stories-requisitos.md diretamente para corrigir os problemas encontrados.

## Processo obrigatorio

### Passo 1: Leitura completa
- Leia discovery-notes.md e os user stories/requisitos gerados
- Compare tudo: o que foi discutido no discovery esta refletido nos requisitos?
- NUNCA comece a reportar antes de ler tudo

### Passo 2: Analise e salvar no arquivo persistente
Salve sua analise no arquivo de relatorio persistente (caminho fornecido no prompt) usando IDs sequenciais (P1, P2, P3...) e marcadores de status:

Formato: `- **P1** [PENDENTE] [CATEGORIA] descricao do problema`

Categorias de problemas:

1. **Lacunas** - Funcionalidades mencionadas no discovery que nao viraram requisito
2. **Ambiguidades** - Stories ou requisitos que um dev interpretaria de formas diferentes
3. **Conflitos** - Requisitos que se contradizem entre si
4. **Excessos** - Requisitos que nao foram mencionados no discovery (scope creep)
5. **Criterios fracos** - Criterios de aceite subjetivos ou nao verificaveis
6. **Organizacao** - Stories no dominio errado, requisitos duplicados, IDs inconsistentes

### Passo 3: Apresentar e discutir
- Apresente o relatorio no chat com a mesma numeracao do arquivo
- Para cada problema, apresente sua sugestao concreta de correcao
- Discuta cada ponto com o usuario
- Aceite feedback, ajuste se necessario
- NUNCA tome decisoes sozinho

### Passo 4: Edicao incremental e direta no arquivo
- A cada ponto aprovado pelo usuario, edite o arquivo stories-requisitos.md IMEDIATAMENTE usando Write ou Edit
- Nao pergunte se pode editar: quando o usuario aprovar ou concordar, edite sem hesitar
- Atualize o status no arquivo persistente: [PENDENTE] -> [APLICADO] ou [REJEITADO]
- Confirme no chat o que foi editado e mostre o trecho alterado

### Regra critica de continuidade
- A CADA TURNO, leia o arquivo de relatorio persistente ANTES de responder
- Este arquivo e sua UNICA fonte de verdade sobre o que ja foi analisado

## Regras absolutas

- NUNCA mencione Enricher, Planner, Coder ou qualquer outro agente
- NUNCA invente requisitos que nao estavam no discovery
- NUNCA edite sem aprovacao explicita do usuario
- NUNCA reescreva tudo de uma vez - edicoes cirurgicas e precisas

## Finalizando a validacao

Quando todos os problemas identificados tiverem sido resolvidos e o usuario confirmar que esta satisfeito com o resultado, inclua o marcador [PHASE_COMPLETE] ao final da sua mensagem de encerramento.
Sempre instrua o usuario: "Se nao ha mais nada para ajustar, clique no botao Aprovar para avancar para a proxima fase."

## Idioma

Toda comunicacao deve ser em portugues brasileiro.
```
