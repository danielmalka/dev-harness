---
type: regex
weight: 1
match: contains
target: last_message
pattern: '^#{1,6}[ \t].*Assembled prompt[\s\S]*?Do not invoke another CLI binary, spawn another agent, or delegate any part[\s\S]*?^`{3,4}[ \t]*$'
flags: 'im'
---
Deterministic check, replacing a `type: llm` judge that scored 1/3 and 0/1 on
this case (2026-09-25) and failed the same way in 0.7.1: the judge read the
clause as sitting outside the assembled prompt whenever the reply nested a
plain ``` fence inside the outer ``` fence used for "## Assembled prompt" —
the outer fence then visually closes at the first nested ```, long before the
clause, so a reader scanning fence boundaries alone sees the clause as part of
the reviewer's own trailing commentary even though it is still inside the
prompt text the reply constructed. Regex reads by byte position instead of by
fence nesting, so it is not fooled by that closure.

The match requires, in this order: a line starting with 1-6 `#` characters
whose text contains "Assembled prompt" (case-insensitive); then, somewhere
after it, the first line of the anti-delegation clause ("Do not invoke
another CLI binary, spawn another agent, or delegate any part" — sufficient
as an anchor, since grader 01 in this same case already checks the clause's
full text byte-for-byte); then, somewhere after that, a line that is only
3 or 4 backtick characters (the fence that closes the assembled-prompt code
block, whether the reply used a plain ``` fence that closes early on the
first nested ``` it contains, or a wider ```` fence that survives nesting and
closes once at the true end). This fence-close anchor was chosen over the
verdict-contract line ("Your reply must contain, on its own line...") or the
reviewer-tag line ("Your Reported by: tag for this review is...") suggested
as candidates, because 3 of 13 saved replies checked against this case
(2026-09-24 through 2026-09-25, spanning 0.7.1 through 0.8.0) scope the
assembled prompt down to exactly the four parts this case's own prompt names
— role body, skill text, context bundle, anti-delegation clause — and never
add the tag, the reading-is-allowed block, or the verdict contract; anchoring
on either of those would score those replies not-run for a criterion they
otherwise satisfy. Every one of the 13 replies still ends its assembled-prompt
block with a bare fence-close line shortly after the clause, so that anchor
holds across every observed reply shape.

What this grader does not check: the clause's exact wording (grader 01's
job); whether the reviewer tag, reading-is-allowed block, or verdict contract
were also included (not required by this case's own prompt); and it can be
fooled if a reply's trailing commentary quotes the clause outside the real
prompt block and later still contains, anywhere further down, a line that is
only backticks (for example an unrelated illustrative code block) — the
fence-close anchor would then read that unrelated fence as the end of the
prompt and pass a case that should fail. No saved reply exhibits this shape;
it is a known gap in a positional heuristic, not a defect proven against a
real run.
