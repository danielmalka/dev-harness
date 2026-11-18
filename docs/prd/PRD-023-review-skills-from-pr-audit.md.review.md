# Revisão: PRD-023

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-023-review-skills-from-pr-audit.md |
| Fonte | Escopo aprovado pelo dono, chat 2026-10-16 (despacho do Coordenador); decisões do dono 2026-10-16 rodada 2 (cópia permitida; emenda security-audit, itens 7-13 + refutador + exclusões); AGENTS.md |
| Rodada | 2 de 2 |
| Veredito | approved |
| Atualizado | 2026-10-16 |

## Achados

- **P1** [fechado pelo dono] [conflict] As frases da §2 quase copiam o pr-audit.
  - Resolução: o dono liberou akitaonrails/my-skills para uso literal (2026-10-16); texto de licença removido; a regra de procedência da R14 substitui "nenhum texto copiado" (§6:134).
  - Reportado por: claude
- **P2** [applied] [conflict] O aumento de versão era opcional e uma decisão em aberto.
  - Resolução: §2:61, R21 e §4:96-99 fixam 0.30.0 pela regra de tag do AGENTS.md; §6:137 e Apêndice:156 dizem ambos "nenhuma".
  - Reportado por: claude
- **P3** [applied] [weak criterion] Lista de dist da R9, prova de clone limpo, generated_at.
  - Resolução: a R16 lista os caminhos de dist; a R17 permite a linha generated_at; a R18 é a prova de clone limpo.
  - Reportado por: claude
- **P4** [applied] [ambiguity] Nomes de arquivo do repositório no texto da skill entregue.
  - Resolução: §2:30 e o item 2 mantêm exemplos genéricos; a R2 faz grep de que `status-cases|fixture_parity` estão ausentes de .skills.
  - Reportado por: claude
- **P5** [applied] [weak criterion] Greps que já davam hit; limites por item.
  - Resolução: padrões mais estreitos; a R19 tem limites por arquivo; a R20 cobre os quatro arquivos.
  - Reportado por: claude
- **P6** [applied] [ambiguity] Sobreposição do item 5 com as linhas Dependencies e Configuration.
  - Resolução: a §2:43 estende checklist.md:42-43, sem linha duplicada; a R5 checa.
  - Reportado por: claude
- **P7** [aplicado em parte] [organization] (`non-blocking`) Tamanho das frases.
  - Resolução: redação da R10 corrigida (R17). A §2:57 ainda tem uma frase procedural de 22 palavras.
  - Reportado por: claude
- **P8** [pending] [organization] (`non-blocking`) Os itens da §2 estão na ordem 1-5, 7-13, 6; os "lugares que a seção 2 nomeia" da R5 não têm lugar no SKILL.md para os itens de correção de segurança, prova e refutador do item 5.
  - Local: §2, R5
  - Sugestão: restaure a ordem numérica; nomeie o passo de cada item do SKILL.md do item 5
  - Reportado por: claude
- **P9** [pending] [weak criterion] (`non-blocking`) R7, R8, R10 e R13 fazem grep de uma âncora só; motivo de descarte, ausência de gate, mocks que escondem efeitos colaterais, texto recuperado como dado e chaves de cache/índice com escopo são checados só por leitura.
  - Local: §3 R7, R8, R10, R13
  - Sugestão: acrescente cláusulas "e a linha também diz ..."
  - Reportado por: claude
- **P10** [pending] [organization] (`non-blocking`) A nota da R5 põe o hit existente de `author's reasoning` em security-review, mas ele está em code-review/SKILL.md:178; R19 "no máximo" contra "ficam abaixo de".
  - Local: §3 R5, R19
  - Reportado por: claude

## Não levantado

- Mapa da emenda: item 7 -> R7, 8 -> R8, 9 -> R9, 10 -> R10, 11 -> R11, 12 -> R12 (passo 10), 13 -> R13; refutador (despacho do Coordenador, só o achado e a evidência do código, o revisor nunca cria agentes) -> §2:46 e R5.
- Todas as exclusões do security-audit estão em §5:120-126 (auditoria de 6 fases, triagem de alertas/PATCH, execução de scanners/gate de testes, ecosystem-checks.md, security-surface.sh / nenhum subcomando do dh, template de saída, política de divulgação); a R22 checa script, subcomando, ecosystem-checks.
- Os itens 1-6 do pr-audit e suas exclusões continuam mapeados (§2, R1-R6, §5:111-116); contagens 13 / 19 + 1 na R15.
- Os greps anteriores à mudança retornam 0 nos três diretórios de skill para: `sibling call`, `hand-maintained`, `belongs in the tree`, `new entry point`, `assertion failure`, `narrowed`, `baseline`, `weaken`, `hardening`, `tool policy`, `control case`, `fail closed`, `reporter` e `author's reasoning` em security-review. `fingerprint` dá hit só em security-review/SKILL.md:117, enquanto a R11 faz grep em checklist.md.
- Os limites cabem na contagem de itens: code-review 14 para os itens 1-4 mais as fontes; regression-testing 4 para o item 6 mais as fontes; SKILL de security-review 18 para os itens do 5, 7, 9, 12 mais as fontes; checklist 24 para as linhas novas e estendidas.
- Os fatos de release se sustentam: `defaultVersion` em build.go:20, `version` do marketplace.json em :14 e `ref` em :24; as versões de plugin.json e harness-manifest.json de dist vêm do build. O harness-manifest.json da raiz é regenerado pelo build (build.go:238), não editado à mão. GENERATED.txt é estático.
- Nenhum resto de licenciamento; §6 e Apêndice dizem ambos que não há decisão em aberto; a evidência da §4 Documentação não mudou e se sustenta.
- A forma segue templates/en/PRD.md mais o precedente da §5/§6 do PRD-022; em inglês; as regras trazem linhas Verificável.

## Emenda de docs, 1 de 1 (2026-10-16)

- Status da revisão: `approved` (0 `blocking`). Escopo: só a §4 Documentação e a R16.
- **P11** [pending] [organization] (`non-blocking`) O grep de Evidência da §4:92 diz que "lista só estes arquivos", mas depois desta mudança README.md e README.en.md (linha 20, resumo do PRD-023) também casam.
  - Local: §4 linha 92
  - Sugestão: "na base, antes desta mudança, lista só estes arquivos: ..."
  - Reportado por: claude

### Não levantado

- A §4:100 acrescenta README.md e README.en.md (linha 20, versão 0.30.0 / tag v0.30.0, resumo que nomeia o PRD-023). Isso casa com o precedente de release e com a árvore de trabalho; os READMEs saíram de "Não afetados" (§4:106).
- A R16 continua consistente: "os arquivos de release da seção 4" agora inclui os READMEs; eles não são copiados para dist/, então a lista de caminhos de dist da R16 não muda.
- O grep de Evidência é evidência anterior à mudança, não uma regra Verificável; nenhuma outra linha da §4 mudou.
