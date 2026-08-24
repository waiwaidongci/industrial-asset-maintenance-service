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
	out := r.values
	r.mu.RUnlock()
	filtered := out[:0]
	for _, v := range out {
		if taskID == "" || v.TaskID == taskID {
			filtered = append(filtered, v)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].OccurredAt.Before(filtered[j].OccurredAt) })
	return filtered, nil
}
