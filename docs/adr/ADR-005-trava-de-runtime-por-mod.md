# ADR-005 · Trava de runtime por mod do plugin

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-05 |
| Decisor | Daniel Lemos (dono) |
| Origem | PRD-011; RISK-001 (subagente ignorou `disallowedTools` e despachou outro agente) |
| Reversibilidade | barata, mas ver ADR-006: o `runtime-guard.ts` é encadeado por `panel.ts`, então reverter exige tirar o import e o encadeamento de lá |
| Alterado por | ADR-006 (composição em `panel.ts`; `runtime-guard.ts` não é mais declarado em `modules`); ADR-007 (a trava casa também `<home>/projects/*/`); PRD-027 (exceção estreita para `gh` de escrita e `git push`, ver Emenda de 2026-10-21) |

## 1. Contexto

- Pergunta: como impedir que um subagente despache agentes ou escreva `MEMORY.md`, `EPOCHAL.md` e `RISKS.md`, se até aqui a regra só existe em prosa.
- Fato: `disallowedTools` no frontmatter não foi imposto pelo runtime (RISK-001, 2026-09-21).
- Fato: o Claude Code oferece mods em plugins (hooks de função em TypeScript, API de acesso antecipado) que interceptam `agent.spawn` e `tool.call`.

## 2. Decisão

- Escolhida: um mod no plugin (`hooks/runtime-guard.ts`, declarado em `modules` no `hooks.json`) com duas travas. `agent.spawn` é negado quando há `parentAgentId` (só o Coordenador despacha). `tool.call` de `Write` e `Edit` é negado quando há `agentId` e o alvo termina em `.harness/MEMORY.md`, `EPOCHAL.md` ou `RISKS.md`.
- Falha aberta: erro no mod libera a chamada; a prosa dos agentes continua como segunda camada.
- Os hooks clássicos (`dh snapshot event`) ficam intactos.
- Descartada: voltar a `disallowedTools` (não imposto) e travar `Bash` por análise de comando (frágil, falso positivo).

## 3. Consequências

- Só vale no Claude Code. Em outros runtimes continua a prosa.
- Lacuna conhecida: escrita por `Bash` (redirect, `sed -i`) não é bloqueada.
- A API de mods é de acesso antecipado e pode mudar entre versões do Claude Code: rodar `claude plugin test` a cada release.
- Primeiro artefato de runtime que não é Go. Ele roda no motor do Claude Code, não como script nosso; o ADR-001 continua valendo para os scripts do kit (`dh`).

## 4. Revisar quando

- A API de mods mudar ou sair do acesso antecipado.
- Aparecer escrita em `.harness/` por `Bash` de subagente (então avaliar travar `tool.call` de `Bash`).
- O kit ganhar adaptador para outro runtime.

## Emenda de 2026-10-21 (PRD-027)

- Exceção estreita à frase "travar `Bash` por análise de comando" da seção 2: o mod nega `tool.call` de `Bash`, `Monitor` e `PowerShell` (qualquer ferramenta que receba um comando de shell) quando há `agentId`, o tipo do subagente tem prefixo `dh:` e não é `dh:release-manager` (ou não consta do mapa de tipos: falha fechada), e o comando cita `gh`, ou `git` com uma palavra de escrita (`push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`; o `git` também casa antes de `-`, ver a emenda de 2026-10-22). Passa só (a) uma leitura `gh` simples da lista do PRD-027 §2.2, sem metacaractere de shell, ou (b) um `grep`, `egrep` ou `fgrep` simples, sem `; | & $ ( ) < >`, crase ou quebra de linha. Todo o resto é negado, e a mensagem manda usar a ferramenta Grep. Não há parser: nenhuma divisão, recursão ou limite de tamanho. Sessão principal, `dh:release-manager` e subagentes de outros plugins não são afetados. Erro dentro da trava libera a chamada. Em `PowerShell`, o texto de busca também perde o acento grave, e todo comando disparado é negado.
- Falsos positivos aceitos (PRD-027 §2.3): comando composto que cita `gh` (`cd x && gh pr view 1`; rode a leitura sozinha), leitura `gh` com aspas ou outros caracteres especiais (`gh pr view 1 --jq '.x'`), `git` que cita uma palavra de escrita (`git log --grep=push`), busca que não é grep e cita `gh` ou palavra de escrita (`rg 'git push'`), caminho com parte `gh` (`ls ~/.config/gh`), e, em geral, todo comando disparado fora dos dois casos, mesmo inofensivo (`echo gh`), inclusive todo `gh` em `PowerShell`.
- Motivo: em `/dh:review pr <n>` o `qa-verifier` lê texto de PR de terceiro, tem `Bash` e a credencial do `gh`, e a sessão roda em modo bypass; a proibição só em prosa não bastava (S3 residual do PRD-026).
- Melhor esforço, por análise de texto. Contornos conhecidos, aceitos e não corrigidos (PRD-027 §2.4): apelidos e funções de shell definidos antes; nome de comando montado em tempo de execução (variáveis, `eval`, truques de string que não deixam a palavra `gh` ou `git`); script gravado em disco e executado, inclusive por ferramenta de build (`make push`); outros interpretadores (`python -c`, `node -e`) e texto enviado a um deles; alias do git gravado na configuração por uma chamada anterior; `curl` com token; ferramentas MCP do GitHub; subagente que não é do kit; processos filhos de uma CLI externa. A prosa dos agentes continua como segunda camada. Sandbox de SO fica fora do escopo.
- Emenda de 2026-10-22 (PRD-027, emendas 6 a 8, após as revisões de segurança 5 a 7): o gatilho é avaliado sobre o texto do comando em minúsculas. Em `Bash`, o gatilho vale se um de dois textos disparar: A tira todo `\`, `'` e `"`; B troca cada par `\\` por `/` e tira o resto. A exceção de leitura `gh` só vale se nenhum texto disparou em `git`. Hipótese não verificada: `GH_ENTERPRISE_TOKEN` pode ir a outro host por `-R HOST/OWNER/REPO`, argumento `HOST/OWNER/REPO` ou URL completa nos grupos de leitura fora de `gh api`. Em `PowerShell`, todo `\` vira `/` e `'`, `"` e acento grave saem. A palavra `git` termina antes de `-`, então `git-<nome>` com palavra de escrita dispara (`git-lfs push`), e saem da lista as formas com traço (`git-push`, `git-send-pack`, `git-subtree`, `git-p4`, `git-svn`). A palavra `gh` mantém a fronteira: `gh-pages` não dispara. `gh api` perde `--hostname`, e um `gh api` com algum token que contém `://` não é leitura. Falso positivo aceito: leitura `gh` com o nome do binário em maiúscula (`GH pr view 1`) é negada. Desvio aceito (emenda 8, decisão do dono): barra invertida seguida de quebra de linha dentro do nome do comando passa sem disparar.
- Fonte: PRD-027 (`docs/prd/PRD-027-gh-write-guard.md`). Os demais itens da decisão (Write/Edit em `.harness/`, `agent.spawn`) não mudam.
