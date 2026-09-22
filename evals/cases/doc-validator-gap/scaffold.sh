#!/usr/bin/env bash
# Copies this case's fixtures into the throwaway workspace cwd (runner: claude plugin eval --scaffold).
set -euo pipefail
rm -rf ./fixtures && cp -r "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/fixtures" ./fixtures
