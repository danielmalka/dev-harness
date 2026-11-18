# Caderno de notas

## Regras de curadoria
- Reordene por prioridade a cada leitura.
- Mantenha só notas recorrentes e de alto valor.
- No máximo 10 itens por categoria.
- Cada item traz a data e um "Faça em vez disso".

## Guardas de comportamento do domínio
1. **[2026-09-20] Os docs HTML internos usam papel de tinta e sinal neon, nunca um verde lavado**
   Faça em vez disso: gere o HTML interno por `.skills/doc-template-html`. Carvão `#0C0D10`, papel `#F1EFE8`, neon `#3DFF9A` só nos destaques, areia `#E2D2A8` em links e faixas. A marca é `Harness.`. Não reintroduza roxo, ciano nem um segundo verde lavado.

## Diretivas do dono
1. **[2026-09-20] As fontes canônicas do kit ficam em pastas com ponto**
   Faça em vez disso: escreva agentes em `.agents/`, skills em `.skills/`, comandos em `.commands/`, mais `templates/` e `profiles/`. Regenere `dist/` com `scripts/build-claude-code.sh`. Nunca edite à mão o pacote gerado.

2. **[2026-09-20] A licença é MIT**
   Faça em vez disso: mantenha `LICENSE` e o `plugin.json` gerado em MIT. Não introduza uma segunda licença sem perguntar.

3. **[2026-09-20] Os scripts de manutenção são bash + Python 3**
   Faça em vez disso: use `scripts/validate.py`, `scripts/doctor.sh` e `scripts/build-claude-code.sh`. Não acrescente um empacotador em Node, a menos que um consumidor precise.
