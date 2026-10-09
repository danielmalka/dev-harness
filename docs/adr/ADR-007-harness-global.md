# ADR-007 · Pasta `.harness` global com resolvedor único em Go

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-08 |
| Decisor | Daniel Lemos (dono) |
| Origem | PRD-014; T-1601; `.harness/projeto-harness-global.md` §8 (regra final aprovada); proposta do dono, conversa de 2026-10-08, registrada em `.harness/projeto-harness-global.md` |
| Reversibilidade | cara: muda onde vivem memória, sessões e token do dashboard, remove `DH_DASHBOARD_ROOTS`/`DH_DASHBOARD_SPRITES` e exige versão nova da extensão VS Code |
| Relação | altera ADR-005 §2 (a trava casa também `<home>/projects/*/`) e ADR-006 D2 (`language` da pasta resolvida) e D5 (escrita em `<home>/projects/*/` não marca o gate), PRD-014 R19 |

## 1. Contexto

- Pergunta: onde fica o `.harness/` de um projeto quando o dono não quer a pasta no repositório, e como o kit, o mod e o dashboard acham o mesmo lugar.
- Fato: hoje cada projeto tem `<repo>/.harness/`, e o `setup` pede para commitar ou ignorar a pasta; o dono quer que a pasta interna seja versionada sem aviso e que repositórios públicos possam não tê-la (proposta do dono, conversa de 2026-10-08, registrada em `.harness/projeto-harness-global.md`).
- Fato: hoje há três lugares fora do repositório sem dono único: snapshots em `~/.claude/dev-harness/sessions`, token em `~/.config/dev-harness/dashboard` e a lista de raízes em `DH_DASHBOARD_ROOTS` (PRD-014 §1).
- Fato: a trava de memória casa pelo sufixo `.harness/<arquivo>` (ADR-005); uma pasta global sem esse sufixo ficaria sem trava (PRD-014, riscos).
- Fato (sondas de 2026-10-08, relatadas pelo Coordenador no despacho da T-1601): H1, `additionalDirectories` aceita `~`, provado só sob `claude -p --permission-mode acceptEdits` (uso interativo na máquina do trabalho a confirmar); H2, o evento de mod `prompt.compose` devolve uma seção com escopo de sessão, enquanto `SessionStartResult` não carrega contexto; H3, o mod lê `$.env.get('DH_HOME')` e `$.env.get('HOME')`.
- Forças: um único registro de projetos (pedido do dono); nenhuma dependência Go nova (ADR-001); rastro zero no repositório em modo global; a trava não pode perder cobertura.

## 2. Opções consideradas

| Dimensão | A · pasta global com mapa no `config.yaml` | B · `.harness/` no repo + excludes globais do git | C · ponteiro no `.git/config` | D · chave pela URL do remoto | E · manter como está |
|---|---|---|---|---|---|
| Fonte única da verdade | sim, `config.yaml` | não, cada repo | não, espalhada nos repos | sim, mas chave instável | não |
| Rastro no repositório | zero em modo global | pasta na árvore, ignorada | nenhum na árvore, mas no `.git` | zero | pasta na árvore |
| Sobrevive a mover/renomear o clone | não, `dh link` religa | sim, a pasta vai junto | sim, `.git/config` vai junto | sim, se o remoto não mudar | sim |
| Sobrevive a `git clean -xdf` | sim | não, a pasta some | sim | sim | depende do ignore |
| Repo sem remoto / clones do mesmo remoto | funciona / nomes distintos | funciona | funciona | quebra / colidem | funciona |
| Trava de memória | precisa casar `<home>/projects/*/` | sem mudança | precisa casar o alvo do ponteiro | precisa casar `<home>` | sem mudança |
| Custo de construir | resolvedor, `dh link`, mapa, dashboard | quase nenhum | resolvedor e escrita em `.git/config` | resolvedor e normalização de URL | nenhum |

## 3. Decisão

- Escolhida: A.
  - `<home>` segue uma ordem só, implementada igual em `harness.Home()` (Go) e `homeFrom` (mod): `DH_HOME`; senão `HOME` + `/.harness`; senão `USERPROFILE` + `\.harness`; senão a pasta do usuário do SO + `.harness`; sem nenhum, o home é desconhecido, os comandos falham com mensagem clara e a trava confere só o sufixo interno (dívida em §4). O `<home>` sempre existe e o `dh` o cria no primeiro uso. Contém `config.yaml` (mapa `projects:` de caminho absoluto para nome, para todos os projetos, mais os padrões), `projects/<nome>/` com os arquivos direto dentro, `sessions/` e `dashboard/`.
  - Um único resolvedor em Go (`internal/harness`, exposto por `dh harness-path`): `<repo>/.harness/`, senão `<home>/projects/<nome>/` pelo mapa, senão nenhum. Com os dois presentes, vence o interno e o `doctor` avisa. Regras e verificações em PRD-014 R1–R7.
  - `config.yaml` só é escrito pelo Go, num subconjunto mínimo de YAML, sem dependência nova; o ADR-001 não muda (PRD-014 R2).
  - Mods: a trava casa os três arquivos em `<repo>/.harness/` e em `<home>/projects/*/`, com `<home>` pela ordem acima, lida por `$.env.get` (sonda H3; PRD-014 R8). A sessão recebe o diretório resolvido por `prompt.compose`, como seção de escopo de sessão (sonda H2), porque `SessionStartResult` não carrega contexto. Permissão de escrita fora do projeto via `additionalDirectories` com `~/.harness` (sonda H1, só em `-p` com `acceptEdits`; PRD-014 R9).
- Motivo decisivo: um registro só, fora do repositório, que vale para todos os projetos e que o resolvedor, o mod e o dashboard leem do mesmo jeito.
- Descartadas:
  - B, excludes globais do git: a pasta continua na árvore, some com `git clean` e não dá fonte única da verdade.
  - C, ponteiro no `.git/config`: rejeitada pelo dono na conversa de 2026-10-08; divergência de nome entre a pasta do repo e a do projeto fica difícil de achar (exemplo dele: pasta local `rabirosca`, `.git/config` apontando para `rabiscui`), e ele quer um registro só no `config.yaml`. Além disso, o mod não lê o `.git/config` e o ponteiro some ao reclonar (PRD-014, alternativas descartadas).
  - D, chave pela URL do remoto: repositório sem remoto fica sem chave, e dois clones do mesmo remoto colidem.
  - Layout `projects/<nome>/.harness/` (variante de A, que constava na §7): aninhamento redundante; o dono prefere uma pasta por projeto com os arquivos direto dentro. O custo para a trava, que deixa de casar só pelo sufixo `.harness/`, foi aceito e é coberto por testes.
  - Manter `DH_DASHBOARD_ROOTS` (e `DH_DASHBOARD_SPRITES`): segundo o dono (conversa de 2026-10-08), só ele usa kit e extensão juntos; o dashboard passa a ler o `config.yaml`, e a extensão muda em release própria (PRD-014 R11–R13).
  - E, manter como está: preserva o aviso de commitar ou ignorar e os três lugares sem dono que motivaram a mudança.
- Preferência declarada: layout plano `projects/<nome>/` e um registro único, ambos escolhas do dono.

## 4. Consequências

- Fica mais fácil: usar o kit num repositório público sem pasta no repositório; achar todos os projetos, sessões e o token num lugar só; o dashboard listar projetos sem variável de ambiente.
- Fica mais difícil: mover ou renomear um repositório desfaz a ligação até o `dh link` corrigir (o `doctor` aponta); a trava passa a ter dois padrões de caminho.
- Dívida assumida: a lacuna de escrita por `Bash` do ADR-005 continua, agora também em `<home>`; sem `DH_HOME`, `HOME` e `USERPROFILE`, a trava não confere o padrão global (o sufixo interno `.harness/` continua conferido), em falha aberta como no ADR-005; cada ambiente (WSL, Windows, macOS) tem o seu `<home>` e o kit não traduz caminhos.
- Contratos afetados: `dh harness-path` é o contrato entre Go e os leitores (mod, comandos, dashboard); `/api/state` e `dh projects --json` mudam para a extensão VS Code (PRD-014 R13, R17, R18); quebras da 0.21.0 no CHANGELOG (PRD-014 R12, R16).

## 5. Validação

- Check que prova que funciona: testes Go do resolvedor nos três modos e no caso "os dois existem"; teste do mod que nega subagente em `<home>/projects/x/MEMORY.md` e em `<repo>/.harness/MEMORY.md`, permite o Coordenador e não trava nome parecido fora dessas pastas; teste de integração de rastro zero (`git status` limpo depois de `dh link`). Verificações exatas em PRD-014 R5, R7 e R8.
- Sinal de que está falhando: `doctor` mostrando `mode: none` num projeto já ligado, sessão sem o caminho resolvido no contexto, ou escrita de subagente em `<home>/projects/*/MEMORY.md` que não foi negada; na máquina do trabalho, sessão interativa pedindo permissão ou recusando escrever em `~/.harness` mesmo com `additionalDirectories` configurado.

## 6. Revisar quando

- Um repositório com mais de uma pessoa usando o `dh` precisar compartilhar memória e tarefas: hoje isso só existe no modo `repo`, e o modo global não cobre esse caso.
- Surgir um segundo adaptador de runtime além do Claude Code.
- O Claude Code oferecer uma pasta durável de dados do plugin exportada para o `Bash`.

## 7. Referências

- `.harness/projeto-harness-global.md` §8 (regra final) e §7 (decisões anteriores, substituídas onde conflitam)
- Proposta do dono, conversa de 2026-10-08, registrada em `.harness/projeto-harness-global.md`
- `docs/prd/PRD-014-harness-global.md`
- `docs/adr/ADR-001-runtime-go.md`, `docs/adr/ADR-005-trava-de-runtime-por-mod.md`, `docs/adr/ADR-006-painel-e-travas-de-commit.md`
