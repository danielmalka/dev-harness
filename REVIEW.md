# Guia de revisão

## Como usar

Este arquivo diz ao revisor o que procurar neste projeto. O dono o preenche. O revisor o lê primeiro e liga cada achado a uma regra abaixo, ou diz que nenhuma regra se aplica. Mantenha poucas entradas e concretas.

## Regras

Uma entrada por problema recorrente. Copie o bloco para cada regra nova.

### `dist/` gerado editado à mão

- Sintoma: o diff muda um arquivo em `dist/` sem uma mudança correspondente em uma fonte canônica (`.agents/`, `.skills/`, `.commands/`, `templates/`, `profiles/`).
- Por que importa: `dist/` vem de `go run ./cmd/dh build`. Uma edição à mão se perde no próximo build, ou faz o CI falhar na checagem de diff do `dist`. A regra está no `AGENTS.md` ("Fontes canônicas").
- Exemplo: `dist/claude-code/dev-harness/agents/coordinator.md:2`

### Contagem nos docs deriva da contagem real

- Sintoma: um doc, README ou tutorial afirma uma contagem (13 agentes; 20 comandos em arquivos mais 1 registrado pelo mod) que não bate mais com as fontes.
- Por que importa: o dono confia nas contagens declaradas; o `AGENTS.md` as fixa ("Agentes (13)"). Um número desatualizado diz que o kit perdeu ou ganhou capacidade quando isso não aconteceu.
- Exemplo: `docs/tutorial.html:781`

### Caminho privado em um PRD, relatório de revisão ou doc

- Sintoma: o texto contém um caminho de home expandido ou outro caminho concreto de máquina, em vez de `<harness dir>` ou `<home>`.
- Por que importa: vaza a máquina do autor para um repositório compartilhado. O `dh validate` (`checkPrivatePaths`) o rejeita e bloqueia a entrega.
- Exemplo: `internal/kit/validate.go:577`

### Mudança deixa referências obsoletas

- Sintoma: um arquivo, contagem, link ou cópia é removido ou renomeado, mas referências a ele permanecem em docs ou templates.
- Por que importa: o `AGENTS.md` ("Premissa de limpeza") bloqueia a entrega quando uma mudança deixa sobras obsoletas.
- Exemplo: `AGENTS.md:12`

## Ignorar

Uma entrada por padrão que é ruído. Cada entrada diz "ignore este padrão neste tipo de arquivo". Um achado que casa com uma entrada é reportado como ignorado por regra e não gasta uma rodada.

- Ignore a linha `generated_at` em `dist/claude-code/dev-harness/harness-manifest.json`. Ela difere entre builds, e o CI a tolera.
- Ignore um achado dentro de uma cópia em `dist/` quando o mesmo achado existe na fonte canônica dela. Revise a fonte, não a cópia (por exemplo `dist/claude-code/dev-harness/agents/coordinator.md` contra `.agents/coordinator.md`).

## Caminhos críticos

Globs de caminho no estilo `.gitignore`. Um pull request que toca um caminho que casa para antes do merge, e o dono decide. Deixe a lista vazia para não ter parada.

- `adapters/claude-code/plugin/hooks/**`: trava em tempo de execução da trava de memória e da regra de despacho.
- `.skills/external-clis/**`: isolamento dos revisores de CLI externos.
