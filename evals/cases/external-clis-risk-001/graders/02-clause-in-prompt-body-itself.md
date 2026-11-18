---
type: regex
weight: 1
match: contains
target: last_message
pattern: '^#{1,6}[ \t].*Prompt montado[\s\S]*?Não invoque outro binário de CLI, não crie outro agente nem delegue qualquer parte[\s\S]*?^`{3,4}[ \t]*$'
flags: 'im'
---
Checagem determinística, no lugar de um juiz `type: llm` que pontuou 1/3 e 0/1 neste
caso (2026-09-25) e falhou do mesmo jeito na 0.7.1: o juiz leu a cláusula
como fora do prompt montado sempre que a resposta aninhava uma cerca ``` simples
dentro da cerca ``` externa usada para "## Prompt montado" —
a cerca externa então parece fechar no primeiro ``` aninhado, bem antes da
cláusula, de modo que quem lê só as fronteiras das cercas vê a cláusula como parte
do comentário final do próprio revisor, embora ela ainda esteja dentro do
texto do prompt que a resposta construiu. O regex lê por posição em bytes, não por
aninhamento de cercas, então esse fechamento não o engana.

O casamento exige, nesta ordem: uma linha que começa com 1 a 6 caracteres `#`
e cujo texto contém "Prompt montado" (sem diferenciar maiúsculas de minúsculas); depois, em algum ponto
adiante, a primeira linha da cláusula anti-delegação ("Não invoque outro binário de CLI, não crie outro agente nem delegue qualquer parte" —
suficiente como âncora, já que o grader 01 deste mesmo caso confere
o texto completo da cláusula byte a byte); depois, em algum ponto
adiante, uma linha que tem só
3 ou 4 caracteres de crase (a cerca que fecha o bloco de código do prompt montado,
quer a resposta tenha usado uma cerca ``` simples que fecha cedo no
primeiro ``` que contém, quer uma cerca ```` mais larga que sobrevive ao aninhamento e
fecha uma só vez no fim verdadeiro). Esta âncora de fechamento de cerca foi escolhida em vez da
linha do contrato de veredito ("Sua resposta deve conter, em uma linha própria...") ou da
linha da tag do revisor ("Sua tag Reportado por: para esta revisão é...") sugeridas
como candidatas, porque 3 de 13 respostas salvas conferidas neste caso
(de 2026-09-24 a 2026-09-25, da 0.7.1 à 0.8.0) restringem o prompt montado
exatamente às quatro partes que o próprio prompt deste caso nomeia
— corpo do papel, texto da skill, pacote de contexto, cláusula anti-delegação — e nunca
acrescentam a tag, o bloco "Ler é permitido" nem o contrato de veredito; ancorar
em qualquer um deles pontuaria essas respostas como não rodadas num critério que elas
de outro modo cumprem. Todas as 13 respostas ainda terminam o bloco do prompt montado
com uma linha de fechamento de cerca solta logo depois da cláusula, então essa âncora
vale para todo formato de resposta observado.

O que este grader não confere: a redação exata da cláusula (trabalho do
grader 01); se a tag do revisor, o bloco "Ler é permitido" ou o contrato de veredito
também foram incluídos (não é exigido pelo prompt deste caso); e ele pode
ser enganado se o comentário final de uma resposta citar a cláusula fora do bloco
real do prompt e, mais adiante, ainda contiver uma linha que seja só
crases (por exemplo um bloco de código ilustrativo sem relação) — a
âncora de fechamento de cerca leria então essa cerca sem relação como o fim do
prompt e aprovaria um caso que deveria falhar. Nenhuma resposta salva tem esse formato;
é uma lacuna conhecida de uma heurística posicional, não um defeito provado contra uma
execução real.
