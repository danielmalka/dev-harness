# Pegadinhas (só acréscimo)

Padrões de falha generalizados, observados ao invocar essas CLIs como revisores, agrupados pelas classes de falha de `SKILL.md`, "Classificação de falhas". Uma entrada por erro: o que aconteceu e o que fazer em vez disso. Acrescente uma entrada nova sob a sua classe; nunca reescreva nem reordene uma antiga. Nenhum caminho de máquina e nenhum nome de projeto consumidor cabe aqui; ver os metadados de `SKILL.md` para a procedência.

## not_installed

- `command -v <binary>` não encontra um binário que está instalado quando o diretório dele fica só no `PATH` de um shell de login. · Consulte `binaries.<binary>` em `.harness/local.yaml` antes de classificar; uma falha de `command -v` sozinha não é `not_installed`.

## auth_required

- Nenhuma registrada.

## rate_limited

- Nenhuma registrada.

## `timed_out`

- Uma chamada `agy --print` escreveu a resposta completa em cerca de 2 minutos, mas o processo continuou vivo além de `--print-timeout 15m` e precisou de SIGKILL (medido em 2026-10-16). · Uma resposta completa em disco que traga o sinal de veredito da etapa é o resultado; pare o processo pelo pid registrado no início da chamada; ver [agy.md](agy.md).
- Um diff grande enviado com `--prompt-file` em esforço de raciocínio alto estourou o timeout. · É uma falha, nunca uma aprovação complacente.

## no_signal

- Uma chamada de revisão de CLI saiu com 0 depois de imprimir só um preâmbulo de raciocínio, sem bloco de veredito, num diff grande rodado em esforço de raciocínio alto. · Saída 0 sem sinal de veredito é uma falha exatamente como uma saída diferente de zero ou um timeout; nunca é lida como aprovação implícita.
- Uma execução real de revisor devolveu um veredito bloqueante correto escrito como `- **Request changes**` e nenhuma linha `Status da revisão:` (2026-09-22). · Declare o sinal de veredito no prompt como exigência, literalmente, como diz o contrato de veredito em `SKILL.md`.
- As mesmas execuções de revisor assinaram os achados com tags que não eram o slug despachado. · Diga ao revisor, no prompt, a string exata de `Reportado por:`.
- `agy` imprimiu a ajuda de `--add-dir` e ignorou o pacote quando `--add-dir PATH` veio antes do prompt. · Embuta os arquivos no prompt; ver [agy.md](agy.md).
- `codex exec` fora de um repositório git confiável saiu com 0 e stdout vazio. · Passe `--skip-git-repo-check`; ver [codex.md](codex.md).
- Duas chamadas simultâneas ao mesmo binário de CLI, da mesma conta, fizeram a segunda falhar com um erro de trava de banco de dados. · Serialize as chamadas a um binário sabidamente intolerante a execuções simultâneas, cada uma em seu próprio lote; ver na receita do binário a quais binários isso se aplica.

## tree_changed

- Uma CLI a quem se deu um resumo explicitamente somente leitura e só de validação rodou por conta própria um script de build do repositório e reescreveu um diretório de saída versionado (o resultado reconstruído por acaso coincidia com as fontes, então não houve dano, mas a execução não foi somente leitura na prática). · Um resumo somente leitura deve proibir rodar scripts do repositório, não só editar arquivos; confira os horários de modificação de qualquer diretório de saída versionado depois da chamada, não só o `git status` dos arquivos rastreados.
- Um `dh build` rodado pelo Coordenador enquanto dois revisores de CLI trabalhavam fez a comparação posterior à chamada divergir no manifesto e em `dist/` (2026-09-22). · Nada mais escreve na árvore enquanto uma chamada está em andamento; ver [live-tree-compare.md](live-tree-compare.md).

## isolation_unavailable

- Nenhuma registrada.

## Sem classe (a chamada rodou; o prompt ou o agendamento estava errado)

- Uma CLI rodando como implementadora sob uma flag de sandbox de acesso total invocou por conta própria um segundo binário de CLI, dentro do mesmo checkout, delegando parte do trabalho sem ser solicitada. · A cláusula anti-delegação deve ser copiada literalmente para o corpo de todo prompt de revisor, nunca deixada só na instrução de despacho em volta (RISK-001); quando houver suspeita de delegação, confira os processos em execução durante a chamada.
- `codex exec -s read-only` impede um `commands.test` da cadeia de ferramentas Go de compilar de fato. · Trate a aprovação desse revisor como leitura e raciocínio sobre o código, não como validação de comportamento; ver [codex.md](codex.md).
- `pgrep -f` ou `pkill -f` com um padrão que contém o comando ou o caminho do próprio delegado casou com o shell do próprio delegado (saída 144, esperas sem fim; visto em 2026-10-21). · Mate pelo pid registrado quando a chamada começou, ou use um padrão com colchetes que não case com o próprio texto (`'[o]pencode run'`).
