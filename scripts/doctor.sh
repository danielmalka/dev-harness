#!/usr/bin/env bash
# Observe a Dev Harness package or checkout. Never installs anything.
# Usage: doctor.sh [TARGET]
set -euo pipefail

SELF="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET="${1:-}"

if [ -z "$TARGET" ]; then
  if [ -d "$SELF/../agents" ] && [ -d "$SELF/../.claude-plugin" ]; then
    TARGET="$(cd "$SELF/.." && pwd)"
  elif [ -d "$SELF/../.agents" ]; then
    TARGET="$(cd "$SELF/.." && pwd)"
  else
    echo "usage: doctor.sh [plugin-dir-or-kit-root]" >&2
    exit 2
  fi
else
  TARGET="$(cd "$TARGET" && pwd)"
fi

VALIDATE="$SELF/validate.py"
if [ ! -f "$VALIDATE" ]; then
  echo "missing validator: $VALIDATE" >&2
  exit 1
fi

echo "## Mode"
echo "diagnosis (no writes)"
echo
echo "## Kit"
echo "target: $TARGET"
echo "validator: $VALIDATE"

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3: missing"
  echo "status: blocked (validator cannot run)"
  exit 1
fi

set +e
python3 "$VALIDATE" "$TARGET"
VALIDATE_STATUS=$?
set -e

echo
echo "## Environment"
probe() {
  local layer="$1"
  local check="$2"
  shift 2
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '| %s | %s | %s | missing | command not found |\n' "$layer" "$check" "$*"
    return
  fi
  local out
  if out="$("$@" 2>&1)"; then
    local trimmed
    trimmed="$(printf '%s' "$out" | tr '\n' ' ' | cut -c1-120)"
    printf '| %s | %s | %s | present | %s |\n' "$layer" "$check" "$*" "$trimmed"
  else
    printf '| %s | %s | %s | unknown | probe failed |\n' "$layer" "$check" "$*"
  fi
}

echo "| Layer | Check | Command or observation | Status | Evidence |"
echo "|---|---|---|---|---|"
probe runtime python3 python3 --version
probe runtime git git --version
if command -v claude >/dev/null 2>&1; then
  probe runtime claude claude --version
else
  echo "| runtime | claude | claude --version | missing | command not found |"
fi

echo
echo "## Harness records in current directory"
CWD="$(pwd)"
if [ -d "$CWD/.harness" ]; then
  echo ".harness/: present"
  for f in project.yaml MEMORY.md EPOCHAL.md RISKS.md; do
    if [ -e "$CWD/.harness/$f" ]; then
      echo ".harness/$f: present"
    else
      echo ".harness/$f: missing"
    fi
  done
else
  echo ".harness/: missing"
fi

echo
echo "## Limits"
echo "- Nothing was installed, upgraded or configured."
echo "- Claude Code discovery of commands and agents is unverified unless this session loaded the plugin."
echo "- A missing claude binary blocks loading the plugin; it does not block reading the docs."
echo
if [ "$VALIDATE_STATUS" -eq 0 ]; then
  echo "status: passed"
else
  echo "status: failed (validator errors above)"
fi
exit "$VALIDATE_STATUS"
