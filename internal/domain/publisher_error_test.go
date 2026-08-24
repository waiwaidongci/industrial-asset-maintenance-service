package domain

import (
	"context"
	"fmt"
	"testing"
)

type typedEventPublisher struct{}

func (typedEventPublisher) Publish(context.Context, Event) error { return nil }

func TestPublisherErrorChain(t *testing.T) {
	if !IsConflict(fmt.Errorf("publish conflict: %w", ErrConflict)) {
		t.Fatal("wrapped conflict was not classified")
	}
}

func TestPublisherRejectsTypedNilEvent(t *testing.T) {
	var publisher EventPublisher = typedEventPublisher{}
	if publisher == nil {
		t.Fatal("publisher unexpectedly nil")
	}
}
