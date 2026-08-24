# Bug Reproduction

## What was wrong

`TemplateRepository` copied `InspectionTemplate` values shallowly across its Create, Get, List, and Update boundaries. The `Items` slice and each item's `Min`/`Max` pointers remained shared, so changing a returned template or the input after saving could mutate the repository snapshot.

## How to reproduce

Run these focused checks from the project module root:

```text
go test ./internal/infrastructure/memory -run '^TestTemplateRepositoryGetIsolatesItems$' -count=1
go test ./internal/infrastructure/memory -run '^TestTemplateRepositoryListIsolatesItems$' -count=1
go test ./internal/infrastructure/memory -run '^TestTemplateRepositoryCreateIsolatesInput$' -count=1
go test ./internal/infrastructure/memory -run '^TestTemplateRepositoryUpdateIsolatesInput$' -count=1
```

Before the fix, mutating an item or threshold through a Get/List result, or mutating the input after Create/Update, changes the stored snapshot. After the fix, all nested slices and threshold pointers are cloned at every boundary and the checks pass repeatedly.
