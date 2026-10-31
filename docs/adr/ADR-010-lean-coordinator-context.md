# ADR-010 · Contexto enxuto do Coordenador: respostas curtas, relatório em arquivo, mapa do repositório

| Campo | Valor |
|---|---|
| Status | aceito |
| Data | 2026-10-31 |
| Decisor | Daniel Malka (dono) |
| Origem | PRD-032 (R1 a R12, D1 a D8, emenda de 2026-10-31); `docs/prd/PRD-032-lean-coordinator-context.md` |
| Reversibilidade | barata: edição de fonte do kit e uma checagem do `dh doctor` |
| Relação | complementa o ADR-009; não muda regras de revisão, estados de achado nem tetos de rodada (PRD-031) |

## 1. Contexto

- Pergunta: o que enche o contexto do Coordenador sem ajudar a decidir.
- Fato: respostas de especialistas voltam longas, `coordinator.md` carrega seções de um só comando, `/dh:auto` lê os seis comandos encadeados por inteiro e o `MEMORY.md` deste projeto passa de 185 KB.
- Fato: sem mapa persistente, cada tarefa remapeia o repositório.

## 2. Decisão

- R1 Resposta de especialista até 200 palavras (padrão; o Coordenador pode trocar e diz o motivo). Relatório com mais de 40 linhas vai a `<pasta da tarefa>/reports/<papel>-<n>.md`; a resposta traz veredito, contagens e no máximo 5 achados bloqueantes.
- R2 Escrevem o próprio relatório `builder`, `qa-verifier`, `docs-guide` e `debugger`. `reviewer` e `repo-scout` não têm Write e respondem inline. CLI revisora nunca escreve: quem roda a chamada grava o arquivo.
- R3 O Coordenador lê o arquivo de relatório só para julgar achado bloqueante, disputa ou decisão que a resposta não sustenta. O `MEMORY.md` guarda caminho e veredito.
- R4 Com mais de um revisor na etapa, o `docs-guide` (sonnet) funde as saídas em `reports/<etapa>-merged-<rodada>.md`. Nunca remove achado bloqueante, nunca baixa severidade, mantém cada veredito literal e não julga.
- R5 e R6 `/dh:auto` lê só `Routing`, `Prerequisites` e `Limits` de cada comando encadeado, por uma tabela; seções de um só comando saem do `coordinator.md` para o comando que as usa.
- R7 e R8 Entrada do `MEMORY.md` até 1200 caracteres e arquivo até 40000 bytes. Acima, o Coordenador oferece `/dh:consolidate-memory` uma vez; o `dh doctor` imprime o tamanho e avisa. Registros antigos não são reescritos.
- R9 a R12 O `/dh:setup` ganha uma segunda fase opcional (`bootstrap` atualiza só o mapa). O `docs-guide` escreve `<harness dir>/REPO-MAP.md` (até 200 linhas, sem template nem chave no `project.yaml`). `repo-scout`, `builder` e `implementation-planner` o leem como dica e verificam contra a árvore; o sha de `HEAD` no cabeçalho marca mapa velho. É o único registro que o setup pode sobrescrever.
- Contagens intactas: 13 agentes, 22 skills, 19 + 1 comandos.

### Reversão de D1 (emenda de 2026-10-31)

- Proposta inicial (D1): o `reviewer` ganharia Write limitado, em prosa, ao arquivo de relatório.
- Revertida após revisão de segurança: achado médio confirmado. Texto não confiável de PR podia chegar ao `project.yaml` ou ao `REVIEW.md` por uma escrita do `reviewer`, porque um limite só em prosa não impede a escrita.
- Efeito: o `reviewer` segue somente-leitura e devolve tudo inline (R1, R2).
- Candidato futuro, fora do escopo: regra do `runtime-guard` que permita uma escrita estreita do `reviewer`, imposta pelo mod (ADR-005).

## 3. Consequências

- Fica mais fácil: menos bytes por turno do Coordenador; mapa reaproveitado entre tarefas.
- Fica mais difícil: o `reviewer` perde detalhe além de 5 achados bloqueantes na resposta; o Coordenador pode pedir revisão mais estreita. O resumo pode esconder achado se o Coordenador pular o arquivo, mitigado pelo limite de 5 e pela proibição de baixar severidade na fusão.
- Mapa velho engana o plano; mitigado por ser dica, pelo sha e pela verificação.
- Contratos afetados: `.agents/coordinator.md`, `.agents/reviewer.md`, `.commands/setup.md`, `.commands/auto.md` e a skill `project-onboarding` são o procedimento; `AGENTS.md` só enuncia a regra vigente. Sem espelho em `docs/en/` (não há `docs/en/adr/`).

## 4. Alternativas

| Alternativa | Por que não |
|---|---|
| `reviewer` com Write em prosa (D1) | Exposição confirmada; ver a reversão acima |
| Novo comando `/dh:bootstrap` | Muda a contagem 19 + 1 em vários documentos; o `/dh:setup` já cuida do conhecimento inicial do projeto (pode voltar no futuro) |
| Reescrever o `MEMORY.md` existente | O disco é a verdade; só oferta de `/dh:consolidate-memory` |
| Medir ganho com eval paga | Só a pedido do dono; a medida é em bytes, não em tokens |

## 5. Revisar quando

- Respostas inline do `reviewer` esconderem achado que o arquivo traria.
- O `runtime-guard` ganhar a regra de escrita estreita.
- Mapa velho causar plano errado mais de uma vez.
