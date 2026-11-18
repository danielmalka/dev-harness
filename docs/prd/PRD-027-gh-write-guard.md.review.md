# Revisão: PRD-027

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Documento | docs/prd/PRD-027-gh-write-guard.md |
| Fonte | Decisões do dono D1-D12 (2026-10-21), descoberta DISC-007, medição M2 |
| Rodada | 1 de 2 |
| Revisor | cli:agy/gemini-3.8-flash-high |
| Veredito | `approved` |
| Atualizado | 2026-10-21 |

## Achados

(Nenhum achado em aberto.)

## Não levantado

- Lacunas: Os três defeitos do lote 0.32.1 (a trava de escrita do `gh` no Bash de subagente, a receita do opencode v2.0.32 e a pegadinha de autocasamento de `pgrep`/`pkill -f`) mapeiam, a partir da DISC-007, do chat do dono e da M2, para as descrições de problema G1-G3 e as regras R1-R14.
- Lacunas: A decisão D1 (subagentes `dh:*` do kit negados, exceto `dh:release-manager`; falha fechada para um `agentId` desconhecido; sessão principal e subagentes fora do kit não afetados) está coberta nas seções 2.1, D1, R4, R5 e R7. O escopo de limpeza do AGENTS.md ("Premissa de limpeza") está coberto na R15 e na seção 4.
- Ambiguidades: A divisão em segmentos de comando nos separadores de shell, o salto de atribuições de ambiente e de invólucros (`env`, `command`, `exec`, `sudo`) e o tratamento de argumentos estão definidos em 2.3 e R2. O contrato de função pura de `hooks/lib/` (R1) fixa entrada, saída e zero de E/S. Os seis contextos de execução do teste de motor (R5) não se sobrepõem. A falha aberta numa exceção da trava e a falha fechada num `agentId` desconhecido ficam separadas (R7).
- Conflitos: A emenda da ADR-005 (R9) resolve a exceção estreita para os verbos de escrita do `gh` e para `git push` contra a rejeição geral de análise de comando da ADR-005. A linha "Trava de `gh`" do `AGENTS.md` casa com 2.1, D1 e R5. As contagens invariantes (13 agentes, 19 + 1 comandos, 14 templates) casam com o inventário (R16). A sincronização da versão 0.32.1 (R17) e a prova de clone limpo para mudanças em Go concordam com o AGENTS.md. A parada de merge de caminho crítico (R18) casa com as entradas do `REVIEW.md` tocadas por T-2701 e T-2702.
- Excessos: Os itens fora de escopo (sandbox do SO, bloqueio de curl e de MCP, restrição das ferramentas do `qa-verifier`, `dh review run`, novos evals ou sondas, mudanças em `reviewers:`, runtimes que não sejam o Claude Code) ficam na seção 5. Os contornos conhecidos estão documentados em 2.4 e R9 sem escopo adicional.
- Critérios fracos: Todas as regras R1-R19 declaram uma checagem determinística e observável. A M1 (R8) tem condições explícitas de passa e falha sob o modo de permissão de contorno.
- Organização: O documento segue `templates/en/PRD.md`. As regras R1-R19 são numeradas em sequência. Os tickets sugeridos T-2701 a T-2704 têm lanes, conjuntos de escrita e dependências corretos (ADR-003, PRD-007).
- Estilo de redação do PRD: conforme as regras do ASD-STE100 em `.skills/requirements-discovery/SKILL.md`.

## Registro de rodadas

- Rodada 1 de 2, 2026-10-21: `approved`, 0 achados bloqueantes. Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Esclarecimento do Coordenador após a rodada 1, 2026-10-21 (não revalidado, dentro da D3): a seção 2.3 e a R2 também saltam as opções globais entre o binário e o subcomando (`gh -R`/`--repo`, `git -C`, `-c`, `--opt=value`); status definido como `approved`.
- Revalidação da emenda (rodada única), 2026-10-21: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 1 da revisão de segurança: seções 2.1-2.4, R1, R2, R3, R5, R6, R9, D1, D2, D3, D6 e os riscos do apêndice (lista de permissão de leitura para o `gh`, casamento em duas passagens, Bash/Monitor/PowerShell). Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Revalidação da emenda 2 (rodada única), 2026-10-21: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 2 da revisão de segurança: seções 2.2, 2.3, 2.4, R2, R9(a), D2, D3 e o risco de falso positivo do apêndice (duas passagens reais, checagem em cada palavra `gh`/`git` com uma exceção de busca/impressão, novas negações de `git`, cabeçalho de sobrescrita de método). Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Revalidação da emenda 3 (rodada única), 2026-10-21: `approved`, 0 achados bloqueantes (analisador léxico com pilha de contexto, limites de recursão e de tamanho, exceção de busca mais estreita, negação do PowerShell sem strings). Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo. Substituída pela emenda 4 antes da implementação.
- Revalidação da emenda 4 (rodada única), 2026-10-21: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Decisão do dono: o analisador de shell é substituído por uma regra curta (palavras de gatilho, uma leitura simples do gh ou um grep simples, sem metacaracteres de shell). Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Revalidação da emenda 5 (rodada única), 2026-10-21: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 4 da revisão de segurança: leitura do gh conferida em ordem, lista de flags de `auth status`, limites de palavra, formas `git-<verb>`, PowerShell sem exceções, regra geral de falso positivo. Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Revalidação da emenda 6 (rodada única), 2026-10-22: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 5 da revisão de segurança: texto de gatilho em minúsculas, tratamento de `\` por shell (separador de caminho do PowerShell), `git` termina antes de `-` (lista de hífens removida), `gh api` sem `--hostname` e sem tokens `://`, leitura do gh em maiúsculas como falso positivo aceito. Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Correção do Coordenador após a emenda 6, 2026-10-22 (não revalidada, só redação): o exemplo de caminho do Windows em 2.3 e na R2 usa um marcador `<drive>` no lugar de uma letra de unidade real, porque o `dh validate` rejeita um caminho com letra de unidade como caminho privado; a regra e o caso não mudam.
- Revalidação da emenda 7 (rodada única), 2026-10-22: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 6 da revisão de segurança (escolha "a" do dono: uma rodada extra): dois textos de gatilho do Bash (toda barra invertida removida, ou `\\` mapeada para `/`), disparado quando qualquer um dispara; exceção de leitura do gh só quando nenhum texto disparou em `git`; casos de crase invertida na R2; 2.4 registra a hipótese `GH_ENTERPRISE_TOKEN`. Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
- Revalidação da emenda 8 (rodada única), 2026-10-22: `approved`, 0 achados bloqueantes, nenhum achado em aberto. Emendado após a rodada 7 da revisão de segurança (escolha "a" do dono): 2.4 nomeia como contorno aceito a forma barra invertida seguida de nova linha dentro de um nome de comando; cabeçalho atualizado. Revisor: cli:agy/gemini-3.8-flash-high. Checagem do workspace: limpo.
