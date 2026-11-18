# PRD-030 · Orientação de diagramas para o docs-guide e importações fixadas do Mermaid

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-28 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-28 / 2026-10-28 |
| Tickets | a definir pelo planejador |
| Origem | Decisão do dono em 2026-10-28 (opção B + fixação), tarefa AUTO-3 |
| Depende de | PRD-029 (base de versão 0.34.0, entregue) |
| Entrega | Uma pull request, versão 0.35.0 |

## 1. Problema

Três fatos, uma capacidade: como o kit orienta diagramas em documentos.

- Fato: o card do roadmap "Skills de diagrama especializadas (Excalidraw e Mermaid)" (`docs/roadmap.html:1726-1741`, `docs/en/roadmap.html:1726-1741`) está planejado e não tem PRD.
- Fato: `.agents/docs-guide.md` não diz nada sobre diagramas. Uma execução do `docs-guide` não tem regra sobre quando desenhar nem em qual formato.
- Fato: quatro páginas HTML escritas à mão importam o Mermaid com a tag flutuante `mermaid@11`. A regra do kit fixa uma versão exata (`.skills/doc-template-html/scripts/stamp.sh:89` usa `11.17.2`). O texto do card do roadmap (linha 1736 do pt, linha 1737 do en) ainda diz que o Mermaid é a única exceção de CDN. O `AGENTS.md` (linha 29) diz que são duas: Mermaid e Google Fonts.

Decisão (dono, 2026-10-28): o kit não importa a skill externa `cathrynlavery/diagram-design`. Fato: ela emite SVG embutido, precisa de Python e Playwright, e tem 45 tipos de diagrama. Isso conflita com a regra `Never embed pre-rendered SVG` da `doc-template-html` e com a ADR-003 (kit enxuto). O dono pode instalá-la fora do kit.

## 2. Solução

Depois da entrega, o `docs-guide` tem uma regra curta sobre diagramas. O Mermaid é o formato padrão, escrito em `<pre class="mermaid">`. O Excalidraw só é permitido como um arquivo `.excalidraw` separado que o documento linka. As quatro páginas importam uma única versão exata do Mermaid. O roadmap mostra o card como entregue na 0.35.0. Nenhuma skill, agente ou comando é acrescentado.

## 3. Regras

- R1 `.agents/docs-guide.md` traz uma instrução curta. É uma frase acrescentada no fim do passo 2 do Procedimento (mesmo item da lista), então nenhum passo é renumerado. Ela diz: use um diagrama Mermaid (`<pre class="mermaid">`, como a `doc-template-html` define) onde ele ajuda o leitor a seguir um fluxo ou uma estrutura. Use um diagrama Excalidraw só como um arquivo `.excalidraw` separado que o documento linka. Nunca o embuta. Sem script, sem exportação (decisão do dono em 2026-09-25). Nenhuma skill de diagrama é acrescentada.
  - Quando o `docs-guide` lê o Procedimento, então a frase termina o passo 2 e os números dos passos não mudam.
  - Quando um documento precisa de uma imagem de fluxo ou de estrutura, então a instrução cita o Mermaid primeiro.
  - Quando o Excalidraw é usado, então a instrução permite só um link para um arquivo separado.
- R2 As quatro importações usam a versão exata `mermaid@11.17.2` com o nome de arquivo `mermaid.esm.min.mjs`. Os arquivos são `docs/roadmap.html:1897`, `docs/en/roadmap.html:1901`, `docs/tutoriais/templates.html:983` e `docs/en/tutorials/templates.html:983`. A versão é a mesma de `stamp.sh:89`. A URL responde HTTP 200.
  - Quando um dos quatro arquivos é lido, então sua importação nomeia `mermaid@11.17.2/dist/mermaid.esm.min.mjs`.
  - Quando a URL é requisitada, então a resposta é HTTP 200.
  - Quando as quatro versões são comparadas com `stamp.sh:89`, então são iguais.
  - Limite aceito (decisão): essas importações não têm Subresource Integrity. Um `import` ES estático não pode carregá-lo. A fixação em uma versão exata é a mitigação.
- R3 Nenhuma referência `mermaid@11/` (flutuante) permanece nos arquivos rastreados, exceto sob `dist/` e no histórico em `docs/prd/`.
  - Quando `grep -rn "mermaid@11/"` roda na árvore sem `dist/` e `docs/prd/`, então não devolve nenhuma linha.
  - Quando `go run ./cmd/dh build` roda, então `dist/` não tem nenhuma tag flutuante que as fontes não tenham mais.
  - Quando uma página nova importa o Mermaid, então ela usa a versão exata de R2.
- R4 O card do roadmap, em pt e en, é mantido em sincronia à mão (o `dh validate` ignora os cards de detalhe). Seu status é entregue na versão 0.35.0. Sua descrição diz o que foi entregue. Ela cita a instrução do `docs-guide`, o Excalidraw como arquivo separado linkado, nenhuma skill nova e o `diagram-design` externo deixado ao dono. As palavras "única exceção de CDN" (pt 1736) e `the only allowed CDN exception` (en 1737) são substituídas. O texto novo usa a redação de duas exceções de `AGENTS.md:29` (Mermaid e Google Fonts).
  - Quando o card é lido em pt e en, então os dois mostram entregue, 0.35.0 e os mesmos fatos.
  - Quando o card é procurado pelas palavras antigas da exceção, então nenhuma permanece.
  - Quando o card é lido, então ele cita duas exceções de CDN.
- R5 A versão 0.35.0 é definida em `internal/build/build.go:20`, em `.claude-plugin/marketplace.json` (versão e ref `v0.35.0`) e em `harness-manifest.json`. `CHANGELOG.md` tem uma entrada datada de 2026-10-28. `docs/roadmap.html` e `docs/en/roadmap.html` têm chips atualizados e um novo card de versão `Entregue` / `Delivered`. O card é construído como o card da 0.34.0 (linhas 703-709). Sua faixa de `Sources` vai até a 0.35.0. `dist/` é regenerado só por `go run ./cmd/dh build`, nunca à mão.
  - Quando os arquivos citados nesta regra são lidos, então cada campo de versão diz 0.35.0 (o ref do marketplace diz v0.35.0).
  - Quando o roadmap é lido em pt e en, então o card novo existe com a mesma estrutura da 0.34.0.
  - Quando `dist/` é comparado com a saída do build, então não há diferença.
- R6 As checagens passam. As contagens ficam em 13 agentes, 19 comandos mais 1 comando de mod e 22 skills.
  - Quando `dh validate` roda, então passa.
  - Quando `go test ./...` roda, então passa.
  - Quando `go run ./cmd/dh build` roda depois da regeneração, então `dist/` não mostra diferença e as contagens não mudam.

Fora de escopo: qualquer skill ou agente novo; geração de SVG; exportação ou embutimento do Excalidraw; SRI para importações ESM; uma mudança no texto de histórico da ADR-003 ou do PRD-017; uma oferta de relatório de deriva; qualquer mudança em `stamp.sh` ou em `internal/dashboard/pdocs.go`.

## 4. Docs

- `.agents/docs-guide.md`: o fluxo documentado que foi tocado (comportamento do `docs-guide`, R1).
- `docs/roadmap.html` e `docs/en/roadmap.html`: card e card de versão (R4, R5).
- `CHANGELOG.md`: entrada da 0.35.0 (R5).
- `docs/tutoriais/templates.html` e `docs/en/tutorials/templates.html`: só a fixação da importação (R2).

A regra de diagrama da `doc-template-html` não muda. Os tutoriais e o README não precisam de texto novo além das importações fixadas de R2.

## Resumo executado

- **Entregue:** O `docs-guide` agora traz uma instrução de diagrama (Mermaid primeiro; Excalidraw só como um arquivo `.excalidraw` separado e linkado, nunca embutido), as quatro importações escritas à mão do Mermaid estão fixadas em `mermaid@11.17.2`, o card do roadmap está entregue e corrigido, e o kit está na 0.35.0; nenhuma skill ou agente novo foi acrescentado.
- **Regras:**

  | Regra | Veredito | Evidência |
  |---|---|---|
  | R1 | Honrada | `.agents/docs-guide.md`: a frase fica no fim do passo 2 do Procedimento, no mesmo item da lista; a linha `## Limits` ainda cita o passo 5 do Procedimento. Conferido em `tasks/T-3002/evidence/` (harness dir). |
  | R2 | Honrada | `docs/roadmap.html:1895`, `docs/en/roadmap.html:1898`, `docs/tutoriais/templates.html:983`, `docs/en/tutorials/templates.html:983` importam `mermaid@11.17.2/dist/mermaid.esm.min.mjs`, a mesma versão de `.skills/doc-template-html/scripts/stamp.sh:89`. O QA observou HTTP 200 para a URL. Limite aceito: uma importação ES estática não tem Subresource Integrity; a fixação exata é a mitigação. |
  | R3 | Honrada | `git grep 'mermaid@11/'` fora de `dist/` e `docs/prd/` não tem nenhuma ocorrência (conferido de novo na hora da documentação). |
  | R4 | Honrada | O card do roadmap está entregue em pt e en com os quatro fatos e as duas exceções de CDN (Mermaid e Google Fonts). A redação antiga da exceção sumiu dos dois roadmaps. |
  | R5 | Honrada | Versão 0.35.0 em `internal/build/build.go:20`, `.claude-plugin/marketplace.json` (`version` 0.35.0, `ref` v0.35.0), `harness-manifest.json`, a entrada do `CHANGELOG.md` datada de 2026-10-28, os chips, a faixa de `Sources` e um novo card de versão Entregue em pt e en. Uma reconstrução numa cópia reproduz `dist/` byte a byte, exceto a linha de timestamp `generated_at` do `harness-manifest.json`, que é saída esperada do build (`tasks/T-3002/evidence/host-checks/repro-diff.txt`, harness dir). |
  | R6 | Honrada | `go test ./...` com saída 0 e `dh validate .` com saída 0 (erros 0, 22 skills), `tasks/T-3002/evidence/host-checks/` (harness dir). Contagens: 13 agentes, 19 comandos mais 1 comando de mod, 22 skills. Vale a mesma nota sobre `dist/` de R5. |
- **Tickets:** T-3001 (backend, harness-maintainer): `done`. Rodada 1 da revisão de código por `cli:opencode/opencode-go/muse-spark-1.3-contributor`: `approve`, 0 achados. Rodada 2 de confirmação por `cli:codex/gpt-6-luna`: `request changes` com um achado F2-1 (`Major`, `R6 check evidence not observed in the reviewer's read-only sandbox`); o Coordenador o classificou como `noise` e o resolveu com evidência do host no mesmo estado avaliado. T-3002 (teste, qa-verifier): verde, matriz R1 a R6. Evidência: `tasks/T-3001/evidence/review-r1b/`, `tasks/T-3001/evidence/review-r2/`, `tasks/T-3002/evidence/` (harness dir).
- **Docs:** `.agents/docs-guide.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/tutoriais/templates.html`, `docs/en/tutorials/templates.html`, `CHANGELOG.md`.
- **Fora de escopo:** qualquer skill ou agente novo; geração de SVG; exportação ou embutimento do Excalidraw; SRI para importações ESM; uma mudança no texto de histórico da ADR-003 ou do PRD-017; uma oferta de relatório de deriva; qualquer mudança em `stamp.sh` ou em `internal/dashboard/pdocs.go`. Também: o card de detalhe entregue continua na seção "Planejado" do roadmap (anotado pelo QA como um possível `Minor`, decisão do dono pendente).
