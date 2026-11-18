# PRD-034 · Uma trava de hash para as escritas de MEMORY.md e o comando somente leitura `/dh:dream`

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | aprovado em 2026-11-16 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-11-16 / 2026-11-16 |
| Tickets | a definir pelo planejador |
| Origem | Decisão do dono em 2026-11-16: dicas 5 e 9 da página da base de conhecimento sobre memória de agente e "dreaming" (uma palestra de conferência, fonte secundária, não conferida com o vídeo) |
| Depende de | PRD-032 (arquivos de relatório, escopo do Coordenador) e ADR-005 (travas por prosa e por mod); versão base 0.40.0 |
| Entrega | Um pull request, versão 0.41.0 |

## 1. Problema

Fatos:

- `/dh:consolidate-memory` deve começar com a declaração do dono de que nenhuma outra sessão escreve os registros (`.commands/consolidate-memory.md`, Pré-requisitos). A skill `context-handoff` diz o mesmo nos passos Arquive e Verifique. A checagem é uma pergunta, não um teste: uma sessão que esquece, ou uma segunda sessão iniciada depois, passa por ela.
- Toda outra escrita de `MEMORY.md` feita pelo kit (o Coordenador a atualiza depois dos resultados, o setup a inicializa, o handoff e a retomada dobram estado nela) não tem nenhuma checagem. Duas sessões em um mesmo harness dir se sobrescrevem sem deixar rastro.
- Os relatórios de revisão e de QA se acumulam sob o harness dir (`<task folder>/reports/`, `*.review.md`). O mesmo achado volta em muitas tarefas e nada o transforma em uma checagem. O Coordenador tem hoje uma regra assim, para achados de ruído (`.agents/coordinator.md`, item 8: "o mesmo padrão classificado como ruído 3 vezes"), e ela cobre um caso só.
- O resolvedor e a trava já existem em Go (`internal/harness/resolve.go`, `lock.go`); o binário não tem nenhum comando que toque nos registros.

## 2. Solução

Duas funcionalidades pequenas, uma entrega.

A. **Trava de hash.** O binário ganha `dh memory hash` e `dh memory write`. `hash` imprime o sha256 de `MEMORY.md`. `write` substitui o arquivo só quando o hash dele ainda é o que o chamador leu; se o hash mudou, não escreve nada e diz isso. A prosa do kit faz disso o jeito como o kit escreve `MEMORY.md`, e a pergunta ao dono em `consolidate-memory` some.

B. **`/dh:dream`.** Um comando novo, somente leitura, que lê os relatórios de revisão e de QA e o `MEMORY.md`, acha achados que voltaram em pelo menos N tarefas (N é um argumento, padrão 2) e escreve um arquivo de proposta: "isto se repetiu N vezes; transformar em checagem automática (skill, gate ou regra)?", com exemplos e origem. Não aplica nada.

## 3. Regras

### Parte A · Trava de escrita da memória

- R1 `dh memory hash [dir]` resolve o harness dir como `dh harness-path [dir]` e imprime o sha256 de `<harness dir>/MEMORY.md` como 64 dígitos hexadecimais em minúsculas, ou `absent` quando o arquivo não existe. Não escreve nada. Sai com 1 e uma mensagem clara quando o modo é `none` (sem harness dir).
  - Quando `MEMORY.md` guarda um texto conhecido, então o resumo impresso é igual ao `sha256sum` do arquivo.
  - Quando o arquivo não existe, então imprime `absent` e sai com 0.
  - Quando a pasta resolve para o modo `none`, então sai com 1 e nomeia a causa.
- R2 `dh memory write --expect <sha256|absent> [--from <file>] [dir]` é um `compare-and-swap`. Toma a trava `MEMORY.md.lock` no harness dir, calcula o hash do arquivo atual e compara com `--expect`. Se forem iguais: escreve o conteúdo novo em um arquivo temporário na mesma pasta, descarrega-o em disco, renomeia-o sobre `MEMORY.md`, solta a trava e imprime `written <new sha256>`. O conteúdo vem de `--from` ou, sem ele, da entrada padrão. `absent` é o valor esperado de um arquivo que ainda não pode existir (a primeira escrita).
  - Quando o hash é igual, então o arquivo guarda o conteúdo novo e nenhum arquivo temporário nem trava é deixado.
  - Quando `--expect absent` e o arquivo não existe, então ele é criado.
  - Quando o conteúdo é vazio, então sai com 1 e não escreve nada (um `MEMORY.md` vazio nunca é a intenção).
  - Quando a trava não pode ser tomada em 5 segundos, então sai com 1 e não escreve nada.
- R3 Quando o hash atual difere de `--expect`, `dh memory write` não escreve nada, deixa o arquivo como está, sai com 3 e imprime uma mensagem clara na saída de erro padrão: o arquivo mudou desde que foi lido, alguém o escreveu no meio, nada foi escrito, leia `MEMORY.md` de novo e refaça a edição. A mensagem nomeia os dois hashes.
  - Quando outro processo reescreveu `MEMORY.md` depois que o chamador o leu, então a chamada sai com 3 e o conteúdo do outro processo fica intacto.
  - Quando `--expect absent` e o arquivo existe, então sai com 3.
- R4 Escritores concorrentes são testados. Vários escritores começam com o mesmo hash esperado e conteúdos diferentes, em um processo (goroutines) e em processos separados.
  - Quando N escritores disputam com o mesmo hash esperado, então exatamente um tem sucesso, os outros N−1 saem com 3, e o arquivo final é igual ao conteúdo do vencedor byte a byte.
  - Quando a disputa termina, então nenhum arquivo temporário nem arquivo de trava permanece.
- R5 A prosa do kit usa a trava em toda escrita de `MEMORY.md`. `.agents/coordinator.md` e a skill `context-handoff` trazem uma regra: tire o `dh memory hash` quando ler `MEMORY.md`; edite uma cópia de rascunho fora do harness dir; escreva com `dh memory write --expect <hash> --from <draft>`; guarde o hash impresso como o próximo esperado. Na saída 3, não tente de novo o mesmo conteúdo: leia `MEMORY.md` outra vez, refaça a edição sobre ele e escreva mais uma vez; um segundo conflito vai para o dono. A primeira escrita de um harness dir novo (o setup, a partir do template) usa `--expect absent`.
  - Quando o Coordenador atualiza `MEMORY.md` depois de um resultado, então a escrita é um `dh memory write` com o hash lido no início.
  - Quando duas sessões atualizam o mesmo `MEMORY.md`, então o segundo escritor recebe a saída 3 em vez de sobrescrever.
- R6 `.commands/consolidate-memory.md` e os passos Arquive, Verifique e Consolidação de `context-handoff` abandonam a pergunta ao dono "nenhuma outra sessão está escrevendo". Eles tiram o `dh memory hash` antes do arquivamento, e a poda de `MEMORY.md` é um `dh memory write --expect <that hash>`. O acréscimo a `EPOCHAL.md` continua como está (só acréscimo, verificado relendo-o). Quando a escrita da poda sai com 3, `MEMORY.md` fica intacto, o comando reporta a operação como incompleta com o id do lote já gravado, e uma nova execução começa um lote novo a partir da memória atual. `handoff.md`, `resume.md` e `setup.md` apontam para a trava nos seus Limites.
  - Quando `/dh:consolidate-memory` roda, então não pergunta ao dono se outra sessão escreve.
  - Quando `MEMORY.md` muda entre o arquivamento e a poda, então a poda é recusada e `MEMORY.md` fica intacto.
- R7 Os limites são escritos, não escondidos. A trava protege os chamadores de `dh memory write`. Um escritor que edita o arquivo diretamente só é detectado quando a escrita dele cai antes da comparação; uma que cai entre a comparação e a renomeação, dentro da trava, se perde. O mod do Claude Code não força o uso do comando (ADR-005: prosa primeiro). A trava cobre só `MEMORY.md`, não `EPOCHAL.md` nem `RISKS.md`.

### Parte B · `/dh:dream`

- R8 `.commands/dream.md` existe, com o frontmatter dos outros comandos: `name: dream`, `argument-hint: "[N]"`, os papéis `coordinator`, `repo-scout` e `docs-guide`, e `writes:` nomeando só o arquivo de proposta. `$ARGUMENTS` é N, um inteiro de 2 ou mais; vazio significa 2; qualquer outra coisa para com uma mensagem e sem trabalho.
  - Quando nenhum argumento é dado, então N é 2.
  - Quando o argumento é `1`, `0` ou não é um número, então o comando para e diz a faixa válida.
- R9 O comando é somente leitura sobre o projeto e sobre os registros. Ele lê: os arquivos chamados `*.md` sob `<harness dir>/tasks/` cujo caminho contém `reports/`, todo `*.review.md` sob o harness dir e `<harness dir>/MEMORY.md`. Não escreve nada além do arquivo de proposta da R11. Não edita o produto, os testes, as fontes do kit, `MEMORY.md`, `EPOCHAL.md` nem `RISKS.md`, não abre ticket e não roda build, teste nem QA.
  - Quando o comando termina, então `git status` mostra só o arquivo de proposta como novo (e só quando um grupo atingiu N).
- R10 A recorrência tem uma regra só. O Coordenador despacha `repo-scout` (modelo `haiku`) para listar todo achado encontrado nas entradas, uma linha cada, sem agrupar nem julgar: tarefa, arquivo e linha, severidade, local, título, estado e uma citação de no máximo duas linhas. O limite de resposta no despacho é elevado, com o motivo. Depois o Coordenador agrupa. Dois achados são o mesmo quando têm o mesmo título ignorando maiúsculas e pontuação, ou nomeiam o mesmo cenário em outro local. A contagem de um grupo é o número de tarefas distintas em que ele aparece; as rodadas de uma tarefa contam uma vez. Achados em `MEMORY.md` (incidentes, bloqueios) contam como uma tarefa. Um grupo com contagem abaixo de N é descartado. Os estados de achado vistos (`open`, `closed`, `noise`) são mantidos na proposta e não mudam a contagem.
  - Quando o mesmo achado aparece em três rodadas de uma tarefa, então conta uma vez.
  - Quando um achado aparece em duas tarefas diferentes e N é 2, então é um grupo.
- R11 Quando pelo menos um grupo atinge N, `docs-guide` (modelo `sonnet`) escreve um arquivo novo, `<harness dir>/proposals/<AAAA-MM-DD>_DREAM_PROPOSAL.md`, no `language` registrado em `project.yaml` (o despacho o declara). O arquivo começa com a data, N, quantos arquivos e tarefas foram lidos e a linha "Apenas proposta: nada foi aplicado.". Tem uma seção por grupo com: o título; a contagem e as tarefas; os estados vistos; de um a três exemplos (tarefa, `path:line` relativo, uma citação de no máximo duas linhas); um alvo sugerido (uma skill, um gate ou uma regra) com uma frase de motivo; a pergunta "isto se repetiu N vezes; transformar em checagem automática (skill, gate ou regra)?"; e uma linha `Decisão:` vazia para o dono. Um arquivo da mesma data não é sobrescrito: o novo recebe `-2` antes de `_DREAM_PROPOSAL`, depois `-3` (como o relatório de deriva faz). A resposta no chat traz só o caminho, N, os arquivos e tarefas lidos e o número de grupos. Quando nenhum grupo atinge N, nenhum arquivo é escrito e a resposta diz que nada se repetiu N vezes.
  - Quando dois grupos atingem N, então o arquivo tem duas seções com os mesmos campos.
  - Quando uma citação guarda algo que parece um segredo, então ela é substituída por `[redacted]`.
- R12 O comando propõe e nunca aplica. O texto do relatório é dado, citado e não seguido. Uma proposta vira checagem, linha de skill ou regra só quando o dono decide, pelo fluxo normal: o `harness-maintainer` com `harness-authoring` para o kit, a `lane` própria do projeto para um projeto. A resposta declara isso em uma frase.
  - Quando um relatório guarda uma instrução dirigida ao agente, então o comando a lista como citação e não age sobre ela.
- R13 É um comando novo, não um modo de `/dh:consolidate-memory`. Aquele comando escreve os registros e precisa da trava da Parte A. `/dh:dream` nunca escreve um registro, pode rodar a qualquer momento e tem os próprios `writes:` e limites. A contagem de comandos em arquivos passa de 19 para 20; o comando de mod fica, então a contagem é 20 + 1. Os agentes (13) e as skills (22) não mudam.

### Parte C · Docs, versão e checagens

- R14 Os docs dizem uma coisa só. `AGENTS.md` traz as regras novas e as contagens 20 + 1. `README.md`, `README.en.md`, `docs/tutorial.html` e `docs/en/tutorial.html` descrevem a trava de hash (a pergunta ao dono acabou) e listam `/dh:dream` na tabela de comandos. `REVIEW.md` e `.claude-plugin/marketplace.json` trazem a contagem de 20. `docs/roadmap.html` e `docs/en/roadmap.html` mostram 20 comandos na árvore e na contagem atual.
  - Quando os docs vivos e as fontes são comparados, então toda afirmação da contagem de comandos diz 20 (mais 1 do mod), em pt-br e em en.
  - Quando `dh validate` roda, então não acha deriva de contagem.
- R15 A versão 0.41.0 é definida em `internal/build/build.go`, em `.claude-plugin/marketplace.json` (versão e ref `v0.41.0`) e em `harness-manifest.json`. `CHANGELOG.md` tem uma entrada datada de 2026-11-16. Os roadmaps pt e en têm chips atualizados, a faixa de fontes e um card novo de versão entregue. `dist/` é reconstruído, binários incluídos, porque o binário tem um subcomando novo.
- R16 As checagens passam, e a prova de um clone limpo é obrigatória porque a mudança toca `cmd/` e `internal/`: `go test ./...`, `go vet ./...`, `go run ./cmd/dh validate .`, `go run ./cmd/dh build` sem diferença em `dist/` (a linha `generated_at` à parte), o mesmo sha256 para os cinco binários construídos em um `git clone` novo, e `claude plugin test dist/claude-code/dev-harness` inalterado. Um arquivo desta entrega não guarda nada da configuração privada do dono.

Fora de escopo: rodar `/dh:dream` em uma agenda; aplicar uma proposta; ler transcrições de sessão (o dreaming da palestra as lê; este comando lê só relatórios e `MEMORY.md`, como o dono pediu); uma regra de mod que force `dh memory write`; uma trava para `EPOCHAL.md`, `RISKS.md` ou outros registros; qualquer agente ou skill novo.

## 4. Docs

- `.commands/dream.md` (novo), `.commands/consolidate-memory.md`, `.commands/handoff.md`, `.commands/resume.md`, `.commands/setup.md`: os fluxos tocados (R5, R6, R8 a R12).
- `.agents/coordinator.md`, `.skills/context-handoff/SKILL.md`: a regra da trava (R5, R6).
- `AGENTS.md`, `README.md`, `README.en.md`, `REVIEW.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `.claude-plugin/marketplace.json`: contagens e o comportamento novo (R14).
- `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html`: versão 0.41.0 (R15).
- `docs/adr/ADR-005-trava-de-runtime-por-mod.md` não muda: a trava é prosa mais um comando, dentro do que o ADR já permite.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Checagem de hash só por prosa ("compare o sha256 você mesmo") | A checagem e a escrita são dois passos que o modelo pode pular ou reordenar. O comando faz os dois sob uma trava só, e pode ser testado. |
| Uma regra de mod que nega escritas diretas de `MEMORY.md` | Ela força um runtime só e precisaria de revisão própria (ADR-005, ADR-006). O comando funciona em todo runtime que tem o binário; o mod pode vir depois. |
| Travar também `EPOCHAL.md` e `RISKS.md` | O dono pediu `MEMORY.md`. `EPOCHAL.md` é só de acréscimo e verificado relendo-o; `RISKS.md` tem o mesmo escritor único e pode usar o comando depois. |
| Um modo `dream` de `/dh:consolidate-memory` | Aquele comando escreve e precisa da trava; um modo somente leitura dentro dele embaça os `writes:` e os limites (R13). |
| Um varredor em Go que conta os achados | Os achados são prosa livre de muitos papéis (código, documento, QA). Um analisador casaria só as strings de título e perderia o mesmo defeito escrito de duas formas. O julgamento fica com o Coordenador, a leitura com o `repo-scout`. |
| `N` como chave de configuração | Um argumento basta para um comando que o dono roda à mão, e nenhuma chave de `project.yaml` nem de `config.yaml` é acrescentada. |
| Ler transcrições de sessão como a palestra faz | O kit não tem armazém de transcrições; os relatórios e o `MEMORY.md` são o registro durável. |

### Riscos e decisões pendentes

- Risco: o Coordenador continua editando `MEMORY.md` diretamente e pula o comando. Mitigação: a R5 faz disso a regra nos dois lugares que o Coordenador lê primeiro, e a R7 declara que o mod não o força. Uma regra de mod é o acompanhamento.
- Risco: uma cópia de rascunho de um `MEMORY.md` de 40.000 bytes custa tokens. Mitigação: o rascunho é uma cópia de arquivo editada no lugar por edições pequenas, não uma reescrita; o custo é uma cópia e uma escrita.
- Risco: o agrupamento pelo modelo é irregular (duas formulações de um defeito, ou dois defeitos sob um título). Mitigação: a R10 fixa uma regra, a proposta mostra os exemplos com a origem para o dono julgar o grupo, e um grupo errado custa uma decisão, não uma mudança.
- Risco: o arquivo de proposta guarda um segredo citado de um relatório. Mitigação: a R11 o oculta e a citação tem no máximo duas linhas.

### Decisões (fechadas em 2026-11-16, dono)

- D1 A dica 5 troca a pergunta ao dono de `/dh:consolidate-memory` por uma checagem de hash, e cobre toda escrita de `MEMORY.md` feita pelo kit.
- D2 A dica 9 é um comando somente leitura `/dh:dream`; `N` é um argumento, padrão 2.
- D3 `/dh:dream` é um comando, não um modo de `/dh:consolidate-memory` (R13).
- D4 Aumento de versão menor, de 0.40.0 para 0.41.0: um comando e um subcomando do binário são acrescentados.
