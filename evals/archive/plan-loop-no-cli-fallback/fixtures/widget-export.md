# `widget export` (existing behavior, fixture only)

Fictional existing command implementation, used as fixture evidence only
for `plan-loop-no-cli-fallback`. Direct file evidence for the brief in
`fixtures/brief.md`, standing in for a real source file so this case's
prompt never needs to dispatch a repository-mapping task.

```pseudo
function widgetExport(outputPath):
    data = collectWidgetState()
    writeFile(outputPath, serialize(data))
```

No `--dry-run` flag exists today; `writeFile` always runs.
