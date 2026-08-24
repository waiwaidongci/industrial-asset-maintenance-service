package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type TemplateRepository struct {
	mu     sync.RWMutex
	values map[string]domain.InspectionTemplate
}

func NewTemplateRepository() *TemplateRepository {
	return &TemplateRepository{values: map[string]domain.InspectionTemplate{}}
}
func (r *TemplateRepository) Create(ctx context.Context, v domain.InspectionTemplate) (domain.InspectionTemplate, error) {
	if err := ctx.Err(); err != nil {
		return v, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[v.ID]; ok {
		return v, domain.ErrConflict
	}
	for _, x := range r.values {
		if x.Name == v.Name && x.Version == v.Version {
			return v, domain.ErrConflict
		}
	}
	r.values[v.ID] = v
	return v, nil
}
func (r *TemplateRepository) Get(ctx context.Context, id string) (domain.InspectionTemplate, error) {
	if err := ctx.Err(); err != nil {
		return domain.InspectionTemplate{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, domain.ErrNotFound
	}
	return v, nil
}
func (r *TemplateRepository) Update(ctx context.Context, v domain.InspectionTemplate) (domain.InspectionTemplate, error) {
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
func (r *TemplateRepository) List(ctx context.Context, active bool) ([]domain.InspectionTemplate, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.InspectionTemplate{}
	for _, v := range r.values {
		if active && !v.Active {
			continue
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
