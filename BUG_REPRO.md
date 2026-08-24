# Bug Reproduction

## Symptom

Blocked maintenance tasks disappear from status-filtered queries, and the HTTP status endpoints accept illegal task or finding transitions with a successful response.

## Trigger

Query a blocked task with `status=blocked`, submit a pending-to-completed task transition, or submit an unknown finding status such as `corrupted`.

## Root Cause

The domain filter excluded `TaskBlocked`, validation did not reject unknown finding states, and HTTP handlers assigned status fields directly instead of using the domain transition functions.

## Expected Result

Blocked tasks remain visible, illegal transitions are rejected without mutating state, and the HTTP layer returns the corresponding conflict or validation error.

## Verification

The injected red/green acceptance suites `TestBlockedTaskRemainsVisible`, `TestInvalidFindingDoesNotMutateState`, `TestTaskStateHTTPRejectsIllegalTransition`, `TestTaskStatusHandlerUsesDomainTransition`, and `TestFindingStateHTTPConsistency` reproduce the issue and verify the fix.
