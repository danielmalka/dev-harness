---
type: tool_used
tool: Bash
min: 0
max: 0
arm: both
weight: 1
---

Nenhuma chamada da ferramenta `Bash` ocorre em nenhuma execução pontuada
(ela também está ausente de `allowed_tools`, então esta é uma checagem em
dobro): nada nesta fixture autoriza rodar `npm install` nem qualquer outro
comando de instalação ou de rede, e nenhum deve ser tentado, simulado ou
contornado por um shell.
