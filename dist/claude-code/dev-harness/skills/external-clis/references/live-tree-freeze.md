# Congelamento da árvore viva, comparação do clone e higiene do prompt

Este arquivo é o procedimento de comparação para uma chamada de CLI que roda num
clone descartável (ciclo de vida em [clone-isolation.md](clone-isolation.md)). Uma chamada sem
clone (`codex`, e a revisão do `agy` na árvore viva) usa
[live-tree-compare.md](live-tree-compare.md) em vez deste, que também explica por que
a revisão do `agy` na árvore viva de um projeto cujo `.git` não é um diretório termina em
`isolation_unavailable` e não roda.

O delegado executa todos os passos deste arquivo: um `qa-verifier` despachado, ou
`builder` na raia `infra` (`SKILL.md`, "Quem executa a chamada"). Ele devolve o
sinal de veredito e o caminho relativo da saída bruta; o Coordenador não executa
nenhum destes passos.

A qualificação para toda menção a este mecanismo, dita por inteiro:

clone descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas; escritas fora do projeto não são detectadas

Duas comparações rodam, em duas árvores. A comparação da árvore viva detecta e
reverte o que a CLI escreveu onde não devia. A comparação do clone
julga se a CLI alterou o conteúdo que recebeu. Nenhuma das duas impede
nada.

## Ordem

1. Precondições e congelamento (seção 1). A precondição é o passo 1 de
   clone-isolation.md: `.git` um diretório de verdade, ou um `gitfile` com um
   `linked_worktree` não nulo e a checagem `test -f`/`test ! -L` logo antes do clone;
   qualquer outra coisa é `isolation_unavailable`. Quando chamadas em clone se sobrepõem no tempo,
   tire um congelamento antes da primeira delas. Salve as cópias fora da
   árvore do projeto e fora do clone: crie um diretório de congelamento com
   `mktemp -d "<tmp root>/dh-freeze-<call-id>-XXXXXXXX"` (`<tmp root>` e
   `<call-id>` como em clone-isolation.md, Termos; o call-id é o da primeira
   chamada; modo 0700) e escreva o marcador `<freeze dir>/dh-freeze.lock` com as
   mesmas duas linhas do marcador do clone (`pid=` e `started=`,
   clone-isolation.md passo 1). Depois calcule o digest do congelamento (fim da
   seção 1). Se qualquer passo deste congelamento falhar (uma ferramenta ausente inclusive, ver a nota sobre
   ferramentas GNU em Notação), a chamada é `not-run`, `isolation_unavailable`: nenhuma
   CLI inicia e nenhuma comparação roda. Remova o diretório de congelamento parcial sob
   a guarda do passo 8.
2. Crie, popule e case todo clone sobreposto (clone-isolation.md passos
   1 e 2) e congele cada clone (seção 3). Chamadas sobrepostas compartilham o único
   diretório de congelamento, com um clone e um arquivo de prompt para cada uma.
3. Escreva o prompt de cada chamada (seção 5), todos antes de qualquer CLI iniciar.
4. Inicie cada CLI como um único comando de shell que começa com o nome do binário, em
   uma chamada de ferramenta própria, com a flag de diretório de trabalho do próprio binário quando a receita
   dele mostrar uma (`grok --cwd`, `mcode exec --cwd`); `claude` e `opencode` não têm
   essa flag e mantêm `cd "$CLONE" && <binary> ...` como as duas exceções nomeadas
   (`references/claude.md`, `references/opencode.md`; a forma começa com `cd`,
   não casa com uma regra de permissão `Bash(<binary>:*)` e pode precisar da aprovação
   do dono; um `cd` separado numa chamada própria não persiste). A forma do opencode
   é `cd "<clone>" && <opencode binary> run --auto -m
   <provider/model>[#<variant>] --file=<prompt file> -- "Siga o prompt
   anexado."`. Nenhum `timeout`, `setsid`, subshell,
   `&&`, `;` ou pipe embrulha ou encadeia a chamada, e o congelamento e as comparações
   posteriores são comandos separados (uma regra de permissão `Bash(<binary>:*)`
   casa só com um comando que começa com o binário). A flag de tempo do próprio
   binário, ou o timeout da ferramenta Bash, limita a chamada. Depois, como um comando
   separado rodado quando a resposta está em disco ou o timeout passou (antes de qualquer
   comparação), pare os processos remanescentes desta chamada pelo pid: `kill <pid>`. O
   pid vem da própria chamada: a chamada roda em segundo plano
   (`run_in_background`) e o pid dela é registrado quando ela começa, ou o delegado
   lista `pgrep -x <binary>` antes da chamada e depois dela e mata só o pid
   novo. Nunca `pkill -f` com o caminho do clone ou o texto do comando: o padrão
   casa com o shell do próprio delegado (`references/gotchas.md`); se um padrão for
   inevitável, use um padrão com colchetes (`'[o]pencode run'`). O pid é único para
   esta chamada, então nunca casa com as sessões do próprio dono do mesmo binário; um
   pid que já se foi é ignorado. Isso estreita o residual (b) de clone-isolation.md
   seção 7 (um processo que sobrevive à chamada); não o fecha, já que um
   filho que tem outro pid escapa. Reporte-o como `leftover processes` na
   Checagem do espaço de trabalho.
5. Quando uma chamada termina, recompare o clone dela (seção 3). Essa comparação continua por clone.
6. Depois que a última chamada sobreposta terminou, e só então, na árvore viva:
   verifique o digest do congelamento, extraia as cópias, compare e restaure os caminhos de
   confiança com ferramentas que não sejam do git (seção 2, parte A) e, só quando a parte A não relatar
   nenhuma diferença restante, rode a comparação baseada em git (seção 2, parte B).
7. Reverta o que a comparação achou (seção 4).
8. Remova cada clone (clone-isolation.md passo 3), depois o diretório de congelamento,
   arquivos de prompt inclusive, com a mesma guarda (o pai é igual a `<tmp root>`,
   o nome-base começa com `dh-freeze-`, `test -d`, `test ! -L`). Nada nunca é
   varrido por idade por conta própria.

Chamadas em clone são despachadas em grupos fechados: uma nova chamada em clone nunca começa (o
clone dela nunca é criado) enquanto uma CLI de um grupo anterior ainda está rodando. Um
grupo é o conjunto de chamadas que se sobrepõem no tempo sob um congelamento.

Nada mais pode escrever na árvore viva nem no `<harness dir>` entre o passo 1 e o passo 6, delegado e Coordenador inclusive (o diretório de congelamento é compartilhado por chamadas sobrepostas, então isto cobre todas elas): um build, uma
regeneração de pacote ou uma execução de testes iniciada em paralelo muda arquivos que a
comparação depois atribui à CLI (o caso medido está em
`live-tree-compare.md`). Chamadas em clone sobrepostas são governadas pela mesma regra: a comparação da árvore viva e a
reversão rodam só depois que a última terminou, porque a reversão de uma irmã
esconderia o escritor de outro modo. Qualquer diferença na árvore viva marca toda chamada sobreposta como
`not-run`, `tree_changed`, como `coordinator.md` já faz para chamadas sem
clone. O limite de dois especialistas simultâneos não muda. A reversão na
seção 4 só é segura quando nenhum processo iniciado pela chamada ainda está rodando; um
processo que sobrevive à chamada pode escrever de novo depois da reversão, e isso
fica não detectado (clone-isolation.md seção 7, item (b)). Espere, ou rode
o trabalho paralelo num checkout separado.

Motivos de `not-run` para uma chamada em clone: `tree_changed` ("mudança detectada"),
deste arquivo, e `isolation_unavailable` ("isolamento indisponivel"), de
clone-isolation.md passo 4. Os dois descartam o veredito, nenhum conta como
aprovação ou rejeição, e nenhum gasta uma rodada de correção.

## Notação

`R5GIT` representa
`GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 git -c core.fsmonitor=false -c core.untrackedCache=false -c core.hooksPath=/dev/null`,
isto é, o prefixo `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1 -c core.fsmonitor=false -c core.untrackedCache=false -c core.hooksPath=/dev/null`
posto em `git`. Todo comando git da comparação da árvore viva, da reversão
e da comparação do clone roda como `R5GIT <subcommand>`, de modo que uma configuração
global ou de sistema não possa nem mudar o resultado nem executar um comando, e nenhum hook
rode. A sobreposição de `core.hooksPath` é necessária porque, de outro modo, um comando git que
atualiza uma ref ou o índice roda os hooks do próprio repositório vivo
(medido: `update-ref` rodou duas vezes um hook `reference-transaction` plantado; com
a sobreposição não rodou nenhum). `<project>` é a árvore viva.

Este arquivo supõe ferramentas GNU em alguns pontos: `sha256sum` (macOS: `shasum -a 256`),
`stat -c` (BSD: `stat -f`) e `tar --null -T`.
Uma ferramenta ausente durante o congelamento, sem equivalente, é
`isolation_unavailable`.

Regra de hash para um caminho, usada em todo o restante: um arquivo regular é hasheado
com `R5GIT hash-object --no-filters -- <path>`; um symlink é registrado como o texto de
`readlink` dele e nunca é seguido; um caminho ausente no disco registra o marcador
`deleted`. Um diretório nunca é hasheado, só expandido para os arquivos e
symlinks sob ele (`find <dir> \( -type f -o -type l \)`, que não
segue symlinks). Todo caminho registrado também registra o modo (`stat -c %a` no
GNU, `stat -f %Lp` no BSD). Entradas que não são nem arquivo regular nem
symlink (diretórios, fifos, sockets) são tratadas pela lista de congelamento da
seção 1, não por esta regra.

## 1. Congelar a árvore viva

Prevenção do RISK-002: nada escreve na árvore nem no `<harness dir>` entre o congelamento e a comparação, Coordenador inclusive. A detecção e a reversão abaixo cobrem o `<harness dir>` resolvido no modo `repo`, onde quer que esteja (dentro do projeto, ou o `.harness/` do checkout principal quando o projeto é um worktree vinculado); no modo `global` ele fica fora do congelamento deste procedimento (residual do RISK-002).

Rode a partir de `<project>`. Registre tudo isto, antes da chamada:

- Os caminhos de confiança, como bytes e modo, salvos com ferramentas que não sejam do git (`cp -p`, ou
  `tar` a partir de uma lista NUL explícita, abaixo): tudo sob o `<harness dir>` resolvido no modo `repo`, onde quer que esteja (dentro do projeto, ou o `.harness/` do checkout principal quando o projeto é um worktree vinculado), e sob `.harness/` quando existir. No modo `global` o harness dir (`<home>/projects/<name>/`) fica fora do congelamento deste procedimento, um residual declarado do RISK-002; só o congelamento de `live-tree-compare.md`, para uma chamada sem clone, o salva por bytes nos dois modos;
  `.claude/settings.json` e `.claude/settings.local.json` (cada um quando
  presente, e a lista dos ausentes); do diretório comum `<C>`, `hooks/`,
  `config`, `config.worktree` (quando presente), `info/`, `HEAD` e `packed-refs`; para um worktree vinculado, também o `config.worktree` do diretório admin (lido do campo `gitdir` de `linked_worktree`; parte (b) abaixo); e todo arquivo `.gitattributes`
  na árvore de trabalho (achado com `find . -path ./.git -prune -o -name .gitattributes -print`).
  `<C>` é `<project>/.git` para um checkout normal e `linked_worktree.commondir`
  (de `dh harness-path --json "<project>"`) para um worktree vinculado, em que
  `<project>/.git` é só um `gitfile` e não guarda nenhum desses caminhos; um item
  ausente nos dois instantâneos não é uma diferença e nunca um erro. Num worktree vinculado o harness dir é o que `dh harness-path --json` dá
  em `dir` (o `.harness/` do checkout principal no modo `repo`), não a cópia ignorada
  de `.harness/` dentro do worktree; no modo `global` ele continua fora do congelamento deste procedimento.
  Também são itens de confiança, comparados por bytes só com ferramentas que não sejam do git (`cat`, `cmp`,
  `stat`; seção 2, parte A): num checkout normal, que `<project>/.git` ainda
  é um diretório de verdade (não um arquivo nem um symlink); para um worktree vinculado, os
  bytes do `gitfile` `<project>/.git`, e `<admin>/commondir` e
  `<admin>/gitdir`; e que `<C>/commondir` e `<project>/.git/commondir` estão
  ausentes quando estavam ausentes no congelamento (um plantado faz o git
  redescobrir outro repositório, com a configuração e o
  `info/attributes` desse repositório). Eles são salvos antes da chamada como os outros caminhos de confiança.
  Os arquivos `config.worktree` são caminhos de confiança como `<C>/config`: com
  `extensions.worktreeConfig` ligado, o git os lê (`core.attributesFile`,
  `filter.<x>.smudge`, `include.path`, `core.worktree`), então um plantado poderia rodar
  um comando no git do próprio delegado. Eles são salvos antes da chamada e
  restaurados com ferramentas que não sejam do git na parte A da seção 2, antes de qualquer comando git.
- `R5GIT status --porcelain -z -uall --ignored`, expandido para arquivos (um diretório
  ignorado como `!! build/` é percorrido, sem seguir symlinks), com
  o hash de cada caminho listado, ignorado ou não. Todo caminho ignorado é hasheado
  (decisão D2 do dono); um limite de custo só é definido por uma emenda depois de uma
  medição.
- O hash de todo arquivo rastreado (`R5GIT ls-files -z`), não só os
  listados.
- `R5GIT ls-files -v --stage`, salvo como texto. A primeira letra é a flag
  do índice: `H` nenhuma flag, `S` só skip-worktree, `h` só assume-unchanged,
  `s` skip-worktree e assume-unchanged juntos. Cada entrada também guarda o modo
  e o id do blob.
- `R5GIT for-each-ref --format='%(refname) %(objectname)'`, salvo como texto.
- Diretórios git, duas partes; `config.worktree` (os dois diretórios) é um caminho de confiança restaurado
  (acima), e o mesmo vale para o `gitfile`, `<admin>/commondir` e `<admin>/gitdir`;
  o resto é só detectado. Parte (a), todo projeto: no diretório comum
  `<C>`, além dos caminhos de confiança acima, os arquivos sob `refs/` (bytes e
  modo) e as entradas sob `worktrees/` quando existirem, com os arquivos de
  cada entrada que a parte (b) lista. Parte (b), só quando `<project>/.git` é um
  `gitfile`: os bytes do `gitfile` e, no diretório admin do worktree (lido do
  campo `gitdir` de `linked_worktree`), `HEAD`, `commondir`, `gitdir`,
  `config.worktree`, os arquivos sob `refs/` e a presença de `locked`. Um
  item ausente nos dois instantâneos não é uma diferença e nunca um erro. Os arquivos
  `index` brutos (do diretório comum e de cada diretório admin) ficam fora do congelamento.
  Uma alteração nesses itens durante a chamada é `tree_changed`, como num checkout
  normal, e o veredito é descartado. Exceto os caminhos de confiança restaurados acima, eles não são
  salvos para reversão: cada um é reportado como `not revertible: <path>`, como um
  caminho ignorado alterado.
- Toda entrada sob o projeto que não é arquivo regular nem symlink,
  excluindo `.git`: `find . -path ./.git -prune -o ! -type f ! -type l -print`
  (diretórios inclusive), salva como texto com o modo de cada entrada
  (`stat -c %a` no GNU, `stat -f %Lp` no BSD). Um novo diretório vazio, FIFO ou socket
  então aparece como um caminho criado, e a seção 4 o remove com `rm -rf --`
  depois da checagem de ancestrais. Um diretório que sumiu é recriado com
  `mkdir -p`; um diretório cujo modo mudou recebe `chmod <saved mode> <dir>`
  durante a reversão; um FIFO ou socket que sumiu é reportado como
  `not revertible: <path>`.

Salve os bytes e o modo só dos caminhos de que uma reversão pode precisar: todo
caminho não ignorado listado por `status` (modificado ou não rastreado), todo caminho rastreado
cuja letra salva de `ls-files -v` não é `H` (um arquivo marcado `h`, `S` ou `s`
pode guardar uma edição não commitada que o `status` esconde) e os caminhos de confiança. Um
caminho rastreado é "limpo" só quando a letra dele é `H` e o `status` não o
lista; um caminho rastreado limpo não é copiado, é restaurado do `HEAD` (seção 4,
item 3). Escreva a lista separada por NUL dos
caminhos salvos em `<freeze dir>/paths.nul`, deixando de fora todo caminho registrado como
`deleted` (o tar falha com `Cannot stat` num caminho que não existe, saída 2; um
caminho apagado não tem nada a salvar e é registrado como `deleted` na lista de
hashes), depois rode
`tar --null -T "<freeze dir>/paths.nul" -cf "<freeze dir>/live.tar"` a partir de
`<project>` (o `tar` guarda um symlink como symlink; o tar GNU aceita isto, e
que o bsdtar também aceite é uma hipótese, não medida aqui). Um caminho ignorado é
hasheado e o modo dele registrado, mas os bytes dele não são salvos: um caminho ignorado alterado ou
apagado não pode ser revertido. Isso descarta o veredito e é
reportado como `not revertible: <path>`. Registre também, para todo diretório ancestral
de um caminho salvo, se o congelamento o viu como diretório de verdade
(`test -d` e `test ! -L`).

Logo depois que os registros acima estão completos, calcule o `sha256sum` de
`<freeze dir>/live.tar` e de todo instantâneo de texto salvo em `<freeze dir>` (o
manifesto: a lista de hashes, a lista do `status`, `ls-files -v --stage`,
`for-each-ref`, a lista de entradas que não são arquivos e os registros de ancestrais). Guarde
o digest na saída da própria sessão do delegado, nunca num arquivo que a CLI
possa alcançar. Isso detecta adulteração das cópias; não a impede.

## 2. Comparar a árvore viva, depois da chamada

Antes da parte A, verifique o congelamento: recalcule o `sha256sum` de `live.tar` e dos
arquivos do manifesto e compare com o digest que o delegado guardou. Uma divergência
é `tree_changed` mais `freeze tampered`: nada é restaurado das
cópias, nenhum comando git roda, e o relatório lista só a divergência do digest e
uma listagem de `find` dos caminhos de confiança vivos para o dono restaurar à mão (sem
`cmp` contra as cópias não confiáveis). Isto é detecção, não prevenção: uma CLI
que roda como o mesmo usuário ainda pode escrever no diretório de congelamento, e
só a adulteração é percebida. Quando o digest confere, extraia as cópias:
primeiro exija `test -d "<freeze dir>" && test ! -L "<freeze dir>"` no caminho guardado
do diretório de congelamento; uma checagem que falha é `freeze tampered`. Crie o
diretório de extração novo, `mktemp -d "<freeze dir>/copy-XXXXXXXX"`, guarde o caminho
literal dele como `<copy dir>` e extraia com
`tar -xpf "<freeze dir>/live.tar" -C "<copy dir>"`. Nunca faça `mkdir -p` num nome
fixo aqui: um symlink plantado num nome fixo mandaria a extração
para outro lugar. Toda restauração abaixo lê de `<copy dir>`.

Checagem de ancestrais, antes de qualquer `find` ou `rm` na parte A ou na parte B. Todo caminho
envolvido (todo caminho salvo, criado ou listado, e as raízes de confiança `<C>`,
`<C>/hooks`, `<C>/info`, o `<harness dir>` resolvido no modo `repo` (`dir` de `dh harness-path --json "<project>"`: `<project>/.harness` num checkout normal, `<main>/.harness` num worktree vinculado), `.harness` quando existir, e `.claude`) deve ter cada ancestral
ainda um diretório de verdade (`test -d` e `test ! -L`). Um ancestral que o congelamento
registrou como diretório de verdade e que agora é um symlink ou um arquivo regular
foi plantado: remova só esse ancestral (`rm -- <ancestor>`), tire da lista
de criados todo caminho abaixo dele (eles foram alcançados pelo link plantado e não estão
no projeto) e recrie o ancestral com `mkdir`. `rm -rf` não
segue o componente final do argumento, e os componentes intermediários são
os que esta checagem cobre. Nunca passe um diretório de confiança com barra final
para `find` ou `rm` (uma barra final faz os dois seguirem um symlink).

Parte A, caminhos de confiança primeiro, só com ferramentas que não sejam do git. A configuração
local do repositório pode executar um comando durante `status`, `hash-object` (sem `--no-filters`),
`checkout` e similares (um driver `filter.<name>.clean` ou `smudge`, selecionado
por um arquivo `.gitattributes`), então um `<C>/config` ou `.gitattributes` hostil
deve ser posto de volta antes de qualquer comando git rodar. Usando só `tar`, `cmp`, `find`, `rm`, `mv`, `cat`, `mktemp`, `chmod`, `mkdir` e
`ln`, compare todo caminho de confiança com a cópia em `<copy dir>`, liste os diretórios
de confiança (o `<harness dir>` resolvido no modo `repo` (`dir` de `dh harness-path --json "<project>"`: `<project>/.harness` num checkout normal, `<main>/.harness` num worktree vinculado), `.harness` quando existir, `<C>/hooks`, `<C>/info`) com `find` e compare
a lista, e restaure toda diferença de uma vez como na seção 4, itens 1 e 2, só para os
caminhos de confiança e os arquivos `.gitattributes`. A mesma comparação
de bytes sem git cobre os itens de descoberta do repositório da seção 1: os bytes do
`gitfile` (ou, num checkout normal, que `.git` ainda é um diretório de verdade, não um
arquivo nem um symlink), `<admin>/commondir` e `<admin>/gitdir`, e que
`<C>/commondir` e `<project>/.git/commondir` ainda estão ausentes quando estavam
ausentes no congelamento. Numa diferença, restaure da cópia do congelamento ou remova
o arquivo plantado antes de qualquer comando git. Se não puder ser restaurado nem removido,
nenhum comando git roda: pule a parte B e a reversão, e registre o
veredito como descartado, `not-run`, com o motivo `not revertible`. Nenhum comando git
roda na parte A. Uma diferença aqui é `tree_changed`.

Ordem: toda comparação de bytes e restauração da parte A acontece antes do primeiro comando
git da parte B; os comandos git da parte B e da reversão (itens 3 a
5 da seção 4) rodam só depois que a parte A não relata nenhuma diferença restante.

Parte B, depois a comparação baseada em git. Recalcule todo outro item da seção 1
e compare cada um com o congelamento. Qualquer diferença é `tree_changed`: um caminho
criado (ausente, agora presente), apagado ou alterado (um hash, modo ou
texto de link diferente); uma entrada de índice cuja letra de flag ou id de blob difere; uma ref
criada, movida ou apagada (`for-each-ref` difere). O veredito é descartado
independentemente do que a CLI disse, o restante da reversão da seção 4 (itens
1 e 2 para os outros caminhos, depois itens 3 a 5) roda, e o relatório diz
`tree_changed`.

## 3. Congelar e comparar o clone

Logo depois do casamento do passo 2 de clone-isolation.md, registre o conjunto julgado: os
caminhos rastreados (`R5GIT -C <clone path> ls-files -z`) e os caminhos
não rastreados e não ignorados (`R5GIT -C <clone path> ls-files -z --others --exclude-standard`),
cada um com o hash dele (`deleted` quando ausente).

Antes de o clone ser removido, recalcule o hash de cada caminho desse conjunto
inicial:

- Um hash alterado, ou um caminho do conjunto removido, é `tree_changed`. O
  veredito é descartado e nada é revertido: o clone é jogado fora.
- Caminhos ignorados do clone e qualquer arquivo novo não rastreado estão fora do conjunto.
  Nenhum chega à árvore viva nem muda o conteúdo julgado, então nenhum
  descarta o veredito. Cada arquivo novo não rastreado (um caminho em `ls-files` agora que
  não está no conjunto inicial) é registrado como nota no relatório da chamada.
- O `.git` do clone, inclusive `dh-cli.lock`, não é comparado.

## 4. Reverter a árvore viva

Reverta a partir das cópias do congelamento e, para um caminho rastreado limpo, a partir do `HEAD` (item
3). Nunca escreva através de um caminho existente: o que está lá agora é removido primeiro,
sem seguir um symlink, e o conteúdo salvo é posto por uma renomeação.
Os itens 1 e 2 primeiro restauram os caminhos de confiança e os arquivos `.gitattributes`, na
parte A da seção 2, sem comando git; na parte B eles rodam de novo para todo
outro caminho. Os itens 3 a 5 rodam só na parte B, depois da parte A. Depois, recalcule
a seção 1 e confirme que é igual ao congelamento; se não for, reporte `reversion
incomplete` ao lado de `tree_changed`.

1. **Um caminho criado** (ausente antes): `rm -rf -- <path>`, que não
   segue um symlink. Um caminho rastreado salvo com a letra `S` cujo arquivo no worktree
   estava ausente antes da chamada é registrado `deleted`, então um arquivo que a CLI cria
   ali é removido por este item. Todos os caminhos criados são removidos antes de qualquer caminho ser
   restaurado, de modo que um caminho criado nunca bloqueia uma restauração.
2. **Um caminho alterado ou apagado, com cópia salva.**
   - Primeiro os ancestrais (a checagem de ancestrais da seção 2 já rodou;
     repita-a aqui para cada caminho), da raiz do projeto para baixo. O git não rastreia
     através de um symlink, então todo ancestral de um caminho rastreado era um diretório
     de verdade antes da chamada por definição; para qualquer outro caminho o congelamento
     registrou se o ancestral era um diretório de verdade. Um ancestral que era um
     diretório de verdade e agora é um symlink (`test -L`) ou um arquivo regular
     foi plantado: remova-o com `rm -- <ancestor>`. Um `.git` (ou qualquer outro
     ancestral) que já era um symlink antes da chamada é deixado em paz. Depois
     `mkdir -p -- <parent>` para todo ancestral ausente.
   - Um caminho com bytes e modo salvos, escrito como um subshell encadeado só, para que o
     nome temporário não se perca:
     `( umask 077; tmp=$(mktemp <parent>/.dh-restore-XXXXXXXX) && cat -- <saved copy> > "$tmp" && chmod <saved mode> "$tmp" && rm -rf -- <path> && mv -f "$tmp" <path> )`.
     Uma renomeação substitui a própria entrada de destino e não segue um symlink
     final: um symlink plantado sobre `.harness/project.yaml` é removido e o alvo dele
     fica intocado (medido: saída 0, alvo intocado, arquivo regular com
     o modo salvo).
   - Um caminho que era um symlink antes da chamada, na mesma forma, depois de
     remover o nome que o `mktemp` criou (`ln -s` falha num nome existente):
     `( tmp=$(mktemp <parent>/.dh-restore-XXXXXXXX) && rm -f "$tmp" && ln -s <saved link text> "$tmp" && rm -rf -- <path> && mv -f "$tmp" <path> )`.
3. **Um caminho rastreado limpo** (letra salva `H`, não listado por `status`, portanto sem
   cópia salva) que agora está alterado ou apagado. Só na parte B, depois de a parte A ter
   restaurado os hooks, a config e os atributos, com a
   sobreposição de `core.hooksPath` em vigor. Trate os ancestrais como no item 2, depois
   remova o que está lá, `rm -rf -- <path>`, depois limpe as duas flags do índice, um
   comando cada (`R5GIT update-index --no-skip-worktree -- <path>`, depois
   `R5GIT update-index --no-assume-unchanged -- <path>`), e só então rode
   `R5GIT checkout HEAD -- <path>`. O `checkout -- <path>` simples não é usado: ele
   restaura o blob do índice, e uma entrada de índice adulterada devolve o conteúdo
   adulterado (medido). `checkout HEAD -- <path>` restaura o conteúdo do `HEAD`,
   reinicia a entrada do índice para o `HEAD`, deixando `H`. Medido (T-1015, git 2.53.0):
   `checkout HEAD -- <path>` num caminho removido cuja entrada de índice carrega
   skip-worktree falha com saída 1,
   `pathspec did not match any file(s) known to git`; limpar as flags antes, como acima, funciona e deixa `H`. O item 4 então
   define as letras salvas. Para um caminho assim o índice é igual ao `HEAD` por
   construção, então o `HEAD` é o estado dele anterior à chamada. Qualquer filtro que rode ali é
   a configuração anterior à chamada do próprio dono, restaurada na parte A; uma
   restauração por `cat-file` quebraria repositórios filtrados como git-lfs.
4. **Entradas de índice, depois flags**, depois dos itens 2 e 3. As letras e os ids de blob
   salvos têm a palavra final: compare o `ls-files -v --stage` salvo com
   o novo, caminho a caminho:
   - primeiro as entradas: uma entrada de stage cujo modo ou id de blob difere é
     posta de volta com `R5GIT update-index --cacheinfo <mode>,<blob id>,<path>` usando os
     valores salvos; uma entrada que não existia antes é removida com
     `R5GIT update-index --force-remove -- <path>`;
   - depois as flags, uma flag por invocação (medido: `--no-skip-worktree
     --no-assume-unchanged` numa invocação deixa `S` ligado; invocações
     separadas o limpam). Limpe primeiro: `R5GIT update-index --no-skip-worktree -- <path>`,
     depois, como comando próprio, `R5GIT update-index --no-assume-unchanged -- <path>`.
     Então defina o que a letra salva tinha, cada uma em seu próprio comando: para `S`,
     `R5GIT update-index --skip-worktree -- <path>`; para `h`,
     `R5GIT update-index --assume-unchanged -- <path>`; para `s` as duas, em dois
     comandos; para `H` nada mais.
5. **Refs**, por último, a partir da lista salva de `for-each-ref`:
   - uma ref criada pela chamada: `R5GIT update-ref -d <refname> <new id>`;
   - uma ref movida pela chamada: `R5GIT update-ref <refname> <old id> <new id>`;
   - uma ref apagada pela chamada: `R5GIT update-ref <refname> <old id>`.
   `update-ref` pode reescrever `<C>/packed-refs`. O recálculo final inclui
   os caminhos de confiança, então um `packed-refs` que difere depois deste item é
   restaurado de novo da cópia do congelamento (item 2, sem comando git).

A reversão restaura aquilo para que o congelamento da árvore viva salvou cópias, `config.worktree` inclusive. Os
outros itens do diretório git da seção 1 são detectados e reportados como
`not revertible: <path>`, nunca revertidos. Ela não restaura o que esse congelamento não nomeia:
objetos que um `git tag -a` escreveu, reflogs e os demais caminhos vivos do diretório git listados
em clone-isolation.md seção 7 ficam como a CLI os deixou, e o veredito já
está descartado.

## 5. O prompt

O prompt leva o contexto como hoje: papel, skill e pacote. Uma chamada em clone
o escreve num arquivo criado por `( umask 077; mktemp "<freeze dir>/prompt-XXXXXXXX.md" )`,
único por chamada, dentro do diretório de congelamento: ele é removido e varrido junto
com esse diretório, e nenhum arquivo de prompt é jamais varrido por idade por conta própria. (Uma chamada
sem clone, a revisão do `agy`, cria o prompt dela com
`( umask 077; mktemp "<tmp root>/dh-prompt-XXXXXXXX" )` e o remove na própria
limpeza.) Ele
mantém o bloco `commands.test`/`commands.lint` da cláusula
anti-delegação, e esses comandos rodam no clone.

A regra do caminho vivo ausente é higiene, não uma fronteira. O prompt não deve conter o caminho absoluto
da árvore viva; confira com `grep -F -- "<absolute path of the project>" <prompt file>`,
que não deve achar nada. A CLI pode achar esse caminho por outros meios (as
comparações acima detectam o que ela faz com ele), então isto não impede uma
escrita.

Um artefato que um teste legítimo gera no clone (um arquivo ignorado
ou um arquivo novo não rastreado) está fora do conjunto julgado e não descarta o
veredito; um teste que reescreve um arquivo rastreado o descarta como qualquer outra alteração.
O resultado de uma checagem que uma CLI reporta nunca é evidência: um `qa-verifier` despachado
roda as checagens de novo. Um `_test.go` ignorado plantado no clone só pode enganar a
resposta da própria CLI, que já é entrada não confiável. Se um comando falha por um
arquivo ignorado ausente (uma dependência, um `.env`), a CLI reporta "não validado",
e o Coordenador lê esse relato como nem uma falha do código nem uma
aprovação.

## 6. O que o relatório leva

A linha "Checagem do espaço de trabalho" do bloco de resposta em `SKILL.md` ("Quem executa a chamada") recebe, para uma
chamada em clone: `clean` ou `reverted (<what changed, now reverted>)` para a
árvore viva; `clone removed: <path>` quando o clone se foi; e
`notes: <new untracked files in the clone>` quando a seção 3 achou algum; e
`leftover processes: stopped` ou `none`, com a string casada. Um
caminho ignorado alterado é reportado como `not revertible: <path>`.

O residual de clone-isolation.md seção 7 vale sem mudança: clone
descartável; escritas fora dele são detectadas na árvore viva, nunca impedidas;
escritas fora do projeto não são detectadas. Vale também para o diretório de
congelamento. Ele contém cópias do `<harness dir>`, de `.harness/`, de `.claude/settings*.json` e dos
arquivos modificados ou não rastreados e não ignorados, fica em `<tmp root>`, e uma
CLI que roda como o mesmo usuário pode lê-lo e alterá-lo (o modo 0700 não impede
isso). Adulterar as cópias cai sob o residual (a), uma escrita fora
do projeto, e ler as cópias cai sob (e), ler e devolver
segredos alcançáveis pelo usuário; o dono aceitou os dois em 2026-09-29. O
diretório de congelamento é removido ao fim da chamada, e varrido como órfão se
a chamada for interrompida (clone-isolation.md passo 5).

Residual do congelamento: arquivos de configuração trazidos por `include.path` ou `includeIf` a partir de `<C>/config` que vivem dentro do projeto não são restaurados antes do primeiro comando git da parte B; as sobreposições de `R5GIT` não os neutralizam. Residual aceito do congelamento da árvore viva.
