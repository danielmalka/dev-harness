# Clone creation fails (R11)

Fixture data only. Task `T-2001`, `<call-id>` = `T-2001-grok-3`,
`<tmp root>` = `/tmp`. The stage is a `code` review with `cli:grok/grok-4.7`.

```
test -d <project>/.git && test ! -L <project>/.git -> exit 0
CLONE=$(mktemp -d "/tmp/dh-cli-T-2001-grok-3-XXXXXXXX") -> /tmp/dh-cli-T-2001-grok-3-Vf2yN6qD
git clone --no-hardlinks <project> "$CLONE"  -> exit 128
    fatal: write error: No space left on device
```
