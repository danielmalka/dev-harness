# Review guide

## How to use

This file tells the reviewer what to look for in this project. The owner fills it. The reviewer reads it first and ties each finding to one rule below, or says no rule applies. Keep entries few and concrete. Always write it in English. The review reads this file from the base branch, so a change to it stops the merge for the owner.

## Rules

One entry per recurring problem. Copy the block for each new rule.

### <Rule title>

- Symptom: <what the reviewer sees in the diff>
- Why it matters: <what breaks or what it costs>
- Example: `<path/to/file.ext>:<line>`

## Ignore

One entry per pattern that is noise. Each entry reads "ignore this pattern in this kind of file". A finding that matches an entry is reported as ignored-by-rule and does not spend a round.

- Ignore <pattern> in <kind of file>.

## Critical paths

Path globs in `.gitignore` style. A pull request that touches a match stops before the merge, and the owner decides. Leave the list empty for no stop. Each entry is one glob in backticks; text after it is a reason, not part of the glob.

- `src/payments/**`
