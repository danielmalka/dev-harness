# Revisão: PRD-024

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-024-public-readiness.md |
| Fonte | Escopo aprovado pelo dono, chat 2026-10-16 (despacho do Coordenador); auditoria UND-006; docs/discovery/renova.md §7; AGENTS.md |
| Rodada | 2 de 2 |
| Veredito | changes required |
| Atualizado | 2026-10-16 |

## Achados

- **P1** [applied] [conflict] A R8 nunca passaria porque o PRD continha o caminho do dono.
  - Resolução: rodada 2. O literal aparece só na R8 (§3:97), com o pathspec `':!docs/prd/PRD-024*'` que também cobre este arquivo de revisão. §1, §2 e §6 agora usam "o caminho da home/checkout do dono".
  - Reportado por: claude
- **P2** [applied] [conflict] A regra de contagem descrevia mal o checkDocumentedCounts, e a R7 já passava na base.
  - Resolução: rodada 2. A §2:62 permite no máximo uma linha ancorada com dois ou mais substantivos (quatro âncoras). O grep da R7 foi ampliado. A §3:85 nomeia R7, R11, R17, R18, R19 e R21 como travas de regressão.
  - Reportado por: claude
- **P3** [applied] [weak criterion] Os greps da R4 passavam na base, e a R1 só contava títulos.
  - Resolução: rodada 2. A tabela da §2 fixa 12 títulos en/pt. A R1 checa a ordem. A R4 tem escopo por seção e inclui per-hook/por hook. A R3 exige optional/opcional.
  - Reportado por: claude
- **P4** [applied] [gap] A lista de referências mortas estava incompleta, e a R14 era ambígua sobre RISKS.md.
  - Resolução: rodada 2. A §1:26-27 lista as linhas exatas nos dois arquivos (batem com a base). A §2:73 mantém 760 e 839 e reescreve as Relações do RISK-002. A R14 tem CLAUDE.md vazio e as linhas de RISKS.md nomeadas (ver P9).
  - Reportado por: claude
- **P5** [applied] [ambiguity] (`non-blocking`) Os blocos atuais do README não tinham seção de destino.
  - Resolução: rodada 2. A §2:54-60 mapeia idioma, carga pelo clone, fontes, a nota da 0.8.0, modo global e a linha 20.
  - Reportado por: claude
- **P6** [applied] [weak criterion] (`non-blocking`) Os títulos dos cards (a)-(c) não estavam fixados.
  - Resolução: rodada 2. A §2:67-70 dá os textos exatos dos `<h3>` en/pt, e a R12 os checa.
  - Reportado por: claude
- **P7** [applied] [ambiguity] (`non-blocking`) O grep de `claude` da R10 proibia nomear o Claude Code.
  - Resolução: rodada 2. O padrão agora é `co-authored-by|generated with`, e Claude Code e `.claude-plugin` são permitidos.
  - Reportado por: claude
- **P8** [aplicado em parte] [organization] (`non-blocking`) Tamanho das frases e voz passiva.
  - Resolução: rodada 2. A §1:31 e o card (d) foram divididos. A §2:62 ainda tem uma frase de cerca de 27 palavras (`Minor`).
  - Reportado por: claude
- **P9** [aceito pelo dono] [conflict] O texto do card (d) põe `RISKS.md` nos dois roadmaps, então o "exatamente 2 linhas por arquivo" da R14 falha.
  - Local: §2:70 contra §3 R14
  - Evidência: o card (d) diz que o Coordenador "continua dono das escritas em `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`" (também §6:148). O grep então imprime 760, 839 e a linha do card.
  - Sugestão: permita exatamente 3 linhas por arquivo (760, 839 e a linha do card (d)). Diga que o card Entregue da 0.30.1 não escreve `CLAUDE.md` nem `RISKS.md`.
  - Reportado por: claude
- **P10** [aceito pelo dono] [ambiguity] (`non-blocking`) As "7 linhas de tabela que começam com `|`" da R4 não dizem se o cabeçalho e o separador contam, e uma contagem literal dá 9.
  - Local: §3 R4, item Escritas
  - Sugestão: "7 linhas de corpo (depois do cabeçalho e do separador)".
  - Reportado por: claude

## Não levantado

- As linhas de RISKS.md na base por roadmap são 760, 839 e as duas linhas de Relações do RISK-002, então ficariam exatamente 2 sem o card novo. As linhas de CLAUDE.md são as listadas na §1:26.
- A linha Requisitos (Claude Code; Git) bate com README.en.md:28 na base. Mover Go para Para quem contribui é consistente com os binários pré-compilados.
- As checagens de título da R1, R5 e R6 se sustentam contra os títulos fixados ("contribut" / "contribu").
- O card (d) continua batendo com a decisão do dono, e a implementação dele fica fora de escopo (§5:138).
- Todos os itens Não levantado da rodada 1 continuam valendo: mapa de fontes, decisões do dono, fatos da auditoria e do release, fora de escopo, regras do repositório, forma.

## Disposição do dono (2026-10-16)

- Limite de validação atingido com a P9 em aberto. O dono escolheu aplicar a correção do validador para a P9 (R14: exatamente 3 linhas de `RISKS.md` por roadmap, o card Entregue da 0.30.1 não cita `CLAUDE.md` nem `RISKS.md`) e para a P10 (7 linhas de corpo), sem outra rodada.

## Emenda no meio do build, 1 de 1 (2026-10-16)

- Status da revisão: `changes required` (1 `blocking`, vai para o dono). Escopo: R8, R15, R22, §4:120, §6:151-152.
- **P11** [aceito pelo dono] [conflict] A cláusula de marketplace da R15 supõe uma contagem desatualizada que não está desatualizada.
  - Local: §3 R15, §4:120, §6:152
  - Evidência: marketplace.json:8 diz "coordinator, 12 specialist agents". O AGENTS.md diz 13 agentes contando o coordinator, então o texto é verdadeiro. A R15 proíbe a string e permite um "13 specialist agents" errado. O checkDocumentedCounts não varre o marketplace.json.
  - Sugestão: remova a cláusula, ou exija "13 agentes no total ou coordinator mais 12 especialistas".
  - Reportado por: claude
- **P12** [aceito pelo dono] [weak criterion] (`non-blocking`) As checagens da R8 não cobrem arquivos HTML, e a nota "varre os arquivos rastreados" é imprecisa.
  - Local: §3 R8
  - Evidência: validate.go:579 pula `.html` e :584 pula `evals/baselines` e `evals/results`. O git grep cobre só os baselines. Os roadmaps são editados neste PR.
  - Sugestão: `git grep -l '/ho[m]e/ma[l]ka'` sobre a árvore inteira não imprime nada.
  - Reportado por: claude
- **P13** [aceito pelo dono] [ambiguity] (`non-blocking`) A R22 pode remover o passo de desabilitar que a R4 exige, e não diz se a forma de CLI prova a forma de barra.
  - Local: §3 R22, R4
  - Sugestão: recorra a `/plugin uninstall` em Desligar, e diga que um comando é verificado na forma que o README imprime.
  - Reportado por: claude
- **P14** [aceito pelo dono] [organization] (`non-blocking`) A §6:152 não está no formato de decisão "Dono, <data>:".
  - Local: §6:152
  - Reportado por: claude

### Não levantado

- O PRD-024 e este arquivo de revisão não contêm nenhuma substring de caminho privado (o padrão do `checkPrivatePaths`), então o checkPrivatePaths passa neles.
- A numeração da R22 é única (R1-R21 são contíguas). A R22 é checável pelo relatório de QA.
- A R8 mantém o laço de parse e a checagem de 28 arquivos `<repo>`. A exclusão por pathspec foi removida.
- As correções de P9 e P10 estão na R14 e na R4. Nada fora do escopo declarado da emenda mudou.

### Disposição do dono (2026-10-16)

- Rodada única terminada com a P11 em aberto. O dono escolheu remover a cláusula da descrição do marketplace (P11) e aplicar as correções do validador para P12, P13 e P14, sem outra rodada.
