# Bug Reproduction

## What Happens

Finding lifecycle validation is inconsistent. An `open` finding can jump directly to `resolved`, active duplicate findings can silently return the existing record, and unknown statuses are accepted on some repository paths but rejected on others.

## How To Trigger

Run the focused checks:

```bash
go test ./internal/domain -run '^TestOpenFindingRequiresAcknowledgementBeforeResolve$' -count=1
go test ./internal/domain -run '^TestAcknowledgedFindingCanResolve$' -count=1
go test ./internal/infrastructure/memory -run '^TestDuplicateActiveFindingReturnsConflict$' -count=1
go test ./internal/infrastructure/memory -run '^TestTerminalFindingAllowsReplacement$' -count=1
go test ./internal/infrastructure/memory -run '^TestCreateRejectsUnknownFindingStatus$' -count=1
go test ./internal/infrastructure/memory -run '^TestUpdateRejectsUnknownFindingStatus$' -count=1
go test ./internal/infrastructure/memory -run '^TestListRejectsStoredUnknownFindingStatus$' -count=1
```

On the affected build, the lifecycle and repository checks fail until transition rules, duplicate handling, and status validation are made consistent across the domain and memory repository layers.
