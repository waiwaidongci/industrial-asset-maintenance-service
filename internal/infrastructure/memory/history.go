package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"sync"
)

type HistoryRepository struct {
	mu     sync.RWMutex
	values []domain.HistoryEntry
}

func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{values: []domain.HistoryEntry{}}
}
func (r *HistoryRepository) Append(ctx context.Context, v domain.HistoryEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values = append(r.values, v)
	return nil
}
func (r *HistoryRepository) List(ctx context.Context, taskID string) ([]domain.HistoryEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []domain.HistoryEntry{}
	for _, v := range r.values {
		if taskID == "" || v.TaskID == taskID {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.Before(out[j].OccurredAt) })
	return out, nil
}
