# Bug Reproduction

## What was wrong

Query snapshots crossed ownership boundaries. `HistoryRepository.List` returned its internal slice and shallow-copied `Details` maps, while `SortTasksByPriority` sorted the caller's task slice in place. A caller could therefore mutate history details or reorder tasks and contaminate later summaries.

## How to reproduce

Run these focused checks from the project module root:

```text
go test ./internal/infrastructure/memory -run '^TestHistoryListIsolatesEntries$' -count=1
go test ./internal/application -run '^TestSortTasksByPriorityDoesNotMutateInput$' -count=1
```

Before the fix, changing a returned history entry or its `Details` map, and sorting a task result, changes the repository or caller-owned input. After the fix, history entries and nested maps are cloned at the boundary and sorting works on a copy, so both checks pass repeatedly.
