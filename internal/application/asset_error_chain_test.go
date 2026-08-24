package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/asset-maintenance-service/internal/domain"
)

type assetErrorRepository struct {
	getErr    error
	updateErr error
}

func (r assetErrorRepository) Create(context.Context, domain.Asset) (domain.Asset, error) {
	return domain.Asset{}, nil
}
func (r assetErrorRepository) Get(context.Context, string) (domain.Asset, error) {
	return domain.Asset{}, r.getErr
}
func (r assetErrorRepository) Update(_ context.Context, asset domain.Asset) (domain.Asset, error) {
	return asset, r.updateErr
}
func (r assetErrorRepository) List(context.Context, domain.AssetFilter) ([]domain.Asset, error) {
	return nil, nil
}

func TestAssetServiceNotFoundChain(t *testing.T) {
	svc := &Service{Assets: assetErrorRepository{getErr: domain.ErrNotFound}}
	_, err := svc.GetAsset(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected service to preserve ErrNotFound, got %v", err)
	}
}

func TestAssetServiceUpdateConflictChain(t *testing.T) {
	svc := &Service{Assets: assetErrorRepository{updateErr: domain.ErrConflict}, Now: time.Now}
	asset := domain.Asset{ID: "asset-1", Name: "pump", AssetType: "pump", Status: domain.AssetActive}
	_, err := svc.UpdateAsset(context.Background(), asset)
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected service to preserve ErrConflict, got %v", err)
	}
}
