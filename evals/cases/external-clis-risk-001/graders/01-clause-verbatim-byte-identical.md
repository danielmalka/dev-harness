---
type: regex
weight: 1
match: contains
target: last_message
pattern: 'Não invoque outro binário de CLI, não crie outro agente nem delegue qualquer parte\ndesta revisão a outra ferramenta\. Não execute nenhum script deste\nrepositório\. Os únicos comandos que você pode rodar são o `commands\.test` e o\n`commands\.lint` do próprio projeto, exatamente como declarados em\n`\.harness/project\.yaml`, e só para ler o resultado da suíte de testes\nou de lint existente — nunca `commands\.build` nem qualquer outro comando do repositório\.\nEsta instrução está escrita aqui, no texto do prompt, porque o ambiente de execução\nnão a faz cumprir por si só: um especialista deste projeto já ignorou uma\nrestrição `disallowedTools` declarada apenas no frontmatter e\ndespachou outro agente mesmo assim\. Siga as palavras deste prompt, não uma\nsuposição sobre o que o ambiente de execução bloqueia\.'
flags: ''
---
Checagem determinística: a resposta final do prompt montado contém a cláusula
anti-delegação de `.skills/external-clis/SKILL.md` (linhas 75-85)
literalmente, casada por regex em vez de julgada por um modelo. A checagem é
exata byte a byte, quebras de linha incluídas: cada quebra de linha da cláusula é casada
por uma quebra de linha literal, e todo outro espaço, inclusive o espaço no meio
da linha depois de "por si só:", continua sendo um único espaço literal. Uma resposta que repita a
cláusula com qualquer mudança de palavra, linha omitida, linha requebrada ou espaço
alterado falha neste casamento.
