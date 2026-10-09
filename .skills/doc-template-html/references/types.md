# Recipes by type

Skeletons live in `assets/skeletons/<lang>/`; the placeholder is “ainda não fechado” in pt-br and “not yet settled” in en.

The chrome (sidebar, hero, banner, footer) comes from `scripts/stamp.sh`. What lives here is the `h2` skeleton of each type and the markup of the primitives.

Do not invent a section or a class. A hole becomes “ainda não fechado”. IDs in kebab-case, identical in the nav and in the `h2 id`.

Stamp: `scripts/stamp.sh --type TIPO --title "…" --out path.html`

The HTML it injects is in `assets/skeletons/<tipo>.html`. For one more product line in a catalog, duplicate the line's `h2` block and add an `<a>` to the `.strip` (2 to 4 cells).

Section titles below are the literal pt-BR headings emitted into the document. Keep them as written.

## catalogo

Use: product map, pricing, positioning. Filled example in this skill: `assets/catalogo-preview.html`.

1. Posicionamento (Positioning)
2. Cada linha de produto (promessa, problema, escopo, entra/fora, preço, como atender)
3. Jornada (Journey)
4. Tabela consolidada (Consolidated table)
5. Pacotes compostos, se existirem (Compound packages, if any)
6. Mapeado vs. em aberto (Mapped vs. open)

Primitives this type needs: `.strip`, `.price`, `.tag`, `.card.in` / `.card.out`.

## plano

Use: development planning, sprint, wave, migration.

1. Contexto e por que agora (Context and why now)
2. Resultado esperado (Expected outcome)
3. Fora de escopo (Out of scope)
4. Abordagem (entra / fora) (Approach, in / out)
5. Etapas (`.steps`) (Stages)
6. Riscos (Risks)
7. Critérios de aceite (Acceptance criteria)
8. Próximas ações (Next actions)

## feature

Use: new feature or product slice. Filled example: `assets/modelo.html`.

1. Problema (Problem)
2. Promessa (Promise)
3. Escopo (entra / fora) (Scope, in / out)
4. Como se implementa (`.steps`) (How it is implemented)
5. Métricas de sucesso (Success metrics)
6. Dependências (Dependencies)
7. Riscos (Risks)
8. Próximas ações (Next actions)

## melhoria

Use: evolution of something that already exists.

1. Estado atual (Current state)
2. Dor (Pain)
3. Depois (o que muda) (After, what changes)
4. Escopo (entra / fora) (Scope, in / out)
5. Como se implementa (`.steps`) (How it is implemented)
6. Como medimos (How we measure)
7. Riscos de regressão (Regression risks)
8. Próximas ações (Next actions)

## bug

Use: fix, incident, defect. The stamp already emits `.banner.warn`.

1. Sintoma (Symptom)
2. Repro (Reproduction)
3. Impacto (Impact)
4. Causa (hipótese vs. confirmada) (Cause, hypothesis vs. confirmed)
5. Fix proposto (entra / fora) (Proposed fix, in / out)
6. Verificação (Verification)
7. Risco residual (Residual risk)
8. Próximas ações (Next actions)

## report

Use: report of actions taken. It does not replace `protocolo/templates/report.html` inside a `/protocolo` run.

1. O que foi pedido (What was requested)
2. O que pousou (What landed)
3. Evidência (Evidence)
4. O que ficou de fora (What was left out)
5. Decisões tomadas (Decisions made)
6. Pendências (Open items)
7. Próximo passo (Next step)

## projeto

Use: living overview of a whole project (the `pdocs/index.html` of a repository). Copy `assets/projeto-modelo.<lang>.html` when there is no shell, or run the stamp with `--type projeto`. Markup: `.table-wrap`, a list and `<pre class="mermaid">`, all from Primitives below. Every slot without a fact stays “ainda não fechado”.

1. Motivação (Motivation)
2. Como rodar local (How to run locally)
3. Planejado e desenvolvido (Planned and delivered)
4. Changelog (Changelog)
5. Diagrama macro (Macro diagram)
6. Domínios e bounded contexts (Domains and bounded contexts)
7. Referências (References)

Ids: `motivacao`, `como-rodar`, `planejado-desenvolvido`, `changelog`, `diagrama-macro`, `dominios`, `referencias`, the same in both languages. “Diagrama macro” is mandatory and holds at least one `<pre class="mermaid">` drawn from the repository's first-level folders. “Planejado e desenvolvido” is a `.table-wrap` with `id`, title and `Status`, plus a “gerado em AAAA-MM-DD” line.

## Primitives

Copy the markup. Do not invent a neighboring class.

Classification banner:

```html
<div class="banner">Doc interno. Não publicar.</div>
```

Risk banner, for data that must not leak:

```html
<div class="banner warn">Incidente interno. Não vazar.</div>
```

Hero chips, in this order when present: data, status, dono, path, contato.

```html
<div class="chips">
  <span class="chip">data <b>2026-09-12</b></span>
  <span class="chip">status <b>rascunho</b></span>
  <span class="chip">dono <b>operador</b></span>
  <span class="chip">path <b>docs/exemplo.html</b></span>
  <span class="chip">contato <b>ops@example.test</b></span>
</div>
```

Strip (2 to 4 real anchors; the first cell is the anchor):

```html
<div class="strip">
  <a href="#agents"><span>Âncora</span><strong>Agents</strong>18 papéis</a>
  <a href="#commands"><span>Operação</span><strong>Commands</strong>Fluxos de trabalho</a>
  <a href="#skills"><span>Método</span><strong>Skills</strong>Procedimentos</a>
</div>
```

In / out, always as a pair:

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

Tag, price, process, table, callouts:

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

| block | when |
|---|---|
| `.banner` | classification (internal, draft, do not publish) |
| `.banner.warn` | sensitive data or active risk |
| `.chips` / `.chip` | date, owner, status, path, contact |
| `.strip` | 2–4 product or module anchors; the first cell is larger |
| `.card` | one subject with a title |
| `.card.in` / `.card.out` | in vs. out, always paired inside `.grid-2` |
| `.tag` | product line or section type |
| `.price` | money value; never a status |
| `.steps` | process, implementation, repro |
| `.table-wrap` + `table` | comparison, price, journey |
| `.callout` | operational note |
| `.callout.gold` | operator rule |
| `.callout.stop` | do not do / blocker |
| `code` | path, command, id |

Diagram: Mermaid source in `<pre class="mermaid">`, rendered by the Mermaid CDN. Two CDN exceptions, Mermaid (origin and exact version pinned in `stamp.sh`) and Google Fonts; everything else inline. Never embed pre-rendered SVG; the owner edits the diagram later. A Mermaid diagram needs network: on `file://` without network the page must still show the diagram source, and you never claim the diagram rendered unless you saw it rendered.

```html
<pre class="mermaid">
flowchart LR
    A[Entrada] --> B[Saída]
</pre>
<!-- emitted by stamp.sh after the shell <script>, at the end of the body; origin and exact version are pinned there, never typed here -->
<script src="<origin, exact version and SRI hash pinned in stamp.sh>" crossorigin="anonymous"></script>
<script>
  mermaid.initialize({ startOnLoad: true, theme: "base", themeVariables: { /* theme variables: see stamp.sh */ } });
</script>
```
