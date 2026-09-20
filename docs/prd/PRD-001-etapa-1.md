# PRD-001 · Etapa 1: distribuição pelo marketplace e observabilidade de sessões e agentes no plugin

| Campo | Valor |
|---|---|
| Status | aprovado (grill de 20/09/2026, rodadas 1 a 3) |
| Dono | Daniel Lemos |
| Criado / atualizado | 2026-09-20 / 2026-09-20 |
| Stories | a derivar: ST-001 marketplace, ST-002 binário `dh`, ST-003 snapshot e hooks, ST-004 prova no danlemos |

## 1. Problema

Hoje o kit só entra num projeto por clone e `--plugin-dir`, e a única visão do que os agentes estão fazendo é o painel embutido do Claude Code, que some ao fechar a sessão e não cobre várias sessões. O dono quer instalar sem clonar e ver, numa lista, quantos agentes trabalham, com qual modelo, quantos tokens entram e saem e se cada um está ativo ou ocioso. A extensão VS Code (`dev-harness-vscode`) será só um visualizador; ela precisa de dados que o plugin ainda não produz.

## 2. Resultado esperado

- Sinal 1: instalar o plugin num projeto novo sem clone · hoje: impossível · alvo: `/plugin marketplace add danielmalka/dev-harness` + `/plugin install dev-harness@dev-harness` funcionam em máquina limpa.
- Sinal 2: snapshot por sessão em disco · hoje: nada · alvo: `~/.claude/dev-harness/sessions/<session_id>.json` atualizado a cada tick e a cada hook, com sessão e subagentes.
- Sinal 3: latência do script de snapshot · hoje: n/a · alvo: abaixo de 50 ms por chamada.

## 3. Usuários e cenários

| Usuário | Cenário | Hoje | Depois |
|---|---|---|---|
| Dono, num projeto consumidor | Quer o kit sem clonar | Clona e aponta `--plugin-dir` | Adiciona o marketplace e instala; o projeto pode fixar `enabledPlugins` |
| Dono, durante um `build` | Quer saber quem está rodando, em que modelo e quanto custa | Olha o painel do Claude, que some depois | O painel de subagentes do Claude mostra as linhas do kit e o snapshot fica em disco para a extensão |
| Extensão VS Code | Precisa listar sessões e agentes de várias janelas | Nada | Lê os JSON de `~/.claude/dev-harness/sessions/` |

## 4. Escopo

**Entra**
- RF-01 `marketplace.json` na raiz do repositório dev-harness, com o plugin apontando para `dist/claude-code/dev-harness` do próprio repo e versão vinda de `harness-manifest.json`.
- RF-02 Binário Go `dh` (ADR-001) com subcomandos `snapshot`, `validate`, `build` e `doctor`; builds pré-compilados em `dist/claude-code/dev-harness/bin/<os>_<arch>/dh` para linux amd64 e arm64, darwin amd64 e arm64, windows amd64; wrapper `bin/dh` que escolhe pela plataforma; checksum de cada binário no manifesto.
- RF-03 `subagentStatusLine` embarcado no `settings.json` do plugin chamando `dh snapshot subagents`: grava as linhas de subagentes no snapshot e devolve as linhas formatadas para o painel do Claude (nome · modelo · tokens · estado).
- RF-04 Hooks embarcados em `hooks/hooks.json` do plugin: `SessionStart`, `UserPromptSubmit`, `Stop`, `SubagentStart`, `SubagentStop`, `SessionEnd`, todos chamando `dh snapshot event`, que grava a transição de estado no snapshot da sessão.
- RF-05 `statusLine` opcional, documentada como opt-in: `dh snapshot statusline` imprime uma linha compacta (modelo · contexto % · custo · agentes ativos) e grava os campos da sessão principal no snapshot. Quem já tem status line encadeia ou ignora.
- RF-06 Retenção: snapshot com `state: closed` no `SessionEnd`; `dh snapshot prune` remove arquivos com mais de 7 dias; a extensão pode chamar o prune.
- RF-07 `dh doctor` reporta: binário presente para a plataforma, marketplace resolvido, hooks e status lines ativos, diretório de snapshots gravável, `.harness/` ignorado pelo git do projeto (aviso).
- RF-08 Tutorial 00 e início rápido ensinam o marketplace como caminho principal e `--plugin-dir` como caminho de desenvolvimento; documentam `enabledPlugins` para times.

**Não entra**
- OTEL, collector, histórico de custo ao longo do tempo. Decisão adiada pelo dono.
- Agent SDK ou qualquer sessão controlada pela extensão. A extensão não é runtime.
- Sessões de Grok, Codex ou OpenCode. Só Claude Code na v1.
- A extensão VS Code em si. Nasce depois desta etapa, lendo o formato abaixo.

## 5. Critérios de aceite

- AC-01 (RF-01) Quando um usuário em máquina limpa rodar `/plugin marketplace add danielmalka/dev-harness` e `/plugin install dev-harness@dev-harness`, então os 17 comandos e os 18 agentes ficam disponíveis sem clone.
- AC-02 (RF-02) Quando `dh build` rodar em Linux, macOS ou WSL2, então `dist/` corresponde às fontes byte a byte e o manifesto lista os cinco binários com checksum.
- AC-03 (RF-03, RF-04) Quando um `/dev-harness:build` real despachar implementador, QA e revisor, então o snapshot lista a sessão e cada subagente com `type`, `model`, `tokenCount` e `status`, e o painel do Claude mostra as mesmas linhas.
- AC-04 (RF-04) Quando o dono enviar um prompt, então `state` vira `active` em até um tick; quando o Claude parar, `idle`; quando a sessão encerrar, `closed`.
- AC-05 (RF-02, RF-03) Quando `dh snapshot` for chamado num tick, então responde em menos de 50 ms e nunca grava prompts, respostas ou segredos.
- AC-06 (RF-05) Quando o usuário tiver a própria `statusLine`, então instalar o plugin não a substitui.
- AC-07 (RF-06) Quando uma sessão encerrar, então o arquivo permanece com `state: closed` e `dh snapshot prune` o remove após 7 dias.
- AC-08 (RF-07) Quando faltar binário para a plataforma, então `dh doctor` (ou o wrapper) explica o que falta sem instalar nada.
- AC-09 (RF-08) Quando alguém seguir o tutorial 00 numa máquina limpa, então instala pelo marketplace sem orientação do autor.

## 6. Restrições e impacto técnico

- Dados: snapshots ficam fora do projeto, em `~/.claude/dev-harness/sessions/`; contêm `cwd` e custo, nunca prompts, respostas, tokens de API ou nomes de arquivo editados. Nada entra em `.harness/`.
- Contratos: formato do snapshot (abaixo) é o contrato com a extensão; versão `schema: 1`. `harness-manifest.json` ganha `binaries` com `os`, `arch`, `path`, `sha256`.
- Segurança: hooks e status lines executam um binário commitado; a CI deve reproduzir o build e o manifesto carrega o checksum; o wrapper recusa binário com checksum divergente.
- Operação: CI cross-compila os cinco alvos; `dh build` local usa o Go instalado do mantenedor.
- Stack e dependências: Go 1.22+ (toolchain do mantenedor, `go 1.27` na máquina do dono); biblioteca padrão apenas; sem CGO.
- Portabilidade: nenhum caminho absoluto ou do autor; o wrapper resolve `bin/` relativo ao diretório do plugin (`${CLAUDE_PLUGIN_ROOT}`).

### Formato do snapshot (`schema: 1`)

```json
{
  "schema": 1,
  "session_id": "…",
  "session_name": "…",
  "cwd": "/path/to/project",
  "agent": "coordinator",
  "model": { "id": "…", "display_name": "…" },
  "state": "active | idle | closed",
  "started_at": "2026-09-20T14:03:00-03:00",
  "updated_at": "2026-09-20T14:20:12-03:00",
  "cost": { "total_cost_usd": 0.0, "total_duration_ms": 0 },
  "context": { "used_percentage": 0, "context_window_size": 0, "input_tokens": 0, "output_tokens": 0 },
  "rate_limits": { },
  "tasks": [
    { "id": "…", "name": "…", "type": "dev-harness:qa-verifier", "status": "running | done | …",
      "model": "…", "effort": "…", "tokenCount": 0, "contextWindowSize": 0, "startTime": "…" }
  ],
  "events": [ { "at": "…", "event": "SubagentStart", "agent_type": "…", "agent_id": "…" } ]
}
```

`events` guarda só as últimas 50 transições. Campos ausentes na origem ficam ausentes, nunca inventados.

## 7. Alternativas descartadas

| Alternativa | Por que não |
|---|---|
| OTEL com collector local | Exige serviço rodando; adiado pelo dono |
| Sessões controladas pela extensão via SDK | A extensão viraria runtime e precisaria de UI de chat; o plano exclui |
| Inferir estado pela idade de arquivos (como a `context-kit-extension`) | Impreciso; os hooks entregam o estado de graça |
| Manter Python nos scripts | Interpretador no consumidor e latência; ADR-001 |
| Forkar a `context-kit-extension` | Herda modelo de biblioteca semeada e convenções que colidem com `.harness/` |

## 8. Riscos e decisões pendentes

- Risco: campos do `subagentStatusLine` mudarem entre versões do Claude Code · mitigação: `schema` no snapshot e campos opcionais; `doctor` reporta a versão do Claude.
- Risco: binários commitados incharem o repositório · mitigação: revisar quando `dist/` passar de 50 MB (ADR-001).
- Risco: hook lento travar a interação · mitigação: `dh snapshot` só faz I/O local, sem rede; medir os 50 ms no `doctor`.
- Decidido em 20/09/2026: o marketplace `dev-harness` fixa a versão do plugin por `ref` de tag (`v<versão do manifesto>`).
- Decidido em 20/09/2026: `dh` substitui validate, build e doctor já nesta etapa; não haverá dois runtimes de manutenção.

## 9. Referências

- `docs/adr/ADR-001-runtime-go.md`
- `docs/plano-produto.html` §4.1, §4.2 e AC-01 a AC-05
- Documentação do Claude Code: plugin marketplaces, hooks, status line e `subagentStatusLine` (consultada em 20/09/2026)
- `~/projetos/00-Revaliar/context-kit-extension` (doadora de código para a extensão, por cópia)
