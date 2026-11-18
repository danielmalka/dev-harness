# Receitas por tipo

Os esqueletos ficam em `assets/skeletons/`; o placeholder é “ainda não fechado”.

O chrome (barra lateral, hero, banner, rodapé) vem de `scripts/stamp.sh`. O que fica aqui é o esqueleto de `h2` de cada tipo e o markup das primitivas.

Não invente uma seção nem uma classe. Uma lacuna vira “ainda não fechado”. IDs em kebab-case, idênticos na nav e no `h2 id`.

Carimbo: `scripts/stamp.sh --type TIPO --title "…" --out path.html`

O HTML que ele injeta está em `assets/skeletons/<tipo>.html`. Para mais uma linha de produto em um catálogo, duplique o bloco `h2` da linha e acrescente um `<a>` ao `.strip` (2 a 4 células).

Os títulos de seção abaixo são os títulos literais em pt-BR emitidos no documento. Mantenha-os como estão.

## catalogo

Uso: mapa de produtos, preços, posicionamento. Exemplo preenchido nesta skill: `assets/catalogo-preview.html`.

1. Posicionamento
2. Cada linha de produto (promessa, problema, escopo, entra/fora, preço, como atender)
3. Jornada
4. Tabela consolidada
5. Pacotes compostos, se existirem
6. Mapeado vs. em aberto

Primitivas de que este tipo precisa: `.strip`, `.price`, `.tag`, `.card.in` / `.card.out`.

## plano

Uso: planejamento de desenvolvimento, sprint, onda, migração.

1. Contexto e por que agora
2. Resultado esperado
3. Fora de escopo
4. Abordagem (entra / fora)
5. Etapas (`.steps`)
6. Riscos
7. Critérios de aceite
8. Próximas ações

## feature

Uso: feature nova ou fatia de produto. Exemplo preenchido: `assets/modelo.html`.

1. Problema
2. Promessa
3. Escopo (entra / fora)
4. Como se implementa (`.steps`)
5. Métricas de sucesso
6. Dependências
7. Riscos
8. Próximas ações

## melhoria

Uso: evolução de algo que já existe.

1. Estado atual
2. Dor
3. Depois (o que muda)
4. Escopo (entra / fora)
5. Como se implementa (`.steps`)
6. Como medimos
7. Riscos de regressão
8. Próximas ações

## bug

Uso: correção, incidente, defeito. O carimbo já emite `.banner.warn`.

1. Sintoma
2. Repro
3. Impacto
4. Causa (hipótese vs. confirmada)
5. Fix proposto (entra / fora)
6. Verificação
7. Risco residual
8. Próximas ações

## report

Uso: relatório das ações realizadas. Não substitui `protocolo/templates/report.html` dentro de uma execução de `/protocolo`.

1. O que foi pedido
2. O que pousou
3. Evidência
4. O que ficou de fora
5. Decisões tomadas
6. Pendências
7. Próximo passo

## projeto

Uso: visão geral viva de um projeto inteiro (o `pdocs/index.html` de um repositório). Copie `assets/projeto-modelo.html` quando não houver shell, ou rode o carimbo com `--type projeto`. Markup: `.table-wrap`, uma lista e `<pre class="mermaid">`, tudo das Primitivas abaixo. Todo espaço sem um fato continua “ainda não fechado”.

1. Motivação
2. Como rodar local
3. Planejado e desenvolvido
4. Changelog
5. Diagrama macro
6. Domínios e bounded contexts
7. Referências

Ids: `motivacao`, `como-rodar`, `planejado-desenvolvido`, `changelog`, `diagrama-macro`, `dominios`, `referencias`. “Diagrama macro” é obrigatório e contém pelo menos um `<pre class="mermaid">` desenhado a partir das pastas de primeiro nível do repositório. “Planejado e desenvolvido” é um `.table-wrap` com `id`, título e `Status`, mais uma linha “gerado em AAAA-MM-DD”.

## aprovacao

Uso: página de aprovação de um plano, carimbada em `pdocs`: `docs-guide` compõe o arquivo de dados JSON e um `builder` (lane `infra`) roda `stamp.sh`. Etapas, trechos propostos e cartões de decisão que o dono aprova ou rejeita e depois cola de volta no chat como um único bloco Markdown. Um `.html` plano, sem escrita em servidor, sem script externo (o carimbo deixa a tag do Mermaid fora deste tipo). O carimbo não tem flag de idioma: a página é pt-BR (ADR-011).

1. Etapas
2. Trechos propostos
3. Pontos de decisão
4. Copiar decisões

Ids: `etapas`, `trechos`, `decisoes`, `copiar`.

Contrato de inserção (fonte única para `/dh:plan` e o teste da página). Nunca escreva o plano no HTML à mão. Escreva um arquivo JSON e passe-o:

```bash
bash <skill-folder>/scripts/stamp.sh --type aprovacao --title "…" --data plan.json --out <pdocs>/aprovacao-<slug>.html
```

```json
{
  "steps":     [{ "id": "S1", "title": "…", "description": "…" }],
  "snippets":  [{ "step": "S1", "title": "…", "code": "…" }],
  "decisions": [{ "id": "D1", "title": "…", "description": "…" }]
}
```

Todos os campos são strings simples (sem HTML, sem renderização de Markdown); `snippets[].step` é um `steps[].id`; o `title` de um trecho é opcional. O carimbo checa cada item da lista: um objeto com campos string e um `snippets[].step` que bata com um `steps[].id`; caso contrário sai com código 2 nomeando a lista, o índice e o campo, e não escreve arquivo. Os ids (de etapas e de decisões) casam com `^[A-Za-z0-9_.-]{1,64}$`, são únicos dentro da própria lista e nunca são `__proto__`, `constructor` ou `prototype`; o `title` de etapas, decisões e trechos e o `snippets[].step` não contêm caractere de controle nem separador de linha (`description` e `code` podem manter quebras de linha e tabs). Com `--type aprovacao`, o nome de saída deve casar com `aprovacao-<slug>.html` (slug `[A-Za-z0-9_-]`). O bloco copiado achata o espaço em branco em ids e títulos, e a chave de `localStorage` termina com um hash dos dados embutidos, então um plano revisado começa sem escolhas guardadas. `--data` só é aceito com `--type aprovacao`; sem ele, a página é o esqueleto vazio.

Mecanismo de escape (R48): o carimbo embute o JSON em `<script type="application/json" id="dh-data">` com `<`, `>`, `&`, U+2028 e U+2029 escritos como escapes `\uXXXX`, de modo que nenhum texto consegue fechar a tag nem iniciar markup; o script da página o lê com `JSON.parse` e insere cada string com `textContent`. A página não tem `innerHTML`.

Ganchos estáveis do DOM:

| gancho | elemento |
|---|---|
| `[data-step="<id>"]` | um item de etapa, em ordem |
| `[data-snippet][data-step-ref="<step id>"]` | um trecho; contém `pre > code` e um rótulo que nomeia a etapa |
| `[data-decision-card="<id>"]` | um cartão de decisão; `data-state` é `pending`, `approved` ou `rejected` |
| `[data-action="approve"]`, `[data-action="reject"]` | os dois botões de um cartão (`aria-pressed` marca o escolhido; uma marca de visto ou de cruz e o rótulo de estado o mostram sem cor) |
| `[data-comment]` | o `textarea` de comentário do cartão |
| `[data-reject-error]` | mensagem exibida (`hidden` removido, `role="alert"`) quando o cartão é rejeitado com comentário em branco; o cartão continua `pending` |
| `[data-state-label]` | o texto visível de estado do cartão |
| `[data-copy]` | o botão “Copiar decisões” |
| `[data-copy-status]` | texto de status ao vivo depois de uma tentativa de cópia |
| `[data-fallback]` | `textarea` somente leitura com o bloco, dentro de `[data-fallback-box]`; exibido quando a área de transferência não existe, lança erro ou é recusada |

Bloco copiado, uma entrada por cartão na ordem da página:

```text
## Decisões

### D1 Título
- decisão: approved|rejected|pending
- comentário: texto (linhas seguintes indentadas com quatro espaços; comentário vazio deixa a linha como "- comentário:")
```

O título é `## Decisões` e os rótulos são `decisão` e `comentário`; os valores de veredito continuam `approved`, `rejected`, `pending`. As escolhas persistem em `localStorage` (chave `dh-aprovacao:` mais o título do documento); todo acesso fica dentro de `try/catch`, e a página funciona com o armazenamento bloqueado.

## Primitives

Copie o markup. Não invente uma classe vizinha.

Banner de classificação:

```html
<div class="banner">Doc interno. Não publicar.</div>
```

Banner de risco, para dados que não podem vazar:

```html
<div class="banner warn">Incidente interno. Não vazar.</div>
```

Chips do hero, nesta ordem quando presentes: data, status, dono, path, contato.

```html
<div class="chips">
  <span class="chip">data <b>2026-09-12</b></span>
  <span class="chip">status <b>rascunho</b></span>
  <span class="chip">dono <b>operador</b></span>
  <span class="chip">path <b>docs/exemplo.html</b></span>
  <span class="chip">contato <b>ops@example.test</b></span>
</div>
```

Strip (2 a 4 âncoras reais; a primeira célula é a âncora):

```html
<div class="strip">
  <a href="#agents"><span>Âncora</span><strong>Agents</strong>18 papéis</a>
  <a href="#commands"><span>Operação</span><strong>Commands</strong>Fluxos de trabalho</a>
  <a href="#skills"><span>Método</span><strong>Skills</strong>Procedimentos</a>
</div>
```

Entra / fora, sempre em par:

```html
<div class="grid-2">
  <div class="card in">
    <h3>Entra</h3>
    <ul>
      <li>ainda não fechado</li>
    </ul>
  </div>
  <div class="card out">
    <h3>Fica de fora</h3>
    <ul>
      <li>ainda não fechado</li>
    </ul>
  </div>
</div>
```

Tag, preço, processo, tabela, callouts:

```html
<span class="tag">Kit · operação</span>
<span class="price">R$ 297–897</span>
<ol class="steps">
  <li><strong>Passo</strong> o que fazer</li>
</ol>
<div class="table-wrap">
  <table>
    <thead><tr><th>Coluna</th><th>Coluna</th></tr></thead>
    <tbody><tr><td>ainda não fechado</td><td>ainda não fechado</td></tr></tbody>
  </table>
</div>
<div class="callout">Nota operacional.</div>
<div class="callout gold">Regra do operador.</div>
<div class="callout stop">Não fazer.</div>
```

| bloco | quando |
|---|---|
| `.banner` | classificação (interno, rascunho, não publicar) |
| `.banner.warn` | dado sensível ou risco ativo |
| `.chips` / `.chip` | data, dono, status, path, contato |
| `.strip` | 2–4 âncoras de produto ou módulo; a primeira célula é maior |
| `.card` | um assunto com um título |
| `.card.in` / `.card.out` | entra vs. fora, sempre em par dentro de `.grid-2` |
| `.tag` | linha de produto ou tipo de seção |
| `.price` | valor em dinheiro; nunca um status |
| `.steps` | processo, implementação, reprodução |
| `.table-wrap` + `table` | comparação, preço, jornada |
| `.callout` | nota operacional |
| `.callout.gold` | regra do operador |
| `.callout.stop` | não fazer / bloqueio |
| `code` | caminho, comando, id |

Diagrama: fonte Mermaid em `<pre class="mermaid">`, renderizada pelo CDN do Mermaid. Duas exceções de CDN, Mermaid (origem e versão exata fixadas em `stamp.sh`) e Google Fonts; todo o resto embutido. Nunca embuta SVG pré-renderizado; o dono edita o diagrama depois. Um diagrama Mermaid precisa de rede: em `file://` sem rede a página ainda deve mostrar a fonte do diagrama, e você nunca afirma que o diagrama foi renderizado a menos que o tenha visto renderizado.

```html
<pre class="mermaid">
flowchart LR
    A[Entrada] --> B[Saída]
</pre>
<!-- emitido por stamp.sh depois do <script> da casca, no fim do body; origem e versão exata ficam fixadas lá, nunca digitadas aqui -->
<script src="<origem, versão exata e hash SRI fixados em stamp.sh>" crossorigin="anonymous"></script>
<script>
  mermaid.initialize({ startOnLoad: true, theme: "base", themeVariables: { /* variáveis de tema: veja stamp.sh */ } });
</script>
```
