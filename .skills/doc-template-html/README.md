# doc-template-html

Documento interno do Harness como um `.html` autônomo, aberto no navegador. Não é página de site.

O carimbo é portátil: ele embute o CSS/JS da própria pasta da skill, onde quer que o kit esteja instalado.

## Peça ao agente

No chat: `/dh:document` ou “gera o plano/feature/bug/melhoria/catálogo/report em HTML interno”.

O agente lê `SKILL.md`, escolhe o tipo em `references/types.md`, roda o carimbo e preenche.

## Gere a casca (humano)

Rode o carimbo pelo caminho dele dentro da pasta desta skill (`<skill-folder>/scripts/stamp.sh`), resolvido a partir de onde o kit está instalado. `--out` aceita um caminho absoluto ou relativo ao diretório atual.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type feature \
  --title "Título" \
  --out docs/slug.html
```

Tipos: `catalogo` `plano` `feature` `melhoria` `bug` `report` `projeto` `aprovacao`.

`--type aprovacao` é a página de aprovação do plano: passe o plano como `--data plan.json` (contrato em `references/types.md`); a página tem script próprio e nenhum externo.

`--type bug` já vem com `.banner.warn`.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type bug \
  --title "Webhook cai em 502" \
  --lead "Cliente sem resposta desde 14h." \
  --status rascunho \
  --dono operador \
  --out docs/bug-webhook-502.html
```

Outras flags: `--contato`, `--path`, `--date`, `--eyebrow`, `--banner-text`, `--banner-warn`, `--force`.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type report \
  --title "Relatório de release" \
  --out docs/relatorio-release.html
```

```bash
bash <skill-folder>/scripts/stamp.sh --help
```

Abra o arquivo em `file://`. Preencha todo “ainda não fechado”. Uma lacuna continua “ainda não fechado”.

## O que é cada coisa

| arquivo | papel |
|---|---|
| `SKILL.md` | contrato do agente |
| `scripts/stamp.sh` | casca com CSS/JS embutido |
| `references/types.md` | receitas + markup das primitivas |
| `assets/skeletons/<tipo>.html` | seções que o carimbo injeta |
| `assets/modelo.html` | feature preenchida |
| `assets/catalogo-preview.html` | catálogo preenchido (strip, price) |
| `assets/projeto-modelo.html` | visão geral do projeto preenchida, sem precisar de casca |
| `assets/shell.css` + `shell.js` | tokens; não linke, o carimbo embute |

Paleta: carvão `#0C0D10`, papel `#F1EFE8`, neon `#3DFF9A` nos destaques, areia `#E2D2A8` em links e banners. O carimbo se recusa a sobrescrever os exemplos da skill sem `--force`.

Depois de editar `shell.css` / `shell.js`:

```bash
bash <skill-folder>/scripts/stamp.sh --sync
```

## Não use

- Só Markdown no chat, sem arquivo
- publicar em um site ou linkar a partir de uma home
