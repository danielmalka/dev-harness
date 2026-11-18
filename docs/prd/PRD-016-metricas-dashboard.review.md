# Revisão do PRD-016 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Veredito
- `changes required`

## Achados
### BLOCKING: P1 nenhuma regra obriga a página a mostrar incidentes, memória, consolidação e datas
- Reportado por: claude
### BLOCKING: P2 contagem de lotes do EPOCHAL e `last_at` ambíguos; falta o caso acima de 1 MiB
- Reportado por: claude
### BLOCKING: P3 o que conta como incidente no RISKS não está definido
- Reportado por: claude
### BLOCKING: P4 os rótulos dos templates en não estão listados
- Reportado por: claude
### BLOCKING: P5 ticket aberto conta ou não como "sem data"
- Reportado por: claude
### BLOCKING: P6 a leitura de MEMORY/RISKS é uma exposição nova, não o risco já aceito; exige `/dh:secure`
- Reportado por: claude
### BLOCKING: P7 o teste de motor não roda a página; tem que ser Playwright
- Reportado por: claude
### MINOR: P8 pasta resolvida; P9 critério do Origin; P10 redação do R9 e do apêndice
- Reportado por: claude

# Revisão do PRD-016 — rodada 2 (final)

## Veredito
- `changes required`

## Achados
P1 a P9 foram aplicados; P10 parcialmente.
### BLOCKING: P11 delimitadores da cópia bruta do EPOCHAL não nomeados
- Reportado por: claude
- Resolução do Coordenador: usar os marcadores reais (`<!-- INICIO-BRUTO` / `<!-- FIM-BRUTO`, mais os equivalentes em inglês) e incluir o EPOCHAL deste repositório no teste.
### BLOCKING: P12 o que é "PRD entregue" em inglês (`delivered on`)
- Reportado por: claude
- Resolução do Coordenador: entregue = Status começa por `entregue em` ou `delivered on`. Os campos existentes `Open`/`Delivered` passam a seguir a mesma regra, e entra um teste com um PRD em inglês.
### MINOR: P10 texto do apêndice; P13 valores quando o arquivo é recusado
- Reportado por: claude

Teto atingido; vai ao dono com as correções aplicadas.
