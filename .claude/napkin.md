# Napkin Runbook

## Curation Rules
- Re-prioritize on every read.
- Keep recurring, high-value notes only.
- Max 10 items per category.
- Each item includes date + "Do instead".

## Domain Behavior Guardrails
1. **[2026-09-20] Internal HTML docs use ink paper + neon signal, never a green wash**
   Do instead: generate internal HTML via `.skills/doc-template-html`. Charcoal `#0C0D10`, paper `#F1EFE8`, neon `#3DFF9A` only on highlights, sand `#E2D2A8` on links/banners. Brand mark is `Harness.`. Do not reintroduce purple, cyan, or a second green wash.

## User Directives
1. **[2026-09-20] Canonical kit sources are dotted folders**
   Do instead: write agents in `.agents/`, skills in `.skills/`, commands in `.commands/`, plus `templates/` and `profiles/`. Regenerate `dist/` with `scripts/build-claude-code.sh`. Never hand-edit the generated package.

2. **[2026-09-20] License is MIT**
   Do instead: keep `LICENSE` and generated `plugin.json` on MIT. Do not introduce a second license without asking.

3. **[2026-09-20] Maintenance scripts are bash + Python 3**
   Do instead: use `scripts/validate.py`, `scripts/doctor.sh` and `scripts/build-claude-code.sh`. Do not add a Node packager unless a consumer needs it.
