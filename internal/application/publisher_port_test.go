package application

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
	"testing"
)

type portPublisher struct{}

func (portPublisher) Publish(context.Context, domain.Event) error { return nil }

func TestApplicationPublisherPortPreservesErrorContract(t *testing.T) {
	var _ EventPublisher = portPublisher{}
}
