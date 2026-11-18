---
name: secure
description: Faça uma revisão de segurança delimitada de uma mudança, somente leitura
author: malka
argument-hint: "[escopo, intervalo de diff ou lista de arquivos]"
metadata:
  roles: [coordinator, reviewer]
  skills: [security-review, external-clis]
  writes: nenhum (somente leitura)
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Escopo: $ARGUMENTS. Despache `reviewer` no modo `security` com a ferramenta Agent, modelo `opus`, e exija que ele carregue a skill do kit `security-review`, só quando `reviewers.security` estiver ausente de `.harness/project.yaml` ou contiver `claude`. Entradas da forma `cli:...` em `reviewers.security` seguem o procedimento de `external-clis` descrito em `.agents/coordinator.md` ("Revisores de CLI externos"); tanto o despacho do Claude quanto qualquer entrada `cli:...` contam contra o limite de dois especialistas concorrentes. Passe o diff estável ou a lista explícita de arquivos, as fronteiras que o projeto considera sensíveis, os contratos que o código expõe, qualquer saída de scanner existente e a declaração de que nenhuma execução está autorizada, a menos que o dono a tenha concedido nesta sessão.

## Pré-requisitos
A mudança em revisão como um diff estável ou uma lista explícita de arquivos; sem ela, pergunte ao dono e pare em vez de varrer o repositório. Fronteiras sensíveis e classes de dados vêm do projeto; quando não declaradas, são inferidas, e toda inferência é listada como hipótese, não como fato. Consulte `.harness/RISKS.md` para incidentes anteriores na área tocada. Quando a mudança não tem superfície de autenticação, dados, entrada nem execução, diga isso e pare.

## Saída
O relatório de `security-review`: escopo e fronteiras tocadas, achados cada um classificado como vulnerabilidade demonstrada, exposição confirmada ou hipótese que exige investigação, com origem, caminho, destino, mitigação presente e por que é insuficiente, atacante e ganho, cenário, menor mitigação e a verificação dela; além de hipóteses em aberto, não problemas com a evidência que os descartou e checagens. Registre as exposições confirmadas como incidentes em `.harness/RISKS.md` e o veredito em `.harness/MEMORY.md`.

## Limites
- Somente leitura: nenhuma correção é aplicada e nenhum arquivo é editado por este comando.
- Nunca execute um payload, exploit ou scanner sem autorização registrada nesta sessão, e nunca toque um sistema vivo, compartilhado ou de terceiros.
- Nunca imprima um segredo, token ou dado pessoal no relatório; cite o local e o tipo dele.
- A falta de evidência sozinha não estabelece um não problema; preocupações não comprovadas continuam qualificadas como hipóteses, e a falta de uma cadeia de exploração nunca apaga uma exposição confirmada.
- Não rebaixe esta revisão do padrão Opus para economizar custo.
- Uma lista em `reviewers.security` do `.harness/project.yaml` é a autorização explícita do dono para o provedor e o modelo dessa etapa, para as entradas que ela nomeia e nada mais; os limites de gasto e o orçamento combinado continuam valendo.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`/dh:fix` para uma vulnerabilidade demonstrada ou exposição confirmada, depois rode este comando de novo sobre o novo estado.
