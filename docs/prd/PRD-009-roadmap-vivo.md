# PRD-009 · Roadmap vivo: deriva mecânica barrada por `dh validate`

| Campo | Valor |
|---|---|
| Status | aprovado (2026-09-30) |
| Dono | Daniel Malka |
| Criado / atualizado | 2026-09-30 / 2026-09-30 |
| Tickets | a derivar |

## 1. Problema

`docs/roadmap.html` e o espelho `docs/en/roadmap.html` (cerca de 1.660 linhas cada) são editados à mão, nas duas
línguas, em todo lote. Eles envelhecem sem que nada quebre. Casos já ocorridos, todos achados por um revisor e nunca
por uma checagem (`.harness/MEMORY.md`, lotes 0.7.0, 0.7.1, 0.8.0, 0.11.0 e 0.12.1): a chip de data atrasada, a linha
"Fontes" ainda dizendo "releases 0.1.0 a 0.7.1", cards de Entregue fora de ordem (corrigido só na 0.12.1), referências
cruzadas velhas e cards de débito deixados para trás. Quem paga é o dono e o Coordenador: o roadmap é a fonte que eles
consultam para saber o que foi entregue, e o `README.md` aponta para ele.

Fato: parte dessa deriva é mecânica, isto é, decidível comparando o roadmap com `CHANGELOG.md`, que é a fonte declarada
da seção Entregue. O `dh validate` já faz uma checagem desse tipo, `checkDocumentedCounts`
(`internal/kit/doc_counts.go`), que compara a linha "Contagem atual" dos dois roadmaps com o inventário real. Não
existe checagem da data, da versão, da linha Fontes, da ordem nem da cobertura dos releases, nem da paridade pt/en.
Outra parte é editorial: o texto de Testado? e Limite conhecido, as Relações e a decisão de mover um card de Ideias
para Planejado. Essa parte depende de julgamento e este PRD não a automatiza.

Decisão do dono (30/09/2026): opção A do discovery, isto é, roadmap continua escrito à mão e a deriva mecânica passa a
ser barrada por checagem determinística. Gerar o HTML a partir de dados (opção B) ou só as partes mecânicas (opção C)
ficam fora e podem virar um PRD futuro.

## 2. Solução

Depois da entrega, `dh validate` falha quando os dois roadmaps divergem de `CHANGELOG.md` ou entre si em fatos
mecânicos, e diz onde. O cabeçalho `## <versão> — <data>` mais alto de `CHANGELOG.md` é o canônico. A checagem vive em
`internal/kit/`, ao lado de `checkDocumentedCounts`, com testes em Go, e roda dentro do `validate` que a CI já executa.
O roadmap continua sem geração: ninguém passa a editar HTML por script.

A regra de manutenção passa a estar escrita onde o fluxo a lê: uma linha nomeada na tabela de prontidão do
`delivery-readiness` e uma no despacho do docs-guide ao fim do lote, no Coordenador. Os dois textos são ativos
genéricos do kit, distribuídos a todo projeto consumidor, então as duas linhas são condicionais e não citam caminho
do dev-harness como regra geral: valem "quando o projeto mantém um roadmap" e dizem que, sem roadmap, não se
aplicam. Os caminhos do dev-harness aparecem só como exemplo entre parênteses.

Títulos históricos: verificado nos dois arquivos (30/09/2026), os 22 cards de release da seção Entregue têm título no
formato `<versão> · <texto>` (versão `N.N.N`). Há também dois cards de detalhe sem versão em cada língua ("Comando de
planejamento em loop" e "Validador estendido a story, plano e ADR"; em en, "Loop planning command" e "Validator
extended to story, plan and ADR"). Não há título de release fora do padrão, então nenhuma lista de exceção é
necessária e nenhum card histórico é reescrito. Um card sem prefixo de versão é tratado como card de detalhe.

## 3. Regras

- R1 `dh validate` lê o cabeçalho `## <versão> — <data>` mais alto de `CHANGELOG.md` como versão e data canônicas.
  A checagem só é pulada quando os dois roadmaps (`docs/roadmap.html` e `docs/en/roadmap.html`) estão ausentes
  (pacote gerado, projeto consumidor); se só um dos dois existir, é erro que nomeia o arquivo ausente.
  `CHANGELOG.md` ausente com algum roadmap presente também é erro. Verificável: teste em Go para cada caso, inclusive
  o de um só roadmap presente (falha, nomeando o ausente) e o de nenhum presente (passa); `go run ./cmd/dh validate .`
  na árvore atual sai limpo.
- R2 A chip de data de cada roadmap (`<span class="chip">data <b>…</b>` em pt, `date` em en) é igual à data do R1. Se
  divergir, o erro cita o arquivo, a linha, o valor esperado e o encontrado. Verificável: teste com chip alterada; na
  árvore, trocar a data numa cópia faz `validate` falhar com essa mensagem.
- R3 A chip de status é `<versão> · vivo` em `docs/roadmap.html` e `<versão> · live` em `docs/en/roadmap.html`, com a
  versão do R1. Divergência gera erro no mesmo formato do R2. Verificável: idem R2.
- R4 O intervalo da linha Fontes (`<code>CHANGELOG.md</code> — releases 0.1.0 a X` em pt, `releases 0.1.0 to X` em en)
  termina na versão do R1. Divergência gera erro no mesmo formato. Verificável: idem R2.
- R5 Cada cabeçalho de release de `CHANGELOG.md` (`## <versão> — <data>`) tem exatamente um card na seção Entregue
  (do `<h2 id="entregue">` ao `<h2>` seguinte), em cada roadmap, cujo `<h3>` começa com `<versão> · `. Falta ou
  duplicata gera erro que nomeia a versão e o arquivo. Um card de versão sem cabeçalho correspondente em
  `CHANGELOG.md` também é erro. Verificável: testes com card faltando, duplicado e sobrando; na árvore atual passa
  (22 releases, 22 cards por língua).
- R6 Os cards de versão da seção Entregue aparecem do mais novo para o mais antigo (ordem semântica de `N.N.N`, não
  lexical: 0.10.0 vem antes de 0.9.0). Cards de detalhe sem prefixo de versão podem estar em qualquer posição e não
  quebram a ordem dos cards de versão. Desordem gera erro que nomeia os dois cards fora de ordem. Verificável: teste
  com 0.9.0 acima de 0.10.0 falha; teste com card de detalhe entre dois cards de versão passa.
- R7 Os dois roadmaps têm a mesma lista de versões de cards de versão, na mesma ordem, e as mesmas chips (data,
  versão do status e path) e o mesmo intervalo de Fontes, exceto o idioma da palavra "vivo/live". Diferença gera erro
  que mostra as duas listas. Verificável: teste com um card a menos no en falha.
- R8 Toda falha das checagens R2 a R7 é uma entrada em `report.Errors` no formato `<arquivo>:<linha>: <o que se
  esperava>, <o que se encontrou>`, e o `validate` sai com código diferente de zero. Para card faltando (R5), a linha
  citada é a do `<h2 id="entregue">`; para divergência entre pt e en (R7), o arquivo citado é
  `docs/en/roadmap.html`. Verificável: os testes afirmam a
  substring de cada mensagem; `go run ./cmd/dh validate .` retorna o código de saída.
- R9 A checagem não altera nenhum card histórico: o diff do lote não muda o texto de nenhum card de release existente
  em `docs/roadmap.html` nem em `docs/en/roadmap.html`. Verificável: `git diff` restrito às seções Entregue mostra só
  o card novo, se o lote entregar uma versão, e nada mais.
- R10 A regra de manutenção existe em dois lugares, ambos condicionais e neutros de caminho (ver §2):
  - Linha nomeada na tabela de prontidão do `.skills/delivery-readiness/SKILL.md`: "Roadmap | Met when the project
    keeps a roadmap tied to its changelog: it names the new release, its date and version, and passes the project's
    roadmap check, when it has one (in dev-harness: `docs/roadmap.html`, `docs/en/roadmap.html`, one card per
    release, `dh validate`). Not applicable when the
    project keeps no roadmap."
  - Frase no despacho de fim de lote do docs-guide em `.agents/coordinator.md`, nas duas ocorrências do texto
    existente (linhas 72 e 76 hoje), idênticas entre si: "When the project keeps a roadmap, the docs-guide dispatch
    includes updating it for the release (date, version, the new release recorded) and the project's roadmap check,
    when it has one, must pass; when the project keeps no roadmap, this does not apply."
  - `.agents/release-manager.md` não enumera checklist (o passo 4 só cita "checklist"), então não ganha linha; se o
    ticket achar uma enumeração, a linha condicional entra ali também. Nada é adicionado a `AGENTS.md` além da
    decisão já registrada.
  - Verificável: `grep -n "keeps no roadmap" .skills/delivery-readiness/SKILL.md` acha a linha da tabela;
    `grep -c "keeps no roadmap" .agents/coordinator.md` dá 2; `git diff AGENTS.md` do lote só tem a decisão do dono
    já gravada.
- R11 Fora do escopo, e o diff não introduz: HTML gerado, subcomando novo do `dh`, arquivo de dados do roadmap
  (opções B e C), checagem de movimento entre Planejado e Ideias (editorial), checagem do texto de Testado?, Limite
  conhecido ou Relações, OTEL e Agent SDK. Verificável: o diff não tem novo subcomando em `cmd/dh/main.go`, novo
  arquivo de dados nem gerador de HTML.
- R12 O lote fecha com `go test ./...`, `go run ./cmd/dh validate .` e `go run ./cmd/dh build` verdes, e `dist/`
  regenerado com a prova de clone limpo já exigida em `AGENTS.md`, porque `.skills/` e `.agents/` mudam. O texto do
  Coordenador é embutido em casos de eval (`checkEmbeddedAgentBodies`), então o re-embed é parte do lote; nenhuma
  eval paga roda, pois a mudança é de redação (revisão enxuta de `AGENTS.md`). Verificável: os três comandos, o
  `sha256` dos binários do clone limpo iguais aos da árvore e `validate` sem erro de eval embutida.

## 4. Docs

- docs/roadmap.html
- docs/en/roadmap.html
- CHANGELOG.md
- README.md
- README.en.md
- .skills/delivery-readiness/SKILL.md
- .agents/coordinator.md
- docs/tutoriais/ e docs/en/tutorials/ (só se algum tutorial descrever o que o `dh validate` checa; o ticket confirma por `grep`)

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| B. Gerar os dois HTML a partir de um arquivo de dados por um subcomando `dh` | Maior lote: novo esquema, novo subcomando e regra nova de "não editar o HTML à mão"; o texto dos cards é julgamento, não dado. Pode virar PRD futuro. |
| C. Gerar só as partes mecânicas entre marcadores | Exige marcadores no HTML e ainda deixa o texto editorial à mão; corta menos esforço que B e custa mais que A. |
| Só a regra no checklist, sem checagem | É a hipótese original do card; depende de alguém lembrar, e foi assim que a deriva aconteceu. |
| Checar o movimento Planejado/Ideias | Depende de haver ou não PRD, decisão editorial que um script não julga. |

### Riscos e decisões pendentes

- Risco: a checagem lê HTML por expressão regular e quebra se o formato do card ou da chip mudar · mitigação: os testes
  fixam o formato; uma mudança de estrutura do roadmap atualiza a checagem no mesmo PR.
- Risco: o `CHANGELOG.md` de um release ainda não escrito faz o roadmap parecer atrasado · mitigação: a ordem de
  trabalho do lote é `CHANGELOG.md` primeiro, e o `validate` só compara com o que está lá.
- Fato a confirmar pelo ticket: o texto do docs-guide no `.agents/coordinator.md` aparece em duas linhas (72 e 76); o
  ticket edita as duas de forma idêntica.
