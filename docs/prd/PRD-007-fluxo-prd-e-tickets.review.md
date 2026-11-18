# Revisão · PRD-007 · rodada 1

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

- Revisor: document-validator (claude, modelo opus solicitado; modelo efetivo não verificado)
- Data: 2026-09-29
- Fonte: docs/novo-fluxo-prd.md; CLAUDE.md (decisões do dono 1a–4a de 2026-09-29, revisão enxuta, docs dentro do PR); texto atual do kit; docs/prd/PRD-006-comando-autonomo.md

## Veredito
- `changes required`
- Achados bloqueantes: 8

## Achados
- P1 · `conflict` · `blocking` · §4 RF-16 / AC-32 contra o RF-03. A salvaguarda 4 da proposta dá as emendas de regra ao criador do PRD; o RF-03 concorda, o RF-16 dá ao Coordenador a emenda no meio da execução. Sugestão: o Coordenador redespacha o `product-discovery` para emendar (Docs ou Regras), depois o `document-validator` roda uma rodada de revalidação. Reportado por: claude
- P2 · `gap` · `blocking` · §4 RF-05 / §5 AC-07. O ticket da proposta tem um bloco "Pronto quando" (a salvaguarda 20 depende dele); o RF-05/AC-07 o omitem, enquanto o RF-08/AC-14 o referenciam. Sugestão: acrescentar Pronto quando com o conteúdo mínimo da proposta (o fluxo acontece, cada R<n> copiado tem evidência, os testes relevantes estão verdes, os docs listados estão atualizados). Reportado por: claude
- P3 · `conflict` · `blocking` · §4 RF-05 / §6 Dados. A decisão 1a diz que o novo formato substitui PRD/STORY/TASK; o RF-05 mantém toda seção atual do TASK e acrescenta campos do STORY, deixando duas fontes de aceite (Fluxo contra os AC-xx da §7). Sugestão: o RF-05 nomeia a lista fechada de seções do ticket, declara quais seções atuais do TASK/STORY são removidas e faz de Fluxo + Pronto quando a única fonte de aceite; manter seções técnicas além da 1a vai para a §8 como decisão do dono. Reportado por: claude
- P4 · `conflict` · `blocking` · §4 RF-17, §5 AC-33, §8 segundo risco. (a) Citação errada: a premissa de code review por fatia mora no RF-09, no RF-12 e no AC-14 do PRD-006 e na linha 19 do `.commands/auto.md`, não no RF-01 do PRD-006. (b) Editar um contrato aprovado muda a forma da entrega (CLAUDE.md "perguntar antes de escolher"; 1a "PRDs antigos ficam como estão"). Sugestão: corrigir a citação; mover para a §8 como decisão pendente do dono: (1) emendar o PRD-006, ou (2) deixar o PRD-006 intocado, com uma nota de substituição no PRD-007 mais a edição do auto.md; recomendar uma. Reportado por: claude
- P5 · `conflict` · `blocking` · §4 RF-12 / AC-23. O RF-12 do PRD-006 e o auto.md rodam o `/dh:secure` por fatia; o PRD-007 o coloca na revisão final alegando "do mesmo jeito". Sugestão: declarar gatilho = `.commands/secure.md`, posição = uma vez na revisão final da feature (ou por ticket), e listar a mudança no RF-17. Reportado por: claude
- P6 · `conflict` · `blocking` · §4 RF-10 / AC-19; §7 linha "Nova lane kit/harness". Rotear o trabalho no código-fonte do kit ao `harness-maintainer` sob a lane `backend` quebra a regra lane→especialista da proposta, é uma decisão do autor apresentada como requisito e põe checagens específicas do kit no fluxo genérico. Sugestão: mover para a §8 como decisão do dono: (a) sobreposição do builder restrita a este repositório; (b) nota de roteamento só no CLAUDE.md do dev-harness. Manter o AC-18 como a única regra lane→builder enquanto isso. Reportado por: claude
- P7 · `conflict` · `blocking` · §4 RF-12, §6 Stack, AC-24. A decisão 4a nomeia Playwright headless; o PRD rebaixa para "qualquer navegador headless", e a ui-verification marca linhas como `not-run` sem navegador, então a salvaguarda 18 se degrada em silêncio. Sugestão: declarar Playwright headless e que a ausência dele torna o veredito da feature `partial`/`blocked`, nunca `approved`; ou registrar o agnosticismo de ferramenta como decisão do dono. Reportado por: claude
- P8 · `gap` · `blocking` · §6 Contratos; afirmações do RF-05/RF-06. Arquivos omitidos: `.agents/coordinator.md` (Procedimento 7 QA e depois code review; Fluxos Feature; menção ao template STORY), `.commands/document.md` (STORY), o frontmatter de `.commands/discover.md` `writes: .harness/stories/`, `.skills/implementation-planning/SKILL.md` (o passo 5 dobra os testes na fatia; o passo 10 não checa a cobertura R<n>→teste), `templates/*/PRD.md` em Contratos. Sugestão: listá-los com o trecho que cada um deve mudar; reformular o RF-06 como uma nova passada de autorrevisão no passo 10. Reportado por: claude
- P9 · `excess` · `non-blocking` · RF-01 / AC-01. Manter a tabela de cabeçalho é uma decisão do autor (fundamentada no RF-02 do PRD-006); listar na §8 para o dono confirmar. Reportado por: claude
- P10 · `conflict` · `non-blocking` · RF-14. A regra do kit é "só o Coordenador grava os três registros", não "o Coordenador grava só eles". Sugestão: "O Coordenador grava a seção Resumo executado no PRD; nenhum outro trecho do PRD é escrito pelo Coordenador." Reportado por: claude
- P11 · `ambiguity` · `non-blocking` · RF-14 / AC-29. Não há marcador para "PRD entregue"; nomeie-o (por exemplo `Status` → `entregue em <data>` só quando nenhuma regra estiver Não honrada sem decisão do dono). Reportado por: claude
- P12 · `excess` · `non-blocking` · RF-01: a realocação de alternativas/riscos para os tickets e o RISKS.md não tem rastro; retirá-la e deixar a §8 carregar a decisão do apêndice. Reportado por: claude
- P13 · `organization` · `non-blocking` · §4 "Não entra" item 5: `"named"` → `"nomeia"`. Reportado por: claude

## Não levantado
- Decisões 1a (sem migração), 2a, 3a, 4a (só o product-discovery julga), revisão enxuta para rejeição, revalidação de uma rodada, limite de 6 rodadas por ticket: rastreadas.
- Salvaguarda 9 resolvida adiando o fechamento (RF-07, AC-11/12): fundamentada; usa o status existente `em revisão`; não quebra o passo 6 do implementation-planning.
- Salvaguardas 1–20 mapeadas; 1, 6, 16 cobertas por regras existentes do kit.
- A checagem determinística de cópia literal do RF-15/AC-30 está corretamente como sugestão; a afirmação de escopo do validate.go confere.
- Decomponível em tickets quando P2, P3, P6 estiverem resolvidos; ids RF/AC únicos e com referência cruzada.

## Evidência
- Lido: PRD-007; docs/novo-fluxo-prd.md (espelho em inglês não lido); CLAUDE.md; .agents/coordinator.md; .agents/product-discovery.md (início); .commands/build.md, discover.md, plan.md, review.md, document.md, auto.md (grep); PRD-006 linhas 30–49, AC-03, AC-14; templates/pt-br/TASK.md; .skills/ui-verification/SKILL.md (grep); .skills/implementation-planning/SKILL.md passos 1–11; internal/kit/validate.go checkTemplateRefs.
- Não verificado: se o templatePattern do validate.go casa o "STORY.md" solto em document.md; corpos de qa-verifier.md e harness-maintainer.md.
- Nenhum script rodado, nenhum arquivo editado.

---

# Revisão · PRD-007 · rodada 2 (final)

- Revisor: document-validator (claude, opus solicitado; modelo efetivo não verificado). Data: 2026-09-29.

## Veredito
- `changes required`
- Achados bloqueantes: 2 (ambos novos; P1–P13 da rodada 1 todos aplicados)
- Limite de rodadas atingido (2 de 2): o documento vai ao dono com P14 e P15 em aberto.

## Assentamento da rodada 1
- P1–P13: aplicados. P5 julgado consistente: o /dh:secure continua por ticket depois do gate do qa-verifier, em linha com o RF-12 do PRD-006 e a linha 19 do `.commands/auto.md`. Observação, não um achado: um problema de segurança que só existe quando os tickets se integram não é visto por nenhum revisor de segurança.

## Novos achados
- P14 · `conflict` · `blocking` · §5 AC-07 "contendo exatamente a lista fechada" pressupõe a resposta da decisão aberta da §8 sobre as seções técnicas do TASK; a linha 21 do `.commands/build.md` precisa de um lugar para status/resultados. Sugestão: lista fechada mais só as seções técnicas que o dono mantiver na §8; nenhuma seção de aceite separada; Fluxo + Pronto quando a única fonte de aceite. Reportado por: claude
- P15 · `ambiguity` · `blocking` · §6 Contratos, item de `.agents/coordinator.md`: "fluxo de slice legado" não está definido. Leitura A: a revisão de código por unidade sobrevive só fora de uma feature com tickets (/dh:fix, /dh:refactor, mudanças de instrução do harness), e o build.md a perde por inteiro. Leitura B: o build.md mantém uma branch por fatia para planos anteriores ao PRD-007, contradizendo o AC-15. Sugestão: declarar a Leitura A explicitamente (ou emendar o AC-15 para B). Reportado por: claude
- P16 · `organization` · `non-blocking` · A citação do RF-17 atribuída ao RF-09 do PRD-006 é na verdade o AC-14 (linha 77). Reportado por: claude
- P17 · `gap` · `non-blocking` · Nenhum arquivo da §6 carrega a regra da passada de cliente com Playwright; a §9 ainda chama a ui-verification de agnóstica de ferramenta. Sugestão: listar `.commands/review.md` ou `.skills/ui-verification/SKILL.md` com "Playwright headless; ausência → not-run, veredito parcial/bloqueado". Reportado por: claude

## Evidência
- Lido: PRD-007 corrigido; relatório da rodada 1; grep de `.commands/` por TASK/BUG-/2B; RF-09, RF-12, AC-14 do PRD-006. Não verificado: corpos de qa-verifier.md e harness-maintainer.md; plan-loop.md julgado coberto pela edição do implementation-planning. Nenhum script, edição ou despacho.
