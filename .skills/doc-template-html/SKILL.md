---
name: doc-template-html
description: Use when the user asks for a plan, feature, bug, improvement, catalog or report as a standalone internal HTML file, or runs `/dev-harness:document`, or wants an HTML deliverable opened in the browser. Not for Markdown-only chat replies.
author: malka
metadata:
  provenance: adapted
  sources: ["prior in-house doc-template-html skill (Void & Lumen HTML standard)"]
---

# Doc template HTML (Harness)

Internal working document: one standalone `.html` file, ink-paper palette, opened in the browser. Not a site page.

## When to use

The output of the task is a plan, feature, bug, improvement, catalog or report **as a file**. Chat summarizes; the HTML is the deliverable.

Do not use when:

- the request is a short Markdown reply in chat, with no file
- the request is a kanban card

## Inputs

| Input | If missing |
| --- | --- |
| Subject of the document | Stop. In the main session, ask the owner. As a dispatched specialist, return the missing input to the Coordinator and stop; never proceed on an assumed value. |
| Type, from `references/types.md` | Infer it and declare the inferred type in the hero. |
| Output path | Default `docs/<slug>.html` in the subject repo; confirm when it is not obvious. |
| The facts that fill the sections | The slot stays “ainda não fechado”. Never lorem, never invented metrics. |

## Source of truth

1. Human usage (stamp, flags, types): [README.md](README.md)
2. Shell with inline CSS/JS: [scripts/stamp.sh](scripts/stamp.sh)
3. Tokens: [assets/shell.css](assets/shell.css) + [assets/shell.js](assets/shell.js)
4. Recipes and markup: [references/types.md](references/types.md)
5. Filled feature: [assets/modelo.html](assets/modelo.html)
6. Filled catalog: [assets/catalogo-preview.html](assets/catalogo-preview.html)

`stamp.sh --sync` re-inlines CSS/JS into the skill examples.

Read `types.md` before filling. Generated HTML does not point at the skill folder.

## Visual kit (do not "fix")

Ink paper on cool charcoal. Neon is a signal on numbers, tags, brand mark, in-cards and progress. Everything else stays paper, sand or charcoal.

- voids `#0C0D10` / `#131418` / `#1A1C22`
- text `#F1EFE8` · muted `#9A978E`
- neon primary `#3DFF9A`
- sand accent `#E2D2A8` (links, banner, operator callout)

Typography: Outfit + Plus Jakarta Sans + JetBrains Mono. Do not wash the page in green. Do not bring purple or cyan back.

One file. Google Fonts in `<head>`. No React, Tailwind, emoji, icon CDN, remote image.

Diagrams: Mermaid source in `<pre class="mermaid">`, rendered by the Mermaid CDN after the shell script. Mermaid is the only CDN exception. Never embed pre-rendered SVG.

## Flow

1. Pick the type in `references/types.md`. If the user did not name a type, infer and declare it in the hero.
2. Confirm the path if it is not obvious. Default: `docs/<slug>.html`.
3. Generate the shell. Run the stamp by its path inside this skill folder (`<skill-folder>/scripts/stamp.sh`), resolved from wherever the kit is installed. `--out` takes an absolute path or one relative to the current directory.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type feature \
  --title "Título" \
  --lang pt-br \
  --out docs/slug.html
```

`--type bug` already emits `.banner.warn`. Other flags: `--lang`, `--lead`, `--status`, `--eyebrow`, `--date`, `--dono`, `--contato`, `--path`, `--banner-text`, `--banner-warn`, `--force`. If no shell is available, or the script fails, copy `assets/modelo.html` and swap sections. Do not write `<link>` to `shell.css`.
4. Fill every “ainda não fechado” with a fact. A hole stays “ainda não fechado”. Never lorem, never invented metrics. Markup only from `types.md`.
5. Verify (below).

## Anatomy (fixed order)

The stamp emits this order. Do not reorder.

```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
  meta charset + viewport
  title: "Harness · {título} (não publicar)"
  Google Fonts: JetBrains Mono 400;500 · Outfit 500;600;700;800 · Plus Jakarta Sans 400;500;600;700
  <style> /* assets/shell.css entire */ </style>
</head>
<body>
  progress + skip-link
  aside.side: mark Harness. · doc interno · type chip · nav 01…
  wrap#conteudo: banner, hero, type h2s, footer.doc-foot
  <script> /* assets/shell.js entire */ </script>
</body>
</html>
```

Nav: one `<a href="#id"><b>01</b>Label</a>` per `h2[id]`. Two-digit numbers. Mark: SVG `.mark` (circle neon→sand) + `Harness<span>.</span>`.

Hero chips, in this order when present: data, status, dono, path, contato.

## Content

- Document body and chrome in the project language: pass `--lang en` or `--lang pt-br` to the stamp (default `pt-br`), taking the value from `language` in `.harness/project.yaml` when the document belongs to a consumer project, `pt-br` for this kit's own docs. Short sentences. No em dash mid-sentence.
- `h2` carries the number in `.n` and the title without repeating the number in words.
- In / out, process, table, price, strip, callout: markup in `references/types.md`. Do not invent a class.

## Verify

1. Confirm there is no `href`/`src` to `shell.css` or `shell.js`, and that `<style>` contains `--void:` and `<script>` contains `IntersectionObserver`.
2. Open the HTML: sidebar, hero, in/out, footer.
3. Do not claim the page was seen if only the file exists on disk.
4. A Mermaid diagram needs network. On `file://` without network the page must still show the diagram source; never claim the diagram rendered unless you saw it rendered.

## Output format

One standalone `.html` at the agreed path. In chat: that path relative to the project root (add the absolute path only if the owner asks for something to open), the type, the next step the footer already states, and the list of `ainda não fechado` slots still unresolved. Do not recap the document.

## Common mistakes

| Mistake | Why it hurts | Do instead |
| --- | --- | --- |
| Publishing on a marketing site or linking it from a home page | The banner classifies it as internal and the classification stops meaning anything | Keep it internal and share only the path |
| Swapping the palette, using Inter, adding emoji | The tokens drift and the doc stops reading as kit output | Keep the declared palette and the three fonts |
| Three near-identical “feature” cards | Padding reads as content and hides what is still open | One card per real subject; a hole stays “ainda não fechado” |
| Leaving CSS or JS external, or pointing at the skill folder | The file stops opening alone from `file://` | The stamp inlines both; generated HTML never references the skill folder |
| Overwriting `assets/modelo.html` or `assets/catalogo-preview.html` | The skill's canonical examples are lost | Write elsewhere; `--force` only when the owner asked |

## Related

Roles: docs-guide (author), coordinator or devops-engineer (runs the stamp script when a shell is needed). Command: `/dev-harness:document`. Skills: context-handoff, delivery-readiness.

## Proof case

Given a request for a feature doc, the flow emits one standalone HTML that opens from `file://` with no external CSS/JS beyond the declared Mermaid exception, every `ainda não fechado` either replaced by a sourced fact or left intact, and no claim that the page was visually inspected unless it was.
