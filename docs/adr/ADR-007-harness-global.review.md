# Revisão do ADR-007 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Verdict
- changes required

## Findings
### BLOCKING: P1 motivo de descarte da opção C sem fonte; falta a linha mover/renomear
- Reported by: claude
- Resolução do Coordenador: o motivo "divergência de nome" vem do dono na conversa de 2026-10-08 (exemplo `rabirosca` × `rabiscui`). Atribuir a ele, somar os motivos do PRD e incluir a linha mover/renomear na tabela.
### BLOCKING: P2 o ADR-007 altera o ADR-005 §2 (padrão da trava) e o ADR-006 D2 (leitura do `language`)
- Reported by: claude
- Resolução do Coordenador: o campo Relação passa a dizer "altera". O ADR-005 e o ADR-006 ganham "Alterado por: ADR-007" no cabeçalho.
### MINOR: P3 "só o dono usa kit e extensão" — atribuir ao dono (conversa de 2026-10-08)
- Reported by: claude
### MINOR: P4 `<home>` = `DH_HOME`, senão `HOME` + `/.harness`
- Reported by: claude
### MINOR: P5 H1 provado só em `-p` com acceptEdits; sinal de falha na máquina do trabalho
- Reported by: claude
### MINOR: P6 gatilho de revisão não observável; citação errada da §7
- Reported by: claude

# Revisão do ADR-007 — rodada 2 (final)

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Verdict
- approved

## Findings
P1 a P6 foram aplicados.
### MINOR: P7 a Relação cita a D2 do ADR-006, mas a regra de "não sujar o gate" é a D5
- Reported by: claude
### MINOR: P8 sem `HOME` no Windows, o mod deve recorrer a `USERPROFILE`
- Reported by: claude
