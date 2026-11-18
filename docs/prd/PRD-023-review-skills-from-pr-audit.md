# PRD-023 · Adicionar checagens de revisão e de teste de duas skills de auditoria externas

> Traduzido em 2026-11-17; o original em inglês está no commit `aa1e80a`.

| Campo | Valor |
|---|---|
| Status | entregue em 2026-10-16 |
| Responsável | Daniel Malka |
| Criado / atualizado | 2026-10-16 / 2026-10-16 |
| Tickets | T-2301 (planejado) |
| Origem | Conversa com o dono, 2026-10-16 (escopo decidido na conversa); emendado no mesmo dia com os itens do `security-audit` |
| Depende de | nada |
| Habilita | nada |

## 1. Problema

- Duas skills externas, `pr-audit` e `security-audit` (akitaonrails/my-skills, commit 7b8a5e5), encontram defeitos que as skills de revisão do kit deixam passar. Fato: o Coordenador as comparou com o kit.
- Lacunas em `code-review`:
  - Ela não manda o revisor procurar pontos de chamada irmãos de um bug corrigido. A ideia existe só em `.skills/systematic-debugging/SKILL.md`, passo 8.
  - Ela não manda o revisor achar todos os consumidores de um formato de dados alterado, inclusive os espelhos que as pessoas mantêm à mão.
  - Ela checa se um arquivo adicionado é seguro. Não checa se o arquivo pertence à árvore.
  - Ela não exige que um novo ponto de entrada prove de novo os invariantes dos antigos.
- Lacunas em `security-review`:
  - Ela não tem gate para um diff hostil nem para a cadeia de suprimentos.
  - Ela não tem regra para uma correção de segurança nem para a prova de um relato de segurança.
  - Ela aceita um resultado limpo de scanner mesmo quando o diff alterou o scanner.
  - Ela não tem linha para um controle desligado, uma fronteira de modelo, um segredo vazado nem uma checagem de identidade com falha fechada.
- Lacuna em `regression-testing`: o passo 6 aceita `"fails before the fix"` mesmo quando a falha é um erro de compilação ou de import.

## 2. Solução

Adicione cada ideia a uma skill existente, em poucas linhas curtas. Não adicione skill, agente nem comando. Os exemplos no texto das skills entregues são genéricos. Os nomes de arquivo deste repositório ficam só neste PRD.

`code-review`:

1. Passos 4 e 5: quando uma mudança corrige um sintoma, o revisor procura pontos de chamada irmãos do mesmo helper, guarda ou normalização. Uma correção pontual que deixa quebrado um ponto gêmeo é Major, a menos que o escopo seja justificado.
2. Item `"Contracts"`: quando um formato de dados compartilhado muda, o revisor lista todos os consumidores. Isso inclui os espelhos que as pessoas mantêm à mão, como docs traduzidos, uma fixture e seu teste de paridade, e contagens declaradas. Um espelho que a mudança deveria ter atualizado é um achado. Neste repositório, os exemplos são os espelhos de docs pt/en, `status-cases.fixture.ts` com `fixture_parity_test.go`, e as contagens declaradas.
3. O revisor checa se cada arquivo adicionado pertence à árvore. Uma captura de tela solta, uma saída de rascunho, uma estrutura auxiliar do revisor ou um binário solto é um achado. Isso segue a premissa de limpeza do `AGENTS.md`.
4. Um novo ponto de entrada (comando, mod, adaptador, rota, superfície de CLI) deve provar de novo os invariantes transversais dos antigos. Exemplos: a trava de escrita nos arquivos de memória, a resolução do harness dir, o tratamento de caminhos. Um novo ponto de entrada não herda as correções antigas.

`security-review` (`SKILL.md` e `references/checklist.md`):

5. Um gate de inventário do diff para uma mudança hostil e para a cadeia de suprimentos.
   - Sinais de arquivo e de conteúdo (uma nova linha do checklist): bits de execução, symlinks, submódulos, caracteres bidi Unicode e homóglifos, blobs codificados ou binários sem explicação. Mudanças em workflows de CI: `pull_request_target`, `permissions`, actions não fixadas, segredos que chegam a código não confiável.
   - Os sinais de dependência estendem as linhas existentes do checklist `"Dependencies"` e `"Configuration"` (`checklist.md`, linhas 42 e 43). Eles não adicionam linha duplicada: deriva do lockfile, novos registries ou dependências git, hooks de ciclo de vida de pacote.
   - Uma correção de segurança é o diff de maior risco. O revisor procura mudanças que vão de carona e guardas que o diff enfraquece em outro lugar.
   - Um relato de segurança é provado ou refutado antes de uma correção. A prova reproduz o problema com código confiável. Ela nunca executa o payload de quem reportou.
   - Um refutador independente de um relato ou de uma correção é um despacho do Coordenador: uma segunda execução do revisor, ou um revisor de CLI configurado. Ele recebe só o achado e a evidência do código, nunca o raciocínio do autor. O revisor nunca cria agentes.
7. Um resultado de scanner não é aprovação quando o diff alterou ou estreitou a configuração do scanner, seus ignores, sua linha de base ou seu limiar. Também não é aprovação quando o scanner não barra o merge. Uma entrada de supressão ou de ignore adia o achado. Ela não o resolve. Para dispensar um alerta, o revisor registra um motivo que nomeia a fronteira que o torna seguro. Local: passo 1 do `SKILL.md`, ao lado de `"scanner output is evidence"`, e as linhas Dependencies e Configuration.
8. Uma nova linha do checklist: o diff desliga ou enfraquece um controle de segurança (testes, logs, checagens de TLS ou de assinatura, autenticação, limites, scanners); testes ou mocks que escondem efeitos colaterais reais; ativação atrasada ou condicional (tempo, host, CI, região).
9. Referência rápida: uma lacuna que outra camada imposta já bloqueia é endurecimento (Low), a menos que a cobertura dessa camada seja ela própria incompleta.
10. Uma nova linha em `"Input and interpretation"`: texto não confiável que chega a um modelo, ou uma ferramenta com autoridade, é uma fronteira de confiança. A política de ferramentas é imposta fora do modelo. Texto recuperado e saída de ferramenta são dados.
11. A linha `"Secrets in the tree"`: um segredo vazado é rotacionado ou revogado fora do git antes de o revisor dispensá-lo ou colocá-lo na linha de base. Uma linha de base guarda só impressões digitais exatas e revisadas. Ela nunca guarda uma exclusão ampla por caminho, regra ou padrão.
12. Passo 10 do `SKILL.md`: um teste de regressão de segurança também tem um caso de controle legítimo. Uma negação total não deve passar por autorização correta.
13. A linha `"Authorization"`: a checagem falha fechada diante de identidade ausente ou parcial. As chaves de cache e de índice têm escopo completo.

`regression-testing`:

6. Passo 6: `"fails before the fix"` significa uma falha de asserção no comportamento sob teste. Um erro de compilação ou de import não conta. Quando o novo teste não compila contra a base, o testador prova o vermelho revertendo só a correção numa cópia isolada. Isso usa as regras de isolamento que o passo 6 já tem.

Cada skill alterada adiciona a origem em `metadata.sources` no seu frontmatter. O valor de `provenance` continua `adapted`.

Release: versão 0.30.0.

## 3. Regras

Os greps abaixo não retornam nenhum resultado na árvore no momento da escrita. O autor do ticket checa de novo na base antes da primeira edição.

- R1 O item 1 está em `.skills/code-review/SKILL.md`, nos passos 4 ou 5. Verificável: `grep -n -i 'sibling call' .skills/code-review/SKILL.md` tem um resultado, e a linha diz que um ponto gêmeo deixado quebrado é Major.
- R2 O item 2 está no item `"Contracts"` de `.skills/code-review/SKILL.md`, com exemplos genéricos. Verificável: `grep -n -i 'hand-maintained' .skills/code-review/SKILL.md` tem um resultado nesse item, e `grep -n -E 'status-cases|fixture_parity' .skills` não tem resultado.
- R3 O item 3 está em `.skills/code-review/SKILL.md`. Verificável: `grep -n -i 'belongs in the tree' .skills/code-review/SKILL.md` tem um resultado, e a linha nomeia os artefatos soltos.
- R4 O item 4 está em `.skills/code-review/SKILL.md`. Verificável: `grep -n -i 'new entry point' .skills/code-review/SKILL.md` tem um resultado, e a linha nomeia a trava de escrita nos arquivos de memória.
- R5 O item 5 está em `.skills/security-review/SKILL.md` e em `references/checklist.md`, nos lugares que a seção 2 nomeia. Verificável: `grep -n -E 'pull_request_target|lifecycle|bidi|symlink' .skills/security-review/references/checklist.md` encontra cada termo; o checklist tem uma nova linha para sinais de arquivo e de CI e nenhuma nova linha para dependências; `grep -n -i 'reporter' .skills/security-review/SKILL.md` e `grep -n -i "author's reasoning" .skills/security-review/SKILL.md` têm cada um um resultado que nomeia a regra da prova e a regra do refutador. Nota: `"author's reasoning"` tem hoje um resultado em outro contexto; o novo resultado deve nomear o refutador.
- R6 O item 6 está em `.skills/regression-testing/SKILL.md`, passo 6. Verificável: `grep -n -i 'assertion failure' .skills/regression-testing/SKILL.md` tem um resultado no passo 6, e o passo diz que um erro de compilação ou de import não conta.
- R7 O item 7 está em `.skills/security-review/SKILL.md` (passo 1) e nas linhas Dependencies e Configuration. Verificável: `grep -n -i 'narrowed' .skills/security-review/SKILL.md .skills/security-review/references/checklist.md` tem um resultado, e `grep -n -i 'baseline' .skills/security-review/SKILL.md` tem um resultado que diz que uma supressão adia.
- R8 O item 8 é uma linha do checklist. Verificável: `grep -n -i 'weaken' .skills/security-review/references/checklist.md` tem um resultado numa linha que nomeia a ativação condicional.
- R9 O item 9 está na Referência rápida de `.skills/security-review/SKILL.md`. Verificável: `grep -n -i 'hardening' .skills/security-review/SKILL.md` tem um resultado na Referência rápida.
- R10 O item 10 é uma linha em `"Input and interpretation"`. Verificável: `grep -n -i 'tool policy' .skills/security-review/references/checklist.md` tem um resultado nessa seção.
- R11 O item 11 está na linha `"Secrets in the tree"`. Verificável: `grep -n -i 'fingerprint' .skills/security-review/references/checklist.md` tem um resultado nessa linha, e a linha diz `"rotated or revoked"` antes de baseline.
- R12 O item 12 está no passo 10 de `.skills/security-review/SKILL.md`. Verificável: `grep -n -i 'control case' .skills/security-review/SKILL.md` tem um resultado no passo 10.
- R13 O item 13 está na linha `"Authorization"`. Verificável: `grep -n -i 'fail closed' .skills/security-review/references/checklist.md` tem um resultado nessa linha.
- R14 Procedência. `metadata.sources` de `code-review` nomeia `akitaonrails/my-skills pr-audit @7b8a5e5`. O mesmo campo de `regression-testing` também o nomeia. O mesmo campo de `security-review` nomeia `akitaonrails/my-skills pr-audit @7b8a5e5` e `akitaonrails/my-skills security-audit @7b8a5e5`. `provenance` continua `adapted` nas três. Verificável: `grep -n -E 'my-skills (pr|security)-audit @7b8a5e5' ` nos frontmatters dos três `SKILL.md`, e `grep -n 'provenance: adapted'` nos mesmos arquivos.
- R15 A mudança não adiciona skill, agente nem comando. As contagens continuam em 13 agentes e 19 + 1 comandos. Verificável: `git diff --name-status main...HEAD` não mostra arquivo adicionado em `.agents/` nem em `.commands/`, nem novo `SKILL.md` em `.skills/`; `dh validate` está verde.
- R16 As edições à mão se limitam a: as três skills e `references/checklist.md` em `.skills/`, este PRD e seu arquivo de revisão, os registros do ticket e os arquivos de release da seção 4. Nenhum arquivo em `dist/` é editado à mão. Verificável: `go run ./cmd/dh build` regenera `dist/`, e `git diff --name-only main...HEAD -- dist` lista só `dist/claude-code/dev-harness/skills/{code-review,security-review,regression-testing}/**`, `plugin.json`, `harness-manifest.json` e `dist/*/bin/**`.
- R17 `go run ./cmd/dh validate .` passa com 0 erros. O build roda duas vezes, e a segunda execução não muda nada em `dist/` além da linha `generated_at` de `harness-manifest.json`. Verificável: rode os dois comandos e compare com `git diff`.
- R18 Prova de clone limpo, porque a PR toca `internal/build/build.go` e `dist/*/bin`: rode `git clone . <tmp>`, rode `go run ./cmd/dh build` lá e compare o `sha256sum` dos binários com a árvore. Verificável: todas as somas são iguais.
- R19 Limites de tamanho, contados em linhas adicionadas (`git diff --numstat main...HEAD -- <file>`, primeira coluna): `code-review/SKILL.md` no máximo 14; `regression-testing/SKILL.md` no máximo 4; `security-review/SKILL.md` no máximo 18; `references/checklist.md` no máximo 24. Verificável: os valores do numstat ficam abaixo dos limites. A linha do frontmatter de `metadata.sources` conta como linha alterada.
- R20 Nenhum título novo. Verificável: `git diff -U0 main...HEAD -- .skills | grep -c '^+#'` imprime 0 para os quatro arquivos tocados.
- R21 Versão 0.30.0. `defaultVersion` em `internal/build/build.go`, `version` em `.claude-plugin/marketplace.json` e o `ref` da origem (`v0.30.0`) mudam juntos. Verificável: `grep` dos três valores mostra `0.30.0` ou `v0.30.0`.
- R22 O escopo da seção 5 continua verdadeiro. Verificável: o diff não tem comando `/dh:pr-audit`, nem script, nem novo subcomando de `dh`, nem `references/ecosystem-checks.md`, nem regra de merge, retarget, atribuição ou lote.

## 4. Docs

Evidência: `grep -rlE 'code-review|security-review|regression-testing'` sobre `docs`, os READMEs, `CHANGELOG.md` e `AGENTS.md`, sem PRDs e ADRs, lista só estes arquivos: `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `docs/plano-produto.html`, `CHANGELOG.md`.

Atualize:

- `CHANGELOG.md` (entrada de release 0.30.0). O roadmap é conferido contra ela por `dh validate`.
- `docs/roadmap.html`, `docs/en/roadmap.html` (card 0.30.0; os cards pt e en devem concordar).
- `.claude-plugin/marketplace.json` (`version` e `ref`).
- `internal/build/build.go` (`defaultVersion`).
- `README.md`, `README.en.md` (linha 20: a linha da versão atual vira 0.30.0 com a tag v0.30.0, e o resumo de uma frase do último release nomeia este PRD). Todo commit de release anterior (43a280b, 00c0fb4, e67dba1) mudou esta linha.

Não afetados:

- `docs/tutorial.html`, `docs/en/tutorial.html` (linhas 859 e 860). A tabela de comandos nomeia `code-review` e `security-review` como a saída de `review` e `secure`. Ela não descreve passos de procedimento.
- `docs/plano-produto.html`. Ele nomeia as skills só no catálogo.
- `docs/inicio-rapido.html`, `docs/en/quick-start.html`, `docs/tutoriais/`, `docs/en/tutorials/`. Nenhum resultado para as três skills.

## 5. Fora do escopo

De `pr-audit`:

- Um comando `/dh:pr-audit` para PRs de contribuidores externos. Ele espera um caso real.
- Regras de conclusão pelo mantenedor, merge, retarget e atribuição.
- Lote de várias PRs.
- Execução em sandbox de PRs não confiáveis.
- O item "Unreleased" do changelog.
- As partes do `pr-audit` que pertencem ao próprio projeto dele.

De `security-audit`:

- Seu processo de auditoria completa em 6 fases.
- Triagem e PATCH de alertas do GitHub por `gh api`.
- Rodar scanners ou o gate de testes a partir do revisor.
- `references/ecosystem-checks.md`.
- `scripts/security-surface.sh`. Não há script nem novo subcomando de `dh`.
- As seções do seu template de saída.
- Sua política de divulgação e de privacidade.

Também fora: qualquer outra mudança nas três skills.

## 6. Decisões do dono

- Dono, 2026-10-16: absorver as ideias em skills existentes e não importar nada inteiro. Decisão: `"follow your recommendations"` sobre a investigação do Coordenador.
- Dono, 2026-10-16: nenhuma skill, agente ou comando novo nesta mudança.
- Dono, 2026-10-16: os arquivos em akitaonrails/my-skills (`pr-audit` e `security-audit`) estão liberados para uso. A cópia literal do texto deles é permitida. Cada skill alterada registra a origem em `metadata.sources` (R14), e `provenance` continua `adapted`.
- Dono, 2026-10-16: os itens do `security-audit` (7 a 13, e a regra do refutador no item 5) entram neste PRD.
- Decisão (regra do `AGENTS.md` sobre tags): a mudança sai como versão 0.30.0 (R21).
- Decisões pendentes: nenhuma.

## Apêndice

### Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| Importar `pr-audit` e `security-audit` inteiros | A maior parte do texto deles pertence ao próprio projeto. O kit deve continuar enxuto. |
| Adicionar agora um comando `/dh:pr-audit` | Ele atende PRs de contribuidores externos. Não existe caso real. Ele aumentaria a contagem de comandos (19 + 1). Adiado até haver um caso real. |
| Escrever uma nova seção "audit" em cada skill | O dono quer um kit enxuto. Linhas curtas dentro dos passos e das linhas existentes carregam as mesmas checagens. |
| Manter as duas fontes em dois PRDs | Os itens do `security-audit` caem nos mesmos arquivos do item 5 e dividem um release. Dois PRDs conflitariam em `checklist.md`. |

### Riscos e decisões pendentes

- Risco: as skills crescem e o texto de gatilho fica ruidoso · mitigação: R19 limita as linhas por arquivo e R20 proíbe títulos novos.
- Risco: o texto adicionado repete uma regra que a skill já tem (por exemplo a regra `"do not execute the payload"` do passo 8) · mitigação: o revisor do ticket checa duplicatas, e a nova linha remete à regra existente.
- Risco: os cards pt e en do roadmap divergem · mitigação: `dh validate` checa a deriva (R17).
- Risco: a prova de clone limpo difere da árvore · mitigação: R18 bloqueia o commit até as somas serem iguais.
- Decisões pendentes: nenhuma.

## Resumo executado

- **Entregue:** `code-review`, `regression-testing`, `security-review` e seu checklist carregam as 13 checagens de `pr-audit` e `security-audit` como linhas curtas dentro de passos e linhas existentes, sem skill, agente, comando nem título novo, lançadas como 0.30.0.
- **Regras:**
  - R1-R13 Cumpridas: cada item tem seu resultado de grep no passo ou na linha nomeada (matriz de QA, `qa-verifier`), e a leitura do texto adicionado confirma o sentido completo dos itens 7, 8, 10 e 13; aderência lida por `product-discovery`: os 13 itens aderem.
  - R14 Cumprida: `metadata.sources` das três skills nomeia `akitaonrails/my-skills pr-audit @7b8a5e5` (e `security-audit @7b8a5e5` em `security-review`); `provenance: adapted` não mudou.
  - R15 Cumprida: `git diff --name-status` mostra só arquivos modificados; 13 agentes, 19 + 1 comandos, 21 skills (`dh validate`).
  - R16 Cumprida: as edições à mão se limitam aos quatro arquivos de skill e aos arquivos de release da seção 4 (linha do README emendada, revalidada); o diff de `dist/` se limita às três skills, `plugin.json`, `harness-manifest.json` e aos cinco binários.
  - R17 Cumprida: `dh validate` 0 erros, `go test ./...` verde; o segundo build muda só `generated_at`.
  - R18 Cumprida: prova de clone limpo, 5 de 5 binários com sha256 igual.
  - R19 Cumprida: linhas adicionadas 5 / 2 / 7 / 7 contra os limites 14 / 4 / 18 / 24.
  - R20 Cumprida: nenhuma nova linha `#` nos quatro arquivos.
  - R21 Cumprida: `defaultVersion` 0.30.0, `version` do marketplace 0.30.0 e `ref` `v0.30.0`.
  - R22 Cumprida: nenhum comando, script, subcomando de `dh`, `ecosystem-checks.md` nem regra de merge/retarget/atribuição/lote.
- **Tickets:** T-2301 (backend, harness-maintainer) `approved`. QA por `qa-verifier`: verificado, uma nota menor (o texto do scanner do item 7 fica na linha Configuration e no passo 1, não na linha Dependencies; aceito). Revisão de código por `cli:codex/gpt-6-luna`: nenhum achado no texto; seu único Major (evidência de R15-R18 ausente do pacote de revisão) foi resolvido com a matriz de QA, sem edição.
- **Docs:** `CHANGELOG.md`, `docs/roadmap.html`, `docs/en/roadmap.html`, `README.md`, `README.en.md`, `.claude-plugin/marketplace.json`, `internal/build/build.go`.
- **Fora do escopo:** um comando `/dh:pr-audit` para PRs externas, regras de merge/retarget/atribuição pelo mantenedor, lote de várias PRs, execução em sandbox de PRs não confiáveis, a auditoria completa em seis fases, triagem de alertas do GitHub, rodar scanners a partir do revisor, `ecosystem-checks.md` e `security-surface.sh`.
