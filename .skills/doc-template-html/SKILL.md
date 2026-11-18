---
name: doc-template-html
description: Use quando o usuário pede um plano, feature, bug, melhoria, catálogo, relatório ou visão geral do projeto como um arquivo HTML interno autônomo, ou roda `/dh:document`, ou quer uma entrega em HTML aberta no navegador. Não use para respostas só em Markdown no chat.
author: malka
metadata:
  provenance: adapted
  sources: ["prior in-house doc-template-html skill (Void & Lumen HTML standard)"]
---

# Template de documento HTML (Harness)

Documento interno de trabalho: um único arquivo `.html` autônomo, paleta papel-tinta, aberto no navegador. Não é página de site.

## Quando usar

A saída da tarefa é um plano, feature, bug, melhoria, catálogo ou relatório **como arquivo**. O chat resume; o HTML é a entrega.

Não use quando:

- o pedido é uma resposta curta em Markdown no chat, sem arquivo
- o pedido é um card de kanban

## Entradas

| Entrada | Se faltar |
| --- | --- |
| Assunto do documento | Pare. Na sessão principal, pergunte ao dono. Como especialista despachado, devolva a entrada que falta ao Coordenador e pare; nunca prossiga com um valor presumido. |
| Tipo, de `references/types.md` | Infira-o e declare o tipo inferido no hero. |
| Caminho de saída | Padrão `docs/<slug>.html` no repositório do assunto; confirme quando não for óbvio. |
| Os fatos que preenchem as seções | O espaço continua “ainda não fechado”. Nunca lorem, nunca métricas inventadas. |

## Fonte da verdade

1. Uso humano (carimbo, flags, tipos): [README.md](README.md)
2. Casca com CSS/JS embutido: [scripts/stamp.sh](scripts/stamp.sh)
3. Tokens: [assets/shell.css](assets/shell.css) + [assets/shell.js](assets/shell.js)
4. Receitas e markup: [references/types.md](references/types.md)
5. Feature preenchida: [assets/modelo.html](assets/modelo.html)
6. Catálogo preenchido: [assets/catalogo-preview.html](assets/catalogo-preview.html)
7. Visão geral do projeto preenchida: `assets/projeto-modelo.html`

`stamp.sh --sync` embute de novo o CSS/JS nos exemplos da skill.

Leia `types.md` antes de preencher. O HTML gerado não aponta para a pasta da skill.

## Kit visual (não "corrigir")

Papel de tinta sobre carvão frio. O neon é um sinal em números, tags, marca, cards `.card.in` e progresso. Todo o resto continua papel, areia ou carvão.

- vazios `#0C0D10` / `#131418` / `#1A1C22`
- texto `#F1EFE8` · atenuado `#9A978E`
- neon primário `#3DFF9A`
- destaque areia `#E2D2A8` (links, banner, callout do operador)

Tipografia: Outfit + Plus Jakarta Sans + JetBrains Mono. Não lave a página de verde. Não traga de volta o roxo nem o ciano.

Um arquivo. Sem React, Tailwind, emoji, CDN de ícones nem imagem remota.

Diagramas: fonte Mermaid em `<pre class="mermaid">`, renderizada pelo CDN do Mermaid depois do script da casca. Duas exceções de CDN, Mermaid (origem e versão exata fixadas em `stamp.sh`) e Google Fonts; todo o resto embutido. Nunca embuta SVG pré-renderizado.

## Fluxo

1. Escolha o tipo em `references/types.md`. Se o usuário não nomeou um tipo, infira-o e declare-o no hero.
2. Confirme o caminho se não for óbvio. Padrão: `docs/<slug>.html`. Para `pdocs`, o caminho vem de `dh harness-path --json` (`pdocs`), nunca digitado nem inferido.
3. Gere a casca. Rode o carimbo pelo caminho dele dentro da pasta desta skill (`<skill-folder>/scripts/stamp.sh`), resolvido a partir de onde o kit está instalado. `--out` aceita um caminho absoluto ou relativo ao diretório atual.

```bash
bash <skill-folder>/scripts/stamp.sh \
  --type feature \
  --title "Título" \
  --out docs/slug.html
```

`--type bug` já emite `.banner.warn`. `--type aprovacao` é a página de aprovação do plano: ela recebe o plano como `--data <file.json>` (nunca HTML escrito à mão; contrato em `references/types.md`). Outras flags: `--lead`, `--status`, `--eyebrow`, `--date`, `--dono`, `--contato`, `--path`, `--banner-text`, `--banner-warn`, `--force`. Para `--type aprovacao`, nunca copie à mão: o carimbo é o único caminho, porque ele escapa os dados do plano em código. Se não houver shell disponível, ou o script falhar, copie `assets/modelo.html` e troque as seções; para o tipo `projeto`, copie `assets/projeto-modelo.html`. Não escreva `<link>` para `shell.css`.
4. Preencha todo “ainda não fechado” com um fato. Uma lacuna continua “ainda não fechado”. Nunca lorem, nunca métricas inventadas. Markup só de `types.md`.
5. Verifique (abaixo).

## Anatomia (ordem fixa)

O carimbo emite esta ordem. Não reordene.

```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
  meta charset + viewport
  title: "Harness · {título} (não publicar)"
  Google Fonts: JetBrains Mono 400;500 · Outfit 500;600;700;800 · Plus Jakarta Sans 400;500;600;700
  <style> /* assets/shell.css inteiro */ </style>
</head>
<body>
  progress + skip-link
  aside.side: marca Harness. · doc interno · chip de tipo · nav 01…
  wrap#conteudo: banner, hero, h2s do tipo, footer.doc-foot
  <script> /* assets/shell.js inteiro */ </script>
</body>
</html>
```

Nav: um `<a href="#id"><b>01</b>Rótulo</a>` por `h2[id]`. Números de dois dígitos. Marca: SVG `.mark` (círculo neon→areia) + `Harness<span>.</span>`.

Chips do hero, nesta ordem quando presentes: data, status, dono, path, contato.

## Conteúdo

- O corpo e o chrome do documento são pt-BR (ADR-011): o carimbo não tem flag de idioma. Frases curtas. Sem travessão no meio da frase.
- O `h2` leva o número em `.n` e o título sem repetir o número por extenso.
- Entra / fora, processo, tabela, preço, strip, callout: markup em `references/types.md`. Não invente uma classe.

## Verificação

1. Confirme que não há `href`/`src` para `shell.css` nem `shell.js`, e que `<style>` contém `--void:` e `<script>` contém `IntersectionObserver`.
2. Abra o HTML: barra lateral, hero, entra/fora, rodapé.
3. Não afirme que a página foi vista se só o arquivo existe em disco.
4. Um diagrama Mermaid precisa de rede. Em `file://` sem rede a página ainda deve mostrar a fonte do diagrama; nunca afirme que o diagrama foi renderizado a menos que o tenha visto renderizado.

## Formato de saída

Um `.html` autônomo no caminho combinado. No chat: esse caminho relativo à raiz do projeto (acrescente o caminho absoluto só se o dono pedir algo para abrir), o tipo, o próximo passo que o rodapé já declara e a lista dos espaços `ainda não fechado` ainda não resolvidos. Não recapitule o documento.

## Erros comuns

| Erro | Por que prejudica | Faça em vez disso |
| --- | --- | --- |
| Publicar em um site de marketing ou linkar a partir de uma home | O banner o classifica como interno e a classificação deixa de significar algo | Mantenha-o interno e compartilhe só o caminho |
| Trocar a paleta, usar Inter, acrescentar emoji | Os tokens derivam e o doc deixa de parecer saída do kit | Mantenha a paleta declarada e as três fontes |
| Três cards “feature” quase idênticos | Enchimento parece conteúdo e esconde o que ainda está em aberto | Um card por assunto real; uma lacuna continua “ainda não fechado” |
| Deixar CSS ou JS externo, ou apontar para a pasta da skill | O arquivo deixa de abrir sozinho a partir de `file://` | O carimbo embute os dois; o HTML gerado nunca referencia a pasta da skill |
| Sobrescrever `assets/modelo.html`, `assets/catalogo-preview.html` ou `assets/projeto-modelo.html` | Os exemplos canônicos da skill se perdem | Escreva em outro lugar; `--force` só quando o dono pediu |

## Relacionado

Papéis: `docs-guide` (autor), `coordinator` ou `builder` na faixa `infra` (roda o script do carimbo quando é preciso um shell). Comando: `/dh:document`. Skills: `context-handoff`, `delivery-readiness`.

## Caso de prova

Dado um pedido de doc de feature, o fluxo emite um HTML autônomo que abre a partir de `file://` sem CSS/JS externo além das duas exceções de CDN declaradas (Mermaid e Google Fonts), com todo `ainda não fechado` substituído por um fato com fonte ou deixado intacto, e sem afirmar que a página foi inspecionada visualmente a menos que tenha sido.
