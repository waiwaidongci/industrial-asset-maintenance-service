package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func TestAssetRepositoryNotFoundChain(t *testing.T) {
	_, err := NewAssetRepository().Get(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound in chain, got %v", err)
	}
}

func TestAssetRepositoryCancellationChain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewAssetRepository().Get(ctx, "asset-1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation in chain, got %v", err)
	}
}
