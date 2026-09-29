# Transport double: cli:agy/gemini-3.1-pro-high, lru.go finding only

Pre-written stdout capture, used as fixture data only. No real `agy` binary
is invoked to produce this file or to grade any case that reads it.
Reports a finding in `internal/cache/lru.go` only — nothing about
`internal/service/retry.go`, so this reviewer's own reply names one
location, genuinely unrelated to the other reviewer's own finding: a
different file, a different line, and a different defect class (unbounded
growth vs. a missing backoff sleep).

```
## Review
- Scope: internal/service/retry.go, internal/cache/lru.go
- Comparison: main..fix/retry-and-cache
- Requirement source: none available
- Axis coverage: correctness complete; spec not-run

## Correctness and regressions
### Major
- [internal/cache/lru.go:87] Set no longer evicts the oldest entry before inserting when the cache is already at capacity, so the cache grows without bound instead of staying at its configured size.
  - Scenario: Set is called more times than the cache's configured capacity, each with a distinct key.
  - Impact: unbounded memory growth under sustained cache writes, defeating the purpose of a bounded LRU cache.
  - Recommendation: keep the capacity check and evictOldest call before inserting the new entry.
  - Reported by: cli:agy/gemini-3.1-pro-high

## Spec compliance
No spec available. Axis not run.

## Checks observed
- go test ./... - not-run (sandbox)

## Verdict
- Review status: request changes
- Blocking findings: 1
- Ready from this review: no
- Limitations: spec axis not run; no requirement source provided; internal/service/retry.go not reviewed by this call.
```
