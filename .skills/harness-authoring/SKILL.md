---
name: harness-authoring
description: Use quando um agente, skill, comando ou tutorial deste kit precisar ser criado ou alterado, quando um gatilho disparar de forma ampla ou estreita demais, quando as instruções de um ativo forem ambíguas ou contradisserem as de outro ativo, ou quando o conteúdo do kit exigir uma nota de procedência, versão ou compatibilidade. Não use para mudanças no código do próprio projeto consumidor, para uma preferência pontual de tom, nem para editar à mão o pacote gerado em vez das fontes.
author: malka
metadata:
  provenance: adapted
  sources: ["skill-development", "writing-for-agents", "agent-development", "command-development", "writing-skills", " skill-creator role", "dev-harness harness-maintainer agent"]
---

# Autoria do harness

## Visão geral

Este procedimento produz ou muda um ativo deste kit: um agente, uma skill, um comando ou um tutorial. Ele fixa o contrato que cada tipo de ativo deve, mantém o texto de disparo separado do texto de procedimento e empurra o material pesado para trás de um ponteiro, para que o arquivo carregado continue legível. Um ativo novo precisa de um caso de uso autorizado e de expectativas testáveis. Um reparo precisa de uma falha reproduzida ou de uma contradição concreta de contrato; uma resposta preferida isolada não é evidência de melhora.

## Quando usar

- Um agente, skill, comando ou tutorial novo é necessário e o contrato dele precisa ser escrito certo de primeira.
- Um gatilho dispara em trabalho em que não deveria tocar, ou deixa de disparar em trabalho que é dele.
- Dois ativos dão instruções conflitantes, ou um ativo duplica outro em silêncio.
- A descrição de um ativo resume o fluxo de trabalho dele, e os consumidores agem pela descrição e pulam o corpo.
- O arquivo carregado cresceu além do que precisa e o material de referência deve passar para trás de um ponteiro.
- Faltam ou estão erradas as informações de procedência, licença ou versão do conteúdo do kit.

## Quando não usar

- A mudança pertence ao código de aplicação de um projeto consumidor ou ao arquivo de instruções dele.
- A única evidência é uma resposta isolada de que alguém gostou mais. Colete antes um caso reproduzível.
- O pedido é avaliar se uma mudança ajudou. Isso é `harness-evaluation`.
- A intenção é remendar o pacote gerado em `dist/`. Mude as fontes.

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Um caso de uso concreto e o comportamento esperado; para reparos, a falha observada ou as instruções conflitantes | Derive um caso do pedido autorizado. Para um ativo novo, o comportamento atual pode não existir. Para uma contradição estática, cite as duas instruções e o cenário afetado; não afirme uma reprodução em modelo que você não rodou. |
| O ativo que é dono do comportamento | Rastreie-o: qual gatilho disparou, qual arquivo foi carregado, qual instrução produziu a saída. Nomeie o dono antes de editar qualquer coisa. |
| A versão do kit ou a identidade do checkout do harness | Registre o estado do checkout de que você partiu e marque a versão como desconhecida em vez de adivinhá-la. |
| Resultados de linha de base existentes para o ativo | Capture o estado anterior. Se rodadas em modelo não estiverem disponíveis, distinga a correção estática da melhora comportamental não medida e reporte esse limite. |
| Procedência de qualquer material emprestado: origem e licença | Não o incorpore. Nomear uma fonte não torna o conteúdo dela disponível. Escreva o procedimento você mesmo ou deixe a lacuna registrada. |
| Idioma dos fontes do kit | O kit é só pt-BR (ADR-011): prosa, títulos, `description`, exemplos e mensagens saem em pt-BR; fica em inglês só o identificador (chaves, comandos, nomes de agente e de arquivo, slugs de skill, código), entre crases quando citado no texto. Comentários de código e mensagens de commit do kit seguem o que o repositório define (`AGENTS.md`, `CONTRIBUTING.md`: inglês). |

## Procedimento

1. **Estabeleça o caso no checkout do harness.** Trabalhe nas fontes. Para um reparo, reproduza a falha ou cite o contrato contraditório e um cenário concreto. Para um ativo novo, registre o comportamento autorizado que hoje falta. Preserve o estado anterior e distinga rodadas observadas de análise estática.
2. **Localize o ativo dono.** Siga o caminho que o runtime percorreu: qual descrição casou, qual corpo foi carregado, qual instrução a saída seguiu. Uma falha de roteamento pertence a uma descrição; uma falha de método pertence a um corpo; uma falha de escopo pertence a um limite.
3. **Escolha o tipo de ativo pelo que a mudança precisa fazer.** Um agente é uma parte responsável com contexto e limites próprios. Uma skill é o procedimento reutilizável. Um comando é o ponto de entrada do usuário e carrega só roteamento, parâmetros e contrato, com o método morando na skill. Se a correção precisa de um tipo novo, diga isso explicitamente em vez de esticar um existente.
4. **Escreva e guarde os casos de teste antes de editar qualquer coisa.** Ao menos um caso positivo e um negativo em `evals/cases/`, para que a `harness-evaluation` capture a linha de base no ativo inalterado. Rodá-los e avaliá-los é a `harness-evaluation`.
5. **Escreva ou repare o contrato desse tipo.** Use as tabelas de contrato na Referência rápida. Todo campo ou está preenchido ou está marcado explicitamente como não aplicável, com o motivo.
6. **Escreva a descrição só com condições de disparo.** Terceira pessoa, começando pelas condições que devem carregar o ativo, mais as contraindicações. Nunca resuma o procedimento na descrição: uma descrição que descreve o fluxo de trabalho vira um atalho que o consumidor segue em vez de ler o corpo. Mantenha-a abaixo de 500 caracteres.
7. **Ajuste a forma da orientação à forma da falha.** Use a tabela de falha para forma. Uma regra pulada sob pressão precisa de uma contra-regra explícita e de uma tabela de racionalizações. Saída com a forma errada precisa de uma receita positiva que diga o que a saída é, não de uma lista de proibições. Um elemento ausente precisa de um campo obrigatório no template. Um comportamento que depende de uma condição precisa de um condicional atrelado a algo observável.
8. **Prefira o alvo positivo à proibição.** Nomear o comportamento proibido o torna mais disponível, não menos. Mantenha uma proibição só como proteção rígida que não pode ser formulada de forma positiva, e a combine com o comportamento-alvo.
9. **Aplique divulgação progressiva.** O arquivo carregado contém o que todo caminho por ele precisa. O material que só alguns caminhos alcançam, e qualquer bloco de referência com mais de cerca de 100 linhas, vai para `references/<topic>.md` dentro da pasta da própria skill e é alcançado por um link relativo à pasta, que resolve tanto na árvore de fontes (`.skills/`) quanto no pacote gerado (`skills/`). Mantenha cada significado em um lugar só; a mesma regra em dois arquivos é um defeito de manutenção e infla a importância aparente dela.
10. **Mantenha a mudança mínima e os identificadores estáveis.** Mude a menor coisa que corrige o caso reproduzido. Não reformate ativos vizinhos na mesma passada. Renomear um identificador é uma mudança incompatível e precisa de motivo próprio.
11. **Registre a procedência no frontmatter.** Preserve `author: malka`. Marque o ativo como `authorial` quando o procedimento é seu, ou como `adapted` quando ele sintetiza material externo em algo diferente. Registre as origens identificáveis e verifique as licenças aplicáveis antes de incorporar texto protegido. A política deste kit proíbe skills copiadas literalmente mesmo quando licenciadas; atribuição sozinha não faz de uma cópia uma adaptação.
12. **Declare a compatibilidade e proponha uma versão.** Nomeie o que quebra para um consumidor que já usa o kit: um identificador renomeado, uma forma de saída alterada, um gatilho estreitado, uma referência removida. Proponha a versão e marque-a como proposta.
13. **Regenere o pacote a partir das fontes.** Nunca edite prompts em `dist/` de forma independente. A partir da raiz do kit, rode `go run ./cmd/dh build` (ou `bin/dh build` a partir de um pacote construído). Ele valida as fontes, monta um pacote completo e então substitui `dist/claude-code/dev-harness`. Se fonte e conteúdo gerado discordarem, a fonte vence e o build deve ser rodado de novo.
14. **Reporte, não registre.** Caso, mudança, evidência e compatibilidade voltam ao Coordenador, que é o único que escreve `.harness/MEMORY.md`, `.harness/EPOCHAL.md` e `.harness/RISKS.md`.
15. **Escreva cada passo do procedimento como uma frase imperativa mais, no máximo, uma frase de motivo.** Ponha as referências cruzadas no fim do passo ou da seção, não dentro de uma frase.

## Formato de saída

```
# Mudança de ativo: <id do ativo>

## Caso
- Prompt:
- Esperado:
- Observado:
- Reproduzido: sim / não (<motivo>)

## Ativo dono
- Arquivo: <caminho relativo> | tipo: agente / skill / comando / tutorial
- Por que este ativo é dono do comportamento:

## Mudança
- O que mudou (menor edição):
- O que deliberadamente não mudou:

## Checagem do contrato
| Campo | Estado |
| --- | --- |
| <campo da tabela de contrato> | preenchido / não aplicável (<motivo>) |

## Casos de teste
- Positivo: <caminho relativo em evals/cases/> | agora deve <comportamento>
- Negativo: <caminho relativo em evals/cases/> | ainda deve <comportamento>

## Procedência
- Classe: authorial / adapted
- Origem:
- Licença verificada: sim / não / não aplicável

## Compatibilidade
- Quebra para consumidores existentes:
- Identificadores renomeados:

## Proposta de versão
- <valor> (proposta, precisa de confirmação)

## Evidência
- <caminho relativo> - <o que estabelece>
```

## Referência rápida

Contrato por tipo de ativo. Todo campo é preenchido ou marcado como não aplicável, com o motivo.

| Ativo | Contrato exigido |
| --- | --- |
| Agente | Missão; gatilhos positivos; gatilhos negativos; entradas mínimas e comportamento diante de entrada ausente; procedimento; política de ferramentas (omita `tools` para herdar tudo; lista de permissão só para papéis somente leitura; nada de `disallowedTools` como trava — o runtime não a aplicou (RISK-001, 2026-09-21); a trava é o mod runtime-guard do plugin (ADR-005), e o corpo do prompt de todo especialista ainda declara que só o Coordenador despacha, para os runtimes sem o mod); limites; formato de saída; handoff de contexto; condição de parada; exemplos avaliáveis; modelo padrão explícito (nunca herdar), com sobrescrita delimitada pelo Coordenador e o modelo efetivo registrado ou `unverified` |
| Skill | Frontmatter com name igual à pasta, descrição de disparo, author, procedência; visão geral com o princípio central; quando usar; quando não usar; entradas; procedimento numerado e verificável; formato de saída; referência rápida; erros comuns; quando o caso reproduzido é que o agente conhece a regra e a pula sob pressão, uma tabela de racionalizações de três colunas (o que você vai ficar tentado a pensar, por que está errado, o que fazer) além dos erros comuns; um exemplo desenvolvido; agentes, comandos e skills irmãs relacionados; caso de prova |
| Comando | Roteamento para o papel e a skill donos; parâmetros; contrato da saída esperada; comportamento diante de pré-requisito ausente; próxima ação. O método fica na skill |
| Tutorial | Pré-requisitos; o exercício; a condição de conclusão; o que o leitor deve ser capaz de fazer sem ajuda no fim |

Ajuste a forma da orientação à forma da falha.

| Falha na linha de base | Forma certa | Forma errada |
| --- | --- | --- |
| Conhece a regra, pula sob pressão | Contra-regra explícita mais uma tabela de racionalizações | Redação suave como prefira ou considere |
| Cumpre, mas a saída tem a forma errada | Receita positiva nomeando as partes da saída em ordem | Uma lista de proibições |
| Omite um elemento obrigatório de algo já produzido | Um campo obrigatório no template | Um lembrete em prosa perto do template |
| O comportamento deve depender de uma condição | Condicional atrelado a um predicado observável | Regra incondicional com cláusulas de exceção |
| Dispara em trabalho que não é dele | Afiar as contraindicações na descrição | Acrescentar ressalvas no corpo, que carrega tarde demais |

Onde o material mora.

| Material | Local |
| --- | --- |
| Condições de disparo | Só a descrição do frontmatter |
| Procedimento de que todo caminho precisa | O arquivo carregado |
| Referência que só alguns caminhos alcançam, ou com mais de cerca de 100 linhas | `references/<topic>.md` dentro da pasta da própria skill, com link relativo à pasta |
| Algo que um script pode verificar de forma mecânica | Um script, não prosa |
| Um fato que o ambiente já declara | Deixe-o no ambiente; uma repetição fica velha |

## Erros comuns

| Erro | Por que prejudica | Em vez disso |
| --- | --- | --- |
| Editar um ativo porque uma resposta ficou melhor | As regras do kit são reescritas por gosto e se afastam dos casos para os quais foram feitas | Exija um caso reproduzido com prompt, esperado e observado |
| Descrição que resume o fluxo de trabalho | Os consumidores seguem o resumo e nunca leem o corpo, que vira decoração | Só condições de disparo e contraindicações |
| Editar à mão o pacote gerado | O próximo build reverte a correção em silêncio e as fontes nunca ficam sabendo | Mude as fontes e reconstrua |
| Corrigir o roteamento acrescentando texto ao corpo | O corpo só carrega depois de a decisão de roteamento já ter sido tomada | Corrija a descrição, que é o que a decisão lê |
| Grande reescrita quando uma cláusula estava errada | A mudança não pode ser atribuída ao caso, então ninguém sabe o que corrigiu o quê | A menor edição que corrige o caso reproduzido |
| Renomear identificadores enquanto corrige comportamento | Duas mudanças chegam juntas e toda referência de consumidor quebra por um motivo não relacionado | Mantenha os identificadores estáveis, salvo se a renomeação constituir a correção |
| Copiar texto de outra fonte porque é bom | Licença não verificada e referências não resolvidas viajam com ele | Authorial ou adapted, com origem e licença registradas |
| Nomear uma ferramenta externa que o kit não distribui | O ativo quebra em qualquer máquina que não a tenha | Formule como disponível-ou-alternativa e nomeie a alternativa |
| A mesma regra escrita em dois ativos | Eles divergem, e a contradição aparece no meio da tarefa | Uma fonte única de verdade, referenciada por caminho relativo |
| Caminhos absolutos da máquina do autor | O kit deixa de ser portátil no momento em que se move | Caminhos relativos à raiz do kit ou ao projeto consumidor |
| Só um caso de teste positivo | Um gatilho ampliado passa por ele e começa a capturar trabalho que não é dele | Sempre o combine com um caso negativo |

## Exemplo

Entrada: o papel de mapeamento explora a árvore inteira para uma pergunta sobre um arquivo conhecido.

```
Caso
  Prompt:   "O que src/auth/token.ts exporta?"
  Esperado: ler esse arquivo, responder, parar.
  Observado: mapa da árvore inteira, 40 arquivos lidos.

Ativo dono: .agents/repo-scout.md, campo description.
  O corpo já diz "menor mapa suficiente". A descrição
  casou antes de o corpo ser carregado, então o corpo não podia ajudar.

Mudança (uma linha, só a descrição):
  acrescentar "Não use quando o usuário nomear um único arquivo conhecido e a
  pergunta puder ser respondida a partir dele."

Casos de teste
  Positivo: pergunta ampla sem arquivo nomeado -> o papel ainda carrega.
  Negativo: o prompt acima -> o papel não carrega; o arquivo é lido diretamente.

Compatibilidade: nenhum identificador mudou, gatilho estreitado. Proposta: versão patch.
```

## Relacionado

Papéis: `harness-maintainer`, `coordinator`, `qa-verifier`. Comando: `/dh:improve`. Skills: `harness-evaluation`, `context-handoff`, `delivery-readiness`. Templates: a pasta `templates/` que acompanha o kit (`templates/MEMORY.md`).

## Caso de prova

Dada uma falha de roteamento reproduzida, o fluxo nomeia o ativo dono e o campo responsável, faz uma única edição mínima no checkout do harness e deixa o pacote gerado intocado até que ele seja reconstruído a partir das fontes. O ativo termina com uma descrição que declara só condições de disparo, um contrato completo para o tipo dele, um caso de teste positivo e um negativo guardados para repetição, uma classe de procedência com a origem, e um impacto de compatibilidade declarado com uma versão proposta. Um pedido para mudar um ativo com base em uma única resposta preferida e isolada é recusado com o motivo, e nenhuma regra é editada. Um caso reproduzido de conhecer a regra e pulá-la sob pressão deixa a skill com a tabela de racionalizações de três colunas além dos erros comuns. Novos casos de avaliação não são acrescentados a menos que o dono peça; o caso de prova dentro da skill é o substituto estático.
