package memory

import (
	"context"
	"errors"
	"testing"
)

func TestStrategyListPropagatesCancellation(t *testing.T) {
	r := NewStrategyRepository()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := r.List(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("List error = %v, want context.Canceled", err)
	}
}
