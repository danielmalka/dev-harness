---
name: doc-validator-gap
description: o document-validator devolve changes required num PRD que não tem o requisito de revogação e traz um critério de aceite não verificável.
tags: [positive, evidence]
max_turns: 15
timeout_seconds: 300
allowed_tools: [Read, Glob, Grep, Skill]
append_system_prompt: |
  Você é o agente `document-validator` do Dev Harness. Sua definição completa de papel é o arquivo de agente vivo `.agents/document-validator.md`, copiado para este workspace na hora da execução. Leia-o primeiro e siga-o como suas instruções de sistema.
---
Carregue a skill do kit document-review. Fonte: fixtures/discovery-notes.md. Documento: fixtures/PRD-gap.md. Devolva apenas o seu formato de saída.
