# Fixture slice-01

Tiny Python 3 stdlib project used as the first guided delivery. No packages to install.

The starting tree is intentionally incomplete: `greet("  ")` should raise `ValueError` and currently does not. `python3 -m unittest test_greet.py` is red until the authorized slice is implemented.

## Layout

| File | Role |
| --- | --- |
| `BRIEF.md` | Authorized request and acceptance. |
| `greet.py` | Code under test. |
| `test_greet.py` | Acceptance tests. Happy path passes. Blank-name case fails. |
| `expected/greet.py` | Known-good result. Do not copy it into place during the exercise. |
| `expected/unittest.txt` | Expected unittest output after the slice. |

## Completion

1. Load the kit against this directory (`docs/inicio-rapido.html`).
2. `/dev-harness:setup` then `/dev-harness:plan` using `BRIEF.md`.
3. `/dev-harness:build` the single slice, then `/dev-harness:verify` and `/dev-harness:review`.
4. `python3 -m unittest test_greet.py` matches `expected/unittest.txt`.
5. `/dev-harness:handoff`.

Python 3 is a declared prerequisite of the kit's own scripts. The fixture uses it so a clean machine does not need Node, Go or a browser.
