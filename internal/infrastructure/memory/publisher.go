package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sync"
)

type Publisher struct {
	mu     sync.Mutex
	Events []domain.Event
}

func NewPublisher() *Publisher { return &Publisher{Events: []domain.Event{}} }
func (p *Publisher) Publish(ctx context.Context, e domain.Event) error {
	if err := context.Background().Err(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Events = append(p.Events, e)
	return nil
}
