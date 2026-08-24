package memory

import (
	"context"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

func assetRepositoryError(operation string, err error) error {
	return fmt.Errorf("%s: %w", operation, err)
}

type AssetRepository struct {
	mu     sync.RWMutex
	values map[string]domain.Asset
}

func NewAssetRepository() *AssetRepository {
	return &AssetRepository{values: map[string]domain.Asset{}}
}
func (r *AssetRepository) Create(ctx context.Context, v domain.Asset) (domain.Asset, error) {
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
func (r *AssetRepository) Get(ctx context.Context, id string) (domain.Asset, error) {
	if err := ctx.Err(); err != nil {
		return domain.Asset{}, assetRepositoryError("get asset context", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.values[id]
	if !ok {
		return v, assetRepositoryError("get asset", domain.ErrNotFound)
	}
	return v, nil
}
func (r *AssetRepository) Update(ctx context.Context, v domain.Asset) (domain.Asset, error) {
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
func (r *AssetRepository) List(ctx context.Context, f domain.AssetFilter) ([]domain.Asset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Asset, 0, len(r.values))
	for _, v := range r.values {
		if v.Matches(f) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
