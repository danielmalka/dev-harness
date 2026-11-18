# Revisão: PRD-029

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

## Rodada 1

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-029-independent-review-and-drift-report.md |
| Fonte | Decisões do dono (2026-10-28); descoberta UND-009 (pesquisa em `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Rodada | 1 de 2 |
| Veredito | changes required |
| Atualizado | 2026-10-28 |

## Achados

- **P1** [applied] [conflict] (blocking) A definição de pacote em R1 e R3 contradiz os contratos existentes de montagem do prompt e de despacho.
  - Local: seção 3.1, regras R1 e R3
  - Evidência: R1 e R3 definem o pacote como o texto completo passado a um revisor e o restringem a cinco itens de caminho e metadados, o que removeria as instruções de papel, a skill de revisão, a cláusula anti-delegação e o contrato de veredito exigidos por `.agents/coordinator.md` linhas 66 e 94.
  - Sugestão: Esclarecer que R1 e R3 valem estritamente para a carga do pacote de contexto dentro da montagem do prompt.
  - Resolução: aplicado pelo autor, rodada 1 (Seção 2 Parte A, R1 e R3).
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] (blocking) R26 exclui `fixed` dos achados assentados, o que torna impossíveis as rodadas limpas sob R30 depois de uma correção de bug.
  - Local: seção 3.3 (R26) e seção 3.4 (R30)
  - Evidência: R26 define achados assentados só como `intentional`, `noise` ou `escalated`. Sob R30 ("Nenhum achado bloqueante fica por assentar"), um defeito corrigido continua tecnicamente não assentado, o que bloqueia a convergência para sempre.
  - Sugestão: Definir a condição limpa de R30 como "Nenhum achado bloqueante está `open` ou `reopened`", ou especificar que achados `fixed` confirmados são assentados.
  - Resolução: aplicado pelo autor, rodada 1 (R21, R26, R30).
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] (blocking) A atribuição de ids de achado em R19 não preserva de forma explícita ids estáveis para achados não resolvidos entre rodadas.
  - Local: seção 3.3, regra R19
  - Evidência: R19 diz que os achados recebem `F<round>-<n>` em cada rodada, sem dizer se os achados abertos mantêm os ids anteriores ou são renumerados, o que contradiz o Problema 1.C.
  - Sugestão: Dizer de forma explícita que `F<round>-<n>` vale para achados novos e que os achados levados adiante mantêm os ids originais.
  - Resolução: aplicado pelo autor, rodada 1 (R19).
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [conflict] (blocking) A regra R39 não está marcada como provisória, embora a Decisão D2(b) apresente o papel de comparação como decisão em aberto do dono.
  - Local: preâmbulo da seção 3, seção 3.5 (R39) e Apêndice (D2)
  - Evidência: O Apêndice D2(b) apresenta três opções para quem roda a comparação de deriva, mas o preâmbulo da Seção 3 afirma que R39 vale por inteiro e omite R39 das regras provisórias.
  - Sugestão: Listar R39 como provisória e à espera da Decisão D2(b) no preâmbulo da Seção 3 e nas decisões pendentes do Apêndice.
  - Resolução: aplicado pelo autor, rodada 1 (preâmbulo da Seção 3 e decisões pendentes do Apêndice).
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] (Minor) Violações do limite de tamanho de frase, procedural e descritiva, nas regras de escrita de PRD da ASD-STE100.
  - Local: seção 1 (Problema C), seção 3.1 (R3, R8), seção 3.7 (R55)
  - Evidência: Frases do Problema C (50 palavras), de R3 (34 palavras), de R8 (37 palavras) e de R55 (34 palavras) passam dos limites de 20 palavras (procedural) e 25 palavras (descritiva).
  - Sugestão: Quebrar as frases longas em afirmações de uma instrução só, dentro dos limites de palavras.
  - Resolução: aplicado pelo autor, rodada 1 (Problema 1.C, R3, R8, R55).
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- Fidelidade às decisões do dono Q1 a Q6 e às premissas do Coordenador.
- Leitura do Coordenador sobre Q3 mantida visível como decisão em aberto D1 e não atribuída ao dono.
- Estruturação das decisões em aberto D1, D2 e D3 com opções, consequências e recomendações.
- Marcadores de provisório em R30-R35, R47 e R51.
- Testabilidade das regras de R1 a R60.
- Todas as citações de arquivos do dh e números de linha conferidas contra `source3-cited-dh-passages.md` e achadas 100% corretas.
- Ausência de excesso de escopo (nenhum agente, comando ou dependência de runtime novo).
- Bloco de compreensão corretamente limitado ao `document-validator`.
- Cobertura obrigatória das seis categorias verificada.

## Rodada 2

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-029-independent-review-and-drift-report.md |
| Fonte | Decisões do dono (2026-10-28); descoberta UND-009 (pesquisa em `~/.harness/projects/dev-harness/tasks/UND-009/egm-research.md`) |
| Rodada | 2 de 2 |
| Veredito | approved |
| Atualizado | 2026-10-28 |

A rodada 2 reportou 0 achados bloqueantes. Respostas brutas das duas rodadas: `~/.harness/projects/dev-harness/tasks/UND-009/evidence/validation-r1/agy-reply.raw.md` e `validation-r2/agy-reply.raw.md`. O bloco `## Compreensão` não está neste relatório: o bloco é proposto por este PRD e ainda não estava em vigor quando as rodadas rodaram.

## Achados

- **P1** [applied] [conflict] A definição de pacote em R1 e R3 contradiz os contratos existentes de montagem do prompt e de despacho.
  - Local: seção 3.1, regras R1 e R3
  - Resolução: O autor emendou a Seção 2 Parte A, R1 e R3 para separar a carga do pacote de contexto (restrita por R1-R9) da montagem do prompt (que mantém o corpo do papel, o texto da skill, a cláusula anti-delegação, o contrato de veredito e o bloco de leitura permitida, conforme a linha 94 do coordinator), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] R26 exclui `fixed` dos achados assentados, o que torna impossíveis as rodadas limpas sob R30 depois de uma correção de bug.
  - Local: seção 3.3 (R26) e seção 3.4 (R30)
  - Resolução: O autor atualizou R30 para definir uma rodada limpa como a que não tem nenhum achado Critical ou Major `open`, `reopened` ou `escalated`, o que permite que achados `fixed` confirmados convirjam sem atrito, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] A atribuição de ids de achado em R19 não preserva de forma explícita ids estáveis para achados não resolvidos entre rodadas.
  - Local: seção 3.3, regra R19
  - Resolução: O autor emendou R19 para dizer: "Um achado levado a uma rodada posterior mantém o id original. Um achado nunca recebe um segundo id", rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [conflict] A regra R39 não está marcada como provisória, embora a Decisão D2(b) apresente o papel de comparação como decisão em aberto do dono.
  - Local: preâmbulo da seção 3, seção 3.5 (R39) e Apêndice (D2)
  - Resolução: O autor acrescentou R39 à lista de regras provisórias no preâmbulo da Seção 3 e ao bloco de decisões pendentes no Apêndice, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Violações do limite de tamanho de frase, procedural e descritiva, nas regras de escrita de PRD da ASD-STE100.
  - Local: seção 1 (Problema C), seção 3.1 (R3, R8), seção 3.7 (R55)
  - Resolução: O autor quebrou as frases longas do Problema 1.C, de R3, de R8 e de R55 em afirmações mais curtas, abaixo de 20 e de 25 palavras respectivamente, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] (Minor) Violações do limite de tamanho de frase, procedural e descritiva, nas regras de escrita de PRD da ASD-STE100, em regras recém-editadas.
  - Local: seção 3.2 (R14), seção 3.5 (R43), seção 3.7 (R57, R59)
  - Evidência: Várias frases descritivas e procedurais passam dos limites de tamanho da ASD-STE100 definidos nas regras de escrita de PRD (procedural <= 20 palavras, descritiva <= 25 palavras): a frase 1 de R14 tem 30 palavras; a frase 1 de R43 tem 38 palavras; a frase 1 de R57 tem 32 palavras; a frase 2 de R59 tem 26 palavras.
  - Sugestão: Quebrar essas frases em afirmações mais curtas, abaixo de 20 e de 25 palavras respectivamente.
  - Resolução: reportado como `pending` na resposta da rodada 2; corrigido depois, no registro do Coordenador (só estilo, nenhuma rodada usada).
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- Fidelidade às decisões do dono Q1 a Q6 e às premissas do Coordenador.
- Leitura do Coordenador sobre Q3 mantida visível como decisão em aberto D1 e não atribuída ao dono.
- Estruturação das decisões em aberto D1, D2 e D3 com opções, consequências e recomendações.
- Marcadores de provisório em R30-R35, R39, R47 e R51.
- Testabilidade das regras de R1 a R60.
- Estados de achado assentado, tratamento de oscilação, condições de rodada limpa e isolamento do estado avaliado.
- Todas as citações de arquivos do dh e números de linha conferidas contra `source3-cited-dh-passages.md` e confirmadas 100% corretas.
- Ausência de excesso de escopo (nenhum agente, comando ou dependência de runtime novo).
- Bloco de compreensão corretamente limitado ao `document-validator`.
- Cobertura obrigatória das seis categorias verificada.
