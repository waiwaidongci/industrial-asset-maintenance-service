# Bug Reproduction

## What was wrong

`EvaluateItem` accepted `NaN`, `+Inf`, and `-Inf` from `strconv.ParseFloat`; threshold comparisons with those values did not fail, so an invalid reading could be marked as passing. `ValidateTaskResults` also accepted result IDs that were not present in the inspection template and allowed a `Passed` flag that contradicted the calculated outcome. `StrategyRepository.List` checked `context.Background()` instead of the caller context, so a cancelled request still returned strategy data.

## How to reproduce

Run the focused checks from the project module root:

```text
go test ./internal/domain -run '^TestEvaluateItemRejectsNonFiniteNumbers$' -count=1
go test ./internal/domain -run '^TestValidateTaskResultsRejectsUnknownItems$' -count=1
go test ./internal/domain -run '^TestValidateTaskResultsRejectsInconsistentPassFlag$' -count=1
go test ./internal/infrastructure/memory -run '^TestStrategyListPropagatesCancellation$' -count=1
```

Before the fix, the checks fail with a non-finite value reported as `pass`, unknown item data accepted, contradictory pass flags accepted, or a nil error from a cancelled strategy list request. After the fix, all four checks pass and the invalid inputs are rejected with `domain.ErrInvalid` or `context.Canceled`.
