#!/usr/bin/env python3
"""Static checks for Dev Harness sources or a generated Claude Code package."""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

FRONTMATTER_RE = re.compile(r"\A---\n(.*?)\n---\n", re.S)
NAME_RE = re.compile(r"^name:\s*[|]?\s*(.*)$", re.M)
MODEL_RE = re.compile(r"^model:\s*(\S+)", re.M)
LINK_RE = re.compile(r"\[[^\]]*\]\(([^)]+)\)")
TEMPLATE_RE = re.compile(r"templates/(?:(<[a-z-]+>|[a-z]{2}(?:-[a-z]{2})?)/)?([A-Za-z0-9_.-]+\.md)")
TEMPLATE_LANGS = ("en", "pt-br")
PRIVATE_RE = re.compile(r"(/home/|/Users/|[A-Za-z]:\\)")
SKIP_LINK_PREFIXES = ("http://", "https://", "mailto:", "#")
ALLOWED_MODELS = {"haiku", "sonnet", "opus"}
MIN_AGENTS = 18
MIN_COMMANDS = 17
MIN_SKILLS = 16
REQUIRED_PROFILES = ("base", "go-api", "typescript-web")

# Informational mentions of third-party tools or other machines, not runtime paths.
PRIVATE_ALLOW_FRAGMENTS = (
    "code.claude.com",
)


def die(msg: str, code: int = 1) -> None:
    print(msg, file=sys.stderr)
    sys.exit(code)


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def frontmatter(text: str) -> str:
    m = FRONTMATTER_RE.match(text)
    return m.group(1) if m else ""


def fm_name(fm: str) -> str:
    m = NAME_RE.search(fm)
    if not m:
        return ""
    value = m.group(1).strip()
    if value in {"|", ""}:
        return ""
    return value.strip("\"'")


def detect_layout(target: Path, source_only: bool = False) -> str:
    if (target / "agents").is_dir() and (target / ".claude-plugin").is_dir():
        return "package"
    if (target / ".agents").is_dir() and (target / ".commands").is_dir():
        pkg = target / "dist" / "claude-code" / "dev-harness"
        if not source_only and pkg.is_dir() and (pkg / ".claude-plugin").is_dir():
            return "kit"
        return "source"
    die(f"unrecognized layout: {target}")
    return ""


def layout_dirs(target: Path, layout: str) -> dict[str, Path]:
    if layout == "source":
        return {
            "root": target,
            "agents": target / ".agents",
            "commands": target / ".commands",
            "skills": target / ".skills",
            "templates": target / "templates",
            "profiles": target / "profiles",
            "scripts": target / "scripts",
        }
    if layout == "package":
        return {
            "root": target,
            "agents": target / "agents",
            "commands": target / "commands",
            "skills": target / "skills",
            "templates": target / "templates",
            "profiles": target / "profiles",
            "scripts": target / "scripts",
        }
    pkg = target / "dist" / "claude-code" / "dev-harness"
    return {
        "root": target,
        "agents": target / ".agents",
        "commands": target / ".commands",
        "skills": target / ".skills",
        "templates": target / "templates",
        "profiles": target / "profiles",
        "scripts": target / "scripts",
        "package": pkg,
    }


def list_md(directory: Path) -> list[Path]:
    if not directory.is_dir():
        return []
    return sorted(p for p in directory.glob("*.md") if p.name != ".gitkeep")


def list_skills(directory: Path) -> list[Path]:
    if not directory.is_dir():
        return []
    found = []
    for child in sorted(directory.iterdir()):
        skill = child / "SKILL.md"
        if child.is_dir() and skill.is_file():
            found.append(skill)
    return found


def list_templates(directory: Path) -> dict[str, set[str]]:
    """Templates live in templates/<lang>/NAME.md. Returns {lang: {NAME.md, ...}}."""
    found: dict[str, set[str]] = {}
    if not directory.is_dir():
        return found
    for lang_dir in sorted(directory.iterdir()):
        if lang_dir.is_dir() and lang_dir.name in TEMPLATE_LANGS:
            found[lang_dir.name] = {p.name for p in lang_dir.glob("*.md")}
    return found


def list_profiles(directory: Path) -> list[Path]:
    if not directory.is_dir():
        return []
    return sorted(p for p in directory.glob("*.yaml"))


def inventory(dirs: dict[str, Path]) -> dict[str, list[str]]:
    agents = [p.stem for p in list_md(dirs["agents"])]
    commands = [p.stem for p in list_md(dirs["commands"])]
    skills = [p.parent.name for p in list_skills(dirs["skills"])]
    profiles = [p.stem for p in list_profiles(dirs["profiles"])]
    templates = sorted(f"{lang}/{name}" for lang, names in list_templates(dirs["templates"]).items() for name in sorted(names))
    return {
        "agents": agents,
        "commands": commands,
        "skills": skills,
        "profiles": profiles,
        "templates": templates,
    }


def check_counts(inv: dict[str, list[str]], errors: list[str]) -> None:
    if len(inv["agents"]) < MIN_AGENTS:
        errors.append(f"agents: {len(inv['agents'])} < {MIN_AGENTS}")
    if len(inv["commands"]) < MIN_COMMANDS:
        errors.append(f"commands: {len(inv['commands'])} < {MIN_COMMANDS}")
    if len(inv["skills"]) < MIN_SKILLS:
        errors.append(f"skills: {len(inv['skills'])} < {MIN_SKILLS}")
    for name in REQUIRED_PROFILES:
        if name not in inv["profiles"]:
            errors.append(f"missing profile: {name}")


def check_agents(dirs: dict[str, Path], errors: list[str]) -> None:
    for path in list_md(dirs["agents"]):
        text = read(path)
        fm = frontmatter(text)
        if not fm:
            errors.append(f"{path}: missing frontmatter")
            continue
        name = fm_name(fm)
        if name and name != path.stem:
            errors.append(f"{path}: name {name!r} != {path.stem!r}")
        if not name:
            errors.append(f"{path}: missing name")
        model = MODEL_RE.search(fm)
        if not model:
            errors.append(f"{path}: missing model")
        elif model.group(1) not in ALLOWED_MODELS:
            errors.append(f"{path}: model {model.group(1)!r} not in {sorted(ALLOWED_MODELS)}")


def check_skills(dirs: dict[str, Path], errors: list[str]) -> None:
    for path in list_skills(dirs["skills"]):
        fm = frontmatter(read(path))
        name = fm_name(fm)
        folder = path.parent.name
        if name != folder:
            errors.append(f"{path}: name {name!r} != folder {folder!r}")


def check_commands(dirs: dict[str, Path], errors: list[str]) -> None:
    for path in list_md(dirs["commands"]):
        if not frontmatter(read(path)):
            errors.append(f"{path}: missing frontmatter")


def check_profile_ids(dirs: dict[str, Path], errors: list[str]) -> None:
    for path in list_profiles(dirs["profiles"]):
        text = read(path)
        m = re.search(r"^id:\s*(\S+)", text, re.M)
        if not m:
            errors.append(f"{path}: missing id")
        elif m.group(1) != path.stem:
            errors.append(f"{path}: id {m.group(1)!r} != {path.stem!r}")


def iter_text_files(root: Path) -> list[Path]:
    skip_parts = {".git", "dist", "node_modules", "__pycache__", ".codegraph"}
    files: list[Path] = []
    for path in root.rglob("*"):
        if not path.is_file():
            continue
        if any(part in skip_parts for part in path.parts):
            continue
        if path.suffix.lower() in {".md", ".yaml", ".yml", ".html", ".py", ".sh", ".json", ".txt"}:
            files.append(path)
    return files


def resolve_link(source: Path, href: str) -> Path | None:
    raw = href.split("#", 1)[0].strip()
    if not raw or raw.startswith(SKIP_LINK_PREFIXES):
        return None
    if re.match(r"^[a-z]+:", raw):
        return None
    return (source.parent / raw).resolve()


def check_links(root: Path, errors: list[str]) -> None:
    for path in iter_text_files(root):
        if path.suffix != ".md":
            continue
        try:
            text = read(path)
        except UnicodeDecodeError:
            continue
        for href in LINK_RE.findall(text):
            target = resolve_link(path, href)
            if target is None:
                continue
            if not str(target).startswith(str(root.resolve())):
                # Links that escape the tree are unresolved for portability.
                if href.startswith(("../", "./", "/")) or "/" in href:
                    errors.append(f"{path}: link escapes tree: {href}")
                continue
            if not target.exists():
                errors.append(f"{path}: unresolved link {href}")


def check_template_refs(dirs: dict[str, Path], errors: list[str]) -> None:
    by_lang = list_templates(dirs["templates"])
    for lang in TEMPLATE_LANGS:
        if lang not in by_lang:
            errors.append(f"templates: missing language directory {lang}/")
    langs = [l for l in TEMPLATE_LANGS if l in by_lang]
    if len(langs) > 1:
        base = by_lang[langs[0]]
        for lang in langs[1:]:
            for name in sorted(base ^ by_lang[lang]):
                errors.append(f"templates: {name} is not present in every language ({langs[0]} vs {lang})")
    roots = [dirs["agents"], dirs["commands"], dirs["skills"]]
    for folder in roots:
        if not folder.is_dir():
            continue
        for path in folder.rglob("*"):
            if not path.is_file() or path.suffix not in {".md", ".yaml"}:
                continue
            for lang, name in TEMPLATE_RE.findall(read(path)):
                targets = [lang] if lang and not lang.startswith("<") else langs
                for target in targets:
                    if name not in by_lang.get(target, set()):
                        errors.append(f"{path}: missing template {target}/{name}")


def check_private_paths(root: Path, errors: list[str]) -> None:
    for path in iter_text_files(root):
        rel = path.relative_to(root)
        # Generated HTML may mention the skill path in comments; still flag kit sources.
        if path.suffix == ".html" or path.name == "validate.py":
            continue
        try:
            text = read(path)
        except UnicodeDecodeError:
            continue
        if not PRIVATE_RE.search(text):
            continue
        if any(frag in text for frag in PRIVATE_ALLOW_FRAGMENTS) and "/home/" not in text and "/Users/" not in text:
            continue
        # Flag only lines that look like author-machine paths.
        for i, line in enumerate(text.splitlines(), 1):
            if PRIVATE_RE.search(line) and "code.claude.com" not in line:
                errors.append(f"{rel}:{i}: private path")


def check_plugin(package: Path, errors: list[str]) -> dict:
    plugin = package / ".claude-plugin" / "plugin.json"
    if not plugin.is_file():
        errors.append(f"{plugin}: missing")
        return {}
    data = json.loads(read(plugin))
    if data.get("name") != "dev-harness":
        errors.append("plugin.json name must be dev-harness")
    if data.get("license") != "MIT":
        errors.append("plugin.json license must be MIT")
    if not data.get("version"):
        errors.append("plugin.json missing version")
    return data


def check_coverage(source_dirs: dict[str, Path], package: Path, errors: list[str]) -> None:
    pairs = (
        (source_dirs["agents"], package / "agents"),
        (source_dirs["commands"], package / "commands"),
        (source_dirs["templates"], package / "templates"),
        (source_dirs["profiles"], package / "profiles"),
    )
    for src, dst in pairs:
        src_files = {p.relative_to(src) for p in src.rglob("*") if p.is_file() and p.name != ".gitkeep"}
        if not dst.is_dir():
            errors.append(f"package missing {dst}")
            continue
        dst_files = {p.relative_to(dst) for p in dst.rglob("*") if p.is_file() and p.name != ".gitkeep"}
        missing = src_files - dst_files
        extra = dst_files - src_files
        for item in sorted(missing):
            errors.append(f"package missing {src.name}/{item}")
        for item in sorted(extra):
            # generated files are allowed at package root, not inside copied trees
            errors.append(f"package extra {dst.name}/{item}")
    src_skills = {p.parent.name for p in list_skills(source_dirs["skills"])}
    dst_skills = {p.parent.name for p in list_skills(package / "skills")}
    for name in sorted(src_skills - dst_skills):
        errors.append(f"package missing skill {name}")
    for name in sorted(dst_skills - src_skills):
        errors.append(f"package extra skill {name}")


def print_report(target: Path, layout: str, inv: dict[str, list[str]], errors: list[str]) -> None:
    print(f"target: {target}")
    print(f"layout: {layout}")
    for key, values in inv.items():
        print(f"{key}: {len(values)}")
        print("  " + ", ".join(values) if values else "  (none)")
    if errors:
        print(f"errors: {len(errors)}")
        for err in errors:
            print(f"  - {err}")
        return
    print("errors: 0")
    print("status: passed")


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate Dev Harness sources or package.")
    parser.add_argument("target", nargs="?", default=".", help="kit root, source tree or plugin package")
    parser.add_argument("--json", action="store_true", help="machine-readable inventory plus errors")
    parser.add_argument(
        "--source-only",
        action="store_true",
        help="ignore dist/ even when present (used before a rebuild)",
    )
    args = parser.parse_args()
    target = Path(args.target).resolve()
    if not target.exists():
        die(f"missing target: {target}")
    layout = detect_layout(target, source_only=args.source_only)
    dirs = layout_dirs(target, layout)
    errors: list[str] = []

    check_agents(dirs, errors)
    check_skills(dirs, errors)
    check_commands(dirs, errors)
    check_profile_ids(dirs, errors)
    inv = inventory(dirs)
    check_counts(inv, errors)
    check_links(dirs["root"] if layout != "kit" else dirs["root"], errors)
    check_template_refs(dirs, errors)
    check_private_paths(dirs["root"] if layout != "package" else dirs["root"], errors)

    plugin = {}
    if layout == "package":
        plugin = check_plugin(target, errors)
        for required in ("scripts/doctor.sh", "scripts/validate.py", "harness-manifest.json", "GENERATED.txt"):
            if not (target / required).is_file():
                errors.append(f"package missing {required}")
    elif layout == "kit":
        pkg = dirs["package"]
        plugin = check_plugin(pkg, errors)
        check_coverage(dirs, pkg, errors)
        for required in ("scripts/doctor.sh", "scripts/validate.py", "harness-manifest.json", "GENERATED.txt"):
            if not (pkg / required).is_file():
                errors.append(f"package missing {required}")
        license_file = target / "LICENSE"
        if not license_file.is_file():
            errors.append("LICENSE missing at kit root")
        else:
            text = read(license_file)
            if "MIT License" not in text:
                errors.append("LICENSE is not MIT")

    payload = {
        "target": str(target),
        "layout": layout,
        "inventory": inv,
        "plugin": plugin,
        "errors": errors,
        "status": "failed" if errors else "passed",
    }
    if args.json:
        print(json.dumps(payload, indent=2, ensure_ascii=False))
    else:
        print_report(target, layout, inv, errors)
    return 1 if errors else 0


if __name__ == "__main__":
    sys.exit(main())
