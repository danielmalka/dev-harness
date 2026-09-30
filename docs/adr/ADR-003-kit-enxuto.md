# ADR-003 · Kit enxuto: menos verificação, menos agentes, `AGENTS.md` só com regra vigente

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-09-30 |
| Decisor | Daniel Malka (dono); R11 decidido pelo Coordenador por delegação do dono |
| Origem | PRD-010 (R1, R9, R11, R12); decisão "Kit enxuto" do `AGENTS.md` de 2026-09-30 |
| Reversibilidade | média: os agentes fundidos e as regras enxutas voltam por edição de fonte, mas `dist/`, evals e tutoriais dependem dos nomes novos |

## 1. Contexto

- Problema: o kit cresceu por acréscimo e a verificação virou o custo dominante. Dito pelo dono: "ficamos dias a dias rodando testes e testes que não sei o quanto foram produtivos" e "não quero que vc tire coisas, talvez um merge conforme sugerido já resolva bastante coisa".
- Fatos (PRD-010 §1, medidos em 2026-09-30): 23 casos de eval, 19 deles com o corpo de um agente copiado em `append_system_prompt` e uma checagem do `dh validate` que reprovava qualquer deriva; 19 agentes, cinco construtores com o mesmo esqueleto, dois revisores e dois de desenho; `AGENTS.md` e o Coordenador misturando regra vigente com data, motivo e histórico; prova de clone limpo exigida em todo PR; sem guia completo de uso.
- Restrição: o kit não perde capacidade. Os 19 comandos continuam, e cada capacidade dos agentes fundidos continua alcançável por comando ou despacho.

## 2. Decisão

Escolhas do dono sobre a auditoria de simplificação (numeração da auditoria):

| # | Proposta | Escolha |
|---|---|---|
| 1 | Evals fora do fluxo normal: poucos casos essenciais, sem checagem de cópia do Coordenador nem `resync` | sim |
| 2 | Tirar o `/dh:plan-loop` (proposta da auditoria, não adotada; o comando fica) | não |
| 3 | Tirar `improve`, `verify` e `consolidate-memory` (proposta da auditoria, não adotada; os comandos ficam) | não |
| 4 | Fundir agentes | sim |
| 5 | Fundir o que se sobrepõe por trás dos comandos, sem remover comando | sim |
| 6 | `AGENTS.md` e Coordenador só com regra vigente | sim |
| 7 | Tirar as CLIs externas e o isolamento por clone (proposta da auditoria, não adotada; pedido do próprio dono) | não |
| 8 | Tirar o validador de documento (proposta da auditoria, não adotada; perderia qualidade) | não |
| 9 | Manter docs só em pt-br (proposta da auditoria, não adotada; pt-br e en continuam) | não |
| 10 | Prova de clone limpo só quando o PR mexe em Go ou binários | sim |

Efeitos:

- Evals: 5 casos essenciais em `evals/cases/`, os demais em `evals/archive/`; rodam só a pedido. A checagem de cópia de corpo de agente sai.
- Agentes, de 19 para 13: `api-designer` e `solution-architect` viram `architect`; `backend-builder`, `frontend-builder`, `data-engineer`, `devops-engineer` e `refactorer` viram `builder` (com lane ou modo); `code-reviewer` e `security-reviewer` viram `reviewer` (modos code e security). `debugger` continua separado.
- Verificação padrão (R11, 2026-09-30): checagens determinísticas mais 1 revisor Claude; painel de 3 só a pedido do dono ou em segurança; sem eval nem sonda nova sem pedido.
- Clone limpo (R12): só quando o PR toca `cmd/`, `internal/`, `go.mod` ou `dist/*/bin`.
- `AGENTS.md`: no máximo 18 regras vigentes, sem data nem motivo; a história está na seção abaixo.
- Tutorial completo em `docs/tutorial.html` e `docs/en/tutorial.html`; as quatro páginas antigas viram ponteiros.

Espelho em inglês: os ADRs anteriores (ADR-001, ADR-002) existem só em pt-br em `docs/adr/`, sem `docs/en/adr/`. Este ADR segue o precedente e não tem espelho.

## 3. Histórico das regras

Cada decisão do `AGENTS.md` anterior (versão de `HEAD` antes deste lote) mais as três decisões que estavam só na árvore de trabalho. "Bullet n" é a ordem em `AGENTS.md` (seção "Regras vigentes"). "Só histórico" significa que a regra saiu do arquivo e vive aqui.

| # | Regra antiga | Data | Motivo | Onde vive agora |
|---|---|---|---|---|
| 1 | Idioma de `.skills/`, `.agents/`, `.commands/`, `profiles/` (en), `templates/` espelhado, `docs/` pt-br com espelho en, READMEs | anterior a 09/2026 | Kit distribuído em duas línguas | Bullet 1 |
| 2 | Memória `.harness/MEMORY.md`, `EPOCHAL.md`, `RISKS.md`: só o Coordenador escreve | anterior a 09/2026 | Um único escritor evita registros conflitantes | Bullet 3 |
| 3 | Conflito de regras que muda a forma da entrega: perguntar | anterior a 09/2026 | Não assumir resolução | Bullet 4 |
| 4 | Licença MIT, `plugin.json` repete | anterior a 09/2026 | Redistribuição do kit | Bullet 5 |
| 5 | Fontes canônicas; `dist/` gerado, não se edita | anterior a 09/2026 | Um só lugar de edição | Bullet 6 |
| 6 | Idioma em execução (`language` no `project.yaml`, despachos em inglês) | anterior a 09/2026 | Dono escolhe a língua dos registros | Bullet 2 |
| 7 | O disco é a verdade sobre `.harness/` | anterior a 09/2026 | Não ressuscitar registros apagados pelo dono | Bullet 3 |
| 8 | Extensão VS Code: só visualizador e lançador, não é fork da `context-kit-extension` (arquivada como doadora por cópia) | anterior a 09/2026 | Extensão nunca vira runtime | Bullet 8 (o detalhe da doadora fica só aqui) |
| 9 | Ordem: etapa 1 no plugin (marketplace, statusLine, hooks, snapshots) antes da extensão; painel v1 ao vivo e só Claude Code; OTEL e Agent SDK adiados; validada no danlemos | anterior a 09/2026 | Extensão precisa dos snapshots do plugin | Bullet 8 (a validação no danlemos fica só aqui) |
| 10 | Runtime Go, binário `dh`, binários commitados por plataforma (ADR-001); migração concluída em 20/09/2026, `scripts/` removido | 20/09/2026 | ADR-001 | Bullet 7 (a data da migração fica só aqui) |
| 11 | Etapa 1 especificada em `docs/prd/PRD-001-etapa-1.md`; código só depois da spec aprovada | anterior a 09/2026 | Spec antes de código | Bullet 8 (referência ao PRD-001); o "só depois da spec" já é a regra zero do kit |
| 12 | Premissa: documento que vira contrato passa por revisor adversarial somente leitura; ciclo até `approved`; teto 2, elevado a 6 na 0.8.0; implementada na 0.3.0 (`document-validator`, skill `document-review`); avaliação dos casos `doc-validator-*` pendente | 20/09/2026 (teto: 25/09/2026) | PRD-002 | Bullet 9 (validador, teto de documento em 2 pela revisão enxuta); o pendente de avaliação é só histórico |
| 13 | Revisão em pares por CLI externa: `reviewers:` por etapa, `local.yaml`, skill `external-clis` em prosa (revisitar em 21/10/2026), revisor CLI somente leitura e veredito descartado se a árvore muda | 21/09/2026 | PRD-003 | Bullet 12 |
| 14 | Documentos HTML internos via `doc-template-html`: seções nomeadas pelo dono prevalecem; tipo mais próximo no hero com nota de desvio | 21/09/2026 | Dono nomeia as seções (`docs/roadmap.html`) | Bullet 15 |
| 15 | Caminhos de máquina sempre relativos ao home ou ao projeto | 21/09/2026 | Mesmo arquivo em qualquer máquina | Bullet 12 |
| 16 | Prova de clone limpo com sha256 antes de commitar `dist/` regenerado; esclarecimento de 24/09/2026 (commit local numa branch de feature é candidato; amend só antes do push, nunca `--force`) | 20/09/2026; 24/09/2026 | A 0.2.0 reprovou na CI por carimbo de VCS nos binários | Bullet 11, agora só quando o PR toca `cmd/`, `internal/`, `go.mod` ou `dist/*/bin` (R12) |
| 17 | CI de evals sem `ANTHROPIC_API_KEY`: o job `evals` pula, o "pass" não mede nada; medições rodam localmente; não propor a chave | 25/09/2026 | Custo | Bullet 10 |
| 18 | Diagramas Excalidraw em `.excalidraw` separado, só link; Mermaid via CDN é a única exceção; sem implementação ainda (item "Skills de diagrama" do roadmap) | 25/09/2026 | Manter docs sem embutir script | Bullet 15 |
| 19 | Teto padrão de 6 rodadas em todo loop (plano, documento, código, build/QA, plan-loop) no lugar de 2; implementado na 0.8.0 | 25/09/2026 | "não acho que qualquer loop de iteração se resolva facilmente em 2 ou 3 rodadas, pois os revisores são bem críticos" | Bullet 9 (6 para build/QA e código; documento em 2 pela revisão enxuta) |
| 20 | `cli:claude` fora do planejamento em loop enquanto o kit for plugin do Claude Code; rever ao criar adaptador não-Claude (PRD-005 AA1) | 25/09/2026 | Tudo já roda sobre o Claude | Bullet 13 |
| 21 | Documentação dentro do PR, nunca em PR próprio; docs-guide ao fim de cada lote (os PRs #9, #13 e #14 violaram) | 25/09/2026 | Nada fica defasado | Bullet 15 |
| 22 | Proporcionalidade de revisão e medição: (1) mudança mecânica só com checagem determinística e 1 revisor, painel de 3 para contrato novo, lógica, segurança ou texto de agente; (2) eval paga só quando o lote muda o comportamento medido; (3) teto de 6 é limite, não meta; implementada na 0.9.0 | 25/09/2026 | "não precisamos de revisão do revisor a não ser que o trabalho a ser medido seja o trabalho do revisor" | Só histórico: substituída pela revisão enxuta e pelo R11 (bullets 9 e 10) |
| 23 | Revisão enxuta: documento com 1 validador e até 2 rodadas; código com painel de 3 só para contrato novo, lógica ou segurança; eval com `--runs 1`; teto de 6 para build/QA e código; implementada na 0.9.0 (PLAN-010) | 26/09/2026 | O PRD-006 levou 6 rodadas: "estamos avaliando demais coisas que poderiam ser mais simples" | Bullet 9; o painel de 3 "para contrato novo ou lógica" foi restringido pelo R11 a pedido do dono ou segurança |
| 24 | CLI planejadora só com sandbox (só `codex` e `mcode`; `grok`, `agy`, `opencode` fora até haver `git worktree` descartável) | 28/09/2026 | Congelamento na árvore viva só detecta, não isola (T-801) | Só histórico: SUPERADA pelo PRD-008 (clone descartável); ver linha 30 e o ADR-002 |
| 25 | Exceção do `/dh:plan-loop`: cada onda julgada pela lista completa de `reviewers.document`; fallback sem CLI externa segue o `/dh:plan` (2 rodadas) | 28/09/2026 | PRD-005 exige aprovação unânime com pelo menos um revisor não-Claude; com um só revisor o loop nunca converge (PLAN-011) | Bullet 9 |
| 26 | Novo fluxo de PRD e tickets: (1a) formato novo substitui PRD, STORY e TASK; story no ticket; (2a) implementação com `qa-verifier` e 1 revisor Claude, teste só com `qa-verifier`, revisão final por integração e aderência; (3a) lanes; (4a) tela pelo `qa-verifier` com Playwright; rejeição de PRD pela revisão enxuta; emenda no meio da execução volta ao validador por uma rodada; decisões do §8 do PRD-007 (apêndice opcional, ticket mantém "Resultado" e "Riscos e limites", cabeçalho do PRD mantido, sem segundo revisor em ticket de teste) | 29/09/2026 | PRD-007, `docs/novo-fluxo-prd.md` | Bullet 16; a parte "painel de 3 em contrato novo ou lógica" da 2a ficou substituída pelo R11; detalhes do §8 vivem no PRD-007 |
| 27 | Roteamento neste repositório: ticket só de fonte do kit usa lane `backend` e vai ao `harness-maintainer` com `harness-authoring` | 29/09/2026 | Não criar lane nova no kit genérico | Bullet 16 |
| 28 | Ticket de bug: sem seção "Defeito"; reprodução e observado × esperado no Fluxo; regressão verde no "Pronto quando"; AC-07 literal | 29/09/2026 | PLAN do PRD-007 | Bullet 16 |
| 29 | `AGENTS.md` é o arquivo canônico; `CLAUDE.md` só importa; no kit, `AGENTS.md` vence o `CLAUDE.md`; referências antigas a `CLAUDE.md` em PRDs, ADRs e relatórios são histórico | 29/09/2026 | Padrão entre ferramentas de IA (PRD-008) | Bullet 18 (a nota sobre referências antigas fica só aqui) |
| 30 | Isolamento de CLI sem modo somente-leitura por clone descartável (`git clone --no-hardlinks`), não `git worktree`; RISK-002 "mitigado"; sandbox do SO adiada; `codex` no sandbox próprio, `mcode` no clone até ter sonda, `agy` só com a sonda de cwd | 29/09/2026 | O worktree compartilha o `.git` inteiro (PRD-008; ver ADR-002) | Bullets 13 e 14 |
| 31 | Riscos residuais do clone aceitos até o sandbox do SO: (a) escrita fora do projeto, (b) processos que sobrevivem, rede e push para outros remotos, (e) leitura e devolução de segredos; detecção da árvore viva sem config global nem fsmonitor, cobrindo o índice | 29/09/2026 | Aceite explícito do dono no PRD-008 | Bullet 14 (resumo); a nota sobre git sem config global é só histórico |
| 32 | PRD-008 aprovado; aceitos os caminhos do `.git` vivo fora da foto do R5 (`logs/`, objetos e pacotes, `worktrees/`, `modules/`, `FETCH_HEAD`, `ORIG_HEAD`, `COMMIT_EDITMSG`, bytes crus do índice); R5 hasheia todos os caminhos ignorados, com limite de custo só se uma medição mostrar lentidão | 29/09/2026 | Aprovação do PRD-008 | Bullet 14 (resumido em "caminhos do `.git` vivo fora da foto do R5"); a lista completa e a nota do limite de custo são só histórico |
| 33 | Roadmap vivo, opção A: escrito à mão nas duas línguas, `dh validate` confere data, versão, fontes, um card por release na ordem certa e paridade pt/en contra o `CHANGELOG.md`; atualizar o roadmap entra no checklist de release e no despacho de docs; geração a partir de dados (B ou C) só em PRD futuro | 30/09/2026 | Discovery do PRD-009 | Bullet 15 |
| 34 | PRD-009 aprovado: checagem pulada só sem os dois roadmaps; linha de manutenção condicional no `delivery-readiness` e no despacho de docs do Coordenador | 30/09/2026 | Aprovação do PRD-009; entregue na 0.13.0 (PR #21) | Bullet 15 (detalhes no PRD-009) |
| 35 | Kit enxuto: o kit não perde capacidade; entram evals fora do fluxo, fusão de agentes, fusão por trás dos comandos, `AGENTS.md` e Coordenador só com regra vigente, clone limpo só com Go ou binários | 30/09/2026 | Simplificar também a verificação; "o kit não deve virar um monstro" | Bullets 9 a 11, 17 e 18; esta decisão em si é o §2 deste ADR |

Renomeações de agentes (para ler PRDs, ADRs e CHANGELOG antigos): `api-designer` e `solution-architect` viraram `architect`; `backend-builder`, `frontend-builder`, `data-engineer`, `devops-engineer` e `refactorer` viraram `builder`; `code-reviewer` e `security-reviewer` viraram `reviewer`.
