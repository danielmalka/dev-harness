# Revisão do PRD-017 — rodada 1

- Revisor: document-validator (claude, opus)
- Data: 2026-10-08

## Verdict
- changes required

## Findings
### BLOCKING: P1 o docs-guide não roda o `stamp.sh`; o modelo que ele copia não tem as seções `projeto` nem o Mermaid
- Reported by: claude
### BLOCKING: P2 no modo repo, nenhum mecanismo informa ao agente onde fica `<home>/projects/<nome>/pdocs/`
- Reported by: claude
### BLOCKING: P3 a CSP afrouxada é um risco novo: um script de CDN conseguiria ler `/api/state` e `/api/memory`
- Reported by: claude
### BLOCKING: P4 o grep do R14 procura uma frase em pt-br num arquivo em inglês
- Reported by: claude
### BLOCKING: P5 as duas leituras de "como requisito sempre" sumiram
- Reported by: claude
### MINOR: P6 origem do CDN já fixada só na versão maior; P7 aviso e symlink no R9; P8 pasta que só tem `pdocs/`; P9 AGENTS.md no R13
- Reported by: claude

# Revisão do PRD-017 — rodada 2 (final)

## Verdict
- approved

## Findings
P1 a P9 foram aplicados.
### MINOR: P10 falta prova de que o Mermaid e as fontes carregam sob a CSP com sandbox (Playwright via dashboard)
- Reported by: claude
### MINOR: P11 Limits de document.md e SKILL.md:137 também precisam das exceções
- Reported by: claude
### MINOR: P12 bullets do tutorial duplicados; linha da CSP é `server.go:119`
- Reported by: claude
