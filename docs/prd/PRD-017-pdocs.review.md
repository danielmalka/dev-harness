# Revisão do PRD-017 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Veredito
- `changes required`

## Achados
### BLOCKING: P1 o docs-guide não roda o `stamp.sh`; o modelo que ele copia não tem as seções `projeto` nem o Mermaid
- Reportado por: claude
### BLOCKING: P2 no modo repo, nenhum mecanismo informa ao agente onde fica `<home>/projects/<nome>/pdocs/`
- Reportado por: claude
### BLOCKING: P3 a CSP afrouxada é um risco novo: um script de CDN conseguiria ler `/api/state` e `/api/memory`
- Reportado por: claude
### BLOCKING: P4 o grep do R14 procura uma frase em pt-br num arquivo em inglês
- Reportado por: claude
### BLOCKING: P5 as duas leituras de "como requisito sempre" sumiram
- Reportado por: claude
### MINOR: P6 origem do CDN já fixada só na versão maior; P7 aviso e symlink no R9; P8 pasta que só tem `pdocs/`; P9 AGENTS.md no R13
- Reportado por: claude

# Revisão do PRD-017 — rodada 2 (final)

## Veredito
- `approved`

## Achados
P1 a P9 foram aplicados.
### MINOR: P10 falta prova de que o Mermaid e as fontes carregam sob a CSP com sandbox (Playwright via dashboard)
- Reportado por: claude
### MINOR: P11 Limits de document.md e SKILL.md:137 também precisam das exceções
- Reportado por: claude
### MINOR: P12 bullets do tutorial duplicados; linha da CSP é `server.go:119`
- Reportado por: claude
