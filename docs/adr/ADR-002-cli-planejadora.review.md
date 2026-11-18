# ADR-002 — revisão

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

Documento: `docs/adr/ADR-002-cli-planejadora.md`. Fontes: PRD-005 (RF-03/04/05/13, AC-05/06/15, §6), GRILL G6/G15/G18–G23, PLAN-011 fatia 1 e T-801, RISK-001, CLAUDE.md (AA1; exceção de 2026-09-28), `.skills/external-clis/SKILL.md`, `.agents/coordinator.md`, `internal/kit/evals.go`.
Profundidade: regra enxuta — um validador Claude (`document-validator`, opus solicitado, efetivo não verificado), duas rodadas.

## Rodada 1 (2026-09-28)
- Veredito: `changes required` — 1 `blocking`, 4 `non-blocking`.
- A1 — `conflict` — `blocking` — §5, linha 55. O sinal de falha "as duas cláusulas diferem em algo além de revisão → planejamento" dispara no texto correto fixado em T-801 (linhas 38-51), que diz "esta tarefa de planejamento" e quebra as primeiras linhas de outro modo. Correção: o sinal passa a ser "a cláusula de planejamento difere do texto fixado em T-801 / PLAN-011 fatia 1; comparadas entre si, diferem em algo além de 'esta revisão' → 'esta tarefa de planejamento' (ignorando a quebra de linha)".
- A2 — `ambiguity` — `non-blocking` — §4 linhas 48-49 e §5 linha 54. "o segundo par" sugere uma cópia da cláusula em coordinator.md, mas T-801 publica a cláusula de planejamento só em external-clis "A chamada de planejamento". Correção: dizer que ela é publicada uma única vez, que coordinator.md só a referencia, e que a sincronização manual é revisão × planejamento.
- A3 — `ambiguity` — `non-blocking` — §3 item 3 (linha 41). "referenciado" pode ser lido como valendo para o prompt, enquanto AC-05 exige o bloco "Ler é permitido" literalmente no prompt. Correção: copiado literalmente no prompt; referenciado, não reescrito, na skill.
- A4 — `organization` — `non-blocking` — §3 item 5 (linha 43) e §7. Faltam citações: RF-03, RF-05/G19, RF-13/AC-15; RF-13 está ausente do §7.
- A5 — `ambiguity` — `non-blocking` — §4 linha 47. "os mesmos sinais do revisor de CLI" contradiz o §3 item 4 (a tabela de veredito não se aplica). Correção: o mesmo registro de congelamento (`clean`/`reverted`) e as mesmas classes de não rodado.
- Não levantado: escopo, cabeçalho, forma do template, as citações do §1, G6/AA1/RISK-001, a justificativa do §3 e os itens 1, 2 e 4, a descrição de `checkAntiDelegationClause` e as afirmações no presente (todos verificados).

### Disposições do Coordenador
- A1–A5 aceitos; correção 1 de 2 devolvida ao autor (solution-architect).

## Rodada 2 (2026-09-28) — final
- Veredito: `approved` — 0 `blocking`.
- A1 aplicado: o sinal do §5 foi conferido com T-801:38-51 e SKILL.md:72-82; uma implementação correta passa e a deriva real é sinalizada. A2–A5 aplicados.
- Sem regressão: só os §3 itens 3 e 5, §4:47-49, §5:55 e §7:63/65 mudaram. Nenhum achado novo.
- Nota do validador, não levantada: §5:54 ("fora e depois" do primeiro bloco cercado da seção de revisores) é atendido automaticamente, pois coordinator.md não tem cópia da cláusula.

---

# Revisão · Emenda 2026-09-29 (T-1002, PRD-008)

## Rodada 1
- Revisor: document-validator (claude; opus solicitado, efetivo não verificado). Um validador Claude (decisão do Coordenador em nome do dono: a emenda reproduz decisões do PRD-008, já validado pelo painel completo).
- Veredito: `changes required` — P1 (`blocking`) a lista residual omitia `refs/` como arquivos ao mesmo tempo em que afirmava ser a lista do PRD-008; P2 o §7 não citava o RISK-002; P3 a redação da sonda do agy era ambígua; P4 (anexo) a proposta do RISK-002 abria com "mitigado". Reportado por: claude.
- O Coordenador aplicou os quatro.

## Rodada 2
- Veredito: `approved` — P1–P4 aplicados; nenhum problema bloqueante novo. Reportado por: claude.
