# Bug Reproduction

## What was wrong

Several HTTP observability paths accessed shared state concurrently without synchronization. The rate limiter updated its window counters without using its mutex, `Recover` appended panic snapshots concurrently, and logger and metrics snapshot state were written from concurrent requests. Under `-race`, these paths reported data races and could lose or corrupt counts.

## How to reproduce

Run the focused race checks from the project module root:

```text
go test -race ./internal/platform/httpx -run '^TestRateLimitConcurrentWindow$' -count=1
go test -race ./internal/platform/httpx -run '^TestRecoverCountsEveryPanic$' -count=1
go test -race ./internal/platform/observability -run '^TestMetricsConcurrentRender$' -count=1
go test -race ./internal/platform/observability -run '^TestAccessLogSnapshotIsolation$' -count=1
```

Before the fix, the checks report races in the limiter, panic snapshot, metrics render, or logger creation paths. After the fix, each shared state transition is synchronized and all four race checks pass repeatedly.
