# Revisão · PRD-018 e T-1801

- Review status: approve
- Rodada: 1 de 2
- Reported by: claude (document-validator, opus)

## Findings

### MINOR: formato exato do `--json` (R1)
- Reported by: claude
- `writeJSON` gera JSON compacto; trocar "exatamente" por "objeto com exatamente estas três chaves, nesta ordem, via `writeJSON`".

### MINOR: "gravado como digitado" (R2)
- Reported by: claude
- `quoteVal` põe aspas em `~/sprites`; o critério é o valor lido de volta por `LoadConfig` ser o argumento literal, sem expandir `~`.

### MINOR: `set` sem `config.yaml` existente
- Reported by: claude
- Cria o arquivo; declarar no R2 ou no teste.

### MINOR: origem da decisão
- Reported by: claude
- A forma vem do §7 do renova.md; a decisão 10 (A, adiar o editor) foi superada pelo pedido do dono de 2026-10-12.

### NIT: "removeu" → "deixou sem efeito" (§1)
- Reported by: claude

### MINOR: trava não exportada (R5)
- Reported by: claude
- `lock` não é exportada; o builder precisa de um wrapper exportado em `internal/harness`, como o `Link`.
