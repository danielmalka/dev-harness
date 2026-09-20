# dev-harness

English: [README.en.md](README.en.md)

Kit portátil de desenvolvimento assistido por IA. Clone, carregue no Claude Code, desenvolva.

Licença MIT. Plano de produto: [`docs/plano-produto.html`](docs/plano-produto.html). Início rápido: [`docs/inicio-rapido.html`](docs/inicio-rapido.html).

## Carregar

Pré-requisitos: Git, Claude Code autenticado com acesso a um modelo, Python 3 (validador, doctor e fixture).

No projeto em que você vai trabalhar:

```bash
claude --plugin-dir "/caminho/absoluto/do/clone/dist/claude-code/dev-harness" --agent coordinator
```

Substitua o caminho pelo resultado de `pwd` no clone, mais `/dist/claude-code/dev-harness`. Atenção: dentro de aspas o `~` não expande e o Claude ignora o plugin em silêncio (`--agent coordinator` falha com "not found"). Use `$HOME/...` ou o caminho completo. Se o runtime exigir o nome qualificado, use `--agent dev-harness:coordinator`.

Depois, na sessão:

```text
/dev-harness:doctor
/dev-harness:setup
```

O diagnóstico mecânico também roda sem o Claude:

```bash
bash dist/claude-code/dev-harness/scripts/doctor.sh dist/claude-code/dev-harness
```

Tutorial 00: [`docs/tutoriais/00-primeira-maquina.html`](docs/tutoriais/00-primeira-maquina.html). Fixture da primeira fatia: [`evals/fixtures/slice-01`](evals/fixtures/slice-01).

## Fontes e pacote

| Caminho | Papel |
| --- | --- |
| `.agents/` `.commands/` `.skills/` `templates/` `profiles/` | Fontes canônicas |
| `scripts/build-claude-code.sh` | Gera `dist/claude-code/dev-harness` |
| `dist/` | Pacote gerado. Não editar. |
| `adapters/claude-code/` | Mapeamento para o Claude Code |
| `adapters/generic/` | Uso manual em outra IA, sem paridade |

Templates de trabalho em `templates/en/` e `templates/pt-br/` (versões espelhadas): PRD, Story, Task/Bug, ADR e RFC, além dos três registros de memória (MEMORY, EPOCHAL, RISKS). O `setup` grava `language` em `.harness/project.yaml` e escolhe a pasta. Quando usar cada um: [`docs/tutoriais/templates.html`](docs/tutoriais/templates.html).

Perfis (`base`, `go-api`, `typescript-web`, `php`) estão em inglês em `profiles/`.

## Idioma

- O Coordenador responde no idioma em que você escreve. Inglês é o padrão quando não há sinal; escreva em português e ele responde em português.
- No primeiro `/dev-harness:setup`, o idioma da sessão é gravado em `.harness/project.yaml` como `language: en` ou `language: pt-br`. Esse campo escolhe `templates/<lang>/` para os registros de memória e os artefatos de trabalho (PRD, story, task, ADR, RFC). Para trocar, edite o campo ou peça ao Coordenador; registros já escritos não são traduzidos.
- Entre agentes tudo é inglês: despachos, respostas dos especialistas, código, commits e casos de eval. `project.yaml` e perfis também.
- Documentos HTML gerados pela skill `doc-template-html` seguem o mesmo campo: `bash .skills/doc-template-html/scripts/stamp.sh --lang en ...` (padrão `pt-br`).
- Documentação do kit: pt-br em `docs/` e inglês em `docs/en/`; cada página tem link para a outra versão. Este README tem versão em [`README.en.md`](README.en.md). O plano de produto é interno e existe só em pt-br.

## Manutenção

```bash
python3 scripts/validate.py --source-only .
bash scripts/build-claude-code.sh
python3 scripts/validate.py .
```
