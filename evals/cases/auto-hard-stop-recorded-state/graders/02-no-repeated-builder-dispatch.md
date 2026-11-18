---
type: tool_used
tool: Agent
input_match: '"subagent_type":[ \t]*"(dh:)?builder"'
min: 0
max: 1
weight: 1
---

Um papel que poderia tentar a implementação da Fatia 1 (`builder`) é
despachado no máximo uma vez. O `.commands/auto.md` exige que, numa parada
obrigatória, a execução termine "sem exceção ... sem laço de espera, sem
nova tentativa automática" — uma execução que despacha pela segunda vez o
mesmo trabalho de implementação bloqueado, esperando outro resultado
depois de já saber que a dependência falta, reprova neste grader. Zero
despachos também passa: reconhecer a dependência `qrcode` ausente
diretamente no texto da Fatia 1 do `fixtures/PLAN.md`, antes de despachar
qualquer builder, é o caminho mais barato e igualmente conforme.
O `input_match` está ancorado no campo serializado `subagent_type`, então
um prompt que apenas menciona um builder não conta como despacho.
