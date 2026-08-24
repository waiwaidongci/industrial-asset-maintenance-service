package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type PlanRepository struct {
	mu     sync.RWMutex
	values map[string]domain.MaintenancePlan
}

func NewPlanRepository() *PlanRepository {
	return &PlanRepository{values: map[string]domain.MaintenancePlan{}}
}
func (r *PlanRepository) Create(ctx context.Context, v domain.MaintenancePlan) (domain.MaintenancePlan, error) {
	if err := ctx.Err(); err != nil {
		return v, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; ok {
		return v, domain.ErrConflict
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *PlanRepository) Get(ctx context.Context, id string) (domain.MaintenancePlan, error) {
	if err := ctx.Err(); err != nil {
		return domain.MaintenancePlan{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, domain.ErrNotFound
	}
	return v, nil
}
func (r *PlanRepository) Update(ctx context.Context, v domain.MaintenancePlan) (domain.MaintenancePlan, error) {
	if err := ctx.Err(); err != nil {
		return v, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; !ok {
		return v, domain.ErrNotFound
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *PlanRepository) List(ctx context.Context, f domain.PlanFilter) ([]domain.MaintenancePlan, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.MaintenancePlan{}
	for _, v := range r.values {
		if v.Matches(f) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].NextRunAt.Before(out[j].NextRunAt) })
	return out, nil
}
