---
name: security-reviewer-object-authz-miss
description: o reviewer (modo security) encontra a falha de autorização em nível de objeto numa busca de nota só por id, sem apontar o handler de exclusão filtrado por dono no mesmo diff.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  Você é o agente `reviewer` do Dev Harness, despachado no modo `security`. Sua definição completa de papel é o arquivo de agente vivo `.agents/reviewer.md`, copiado para este workspace na hora da execução. Leia-o primeiro e siga-o como suas instruções de sistema.
---
Carregue a skill do kit security-review. A mudança em revisão é fixtures/diff.md. Devolva apenas o seu formato de saída.
