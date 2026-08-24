package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type FindingRepository struct {
	mu     sync.RWMutex
	values map[string]domain.Finding
}

func NewFindingRepository() *FindingRepository {
	return &FindingRepository{values: map[string]domain.Finding{}}
}
func (r *FindingRepository) Create(ctx context.Context, v domain.Finding) (domain.Finding, error) {
	if err := ctx.Err(); err != nil {
		return v, err
	}
	if !domain.IsKnownFindingStatus(v.Status) {
		return v, domain.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; ok {
		return v, domain.ErrConflict
	}
	for _, x := range r.values {
		if x.TaskID != v.TaskID || x.Title != v.Title {
			continue
		}
		if domain.IsActiveFindingStatus(x.Status) {
			return v, domain.ErrConflict
		}
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *FindingRepository) Get(ctx context.Context, id string) (domain.Finding, error) {
	if err := ctx.Err(); err != nil {
		return domain.Finding{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, domain.ErrNotFound
	}
	return v, nil
}
func (r *FindingRepository) Update(ctx context.Context, v domain.Finding) (domain.Finding, error) {
	if err := ctx.Err(); err != nil {
		return v, err
	}
	if !domain.IsKnownFindingStatus(v.Status) {
		return v, domain.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; !ok {
		return v, domain.ErrNotFound
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *FindingRepository) List(ctx context.Context, f domain.FindingFilter) ([]domain.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.Finding{}
	for _, v := range r.values {
		if !domain.IsKnownFindingStatus(v.Status) {
			return nil, domain.ErrInvalid
		}
		if v.Matches(f) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ReportedAt.After(out[j].ReportedAt) })
	return out, nil
}
