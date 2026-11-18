# Revisão · PRD-009 · rodada 1

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

- Revisor: document-validator (claude, opus solicitado; modelo efetivo não verificado)
- Data: 2026-09-30
- Fonte: pedido do dono (cartão "Este próprio roadmap como artefato vivo"); decisão do AGENTS.md de 2026-09-30 (opção A); decisões de acompanhamento do Coordenador em nome do dono.

## Veredito
- `changes required`
- Achados bloqueantes: 2

## Achados
- P1 · `conflict` · `blocking` · §3 R10 (e §2). Os ativos genéricos do kit (a tabela de prontidão da `delivery-readiness`, o despacho do docs-guide pelo Coordenador) carregariam uma linha de roadmap específica do dev-harness que nomeia `docs/roadmap.html` e `dh validate`, enquanto a R1 diz que os projetos consumidores não têm roadmap; o gate "Atendida quando" não tem texto de não aplicável para esse caso. Sugestão: linha condicional e neutra quanto ao caminho ("quando o projeto mantém um roadmap ligado ao seu changelog …; não aplicável quando o projeto não mantém roadmap"), frase exata do Coordenador com a mesma condição, e o Verificável da R10 nomeia o caso sem roadmap. Reportado por: claude
- P2 · `ambiguity` · `blocking` · §3 R1. Só um dos dois roadmaps presente não está especificado (ignorar vs. erro). Sugestão: erro que nomeia o arquivo ausente; ignorar só quando ambos estão ausentes; teste em Go para isso. Reportado por: claude
- P3 · `ambiguity` · `non-blocking` · §3 R8 vs R5/R7. Linha para um cartão ausente e arquivo para uma divergência pt/en não especificados. Sugestão: cartão ausente → a linha `<h2 id="entregue">`; R7 → citar `docs/en/roadmap.html`. Reportado por: claude
- P4 · `gap` · `non-blocking` · §4 Docs. O espelho `docs/en/tutorials/` falta na linha condicional. Reportado por: claude

## Não levantado
- 22 títulos do CHANGELOG ↔ 22 cartões de versão por idioma, já do mais novo para o mais antigo; dois cartões de detalhe por idioma; os formatos de R2–R4 batem com a árvore; linhas 72/76 do coordinator.md; o release-manager não tem checklist enumerada; a reincorporação da R12 e a ausência de eval paga são coerentes com o AGENTS.md; nenhuma contradição com o AGENTS.md.

---

# Revisão · PRD-009 · rodada 2 (final)

- Revisor: document-validator (claude, opus solicitado; modelo efetivo não verificado). Data: 2026-09-30.

## Veredito
- `approved`
- Achados bloqueantes: 0

## Assentamento da rodada 1
- P1–P4: aplicados e fechados (R10/§2 condicional e neutra quanto ao caminho, com frases e greps exatos; R1 ignora só quando ambos os roadmaps estão ausentes, um presente é erro; R8 cita a linha `<h2 id="entregue">` e `docs/en/roadmap.html`; a linha de Docs lista `docs/en/tutorials/`).

## Novos achados
- P5 · `ambiguity` · `non-blocking` · §3 R10. "passa na checagem do projeto" e "um cartão por release" pressupõem uma checagem e um roadmap em cartões em todo consumidor. Sugestão: "a checagem de roadmap do projeto, quando houver uma"; "o novo release registrado", mantendo "um cartão por release" só dentro do parêntese do dev-harness. Reportado por: claude
- Coordenador: P5 aplicado literalmente à R10 após a rodada 2 (só a redação; o software construído é o mesmo).
