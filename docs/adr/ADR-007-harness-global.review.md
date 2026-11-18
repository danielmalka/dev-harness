# Revisão do ADR-007 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Veredito
- `changes required`

## Achados
### BLOCKING: P1 motivo de descarte da opção C sem fonte; falta a linha mover/renomear
- Reportado por: claude
- Resolução do Coordenador: o motivo "divergência de nome" vem do dono na conversa de 2026-10-08 (exemplo `rabirosca` × `rabiscui`). Atribuir a ele, somar os motivos do PRD e incluir a linha mover/renomear na tabela.
### BLOCKING: P2 o ADR-007 altera o ADR-005 §2 (padrão da trava) e o ADR-006 D2 (leitura do `language`)
- Reportado por: claude
- Resolução do Coordenador: o campo Relação passa a dizer "altera". O ADR-005 e o ADR-006 ganham "Alterado por: ADR-007" no cabeçalho.
### MINOR: P3 "só o dono usa kit e extensão" — atribuir ao dono (conversa de 2026-10-08)
- Reportado por: claude
### MINOR: P4 `<home>` = `DH_HOME`, senão `HOME` + `/.harness`
- Reportado por: claude
### MINOR: P5 H1 provado só em `-p` com acceptEdits; sinal de falha na máquina do trabalho
- Reportado por: claude
### MINOR: P6 gatilho de revisão não observável; citação errada da §7
- Reportado por: claude

# Revisão do ADR-007 — rodada 2 (final)

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Veredito
- `approved`

## Achados
P1 a P6 foram aplicados.
### MINOR: P7 a Relação cita a D2 do ADR-006, mas a regra de "não sujar o gate" é a D5
- Reportado por: claude
### MINOR: P8 sem `HOME` no Windows, o mod deve recorrer a `USERPROFILE`
- Reportado por: claude
