# Revisão · PRD-018 e T-1801

- Status da revisão: `approve`
- Rodada: 1 de 2
- Reportado por: claude (document-validator, opus)

## Achados

### MINOR: formato exato do `--json` (R1)
- Reportado por: claude
- `writeJSON` gera JSON compacto; trocar "exatamente" por "objeto com exatamente estas três chaves, nesta ordem, via `writeJSON`".

### MINOR: "gravado como digitado" (R2)
- Reportado por: claude
- `quoteVal` põe aspas em `~/sprites`; o critério é o valor lido de volta por `LoadConfig` ser o argumento literal, sem expandir `~`.

### MINOR: `set` sem `config.yaml` existente
- Reportado por: claude
- Cria o arquivo; declarar no R2 ou no teste.

### MINOR: origem da decisão
- Reportado por: claude
- A forma vem do §7 do renova.md; a decisão 10 (A, adiar o editor) foi superada pelo pedido do dono de 2026-10-12.

### NIT: "removeu" → "deixou sem efeito" (§1)
- Reportado por: claude

### MINOR: trava não exportada (R5)
- Reportado por: claude
- `lock` não é exportada; o builder precisa de um wrapper exportado em `internal/harness`, como o `Link`.
