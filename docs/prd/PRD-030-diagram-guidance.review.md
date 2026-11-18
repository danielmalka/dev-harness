# Revisão: PRD-030

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

## Rodada 1

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-030-diagram-guidance.md |
| Fonte | docs/roadmap.html:1726-1741, docs/en/roadmap.html:1726-1741, .agents/docs-guide.md, .skills/doc-template-html/SKILL.md, .skills/doc-template-html/scripts/stamp.sh, AGENTS.md, docs/prd/PRD-029-independent-review-and-drift-report.md, palavras do dono |
| Rodada | 1 de 2 |
| Veredito | changes required |
| Atualizado | 2026-10-28 |

Revisor: cli:agy/gemini-3.8-flash-high (etapa de documento). Achados bloqueantes: 1. Resposta bruta: `tasks/AUTO-3/evidence/validation-prd-r1/reply.md` (harness dir).

## Compreensão
- O que o documento propõe: O documento estabelece, no docs-guide, uma orientação concisa para escrever diagramas Mermaid dentro de blocos `<pre class="mermaid">` e linkar diagramas Excalidraw como arquivos `.excalidraw` separados, sem embuti-los. Ele atualiza quatro páginas HTML escritas à mão para trocar as importações flutuantes de versão do Mermaid por uma referência fixada na versão 11.17.2, igual à versão fixada no `stamp.sh`. Ele garante que nenhuma referência flutuante ao Mermaid sem rastreio fique fora das saídas de distribuição e dos registros históricos, e atualiza o card de diagrama do roadmap para refletir o escopo entregue. Por fim, alinha os identificadores de versão em 0.35.0, regenera os arquivos de distribuição do pacote e confirma que as contagens de agentes, comandos e skills continuam as mesmas.
- Como o sistema atual funciona, como o documento o declara: O roadmap tem um card planejado para skills especializadas de diagrama em Excalidraw e Mermaid, mas o `docs-guide` hoje não traz regra sobre quando nem como desenhar diagramas. Embora o `stamp.sh` e os `pdocs` do dashboard já fixem o Mermaid na versão 11.17.2, quatro páginas HTML de documentação escritas à mão ainda carregam o Mermaid por uma tag flutuante de versão maior. Além disso, o texto do card do roadmap descreve o Mermaid como a única exceção de CDN permitida, atrasado em relação ao `AGENTS.md`, onde o Google Fonts é reconhecido como segunda exceção. Uma skill externa de design de diagramas foi avaliada, mas rejeitada, porque a emissão de SVG inline e as dependências externas de runtime conflitam com os princípios do kit.
- Seções não parafraseadas: nenhuma

## Achados

- **P1** [applied] [conflict] (blocking) O critério de aceite "todos dizem 0.35.0" contradiz o ref "v0.35.0" e deixa indefinidos os cinco lugares da versão.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R5 (linhas 46-47)
  - Evidência: R5 linha 46 especifica `ref v0.35.0` em `.claude-plugin/marketplace.json`, enquanto o critério da linha 47 afirma que os cinco lugares da versão dizem `0.35.0`. Além disso, R5 descreve 7 alvos distintos de versão no repositório.
  - Sugestão: Nomear de forma explícita cada local de alvo e a string exata esperada (`0.35.0` ou `v0.35.0`) nos critérios de aceite.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] A posição da instrução depois do passo 2 do Procedimento deixa ambíguos a numeração dos passos e a referência cruzada da linha 75.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R1 (linha 29)
  - Evidência: Renumerar os passos do Procedimento tira de sincronia a referência em `.agents/docs-guide.md:75` ("passo 5 do Procedimento").
  - Sugestão: Esclarecer que a instrução entra como parágrafo sem número ou subitem dentro do passo 2, sem renumerar os passos 1 a 5.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [organization] A seção 4 Docs omite da lista de itens os dois arquivos de tutorial modificados por R2.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 4 (linhas 57-63)
  - Evidência: `docs/tutoriais/templates.html` e `docs/en/tutorials/templates.html` são modificados em R2 e citados na prosa, mas faltam na lista.
  - Sugestão: Acrescentar os dois arquivos de tutorial à lista de itens da Seção 4.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [organization] O campo Status do cabeçalho está marcado como aprovado antes de a validação terminar.
  - Local: docs/prd/PRD-030-diagram-guidance.md linha 5
  - Evidência: `templates/en/PRD.md` prescreve `in review` ou `draft` enquanto o documento passa por validação.
  - Sugestão: Mudar `Status` para `in review`.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Frases de R4 e R5 passam do limite de 25 palavras das regras da ASD-STE100.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R4 (linha 42) e R5 (linha 46)
  - Evidência: Frases descritivas de R4 e R5 têm 28 palavras cada uma.
  - Sugestão: Dividir cada frase em duas frases mais curtas, abaixo de 25 palavras.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- Lacunas: Todos os pontos de `docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`, das palavras do dono e de `AGENTS.md:29` estão mapeados em R1-R5 e em Fora do escopo.
- Excessos: Nenhum agente, skill ou comando não pedido foi introduzido.
- Regras de diagrama (R1): O Mermaid em `<pre class="mermaid">` como padrão e o Excalidraw como arquivo `.excalidraw` separado e linkado, sem script nem embutido, honram corretamente a decisão do dono de 2026-09-25 e as regras da `doc-template-html`.
- Importações fixadas (R2): A versão 11.17.2 confere com `stamp.sh:89` e `internal/dashboard/pdocs.go:18` e mira corretamente as quatro páginas HTML escritas à mão com tags flutuantes.
- Expurgo das referências flutuantes (R3): A verificação de árvore limpa por `grep` exclui corretamente os arquivos gerados em `dist/` e os arquivos históricos em `docs/prd/`.
- Sincronia do card do roadmap (R4): Atualizar o status para entregue na 0.35.0, remover "única exceção de CDN" em favor das duas exceções de `AGENTS.md:29` e manter pt e en em sincronia honra corretamente as regras de manutenção do roadmap.
- Checagens determinísticas (R6): Rodar `dh validate`, `go test ./...` e `dh build` e preservar as contagens invariantes (13 agentes, 19+1 comandos, 22 skills) confere com os padrões de verificação do repositório.

## Rodada 2

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-030-diagram-guidance.md |
| Fonte | docs/roadmap.html:1726-1741, docs/en/roadmap.html:1726-1741, .agents/docs-guide.md, .skills/doc-template-html/SKILL.md, .skills/doc-template-html/scripts/stamp.sh, AGENTS.md, docs/prd/PRD-029-independent-review-and-drift-report.md, palavras do dono |
| Rodada | 2 de 2 |
| Veredito | approved |
| Atualizado | 2026-10-28 |

Revisor: cli:agy/gemini-3.8-flash-high (etapa de documento). Achados bloqueantes: 0. Resposta bruta: `tasks/AUTO-3/evidence/validation-prd-r2/reply.md` (harness dir). As correções do PRD entre as rodadas foram feitas pelo product-discovery.

## Compreensão
- O que o documento propõe: O PRD acrescenta ao `docs-guide` uma instrução concisa que estabelece o Mermaid dentro de `<pre class="mermaid">` como formato padrão de diagrama e restringe os diagramas Excalidraw a arquivos `.excalidraw` externos linkados. Ele atualiza quatro arquivos HTML de documentação escritos à mão, trocando as importações flutuantes do Mermaid por uma referência fixada na versão exata 11.17.2, igual à do `stamp.sh`. Ele garante que nenhuma importação flutuante do Mermaid continue rastreada no código e atualiza o card de diagrama do roadmap para refletir o escopo entregue na versão 0.35.0. Ele atualiza os identificadores de versão do repositório e as entradas do changelog para 0.35.0, regenera os arquivos de distribuição do pacote e verifica que as contagens de agentes, comandos e skills continuam estritamente as mesmas.
- Como o sistema atual funciona, como o documento o declara: O roadmap lista um card planejado para skills especializadas de diagrama em Excalidraw e Mermaid, enquanto o `docs-guide` hoje não traz instruções sobre quando nem como incluir diagramas. Embora o `stamp.sh` e os `pdocs` do dashboard fixem o Mermaid na versão 11.17.2, quatro arquivos HTML de documentação escritos à mão continuam importando o Mermaid por uma tag flutuante de versão maior. O card do roadmap ainda descreve o Mermaid como a única exceção de CDN permitida, atrás do `AGENTS.md`, que especifica duas exceções de CDN, entre elas o Google Fonts. Uma skill externa de design de diagramas foi avaliada, mas rejeitada pelo dono, porque a emissão de SVG inline e as dependências externas contradizem a arquitetura do kit e a ADR-003.
- Seções não parafraseadas: nenhuma

## Achados

- **P1** [applied] [conflict] (blocking) O critério de aceite "todos dizem 0.35.0" contradiz o ref "v0.35.0" e deixa indefinidos os cinco lugares da versão.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R5 (linhas 46-47)
  - Resolução: O autor atualizou o critério de aceite para citar de forma explícita os arquivos nomeados em R5 e levou em conta o prefixo `v0.35.0` no ref do marketplace ("Quando os arquivos citados nesta regra são lidos, então cada campo de versão diz 0.35.0 (o ref do marketplace diz v0.35.0)"), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [ambiguity] A posição da instrução depois do passo 2 do Procedimento deixa ambíguos a numeração dos passos e a referência cruzada da linha 75.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R1 (linha 29)
  - Resolução: O autor esclareceu que a instrução é acrescentada no fim do passo 2 do Procedimento, no mesmo item da lista, e afirmou de forma explícita que os números dos passos não mudam, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [organization] A seção 4 Docs omite da lista de itens os dois arquivos de tutorial modificados por R2.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 4, linha 62
  - Resolução: O autor acrescentou os dois arquivos de tutorial à lista da Seção 4 ("docs/tutoriais/templates.html e docs/en/tutorials/templates.html: só a fixação da importação (R2)"), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [pending] [organization] O campo Status do cabeçalho está marcado como aprovado antes de a validação terminar.
  - Local: docs/prd/PRD-030-diagram-guidance.md linha 5
  - Evidência: A tabela do cabeçalho declara `Status: approved (...)`, enquanto `templates/en/PRD.md` prescreve `draft` ou `in review` enquanto o documento passa por validação.
  - Sugestão: Manter como registro de organização não bloqueante; o status muda formalmente quando o Coordenador aceitar. Nota (Coordenador, depois desta rodada): a linha Status do cabeçalho foi definida como `delivered on 2026-10-28` no fim do ciclo do PRD; o estado de P4 para o revisor não mudou.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] Frases de R4 e R5 passam do limite de 25 palavras das regras da ASD-STE100.
  - Local: docs/prd/PRD-030-diagram-guidance.md seção 3, R4 (linha 42) e R5 (linha 46)
  - Resolução: O autor dividiu as frases descritivas de R4 e R5 em afirmações concisas, de menos de 25 palavras cada uma, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- Lacunas: Todos os pontos de origem de `docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`, `AGENTS.md:29` e do pedido do dono ("Vamos seguir com a B e pode fixar a versão do mermaid junto") são rastreados até R1-R6 e Fora do escopo.
- Ambiguidades: Todas as instruções e regras trazem implementações únicas e sem ambiguidade. A posição da instrução do docs-guide no fim do passo 2 do Procedimento preserva a numeração dos passos e as referências cruzadas de `.agents/docs-guide.md:75`.
- Conflitos: As especificações de string de versão em R5 conferem com as exigências de tag do Git (`v0.35.0` para o ref do marketplace e `0.35.0` para as versões do pacote). A versão fixada `11.17.2` confere com `stamp.sh:89` e `internal/dashboard/pdocs.go:18`. As afirmações sobre as exceções de CDN conferem com `AGENTS.md:29`.
- Excessos: Nenhum agente, skill ou comando novo é introduzido; o Fora do escopo rejeita de forma explícita embutir, exportar SVG e depender de diagramas externos, de acordo com a decisão B do dono e com a ADR-003.
- Critérios fracos: Todos os critérios de aceite sob R1-R6 usam estruturas determinísticas "Quando ..., então ...", verificáveis por casamento exato de string, status HTTP, conteúdo de arquivo ou execução de ferramenta.
- Organização: As listas de requisitos e de impacto na documentação estão bem estruturadas e categorizadas.
- Estilo de escrita do PRD: Todas as frases descritivas respeitam o limite de 25 palavras especificado pelas diretrizes da ASD-STE100.
