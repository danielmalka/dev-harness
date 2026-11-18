---
name: doctor
description: Diagnostique a prontidão do kit, do runtime e do projeto
author: malka
argument-hint: ""
metadata:
  roles: [coordinator]
  skills: [project-onboarding]
  writes: nenhum (somente leitura)
---
## Papel
Atue como o `coordinator` do kit nesta sessão (leia a definição do agente `coordinator` que acompanha o kit se esta sessão não foi iniciada com ele). Leia `.harness/MEMORY.md`, se existir, antes de qualquer outra coisa.

## Roteamento
Rode o diagnóstico você mesmo; não despache um especialista. Carregue a skill do kit `project-onboarding` e siga só o modo de diagnóstico dela, nunca o modo de configuração. Quando houver um shell disponível, rode o `bin/dh doctor` que acompanha o kit, a partir do diretório do plugin carregado (o caminho de `--plugin-dir`) ou do checkout do kit `bin/dh doctor`, e incorpore a saída dele como inventário mecânico e sonda do ambiente. Se o script estiver ausente ou não puder rodar, diga isso e continue com o diagnóstico observacional.

## Pré-requisitos
Nada além de um pacote do kit resolvível. Liste os agentes, comandos e skills realmente descobertos, por nome, a partir dos agentes, skills e templates que acompanham este kit; se a localização do kit não puder ser resolvida, reporte o kit como não verificado e pare. Não pergunte nada ao dono: este comando só observa.

## Saída
O relatório de diagnóstico de `project-onboarding`: inventário do kit (comandos, papéis, skills encontrados por nome, referências não resolvidas), a tabela do ambiente por camada com comando ou observação, status e evidência, impacto na capacidade, o `mode` e o `dir` que `dh doctor` imprime (e os avisos dele: as duas pastas presentes, projeto não registrado, `local.yaml` não ignorado, dica de `additionalDirectories`), a linha `tamanho do MEMORY.md: <n> bytes` e, acima de 40000 bytes, o aviso `aviso: o MEMORY.md tem <n> bytes, acima do orçamento de 40000: rode /dh:consolidate-memory` (ofereça `/dh:consolidate-memory` uma vez; nunca o rode), e, só quando o `FACTS.md` existe, a linha `tamanho do FACTS.md: <n> linhas, <b> bytes` com o aviso `aviso: o FACTS.md tem …, acima do teto de 150` (linhas) ou `… acima do teto de 12288 (12 KB)` (bytes), que pede ao Coordenador a curadoria do arquivo (repasse o aviso; nunca edite o `FACTS.md`), se o harness dir e o `project.yaml` dele existem, quais ferramentas estão disponíveis, limitações e próximo passo. Inclua o status de `bin/dh doctor` quando ele rodou. Não escreva arquivos. Quando `.harness/MEMORY.md` já existe, registre que o diagnóstico rodou e o que foi achado ausente.

## Limites
- Nunca instale, atualize nem configure nada; reporte só o que foi observado.
- Nunca reporte uma camada como presente sem a saída observada que a mostra; o desconhecido continua desconhecido.
- Não crie o harness dir nem qualquer registro aqui; `/dh:setup` cuida disso.
- Só o Coordenador escreve `.harness/MEMORY.md`, `EPOCHAL.md` e `RISKS.md`.

## Próximo
`dh link` quando `mode: none` mas já existe uma pasta global para este projeto; `/dh:setup` quando nada existe ou o perfil do projeto está ausente; caso contrário, o comando de trabalho da tarefa em curso.
