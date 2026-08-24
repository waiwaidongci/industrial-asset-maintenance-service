# Bug Reproduction

## What was wrong

The publisher contract used an untyped context parameter, the memory publisher checked `context.Background()`, and error classifiers compared wrapped errors directly. As a result, cancellation and conflict information was lost and cancelled events could still be appended.

## How to reproduce

Run the focused checks from the module root:

```text
go test ./internal/infrastructure/memory -run '^TestPublisherPreservesCancellation$' -count=1
go test ./internal/infrastructure/memory -run '^TestPublishFailureStopsCommit$' -count=1
go test ./internal/domain -run '^TestPublisherErrorChain$' -count=1
go test ./internal/domain -run '^TestPublisherRejectsTypedNilEvent$' -count=1
go test ./internal/application -run '^TestApplicationPublisherPortPreservesErrorContract$' -count=1
```

Before the fix, cancelled publishing returned nil and appended an event, while wrapped conflict errors were not classified. After the fix, cancellation is propagated through `context.Context`, events are not committed on failure, and `errors.Is` recognizes wrapped domain errors.
