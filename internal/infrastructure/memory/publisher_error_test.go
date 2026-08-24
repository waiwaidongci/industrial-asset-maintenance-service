package memory

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"testing"
)

func TestPublisherPreservesCancellation(t *testing.T) {
	p := NewPublisher()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.Publish(ctx, domain.Event{Type: domain.EventTaskCreated}); err != context.Canceled {
		t.Fatalf("Publish error = %v, want context.Canceled", err)
	}
}

func TestPublishFailureStopsCommit(t *testing.T) {
	p := NewPublisher()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = p.Publish(ctx, domain.Event{Type: domain.EventTaskCreated})
	if len(p.Events) != 0 {
		t.Fatalf("cancelled event was committed: %d", len(p.Events))
	}
}
