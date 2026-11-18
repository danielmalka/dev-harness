# Rodar uma chamada de CLI num clone descartável

Este arquivo é o ciclo de vida exato de uma chamada de CLI que roda num clone
descartável: criar, popular, verificar, usar, remover, varrer órfãos. Ele vale
para uma chamada de revisão ou de planejamento a um binário sem modo somente
leitura documentado e medido. `codex` mantém o próprio sandbox e nunca o usa.

Toda menção a este mecanismo traz uma qualificação, dita por inteiro:

clone descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas; escritas fora do projeto não são detectadas

O clone é um
diretório de trabalho, não um sandbox. A CLI roda como o mesmo usuário do sistema operacional e ainda
alcança a árvore viva e o resto da máquina por caminho absoluto. Nada
aqui impede uma escrita; ele só põe o diretório de trabalho da chamada em outro
lugar e deixa o delegado jogar esse diretório fora.

Este arquivo define só o ciclo de vida do clone. Quais binários podem planejar é decidido
em "A chamada de planejamento" de `SKILL.md`, não aqui.

## Termos

`<project>` é a raiz do projeto (a árvore viva). `$CLONE` é o caminho exato do
clone. Variáveis de shell não sobrevivem entre chamadas de ferramenta separadas, então a
chamada guarda o caminho literal que o `mktemp` imprimiu, e todo comando
posterior (verificação, comparação, remoção) usa essa string literal, não `$CLONE`.
`<call-id>` é escolhido uma vez pelo Coordenador antes do passo 1, como
`<task id>-<binary>-<n>` com `<n>` um contador por lote, e pode conter só
`[A-Za-z0-9_-]`; qualquer outra coisa é `isolation_unavailable`. Ele é usado literalmente
nos dois nomes de uma chamada: `dh-cli-<call-id>-XXXXXXXX` (o clone) e
`dh-freeze-<call-id>-XXXXXXXX` (o diretório de congelamento, que também guarda o
arquivo de prompt da chamada). Chamadas que se sobrepõem no tempo compartilham um diretório de congelamento,
com o nome da primeira chamada, com um clone e um arquivo de prompt para cada uma.

`<tmp root>` é `(cd "${TMPDIR:-/tmp}" && pwd -P)`, resolvido uma vez quando a chamada (ou
o lote) começa e guardado como caminho literal. Todo template de `mktemp`, toda
guarda de remoção e a varredura de órfãos usam `<tmp root>`, nunca a variável
não resolvida.

`$PROMPT` é o caminho literal do arquivo de prompt da chamada, escrito no comando como essa string literal, como `$CLONE`. Uma chamada em clone o cria com `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )` dentro do diretório de congelamento da chamada ([live-tree-freeze.md](live-tree-freeze.md), seção 5); ele é único por chamada e é removido junto com esse diretório. Uma chamada sem clone (`codex`, a revisão do `agy`) o cria com `( umask 077; mktemp "<tmp root>/dh-prompt-XXXXXXXX" )` e o remove na própria limpeza (uma chamada interrompida pode deixá-lo para trás, modo 0600; ele não é varrido e fica para o dono).

## Requisito do projeto

O projeto deve ter `.harness/` e `.claude/settings*.json` ignorados pelo git. Se
não tiver, a árvore viva os guarda como arquivos rastreados ou não rastreados, o clone
os levaria, as checagens de ausência do passo 2 falhariam e toda chamada terminaria como
`isolation_unavailable`. Isso nunca é consertado apagando-os do clone.
Os worktrees de desenvolvimento do kit ficam em `<project>/.worktrees/`, que deve
estar ignorado pelo git também (`/dh:setup` acrescenta a linha): senão os arquivos deles
ficam não rastreados, o clone os copiaria e a comparação os julgaria.

## 1. Criar o clone

Antes de criar qualquer coisa, o projeto deve passar por um de dois ramos, escolhido pelo
que `<project>/.git` é:

- (a) `<project>/.git` é um diretório: `test -d "<project>/.git" && test ! -L "<project>/.git"`
  deve passar, e nenhuma outra checagem se aplica.
- (b) `<project>/.git` não é um diretório (um worktree vinculado tem um `gitfile`):
  `dh harness-path --json "<project>"` deve ter sucesso com um
  `linked_worktree` não nulo e, no mesmo shell, logo antes do clone,
  `test -f "<project>/.git" && test ! -L "<project>/.git"` deve passar. O formato
  do `gitfile` é validado pelo `dh` (Go), nunca repetido nem reconferido aqui:
  o shell só consome `linked_worktree`.
- (c) Qualquer outra coisa é `isolation_unavailable`: `dh` ausente ou falhando,
  `linked_worktree` nulo, `.git` um symlink, um repositório bare, um submódulo, um
  formato inválido, conteúdo fora do formato do `gitfile`, ou um `.git` que não é
  um arquivo regular na checagem logo antes do clone.

O sandbox continua sendo um novo `git clone --no-hardlinks "<project>"` numa
pasta temporária descartável, nos dois ramos. O kit nunca usa o
worktree como sandbox e nunca usa `git worktree` para o sandbox.
Hipótese a provar (T-2805): o clone de um worktree vinculado começa no
`HEAD` do worktree. Se não começar, o casamento do passo 2 falha e a chamada termina
em `isolation_unavailable`; a hipótese só pode bloquear uma chamada, nunca abrir uma brecha.

Rode estes em ordem. Todo comando git deste arquivo, em todo passo, roda com
o prefixo `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1` (os comandos
escritos neste passo inclusive). O `mktemp` dá o sufixo aleatório, o modo 0700 e
a unicidade. `<tmp root>` é a única raiz em que o kit cria clones e
a única raiz que a varredura de órfãos (passo 5) lê.

```
CLONE=$(mktemp -d "<tmp root>/dh-cli-<call-id>-XXXXXXXX")
GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git clone --no-hardlinks "<project>" "$CLONE"
GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" remote remove origin
rm -rf "$CLONE/.git/logs"
printf 'pid=%s\nstarted=%s\n' "$PPID" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$CLONE/.git/dh-cli.lock"
```

`$PPID` é o processo que hospeda o shell do delegado, não o próprio shell
de vida curta. `started` é ISO 8601 em UTC; `date -Iseconds` (GNU) também funciona,
mas a forma `date -u +%Y-%m-%dT%H:%M:%SZ` é portátil. O clone tem o próprio
`.git` (hooks, config, refs, objetos, `info/exclude`), não guarda nenhum arquivo ignorado
da árvore viva e nada não commitado: nasce limpo no `HEAD`. Até
o terceiro comando rodar, `origin` ainda aponta para a árvore viva; o reflog
sob `.git/logs` registra o caminho vivo, e é por isso que o quarto comando
o apaga.

Verifique antes de seguir. Cada checagem deve passar:

- `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" remote` não imprime nada.
- `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" rev-parse --git-common-dir` resolve para `$CLONE/.git`, o
  diretório do próprio clone.
- `test ! -e "$CLONE/.git/objects/info/alternates"`.
- `test ! -e "$CLONE/.git/logs"`.

## 2. Popular o clone e casá-lo

O clone deve guardar o estado sob julgamento: o commit atual, mais as
modificações rastreadas e os arquivos novos não ignorados da árvore viva no
congelamento. Rode os comandos da árvore viva a partir de `<project>`, cada um com o prefixo
`GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false`.
Todo comando git rodado contra o clone (`apply`, `status`, `hash-object`)
recebe o mesmo prefixo `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1`, para que uma
configuração global ou de sistema não possa mudar nem rodar nada ali.

1. Parte preparada (staged): salve
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false diff --cached --binary --no-ext-diff --no-textconv --src-prefix=a/ --dst-prefix=b/`
   num arquivo. Se o diff estiver vazio, pule este passo: `git apply` sai com 128
   (`No valid patches in input`) numa entrada vazia, então um diff vazio nunca é
   enviado por pipe. (`git apply --allow-empty`, do git 2.35, é uma alternativa;
   o pulo funciona em qualquer versão do git e é a regra.) Senão, aplique-o com
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" apply --index`.
2. Parte não preparada (unstaged): o mesmo comando sem `--cached`, com o mesmo pulo em
   diff vazio, aplicado com
   `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -C "$CLONE" apply` (sem `--index`).
3. Arquivos novos: a partir de `<project>`, liste-os com
   `git ls-files -z --others --exclude-standard`. Qualquer entrada terminada em `/` (um
   repositório aninhado aparece como `sub/`) é `isolation_unavailable` por si só,
   antes de qualquer coisa ser copiada; isso não fica a cargo do casamento. Se a lista estiver
   vazia, pule a cópia. Senão rode, sob `set -o pipefail` para que um `ls-files` que falha
   não seja mascarado pelos estágios do `tar`,
   `git ls-files -z --others --exclude-standard | tar --null -T - -cf - | tar -xpf - -C "$CLONE"`.
   A lista separada por NUL sobrevive a qualquer nome de arquivo, o `tar` guarda um symlink como
   symlink em vez de segui-lo, e `-p` mantém o modo. O tar GNU aceita
   este pipeline (uma lista vazia também sai com 0 nele; o pulo é mantido mesmo assim);
   que o bsdtar (macOS, BSD) o aceite é uma hipótese, não medida aqui.
   Arquivos ignorados nunca são copiados.

Depois case. Tudo isto deve valer:

- O conjunto julgado (caminhos rastreados mais caminhos não rastreados e não ignorados) é o mesmo
  conjunto de caminhos no clone e na árvore viva. Um caminho ausente no disco registra o
  marcador `deleted` do seu lado, e os marcadores também devem ser iguais.
- Por caminho, por tipo. Um arquivo regular: o mesmo hash de `git hash-object --no-filters` e o
  mesmo modo (`stat -c %a` no GNU, `stat -f %Lp` no BSD), para caminhos rastreados e
  não rastreados igualmente. Um symlink: o mesmo texto de link de `readlink`, nunca
  um hash do alvo para o qual aponta. Um submódulo: o mesmo id de commit.
- `git status --porcelain -uall` (nunca `--ignored`) imprime as mesmas linhas no
  clone (rodado com o prefixo do lado do clone acima) e na árvore viva, inclusive as duas letras de status XY, de modo que o
  clone reproduz o índice (parte preparada e parte não preparada), não só o
  conteúdo. Só o marcador `.git/dh-cli.lock` é ignorado, e ele fica dentro de
  `.git`, então o status não o mostra.
- `test ! -e "$CLONE/.harness"` e
  `test ! -e "$CLONE/.claude/settings.local.json"` passam.

Um estado que não se reproduz exatamente (uma renomeação, uma mudança de modo, um symlink, um
submódulo, um repositório aninhado (recusado antes da cópia), um arquivo binário, ou um XY que a população
em dois passos não consegue recriar) é `isolation_unavailable` (passo 4). A divergência
é o gate. Nunca ajuste o clone à mão até ele casar, e nunca
continue na árvore viva em vez dele.

## 3. Usar o clone e removê-lo

Fixe o diretório de trabalho da CLI no caminho guardado do clone com a flag de
diretório de trabalho do próprio binário, escrita dentro do único comando que começa com
o nome do binário (`grok --cwd`, `mcode exec --cwd`; fato, lido na ajuda de cada
binário em 2026-10-16). `claude` e `opencode` não têm essa flag
(`opencode run` medido em 2026-10-21 na v2.0.32), então as receitas deles mantêm
`cd "$CLONE" && <binary> ...` como as duas exceções nomeadas; essa forma começa com
`cd`, não casa com uma regra de permissão `Bash(claude:*)` ou `Bash(opencode:*)`, então
a chamada pode precisar da aprovação do dono ou ser `not-run`. A forma do opencode é
`cd "<clone>" && <opencode binary> run --auto -m <provider/model>[#<variant>]
--file=<prompt file> -- "Siga o prompt anexado."`. O diretório de trabalho de
uma chamada Bash não persiste para a chamada seguinte (medido em 2026-10-16), então a
flag ou o `cd` pertence ao mesmo comando que inicia o binário. Medido
antes (T-1015): um lançamento só com flag pode mostrar um processo invólucro com um cwd
dentro da árvore viva antes de o binário mudar de diretório; a chamada por isso
não tem `timeout` nem outro invólucro antes do binário. O
prompt não carrega o caminho absoluto da árvore viva.

Quando a chamada termina (sucesso, falha de transporte, timeout ou `tree_changed`),
e depois de rodadas as comparações da árvore viva e do clone, o
delegado remove o clone com `rm -rf` no caminho exato guardado, e só
nele. Antes, a guarda; tudo isto deve valer: o diretório pai do
caminho guardado é igual a `<tmp root>`; o nome-base começa com `dh-cli-` e
não contém `/` nem `..`; `test -d` passa; `test ! -L` passa. Um caminho que
falha na guarda nunca é passado ao `rm`. Depois confirme que o diretório não existe mais. Se
o `rm` falhar, reporte e deixe o diretório para a varredura de órfãos. Nada
do clone é copiado nem mesclado na árvore viva: o veredito ou o texto do plano
volta como a resposta, como antes.

## 4. Isolamento indisponível

A chamada é `not-run`, com o motivo `isolation_unavailable` ("isolamento
indisponivel" no idioma do dono), quando qualquer uma destas condições vale: o projeto
não é um repositório git; o ramo do passo 1 falha (`dh` ausente ou falhando,
`linked_worktree` nulo com um `.git` que não é um diretório, `.git` um symlink,
bare, submódulo, formato inválido, conteúdo fora do formato do `gitfile`, ou um
`.git` que não é um arquivo regular na checagem logo antes do clone; um
worktree vinculado válido não está nesta lista); qualquer passo do congelamento
(`references/live-tree-freeze.md` seção 1) falha (nenhuma CLI inicia e nenhuma comparação roda; o diretório de congelamento parcial é removido sob a guarda do passo 3);
`HEAD` ainda não tem commit; falta disco; `<call-id>`
tem um caractere fora de `[A-Za-z0-9_-]`; uma checagem do passo 1 falha;
a população ou o casamento do passo 2 falha. Remova qualquer clone parcial no seu
caminho exato guardado sob a guarda exata do passo 3, inclusive um que não tem marcador
porque a criação falhou, reporte o motivo e pare. Não há
recuo: um binário sem modo somente leitura nunca roda na árvore viva no
lugar do clone. A rodada não é gasta, e a chamada não conta como aprovação nem
como rejeição.

## 5. Varredura de órfãos

Execuções interrompidas podem deixar clones para trás. Uma vez por lote de despachos, rodada pelo
delegado (um `qa-verifier` despachado) antes da primeira chamada de CLI do lote e nunca dentro de uma
chamada, liste os filhos diretos de `<tmp root>` e remova um diretório
só quando as cinco condições valerem. A varredura cobre dois tipos: um clone
(`dh-cli-`, marcador `.git/dh-cli.lock`) e um diretório de congelamento (`dh-freeze-`,
marcador `dh-freeze.lock` no topo dele, escrito por
`references/live-tree-freeze.md`); as mesmas cinco condições valem para os dois,
com o prefixo e o marcador do seu tipo:

1. O nome começa com `dh-cli-` (ou `dh-freeze-`), é filho direto de `<tmp root>`
   e é um diretório de verdade, não um symlink (`test -d` e `test ! -L`).
2. Pertence ao usuário atual (`stat -c %U` no GNU, ou `find -user`).
3. Contém o marcador dele: `.git/dh-cli.lock` (ou `dh-freeze.lock`).
4. O `pid` no marcador não existe mais (`kill -0 <pid>` falha com
   `no such process`). Um `kill -0` que falha com EPERM significa que o processo existe,
   então o PID conta como vivo e o diretório fica.
5. O `started` no marcador é mais antigo que o timeout da chamada (15 minutos por
   padrão) mais 5 minutos de margem.

As condições são conjuntivas. Nunca toque num diretório que falhe em qualquer uma
delas: sem prefixo, um symlink, sem marcador, outro usuário, um PID vivo ou um
`started` recente é deixado em paz, inclusive um diretório de outro projeto.

## 6. Chamadas em paralelo

Duas chamadas em clones ao mesmo tempo recebem cada uma o próprio `<call-id>` e o
próprio `$CLONE` (resultados distintos de `mktemp`). Cada chamada remove só o
próprio clone, e nunca toca outro diretório `dh-cli-`; só a varredura do passo 5
o faz, antes do lote. Quando essas chamadas se sobrepõem no tempo, a ordem é: um congelamento da árvore viva tirado
antes da primeira delas; depois todo clone criado, populado e casado;
depois todo prompt escrito; só então a primeira CLI iniciada. As chamadas são despachadas em grupos fechados: uma nova chamada em clone nunca começa enquanto uma
CLI de um grupo anterior ainda está rodando. Elas compartilham um
diretório de congelamento (com o nome da primeira chamada), cada uma com o próprio clone e
arquivo de prompt. A comparação da árvore viva e a reversão rodam só depois que a última terminou; qualquer diferença na árvore viva marca toda chamada sobreposta como `not-run`,
`tree_changed`, e a comparação do clone continua por clone (`references/live-tree-freeze.md`,
Ordem). O limite de dois especialistas simultâneos e a serialização de
`opencode` não mudam.

## 7. O que fica fora do clone (residual, dito por inteiro)

O clone é um diretório de trabalho, não um sandbox: clone descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas; escritas fora do projeto não são detectadas.

O que a comparação da árvore viva cobre: a lista de
`git status --porcelain -uall --ignored` com o hash de cada caminho; o hash
de todo arquivo rastreado; a saída de `git ls-files -v --stage` (flags do índice
e ids de blob, o que pega `--skip-worktree` e `--assume-unchanged`); os
bytes do harness dir resolvido no modo `repo`, onde quer que esteja (dentro do
projeto, ou o `.harness/` do checkout principal, não a cópia dentro do worktree,
quando o projeto é um worktree vinculado; no modo `global` ele fica fora do congelamento deste procedimento de clone,
enquanto o congelamento de [live-tree-compare.md](live-tree-compare.md), para uma chamada sem clone, o salva
nos dois modos), de `.claude/settings*.json` e, a partir do diretório comum `<C>`
(`<project>/.git` para um checkout normal, `linked_worktree.commondir` para um
worktree vinculado), de `hooks/`, `config`, `config.worktree` (quando presente), `info/`,
`HEAD` e `packed-refs`, e o `config.worktree` do diretório admin de um worktree
vinculado (restaurado antes de qualquer comando git, como `<C>/config`); a
saída de `git for-each-ref`; e, também a partir de `<C>`, os arquivos sob `refs/` e as entradas sob `worktrees/` quando existirem.
Quando `<project>/.git` é um `gitfile`, também os bytes do `gitfile` e, no
diretório admin do worktree, `HEAD`, `commondir`, `gitdir`, os arquivos
sob `refs/` e a presença de `locked`
(`references/live-tree-freeze.md` seção 1). Os itens de descoberta do repositório
(o `.git` de um checkout normal ainda um diretório de verdade; os bytes do `gitfile`,
`<admin>/commondir` e `<admin>/gitdir`; `<C>/commondir` e
`<project>/.git/commondir` ainda ausentes) são comparados por bytes com ferramentas que não são
do git e restaurados ou removidos antes do primeiro comando git; se isso for
impossível, nenhum comando git roda e o veredito é descartado, `not-run`,
`not revertible`. Um item ausente nos dois instantâneos
não é uma diferença e nunca um erro. Os caminhos de confiança de `<C>`, os dois arquivos `config.worktree`, o `gitfile`, `<admin>/commondir` e `<admin>/gitdir` são restaurados; os arquivos sob `refs/`, as
entradas de `worktrees/` e os demais itens do diretório admin são só detectados: uma
alteração é `tree_changed` e é reportada como `not revertible: <path>`.

Não coberto pela comparação da árvore viva, portanto não detectado (esta lista de caminhos vivos do diretório git é a
que o autor conhece e não é exaustiva; qualquer caminho vivo do diretório git que essa comparação não nomeie
está fora da cobertura): `index` como bytes brutos (só o conteúdo lógico dele
é lido por `ls-files`; o mesmo vale para o `index` de cada diretório admin de worktree),
`logs/` (reflogs), objetos soltos e packs, `modules/`, `FETCH_HEAD`,
`ORIG_HEAD` e `COMMIT_EDITMSG`.

Residuais de aceitar um worktree vinculado, não detectados ou não impedidos (ADR-008 D4):

- um escritor que roda como o mesmo usuário pode forjar um arranjo inteiro coerente, e
  a validação não o distingue de um worktree legítimo (a mesma confiança dada
  ao `.git/config` de um checkout normal);
- a árvore do checkout principal e de outros worktrees está fora do
  `<project>` da chamada e cai em "uma escrita fora do diretório do
  projeto" em (a) abaixo. Quando o `<project>` é o checkout principal, os
  worktrees na pasta `.worktrees/` ignorada dele estão dentro dele: o
  status com `--ignored` os lista, então uma escrita de `builder` ali durante a chamada
  descarta o veredito (um alarme falso, nunca uma falha de detecção; revise o caminho do worktree
  em vez do checkout principal enquanto um worktree está sendo escrito);
- o `gitfile` pode ser trocado entre a validação e o clone; isso é
  detectado só se a troca persistir até o segundo instantâneo;
- o harness dir global fica fora do congelamento deste procedimento de clone.

RISK-002 continua mitigado: o congelamento detecta, nunca impede.

Aceito pelo dono em 2026-09-29, até o lote do sandbox do sistema operacional, e não
detectado:

- (a) uma escrita fora do diretório do projeto: dotfiles sob o diretório
  pessoal, `~/.ssh`, a configuração global do git, credenciais e a
  configuração da própria CLI, outros repositórios, outras pastas sob `/tmp`;
- (b) processos que sobrevivem à chamada, chamadas de rede e `git push` para qualquer
  remoto que não seja o repositório vivo;
- (e) ler segredos (`.env`, `~/.aws`, qualquer coisa que o usuário alcance) e
  devolvê-los na resposta, porque o clone só os esconde de um
  caminho relativo.

O diretório de congelamento de `references/live-tree-freeze.md` é coberto por (a) e
(e) acima. Ele guarda cópias de `.harness/`, `.claude/settings*.json` e dos
arquivos modificados ou não rastreados e não ignorados, fica em `<tmp root>`, e uma CLI
rodando como o mesmo usuário pode lê-lo e alterá-lo; o modo 0700 não impede isso.
Adulterar as cópias é uma escrita fora do projeto (a), lê-las é
ler segredos (e). O mesmo vale entre chamadas sobrepostas: uma CLI em execução
pode ler ou alterar o arquivo de prompt de uma chamada irmã ou o diretório de congelamento compartilhado
(residual (a)). Escrever todo prompt antes de qualquer CLI iniciar estreita esta janela;
não a fecha. Um comando separado depois da chamada (um kill pelo pid registrado quando a chamada
começou, ou por um padrão com colchetes como `pkill -f '[o]pencode run'`; Ordem, passo
4 de `live-tree-freeze.md`) para os processos remanescentes desta chamada e estreita
de novo o residual (b); não o fecha. Uma
revisão `agy` interrompida pode deixar o arquivo `dh-prompt-XXXXXXXX` dela (modo 0600)
sob `<tmp root>`: ele não é varrido e fica para o dono remover.

Também não detectados: (c) uma escrita na árvore viva que é desfeita antes do
segundo instantâneo ou que preserva os bytes; (f) uma escrita concorrente de outro
processo, o Coordenador inclusive, que `live-tree-compare.md` já manda
evitar. Só um isolamento do sistema operacional fecharia isso; é um lote
posterior.
