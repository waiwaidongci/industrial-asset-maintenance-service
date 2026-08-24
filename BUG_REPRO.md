# Bug Reproduction

## What Happens

Asset repository and application errors are converted to plain text while crossing layers. The HTTP error handler also compares sentinel errors directly. As a result, missing assets and update conflicts cannot be recognized with `errors.Is`, and public asset endpoints return HTTP 500 instead of the appropriate client-facing status.

## How To Trigger

Run the focused error-chain and HTTP checks:

```bash
go test ./internal/infrastructure/memory -run '^TestAssetRepositoryNotFoundChain$' -count=1
go test ./internal/infrastructure/memory -run '^TestAssetRepositoryCancellationChain$' -count=1
go test ./internal/application -run '^TestAssetServiceNotFoundChain$' -count=1
go test ./internal/application -run '^TestAssetServiceUpdateConflictChain$' -count=1
go test ./internal/adapter/http -run '^TestMissingAssetHTTPStatus$' -count=1
```

## Observed Error

On the affected build, the checks report that wrapped repository and service errors no longer match `domain.ErrNotFound`, `domain.ErrConflict`, or cancellation errors. The public missing-asset request consequently returns status 500 where status 404 is expected. Each focused command exits with status 1 until the error chain and HTTP classification are preserved across all three layers.
