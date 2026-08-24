# Bug Reproduction

## What was wrong

Shutdown configuration was not propagated consistently. YAML and environment values could be ignored or accepted outside the supported bounds, the server used the wrong timeout when creating its shutdown context, and the scheduler replaced the cancellation chain with a background context. A signal could therefore leave workers running past the configured shutdown deadline.

## How to reproduce

Run these focused checks from the project module root:

```text
go test ./internal/platform/config -run '^TestShutdownTimeoutFromYAML$' -count=1
go test ./internal/platform/config -run '^TestShutdownTimeoutFromEnv$' -count=1
go test ./internal/platform/config -run '^TestRejectsZeroShutdownTimeout$' -count=1
go test ./internal/platform/config -run '^TestCapsShutdownTimeout$' -count=1
go test ./cmd/server -run '^TestServerUsesConfiguredShutdownTimeout$' -count=1
go test ./cmd/scheduler -run '^TestSchedulerStopsOnSignal$' -count=1
```

Before the fix, configuration and shutdown checks fail because the configured duration is lost or the wrong context is used; the scheduler cancellation check also fails because its worker continues after a signal. After the fix, values are validated and propagated through the server and scheduler shutdown path, and all checks pass across repeated runs.
