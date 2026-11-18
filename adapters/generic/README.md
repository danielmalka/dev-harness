# Adapter genérico

Como usar este kit com um runtime de IA que não é o Claude Code. Não há integração nativa. Não anuncie uma.

## O que é transferido

| Ativo | Como carregar |
| --- | --- |
| `.agents/coordinator.md` | Instruções de sistema da sessão principal. |
| Outros arquivos de `.agents/` | Prompts de especialista. Cole ou anexe um papel por tarefa delegada. Passe a família do modelo como pedido; o runtime pode ignorá-la. |
| `.skills/<id>/SKILL.md` | Procedimento a carregar quando esse trabalho começar. Siga os links relativos dentro da pasta da skill. |
| `.commands/<id>.md` | Receita voltada ao usuário. Digite os passos; não há menu de comandos de barra. |
| `templates/` | Copie para o `.harness/` do consumidor quando a configuração estiver autorizada. |
| `profiles/` | Mesmas regras de correspondência do Claude. |

## O que não é transferido

- Os nomes de comando `/dh:*`.
- `--plugin-dir` e `plugin.json`.
- As famílias de `model` do Claude como garantia de qual modelo vai rodar.
- A ferramenta Agent. Turnos sequenciais de especialista substituem o despacho paralelo. Mantenha o Coordenador como o único a escrever em MEMORY, EPOCHAL e RISKS.
- Listas de permissão de ferramentas no frontmatter, a menos que o outro runtime tenha um equivalente.

## Sessão manual

1. Abra o projeto consumidor.
2. Carregue `coordinator.md` como as instruções principais.
3. Rode os passos de diagnóstico em `.skills/project-onboarding/SKILL.md`. Rode também `bin/dh doctor` a partir de um pacote compilado (não precisa de interpretador).
4. Para cada tarefa, siga o arquivo de comando correspondente, depois a skill que ele nomeia, depois o agente especialista que ele nomeia.
5. Registre como instruções os limites que o runtime não consegue impor (revisão somente leitura, sem Bash para revisores) e confira os arquivos escritos depois.

Condição de conclusão: uma fatia de fixture (`evals/fixtures/slice-01`) concluída com evidência, com toda limitação do runtime hospedeiro declarada no handoff. Não afirme paridade de funcionalidades com o Claude Code.
