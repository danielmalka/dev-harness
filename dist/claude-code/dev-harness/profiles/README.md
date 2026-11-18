# Perfis

Convenções de projeto compartilhadas e portáteis. Sem segredos, sem caminhos absolutos, sem binários específicos de uma máquina.

Um perfil é um ponto de partida para o `.harness/project.yaml` de um projeto consumidor. O arquivo do consumidor não tem opção de idioma (ADR-011); ele pode trazer o opcional `conventions.code_language: en | pt-br` (idioma das mensagens de commit e dos comentários de código). O `/dh:setup` indica a correspondência mais próxima e copia só os campos que têm evidência no repositório. Quando nada corresponde, o perfil registrado é `none, custom`.

## Arquivos

| Arquivo | Quando se aplica |
| --- | --- |
| `base.yaml` | Todo projeto. Regras compartilhadas do harness. |
| `go-api.yaml` | Módulos Go que expõem uma API HTTP ou RPC. |
| `typescript-web.yaml` | Apps TypeScript com interface de usuário renderizada. |
| `php.yaml` | Projetos PHP gerenciados pelo Composer, com ou sem framework ou CMS. |
| `kotlin.yaml` | Serviços Kotlin na JVM (Gradle ou Maven), com seções Spring Boot, Quarkus ou Micronaut condicionadas a evidência no build. |
| `python.yaml` | Serviços e pacotes Python, com seções Django e FastAPI condicionadas a evidência nas dependências. |
| `typescript-api.yaml` | Serviços TypeScript ou Node sem UI renderizada. |
| `examples/project.yaml` | Formato de um `.harness/project.yaml` de consumidor. Não é um perfil de stack. |

Electron e outras stacks ficam de fora até um consumidor real precisar delas. O `php` entrou em 2026-09-20 para o primeiro consumidor PHP (um site Craft CMS). `kotlin`, `python` e `typescript-api` entraram em 2026-10-09 para os projetos de trabalho do dono, os primeiros consumidores Kotlin, Python e backend TypeScript.

## Esquema

```yaml
id: kebab-case-id          # deve coincidir com o nome do arquivo sem .yaml
name: Rótulo legível
summary: Uma frase
extends: base              # opcional; omitido no base
stack:                     # dicas usadas para casar com um repositório; não são instruções de instalação
  - go
commands:                  # null significa "não adivinhe"
  test: go test ./...
  lint: null
  build: go build ./...
  run: go run .
conventions:               # itens que o Coordenador e os especialistas devem seguir
  - ...
limits:                    # paradas obrigatórias
  - ...
```

Regras:

- `id` é igual ao nome do arquivo sem a extensão.
- Os valores de comando são copiados para o `.harness/project.yaml` só quando o mesmo comando é observado no repositório consumidor (Makefile, `package.json`, scripts do `go.mod`). Caso contrário, deixe o campo sem valor e rotule-o como desconhecido.
- `extends` mescla convenções e limites; um campo do filho substitui o campo do pai de mesmo nome. Não instala ferramentas.
- Os perfis nunca listam caminhos da máquina do autor, tokens nem servidores MCP.

## Correspondência

1. Leia a árvore do consumidor (manifestos, linguagens, presença de UI).
2. Escolha `go-api`, `typescript-web`, `php`, `kotlin`, `python` ou `typescript-api` quando a evidência de stack é direta.
   - `package.json` com UI renderizada presente (dependência de framework de UI como react, vue, svelte, next, nuxt ou angular, `index.html` ou uma pasta `public/`): `typescript-web`. Um serviço sem UI: `typescript-api`.
   - Monorepo com UI e API: uma regra só. O `setup` rodado na raiz escolhe `typescript-web` e registra em `notes:` do `project.yaml` a pasta da API e o perfil `typescript-api`; o `setup` rodado dentro da pasta da API escolhe `typescript-api`.
3. Caso contrário, use `base` e registre `none, custom` para os comandos específicos de stack.
4. As instruções existentes do projeto prevalecem sobre o perfil.

## Campos só do consumidor em `.harness/project.yaml`

Estes campos existem só no `.harness/project.yaml` de um consumidor, nunca em um perfil de stack `profiles/*.yaml` — `base.yaml`, `go-api.yaml`, `typescript-web.yaml`, `php.yaml`, `kotlin.yaml`, `python.yaml` e `typescript-api.yaml` nunca os trazem.

`reviewers:` é opcional. Ausente, o fluxo é só Claude e não há outra CLI para rotacionar. Quando presente, mapeia uma etapa para uma lista de revisores dessa etapa. Na etapa `code`, a primeira entrada de uma tarefa é `claude` quando a lista o inclui, senão a primeira entrada; cada entrada seguinte da mesma tarefa avança uma entrada. A etapa `security` continua rodando a lista inteira. A lista `code` inteira só roda quando o dono pede essa etapa.

```yaml
reviewers:                       # opcional; ausente = fluxo de hoje, só Claude
  document: [claude, "cli:agy/gemini-3.1-pro-high"]
  code: [claude, "cli:codex/gpt-5.6-sol"]
  security: [claude]
```

Etapas válidas: `document`, `code`, `security`. `verify` fica fora da v1.

Cada entrada da lista de uma etapa é de exatamente uma de três formas — não existe quarta forma:

- `claude` — despacha o agente do próprio kit para a etapa; o modelo já está fixado em `.agents/<role>.md`.
- `"cli:<binary>/<slug>"` — uma string; um revisor de CLI externo com o timeout padrão de 15 minutos.
- `{reviewer: "cli:<binary>/<slug>", timeout_minutes: <n>}` — um mapa; o mesmo revisor de CLI com um timeout personalizado.

`internal/kit/validate.go` não muda para este campo: ele valida a árvore de fontes do próprio kit (`.agents/`, `.skills/`, `.commands/`, `profiles/`), não o `.harness/project.yaml` de um consumidor, então não há nada a conferir aqui.
