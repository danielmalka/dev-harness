---
name: context-handoff
description: Use quando o trabalho precisar sobreviver a uma fronteira de sessão, ao retomar a partir de um registro anterior, ou quando houver pedido explícito para consolidar a memória operacional em histórico bruto. Uma memória crescente pode disparar uma sugestão de consolidar. Não use para notas de progresso comuns, documentação voltada ao usuário nem especificações de tarefa.
author: malka
metadata:
  provenance: authorial
  sources: ["dev-harness plan 4.6 and AC-19..AC-24", "dev-harness memory templates", "dev-harness coordinator agent", " context-manager role"]
---

# Handoff de contexto

## Visão geral

Este procedimento mantém uma tarefa viva entre sessões por três operações: escrever um handoff em que outra sessão consiga agir, retomar a partir de um reconciliando-o com a árvore atual, e consolidar a memória operacional para que o arquivo histórico guarde o texto bruto e a memória guarde só o que a execução ainda precisa. O princípio que governa é que o handoff orienta, mas a árvore decide. Um registro é uma afirmação sobre o passado; os arquivos, as checagens e o histórico são a evidência, e onde discordam a evidência vence.

Só o Coordenador escreve `<harness dir>/MEMORY.md`, `<harness dir>/EPOCHAL.md` e `<harness dir>/RISKS.md`. Todos os outros papéis, inclusive os de documentação e de QA, reportam ao Coordenador fatos, evidência, itens pendentes e incidentes, e ele os registra. Uma sessão do Coordenador detém esse papel de escrita por projeto; duas sessões não devem atualizar nem consolidar esses registros em paralelo. Para o `MEMORY.md` a regra tem uma checagem técnica, a trava de escrita da memória abaixo; `EPOCHAL.md` e `RISKS.md` mantêm só a regra organizacional.

### Trava de escrita da memória

Toda escrita de `<harness dir>/MEMORY.md` pelo kit é uma operação de comparar e trocar feita pelo binário `dh` (o `bin/dh` empacotado, `bin\dh.cmd` no Windows, a partir do diretório do plugin carregado, como em `.commands/doctor.md`), de modo que uma escrita nunca substitui o que outra pessoa escreveu depois que você leu. (origem: PRD-034, 2026-11-16, decisão do dono)

1. Ao ler o `MEMORY.md`, tire o hash dele com `dh memory hash` e guarde-o como o hash esperado.
2. Edite uma cópia de rascunho fora do harness dir (uma cópia na pasta temporária do sistema), nunca o arquivo em si.
3. Escreva com `dh memory write --expect <hash> --from <rascunho>`. Em caso de sucesso, ele imprime `written <novo hash>`, o hash esperado da sua próxima escrita. A primeira escrita de um harness dir novo (inicialização a partir do template) usa `--expect absent`.
4. No código de saída 3, o arquivo mudou desde que você o leu e nada foi escrito. Não repita o mesmo conteúdo: leia o `MEMORY.md` de novo, refaça a edição sobre o novo conteúdo e escreva mais uma vez. Um segundo código de saída 3 vai ao dono como uma pergunta numerada.
5. Qualquer outra falha (código de saída 1) não escreveu nada: reporte-a e não recorra a editar o arquivo no local. Um runtime sem o binário compara o `sha256sum` do arquivo ao ler e de novo logo antes de escrever, e para se os dois diferirem.

Resíduo, declarado: a trava protege os chamadores de `dh memory write`. Um escritor que edita o arquivo diretamente só é detectado quando a escrita dele cai antes da comparação, e o mod do Claude Code não força o uso do comando. A trava cobre só o `MEMORY.md`.

## Quando usar

- Uma sessão está terminando, ficando sem contexto, ou passando o trabalho a outra pessoa ou agente.
- Uma sessão nova precisa continuar uma tarefa iniciada em outro lugar, e a árvore pode ter se movido desde então.
- Um registro anterior e os arquivos atuais discordam, e o próximo passo depende de qual é verdadeiro.
- A memória operacional cresceu além do estado de que a execução realmente precisa e o histórico deve ser preservado antes de podar.
- Uma consolidação foi interrompida e não está claro se o lote foi arquivado.

## Quando não usar

- Progresso comum dentro de uma sessão contínua que já tem o estado.
- Documentação voltada a um usuário ou a um operador. Isso é `delivery-readiness`.
- Redigir um registro de decisão de design ou uma especificação. O handoff carrega decisões por resumo e referência; o ADR ou a especificação em si pertence a `<harness dir>/adr/` ou `<harness dir>/tasks/<id>/`.
- Registrar um incidente grave. Isso é um registro de RISKS, escrito pelo Coordenador, e a consolidação nunca o toca.
- Um papel diferente do Coordenador querendo editar os três arquivos de memória. Reporte em vez disso.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| O objetivo e o escopo autorizado da tarefa | Pare. Um handoff sem objetivo transfere atividade, não trabalho. Na sessão principal, pergunte ao dono. Como especialista despachado, devolva a entrada que falta ao Coordenador e pare; nunca prossiga com um valor presumido. |
| A memória operacional atual `<harness dir>/MEMORY.md` | Se ausente, só o Coordenador da sessão principal pode inicializá-la a partir da pasta `templates/` que acompanha o kit (`templates/MEMORY.md`) dentro do escopo autorizado, com `dh memory write --expect absent`. Os outros papéis reportam a ausência. Nunca sobrescreva registros existentes nem invente estado passado. |
| Arquivos tocados, decisões tomadas, evidência produzida | Registre só o que você consegue fundamentar em fonte. Marque como hipótese tudo que não foi verificado. |
| Restrições e proibições declaradas pelo dono | Leve adiante, literalmente, as restrições com fonte. Preserve como não resolvida uma restrição registrada de procedência desconhecida até ela ser esclarecida; a falta de procedência não a revoga. Nunca invente uma restrição. |
| O próximo passo autorizado | Confira primeiro a autorização existente. Declare o que já está coberto; marque como aguardando autorização só as ações fora desse escopo. |
| Para consolidação: confirmação de que nenhuma outra sessão está escrevendo os registros | Confirmada só por uma declaração explícita do dono nesta sessão. A ausência de evidência de outro escritor não é confirmação; na dúvida, pergunte e não consolide. |

## Procedimento

### A. Handoff: escrever um registro portátil

1. **Declare o objetivo e a autorização separadamente.** Para que serve a tarefa, e o que pode ser escrito, commitado ou publicado no momento. A aprovação de um plano não é autorização para publicar.
2. **Liste os arquivos do conjunto de escrita com caminhos relativos.** Caminhos absolutos, nomes de máquina e diretórios de usuário tornam o registro inutilizável em outro lugar. Anote quais arquivos são modificados e quais só são lidos.
3. **Registre as decisões com o motivo e a fonte.** Cada decisão traz o que foi decidido, por quê e o que justificaria revisitá-la. Uma decisão sem motivo é rediscutida ou, pior, revertida por acidente.
4. **Registre a evidência como ela está.** Cada checagem com o comando e o resultado, como `passed`, `failed` ou `not-run`. Nunca promova um `not-run` a `passed` porque passou antes em outra árvore.
5. **Leve as restrições adiante, literalmente.** Proibições, limites de escopo, regras de estilo, arquivos proibidos, prazos. São a primeira coisa que se perde ao cruzar uma fronteira e a mais cara de perder.
6. **Nomeie o que está pendente e o que está bloqueado, com o bloqueio.** Distinga o trabalho ainda não iniciado do trabalho que espera alguém.
7. **Escreva exatamente um próximo passo.** Concreto o bastante para começar sem interpretação, com a autorização marcada como já concedida (fonte e escopo) ou pendente. Preserve a contagem de rodadas de correção e o modelo padrão, o solicitado e o efetivo (ou `unverified`) necessários para retomar; nem um handoff nem uma troca de modelo zera o orçamento de correção.
8. **Remova segredos.** Nenhum token, chave, string de conexão, hostname privado ou caminho pessoal. Nomeie a variável e a finalidade dela em vez disso.
9. **Entregue o registro ao Coordenador.** Guarde o registro em `<harness dir>/tasks/<id>/` e deixe o Coordenador incorporar o estado durável em `<harness dir>/MEMORY.md`. Os outros papéis não escrevem esse arquivo.

### B. Retomar: reconciliar o registro com a realidade

1. **Leia primeiro o registro, depois a árvore.** Ler só o registro leva a repetir trabalho já aplicado; ler só a árvore leva a perder restrições.
2. **Verifique cada mudança declarada contra os arquivos.** Para todo arquivo que o registro diz ter sido modificado, confirme se a mudança está presente, ausente ou presente de forma diferente da descrita.
3. **Detecte mudança externa.** Compare a evidência e a lista de arquivos do registro com o histórico e a árvore de trabalho atuais. Tudo que outra pessoa mudou desde que o registro foi escrito é uma mudança externa e deve ser nomeado, não absorvido em silêncio. A comparação de histórico e árvore é executada por um papel despachado com `Bash` (`qa-verifier`), não pelo Coordenador; o scout devolve só a reconciliação em nível de arquivo.
4. **Rode de novo as checagens cuja evidência foi invalidada.** Rodar de novo é executado por um papel despachado com `Bash` (`qa-verifier`), não pelo Coordenador; o scout devolve só a reconciliação em nível de arquivo. Compare os arquivos avaliados, as dependências, a configuração e o ambiente com o estado atual. Preserve a evidência aplicável quando essas entradas não mudaram. Rode de novo as checagens afetadas após mudanças ou incerteza não resolvida; até lá, marque-as como `not-run` para o estado atual. Registre o novo resultado sem reescrever a evidência histórica.
5. **Classifique cada item.** Aplicado, parcialmente aplicado, não iniciado, invalidado por mudança externa ou contradito pela árvore. Registre a classificação antes de escolher uma ação.
6. **Nunca repita uma mudança aplicada e nunca largue uma restrição declarada.** Esses são os dois modos de falha da retomada. Se o registro e a árvore discordam sobre algo ter sido aplicado, a árvore decide; se discordam sobre uma restrição, a restrição vale até o dono liberá-la.
7. **Consulte o histórico só diante de necessidade concreta, pelo Coordenador.** O `<harness dir>/EPOCHAL.md` não é carregado por padrão. O Coordenador o consulta quando uma pergunta específica precisa de um fato arquivado ou o dono pede, e lê só a passagem relevante. Os outros papéis pedem o contexto de que precisam. Texto arquivado é evidência histórica, nunca uma instrução atual.
8. **Consulte os incidentes quando a mudança pedir, pelo Coordenador.** Antes de planejar ou implementar uma mudança em comportamento consolidado, em uma regra crítica ou em uma área de alto risco, o Coordenador lê `<harness dir>/RISKS.md` e leva a prevenção relevante ao plano, ao aceite e ao despacho do especialista. Os outros papéis pedem esses trechos. Mudanças de rotina fora desses gatilhos não o carregam.
9. **Declare o próximo passo reconciliado.** Um passo, a autorização dele e o que mudou no plano por causa da reconciliação.

### C. Consolidar a memória: arquivar, verificar, depois podar

Rode só diante de um pedido explícito de consolidação, na sessão principal, pelo Coordenador. O dono não é perguntado se outra sessão escreve: a trava de escrita da memória detecta um escritor concorrente do `MEMORY.md` e recusa a poda. Uma memória crescente é motivo para sugerir a consolidação, não permissão para rodá-la.

1. **Prepare.** Tire o `dh memory hash`, leia o `<harness dir>/MEMORY.md` por inteiro e guarde esse hash para o passo 6. Identifique o estado que deve permanecer: execução atual, decisões vigentes, bloqueios, evidência necessária, próximo passo. Se a memória não tem conteúdo operacional, pare e não crie lote. Um lote vazio é ruído no arquivo histórico para sempre.
2. **Atribua a identidade do lote.** Um id de lote estável no formato `<AAAA-MM-DD>-<n>`, em que `n` começa em 1 naquela data no fuso horário do projeto; se o id já existir em `<harness dir>/EPOCHAL.md`, avance `n`. Nunca reutilize um id para conteúdo diferente. Acrescente um timestamp em ISO 8601 com deslocamento de fuso horário. A data da consolidação não é a data dos eventos que ela contém; as datas originais ficam como foram escritas, e timestamps antigos desconhecidos continuam declarados desconhecidos em vez de inventados.
3. **Arquive.** Anexe ao fim de `<harness dir>/EPOCHAL.md` uma cópia literal completa do `<harness dir>/MEMORY.md` atual, com um delimitador inequívoco de início e de fim (uma linha começando com `<!-- INICIO-BRUTO` ou `<!-- BEGIN-RAW` e uma linha começando com `<!-- FIM-BRUTO` ou `<!-- END-RAW`, no formato `<!-- INICIO-BRUTO lote <id> -->` que os templates e o EPOCHAL existente usam, conforme o idioma do projeto), o id do lote, o timestamp, a fonte, e as tarefas e o intervalo de datas cobertos. Não resuma, reordene, reformate nem corrija em silêncio o texto bruto. Não carregue o arquivo histórico inteiro para anexar um lote; anexe ao fim.
4. **Mantenha o comentário fora do bloco bruto.** Notas do Coordenador e correções posteriores ficam fora do bloco arquivado, rotuladas como comentário e datadas. Uma correção é um novo registro ligado, nunca uma edição dentro do lote.
5. **Verifique antes de tocar na fonte.** Releia o lote escrito e compare-o na íntegra com a memória de onde veio, e tire o `dh memory hash` de novo: ele deve ser igual ao hash do passo 1. Em qualquer falha, divergência, hash alterado ou arquivo de destino inválido, deixe o `<harness dir>/MEMORY.md` intocado e reporte a situação. Este e a trava do passo 6 são os gates que protegem a memória.
6. **Poda, só depois de confirmado o arquivamento.** Escreva a memória podada com `dh memory write --expect <hash do passo 1>`; no código de saída 3 a memória mudou desde o passo 1, então deixe-a intacta, reporte a operação como incompleta com o id do lote armazenado, e deixe uma nova execução começar um novo lote a partir da memória atual. Mantenha a execução atual, as decisões vigentes, os bloqueios, a evidência necessária, o próximo passo e um histórico recente breve. Registre o id do lote e a data da consolidação na seção de última consolidação da memória. Os eventos encerrados continuam acessíveis no arquivo histórico. Um incidente grave aberto continua visível na memória como bloqueio, mesmo estando também registrado no arquivo de incidentes.
7. **Retome com segurança após uma interrupção.** A operação abrange dois arquivos e não é atômica. Localize o lote anterior pelo id estável antes de anexar qualquer coisa. Se estiver completo, verifique-o e confira que a fonte não mudou antes de podar; se a MEMORY já registra essa consolidação concluída, não pode de novo. Se o lote estiver parcial, malformado, ambíguo ou diferir da MEMORY atual, deixe a MEMORY intacta e reporte a recuperação como bloqueada. Não reinicie às cegas com um id novo, não duplique um lote nem edite conteúdo bruto histórico. Nunca remova da memória nada que não esteja provado preservado.
8. **Deixe o registro de incidentes em paz.** O `<harness dir>/RISKS.md` não é consolidado, podado, atualizado nem movido para o arquivo histórico. Os incidentes resolvidos e a prevenção deles permanecem intactos; atualizações de incidente pertencem a uma operação separada do Coordenador.
9. **Entregue.** Reporte o lote criado, o que permaneceu ativo na memória e qualquer item pendente. A consolidação é um comando explícito. O Coordenador pode sugeri-la quando a memória perde concisão, e nunca faz uma limpeza silenciosa.

## Formato de saída

Registro de handoff.

```
# Handoff: <id da tarefa>

## Objetivo
- Meta:
- Autorização: <o que pode ser escrito, commitado, publicado>

## Restrições
- <restrição, levada adiante literalmente>

## Conjunto de escrita
- <caminho relativo> | modificado | somente leitura

## Decisões
- <decisão> | motivo | fonte | rever se <condição>

## Evidência
| Checagem | Comando | Resultado |
| --- | --- | --- |
| <nome> | <comando do perfil do projeto> | passed / failed / `not-run` (<motivo>) |

## Pendentes e bloqueados
- <item> | responsável | bloqueio | condição de conclusão

## Próximo passo autorizado
- <um passo concreto> | autorização: já concedida (<fonte e escopo>) | pendente (<o quê>)
```

Reconciliação da retomada.

```
# Retomada: <id da tarefa>

## Reconciliação
| Item do handoff | Estado na árvore | Classificação | Ação |
| --- | --- | --- | --- |
| <item> | <observado> | aplicado / parcial / não iniciado / invalidado / contradito | <ação> |

## Mudanças externas detectadas
- <mudança> | fonte | efeito no plano

## Restrições ainda em vigor
- <restrição>

## Checagens rodadas de novo
| Checagem | Comando | Resultado |

## Próximo passo
- <passo> | autorização: já concedida (<fonte e escopo>) | pendente (<o quê>)
```

Relatório de consolidação.

```
# Consolidação: <id do lote>

- Timestamp: <ISO 8601 com deslocamento>
- Fonte: MEMORY.md
- Tarefas cobertas: <ids> | intervalo de datas: <intervalo conhecido>
- Verificação: lote relido e comparado na íntegra; fonte inalterada: sim / não
- Mantido ativo na memória: execução, decisões, bloqueios, evidência, próximo passo, histórico breve
- Pendente após a consolidação: <itens>
- RISKS.md: intocado
```

## Referência rápida

Política de leitura dos três registros.

| Arquivo | Escrito por | Lido quando |
| --- | --- | --- |
| `<harness dir>/MEMORY.md` | Só o Coordenador | Todo início de execução, toda retomada, antes de todo despacho |
| `<harness dir>/EPOCHAL.md` | Só o Coordenador | Só diante de uma necessidade histórica concreta ou de um pedido explícito; busque por tarefa, lote, data ou assunto e leia só a passagem |
| `<harness dir>/RISKS.md` | Só o Coordenador | Antes de mudar comportamento consolidado, uma regra crítica ou uma área de alto risco; e para registrar ou atualizar um incidente grave |

Gates da consolidação, em ordem. Um gate não atendido interrompe a operação com a memória intacta.

| Gate | Atendido quando |
| --- | --- |
| Pedido | Pedido explícito de consolidação, não inferido do tamanho da memória |
| Escritor único | O dono declarou nesta sessão que nenhuma outra sessão está escrevendo os registros |
| Conteúdo | A memória tem conteúdo operacional; senão nenhum lote é criado |
| Identidade | Id de lote `<AAAA-MM-DD>-<n>` ainda não presente no arquivo histórico, mais timestamp ISO 8601 com deslocamento |
| Arquivamento | Cópia literal completa anexada com delimitadores inequívocos |
| Verificação | Lote relido, comparado na íntegra, e a fonte confirmada como inalterada |
| Poda | Roda só depois de a verificação passar |
| Incidentes | O registro de incidentes está intocado |

## Erros comuns

| Desculpa | Realidade |
| --- | --- |
| "O handoff diz que foi feito, então foi feito" | O handoff é uma afirmação. Confirme nos arquivos antes de construir em cima dele. |
| "A árvore parece diferente, o handoff deve estar velho, vou recomeçar" | Parte dele foi aplicada. Classifique item por item; recomeçar repete trabalho aplicado e perde decisões. |
| "Essa restrição provavelmente era de uma fase anterior" | Uma restrição vale até o dono liberá-la. Largar uma em silêncio é exatamente a falha que a retomada deve impedir. |
| "A memória está longa, vou arrumá-la e arquivar depois" | Podar antes de o arquivamento ser verificado destrói a única cópia. Arquive, verifique, depois pode, nessa ordem. |
| "Resumir a memória antiga deixa o arquivo histórico legível" | O arquivo histórico é história bruta. Um resumo troca evidência por interpretação e não dá para desfazer. |
| "O timestamp da consolidação pode valer pelas datas dos eventos" | Não pode. As datas originais ficam; as desconhecidas continuam declaradas desconhecidas. |
| "A execução foi interrompida, o mais seguro é arquivar de novo" | Um lote duplicado corrompe o histórico. Localize e verifique primeiro o mesmo lote. |
| "Eu sou o papel de documentação, atualizar a memória faz parte do meu trabalho" | Só o Coordenador escreve os três registros. Reporte fatos e evidência em vez disso. |
| "Carregar o histórico inteiro me dá um contexto melhor" | O arquivo histórico entra no contexto só diante de necessidade concreta. Carregá-lo por padrão enterra o estado operacional. |
| "O incidente está resolvido, a consolidação pode limpá-lo" | A consolidação nunca toca o registro de incidentes. Incidentes resolvidos ficam com a prevenção deles. |
| "As checagens passaram antes da quebra, posso levar o resultado adiante" | Preserve-as só se as entradas avaliadas continuam inalteradas; senão rode de novo as checagens afetadas ou marque `not-run`. |
| "Nós dois podemos dividir a consolidação para ir mais rápido" | Uma sessão do Coordenador escreve por projeto. Escritores em paralelo perdem estado que nenhum arquivo histórico guarda. |

## Exemplo

Consolidação interrompida, retomada em uma sessão nova.

```
A memória tem: 2 bloqueios abertos, 1 decisão, 11 entradas de histórico.
O fim do arquivo histórico mostra o lote 2026-04-02-1, delimitadores presentes, marcador de fim presente.

1. Localize o lote 2026-04-02-1 no EPOCHAL.md.
2. Compare-o linha a linha com o MEMORY.md atual: idênticos.
3. Confirme que o MEMORY.md não mudou desde o timestamp do lote: sim.
   -> O passo de arquivamento foi concluído. Não crie um segundo lote.
4. Pode o MEMORY.md: mantenha os 2 bloqueios, a decisão, o próximo passo,
   3 entradas recentes de histórico; registre "última consolidação: 2026-04-02-1".
5. Reporte: lote 2026-04-02-1 verificado, 8 entradas de histórico agora só no arquivo histórico,
   2 bloqueios ainda ativos, RISKS.md intocado.
```

## Relacionado

Papéis: `coordinator`, `docs-guide`, `repo-scout`. Comandos: `/dh:handoff`, `/dh:resume`, `/dh:consolidate-memory`, `/dh:document`. Skills: `delivery-readiness`, `project-onboarding`, `repository-mapping`. Templates: a pasta `templates/` que acompanha o kit (`templates/MEMORY.md`, `templates/EPOCHAL.md`, `templates/RISKS.md`).

## Caso de prova

Uma sessão nova retoma uma tarefa cuja árvore foi alterada por outra pessoa depois que o handoff foi escrito. Ela nomeia a mudança externa, não reaplica uma mudança que já foi aplicada e leva adiante toda restrição declarada. Uma execução de consolidação produz um lote de arquivo histórico cujo conteúdo corresponde integralmente à memória anterior, com timestamp e fuso horário, e deixa na memória as decisões vigentes, os bloqueios e o próximo passo. Interromper a execução antes de o arquivamento ser verificado deixa a memória inalterada; interrompê-la depois de o lote ser escrito e retomar verifica esse mesmo lote em vez de criar um segundo. Uma execução de rotina termina sem que o arquivo histórico jamais entre no contexto, e o registro de incidentes é idêntico byte a byte antes e depois da consolidação.
