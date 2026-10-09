# Revisão do PRD-014 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08
- Fontes: `.harness/projeto-harness-global.md` §8, proposta do dono (conversa de 2026-10-08, registrada em `.harness/projeto-harness-global.md`), decisões do Coordenador, AGENTS.md, código citado.

## Verdict
- changes required

## Findings
### BLOCKING: P1 R4 × R6 conflito ao religar após mover
- Reported by: claude
- Ação: se a entrada com o mesmo nome aponta para caminho que não é mais diretório, `dh link` substitui (religa); se o caminho ainda existe, recusa e exige nome. Teste para cada caso.
### BLOCKING: P2 snapshot-writer.ts fora do R10
- Reported by: claude
- Ação: o mod grava em `<home>/sessions/` (`DH_HOME` ou `~/.harness`) e ignora `DEV_HARNESS_SNAPSHOT_DIR`/`CLAUDE_CONFIG_DIR`; teste do mod com e sem `DH_HOME`.
### BLOCKING: P3 leitores de código sem regra (panel.ts, isHarnessPath, ProjectProgress, status_drift)
- Reported by: claude
- Ação: regra para todos lerem a pasta resolvida; escrita em `<home>/projects/<nome>/` não suja o gate; testes nos dois modos.
### BLOCKING: P4 caminho resolvido na abertura da sessão sem regra
- Reported by: claude
- Ação: regra com o evento `prompt.compose` (provado na sonda H2) e teste de motor nos três modos.
### BLOCKING: P5 projetos repo pré-0.21.0 somem do dashboard
- Reported by: claude
- Ação: quinto item de quebra no R12 (rodar `dh link`/`setup` em cada um; doctor aponta); R16 reescrito.
### BLOCKING: P6 forma de `/api/state.projects[]` × `dh projects --json`
- Reported by: claude
- Ação: `/api/state.projects[]` ganha `mode` e `harness`; CLI imprime os quatro campos, mesma ordem; listar no R13.
### MINOR: P7 YAML com chaves aninhadas e caminho Windows
- Reported by: claude
### MINOR: P8 aviso do doctor com as duas pastas fora da verificação
- Reported by: claude
### MINOR: P9 "pasta global criada vazia" × R1
- Reported by: claude

# Revisão do PRD-014 — rodada 2 (final)

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Verdict
- changes required

## Findings
P1 a P9 foram aplicados sem contradição nova.
### BLOCKING: P10 `status_drift` é do `dh validate` (`internal/kit/status_drift.go`), não do dashboard
- Reported by: claude
### BLOCKING: P11 o aviso do doctor sobre projeto `repo` não registrado está sem critério
- Reported by: claude
### BLOCKING: P12 projeto registrado que resolve `none` não tem destino em `projects[]`
- Reported by: claude
### MINOR: P13 cache do `prompt.compose` indefinido
- Reported by: claude
### MINOR: P14 Status "aprovado" antes do dono ver o PRD; ordem das regras
- Reported by: claude
### MINOR: P15 releitura do config sem fonte citada
- Reported by: claude
- Fonte: contrato acordado com a sessão da extensão VS Code em 2026-10-08 (decisão do Coordenador).

Teto de 2 rodadas atingido. Vai ao dono com P10 a P12 abertos.
