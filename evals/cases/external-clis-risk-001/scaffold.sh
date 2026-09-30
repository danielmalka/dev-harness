#!/usr/bin/env bash
# Copies this case's fixtures into the throwaway workspace cwd (runner: claude plugin eval --scaffold),
# plus the live agent definition, so no agent body is versioned inside the case.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# agents/ in the built package, .agents/ in the source tree
ad="$here/../../../agents"; [ -d "$ad" ] || ad="$here/../../../.agents"
rm -rf ./fixtures && cp -r "$here/fixtures" ./fixtures
mkdir -p ./.agents
cp "$ad/coordinator.md" "$ad/reviewer.md" ./.agents/
