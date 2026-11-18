# Comparação na árvore viva sem clone

Este arquivo é o procedimento de congelar, comparar e reverter para uma chamada de CLI que roda na árvore viva, sem clone descartável: `codex` sob a flag somente leitura dele, e a revisão do `agy`. Uma chamada em clone usa [live-tree-freeze.md](live-tree-freeze.md) em vez deste. A comparação não é um sandbox: ela só detecta uma alteração depois do fato. Um revisor é somente leitura por contrato mesmo quando o binário dele não tem flag somente leitura.

## Período de silêncio

Nada mais pode escrever na árvore nem no `<harness dir>` entre o congelamento e a comparação, Coordenador e delegado inclusive. Um build, uma regeneração de pacote ou uma execução de testes iniciada em paralelo muda arquivos que a comparação depois atribui ao revisor. Medido em 2026-09-22: um `dh build` rodado pelo Coordenador enquanto dois revisores de CLI trabalhavam fez a comparação divergir em `harness-manifest.json` e em `dist/`, o que, pela regra abaixo, descarta um veredito que o revisor nunca ameaçou. Espere, ou rode o trabalho paralelo num checkout separado.

## Congelamento (antes da chamada)

- `git status --porcelain -uall` para a lista de caminhos, depois `git hash-object "<path>"` para cada caminho listado (modificado ou não rastreado). Um caminho ausente no disco registra o marcador `deleted` no lugar do hash.
- Os bytes de todo arquivo sob o `<harness dir>` (modo `global`: `<home>/projects/<name>/`, fora da árvore, e o Coordenador também não escreve nada ali durante a chamada; modo `repo`: `.harness/`, ou o `.harness/` do checkout principal quando o projeto é um worktree vinculado), sob `.harness/` quando existir, e de todo arquivo sob `.git/hooks/`, de `.git/config`, e de `.claude/settings.json` e `.claude/settings.local.json` quando presentes. O git os ignora, então a lista de caminhos não vê uma alteração ali, e no entanto eles guardam os comandos, a lista de revisores e as permissões em que os passos posteriores confiam.
- Os itens de descoberta do repositório: os bytes do `gitfile` (ou o fato de `.git` ser um diretório de verdade) e, num worktree vinculado, os bytes de `<admin>/commondir` e `<admin>/gitdir`, mais a ausência registrada de `<C>/commondir` e `<project>/.git/commondir`.

## Comparar e reverter (depois da chamada)

Recalcule tudo o que está acima. Qualquer diferença é a classe de falha `tree_changed`: uma lista de caminhos alterada, um hash alterado para qualquer caminho, ou, no instantâneo de bytes, um arquivo alterado, um arquivo criado durante a chamada (de ausente a presente) ou um arquivo apagado. O veredito é descartado, a alteração é revertida e o revisor é `not-run` por esse motivo, seja qual for o que ele disse. Isso também pega uma edição feita dentro de um arquivo que já estava modificado no estado congelado, que só a lista de caminhos não consegue distinguir.

Reverta todo caminho divergente ao estado exato anterior à chamada:

- Um caminho que estava listado (modificado ou não rastreado) ou em instantâneo antes da chamada: o congelamento guarda uma cópia dos bytes e do modo anteriores à chamada fora do checkout, e a reversão escreve esses bytes de volta. Nunca use `git checkout` nem `git restore` para um caminho assim: isso o devolve ao `HEAD` e destrói alterações que já estavam na árvore antes da chamada, e o git não consegue restaurar um arquivo ignorado de jeito nenhum.
- Um caminho que estava limpo e rastreado antes da chamada (em nenhum dos dois conjuntos): `git checkout -- "<path>"` é a restauração correta, porque o `HEAD` é exatamente o estado dele anterior à chamada.
- Um caminho que não existia antes da chamada: apague-o.
- Nunca escreva através da entrada de diretório atual: remova o caminho primeiro (sem seguir um symlink) e recrie-o, para que um symlink plantado durante a chamada seja removido em vez de seguido.

## Worktree vinculado

- A revisão do `agy` na árvore viva termina em `isolation_unavailable` quando `<project>/.git` não é um diretório: a revisão não roda e nenhuma rodada de revisão é gasta, porque este congelamento tira instantâneo de `.git/hooks/` e `.git/config`, que não são pastas num checkout com `gitfile`. O gatilho é `.git` não ser um diretório, um superconjunto deliberado de um worktree vinculado (cobre também um submódulo e um repositório bare, que também não têm essas pastas).
- Para `codex` a partir de um worktree vinculado o procedimento é, no mais, inalterado, mas os caminhos de confiança (`hooks/`, `config`, `config.worktree` quando presente, `info/`, `HEAD`, `packed-refs`, e o `config.worktree` do diretório admin, todos salvos e restaurados com ferramentas que não sejam do git antes de qualquer comando git) são lidos do diretório comum `<C>` (`linked_worktree.commondir` de `dh harness-path --json "<project>"`), não do `gitfile` em `<project>/.git`. Um caminho ausente nos dois instantâneos não é uma diferença, e o sandbox somente leitura é a guarda.
- Antes de qualquer comando git da comparação ou da reversão, compare por bytes, com ferramentas que não sejam do git (`cat`, `cmp`, `stat`), os itens de descoberta do repositório salvos no congelamento (um item registrado como ausente deve continuar ausente). Restaure ou remova primeiro. Se isso for impossível, não rode nenhum comando git e registre o veredito como descartado, `not-run`, `not revertible`. Toda comparação de bytes e restauração desse tipo vem antes do primeiro comando git.
