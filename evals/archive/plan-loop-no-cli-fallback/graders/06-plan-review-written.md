---
type: regex
target: trace
pattern: '"file_path":[ \t]*"[^"]*PLAN\.review\.md"'
flags: ''
match: contains
weight: 1
---

A `Write` (or `Edit`) tool call in the trace targets a path ending in
`PLAN.review.md`. AC-03 (and T-803 §7) names `PLAN.md`, its `TASK.md`
file(s), and `PLAN.review.md` as the exact artifact set the no-CLI
fallback persists, matching `/dh:plan`'s own Output in full — a run that
drafts and validates the plan but never persists the validator's review to
disk satisfies every other grader in this case while silently skipping
one of the three required artifacts, which is exactly the gap grader 05
(last-message description only) cannot catch on its own. The pattern is
anchored to the serialized `file_path` key, not a bare `PLAN.review.md`
substring, so a mention of the filename in a reply's own prose (for
example, naming the path it is about to write) cannot satisfy this
grader without an actual write.
