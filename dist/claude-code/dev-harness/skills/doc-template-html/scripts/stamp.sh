#!/usr/bin/env bash
# Emit a standalone Harness doc: shell.css + shell.js inlined, type skeleton injected.
# Usage:
#   stamp.sh --type feature --title "Limite por plano" --lang pt-br --out docs/limite.html
#   stamp.sh --sync
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
exec python3 - "$ROOT" "$@" <<'PY'
from __future__ import annotations

import argparse
import html
import re
import sys
from datetime import date
from pathlib import Path

ROOT = Path(sys.argv[1])
TYPES = ("catalogo", "plano", "feature", "melhoria", "bug", "report")
LANG = {
    "pt-br": {
        "html_lang": "pt-BR",
        "lead_default": "ainda não fechado",
        "status_default": "rascunho",
        "skip_link": "Ir ao conteúdo",
        "side_sub": "Doc interno",
        "banner": "Doc interno. Não publicar.",
        "banner_warn": "Incidente interno. Não vazar.",
        "title_suffix": "(não publicar)",
        "chip_date": "data",
        "chip_status": "status",
        "chip_owner": "dono",
        "chip_path": "path",
        "chip_contact": "contato",
        "foot_file": "Arquivo",
        "foot_type": "tipo",
        "foot_next": "Próximo passo",
        "next_step": "trocar cada “ainda não fechado” por fato. Buraco continua “ainda não fechado”.",
    },
    "en": {
        "html_lang": "en",
        "lead_default": "not yet settled",
        "status_default": "draft",
        "skip_link": "Skip to content",
        "side_sub": "Internal doc",
        "banner": "Internal doc. Do not publish.",
        "banner_warn": "Internal incident. Do not leak.",
        "title_suffix": "(do not publish)",
        "chip_date": "date",
        "chip_status": "status",
        "chip_owner": "owner",
        "chip_path": "path",
        "chip_contact": "contact",
        "foot_file": "File",
        "foot_type": "type",
        "foot_next": "Next step",
        "next_step": "replace every “not yet settled” with a fact. A gap stays “not yet settled”.",
    },
}
LABELS = {
    "pt-br": {
        "catalogo": "Catálogo",
        "plano": "Plano",
        "feature": "Feature",
        "melhoria": "Melhoria",
        "bug": "Bug",
        "report": "Report",
    },
    "en": {
        "catalogo": "Catalog",
        "plano": "Plan",
        "feature": "Feature",
        "melhoria": "Improvement",
        "bug": "Bug",
        "report": "Report",
    },
}
FONTS = (
    "https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500"
    "&family=Outfit:wght@500;600;700;800"
    "&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap"
)
MARK = """        <svg class="mark" viewBox="0 0 32 32" aria-hidden="true">
          <defs>
            <linearGradient id="mark-grad" x1="4" y1="2" x2="28" y2="30">
              <stop offset="0" stop-color="#3DFF9A"/>
              <stop offset="1" stop-color="#E2D2A8"/>
            </linearGradient>
          </defs>
          <circle cx="16" cy="16" r="15" fill="url(#mark-grad)"/>
        </svg>"""
H2_RE = re.compile(
    r'<h2 id="([^"]+)"><span class="n">(\d+)</span>([^<]+)</h2>'
)
STYLE_RE = re.compile(r"<style>.*?</style>", re.S)
SCRIPT_RE = re.compile(r"<script>.*?</script>", re.S)

PROTECTED = {
    (ROOT / "assets" / "modelo.html").resolve(),
    (ROOT / "assets" / "catalogo-preview.html").resolve(),
    (ROOT / "assets" / "shell.css").resolve(),
    (ROOT / "assets" / "shell.js").resolve(),
}


def die(msg: str, code: int = 2) -> None:
    print(msg, file=sys.stderr)
    sys.exit(code)


def read_shell() -> tuple[str, str]:
    css = (ROOT / "assets" / "shell.css").read_text(encoding="utf-8")
    js = (ROOT / "assets" / "shell.js").read_text(encoding="utf-8")
    if not css.strip() or not js.strip():
        die("assets/shell.css or assets/shell.js is empty")
    return css, js


def parse_args(argv: list[str]) -> argparse.Namespace:
    p = argparse.ArgumentParser(
        prog="stamp.sh",
        description="Emit a standalone Harness HTML document with CSS/JS inlined.",
    )
    p.add_argument("--sync", action="store_true", help="Re-inline CSS/JS into the skill example HTML files")
    p.add_argument("--type", choices=TYPES)
    p.add_argument("--title")
    p.add_argument("--out")
    p.add_argument("--lang", choices=tuple(LANG), default="pt-br")
    p.add_argument("--lead")
    p.add_argument("--status")
    p.add_argument("--eyebrow")
    p.add_argument("--date", default=date.today().isoformat())
    p.add_argument("--path")
    p.add_argument("--dono")
    p.add_argument("--contato")
    p.add_argument("--banner-text")
    p.add_argument("--banner-warn", action="store_true", help="Use .banner.warn (default for type bug)")
    p.add_argument("--force", action="store_true")
    ns = p.parse_args(argv)
    strings = LANG[ns.lang]
    if ns.lead is None:
        ns.lead = strings["lead_default"]
    if ns.status is None:
        ns.status = strings["status_default"]
    return ns


def sync_examples(css: str, js: str) -> None:
    for name in ("modelo.html", "catalogo-preview.html"):
        path = ROOT / "assets" / name
        if not path.is_file():
            continue
        text = path.read_text(encoding="utf-8")
        if not STYLE_RE.search(text) or not SCRIPT_RE.search(text):
            die(f"{path}: no <style> or <script> to synchronize")
        text = STYLE_RE.sub("<style>\n" + css.rstrip() + "\n</style>", text, count=1)
        text = SCRIPT_RE.sub("<script>\n" + js.rstrip() + "\n</script>", text, count=1)
        path.write_text(text, encoding="utf-8")
        print(path)


def nav_from(skel: str) -> str:
    items = H2_RE.findall(skel)
    if not items:
        die("skeleton has no h2 in the form <h2 id=\"...\"><span class=\"n\">NN</span>Title</h2>")
    lines = []
    for id_, num, label in items:
        lines.append(f'        <a href="#{id_}"><b>{num}</b>{html.escape(label.strip())}</a>')
    return "\n".join(lines)


def chips(ns: argparse.Namespace, shown_path: str, strings: dict[str, str]) -> str:
    parts = [
        f'<span class="chip">{strings["chip_date"]} <b>{html.escape(ns.date)}</b></span>',
        f'<span class="chip">{strings["chip_status"]} <b>{html.escape(ns.status)}</b></span>',
    ]
    if ns.dono:
        parts.append(f'<span class="chip">{strings["chip_owner"]} <b>{html.escape(ns.dono)}</b></span>')
    parts.append(f'<span class="chip">{strings["chip_path"]} <b>{html.escape(shown_path)}</b></span>')
    if ns.contato:
        parts.append(f'<span class="chip">{strings["chip_contact"]} <b>{html.escape(ns.contato)}</b></span>')
    return "\n          ".join(parts)


def render(ns: argparse.Namespace, css: str, js: str) -> str:
    strings = LANG[ns.lang]
    skel_path = ROOT / "assets" / "skeletons" / ns.lang / f"{ns.type}.html"
    if not skel_path.is_file():
        die(f"missing skeleton: {skel_path}")
    skel = skel_path.read_text(encoding="utf-8").rstrip() + "\n"
    title = ns.title.strip()
    label = LABELS[ns.lang][ns.type]
    eyebrow = ns.eyebrow or f"{label} · {ns.status}"
    shown_path = ns.path or ns.out
    warn = ns.banner_warn or ns.type == "bug"
    banner_class = ' class="banner warn"' if warn else ' class="banner"'
    banner_text = ns.banner_text or (
        strings["banner_warn"] if warn else strings["banner"]
    )
    esc_title = html.escape(title)
    doc_title = f"Harness · {esc_title} {strings['title_suffix']}"
    next_step = strings["next_step"]
    return f"""<!DOCTYPE html>
<html lang="{strings['html_lang']}">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{doc_title}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="{FONTS}" rel="stylesheet">
  <style>
{css.rstrip()}
  </style>
</head>
<body>
  <div class="progress" id="progress" aria-hidden="true"></div>
  <a class="skip-link" href="#conteudo">{strings['skip_link']}</a>
  <div class="page">
    <aside class="side">
      <div class="side__brand">
{MARK}
        Harness<span>.</span>
      </div>
       <div class="side__sub">{strings['side_sub']}</div>
      <div class="side__type">{html.escape(ns.type)}</div>
      <nav>
{nav_from(skel)}
      </nav>
    </aside>
    <div class="wrap" id="conteudo">
      <div{banner_class}>{html.escape(banner_text)}</div>
      <header class="hero">
        <div class="eyebrow">{html.escape(eyebrow)}</div>
        <h1>{esc_title}</h1>
        <p class="lead">{html.escape(ns.lead)}</p>
        <div class="chips">
           {chips(ns, shown_path, strings)}
        </div>
      </header>

{skel}      <footer class="doc-foot">
        <div>{strings['foot_file']}: <code>{html.escape(shown_path)}</code> · {strings['foot_type']} <code>{html.escape(ns.type)}</code></div>
        <div>{strings['foot_next']}: {next_step}</div>
      </footer>
    </div>
  </div>
  <script>
{js.rstrip()}
  </script>
</body>
</html>
"""


def main() -> None:
    ns = parse_args(sys.argv[2:])
    css, js = read_shell()
    if ns.sync:
        sync_examples(css, js)
        return
    if not (ns.type and ns.title and ns.out):
        die("usage: stamp.sh --type TYPE --title TITLE --lang pt-br|en --out PATH.html\n       stamp.sh --sync")
    out = Path(ns.out).expanduser()
    if not out.is_absolute():
        out = Path.cwd() / out
    resolved = out.resolve()
    if resolved in PROTECTED and not ns.force:
        die(f"refused: {resolved} is a canonical skill file. Pass --force if the owner asked.")
    out.parent.mkdir(parents=True, exist_ok=True)
    html_out = render(ns, css, js)
    if re.search(r'href=["\'][^"\']*shell\.css', html_out) or re.search(
        r'src=["\'][^"\']*shell\.js', html_out
    ):
        die("shell points at an external shell.css/shell.js")
    if "--void:" not in html_out or "IntersectionObserver" not in html_out:
        die("shell is missing the CSS tokens or shell.js")
    out.write_text(html_out, encoding="utf-8")
    print(out)


if __name__ == "__main__":
    main()
PY
