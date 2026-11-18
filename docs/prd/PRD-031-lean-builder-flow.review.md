# Revisão: PRD-031

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

## Rodada 1

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-031-lean-builder-flow.md |
| Fonte | tasks/PRD-031/SOURCE.md (harness dir) |
| Rodada | 1 de 2 |
| Veredito | `changes required` |
| Atualizado | 2026-10-30 |

Revisor: cli:agy/gemini-3.8-flash-high (etapa de documento). Achados bloqueantes: 5. Resposta bruta: `tasks/PRD-031/evidence/document-r1/agy-output.md` (harness dir).

## Compreensão
- O que o documento propõe: O documento estabelece um fluxo de builder mais enxuto, instituindo o desenvolvimento test-first como padrão para as fatias de implementação de código e introduzindo três heurísticas no teste de regressão para evitar testes vazios ou desonestos. Ele padroniza o corpo das mensagens de commit para documentar o motivo, as mudanças de código e as checagens de verificação, e limita os passos de procedimento do harness a imperativos de uma frase, com uma justificativa opcional de uma frase. Impõe um teto de tamanho que obriga quatro arquivos da fase do builder a encolher para no máximo 40% do tamanho atual, sem perder as salvaguardas essenciais registradas num inventário de regras. Por fim, descreve sete simplificações distintas da etapa de revisão, entre elas a conclusão com uma única rodada limpa e os estados de achado consolidados, submetidas à aprovação individual do dono.
- Como o sistema atual funciona, segundo o documento: A fase do builder é regida por quatro arquivos que somam 830 linhas e mais de 107 quilobytes de instruções de procedimento. A execução test-first hoje vale estritamente para o teste de regressão de bug, e as fatias de feature em geral ficam sem mandato de red-first. Os ciclos de revisão carregam um custo grande de tokens e de controle por causa da maquinaria de várias etapas, incluindo seis estados de achado, identificadores estruturados de achado, cursores de rotação de fornecedor e confirmações obrigatórias de duas rodadas limpas consecutivas. A documentação de autoria do harness não limita o tamanho dos passos de procedimento, e os corpos das mensagens de commit não têm padrões estruturais impostos. As dúvidas no meio do build permitem até três ciclos de contexto novo e dependem de classificações complexas de achados em várias vias.
- Seções não parafraseadas: nenhuma

## Achados

- **P1** [pending] [gap] (blocking) O padrão test-first para fatias de implementação foi deixado de fora do regression-testing
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R1 e R3
  - Evidência: A recomendação aprovada 1 do SOURCE.md exige levar o test-first para incremental-implementation e regression-testing; a R3 mantém explicitamente o red-first restrito aos testes de bug.
  - Sugestão: Atualizar a R3 para levar os requisitos e as exceções de test-first para fatias ao regression-testing, junto com as três heurísticas.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [pending] [conflict] (blocking) Cinco estados de achado na fonte afirmados como seis estados de achado no documento
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 1, linha 31 e seção 3, R13
  - Evidência: A seção `Source findings` do SOURCE.md, item 4, diz cinco estados de achado; a seção 1 do PRD e a R13 afirmam seis estados de achado e seis saídas.
  - Sugestão: Alinhar a seção 1 e a R13 ao SOURCE.md, dizendo que os estados de achado encolhem de cinco para três e corrigindo a contagem de saídas para cinco.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [pending] [conflict] (blocking) A forma restante da dúvida no meio do build é atribuída por engano à R15
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R2
  - Evidência: A R2 afirma que a R15 define a forma restante da dúvida no meio do build, mas a R15 define a substituição da oscilação na reentrada da revisão.
  - Sugestão: Remover a referência à R15 da R2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [weak criterion] (blocking) Checagem de grep incompleta para a lista de exceções da R1 e para o requisito de relatório
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R1
  - Evidência: `grep -n "fail"` não consegue observar se a lista de exceções ou a regra de relatório da fatia foi implementada.
  - Sugestão: Usar uma checagem de padrão explícita para as exceções e para a estrutura do relatório.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [pending] [ambiguity] (blocking) Condição "substantivo" de achado indefinida na classificação de duas classes
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R14
  - Evidência: A classificação binária ("acionável ou não") conflita com a escalação indefinida para achados "substantivos" sem nenhum acionável.
  - Sugestão: Definir critérios concretos para o que torna substantivo um achado não acionável.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P6** [pending] [excess] (non-blocking) Redução não pedida do limite de rodadas de 6 para 3
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R16
  - Evidência: O SOURCE.md não pede para mudar o limite de 6 para 3 na build/QA e na revisão.
  - Sugestão: Descartar a R16 ou documentá-la em Decisões como uma proposta não pedida que exige aprovação do dono.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P7** [pending] [organization] (non-blocking) Omissão da etapa do pacote de revisão na hipótese da seção 1
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 1, Hipótese
  - Evidência: Os achados da fonte listam as regras do pacote de revisão como contribuinte de peso, e o PRD as substitui pela dúvida no meio do build, embora remova os pacotes na R17.
  - Sugestão: Incluir as regras do pacote de revisão na lista de etapas da seção 1.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- As medições de linhas e bytes dos arquivos casam com as medições da árvore de trabalho do SOURCE.md.
- A fórmula do orçamento de tamanho (<= 330 linhas, <= 43000 bytes) espelha corretamente a meta de 40%.
- As três heurísticas de teste casam literalmente com a recomendação 2 do SOURCE.md.
- Os requisitos de conteúdo do commit (por quê, o quê, como) casam com a recomendação 3; a D1 justifica bem levá-los para `delivery-readiness`.
- As restrições de tamanho dos passos de procedimento e de referências cruzadas na R5 casam com a recomendação 4.
- As sete salvaguardas essenciais da R8 preservam corretamente os requisitos de despacho.
- As contagens de agentes (13) e de comandos (19+1) são preservadas, sem novas chaves em `project.yaml`.
- O Apêndice exclui corretamente todos os oito itens rejeitados da fonte e a direção de planejamento futuro.
- Custo, risco e recomendações são dados para a R11, a R12, a R15 e a R17.

## Rodada 2

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-031-lean-builder-flow.md |
| Fonte | tasks/PRD-031/SOURCE.md (harness dir) |
| Rodada | 2 de 2 |
| Veredito | `approved` |
| Atualizado | 2026-10-30 |

Revisor: cli:agy/gemini-3.8-flash-high (etapa de documento). Achados bloqueantes: 0. Resposta bruta: `tasks/PRD-031/evidence/document-r2/agy-output.md` (harness dir). O documento foi corrigido entre as rodadas; a revisão da rodada 2 é sobre o documento corrigido.

## Compreensão
- O que o documento propõe: O documento estabelece um fluxo de desenvolvimento test-first obrigatório para as fatias de implementação de código, tanto em `incremental-implementation` quanto em `regression-testing`, isentando apenas as fatias de documentação, configuração e texto. Ele dá ao `regression-testing` três heurísticas para evitar testes vazios ou tautológicos e exige que o corpo das mensagens de commit documente por que, o que e como as mudanças foram verificadas. Limita os novos passos de procedimento do harness a uma única frase imperativa com no máximo uma frase de justificativa, movendo as referências cruzadas para o fim dos passos. Impõe um teto de compressão que restringe quatro arquivos centrais da fase do builder a no máximo 40% do tamanho atual, preservando todas as salvaguardas essenciais por meio de um inventário de regras completo. Por fim, reúne sete simplificações distintas da etapa de revisão, como a conclusão com uma única rodada limpa, os estados de achado reduzidos e os ciclos de dúvida mais estreitos, para aprovação individual do dono.
- Como o sistema atual funciona, segundo o documento: Quatro arquivos de procedimento regem a fase do builder e consomem 830 linhas e 107,478 bytes de texto instrucional. O harness hoje impõe o teste red-first só para a verificação de regressão de bug, e as fatias de feature em geral ficam sem requisito de test-first. Os corpos das mensagens de commit não têm padronização estrutural, e as regras de autoria de skill não impõem limites de brevidade aos passos de procedimento. Os ciclos de revisão acumulam um custo grande de tokens por causa de máquinas de estado de várias etapas, incluindo seis estados de achado, identificadores granulares de achado, cursores de rotação de revisor e verificações obrigatórias de duas rodadas limpas consecutivas. As dúvidas no meio do build permitem até três despachos de contexto novo, baseados em classificações de quatro vias, antes de chegar às revisões comuns de fatia.
- Seções não parafraseadas: nenhuma

## Achados

- **P1** [applied] [gap] (blocking) O padrão test-first para fatias de implementação foi deixado de fora do regression-testing
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R1 e R3
  - Resolução: A R1 agora leva explicitamente o test-first para `incremental-implementation` e `regression-testing` e verifica frases exatas nos dois arquivos, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) Cinco estados de achado na fonte afirmados como seis estados de achado no documento
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 1, linha 27 e seção 3, R13
  - Resolução: A seção 1 cita explicitamente os seis estados do `.agents/coordinator.md`, item 8, alinhada à Correção de 2026-10-30 do SOURCE.md, rodada 2.
  - Nota (Coordenador): a análise da fonte em `tasks/PRD-031/SOURCE.md` dizia cinco estados de achado, e o P2 da rodada 1 do revisor foi levantado contra essa contagem. Uma correção datada no fim desse arquivo (2026-10-30, Coordenador) fixa seis estados, conforme o `.agents/coordinator.md`, item 8, que o PRD usa. O texto do P2 da rodada 1 fica como o revisor o escreveu.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [conflict] (blocking) A forma restante da dúvida no meio do build é atribuída por engano à R15
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R2
  - Resolução: A referência à R15 foi removida da R2, restando só a R14, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [weak criterion] (blocking) Checagem de grep incompleta para a lista de exceções da R1 e para o requisito de relatório
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R1
  - Resolução: Os critérios verificáveis foram atualizados para conferir as frases exatas `expected failure`, `docs-only`, `config-only`, `copy-only`, `test-first` e a linha de relatório `Red:`, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [ambiguity] (blocking) Condição "substantivo" de achado indefinida na classificação de duas classes
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R14
  - Resolução: A redação indefinida "substantivo" foi removida; a classificação é estritamente binária ("acionável ou não"), com escalação explícita se nenhum for acionável, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [excess] (non-blocking) Redução não pedida do limite de rodadas de 6 para 3
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 3, R16
  - Resolução: A R16 foi rotulada explicitamente como uma proposta não pedida do Coordenador, sujeita à aprovação do dono, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P7** [applied] [organization] (non-blocking) Omissão da etapa do pacote de revisão na hipótese da seção 1
  - Local: docs/prd/PRD-031-lean-builder-flow.md seção 1, Hipótese
  - Resolução: A hipótese da seção 1 foi atualizada para incluir as regras do pacote de revisão, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- As medições de linhas e bytes dos arquivos (830 linhas, 107478 bytes) casam com as medições da árvore de trabalho do SOURCE.md.
- O teto do orçamento de tamanho (<= 330 linhas, <= 43000 bytes) espelha corretamente a meta de 40%.
- As três heurísticas de teste casam literalmente com a recomendação 2 do SOURCE.md.
- Os requisitos de conteúdo do commit (por quê, o quê, como) casam com a recomendação 3; a D1 justifica bem levá-los para `delivery-readiness`.
- As restrições de tamanho dos passos de procedimento e de referências cruzadas na R5 casam com a recomendação 4.
- As sete salvaguardas essenciais da R8 preservam corretamente os requisitos de despacho.
- As contagens de agentes (13), de skills (22) e de comandos (19+1) são preservadas, sem novas chaves em `project.yaml`.
- O Apêndice exclui corretamente todos os oito itens rejeitados da fonte e a direção de planejamento futuro.
- Custo, risco e recomendações são dados para a R11, a R12, a R13, a R14, a R15, a R16 e a R17.
