import { expect, test } from 'claude-code/testing'
import { ghDenied, ghWriteVerb } from './gh-guard.ts'

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
  'cd x && gh pr view 1', 'git log --grep push', 'git log --grep=push', 'git commit -m "push the fix"', "rg 'git push' docs",
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
