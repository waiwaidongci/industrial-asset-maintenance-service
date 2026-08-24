package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type StrategyRepository struct {
	mu     sync.RWMutex
	values map[string]domain.MaintenanceStrategy
}

func NewStrategyRepository() *StrategyRepository {
	return &StrategyRepository{values: map[string]domain.MaintenanceStrategy{}}
}
func (r *StrategyRepository) Create(ctx context.Context, v domain.MaintenanceStrategy) (domain.MaintenanceStrategy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; ok {
		return v, domain.ErrConflict
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *StrategyRepository) Get(ctx context.Context, id string) (domain.MaintenanceStrategy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, domain.ErrNotFound
	}
	return v, nil
}
func (r *StrategyRepository) List(ctx context.Context) ([]domain.MaintenanceStrategy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.MaintenanceStrategy{}
	for _, v := range r.values {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
