import { expect, test } from 'claude-code/testing'
import { cloneAsk, ghDenied, ghWriteVerb } from './gh-guard.ts'

// PRD-027 R2. Denied: any string verb.
const denied = [
  // base set
  'gh pr merge 1 --squash', 'gh pr merge --help', 'cd x && git push', 'GH_TOKEN=t gh pr comment 1 -b hi',
  'gh api -X POST repos/o/r/issues/1/comments', 'gh api repos/o/r/issues/1/comments -f body=hi',
  'gh api repos/o/r/pulls/1 -iXPUT', 'gh api repos/o/r/pulls/1 --method PUT', 'gh api -iXPOST repos/o/r/issues/1/comments',
  'gh api -pfbody=hi repos/o/r/issues/1/comments', 'gh api -H "X-HTTP-Method-Override: PUT" repos/o/r/pulls/1/merge',
  // regression set from the review rounds
  'git status; gh pr merge 1', 'echo "$(true; gh pr merge 1)"', `echo "$(git -c 'user.name=a;b' push)"`,
  `git -c 'user.name=a;b' push`, "echo $'it\\'s'; gh pr merge 1", 'g\\it push', '\\git push', 'git \\\npush',
  'echo `g\\\\it push`', 'echo `g\\\\h pr merge 1`', 'g\\\\it push',
  'sudo -u x gh pr merge 1', 'env -u VAR gh pr merge 1', 'env -u git gh pr merge 1', 'env -u git git push',
  'timeout 30 git push', 'nohup git push', 'nice -n 5 git push', 'time git push', 'xargs gh pr merge 1',
  'stdbuf -oL git push', 'setsid git push',
  '{ git push; }', '! git push', 'if x; then git push; fi', 'for i in 1 2; do gh pr merge 1; done', 'x || true; else git push',
  '>/dev/null gh pr merge 1', '<(gh pr merge 1)', 'tee >(gh pr merge 1)', 'gh.exe pr merge 1', 'git.exe push',
  'bash -c "git push"', 'git push>out.txt', 'gh pr merge>out.txt', 'find . -exec git push \\;',
  'find . -execdir gh pr merge 1 {} +', 'watch -n 1 git push', 'script -qc "git push" /dev/null', 'parallel git push ::: a',
  'flock /tmp/l git push', 'chronic git push', 'ssh host git push', 'git send-pack', 'git subtree push --prefix d origin main',
  'git p4 submit', 'git svn dcommit', 'git -c alias.p=push p', 'git config --global alias.p push && git p',
  'gh alias set x y', 'gh auth status -t', 'gh auth status --show-token', "rg --pre 'git push' x .",
  "ag --pager 'gh pr merge 1' x", "echo 'a && git push'",
  // round 4: value flags, auth status flags, word boundaries, dash forms
  'gh api repos/o/r --jq -R --method=DELETE', 'gh api repos/o/r --repo --method=DELETE', 'gh api repos/o/r --jq',
  'gh api repos/o/r -R o/r', 'gh pr view 1 -R', 'gh auth status --show-token=true', 'gh auth status -t=true',
  'gh auth status -at', 'gh auth status -h', 'gh auth status -R o/r', 'gh auth -R o/r status',
  "git -c core.sshCommand='gh pr merge 1 #' fetch", "git --upload-pack='gh pr merge 1 #' fetch x",
  'GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=alias.p GIT_CONFIG_VALUE_0=push git p', '$(git --exec-path)/git-push',
  'git-push', 'git-send-pack x', 'git-subtree push', 'git-p4 submit', 'git-svn dcommit', 'x=gh pr merge 1', 'echo gh',
  "grep -n 'gh;' README.md", 'git log --grep=push', 'npm run x-gh', 'git --upload-pack=gh fetch x', 'echo gh)',
  // other gh families
  'gh issue close 3', 'gh release create v1', 'gh repo delete o/r', 'gh auth token', 'gh workflow run ci.yml',
  'gh api --method DELETE repos/o/r', 'gh api --input x repos/o/r', 'gh secret set X', 'gh secret delete X',
  'gh variable set X', 'gh label create x', 'gh run rerun 1', 'gh run cancel 1', 'gh run delete 1', 'gh workflow enable x',
  'gh workflow disable x', 'gh gist create f', 'gh pr update-branch 1', 'gh pr lock 1', 'gh cache delete 1',
  'gh repo deploy-key add k', 'gh release delete-asset v1 a',
  // new with the short rule
  'ls ~/.config/gh', 'cat ~/.config/gh/hosts.yml', '$(echo gh) pr merge 1', "gh pr view 1 --jq '.x'",
  // amendment 6
  'GH pr merge 1', 'Git push', 'git-lfs push origin main', 'gh api --hostname evil.example repos/o/r',
  'gh api --hostname=evil.example repos/o/r', 'gh api https://evil.example/x', 'GH pr view 1',
  // accepted false positives (PRD 2.3)
  'cd x && gh pr view 1', "rg 'git push' docs",
  // BUG-053 item 4: the write word is the git subcommand, after global options
  'git -C x p4 submit', 'git -c a=b svn dcommit', 'git --git-dir x push', 'git --git-dir=x --work-tree y subtree push',
  'git --no-pager -C x push', 'git --namespace n push', 'echo x | git -C d push',
  // BUG-053 fix-r1 F1: the write word anywhere in the git segment; xargs-fed git; accepted segment false positives
  'git -C "/tmp/a b" push', 'git -c "user.name=a b" push', 'echo push | xargs git', "git $'push'", 'git {push,}',
  'git --attr-source HEAD push', 'git log --grep push', 'git commit -m "push the fix"', 'git 2>&1 push',
  // owner decision 2026-11-04: `git` plus a write word anywhere is denied; the p4 false positive is accepted
  "sed 's/p4/x/' f | git commit-tree", 'git log --oneline | grep p4', 'git status && echo p4', 'grep -n "git" f | sed "s/p4/x/"',
  // F3: a bare glob token stays out of the read form
  'gh api repos/o/r/issues ?F title=x', 'gh api ?x repos/o/r',
  // fix-r2 G1: with `$` or a backtick the verb may come from an expansion, so the write word counts anywhere
  'git $(:; echo push)', 'git $(\necho push)', 'x=push; git $x', 'set -- push; git "$@"', 'for v in push; do git $v; done',
  'git -C $(pwd | cat) push', 'git `echo push`',
]
const allowed = [
  'gh pr view 1 --json title', 'gh pr diff 1', 'gh api repos/o/r/pulls/1', 'gh api repos/o/r/pulls/1 --jq .title',
  'grep -rn "git push" docs/', 'grep -rn "gh pr view" .commands', 'git log --oneline', 'git status && git diff',
  `git -c 'user.name=a;b' status`, 'env -u git git status', 'gh pr checks 1 --watch', 'gh run watch 123', 'git fetch origin',
  'gh pr view 45 --json state,mergedAt', 'gh run list --limit 3', 'go test ./... 2>&1 | tail -n 20',
  'sha256sum $(git ls-files -m)', `git log --format='%h (%s)' -n 3`, 'npm run gh-pages', 'gh auth status', 'gh --version',
  `echo "$(git -c 'user.name=a;b' status)"`, 'gh api -i repos/o/r/pulls/1',
  // one case for each other read of 2.2, and the optional -R / --repo
  'gh', 'gh status', 'gh pr list', 'gh pr status', 'gh issue view 1', 'gh issue list', 'gh issue status', 'gh repo view',
  'gh repo list', 'gh run view 1', 'gh release view v1', 'gh release list', 'gh workflow view x', 'gh workflow list',
  'gh search issues foo', 'gh -R o/r pr view 1', 'gh pr view 1 --repo o/r', 'gh pr view 1 --help',
  'gh api repos/o/r --paginate --include --silent', 'gh api repos/o/r --cache 1h -q .x', 'ghq list',
  'git status --porcelain -uall', '/usr/bin/gh pr view 1', 'gh auth status -h github.com', 'gh auth status --active',
  'gh auth status --hostname=h.example -a', 'gh pr view 1 -R o/r', 'gh pr -R o/r view 1', 'gh run watch 123', 'git log --oneline',
  'npm run gh-pages', 'ghq list', 'git status', 'gh pr view 1', 'gh api repos/o/r --jq .name', 'git log', 'ls gh-pages',
  // BUG-053 items 4 and 5
  'gh api repos/o/r/contents/p?ref=abc123',
  'gh api repos/o/r/contents/p?ref=abc --jq .sha',
  // fix-r2 G1: `$` without a write word, and a stray `git` / `p4` without `$`, still pass
  'git log $(git rev-parse HEAD)', 'git -C $(pwd) status', 'echo $HOME | git commit-tree',
]

test('ghWriteVerb: denied commands return a verb', () => {
  for (const cmd of denied) expect([cmd, typeof ghWriteVerb(cmd)]).toEqual([cmd, 'string'])
})

test('ghWriteVerb: allowed commands return undefined', () => {
  for (const cmd of allowed) expect([cmd, ghWriteVerb(cmd)]).toEqual([cmd, undefined])
})

// PowerShell: the backtick is removed and no exception applies.
test('ghWriteVerb: PowerShell', () => {
  expect(ghWriteVerb('g`h pr merge 1', true)).toBe('gh pr merge')
  expect(ghWriteVerb('gh pr view 1', true)).toBe('gh pr view')
  expect(ghWriteVerb('grep -rn "git push" docs/', true)).toBe('git push')
  expect(ghWriteVerb('gh pr view 1', false)).toBe(undefined)
  expect(ghWriteVerb('C:\\Program Files\\GitHub CLI\\gh.exe pr merge 1', true)).toBe('gh pr merge')
  expect(ghDenied('gh pr view 1', 'a', 'dh:qa-verifier', false, true)).toEqual({ verb: 'gh pr view', unknown: false })
})

test('ghWriteVerb: the verb names the command', () => {
  expect(ghWriteVerb('gh pr merge 1 --squash')).toBe('gh pr merge')
  expect(ghWriteVerb('gh -R o/r pr merge 1')).toBe('gh pr merge')
  expect(ghWriteVerb('gh api -X POST repos/o/r')).toBe('gh api')
  expect(ghWriteVerb('cd x && git push')).toBe('git push')
  expect(ghWriteVerb('ls ~/.config/gh')).toBe('gh')
})

test('ghDenied: scope by agent type', () => {
  expect(ghDenied('git push', undefined, undefined)).toBe(undefined)
  expect(ghDenied('git push', 'a', 'dh:release-manager')).toBe(undefined)
  expect(ghDenied('git push', 'a', 'vault-agent')).toBe(undefined)
  expect(ghDenied('git push', 'a', 'dh:qa-verifier')).toEqual({ verb: 'git push', unknown: false })
  expect(ghDenied('git push', 'a', undefined)).toEqual({ verb: 'git push', unknown: true })
  expect(ghDenied('gh pr view 1', 'a', 'dh:qa-verifier')).toBe(undefined)
})

// BUG-053 item 5: `?` only; `&` and the flag table stay shut.
test('ghWriteVerb: api query string', () => {
  expect(ghWriteVerb('gh api repos/o/r/contents/p?ref=sha&x=1')).toBe('gh api')
  expect(ghWriteVerb('gh api repos/o/r/contents/p?ref=sha -f a=b')).toBe('gh api')
})

// Issue #54: a remote clone or download asks the owner (never denies), naming the repository.
const asks: [string, string][] = [
  ['git clone https://github.com/o/r.git', 'o/r'], ['git clone --depth 1 -b main git@github.com:o/r.git d', 'o/r'],
  ['cd /tmp && git clone ssh://git@host/o/r', 'o/r'], ['gh repo clone o/r', 'o/r'], ['gh repo clone o/r -- --depth 1', 'o/r'],
  ['gh api repos/o/r/tarball/main', 'o/r'], ['gh api /repos/o/r/zipball', 'o/r'], ['gh api repos/o/r/tarball/v1 > x.tgz', 'o/r'],
  ['git fetch https://github.com/o/r.git', 'o/r'], ['git fetch git@github.com:o/r.git main', 'o/r'],
  ['git -C x clone https://h.example/o/r', 'o/r'], ['gh.exe repo clone o/r', 'o/r'],
  // fix-r1 F2 and F4
  ['git clone host:org/repo.git', 'org/repo'], ['git fetch host:org/repo.git', 'org/repo'], ['git clone -c "a=b c" host:org/repo.git', 'org/repo'],
  ['gh repo fork o/r --clone', 'o/r'], ['git pull https://h.example/o/r', 'o/r'], ['git submodule add https://h.example/o/r d', 'o/r'],
  ['git submodule add git@h.example:o/r.git', 'o/r'],
  // fix-r2 G2: scp remote with a single path component; the repository is the first non-option argument
  ['git clone host:repo', 'repo'], ['git clone --depth 1 -b main host:repo d', 'repo'], ['git fetch host:repo', 'repo'],
  ['git pull host:repo main', 'repo'], ['git submodule add -b main host:repo d', 'repo'], ['git clone -c "a=b c" host:repo', 'repo'],
  // fix-r2: clone / submodule add check every non-option argument (an unknown value option cannot hide the repository)
  ['git clone ./local host:d', 'd'], ['git clone --reference-if-able /x https://h/o/r', 'o/r'], ['git clone --unknown-opt v host:o/r', 'o/r'],
  ['git submodule add --unknown-opt v host:o/r', 'o/r'],
  // owner decision 2026-11-04: word-based and fail-safe; no subcommand parsing
  ['git -C fetch clone https://host/o/r', 'o/r'], ['git -c "user.name=\\";x" clone host:o/r', 'o/r'],
  // fix-r4 H1: fork clone flag forms; H2: a remote hidden inside a token
  ['gh repo fork o/r -c', 'o/r'], ['gh repo fork o/r --clone=1', 'o/r'], ['gh repo fork o/r --clone=t', 'o/r'],
  ['git clone {git@github.com:o/r,}', 'o/r'], ['git -c remote.x.url=git@h:o/r fetch x', 'o/r'],
  // fix-r5 J1: an empty repository name still asks; J2: remote add / archive --remote; J3: autocorrect hides the verb
  ['git clone git@host:.git', 'a remote repository'], ['git remote add -f o https://host/a/b', 'a/b'],
  ['git archive --remote=ssh://host/a/b HEAD', 'a/b'], ['git -c help.autocorrect=1 clon https://host/a/b', 'a/b'],
  ['git fetch origin main:main', 'main'], ['git pull origin host:x', 'x'], ['git clone ./local host:d', 'd'],
]
const noAsk = [
  'git clone ../local d', 'git clone /srv/repo.git', 'git fetch origin', 'git fetch', 'gh search code foo',
  'gh api repos/o/r/contents/README.md', 'gh api repos/o/r/contents/p?ref=sha', 'gh pr view 1', 'git log --oneline',
  'grep -rn "git clone https://x.example/o/r" docs/', 'git status',
  'git pull', 'git pull origin main', 'git submodule update', 'gh repo fork o/r', 'gh repo view o/r --clone',
  // fix-r2 G2: later arguments (refspecs, directories) are not the repository; one-letter hosts are drive letters
  'git fetch origin a:bc', 'git remote -v', 'git remote add o ./local', 'git clone C:/x d', 'git clone ./a:b',
]

test('cloneAsk: remote clone and download forms name the repository', () => {
  for (const [cmd, repo] of asks) expect([cmd, cloneAsk(cmd)]).toEqual([cmd, repo])
})

test('cloneAsk: reads and local clones produce no ask', () => {
  for (const cmd of noAsk) expect([cmd, cloneAsk(cmd)]).toEqual([cmd, undefined])
})

test('cloneAsk: PowerShell', () => {
  expect(cloneAsk('git clone https://github.com/o/r', true)).toBe('o/r')
  expect(cloneAsk('git clone C:\\x d', true)).toBe(undefined)
  // fix-r6: a line continuation (backtick or backslash + newline) is a space, so it cannot hide the verb
  for (const nl of ['\n', '\r\n']) {
    expect(cloneAsk('gh repo `' + nl + 'clone o/r', true)).toBe('o/r')
    expect(cloneAsk('gh repo fork o/r `' + nl + '--clone', true)).toBe('o/r')
    expect(cloneAsk('gh api `' + nl + 'repos/o/r/tarball', true)).toBe('o/r')
    expect(cloneAsk('git `' + nl + 'clone https://h/o/r', true)).toBe('o/r')
    expect(cloneAsk('gh repo \\' + nl + 'clone o/r')).toBe('o/r')
    expect(cloneAsk('git clone \\' + nl + 'https://h/o/r')).toBe('o/r')
  }
})
