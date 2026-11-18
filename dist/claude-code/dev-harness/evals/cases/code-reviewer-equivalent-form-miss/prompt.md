---
name: code-reviewer-equivalent-form-miss
description: o reviewer (modo code) pede mudanças num diff em que um classificador de forma de comando reconhece `npm test`, mas trata errado em silêncio o equivalente `npm run test --workspace=server`, sem apontar junto o classificador de lint, que trata corretamente as duas formas.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  Você é o agente `reviewer` do Dev Harness, despachado no modo `code`. Sua definição completa de papel é o arquivo de agente vivo `.agents/reviewer.md`, copiado para este workspace na hora da execução. Leia-o primeiro e siga-o como suas instruções de sistema.
---
Carregue a skill do kit code-review. A mudança em revisão é fixtures/diff.md. Devolva apenas o seu formato de saída.
