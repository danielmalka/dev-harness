# Brief · slice-01

Reject a blank name in `greet`.

## Current behavior

`greet(name)` returns `"Hello, {name}"` for any string, including whitespace-only values.

## Desired behavior

When `name` is empty or only whitespace, `greet` raises `ValueError` with a message that mentions `name`. Non-blank names keep the current greeting.

## Acceptance

| ID | When | Then |
| --- | --- | --- |
| AC-01 | `greet("Ada")` | returns `Hello, Ada` |
| AC-02 | `greet("  ")` or `greet("")` | raises `ValueError` |

## Scope

- Write set: `greet.py`. Tests already live in `test_greet.py`.
- Out: new dependencies, extra public functions, renaming files.

## Checks

```bash
python3 -m unittest test_greet.py
```

Today AC-01 passes and AC-02 fails. After the slice both pass. Typecheck does not apply.
