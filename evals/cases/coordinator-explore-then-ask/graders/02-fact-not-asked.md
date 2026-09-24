---
type: llm
weight: 1
focus: last_message
---
The user's prompt asks for "the exact command that runs behind the
project's test script", not for the name of the script itself. The fact
discoverable by reading `fixtures/package.json` is the runner invocation
`jest --runInBand`, declared as the value of the `test` script. Naming
only the wrapper — "the test script", "npm test", or any other way of
pointing at the script without unpacking what it runs — does not answer
what was asked, so it does not count as stating the fact. Count the fact
as stated when the reply gives `jest --runInBand`, whether alone or
together with the wrapper that invokes it (for example "npm test, which
runs `jest --runInBand`" or "npm test → jest --runInBand"), with or
without a working-directory prefix such as `cd fixtures &&`. Fail this
criterion only if the final reply asks the user what the command is
instead of stating it, or never states `jest --runInBand` in any of
these forms. Any other question the reply may also ask does not by
itself fail this criterion.
