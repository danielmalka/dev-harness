# doc-template-html

Internal Harness document as a standalone `.html`, opened in the browser. Not a site page.

The stamp is portable: it inlines CSS/JS from its own skill folder, wherever the kit is installed.

## Ask the agent

In chat: `/dev-harness:document` or “gera o plano/feature/bug/melhoria/catálogo/report em HTML interno”.

The agent reads `SKILL.md`, picks the type in `references/types.md`, runs the stamp and fills.

## Generate the shell (human)

Run the stamp by its path inside this skill folder (`<skill-folder>/scripts/stamp.sh`), resolved from wherever the kit is installed. `--out` takes an absolute path or one relative to the current directory.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type feature \
  --title "Título" \
  --lang pt-br \
  --out docs/slug.html
```

Types: `catalogo` `plano` `feature` `melhoria` `bug` `report`.

`--type bug` already comes with `.banner.warn`.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type bug \
  --title "Webhook cai em 502" \
  --lead "Cliente sem resposta desde 14h." \
  --status rascunho \
  --dono operador \
  --out docs/bug-webhook-502.html
```

Other flags: `--lang pt-br|en`, `--contato`, `--path`, `--date`, `--eyebrow`, `--banner-text`, `--banner-warn`, `--force`.

Language flag: `--lang pt-br|en` (default: `pt-br`).

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type report \
  --title "Release report" \
  --lang en \
  --out docs/release-report.html
```

```bash
bash <skill-folder>/scripts/stamp.sh --help
```

Open the file in `file://`. Fill every “ainda não fechado”. A hole stays “ainda não fechado”.

## What is what

| file | role |
|---|---|
| `SKILL.md` | agent contract |
| `scripts/stamp.sh` | shell with inline CSS/JS |
| `references/types.md` | recipes + primitive markup |
| `assets/skeletons/<lang>/<tipo>.html` | sections the stamp injects |
| `assets/modelo.html` | filled feature |
| `assets/catalogo-preview.html` | filled catalog (strip, price) |
| `assets/shell.css` + `shell.js` | tokens; do not link, the stamp inlines |

Palette: charcoal `#0C0D10`, paper `#F1EFE8`, neon `#3DFF9A` on highlights, sand `#E2D2A8` on links and banners. Stamp refuses to overwrite the skill examples without `--force`.

After editing `shell.css` / `shell.js`:

```bash
bash <skill-folder>/scripts/stamp.sh --sync
```

## Do not use

- Markdown only in chat, no file
- publishing on a site or linking from a home page
