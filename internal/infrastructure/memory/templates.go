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

// cloneTemplate deep-copies an InspectionTemplate's nested Items slice, including
// the pointer-valued Min/Max fields, so callers cannot mutate the repository
// snapshot by editing the returned template. This is the boundary-isolation
// contract for all read and write paths.
func cloneTemplate(v domain.InspectionTemplate) domain.InspectionTemplate {
	if len(v.Items) == 0 {
		v.Items = nil
		return v
	}
	items := make([]domain.InspectionItem, len(v.Items))
	for i, it := range v.Items {
		items[i] = it
		if it.Min != nil {
			m := *it.Min
			items[i].Min = &m
		}
		if it.Max != nil {
			m := *it.Max
			items[i].Max = &m
		}
	}
	v.Items = items
	return v
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
	r.values[v.ID] = cloneTemplate(v)
	return cloneTemplate(v), nil
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
	return cloneTemplate(v), nil
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
	r.values[v.ID] = cloneTemplate(v)
	return cloneTemplate(v), nil
}
func (r *TemplateRepository) List(ctx context.Context, active bool) ([]domain.InspectionTemplate, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.InspectionTemplate, 0, len(r.values))
	for _, v := range r.values {
		if active && !v.Active {
			continue
		}
		out = append(out, cloneTemplate(v))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}
