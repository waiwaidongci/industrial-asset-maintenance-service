# Bug Reproduction

## What Happens

Dispatch workers update the `WaitGroup` after they start, while result-channel closure depends on the same group. Error and cancellation branches can exit without publishing a result. The scheduler can therefore deadlock or return an incomplete batch. The plan and task repositories also ignore canceled contexts, allowing work to continue after the caller has gone away.

## How To Trigger

Run the focused scheduler and cancellation checks with the race detector:

```bash
go test -race ./internal/application -run '^TestDispatchDueConcurrentPlans$' -count=1
go test -race ./internal/application -run '^TestDispatchDueClosesResults$' -count=1
go test -race ./internal/application -run '^TestDispatchDueCancellation$' -count=1
go test -race ./internal/infrastructure/memory -run '^TestPlanListPropagatesCancellation$' -count=1
go test -race ./internal/infrastructure/memory -run '^TestTaskCreatePropagatesCancellation$' -count=1
go test -race ./internal/infrastructure/memory -run '^TestTaskListPropagatesCancellation$' -count=1
```

## Observed Error

On the affected build, concurrent dispatch can wait indefinitely or omit tasks. Cancellation checks report errors such as `expected canceled due-plan lookup to return an error`, `expected canceled task create, got <nil>`, and `expected canceled task list, got <nil>`. The focused commands exit with status 1 until worker completion, channel closure, unique task creation, and context propagation are handled consistently.
