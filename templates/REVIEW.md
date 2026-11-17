# Guia de revisão

## Como usar

Este arquivo diz ao revisor o que procurar neste projeto. O dono o preenche. O revisor o lê primeiro e liga cada achado a uma regra abaixo, ou diz que nenhuma regra se aplica. Mantenha poucas entradas e concretas. A revisão lê este arquivo da branch base, então uma mudança nele para o merge e vai ao dono.

## Regras

Uma entrada por problema recorrente. Copie o bloco para cada regra nova.

### <Título da regra>

- Sintoma: <o que o revisor vê no diff>
- Por que importa: <o que quebra ou o que custa>
- Exemplo: `<caminho/do/arquivo.ext>:<linha>`

## Ignorar

Uma entrada por padrão que é ruído. Cada entrada diz "ignore este padrão neste tipo de arquivo". Um achado que casa com uma entrada é reportado como ignorado por regra e não gasta uma rodada.

- Ignorar <padrão> em <tipo de arquivo>.

## Caminhos críticos

Globs de caminho no estilo `.gitignore`. Um pull request que toca um caminho que casa para antes do merge, e o dono decide. Deixe a lista vazia para não ter parada. Cada entrada é um glob entre crases; o texto depois dele é o motivo, não faz parte do glob.

- `src/payments/**`
