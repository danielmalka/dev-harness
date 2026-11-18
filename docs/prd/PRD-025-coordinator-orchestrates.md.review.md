# Revisão: PRD-025

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-025-coordinator-orchestrates.md |
| Fonte | Decisões do dono (2026-10-16) e card do roadmap "O Coordenador orquestra, não executa" (docs/roadmap.html ~L1726) |
| Rodada | 2 de 2 |
| Revisor | cli:agy/gemini-3.8-flash-high |
| Veredito | `approved` |
| Atualizado | 2026-10-16 |

## Achados

- **P1** [applied] [conflict] A delegação da saída de consolidação e do local.yaml na Solução contradiz as regras de posse da memória, a R2, a linha 54 e a R7
  - Local: seção 2, linha 52
  - Resolução: O autor atualizou a tabela da seção 2, atribuindo candidate-*.md ao docs-guide e mantendo a escrita da consolidação e do local.yaml com o Coordenador (despachando ao qa-verifier só o git check-ignore, conforme a linha 48), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P2** [applied] [conflict] A lista de permissão da R5 que reusa protectedHit exclui outros arquivos do harness e contradiz a R7
  - Local: seção 3, linha 76 (Regra R5)
  - Resolução: O autor trocou o reúso de protectedHit por uma checagem dedicada que cobre qualquer caminho dentro do harness dir (.harness/ ou <home>/projects/<name>/), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P3** [applied] [ambiguity] O mecanismo de aviso não especificado e a semântica de arquivo novo no Write produzem implementações incompatíveis
  - Local: seção 3, linha 76 (Regra R5) e Apêndice, linha 156
  - Resolução: O autor definiu os gatilhos do Write (mais de 5 linhas ou criação de arquivo novo) e esclareceu o contrato de teste como uma função pura observável em hooks/lib/ que devolve o texto do aviso ou undefined, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P4** [applied] [gap] A seção 4 Docs omite .agents/docs-guide.md e .skills/document-review/SKILL.md
  - Local: seção 4, linhas 90-104 e seção 3, linha 79 (Regra R8)
  - Resolução: O autor acrescentou .agents/docs-guide.md e .skills/document-review/SKILL.md à lista de atualização da seção 4, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P5** [applied] [organization] O ticket sugerido T-2503 usa a lane teste com um marcador (...)
  - Local: seção Apêndice, linha 137
  - Resolução: O autor trocou as reticências do marcador pela lane teste concreta (unitário e integração), rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high
- **P6** [applied] [organization] Frases descritivas na R5 e na R10 passam dos limites de palavras do ASD-STE100
  - Local: seção 3, linhas 76 e 81
  - Resolução: O autor quebrou as frases descritivas da R5 e da R10 em frases concisas de menos de 25 palavras, rodada 2.
  - Reportado por: cli:agy/gemini-3.8-flash-high

## Não levantado

- R1: O Coordenador não lê o repositório, não faz grep, não faz build, não valida nem roda checagens por conta própria; despacha repo-scout e qa-verifier.
- R2: O Coordenador mantém os próprios registros (MEMORY.md, EPOCHAL.md, RISKS.md) e a mudança de status do ticket com `Concluído em / Done on`.
- R3: Proposta de limite de 5 linhas alteradas para edições diretas em arquivos existentes, sem arquivo novo.
- R4 e R11: A mecânica de revisão por CLI é despachada ao qa-verifier ou à lane infra do builder, com a regra de congelamento obrigando o delegado e o Coordenador consolidando os vereditos.
- R5, R6 e R7: Especificações da regra de aviso do mod: só avisa, nunca nega, falha aberta, ignora especialistas (agentId presente), exclui o Bash e preserva as regras de negação da proteção de memória.
- R8: O docs-guide escreve o Resumo executado do PRD e persiste <document>.review.md.
- R9: Definição de "Julgar" restrita a ler relatórios e a evidência citada, com as conferências por amostragem convertidas em despachos.
- R10: Acréscimo de uma regra de escopo de uma linha ao AGENTS.md, que prevalece sobre o trabalho trivial feito em linha.
- R12: Contagens invariantes de 13 agentes e 19 + 1 comandos sem mudança.
- R13: Versão 0.31.0 em internal/build/build.go, marketplace.json e na ref do git.
- R14: Invariante de que o dist/ é regenerado exclusivamente por go run ./cmd/dh build.
- R15 e R16: Comandos de verificação determinísticos (dh validate, go test, prova de clone limpo por causa de internal/build/build.go).
- R17: Proibição de caminhos privados da máquina, verificada contra checkPrivatePaths.
- R18: Chamadas de CLI externa executadas como comandos únicos e isolados que começam pelo nome do binário, sem invólucros nem encadeamento.

## Emenda pós-validação (2026-10-16)

- A R18 mudou depois da rodada 2 (`approved`). Um qa-verifier despachado mediu que o diretório de trabalho de um subagente não persiste entre chamadas do Bash; por isso a sequência de três comandos com `cd` para as chamadas de clone foi trocada pela flag de diretório de trabalho do próprio binário, com `cd "$CLONE" && <binary>` como a única exceção para um binário sem essa flag. A receita do agy ganhou a regra de tratar uma resposta completa em disco como o resultado e de parar um processo que sobrevive ao `--print-timeout`.
- Não revalidado: o limite de duas rodadas está usado. Disposição do dono solicitada.
