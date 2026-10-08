# PRD-013 · Validação

Documento: `docs/prd/PRD-013-extensao-vscode.md`. Validador: `document-validator` (opus), skill `document-review`.

## Rodada 1 de 2 (2026-10-07)

- Review status: changes required
- Blocking findings: 4

## Findings

### BLOCKING: P1 [conflict] R3 fixes the `127.0.0.1:<port>` origin while H2/R7 rely on `asExternalUri`/Remote WSL forwarding, which can yield another origin
- Location: §3 R3; Appendix H2; risk Remote WSL; §1
- Reported by: claude
- Problem: R3 forbids any origin but `127.0.0.1:<port>`; `asExternalUri` may return `localhost:<other port>` or a forwarded host. The dashboard Host guard (`internal/dashboard/server.go:104-111`) checks only the hostname, so forwarded localhost ports pass, but a non-loopback forwarded host gets 403. The webview CSP `frame-src` is not named.
- Suggested fix: embedded URL from `asExternalUri(http://127.0.0.1:<port>)` built in one place; webview CSP `frame-src` lists only that origin; extension-side requests (probe, status, proxy) go to `http://127.0.0.1:<port>`; H2 "Como cairia" names the 403 path to R7.

### BLOCKING: P2 [gap] Where the extension runs under Remote WSL is unstated; R4, R5, R8, R12 depend on it
- Location: §3 R4, R5, R8, R12; §1
- Reported by: claude
- Problem: `dh`, `claude` and the plugin live on the WSL side; a UI extension (Windows side) would look in the Windows home, start the wrong dashboard and open a Windows terminal.
- Suggested fix: rule declaring `extensionKind: ["workspace"]`; lookup, probe, status and terminal run on the remote side; verified in `package.json` and in the probe ticket.

### BLOCKING: P3 [conflict] R9 forbids reading snapshots; AGENTS.md says the extension reads them; §4 does not drop that clause
- Location: §3 R9, R14; §4; pending decision 1
- Reported by: claude
- Problem: R9 is written as settled while decision 1 is open; the AGENTS.md update does not remove "Lê snapshots".
- Suggested fix: make R9 conditional on decision 1 (recommendation assumed) and have §4 say the AGENTS.md rule drops the snapshot clause.

### BLOCKING: P4 [ambiguity] R12 does not decide whether Enter is pressed
- Location: §3 R12; H4; pending decision 5
- Reported by: claude
- Problem: `Terminal.sendText` adds a newline by default; the exact sequence (and how `claude` starts) is undefined, so the check cannot be written.
- Suggested fix: exact sequence: `createTerminal({cwd})`, `claude` + Enter, then after the H4 mechanism `/dh:<nome>` with `addNewLine: false` (decision 5 recommendation; `true` if the owner picks auto-Enter); mock records both calls.

### NON-BLOCKING: P5 [factual error] §1 describes the Host guard as port-bound
- Reported by: claude
- Suggested fix: "valida o hostname do `Host` (`127.0.0.1` ou `localhost`, qualquer porta; `[::1]` recusado)".

### NON-BLOCKING: P6 [ambiguity] R8 does not map `/api/state` values to display text
- Reported by: claude
- Suggested fix: table of the six `avatar.state` values to text, integer rounding, language, whether limit age is shown.

### NON-BLOCKING: P7 [weak criterion] R1/R2 checks name actors who cannot run them
- Reported by: claude
- Suggested fix: manual check by the owner in VS Code Linux and Remote WSL (screenshots in the probe ticket); automated `@vscode/test-electron` under `xvfb-run` for view registration and iframe URL.

### NON-BLOCKING: P8 [organization] Pending decisions 4 and 6 block R11, written as decided
- Reported by: claude
- Suggested fix: mark R11 "decisões pendentes 4 e 6; recomendação assumida" and add the no-workspace/no-roots case as an R11 check.

## Rodada 2 de 2 (2026-10-07)

- Review status: changes required
- Blocking findings: 1 (P9, new); P1–P8 resolved. Cap reached.

## Findings

### BLOCKING: P9 [conflict] R8 maps English keys, but `/api/state` returns `avatar.state` as jevmon ids in Portuguese
- Location: §3 R8
- Reported by: claude
- Problem: `avatar.state` carries `esperando`, `erro`, `trabalhando`, `concluido`, `atencao`, `parado` (`internal/dashboard/sessions.go:73-81`; page table `AVLABEL`, `internal/dashboard/web/index.html:46`). English keys are the session `activity` values. Round-1 P6 caused it.
- Suggested fix: map the jevmon ids (as `AVLABEL`); unknown value shows the raw id; the test uses a real `/api/state` response captured from `dh dashboard` as fixture.

### NON-BLOCKING: P10 [ambiguity] R12 and H4 place the shell-integration wait differently
- Reported by: claude
- Suggested fix: wait for shell integration (or ~1 s without it), `sendText("claude", true)`, fixed 3 s wait (H4), `sendText("/dh:<nome>", false)`.

## Disposição do Coordenador (2026-10-07)

- P9 é erro de fato com correção única, conferida pelo Coordenador em `internal/dashboard/sessions.go:73-81`. Com o teto de 2 rodadas atingido e o dono ausente sob autorização de levar a demanda até o fim, a correção foi aplicada pelo autor sem terceira rodada; o ponto fica listado para o dono no relatório final. P10 aplicado junto.
