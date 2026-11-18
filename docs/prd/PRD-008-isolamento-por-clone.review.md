# Revisão · PRD-008 · rodada 1

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

- Painel (reviewers.document; o PRD toca em segurança): claude (document-validator, opus solicitado, efetivo não verificado) — `changes required`. cli:grok/grok-4.7 e cli:codex/gpt-5.6-luna — `not-run`: o Coordenador gravou na árvore (migração do AGENTS.md e registros de .harness) enquanto as chamadas deles estavam em andamento, o que invalida a comparação do congelamento; o grok foi interrompido no meio da chamada e o codex nunca começou. O diff de congelamento posterior à parada mostra só as quatro gravações conhecidas do Coordenador. Nenhuma rodada de correção é gasta com eles; os dois rodam na rodada 2.
- Data: 2026-09-29.

## Veredito
- `changes required`
- Achados bloqueantes: 5

## Achados
- P1 · conflict · blocking · §2, Apêndice (D1, primeiro risco). O residual é subestimado como ".git/hooks e .git/config". O diretório comum compartilhado também guarda `refs/`, `packed-refs`, `objects/`, `info/exclude` (graváveis de dentro do worktree sem caminhos absolutos, por exemplo `git update-ref`, acrescentando a info/exclude, o que cega a lista de caminhos da árvore viva); e a árvore viva continua alcançável por caminho absoluto (o arquivo .git do worktree e o `git worktree list` a revelam), de modo que o vetor da rodada 9 do RISK-002 (.gitignore autoignorado + _test.go na árvore viva) continua fora do R8. Sugestão: declarar o residual por inteiro, dizer quais itens o R8 detecta e quais não, e ou estender o R8 (`for-each-ref`, `packed-refs`, `info/` do diretório comum, `--ignored` na lista da árvore viva) ou registrar cada item não detectado como residual aceito sob a D1. Reportado por: claude
- P2 · conflict · blocking · §3 R14, R15. A avaliação grava dentro do worktree, onde o R3 garante que .harness/ não existe — ela passa por construção e nunca exercita a ameaça do RISK-002; a ação ".git/hooks" não está especificada (a escrita relativa falha com ENOTDIR); nenhum caso grava na árvore viva por caminho absoluto. Sugestão: acrescentar ações de fixture (i) escrita absoluta em .harness/project.yaml da árvore viva, (ii) escrita em $(git rev-parse --git-common-dir)/hooks/<name>, (iii) plantar .gitignore + _test.go na árvore viva, com o veredito esperado para cada uma. Reportado por: claude
- P3 · ambiguity · blocking · §3 R7 (com R12, R14). Não está definido se a comparação dentro do worktree inclui caminhos ignorados: a leitura A aceita arquivos ignorados plantados, a leitura B descarta toda chamada que rodou commands.test (R12). Sugestão: nomear a comparação exata e como os artefatos de teste são tratados; declarar o veredito e o status `not-run` por ação de fixture. Reportado por: claude
- P4 · conflict · blocking · §2 ("isolamento verificado"), R10, R4. Afirma-se prevenção onde só existe detecção; os textos que o R10 emenda justificam os planejadores por "deve ser incapaz de gravar" e "o congelamento é um detector, não uma sandbox", e o R10 não diz como essa justificativa muda; o R4 é higiene, não um controle. Sugestão: reescrever o §2; fazer o R10 exigir que os três textos declarem a distinção de só detectar; marcar o R4 como higiene. Reportado por: claude
- P5 · ambiguity · blocking · §3 R2, R5, R6, R13, R8. "caminho do kit" e "órfão" não estão definidos; o despacho paralelo B poderia podar o worktree vivo de A ou um worktree do dono; `git worktree lock` não é usado; o R8 pode sinalizar o próprio worktree do kit. Sugestão: nomear o prefixo `dh-cli-<call-id>` mais `git worktree lock --reason <call-id>`; órfão = com prefixo e destravado por uma chamada viva; nunca tocar em worktrees sem prefixo; o R8 exclui o worktree da própria chamada ou roda depois do R5. Reportado por: claude
- P6 · organization · non-blocking · D2 (metade do codex), D6 assentada pelo R16 e por .commands/plan-loop.md; D1 vs R15 provisório; D5 é um resultado de sonda. Reportado por: claude
- P7 · organization · non-blocking · H1 confirmada pela checagem do Coordenador — mover para Fatos. Reportado por: claude

## Não levantado
- Enunciado do problema, rastreabilidade do objetivo, status somente leitura por binário, F1/F4/F5, contagem do R10, preservação do contrato do R11, limite do R13, regra de custo do R14, falha fechada do R3/R9, fronteiras do R16, lista e formato de Docs: tudo rastreado e consistente.

---

# Revisão · PRD-008 · rodada 2 (final)

- Painel: claude — `changes required` (P1–P7 aplicados; novos P8, P9). cli:grok/grok-4.7 (desacoplado, congelamento limpo; veredito recuperado, linha sem "- ") — `changes required`. cli:codex/gpt-5.6-luna (congelamento limpo; veredito recuperado) — `changes required`. Saídas brutas das CLIs: .harness/tasks/PRD-008/evidence/`{grok,codex}-r2.out`. Limite de rodadas atingido (2 de 2).

## Veredito
- `changes required`

## Achados (mesclados; pontos equivalentes compartilham uma entrada)
### Metadados .git vivos fora do instantâneo (.git/index, logs/, worktrees/) — skip-worktree e fsmonitor global podem esconder edições de arquivos rastreados
- Reportado por: claude — "a CLI edita um arquivo rastreado limpo ... depois roda `git -C <live> update-index --skip-worktree` ... Um `core.fsmonitor = <command>` em `~/.gitconfig` é executado pelo `git status` posterior à chamada."
- Reportado por: cli:codex/gpt-5.6-luna — "Metadados vivos graváveis como `.git/index`, `.git/logs/`, `.git/worktrees/` ... não estão listados como detectados nem como residual aceito."
### Aceitação pelo dono dos residuais (a), (b), (e) atribuída sem registro
- Reportado por: claude
### R2 "diretório temporário do sistema" vs R9 "diretório temporário do kit" não definidos como a mesma raiz
- Reportado por: cli:grok/grok-4.7; cli:codex/gpt-5.6-luna
### O marcador de trava na raiz do clone quebra a correspondência do R3 (deveria ficar em `<clone>/.git/`)
- Reportado por: cli:grok/grok-4.7
### A correspondência do R3 compara com a fotografia `--ignored` do R5 (ambíguo; deve usar só o conjunto julgado)
- Reportado por: cli:grok/grok-4.7
### agy sem a sonda de cwd: fica na árvore viva para revisão, não planeja (R1 vs redação do §2)
- Reportado por: cli:grok/grok-4.7
### O §1 afirma detecção de qualquer coisa fora do clone; só a árvore viva é detectada
- Reportado por: cli:grok/grok-4.7
### Residual (b) "git push não detectado" vs R5 detectando movimentos de ref viva
- Reportado por: cli:grok/grok-4.7
### ADR-002 §3 "nunca grava em disco / qualquer mudança descarta" não coberto pelas emendas do R12/R14
- Reportado por: cli:grok/grok-4.7
### Um _test.go ignorado plantado no clone pode influenciar a própria execução de `commands.test` do revisor
- Reportado por: cli:codex/gpt-5.6-luna
### (não bloqueante) D1 marcada como "bloqueia R6" embora o R6 já fixe "só nota"
- Reportado por: cli:grok/grok-4.7

## Disposição do dono (2026-09-29)
- O dono aceitou os residuais (a) escritas fora do projeto, (b) processos que sobrevivem / rede / push para remotos que não são o vivo, (e) leitura e devolução de segredos, até o lote de sandbox do SO; e aceitou a recomendação do Coordenador de fechar a brecha do índice/fsmonitor. Todos os outros achados são mecânicos e são aplicados em uma correção final, conferida pelo Coordenador (sem terceira rodada de validação).
