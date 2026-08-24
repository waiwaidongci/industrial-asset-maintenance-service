package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type TaskRepository struct {
	mu     sync.RWMutex
	values map[string]domain.MaintenanceTask
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{values: map[string]domain.MaintenanceTask{}}
}
func (r *TaskRepository) Create(ctx context.Context, v domain.MaintenanceTask) (domain.MaintenanceTask, error) {
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
func (r *TaskRepository) Get(ctx context.Context, id string) (domain.MaintenanceTask, error) {
	if err := ctx.Err(); err != nil {
		return domain.MaintenanceTask{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, domain.ErrNotFound
	}
	return v, nil
}
func (r *TaskRepository) Update(ctx context.Context, v domain.MaintenanceTask) (domain.MaintenanceTask, error) {
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
func (r *TaskRepository) List(ctx context.Context, f domain.TaskFilter) ([]domain.MaintenanceTask, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.MaintenanceTask{}
	for _, v := range r.values {
		if v.Matches(f) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ScheduledAt.Before(out[j].ScheduledAt) })
	return out, nil
}
